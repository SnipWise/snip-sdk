package chatserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/chat"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
)

type ChatAgentServer struct {
	agent *chat.ChatAgent

	serverConfig *ConfigHTTP
	httpServer   *http.Server
	serverCancel context.CancelFunc

	logger logger.Logger

	ctx context.Context
}

/* TODO: & QUESTION:
- logger?
*/

func NewChatAgentServer(
	ctx context.Context,
	agentConfig agents.AgentConfig,
	modelConfig models.ModelConfig,
	opts ...ChatAgentServerOption) (*ChatAgentServer, error) {
	//ctx = context.Background()
	agent, err := chat.NewChatAgent(
		ctx,
		agentConfig,
		modelConfig,
		chat.EnableChatFlow(),
		chat.EnableChatStreamFlow(),
		chat.EnableChatFlowWithMemory(),
		chat.EnableChatStreamFlowWithMemory(),
	)

	if err != nil {
		return nil, err
	}

	chatAgentServer := &ChatAgentServer{
		agent:  agent,
		ctx:    ctx,
		logger: &logger.NoOpLogger{}, // Initialize with a no-op logger by default
	}
	for _, opt := range opts {
		opt(chatAgentServer)
	}

	return chatAgentServer, nil
}

func (cas *ChatAgentServer) GetName() string {
	return cas.agent.Name
}

func (cas *ChatAgentServer) Kind() agents.AgentKind {
	return agents.ChatServer
}

func (cas *ChatAgentServer) GetMessages() []*ai.Message {
	return cas.agent.Messages
}

func (cas *ChatAgentServer) GetCurrentContextSize() int {
	return cas.agent.GetCurrentContextSize()
}

func (cas *ChatAgentServer) AddSystemMessage(context string) error {
	return cas.agent.AddSystemMessage(context)
}

func (cas *ChatAgentServer) ReplaceMessagesWith(messages []*ai.Message) error {
	return cas.agent.ReplaceMessagesWith(messages)
}

func (cas *ChatAgentServer) ReplaceMessagesWithSystemMessages(systemMessages []string) error {
	return cas.agent.ReplaceMessagesWithSystemMessages(systemMessages)
}

func (cas *ChatAgentServer) GetInfo() (agents.AgentInfo, error) {
	return cas.agent.GetInfo()
}

func (cas *ChatAgentServer) AskWithMemory(question string) (agents.ChatResponse, error) {
	return cas.agent.AskWithMemory(question)
}

func (cas *ChatAgentServer) AskStreamWithMemory(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	return cas.agent.AskStreamWithMemory(question, callback)
}

func (cas *ChatAgentServer) Ask(question string) (agents.ChatResponse, error) {
	return cas.agent.Ask(question)
}

func (cas *ChatAgentServer) AskStream(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	return cas.agent.AskStream(question, callback)
}

// Serve starts the HTTP server with the configured endpoints for the agent's flows
// The server automatically handles SIGINT (Ctrl+C) and SIGTERM signals for graceful shutdown
// Use the Stop() method to manually shutdown the server
func (cas *ChatAgentServer) Serve() error {
	if cas.serverConfig == nil {
		return fmt.Errorf("server configuration is not set, use EnableServer option")
	}

	mux := http.NewServeMux()

	// Set default values for paths if not provided
	if cas.serverConfig.ChatFlowPath == "" {
		cas.serverConfig.ChatFlowPath = DefaultChatFlowPath
	}
	if cas.serverConfig.ChatStreamFlowPath == "" {
		cas.serverConfig.ChatStreamFlowPath = DefaultChatStreamFlowPath
	}
	if cas.serverConfig.InformationPath == "" {
		cas.serverConfig.InformationPath = DefaultInformationPath
	}
	if cas.serverConfig.ShutdownPath == "" {
		cas.serverConfig.ShutdownPath = "-"
	}
	if cas.serverConfig.CancelStreamPath == "" {
		cas.serverConfig.CancelStreamPath = DefaultCancelStreamPath
	}
	if cas.serverConfig.AddContextPath == "" {
		cas.serverConfig.AddContextPath = DefaultAddSystemMessagePath
	}
	if cas.serverConfig.HealthcheckPath == "" {
		cas.serverConfig.HealthcheckPath = DefaultHealthcheckPath
	}
	if cas.serverConfig.GetMessagesPath == "" {
		cas.serverConfig.GetMessagesPath = DefaultGetMessagesPath
	}

	// Register healthcheck endpoint
	healthcheckPath := cas.serverConfig.HealthcheckPath
	mux.HandleFunc("GET "+healthcheckPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	cas.logger.Info("Registered endpoint: GET %s", healthcheckPath)

	// Register agent information endpoint
	informationPath := cas.serverConfig.InformationPath
	mux.HandleFunc("GET "+informationPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		info := agents.AgentInfo{
			Name:    cas.agent.Name,
			ModelID: cas.agent.ModelID,
			Config:  cas.agent.Config,
		}
		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	cas.logger.Info("Registered endpoint: GET %s", informationPath)

	// IMPORTANT: with memory flows
	// Register chat flow endpoint if available
	if cas.agent.GetChatFlowWithMemory() != nil && cas.serverConfig.ChatFlowHandler != nil {
		chatFlowPath := cas.serverConfig.ChatFlowPath
		mux.HandleFunc("POST "+chatFlowPath, cas.serverConfig.ChatFlowHandler)
		cas.logger.Info("Registered endpoint: POST %s", chatFlowPath)
	}
	// IMPORTANT: with memory flows
	// Register chat stream flow endpoint if available
	if cas.agent.GetChatStreamFlowWithMemory() != nil && cas.serverConfig.ChatStreamFlowHandler != nil {
		chatStreamFlowPath := cas.serverConfig.ChatStreamFlowPath
		mux.HandleFunc("POST "+chatStreamFlowPath, cas.serverConfig.ChatStreamFlowHandler)
		cas.logger.Info("Registered endpoint: POST %s", chatStreamFlowPath)
	}

	// Create server context with cancel
	//serverCtx, cancel := context.WithCancel(cas.agent.ctx)
	// NOTE: TODO: to be checked
	serverCtx, cancel := context.WithCancel(cas.ctx)

	cas.serverCancel = cancel

	// Register shutdown endpoint if enabled
	shutdownPath := cas.serverConfig.ShutdownPath
	if shutdownPath != "-" {
		if shutdownPath == "" {
			shutdownPath = DefaultShutdownPath
		}
		mux.HandleFunc("POST "+shutdownPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"shutting down"}`))

			cas.logger.Info("Shutdown requested via HTTP endpoint")

			// Trigger shutdown asynchronously to allow response to be sent
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()
		})
		cas.logger.Info("Registered endpoint: POST %s", shutdownPath)
	}

	// Register cancel stream endpoint
	cancelStreamPath := cas.serverConfig.CancelStreamPath
	if cancelStreamPath != "" {
		mux.HandleFunc("POST "+cancelStreamPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			cancelFunc := cas.agent.GetStreamCancel()
			if cancelFunc != nil {
				cancelFunc()
				cas.logger.Info("Streaming completion cancelled via HTTP endpoint")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"stream cancelled"}`))
			} else {
				cas.logger.Info("No active stream to cancel")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"no active stream"}`))
			}
		})
		cas.logger.Info("Registered endpoint: POST %s", cancelStreamPath)
	}

	// Register add context endpoint
	addContextPath := cas.serverConfig.AddContextPath
	if addContextPath != "" {
		mux.HandleFunc("POST "+addContextPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Parse request body
			var req struct {
				Context string `json:"context"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				cas.logger.Error("Error decoding add context request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"invalid request body"}`))
				return
			}

			// Add context to messages
			if err := cas.AddSystemMessage(req.Context); err != nil {
				cas.logger.Error("Error adding context to messages: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to add context"}`))
				return
			}

			cas.logger.Info("Context added to messages via HTTP endpoint")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		})
		cas.logger.Info("Registered endpoint: POST %s", addContextPath)
	}

	// Register get messages endpoint
	getMessagesPath := cas.serverConfig.GetMessagesPath
	if getMessagesPath != "" {
		mux.HandleFunc("GET "+getMessagesPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Encode messages to JSON
			if err := json.NewEncoder(w).Encode(cas.agent.Messages); err != nil {
				cas.logger.Error("Error encoding messages: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to encode messages"}`))
				return
			}

			cas.logger.Debug("Messages retrieved via HTTP endpoint")
		})
		cas.logger.Info("Registered endpoint: GET %s", getMessagesPath)
	}

	cas.httpServer = &http.Server{
		Addr:    cas.serverConfig.Address,
		Handler: mux,
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		cas.logger.Info("Starting HTTP server on %s (Press Ctrl+C to stop)", cas.serverConfig.Address)
		serverErrors <- cas.httpServer.ListenAndServe()
	}()

	// Wait for either context cancellation, signal, or server error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	case sig := <-sigChan:
		cas.logger.Info("Received signal: %v", sig)
		return cas.Stop()
	case <-serverCtx.Done():
		return cas.Stop()
	}
}

// Stop gracefully shuts down the HTTP server with a 5-second timeout
func (cas *ChatAgentServer) Stop() error {
	if cas.httpServer == nil {
		return fmt.Errorf("server is not running")
	}

	cas.logger.Info("Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cas.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error during shutdown: %w", err)
	}

	cas.logger.Info("Server stopped")
	return nil
}

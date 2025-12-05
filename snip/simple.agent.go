package snip

/*
This is a simple agent
the conversation history is stored in memory
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/snipwise/snip-sdk/conversion"
	"github.com/snipwise/snip-sdk/env"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

type Agent struct {
	ctx                context.Context
	Name               string
	SystemInstructions string
	ModelID            string
	Messages           []*ai.Message

	Config ModelConfig

	genKitInstance *genkit.Genkit

	chatStreamFlow *core.Flow[*ChatRequest, *ChatResponse, string]
	chatFlow       *core.Flow[*ChatRequest, *ChatResponse, struct{}]

	serverConfig *ConfigHTTP
	httpServer   *http.Server
	serverCancel context.CancelFunc

	// streamCancel cancels the current streaming completion
	streamCancel context.CancelFunc
	streamCtx    context.Context
}

func (agent *Agent) GetName() string {
	return agent.Name
}

func (agent *Agent) GetMessages() []*ai.Message {
	return agent.Messages
}

func (agent *Agent) Kind() AgentKind {
	return Basic
}

func (agent *Agent) AddSystemMessage(context string) error {
	// Add a system message to the conversation history
	agent.Messages = append(agent.Messages, ai.NewSystemTextMessage(strings.TrimSpace(context)))
	return nil
}

func (agent *Agent) GetInfo() (AgentInfo, error) {
	return AgentInfo{
		Name:    agent.Name,
		Config:  agent.Config,
		ModelID: agent.ModelID,
	}, nil
}

// ping checks if the model is available by sending a simple prompt
// func ping(ctx context.Context, genKitInstance *genkit.Genkit, modelID string) error {
// 	log.Println("⏳ model availability check in progress...")
// 	_, err := genkit.Generate(
// 		ctx,
// 		genKitInstance,
// 		ai.WithModelName("openai/"+modelID),
// 		ai.WithPrompt(""),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("model not available: %w", err)
// 	}
// 	return nil
// }

func NewAgent(ctx context.Context, agentConfig AgentConfig, modelConfig ModelConfig, opts ...AgentOption) (*Agent, error) {
	oaiPlugin := &oai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(agentConfig.EngineURL),
		},
	}

	genKitInstance := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	// Check if model is available
	if !IsModelAvailable(ctx, agentConfig.EngineURL, agentConfig.ModelID) {
		return nil, fmt.Errorf("model %s is not available at %s", agentConfig.ModelID, agentConfig.EngineURL)
	} else {
		log.Printf("✅ Model %s is available at %s", agentConfig.ModelID, agentConfig.EngineURL)
	}
	// err := ping(ctx, genKitInstance, agentConfig.ModelID)
	// if err != nil {
	// 	return nil, err
	// }

	agent := &Agent{
		Name:               agentConfig.Name,
		SystemInstructions: agentConfig.SystemInstructions,
		ModelID:            agentConfig.ModelID,
		Messages:           []*ai.Message{},
		Config:             modelConfig,

		ctx:            ctx,
		genKitInstance: genKitInstance,
	}

	// Apply all options
	for _, opt := range opts {
		opt(agent)
	}

	return agent, nil

}

func (agent *Agent) Ask(question string) (ChatResponse, error) {
	if agent.chatFlow == nil {
		return ChatResponse{}, fmt.Errorf("chat flow is not initialized")
	}
	resp, err := agent.chatFlow.Run(agent.ctx, &ChatRequest{
		UserMessage: question,
	})
	if err != nil {
		return ChatResponse{}, err
	}
	return *resp, nil

}

func (agent *Agent) AskStream(question string, callback func(string) error) (string, error) {
	if agent.chatStreamFlow == nil {
		return "", fmt.Errorf("chat stream flow is not initialized")
	}
	// Streaming channel of results
	streamCh := agent.chatStreamFlow.Stream(agent.ctx, &ChatRequest{
		UserMessage: question,
	})

	finalAnswer := ""
	for result, err := range streamCh {
		// Check for errors from the stream
		if err != nil {
			// Return both the partial answer and the error
			return finalAnswer, fmt.Errorf("streaming error: %w", err)
		}

		// Check for nil result (defensive programming)
		if result == nil {
			continue
		}

		if !result.Done {
			finalAnswer += result.Stream
			err := callback(result.Stream)
			if err != nil {
				return finalAnswer, err
			}
		}
	}

	return finalAnswer, nil
}

// Serve starts the HTTP server with the configured endpoints for the agent's flows
// The server automatically handles SIGINT (Ctrl+C) and SIGTERM signals for graceful shutdown
// Use the Stop() method to manually shutdown the server
func (agent *Agent) Serve() error {
	if agent.serverConfig == nil {
		return fmt.Errorf("server configuration is not set, use EnableServer option")
	}

	mux := http.NewServeMux()

	// Set default values for paths if not provided
	if agent.serverConfig.ChatFlowPath == "" {
		agent.serverConfig.ChatFlowPath = DefaultChatFlowPath
	}
	if agent.serverConfig.ChatStreamFlowPath == "" {
		agent.serverConfig.ChatStreamFlowPath = DefaultChatStreamFlowPath
	}
	if agent.serverConfig.InformationPath == "" {
		agent.serverConfig.InformationPath = DefaultInformationPath
	}
	if agent.serverConfig.ShutdownPath == "" {
		agent.serverConfig.ShutdownPath = "-"
	}
	if agent.serverConfig.CancelStreamPath == "" {
		agent.serverConfig.CancelStreamPath = DefaultCancelStreamPath
	}
	if agent.serverConfig.AddContextPath == "" {
		agent.serverConfig.AddContextPath = DefaultAddSystemMessagePath
	}
	if agent.serverConfig.HealthcheckPath == "" {
		agent.serverConfig.HealthcheckPath = DefaultHealthcheckPath
	}
	if agent.serverConfig.GetMessagesPath == "" {
		agent.serverConfig.GetMessagesPath = DefaultGetMessagesPath
	}

	// Register healthcheck endpoint
	healthcheckPath := agent.serverConfig.HealthcheckPath
	mux.HandleFunc("GET "+healthcheckPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	log.Printf("[%s] Registered endpoint: GET %s", agent.Name, healthcheckPath)

	// Register agent information endpoint
	informationPath := agent.serverConfig.InformationPath
	mux.HandleFunc("GET "+informationPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		info := AgentInfo{
			Name:    agent.Name,
			ModelID: agent.ModelID,
			Config:  agent.Config,
		}
		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Printf("[%s] Registered endpoint: GET %s", agent.Name, informationPath)

	// Register chat flow endpoint if available
	if agent.chatFlow != nil && agent.serverConfig.ChatFlowHandler != nil {
		chatFlowPath := agent.serverConfig.ChatFlowPath
		mux.HandleFunc("POST "+chatFlowPath, agent.serverConfig.ChatFlowHandler)
		log.Printf("[%s] Registered endpoint: POST %s", agent.Name, chatFlowPath)
	}

	// Register chat stream flow endpoint if available
	if agent.chatStreamFlow != nil && agent.serverConfig.ChatStreamFlowHandler != nil {
		chatStreamFlowPath := agent.serverConfig.ChatStreamFlowPath
		mux.HandleFunc("POST "+chatStreamFlowPath, agent.serverConfig.ChatStreamFlowHandler)
		log.Printf("[%s] Registered endpoint: POST %s", agent.Name, chatStreamFlowPath)
	}

	// Create server context with cancel
	serverCtx, cancel := context.WithCancel(agent.ctx)
	agent.serverCancel = cancel

	// Register shutdown endpoint if enabled
	shutdownPath := agent.serverConfig.ShutdownPath
	if shutdownPath != "-" {
		if shutdownPath == "" {
			shutdownPath = DefaultShutdownPath
		}
		mux.HandleFunc("POST "+shutdownPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"shutting down"}`))

			//log.Println("Shutdown requested via HTTP endpoint")
			log.Printf("[%s] Shutdown requested via HTTP endpoint", agent.Name)

			// Trigger shutdown asynchronously to allow response to be sent
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()
		})
		log.Printf("[%s] Registered endpoint: POST %s", agent.Name, shutdownPath)
	}

	// Register cancel stream endpoint
	cancelStreamPath := agent.serverConfig.CancelStreamPath
	if cancelStreamPath != "" {
		mux.HandleFunc("POST "+cancelStreamPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if agent.streamCancel != nil {
				agent.streamCancel()
				log.Printf("[%s] Streaming completion cancelled via HTTP endpoint", agent.Name)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"stream cancelled"}`))
			} else {
				log.Printf("[%s] No active stream to cancel", agent.Name)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"no active stream"}`))
			}
		})
		log.Printf("[%s] Registered endpoint: POST %s", agent.Name, cancelStreamPath)
	}

	// Register add context endpoint
	addContextPath := agent.serverConfig.AddContextPath
	if addContextPath != "" {
		mux.HandleFunc("POST "+addContextPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Parse request body
			var req struct {
				Context string `json:"context"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Printf("[%s] Error decoding add context request: %v", agent.Name, err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"invalid request body"}`))
				return
			}

			// Add context to messages
			if err := agent.AddSystemMessage(req.Context); err != nil {
				log.Printf("[%s] Error adding context to messages: %v", agent.Name, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to add context"}`))
				return
			}

			log.Printf("[%s] Context added to messages via HTTP endpoint", agent.Name)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		})
		log.Printf("[%s] Registered endpoint: POST %s", agent.Name, addContextPath)
	}

	// Register get messages endpoint
	getMessagesPath := agent.serverConfig.GetMessagesPath
	if getMessagesPath != "" {
		mux.HandleFunc("GET "+getMessagesPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Encode messages to JSON
			if err := json.NewEncoder(w).Encode(agent.Messages); err != nil {
				log.Printf("[%s] Error encoding messages: %v", agent.Name, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to encode messages"}`))
				return
			}

			log.Printf("[%s] Messages retrieved via HTTP endpoint", agent.Name)
		})
		log.Printf("[%s] Registered endpoint: GET %s", agent.Name, getMessagesPath)
	}

	agent.httpServer = &http.Server{
		Addr:    agent.serverConfig.Address,
		Handler: mux,
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("[%s] Starting HTTP server on %s (Press Ctrl+C to stop)", agent.Name, agent.serverConfig.Address)
		serverErrors <- agent.httpServer.ListenAndServe()
	}()

	// Wait for either context cancellation, signal, or server error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	case sig := <-sigChan:
		log.Printf("\nReceived signal: %v", sig)
		return agent.Stop()
	case <-serverCtx.Done():
		return agent.Stop()
	}
}

// Stop gracefully shuts down the HTTP server with a 5-second timeout
func (agent *Agent) Stop() error {
	if agent.httpServer == nil {
		return fmt.Errorf("server is not running")
	}

	//log.Println("Shutting down server gracefully...")
	log.Printf("[%s] Shutting down server gracefully...", agent.Name)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := agent.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error during shutdown: %w", err)
	}

	//log.Println("Server stopped")
	log.Printf("[%s] Server stopped", agent.Name)
	return nil
}

func displayConversationHistory(agent *Agent) {
	// For debugging: print conversation history
	shouldIDisplay := env.GetEnvOrDefault("LOG_MESSAGES", "false")

	if conversion.StringToBool(shouldIDisplay) {

		fmt.Println()
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println("🗒️ Conversation history:")
		for _, msg := range agent.Messages {
			content := msg.Content[0].Text
			if len(content) > 80 {
				fmt.Println("📝", msg.Role, ":", content[:80]+"...")
			} else {
				fmt.Println("📝", msg.Role, ":", content)
			}
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}

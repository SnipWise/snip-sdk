package chat

/*
This is a simple agent
the conversation history is stored in memory
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/snipwise/snip-sdk/snip"
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/models"
	openaihelpers "github.com/snipwise/snip-sdk/snip/openai-helpers"
	"github.com/snipwise/snip-sdk/snip/toolbox/conversion"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

type ChatAgent struct {
	ctx                context.Context
	Name               string
	SystemInstructions string
	ModelID            string
	Messages           []*ai.Message

	Config models.ModelConfig

	genKitInstance *genkit.Genkit

	chatStreamFlowWithMemory *core.Flow[*agents.ChatRequest, *agents.ChatResponse, agents.ChatResponse]
	chatFlowWithMemory       *core.Flow[*agents.ChatRequest, *agents.ChatResponse, struct{}]

	chatFlow       *core.Flow[*agents.ChatRequest, *agents.ChatResponse, struct{}]
	chatStreamFlow *core.Flow[*agents.ChatRequest, *agents.ChatResponse, agents.ChatResponse]

	serverConfig *ConfigHTTP
	httpServer   *http.Server
	serverCancel context.CancelFunc

	// streamCancel cancels the current streaming completion
	streamCancel context.CancelFunc
	streamCtx    context.Context

	compressorAgent snip.AICompressorAgent

	logger logger.Logger
}

func NewChatAgent(
	ctx context.Context, 
	agentConfig agents.AgentConfig, 
	modelConfig models.ModelConfig, 
	opts ...AgentOption) (*ChatAgent, error) {
		
	oaiPlugin := &oai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(agentConfig.EngineURL),
		},
	}

	genKitInstance := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	// Check if model is available
	if !openaihelpers.IsModelAvailable(ctx, agentConfig.EngineURL, agentConfig.ModelID) {
		return nil, fmt.Errorf("model %s is not available at %s", agentConfig.ModelID, agentConfig.EngineURL)
	}

	agent := &ChatAgent{
		Name:               agentConfig.Name,
		SystemInstructions: agentConfig.SystemInstructions,
		ModelID:            agentConfig.ModelID,
		Messages:           []*ai.Message{},
		Config:             modelConfig,

		ctx:            ctx,
		genKitInstance: genKitInstance,
		logger:         logger.GetLoggerFromEnvWithPrefix(agentConfig.Name), // Default logger from env
	}

	// Apply all options (can override logger)
	for _, opt := range opts {
		opt(agent)
	}

	// Log model availability
	agent.logger.Info("✅ Model %s is available at %s", agentConfig.ModelID, agentConfig.EngineURL)

	return agent, nil

}

func (agent *ChatAgent) GetName() string {
	return agent.Name
}

func (agent *ChatAgent) GetMessages() []*ai.Message {
	return agent.Messages
}

func (agent *ChatAgent) GetCurrentContextSize() int {
	totalContextSize := len(agent.SystemInstructions)
	for _, msg := range agent.Messages {
		for _, content := range msg.Content {
			totalContextSize += len(content.Text)
		}
	}
	return totalContextSize
}

func (agent *ChatAgent) Kind() agents.AgentKind {
	return agents.Chat
}

func (agent *ChatAgent) AddSystemMessage(context string) error {
	// Add a system message to the conversation history
	agent.Messages = append(agent.Messages, ai.NewSystemTextMessage(strings.TrimSpace(context)))
	return nil
}

func (agent *ChatAgent) ReplaceMessagesWith(messages []*ai.Message) error {
	// Replace the entire conversation history with new messages
	if messages == nil {
		return fmt.Errorf("messages cannot be nil")
	}
	agent.Messages = messages
	return nil
}

func (agent *ChatAgent) ReplaceMessagesWithSystemMessages(systemMessages []string) error {
	// Replace the entire conversation history with system messages
	if systemMessages == nil {
		return fmt.Errorf("systemMessages cannot be nil")
	}

	// Create new message slice with system messages
	newMessages := make([]*ai.Message, 0, len(systemMessages))
	for _, msg := range systemMessages {
		newMessages = append(newMessages, ai.NewSystemTextMessage(strings.TrimSpace(msg)))
	}

	agent.Messages = newMessages
	return nil
}

func (agent *ChatAgent) GetInfo() (agents.AgentInfo, error) {
	return agents.AgentInfo{
		Name:    agent.Name,
		Config:  agent.Config,
		ModelID: agent.ModelID,
	}, nil
}

// IMPORTANT: this function uses the chat flow with memory
func (agent *ChatAgent) AskWithMemory(question string) (agents.ChatResponse, error) {
	if agent.chatFlowWithMemory == nil {
		return agents.ChatResponse{}, fmt.Errorf("chat flow is not initialized")
	}
	resp, err := agent.chatFlowWithMemory.Run(agent.ctx, &agents.ChatRequest{
		UserMessage: question,
	})
	if err != nil {
		return agents.ChatResponse{}, err
	}
	return *resp, nil

}

// IMPORTANT: this function uses the chat stream flow with memory
func (agent *ChatAgent) AskStreamWithMemory(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	if agent.chatStreamFlowWithMemory == nil {
		return agents.ChatResponse{}, fmt.Errorf("chat stream flow is not initialized")
	}
	// Streaming channel of results
	streamCh := agent.chatStreamFlowWithMemory.Stream(agent.ctx, &agents.ChatRequest{
		UserMessage: question,
	})

	finalAnswer := ""
	var finalResponse agents.ChatResponse
	for result, err := range streamCh {
		// Check for errors from the stream
		if err != nil {
			// Return both the partial answer and the error
			return agents.ChatResponse{Text: finalAnswer}, fmt.Errorf("streaming error: %w", err)
		}

		// Check for nil result (defensive programming)
		if result == nil {
			continue
		}

		if !result.Done {
			finalAnswer += result.Stream.Text
			err := callback(result.Stream)
			if err != nil {
				return agents.ChatResponse{Text: finalAnswer}, err
			}
		} else {
			// Store the final response with all metadata
			finalResponse = *result.Output
		}
	}

	return finalResponse, nil
}

// IMPORTANT: this function uses the chat flow WITHOUT memory
func (agent *ChatAgent) Ask(question string) (agents.ChatResponse, error) {
	if agent.chatFlow == nil {
		return agents.ChatResponse{}, fmt.Errorf("chat flow is not initialized")
	}
	resp, err := agent.chatFlow.Run(agent.ctx, &agents.ChatRequest{
		UserMessage: question,
	})
	if err != nil {
		return agents.ChatResponse{}, err
	}
	return *resp, nil
}

// IMPORTANT: this function uses the chat stream flow WITHOUT memory
func (agent *ChatAgent) AskStream(question string, callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	if agent.chatStreamFlow == nil {
		return agents.ChatResponse{}, fmt.Errorf("chat stream flow is not initialized")
	}
	// Streaming channel of results
	streamCh := agent.chatStreamFlow.Stream(agent.ctx, &agents.ChatRequest{
		UserMessage: question,
	})

	finalAnswer := ""
	var finalResponse agents.ChatResponse
	for result, err := range streamCh {
		// Check for errors from the stream
		if err != nil {
			// Return both the partial answer and the error
			return agents.ChatResponse{Text: finalAnswer}, fmt.Errorf("streaming error: %w", err)
		}

		// Check for nil result (defensive programming)
		if result == nil {
			continue
		}

		if !result.Done {
			finalAnswer += result.Stream.Text
			err := callback(result.Stream)
			if err != nil {
				return agents.ChatResponse{Text: finalAnswer}, err
			}
		} else {
			// Store the final response with all metadata
			finalResponse = *result.Output
		}
	}

	return finalResponse, nil
}

// Serve starts the HTTP server with the configured endpoints for the agent's flows
// The server automatically handles SIGINT (Ctrl+C) and SIGTERM signals for graceful shutdown
// Use the Stop() method to manually shutdown the server
func (agent *ChatAgent) Serve() error {
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
	agent.logger.Info("Registered endpoint: GET %s", healthcheckPath)

	// Register agent information endpoint
	informationPath := agent.serverConfig.InformationPath
	mux.HandleFunc("GET "+informationPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		info := agents.AgentInfo{
			Name:    agent.Name,
			ModelID: agent.ModelID,
			Config:  agent.Config,
		}
		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	agent.logger.Info("Registered endpoint: GET %s", informationPath)

	// IMPORTANT: with memory flows
	// Register chat flow endpoint if available
	if agent.chatFlowWithMemory != nil && agent.serverConfig.ChatFlowHandler != nil {
		chatFlowPath := agent.serverConfig.ChatFlowPath
		mux.HandleFunc("POST "+chatFlowPath, agent.serverConfig.ChatFlowHandler)
		agent.logger.Info("Registered endpoint: POST %s", chatFlowPath)
	}
	// IMPORTANT: with memory flows
	// Register chat stream flow endpoint if available
	if agent.chatStreamFlowWithMemory != nil && agent.serverConfig.ChatStreamFlowHandler != nil {
		chatStreamFlowPath := agent.serverConfig.ChatStreamFlowPath
		mux.HandleFunc("POST "+chatStreamFlowPath, agent.serverConfig.ChatStreamFlowHandler)
		agent.logger.Info("Registered endpoint: POST %s", chatStreamFlowPath)
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

			agent.logger.Info("Shutdown requested via HTTP endpoint")

			// Trigger shutdown asynchronously to allow response to be sent
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()
		})
		agent.logger.Info("Registered endpoint: POST %s", shutdownPath)
	}

	// Register cancel stream endpoint
	cancelStreamPath := agent.serverConfig.CancelStreamPath
	if cancelStreamPath != "" {
		mux.HandleFunc("POST "+cancelStreamPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if agent.streamCancel != nil {
				agent.streamCancel()
				agent.logger.Info("Streaming completion cancelled via HTTP endpoint")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"stream cancelled"}`))
			} else {
				agent.logger.Info("No active stream to cancel")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"no active stream"}`))
			}
		})
		agent.logger.Info("Registered endpoint: POST %s", cancelStreamPath)
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
				agent.logger.Error("Error decoding add context request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"invalid request body"}`))
				return
			}

			// Add context to messages
			if err := agent.AddSystemMessage(req.Context); err != nil {
				agent.logger.Error("Error adding context to messages: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to add context"}`))
				return
			}

			agent.logger.Info("Context added to messages via HTTP endpoint")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		})
		agent.logger.Info("Registered endpoint: POST %s", addContextPath)
	}

	// Register get messages endpoint
	getMessagesPath := agent.serverConfig.GetMessagesPath
	if getMessagesPath != "" {
		mux.HandleFunc("GET "+getMessagesPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Encode messages to JSON
			if err := json.NewEncoder(w).Encode(agent.Messages); err != nil {
				agent.logger.Error("Error encoding messages: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"failed to encode messages"}`))
				return
			}

			agent.logger.Debug("Messages retrieved via HTTP endpoint")
		})
		agent.logger.Info("Registered endpoint: GET %s", getMessagesPath)
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
		agent.logger.Info("Starting HTTP server on %s (Press Ctrl+C to stop)", agent.serverConfig.Address)
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
		agent.logger.Info("Received signal: %v", sig)
		return agent.Stop()
	case <-serverCtx.Done():
		return agent.Stop()
	}
}

// Stop gracefully shuts down the HTTP server with a 5-second timeout
func (agent *ChatAgent) Stop() error {
	if agent.httpServer == nil {
		return fmt.Errorf("server is not running")
	}

	agent.logger.Info("Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := agent.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error during shutdown: %w", err)
	}

	agent.logger.Info("Server stopped")
	return nil
}

// CompressContext compresses the conversation history using the configured compressor agent
// Returns an error if no compressor agent is configured
// The compression result is returned as a ChatResponse
// After compression, the agent's messages are replaced with a single system message containing the compressed context
func (agent *ChatAgent) CompressContext() (agents.ChatResponse, error) {
	if agent.compressorAgent == nil {
		return agents.ChatResponse{}, fmt.Errorf("no compressor agent configured, use EnableContextCompression option")
	}

	response, err := agent.compressorAgent.CompressMessages(agent.Messages)
	if err != nil {
		return agents.ChatResponse{}, err
	}

	// Replace the agent's messages with the compressed context
	compressedMessages := []*ai.Message{
		ai.NewSystemTextMessage(strings.TrimSpace(response.Text)),
	}
	if err := agent.ReplaceMessagesWith(compressedMessages); err != nil {
		return agents.ChatResponse{}, err
	}

	return response, nil
}

// CompressContextStream compresses the conversation history using streaming with the configured compressor agent
// Returns an error if no compressor agent is configured
// The callback function is called for each streamed chunk
// The final compression result is returned as a ChatResponse
// After compression, the agent's messages are replaced with a single system message containing the compressed context
func (agent *ChatAgent) CompressContextStream(callback func(agents.ChatResponse) error) (agents.ChatResponse, error) {
	if agent.compressorAgent == nil {
		return agents.ChatResponse{}, fmt.Errorf("no compressor agent configured, use EnableContextCompression option")
	}

	response, err := agent.compressorAgent.CompressMessagesStream(agent.Messages, callback)
	if err != nil {
		return agents.ChatResponse{}, err
	}

	// Replace the agent's messages with the compressed context
	compressedMessages := []*ai.Message{
		ai.NewSystemTextMessage(strings.TrimSpace(response.Text)),
	}
	if err := agent.ReplaceMessagesWith(compressedMessages); err != nil {
		return agents.ChatResponse{}, err
	}

	return response, nil
}

func displayConversationHistory(agent *ChatAgent) {
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

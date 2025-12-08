package chat

/*
This is a simple agent
the conversation history is stored in memory
*/

import (
	"context"
	"fmt"
	"strings"

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

	// streamCancel cancels the current streaming completion
	streamCancel context.CancelFunc
	streamCtx    context.Context

	logger logger.Logger
}

func NewChatAgent(
	ctx context.Context,
	agentConfig agents.AgentConfig,
	modelConfig models.ModelConfig,
	opts ...ChatAgentOption) (*ChatAgent, error) {

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

func (agent *ChatAgent) GetStreamCancel() context.CancelFunc {
	return agent.streamCancel
}

func (agent *ChatAgent) GetName() string {
	return agent.Name
}

func (agent *ChatAgent) Kind() agents.AgentKind {
	return agents.Chat
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

func (agent *ChatAgent) GetChatFlowWithMemory() *core.Flow[*agents.ChatRequest, *agents.ChatResponse, struct{}] {
	return agent.chatFlowWithMemory
}

func (agent *ChatAgent) GetChatStreamFlowWithMemory() *core.Flow[*agents.ChatRequest, *agents.ChatResponse, agents.ChatResponse] {
	return agent.chatStreamFlowWithMemory
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

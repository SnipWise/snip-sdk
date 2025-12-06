package chat

import (
	"github.com/snipwise/snip-sdk/snip/agents"

	"context"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// EnableChatFlowWithMemory initializes the chat flow for the agent
func EnableChatFlowWithMemory() AgentOption {
	return func(agent *ChatAgent) {
		initializeChatFlowWithMemory(agent)
	}
}

func initializeChatFlowWithMemory(agent *ChatAgent) {
	chatFlowWithMemory := genkit.DefineFlow(agent.genKitInstance, agent.Name+"-chat-flow-with-memory",
		func(ctx context.Context, input *agents.ChatRequest) (*agents.ChatResponse, error) {

			// === COMPLETION ===
			resp, err := genkit.Generate(ctx, agent.genKitInstance,
				ai.WithModelName("openai/"+agent.ModelID),
				ai.WithSystem(agent.SystemInstructions),
				ai.WithPrompt(input.UserMessage),
				ai.WithConfig(agent.Config.ToOpenAIParams()),
				ai.WithMessages(
					agent.Messages...,
				),
			)
			if err != nil {
				return nil, err
			}
			// === CONVERSATIONAL MEMORY ===

			// USER MESSAGE: append user message to history
			agent.Messages = append(agent.Messages, ai.NewUserTextMessage(strings.TrimSpace(input.UserMessage)))
			// ASSISTANT MESSAGE: append assistant response to history
			agent.Messages = append(agent.Messages, ai.NewModelTextMessage(strings.TrimSpace(resp.Text())))

			// DEBUG: print conversation history
			displayConversationHistory(agent)

			return &agents.ChatResponse{
				Text:          resp.Text(),
				FinishReason:  string(resp.FinishReason),
				FinishMessage: resp.FinishMessage,
			}, nil
		})

	agent.chatFlowWithMemory = chatFlowWithMemory

}

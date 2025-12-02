package smart

import (
	"context"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// EnableChatFlowWithMemory initializes the chat flow for the agent
func EnableChatFlowWithMemory() AgentOption {
	return func(agent *Agent) {
		initializeChatFlow(agent)
	}
}

func initializeChatFlow(agent *Agent) {
	chatFlow := genkit.DefineFlow(agent.genKitInstance, agent.Name+"-chat-flow",
		func(ctx context.Context, input *ChatRequest) (*ChatResponse, error) {

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

			return &ChatResponse{Response: resp.Text()}, nil
		})
	agent.chatFlow = chatFlow

}
package chat

import (
	"context"

	"github.com/snipwise/snip-sdk/snip/agents"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// EnableChatFlowWithMemory initializes the chat flow for the agent
func EnableChatFlow() ChatAgentOption {
	return func(agent *ChatAgent) {
		initializeChatFlow(agent)
	}
}

func initializeChatFlow(agent *ChatAgent) {
	chatFlow := genkit.DefineFlow(agent.genKitInstance, agent.Name+"-chat-flow",
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

			return &agents.ChatResponse{
				Text:             resp.Text(),
				Content:          resp.Message.Content,
				Role:             resp.Message.Role,
				FinishReason:     string(resp.FinishReason),
				FinishMessage:    resp.FinishMessage,
				ReasoningContent: resp.Reasoning(),
			}, nil
		})

	agent.chatFlow = chatFlow

}

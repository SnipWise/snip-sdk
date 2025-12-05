package snip

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// EnableChatFlowWithMemory initializes the chat flow for the agent
func EnableChatFlow() AgentOption {
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
				// ai.WithMessages(
				// 	agent.Messages...,
				// ),
			)
			if err != nil {
				return nil, err
			}

			return &ChatResponse{
				Text:          resp.Text(),
				FinishReason:  string(resp.FinishReason),
				FinishMessage: resp.FinishMessage,
			}, nil
		})

	agent.chatFlow = chatFlow

}

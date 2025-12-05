package snip

import (
	"context"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
)

// EnableChatStreamFlowWithMemory initializes the chat stream flow for the agent
func EnableChatStreamFlowWithMemory() AgentOption {
	return func(agent *Agent) {
		initializeChatStreamFlow(agent)
	}
}

func initializeChatStreamFlow(agent *Agent) {

	chatStreamFlow := genkit.DefineStreamingFlow(agent.genKitInstance, agent.Name+"-chat-stream-flow",
		func(ctx context.Context, input *ChatRequest, callback core.StreamCallback[string]) (*ChatResponse, error) {

			// Create a cancellable context for this streaming request
			streamCtx, streamCancel := context.WithCancel(ctx)
			agent.streamCtx = streamCtx
			agent.streamCancel = streamCancel

			// === COMPLETION ===
			resp, err := genkit.Generate(streamCtx, agent.genKitInstance,
				ai.WithModelName("openai/"+agent.ModelID),
				ai.WithSystem(agent.SystemInstructions),
				ai.WithPrompt(input.UserMessage),
				ai.WithConfig(agent.Config.ToOpenAIParams()),
				ai.WithMessages(
					agent.Messages...,
				),
				ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
					// Check if the context has been cancelled
					select {
					case <-streamCtx.Done():
						return streamCtx.Err()
					default:
						return callback(ctx, chunk.Text())
					}
				}),
			)
			if err != nil {
				// Clean up the stream context
				agent.streamCancel = nil
				agent.streamCtx = nil
				return nil, err
			}
			// === CONVERSATIONAL MEMORY ===

			// USER MESSAGE: append user message to history
			agent.Messages = append(agent.Messages, ai.NewUserTextMessage(strings.TrimSpace(input.UserMessage)))
			// ASSISTANT MESSAGE: append assistant response to history
			agent.Messages = append(agent.Messages, ai.NewModelTextMessage(strings.TrimSpace(resp.Text())))

			// DEBUG: print conversation history
			displayConversationHistory(agent)

			// Clean up the stream context
			agent.streamCancel = nil
			agent.streamCtx = nil

			return &ChatResponse{Response: resp.Text()}, nil
		})
	agent.chatStreamFlow = chatStreamFlow
}
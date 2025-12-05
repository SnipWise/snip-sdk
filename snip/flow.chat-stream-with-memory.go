package snip

import (
	"context"
	"fmt"
	"log"
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
		func(ctx context.Context, input *ChatRequest, callback core.StreamCallback[ChatResponse]) (*ChatResponse, error) {

			// Create a cancellable context for this streaming request
			streamCtx, streamCancel := context.WithCancel(ctx)
			agent.streamCtx = streamCtx
			agent.streamCancel = streamCancel

			// === DEBUG: CONTEXT SIZE ===
			// Log total context size for debugging
			totalContextSize := len(agent.SystemInstructions) + len(input.UserMessage)
			for _, msg := range agent.Messages {
				for _, content := range msg.Content {
					totalContextSize += len(content.Text)
				}
			}
			log.Printf("[%s] Total context size: %d characters, %d messages in history",
				agent.Name, totalContextSize, len(agent.Messages))

			// === End of DEBUG: CONTEXT SIZE ===

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
						// Send ChatResponse with the chunk text
						return callback(ctx, ChatResponse{
							Text: chunk.Text(),
						})
					}
				}),
			)
			if err != nil {
				// Clean up the stream context
				agent.streamCancel = nil
				agent.streamCtx = nil

				// Log detailed error information
				log.Printf("[%s] ❌ Generation error: %v", agent.Name, err)
				log.Printf("[%s] Context details - Total size: %d chars, Messages: %d",
					agent.Name, totalContextSize, len(agent.Messages))

				return nil, fmt.Errorf("generation failed (context size: %d chars): %w", totalContextSize, err)
			}

			// Send a final callback with complete metadata (FinishReason and FinishMessage)
			finalChunk := ChatResponse{
				Text:          "", // Empty text since all text was already streamed
				FinishReason:  string(resp.FinishReason),
				FinishMessage: resp.FinishMessage,
			}
			if callbackErr := callback(ctx, finalChunk); callbackErr != nil {
				log.Printf("[%s] ⚠️ Error in final callback: %v", agent.Name, callbackErr)
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

			return &ChatResponse{
				Text: resp.Text(),
				FinishReason: string(resp.FinishReason),
				FinishMessage: resp.FinishMessage,
			}, nil
		})
	agent.chatStreamFlow = chatStreamFlow
}

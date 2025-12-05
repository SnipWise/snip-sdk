package main

import (
	"context"
	"fmt"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/snip"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	agent0, err := snip.NewAgent(ctx,
		snip.AgentConfig{
			Name:               "Local Agent",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		snip.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		snip.EnableChatStreamFlowWithMemory(),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	finalResponse, err := agent0.AskStream("What is the capital of France?",
		func(chunk snip.ChatResponse) error {
			// During streaming, Text contains the chunk content
			// At the end, a final chunk is sent with FinishReason and FinishMessage
			if chunk.Text != "" {
				fmt.Print(chunk.Text)
			}
			// Check if this is the final chunk with metadata
			if chunk.FinishReason != "" {
				fmt.Printf("\n\n[Final chunk received]")
				fmt.Printf("\n  - FinishReason: %s", chunk.FinishReason)
				fmt.Printf("\n  - FinishMessage: %s", chunk.FinishMessage)
			}
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}
	fmt.Println()
	fmt.Println(finalResponse)

	fmt.Println("\n---")

	finalResponse, err = agent0.AskStream("What is the capital of Belgium?",
		func(chunk snip.ChatResponse) error {
			// During streaming, Text contains the chunk content
			// At the end, a final chunk is sent with FinishReason and FinishMessage
			if chunk.Text != "" {
				fmt.Print(chunk.Text)
			}
			// Check if this is the final chunk with metadata
			if chunk.FinishReason != "" {
				fmt.Printf("\n\n[Final chunk received]")
				fmt.Printf("\n  - FinishReason: %s", chunk.FinishReason)
				fmt.Printf("\n  - FinishMessage: %s", chunk.FinishMessage)
			}
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println(finalResponse)
}

package main

import (
	"context"
	"fmt"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/smart"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	agent0, err := smart.NewAgent(ctx,
		smart.AgentConfig{
			Name:               "Local Agent",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		smart.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		smart.EnableChatStreamFlowWithMemory(),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	_, err = agent0.AskStream("What is the capital of France?",
		func(chunk string) error {
			fmt.Print(chunk)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}

	fmt.Println("\n---")

	_, err = agent0.AskStream("What is the capital of Belgium?",
		func(chunk string) error {
			fmt.Print(chunk)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}

	fmt.Println()
}

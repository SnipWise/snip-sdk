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
		snip.EnableChatFlowWithMemory(),
		snip.EnableChatStreamFlowWithMemory(),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	//agent0.ModelID = "ai/qwen2.5:latest"

	response, err := agent0.AskWithMemory("What is the capital of France?")
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}
	fmt.Printf("Response from Local Agent: %s\n", response)

	_, err = agent0.AskStreamWithMemory("What is the capital of Belgium?",
		func(chunk snip.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}
	fmt.Println()
}
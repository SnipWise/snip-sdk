package main

import (
	"context"
	"fmt"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/files"

	"github.com/snipwise/snip-sdk/smart"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/lucy-gguf:q4_k_m")
	//chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:latest")

	systemInstructions, err := files.ReadTextFile(env.GetEnvOrDefault("SYSTEM_INSTRUCTION_PATH", "./system-instructions.md"))
	if err != nil {
		fmt.Printf("Error reading system instructions: %v\n", err)
		return
	}
	knowledgeBase, err := files.ReadTextFile(env.GetEnvOrDefault("KNOWLEDGE_BASE_PATH", "./knowledge-base.md"))
	if err != nil {
		fmt.Printf("Error reading knowledge base: %v\n", err)
		return
	}

	agent0 := smart.NewAgent(ctx,
		"Bob_Agentic_Agent",
		systemInstructions,
		chatModelId,
		engineURL,
		smart.Config{
			Temperature: 0.5,
			TopP:        0.9,
		},
		smart.EnableChatStreamFlowWithMemory(),
	)

	_, err = agent0.AskStream("What is the best pizza of the world?",
		func(chunk string) error {
			fmt.Print(chunk)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}

	fmt.Println("\n--- Adding knowledge base context ---")

	err = agent0.AddSystemMessage(knowledgeBase)
	if err != nil {
		fmt.Printf("Error adding system message: %v\n", err)
		return
	}
	_, err = agent0.AskStream("Who invented Hawaiian pizza?",
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

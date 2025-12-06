package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/snipwise/snip-sdk/env"
	"github.com/snipwise/snip-sdk/files"

	"github.com/snipwise/snip-sdk/snip"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	//chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/lucy-gguf:q4_k_m")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	//compressorModelId := env.GetEnvOrDefault("COMPRESSOR_MODEL", "hf.co/menlo/jan-nano-128k-gguf:q4_k_m")
	compressorModelId := env.GetEnvOrDefault("COMPRESSOR_MODEL", "ai/qwen2.5:1.5B-F16")


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

	// Create a compressor agent for context compression
	compressor, err := snip.NewCompressorAgent(
		ctx,
		snip.AgentConfig{
			Name:               "MessageCompressor",
			SystemInstructions: "",
			ModelID:            compressorModelId,
			EngineURL:          engineURL,
		},
		snip.ModelConfig{
			Temperature: 0.3, // Lower temperature for more consistent compression
		},
	)
	if err != nil {
		fmt.Printf("Error creating compressor agent: %v\n", err)
		return
	}

	compressor.SetCompressionPrompt(snip.DefaultCompressionPrompts.Minimalist)

	// Create the main chat agent
	agent0, err := snip.NewAgent(ctx,
		snip.AgentConfig{
			Name:               "Bob_Agentic_Agent",
			SystemInstructions: systemInstructions,
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		snip.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		snip.EnableChatStreamFlowWithMemory(),
		snip.EnableContextCompression(compressor),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	// run a go routine to dispaly a waiting animation

	// Add the knowledge base as a system message
	err = agent0.AddSystemMessage(knowledgeBase)
	if err != nil {
		fmt.Printf("Error adding knowledge base to agent: %v\n", err)
		return
	}


	// First question
	fmt.Println("=== First Question ===")
	answer, err := agent0.AskStreamWithMemory("What is the best pizza of the world?",
		func(chunk snip.ChatResponse) error {
			// si text pas vide et anim pas stop -> stop anim
			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}

	fmt.Println("\n✋ FinishReason:", answer.FinishReason)
	if answer.IsFinishReasonLength() {
		fmt.Println("⚠️ The answer was cut off due to length limits.")
	}
	if answer.IsFinishReasonStop() {
		fmt.Println("✅ The answer was completed successfully.")
	}
	fmt.Println(strings.Repeat("*", 50))
	fmt.Println("📏 Current Context Size:", agent0.GetCurrentContextSize())
	fmt.Println(strings.Repeat("*", 50))

	resp, err := agent0.CompressContextStream(func(chunk snip.ChatResponse) error {
		fmt.Print(chunk.Text)
		return nil
	})
	if err != nil {
		fmt.Printf("Error compressing context: %v\n", err)
		return
	}
	fmt.Println("\n✅ Compressed Context Summary:", resp.Text)
	fmt.Println(strings.Repeat("*", 50))
	fmt.Println("📏 Current Context Size After Compression:", agent0.GetCurrentContextSize())
	fmt.Println(strings.Repeat("*", 50))

	// Second question
	fmt.Println("=== Second Question ===")
	answer, err = agent0.AskStreamWithMemory("Who invented Hawaiian pizza?",
		func(chunk snip.ChatResponse) error {
			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		return
	}

	fmt.Println("\n✋ FinishReason:", answer.FinishReason)
	if answer.IsFinishReasonLength() {
		fmt.Println("⚠️ The answer was cut off due to length limits.")
	}
	if answer.IsFinishReasonStop() {
		fmt.Println("✅ The answer was completed successfully.")
	}
	fmt.Println(strings.Repeat("*", 50))
	fmt.Println("📏 Current Context Size:", agent0.GetCurrentContextSize())
	fmt.Println(strings.Repeat("*", 50))

	resp, err = agent0.CompressContextStream(func(chunk snip.ChatResponse) error {
		fmt.Print(chunk.Text)
		return nil
	})
	if err != nil {
		fmt.Printf("Error compressing context: %v\n", err)
		return
	}
	fmt.Println("\n✅ Compressed Context Summary:", resp.Text)
	fmt.Println(strings.Repeat("*", 50))
	fmt.Println("📏 Current Context Size After Compression:", agent0.GetCurrentContextSize())
	fmt.Println(strings.Repeat("*", 50))

}

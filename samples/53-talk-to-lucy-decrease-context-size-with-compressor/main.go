package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/toolbox/files"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/chat"
	"github.com/snipwise/snip-sdk/snip/compressor"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/ui/spinner"
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
	compressorAgent, err := compressor.NewCompressorAgent(
		ctx,
		agents.AgentConfig{
			Name:               "MessageCompressor",
			SystemInstructions: "",
			ModelID:            compressorModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.3, // Lower temperature for more consistent compression
		},
	)
	if err != nil {
		fmt.Printf("Error creating compressor agent: %v\n", err)
		return
	}

	// Override the default compression prompt with a minimalist one
	compressorAgent.SetCompressionPrompt(compressor.DefaultCompressionPrompts.Minimalist)

	// Create the main chat agent
	agent0, err := chat.NewChatAgent(ctx,
		agents.AgentConfig{
			Name:               "Bob_Agentic_Agent",
			SystemInstructions: systemInstructions,
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		chat.EnableChatStreamFlowWithMemory(),
		chat.EnableContextCompression(compressorAgent),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	// Add the knowledge base as a system message
	err = agent0.AddSystemMessage(knowledgeBase)
	if err != nil {
		fmt.Printf("Error adding knowledge base to agent: %v\n", err)
		return
	}

	// First question
	thinkingSpinner := spinner.New("").SetSuffix("thinking...").SetFrames(spinner.FramesDots)
	thinkingSpinner.Start()
	fmt.Println("=== First Question ===")
	answer, err := agent0.AskStreamWithMemory("What is the best pizza of the world?",
		func(chunk agents.ChatResponse) error {

			if thinkingSpinner.IsRunning() && chunk.FinishReason == "" {
				thinkingSpinner.Success("Let's go!")
				thinkingSpinner.Stop()
			}

			fmt.Print(chunk.Text)
			return nil
		},
	)
	if err != nil {
		if thinkingSpinner.IsRunning() {
			thinkingSpinner.Error("Failed to get response!")
		}
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

	compressSpinner := spinner.New("").SetSuffix("compressing...").SetFrames(spinner.FramesCircle)
	compressSpinner.Start()

	resp, err := agent0.CompressContextStream(func(chunk agents.ChatResponse) error {
		//fmt.Print(chunk.Text)
		if compressSpinner.IsRunning() && chunk.FinishReason == "stop" {
			compressSpinner.Success("Done compressing!")
			compressSpinner.Stop()
		}
		return nil
	})
	if err != nil {
		if compressSpinner.IsRunning() {
			compressSpinner.Error("Failed to compress context!")
		}
		fmt.Printf("Error compressing context: %v\n", err)
		return
	}

	//fmt.Println("\n✅ Compressed Context Summary:", resp.Text)
	fmt.Println(strings.Repeat("*", 50))
	fmt.Println("📏 Current Context Size After Compression:", agent0.GetCurrentContextSize())
	fmt.Println(strings.Repeat("*", 50))

	// Second question
	fmt.Println("=== Second Question ===")
	answer, err = agent0.AskStreamWithMemory("Who invented Hawaiian pizza?",
		func(chunk agents.ChatResponse) error {
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

	compressSpinner.SetFrames(spinner.FramesPulsingStar)
	compressSpinner.Start()

	resp, err = agent0.CompressContextStream(func(chunk agents.ChatResponse) error {
		//fmt.Print(chunk.Text)
		if compressSpinner.IsRunning() && chunk.FinishReason == "stop" {
			compressSpinner.Success("Done compressing!")
			compressSpinner.Stop()
		}
		return nil
	})
	if err != nil {
		if compressSpinner.IsRunning() {
			compressSpinner.Error("Failed to compress context!")
		}
		fmt.Printf("Error compressing context: %v\n", err)
		return
	}
	fmt.Println("\n✅ Compressed Context Summary:", resp.Text)
	fmt.Println(strings.Repeat("*", 50))
	fmt.Println("📏 Current Context Size After Compression:", agent0.GetCurrentContextSize())
	fmt.Println(strings.Repeat("*", 50))

}

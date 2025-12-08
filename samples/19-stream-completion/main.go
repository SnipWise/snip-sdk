package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/chat"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/ui/display"
	"github.com/snipwise/snip-sdk/snip/ui/spinner"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	//chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/lucy-gguf:q4_k_m")
	//chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "ai/qwen3:latest")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	loadingSpinner := spinner.NewWithColor("").SetSuffix("loading model...").SetFrames(spinner.FramesPulsingStar)
	loadingSpinner.SetSuffixColor(spinner.ColorBold).SetFrameColor(spinner.ColorYellow)

	thinkingSpinner := spinner.NewWithColor("").SetSuffix("thinking...").SetFrames(spinner.FramesDots)
	thinkingSpinner.SetSuffixColor(spinner.ColorPurple).SetFrameColor(spinner.ColorRed)

	agent0, err := chat.NewChatAgent(ctx,
		agents.AgentConfig{
			Name:               "Riker",
			SystemInstructions: "You are a Star Trek expert assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature:     0.5,
			TopP:            0.9,
			ReasoningEffort: "high",
		},
		chat.EnableChatStreamFlow(),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	loadingSpinner.Start()

	finalResponse, err := agent0.AskStream("/think Who is James T Kirk?",
		func(chunk agents.ChatResponse) error {
			if loadingSpinner.IsRunning() && chunk.FinishReason == "" {
				loadingSpinner.Success("Model loaded!")
				loadingSpinner.Stop()
			}

			// During streaming, Text contains the chunk content
			// At the end, a final chunk is sent with FinishReason and FinishMessage
			if chunk.Text != "" {
				//fmt.Println(chunk.GetContent(), chunk.Text)
				fmt.Print(chunk.Text)
				for _, part := range chunk.Content {
					if part.IsReasoning() {
						display.Color("Green text ", display.ColorGreen)
					}
				}

			}

			// Check if this is the final chunk with metadata
			if chunk.FinishReason != "" {
				fmt.Println()
				fmt.Println(strings.Repeat("-", 60))
				fmt.Println("[Final chunk received]")
				fmt.Println("  - FinishReason:", chunk.FinishReason)
				fmt.Println(strings.Repeat("-", 60))
			}
			return nil
		},
	)
	if err != nil {
		if loadingSpinner.IsRunning() {
			loadingSpinner.Error("Failed to load model!")
		}
		fmt.Printf("Error asking question: %v\n", err)
		return
	}
	fmt.Println("✋ FinishReason", finalResponse.FinishReason)

	//fmt.Println(finalResponse)
	display.Color("Reasoning content:\n"+finalResponse.ReasoningContent, display.ColorGreen)

	fmt.Println(strings.Repeat("-", 60))

}

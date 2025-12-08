package main

import (
	"context"
	"fmt"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/chatserver"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/ui/spinner"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	thinkingSpinner := spinner.New("").SetSuffix("thinking...").SetFrames(spinner.FramesDots)
	generatingSpinner := spinner.New("").SetSuffix("generating...").SetFrames(spinner.FramesPulsingStar)

	agent0, err := chatserver.NewChatAgentServer(ctx,
		agents.AgentConfig{
			Name:               "Local Agent",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	generatingSpinner.Start()

	response, err := agent0.AskWithMemory("What is the capital of France?")
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		generatingSpinner.Error("Failed!")
		return
	}
	generatingSpinner.Success("Done!")

	fmt.Printf("Response from Local Agent: %s\n", response.Text)

	thinkingSpinner.Start()

	_, err = agent0.AskStreamWithMemory("What is the capital of Belgium? And tell me something about its history.",
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
		fmt.Printf("Error asking question: %v\n", err)
		return
	}
	fmt.Println()
}

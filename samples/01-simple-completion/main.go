package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/chat"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/ui/spinner"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	generatingSpinner := spinner.New("").SetSuffix("generating...").SetFrames(spinner.FramesPulsingStar)

	agent0, err := chat.NewChatAgent(ctx,
		agents.AgentConfig{
			Name:               "Riker",
			SystemInstructions: "You are a Star Trek expert assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		chat.EnableChatFlowWithMemory(),
	)
	if err != nil {
		fmt.Printf("Error creating agent: %v\n", err)
		return
	}

	generatingSpinner.Start()

	response, err := agent0.AskWithMemory("What is the name of the ship of Jean-Luc Picard?")
	if err != nil {
		fmt.Printf("Error asking question: %v\n", err)
		generatingSpinner.Error("Failed!")
		return
	}

	generatingSpinner.Success("Done!")

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📝 Response from Local %s Agent:\n", agent0.Kind())
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(response.Text)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("✋ FinishReason", response.FinishReason)
	fmt.Println(strings.Repeat("=", 60))

}

/*
	sp := spinner.New("waiting...")
	sp.Start()
	time.Sleep(3 * time.Second)
	sp.Success("Done!")

	if sp.IsRunning() {
		sp.Stop()
	}
	//sp.Stop()

	sp.SetSuffix("tada").SetPrefix("").SetFrames(spinner.FramesBraille)
	sp.Start()
	time.Sleep(3 * time.Second)
	sp.Error("Failed!")
*/

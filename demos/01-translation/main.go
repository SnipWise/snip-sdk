package main

import (
	"context"
	"fmt"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/chat"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/ui/display"

	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/toolbox/files"
	"github.com/snipwise/snip-sdk/snip/ui/prompt"
	"github.com/snipwise/snip-sdk/snip/ui/spinner"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	chatModelId := env.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:1.5B-F16")

	display.Titlef("Document Translation Assistant [%s]", chatModelId)
	
	systemInstructions := `
	You are a translation assistant that helps users translate documents French to English.
	`
	userInstructions := `
	Read the french document above.
    Translate its content from French to English.
	`

	loadingSpinner := spinner.NewWithColor("").SetSuffix("loading model...").SetFrames(spinner.FramesPulsingStar)
	loadingSpinner.SetSuffixColor(spinner.ColorBold).SetFrameColor(spinner.ColorYellow)

	thinkingSpinner := spinner.NewWithColor("").SetSuffix("thinking...").SetFrames(spinner.FramesDots)
	thinkingSpinner.SetSuffixColor(spinner.ColorPurple).SetFrameColor(spinner.ColorRed)

	filesNames, err := files.GetAllFilesInDirectory("./documents")

	if err != nil {
		display.Error(err.Error())
		return
	}
	choices := []prompt.Choice{}

	for _, fileName := range filesNames {
		choices = append(choices, prompt.Choice{
			Value: fileName,
			Label: fileName,
		})

	}
	selectPrompt := prompt.NewColorSelect("Choose a document to translate:", choices).
		SetColors(
			prompt.ColorBrightCyan,   // message color
			prompt.ColorWhite,        // choice color
			prompt.ColorBrightYellow, // default color
			prompt.ColorGray,         // number color
			prompt.ColorRed,          // error color
		).
		SetSymbols("❯", "★", "✗")

	selected, err := selectPrompt.Run()
	if err != nil {
		display.Error(err.Error())
		return
	}

	fmt.Printf(
		"%s✓ You selected: %s%s%s\n",
		prompt.ColorGreen, prompt.ColorBold, selected, prompt.ColorReset,
	)

	content, err := files.ReadTextFile(selected)
	if err != nil {
		display.Errorf("Failed when loadinf file: %s\n", err.Error())
		return
	}

	loadingSpinner.Start()

	translator, err := chat.NewChatAgent(ctx,
		agents.AgentConfig{
			Name:               "Bob",
			SystemInstructions: systemInstructions,
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		chat.EnableChatStreamFlow(),
	)

	if err != nil {
		if loadingSpinner.IsRunning() {
			loadingSpinner.Error("Error creating agent!" + err.Error())
		}
		return
	}

	_ = translator.AddSystemMessage("DOCUMENT:\n" + content)

	_, err = translator.AskStream(userInstructions,
		func(chunk agents.ChatResponse) error {
			if loadingSpinner.IsRunning() && chunk.FinishReason == "" {
				loadingSpinner.Success("Model loaded! Starting translation...")
				loadingSpinner.Stop()
			}
			if chunk.Text != "" {
				fmt.Print(chunk.Text)
			}
			return nil
		},
	)
	if err != nil {
		display.Errorf("Failed with translation: %s\n", err.Error())
		return
	}

}

package main

import (
	"context"
	"fmt"

	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/structured"

	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/ui/spinner"
)

func main() {

	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	modelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	generatingSpinner := spinner.NewWithColor("").SetSuffix("generating...").SetFrames(spinner.FramesPulsingStar)
	generatingSpinner.SetFrameColor(spinner.ColorBrightRed)

	// Structure for final flow output
	type Intent struct {
		Action string `json:"intent"`
		//Action    string `json:"action"`
		Character string `json:"name"`
		Known     bool   `json:"known"`
	}

	dungeonMaster, err := structured.NewStructuredAgent[Intent](
		ctx,
		agents.AgentConfig{
			Name: "DungeonMaster",
			SystemInstructions: `
			You are helping the dungeon master of a D&D game.
			Detect if the user want to speak to one of the following NPCs: 
			Thrain (dwarf blacksmith), 
			Liora (elven mage), 
			Galdor (human rogue), 
			Elara (halfling ranger), 
			Shesepankh (tiefling warlock).

			If the user's message does not explicitly mention wanting to speak to one of these NPCs, respond with:
			action: speak
			character: <NPC name>
			known: false

			Otherwise, respond with:
			action: speak
			character: <NPC name> 
			Where <NPC name> is the name of the NPC the user wants to speak to: Thrain, Liora, Galdor, Elara, or Shesepankh.
			known: true			
			`,
			ModelID:   modelId,
			EngineURL: engineURL,
		},
		models.ModelConfig{
			Temperature: 0.7,
			TopP:        0.9,
		},
	)
	if err != nil {
		fmt.Printf("Error creating structured agent: %v\n", err)
		return
	}

	testMessages := []string{
		"I want to chat with Thrain and learn about his blacksmith skills.",
		"I want to meet a dwarf blacksmith.",
		"I want to speak about spells and magic.",
		"I want to speak to Bob Morane.",
	}

	generatingSpinner.Start()

	for _, message := range testMessages {

		intent, err := dungeonMaster.GenerateStructuredData(message)

		if err != nil {
			generatingSpinner.Error("Failed!")
			fmt.Printf("😡 Error running intent flow: %v\n", err)
			return
		}
		if !intent.Known {
			generatingSpinner.Error("Unknown!")
			fmt.Println("🙀 NPC", intent.Character, "not recognized!")
			continue
		}
		generatingSpinner.Success("Done!")

		fmt.Println("🙂 Detected Intent: action ->", intent.Action, "character ->", intent.Character)
		generatingSpinner.Start()

	}

}

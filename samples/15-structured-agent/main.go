package main

import (
	"context"
	"fmt"
	"strings"

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
	type NonPlayerCharacter struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		Race        string `json:"race"`
		Background  string `json:"background"`
		Personality string `json:"personality"`
		Abilities   string `json:"abilities"`
	}

	dungeonMaster, err := structured.NewStructuredAgent[NonPlayerCharacter](
		ctx,
		agents.AgentConfig{
			Name:               "DungeonMaster",
			SystemInstructions: "You are the dungeon master of a D&D game.",
			ModelID:            modelId,
			EngineURL:          engineURL,
		},
		models.ModelConfig{
			Temperature: 0.7,
			TopP:        0.9,
		},
		//structured.WithLogLevel[NonPlayerCharacter](logger.LevelDebug),
	)
	if err != nil {
		fmt.Printf("Error creating structured agent: %v\n", err)
		return
	}

	generatingSpinner.Start()

	npc, err := dungeonMaster.GenerateStructuredData("Generate a D&D NPC Elf name and all its characteristics.")
	if err != nil {
		fmt.Printf("Error generating structured data: %v\n", err)
		generatingSpinner.Error("Failed!")
		return
	}

	generatingSpinner.Success("Done!")

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📝 Response from Structured %s Agent:\n", dungeonMaster.Kind())
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("Generated D&D NPC Elf:")
	fmt.Println("Name:", npc.Name)
	fmt.Println("Role:", npc.Role)
	fmt.Println("Race:", npc.Race)
	fmt.Println("Personality:", npc.Personality)
	fmt.Println("Abilities:", npc.Abilities)
	fmt.Println("Background:", npc.Background)
	fmt.Println(strings.Repeat("=", 60))

}

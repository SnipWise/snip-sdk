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
		Action    string `json:"intent"`
		Character string `json:"character"`
		Topic     string `json:"topic"`
		Initiator string `json:"initiator"` // "player" or "npc"
	}

	// ✋ This time we expect a **slice** of Intent
	dungeonMaster, err := structured.NewStructuredAgent[[]Intent](
		ctx,
		agents.AgentConfig{
			Name: "DungeonMaster",
			SystemInstructions: `
			You are helping the dungeon master of a D&D game.
			Detect if the user want to speak (chat/talk/meet) to one of the following NPCs: 
			Thrain (dwarf blacksmith), 
			Liora (elven mage), 
			Galdor (human rogue), 
			Elara (halfling ranger), 
			Shesepankh (tiefling warlock).

			if the user want to speak to a dwarf blacksmith, they mean Thrain.
			if the user want to speak to an elven mage, they mean Liora.
			if the user want to speak to a human rogue, they mean Galdor.
			if the user want to speak to a halfling ranger, they mean Elara.
			if the user want to speak to a tiefling warlock, they mean Shesepankh.

			If the user's message does not explicitly mention wanting to speak to one of these NPCs, respond with:
			action: speak
			character: <NPC name>
			topic: <topic>
			initiator: player (user)

			Otherwise, respond with:
			action: speak
			character: <NPC name> 
			topic: <topic>
			initiator: player (user)
			Where <NPC name> is the name of the NPC the user wants to speak to: Thrain, Liora, Galdor, Elara, or Shesepankh.
						
			Detect if the user wants to assign a task (perform/do/prepare/help/fix/scout/research/walk) to the NPC.
			
			NPCs available:
			- Thrain (dwarf blacksmith)
			- Liora (elven mage)
			- Galdor (human rogue)
			- Elara (halfling ranger)
			- Shesepankh (tiefling warlock)

			Respond with:
			action: task
			character: <NPC name> 
			topic: <task>
			initiator: npc
			Where <NPC name> is the name of the NPC the user wants to speak to: Thrain, Liora, Galdor, Elara, or Shesepankh.
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

	generatingSpinner.Start()

	intents, err := dungeonMaster.GenerateStructuredData(`
		I want to chat with Shesepankh about dark magic.
			Then Shesepankh needs to research forbidden spells.
		I must speak with Elara about our next adventure.
			Then Elara needs to scout the northern woods.
		We need to speak with Liora about ancient artifacts.
			Liora walks to the ancient ruins.
			Liora prepares protective charms.
	`)

	if err != nil {
		generatingSpinner.Error("Failed!")
		fmt.Printf("😡 Error running intent flow: %v\n", err)
		return
	}
	generatingSpinner.Success("Done!")

	for idx, intent := range *intents {

		switch intent.Action {
		case "speak":
			fmt.Printf("%d [SPEAK] %s wants to talk with %s about: %s\n",
				idx, intent.Initiator, intent.Character, intent.Topic)
		case "task":
			fmt.Printf("%d   - [TASK] %s will %s\n",
				idx, intent.Character, intent.Topic)
		default:
			fmt.Printf("%d [UNKNOWN] action: %s, character: %s, topic: %s, initiator: %s\n",
				idx, intent.Action, intent.Character, intent.Topic, intent.Initiator)
		}
	}
}

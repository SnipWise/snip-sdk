package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/snipwise/snip-sdk/snip/agents"
	"github.com/snipwise/snip-sdk/snip/models"
	"github.com/snipwise/snip-sdk/snip/toolbox/env"
	"github.com/snipwise/snip-sdk/snip/toolbox/logger"
	"github.com/snipwise/snip-sdk/snip/tools"
)

func main() {
	ctx := context.Background()
	engineURL := env.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	toolModelId := env.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	// Create a tools agent
	dungeonMaster, err := tools.NewToolsAgent(
		ctx,
		agents.AgentConfig{
			Name: "DungeonMaster",
			SystemInstructions: `
				You are a helpful D&D assistant that can roll dice and generate character names.
				Use the appropriate tools when asked to roll dice or generate character names.
			`,
			ModelID:   toolModelId,
			EngineURL: engineURL,
		},
		models.ModelConfig{
			Temperature: 0.0,
		},
		tools.EnableToolCallFlow(),
		tools.WithLogLevel(logger.LevelDebug),
	)
	if err != nil {
		log.Fatalf("Error creating tools agent: %v", err)
	}

	tools.AddToolToAgent(
		dungeonMaster,
		"roll_dice",
		"Roll n dice with n faces each",
		func(ctx *ai.ToolContext, input DiceRollInput) (DiceRollResult, error) {
			return rollDice(input.NumDice, input.NumFaces), nil
		},
	)
	tools.AddToolToAgent(
		dungeonMaster,
		"generate_character_name",
		"Generate a D&D character name for a specific race",
		func(ctx *ai.ToolContext, input CharacterNameInput) (CharacterNameResult, error) {
			return generateCharacterName(input.Race), nil
		},
	)

	response, err := dungeonMaster.RunToolCalls(
		`
		Roll 3 dices with 6 faces each.
		Then generate a character name for an elf.
		Finally, roll 2 dices with 8 faces each.
		After that, generate a character name for a dwarf.
		`,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Final Response:", response)

}

type DiceRollInput struct {
	NumDice  int `json:"num_dice"`
	NumFaces int `json:"num_faces"`
}

type DiceRollResult struct {
	Rolls []int `json:"rolls"`
	Total int   `json:"total"`
}

type CharacterNameInput struct {
	Race string `json:"race"`
}

type CharacterNameResult struct {
	Name string `json:"name"`
	Race string `json:"race"`
}

func rollDice(numDice, numFaces int) DiceRollResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rolls := make([]int, numDice)
	total := 0

	for i := 0; i < numDice; i++ {
		roll := r.Intn(numFaces) + 1
		rolls[i] = roll
		total += roll
	}

	return DiceRollResult{
		Rolls: rolls,
		Total: total,
	}
}

func generateCharacterName(race string) CharacterNameResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	namesByRace := map[string][]string{
		"elf":      {"Aerdrie", "Ahvonna", "Aramil", "Aranea", "Berrian", "Caelynn", "Carric", "Dayereth", "Enna", "Galinndan"},
		"dwarf":    {"Adrik", "Baern", "Darrak", "Eberk", "Fargrim", "Gardain", "Harbek", "Kildrak", "Morgran", "Thorek"},
		"human":    {"Aerdrie", "Aramil", "Berris", "Cithreth", "Dayereth", "Enna", "Galinndan", "Hadarai", "Immeral", "Lamlis"},
		"halfling": {"Alton", "Ander", "Bernie", "Bobbin", "Cade", "Callus", "Corrin", "Dannad", "Garret", "Lindal"},
		"orc":      {"Gash", "Gell", "Henk", "Holg", "Imsh", "Keth", "Krusk", "Mhurren", "Ront", "Shump"},
		"tiefling": {"Akmenos", "Amnon", "Barakas", "Damakos", "Ekemon", "Iados", "Kairon", "Leucis", "Melech", "Mordai"},
	}

	raceLower := strings.ToLower(race)
	names, exists := namesByRace[raceLower]
	if !exists {
		names = namesByRace["human"] // Default to human names
	}

	selectedName := names[r.Intn(len(names))]

	return CharacterNameResult{
		Name: selectedName,
		Race: race,
	}
}

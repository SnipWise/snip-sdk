package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

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
		tools.EnableAutoToolCallFlow(),
		tools.WithLogLevel(logger.LevelDebug), // ✋ logs activated
		tools.WithToolExecution(tools.ToolExecution{
			OnExecuted: func(toolName string, toolInput any, toolCallRef string, output any, err error) {
				if err != nil {
					log.Printf("❌ Tool %q execution failed: %v", toolName, err)
					return
				}
				log.Printf("✅ Tool %q executed successfully with input: %v | output: %v", toolName, toolInput, output)
			},
		}),
	)
	if err != nil {
		log.Fatalf("Error creating tools agent: %v", err)
	}

	tools.AddToolToAgent(
		dungeonMaster,
		"roll_dice",
		"Roll n dice with n faces each",
		func(input DiceRollInput) (DiceRollResult, error) {
			return rollDice(input.NumDice, input.NumFaces), nil
		},
	)
	tools.AddToolToAgent(
		dungeonMaster,
		"generate_character_name",
		"Generate a D&D character name for a specific race",
		func(input CharacterNameInput) (CharacterNameResult, error) {
			return generateCharacterName(input.Race), nil
		},
	)

	// TODO: RuntToolCallsWithCallback
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

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Tool calls:", response.List)

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Final Response:!")
	fmt.Println(response.Text)
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n📋 Tool Calls Results:")

	// Display tool call results
	for _, toolResult := range response.List {
		for toolName, output := range toolResult {
			fmt.Printf("\n🔧 Tool: %s\n", toolName)

			switch toolName {
			case "roll_dice":
				result, err := tools.Transform[DiceRollResult](output)
				if err != nil {
					fmt.Printf("❌ Error transforming dice roll: %v\n", err)
					continue
				}
				fmt.Printf("   🎲 Rolls: %v | Total: %d\n", result.Rolls, result.Total)

			case "generate_character_name":
				result, err := tools.Transform[CharacterNameResult](output)
				if err != nil {
					fmt.Printf("❌ Error transforming character name: %v\n", err)
					continue
				}
				fmt.Printf("   🧙 %s character: %s\n", result.Race, result.Name)
			}
		}
	}

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

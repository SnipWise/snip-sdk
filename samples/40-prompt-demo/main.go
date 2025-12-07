package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/snipwise/snip-sdk/snip/ui/prompt"
)

func main() {
	fmt.Println("=== Prompt Package Demo ===")

	// 1. Simple input
	fmt.Println("1. Simple Input")
	input := prompt.New("What is your name?")
	name, err := input.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Hello, %s!\n\n", name)

	// 2. Input with default value
	fmt.Println("2. Input with Default Value")
	input2 := prompt.New("What is your favorite color?").
		SetDefault("blue")
	color, err := input2.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Your favorite color is %s\n\n", color)

	// 3. Input with validation
	fmt.Println("3. Input with Validation")
	emailInput := prompt.New("Enter your email").
		SetValidator(func(value string) error {
			if !strings.Contains(value, "@") {
				return fmt.Errorf("invalid email format - must contain @")
			}
			return nil
		})
	email, err := emailInput.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Email: %s\n\n", email)

	// 4. Simple confirmation
	fmt.Println("4. Simple Confirmation")
	confirm := prompt.NewConfirm("Do you want to continue?")
	result, err := confirm.Run()
	if err != nil {
		log.Fatal(err)
	}
	if result {
		fmt.Println("✓ Continuing...")
	} else {
		fmt.Println("✗ Cancelled.")
		return
	}

	// 5. Select (single choice)
	fmt.Println("5. Select (Single Choice)")
	choices := []prompt.Choice{
		{Label: "Python", Value: "python"},
		{Label: "Go", Value: "go"},
		{Label: "JavaScript", Value: "javascript"},
		{Label: "Rust", Value: "rust"},
	}
	selectPrompt := prompt.NewSelect("Choose your favorite language", choices).
		SetDefault("go")
	selected, err := selectPrompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ You selected: %s\n\n", selected)

	// 6. Multi-choice
	fmt.Println("6. Multi-Choice")
	toolChoices := []prompt.Choice{
		{Label: "Docker", Value: "docker"},
		{Label: "Kubernetes", Value: "k8s"},
		{Label: "Terraform", Value: "terraform"},
		{Label: "Ansible", Value: "ansible"},
	}
	multiPrompt := prompt.NewMultiChoice("Select the tools you use (comma-separated)", toolChoices).
		SetDefaults([]string{"docker"})
	selectedTools, err := multiPrompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ You selected: %v\n\n", selectedTools)

	// 7. Tool execution confirmation (like in flow.too.call.go)
	fmt.Println("7. Tool Execution Confirmation")
	toolPrompt := prompt.NewToolExecutionPrompt("get_weather")
	decision, err := toolPrompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch decision {
	case prompt.ToolExecute:
		fmt.Println("✓ Executing tool 'get_weather'...")
		fmt.Println("  (Simulated) Weather: Sunny, 25°C")
	case prompt.ToolSkip:
		fmt.Println("✗ Tool execution skipped.")
	case prompt.ToolQuit:
		fmt.Println("✗ Exiting the program.")
		return
	}

	fmt.Println("\n=== Demo Complete ===")
}

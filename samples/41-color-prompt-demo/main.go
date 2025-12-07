package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/snipwise/snip-sdk/snip/ui/prompt"
)

func main() {
	fmt.Println("=== Color Prompt Package Demo ===")

	// 1. Simple colored input
	fmt.Println("\n1. Simple Colored Input")
	input := prompt.NewWithColor("What is your name?").
		SetMessageColor(prompt.ColorBrightCyan).
		SetInputColor(prompt.ColorBrightWhite)

	name, err := input.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s%s Hello, %s!%s\n",
		prompt.ColorGreen, prompt.ColorBold, name, prompt.ColorReset)

	// 2. Input with default value and custom colors
	fmt.Println("\n2. Input with Default Value and Custom Colors")
	colorInput := prompt.NewWithColor("What is your favorite color?").
		SetDefault("blue").
		SetColors(
			prompt.ColorBrightMagenta, // message color
			prompt.ColorGray,          // default color
			prompt.ColorYellow,        // input color
			prompt.ColorRed,           // error color
			prompt.ColorGreen,         // success color
		)

	color, err := colorInput.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s✓ Your favorite color is %s%s\n",
		prompt.ColorGreen, color, prompt.ColorReset)

	// 3. Input with validation and custom symbols
	fmt.Println("\n3. Input with Validation and Custom Symbols")
	emailInput := prompt.NewWithColor("Enter your email").
		SetValidator(func(value string) error {
			if !strings.Contains(value, "@") {
				return fmt.Errorf("invalid email format - must contain @")
			}
			return nil
		}).
		SetMessageColor(prompt.ColorBrightBlue).
		SetErrorColor(prompt.ColorBrightRed).
		SetSymbols("→", "✓", "✗")

	email, err := emailInput.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s✓ Email: %s%s\n",
		prompt.ColorGreen, email, prompt.ColorReset)

	// 4. Colored confirmation
	fmt.Println("\n4. Colored Confirmation")
	confirm := prompt.NewColorConfirm("Do you want to continue?").
		SetMessageColor(prompt.ColorBrightYellow).
		SetSuccessColor(prompt.ColorGreen).
		SetErrorColor(prompt.ColorRed)

	result, err := confirm.Run()
	if err != nil {
		log.Fatal(err)
	}
	if result {
		fmt.Printf("%s✓ Continuing...%s\n", prompt.ColorGreen, prompt.ColorReset)
	} else {
		fmt.Printf("%s✗ Cancelled.%s\n", prompt.ColorRed, prompt.ColorReset)
		return
	}

	// 5. Colored select (single choice)
	fmt.Println("\n5. Colored Select (Single Choice)")
	choices := []prompt.Choice{
		{Label: "Python", Value: "python"},
		{Label: "Go", Value: "go"},
		{Label: "JavaScript", Value: "javascript"},
		{Label: "Rust", Value: "rust"},
	}

	selectPrompt := prompt.NewColorSelect("Choose your favorite language", choices).
		SetDefault("go").
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
		log.Fatal(err)
	}
	fmt.Printf("%s✓ You selected: %s%s%s\n",
		prompt.ColorGreen, prompt.ColorBold, selected, prompt.ColorReset)

	// 6. Colored multi-choice
	fmt.Println("\n6. Colored Multi-Choice")
	toolChoices := []prompt.Choice{
		{Label: "Docker", Value: "docker"},
		{Label: "Kubernetes", Value: "k8s"},
		{Label: "Terraform", Value: "terraform"},
		{Label: "Ansible", Value: "ansible"},
	}

	multiPrompt := prompt.NewColorMultiChoice("Select the tools you use", toolChoices).
		SetDefaults([]string{"docker"}).
		SetMessageColor(prompt.ColorBrightMagenta).
		SetChoiceColor(prompt.ColorWhite).
		SetDefaultColor(prompt.ColorBrightGreen).
		SetNumberColor(prompt.ColorGray)

	selectedTools, err := multiPrompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s✓ You selected: %v%s\n",
		prompt.ColorGreen, selectedTools, prompt.ColorReset)

	// 7. Rainbow style example
	fmt.Println("\n7. Rainbow Style Example")
	rainbowInput := prompt.NewWithColor("What's your favorite emoji?").
		SetMessageColor(prompt.ColorBrightMagenta).
		SetInputColor(prompt.ColorBrightYellow).
		SetSymbols("🌈", "✨", "❌")

	emoji, err := rainbowInput.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s✨ Nice choice: %s%s\n",
		prompt.ColorBrightMagenta, emoji, prompt.ColorReset)

	// 8. Professional style (minimal colors)
	fmt.Println("\n8. Professional Style (Minimal Colors)")
	proInput := prompt.NewWithColor("Enter your company name").
		SetMessageColor(prompt.ColorBlue).
		SetInputColor(prompt.ColorWhite).
		SetDefaultColor(prompt.ColorGray).
		SetErrorColor(prompt.ColorRed).
		SetDefault("Acme Corp")

	company, err := proInput.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s✓ Company: %s%s\n",
		prompt.ColorBlue, company, prompt.ColorReset)


	fmt.Printf("\n%s=== Demo Complete ===%s\n",
		prompt.ColorBrightGreen, prompt.ColorReset)
}

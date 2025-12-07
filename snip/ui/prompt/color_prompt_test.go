package prompt

import (
	"fmt"
	"strings"
	"testing"
)

// TestColorInputCreation tests that colored input prompts can be created
func TestColorInputCreation(t *testing.T) {
	input := NewWithColor("What is your name?")
	if input == nil {
		t.Fatal("Expected input to be created")
	}
	if input.message != "What is your name?" {
		t.Errorf("Expected message to be 'What is your name?', got '%s'", input.message)
	}
	// Check default colors
	if input.messageColor != ColorCyan {
		t.Errorf("Expected default message color to be ColorCyan")
	}
	if input.promptSymbol != "❯" {
		t.Errorf("Expected default prompt symbol to be '❯', got '%s'", input.promptSymbol)
	}
}

// TestColorInputWithColors tests color customization
func TestColorInputWithColors(t *testing.T) {
	input := NewWithColor("Test").
		SetMessageColor(ColorRed).
		SetDefaultColor(ColorGreen).
		SetInputColor(ColorBlue).
		SetErrorColor(ColorYellow).
		SetSuccessColor(ColorMagenta)

	if input.messageColor != ColorRed {
		t.Error("Message color not set correctly")
	}
	if input.defaultColor != ColorGreen {
		t.Error("Default color not set correctly")
	}
	if input.inputColor != ColorBlue {
		t.Error("Input color not set correctly")
	}
	if input.errorColor != ColorYellow {
		t.Error("Error color not set correctly")
	}
	if input.successColor != ColorMagenta {
		t.Error("Success color not set correctly")
	}
}

// TestColorInputSetColorsMethod tests the SetColors method
func TestColorInputSetColorsMethod(t *testing.T) {
	input := NewWithColor("Test").
		SetColors(ColorRed, ColorGreen, ColorBlue, ColorYellow, ColorMagenta)

	if input.messageColor != ColorRed {
		t.Error("Message color not set correctly via SetColors")
	}
	if input.defaultColor != ColorGreen {
		t.Error("Default color not set correctly via SetColors")
	}
	if input.inputColor != ColorBlue {
		t.Error("Input color not set correctly via SetColors")
	}
	if input.errorColor != ColorYellow {
		t.Error("Error color not set correctly via SetColors")
	}
	if input.successColor != ColorMagenta {
		t.Error("Success color not set correctly via SetColors")
	}
}

// TestColorInputSetSymbols tests custom symbols
func TestColorInputSetSymbols(t *testing.T) {
	input := NewWithColor("Test").
		SetSymbols("→", "✓", "✗")

	if input.promptSymbol != "→" {
		t.Errorf("Expected prompt symbol to be '→', got '%s'", input.promptSymbol)
	}
	if input.successSymbol != "✓" {
		t.Errorf("Expected success symbol to be '✓', got '%s'", input.successSymbol)
	}
	if input.errorSymbol != "✗" {
		t.Errorf("Expected error symbol to be '✗', got '%s'", input.errorSymbol)
	}
}

// TestColorConfirmCreation tests colored confirm creation
func TestColorConfirmCreation(t *testing.T) {
	confirm := NewColorConfirm("Continue?")
	if confirm == nil {
		t.Fatal("Expected confirm to be created")
	}
	if confirm.message != "Continue?" {
		t.Errorf("Expected message to be 'Continue?', got '%s'", confirm.message)
	}
	if confirm.messageColor != ColorCyan {
		t.Error("Expected default message color to be ColorCyan")
	}
}

// TestColorConfirmWithColors tests color customization for confirm
func TestColorConfirmWithColors(t *testing.T) {
	confirm := NewColorConfirm("Test").
		SetColors(ColorRed, ColorGreen, ColorBlue, ColorYellow)

	if confirm.messageColor != ColorRed {
		t.Error("Message color not set correctly")
	}
	if confirm.optionColor != ColorGreen {
		t.Error("Option color not set correctly")
	}
	if confirm.successColor != ColorBlue {
		t.Error("Success color not set correctly")
	}
	if confirm.errorColor != ColorYellow {
		t.Error("Error color not set correctly")
	}
}

// TestColorSelectCreation tests colored select creation
func TestColorSelectCreation(t *testing.T) {
	choices := []Choice{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}
	selectPrompt := NewColorSelect("Choose", choices)

	if selectPrompt == nil {
		t.Fatal("Expected select to be created")
	}
	if len(selectPrompt.choices) != 2 {
		t.Errorf("Expected 2 choices, got %d", len(selectPrompt.choices))
	}
	if selectPrompt.messageColor != ColorCyan {
		t.Error("Expected default message color to be ColorCyan")
	}
	if selectPrompt.defaultSymbol != "●" {
		t.Errorf("Expected default symbol to be '●', got '%s'", selectPrompt.defaultSymbol)
	}
}

// TestColorSelectWithColors tests color customization for select
func TestColorSelectWithColors(t *testing.T) {
	choices := []Choice{{Label: "Test", Value: "test"}}
	selectPrompt := NewColorSelect("Test", choices).
		SetColors(ColorRed, ColorGreen, ColorBlue, ColorYellow, ColorMagenta)

	if selectPrompt.messageColor != ColorRed {
		t.Error("Message color not set correctly")
	}
	if selectPrompt.choiceColor != ColorGreen {
		t.Error("Choice color not set correctly")
	}
	if selectPrompt.defaultColor != ColorBlue {
		t.Error("Default color not set correctly")
	}
	if selectPrompt.numberColor != ColorYellow {
		t.Error("Number color not set correctly")
	}
	if selectPrompt.errorColor != ColorMagenta {
		t.Error("Error color not set correctly")
	}
}

// TestColorMultiChoiceCreation tests colored multi-choice creation
func TestColorMultiChoiceCreation(t *testing.T) {
	choices := []Choice{
		{Label: "Tool 1", Value: "t1"},
		{Label: "Tool 2", Value: "t2"},
	}
	multiPrompt := NewColorMultiChoice("Select tools", choices)

	if multiPrompt == nil {
		t.Fatal("Expected multi-choice to be created")
	}
	if len(multiPrompt.choices) != 2 {
		t.Errorf("Expected 2 choices, got %d", len(multiPrompt.choices))
	}
	if multiPrompt.messageColor != ColorCyan {
		t.Error("Expected default message color to be ColorCyan")
	}
}

// TestColorMultiChoiceWithColors tests color customization for multi-choice
func TestColorMultiChoiceWithColors(t *testing.T) {
	choices := []Choice{{Label: "Test", Value: "test"}}
	multiPrompt := NewColorMultiChoice("Test", choices).
		SetColors(ColorRed, ColorGreen, ColorBlue, ColorYellow, ColorMagenta)

	if multiPrompt.messageColor != ColorRed {
		t.Error("Message color not set correctly")
	}
	if multiPrompt.choiceColor != ColorGreen {
		t.Error("Choice color not set correctly")
	}
	if multiPrompt.defaultColor != ColorBlue {
		t.Error("Default color not set correctly")
	}
	if multiPrompt.numberColor != ColorYellow {
		t.Error("Number color not set correctly")
	}
	if multiPrompt.errorColor != ColorMagenta {
		t.Error("Error color not set correctly")
	}
}

// TestColorConstants tests that color constants are defined
func TestColorConstants(t *testing.T) {
	colors := []string{
		ColorReset, ColorBold, ColorDim, ColorItalic, ColorUnderline,
		ColorBlack, ColorRed, ColorGreen, ColorYellow, ColorBlue,
		ColorMagenta, ColorCyan, ColorWhite, ColorGray,
		ColorBrightRed, ColorBrightGreen, ColorBrightYellow,
		ColorBrightBlue, ColorBrightMagenta, ColorBrightCyan,
	}

	for _, color := range colors {
		if color == "" {
			t.Error("Color constant should not be empty")
		}
		if !strings.HasPrefix(color, "\033[") {
			t.Errorf("Color constant should start with ANSI escape sequence, got: %q", color)
		}
	}
}

// TestFluentAPIColorInput tests the fluent API for ColorInput
func TestFluentAPIColorInput(t *testing.T) {
	input := NewWithColor("Test").
		SetDefault("default").
		SetValidator(func(s string) error { return nil }).
		SetMessageColor(ColorRed).
		SetDefaultColor(ColorGreen).
		SetInputColor(ColorBlue).
		SetErrorColor(ColorYellow).
		SetSuccessColor(ColorMagenta).
		SetSymbols("→", "✓", "✗")

	if input.defaultValue != "default" {
		t.Error("Fluent API failed for default value")
	}
	if input.messageColor != ColorRed {
		t.Error("Fluent API failed for message color")
	}
	if input.promptSymbol != "→" {
		t.Error("Fluent API failed for prompt symbol")
	}
}

// TestFluentAPIColorConfirm tests the fluent API for ColorConfirm
func TestFluentAPIColorConfirm(t *testing.T) {
	confirm := NewColorConfirm("Test").
		SetDefault(true).
		SetMessageColor(ColorRed).
		SetOptionColor(ColorGreen).
		SetSuccessColor(ColorBlue).
		SetErrorColor(ColorYellow).
		SetSymbols("→", "✓", "✗")

	if !confirm.defaultValue {
		t.Error("Fluent API failed for default value")
	}
	if confirm.messageColor != ColorRed {
		t.Error("Fluent API failed for message color")
	}
	if confirm.promptSymbol != "→" {
		t.Error("Fluent API failed for prompt symbol")
	}
}

// TestFluentAPIColorSelect tests the fluent API for ColorSelect
func TestFluentAPIColorSelect(t *testing.T) {
	choices := []Choice{{Label: "A", Value: "a"}}
	selectPrompt := NewColorSelect("Test", choices).
		SetDefault("a").
		SetMessageColor(ColorRed).
		SetChoiceColor(ColorGreen).
		SetDefaultColor(ColorBlue).
		SetNumberColor(ColorYellow).
		SetErrorColor(ColorMagenta).
		SetSymbols("→", "★", "✗")

	if selectPrompt.defaultValue != "a" {
		t.Error("Fluent API failed for default value")
	}
	if selectPrompt.messageColor != ColorRed {
		t.Error("Fluent API failed for message color")
	}
	if selectPrompt.defaultSymbol != "★" {
		t.Error("Fluent API failed for default symbol")
	}
}

// TestColorInputWithValidator tests that validators work with colored input
func TestColorInputWithValidator(t *testing.T) {
	validator := func(value string) error {
		if len(value) < 3 {
			return fmt.Errorf("too short")
		}
		return nil
	}

	input := NewWithColor("Test").SetValidator(validator)

	// Test the validator
	err := input.validator("ab")
	if err == nil {
		t.Error("Expected validator to return error for short value")
	}

	err = input.validator("abc")
	if err != nil {
		t.Errorf("Expected validator to accept value 'abc', got error: %v", err)
	}
}

// TestBrightColors tests bright color constants
func TestBrightColors(t *testing.T) {
	brightColors := []struct {
		name  string
		value string
	}{
		{"ColorBrightBlack", ColorBrightBlack},
		{"ColorBrightRed", ColorBrightRed},
		{"ColorBrightGreen", ColorBrightGreen},
		{"ColorBrightYellow", ColorBrightYellow},
		{"ColorBrightBlue", ColorBrightBlue},
		{"ColorBrightMagenta", ColorBrightMagenta},
		{"ColorBrightCyan", ColorBrightCyan},
		{"ColorBrightWhite", ColorBrightWhite},
	}

	for _, bc := range brightColors {
		if bc.value == "" {
			t.Errorf("%s should not be empty", bc.name)
		}
		if !strings.Contains(bc.value, "9") {
			t.Errorf("%s should contain '9' for bright colors, got: %q", bc.name, bc.value)
		}
	}
}

// TestBackgroundColors tests background color constants
func TestBackgroundColors(t *testing.T) {
	bgColors := []struct {
		name  string
		value string
	}{
		{"BgBlack", BgBlack},
		{"BgRed", BgRed},
		{"BgGreen", BgGreen},
		{"BgYellow", BgYellow},
		{"BgBlue", BgBlue},
		{"BgMagenta", BgMagenta},
		{"BgCyan", BgCyan},
		{"BgWhite", BgWhite},
	}

	for _, bg := range bgColors {
		if bg.value == "" {
			t.Errorf("%s should not be empty", bg.name)
		}
		if !strings.Contains(bg.value, "4") {
			t.Errorf("%s should contain '4' for background colors, got: %q", bg.name, bg.value)
		}
	}
}

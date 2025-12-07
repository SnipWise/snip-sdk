package prompt

import (
	"fmt"
	"strings"
	"testing"
)

// TestInputCreation tests that input prompts can be created
func TestInputCreation(t *testing.T) {
	input := New("What is your name?")
	if input == nil {
		t.Fatal("Expected input to be created")
	}
	if input.message != "What is your name?" {
		t.Errorf("Expected message to be 'What is your name?', got '%s'", input.message)
	}
}

// TestInputWithDefault tests that default values are set correctly
func TestInputWithDefault(t *testing.T) {
	input := New("What is your name?").SetDefault("John")
	if input.defaultValue != "John" {
		t.Errorf("Expected default value to be 'John', got '%s'", input.defaultValue)
	}
}

// TestInputWithValidator tests that validators are set correctly
func TestInputWithValidator(t *testing.T) {
	validator := func(value string) error {
		if len(value) < 3 {
			return fmt.Errorf("value too short")
		}
		return nil
	}

	input := New("Enter text").SetValidator(validator)
	if input.validator == nil {
		t.Fatal("Expected validator to be set")
	}

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

// TestConfirmCreation tests that confirm prompts can be created
func TestConfirmCreation(t *testing.T) {
	confirm := NewConfirm("Do you want to continue?")
	if confirm == nil {
		t.Fatal("Expected confirm to be created")
	}
	if confirm.message != "Do you want to continue?" {
		t.Errorf("Expected message to be 'Do you want to continue?', got '%s'", confirm.message)
	}
	if confirm.defaultValue != false {
		t.Error("Expected default value to be false")
	}
}

// TestConfirmWithDefault tests that confirm default values work
func TestConfirmWithDefault(t *testing.T) {
	confirm := NewConfirm("Continue?").SetDefault(true)
	if !confirm.defaultValue {
		t.Error("Expected default value to be true")
	}
}

// TestSelectCreation tests that select prompts can be created
func TestSelectCreation(t *testing.T) {
	choices := []Choice{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}
	selectPrompt := NewSelect("Choose an option", choices)

	if selectPrompt == nil {
		t.Fatal("Expected select to be created")
	}
	if len(selectPrompt.choices) != 2 {
		t.Errorf("Expected 2 choices, got %d", len(selectPrompt.choices))
	}
}

// TestSelectWithDefault tests that select default values work
func TestSelectWithDefault(t *testing.T) {
	choices := []Choice{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}
	selectPrompt := NewSelect("Choose", choices).SetDefault("opt2")

	if selectPrompt.defaultValue != "opt2" {
		t.Errorf("Expected default to be 'opt2', got '%s'", selectPrompt.defaultValue)
	}
}

// TestGetDefaultLabel tests the helper method
func TestGetDefaultLabel(t *testing.T) {
	choices := []Choice{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
		{Label: "Option 3", Value: "opt3"},
	}
	selectPrompt := NewSelect("Choose", choices).SetDefault("opt2")

	label := selectPrompt.getDefaultLabel()
	if label != "2" {
		t.Errorf("Expected default label to be '2', got '%s'", label)
	}
}

// TestMultiChoiceCreation tests that multi-choice prompts can be created
func TestMultiChoiceCreation(t *testing.T) {
	choices := []Choice{
		{Label: "Tool 1", Value: "t1"},
		{Label: "Tool 2", Value: "t2"},
	}
	multiPrompt := NewMultiChoice("Select tools", choices)

	if multiPrompt == nil {
		t.Fatal("Expected multi-choice to be created")
	}
	if len(multiPrompt.choices) != 2 {
		t.Errorf("Expected 2 choices, got %d", len(multiPrompt.choices))
	}
}

// TestMultiChoiceWithDefaults tests that multi-choice defaults work
func TestMultiChoiceWithDefaults(t *testing.T) {
	choices := []Choice{
		{Label: "Tool 1", Value: "t1"},
		{Label: "Tool 2", Value: "t2"},
		{Label: "Tool 3", Value: "t3"},
	}
	multiPrompt := NewMultiChoice("Select tools", choices).
		SetDefaults([]string{"t1", "t3"})

	if len(multiPrompt.defaultValues) != 2 {
		t.Errorf("Expected 2 default values, got %d", len(multiPrompt.defaultValues))
	}
	if multiPrompt.defaultValues[0] != "t1" || multiPrompt.defaultValues[1] != "t3" {
		t.Error("Expected default values to be ['t1', 't3']")
	}
}



// TestEmailValidator provides a reusable email validator
func TestEmailValidator(t *testing.T) {
	validator := func(value string) error {
		if !strings.Contains(value, "@") {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}

	tests := []struct {
		input     string
		wantError bool
	}{
		{"test@example.com", false},
		{"invalid", true},
		{"@example.com", false},
		{"test@", false},
		{"", true},
	}

	for _, tt := range tests {
		err := validator(tt.input)
		if (err != nil) != tt.wantError {
			t.Errorf("validator(%q) error = %v, wantError %v", tt.input, err, tt.wantError)
		}
	}
}


// TestFluentAPI tests the fluent API pattern
func TestFluentAPI(t *testing.T) {
	// Test that methods return the object for chaining
	input := New("Test").
		SetDefault("default").
		SetValidator(func(s string) error { return nil })

	if input.defaultValue != "default" {
		t.Error("Fluent API failed for Input")
	}

	confirm := NewConfirm("Test").
		SetDefault(true)

	if !confirm.defaultValue {
		t.Error("Fluent API failed for Confirm")
	}

	choices := []Choice{{Label: "A", Value: "a"}}
	selectPrompt := NewSelect("Test", choices).
		SetDefault("a")

	if selectPrompt.defaultValue != "a" {
		t.Error("Fluent API failed for Select")
	}

	multiPrompt := NewMultiChoice("Test", choices).
		SetDefaults([]string{"a"})

	if len(multiPrompt.defaultValues) != 1 {
		t.Error("Fluent API failed for MultiChoice")
	}
}

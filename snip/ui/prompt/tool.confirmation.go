package prompt

import (
	"fmt"
)

// ToolConfirmation represents the user's decision about tool execution
type ToolConfirmation int

const (
	// ToolExecute indicates the user wants to execute the tool
	ToolExecute ToolConfirmation = iota
	// ToolSkip indicates the user wants to skip the tool
	ToolSkip
	// ToolQuit indicates the user wants to quit the program
	ToolQuit
)

// ToolExecutionPrompt prompts the user to confirm tool execution
type ToolExecutionPrompt struct {
	toolName string
}

// NewToolExecutionPrompt creates a new tool execution confirmation prompt
func NewToolExecutionPrompt(toolName string) *ToolExecutionPrompt {
	return &ToolExecutionPrompt{
		toolName: toolName,
	}
}

// Run displays the tool execution prompt and returns the user's decision
func (t *ToolExecutionPrompt) Run() (ToolConfirmation, error) {
	choices := []Choice{
		{Label: "Execute the tool", Value: "execute"},
		{Label: "Skip this tool", Value: "skip"},
		{Label: "Quit the program", Value: "quit"},
	}

	selectPrompt := NewSelect(
		fmt.Sprintf("What do you want to do with tool %q?", t.toolName),
		choices,
	).SetDefault("execute")

	result, err := selectPrompt.Run()
	if err != nil {
		return ToolSkip, err
	}

	switch result {
	case "execute":
		return ToolExecute, nil
	case "skip":
		return ToolSkip, nil
	case "quit":
		return ToolQuit, nil
	default:
		return ToolSkip, nil
	}
}

// RunSimple is a simplified version using Confirm instead of Select
func (t *ToolExecutionPrompt) RunSimple() (ToolConfirmation, error) {
	confirm := NewConfirm(
		fmt.Sprintf("Do you want to execute tool %q?", t.toolName),
	).SetDefault(true)

	result, err := confirm.Run()
	if err != nil {
		return ToolSkip, err
	}

	if result {
		return ToolExecute, nil
	}
	return ToolSkip, nil
}

// ColorToolExecutionPrompt is a colored version of ToolExecutionPrompt
type ColorToolExecutionPrompt struct {
	toolName     string
	messageColor string
	choiceColor  string
	defaultColor string
}

// NewColorToolExecutionPrompt creates a new colored tool execution confirmation prompt
func NewColorToolExecutionPrompt(toolName string) *ColorToolExecutionPrompt {
	return &ColorToolExecutionPrompt{
		toolName:     toolName,
		messageColor: ColorBrightYellow,
		choiceColor:  ColorWhite,
		defaultColor: ColorBrightGreen,
	}
}

// SetMessageColor sets the color of the prompt message
func (t *ColorToolExecutionPrompt) SetMessageColor(color string) *ColorToolExecutionPrompt {
	t.messageColor = color
	return t
}

// SetChoiceColor sets the color of choice labels
func (t *ColorToolExecutionPrompt) SetChoiceColor(color string) *ColorToolExecutionPrompt {
	t.choiceColor = color
	return t
}

// SetDefaultColor sets the color of the default choice indicator
func (t *ColorToolExecutionPrompt) SetDefaultColor(color string) *ColorToolExecutionPrompt {
	t.defaultColor = color
	return t
}

// Run displays the colored tool execution prompt and returns the user's decision
func (t *ColorToolExecutionPrompt) Run() (ToolConfirmation, error) {
	choices := []Choice{
		{Label: "Execute the tool", Value: "execute"},
		{Label: "Skip this tool", Value: "skip"},
		{Label: "Quit the program", Value: "quit"},
	}

	selectPrompt := NewColorSelect(
		fmt.Sprintf("What do you want to do with tool %q?", t.toolName),
		choices,
	).SetDefault("execute").
		SetMessageColor(t.messageColor).
		SetChoiceColor(t.choiceColor).
		SetDefaultColor(t.defaultColor)

	result, err := selectPrompt.Run()
	if err != nil {
		return ToolSkip, err
	}

	switch result {
	case "execute":
		return ToolExecute, nil
	case "skip":
		return ToolSkip, nil
	case "quit":
		return ToolQuit, nil
	default:
		return ToolSkip, nil
	}
}

// RunSimple is a simplified colored version using ColorConfirm instead of ColorSelect
func (t *ColorToolExecutionPrompt) RunSimple() (ToolConfirmation, error) {
	confirm := NewColorConfirm(
		fmt.Sprintf("Do you want to execute tool %q?", t.toolName),
	).SetDefault(true).
		SetMessageColor(t.messageColor).
		SetSuccessColor(ColorGreen)

	result, err := confirm.Run()
	if err != nil {
		return ToolSkip, err
	}

	if result {
		return ToolExecute, nil
	}
	return ToolSkip, nil
}

// Example usage to replace the code in flow.too.call.go:
//
// Instead of:
//   var response string
//   for {
//       fmt.Printf("Do you want to execute tool %q? (y/n/q): ", req.Name)
//       _, err := fmt.Scanln(&response)
//       if err != nil {
//           fmt.Println("Error reading input:", err)
//           continue
//       }
//       response = strings.ToLower(strings.TrimSpace(response))
//
//       switch response {
//       case "q":
//           fmt.Println("Exiting the program.")
//           stopped = true
//           return
//       case "y":
//           // Execute tool
//       case "n":
//           // Skip tool
//       default:
//           fmt.Println("Invalid input. Please enter 'y', 'n', or 'q'.")
//       }
//   }
//
// Use:
//   toolPrompt := prompt.NewToolExecutionPrompt(req.Name)
//   decision, err := toolPrompt.Run()
//   if err != nil {
//       toolsAgent.logger.Error("Error reading input: %v", err)
//       continue
//   }
//
//   switch decision {
//   case prompt.ToolExecute:
//       output, err := tool.RunRaw(ctx, req.Input)
//       // ... handle execution
//   case prompt.ToolSkip:
//       fmt.Println("Tool execution skipped.")
//       continue
//   case prompt.ToolQuit:
//       fmt.Println("Exiting the program.")
//       stopped = true
//       return
//   }
//
// Or with colors:
//   toolPrompt := prompt.NewColorToolExecutionPrompt(req.Name).
//       SetMessageColor(prompt.ColorBrightYellow).
//       SetDefaultColor(prompt.ColorBrightGreen)
//   decision, err := toolPrompt.Run()
//   // ... same switch statement as above

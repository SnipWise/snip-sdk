package display

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput captures stdout output for testing
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestColorConstants tests that color constants are defined
func TestColorConstants(t *testing.T) {
	colors := []string{
		ColorReset, ColorBold, ColorDim, ColorItalic,
		ColorBlack, ColorRed, ColorGreen, ColorYellow,
		ColorBlue, ColorMagenta, ColorCyan, ColorWhite,
	}

	for _, color := range colors {
		if color == "" {
			t.Error("Color constant should not be empty")
		}
	}
}

// TestSymbolConstants tests that symbol constants are defined
func TestSymbolConstants(t *testing.T) {
	symbols := []string{
		SymbolSuccess, SymbolError, SymbolWarning,
		SymbolInfo, SymbolDebug, SymbolArrow, SymbolBullet,
	}

	for _, symbol := range symbols {
		if symbol == "" {
			t.Error("Symbol constant should not be empty")
		}
	}
}

// TestPrint tests the Print function
func TestPrint(t *testing.T) {
	output := captureOutput(func() {
		Print("test")
	})
	if output != "test" {
		t.Errorf("Expected 'test', got '%s'", output)
	}
}

// TestPrintln tests the Println function
func TestPrintln(t *testing.T) {
	output := captureOutput(func() {
		Println("test")
	})
	if !strings.Contains(output, "test") {
		t.Errorf("Expected output to contain 'test', got '%s'", output)
	}
}

// TestPrintf tests the Printf function
func TestPrintf(t *testing.T) {
	output := captureOutput(func() {
		Printf("test %d", 42)
	})
	if !strings.Contains(output, "test 42") {
		t.Errorf("Expected output to contain 'test 42', got '%s'", output)
	}
}

// TestColor tests the Color function
func TestColor(t *testing.T) {
	output := captureOutput(func() {
		Color("test", ColorRed)
	})
	if !strings.Contains(output, "test") {
		t.Errorf("Expected output to contain 'test', got '%s'", output)
	}
	if !strings.Contains(output, ColorRed) {
		t.Error("Expected output to contain color code")
	}
	if !strings.Contains(output, ColorReset) {
		t.Error("Expected output to contain reset code")
	}
}

// TestColorln tests the Colorln function
func TestColorln(t *testing.T) {
	output := captureOutput(func() {
		Colorln("test", ColorGreen)
	})
	if !strings.Contains(output, "test") {
		t.Errorf("Expected output to contain 'test', got '%s'", output)
	}
	if !strings.Contains(output, ColorGreen) {
		t.Error("Expected output to contain color code")
	}
}

// TestSuccess tests the Success function
func TestSuccess(t *testing.T) {
	output := captureOutput(func() {
		Success("Operation completed")
	})
	if !strings.Contains(output, "Operation completed") {
		t.Errorf("Expected output to contain 'Operation completed', got '%s'", output)
	}
	if !strings.Contains(output, SymbolSuccess) {
		t.Error("Expected output to contain success symbol")
	}
}

// TestSuccessf tests the Successf function
func TestSuccessf(t *testing.T) {
	output := captureOutput(func() {
		Successf("Completed %d tasks", 5)
	})
	if !strings.Contains(output, "Completed 5 tasks") {
		t.Errorf("Expected formatted output, got '%s'", output)
	}
}

// TestError tests the Error function
func TestError(t *testing.T) {
	output := captureOutput(func() {
		Error("Something went wrong")
	})
	if !strings.Contains(output, "Something went wrong") {
		t.Errorf("Expected output to contain error message, got '%s'", output)
	}
	if !strings.Contains(output, SymbolError) {
		t.Error("Expected output to contain error symbol")
	}
}

// TestWarning tests the Warning function
func TestWarning(t *testing.T) {
	output := captureOutput(func() {
		Warning("Be careful")
	})
	if !strings.Contains(output, "Be careful") {
		t.Errorf("Expected output to contain warning message, got '%s'", output)
	}
	if !strings.Contains(output, SymbolWarning) {
		t.Error("Expected output to contain warning symbol")
	}
}

// TestInfo tests the Info function
func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		Info("For your information")
	})
	if !strings.Contains(output, "For your information") {
		t.Errorf("Expected output to contain info message, got '%s'", output)
	}
	if !strings.Contains(output, SymbolInfo) {
		t.Error("Expected output to contain info symbol")
	}
}

// TestDebug tests the Debug function
func TestDebug(t *testing.T) {
	output := captureOutput(func() {
		Debug("Debug information")
	})
	if !strings.Contains(output, "Debug information") {
		t.Errorf("Expected output to contain debug message, got '%s'", output)
	}
	if !strings.Contains(output, SymbolDebug) {
		t.Error("Expected output to contain debug symbol")
	}
}

// TestHeader tests the Header function
func TestHeader(t *testing.T) {
	output := captureOutput(func() {
		Header("Main Title")
	})
	if !strings.Contains(output, "Main Title") {
		t.Errorf("Expected output to contain 'Main Title', got '%s'", output)
	}
}

// TestTitle tests the Title function
func TestTitle(t *testing.T) {
	output := captureOutput(func() {
		Title("Section Title")
	})
	if !strings.Contains(output, "Section Title") {
		t.Errorf("Expected output to contain 'Section Title', got '%s'", output)
	}
	// Should contain separator
	if !strings.Contains(output, "─") {
		t.Error("Expected output to contain separator")
	}
}

// TestSeparator tests the Separator function
func TestSeparator(t *testing.T) {
	output := captureOutput(func() {
		Separator()
	})
	if !strings.Contains(output, "─") {
		t.Error("Expected output to contain separator character")
	}
}

// TestBullet tests the Bullet function
func TestBullet(t *testing.T) {
	output := captureOutput(func() {
		Bullet("Item 1")
	})
	if !strings.Contains(output, "Item 1") {
		t.Errorf("Expected output to contain 'Item 1', got '%s'", output)
	}
	if !strings.Contains(output, SymbolBullet) {
		t.Error("Expected output to contain bullet symbol")
	}
}

// TestArrow tests the Arrow function
func TestArrow(t *testing.T) {
	output := captureOutput(func() {
		Arrow("Next step")
	})
	if !strings.Contains(output, "Next step") {
		t.Errorf("Expected output to contain 'Next step', got '%s'", output)
	}
	if !strings.Contains(output, SymbolArrow) {
		t.Error("Expected output to contain arrow symbol")
	}
}

// TestBox tests the Box function
func TestBox(t *testing.T) {
	output := captureOutput(func() {
		Box("Important")
	})
	if !strings.Contains(output, "Important") {
		t.Errorf("Expected output to contain 'Important', got '%s'", output)
	}
	if !strings.Contains(output, "┌") || !strings.Contains(output, "┐") {
		t.Error("Expected output to contain box corners")
	}
}

// TestBanner tests the Banner function
func TestBanner(t *testing.T) {
	output := captureOutput(func() {
		Banner("Welcome")
	})
	if !strings.Contains(output, "Welcome") {
		t.Errorf("Expected output to contain 'Welcome', got '%s'", output)
	}
	if !strings.Contains(output, "╔") || !strings.Contains(output, "╗") {
		t.Error("Expected output to contain banner borders")
	}
}

// TestList tests the List function
func TestList(t *testing.T) {
	output := captureOutput(func() {
		List(1, "First item")
	})
	if !strings.Contains(output, "1.") {
		t.Error("Expected output to contain list number")
	}
	if !strings.Contains(output, "First item") {
		t.Errorf("Expected output to contain 'First item', got '%s'", output)
	}
}

// TestStep tests the Step function
func TestStep(t *testing.T) {
	output := captureOutput(func() {
		Step(1, 3, "First step")
	})
	if !strings.Contains(output, "[1/3]") {
		t.Error("Expected output to contain step indicator [1/3]")
	}
	if !strings.Contains(output, "First step") {
		t.Errorf("Expected output to contain 'First step', got '%s'", output)
	}
}

// TestProgress tests the Progress function
func TestProgress(t *testing.T) {
	output := captureOutput(func() {
		Progress("Loading")
	})
	if !strings.Contains(output, "Loading") {
		t.Errorf("Expected output to contain 'Loading', got '%s'", output)
	}
	if !strings.Contains(output, "⏳") {
		t.Error("Expected output to contain hourglass emoji")
	}
}

// TestDone tests the Done function
func TestDone(t *testing.T) {
	output := captureOutput(func() {
		Done("Task completed")
	})
	if !strings.Contains(output, "Task completed") {
		t.Errorf("Expected output to contain 'Task completed', got '%s'", output)
	}
	if !strings.Contains(output, "✓") {
		t.Error("Expected output to contain checkmark")
	}
}

// TestIndent tests the Indent function
func TestIndent(t *testing.T) {
	output := captureOutput(func() {
		Indent(2, "Indented text")
	})
	if !strings.Contains(output, "Indented text") {
		t.Errorf("Expected output to contain 'Indented text', got '%s'", output)
	}
	// Should have 4 spaces (2 levels * 2 spaces)
	if !strings.HasPrefix(output, "    ") {
		t.Error("Expected output to be indented with 4 spaces")
	}
}

// TestTable tests the Table function
func TestTable(t *testing.T) {
	output := captureOutput(func() {
		Table("Name", "John Doe")
	})
	if !strings.Contains(output, "Name") {
		t.Error("Expected output to contain key")
	}
	if !strings.Contains(output, "John Doe") {
		t.Error("Expected output to contain value")
	}
}

// TestKeyValue tests the KeyValue function
func TestKeyValue(t *testing.T) {
	output := captureOutput(func() {
		KeyValue("Status", "Active")
	})
	if !strings.Contains(output, "Status") {
		t.Error("Expected output to contain key")
	}
	if !strings.Contains(output, "Active") {
		t.Error("Expected output to contain value")
	}
	if !strings.Contains(output, ":") {
		t.Error("Expected output to contain colon separator")
	}
}

// TestStyledFunction tests the Styled function
func TestStyledFunction(t *testing.T) {
	output := captureOutput(func() {
		Styled("Bold Text", ColorBold)
	})
	if !strings.Contains(output, "Bold Text") {
		t.Errorf("Expected output to contain 'Bold Text', got '%s'", output)
	}
	if !strings.Contains(output, ColorBold) {
		t.Error("Expected output to contain bold code")
	}
}

// TestStyledlnFunction tests the Styledln function
func TestStyledlnFunction(t *testing.T) {
	output := captureOutput(func() {
		Styledln("Styled Text", ColorItalic, ColorBlue)
	})
	if !strings.Contains(output, "Styled Text") {
		t.Errorf("Expected output to contain 'Styled Text', got '%s'", output)
	}
}

// TestObjectStructure tests object-like output
func TestObjectStructure(t *testing.T) {
	output := captureOutput(func() {
		ObjectStart("User")
		Field("name", "John")
		Field("age", "30")
		ObjectEnd()
	})
	if !strings.Contains(output, "User") {
		t.Error("Expected output to contain object name")
	}
	if !strings.Contains(output, "name") {
		t.Error("Expected output to contain field name")
	}
	if !strings.Contains(output, "John") {
		t.Error("Expected output to contain field value")
	}
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Error("Expected output to contain braces")
	}
}

// TestFormattedFunctions tests formatted function variants
func TestFormattedFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
		want string
	}{
		{"Headerf", func() { Headerf("Header %d", 1) }, "Header 1"},
		{"Subheaderf", func() { Subheaderf("Sub %s", "title") }, "Sub title"},
		{"Errorf", func() { Errorf("Error %d", 404) }, "Error 404"},
		{"Warningf", func() { Warningf("Warning %s", "test") }, "Warning test"},
		{"Infof", func() { Infof("Info %d", 42) }, "Info 42"},
		{"Debugf", func() { Debugf("Debug %s", "msg") }, "Debug msg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(tt.fn)
			if !strings.Contains(output, tt.want) {
				t.Errorf("Expected output to contain '%s', got '%s'", tt.want, output)
			}
		})
	}
}

// TestBrightColors tests bright color constants
func TestBrightColors(t *testing.T) {
	colors := []string{
		ColorBrightRed, ColorBrightGreen, ColorBrightBlue,
		ColorBrightYellow, ColorBrightCyan, ColorBrightMagenta,
	}

	for _, color := range colors {
		if color == "" {
			t.Error("Bright color constant should not be empty")
		}
		if !strings.Contains(color, "9") {
			t.Errorf("Bright color should contain '9', got: %q", color)
		}
	}
}

// TestBackgroundColors tests background color constants
func TestBackgroundColors(t *testing.T) {
	colors := []string{
		BgRed, BgGreen, BgBlue, BgYellow,
	}

	for _, color := range colors {
		if color == "" {
			t.Error("Background color constant should not be empty")
		}
		if !strings.Contains(color, "4") {
			t.Errorf("Background color should contain '4', got: %q", color)
		}
	}
}

// Example functions for documentation

func ExampleSuccess() {
	Success("Operation completed successfully")
}

func ExampleError() {
	Error("Failed to connect to server")
}

func ExampleWarning() {
	Warning("Disk space is running low")
}

func ExampleInfo() {
	Info("Server started on port 8080")
}

func ExampleHeader() {
	Header("Application Settings")
}

func ExampleTitle() {
	Title("Configuration")
}

func ExampleBox() {
	Box("Important Notice")
}

func ExampleBanner() {
	Banner("Welcome to MyApp")
}

func ExampleStep() {
	Step(1, 5, "Initializing application")
	Step(2, 5, "Loading configuration")
	Step(3, 5, "Connecting to database")
}

func ExampleTable() {
	Table("Name", "John Doe")
	Table("Email", "john@example.com")
	Table("Status", "Active")
}

func ExampleObjectStart() {
	ObjectStart("User")
	Field("id", "12345")
	Field("name", "Alice")
	Field("role", "Admin")
	ObjectEnd()
}

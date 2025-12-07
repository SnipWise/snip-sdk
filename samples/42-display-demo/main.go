package main

import (
	"github.com/snipwise/snip-sdk/snip/ui/display"
)

func main() {
	// Banner
	display.Banner("Display Package Demo")

	// 1. Basic Messages
	display.Title("1. Basic Messages")
	display.Print("This is a simple print. ")
	display.Println("This is println.")
	display.Printf("This is printf with value: %d\n", 42)
	display.NewLine()

	// 2. Colored Messages
	display.Title("2. Colored Messages")
	display.Color("Red text ", display.ColorRed)
	display.Color("Green text ", display.ColorGreen)
	display.Color("Blue text\n", display.ColorBlue)
	display.Colorln("Cyan line", display.ColorCyan)
	display.Colorln("Yellow line", display.ColorYellow)
	display.NewLine()

	// 3. Text Styles
	display.Title("3. Text Styles")
	display.Boldln("This is bold text")
	display.Italicln("This is italic text")
	display.Underlineln("This is underlined text")
	display.NewLine()

	// 4. Status Messages
	display.Title("4. Status Messages")
	display.Success("Operation completed successfully")
	display.Error("Something went wrong")
	display.Warning("This is a warning message")
	display.Info("For your information")
	display.Debug("Debug information")
	display.NewLine()

	// 5. Formatted Status Messages
	display.Title("5. Formatted Status Messages")
	display.Successf("Processed %d files", 10)
	display.Errorf("Failed with code: %d", 404)
	display.Warningf("Disk usage: %d%%", 85)
	display.Infof("Server running on port %d", 8080)
	display.Debugf("Memory usage: %dMB", 256)
	display.NewLine()

	// 6. Headers and Titles
	display.Title("6. Headers and Titles")
	display.Header("Main Header")
	display.Subheader("Subheader")
	display.Title("Section Title")
	display.NewLine()

	// 7. Lists and Bullets
	display.Title("7. Lists and Bullets")
	display.Bullet("First bullet point")
	display.Bullet("Second bullet point")
	display.Bullet("Third bullet point")
	display.NewLine()
	display.List(1, "First numbered item")
	display.List(2, "Second numbered item")
	display.List(3, "Third numbered item")
	display.NewLine()
	display.ColoredBullet("Red bullet", display.ColorRed)
	display.ColoredBullet("Green bullet", display.ColorGreen)
	display.ColoredBullet("Blue bullet", display.ColorBlue)
	display.NewLine()

	// 8. Arrows and Progress
	display.Title("8. Arrows and Progress")
	display.Arrow("Next step")
	display.Arrow("Another step")
	display.Progress("Loading data")
	display.Done("Task completed")
	display.NewLine()

	// 9. Steps
	display.Title("9. Multi-Step Process")
	display.Step(1, 5, "Initializing application")
	display.Step(2, 5, "Loading configuration")
	display.Step(3, 5, "Connecting to database")
	display.Step(4, 5, "Starting services")
	display.Step(5, 5, "Ready to serve requests")
	display.NewLine()

	// 10. Boxes and Banners
	display.Title("10. Boxes and Banners")
	display.Box("Important Notice")
	display.ColoredBox("Colored Box", display.ColorBrightMagenta)
	display.Banner("Special Announcement")

	// 11. Separators
	display.Title("11. Separators")
	display.Println("Above separator")
	display.Separator()
	display.Println("Below separator")
	display.SeparatorWithChar("=", 40)
	display.Println("Custom separator")
	display.NewLine()

	// 12. Indentation
	display.Title("12. Indentation")
	display.Println("Level 0")
	display.Indent(1, "Level 1")
	display.Indent(2, "Level 2")
	display.Indent(3, "Level 3")
	display.ColoredIndent(1, "Colored indent", display.ColorCyan)
	display.NewLine()

	// 13. Tables and Key-Value Pairs
	display.Title("13. Tables and Key-Value Pairs")
	display.Subheader("Table Format:")
	display.Table("Name", "John Doe")
	display.Table("Email", "john@example.com")
	display.Table("Status", "Active")
	display.Table("Role", "Administrator")
	display.NewLine()
	display.Subheader("Key-Value Format:")
	display.KeyValue("Server", "localhost:8080")
	display.KeyValue("Database", "PostgreSQL")
	display.KeyValue("Environment", "Production")
	display.NewLine()

	// 14. Object-like Structure
	display.Title("14. Object-like Structure")
	display.ObjectStart("User")
	display.Field("id", "12345")
	display.Field("name", "Alice Smith")
	display.Field("email", "alice@example.com")
	display.Field("role", "Developer")
	display.ObjectEnd()
	display.NewLine()

	// 15. Styled Messages
	display.Title("15. Styled Messages")
	display.Styledln("Bold Red", display.ColorBold, display.ColorRed)
	display.Styledln("Italic Cyan", display.ColorItalic, display.ColorCyan)
	display.Styledln("Underline Green", display.ColorUnderline, display.ColorGreen)
	display.Styledln("Bold Italic Blue", display.ColorBold, display.ColorItalic, display.ColorBlue)
	display.NewLine()

	// 16. Bright Colors
	display.Title("16. Bright Colors")
	display.Colorln("Bright Red", display.ColorBrightRed)
	display.Colorln("Bright Green", display.ColorBrightGreen)
	display.Colorln("Bright Yellow", display.ColorBrightYellow)
	display.Colorln("Bright Blue", display.ColorBrightBlue)
	display.Colorln("Bright Magenta", display.ColorBrightMagenta)
	display.Colorln("Bright Cyan", display.ColorBrightCyan)
	display.NewLine()

	// 17. Highlighted Text
	display.Title("17. Highlighted Text")
	display.Highlight("White on Red", display.ColorWhite, display.BgRed)
	display.Highlight("Black on Yellow", display.ColorBlack, display.BgYellow)
	display.Highlight("White on Blue", display.ColorWhite, display.BgBlue)
	display.NewLine()

	// 18. Real-World Example
	display.Title("18. Real-World Example: Deployment Status")
	display.Header("Deployment Pipeline")
	display.NewLine()

	display.Step(1, 4, "Building application")
	display.Indent(1, "Compiling source code...")
	display.Indent(1, "Running tests...")
	display.Done("Build successful")
	display.NewLine()

	display.Step(2, 4, "Creating Docker image")
	display.Indent(1, "Pulling base image...")
	display.Indent(1, "Building layers...")
	display.Done("Image created: myapp:v1.2.3")
	display.NewLine()

	display.Step(3, 4, "Pushing to registry")
	display.Progress("Uploading image")
	display.Done("Image pushed successfully")
	display.NewLine()

	display.Step(4, 4, "Deploying to production")
	display.Indent(1, "Updating containers...")
	display.Warning("Old containers will be removed")
	display.Done("Deployment complete")
	display.NewLine()

	display.Separator()
	display.Subheader("Deployment Summary:")
	display.Table("Application", "myapp")
	display.Table("Version", "v1.2.3")
	display.Table("Environment", "Production")
	display.Table("Status", "Running")
	display.NewLine()

	display.Success("All operations completed successfully!")
	display.NewLine()

	// Final banner
	display.Banner("Demo Complete")
}

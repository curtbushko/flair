// Example: charm-lipgloss demonstrates using flair with lipgloss for styled text output.
//
// This example shows how to load a theme and use pre-configured lipgloss styles
// for common UI patterns like surfaces, text, status messages, and components.
//
// Run with: go run ./examples/charm-lipgloss
package main

import (
	"fmt"
	"strings"

	flairlip "github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func main() {
	// Load the currently selected theme, or fall back to tokyo-night-dark.
	theme := flair.MustLoad()
	styles := flairlip.NewStyles(theme)

	fmt.Println()
	fmt.Println(styles.Raised.Padding(0, 2).Render("Flair + Lipgloss Example"))
	fmt.Println()

	// Text styles section
	fmt.Println(styles.Text.Bold(true).Render("Text Styles"))
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println(styles.Text.Render("Primary text - main content"))
	fmt.Println(styles.Secondary.Render("Secondary text - supporting content"))
	fmt.Println(styles.Muted.Render("Muted text - hints and metadata"))
	fmt.Println()

	// Status messages section
	fmt.Println(styles.Text.Bold(true).Render("Status Messages"))
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println(styles.Error.Render("  Error: Connection failed"))
	fmt.Println(styles.Warning.Render("  Warning: Disk space low"))
	fmt.Println(styles.Success.Render("  Success: File saved"))
	fmt.Println(styles.Info.Render("  Info: 3 items selected"))
	fmt.Println()

	// Surface styles section
	fmt.Println(styles.Text.Bold(true).Render("Surface Styles"))
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println(styles.Background.Padding(0, 2).Render("Background surface"))
	fmt.Println(styles.Raised.Padding(0, 2).Render("Raised surface (cards, headers)"))
	fmt.Println(styles.Sunken.Padding(0, 2).Render("Sunken surface (inputs, wells)"))
	fmt.Println()

	// Component styles section
	fmt.Println(styles.Text.Bold(true).Render("Component Styles"))
	fmt.Println(strings.Repeat("-", 40))
	fmt.Print("Buttons: ")
	fmt.Print(styles.Button.Render(" Cancel "))
	fmt.Print("  ")
	fmt.Println(styles.ButtonFocused.Render(" Submit "))

	fmt.Print("List:    ")
	fmt.Print(styles.ListItem.Render("Item 1"))
	fmt.Print("  ")
	fmt.Print(styles.ListSelected.Render("> Item 2"))
	fmt.Print("  ")
	fmt.Println(styles.ListItem.Render("Item 3"))

	fmt.Println()
	fmt.Println(styles.Muted.Render("Theme: " + theme.Name()))
	fmt.Println()
}

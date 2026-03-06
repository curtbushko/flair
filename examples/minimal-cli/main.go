// Example: minimal-cli demonstrates basic CLI theming with flair.
//
// This example shows how to:
// - Load a theme using flair.Default() (respects user selection or falls back)
// - Create lipgloss styles from the theme
// - Print styled output to the terminal
//
// Run with: go run ./examples/minimal-cli
package main

import (
	"fmt"
	"log"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func main() {
	// Load the currently selected theme, or fall back to tokyo-night-dark.
	// This is the recommended way to load a theme for most applications.
	theme, err := flair.Default()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using theme: %s (%s variant)\n\n", theme.Name(), theme.Variant())

	// Create lipgloss styles from the theme.
	// NewStyles builds pre-configured styles for common UI elements.
	styles := lipgloss.NewStyles(theme)

	// Render some styled output demonstrating different style categories.

	// Text styles
	fmt.Println("Text Styles:")
	fmt.Println(styles.Text.Render("  Primary text - main body content"))
	fmt.Println(styles.Secondary.Render("  Secondary text - descriptions and hints"))
	fmt.Println(styles.Muted.Render("  Muted text - disabled or placeholder content"))
	fmt.Println()

	// Status styles
	fmt.Println("Status Messages:")
	fmt.Println(styles.Error.Render("  Error: Something went wrong"))
	fmt.Println(styles.Warning.Render("  Warning: Proceed with caution"))
	fmt.Println(styles.Success.Render("  Success: Operation completed"))
	fmt.Println(styles.Info.Render("  Info: Here's some information"))
	fmt.Println()

	// Access theme colors directly for custom styling
	fmt.Println("Theme Colors:")
	surface := theme.Surface()
	text := theme.Text()
	fmt.Printf("  Background: %s\n", surface.Background.Hex())
	fmt.Printf("  Foreground: %s\n", text.Primary.Hex())
}

// Example: zero-setup demonstrates CLI theming with no configuration required.
//
// This example shows how to:
// - Use flair.LoadBuiltin() to load themes directly from embedded palettes
// - Work without any ~/.config/flair directory
// - List available built-in themes
//
// This is ideal for distributing CLIs that work out of the box.
//
// Run with: go run ./examples/zero-setup
package main

import (
	"fmt"
	"log"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func main() {
	// List all available built-in themes.
	// These are embedded in the binary and always available.
	builtins := flair.ListBuiltins()
	fmt.Println("Available built-in themes:")
	for _, name := range builtins {
		fmt.Printf("  - %s\n", name)
	}
	fmt.Println()

	// Load a specific built-in theme directly.
	// No filesystem access required - the palette is embedded.
	theme, err := flair.LoadBuiltin("catppuccin-mocha")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded theme: %s\n", theme.Name())
	fmt.Printf("Variant: %s\n\n", theme.Variant())

	// Create styles from the theme.
	styles := lipgloss.NewStyles(theme)

	// Render styled output.
	fmt.Println(styles.Text.Render("This is primary text"))
	fmt.Println(styles.Success.Render("Zero setup required!"))
	fmt.Println()

	// Access syntax highlighting colors.
	syntax := theme.Syntax()
	fmt.Println("Syntax colors:")
	fmt.Printf("  Keyword:  %s\n", syntax.Keyword.Hex())
	fmt.Printf("  String:   %s\n", syntax.String.Hex())
	fmt.Printf("  Function: %s\n", syntax.Function.Hex())
	fmt.Printf("  Comment:  %s\n", syntax.Comment.Hex())
}

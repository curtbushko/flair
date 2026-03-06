// Example: embed-viewer demonstrates embedding the flair style viewer in a CLI.
//
// This example shows how to:
// - Use viewer.Run() with custom options
// - Provide an OnSelect callback to respond to theme selection
// - Configure the viewer with a theme loader
//
// Run with: go run ./examples/embed-viewer
package main

import (
	"fmt"
	"log"

	"github.com/curtbushko/flair/pkg/flair"
	"github.com/curtbushko/flair/pkg/flair/viewer"
)

func main() {
	fmt.Println("Flair Theme Viewer")
	fmt.Println("==================")
	fmt.Println("Use arrow keys to navigate, Tab to switch pages, Enter to select, q to quit.")
	fmt.Println()

	// Get the list of built-in themes.
	themes := flair.ListBuiltins()
	if len(themes) == 0 {
		log.Fatal("No built-in themes available")
	}

	// Run the viewer with custom options.
	// The viewer is a bubbletea-based TUI that displays theme previews.
	err := viewer.Run(viewer.Options{
		// List of themes to display in the viewer.
		Themes: themes,

		// Pre-select a specific theme on startup.
		InitialTheme: "tokyo-night-dark",

		// Callback when user confirms selection with Enter.
		OnSelect: func(name string) {
			fmt.Printf("\nSelected theme: %s\n", name)
		},

		// ThemeLoader provides palette and token data for preview.
		// Use NewBuiltinThemeLoader() for embedded themes.
		ThemeLoader: viewer.NewBuiltinThemeLoader(),
	})

	if err != nil {
		log.Fatal(err)
	}
}

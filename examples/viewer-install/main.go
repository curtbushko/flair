// Example: viewer-install demonstrates the viewer with install-on-select behavior.
//
// This example shows how to:
// - Use viewer.RunBuiltins() for zero-config viewer startup
// - Provide OnSelect and OnInstall callbacks
// - Install themes to ~/.config/flair when selected
//
// Run with: go run ./examples/viewer-install
package main

import (
	"fmt"
	"log"

	"github.com/curtbushko/flair/pkg/flair"
	"github.com/curtbushko/flair/pkg/flair/viewer"
)

func main() {
	fmt.Println("Flair Theme Viewer with Install")
	fmt.Println("================================")
	fmt.Println("Navigate with arrows, Enter to select and install, q to quit.")
	fmt.Println()

	// Create a store for installing themes.
	store := flair.NewStore()

	// RunBuiltins is a zero-config function that shows all built-in themes.
	// It does not require any themes to be installed beforehand.
	err := viewer.RunBuiltins(viewer.RunBuiltinsOptions{
		// Use alternate screen buffer (default is true).
		WithAltScreen: true,

		// Pre-select the currently selected theme if any.
		InitialTheme: func() string {
			name, _ := flair.SelectedTheme()
			return name
		}(),

		// Called when user confirms selection with Enter.
		OnSelect: func(name string) {
			fmt.Printf("\nSelected: %s\n", name)
		},

		// Called to install the selected theme.
		// This writes theme files to ~/.config/flair/<themename>/
		OnInstall: func(name string) error {
			fmt.Printf("Installing %s...\n", name)

			// Install the theme files.
			if err := store.Install(name); err != nil {
				return err
			}

			// Select the theme (creates symlinks).
			if err := store.Select(name); err != nil {
				return err
			}

			fmt.Printf("Theme %s installed and selected!\n", name)
			fmt.Printf("Config directory: %s\n", store.ConfigDir())
			return nil
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

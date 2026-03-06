// Example: charm-huh demonstrates using flair with huh for themed forms.
//
// This example shows a simple form with text input, select, and confirm fields
// that all use the flair theme colors.
//
// Run with: go run ./examples/charm-huh
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"

	flairhuh "github.com/curtbushko/flair/pkg/charm/huh"
	"github.com/curtbushko/flair/pkg/flair"
)

func main() {
	// Load theme and create huh theme.
	theme := flair.MustLoad()
	huhTheme := flairhuh.NewTheme(theme)

	var (
		name     string
		language string
		confirm  bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's your name?").
				Placeholder("Enter your name").
				Value(&name),

			huh.NewSelect[string]().
				Title("Favorite programming language?").
				Options(
					huh.NewOption("Go", "go"),
					huh.NewOption("Rust", "rust"),
					huh.NewOption("Python", "python"),
					huh.NewOption("TypeScript", "typescript"),
				).
				Value(&language),

			huh.NewConfirm().
				Title("Ready to continue?").
				Affirmative("Yes!").
				Negative("No").
				Value(&confirm),
		),
	).WithTheme(huhTheme)

	err := form.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println()
	if confirm {
		fmt.Printf("Hello, %s! You chose %s.\n", name, language)
	} else {
		fmt.Println("Maybe next time!")
	}
}

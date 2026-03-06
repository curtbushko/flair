// Example: charm-list demonstrates using flair with bubbles/list for themed lists.
//
// This example shows how to create themed list styles and delegates.
// Note: This is a demonstration of the styling API. For a full interactive
// list, see the bubbles documentation.
//
// Run with: go run ./examples/charm-list
package main

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/curtbushko/flair/pkg/flair"
)

func main() {
	// Load theme.
	theme := flair.MustLoad()
	colors := theme.Colors()

	// Build list styles using theme colors.
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors["text.primary"].Hex())).
		MarginBottom(1)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors["text.primary"].Hex())).
		PaddingLeft(2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors["accent.primary"].Hex())).
		Background(lipgloss.Color(colors["surface.background.selection"].Hex())).
		Bold(true).
		PaddingLeft(2)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors["text.muted"].Hex())).
		PaddingLeft(4)

	// Sample list items.
	items := []struct {
		title, desc string
		selected    bool
	}{
		{"Tokyo Night", "A dark theme inspired by Tokyo's night lights", false},
		{"Gruvbox Dark", "Retro groove color scheme with warm tones", true},
		{"Catppuccin Mocha", "Soothing pastel theme with rich colors", false},
		{"Dracula", "Dark theme with vibrant accent colors", false},
		{"Nord", "Arctic, north-bluish color palette", false},
	}

	// Render the list.
	fmt.Println()
	fmt.Println(titleStyle.Render("Select a Theme"))
	fmt.Println(strings.Repeat("-", 50))

	for _, item := range items {
		var line string
		if item.selected {
			line = selectedStyle.Render("> " + item.title)
		} else {
			line = itemStyle.Render("  " + item.title)
		}
		fmt.Println(line)
		fmt.Println(descStyle.Render(item.desc))
	}

	fmt.Println()
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors["text.muted"].Hex()))
	fmt.Println(helpStyle.Render("Theme: " + theme.Name()))
	fmt.Println()
}

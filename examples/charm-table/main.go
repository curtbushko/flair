// Example: charm-table demonstrates using flair with bubbles/table for themed tables.
//
// This example shows how to create themed table styles.
// Note: This is a demonstration of the styling API. For a full interactive
// table, see the bubbles documentation.
//
// Run with: go run ./examples/charm-table
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

	// Build table styles using theme colors.
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors["text.secondary"].Hex())).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(colors["border.default"].Hex()))

	cellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors["text.primary"].Hex())).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors["text.inverse"].Hex())).
		Background(lipgloss.Color(colors["accent.primary"].Hex())).
		Bold(true).
		Padding(0, 1)

	// Table data.
	headers := []string{"Theme", "Variant", "Colors", "Status"}
	widths := []int{20, 10, 8, 12}

	rows := [][]string{
		{"Tokyo Night", "Dark", "24", "Installed"},
		{"Gruvbox", "Dark", "24", "Installed"},
		{"Catppuccin", "Mocha", "24", "Selected"},
		{"Dracula", "Dark", "24", "Available"},
		{"Nord", "Dark", "16", "Available"},
	}

	// Render the table.
	fmt.Println()
	titleStyle := lipgloss.NewStyle().Bold(true).MarginBottom(1)
	fmt.Println(titleStyle.Render("Theme Table"))
	fmt.Println()

	// Render header.
	var headerRow strings.Builder
	for i, h := range headers {
		cell := headerStyle.Width(widths[i]).Render(h)
		headerRow.WriteString(cell)
	}
	fmt.Println(headerRow.String())

	// Render rows.
	for rowIdx, row := range rows {
		var rowStr strings.Builder
		isSelected := rowIdx == 2 // Catppuccin is selected

		for i, cell := range row {
			var style lipgloss.Style
			if isSelected {
				style = selectedStyle.Width(widths[i])
			} else {
				style = cellStyle.Width(widths[i])
			}
			rowStr.WriteString(style.Render(cell))
		}
		fmt.Println(rowStr.String())
	}

	fmt.Println()
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors["text.muted"].Hex()))
	fmt.Println(helpStyle.Render("Theme: " + theme.Name()))
	fmt.Println()
}

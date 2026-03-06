// Example: charm-viewport demonstrates using flair for styled content areas.
//
// This example shows how to create themed viewport/content styles.
// Note: This is a demonstration of the styling API. For a full scrollable
// viewport, see the bubbles documentation.
//
// Run with: go run ./examples/charm-viewport
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

	// Build viewport styles using theme colors.
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors["text.primary"].Hex()))

	contentStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(colors["surface.background.sunken"].Hex())).
		Foreground(lipgloss.Color(colors["text.primary"].Hex())).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors["border.default"].Hex())).
		Padding(1, 2).
		Width(60)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors["text.secondary"].Hex())).
		MarginBottom(1)

	codeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors["syntax.string"].Hex()))

	// Build content.
	var content strings.Builder
	content.WriteString(headerStyle.Render("Theme Information"))
	content.WriteString("\n")
	content.WriteString(fmt.Sprintf("Name:    %s\n", theme.Name()))
	content.WriteString(fmt.Sprintf("Variant: %s\n\n", theme.Variant()))

	content.WriteString(headerStyle.Render("Sample Colors"))
	content.WriteString("\n")
	sampleColors := []string{
		"surface.background",
		"text.primary",
		"status.error",
		"status.success",
		"syntax.keyword",
	}
	for _, name := range sampleColors {
		if c, ok := colors[name]; ok {
			content.WriteString(fmt.Sprintf("%-22s %s\n", name+":", codeStyle.Render(c.Hex())))
		}
	}

	content.WriteString("\n")
	content.WriteString(headerStyle.Render("Description"))
	content.WriteString("\n")
	content.WriteString("This viewport demonstrates themed content areas\n")
	content.WriteString("using flair colors for borders, backgrounds,\n")
	content.WriteString("and text styling.")

	// Render.
	fmt.Println()
	fmt.Println(titleStyle.Render("Flair + Viewport Example"))
	fmt.Println()
	fmt.Println(contentStyle.Render(content.String()))
	fmt.Println()

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors["text.muted"].Hex()))
	fmt.Println(helpStyle.Render("Theme: " + theme.Name()))
	fmt.Println()
}

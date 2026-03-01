package lipgloss

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildBorderDefault creates a lipgloss style with the default border color.
// Uses a rounded border with border.default color.
func BuildBorderDefault(theme *flair.Theme) lipgloss.Style {
	borderDefault := getColor(theme, "border.default", flair.Color{R: 86, G: 95, B: 137})
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderDefault.Hex()))
}

// BuildBorderFocus creates a lipgloss style with the focus border color.
// Uses a rounded border with border.focus color for focused elements.
func BuildBorderFocus(theme *flair.Theme) lipgloss.Style {
	borderFocus := getColor(theme, "border.focus", flair.Color{R: 122, G: 162, B: 247})
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderFocus.Hex()))
}

// BuildBorderMuted creates a lipgloss style with the muted border color.
// Uses a rounded border with border.muted color for less prominent elements.
func BuildBorderMuted(theme *flair.Theme) lipgloss.Style {
	borderMuted := getColor(theme, "border.muted", flair.Color{R: 59, G: 66, B: 97})
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderMuted.Hex()))
}

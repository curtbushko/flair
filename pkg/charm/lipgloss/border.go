package lipgloss

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildBorderDefault creates a lipgloss style with the default border color.
//
// This style uses a rounded border with the border.default color from the theme.
// It is suitable for general-purpose container borders.
func BuildBorderDefault(theme *flair.Theme) lipgloss.Style {
	borderDefault := getColor(theme, "border.default", flair.Color{R: 86, G: 95, B: 137})
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderDefault.Hex()))
}

// BuildBorderFocus creates a lipgloss style with the focus border color.
//
// This style uses a rounded border with the border.focus color from the theme.
// It is suitable for indicating focused or active elements.
func BuildBorderFocus(theme *flair.Theme) lipgloss.Style {
	borderFocus := getColor(theme, "border.focus", flair.Color{R: 122, G: 162, B: 247})
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderFocus.Hex()))
}

// BuildBorderMuted creates a lipgloss style with the muted border color.
//
// This style uses a rounded border with the border.muted color from the theme.
// It is suitable for subtle, less prominent element borders.
func BuildBorderMuted(theme *flair.Theme) lipgloss.Style {
	borderMuted := getColor(theme, "border.muted", flair.Color{R: 59, G: 66, B: 97})
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderMuted.Hex()))
}

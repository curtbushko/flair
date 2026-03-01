package lipgloss

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildStateHover creates a lipgloss style for hovered elements.
// Uses highlight surface background for subtle visual feedback.
func BuildStateHover(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Highlight.Hex()))
}

// BuildStateActive creates a lipgloss style for active/pressed elements.
// Uses accent primary foreground to indicate active state.
func BuildStateActive(theme *flair.Theme) lipgloss.Style {
	accentPrimary := getColor(theme, "accent.primary", flair.Color{R: 122, G: 162, B: 247})
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentPrimary.Hex())).
		Background(lipgloss.Color(text.Inverse.Hex()))
}

// BuildStateDisabled creates a lipgloss style for disabled elements.
// Uses muted text foreground to indicate unavailability.
func BuildStateDisabled(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Muted.Hex()))
}

// BuildStateSelected creates a lipgloss style for selected elements.
// Uses selection surface background with primary text.
func BuildStateSelected(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Selection.Hex()))
}

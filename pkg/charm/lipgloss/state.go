package lipgloss

import (
	"charm.land/lipgloss/v2"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildStateHover creates a lipgloss style for hovered elements.
//
// This style uses the highlight surface background with primary text foreground,
// providing subtle visual feedback for hover states.
func BuildStateHover(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Highlight.Hex()))
}

// BuildStateActive creates a lipgloss style for active/pressed elements.
//
// This style uses accent primary foreground with inverse background,
// clearly indicating that an element is currently active or pressed.
func BuildStateActive(theme *flair.Theme) lipgloss.Style {
	accentPrimary := getColor(theme, "accent.primary", flair.Color{R: 122, G: 162, B: 247})
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentPrimary.Hex())).
		Background(lipgloss.Color(text.Inverse.Hex()))
}

// BuildStateDisabled creates a lipgloss style for disabled elements.
//
// This style uses muted text foreground to visually indicate that
// an element is unavailable for interaction.
func BuildStateDisabled(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Muted.Hex()))
}

// BuildStateSelected creates a lipgloss style for selected elements.
//
// This style uses the selection surface background with primary text foreground,
// suitable for highlighting selected items in lists or tables.
func BuildStateSelected(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Selection.Hex()))
}

package lipgloss

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildSurfaceBackground creates a lipgloss style with the surface background color.
func BuildSurfaceBackground(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Background.Hex()))
}

// BuildSurfaceRaised creates a lipgloss style with the raised surface color.
func BuildSurfaceRaised(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Raised.Hex()))
}

// BuildSurfaceSunken creates a lipgloss style with the sunken surface color.
func BuildSurfaceSunken(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Sunken.Hex()))
}

// BuildSurfaceOverlay creates a lipgloss style with the overlay surface color.
func BuildSurfaceOverlay(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Overlay.Hex()))
}

// BuildSurfacePopup creates a lipgloss style with the popup surface color.
func BuildSurfacePopup(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Popup.Hex()))
}

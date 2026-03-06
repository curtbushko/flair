package lipgloss

import (
	"charm.land/lipgloss/v2"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildSurfaceBackground creates a lipgloss style with the primary surface background color.
//
// This style is suitable for the main application background.
func BuildSurfaceBackground(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Background.Hex()))
}

// BuildSurfaceRaised creates a lipgloss style with the raised surface color.
//
// This style is suitable for elevated UI elements like cards and panels.
func BuildSurfaceRaised(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Raised.Hex()))
}

// BuildSurfaceSunken creates a lipgloss style with the sunken surface color.
//
// This style is suitable for inset areas like input field backgrounds.
func BuildSurfaceSunken(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Sunken.Hex()))
}

// BuildSurfaceOverlay creates a lipgloss style with the overlay surface color.
//
// This style is suitable for modal overlays and dialogs.
func BuildSurfaceOverlay(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Overlay.Hex()))
}

// BuildSurfacePopup creates a lipgloss style with the popup surface color.
//
// This style is suitable for popup menus and tooltips.
func BuildSurfacePopup(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Popup.Hex()))
}

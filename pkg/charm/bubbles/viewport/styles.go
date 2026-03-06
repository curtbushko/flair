package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// NewStyle creates a lipgloss.Style for viewport from a flair.Theme.
//
// NewStyle applies theme colors to create a style suitable for viewport background
// and text rendering. If theme is nil, it returns an empty lipgloss.Style.
//
// The style uses surface.background for the background color and text.primary
// for the foreground color.
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	style := viewport.NewStyle(theme)
//	myViewport.Style = style
func NewStyle(theme *flair.Theme) lipgloss.Style {
	if theme == nil {
		return lipgloss.NewStyle()
	}

	// Get colors from theme with fallbacks.
	surfaceBg := getColorHex(theme, "surface.background", "#1a1b26")
	textPrimary := getColorHex(theme, "text.primary", "#c0caf5")

	return lipgloss.NewStyle().
		Background(lipgloss.Color(surfaceBg)).
		Foreground(lipgloss.Color(textPrimary))
}

// NewModel creates a viewport.Model with themed style from a flair.Theme.
//
// NewModel creates a new viewport with the specified dimensions and applies
// the theme's style to it. If theme is nil, it returns a default viewport.Model.
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	vp := viewport.NewModel(theme, 80, 24)
//	vp.SetContent("Hello, world!")
func NewModel(theme *flair.Theme, width, height int) viewport.Model {
	model := viewport.New(width, height)

	if theme != nil {
		model.Style = NewStyle(theme)
	}

	return model
}

// Default returns an empty lipgloss.Style.
//
// This is a convenience function for when no theming is needed.
func Default() lipgloss.Style {
	return lipgloss.NewStyle()
}

// getColorHex retrieves a color from the theme by path, returning the hex string.
// If the color is not found, it returns the fallback value.
func getColorHex(theme *flair.Theme, path, fallback string) string {
	if c, ok := theme.Color(path); ok {
		return c.Hex()
	}
	return fallback
}

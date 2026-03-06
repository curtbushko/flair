package bubbletea

import (
	"charm.land/lipgloss/v2"

	"github.com/curtbushko/flair/pkg/flair"
)

// Styles contains themed lipgloss styles for bubbletea applications.
//
// Styles provides pre-configured [lipgloss.Style] values organized by semantic
// category (Surface, Text, Status, etc.). Each style uses colors from a flair
// theme, ensuring visual consistency across the application.
//
// Create Styles using [NewStyles] with a [flair.Theme], or use [Default] to
// load styles from the currently selected theme.
//
// All styles can be further customized using lipgloss methods:
//
//	// Add padding to a background style
//	panel := styles.Surface.Background.Padding(1, 2)
//
//	// Make an error message bold
//	errMsg := styles.Status.Error.Bold(true).Render("Error!")
type Styles struct {
	// Surface contains background-related styles.
	Surface SurfaceStyles

	// Text contains text/foreground-related styles.
	Text TextStyles

	// Status contains status message styles (error, warning, success, info).
	Status StatusStyles

	// Accent contains accent color styles.
	Accent AccentStyles

	// Border contains border color styles.
	Border BorderStyles
}

// SurfaceStyles contains lipgloss styles for surface/background elements.
type SurfaceStyles struct {
	// Background is the primary background style.
	Background lipgloss.Style

	// Raised is for elevated surfaces like cards and panels.
	Raised lipgloss.Style

	// Sunken is for inset areas like input fields.
	Sunken lipgloss.Style

	// Selection is for selected/highlighted backgrounds.
	Selection lipgloss.Style

	// Overlay is for modal overlays.
	Overlay lipgloss.Style

	// Popup is for popup menus and tooltips.
	Popup lipgloss.Style
}

// TextStyles contains lipgloss styles for text/foreground elements.
type TextStyles struct {
	// Primary is the main text style for body content.
	Primary lipgloss.Style

	// Secondary is for less prominent text.
	Secondary lipgloss.Style

	// Muted is for disabled or placeholder text.
	Muted lipgloss.Style

	// Subtle is for very low-contrast text.
	Subtle lipgloss.Style

	// Inverse is for text on accent backgrounds.
	Inverse lipgloss.Style
}

// StatusStyles contains lipgloss styles for status messages.
type StatusStyles struct {
	// Error is for error messages and indicators.
	Error lipgloss.Style

	// Warning is for warning messages.
	Warning lipgloss.Style

	// Success is for success messages.
	Success lipgloss.Style

	// Info is for informational messages.
	Info lipgloss.Style

	// Hint is for hints and suggestions.
	Hint lipgloss.Style
}

// AccentStyles contains lipgloss styles for accent colors.
type AccentStyles struct {
	// Primary is the main accent color style.
	Primary lipgloss.Style

	// Secondary is a secondary accent color style.
	Secondary lipgloss.Style
}

// BorderStyles contains lipgloss styles for borders.
type BorderStyles struct {
	// Default is the default border color style.
	Default lipgloss.Style

	// Focus is the border color style for focused elements.
	Focus lipgloss.Style
}

// NewStyles creates a new [Styles] from a [flair.Theme].
//
// NewStyles initializes all style categories with colors from the theme.
// If theme is nil, NewStyles returns nil.
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	styles := bubbletea.NewStyles(theme)
//	bg := styles.Surface.Background.Padding(1, 2)
//	panel := bg.Render("Content")
func NewStyles(theme *flair.Theme) *Styles {
	if theme == nil {
		return nil
	}

	surface := theme.Surface()
	text := theme.Text()
	status := theme.Status()

	// Get accent and border colors with fallbacks.
	accentPrimary := getColorHex(theme, "accent.primary", "#7aa2f7")
	accentSecondary := getColorHex(theme, "accent.secondary", "#bb9af7")
	borderDefault := getColorHex(theme, "border.default", "#565f89")
	borderFocus := getColorHex(theme, "border.focus", "#7aa2f7")

	return &Styles{
		Surface: SurfaceStyles{
			Background: lipgloss.NewStyle().Background(lipgloss.Color(surface.Background.Hex())),
			Raised:     lipgloss.NewStyle().Background(lipgloss.Color(surface.Raised.Hex())),
			Sunken:     lipgloss.NewStyle().Background(lipgloss.Color(surface.Sunken.Hex())),
			Selection:  lipgloss.NewStyle().Background(lipgloss.Color(surface.Selection.Hex())),
			Overlay:    lipgloss.NewStyle().Background(lipgloss.Color(surface.Overlay.Hex())),
			Popup:      lipgloss.NewStyle().Background(lipgloss.Color(surface.Popup.Hex())),
		},
		Text: TextStyles{
			Primary:   lipgloss.NewStyle().Foreground(lipgloss.Color(text.Primary.Hex())),
			Secondary: lipgloss.NewStyle().Foreground(lipgloss.Color(text.Secondary.Hex())),
			Muted:     lipgloss.NewStyle().Foreground(lipgloss.Color(text.Muted.Hex())),
			Subtle:    lipgloss.NewStyle().Foreground(lipgloss.Color(text.Subtle.Hex())),
			Inverse:   lipgloss.NewStyle().Foreground(lipgloss.Color(text.Inverse.Hex())),
		},
		Status: StatusStyles{
			Error:   lipgloss.NewStyle().Foreground(lipgloss.Color(status.Error.Hex())),
			Warning: lipgloss.NewStyle().Foreground(lipgloss.Color(status.Warning.Hex())),
			Success: lipgloss.NewStyle().Foreground(lipgloss.Color(status.Success.Hex())),
			Info:    lipgloss.NewStyle().Foreground(lipgloss.Color(status.Info.Hex())),
			Hint:    lipgloss.NewStyle().Foreground(lipgloss.Color(status.Hint.Hex())),
		},
		Accent: AccentStyles{
			Primary:   lipgloss.NewStyle().Foreground(lipgloss.Color(accentPrimary)),
			Secondary: lipgloss.NewStyle().Foreground(lipgloss.Color(accentSecondary)),
		},
		Border: BorderStyles{
			Default: lipgloss.NewStyle().BorderForeground(lipgloss.Color(borderDefault)),
			Focus:   lipgloss.NewStyle().BorderForeground(lipgloss.Color(borderFocus)),
		},
	}
}

// Default loads themed styles from the currently selected flair theme.
//
// Default uses [flair.Default] to load the theme, which falls back to the
// built-in default theme (tokyo-night-dark) if no theme is selected.
//
// Example:
//
//	styles, err := bubbletea.Default()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	bg := styles.Surface.Background
func Default() (*Styles, error) {
	theme, err := flair.Default()
	if err != nil {
		return nil, err
	}
	return NewStyles(theme), nil
}

// DefaultFrom loads themed styles from the specified config directory.
//
// DefaultFrom is useful for testing or when using a non-standard config
// location. It loads the currently selected theme from the config directory.
//
// Returns an error if no theme is selected or the theme cannot be loaded.
func DefaultFrom(configDir string) (*Styles, error) {
	theme, err := flair.LoadFrom(configDir)
	if err != nil {
		return nil, err
	}
	return NewStyles(theme), nil
}

// getColorHex retrieves a color from the theme by path, returning the hex string.
// If the color is not found, it returns the fallback value.
func getColorHex(theme *flair.Theme, path, fallback string) string {
	if c, ok := theme.Color(path); ok {
		return c.Hex()
	}
	return fallback
}

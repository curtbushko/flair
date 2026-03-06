// Package lipgloss provides themed lipgloss styles built from flair themes.
// It creates pre-configured lipgloss.Style instances for common UI elements
// like surfaces, text, status indicators, borders, and components.
//
// This package is fully independent from flair's internal packages.
// It only depends on pkg/flair and github.com/charmbracelet/lipgloss.
package lipgloss

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// Styles contains pre-configured lipgloss styles for a flair theme.
//
// Styles provides ready-to-use [lipgloss.Style] instances for common UI elements
// organized by category. Create a Styles instance using [NewStyles] or [Default].
//
// All styles are derived from the semantic color tokens in the source theme,
// ensuring visual consistency across your TUI application.
//
// Example:
//
//	styles := lipgloss.NewStyles(theme)
//	fmt.Println(styles.Error.Render("Error: operation failed"))
//	fmt.Println(styles.Success.Render("Success: file saved"))
type Styles struct {
	// Surface styles for backgrounds
	Background lipgloss.Style // Primary background
	Raised     lipgloss.Style // Elevated surfaces like cards
	Sunken     lipgloss.Style // Inset areas like input fields
	Overlay    lipgloss.Style // Modal overlays
	Popup      lipgloss.Style // Popup menus and tooltips

	// Text styles for foregrounds
	Text      lipgloss.Style // Primary text
	Secondary lipgloss.Style // Less prominent text
	Muted     lipgloss.Style // Disabled or placeholder text
	Inverse   lipgloss.Style // Text on accent backgrounds

	// Status styles for messages
	Error   lipgloss.Style // Error messages
	Warning lipgloss.Style // Warning messages
	Success lipgloss.Style // Success messages
	Info    lipgloss.Style // Informational messages

	// Border styles for containers
	Border      lipgloss.Style // Default border
	BorderFocus lipgloss.Style // Focused element border
	BorderMuted lipgloss.Style // Subtle border

	// Component styles for interactive elements
	Button        lipgloss.Style // Unfocused button
	ButtonFocused lipgloss.Style // Focused button
	Input         lipgloss.Style // Unfocused input field
	InputFocused  lipgloss.Style // Focused input field
	ListItem      lipgloss.Style // List item
	ListSelected  lipgloss.Style // Selected list item
	Table         lipgloss.Style // Table cell
	TableHeader   lipgloss.Style // Table header
	Dialog        lipgloss.Style // Dialog box

	// State styles for element states
	Hover    lipgloss.Style // Hovered element
	Active   lipgloss.Style // Active/pressed element
	Disabled lipgloss.Style // Disabled element
	Selected lipgloss.Style // Selected element
}

// buildForegroundStyle creates a lipgloss style with the given hex color as foreground.
func buildForegroundStyle(hexColor string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(hexColor))
}

// BuildStatusError creates a lipgloss style with the error status foreground color.
//
// This style is suitable for rendering error messages in TUI applications.
func BuildStatusError(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Error.Hex())
}

// BuildStatusWarning creates a lipgloss style with the warning status foreground color.
//
// This style is suitable for rendering warning messages in TUI applications.
func BuildStatusWarning(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Warning.Hex())
}

// BuildStatusSuccess creates a lipgloss style with the success status foreground color.
//
// This style is suitable for rendering success messages in TUI applications.
func BuildStatusSuccess(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Success.Hex())
}

// BuildStatusInfo creates a lipgloss style with the info status foreground color.
//
// This style is suitable for rendering informational messages in TUI applications.
func BuildStatusInfo(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Info.Hex())
}

// BuildTextPrimary creates a lipgloss style with the primary text foreground color.
//
// This style is suitable for main body text content.
func BuildTextPrimary(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Primary.Hex())
}

// BuildTextSecondary creates a lipgloss style with the secondary text foreground color.
//
// This style is suitable for less prominent text like descriptions or hints.
func BuildTextSecondary(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Secondary.Hex())
}

// BuildTextMuted creates a lipgloss style with the muted text foreground color.
//
// This style is suitable for disabled text or placeholder content.
func BuildTextMuted(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Muted.Hex())
}

// BuildTextInverse creates a lipgloss style with the inverse text foreground color.
//
// This style is suitable for text on accent-colored backgrounds.
func BuildTextInverse(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Inverse.Hex())
}

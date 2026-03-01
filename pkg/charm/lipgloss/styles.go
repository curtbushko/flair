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
// Each field provides a ready-to-use style for common UI elements.
type Styles struct {
	// Surface styles
	Background lipgloss.Style
	Raised     lipgloss.Style
	Sunken     lipgloss.Style
	Overlay    lipgloss.Style
	Popup      lipgloss.Style

	// Text styles
	Text      lipgloss.Style
	Secondary lipgloss.Style
	Muted     lipgloss.Style
	Inverse   lipgloss.Style

	// Status styles
	Error   lipgloss.Style
	Warning lipgloss.Style
	Success lipgloss.Style
	Info    lipgloss.Style

	// Border styles
	Border      lipgloss.Style
	BorderFocus lipgloss.Style
	BorderMuted lipgloss.Style

	// Component styles
	Button        lipgloss.Style
	ButtonFocused lipgloss.Style
	Input         lipgloss.Style
	InputFocused  lipgloss.Style
	ListItem      lipgloss.Style
	ListSelected  lipgloss.Style
	Table         lipgloss.Style
	TableHeader   lipgloss.Style
	Dialog        lipgloss.Style

	// State styles
	Hover    lipgloss.Style
	Active   lipgloss.Style
	Disabled lipgloss.Style
	Selected lipgloss.Style
}

// buildForegroundStyle creates a lipgloss style with the given hex color as foreground.
func buildForegroundStyle(hexColor string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(hexColor))
}

// BuildStatusError creates a lipgloss style with the error status foreground color.
func BuildStatusError(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Error.Hex())
}

// BuildStatusWarning creates a lipgloss style with the warning status foreground color.
func BuildStatusWarning(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Warning.Hex())
}

// BuildStatusSuccess creates a lipgloss style with the success status foreground color.
func BuildStatusSuccess(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Success.Hex())
}

// BuildStatusInfo creates a lipgloss style with the info status foreground color.
func BuildStatusInfo(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Status().Info.Hex())
}

// BuildTextPrimary creates a lipgloss style with the primary text foreground color.
func BuildTextPrimary(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Primary.Hex())
}

// BuildTextSecondary creates a lipgloss style with the secondary text foreground color.
func BuildTextSecondary(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Secondary.Hex())
}

// BuildTextMuted creates a lipgloss style with the muted text foreground color.
func BuildTextMuted(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Muted.Hex())
}

// BuildTextInverse creates a lipgloss style with the inverse text foreground color.
func BuildTextInverse(theme *flair.Theme) lipgloss.Style {
	return buildForegroundStyle(theme.Text().Inverse.Hex())
}

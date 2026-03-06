package lipgloss

import (
	"charm.land/lipgloss/v2"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildButton creates a lipgloss style for unfocused buttons.
//
// This style uses a raised surface background with primary text foreground
// and horizontal padding. Suitable for standard button elements.
func BuildButton(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Raised.Hex())).
		Padding(0, 2)
}

// BuildButtonFocused creates a lipgloss style for focused buttons.
//
// This style uses an accent primary background with inverse text, bold weight,
// and horizontal padding. Suitable for highlighting the focused button.
func BuildButtonFocused(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	accentPrimary := getColor(theme, "accent.primary", flair.Color{R: 122, G: 162, B: 247})
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Inverse.Hex())).
		Background(lipgloss.Color(accentPrimary.Hex())).
		Padding(0, 2).
		Bold(true)
}

// BuildInput creates a lipgloss style for unfocused input fields.
//
// This style uses a sunken surface background with primary text foreground
// and a rounded default border. Suitable for text input fields.
func BuildInput(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	surface := theme.Surface()
	borderDefault := getColor(theme, "border.default", flair.Color{R: 86, G: 95, B: 137})
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Sunken.Hex())).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderDefault.Hex()))
}

// BuildInputFocused creates a lipgloss style for focused input fields.
//
// This style uses a sunken surface background with primary text foreground
// and a rounded focus border. Suitable for indicating the active input field.
func BuildInputFocused(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	surface := theme.Surface()
	borderFocus := getColor(theme, "border.focus", flair.Color{R: 122, G: 162, B: 247})
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Sunken.Hex())).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderFocus.Hex()))
}

// BuildListItem creates a lipgloss style for list items.
//
// This style uses primary text foreground. Suitable for normal list items.
func BuildListItem(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex()))
}

// BuildListSelected creates a lipgloss style for selected list items.
//
// This style uses accent primary foreground with bold weight.
// Suitable for highlighting the currently selected list item.
func BuildListSelected(theme *flair.Theme) lipgloss.Style {
	accentPrimary := getColor(theme, "accent.primary", flair.Color{R: 122, G: 162, B: 247})
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentPrimary.Hex())).
		Bold(true)
}

// BuildTable creates a lipgloss style for table cells.
//
// This style uses primary text foreground. Suitable for table body cells.
func BuildTable(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex()))
}

// BuildTableHeader creates a lipgloss style for table headers.
//
// This style uses secondary text foreground with bold weight.
// Suitable for table column headers.
func BuildTableHeader(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Secondary.Hex())).
		Bold(true)
}

// BuildDialog creates a lipgloss style for dialog boxes.
//
// This style uses overlay surface background with a rounded default border
// and padding. Suitable for modal dialogs and confirmation boxes.
func BuildDialog(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	borderDefault := getColor(theme, "border.default", flair.Color{R: 86, G: 95, B: 137})
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Overlay.Hex())).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderDefault.Hex())).
		Padding(1, 2)
}

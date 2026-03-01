package lipgloss

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// BuildButton creates a lipgloss style for unfocused buttons.
// Uses raised surface background with primary text and horizontal padding.
func BuildButton(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	surface := theme.Surface()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(surface.Raised.Hex())).
		Padding(0, 2)
}

// BuildButtonFocused creates a lipgloss style for focused buttons.
// Uses accent primary background with inverse text, bold, and horizontal padding.
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
// Uses sunken surface background with primary text and default border.
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
// Uses sunken surface background with primary text and focus border.
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
// Uses primary text foreground.
func BuildListItem(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex()))
}

// BuildListSelected creates a lipgloss style for selected list items.
// Uses accent primary foreground with bold.
func BuildListSelected(theme *flair.Theme) lipgloss.Style {
	accentPrimary := getColor(theme, "accent.primary", flair.Color{R: 122, G: 162, B: 247})
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentPrimary.Hex())).
		Bold(true)
}

// BuildTable creates a lipgloss style for table cells.
// Uses primary text foreground.
func BuildTable(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Primary.Hex()))
}

// BuildTableHeader creates a lipgloss style for table headers.
// Uses secondary text foreground with bold.
func BuildTableHeader(theme *flair.Theme) lipgloss.Style {
	text := theme.Text()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(text.Secondary.Hex())).
		Bold(true)
}

// BuildDialog creates a lipgloss style for dialog boxes.
// Uses overlay surface background with default border and padding.
func BuildDialog(theme *flair.Theme) lipgloss.Style {
	surface := theme.Surface()
	borderDefault := getColor(theme, "border.default", flair.Color{R: 86, G: 95, B: 137})
	return lipgloss.NewStyle().
		Background(lipgloss.Color(surface.Overlay.Hex())).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderDefault.Hex())).
		Padding(1, 2)
}

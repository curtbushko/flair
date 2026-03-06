package table

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// NewStyles creates table.Styles from a flair.Theme.
//
// NewStyles applies theme colors to the table component's header, cell,
// and selected row styles. If theme is nil, it returns table.DefaultStyles().
//
// The header is styled with secondary text color and bold formatting.
// Selected rows use the selection background with accent foreground.
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	styles := table.NewStyles(theme)
//	myTable := bubbletable.New(bubbletable.WithStyles(styles))
func NewStyles(theme *flair.Theme) table.Styles {
	if theme == nil {
		return table.DefaultStyles()
	}

	// Get colors from theme with fallbacks.
	textPrimary := getColorHex(theme, "text.primary", "#c0caf5")
	textSecondary := getColorHex(theme, "text.secondary", "#a9b1d6")
	accentPrimary := getColorHex(theme, "accent.primary", "#7aa2f7")
	selectionBg := getColorHex(theme, "surface.background.selection", "#33467c")
	surfaceRaised := getColorHex(theme, "surface.background.raised", "#24283b")
	borderColor := getColorHex(theme, "border.default", "#565f89")

	s := table.DefaultStyles()

	// Header styling: secondary text, bold, with border.
	s.Header = s.Header.
		Foreground(lipgloss.Color(textSecondary)).
		Background(lipgloss.Color(surfaceRaised)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		BorderBottom(true).
		Bold(true)

	// Cell styling: primary text color.
	s.Cell = s.Cell.
		Foreground(lipgloss.Color(textPrimary))

	// Selected row styling: accent foreground with selection background.
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(accentPrimary)).
		Background(lipgloss.Color(selectionBg)).
		Bold(true)

	return s
}

// Default returns the default table.Styles.
//
// This is a convenience function equivalent to table.DefaultStyles().
func Default() table.Styles {
	return table.DefaultStyles()
}

// getColorHex retrieves a color from the theme by path, returning the hex string.
// If the color is not found, it returns the fallback value.
func getColorHex(theme *flair.Theme, path, fallback string) string {
	if c, ok := theme.Color(path); ok {
		return c.Hex()
	}
	return fallback
}

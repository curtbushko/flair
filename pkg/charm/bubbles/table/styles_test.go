package table_test

import (
	"testing"

	bubbletable "github.com/charmbracelet/bubbles/table"

	"github.com/curtbushko/flair/pkg/charm/bubbles/table"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestNewStyles_CreatesTableStyles(t *testing.T) {
	// Arrange: Create a theme with text, surface, border colors.
	colors := map[string]flair.Color{
		"surface.background":           {R: 26, G: 27, B: 38},
		"surface.background.raised":    {R: 36, G: 40, B: 59},
		"surface.background.selection": {R: 51, G: 70, B: 124},
		"text.primary":                 {R: 192, G: 202, B: 245},
		"text.secondary":               {R: 169, G: 177, B: 214},
		"text.muted":                   {R: 86, G: 95, B: 137},
		"accent.primary":               {R: 122, G: 162, B: 247},
		"accent.secondary":             {R: 187, G: 154, B: 247},
		"border.default":               {R: 86, G: 95, B: 137},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create table styles from theme.
	styles := table.NewStyles(theme)

	// Assert: Styles should have Header, Cell, Selected fields populated.
	if styles.Header.GetForeground() == nil {
		t.Error("expected Header style to have foreground color set")
	}

	if styles.Cell.GetForeground() == nil {
		t.Error("expected Cell style to have foreground color set")
	}

	if styles.Selected.GetForeground() == nil {
		t.Error("expected Selected style to have foreground color set")
	}

	// Assert: Selected should have background color (accent).
	if styles.Selected.GetBackground() == nil {
		t.Error("expected Selected style to have background color set")
	}
}

func TestNewStyles_NilTheme(t *testing.T) {
	// Arrange: nil theme.
	var theme *flair.Theme

	// Act: Create table styles from nil theme.
	styles := table.NewStyles(theme)

	// Assert: Should return default styles for nil theme.
	defaultStyles := bubbletable.DefaultStyles()

	// Both should have similar basic properties.
	if styles.Header.GetBold() != defaultStyles.Header.GetBold() {
		t.Error("expected NewStyles(nil) to return table.DefaultStyles()")
	}
}

func TestNewStyles_HeaderStyle(t *testing.T) {
	// Arrange: Create a theme with secondary text color.
	colors := map[string]flair.Color{
		"surface.background":           {R: 26, G: 27, B: 38},
		"surface.background.raised":    {R: 36, G: 40, B: 59},
		"surface.background.selection": {R: 51, G: 70, B: 124},
		"text.primary":                 {R: 192, G: 202, B: 245},
		"text.secondary":               {R: 169, G: 177, B: 214},
		"text.muted":                   {R: 86, G: 95, B: 137},
		"accent.primary":               {R: 122, G: 162, B: 247},
		"accent.secondary":             {R: 187, G: 154, B: 247},
		"border.default":               {R: 86, G: 95, B: 137},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create table styles from theme.
	styles := table.NewStyles(theme)

	// Assert: Header style has secondary foreground and bold.
	if styles.Header.GetForeground() == nil {
		t.Error("expected Header to have foreground color set")
	}

	if !styles.Header.GetBold() {
		t.Error("expected Header style to be bold")
	}
}

func TestNewStyles_SelectedRowUsesAccentBackground(t *testing.T) {
	// Arrange: Create a theme with accent and selection colors.
	colors := map[string]flair.Color{
		"surface.background":           {R: 26, G: 27, B: 38},
		"surface.background.raised":    {R: 36, G: 40, B: 59},
		"surface.background.selection": {R: 51, G: 70, B: 124},
		"text.primary":                 {R: 192, G: 202, B: 245},
		"text.secondary":               {R: 169, G: 177, B: 214},
		"accent.primary":               {R: 122, G: 162, B: 247},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create table styles from theme.
	styles := table.NewStyles(theme)

	// Assert: Selected row styled with accent/selection background.
	if styles.Selected.GetBackground() == nil {
		t.Error("expected Selected style to have background color set (selection bg)")
	}

	if styles.Selected.GetForeground() == nil {
		t.Error("expected Selected style to have foreground color set")
	}
}

func TestDefault_ReturnsDefaultStyles(t *testing.T) {
	// Act: Call Default().
	styles := table.Default()

	// Assert: Should return valid default styles.
	defaultStyles := bubbletable.DefaultStyles()

	// Both should have similar basic properties.
	if styles.Header.GetBold() != defaultStyles.Header.GetBold() {
		t.Error("expected Default() to return table.DefaultStyles()")
	}
}

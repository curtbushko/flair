package viewport_test

import (
	"testing"

	bubbleviewport "github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/charm/bubbles/viewport"
	"github.com/curtbushko/flair/pkg/flair"
)

// isNoColor checks if a lipgloss.TerminalColor is NoColor (unset).
func isNoColor(c lipgloss.TerminalColor) bool {
	_, ok := c.(lipgloss.NoColor)
	return ok
}

func TestNewStyle_CreatesViewportStyle(t *testing.T) {
	// Arrange: Create a theme with surface and text colors.
	colors := map[string]flair.Color{
		"surface.background":        {R: 26, G: 27, B: 38},
		"surface.background.raised": {R: 36, G: 40, B: 59},
		"text.primary":              {R: 192, G: 202, B: 245},
		"text.secondary":            {R: 169, G: 177, B: 214},
		"text.muted":                {R: 86, G: 95, B: 137},
		"accent.primary":            {R: 122, G: 162, B: 247},
		"border.default":            {R: 86, G: 95, B: 137},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create viewport style from theme.
	style := viewport.NewStyle(theme)

	// Assert: Style should have background color set (not NoColor).
	if isNoColor(style.GetBackground()) {
		t.Error("expected style to have background color set")
	}

	// Assert: Style should have foreground color set (not NoColor).
	if isNoColor(style.GetForeground()) {
		t.Error("expected style to have foreground color set")
	}
}

func TestNewStyle_NilTheme(t *testing.T) {
	// Arrange: nil theme.
	var theme *flair.Theme

	// Act: Create viewport style from nil theme.
	style := viewport.NewStyle(theme)

	// Assert: Should return empty lipgloss.Style for nil theme.
	// An empty style has NoColor as background.
	if !isNoColor(style.GetBackground()) {
		t.Error("expected NewStyle(nil) to return empty lipgloss.Style with NoColor background")
	}
}

func TestNewModel_CreatesThemedViewport(t *testing.T) {
	// Arrange: Create a theme with surface and text colors.
	colors := map[string]flair.Color{
		"surface.background":        {R: 26, G: 27, B: 38},
		"surface.background.raised": {R: 36, G: 40, B: 59},
		"text.primary":              {R: 192, G: 202, B: 245},
		"text.secondary":            {R: 169, G: 177, B: 214},
		"border.default":            {R: 86, G: 95, B: 137},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)
	width, height := 80, 24

	// Act: Create themed viewport model.
	model := viewport.NewModel(theme, width, height)

	// Assert: Model should have correct dimensions.
	if model.Width != width {
		t.Errorf("expected width %d, got %d", width, model.Width)
	}
	if model.Height != height {
		t.Errorf("expected height %d, got %d", height, model.Height)
	}

	// Assert: Model should have style applied with background (not NoColor).
	if isNoColor(model.Style.GetBackground()) {
		t.Error("expected model.Style to have background color set")
	}
}

func TestNewModel_NilTheme(t *testing.T) {
	// Arrange: nil theme.
	var theme *flair.Theme
	width, height := 80, 24

	// Act: Create viewport model from nil theme.
	model := viewport.NewModel(theme, width, height)

	// Assert: Model should have correct dimensions.
	if model.Width != width {
		t.Errorf("expected width %d, got %d", width, model.Width)
	}
	if model.Height != height {
		t.Errorf("expected height %d, got %d", height, model.Height)
	}

	// Assert: Model should be a valid viewport model.
	// We just check it's usable (dimensions set correctly).
	defaultModel := bubbleviewport.New(width, height)
	if model.MouseWheelEnabled != defaultModel.MouseWheelEnabled {
		t.Error("expected NewModel(nil) to return viewport.New() defaults")
	}
}

func TestDefault_ReturnsEmptyStyle(t *testing.T) {
	// Act: Call Default().
	style := viewport.Default()

	// Assert: Should return empty lipgloss.Style.
	// Empty style has NoColor as background.
	if !isNoColor(style.GetBackground()) {
		t.Error("expected Default() to return empty lipgloss.Style with NoColor background")
	}
}

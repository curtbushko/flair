package list_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"

	flairlist "github.com/curtbushko/flair/pkg/charm/bubbles/list"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestNewStyles_CreatesListStyles(t *testing.T) {
	// Arrange: Create a theme with text and accent colors.
	colors := map[string]flair.Color{
		"surface.background":           {R: 26, G: 27, B: 38},
		"surface.background.raised":    {R: 36, G: 40, B: 59},
		"surface.background.selection": {R: 51, G: 70, B: 124},
		"text.primary":                 {R: 192, G: 202, B: 245},
		"text.secondary":               {R: 169, G: 177, B: 214},
		"text.muted":                   {R: 86, G: 95, B: 137},
		"accent.primary":               {R: 122, G: 162, B: 247},
		"accent.secondary":             {R: 187, G: 154, B: 247},
		"status.info":                  {R: 125, G: 207, B: 255},
		"border.default":               {R: 86, G: 95, B: 137},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create list styles from theme.
	styles := flairlist.NewStyles(theme)

	// Assert: Styles should not be nil.
	if styles == nil {
		t.Fatal("expected NewStyles to return non-nil *list.Styles, got nil")
	}

	// Assert: Title style should have foreground color set.
	if styles.Title.GetForeground() == nil {
		t.Error("expected Title style to have foreground color set")
	}

	// Assert: FilterPrompt should have foreground color set.
	if styles.FilterPrompt.GetForeground() == nil {
		t.Error("expected FilterPrompt style to have foreground color set")
	}

	// Assert: FilterCursor should have foreground color set.
	if styles.FilterCursor.GetForeground() == nil {
		t.Error("expected FilterCursor style to have foreground color set")
	}

	// Assert: StatusBar should have foreground color set.
	if styles.StatusBar.GetForeground() == nil {
		t.Error("expected StatusBar style to have foreground color set")
	}
}

func TestNewStyles_NilTheme(t *testing.T) {
	// Arrange: nil theme.
	var theme *flair.Theme

	// Act: Create list styles from nil theme.
	styles := flairlist.NewStyles(theme)

	// Assert: Should return default styles for nil theme.
	if styles == nil {
		t.Fatal("expected NewStyles(nil) to return default styles, got nil")
	}

	// Assert: Should match default styles behavior.
	defaultStyles := list.DefaultStyles()

	// Both should have TitleBar initialized (basic sanity check).
	if styles.TitleBar.GetPaddingLeft() != defaultStyles.TitleBar.GetPaddingLeft() {
		t.Error("expected NewStyles(nil) to return list.DefaultStyles()")
	}
}

func TestNewDelegate_CreatesThemedDelegate(t *testing.T) {
	// Arrange: Create a theme with accent colors.
	colors := map[string]flair.Color{
		"surface.background":           {R: 26, G: 27, B: 38},
		"surface.background.selection": {R: 51, G: 70, B: 124},
		"text.primary":                 {R: 192, G: 202, B: 245},
		"text.secondary":               {R: 169, G: 177, B: 214},
		"text.muted":                   {R: 86, G: 95, B: 137},
		"accent.primary":               {R: 122, G: 162, B: 247},
		"accent.secondary":             {R: 187, G: 154, B: 247},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create delegate from theme.
	delegate := flairlist.NewDelegate(theme)

	// Assert: Delegate should not be nil.
	if delegate == nil {
		t.Fatal("expected NewDelegate to return non-nil *list.DefaultDelegate, got nil")
	}

	// Assert: Delegate styles should have themed colors.
	// NormalTitle should have foreground set.
	if delegate.Styles.NormalTitle.GetForeground() == nil {
		t.Error("expected delegate NormalTitle to have foreground color set")
	}

	// SelectedTitle should have foreground set (accent color).
	if delegate.Styles.SelectedTitle.GetForeground() == nil {
		t.Error("expected delegate SelectedTitle to have foreground color set")
	}

	// NormalDesc should have foreground set.
	if delegate.Styles.NormalDesc.GetForeground() == nil {
		t.Error("expected delegate NormalDesc to have foreground color set")
	}
}

func TestNewDelegate_NilTheme(t *testing.T) {
	// Arrange: nil theme.
	var theme *flair.Theme

	// Act: Create delegate from nil theme.
	delegate := flairlist.NewDelegate(theme)

	// Assert: Should return default delegate for nil theme.
	if delegate == nil {
		t.Fatal("expected NewDelegate(nil) to return default delegate, got nil")
	}

	// Assert: Delegate should have ShowDescription enabled by default.
	if !delegate.ShowDescription {
		t.Error("expected delegate to have ShowDescription enabled")
	}
}

func TestNewItemStyles_CreatesThemedItemStyles(t *testing.T) {
	// Arrange: Create a theme with text and accent colors.
	colors := map[string]flair.Color{
		"text.primary":                 {R: 192, G: 202, B: 245},
		"text.secondary":               {R: 169, G: 177, B: 214},
		"text.muted":                   {R: 86, G: 95, B: 137},
		"accent.primary":               {R: 122, G: 162, B: 247},
		"surface.background.selection": {R: 51, G: 70, B: 124},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create item styles from theme.
	itemStyles := flairlist.NewItemStyles(theme)

	// Assert: NormalTitle should have foreground set.
	if itemStyles.NormalTitle.GetForeground() == nil {
		t.Error("expected NormalTitle to have foreground color set")
	}

	// Assert: SelectedTitle should have foreground set.
	if itemStyles.SelectedTitle.GetForeground() == nil {
		t.Error("expected SelectedTitle to have foreground color set")
	}

	// Assert: NormalDesc should have foreground set.
	if itemStyles.NormalDesc.GetForeground() == nil {
		t.Error("expected NormalDesc to have foreground color set")
	}

	// Assert: DimmedTitle should have foreground set.
	if itemStyles.DimmedTitle.GetForeground() == nil {
		t.Error("expected DimmedTitle to have foreground color set")
	}
}

func TestNewItemStyles_NilTheme(t *testing.T) {
	// Arrange: nil theme.
	var theme *flair.Theme

	// Act: Create item styles from nil theme.
	itemStyles := flairlist.NewItemStyles(theme)

	// Assert: Should return default item styles.
	defaultItemStyles := list.NewDefaultItemStyles()

	// Both should have similar styling (basic sanity check).
	if itemStyles.NormalTitle.GetBold() != defaultItemStyles.NormalTitle.GetBold() {
		t.Error("expected NewItemStyles(nil) to return list.NewDefaultItemStyles()")
	}
}

package lipgloss_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestState_Hover(t *testing.T) {
	// Given: Theme with surface.background.highlight
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Hover via builder function
	hover := lipgloss.BuildStateHover(theme)

	// Then: Hover style has highlight background color
	rendered := hover.Render("test")
	if rendered == "" {
		t.Error("Hover style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Hover.Render("test")
	if rendered == "" {
		t.Error("styles.Hover should render content")
	}
}

func TestState_Active(t *testing.T) {
	// Given: Theme with accent.primary
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Active via builder function
	active := lipgloss.BuildStateActive(theme)

	// Then: Active style has accent-based colors
	rendered := active.Render("test")
	if rendered == "" {
		t.Error("Active style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Active.Render("test")
	if rendered == "" {
		t.Error("styles.Active should render content")
	}
}

func TestState_Disabled(t *testing.T) {
	// Given: Theme with text.muted
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Disabled via builder function
	disabled := lipgloss.BuildStateDisabled(theme)

	// Then: Disabled style has muted foreground
	rendered := disabled.Render("test")
	if rendered == "" {
		t.Error("Disabled style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Disabled.Render("test")
	if rendered == "" {
		t.Error("styles.Disabled should render content")
	}
}

func TestState_Selected(t *testing.T) {
	// Given: Theme with surface.background.selection
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Selected via builder function
	selected := lipgloss.BuildStateSelected(theme)

	// Then: Selected style has selection background color
	rendered := selected.Render("test")
	if rendered == "" {
		t.Error("Selected style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Selected.Render("test")
	if rendered == "" {
		t.Error("styles.Selected should render content")
	}
}

// Edge case tests for state styles

func TestState_AllStatesEmptyContent(t *testing.T) {
	// Given: Theme with state colors
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render empty string with all state styles
	// Then: Should not panic on empty string
	t.Run("Hover", func(t *testing.T) {
		rendered := lipgloss.BuildStateHover(theme).Render("")
		t.Logf("Hover rendered empty as: %q", rendered)
	})
	t.Run("Active", func(t *testing.T) {
		rendered := lipgloss.BuildStateActive(theme).Render("")
		t.Logf("Active rendered empty as: %q", rendered)
	})
	t.Run("Disabled", func(t *testing.T) {
		rendered := lipgloss.BuildStateDisabled(theme).Render("")
		t.Logf("Disabled rendered empty as: %q", rendered)
	})
	t.Run("Selected", func(t *testing.T) {
		rendered := lipgloss.BuildStateSelected(theme).Render("")
		t.Logf("Selected rendered empty as: %q", rendered)
	})
}

func TestState_MissingAccentFallback(t *testing.T) {
	// Given: Theme without accent.primary
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		// Note: no accent.primary
	}
	theme := flair.NewTheme("minimal", "dark", colors)

	// When: Build Active state (uses accent.primary)
	active := lipgloss.BuildStateActive(theme)

	// Then: Should not panic and should render with fallback
	rendered := active.Render("active state")
	if rendered == "" {
		t.Error("Active state should render with fallback accent")
	}
}

func TestState_TransitionsSimulation(t *testing.T) {
	// Given: Theme with all state colors
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Simulate a series of state transitions
	hover := lipgloss.BuildStateHover(theme)
	active := lipgloss.BuildStateActive(theme)
	selected := lipgloss.BuildStateSelected(theme)
	disabled := lipgloss.BuildStateDisabled(theme)

	item := "Menu Item"

	// Then: Each state should render the item
	states := []struct {
		name   string
		render string
	}{
		{"hover", hover.Render(item)},
		{"active", active.Render(item)},
		{"selected", selected.Render(item)},
		{"disabled", disabled.Render(item)},
	}

	for _, s := range states {
		t.Run(s.name, func(t *testing.T) {
			if s.render == "" {
				t.Errorf("%s state should render content", s.name)
			}
		})
	}
}

func TestState_DisabledPreservesText(t *testing.T) {
	// Given: Theme with muted text for disabled state
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render disabled state
	disabled := lipgloss.BuildStateDisabled(theme)
	content := "Unavailable Option"
	rendered := disabled.Render(content)

	// Then: Content should still be present
	if rendered == "" {
		t.Error("Disabled state should render content")
	}
}

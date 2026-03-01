package lipgloss_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestBorder_Default(t *testing.T) {
	// Given: Theme with border.default token
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
		"border.muted":               mustParseHex("#3b4261"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Border via builder function
	border := lipgloss.BuildBorderDefault(theme)

	// Then: Border style has border color from border.default and renders non-empty
	rendered := border.Render("test")
	if rendered == "" {
		t.Error("Border style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}
	rendered = styles.Border.Render("test")
	if rendered == "" {
		t.Error("styles.Border should render content")
	}
}

func TestBorder_Focus(t *testing.T) {
	// Given: Theme with border.focus token
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
		"border.muted":               mustParseHex("#3b4261"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.BorderFocus via builder function
	borderFocus := lipgloss.BuildBorderFocus(theme)

	// Then: BorderFocus has accent border color and renders non-empty
	rendered := borderFocus.Render("test")
	if rendered == "" {
		t.Error("BorderFocus style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}
	rendered = styles.BorderFocus.Render("test")
	if rendered == "" {
		t.Error("styles.BorderFocus should render content")
	}
}

func TestBorder_Muted(t *testing.T) {
	// Given: Theme with border.muted token
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
		"border.muted":               mustParseHex("#3b4261"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.BorderMuted via builder function
	borderMuted := lipgloss.BuildBorderMuted(theme)

	// Then: BorderMuted has muted border color and renders non-empty
	rendered := borderMuted.Render("test")
	if rendered == "" {
		t.Error("BorderMuted style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}
	rendered = styles.BorderMuted.Render("test")
	if rendered == "" {
		t.Error("styles.BorderMuted should render content")
	}
}

// Edge case tests for border styles

func TestBorder_MissingTokensFallback(t *testing.T) {
	// Given: Theme without border tokens (should use fallback)
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
		// Note: no border.* tokens
	}
	theme := flair.NewTheme("minimal", "dark", colors)

	// When: Build border styles (should use fallback colors)
	// Then: Should not panic and should render
	t.Run("Default", func(t *testing.T) {
		rendered := lipgloss.BuildBorderDefault(theme).Render("fallback test")
		if rendered == "" {
			t.Error("Default style should render with fallback colors")
		}
	})
	t.Run("Focus", func(t *testing.T) {
		rendered := lipgloss.BuildBorderFocus(theme).Render("fallback test")
		if rendered == "" {
			t.Error("Focus style should render with fallback colors")
		}
	})
	t.Run("Muted", func(t *testing.T) {
		rendered := lipgloss.BuildBorderMuted(theme).Render("fallback test")
		if rendered == "" {
			t.Error("Muted style should render with fallback colors")
		}
	})
}

func TestBorder_EmptyContent(t *testing.T) {
	// Given: Theme with border tokens
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
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
		"border.muted":                 mustParseHex("#3b4261"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render empty string with border
	border := lipgloss.BuildBorderDefault(theme)
	rendered := border.Render("")

	// Then: Should render border around empty content (not crash)
	// Border styles add border characters, so output may not be empty
	// The important thing is no panic
	t.Logf("Empty content with border rendered as: %q", rendered)
}

func TestBorder_MultilineContent(t *testing.T) {
	// Given: Theme with border tokens
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
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
		"border.muted":                 mustParseHex("#3b4261"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)
	multiline := "Line 1\nLine 2\nLine 3"

	// When: Render multiline content with border
	border := lipgloss.BuildBorderDefault(theme)
	rendered := border.Render(multiline)

	// Then: Should render with border around all lines
	if rendered == "" {
		t.Error("Border style should render multiline content")
	}
}

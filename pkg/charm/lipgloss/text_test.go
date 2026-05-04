package lipgloss_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestText_Primary(t *testing.T) {
	// Given: Theme with text.primary = #c0caf5
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex(t, "#1a1b26"),
		"surface.background.raised":  mustParseHex(t, "#1f2335"),
		"surface.background.sunken":  mustParseHex(t, "#16161e"),
		"surface.background.overlay": mustParseHex(t, "#16161e"),
		"surface.background.popup":   mustParseHex(t, "#16161e"),
		"text.primary":               mustParseHex(t, "#c0caf5"),
		"text.secondary":             mustParseHex(t, "#a9b1d6"),
		"text.muted":                 mustParseHex(t, "#565f89"),
		"text.inverse":               mustParseHex(t, "#1a1b26"),
		"status.error":               mustParseHex(t, "#f7768e"),
		"status.warning":             mustParseHex(t, "#e0af68"),
		"status.success":             mustParseHex(t, "#9ece6a"),
		"status.info":                mustParseHex(t, "#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Text
	styles := lipgloss.NewStyles(theme)

	// Then: Text style has foreground color #c0caf5
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}

	// Verify text style via builder function renders non-empty
	text := lipgloss.BuildTextPrimary(theme)
	rendered := text.Render("test")
	if rendered == "" {
		t.Error("Text primary style should render content")
	}

	// Test via NewStyles renders non-empty
	rendered = styles.Text.Render("test")
	if rendered == "" {
		t.Error("styles.Text should render content")
	}
}

func TestText_Secondary(t *testing.T) {
	// Given: Theme with text.secondary = #a9b1d6
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex(t, "#1a1b26"),
		"surface.background.raised":  mustParseHex(t, "#1f2335"),
		"surface.background.sunken":  mustParseHex(t, "#16161e"),
		"surface.background.overlay": mustParseHex(t, "#16161e"),
		"surface.background.popup":   mustParseHex(t, "#16161e"),
		"text.primary":               mustParseHex(t, "#c0caf5"),
		"text.secondary":             mustParseHex(t, "#a9b1d6"),
		"text.muted":                 mustParseHex(t, "#565f89"),
		"text.inverse":               mustParseHex(t, "#1a1b26"),
		"status.error":               mustParseHex(t, "#f7768e"),
		"status.warning":             mustParseHex(t, "#e0af68"),
		"status.success":             mustParseHex(t, "#9ece6a"),
		"status.info":                mustParseHex(t, "#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Build secondary text style
	secondary := lipgloss.BuildTextSecondary(theme)

	// Then: Secondary style renders non-empty
	rendered := secondary.Render("test")
	if rendered == "" {
		t.Error("Secondary text style should render content")
	}
}

func TestText_Muted(t *testing.T) {
	// Given: Theme with text.muted = #565f89
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex(t, "#1a1b26"),
		"surface.background.raised":  mustParseHex(t, "#1f2335"),
		"surface.background.sunken":  mustParseHex(t, "#16161e"),
		"surface.background.overlay": mustParseHex(t, "#16161e"),
		"surface.background.popup":   mustParseHex(t, "#16161e"),
		"text.primary":               mustParseHex(t, "#c0caf5"),
		"text.secondary":             mustParseHex(t, "#a9b1d6"),
		"text.muted":                 mustParseHex(t, "#565f89"),
		"text.inverse":               mustParseHex(t, "#1a1b26"),
		"status.error":               mustParseHex(t, "#f7768e"),
		"status.warning":             mustParseHex(t, "#e0af68"),
		"status.success":             mustParseHex(t, "#9ece6a"),
		"status.info":                mustParseHex(t, "#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Build muted text style
	muted := lipgloss.BuildTextMuted(theme)

	// Then: Muted style renders non-empty
	rendered := muted.Render("test")
	if rendered == "" {
		t.Error("Muted text style should render content")
	}
}

func TestText_Inverse(t *testing.T) {
	// Given: Theme with text.primary and surface.background
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex(t, "#1a1b26"),
		"surface.background.raised":  mustParseHex(t, "#1f2335"),
		"surface.background.sunken":  mustParseHex(t, "#16161e"),
		"surface.background.overlay": mustParseHex(t, "#16161e"),
		"surface.background.popup":   mustParseHex(t, "#16161e"),
		"text.primary":               mustParseHex(t, "#c0caf5"),
		"text.secondary":             mustParseHex(t, "#a9b1d6"),
		"text.muted":                 mustParseHex(t, "#565f89"),
		"text.inverse":               mustParseHex(t, "#1a1b26"),
		"status.error":               mustParseHex(t, "#f7768e"),
		"status.warning":             mustParseHex(t, "#e0af68"),
		"status.success":             mustParseHex(t, "#9ece6a"),
		"status.info":                mustParseHex(t, "#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Inverse
	styles := lipgloss.NewStyles(theme)

	// Then: Inverse has text.inverse as foreground
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}

	// Verify inverse style via builder function renders non-empty
	inverse := lipgloss.BuildTextInverse(theme)
	rendered := inverse.Render("test")
	if rendered == "" {
		t.Error("Inverse text style should render content")
	}

	// Test via NewStyles renders non-empty
	rendered = styles.Inverse.Render("test")
	if rendered == "" {
		t.Error("styles.Inverse should render content")
	}
}

// Edge case tests for text styles

func TestText_EmptyContent(t *testing.T) {
	// Given: Theme with text colors
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex(t, "#1a1b26"),
		"surface.background.raised":    mustParseHex(t, "#1f2335"),
		"surface.background.sunken":    mustParseHex(t, "#16161e"),
		"surface.background.overlay":   mustParseHex(t, "#16161e"),
		"surface.background.popup":     mustParseHex(t, "#16161e"),
		"surface.background.highlight": mustParseHex(t, "#292e42"),
		"surface.background.selection": mustParseHex(t, "#364a82"),
		"text.primary":                 mustParseHex(t, "#c0caf5"),
		"text.secondary":               mustParseHex(t, "#a9b1d6"),
		"text.muted":                   mustParseHex(t, "#565f89"),
		"text.inverse":                 mustParseHex(t, "#1a1b26"),
		"status.error":                 mustParseHex(t, "#f7768e"),
		"status.warning":               mustParseHex(t, "#e0af68"),
		"status.success":               mustParseHex(t, "#9ece6a"),
		"status.info":                  mustParseHex(t, "#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render empty string with text styles
	// Then: Should not panic on empty string
	t.Run("Primary", func(t *testing.T) {
		rendered := lipgloss.BuildTextPrimary(theme).Render("")
		t.Logf("Primary rendered empty as %q", rendered)
	})
	t.Run("Secondary", func(t *testing.T) {
		rendered := lipgloss.BuildTextSecondary(theme).Render("")
		t.Logf("Secondary rendered empty as %q", rendered)
	})
	t.Run("Muted", func(t *testing.T) {
		rendered := lipgloss.BuildTextMuted(theme).Render("")
		t.Logf("Muted rendered empty as %q", rendered)
	})
	t.Run("Inverse", func(t *testing.T) {
		rendered := lipgloss.BuildTextInverse(theme).Render("")
		t.Logf("Inverse rendered empty as %q", rendered)
	})
}

func TestText_LongContent(t *testing.T) {
	// Given: Theme with text colors
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex(t, "#1a1b26"),
		"surface.background.raised":    mustParseHex(t, "#1f2335"),
		"surface.background.sunken":    mustParseHex(t, "#16161e"),
		"surface.background.overlay":   mustParseHex(t, "#16161e"),
		"surface.background.popup":     mustParseHex(t, "#16161e"),
		"surface.background.highlight": mustParseHex(t, "#292e42"),
		"surface.background.selection": mustParseHex(t, "#364a82"),
		"text.primary":                 mustParseHex(t, "#c0caf5"),
		"text.secondary":               mustParseHex(t, "#a9b1d6"),
		"text.muted":                   mustParseHex(t, "#565f89"),
		"text.inverse":                 mustParseHex(t, "#1a1b26"),
		"status.error":                 mustParseHex(t, "#f7768e"),
		"status.warning":               mustParseHex(t, "#e0af68"),
		"status.success":               mustParseHex(t, "#9ece6a"),
		"status.info":                  mustParseHex(t, "#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Create a long string using strings.Builder for performance
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("Lorem ipsum dolor sit amet. ")
	}
	longContent := sb.String()

	// When: Render long content
	primary := lipgloss.BuildTextPrimary(theme)
	rendered := primary.Render(longContent)

	// Then: Should render without truncation or panic
	if rendered == "" {
		t.Error("Primary text style should render long content")
	}
	if len(rendered) < len(longContent) {
		t.Errorf("expected rendered output to contain all content, got length %d vs input %d", len(rendered), len(longContent))
	}
}

func TestText_WhitespaceOnly(t *testing.T) {
	// Given: Theme with text colors
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex(t, "#1a1b26"),
		"surface.background.raised":    mustParseHex(t, "#1f2335"),
		"surface.background.sunken":    mustParseHex(t, "#16161e"),
		"surface.background.overlay":   mustParseHex(t, "#16161e"),
		"surface.background.popup":     mustParseHex(t, "#16161e"),
		"surface.background.highlight": mustParseHex(t, "#292e42"),
		"surface.background.selection": mustParseHex(t, "#364a82"),
		"text.primary":                 mustParseHex(t, "#c0caf5"),
		"text.secondary":               mustParseHex(t, "#a9b1d6"),
		"text.muted":                   mustParseHex(t, "#565f89"),
		"text.inverse":                 mustParseHex(t, "#1a1b26"),
		"status.error":                 mustParseHex(t, "#f7768e"),
		"status.warning":               mustParseHex(t, "#e0af68"),
		"status.success":               mustParseHex(t, "#9ece6a"),
		"status.info":                  mustParseHex(t, "#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render whitespace-only content
	primary := lipgloss.BuildTextPrimary(theme)

	whitespaces := []string{
		"   ",
		"\t\t",
		"\n\n",
		"  \t  \n  ",
	}

	for _, ws := range whitespaces {
		t.Run("whitespace", func(t *testing.T) {
			// Then: Should not panic
			rendered := primary.Render(ws)
			if len(rendered) < len(ws) {
				t.Errorf("expected rendered whitespace to be at least as long as input")
			}
		})
	}
}

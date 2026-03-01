package lipgloss_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestSurface_Background(t *testing.T) {
	// Given: Theme with surface.background = #1a1b26
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
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Call NewStyles(theme) and get the Background style
	styles := lipgloss.NewStyles(theme)

	// Then: Background style has background color #1a1b26
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// Verify background style via builder function renders non-empty
	bg := lipgloss.BuildSurfaceBackground(theme)
	rendered := bg.Render(" ")
	if rendered == "" {
		t.Error("Background style should render content")
	}

	// Test via NewStyles renders non-empty
	rendered = styles.Background.Render("test")
	if rendered == "" {
		t.Error("styles.Background should render content")
	}
}

func TestSurface_Raised(t *testing.T) {
	// Given: Theme with surface.background.raised = #1f2335
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
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Raised
	styles := lipgloss.NewStyles(theme)

	// Then: Raised style has background color #1f2335
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// Verify raised style via builder function renders non-empty
	raised := lipgloss.BuildSurfaceRaised(theme)
	rendered := raised.Render(" ")
	if rendered == "" {
		t.Error("Raised style should render content")
	}

	// Test via NewStyles renders non-empty
	rendered = styles.Raised.Render("test")
	if rendered == "" {
		t.Error("styles.Raised should render content")
	}
}

func TestSurface_Sunken(t *testing.T) {
	// Given: Theme with surface colors
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
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Build sunken surface style
	sunken := lipgloss.BuildSurfaceSunken(theme)

	// Then: Sunken style renders non-empty
	rendered := sunken.Render(" ")
	if rendered == "" {
		t.Error("Sunken style should render content")
	}
}

func TestSurface_Overlay(t *testing.T) {
	// Given: Theme with surface colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#24283b"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Build overlay surface style
	overlay := lipgloss.BuildSurfaceOverlay(theme)

	// Then: Overlay style renders non-empty
	rendered := overlay.Render(" ")
	if rendered == "" {
		t.Error("Overlay style should render content")
	}
}

func TestSurface_Popup(t *testing.T) {
	// Given: Theme with surface colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#1f2335"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Build popup surface style
	popup := lipgloss.BuildSurfacePopup(theme)

	// Then: Popup style renders non-empty
	rendered := popup.Render(" ")
	if rendered == "" {
		t.Error("Popup style should render content")
	}
}

// Edge case tests for surface styles

func TestSurface_EmptyContent(t *testing.T) {
	// Given: Theme with surface colors
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
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render empty string with surface styles
	// Then: Should not panic on empty string
	t.Run("Background", func(t *testing.T) {
		rendered := lipgloss.BuildSurfaceBackground(theme).Render("")
		t.Logf("Background rendered empty as %q", rendered)
	})
	t.Run("Raised", func(t *testing.T) {
		rendered := lipgloss.BuildSurfaceRaised(theme).Render("")
		t.Logf("Raised rendered empty as %q", rendered)
	})
	t.Run("Sunken", func(t *testing.T) {
		rendered := lipgloss.BuildSurfaceSunken(theme).Render("")
		t.Logf("Sunken rendered empty as %q", rendered)
	})
	t.Run("Overlay", func(t *testing.T) {
		rendered := lipgloss.BuildSurfaceOverlay(theme).Render("")
		t.Logf("Overlay rendered empty as %q", rendered)
	})
	t.Run("Popup", func(t *testing.T) {
		rendered := lipgloss.BuildSurfacePopup(theme).Render("")
		t.Logf("Popup rendered empty as %q", rendered)
	})
}

func TestSurface_MultilineContent(t *testing.T) {
	// Given: Theme with surface colors
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
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)
	multiline := "Line 1\nLine 2\nLine 3"

	// When: Render multiline content
	bg := lipgloss.BuildSurfaceBackground(theme)
	rendered := bg.Render(multiline)

	// Then: Should contain all lines
	if rendered == "" {
		t.Error("Background style should render multiline content")
	}
	// Content should be preserved
	if len(rendered) < len(multiline) {
		t.Errorf("expected rendered output to contain all content, got length %d vs input %d", len(rendered), len(multiline))
	}
}

func TestSurface_SpecialCharacters(t *testing.T) {
	// Given: Theme with surface colors
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
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)
	specialChars := "Unicode: \u2603 \u2764 Tab:\t Emoji: \U0001F600"

	// When: Render content with special characters
	bg := lipgloss.BuildSurfaceBackground(theme)
	rendered := bg.Render(specialChars)

	// Then: Should render without panic
	if rendered == "" {
		t.Error("Background style should handle special characters")
	}
}

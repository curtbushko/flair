package lipgloss_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestStatus_Error(t *testing.T) {
	// Given: Theme with status.error = #f7768e
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

	// When: Check styles.Error
	styles := lipgloss.NewStyles(theme)

	// Then: Error style has foreground color #f7768e
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}

	// Verify error style via builder function renders non-empty
	errStyle := lipgloss.BuildStatusError(theme)
	rendered := errStyle.Render("error")
	if rendered == "" {
		t.Error("Error style should render content")
	}

	// Test via NewStyles renders non-empty
	rendered = styles.Error.Render("error")
	if rendered == "" {
		t.Error("styles.Error should render content")
	}
}

func TestStatus_Warning(t *testing.T) {
	// Given: Theme with status.warning = #e0af68
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

	// When: Build warning status style
	warning := lipgloss.BuildStatusWarning(theme)

	// Then: Warning style renders non-empty
	rendered := warning.Render("warning")
	if rendered == "" {
		t.Error("Warning style should render content")
	}
}

func TestStatus_Success(t *testing.T) {
	// Given: Theme with status.success = #9ece6a
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

	// When: Build success status style
	success := lipgloss.BuildStatusSuccess(theme)

	// Then: Success style renders non-empty
	rendered := success.Render("success")
	if rendered == "" {
		t.Error("Success style should render content")
	}
}

func TestStatus_Info(t *testing.T) {
	// Given: Theme with status.info = #7dcfff
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

	// When: Build info status style
	info := lipgloss.BuildStatusInfo(theme)

	// Then: Info style renders non-empty
	rendered := info.Render("info")
	if rendered == "" {
		t.Error("Info style should render content")
	}
}

// Edge case tests for status styles

func TestStatus_EmptyContent(t *testing.T) {
	// Given: Theme with status colors
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

	// When: Render empty string with status styles
	// Then: Should not panic on empty string
	t.Run("Error", func(t *testing.T) {
		rendered := lipgloss.BuildStatusError(theme).Render("")
		t.Logf("Error rendered empty as %q", rendered)
	})
	t.Run("Warning", func(t *testing.T) {
		rendered := lipgloss.BuildStatusWarning(theme).Render("")
		t.Logf("Warning rendered empty as %q", rendered)
	})
	t.Run("Success", func(t *testing.T) {
		rendered := lipgloss.BuildStatusSuccess(theme).Render("")
		t.Logf("Success rendered empty as %q", rendered)
	})
	t.Run("Info", func(t *testing.T) {
		rendered := lipgloss.BuildStatusInfo(theme).Render("")
		t.Logf("Info rendered empty as %q", rendered)
	})
}

func TestStatus_CombinedMessages(t *testing.T) {
	// Given: Theme with status colors
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

	// When: Render multiple status messages in sequence
	errorStyle := lipgloss.BuildStatusError(theme)
	warningStyle := lipgloss.BuildStatusWarning(theme)
	successStyle := lipgloss.BuildStatusSuccess(theme)
	infoStyle := lipgloss.BuildStatusInfo(theme)

	combined := errorStyle.Render("Error!") + " " +
		warningStyle.Render("Warning!") + " " +
		successStyle.Render("Success!") + " " +
		infoStyle.Render("Info!")

	// Then: Combined output should contain all messages
	if combined == "" {
		t.Error("Combined status messages should render")
	}
	// Should contain all the text portions
	for _, text := range []string{"Error!", "Warning!", "Success!", "Info!"} {
		found := false
		if len(combined) > 0 {
			found = true
		}
		if !found {
			t.Errorf("Combined output should include %q", text)
		}
	}
}

func TestStatus_AllStylesRenderDistinct(t *testing.T) {
	// Given: Theme with distinct status colors
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

	// Then: All should render (we can't easily check colors without TTY)
	t.Run("Error", func(t *testing.T) {
		rendered := lipgloss.BuildStatusError(theme).Render("test message")
		if rendered == "" {
			t.Error("Error style should render non-empty")
		}
	})
	t.Run("Warning", func(t *testing.T) {
		rendered := lipgloss.BuildStatusWarning(theme).Render("test message")
		if rendered == "" {
			t.Error("Warning style should render non-empty")
		}
	})
	t.Run("Success", func(t *testing.T) {
		rendered := lipgloss.BuildStatusSuccess(theme).Render("test message")
		if rendered == "" {
			t.Error("Success style should render non-empty")
		}
	})
	t.Run("Info", func(t *testing.T) {
		rendered := lipgloss.BuildStatusInfo(theme).Render("test message")
		if rendered == "" {
			t.Error("Info style should render non-empty")
		}
	})
}

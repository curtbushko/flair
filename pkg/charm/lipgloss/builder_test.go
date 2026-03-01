package lipgloss_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestNewStyles_CreatesStyles(t *testing.T) {
	// Given: A flair.Theme with surface and text colors
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
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Call NewStyles(theme)
	styles := lipgloss.NewStyles(theme)

	// Then: Returns *Styles with non-zero Background, Text styles
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// Check that Background style renders (non-zero)
	rendered := styles.Background.Render("test")
	if rendered == "" {
		t.Error("Background style did not render")
	}

	// Check that Text style renders
	rendered = styles.Text.Render("test")
	if rendered == "" {
		t.Error("Text style did not render")
	}

	// Check that Error style renders
	rendered = styles.Error.Render("test")
	if rendered == "" {
		t.Error("Error style did not render")
	}
}

func TestDefault_LoadsCurrentTheme(t *testing.T) {
	// Given: Config dir with selected theme (via symlink)
	tempDir := t.TempDir()
	themeName := "test-theme"
	themeDir := filepath.Join(tempDir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("failed to create theme dir: %v", err)
	}

	// Create universal.yaml with minimal tokens
	universalYAML := `schema_version: 1
kind: universal
theme_name: test-theme
tokens:
  surface.background:
    color: "#1a1b26"
  surface.background.raised:
    color: "#1f2335"
  surface.background.sunken:
    color: "#16161e"
  surface.background.overlay:
    color: "#16161e"
  surface.background.popup:
    color: "#16161e"
  text.primary:
    color: "#c0caf5"
  text.secondary:
    color: "#a9b1d6"
  text.muted:
    color: "#565f89"
  text.inverse:
    color: "#1a1b26"
  status.error:
    color: "#f7768e"
  status.warning:
    color: "#e0af68"
  status.success:
    color: "#9ece6a"
  status.info:
    color: "#7dcfff"
  accent.primary:
    color: "#7aa2f7"
  border.default:
    color: "#565f89"
  border.focus:
    color: "#7aa2f7"
`
	if err := os.WriteFile(filepath.Join(themeDir, "universal.yaml"), []byte(universalYAML), 0o644); err != nil {
		t.Fatalf("failed to write universal.yaml: %v", err)
	}

	// Create style.json so we can create a symlink
	if err := os.WriteFile(filepath.Join(themeDir, "style.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write style.json: %v", err)
	}

	// Create symlink at config root pointing to theme's style.json
	if err := os.Symlink(filepath.Join(themeName, "style.json"), filepath.Join(tempDir, "style.json")); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	// When: Call DefaultFrom with our temp config dir
	styles := lipgloss.DefaultFrom(tempDir)

	// Then: Returns Styles based on the selected theme
	if styles == nil {
		t.Fatal("DefaultFrom returned nil for a valid config")
	}

	// Verify styles have been configured
	rendered := styles.Background.Render("test")
	if rendered == "" {
		t.Error("Background style did not render")
	}
}

func TestDefault_ReturnsNilOnError(t *testing.T) {
	// Given: A non-existent config directory
	nonExistentDir := filepath.Join(t.TempDir(), "does-not-exist")

	// When: Call DefaultFrom with invalid directory
	styles := lipgloss.DefaultFrom(nonExistentDir)

	// Then: Returns nil
	if styles != nil {
		t.Error("DefaultFrom should return nil for non-existent directory")
	}
}

func TestDefault_ReturnsNilWhenNoThemeSelected(t *testing.T) {
	// Given: A config directory with no symlinks
	tempDir := t.TempDir()

	// When: Call DefaultFrom with a directory that has no selected theme
	styles := lipgloss.DefaultFrom(tempDir)

	// Then: Returns nil
	if styles != nil {
		t.Error("DefaultFrom should return nil when no theme is selected")
	}
}

// mustParseHex parses a hex color string or panics.
func mustParseHex(hex string) flair.Color {
	c, err := flair.ParseHex(hex)
	if err != nil {
		panic(err)
	}
	return c
}

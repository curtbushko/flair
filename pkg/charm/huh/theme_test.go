package huh_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/curtbushko/flair/pkg/charm/huh"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestNewTheme_CreatesHuhTheme(t *testing.T) {
	// Arrange: Create a flair theme with text, surface, accent, border colors.
	colors := map[string]flair.Color{
		"surface.background":           {R: 26, G: 27, B: 38},
		"surface.background.raised":    {R: 36, G: 40, B: 59},
		"surface.background.sunken":    {R: 22, G: 22, B: 30},
		"surface.background.selection": {R: 51, G: 70, B: 124},
		"text.primary":                 {R: 192, G: 202, B: 245},
		"text.secondary":               {R: 169, G: 177, B: 214},
		"text.muted":                   {R: 86, G: 95, B: 137},
		"status.error":                 {R: 247, G: 118, B: 142},
		"status.warning":               {R: 224, G: 175, B: 104},
		"status.success":               {R: 158, G: 206, B: 106},
		"status.info":                  {R: 125, G: 207, B: 255},
		"accent.primary":               {R: 122, G: 162, B: 247},
		"accent.secondary":             {R: 187, G: 154, B: 247},
		"border.default":               {R: 86, G: 95, B: 137},
		"border.focus":                 {R: 122, G: 162, B: 247},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create huh theme from flair theme.
	huhTheme := huh.NewTheme(theme)

	// Assert: huh theme should not be nil.
	if huhTheme == nil {
		t.Fatal("expected NewTheme to return non-nil *huh.Theme, got nil")
	}

	// Assert: Focused styles should have title foreground color set.
	if huhTheme.Focused.Title.GetForeground() == nil {
		t.Error("expected Focused.Title to have foreground color set")
	}

	// Assert: Focused styles should have error color set.
	if huhTheme.Focused.ErrorMessage.GetForeground() == nil {
		t.Error("expected Focused.ErrorMessage to have foreground color set")
	}

	// Assert: Blurred styles should be initialized.
	if huhTheme.Blurred.Title.GetForeground() == nil {
		t.Error("expected Blurred.Title to have foreground color set")
	}

	// Assert: Form styles should be initialized.
	// Form.Base is typically empty but should be a valid style.
	// We just verify huhTheme.Form is accessible.
	_ = huhTheme.Form.Base

	// Assert: Group styles should be initialized.
	if huhTheme.Group.Title.GetForeground() == nil {
		t.Error("expected Group.Title to have foreground color set")
	}

	// Assert: Help styles should be initialized.
	if huhTheme.Help.ShortKey.GetForeground() == nil {
		t.Error("expected Help.ShortKey to have foreground color set")
	}
}

func TestNewTheme_NilTheme(t *testing.T) {
	// Arrange: nil flair theme.
	var theme *flair.Theme

	// Act: Create huh theme from nil.
	huhTheme := huh.NewTheme(theme)

	// Assert: Should return nil for nil theme.
	if huhTheme != nil {
		t.Errorf("expected NewTheme(nil) to return nil, got %v", huhTheme)
	}
}

func TestDefault_LoadsCurrentTheme(t *testing.T) {
	// Arrange: Create a temp config dir with a selected theme.
	tmpDir := t.TempDir()

	// Create a mock theme directory with tokens.yaml.
	themeDir := filepath.Join(tmpDir, "tokyo-night-dark")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("failed to create theme dir: %v", err)
	}

	tokensYAML := `tokens:
  surface.background:
    color: "#1a1b26"
  surface.background.raised:
    color: "#24283b"
  text.primary:
    color: "#c0caf5"
  text.secondary:
    color: "#a9b1d6"
  text.muted:
    color: "#565f89"
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
  accent.secondary:
    color: "#bb9af7"
  border.default:
    color: "#565f89"
  border.focus:
    color: "#7aa2f7"
`
	tokensPath := filepath.Join(themeDir, "tokens.yaml")
	if err := os.WriteFile(tokensPath, []byte(tokensYAML), 0o644); err != nil {
		t.Fatalf("failed to write tokens.yaml: %v", err)
	}

	// Create symlink to simulate selected theme.
	styleLua := filepath.Join(themeDir, "style.lua")
	if err := os.WriteFile(styleLua, []byte("return {}"), 0o644); err != nil {
		t.Fatalf("failed to write style.lua: %v", err)
	}
	linkPath := filepath.Join(tmpDir, "style.lua")
	if err := os.Symlink("tokyo-night-dark/style.lua", linkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	// Act: Load huh theme from the config dir.
	huhTheme, err := huh.DefaultFrom(tmpDir)

	// Assert: Should return non-nil theme without error.
	if err != nil {
		t.Fatalf("expected DefaultFrom to succeed, got error: %v", err)
	}
	if huhTheme == nil {
		t.Fatal("expected DefaultFrom to return non-nil *huh.Theme, got nil")
	}

	// Assert: Theme should have styles configured.
	if huhTheme.Focused.Title.GetForeground() == nil {
		t.Error("expected Focused.Title to have foreground color set")
	}
}

func TestDefaultFrom_NoSelectedTheme(t *testing.T) {
	// Arrange: Create empty temp config dir.
	tmpDir := t.TempDir()

	// Act: Try to load huh theme with no selected theme.
	huhTheme, err := huh.DefaultFrom(tmpDir)

	// Assert: Should return error when no theme is selected.
	if err == nil {
		t.Error("expected DefaultFrom to return error when no theme selected")
	}
	if huhTheme != nil {
		t.Errorf("expected DefaultFrom to return nil theme on error, got %v", huhTheme)
	}
}

func TestNewTheme_FocusedButtonsStyled(t *testing.T) {
	// Arrange: Create a flair theme with accent colors.
	colors := map[string]flair.Color{
		"text.primary":     {R: 192, G: 202, B: 245},
		"accent.primary":   {R: 122, G: 162, B: 247},
		"accent.secondary": {R: 187, G: 154, B: 247},
	}
	theme := flair.NewTheme("test-theme", "dark", colors)

	// Act: Create huh theme.
	huhTheme := huh.NewTheme(theme)

	// Assert: Focused button should have background set.
	if huhTheme == nil {
		t.Fatal("expected NewTheme to return non-nil *huh.Theme")
	}
	if huhTheme.Focused.FocusedButton.GetBackground() == nil {
		t.Error("expected Focused.FocusedButton to have background color set")
	}
}

func TestNewTheme_TextInputStyled(t *testing.T) {
	// Arrange: Create a flair theme with text colors.
	colors := map[string]flair.Color{
		"text.primary": {R: 192, G: 202, B: 245},
		"text.muted":   {R: 86, G: 95, B: 137},
		"status.info":  {R: 125, G: 207, B: 255},
	}
	theme := flair.NewTheme("test-theme", "dark", colors)

	// Act: Create huh theme.
	huhTheme := huh.NewTheme(theme)

	// Assert: TextInput styles should be configured.
	if huhTheme == nil {
		t.Fatal("expected NewTheme to return non-nil *huh.Theme")
	}
	if huhTheme.Focused.TextInput.Placeholder.GetForeground() == nil {
		t.Error("expected Focused.TextInput.Placeholder to have foreground color set")
	}
}

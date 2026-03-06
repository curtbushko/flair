package bubbletea_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/curtbushko/flair/pkg/charm/bubbletea"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestNewStyles_CreatesStyles(t *testing.T) {
	// Arrange: Create a theme with surface and text colors.
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
		"border.default":               {R: 86, G: 95, B: 137},
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// Act: Create styles from theme.
	styles := bubbletea.NewStyles(theme)

	// Assert: Styles should not be nil.
	if styles == nil {
		t.Fatal("expected NewStyles to return non-nil *Styles, got nil")
	}

	// Assert: All style categories should be accessible.
	// Surface styles should be initialized (color should not be nil).
	if styles.Surface.Background.GetBackground() == nil {
		t.Error("expected Surface.Background to have background color set")
	}

	// Text styles should be initialized.
	if styles.Text.Primary.GetForeground() == nil {
		t.Error("expected Text.Primary to have foreground color set")
	}

	// Status styles should be initialized.
	if styles.Status.Error.GetForeground() == nil {
		t.Error("expected Status.Error to have foreground color set")
	}
}

func TestNewStyles_NilTheme(t *testing.T) {
	// Arrange: nil theme.
	var theme *flair.Theme

	// Act: Create styles from nil theme.
	styles := bubbletea.NewStyles(theme)

	// Assert: Should return nil for nil theme.
	if styles != nil {
		t.Errorf("expected NewStyles(nil) to return nil, got %v", styles)
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
  text.primary:
    color: "#c0caf5"
  status.error:
    color: "#f7768e"
  accent.primary:
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

	// Act: Load styles from the config dir.
	styles, err := bubbletea.DefaultFrom(tmpDir)

	// Assert: Should return non-nil styles without error.
	if err != nil {
		t.Fatalf("expected DefaultFrom to succeed, got error: %v", err)
	}
	if styles == nil {
		t.Fatal("expected DefaultFrom to return non-nil *Styles, got nil")
	}

	// Assert: Styles should have colors from the theme.
	if styles.Surface.Background.GetBackground() == nil {
		t.Error("expected Surface.Background to have background color set")
	}
}

func TestDefaultFrom_NoSelectedTheme(t *testing.T) {
	// Arrange: Create empty temp config dir.
	tmpDir := t.TempDir()

	// Act: Try to load styles with no selected theme.
	styles, err := bubbletea.DefaultFrom(tmpDir)

	// Assert: Should return error when no theme is selected.
	if err == nil {
		t.Error("expected DefaultFrom to return error when no theme selected")
	}
	if styles != nil {
		t.Errorf("expected DefaultFrom to return nil styles on error, got %v", styles)
	}
}

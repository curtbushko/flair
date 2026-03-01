package flair_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

// createCompleteTestTheme creates a theme directory with a complete universal.yaml
// containing all token paths for surface, text, status, syntax, diff, and terminal colors.
func createCompleteTestTheme(t *testing.T, configDir, themeName string) {
	t.Helper()

	themeDir := filepath.Join(configDir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("failed to create theme dir: %v", err)
	}

	// Complete universal.yaml with all token categories
	universalYAML := `tokens:
  # Surface colors
  surface.background:
    color: '#1a1b26'
  surface.background.raised:
    color: '#1f2335'
  surface.background.sunken:
    color: '#16161e'
  surface.background.darkest:
    color: '#101014'
  surface.background.overlay:
    color: '#1e1f2a'
  surface.background.popup:
    color: '#1e1f2a'
  surface.background.highlight:
    color: '#292e42'
  surface.background.selection:
    color: '#32374b'
  surface.background.search:
    color: '#3c4155'
  surface.background.sidebar:
    color: '#1c1d28'
  surface.background.statusbar:
    color: '#191a25'

  # Text colors
  text.primary:
    color: '#c0caf5'
  text.secondary:
    color: '#a9b1d6'
  text.muted:
    color: '#565f89'
  text.subtle:
    color: '#465078'
  text.inverse:
    color: '#1a1b26'
  text.overlay:
    color: '#c8d2fa'
  text.sidebar:
    color: '#a9b1d6'

  # Status colors
  status.error:
    color: '#f7768e'
  status.warning:
    color: '#e0af68'
  status.success:
    color: '#9ece6a'
  status.info:
    color: '#7dcfff'
  status.hint:
    color: '#7dcfff'
  status.todo:
    color: '#7aa2f7'

  # Syntax colors
  syntax.keyword:
    color: '#bb9af7'
    bold: true
  syntax.string:
    color: '#9ece6a'
  syntax.function:
    color: '#7aa2f7'
  syntax.comment:
    color: '#565f89'
    italic: true
  syntax.variable:
    color: '#c0caf5'
  syntax.constant:
    color: '#ff9e64'
  syntax.operator:
    color: '#89ddeb'
  syntax.type:
    color: '#e0af68'
  syntax.number:
    color: '#ff9e64'
  syntax.tag:
    color: '#f7768e'
  syntax.property:
    color: '#9ece6a'
  syntax.parameter:
    color: '#e0af68'
  syntax.regexp:
    color: '#7dcfff'
  syntax.escape:
    color: '#bb9af7'
  syntax.constructor:
    color: '#c8acf8'

  # Diff colors
  diff.added.fg:
    color: '#9ece6a'
  diff.added.bg:
    color: '#283c28'
  diff.added.sign:
    color: '#9ece6a'
  diff.deleted.fg:
    color: '#f7768e'
  diff.deleted.bg:
    color: '#3c2828'
  diff.deleted.sign:
    color: '#f7768e'
  diff.changed.fg:
    color: '#89ddeb'
  diff.changed.bg:
    color: '#28323c'
  diff.ignored:
    color: '#565f89'

  # Terminal colors (ANSI 0-15)
  terminal.black:
    color: '#1f2335'
  terminal.red:
    color: '#f7768e'
  terminal.green:
    color: '#9ece6a'
  terminal.yellow:
    color: '#e0af68'
  terminal.blue:
    color: '#7aa2f7'
  terminal.magenta:
    color: '#bb9af7'
  terminal.cyan:
    color: '#7dcfff'
  terminal.white:
    color: '#c0caf5'
  terminal.brblack:
    color: '#565f89'
  terminal.brred:
    color: '#ff899d'
  terminal.brgreen:
    color: '#afd67a'
  terminal.bryellow:
    color: '#e9c582'
  terminal.brblue:
    color: '#8db6fa'
  terminal.brmagenta:
    color: '#c8acf8'
  terminal.brcyan:
    color: '#97d8f8'
  terminal.brwhite:
    color: '#c8d3f5'

  # Accent colors for Get() test
  accent.primary:
    color: '#7aa2f7'
`
	if err := os.WriteFile(filepath.Join(themeDir, "universal.yaml"), []byte(universalYAML), 0o644); err != nil {
		t.Fatalf("failed to write universal.yaml: %v", err)
	}

	// Create style.json for symlink tests
	if err := os.WriteFile(filepath.Join(themeDir, "style.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write style.json: %v", err)
	}
}

func TestIntegration_LoadAndAccessColors(t *testing.T) {
	// Arrange: Create a mock theme directory with complete universal.yaml containing all token paths
	configDir := t.TempDir()
	themeName := "tokyo-night-dark"

	createCompleteTestTheme(t, configDir, themeName)

	// Act: Load theme, call all accessor methods (Surface, Text, Status, Syntax, Diff, Terminal)
	theme, err := flair.LoadNamedFrom(configDir, themeName)
	if err != nil {
		t.Fatalf("LoadNamedFrom() error = %v, want nil", err)
	}

	// Assert: All color fields are populated correctly from the YAML

	// Verify theme metadata
	if theme.Name() != themeName {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), themeName)
	}
	if theme.Variant() != "dark" {
		t.Errorf("theme.Variant() = %q, want %q", theme.Variant(), "dark")
	}
	if !theme.HasColors() {
		t.Error("theme.HasColors() = false, want true")
	}

	// Test Surface colors
	t.Run("Surface", func(t *testing.T) {
		surface := theme.Surface()
		assertColorHex(t, "Surface.Background", surface.Background, "#1a1b26")
		assertColorHex(t, "Surface.Raised", surface.Raised, "#1f2335")
		assertColorHex(t, "Surface.Sunken", surface.Sunken, "#16161e")
		assertColorHex(t, "Surface.Darkest", surface.Darkest, "#101014")
		assertColorHex(t, "Surface.Overlay", surface.Overlay, "#1e1f2a")
		assertColorHex(t, "Surface.Popup", surface.Popup, "#1e1f2a")
		assertColorHex(t, "Surface.Highlight", surface.Highlight, "#292e42")
		assertColorHex(t, "Surface.Selection", surface.Selection, "#32374b")
		assertColorHex(t, "Surface.Search", surface.Search, "#3c4155")
		assertColorHex(t, "Surface.Sidebar", surface.Sidebar, "#1c1d28")
		assertColorHex(t, "Surface.Statusbar", surface.Statusbar, "#191a25")
	})

	// Test Text colors
	t.Run("Text", func(t *testing.T) {
		text := theme.Text()
		assertColorHex(t, "Text.Primary", text.Primary, "#c0caf5")
		assertColorHex(t, "Text.Secondary", text.Secondary, "#a9b1d6")
		assertColorHex(t, "Text.Muted", text.Muted, "#565f89")
		assertColorHex(t, "Text.Subtle", text.Subtle, "#465078")
		assertColorHex(t, "Text.Inverse", text.Inverse, "#1a1b26")
		assertColorHex(t, "Text.Overlay", text.Overlay, "#c8d2fa")
		assertColorHex(t, "Text.Sidebar", text.Sidebar, "#a9b1d6")
	})

	// Test Status colors
	t.Run("Status", func(t *testing.T) {
		status := theme.Status()
		assertColorHex(t, "Status.Error", status.Error, "#f7768e")
		assertColorHex(t, "Status.Warning", status.Warning, "#e0af68")
		assertColorHex(t, "Status.Success", status.Success, "#9ece6a")
		assertColorHex(t, "Status.Info", status.Info, "#7dcfff")
		assertColorHex(t, "Status.Hint", status.Hint, "#7dcfff")
		assertColorHex(t, "Status.Todo", status.Todo, "#7aa2f7")
	})

	// Test Syntax colors
	t.Run("Syntax", func(t *testing.T) {
		syntax := theme.Syntax()
		assertColorHex(t, "Syntax.Keyword", syntax.Keyword, "#bb9af7")
		assertColorHex(t, "Syntax.String", syntax.String, "#9ece6a")
		assertColorHex(t, "Syntax.Function", syntax.Function, "#7aa2f7")
		assertColorHex(t, "Syntax.Comment", syntax.Comment, "#565f89")
		assertColorHex(t, "Syntax.Variable", syntax.Variable, "#c0caf5")
		assertColorHex(t, "Syntax.Constant", syntax.Constant, "#ff9e64")
		assertColorHex(t, "Syntax.Operator", syntax.Operator, "#89ddeb")
		assertColorHex(t, "Syntax.Type", syntax.Type, "#e0af68")
		assertColorHex(t, "Syntax.Number", syntax.Number, "#ff9e64")
		assertColorHex(t, "Syntax.Tag", syntax.Tag, "#f7768e")
		assertColorHex(t, "Syntax.Property", syntax.Property, "#9ece6a")
		assertColorHex(t, "Syntax.Parameter", syntax.Parameter, "#e0af68")
		assertColorHex(t, "Syntax.Regexp", syntax.Regexp, "#7dcfff")
		assertColorHex(t, "Syntax.Escape", syntax.Escape, "#bb9af7")
		assertColorHex(t, "Syntax.Constructor", syntax.Constructor, "#c8acf8")
	})

	// Test Diff colors
	t.Run("Diff", func(t *testing.T) {
		diff := theme.Diff()
		assertColorHex(t, "Diff.AddedFg", diff.AddedFg, "#9ece6a")
		assertColorHex(t, "Diff.AddedBg", diff.AddedBg, "#283c28")
		assertColorHex(t, "Diff.AddedSign", diff.AddedSign, "#9ece6a")
		assertColorHex(t, "Diff.DeletedFg", diff.DeletedFg, "#f7768e")
		assertColorHex(t, "Diff.DeletedBg", diff.DeletedBg, "#3c2828")
		assertColorHex(t, "Diff.DeletedSign", diff.DeletedSign, "#f7768e")
		assertColorHex(t, "Diff.ChangedFg", diff.ChangedFg, "#89ddeb")
		assertColorHex(t, "Diff.ChangedBg", diff.ChangedBg, "#28323c")
		assertColorHex(t, "Diff.Ignored", diff.Ignored, "#565f89")
	})

	// Test Terminal colors
	t.Run("Terminal", func(t *testing.T) {
		terminal := theme.Terminal()
		expectedHex := [16]string{
			"#1f2335", // 0: black
			"#f7768e", // 1: red
			"#9ece6a", // 2: green
			"#e0af68", // 3: yellow
			"#7aa2f7", // 4: blue
			"#bb9af7", // 5: magenta
			"#7dcfff", // 6: cyan
			"#c0caf5", // 7: white
			"#565f89", // 8: bright black
			"#ff899d", // 9: bright red
			"#afd67a", // 10: bright green
			"#e9c582", // 11: bright yellow
			"#8db6fa", // 12: bright blue
			"#c8acf8", // 13: bright magenta
			"#97d8f8", // 14: bright cyan
			"#c8d3f5", // 15: bright white
		}
		for i, expected := range expectedHex {
			if terminal[i].Hex() != expected {
				t.Errorf("Terminal[%d] = %s, want %s", i, terminal[i].Hex(), expected)
			}
		}
	})

	// Test Get() method
	t.Run("Get", func(t *testing.T) {
		color, ok := theme.Get("accent.primary")
		if !ok {
			t.Error("Get(\"accent.primary\") ok = false, want true")
		}
		if color.Hex() != "#7aa2f7" {
			t.Errorf("Get(\"accent.primary\") = %s, want #7aa2f7", color.Hex())
		}

		// Test missing token
		_, ok = theme.Get("nonexistent.token")
		if ok {
			t.Error("Get(\"nonexistent.token\") ok = true, want false")
		}
	})
}

func TestIntegration_LoadWithSymlink(t *testing.T) {
	// Arrange: Create a config dir with theme and symlink selection
	configDir := t.TempDir()
	themeName := "test-theme-light"

	createCompleteTestTheme(t, configDir, themeName)

	// Create symlink for theme selection
	link := filepath.Join(configDir, "style.json")
	target := filepath.Join(themeName, "style.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	// Act: Load via LoadFrom (should follow symlink)
	theme, err := flair.LoadFrom(configDir)

	// Assert
	if err != nil {
		t.Fatalf("LoadFrom() error = %v, want nil", err)
	}
	if theme.Name() != themeName {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), themeName)
	}
	if theme.Variant() != "light" {
		t.Errorf("theme.Variant() = %q, want %q", theme.Variant(), "light")
	}
}

func TestIntegration_ListAndLoadMultipleThemes(t *testing.T) {
	// Arrange: Create multiple themes
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "flair")

	themes := []string{"gruvbox-dark", "tokyo-night-storm", "catppuccin-mocha"}
	for _, name := range themes {
		createCompleteTestTheme(t, configDir, name)
	}

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Act: List themes
	listed, err := flair.ListThemes()
	if err != nil {
		t.Fatalf("ListThemes() error = %v", err)
	}

	// Assert: All themes are listed
	if len(listed) != len(themes) {
		t.Errorf("ListThemes() returned %d themes, want %d", len(listed), len(themes))
	}

	// Verify each theme can be loaded
	for _, name := range themes {
		theme, err := flair.LoadNamed(name)
		if err != nil {
			t.Errorf("LoadNamed(%q) error = %v", name, err)
			continue
		}
		if theme.Name() != name {
			t.Errorf("theme.Name() = %q, want %q", theme.Name(), name)
		}
		if !theme.HasColors() {
			t.Errorf("theme %q has no colors", name)
		}
	}
}

func TestLoad_MissingUniversalYaml(t *testing.T) {
	// Arrange: Theme directory exists but universal.yaml is missing
	configDir := t.TempDir()
	themeName := "incomplete-theme"

	themeDir := filepath.Join(configDir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("failed to create theme dir: %v", err)
	}
	// Don't create universal.yaml

	// Act: Call LoadNamedFrom()
	_, err := flair.LoadNamedFrom(configDir, themeName)

	// Assert: Returns descriptive error about missing file
	if err == nil {
		t.Fatal("LoadNamedFrom() error = nil, want error for missing universal.yaml")
	}

	// Error should wrap ErrThemeNotFound
	if !errors.Is(err, flair.ErrThemeNotFound) {
		t.Errorf("error should wrap ErrThemeNotFound, got: %v", err)
	}
}

func TestLoad_MalformedYaml(t *testing.T) {
	// Arrange: universal.yaml contains invalid YAML syntax
	configDir := t.TempDir()
	themeName := "malformed-theme"

	themeDir := filepath.Join(configDir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("failed to create theme dir: %v", err)
	}

	malformedYAML := `tokens:
  surface.background:
    color: '#1a1b26'
  invalid yaml here: [
    - this is not valid
  text.primary:
    color: '#c0caf5'
`
	if err := os.WriteFile(filepath.Join(themeDir, "universal.yaml"), []byte(malformedYAML), 0o644); err != nil {
		t.Fatalf("failed to write malformed universal.yaml: %v", err)
	}

	// Act: Call LoadNamedFrom()
	_, err := flair.LoadNamedFrom(configDir, themeName)

	// Assert: Returns parse error
	if err == nil {
		t.Fatal("LoadNamedFrom() error = nil, want error for malformed YAML")
	}

	// Error message should mention parsing
	errStr := err.Error()
	if !containsAny(errStr, "parse", "yaml", "unmarshal") {
		t.Errorf("error should mention parsing, got: %v", err)
	}
}

func TestListThemes_EmptyDir(t *testing.T) {
	// Arrange: Config directory exists but contains no themes
	configDir := t.TempDir()

	// Act: Call ListThemesFrom()
	themes, err := flair.ListThemesFrom(configDir)

	// Assert: Returns empty slice (not error)
	if err != nil {
		t.Fatalf("ListThemesFrom() error = %v, want nil", err)
	}

	if len(themes) != 0 {
		t.Errorf("ListThemesFrom() returned %d themes, want 0", len(themes))
	}
}

func TestListThemes_NonexistentDir(t *testing.T) {
	// Arrange: Config directory does not exist
	configDir := filepath.Join(t.TempDir(), "nonexistent")

	// Act: Call ListThemesFrom()
	themes, err := flair.ListThemesFrom(configDir)

	// Assert: Returns empty slice and no error (graceful handling)
	if err != nil {
		t.Fatalf("ListThemesFrom() error = %v, want nil", err)
	}

	if len(themes) != 0 {
		t.Errorf("ListThemesFrom() returned %v, want empty slice", themes)
	}
}

func TestLoad_XDGConfigHome(t *testing.T) {
	// Arrange: Set XDG_CONFIG_HOME to temp dir, create theme there
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "flair")
	themeName := "xdg-test-theme"

	createCompleteTestTheme(t, configDir, themeName)

	// Create symlink for selection
	link := filepath.Join(configDir, "style.json")
	target := filepath.Join(themeName, "style.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Act: Call Load() without explicit path
	theme, err := flair.Load()

	// Assert: Loads theme from XDG location
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if theme.Name() != themeName {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), themeName)
	}
}

func TestLoad_InvalidHexColor(t *testing.T) {
	// Arrange: universal.yaml contains invalid hex color
	configDir := t.TempDir()
	themeName := "invalid-color-theme"

	themeDir := filepath.Join(configDir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("failed to create theme dir: %v", err)
	}

	invalidColorYAML := `tokens:
  surface.background:
    color: '#invalid'
`
	if err := os.WriteFile(filepath.Join(themeDir, "universal.yaml"), []byte(invalidColorYAML), 0o644); err != nil {
		t.Fatalf("failed to write universal.yaml: %v", err)
	}

	// Act: Call LoadNamedFrom()
	_, err := flair.LoadNamedFrom(configDir, themeName)

	// Assert: Returns error about invalid color
	if err == nil {
		t.Fatal("LoadNamedFrom() error = nil, want error for invalid hex color")
	}

	errStr := err.Error()
	if !containsAny(errStr, "invalid", "hex", "color") {
		t.Errorf("error should mention invalid hex color, got: %v", err)
	}
}

// assertColorHex is a test helper that asserts a color matches an expected hex string.
func assertColorHex(t *testing.T, name string, color flair.Color, expectedHex string) {
	t.Helper()
	if color.Hex() != expectedHex {
		t.Errorf("%s = %s, want %s", name, color.Hex(), expectedHex)
	}
}

// containsAny checks if the string contains any of the substrings (case-insensitive).
func containsAny(s string, subs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range subs {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

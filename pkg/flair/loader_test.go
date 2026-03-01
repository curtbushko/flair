package flair_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

// createTestTheme creates a minimal theme directory structure for testing.
// Returns the config directory path.
func createTestTheme(t *testing.T, configDir, themeName string) {
	t.Helper()

	themeDir := filepath.Join(configDir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("failed to create theme dir: %v", err)
	}

	universalYAML := `tokens:
    surface.background:
        color: '#1a1b26'
    text.primary:
        color: '#c0caf5'
    syntax.keyword:
        color: '#bb9af7'
        bold: true
`
	if err := os.WriteFile(filepath.Join(themeDir, "universal.yaml"), []byte(universalYAML), 0o644); err != nil {
		t.Fatalf("failed to write universal.yaml: %v", err)
	}
}

// createStyleSymlink creates a symlink at configDir/style.json pointing to themeName/style.json.
func createStyleSymlink(t *testing.T, configDir, themeName string) {
	t.Helper()

	link := filepath.Join(configDir, "style.json")
	target := filepath.Join(themeName, "style.json")

	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}
}

func TestLoad_SelectedTheme(t *testing.T) {
	// Arrange: Create a temp config dir with a theme and style.json symlink pointing to it.
	configDir := t.TempDir()
	themeName := "tokyo-night"

	createTestTheme(t, configDir, themeName)

	// Create the style.json file in theme dir (even if empty, just needs to exist for symlink)
	themeDir := filepath.Join(configDir, themeName)
	if err := os.WriteFile(filepath.Join(themeDir, "style.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write style.json: %v", err)
	}

	// Create symlink at configDir/style.json -> tokyo-night/style.json
	createStyleSymlink(t, configDir, themeName)

	// Act: Call LoadFrom() with config dir.
	theme, err := flair.LoadFrom(configDir)

	// Assert: Returns Theme with correct name and colors loaded from universal.yaml.
	if err != nil {
		t.Fatalf("LoadFrom() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("LoadFrom() returned nil theme")
	}

	if theme.Name() != themeName {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), themeName)
	}

	// Verify colors were loaded
	bg, ok := theme.Color("surface.background")
	if !ok {
		t.Error("theme.Color(\"surface.background\") not found")
	} else if bg.Hex() != "#1a1b26" {
		t.Errorf("surface.background = %v, want #1a1b26", bg.Hex())
	}

	fg, ok := theme.Color("text.primary")
	if !ok {
		t.Error("theme.Color(\"text.primary\") not found")
	} else if fg.Hex() != "#c0caf5" {
		t.Errorf("text.primary = %v, want #c0caf5", fg.Hex())
	}
}

func TestLoadNamed_ExistingTheme(t *testing.T) {
	// Arrange: Create a temp config dir with theme 'tokyonight' directory containing universal.yaml.
	configDir := t.TempDir()
	themeName := "tokyonight"

	createTestTheme(t, configDir, themeName)

	// Act: Call LoadNamed('tokyonight') with config dir set via env.
	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	// Adjust config path since XDG_CONFIG_HOME/flair is expected
	actualConfigDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, actualConfigDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}
	_ = actualConfigDir // renamed dir is accessed via XDG_CONFIG_HOME

	theme, err := flair.LoadNamed(themeName)

	// Assert: Returns Theme with name 'tokyonight' and colors from universal.yaml.
	if err != nil {
		t.Fatalf("LoadNamed() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("LoadNamed() returned nil theme")
	}

	if theme.Name() != themeName {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), themeName)
	}

	// Verify at least one color was loaded
	if !theme.HasColors() {
		t.Error("theme has no colors")
	}
}

func TestLoadNamed_NotFound(t *testing.T) {
	// Arrange: Create an empty temp config dir.
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "flair")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Act: Call LoadNamed('nonexistent').
	_, err := flair.LoadNamed("nonexistent")

	// Assert: Returns error indicating theme not found.
	if err == nil {
		t.Fatal("LoadNamed() error = nil, want error for nonexistent theme")
	}
}

func TestListThemes_MultipleThemes(t *testing.T) {
	// Arrange: Create config dir with 'tokyo-night-dark' and 'gruvbox' theme directories.
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "flair")

	createTestTheme(t, configDir, "tokyo-night-dark")
	createTestTheme(t, configDir, "gruvbox")

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Act: Call ListThemes().
	themes, err := flair.ListThemes()

	// Assert: Returns slice containing both theme names.
	if err != nil {
		t.Fatalf("ListThemes() error = %v, want nil", err)
	}

	if len(themes) != 2 {
		t.Errorf("ListThemes() returned %d themes, want 2", len(themes))
	}

	// Check that both themes are present (order may vary but should be sorted).
	found := make(map[string]bool)
	for _, th := range themes {
		found[th] = true
	}

	if !found["tokyo-night-dark"] {
		t.Error("ListThemes() missing 'tokyo-night-dark'")
	}
	if !found["gruvbox"] {
		t.Error("ListThemes() missing 'gruvbox'")
	}
}

func TestSelectedTheme_FollowsSymlink(t *testing.T) {
	// Arrange: Create config dir with style.json symlink pointing to tokyonight/style.json.
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "flair")
	themeName := "tokyonight"

	createTestTheme(t, configDir, themeName)

	// Create the style.json file in theme dir
	themeDir := filepath.Join(configDir, themeName)
	if err := os.WriteFile(filepath.Join(themeDir, "style.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write style.json: %v", err)
	}

	// Create symlink
	createStyleSymlink(t, configDir, themeName)

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Act: Call SelectedTheme().
	selected, err := flair.SelectedTheme()

	// Assert: Returns 'tokyonight'.
	if err != nil {
		t.Fatalf("SelectedTheme() error = %v, want nil", err)
	}

	if selected != themeName {
		t.Errorf("SelectedTheme() = %q, want %q", selected, themeName)
	}
}

func TestLoad_UsesXDGConfigHome(t *testing.T) {
	// Arrange: Set XDG_CONFIG_HOME and create theme there.
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "flair")
	themeName := "test-theme"

	createTestTheme(t, configDir, themeName)

	// Create symlink for selection
	themeDir := filepath.Join(configDir, themeName)
	if err := os.WriteFile(filepath.Join(themeDir, "style.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write style.json: %v", err)
	}
	createStyleSymlink(t, configDir, themeName)

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Act: Call Load() (should use XDG_CONFIG_HOME).
	theme, err := flair.Load()

	// Assert: Theme loaded from XDG_CONFIG_HOME location.
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if theme.Name() != themeName {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), themeName)
	}
}

func TestLoadFrom_NoSelectedTheme(t *testing.T) {
	// Arrange: Create a config dir with a theme but no symlink.
	configDir := t.TempDir()
	createTestTheme(t, configDir, "some-theme")

	// Act: Call LoadFrom() without a symlink.
	_, err := flair.LoadFrom(configDir)

	// Assert: Returns an error indicating no theme is selected.
	if err == nil {
		t.Fatal("LoadFrom() error = nil, want error for no selected theme")
	}
}

func TestListThemesFrom(t *testing.T) {
	// Arrange: Create config dir with themes.
	configDir := t.TempDir()
	createTestTheme(t, configDir, "alpha")
	createTestTheme(t, configDir, "beta")

	// Act: Call ListThemesFrom().
	themes, err := flair.ListThemesFrom(configDir)

	// Assert: Returns both themes, sorted.
	if err != nil {
		t.Fatalf("ListThemesFrom() error = %v", err)
	}

	if len(themes) != 2 {
		t.Fatalf("ListThemesFrom() returned %d themes, want 2", len(themes))
	}

	// Should be sorted alphabetically.
	if themes[0] != "alpha" || themes[1] != "beta" {
		t.Errorf("ListThemesFrom() = %v, want [alpha, beta]", themes)
	}
}

func TestLoadNamedFrom(t *testing.T) {
	// Arrange: Create a theme in a custom config dir.
	configDir := t.TempDir()
	themeName := "custom-theme"
	createTestTheme(t, configDir, themeName)

	// Act: Call LoadNamedFrom().
	theme, err := flair.LoadNamedFrom(configDir, themeName)

	// Assert: Theme loaded successfully.
	if err != nil {
		t.Fatalf("LoadNamedFrom() error = %v", err)
	}

	if theme.Name() != themeName {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), themeName)
	}
}

func TestSelectedThemeFrom(t *testing.T) {
	// Arrange: Create config dir with a symlinked theme.
	configDir := t.TempDir()
	themeName := "selected-theme"
	createTestTheme(t, configDir, themeName)

	themeDir := filepath.Join(configDir, themeName)
	if err := os.WriteFile(filepath.Join(themeDir, "style.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write style.json: %v", err)
	}
	createStyleSymlink(t, configDir, themeName)

	// Act: Call SelectedThemeFrom().
	selected, err := flair.SelectedThemeFrom(configDir)

	// Assert: Returns correct theme name.
	if err != nil {
		t.Fatalf("SelectedThemeFrom() error = %v", err)
	}

	if selected != themeName {
		t.Errorf("SelectedThemeFrom() = %q, want %q", selected, themeName)
	}
}

func TestSelectedThemeFrom_NoSymlink(t *testing.T) {
	// Arrange: Create config dir with no symlinks.
	configDir := t.TempDir()

	// Act: Call SelectedThemeFrom().
	selected, err := flair.SelectedThemeFrom(configDir)

	// Assert: Returns empty string and no error.
	if err != nil {
		t.Fatalf("SelectedThemeFrom() error = %v", err)
	}

	if selected != "" {
		t.Errorf("SelectedThemeFrom() = %q, want empty string", selected)
	}
}

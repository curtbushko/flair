package flair_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

const testThemeTokyoNight = "tokyo-night-dark"

func TestNewStoreAt(t *testing.T) {
	// Arrange: Create a temp directory.
	tempDir := t.TempDir()

	// Act: Create a store at that directory.
	store := flair.NewStoreAt(tempDir)

	// Assert: Store has the correct config dir.
	if store.ConfigDir() != tempDir {
		t.Errorf("store.ConfigDir() = %q, want %q", store.ConfigDir(), tempDir)
	}
}

func TestNewStore_DefaultDir(t *testing.T) {
	// Arrange: Set XDG_CONFIG_HOME to a temp directory.
	tempDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tempDir)

	// Act: Create a store with default config.
	store := flair.NewStore()

	// Assert: Store config dir is XDG_CONFIG_HOME/flair.
	expected := filepath.Join(tempDir, "flair")
	if store.ConfigDir() != expected {
		t.Errorf("store.ConfigDir() = %q, want %q", store.ConfigDir(), expected)
	}
}

func TestStore_Install(t *testing.T) {
	// Arrange: Create an empty temp config dir.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	// Act: Install a built-in theme.
	err := store.Install(testThemeTokyoNight)

	// Assert: Theme files written to configDir/tokyo-night-dark/.
	if err != nil {
		t.Fatalf("store.Install() error = %v, want nil", err)
	}

	themeDir := filepath.Join(configDir, testThemeTokyoNight)
	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		t.Fatalf("theme directory %q not created", themeDir)
	}

	// Check tokens.yaml exists.
	tokensPath := filepath.Join(themeDir, "tokens.yaml")
	if _, err := os.Stat(tokensPath); os.IsNotExist(err) {
		t.Errorf("tokens.yaml not created at %q", tokensPath)
	}

	// Check output files exist.
	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		path := filepath.Join(themeDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("output file %q not created", f)
		}
	}
}

func TestStore_Install_NotFound(t *testing.T) {
	// Arrange: Create an empty temp config dir.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	// Act: Try to install a non-existent theme.
	err := store.Install("nonexistent-theme-that-does-not-exist")

	// Assert: Returns error.
	if err == nil {
		t.Fatal("store.Install() error = nil, want error for nonexistent theme")
	}
}

func TestStore_InstallAll(t *testing.T) {
	// Arrange: Create an empty temp config dir.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	// Act: Install all built-in themes.
	err := store.InstallAll()

	// Assert: All 189 themes installed.
	if err != nil {
		t.Fatalf("store.InstallAll() error = %v, want nil", err)
	}

	// Count installed themes by counting directories with tokens.yaml.
	entries, err := os.ReadDir(configDir)
	if err != nil {
		t.Fatalf("failed to read config dir: %v", err)
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			tokensPath := filepath.Join(configDir, entry.Name(), "tokens.yaml")
			if _, err := os.Stat(tokensPath); err == nil {
				count++
			}
		}
	}

	if count != 189 {
		t.Errorf("InstallAll() installed %d themes, want 189", count)
	}
}

func TestStore_Select(t *testing.T) {
	// Arrange: Install a theme first.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install(testThemeTokyoNight); err != nil {
		t.Fatalf("failed to install theme: %v", err)
	}

	// Act: Select the theme.
	err := store.Select(testThemeTokyoNight)

	// Assert: Symlinks created at configDir/style.lua, style.css, etc.
	if err != nil {
		t.Fatalf("store.Select() error = %v, want nil", err)
	}

	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		link := filepath.Join(configDir, f)
		target, err := os.Readlink(link)
		if err != nil {
			t.Errorf("symlink %q not created: %v", f, err)
			continue
		}

		expectedTarget := filepath.Join(testThemeTokyoNight, f)
		if target != expectedTarget {
			t.Errorf("symlink %q points to %q, want %q", f, target, expectedTarget)
		}
	}
}

func TestStore_Select_NotInstalled(t *testing.T) {
	// Arrange: Create an empty config dir (no themes installed).
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	// Act: Try to select a theme that isn't installed.
	err := store.Select("nonexistent-theme")

	// Assert: Returns error.
	if err == nil {
		t.Fatal("store.Select() error = nil, want error for non-installed theme")
	}
}

func TestStore_Load(t *testing.T) {
	// Arrange: Install and select a theme.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install(testThemeTokyoNight); err != nil {
		t.Fatalf("failed to install theme: %v", err)
	}
	if err := store.Select(testThemeTokyoNight); err != nil {
		t.Fatalf("failed to select theme: %v", err)
	}

	// Act: Load the selected theme.
	theme, err := store.Load()

	// Assert: Returns *Theme with correct name and colors.
	if err != nil {
		t.Fatalf("store.Load() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("store.Load() returned nil theme")
	}

	if theme.Name() != testThemeTokyoNight {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), testThemeTokyoNight)
	}

	// Verify colors were loaded.
	if !theme.HasColors() {
		t.Error("theme has no colors")
	}

	// Check for specific semantic tokens.
	expectedTokens := []string{
		"surface.background",
		"text.primary",
		"syntax.keyword",
	}
	for _, token := range expectedTokens {
		if _, ok := theme.Color(token); !ok {
			t.Errorf("theme missing expected token %q", token)
		}
	}
}

func TestStore_Load_NoSelection(t *testing.T) {
	// Arrange: Create an empty config dir (no theme selected).
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	// Act: Try to load when nothing is selected.
	_, err := store.Load()

	// Assert: Returns error.
	if err == nil {
		t.Fatal("store.Load() error = nil, want error when no theme selected")
	}
}

func TestStore_List(t *testing.T) {
	// Arrange: Install two themes.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install(testThemeTokyoNight); err != nil {
		t.Fatalf("failed to install tokyo-night-dark: %v", err)
	}
	if err := store.Install("gruvbox-dark"); err != nil {
		t.Fatalf("failed to install gruvbox-dark: %v", err)
	}

	// Act: List installed themes.
	themes, err := store.List()

	// Assert: Returns sorted slice of theme names.
	if err != nil {
		t.Fatalf("store.List() error = %v, want nil", err)
	}

	if len(themes) != 2 {
		t.Fatalf("store.List() returned %d themes, want 2", len(themes))
	}

	// Should be sorted.
	if !sort.StringsAreSorted(themes) {
		t.Errorf("store.List() returned unsorted list: %v", themes)
	}

	// Check both themes are present.
	found := make(map[string]bool)
	for _, th := range themes {
		found[th] = true
	}

	if !found[testThemeTokyoNight] {
		t.Error("store.List() missing 'tokyo-night-dark'")
	}
	if !found["gruvbox-dark"] {
		t.Error("store.List() missing 'gruvbox-dark'")
	}
}

func TestStore_List_Empty(t *testing.T) {
	// Arrange: Create an empty config dir.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	// Act: List themes in empty dir.
	themes, err := store.List()

	// Assert: Returns empty slice, no error.
	if err != nil {
		t.Fatalf("store.List() error = %v, want nil", err)
	}

	if len(themes) != 0 {
		t.Errorf("store.List() returned %d themes, want 0", len(themes))
	}
}

func TestStore_Selected(t *testing.T) {
	// Arrange: Install and select a theme.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install(testThemeTokyoNight); err != nil {
		t.Fatalf("failed to install theme: %v", err)
	}
	if err := store.Select(testThemeTokyoNight); err != nil {
		t.Fatalf("failed to select theme: %v", err)
	}

	// Act: Get selected theme name.
	selected, err := store.Selected()

	// Assert: Returns the selected theme name.
	if err != nil {
		t.Fatalf("store.Selected() error = %v, want nil", err)
	}

	if selected != testThemeTokyoNight {
		t.Errorf("store.Selected() = %q, want %q", selected, testThemeTokyoNight)
	}
}

func TestStore_Selected_NoSelection(t *testing.T) {
	// Arrange: Create an empty config dir.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	// Act: Get selected theme when none is selected.
	selected, err := store.Selected()

	// Assert: Returns empty string, no error.
	if err != nil {
		t.Fatalf("store.Selected() error = %v, want nil", err)
	}

	if selected != "" {
		t.Errorf("store.Selected() = %q, want empty string", selected)
	}
}

func TestStore_LoadNamed(t *testing.T) {
	// Arrange: Install a theme.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install("gruvbox-dark"); err != nil {
		t.Fatalf("failed to install theme: %v", err)
	}

	// Act: Load a specific theme by name.
	theme, err := store.LoadNamed("gruvbox-dark")

	// Assert: Returns the theme.
	if err != nil {
		t.Fatalf("store.LoadNamed() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("store.LoadNamed() returned nil theme")
	}

	if theme.Name() != "gruvbox-dark" {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), "gruvbox-dark")
	}
}

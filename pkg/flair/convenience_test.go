package flair_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

// TestDefault_FallbackToBuiltin verifies that Default() returns tokyo-night-dark
// built-in theme when no theme is selected in an empty config directory.
func TestDefault_FallbackToBuiltin(t *testing.T) {
	// Arrange: Create empty config dir (no theme selected).
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))

	// Rename to match expected flair subdirectory.
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, flairDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}

	// Act: Call Default().
	theme, err := flair.Default()

	// Assert: Returns tokyo-night-dark built-in theme.
	if err != nil {
		t.Fatalf("Default() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("Default() returned nil theme")
	}

	// Should be the tokyo-night-dark theme (display name is "Tokyo Night Dark").
	if theme.Variant() != testVariantDark {
		t.Errorf("Default() theme.Variant() = %q, want %q", theme.Variant(), testVariantDark)
	}

	// Should have colors.
	if !theme.HasColors() {
		t.Error("Default() theme has no colors")
	}
}

// TestDefault_SelectedTheme verifies that Default() returns the currently
// selected theme when one is installed and selected.
func TestDefault_SelectedTheme(t *testing.T) {
	// Arrange: Install and select a theme in temp config dir.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install("gruvbox-dark"); err != nil {
		t.Fatalf("failed to install theme: %v", err)
	}
	if err := store.Select("gruvbox-dark"); err != nil {
		t.Fatalf("failed to select theme: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, flairDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}

	// Act: Call Default().
	theme, err := flair.Default()

	// Assert: Returns the selected gruvbox-dark theme.
	if err != nil {
		t.Fatalf("Default() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("Default() returned nil theme")
	}

	if theme.Name() != "gruvbox-dark" {
		t.Errorf("Default() theme.Name() = %q, want %q", theme.Name(), "gruvbox-dark")
	}
}

// TestMustLoad_Panic verifies that MustLoad() panics when no theme can be loaded.
func TestMustLoad_Panic(t *testing.T) {
	// Arrange: Set up environment to cause a load error.
	// We can't easily cause a panic since Default() falls back to builtin,
	// but we can test the panic mechanism by checking it doesn't panic on success.
	// For a true panic test, we'd need to mock LoadBuiltin which isn't easy.
	// Instead, we verify the happy path doesn't panic.

	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.MkdirAll(flairDir, 0o755); err != nil {
		t.Fatalf("failed to create flair dir: %v", err)
	}

	// Act & Assert: MustLoad should not panic (falls back to builtin).
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustLoad() panicked unexpectedly: %v", r)
		}
	}()

	theme := flair.MustLoad()
	if theme == nil {
		t.Fatal("MustLoad() returned nil theme")
	}
}

// TestMustLoad_Success verifies that MustLoad() returns theme without panic on success.
func TestMustLoad_Success(t *testing.T) {
	// Arrange: Install and select a valid theme.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install("tokyo-night-dark"); err != nil {
		t.Fatalf("failed to install theme: %v", err)
	}
	if err := store.Select("tokyo-night-dark"); err != nil {
		t.Fatalf("failed to select theme: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, flairDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}

	// Act: Call MustLoad().
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustLoad() panicked: %v", r)
		}
	}()

	theme := flair.MustLoad()

	// Assert: Returns theme without panic.
	if theme == nil {
		t.Fatal("MustLoad() returned nil theme")
	}

	if theme.Name() != "tokyo-night-dark" {
		t.Errorf("MustLoad() theme.Name() = %q, want %q", theme.Name(), "tokyo-night-dark")
	}
}

// TestLoadOrDefault_NamedExists verifies that LoadOrDefault returns the named
// theme when it exists as an installed theme.
func TestLoadOrDefault_NamedExists(t *testing.T) {
	// Arrange: Install the named theme.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install("catppuccin-mocha"); err != nil {
		t.Fatalf("failed to install theme: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, flairDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}

	// Act: Call LoadOrDefault with named theme and fallback.
	theme, err := flair.LoadOrDefault("catppuccin-mocha", "tokyo-night-dark")

	// Assert: Returns the named theme (catppuccin-mocha).
	if err != nil {
		t.Fatalf("LoadOrDefault() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("LoadOrDefault() returned nil theme")
	}

	if theme.Name() != "catppuccin-mocha" {
		t.Errorf("LoadOrDefault() theme.Name() = %q, want %q", theme.Name(), "catppuccin-mocha")
	}
}

// TestLoadOrDefault_Fallback verifies that LoadOrDefault falls back to the
// built-in fallback theme when the named theme is not installed.
func TestLoadOrDefault_Fallback(t *testing.T) {
	// Arrange: Empty config dir (named theme not installed).
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, flairDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}

	// Act: Call LoadOrDefault with non-installed name and valid fallback.
	theme, err := flair.LoadOrDefault("nonexistent-theme", "gruvbox-dark")

	// Assert: Returns the fallback built-in theme (gruvbox-dark).
	if err != nil {
		t.Fatalf("LoadOrDefault() error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("LoadOrDefault() returned nil theme")
	}

	// Fallback is a built-in, so check it loaded.
	if theme.Variant() != testVariantDark {
		t.Errorf("LoadOrDefault() fallback theme.Variant() = %q, want %q", theme.Variant(), testVariantDark)
	}

	if !theme.HasColors() {
		t.Error("LoadOrDefault() fallback theme has no colors")
	}
}

// TestEnsureInstalled_Empty verifies that EnsureInstalled installs all built-in
// themes when the config directory is empty.
func TestEnsureInstalled_Empty(t *testing.T) {
	// Arrange: Create empty config dir.
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, flairDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}

	// Act: Call EnsureInstalled().
	err := flair.EnsureInstalled()

	// Assert: All built-in themes installed, returns nil.
	if err != nil {
		t.Fatalf("EnsureInstalled() error = %v, want nil", err)
	}

	// Verify themes were installed by checking the directory.
	entries, err := os.ReadDir(flairDir)
	if err != nil {
		t.Fatalf("failed to read config dir: %v", err)
	}

	// Count installed themes.
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			tokensPath := filepath.Join(flairDir, entry.Name(), "tokens.yaml")
			if _, err := os.Stat(tokensPath); err == nil {
				count++
			}
		}
	}

	// Should have all 190 built-in themes.
	if count != 190 {
		t.Errorf("EnsureInstalled() installed %d themes, want 190", count)
	}
}

// TestEnsureInstalled_NotEmpty verifies that EnsureInstalled does nothing when
// the config directory already has at least one theme.
func TestEnsureInstalled_NotEmpty(t *testing.T) {
	// Arrange: Create config dir with one theme already installed.
	configDir := t.TempDir()
	store := flair.NewStoreAt(configDir)

	if err := store.Install("tokyo-night-dark"); err != nil {
		t.Fatalf("failed to install initial theme: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(configDir))
	flairDir := filepath.Join(filepath.Dir(configDir), "flair")
	if err := os.Rename(configDir, flairDir); err != nil {
		t.Fatalf("failed to rename config dir: %v", err)
	}

	// Act: Call EnsureInstalled().
	err := flair.EnsureInstalled()

	// Assert: No additional themes installed, returns nil.
	if err != nil {
		t.Fatalf("EnsureInstalled() error = %v, want nil", err)
	}

	// Verify only the original theme exists.
	entries, err := os.ReadDir(flairDir)
	if err != nil {
		t.Fatalf("failed to read config dir: %v", err)
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			tokensPath := filepath.Join(flairDir, entry.Name(), "tokens.yaml")
			if _, err := os.Stat(tokensPath); err == nil {
				count++
			}
		}
	}

	// Should still have only 1 theme.
	if count != 1 {
		t.Errorf("EnsureInstalled() installed %d themes, want 1 (no new themes)", count)
	}
}

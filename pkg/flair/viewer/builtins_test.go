package viewer

import (
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

// TestBuiltinThemeLoader_LoadPalette verifies BuiltinThemeLoader loads palette from embedded themes.
func TestBuiltinThemeLoader_LoadPalette(t *testing.T) {
	loader := NewBuiltinThemeLoader()

	// Load a known built-in theme.
	palette, err := loader.LoadPalette("tokyo-night-dark")
	if err != nil {
		t.Fatalf("LoadPalette failed: %v", err)
	}

	// Verify 24 colors are present (non-empty).
	for i, c := range palette.Colors {
		if c == "" {
			t.Errorf("palette.Colors[%d] is empty, expected hex color", i)
		}
	}

	// Verify at least base00 is a valid hex color (starts with #).
	if len(palette.Colors[0]) < 4 || palette.Colors[0][0] != '#' {
		t.Errorf("palette.Colors[0] = %q, want valid hex color", palette.Colors[0])
	}
}

// TestBuiltinThemeLoader_LoadPalette_NotFound verifies error for non-existent theme.
func TestBuiltinThemeLoader_LoadPalette_NotFound(t *testing.T) {
	loader := NewBuiltinThemeLoader()

	_, err := loader.LoadPalette("nonexistent-theme")
	if err == nil {
		t.Error("LoadPalette should return error for non-existent theme")
	}
}

// TestBuiltinThemeLoader_LoadTokens verifies BuiltinThemeLoader loads tokens from embedded themes.
func TestBuiltinThemeLoader_LoadTokens(t *testing.T) {
	loader := NewBuiltinThemeLoader()

	tokens, err := loader.LoadTokens("tokyo-night-dark")
	if err != nil {
		t.Fatalf("LoadTokens failed: %v", err)
	}

	// Verify all token groups are populated.
	if len(tokens.Surface) == 0 {
		t.Error("tokens.Surface is empty")
	}
	if len(tokens.Text) == 0 {
		t.Error("tokens.Text is empty")
	}
	if len(tokens.Status) == 0 {
		t.Error("tokens.Status is empty")
	}
	if len(tokens.Syntax) == 0 {
		t.Error("tokens.Syntax is empty")
	}
	if len(tokens.Diff) == 0 {
		t.Error("tokens.Diff is empty")
	}

	// Verify specific token exists.
	if _, ok := tokens.Surface["surface.background"]; !ok {
		t.Error("expected surface.background token")
	}
	if _, ok := tokens.Text["text.primary"]; !ok {
		t.Error("expected text.primary token")
	}
}

// TestBuiltinThemeLoader_LoadTokens_NotFound verifies error for non-existent theme.
func TestBuiltinThemeLoader_LoadTokens_NotFound(t *testing.T) {
	loader := NewBuiltinThemeLoader()

	_, err := loader.LoadTokens("nonexistent-theme")
	if err == nil {
		t.Error("LoadTokens should return error for non-existent theme")
	}
}

// TestBuiltinThemeLoader_ImplementsThemeLoader verifies interface compliance.
func TestBuiltinThemeLoader_ImplementsThemeLoader(t *testing.T) {
	var _ ThemeLoader = &BuiltinThemeLoader{}
}

// TestRunBuiltins_NoConfigDir verifies RunBuiltins works without filesystem config.
func TestRunBuiltins_NoConfigDir(t *testing.T) {
	// Set empty config home to ensure no filesystem dependency.
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	// RunBuiltins with DryRun should not error.
	err := RunBuiltins(RunBuiltinsOptions{
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("RunBuiltins failed: %v", err)
	}
}

// TestRunBuiltins_UsesListBuiltins verifies RunBuiltins uses flair.ListBuiltins.
func TestRunBuiltins_UsesListBuiltins(t *testing.T) {
	// Get the expected themes from ListBuiltins.
	expected := flair.ListBuiltins()
	if len(expected) == 0 {
		t.Skip("no built-in themes available")
	}

	// RunBuiltins should use these themes.
	// We verify indirectly by checking it doesn't error with DryRun.
	err := RunBuiltins(RunBuiltinsOptions{
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("RunBuiltins failed: %v", err)
	}
}

// TestOnInstall_Callback verifies OnInstall is called when user confirms selection.
func TestOnInstall_Callback(t *testing.T) {
	var installedTheme string
	loader := NewBuiltinThemeLoader()

	m := NewModel(Options{
		Themes:      flair.ListBuiltins()[:3], // Use first 3 themes for test
		ThemeLoader: loader,
		OnInstall: func(name string) error {
			installedTheme = name
			return nil
		},
	})

	// Navigate to second theme and press Enter.
	m.cursor = 1
	m.confirmSelection()

	// OnInstall should have been called with the theme name.
	if installedTheme == "" {
		t.Error("OnInstall was not called")
	}
	if installedTheme != m.themes[1] {
		t.Errorf("OnInstall called with %q, want %q", installedTheme, m.themes[1])
	}
}

// TestOnInstall_NotCalledWithoutOption verifies OnInstall not called when not set.
func TestOnInstall_NotCalledWithoutOption(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1", "theme2"},
	})

	// Should not panic when OnInstall is nil.
	m.cursor = 0
	m.confirmSelection()

	// Verify selection was made.
	if m.selectedTheme != "theme1" {
		t.Errorf("selectedTheme = %q, want theme1", m.selectedTheme)
	}
}

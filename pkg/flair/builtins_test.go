package flair

import (
	"os"
	"sort"
	"testing"
)

func TestListBuiltins(t *testing.T) {
	names := ListBuiltins()

	// Should return non-empty list
	if len(names) == 0 {
		t.Fatal("ListBuiltins() returned empty list, expected at least one palette")
	}

	// Should be sorted
	if !sort.StringsAreSorted(names) {
		t.Errorf("ListBuiltins() returned unsorted list: %v", names)
	}

	// Should contain known palettes
	known := map[string]bool{
		"tokyo-night-dark": false,
		"gruvbox-dark":     false,
		"catppuccin-mocha": false,
	}

	for _, name := range names {
		if _, ok := known[name]; ok {
			known[name] = true
		}
	}

	for name, found := range known {
		if !found {
			t.Errorf("ListBuiltins() missing expected palette %q", name)
		}
	}
}

func TestLoadBuiltin_Valid(t *testing.T) {
	theme, err := LoadBuiltin("tokyo-night-dark")
	if err != nil {
		t.Fatalf("LoadBuiltin(\"tokyo-night-dark\") error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("LoadBuiltin(\"tokyo-night-dark\") returned nil theme")
	}

	// Check theme name
	if theme.Name() != "Tokyo Night Dark" {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), "Tokyo Night Dark")
	}

	// Check variant
	if theme.Variant() != "dark" {
		t.Errorf("theme.Variant() = %q, want %q", theme.Variant(), "dark")
	}

	// Check that theme has colors (semantic tokens)
	if !theme.HasColors() {
		t.Error("theme.HasColors() = false, want true")
	}

	// Check for specific semantic tokens
	expectedTokens := []string{
		"surface.background",
		"text.primary",
		"syntax.keyword",
		"status.error",
		"accent.primary",
	}

	for _, token := range expectedTokens {
		if _, ok := theme.Color(token); !ok {
			t.Errorf("theme missing expected token %q", token)
		}
	}
}

func TestLoadBuiltin_NotFound(t *testing.T) {
	_, err := LoadBuiltin("nonexistent-palette-that-does-not-exist")
	if err == nil {
		t.Fatal("LoadBuiltin(\"nonexistent-palette-that-does-not-exist\") error = nil, want error")
	}
}

func TestHasBuiltin(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"tokyo-night-dark", true},
		{"gruvbox-dark", true},
		{"catppuccin-mocha", true},
		{"nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasBuiltin(tt.name)
			if got != tt.want {
				t.Errorf("HasBuiltin(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestLoadBuiltin_NoFilesystem(t *testing.T) {
	// Set empty XDG_CONFIG_HOME to simulate no config directory
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	defer func() {
		if originalXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	// Should still work because built-ins are embedded
	theme, err := LoadBuiltin("tokyo-night-dark")
	if err != nil {
		t.Fatalf("LoadBuiltin() with empty XDG_CONFIG_HOME error = %v, want nil", err)
	}

	if theme == nil {
		t.Fatal("LoadBuiltin() with empty XDG_CONFIG_HOME returned nil theme")
	}

	if theme.Name() != "Tokyo Night Dark" {
		t.Errorf("theme.Name() = %q, want %q", theme.Name(), "Tokyo Night Dark")
	}
}

func TestListBuiltins_Count(t *testing.T) {
	names := ListBuiltins()

	// We should have exactly 189 built-in palettes
	if len(names) != 189 {
		t.Errorf("ListBuiltins() returned %d palettes, want 189", len(names))
	}
}

package application_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/application"
)

func TestListThemesUseCase_ReturnsInstalledThemes(t *testing.T) {
	store := newStubThemeStore()

	// Set up two themes with palette.yaml and all output files.
	for _, theme := range []string{"catppuccin", "dracula"} {
		if err := store.EnsureThemeDir(theme); err != nil {
			t.Fatal(err)
		}
		w, err := store.OpenWriter(theme, "palette.yaml")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("content"))
		_ = w.Close()
		for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
			w, err := store.OpenWriter(theme, f)
			if err != nil {
				t.Fatal(err)
			}
			_, _ = w.Write([]byte("content"))
			_ = w.Close()
		}
	}

	builtins := newStubPaletteSource()

	uc := application.NewListThemesUseCase(store, builtins)
	themes, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(themes) != 2 {
		t.Fatalf("expected 2 themes, got %d", len(themes))
	}

	// Verify theme names are present and marked as Generated.
	nameMap := make(map[string]application.ThemeInfo)
	for _, info := range themes {
		nameMap[info.Name] = info
	}
	if _, ok := nameMap["catppuccin"]; !ok {
		t.Error("expected catppuccin in list")
	}
	if _, ok := nameMap["dracula"]; !ok {
		t.Error("expected dracula in list")
	}
	for _, info := range themes {
		if !info.Generated {
			t.Errorf("theme %q should be marked Generated", info.Name)
		}
	}
}

func TestListThemesUseCase_MarksSelected(t *testing.T) {
	store := newStubThemeStore()

	// Set up two themes with palette.yaml and all output files.
	for _, theme := range []string{"alpha", "beta"} {
		if err := store.EnsureThemeDir(theme); err != nil {
			t.Fatal(err)
		}
		w, err := store.OpenWriter(theme, "palette.yaml")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("content"))
		_ = w.Close()
		for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
			w, err := store.OpenWriter(theme, f)
			if err != nil {
				t.Fatal(err)
			}
			_, _ = w.Write([]byte("content"))
			_ = w.Close()
		}
	}

	// Select "beta" as active theme.
	if err := store.Select("beta"); err != nil {
		t.Fatal(err)
	}

	builtins := newStubPaletteSource()

	uc := application.NewListThemesUseCase(store, builtins)
	themes, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var selectedCount int
	for _, info := range themes {
		if info.Name == "beta" && info.Selected {
			selectedCount++
		}
		if info.Name == "alpha" && info.Selected {
			t.Error("alpha should not be selected")
		}
	}
	if selectedCount != 1 {
		t.Errorf("expected exactly 1 selected theme (beta), got %d", selectedCount)
	}
}

func TestListThemesUseCase_EmptyDirShowsBuiltins(t *testing.T) {
	store := newStubThemeStore()
	builtins := newStubPaletteSource()
	builtins.palettes["catppuccin-mocha"] = []byte("yaml")
	builtins.palettes["gruvbox-dark"] = []byte("yaml")

	uc := application.NewListThemesUseCase(store, builtins)
	themes, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should show built-in palettes even when no themes are generated.
	if len(themes) != 2 {
		t.Fatalf("expected 2 themes (builtins), got %d", len(themes))
	}

	for _, info := range themes {
		if info.Generated {
			t.Errorf("theme %q should not be marked Generated", info.Name)
		}
	}
}

func TestListThemesUseCase_FiltersNonThemeDirs(t *testing.T) {
	store := newStubThemeStore()

	// Create non-theme directories (no palette.yaml or output files).
	if err := store.EnsureThemeDir("generated"); err != nil {
		t.Fatal(err)
	}
	if err := store.EnsureThemeDir("styles"); err != nil {
		t.Fatal(err)
	}

	builtins := newStubPaletteSource()

	uc := application.NewListThemesUseCase(store, builtins)
	themes, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Non-theme directories should be filtered out.
	for _, info := range themes {
		if info.Name == "generated" || info.Name == "styles" {
			t.Errorf("non-theme directory %q should not appear in listing", info.Name)
		}
	}
}

func TestListThemesUseCase_MergesBuiltinsWithGenerated(t *testing.T) {
	store := newStubThemeStore()

	// Generate one theme that matches a built-in name.
	if err := store.EnsureThemeDir("tokyo-night-dark"); err != nil {
		t.Fatal(err)
	}
	writeStubFile(t, store, "tokyo-night-dark", "palette.yaml")
	for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
		writeStubFile(t, store, "tokyo-night-dark", f)
	}

	builtins := newStubPaletteSource()
	builtins.palettes["tokyo-night-dark"] = []byte("yaml")
	builtins.palettes["catppuccin-mocha"] = []byte("yaml")
	builtins.palettes["gruvbox-dark"] = []byte("yaml")

	uc := application.NewListThemesUseCase(store, builtins)
	themes, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should show 3 themes: 1 generated + 2 available builtins.
	if len(themes) != 3 {
		t.Fatalf("expected 3 themes, got %d", len(themes))
	}

	nameMap := make(map[string]application.ThemeInfo)
	for _, info := range themes {
		nameMap[info.Name] = info
	}

	// tokyo-night-dark should be Generated.
	if info, ok := nameMap["tokyo-night-dark"]; !ok {
		t.Error("expected tokyo-night-dark in list")
	} else if !info.Generated {
		t.Error("tokyo-night-dark should be marked Generated")
	}

	// catppuccin-mocha should NOT be Generated.
	if info, ok := nameMap["catppuccin-mocha"]; !ok {
		t.Error("expected catppuccin-mocha in list")
	} else if info.Generated {
		t.Error("catppuccin-mocha should not be marked Generated")
	}

	// gruvbox-dark should NOT be Generated.
	if info, ok := nameMap["gruvbox-dark"]; !ok {
		t.Error("expected gruvbox-dark in list")
	} else if info.Generated {
		t.Error("gruvbox-dark should not be marked Generated")
	}
}

func TestListThemesUseCase_ListBuiltins(t *testing.T) {
	store := newStubThemeStore()
	builtins := newStubPaletteSource()
	builtins.palettes["tokyo-night-dark"] = []byte("yaml")
	builtins.palettes["catppuccin-mocha"] = []byte("yaml")

	uc := application.NewListThemesUseCase(store, builtins)
	names := uc.ListBuiltins()

	if len(names) != 2 {
		t.Fatalf("expected 2 builtins, got %d", len(names))
	}

	found := make(map[string]bool)
	for _, n := range names {
		found[n] = true
	}
	if !found["tokyo-night-dark"] {
		t.Error("expected tokyo-night-dark in builtins")
	}
	if !found["catppuccin-mocha"] {
		t.Error("expected catppuccin-mocha in builtins")
	}
}

func TestListThemesUseCase_Complete(t *testing.T) {
	store := newStubThemeStore()

	// Set up a complete theme (palette.yaml + all output files).
	if err := store.EnsureThemeDir("complete"); err != nil {
		t.Fatal(err)
	}
	writeStubFile(t, store, "complete", "palette.yaml")
	for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
		writeStubFile(t, store, "complete", f)
	}

	// Set up an incomplete theme (palette.yaml + missing style.json).
	if err := store.EnsureThemeDir("incomplete"); err != nil {
		t.Fatal(err)
	}
	writeStubFile(t, store, "incomplete", "palette.yaml")
	for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss"} {
		writeStubFile(t, store, "incomplete", f)
	}

	builtins := newStubPaletteSource()

	uc := application.NewListThemesUseCase(store, builtins)
	themes, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, info := range themes {
		if info.Name == "complete" && !info.Complete {
			t.Error("expected complete theme to have Complete=true")
		}
		if info.Name == "incomplete" && info.Complete {
			t.Error("expected incomplete theme to have Complete=false")
		}
	}
}

// writeStubFile writes a stub file to the in-memory theme store.
func writeStubFile(t *testing.T, store *stubThemeStore, theme, filename string) {
	t.Helper()
	w, err := store.OpenWriter(theme, filename)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write([]byte("content"))
	_ = w.Close()
}

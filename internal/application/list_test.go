package application_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/application"
)

func TestListThemesUseCase_ReturnsInstalledThemes(t *testing.T) {
	store := newStubThemeStore()

	// Set up two themes with all output files.
	for _, theme := range []string{"catppuccin", "dracula"} {
		if err := store.EnsureThemeDir(theme); err != nil {
			t.Fatal(err)
		}
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

	// Verify theme names are present (order from ListThemes).
	names := make(map[string]bool)
	for _, info := range themes {
		names[info.Name] = true
	}
	if !names["catppuccin"] {
		t.Error("expected catppuccin in list")
	}
	if !names["dracula"] {
		t.Error("expected dracula in list")
	}
}

func TestListThemesUseCase_MarksSelected(t *testing.T) {
	store := newStubThemeStore()

	// Set up two themes with all output files.
	for _, theme := range []string{"alpha", "beta"} {
		if err := store.EnsureThemeDir(theme); err != nil {
			t.Fatal(err)
		}
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

func TestListThemesUseCase_EmptyDir(t *testing.T) {
	store := newStubThemeStore()
	builtins := newStubPaletteSource()

	uc := application.NewListThemesUseCase(store, builtins)
	themes, err := uc.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(themes) != 0 {
		t.Errorf("expected 0 themes, got %d", len(themes))
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

	// Set up a complete theme.
	if err := store.EnsureThemeDir("complete"); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
		w, err := store.OpenWriter("complete", f)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("content"))
		_ = w.Close()
	}

	// Set up an incomplete theme (missing style.json).
	if err := store.EnsureThemeDir("incomplete"); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss"} {
		w, err := store.OpenWriter("incomplete", f)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("content"))
		_ = w.Close()
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

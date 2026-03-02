package application_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/application"
)

func TestSelectThemeUseCase_Success(t *testing.T) {
	store := newStubThemeStore()

	// Set up a theme with all 5 output files.
	theme := "tokyonight"
	if err := store.EnsureThemeDir(theme); err != nil {
		t.Fatal(err)
	}
	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		w, err := store.OpenWriter(theme, f)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("content"))
		_ = w.Close()
	}

	uc := application.NewSelectThemeUseCase(store, nil, nil)

	err := uc.Execute(theme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify store.Select was called by checking the selected theme.
	selected, err := store.SelectedTheme()
	if err != nil {
		t.Fatalf("SelectedTheme() error: %v", err)
	}
	if selected != theme {
		t.Errorf("SelectedTheme() = %q, want %q", selected, theme)
	}
}

func TestSelectThemeUseCase_ThemeNotFound(t *testing.T) {
	store := newStubThemeStore()

	uc := application.NewSelectThemeUseCase(store, nil, nil)

	err := uc.Execute("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent theme, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error = %v, want it to mention theme name", err)
	}
}

func TestSelectThemeUseCase_IncompleteTheme(t *testing.T) {
	store := newStubThemeStore()

	// Set up a theme with only some output files (missing style.qss and style.json).
	theme := "incomplete"
	if err := store.EnsureThemeDir(theme); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"style.lua", "style.css", "gtk.css"} {
		w, err := store.OpenWriter(theme, f)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte("content"))
		_ = w.Close()
	}

	uc := application.NewSelectThemeUseCase(store, nil, nil)

	err := uc.Execute(theme)
	if err == nil {
		t.Fatal("expected error for incomplete theme, got nil")
	}

	// Error should mention missing files.
	errMsg := err.Error()
	if !strings.Contains(errMsg, "style.qss") {
		t.Errorf("error = %v, want it to mention style.qss", err)
	}
	if !strings.Contains(errMsg, "style.json") {
		t.Errorf("error = %v, want it to mention style.json", err)
	}
}

func TestSelectThemeUseCase_AutoGeneratesBuiltin(t *testing.T) {
	store := newStubThemeStore()

	// Set up a stub builtins source with a palette.
	builtins := newStubPaletteSource()
	builtins.palettes["catppuccin-mocha"] = []byte("palette-yaml")

	// Build a GenerateThemeUseCase that the select UC can delegate to.
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	targets := makeStubTargets()
	generateUC := application.NewGenerateThemeUseCase(
		&stubGenParser{palette: pal},
		&stubGenTokenizer{tokenSet: ts},
		targets,
		store,
		builtins,
	)

	uc := application.NewSelectThemeUseCase(store, builtins, generateUC)

	// Theme does not exist on disk yet — select should auto-generate it.
	err := uc.Execute("catppuccin-mocha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the theme was generated (output files exist).
	for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
		if !store.FileExists("catppuccin-mocha", f) {
			t.Errorf("expected file %q to be generated", f)
		}
	}

	// Verify the theme was selected.
	selected, err := store.SelectedTheme()
	if err != nil {
		t.Fatalf("SelectedTheme() error: %v", err)
	}
	if selected != "catppuccin-mocha" {
		t.Errorf("SelectedTheme() = %q, want %q", selected, "catppuccin-mocha")
	}
}

func TestSelectThemeUseCase_NonBuiltinStillFails(t *testing.T) {
	store := newStubThemeStore()
	builtins := newStubPaletteSource() // empty — no built-ins

	uc := application.NewSelectThemeUseCase(store, builtins, nil)

	err := uc.Execute("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent theme, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error = %v, want it to mention theme name", err)
	}
}

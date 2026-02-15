package application_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/application"
)

// writePaletteToStore writes a valid palette.yaml into the stub store for the given theme.
func writePaletteToStore(t *testing.T, store *stubThemeStore, themeName string) {
	t.Helper()

	paletteYAML := `system: "base24"
name: "Test Palette"
author: "test"
variant: "dark"
palette:
  base00: "1a1b26"
  base01: "1f2335"
  base02: "292e42"
  base03: "565f89"
  base04: "a9b1d6"
  base05: "c0caf5"
  base06: "c0caf5"
  base07: "c8d3f5"
  base08: "f7768e"
  base09: "ff9e64"
  base0A: "e0af68"
  base0B: "9ece6a"
  base0C: "7dcfff"
  base0D: "7aa2f7"
  base0E: "bb9af7"
  base0F: "db4b4b"
  base10: "16161e"
  base11: "101014"
  base12: "ff899d"
  base13: "e9c582"
  base14: "afd67a"
  base15: "97d8f8"
  base16: "8db6fa"
  base17: "c8acf8"
`

	if err := store.EnsureThemeDir(themeName); err != nil {
		t.Fatalf("ensure theme dir: %v", err)
	}
	w, err := store.OpenWriter(themeName, "palette.yaml")
	if err != nil {
		t.Fatalf("open palette writer: %v", err)
	}
	_, _ = w.Write([]byte(paletteYAML))
	_ = w.Close()
}

// writeUniversalToStore writes a minimal universal.yaml into the stub store.
func writeUniversalToStore(t *testing.T, store *stubThemeStore, themeName string) {
	t.Helper()

	universalYAML := `tokens:
  surface.background:
    color: "#1a1b26"
  surface.background.raised:
    color: "#1f2335"
  text.primary:
    color: "#c0caf5"
  text.muted:
    color: "#565f89"
  status.error:
    color: "#f7768e"
  status.warning:
    color: "#e0af68"
  status.success:
    color: "#9ece6a"
  status.info:
    color: "#7dcfff"
  syntax.keyword:
    color: "#bb9af7"
    bold: true
`

	w, err := store.OpenWriter(themeName, "universal.yaml")
	if err != nil {
		t.Fatalf("open universal writer: %v", err)
	}
	_, _ = w.Write([]byte(universalYAML))
	_ = w.Close()
}

func TestPreviewThemeUseCase_OutputContainsPaletteColors(t *testing.T) {
	store := newStubThemeStore()
	writePaletteToStore(t, store, "test-dark")
	writeUniversalToStore(t, store, "test-dark")

	parser := yamlparser.NewParser()
	uc := application.NewPreviewThemeUseCase(store, parser, fileio.ReadUniversal)

	var buf bytes.Buffer
	err := uc.Execute("test-dark", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Should contain hex values for palette slots.
	expectedSlots := []string{
		"base00", "base01", "base02", "base03",
		"base04", "base05", "base06", "base07",
		"base08", "base09", "base0A", "base0B",
		"base0C", "base0D", "base0E", "base0F",
		"base10", "base11", "base12", "base13",
		"base14", "base15", "base16", "base17",
	}
	for _, slot := range expectedSlots {
		if !strings.Contains(output, slot) {
			t.Errorf("output should contain slot name %q", slot)
		}
	}
}

func TestPreviewThemeUseCase_OutputContainsANSI(t *testing.T) {
	store := newStubThemeStore()
	writePaletteToStore(t, store, "test-dark")
	writeUniversalToStore(t, store, "test-dark")

	parser := yamlparser.NewParser()
	uc := application.NewPreviewThemeUseCase(store, parser, fileio.ReadUniversal)

	var buf bytes.Buffer
	err := uc.Execute("test-dark", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// ANSI escape sequences use ESC (0x1b) followed by '['.
	if !strings.Contains(output, "\x1b[") {
		t.Error("output should contain ANSI escape sequences")
	}

	// Should contain 24-bit color escape codes: ESC[38;2;R;G;Bm or ESC[48;2;R;G;Bm
	if !strings.Contains(output, ";2;") {
		t.Error("output should contain 24-bit (truecolor) ANSI escape codes")
	}
}

func TestPreviewThemeUseCase_OutputContainsTokenPreview(t *testing.T) {
	store := newStubThemeStore()
	writePaletteToStore(t, store, "test-dark")
	writeUniversalToStore(t, store, "test-dark")

	parser := yamlparser.NewParser()
	uc := application.NewPreviewThemeUseCase(store, parser, fileio.ReadUniversal)

	var buf bytes.Buffer
	err := uc.Execute("test-dark", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Key semantic tokens should appear in the output.
	keyTokens := []string{
		"surface.background",
		"text.primary",
	}
	for _, tok := range keyTokens {
		if !strings.Contains(output, tok) {
			t.Errorf("output should contain semantic token %q", tok)
		}
	}
}

func TestPreviewThemeUseCase_ThemeNotFound(t *testing.T) {
	store := newStubThemeStore()
	parser := yamlparser.NewParser()
	uc := application.NewPreviewThemeUseCase(store, parser, fileio.ReadUniversal)

	var buf bytes.Buffer
	err := uc.Execute("nonexistent-theme", &buf)
	if err == nil {
		t.Fatal("expected error for non-existent theme, got nil")
	}
}

package domain_test

import (
	"errors"
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

// makeBase24Colors creates a full map of 24 hex color values for testing.
func makeBase24Colors() map[string]string {
	return map[string]string{
		"base00": "#1a1b26", "base01": "#16161e", "base02": "#2f3549", "base03": "#444b6a",
		"base04": "#787c99", "base05": "#a9b1d6", "base06": "#cbccd1", "base07": "#d5d6db",
		"base08": "#f7768e", "base09": "#ff9e64", "base0A": "#e0af68", "base0B": "#9ece6a",
		"base0C": "#2ac3de", "base0D": "#7aa2f7", "base0E": "#bb9af7", "base0F": "#ab6f60",
		"base10": "#14141e", "base11": "#111118", "base12": "#e06c75", "base13": "#d19a66",
		"base14": "#98c379", "base15": "#56b6c2", "base16": "#61afef", "base17": "#c678dd",
	}
}

// makeBase16Colors creates a map of only the first 16 base colors.
func makeBase16Colors() map[string]string {
	return map[string]string{
		"base00": "#1a1b26", "base01": "#16161e", "base02": "#2f3549", "base03": "#444b6a",
		"base04": "#787c99", "base05": "#a9b1d6", "base06": "#cbccd1", "base07": "#d5d6db",
		"base08": "#f7768e", "base09": "#ff9e64", "base0A": "#e0af68", "base0B": "#9ece6a",
		"base0C": "#2ac3de", "base0D": "#7aa2f7", "base0E": "#bb9af7", "base0F": "#ab6f60",
	}
}

func TestNewPalette_Full24Colors(t *testing.T) {
	colors := makeBase24Colors()

	pal, err := domain.NewPalette("tokyonight", "folke", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}
	if pal.Name != "tokyonight" {
		t.Errorf("Name = %q, want %q", pal.Name, "tokyonight")
	}
	if pal.Author != "folke" {
		t.Errorf("Author = %q, want %q", pal.Author, "folke")
	}
	if pal.Variant != "dark" {
		t.Errorf("Variant = %q, want %q", pal.Variant, "dark")
	}
	if pal.System != "base24" {
		t.Errorf("System = %q, want %q", pal.System, "base24")
	}
	if pal.Slug != "tokyonight-dark" {
		t.Errorf("Slug = %q, want %q", pal.Slug, "tokyonight-dark")
	}

	// Verify all 24 colors were stored
	for i := 0; i < 24; i++ {
		c := pal.Base(i)
		if c.IsNone {
			t.Errorf("Base(%d) returned NoneColor, expected a valid color", i)
		}
	}
}

func TestPalette_Slot_ByName(t *testing.T) {
	colors := makeBase24Colors()
	pal, err := domain.NewPalette("test", "author", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	c, err := pal.Slot("base0D")
	if err != nil {
		t.Fatalf("Slot(base0D) unexpected error: %v", err)
	}

	// #7aa2f7 -> R=122, G=162, B=247
	if c.R != 122 || c.G != 162 || c.B != 247 {
		t.Errorf("Slot(base0D) = {R:%d, G:%d, B:%d}, want {R:122, G:162, B:247}", c.R, c.G, c.B)
	}
}

func TestPalette_Base_ByIndex(t *testing.T) {
	colors := makeBase24Colors()
	pal, err := domain.NewPalette("test", "author", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	// base0D is at index 13
	byIndex := pal.Base(13)

	slotColor, err := pal.Slot("base0D")
	if err != nil {
		t.Fatalf("Slot(base0D) unexpected error: %v", err)
	}

	if !byIndex.Equal(slotColor) {
		t.Errorf("Base(13) = %s, Slot(base0D) = %s, expected equal", byIndex.Hex(), slotColor.Hex())
	}
}

func TestNewPalette_Base16Fallbacks(t *testing.T) {
	colors := makeBase16Colors()
	pal, err := domain.NewPalette("test", "author", "dark", "base16", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	// Fallback rules: base10=base00, base11=base00, base12=base08,
	// base13=base0A, base14=base0B, base15=base0C, base16=base0D, base17=base0E
	fallbacks := []struct {
		slot     string
		fallback string
	}{
		{"base10", "base00"},
		{"base11", "base00"},
		{"base12", "base08"},
		{"base13", "base0A"},
		{"base14", "base0B"},
		{"base15", "base0C"},
		{"base16", "base0D"},
		{"base17", "base0E"},
	}

	for _, fb := range fallbacks {
		got, err := pal.Slot(fb.slot)
		if err != nil {
			t.Fatalf("Slot(%s) unexpected error: %v", fb.slot, err)
		}
		want, err := pal.Slot(fb.fallback)
		if err != nil {
			t.Fatalf("Slot(%s) unexpected error: %v", fb.fallback, err)
		}
		if !got.Equal(want) {
			t.Errorf("Slot(%s) = %s, want Slot(%s) = %s (fallback)", fb.slot, got.Hex(), fb.fallback, want.Hex())
		}
	}
}

func TestNewPalette_MissingSlot(t *testing.T) {
	// Map missing base00
	colors := makeBase16Colors()
	delete(colors, "base00")

	_, err := domain.NewPalette("test", "author", "dark", "base24", colors)
	if err == nil {
		t.Fatal("NewPalette() expected error for missing base00, got nil")
	}
	var pe *domain.ParseError
	if !errors.As(err, &pe) {
		t.Errorf("error type = %T, want *domain.ParseError", err)
	}
}

func TestPalette_SlotNames(t *testing.T) {
	colors := makeBase24Colors()
	pal, err := domain.NewPalette("test", "author", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	names := pal.SlotNames()
	if len(names) != 24 {
		t.Fatalf("SlotNames() returned %d names, want 24", len(names))
	}

	expected := []string{
		"base00", "base01", "base02", "base03",
		"base04", "base05", "base06", "base07",
		"base08", "base09", "base0A", "base0B",
		"base0C", "base0D", "base0E", "base0F",
		"base10", "base11", "base12", "base13",
		"base14", "base15", "base16", "base17",
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("SlotNames()[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

func TestPalette_Slot_Unknown(t *testing.T) {
	colors := makeBase24Colors()
	pal, err := domain.NewPalette("test", "author", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	_, err = pal.Slot("baseZZ")
	if err == nil {
		t.Fatal("Slot(baseZZ) expected error, got nil")
	}
}

func TestPalette_WithOverrides(t *testing.T) {
	// Arrange: Create a Palette with Overrides map containing syntax.keyword override
	colors := makeBase24Colors()
	pal, err := domain.NewPalette("test", "author", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	// Create override with color and italic
	overrideColor, _ := domain.ParseHex("#ff5555")
	override := domain.TokenOverride{
		Color:  &overrideColor,
		Italic: true,
	}

	// Set the override on the palette
	pal.Overrides = map[string]domain.TokenOverride{
		"syntax.keyword": override,
	}

	// Act: Access the override
	got, exists := pal.Overrides["syntax.keyword"]

	// Assert: Returns the TokenOverride struct with correct values
	if !exists {
		t.Fatal("pal.Overrides[\"syntax.keyword\"] does not exist")
	}
	if !got.HasColor() {
		t.Error("override.HasColor() = false, want true")
	}
	if !got.Color.Equal(overrideColor) {
		t.Errorf("override.Color = %s, want %s", got.Color.Hex(), overrideColor.Hex())
	}
	if !got.Italic {
		t.Error("override.Italic = false, want true")
	}
}

func TestPalette_OverridesNilByDefault(t *testing.T) {
	// Verify that Overrides is nil when no overrides are set (default)
	colors := makeBase24Colors()
	pal, err := domain.NewPalette("test", "author", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	if pal.Overrides != nil {
		t.Errorf("pal.Overrides = %v, want nil", pal.Overrides)
	}
}

package flair_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

// validPaletteYAML is a complete base24 palette YAML for testing.
const validPaletteYAML = `system: "base24"
name: "Tokyo Night Dark"
author: "Test Author"
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

// TestParsePalette_ValidYAML tests parsing a valid base24 palette YAML.
func TestParsePalette_ValidYAML(t *testing.T) {
	r := strings.NewReader(validPaletteYAML)
	pal, err := flair.ParsePalette(r)
	if err != nil {
		t.Fatalf("ParsePalette() unexpected error: %v", err)
	}

	// Verify metadata.
	if got := pal.Name(); got != "Tokyo Night Dark" {
		t.Errorf("Name() = %q, want %q", got, "Tokyo Night Dark")
	}
	if got := pal.Author(); got != "Test Author" {
		t.Errorf("Author() = %q, want %q", got, "Test Author")
	}
	if got := pal.Variant(); got != "dark" {
		t.Errorf("Variant() = %q, want %q", got, "dark")
	}

	// Verify colors via Base() method.
	// base00 = "1a1b26" -> R=0x1a, G=0x1b, B=0x26.
	c00 := pal.Base(0x00)
	if c00.R != 0x1a || c00.G != 0x1b || c00.B != 0x26 {
		t.Errorf("Base(0x00) = {%d, %d, %d}, want {26, 27, 38}", c00.R, c00.G, c00.B)
	}

	// base0D = "7aa2f7" -> R=0x7a, G=0xa2, B=0xf7.
	c0d := pal.Base(0x0D)
	if c0d.R != 0x7a || c0d.G != 0xa2 || c0d.B != 0xf7 {
		t.Errorf("Base(0x0D) = {%d, %d, %d}, want {122, 162, 247}", c0d.R, c0d.G, c0d.B)
	}

	// base17 = "c8acf8" -> R=0xc8, G=0xac, B=0xf8.
	c17 := pal.Base(0x17)
	if c17.R != 0xc8 || c17.G != 0xac || c17.B != 0xf8 {
		t.Errorf("Base(0x17) = {%d, %d, %d}, want {200, 172, 248}", c17.R, c17.G, c17.B)
	}
}

// TestParsePalette_MissingColors tests that parsing fails when required colors are missing.
func TestParsePalette_MissingColors(t *testing.T) {
	// YAML missing base10-base17 slots.
	yaml := `system: "base24"
name: "Incomplete Palette"
author: "Test Author"
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
`
	r := strings.NewReader(yaml)
	_, err := flair.ParsePalette(r)
	if err == nil {
		t.Fatal("ParsePalette() expected error for missing colors, got nil")
	}

	// Error message should indicate missing color.
	if !strings.Contains(err.Error(), "missing") && !strings.Contains(err.Error(), "base10") {
		t.Errorf("error message should mention missing color, got: %v", err)
	}
}

// TestParsePalette_InvalidHex tests that parsing fails for invalid hex color values.
func TestParsePalette_InvalidHex(t *testing.T) {
	yaml := `system: "base24"
name: "Bad Hex Palette"
author: "Test Author"
variant: "dark"
palette:
  base00: "GGGGGG"
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
	r := strings.NewReader(yaml)
	_, err := flair.ParsePalette(r)
	if err == nil {
		t.Fatal("ParsePalette() expected error for invalid hex, got nil")
	}

	// Error message should indicate invalid hex.
	if !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "hex") {
		t.Errorf("error message should mention invalid hex, got: %v", err)
	}
}

// TestPalette_Base tests the Base method with valid indices.
func TestPalette_Base(t *testing.T) {
	r := strings.NewReader(validPaletteYAML)
	pal, err := flair.ParsePalette(r)
	if err != nil {
		t.Fatalf("ParsePalette() unexpected error: %v", err)
	}

	// Test base00 = "1a1b26" -> R=0x1a, G=0x1b, B=0x26.
	t.Run("base00", func(t *testing.T) {
		c := pal.Base(0x00)
		assertColor(t, c, 0x1a, 0x1b, 0x26)
	})

	// Test base0D = "7aa2f7" -> R=0x7a, G=0xa2, B=0xf7.
	t.Run("base0D", func(t *testing.T) {
		c := pal.Base(0x0D)
		assertColor(t, c, 0x7a, 0xa2, 0xf7)
	})

	// Test base17 = "c8acf8" -> R=0xc8, G=0xac, B=0xf8.
	t.Run("base17", func(t *testing.T) {
		c := pal.Base(0x17)
		assertColor(t, c, 0xc8, 0xac, 0xf8)
	})
}

// assertColor is a test helper that asserts a color matches expected RGB values.
func assertColor(t *testing.T, c *flair.Color, wantR, wantG, wantB uint8) {
	t.Helper()
	if c == nil {
		t.Fatal("color is nil, want non-nil")
		return // unreachable but satisfies staticcheck
	}
	if c.R != wantR || c.G != wantG || c.B != wantB {
		t.Errorf("color = {R:%d, G:%d, B:%d}, want {R:%d, G:%d, B:%d}",
			c.R, c.G, c.B, wantR, wantG, wantB)
	}
}

// TestPalette_Base_OutOfRange tests that Base returns nil for out-of-range indices.
func TestPalette_Base_OutOfRange(t *testing.T) {
	r := strings.NewReader(validPaletteYAML)
	pal, err := flair.ParsePalette(r)
	if err != nil {
		t.Fatalf("ParsePalette() unexpected error: %v", err)
	}

	tests := []struct {
		name  string
		index int
	}{
		{"negative index", -1},
		{"out of range", 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pal.Base(tt.index)
			if c != nil {
				t.Errorf("Base(%d) = %v, want nil", tt.index, c)
			}
		})
	}
}

// TestPalette_Accessors tests the Name, Author, and Variant accessor methods.
func TestPalette_Accessors(t *testing.T) {
	r := strings.NewReader(validPaletteYAML)
	pal, err := flair.ParsePalette(r)
	if err != nil {
		t.Fatalf("ParsePalette() unexpected error: %v", err)
	}

	tests := []struct {
		name   string
		method func() string
		want   string
	}{
		{"Name", pal.Name, "Tokyo Night Dark"},
		{"Author", pal.Author, "Test Author"},
		{"Variant", pal.Variant, "dark"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method(); got != tt.want {
				t.Errorf("%s() = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

package domain_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

// makeValidDarkPalette creates a valid dark palette with correct luminance
// ordering: base00 (dark) < base01 < ... < base07 (light).
// Uses Tokyo Night Dark reference values.
func makeValidDarkPalette(t *testing.T) *domain.Palette {
	t.Helper()
	colors := map[string]string{
		"base00": "1a1b26", // Default Background (darkest neutral)
		"base01": "1f2335", // Lighter Background
		"base02": "292e42", // Selection Background
		"base03": "565f89", // Comments
		"base04": "a9b1d6", // Dark Foreground
		"base05": "c0caf5", // Default Foreground
		"base06": "c0caf5", // Light Foreground
		"base07": "c8d3f5", // Lightest Foreground
		"base08": "f7768e", // Red
		"base09": "ff9e64", // Orange
		"base0A": "e0af68", // Yellow
		"base0B": "9ece6a", // Green
		"base0C": "7dcfff", // Cyan
		"base0D": "7aa2f7", // Blue
		"base0E": "bb9af7", // Magenta
		"base0F": "db4b4b", // Brown/Dark Red
		"base10": "16161e", // Darker Background
		"base11": "101014", // Darkest Background
		"base12": "ff899d", // Bright Red (brighter than base08)
		"base13": "e9c582", // Bright Yellow (brighter than base0A)
		"base14": "afd67a", // Bright Green (brighter than base0B)
		"base15": "97d8f8", // Bright Cyan (brighter than base0C)
		"base16": "8db6fa", // Bright Blue (brighter than base0D)
		"base17": "c8acf8", // Bright Magenta (brighter than base0E)
	}
	pal, err := domain.NewPalette("tokyo-night", "test", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("makeValidDarkPalette: %v", err)
	}
	return pal
}

func TestValidatePalette_ValidDark(t *testing.T) {
	pal := makeValidDarkPalette(t)

	violations := domain.ValidatePalette(pal)
	if len(violations) != 0 {
		t.Errorf("ValidatePalette() returned %d violations for valid dark palette, want 0:\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

func TestValidatePalette_MissingSlot(t *testing.T) {
	pal := makeValidDarkPalette(t)
	// Set one slot to NoneColor to simulate a missing color
	pal.Colors[5] = domain.NoneColor() // base05

	violations := domain.ValidatePalette(pal)
	if len(violations) == 0 {
		t.Fatal("ValidatePalette() returned no violations, expected violation for missing slot")
	}

	found := false
	for _, v := range violations {
		if strings.Contains(v, "base05") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("violations = %v, expected mention of base05", violations)
	}
}

func TestValidatePalette_DarkLuminanceViolation(t *testing.T) {
	pal := makeValidDarkPalette(t)
	// Swap base00 (background) and base05 (foreground) so base00 is lighter
	// than base05, violating the dark palette rule: base00.Luminance < base05.Luminance
	pal.Colors[0], pal.Colors[5] = pal.Colors[5], pal.Colors[0]

	violations := domain.ValidatePalette(pal)
	if len(violations) == 0 {
		t.Fatal("ValidatePalette() returned no violations, expected luminance ordering violation")
	}

	found := false
	for _, v := range violations {
		if strings.Contains(v, "luminance") || strings.Contains(v, "Luminance") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("violations = %v, expected mention of luminance ordering", violations)
	}
}

func TestValidatePalette_LightLuminanceViolation(t *testing.T) {
	// Create a light palette where base00 should be lighter than base05
	// But set base00 to be darker than base05
	colors := map[string]string{
		"base00": "1a1b26", // Dark -- wrong for light theme (should be light)
		"base01": "1f2335",
		"base02": "292e42",
		"base03": "565f89",
		"base04": "a9b1d6",
		"base05": "c0caf5", // Light -- should be darker than base00 in light theme
		"base06": "c0caf5",
		"base07": "c8d3f5",
		"base08": "f7768e",
		"base09": "ff9e64",
		"base0A": "e0af68",
		"base0B": "9ece6a",
		"base0C": "7dcfff",
		"base0D": "7aa2f7",
		"base0E": "bb9af7",
		"base0F": "db4b4b",
		"base10": "16161e",
		"base11": "101014",
		"base12": "ff899d",
		"base13": "e9c582",
		"base14": "afd67a",
		"base15": "97d8f8",
		"base16": "8db6fa",
		"base17": "c8acf8",
	}
	pal, err := domain.NewPalette("light-test", "test", "light", "base24", colors)
	if err != nil {
		t.Fatalf("NewPalette() unexpected error: %v", err)
	}

	violations := domain.ValidatePalette(pal)
	if len(violations) == 0 {
		t.Fatal("ValidatePalette() returned no violations, expected luminance ordering violation for light palette")
	}

	found := false
	for _, v := range violations {
		if strings.Contains(v, "luminance") || strings.Contains(v, "Luminance") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("violations = %v, expected mention of luminance ordering", violations)
	}
}

func TestValidatePalette_MonotonicityWarning(t *testing.T) {
	pal := makeValidDarkPalette(t)
	// Swap base01 and base02 so that luminance dips:
	// base01 should be lighter than base02 after the swap,
	// breaking the monotonic increase base00 -> base07.
	// Actually, we need base02 to be darker than base01.
	// Set base02 to something very dark (darker than base01).
	dark, err := domain.ParseHex("0a0a0a")
	if err != nil {
		t.Fatalf("ParseHex: %v", err)
	}
	pal.Colors[2] = dark // base02 is now darker than base01

	violations := domain.ValidatePalette(pal)
	if len(violations) == 0 {
		t.Fatal("ValidatePalette() returned no violations, expected monotonicity warning")
	}

	found := false
	for _, v := range violations {
		if strings.Contains(v, "monoton") || strings.Contains(v, "Monoton") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("violations = %v, expected mention of monotonicity", violations)
	}
}

func TestValidatePalette_BrightVariantWarning(t *testing.T) {
	pal := makeValidDarkPalette(t)
	// Make base12 (bright red) dimmer than base08 (red)
	// base08 is currently "f7768e" — set base12 to something very dark
	dim, err := domain.ParseHex("110000")
	if err != nil {
		t.Fatalf("ParseHex: %v", err)
	}
	pal.Colors[18] = dim // base12 index is 18 (0-indexed: base10=16, base11=17, base12=18)

	violations := domain.ValidatePalette(pal)
	if len(violations) == 0 {
		t.Fatal("ValidatePalette() returned no violations, expected bright variant warning")
	}

	found := false
	for _, v := range violations {
		if strings.Contains(v, "bright") || strings.Contains(v, "Bright") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("violations = %v, expected mention of bright variant", violations)
	}
}

package mapper_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/deriver"
	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// tokyoNightDarkPalette returns the reference Tokyo Night Dark base24 palette
// used as the canonical test fixture.
func tokyoNightDarkPalette(t *testing.T) *domain.Palette {
	t.Helper()

	colors := map[string]string{
		"base00": "1a1b26",
		"base01": "1f2335",
		"base02": "292e42",
		"base03": "565f89",
		"base04": "a9b1d6",
		"base05": "c0caf5",
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

	pal, err := domain.NewPalette("Tokyo Night Dark", "Michael Ball", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("failed to create test palette: %v", err)
	}
	return pal
}

// buildResolvedTheme constructs a ResolvedTheme from the Tokyo Night Dark palette
// using the default deriver.
func buildResolvedTheme(t *testing.T) *domain.ResolvedTheme {
	t.Helper()
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)
	return &domain.ResolvedTheme{
		Name:    "tokyo-night-dark",
		Variant: "dark",
		Palette: pal,
		Tokens:  ts,
	}
}

// mustParseHex is a test helper that parses a hex color or fails the test.
func mustParseHex(t *testing.T, hex string) domain.Color {
	t.Helper()
	c, err := domain.ParseHex(hex)
	if err != nil {
		t.Fatalf("failed to parse hex %q: %v", hex, err)
	}
	return c
}

// isValidHex checks whether a string is a valid 7-character hex color (#rrggbb).
func isValidHex(s string) bool {
	if len(s) != 7 || s[0] != '#' {
		return false
	}
	for _, c := range s[1:] {
		isDigit := c >= '0' && c <= '9'
		isLower := c >= 'a' && c <= 'f'
		isUpper := c >= 'A' && c <= 'F'
		if !isDigit && !isLower && !isUpper {
			return false
		}
	}
	return true
}

// TestStylixMapper_Interface verifies that the Stylix mapper implements
// ports.Mapper and Name() returns "stylix".
func TestStylixMapper_Interface(t *testing.T) {
	m := mapper.NewStylix()

	// Compile-time interface check.
	var _ ports.Mapper = m

	if name := m.Name(); name != "stylix" {
		t.Errorf("Name() = %q, want %q", name, "stylix")
	}
}

// TestStylixMapper_KeyCount verifies that the mapper produces at least 60
// key-value pairs from the Tokyo Night Dark palette.
func TestStylixMapper_KeyCount(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewStylix()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	st, ok := result.(*ports.StylixTheme)
	if !ok {
		t.Fatalf("Map() returned %T, want *ports.StylixTheme", result)
	}

	if len(st.Values) < 60 {
		t.Errorf("StylixTheme.Values has %d entries, want >= 60", len(st.Values))
		// Print all keys for debugging
		for k := range st.Values {
			t.Logf("  key: %s", k)
		}
	}

	// All values must be valid hex colors.
	for key, val := range st.Values {
		if !isValidHex(val) {
			t.Errorf("key %q has invalid hex value %q", key, val)
		}
	}
}

// TestStylixMapper_PalettePassthrough verifies that all 24 base palette
// slots are included in the output with correct hex values.
func TestStylixMapper_PalettePassthrough(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewStylix()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	st := result.(*ports.StylixTheme)

	// Verify all 24 base slots are present.
	expectedSlots := map[string]string{
		"base00": "#1a1b26",
		"base01": "#1f2335",
		"base02": "#292e42",
		"base03": "#565f89",
		"base04": "#a9b1d6",
		"base05": "#c0caf5",
		"base06": "#c0caf5",
		"base07": "#c8d3f5",
		"base08": "#f7768e",
		"base09": "#ff9e64",
		"base0A": "#e0af68",
		"base0B": "#9ece6a",
		"base0C": "#7dcfff",
		"base0D": "#7aa2f7",
		"base0E": "#bb9af7",
		"base0F": "#db4b4b",
		"base10": "#16161e",
		"base11": "#101014",
		"base12": "#ff899d",
		"base13": "#e9c582",
		"base14": "#afd67a",
		"base15": "#97d8f8",
		"base16": "#8db6fa",
		"base17": "#c8acf8",
	}

	for slot, wantHex := range expectedSlots {
		got, ok := st.Values[slot]
		if !ok {
			t.Errorf("missing base palette slot %q in StylixTheme.Values", slot)
			continue
		}
		if !strings.EqualFold(got, wantHex) {
			t.Errorf("slot %q = %q, want %q", slot, got, wantHex)
		}
	}
}

// TestStylixMapper_SemanticTokens verifies that semantic token keys are
// present in the output with correct hex color values.
func TestStylixMapper_SemanticTokens(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewStylix()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	st := result.(*ports.StylixTheme)

	// Check a representative set of semantic keys.
	// The mapper should convert dot-separated token paths to kebab-case keys.
	semanticChecks := []struct {
		key     string
		wantHex string
	}{
		{"surface-bg", "#1a1b26"},
		{"text-primary", "#c0caf5"},
		{"syntax-keyword", "#bb9af7"},
		{"syntax-string", "#9ece6a"},
		{"syntax-function", "#7aa2f7"},
		{"status-error", "#ff899d"},
		{"status-warning", "#e9c582"},
		{"status-success", "#afd67a"},
		{"accent-primary", "#7aa2f7"},
		{"git-added", "#9ece6a"},
		{"git-deleted", "#f7768e"},
		{"terminal-red", "#f7768e"},
		{"terminal-blue", "#7aa2f7"},
	}

	for _, tc := range semanticChecks {
		t.Run(tc.key, func(t *testing.T) {
			got, ok := st.Values[tc.key]
			if !ok {
				// Print available keys for debugging
				keys := make([]string, 0, len(st.Values))
				for k := range st.Values {
					keys = append(keys, k)
				}
				t.Fatalf("key %q not found in StylixTheme.Values. Available keys: %v", tc.key, keys)
			}
			want := mustParseHex(t, tc.wantHex)
			gotColor := mustParseHex(t, got)
			if !gotColor.Equal(want) {
				t.Errorf("key %q = %q, want %q", tc.key, got, tc.wantHex)
			}
		})
	}

	// Verify no empty values
	for key, val := range st.Values {
		if val == "" {
			t.Errorf("key %q has empty value", key)
		}
	}

	// Print total count for informational purposes
	t.Logf("Total Stylix theme values: %d", len(st.Values))
}

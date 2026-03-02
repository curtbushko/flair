package fileio_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/domain"
)

func TestWritePalette_ContainsAllFields(t *testing.T) {
	// Create a palette with known values.
	var colors [24]domain.Color
	for i := range colors {
		// Use a simple pattern: #00XXYY where XX=i*10, YY=i*11
		colors[i] = domain.Color{R: uint8(i * 10), G: uint8(i * 11), B: uint8(i * 12)}
	}
	pal := &domain.Palette{
		System:  "base24",
		Name:    "Test Theme",
		Author:  "Test Author",
		Variant: "dark",
		Colors:  colors,
	}

	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// Verify metadata fields are present.
	if !strings.Contains(output, `system: "base24"`) {
		t.Error("output should contain system field")
	}
	if !strings.Contains(output, `name: "Test Theme"`) {
		t.Error("output should contain name field")
	}
	if !strings.Contains(output, `author: "Test Author"`) {
		t.Error("output should contain author field")
	}
	if !strings.Contains(output, `variant: "dark"`) {
		t.Error("output should contain variant field")
	}
	if !strings.Contains(output, "palette:") {
		t.Error("output should contain palette section")
	}

	// Verify all 24 base slots are present.
	expectedSlots := []string{
		"base00:", "base01:", "base02:", "base03:",
		"base04:", "base05:", "base06:", "base07:",
		"base08:", "base09:", "base0A:", "base0B:",
		"base0C:", "base0D:", "base0E:", "base0F:",
		"base10:", "base11:", "base12:", "base13:",
		"base14:", "base15:", "base16:", "base17:",
	}
	for _, slot := range expectedSlots {
		if !strings.Contains(output, slot) {
			t.Errorf("output should contain %s", slot)
		}
	}
}

func TestWritePalette_HexFormat(t *testing.T) {
	// Create a palette with a specific color to verify hex formatting.
	var colors [24]domain.Color
	colors[0] = domain.Color{R: 0x1a, G: 0x1b, B: 0x26}
	for i := 1; i < 24; i++ {
		colors[i] = domain.Color{R: 0xff, G: 0xff, B: 0xff}
	}
	pal := &domain.Palette{
		System:  "base24",
		Name:    "Test",
		Author:  "Test",
		Variant: "dark",
		Colors:  colors,
	}

	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// base00 should have the specific color.
	if !strings.Contains(output, `base00: "#1a1b26"`) {
		t.Errorf("output should contain base00 with correct hex value, got: %s", output)
	}
}

func TestWritePalette_WithOverrides(t *testing.T) {
	// Arrange: Create a palette with overrides containing syntax.keyword override.
	var colors [24]domain.Color
	for i := range colors {
		colors[i] = domain.Color{R: 0xff, G: 0xff, B: 0xff}
	}

	overrideColor := domain.Color{R: 0xff, G: 0x00, B: 0xff}
	pal := &domain.Palette{
		System:  "base24",
		Name:    "Test",
		Author:  "Test",
		Variant: "dark",
		Colors:  colors,
		Overrides: map[string]domain.TokenOverride{
			"syntax.keyword": {
				Color:  &overrideColor,
				Bold:   true,
				Italic: false,
			},
		},
	}

	// Act: Call WritePalette.
	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// Assert: Output YAML contains overrides section with syntax.keyword entry.
	if !strings.Contains(output, "overrides:") {
		t.Errorf("output should contain overrides section, got:\n%s", output)
	}
	if !strings.Contains(output, "syntax.keyword:") {
		t.Errorf("output should contain syntax.keyword override, got:\n%s", output)
	}
	if !strings.Contains(output, `color: "#ff00ff"`) {
		t.Errorf("output should contain color with hex value, got:\n%s", output)
	}
	if !strings.Contains(output, "bold: true") {
		t.Errorf("output should contain bold: true, got:\n%s", output)
	}
}

func TestWritePalette_NoOverrides(t *testing.T) {
	// Arrange: Palette with nil Overrides.
	var colors [24]domain.Color
	for i := range colors {
		colors[i] = domain.Color{R: 0xff, G: 0xff, B: 0xff}
	}
	pal := &domain.Palette{
		System:    "base24",
		Name:      "Test",
		Author:    "Test",
		Variant:   "dark",
		Colors:    colors,
		Overrides: nil,
	}

	// Act: Call WritePalette.
	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// Assert: Output YAML does not contain overrides section.
	if strings.Contains(output, "overrides:") {
		t.Errorf("output should NOT contain overrides section when nil, got:\n%s", output)
	}
}

func TestWritePalette_OverridesSorted(t *testing.T) {
	// Arrange: Palette with overrides for z.token, a.token, m.token.
	var colors [24]domain.Color
	for i := range colors {
		colors[i] = domain.Color{R: 0xff, G: 0xff, B: 0xff}
	}

	color := domain.Color{R: 0x11, G: 0x22, B: 0x33}
	pal := &domain.Palette{
		System:  "base24",
		Name:    "Test",
		Author:  "Test",
		Variant: "dark",
		Colors:  colors,
		Overrides: map[string]domain.TokenOverride{
			"z.token": {Color: &color},
			"a.token": {Color: &color},
			"m.token": {Color: &color},
		},
	}

	// Act: Call WritePalette.
	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// Assert: Overrides appear in alphabetical order (a.token, m.token, z.token).
	aIdx := strings.Index(output, "a.token:")
	mIdx := strings.Index(output, "m.token:")
	zIdx := strings.Index(output, "z.token:")

	if aIdx == -1 || mIdx == -1 || zIdx == -1 {
		t.Fatalf("output should contain all tokens, got:\n%s", output)
	}
	if aIdx >= mIdx || mIdx >= zIdx {
		t.Errorf("overrides should be sorted alphabetically (a < m < z), positions: a=%d, m=%d, z=%d", aIdx, mIdx, zIdx)
	}
}

func TestWritePalette_AllStyleFlags(t *testing.T) {
	// Arrange: Override with all style flags set to true.
	var colors [24]domain.Color
	for i := range colors {
		colors[i] = domain.Color{R: 0xaa, G: 0xbb, B: 0xcc}
	}

	color := domain.Color{R: 0x12, G: 0x34, B: 0x56}
	pal := &domain.Palette{
		System:  "base24",
		Name:    "AllStyles",
		Author:  "Test",
		Variant: "dark",
		Colors:  colors,
		Overrides: map[string]domain.TokenOverride{
			"test.token": {
				Color:         &color,
				Bold:          true,
				Italic:        true,
				Underline:     true,
				Undercurl:     true,
				Strikethrough: true,
			},
		},
	}

	// Act: Write palette.
	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// Assert: All style flags present in output.
	expectedFlags := []string{
		"bold: true",
		"italic: true",
		"underline: true",
		"undercurl: true",
		"strikethrough: true",
	}
	for _, flag := range expectedFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("output should contain %q, got:\n%s", flag, output)
		}
	}
}

func TestPalette_RoundTrip_WithOverrides(t *testing.T) {
	// Arrange: Palette with multiple overrides (color and style).
	var colors [24]domain.Color
	for i := range colors {
		colors[i] = domain.Color{R: uint8(i * 10), G: uint8(i * 11), B: uint8(i * 12)}
	}

	colorA := domain.Color{R: 0xaa, G: 0xbb, B: 0xcc}
	colorB := domain.Color{R: 0x11, G: 0x22, B: 0x33}
	original := &domain.Palette{
		System:  "base24",
		Name:    "RoundTrip",
		Author:  "Tester",
		Variant: "dark",
		Colors:  colors,
		Overrides: map[string]domain.TokenOverride{
			"syntax.keyword": {
				Color:  &colorA,
				Bold:   true,
				Italic: true,
			},
			"syntax.string": {
				Color:     &colorB,
				Underline: true,
			},
			"ui.cursor": {
				Bold:          true,
				Strikethrough: true,
			},
		},
	}

	// Act: Write to buffer, then parse back.
	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, original)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	parser := yaml.NewParser()
	parsed, err := parser.Parse(&buf)
	if err != nil {
		t.Fatalf("Parser.Parse() error = %v", err)
	}

	// Assert: Parsed palette has identical Overrides map.
	if len(parsed.Overrides) != len(original.Overrides) {
		t.Errorf("overrides count mismatch: got %d, want %d",
			len(parsed.Overrides), len(original.Overrides))
	}

	for key, origOvr := range original.Overrides {
		parsedOvr, ok := parsed.Overrides[key]
		if !ok {
			t.Errorf("missing override key %q after round-trip", key)
			continue
		}

		// Compare color.
		if origOvr.Color != nil {
			if parsedOvr.Color == nil {
				t.Errorf("override %q: color should not be nil", key)
			} else if *origOvr.Color != *parsedOvr.Color {
				t.Errorf("override %q: color mismatch: got %v, want %v",
					key, parsedOvr.Color, origOvr.Color)
			}
		} else if parsedOvr.Color != nil {
			t.Errorf("override %q: color should be nil", key)
		}

		// Compare style flags.
		if origOvr.Bold != parsedOvr.Bold {
			t.Errorf("override %q: bold mismatch: got %v, want %v",
				key, parsedOvr.Bold, origOvr.Bold)
		}
		if origOvr.Italic != parsedOvr.Italic {
			t.Errorf("override %q: italic mismatch: got %v, want %v",
				key, parsedOvr.Italic, origOvr.Italic)
		}
		if origOvr.Underline != parsedOvr.Underline {
			t.Errorf("override %q: underline mismatch: got %v, want %v",
				key, parsedOvr.Underline, origOvr.Underline)
		}
		if origOvr.Undercurl != parsedOvr.Undercurl {
			t.Errorf("override %q: undercurl mismatch: got %v, want %v",
				key, parsedOvr.Undercurl, origOvr.Undercurl)
		}
		if origOvr.Strikethrough != parsedOvr.Strikethrough {
			t.Errorf("override %q: strikethrough mismatch: got %v, want %v",
				key, parsedOvr.Strikethrough, origOvr.Strikethrough)
		}
	}
}

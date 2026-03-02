package yaml_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

const wantName = "Tokyo Night Dark"

// Verify Parser satisfies ports.PaletteParser at compile time.
var _ ports.PaletteParser = (*yaml.Parser)(nil)

const validBase24YAML = `system: "base24"
name: "Tokyo Night Dark"
author: "Michael Ball (base24 by curtbushko)"
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

const validBase16YAML = `system: "base16"
name: "Test Base16"
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

func TestParser_Parse_ValidBase24(t *testing.T) {
	parser := yaml.NewParser()
	reader := strings.NewReader(validBase24YAML)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pal.Name != wantName {
		t.Errorf("name = %q, want %q", pal.Name, wantName)
	}
	if pal.Variant != "dark" {
		t.Errorf("variant = %q, want %q", pal.Variant, "dark")
	}
	if pal.System != "base24" {
		t.Errorf("system = %q, want %q", pal.System, "base24")
	}

	// Verify specific colors
	wantBase00, _ := domain.ParseHex("1a1b26")
	if !pal.Base(0).Equal(wantBase00) {
		t.Errorf("base00 = %s, want %s", pal.Base(0).Hex(), wantBase00.Hex())
	}

	wantBase0D, _ := domain.ParseHex("7aa2f7")
	if !pal.Base(13).Equal(wantBase0D) {
		t.Errorf("base0D = %s, want %s", pal.Base(13).Hex(), wantBase0D.Hex())
	}

	wantBase17, _ := domain.ParseHex("c8acf8")
	if !pal.Base(23).Equal(wantBase17) {
		t.Errorf("base17 = %s, want %s", pal.Base(23).Hex(), wantBase17.Hex())
	}

	// All 24 slots should have valid colors (not IsNone)
	for i := 0; i < 24; i++ {
		if pal.Base(i).IsNone {
			t.Errorf("base%02d is IsNone, expected valid color", i)
		}
	}
}

func TestParser_Parse_ValidBase16(t *testing.T) {
	parser := yaml.NewParser()
	reader := strings.NewReader(validBase16YAML)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pal.System != "base16" {
		t.Errorf("system = %q, want %q", pal.System, "base16")
	}

	// Verify base16 fallbacks are applied:
	// base10 should equal base00
	if !pal.Base(16).Equal(pal.Base(0)) {
		t.Errorf("base10 = %s, want base00 = %s (fallback)", pal.Base(16).Hex(), pal.Base(0).Hex())
	}
	// base12 should equal base08
	if !pal.Base(18).Equal(pal.Base(8)) {
		t.Errorf("base12 = %s, want base08 = %s (fallback)", pal.Base(18).Hex(), pal.Base(8).Hex())
	}
	// base16 should equal base0D
	if !pal.Base(22).Equal(pal.Base(13)) {
		t.Errorf("base16 = %s, want base0D = %s (fallback)", pal.Base(22).Hex(), pal.Base(13).Hex())
	}

	// All 24 slots should be populated
	for i := 0; i < 24; i++ {
		if pal.Base(i).IsNone {
			t.Errorf("slot %d is IsNone after base16 fallback", i)
		}
	}
}

func TestParser_Parse_FromTestdataFile(t *testing.T) {
	f, err := os.Open("../../../testdata/tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("failed to open testdata: %v", err)
	}
	defer func() { _ = f.Close() }()

	parser := yaml.NewParser()
	pal, err := parser.Parse(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pal.Name != wantName {
		t.Errorf("name = %q, want %q", pal.Name, wantName)
	}

	wantBase00, _ := domain.ParseHex("1a1b26")
	if !pal.Base(0).Equal(wantBase00) {
		t.Errorf("base00 = %s, want %s", pal.Base(0).Hex(), wantBase00.Hex())
	}

	wantBase17, _ := domain.ParseHex("c8acf8")
	if !pal.Base(23).Equal(wantBase17) {
		t.Errorf("base17 = %s, want %s", pal.Base(23).Hex(), wantBase17.Hex())
	}
}

func TestParser_Parse_InvalidYAML(t *testing.T) {
	parser := yaml.NewParser()
	reader := strings.NewReader("not: [valid: yaml")

	_, err := parser.Parse(reader)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestParser_Parse_MissingPalette(t *testing.T) {
	parser := yaml.NewParser()
	yamlContent := `system: "base24"
name: "Missing Palette"
author: "Test"
variant: "dark"
`
	reader := strings.NewReader(yamlContent)

	_, err := parser.Parse(reader)
	if err == nil {
		t.Fatal("expected error for missing palette section, got nil")
	}

	var parseErr *domain.ParseError
	if !errors.As(err, &parseErr) {
		t.Errorf("expected *domain.ParseError, got %T: %v", err, err)
	}
}

func TestParser_Parse_InvalidHex(t *testing.T) {
	parser := yaml.NewParser()
	yamlContent := `system: "base24"
name: "Bad Hex"
author: "Test"
variant: "dark"
palette:
  base00: "zzzzzz"
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
	reader := strings.NewReader(yamlContent)

	_, err := parser.Parse(reader)
	if err == nil {
		t.Fatal("expected error for invalid hex, got nil")
	}

	var parseErr *domain.ParseError
	if !errors.As(err, &parseErr) {
		t.Errorf("expected *domain.ParseError, got %T: %v", err, err)
	}
}

func TestParser_Parse_BytesReader(t *testing.T) {
	parser := yaml.NewParser()
	reader := bytes.NewReader([]byte(validBase24YAML))

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pal.Name != wantName {
		t.Errorf("name = %q, want %q", pal.Name, wantName)
	}

	// Verify a specific color to prove it actually parsed
	wantBase0D, _ := domain.ParseHex("7aa2f7")
	if !pal.Base(13).Equal(wantBase0D) {
		t.Errorf("base0D = %s, want %s", pal.Base(13).Hex(), wantBase0D.Hex())
	}
}

// --- Override Parsing Tests ---

const validBase24WithOverridesYAML = `system: "base24"
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
overrides:
  syntax.keyword:
    color: "ff00ff"
    italic: true
`

func TestParser_Parse_WithOverrides(t *testing.T) {
	parser := yaml.NewParser()
	reader := strings.NewReader(validBase24WithOverridesYAML)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify overrides map is populated
	if pal.Overrides == nil {
		t.Fatal("expected Overrides map to be non-nil")
	}

	override, ok := pal.Overrides["syntax.keyword"]
	if !ok {
		t.Fatal("expected override for 'syntax.keyword' to exist")
	}

	// Verify color is set correctly
	if !override.HasColor() {
		t.Error("expected override to have color")
	}
	wantColor, _ := domain.ParseHex("ff00ff")
	if !override.Color.Equal(wantColor) {
		t.Errorf("override color = %s, want %s", override.Color.Hex(), wantColor.Hex())
	}

	// Verify italic flag
	if !override.Italic {
		t.Error("expected override Italic to be true")
	}

	// Verify other flags are false
	if override.Bold {
		t.Error("expected override Bold to be false")
	}
}

func TestParser_Parse_OverrideColorOnly(t *testing.T) {
	yamlContent := `system: "base24"
name: "Test"
author: "Test"
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
overrides:
  surface.background:
    color: "000000"
`
	parser := yaml.NewParser()
	reader := strings.NewReader(yamlContent)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	override, ok := pal.Overrides["surface.background"]
	if !ok {
		t.Fatal("expected override for 'surface.background'")
	}

	// Verify color is set
	if !override.HasColor() {
		t.Error("expected override to have color")
	}
	wantColor, _ := domain.ParseHex("000000")
	if !override.Color.Equal(wantColor) {
		t.Errorf("override color = %s, want %s", override.Color.Hex(), wantColor.Hex())
	}

	// Verify style flags are false
	if override.Bold {
		t.Error("expected Bold to be false")
	}
	if override.Italic {
		t.Error("expected Italic to be false")
	}
	if override.Underline {
		t.Error("expected Underline to be false")
	}
	if override.Undercurl {
		t.Error("expected Undercurl to be false")
	}
	if override.Strikethrough {
		t.Error("expected Strikethrough to be false")
	}
}

func TestParser_Parse_OverrideStyleOnly(t *testing.T) {
	yamlContent := `system: "base24"
name: "Test"
author: "Test"
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
overrides:
  syntax.comment:
    bold: true
`
	parser := yaml.NewParser()
	reader := strings.NewReader(yamlContent)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	override, ok := pal.Overrides["syntax.comment"]
	if !ok {
		t.Fatal("expected override for 'syntax.comment'")
	}

	// Verify no color is set
	if override.HasColor() {
		t.Error("expected override to not have color")
	}

	// Verify Bold is true
	if !override.Bold {
		t.Error("expected Bold to be true")
	}
}

func TestParser_Parse_InvalidOverrideColor(t *testing.T) {
	yamlContent := `system: "base24"
name: "Test"
author: "Test"
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
overrides:
  syntax.keyword:
    color: "zzzzzz"
`
	parser := yaml.NewParser()
	reader := strings.NewReader(yamlContent)

	_, err := parser.Parse(reader)
	if err == nil {
		t.Fatal("expected error for invalid override hex color, got nil")
	}

	var parseErr *domain.ParseError
	if !errors.As(err, &parseErr) {
		t.Errorf("expected *domain.ParseError, got %T: %v", err, err)
	}

	// Verify error mentions the override field
	if parseErr != nil && !strings.Contains(parseErr.Field, "overrides") {
		t.Errorf("expected error field to mention 'overrides', got %q", parseErr.Field)
	}
}

func TestParser_Parse_NoOverrides(t *testing.T) {
	parser := yaml.NewParser()
	reader := strings.NewReader(validBase24YAML)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Overrides is nil or empty
	if len(pal.Overrides) != 0 {
		t.Errorf("expected Overrides to be nil or empty, got %d entries", len(pal.Overrides))
	}
}

func TestParser_Parse_MultipleOverrides(t *testing.T) {
	yamlContent := `system: "base24"
name: "Test"
author: "Test"
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
overrides:
  syntax.keyword:
    color: "ff00ff"
    italic: true
  surface.background:
    color: "000000"
  status.error:
    color: "ff0000"
    bold: true
    underline: true
`
	parser := yaml.NewParser()
	reader := strings.NewReader(yamlContent)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all three overrides are present
	if len(pal.Overrides) != 3 {
		t.Errorf("expected 3 overrides, got %d", len(pal.Overrides))
	}

	expectedTokens := []string{"syntax.keyword", "surface.background", "status.error"}
	for _, token := range expectedTokens {
		if _, ok := pal.Overrides[token]; !ok {
			t.Errorf("expected override for %q to exist", token)
		}
	}

	// Verify status.error has all its flags
	statusErr := pal.Overrides["status.error"]
	if !statusErr.Bold {
		t.Error("expected status.error Bold to be true")
	}
	if !statusErr.Underline {
		t.Error("expected status.error Underline to be true")
	}
}

func TestParser_Parse_OverrideAllStyleFlags(t *testing.T) {
	yamlContent := `system: "base24"
name: "Test"
author: "Test"
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
overrides:
  syntax.keyword:
    color: "aabbcc"
    bold: true
    italic: true
    underline: true
    undercurl: true
    strikethrough: true
`
	parser := yaml.NewParser()
	reader := strings.NewReader(yamlContent)

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	override, ok := pal.Overrides["syntax.keyword"]
	if !ok {
		t.Fatal("expected override for 'syntax.keyword'")
	}

	// Verify color
	if !override.HasColor() {
		t.Error("expected override to have color")
	}
	wantColor, _ := domain.ParseHex("aabbcc")
	if !override.Color.Equal(wantColor) {
		t.Errorf("override color = %s, want %s", override.Color.Hex(), wantColor.Hex())
	}

	// Verify all style flags are true
	if !override.Bold {
		t.Error("expected Bold to be true")
	}
	if !override.Italic {
		t.Error("expected Italic to be true")
	}
	if !override.Underline {
		t.Error("expected Underline to be true")
	}
	if !override.Undercurl {
		t.Error("expected Undercurl to be true")
	}
	if !override.Strikethrough {
		t.Error("expected Strikethrough to be true")
	}
}

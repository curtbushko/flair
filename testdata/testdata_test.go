package testdata_test

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

// tokyoNightDarkPalette is the expected palette structure parsed from the YAML.
type paletteFile struct {
	System  string            `yaml:"system"`
	Name    string            `yaml:"name"`
	Author  string            `yaml:"author"`
	Variant string            `yaml:"variant"`
	Palette map[string]string `yaml:"palette"`
}

// expectedPalette contains the reference Tokyo Night Dark base24 palette
// values from PLAN.md.
var expectedPalette = map[string]string{
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

func TestTokyoNightDark_ValidYAML(t *testing.T) {
	data, err := os.ReadFile("tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("failed to read testdata/tokyo-night-dark.yaml: %v", err)
	}

	var pf paletteFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}
}

func TestTokyoNightDark_AllSlots(t *testing.T) {
	data, err := os.ReadFile("tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("failed to read testdata/tokyo-night-dark.yaml: %v", err)
	}

	var pf paletteFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	// All 24 base slots must be present: base00 through base17
	allSlots := []string{
		"base00", "base01", "base02", "base03",
		"base04", "base05", "base06", "base07",
		"base08", "base09", "base0A", "base0B",
		"base0C", "base0D", "base0E", "base0F",
		"base10", "base11", "base12", "base13",
		"base14", "base15", "base16", "base17",
	}

	if len(pf.Palette) != 24 {
		t.Errorf("palette has %d slots, want 24", len(pf.Palette))
	}

	for _, slot := range allSlots {
		if _, ok := pf.Palette[slot]; !ok {
			t.Errorf("missing palette slot: %s", slot)
		}
	}
}

func TestTokyoNightDark_ValuesMatch(t *testing.T) {
	data, err := os.ReadFile("tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("failed to read testdata/tokyo-night-dark.yaml: %v", err)
	}

	var pf paletteFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	for slot, want := range expectedPalette {
		got, ok := pf.Palette[slot]
		if !ok {
			t.Errorf("missing slot %s", slot)
			continue
		}
		if got != want {
			t.Errorf("slot %s = %q, want %q", slot, got, want)
		}
	}
}

func TestTokyoNightDark_Metadata(t *testing.T) {
	data, err := os.ReadFile("tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("failed to read testdata/tokyo-night-dark.yaml: %v", err)
	}

	var pf paletteFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	if pf.System != "base24" {
		t.Errorf("system = %q, want %q", pf.System, "base24")
	}
	if pf.Name != "Tokyo Night Dark" {
		t.Errorf("name = %q, want %q", pf.Name, "Tokyo Night Dark")
	}
	if pf.Variant != "dark" {
		t.Errorf("variant = %q, want %q", pf.Variant, "dark")
	}
}

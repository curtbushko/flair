package flair

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// slotNames lists the 24 base24 slot names in index order.
var slotNames = [24]string{
	"base00", "base01", "base02", "base03",
	"base04", "base05", "base06", "base07",
	"base08", "base09", "base0A", "base0B",
	"base0C", "base0D", "base0E", "base0F",
	"base10", "base11", "base12", "base13",
	"base14", "base15", "base16", "base17",
}

// paletteYAML is the internal representation of the YAML structure.
type paletteYAML struct {
	System  string            `yaml:"system"`
	Name    string            `yaml:"name"`
	Author  string            `yaml:"author"`
	Variant string            `yaml:"variant"`
	Palette map[string]string `yaml:"palette"`
}

// ParsePalette parses a base24 palette from YAML data read from r.
//
// The expected YAML format is:
//
//	name: "Theme Name"
//	author: "Author Name"
//	variant: "dark"
//	palette:
//	  base00: "#1a1b26"
//	  base01: "#16161e"
//	  ...
//	  base17: "#d18616"
//
// ParsePalette validates that all 24 base colors (base00-base17) are present
// and that all hex values are valid 3 or 6-digit hex colors.
//
// Returns an error if:
//   - The YAML is malformed
//   - Any required base color is missing
//   - Any hex value is invalid
//
// Example:
//
//	file, err := os.Open("palette.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	palette, err := flair.ParsePalette(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Loaded:", palette.Name())
func ParsePalette(r io.Reader) (*Palette, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read palette: %w", err)
	}

	var py paletteYAML
	if err := yaml.Unmarshal(data, &py); err != nil {
		return nil, fmt.Errorf("parse palette YAML: %w", err)
	}

	pal := &Palette{
		name:    py.Name,
		author:  py.Author,
		variant: py.Variant,
	}

	// Parse and validate all 24 color slots.
	for i, slotName := range slotNames {
		hexVal, ok := py.Palette[slotName]
		if !ok {
			return nil, fmt.Errorf("missing required color slot %s", slotName)
		}

		c, err := ParseHex(hexVal)
		if err != nil {
			return nil, fmt.Errorf("invalid hex color for %s: %w", slotName, err)
		}
		pal.colors[i] = c
	}

	return pal, nil
}

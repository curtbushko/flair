// Package yaml provides a YAML palette parser adapter that reads
// tinted-theming common format YAML from an io.Reader and returns
// a domain.Palette. It implements ports.PaletteParser.
package yaml

import (
	"fmt"
	"io"

	"github.com/curtbushko/flair/internal/domain"

	yamlv3 "gopkg.in/yaml.v3"
)

// overrideYAML is the intermediate structure for parsing a token override
// from the YAML overrides section.
type overrideYAML struct {
	Color         string `yaml:"color,omitempty"`
	Bold          bool   `yaml:"bold,omitempty"`
	Italic        bool   `yaml:"italic,omitempty"`
	Underline     bool   `yaml:"underline,omitempty"`
	Undercurl     bool   `yaml:"undercurl,omitempty"`
	Strikethrough bool   `yaml:"strikethrough,omitempty"`
}

// paletteYAML is the intermediate structure matching the tinted-theming
// common YAML format (spec 0.11+), extended with optional token overrides.
type paletteYAML struct {
	System    string                  `yaml:"system"`
	Name      string                  `yaml:"name"`
	Author    string                  `yaml:"author"`
	Variant   string                  `yaml:"variant"`
	Palette   map[string]string       `yaml:"palette"`
	Overrides map[string]overrideYAML `yaml:"overrides,omitempty"`
}

// Parser implements ports.PaletteParser by decoding YAML from an io.Reader
// and delegating palette construction to domain.NewPalette.
type Parser struct{}

// NewParser returns a new YAML palette parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads palette YAML from r and returns a domain.Palette.
// It expects the tinted-theming common format with system, name, author,
// variant, and palette fields. Returns a *domain.ParseError for missing
// fields, invalid hex values, or malformed YAML.
func (p *Parser) Parse(r io.Reader) (*domain.Palette, error) {
	var raw paletteYAML
	decoder := yamlv3.NewDecoder(r)
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	if len(raw.Palette) == 0 {
		return nil, &domain.ParseError{
			Field:   "palette",
			Message: "missing or empty palette section",
		}
	}

	// Determine system: default to base24 if enough colors provided, base16 otherwise.
	system := raw.System
	if system == "" {
		if len(raw.Palette) >= 24 {
			system = "base24"
		} else {
			system = "base16"
		}
	}

	pal, err := domain.NewPalette(raw.Name, raw.Author, raw.Variant, system, raw.Palette)
	if err != nil {
		return nil, err
	}

	// Parse overrides if present
	if len(raw.Overrides) > 0 {
		overrides, err := parseOverrides(raw.Overrides)
		if err != nil {
			return nil, err
		}
		pal.Overrides = overrides
	}

	return pal, nil
}

// parseOverrides converts the YAML override map to domain.TokenOverride map.
// Returns a *domain.ParseError if any override has an invalid hex color.
func parseOverrides(rawOverrides map[string]overrideYAML) (map[string]domain.TokenOverride, error) {
	overrides := make(map[string]domain.TokenOverride, len(rawOverrides))

	for tokenPath, raw := range rawOverrides {
		override, err := domain.NewTokenOverride(
			raw.Color,
			raw.Bold,
			raw.Italic,
			raw.Underline,
			raw.Undercurl,
			raw.Strikethrough,
		)
		if err != nil {
			return nil, &domain.ParseError{
				Field:   fmt.Sprintf("overrides.%s.color", tokenPath),
				Message: fmt.Sprintf("invalid hex color %q for override %s", raw.Color, tokenPath),
				Cause:   err,
			}
		}
		overrides[tokenPath] = *override
	}

	return overrides, nil
}

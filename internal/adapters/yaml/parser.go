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

// paletteYAML is the intermediate structure matching the tinted-theming
// common YAML format (spec 0.11+).
type paletteYAML struct {
	System  string            `yaml:"system"`
	Name    string            `yaml:"name"`
	Author  string            `yaml:"author"`
	Variant string            `yaml:"variant"`
	Palette map[string]string `yaml:"palette"`
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

	return pal, nil
}

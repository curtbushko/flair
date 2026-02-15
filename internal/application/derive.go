// Package application contains use cases that orchestrate domain logic
// through port interfaces. It depends only on ports and domain, never
// on concrete adapters (hexagonal architecture).
package application

import (
	"fmt"
	"io"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// DeriveThemeUseCase orchestrates palette parsing and token derivation
// to produce a fully resolved theme. It depends on port interfaces,
// keeping the application layer adapter-agnostic.
type DeriveThemeUseCase struct {
	parser  ports.PaletteParser
	deriver ports.TokenDeriver
}

// NewDeriveThemeUseCase returns a new DeriveThemeUseCase wired to the
// given parser and deriver ports.
func NewDeriveThemeUseCase(parser ports.PaletteParser, deriver ports.TokenDeriver) *DeriveThemeUseCase {
	return &DeriveThemeUseCase{
		parser:  parser,
		deriver: deriver,
	}
}

// Execute reads palette YAML from r, parses it into a domain Palette,
// derives the full semantic token set, and returns a ResolvedTheme.
func (uc *DeriveThemeUseCase) Execute(r io.Reader) (*domain.ResolvedTheme, error) {
	palette, err := uc.parser.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parse palette: %w", err)
	}

	tokens := uc.deriver.Derive(palette)

	return &domain.ResolvedTheme{
		Name:    palette.Name,
		Variant: palette.Variant,
		Palette: palette,
		Tokens:  tokens,
	}, nil
}

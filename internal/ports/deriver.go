package ports

import "github.com/curtbushko/flair/internal/domain"

// TokenDeriver derives the full semantic token set from a base24 palette.
type TokenDeriver interface {
	Derive(p *domain.Palette) *domain.TokenSet
}

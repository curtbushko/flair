package ports

import "github.com/curtbushko/flair/internal/domain"

// Tokenizer transforms a base24 palette into a semantic token set.
type Tokenizer interface {
	Tokenize(p *domain.Palette) *domain.TokenSet
}

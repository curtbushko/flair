package flair

import (
	"github.com/curtbushko/flair/internal/adapters/tokenizer"
	"github.com/curtbushko/flair/internal/domain"
)

// Tokenize derives semantic color tokens from a base24 palette.
//
// Tokenize transforms a raw 24-color palette into a Theme containing
// semantic tokens organized into categories:
//
//   - Surface: Background colors for UI elements
//   - Text: Foreground colors for text content
//   - Status: Colors for error, warning, success, info states
//   - Diff: Colors for version control diffs
//   - Syntax: Colors for code syntax highlighting
//   - Markup: Colors for documentation and markdown
//   - Comment: Colors for LSP diagnostic-style comments
//   - Accent/Border: UI accent and border colors
//   - Git: Git status indicator colors
//   - Terminal: 16 ANSI terminal colors
//   - Statusline: Colors for statusline sections
//
// This is a thin façade over the internal tokenizer adapter; derivation
// rules are defined in one place (internal/adapters/tokenizer). Style
// flags (bold, italic, etc.) attached to internal tokens are dropped
// because the public Theme API exposes color values only.
//
// Tokenize returns nil if the palette is nil.
func Tokenize(p *Palette) *Theme {
	if p == nil {
		return nil
	}

	dp := toDomainPalette(p)
	ts := tokenizer.New().Tokenize(dp)

	colors := make(map[string]Color, ts.Len())
	for _, path := range ts.Paths() {
		tok, _ := ts.Get(path)
		colors[path] = fromDomainColor(tok.Color)
	}

	return NewTheme(p.Name(), p.Variant(), colors)
}

// toDomainPalette converts a public Palette into a domain.Palette that the
// internal tokenizer can consume. The public Palette does not carry a system
// tag or overrides, so we assume base24 (its 24 slots are already resolved)
// and pass no overrides.
func toDomainPalette(p *Palette) *domain.Palette {
	dp := &domain.Palette{
		Name:    p.Name(),
		Author:  p.Author(),
		Variant: p.Variant(),
		System:  "base24",
		Slug:    p.Name() + "-" + p.Variant(),
	}
	for i := 0; i < 24; i++ {
		c := p.Base(i)
		if c == nil {
			continue
		}
		dp.Colors[i] = domain.Color{R: c.R, G: c.G, B: c.B}
	}
	return dp
}

// fromDomainColor converts a domain.Color to the public Color type,
// discarding the internal IsNone sentinel flag.
func fromDomainColor(c domain.Color) Color {
	return Color{R: c.R, G: c.G, B: c.B}
}

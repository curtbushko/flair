// Package tokenizer implements the default token derivation rules that transform
// a base24 palette into a semantic token set.
package tokenizer

import (
	"github.com/curtbushko/flair/internal/domain"
)

// DefaultTokenizer implements ports.Tokenizer using the standard derivation
// rules defined in the PLAN.md token inventory.
type DefaultTokenizer struct{}

// New returns a new DefaultTokenizer.
func New() *DefaultTokenizer {
	return &DefaultTokenizer{}
}

// Tokenize transforms a base24 palette into a complete semantic token set.
// After derivation, any overrides defined in the palette are applied.
// Overridden colors replace derived colors; overridden style flags are
// merged (OR'd) with derived styles.
func (d *DefaultTokenizer) Tokenize(p *domain.Palette) *domain.TokenSet {
	ts := domain.NewTokenSet()
	deriveSurface(p, ts)
	deriveText(p, ts)
	deriveStatus(p, ts)
	deriveDiff(p, ts)
	deriveSyntax(p, ts)
	deriveMarkup(p, ts)
	deriveComment(p, ts)
	deriveAccentBorderState(p, ts)
	deriveGit(p, ts)
	deriveTerminal(p, ts)
	deriveStatusline(p, ts)
	applyOverrides(p, ts)
	return ts
}

// applyOverrides applies any token overrides from the palette to the token set.
// For each override path that exists in the token set, the override is applied
// using TokenOverride.Apply(), which replaces color (if set) and merges style
// flags with OR logic. Invalid/unknown paths are silently ignored.
func applyOverrides(p *domain.Palette, ts *domain.TokenSet) {
	if p.Overrides == nil {
		return
	}

	for path, override := range p.Overrides {
		tok, ok := ts.Get(path)
		if !ok {
			// Path doesn't exist in token set; silently ignore
			continue
		}

		// Apply the override to the existing token
		updated := override.Apply(tok)
		ts.Set(path, updated)
	}
}

// deriveSurface derives the 11 surface tokens from the palette.
func deriveSurface(p *domain.Palette, ts *domain.TokenSet) {
	// Direct palette mappings
	ts.Set("surface.background", domain.Token{Color: p.Base(0x00)})
	ts.Set("surface.background.raised", domain.Token{Color: p.Base(0x01)})
	ts.Set("surface.background.sunken", domain.Token{Color: p.Base(0x10)})
	ts.Set("surface.background.darkest", domain.Token{Color: p.Base(0x11)})
	ts.Set("surface.background.highlight", domain.Token{Color: p.Base(0x02)})

	// Blended tokens
	ts.Set("surface.background.selection", domain.Token{
		Color: domain.BlendBg(p.Base(0x0D), p.Base(0x00), 0.30),
	})
	ts.Set("surface.background.search", domain.Token{
		Color: domain.BlendBg(p.Base(0x0A), p.Base(0x00), 0.30),
	})

	// base10 aliases
	base10 := p.Base(0x10)
	ts.Set("surface.background.overlay", domain.Token{Color: base10})
	ts.Set("surface.background.popup", domain.Token{Color: base10})
	ts.Set("surface.background.sidebar", domain.Token{Color: base10})
	ts.Set("surface.background.statusbar", domain.Token{Color: base10})
}

// deriveText derives the 7 text tokens from the palette.
func deriveText(p *domain.Palette, ts *domain.TokenSet) {
	ts.Set("text.primary", domain.Token{Color: p.Base(0x05)})
	ts.Set("text.secondary", domain.Token{Color: p.Base(0x04)})
	ts.Set("text.muted", domain.Token{Color: p.Base(0x03)})
	ts.Set("text.subtle", domain.Token{
		Color: domain.BlendBg(p.Base(0x03), p.Base(0x00), 0.50),
	})
	ts.Set("text.inverse", domain.Token{Color: p.Base(0x00)})
	ts.Set("text.overlay", domain.Token{Color: p.Base(0x06)})
	ts.Set("text.sidebar", domain.Token{Color: p.Base(0x04)})
}

// deriveStatus derives the 6 status tokens from the palette.
func deriveStatus(p *domain.Palette, ts *domain.TokenSet) {
	ts.Set("status.error", domain.Token{Color: p.Base(0x12)})
	ts.Set("status.warning", domain.Token{Color: p.Base(0x13)})
	ts.Set("status.success", domain.Token{Color: p.Base(0x14)})
	ts.Set("status.info", domain.Token{Color: p.Base(0x14)})
	ts.Set("status.hint", domain.Token{Color: p.Base(0x09)})
	ts.Set("status.todo", domain.Token{Color: p.Base(0x0D)})
}

// deriveDiff derives the 10 diff tokens from the palette.
func deriveDiff(p *domain.Palette, ts *domain.TokenSet) {
	ts.Set("diff.added.fg", domain.Token{Color: p.Base(0x05)})
	ts.Set("diff.added.bg", domain.Token{
		Color: domain.BlendBg(p.Base(0x14), p.Base(0x00), 0.25),
	})
	ts.Set("diff.added.sign", domain.Token{Color: p.Base(0x14)})
	ts.Set("diff.deleted.fg", domain.Token{Color: p.Base(0x05)})
	ts.Set("diff.deleted.bg", domain.Token{
		Color: domain.BlendBg(p.Base(0x12), p.Base(0x00), 0.25),
	})
	ts.Set("diff.deleted.sign", domain.Token{Color: p.Base(0x12)})
	ts.Set("diff.changed.fg", domain.Token{Color: p.Base(0x05)})
	ts.Set("diff.changed.bg", domain.Token{
		Color: domain.BlendBg(p.Base(0x16), p.Base(0x00), 0.15),
	})
	ts.Set("diff.changed.sign", domain.Token{Color: p.Base(0x16)})
	ts.Set("diff.ignored", domain.Token{Color: p.Base(0x03)})
}

// deriveSyntax derives the 25 syntax tokens from the palette.
// syntax.comment includes the Italic style flag.
// syntax.deprecated includes the Strikethrough style flag.
func deriveSyntax(p *domain.Palette, ts *domain.TokenSet) {
	// Core syntax tokens
	ts.Set("syntax.keyword", domain.Token{Color: p.Base(0x0E)})
	ts.Set("syntax.string", domain.Token{Color: p.Base(0x0B)})
	ts.Set("syntax.function", domain.Token{Color: p.Base(0x0D)})
	ts.Set("syntax.comment", domain.Token{Color: p.Base(0x03), Italic: true})
	ts.Set("syntax.variable", domain.Token{Color: p.Base(0x05)})
	ts.Set("syntax.constant", domain.Token{Color: p.Base(0x09)})
	ts.Set("syntax.operator", domain.Token{Color: p.Base(0x16)})
	ts.Set("syntax.type", domain.Token{Color: p.Base(0x0A)})
	ts.Set("syntax.number", domain.Token{Color: p.Base(0x09)})
	ts.Set("syntax.tag", domain.Token{Color: p.Base(0x08)})
	ts.Set("syntax.property", domain.Token{Color: p.Base(0x15)}) // bright cyan (was 0x0B)
	ts.Set("syntax.parameter", domain.Token{Color: p.Base(0x0A)})
	ts.Set("syntax.regexp", domain.Token{Color: p.Base(0x0C)})
	ts.Set("syntax.escape", domain.Token{Color: p.Base(0x17)}) // bright purple (was 0x0E)
	ts.Set("syntax.constructor", domain.Token{Color: p.Base(0x17)})

	// New tokens for color variety
	ts.Set("syntax.boolean", domain.Token{Color: p.Base(0x09)})
	ts.Set("syntax.function.builtin", domain.Token{Color: p.Base(0x0C)})
	ts.Set("syntax.type.builtin", domain.Token{Color: p.Base(0x0C)})
	ts.Set("syntax.module", domain.Token{Color: p.Base(0x0A)})
	ts.Set("syntax.module.builtin", domain.Token{Color: p.Base(0x05)})
	ts.Set("syntax.string.documentation", domain.Token{Color: p.Base(0x0A)})
	ts.Set("syntax.label", domain.Token{Color: p.Base(0x0D)})
	ts.Set("syntax.punctuation", domain.Token{Color: p.Base(0x16)})
	ts.Set("syntax.deprecated", domain.Token{Color: p.Base(0x0F), Strikethrough: true})
	ts.Set("syntax.macro", domain.Token{Color: p.Base(0x13)})
}

// deriveMarkup derives the markup semantic tokens from the palette.
// Includes per-heading-level colors for markdown.
func deriveMarkup(p *domain.Palette, ts *domain.TokenSet) {
	ts.Set("markup.heading", domain.Token{Color: p.Base(0x0D), Bold: true})
	// Per-level markdown heading colors
	ts.Set("markup.heading.1", domain.Token{Color: p.Base(0x0D), Bold: true}) // blue
	ts.Set("markup.heading.2", domain.Token{Color: p.Base(0x0A), Bold: true}) // yellow
	ts.Set("markup.heading.3", domain.Token{Color: p.Base(0x14), Bold: true}) // green
	ts.Set("markup.heading.4", domain.Token{Color: p.Base(0x08), Bold: true}) // red
	ts.Set("markup.heading.5", domain.Token{Color: p.Base(0x17), Bold: true}) // purple
	ts.Set("markup.heading.6", domain.Token{Color: p.Base(0x0D), Bold: true}) // blue
	ts.Set("markup.link", domain.Token{Color: p.Base(0x0C)})
	ts.Set("markup.code", domain.Token{Color: p.Base(0x0B)})
	ts.Set("markup.bold", domain.Token{Color: p.Base(0x05), Bold: true})
	ts.Set("markup.italic", domain.Token{Color: p.Base(0x05), Italic: true})
	ts.Set("markup.strikethrough", domain.Token{Color: p.Base(0x03), Strikethrough: true})
	ts.Set("markup.quote", domain.Token{Color: p.Base(0x03), Italic: true})
	ts.Set("markup.list.bullet", domain.Token{Color: p.Base(0x09)})
	ts.Set("markup.list.checked", domain.Token{Color: p.Base(0x0B)})
	ts.Set("markup.list.unchecked", domain.Token{Color: p.Base(0x0D)})
}

// deriveComment derives comment variant tokens from the palette.
func deriveComment(p *domain.Palette, ts *domain.TokenSet) {
	ts.Set("comment.error", domain.Token{Color: p.Base(0x12)})   // bright red
	ts.Set("comment.warning", domain.Token{Color: p.Base(0x0A)}) // yellow
	ts.Set("comment.info", domain.Token{Color: p.Base(0x16)})    // bright blue
	ts.Set("comment.hint", domain.Token{Color: p.Base(0x09)})    // orange
	ts.Set("comment.note", domain.Token{Color: p.Base(0x09)})    // orange
	ts.Set("comment.todo", domain.Token{Color: p.Base(0x0D)})    // blue
}

// deriveAccentBorderState derives the 11 accent, border, scrollbar, and state
// tokens from the palette. Alias tokens (state.hover, state.disabled.fg)
// produce the same color as their referenced semantic tokens.
func deriveAccentBorderState(p *domain.Palette, ts *domain.TokenSet) {
	// Accent (3 tokens)
	ts.Set("accent.primary", domain.Token{Color: p.Base(0x0D)})
	ts.Set("accent.secondary", domain.Token{Color: p.Base(0x0E)})
	ts.Set("accent.foreground", domain.Token{Color: p.Base(0x00)})

	// Border (3 tokens)
	ts.Set("border.default", domain.Token{Color: p.Base(0x09)}) // orange
	ts.Set("border.focus", domain.Token{
		Color: domain.BlendBg(p.Base(0x0D), p.Base(0x00), 0.70),
	})
	ts.Set("border.muted", domain.Token{Color: p.Base(0x01)})

	// Scrollbar (2 tokens)
	ts.Set("scrollbar.thumb", domain.Token{
		Color: domain.BlendBg(p.Base(0x03), p.Base(0x00), 0.40),
	})
	ts.Set("scrollbar.track", domain.Token{Color: p.Base(0x01)})

	// State (3 tokens)
	// state.hover aliases surface.background.highlight (base02)
	ts.Set("state.hover", domain.Token{Color: p.Base(0x02)})
	ts.Set("state.active", domain.Token{
		Color: domain.BlendBg(p.Base(0x0D), p.Base(0x00), 0.20),
	})
	// state.disabled.fg aliases text.muted (base03)
	ts.Set("state.disabled.fg", domain.Token{Color: p.Base(0x03)})
}

// deriveGit derives the 4 git status tokens from the palette.
func deriveGit(p *domain.Palette, ts *domain.TokenSet) {
	ts.Set("git.added", domain.Token{Color: p.Base(0x0B)})
	ts.Set("git.modified", domain.Token{Color: p.Base(0x0D)})
	ts.Set("git.deleted", domain.Token{Color: p.Base(0x08)})
	ts.Set("git.ignored", domain.Token{Color: p.Base(0x03)})
}

// deriveTerminal derives the 16 terminal ANSI color tokens from the palette.
// Uses a table-driven approach mapping ANSI color names to base24 slots.
func deriveTerminal(p *domain.Palette, ts *domain.TokenSet) {
	type ansiMapping struct {
		name string
		base int
	}

	mappings := []ansiMapping{
		// Normal colors
		{"terminal.black", 0x01},
		{"terminal.red", 0x08},
		{"terminal.green", 0x0B},
		{"terminal.yellow", 0x0A},
		{"terminal.blue", 0x0D},
		{"terminal.magenta", 0x0E},
		{"terminal.cyan", 0x0C},
		{"terminal.white", 0x05},
		// Bright colors
		{"terminal.brblack", 0x03},
		{"terminal.brred", 0x12},
		{"terminal.brgreen", 0x14},
		{"terminal.bryellow", 0x13},
		{"terminal.brblue", 0x16},
		{"terminal.brmagenta", 0x17},
		{"terminal.brcyan", 0x15},
		{"terminal.brwhite", 0x07},
	}

	for _, m := range mappings {
		ts.Set(m.name, domain.Token{Color: p.Base(m.base)})
	}
}

// deriveStatusline derives the 6 statusline tokens from the palette.
// These provide foreground and background colors for statusline sections A, B, C.
// Section A: mode indicator (bold bg with dark fg)
// Section B: branch/info (subtle bg with light text)
// Section C: filename area (darkest bg with light text)
func deriveStatusline(p *domain.Palette, ts *domain.TokenSet) {
	ts.Set("statusline.a.bg", domain.Token{Color: p.Base(0x04)}) // secondary text color as bg
	ts.Set("statusline.a.fg", domain.Token{Color: p.Base(0x00)}) // darkest for contrast
	ts.Set("statusline.b.bg", domain.Token{Color: p.Base(0x10)}) // sunken surface
	ts.Set("statusline.b.fg", domain.Token{Color: p.Base(0x05)}) // primary text
	ts.Set("statusline.c.bg", domain.Token{Color: p.Base(0x01)}) // raised surface
	ts.Set("statusline.c.fg", domain.Token{Color: p.Base(0x05)}) // primary text
}

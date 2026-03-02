package application

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// TokensReader reads a tokens.yaml from a reader and returns a TokenSet.
// The composition root wires this to fileio.ReadTokens.
type TokensReader func(r io.Reader) (*domain.TokenSet, error)

// PreviewThemeUseCase reads palette.yaml and tokens.yaml from a theme
// directory and renders an ANSI-colored preview to an io.Writer.
// If tokens.yaml doesn't exist, it tokenizes palette.yaml
// or uses a built-in palette.
type PreviewThemeUseCase struct {
	store        ports.ThemeStore
	parser       ports.PaletteParser
	tokensReader TokensReader
	tokenizer    ports.Tokenizer
	builtins     ports.PaletteSource
}

// NewPreviewThemeUseCase returns a new PreviewThemeUseCase wired to the
// given store, palette parser, tokens reader, tokenizer, and builtins source.
func NewPreviewThemeUseCase(
	store ports.ThemeStore,
	parser ports.PaletteParser,
	tr TokensReader,
	tokenizer ports.Tokenizer,
	builtins ports.PaletteSource,
) *PreviewThemeUseCase {
	return &PreviewThemeUseCase{
		store:        store,
		parser:       parser,
		tokensReader: tr,
		tokenizer:    tokenizer,
		builtins:     builtins,
	}
}

// Execute reads palette.yaml and tokens.yaml from the named theme and
// writes an ANSI-colored preview to w. If tokens.yaml doesn't exist,
// tokens are derived on-the-fly from palette.yaml or a built-in palette.
func (uc *PreviewThemeUseCase) Execute(themeName string, w io.Writer) error {
	// Try to read and parse palette.yaml (optional — may be a comment-only stub).
	palette := uc.tryParsePalette(themeName)

	// Try to read tokens.yaml first.
	tokens, err := uc.tryReadTokens(themeName)
	if err != nil {
		// tokens.yaml doesn't exist or can't be read — derive tokens instead.
		tokens, palette, err = uc.deriveTokens(themeName, palette)
		if err != nil {
			return fmt.Errorf("preview %q: %w", themeName, err)
		}
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Theme: %s\n\n", themeName))

	// Render palette color swatches when available.
	if palette != nil {
		renderPaletteSection(&sb, palette)
	}

	// Render semantic token groups.
	renderTokenSection(&sb, tokens)

	_, err = fmt.Fprint(w, sb.String())
	return err
}

// tryReadTokens attempts to read tokens.yaml from the store.
// Returns nil tokens and an error if the file doesn't exist or can't be read.
func (uc *PreviewThemeUseCase) tryReadTokens(themeName string) (*domain.TokenSet, error) {
	tokensRC, err := uc.store.OpenReader(themeName, "tokens.yaml")
	if err != nil {
		return nil, err
	}
	defer func() { _ = tokensRC.Close() }()

	return uc.tokensReader(tokensRC)
}

// deriveTokens derives tokens from a palette. If palette is nil, it tries
// to parse palette.yaml from the store. If that fails, it checks if themeName
// is a built-in palette and derives from that.
func (uc *PreviewThemeUseCase) deriveTokens(themeName string, palette *domain.Palette) (*domain.TokenSet, *domain.Palette, error) {
	// If we already have a palette, derive from it.
	if palette != nil {
		return uc.tokenizer.Tokenize(palette), palette, nil
	}

	// Try to get from built-in palettes.
	if uc.builtins != nil && uc.builtins.Has(themeName) {
		r, err := uc.builtins.Get(themeName)
		if err != nil {
			return nil, nil, fmt.Errorf("get built-in palette: %w", err)
		}
		palette, err = uc.parser.Parse(r)
		if err != nil {
			return nil, nil, fmt.Errorf("parse built-in palette: %w", err)
		}
		return uc.tokenizer.Tokenize(palette), palette, nil
	}

	return nil, nil, errors.New("no palette or tokens.yaml found")
}

// tryParsePalette attempts to open and parse palette.yaml from the store.
// Returns nil if the file is missing or cannot be parsed (e.g. comment-only stub).
func (uc *PreviewThemeUseCase) tryParsePalette(themeName string) *domain.Palette {
	rc, err := uc.store.OpenReader(themeName, "palette.yaml")
	if err != nil {
		return nil
	}
	defer func() { _ = rc.Close() }()

	palette, err := uc.parser.Parse(rc)
	if err != nil {
		return nil
	}
	return palette
}

// Parser returns the palette parser used by this use case.
func (uc *PreviewThemeUseCase) Parser() ports.PaletteParser {
	return uc.parser
}

// Tokenizer returns the tokenizer used by this use case.
func (uc *PreviewThemeUseCase) Tokenizer() ports.Tokenizer {
	return uc.tokenizer
}

// renderPaletteSection writes the palette color swatches with ANSI escape codes.
func renderPaletteSection(sb *strings.Builder, palette *domain.Palette) {
	sb.WriteString("Palette Colors\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	slotNames := palette.SlotNames()
	for _, slot := range slotNames {
		c, _ := palette.Slot(slot)
		swatch := fgBgSwatch(c)
		fmt.Fprintf(sb, "  %-8s %s  %s\n", slot, swatch, fgHex(c))
	}
	sb.WriteString("\n")
}

// renderTokenSection writes the semantic token groups with ANSI escape codes.
func renderTokenSection(sb *strings.Builder, tokens *domain.TokenSet) {
	sb.WriteString("Semantic Tokens\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	tokenPaths := tokens.Paths()
	for _, path := range tokenPaths {
		tok, _ := tokens.Get(path)
		if tok.Color.IsNone {
			fmt.Fprintf(sb, "  %-40s (none)\n", path)
			continue
		}
		c := tok.Color
		fmt.Fprintf(sb, "  %-40s %s  %s\n", path, fgBgSwatch(c), fgHex(c))
	}
}

// fgBgSwatch returns an ANSI 24-bit background color swatch (4 spaces wide).
func fgBgSwatch(c domain.Color) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm    \x1b[0m", c.R, c.G, c.B)
}

// fgHex returns the hex color value rendered in its own foreground color.
func fgHex(c domain.Color) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", c.R, c.G, c.B, c.Hex())
}

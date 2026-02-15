package application

import (
	"fmt"
	"io"
	"strings"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// UniversalReader reads a universal.yaml from a reader and returns a TokenSet.
// The composition root wires this to fileio.ReadUniversal.
type UniversalReader func(r io.Reader) (*domain.TokenSet, error)

// PreviewThemeUseCase reads palette.yaml and universal.yaml from a theme
// directory and renders an ANSI-colored preview to an io.Writer.
type PreviewThemeUseCase struct {
	store           ports.ThemeStore
	parser          ports.PaletteParser
	universalReader UniversalReader
}

// NewPreviewThemeUseCase returns a new PreviewThemeUseCase wired to the
// given store, palette parser, and universal reader function.
func NewPreviewThemeUseCase(store ports.ThemeStore, parser ports.PaletteParser, ur UniversalReader) *PreviewThemeUseCase {
	return &PreviewThemeUseCase{store: store, parser: parser, universalReader: ur}
}

// Execute reads palette.yaml and universal.yaml from the named theme and
// writes an ANSI-colored preview to w. Returns an error if universal.yaml
// cannot be read. The palette section is shown when palette.yaml is valid.
func (uc *PreviewThemeUseCase) Execute(themeName string, w io.Writer) error {
	// Try to read and parse palette.yaml (optional — may be a comment-only stub).
	palette := uc.tryParsePalette(themeName)

	// Read and parse universal.yaml (required).
	univRC, err := uc.store.OpenReader(themeName, "universal.yaml")
	if err != nil {
		return fmt.Errorf("preview %q: %w", themeName, err)
	}
	defer func() { _ = univRC.Close() }()

	tokens, err := uc.universalReader(univRC)
	if err != nil {
		return fmt.Errorf("preview %q: read universal: %w", themeName, err)
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

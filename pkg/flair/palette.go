package flair

// Palette represents a complete base24 color palette with 24 indexed color slots.
//
// A base24 palette extends the base16 scheme with 8 additional slots for
// enhanced color coverage. The slots are organized as follows:
//
//	00-07: Background and foreground shades (darkest to lightest)
//	08-0F: Accent colors (red, orange, yellow, green, cyan, blue, purple, brown)
//	10-17: Extended colors (darker bg, bright variants)
//
// Palette provides a public, immutable view of a color palette that can be used
// by external consumers without depending on internal packages. Palette instances
// are created via [ParsePalette] and are safe for concurrent use since they are
// read-only after construction.
//
// To convert a Palette to a Theme with semantic tokens, use [Tokenize].
//
// Example:
//
//	palette, err := flair.ParsePalette(reader)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	theme := flair.Tokenize(palette)
type Palette struct {
	name    string
	author  string
	variant string
	colors  [24]Color
}

// Name returns the name of the palette (e.g., "Tokyo Night").
func (p *Palette) Name() string {
	return p.name
}

// Author returns the author of the palette.
func (p *Palette) Author() string {
	return p.author
}

// Variant returns the variant of the palette (e.g., "dark", "light", "storm").
func (p *Palette) Variant() string {
	return p.variant
}

// Base returns the color at the given base24 index (0x00-0x17, or 0-23 decimal).
//
// Common base24 indices:
//
//	0x00: Background
//	0x05: Foreground
//	0x08: Red
//	0x0B: Green
//	0x0D: Blue
//	0x0E: Purple
//
// Base returns nil for out-of-range indices.
//
// Example:
//
//	bg := palette.Base(0x00)   // Background color
//	red := palette.Base(0x08)  // Red accent
func (p *Palette) Base(index int) *Color {
	if index < 0 || index >= len(p.colors) {
		return nil
	}
	// Return a copy to maintain immutability.
	c := p.colors[index]
	return &c
}

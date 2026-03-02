package domain

// TokenOverride represents optional overrides that can be applied to a Token.
// It allows users to customize specific tokens in their theme by overriding
// color and/or style flags. Nil Color means "no color override".
type TokenOverride struct {
	Color         *Color
	Bold          bool
	Italic        bool
	Underline     bool
	Undercurl     bool
	Strikethrough bool
}

// NewTokenOverride creates a TokenOverride from a hex color string and style flags.
// If colorHex is empty, the Color field will be nil (no color override).
// Returns a *ParseError if the hex color string is invalid.
func NewTokenOverride(colorHex string, bold, italic, underline, undercurl, strikethrough bool) (*TokenOverride, error) {
	override := &TokenOverride{
		Bold:          bold,
		Italic:        italic,
		Underline:     underline,
		Undercurl:     undercurl,
		Strikethrough: strikethrough,
	}

	if colorHex != "" {
		c, err := ParseHex(colorHex)
		if err != nil {
			return nil, err
		}
		override.Color = &c
	}

	return override, nil
}

// HasColor returns true if the override has a color set.
func (o TokenOverride) HasColor() bool {
	return o.Color != nil
}

// HasStyle returns true if any style flag (Bold, Italic, Underline,
// Undercurl, Strikethrough) is set on the override.
func (o TokenOverride) HasStyle() bool {
	return o.Bold || o.Italic || o.Underline || o.Undercurl || o.Strikethrough
}

// Apply merges the override into a base Token and returns a new Token.
// Color is replaced only if the override has a color set.
// Style flags from the override are OR'd with the base token's flags.
func (o TokenOverride) Apply(base Token) Token {
	result := Token{
		Color:         base.Color,
		Bold:          base.Bold || o.Bold,
		Italic:        base.Italic || o.Italic,
		Underline:     base.Underline || o.Underline,
		Undercurl:     base.Undercurl || o.Undercurl,
		Strikethrough: base.Strikethrough || o.Strikethrough,
	}

	if o.HasColor() {
		result.Color = *o.Color
	}

	return result
}

package domain

// Token represents a design token that maps a semantic path to a color
// and optional style flags (bold, italic, underline, undercurl, strikethrough).
type Token struct {
	Color         Color
	Bold          bool
	Italic        bool
	Underline     bool
	Undercurl     bool
	Strikethrough bool
}

// HasStyle returns true if any style flag (Bold, Italic, Underline,
// Undercurl, Strikethrough) is set on the token.
func (t Token) HasStyle() bool {
	return t.Bold || t.Italic || t.Underline || t.Undercurl || t.Strikethrough
}

package domain

// ResolvedTheme represents a fully resolved theme with its palette and
// semantic token mappings. It is the aggregate that downstream renderers
// consume to produce output for specific applications (e.g., Neovim, Alacritty).
type ResolvedTheme struct {
	Name    string
	Variant string
	Palette *Palette
	Tokens  *TokenSet
}

// Token retrieves the Token at the given semantic path by delegating to the
// underlying TokenSet. Returns the token and true if found, or a zero Token
// and false if the path does not exist.
func (rt *ResolvedTheme) Token(path string) (Token, bool) {
	return rt.Tokens.Get(path)
}

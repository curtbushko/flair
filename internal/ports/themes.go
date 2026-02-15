package ports

import "github.com/curtbushko/flair/internal/domain"

// VimHighlight represents a single Vim highlight group with optional colors
// and style attributes. Nil color pointers indicate "no color set".
type VimHighlight struct {
	Fg            *domain.Color
	Bg            *domain.Color
	Sp            *domain.Color
	Bold          bool
	Italic        bool
	Underline     bool
	Undercurl     bool
	Strikethrough bool
	Reverse       bool
	Nocombine     bool
	Link          string
}

// VimTheme is the mapped theme DTO for the Vim/Neovim target.
type VimTheme struct {
	Name           string
	Highlights     map[string]VimHighlight
	TerminalColors [16]domain.Color
}

// StylixTheme is the mapped theme DTO for the Stylix target.
type StylixTheme struct {
	Values map[string]string
}

// CSSProperty is a single CSS property-value pair.
type CSSProperty struct {
	Property string
	Value    string
}

// CSSRule is a CSS rule with a selector and its properties.
type CSSRule struct {
	Selector   string
	Properties []CSSProperty
}

// CSSTheme is the mapped theme DTO for the CSS target.
type CSSTheme struct {
	CustomProperties map[string]string
	Rules            []CSSRule
}

// GtkColorDef is a GTK @define-color name-value pair.
type GtkColorDef struct {
	Name  string
	Value string
}

// GtkTheme is the mapped theme DTO for the GTK target.
type GtkTheme struct {
	Colors []GtkColorDef
	Rules  []CSSRule
}

// QssTheme is the mapped theme DTO for the QSS target.
type QssTheme struct {
	Rules []CSSRule
}

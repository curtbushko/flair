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

// LualineModeColors holds the fg/bg colors for a lualine mode section.
type LualineModeColors struct {
	Fg *domain.Color
	Bg *domain.Color
}

// LualineMode holds the colors for sections a, b, c of a lualine mode.
type LualineMode struct {
	A LualineModeColors
	B LualineModeColors
	C LualineModeColors
}

// LualineTheme holds the complete lualine theme with all modes.
type LualineTheme struct {
	Normal   LualineMode
	Insert   LualineMode
	Visual   LualineMode
	Replace  LualineMode
	Command  LualineMode
	Inactive LualineMode
}

// BufferlineColors holds the fg/bg colors and style for a bufferline highlight group.
type BufferlineColors struct {
	Fg     *domain.Color
	Bg     *domain.Color
	Bold   bool
	Italic bool
}

// BufferlineTheme holds the complete bufferline theme with all highlight groups.
type BufferlineTheme struct {
	Background                BufferlineColors
	Fill                      BufferlineColors
	BufferSelected            BufferlineColors
	BufferVisible             BufferlineColors
	CloseButton               BufferlineColors
	CloseButtonSelected       BufferlineColors
	CloseButtonVisible        BufferlineColors
	Diagnostic                BufferlineColors
	DiagnosticSelected        BufferlineColors
	DiagnosticVisible         BufferlineColors
	Duplicate                 BufferlineColors
	DuplicateSelected         BufferlineColors
	DuplicateVisible          BufferlineColors
	Error                     BufferlineColors
	ErrorSelected             BufferlineColors
	ErrorVisible              BufferlineColors
	ErrorDiagnostic           BufferlineColors
	ErrorDiagnosticSelected   BufferlineColors
	ErrorDiagnosticVisible    BufferlineColors
	Hint                      BufferlineColors
	HintSelected              BufferlineColors
	HintVisible               BufferlineColors
	HintDiagnostic            BufferlineColors
	HintDiagnosticSelected    BufferlineColors
	HintDiagnosticVisible     BufferlineColors
	IndicatorSelected         BufferlineColors
	IndicatorVisible          BufferlineColors
	Info                      BufferlineColors
	InfoSelected              BufferlineColors
	InfoVisible               BufferlineColors
	InfoDiagnostic            BufferlineColors
	InfoDiagnosticSelected    BufferlineColors
	InfoDiagnosticVisible     BufferlineColors
	Modified                  BufferlineColors
	ModifiedSelected          BufferlineColors
	ModifiedVisible           BufferlineColors
	Numbers                   BufferlineColors
	NumbersSelected           BufferlineColors
	NumbersVisible            BufferlineColors
	OffsetSeparator           BufferlineColors
	Pick                      BufferlineColors
	PickSelected              BufferlineColors
	PickVisible               BufferlineColors
	Separator                 BufferlineColors
	SeparatorSelected         BufferlineColors
	SeparatorVisible          BufferlineColors
	Tab                       BufferlineColors
	TabClose                  BufferlineColors
	TabSelected               BufferlineColors
	TabSeparator              BufferlineColors
	TabSeparatorSelected      BufferlineColors
	TruncMarker               BufferlineColors
	Warning                   BufferlineColors
	WarningSelected           BufferlineColors
	WarningVisible            BufferlineColors
	WarningDiagnostic         BufferlineColors
	WarningDiagnosticSelected BufferlineColors
	WarningDiagnosticVisible  BufferlineColors
}

// VimTheme is the mapped theme DTO for the Vim/Neovim target.
type VimTheme struct {
	Name           string
	Highlights     map[string]VimHighlight
	TerminalColors [16]domain.Color
	Lualine        *LualineTheme
	Bufferline     *BufferlineTheme
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

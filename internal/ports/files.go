package ports

import "github.com/curtbushko/flair/internal/domain"

// FileHeader is embedded in every YAML file produced by flair.
type FileHeader struct {
	SchemaVersion int             `yaml:"schema_version"`
	Kind          domain.FileKind `yaml:"kind"`
	ThemeName     string          `yaml:"theme_name"`
}

// PaletteFile is the input palette (base24).
type PaletteFile struct {
	FileHeader `yaml:",inline"`
	System     string            `yaml:"system"`
	Author     string            `yaml:"author"`
	Variant    string            `yaml:"variant"`
	Palette    map[string]string `yaml:"palette"`
}

// TokenEntry is a single semantic token in tokens.yaml.
type TokenEntry struct {
	Color         string `yaml:"color"`
	Bold          bool   `yaml:"bold,omitempty"`
	Italic        bool   `yaml:"italic,omitempty"`
	Underline     bool   `yaml:"underline,omitempty"`
	Undercurl     bool   `yaml:"undercurl,omitempty"`
	Strikethrough bool   `yaml:"strikethrough,omitempty"`
}

// TokensFile is the derived semantic token set.
type TokensFile struct {
	FileHeader `yaml:",inline"`
	Tokens     map[string]TokenEntry `yaml:"tokens"`
}

// VimMappingHighlight is a single Vim highlight group in the mapping file.
type VimMappingHighlight struct {
	Fg            string `yaml:"fg,omitempty"`
	Bg            string `yaml:"bg,omitempty"`
	Sp            string `yaml:"sp,omitempty"`
	Bold          bool   `yaml:"bold,omitempty"`
	Italic        bool   `yaml:"italic,omitempty"`
	Underline     bool   `yaml:"underline,omitempty"`
	Undercurl     bool   `yaml:"undercurl,omitempty"`
	Strikethrough bool   `yaml:"strikethrough,omitempty"`
	Reverse       bool   `yaml:"reverse,omitempty"`
	Nocombine     bool   `yaml:"nocombine,omitempty"`
	Link          string `yaml:"link,omitempty"`
}

// BufferlineMappingColors is a single bufferline highlight group in the mapping file.
type BufferlineMappingColors struct {
	Fg     string `yaml:"fg,omitempty"`
	Bg     string `yaml:"bg,omitempty"`
	Bold   bool   `yaml:"bold,omitempty"`
	Italic bool   `yaml:"italic,omitempty"`
}

// BufferlineMappingTheme holds all bufferline highlight groups for the mapping file.
type BufferlineMappingTheme struct {
	Fill              BufferlineMappingColors `yaml:"fill"`
	Background        BufferlineMappingColors `yaml:"background"`
	BufferVisible     BufferlineMappingColors `yaml:"buffer_visible"`
	BufferSelected    BufferlineMappingColors `yaml:"buffer_selected"`
	Separator         BufferlineMappingColors `yaml:"separator"`
	SeparatorVisible  BufferlineMappingColors `yaml:"separator_visible"`
	SeparatorSelected BufferlineMappingColors `yaml:"separator_selected"`
	IndicatorSelected BufferlineMappingColors `yaml:"indicator_selected"`
	Modified          BufferlineMappingColors `yaml:"modified"`
	ModifiedVisible   BufferlineMappingColors `yaml:"modified_visible"`
	ModifiedSelected  BufferlineMappingColors `yaml:"modified_selected"`
	Error             BufferlineMappingColors `yaml:"error"`
	Warning           BufferlineMappingColors `yaml:"warning"`
	Info              BufferlineMappingColors `yaml:"info"`
	Hint              BufferlineMappingColors `yaml:"hint"`
}

// VimMappingFile is the Vim-specific mapping.
type VimMappingFile struct {
	FileHeader     `yaml:",inline"`
	Highlights     map[string]VimMappingHighlight `yaml:"highlights"`
	TerminalColors [16]string                     `yaml:"terminal_colors"`
	Bufferline     *BufferlineMappingTheme        `yaml:"bufferline,omitempty"`
}

// CSSRuleEntry is a CSS selector with its properties.
type CSSRuleEntry struct {
	Selector   string            `yaml:"selector"`
	Properties map[string]string `yaml:"properties"`
}

// CSSMappingFile is the CSS-specific mapping.
type CSSMappingFile struct {
	FileHeader       `yaml:",inline"`
	CustomProperties map[string]string `yaml:"custom_properties"`
	Rules            []CSSRuleEntry    `yaml:"rules"`
}

// GtkMappingFile is the GTK-specific mapping.
type GtkMappingFile struct {
	FileHeader `yaml:",inline"`
	Colors     map[string]string `yaml:"colors"`
	Rules      []CSSRuleEntry    `yaml:"rules"`
}

// QssMappingFile is the QSS-specific mapping.
type QssMappingFile struct {
	FileHeader `yaml:",inline"`
	Rules      []CSSRuleEntry `yaml:"rules"`
}

// StylixMappingFile is the Stylix-specific mapping.
type StylixMappingFile struct {
	FileHeader `yaml:",inline"`
	Values     map[string]string `yaml:"values"`
}

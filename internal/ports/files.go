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

// UniversalToken is a single semantic token in universal.yaml.
type UniversalToken struct {
	Color         string `yaml:"color"`
	Bold          bool   `yaml:"bold,omitempty"`
	Italic        bool   `yaml:"italic,omitempty"`
	Underline     bool   `yaml:"underline,omitempty"`
	Undercurl     bool   `yaml:"undercurl,omitempty"`
	Strikethrough bool   `yaml:"strikethrough,omitempty"`
}

// UniversalFile is the derived semantic token set.
type UniversalFile struct {
	FileHeader `yaml:",inline"`
	Tokens     map[string]UniversalToken `yaml:"tokens"`
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

// VimMappingFile is the Vim-specific mapping.
type VimMappingFile struct {
	FileHeader     `yaml:",inline"`
	Highlights     map[string]VimMappingHighlight `yaml:"highlights"`
	TerminalColors [16]string                     `yaml:"terminal_colors"`
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

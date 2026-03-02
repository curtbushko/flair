package domain

// Schema versions for each file type. Bump when the file format changes.
const (
	SchemaPalette       = 1
	SchemaTokens        = 1
	SchemaVimMapping    = 1
	SchemaCSSMapping    = 1
	SchemaGtkMapping    = 1
	SchemaQssMapping    = 1
	SchemaStylixMapping = 1
)

// FileKind identifies a file type in the pipeline.
type FileKind string

// FileKind constants for each file type produced by flair.
const (
	FileKindPalette       FileKind = "palette"
	FileKindTokens        FileKind = "tokens"
	FileKindVimMapping    FileKind = "vim-mapping"
	FileKindCSSMapping    FileKind = "css-mapping"
	FileKindGtkMapping    FileKind = "gtk-mapping"
	FileKindQssMapping    FileKind = "qss-mapping"
	FileKindStylixMapping FileKind = "stylix-mapping"
)

// CurrentVersion returns the current schema version for a file kind.
// Returns 0 for unknown kinds.
func CurrentVersion(kind FileKind) int {
	switch kind {
	case FileKindPalette:
		return SchemaPalette
	case FileKindTokens:
		return SchemaTokens
	case FileKindVimMapping:
		return SchemaVimMapping
	case FileKindCSSMapping:
		return SchemaCSSMapping
	case FileKindGtkMapping:
		return SchemaGtkMapping
	case FileKindQssMapping:
		return SchemaQssMapping
	case FileKindStylixMapping:
		return SchemaStylixMapping
	default:
		return 0
	}
}

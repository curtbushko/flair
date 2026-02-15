package ports

import "io"

// Generator writes the final output file from a mapped theme.
type Generator interface {
	// Name returns the target name (e.g. "vim", "css", "gtk").
	Name() string

	// DefaultFilename returns the default output filename (e.g. "style.lua", "gtk.css").
	DefaultFilename() string

	// Generate writes the final output to w from a mapped theme.
	Generate(w io.Writer, mapped MappedTheme) error
}

// Target pairs a mapper with its generator and mapping file path.
type Target struct {
	Mapper      Mapper
	Generator   Generator
	MappingFile string // filename in theme dir, e.g. "vim-mapping.yaml"
}

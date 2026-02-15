package ports

import (
	"io"

	"github.com/curtbushko/flair/internal/domain"
)

// Generator writes the final output file from a mapped theme.
type Generator interface {
	// Name returns the target name (e.g. "vim", "css", "gtk").
	Name() string

	// DefaultFilename returns the default output filename (e.g. "style.lua", "gtk.css").
	DefaultFilename() string

	// Generate writes the final output to w from a mapped theme.
	Generate(w io.Writer, mapped MappedTheme) error
}

// MappingFileWriter is a function that serializes a mapped theme to YAML
// via an io.Writer. The composition root wires the specific fileio.WriteXxx
// function for each target. This keeps the application layer adapter-agnostic.
type MappingFileWriter func(w io.Writer, mapped MappedTheme) error

// Target pairs a mapper with its generator and mapping file I/O.
type Target struct {
	Mapper           Mapper
	Generator        Generator
	MappingFile      string            // filename in theme dir, e.g. "vim-mapping.yaml"
	MappingFileKind  domain.FileKind   // schema kind for versioned header
	WriteMappingFile MappingFileWriter // serializes mapped theme to YAML
}

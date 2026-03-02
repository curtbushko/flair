// Package ports defines the interface boundaries for the flair theme pipeline.
// It contains port interfaces (PaletteParser, PaletteSource, Tokenizer, Mapper,
// Generator, ThemeStore), file structs for YAML serialization, and theme DTOs
// shared between mapper and generator adapters.
package ports

import (
	"io"

	"github.com/curtbushko/flair/internal/domain"
)

// PaletteParser reads palette YAML from a reader and returns a domain Palette.
// The caller is responsible for opening/closing the underlying source.
// Works identically on files, embedded built-ins, test buffers, or stdin.
type PaletteParser interface {
	Parse(r io.Reader) (*domain.Palette, error)
}

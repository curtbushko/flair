// Package palettes provides a PaletteSource adapter that delegates to
// pkg/flair/palettes for embedded YAML files. This ensures a single
// source of truth for built-in palettes (DRY principle).
package palettes

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"

	pkgpalettes "github.com/curtbushko/flair/pkg/flair/palettes"
)

// Source implements ports.PaletteSource using embedded YAML files
// from pkg/flair/palettes.
type Source struct{}

// NewSource returns a new built-in palette source.
func NewSource() *Source {
	return &Source{}
}

// List returns the names of all built-in palettes sorted alphabetically.
// Each name is the filename without the .yaml extension.
func (s *Source) List() []string {
	entries, err := fs.ReadDir(pkgpalettes.EmbeddedFS, ".")
	if err != nil {
		return nil
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") {
			names = append(names, strings.TrimSuffix(name, ".yaml"))
		}
	}

	sort.Strings(names)
	return names
}

// Get returns an io.Reader for the named built-in palette's YAML.
// The returned reader wraps the embedded bytes via bytes.NewReader,
// so no file I/O occurs. Returns an error if the name is not found.
func (s *Source) Get(name string) (io.Reader, error) {
	data, err := pkgpalettes.EmbeddedFS.ReadFile(name + ".yaml")
	if err != nil {
		return nil, fmt.Errorf("built-in palette %q not found: %w", name, err)
	}
	return bytes.NewReader(data), nil
}

// Has returns true if the named palette exists as a built-in.
func (s *Source) Has(name string) bool {
	_, err := pkgpalettes.EmbeddedFS.ReadFile(name + ".yaml")
	return err == nil
}

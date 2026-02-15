// Package palettes provides a PaletteSource adapter using go:embed to
// ship built-in palette YAML files with the binary. Each .yaml file in
// this directory is embedded at compile time and accessible via List,
// Get, and Has methods.
package palettes

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"
)

//go:embed *.yaml
var embedded embed.FS

// Source implements ports.PaletteSource using embedded YAML files.
type Source struct{}

// NewSource returns a new built-in palette source.
func NewSource() *Source {
	return &Source{}
}

// List returns the names of all built-in palettes sorted alphabetically.
// Each name is the filename without the .yaml extension.
func (s *Source) List() []string {
	entries, err := fs.ReadDir(embedded, ".")
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
	data, err := embedded.ReadFile(name + ".yaml")
	if err != nil {
		return nil, fmt.Errorf("built-in palette %q not found: %w", name, err)
	}
	return bytes.NewReader(data), nil
}

// Has returns true if the named palette exists as a built-in.
func (s *Source) Has(name string) bool {
	_, err := embedded.ReadFile(name + ".yaml")
	return err == nil
}

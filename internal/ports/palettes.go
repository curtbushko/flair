package ports

import "io"

// PaletteSource provides access to built-in palettes shipped with flair.
type PaletteSource interface {
	// List returns the names of all built-in palettes (e.g. "tokyo-night-dark").
	List() []string

	// Get returns a reader for the named built-in palette's YAML.
	// Returns an error if the name is not found.
	Get(name string) (io.Reader, error)

	// Has returns true if the named palette exists as a built-in.
	Has(name string) bool
}

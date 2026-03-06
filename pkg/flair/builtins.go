package flair

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/curtbushko/flair/pkg/flair/palettes"
)

// ListBuiltins returns a sorted slice of all built-in palette names.
//
// The returned names can be used with [LoadBuiltin] to load the corresponding
// theme without filesystem access. Built-in themes are embedded in the binary
// and always available.
//
// Example:
//
//	for _, name := range flair.ListBuiltins() {
//	    fmt.Println(name)
//	}
//	// Output:
//	// catppuccin-frappe
//	// catppuccin-latte
//	// gruvbox-dark
//	// tokyo-night-dark
//	// ...
func ListBuiltins() []string {
	entries, err := palettes.EmbeddedFS.ReadDir(".")
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
			// Strip .yaml extension to get palette name
			names = append(names, strings.TrimSuffix(name, ".yaml"))
		}
	}

	sort.Strings(names)
	return names
}

// HasBuiltin reports whether a built-in palette with the given name exists.
//
// This function checks the embedded palette files without loading the full theme.
// Use [ListBuiltins] to get all available names.
//
// Example:
//
//	if flair.HasBuiltin("tokyo-night-dark") {
//	    theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	}
func HasBuiltin(name string) bool {
	if name == "" {
		return false
	}
	filename := name + ".yaml"
	_, err := palettes.EmbeddedFS.ReadFile(filename)
	return err == nil
}

// LoadBuiltin loads a built-in palette by name and returns a fully tokenized Theme.
//
// LoadBuiltin does not access the filesystem; all palettes are embedded in the
// binary. This makes it suitable for use in applications that need guaranteed
// theme availability without external dependencies.
//
// The returned theme contains all semantic tokens derived from the base24 palette
// via [Tokenize].
//
// Returns an error if the palette does not exist or cannot be parsed.
// Use [ListBuiltins] to see available palette names.
//
// Example:
//
//	theme, err := flair.LoadBuiltin("tokyo-night-dark")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Loaded %s (%s variant)\n", theme.Name(), theme.Variant())
func LoadBuiltin(name string) (*Theme, error) {
	if name == "" {
		return nil, errors.New("palette name cannot be empty")
	}

	filename := name + ".yaml"
	data, err := palettes.EmbeddedFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("built-in palette %q not found", name)
	}

	palette, err := ParsePalette(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse built-in palette %q: %w", name, err)
	}

	return Tokenize(palette), nil
}

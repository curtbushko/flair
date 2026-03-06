// Package palettes provides embedded palette YAML files for built-in themes.
package palettes

import "embed"

// EmbeddedFS contains all embedded palette YAML files.
// These files are compiled into the binary and do not require
// filesystem access at runtime.
//
//go:embed *.yaml
var EmbeddedFS embed.FS

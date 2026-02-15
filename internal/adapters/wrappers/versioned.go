// Package wrappers provides io.Writer decorators for the flair pipeline.
package wrappers

import (
	"fmt"
	"io"

	"github.com/curtbushko/flair/internal/domain"
)

// VersionedWriter wraps an io.Writer and prepends a YAML header containing
// schema_version, kind, and theme_name on the first Write call.
// Subsequent writes pass through directly to the inner writer.
type VersionedWriter struct {
	inner      io.Writer
	kind       domain.FileKind
	themeName  string
	headerDone bool
}

// NewVersionedWriter creates a VersionedWriter that wraps w. On the first
// call to Write it emits a three-line YAML header (schema_version, kind,
// theme_name) before forwarding the caller's data.
func NewVersionedWriter(w io.Writer, kind domain.FileKind, themeName string) *VersionedWriter {
	return &VersionedWriter{
		inner:     w,
		kind:      kind,
		themeName: themeName,
	}
}

// Write implements io.Writer. On the first invocation it writes the YAML
// header to the inner writer, then writes p. The returned byte count n
// reflects only the caller's data (len(p)), not the header bytes.
func (vw *VersionedWriter) Write(p []byte) (int, error) {
	if !vw.headerDone {
		vw.headerDone = true

		header := fmt.Sprintf("schema_version: %d\nkind: %s\ntheme_name: %s\n",
			domain.CurrentVersion(vw.kind), string(vw.kind), vw.themeName)

		if _, err := io.WriteString(vw.inner, header); err != nil {
			return 0, fmt.Errorf("write versioned header: %w", err)
		}
	}

	if len(p) == 0 {
		return 0, nil
	}

	n, err := vw.inner.Write(p)
	if err != nil {
		return n, fmt.Errorf("write content: %w", err)
	}

	return n, nil
}

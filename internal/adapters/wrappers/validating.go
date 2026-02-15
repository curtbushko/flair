package wrappers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/curtbushko/flair/internal/domain"
)

// peekSize is the maximum number of bytes to read from the inner reader
// to extract the schema_version header line. 512 bytes is more than enough
// for any realistic YAML header.
const peekSize = 512

// ValidatingReader wraps an io.Reader and peeks at the beginning of the
// stream on the first Read to extract the schema_version field.
// If the version does not match domain.CurrentVersion for the given FileKind,
// Read returns a *domain.SchemaVersionError. When the version matches,
// the peeked bytes are replayed seamlessly via io.MultiReader.
type ValidatingReader struct {
	inner     io.Reader
	kind      domain.FileKind
	validated bool
}

// NewValidatingReader creates a ValidatingReader that validates the
// schema_version header of the underlying reader against the current
// version for the given FileKind.
func NewValidatingReader(r io.Reader, kind domain.FileKind) *ValidatingReader {
	return &ValidatingReader{
		inner: r,
		kind:  kind,
	}
}

// Read implements io.Reader. On the first call, it peeks enough bytes to
// parse schema_version. If the version is incompatible, it returns a
// *domain.SchemaVersionError. Otherwise, it replays the peeked bytes
// followed by the remainder of the stream.
func (vr *ValidatingReader) Read(p []byte) (int, error) {
	if !vr.validated {
		if err := vr.validate(); err != nil {
			return 0, err
		}
	}
	return vr.inner.Read(p)
}

// validate peeks bytes from inner, extracts schema_version, checks it
// against the current version, and rebuilds inner as a MultiReader that
// replays the peeked bytes.
func (vr *ValidatingReader) validate() error {
	vr.validated = true

	peek := make([]byte, peekSize)
	n, err := io.ReadAtLeast(vr.inner, peek, 1)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("peek schema_version: %w", err)
	}
	peek = peek[:n]

	version, err := extractSchemaVersion(peek)
	if err != nil {
		return err
	}

	current := domain.CurrentVersion(vr.kind)
	if version != current {
		return &domain.SchemaVersionError{
			Found:        version,
			Expected:     current,
			NeedsUpgrade: version > current,
		}
	}

	// Replay peeked bytes + rest of stream.
	vr.inner = io.MultiReader(bytes.NewReader(peek), vr.inner)
	return nil
}

// extractSchemaVersion scans the peeked bytes for a line matching
// "schema_version: <int>" and returns the parsed integer.
func extractSchemaVersion(data []byte) (int, error) {
	const prefix = "schema_version:"

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			valStr := strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
			v, err := strconv.Atoi(valStr)
			if err != nil {
				return 0, &domain.SchemaVersionError{
					Found:    0,
					Expected: 0,
				}
			}
			return v, nil
		}
	}

	return 0, errors.New("schema_version field not found in header")
}

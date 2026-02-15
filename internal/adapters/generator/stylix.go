// Package generator provides adapters that write final output files from
// mapped theme structs (ports.StylixTheme, etc.).
package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// Stylix implements ports.Generator for the Stylix target.
// It writes a JSON file with sorted keys and 2-space indent from a StylixTheme.
type Stylix struct{}

// NewStylix returns a new Stylix generator.
func NewStylix() *Stylix {
	return &Stylix{}
}

// Name returns the target name for this generator.
func (s *Stylix) Name() string {
	return "stylix"
}

// DefaultFilename returns the default output filename for Stylix.
func (s *Stylix) DefaultFilename() string {
	return "style.json"
}

// Generate writes the StylixTheme as a JSON object with sorted keys and
// 2-space indent to w. The mapped argument must be a *ports.StylixTheme;
// a type assertion failure returns a *domain.GenerateError.
func (s *Stylix) Generate(w io.Writer, mapped ports.MappedTheme) error {
	theme, ok := mapped.(*ports.StylixTheme)
	if !ok {
		return &domain.GenerateError{
			Target:  "stylix",
			Message: fmt.Sprintf("expected *ports.StylixTheme, got %T", mapped),
		}
	}

	// Build an ordered structure for deterministic JSON output.
	// json.MarshalIndent on a map does sort keys alphabetically in Go,
	// but we use an explicit sorted approach for clarity and control.
	ordered := sortedMap(theme.Values)

	data, err := json.MarshalIndent(ordered, "", "  ")
	if err != nil {
		return &domain.GenerateError{
			Target:  "stylix",
			Message: "failed to marshal JSON",
			Cause:   err,
		}
	}

	// Append trailing newline for POSIX compliance.
	data = append(data, '\n')

	if _, err := w.Write(data); err != nil {
		return &domain.GenerateError{
			Target:  "stylix",
			Message: "failed to write output",
			Cause:   err,
		}
	}

	return nil
}

// sortedKeyValue is an ordered key-value pair for JSON serialization.
type sortedKeyValue struct {
	Key   string
	Value string
}

// sortedMap converts a map[string]string into an ordered slice of key-value
// pairs sorted by key, then wraps it in a type that marshals as a JSON object
// with keys in the sorted order.
func sortedMap(m map[string]string) *orderedJSON {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]sortedKeyValue, len(keys))
	for i, k := range keys {
		pairs[i] = sortedKeyValue{Key: k, Value: m[k]}
	}

	return &orderedJSON{pairs: pairs}
}

// orderedJSON implements json.Marshaler to produce a JSON object with keys
// in the order defined by pairs.
type orderedJSON struct {
	pairs []sortedKeyValue
}

// MarshalJSON writes a JSON object with keys in the order of pairs.
func (o *orderedJSON) MarshalJSON() ([]byte, error) {
	var buf []byte
	buf = append(buf, '{')
	for i, p := range o.pairs {
		if i > 0 {
			buf = append(buf, ',')
		}
		key, err := json.Marshal(p.Key)
		if err != nil {
			return nil, err
		}
		val, err := json.Marshal(p.Value)
		if err != nil {
			return nil, err
		}
		buf = append(buf, key...)
		buf = append(buf, ':')
		buf = append(buf, val...)
	}
	buf = append(buf, '}')
	return buf, nil
}

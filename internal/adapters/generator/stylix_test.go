package generator_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/generator"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// sampleStylixTheme returns a StylixTheme with 5 key-value pairs for testing.
func sampleStylixTheme() *ports.StylixTheme {
	return &ports.StylixTheme{
		Values: map[string]string{
			"base00":         "#1a1b26",
			"base0D":         "#7aa2f7",
			"text-primary":   "#c0caf5",
			"syntax-keyword": "#bb9af7",
			"accent-primary": "#7aa2f7",
		},
	}
}

// randomOrderStylixTheme returns a StylixTheme with keys in non-alphabetical order.
func randomOrderStylixTheme() *ports.StylixTheme {
	return &ports.StylixTheme{
		Values: map[string]string{
			"zebra":  "#ffffff",
			"alpha":  "#000000",
			"middle": "#808080",
			"base00": "#1a1b26",
			"accent": "#7aa2f7",
		},
	}
}

// TestStylixGenerator_Interface verifies that the Stylix generator implements
// ports.Generator and returns the expected Name() and DefaultFilename().
func TestStylixGenerator_Interface(t *testing.T) {
	g := generator.NewStylix()

	// Compile-time interface check.
	var _ ports.Generator = g

	if name := g.Name(); name != "stylix" {
		t.Errorf("Name() = %q, want %q", name, "stylix")
	}

	if filename := g.DefaultFilename(); filename != "style.json" {
		t.Errorf("DefaultFilename() = %q, want %q", filename, "style.json")
	}
}

// TestStylixGenerator_ValidJSON verifies that Generate produces valid JSON
// that can be unmarshalled to map[string]string.
func TestStylixGenerator_ValidJSON(t *testing.T) {
	g := generator.NewStylix()
	theme := sampleStylixTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.Bytes()
	if len(output) == 0 {
		t.Fatal("Generate() produced empty output")
	}

	var result map[string]string
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, output)
	}

	// Verify all keys from the theme are present.
	for key, wantVal := range theme.Values {
		gotVal, ok := result[key]
		if !ok {
			t.Errorf("missing key %q in JSON output", key)
			continue
		}
		if gotVal != wantVal {
			t.Errorf("key %q = %q, want %q", key, gotVal, wantVal)
		}
	}
}

// TestStylixGenerator_SortedKeys verifies that JSON keys appear in
// alphabetical order in the output.
func TestStylixGenerator_SortedKeys(t *testing.T) {
	g := generator.NewStylix()
	theme := randomOrderStylixTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Extract keys in order of appearance from the JSON output.
	lines := strings.Split(output, "\n")
	var keys []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "\"") && strings.Contains(trimmed, ":") {
			// Extract the key portion between the first pair of quotes.
			endQuote := strings.Index(trimmed[1:], "\"")
			if endQuote >= 0 {
				key := trimmed[1 : endQuote+1]
				keys = append(keys, key)
			}
		}
	}

	if len(keys) == 0 {
		t.Fatal("no keys found in JSON output")
	}

	// Verify keys are sorted alphabetically.
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("keys not sorted: %q appears before %q", keys[i-1], keys[i])
		}
	}
}

// TestStylixGenerator_TwoSpaceIndent verifies that the JSON output uses
// 2-space indent (not tabs, not 4-space).
func TestStylixGenerator_TwoSpaceIndent(t *testing.T) {
	g := generator.NewStylix()
	theme := sampleStylixTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(output, "\n")

	foundIndented := false
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		// Check lines that start with whitespace (indented lines).
		if line[0] == ' ' || line[0] == '\t' {
			foundIndented = true

			// Must NOT contain tabs.
			if strings.Contains(line, "\t") {
				t.Errorf("line contains tab character: %q", line)
			}

			// Must start with exactly 2 spaces (for JSON object members).
			if !strings.HasPrefix(line, "  ") {
				t.Errorf("line does not start with 2-space indent: %q", line)
			}

			// Must NOT start with 4 spaces (would indicate 4-space indent).
			if strings.HasPrefix(line, "    ") {
				t.Errorf("line uses 4-space indent instead of 2-space: %q", line)
			}
		}
	}

	if !foundIndented {
		t.Error("no indented lines found in output; expected 2-space indented JSON")
	}
}

// TestStylixGenerator_HexColors verifies that all JSON values contain
// hex color strings with the # prefix.
func TestStylixGenerator_HexColors(t *testing.T) {
	g := generator.NewStylix()
	theme := &ports.StylixTheme{
		Values: map[string]string{
			"bg":     "#1a1b26",
			"fg":     "#c0caf5",
			"accent": "#7aa2f7",
		},
	}

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	for key, val := range result {
		if !strings.HasPrefix(val, "#") {
			t.Errorf("key %q has value %q without # prefix", key, val)
		}
		if len(val) != 7 {
			t.Errorf("key %q has value %q with length %d, want 7", key, val, len(val))
		}
	}
}

// TestStylixGenerator_WrongType verifies that passing a non-StylixTheme
// value as MappedTheme returns a GenerateError.
func TestStylixGenerator_WrongType(t *testing.T) {
	g := generator.NewStylix()

	var buf bytes.Buffer

	// Pass a string instead of *ports.StylixTheme.
	err := g.Generate(&buf, "not a stylix theme")
	if err == nil {
		t.Fatal("Generate() with wrong type should return error, got nil")
	}

	var genErr *domain.GenerateError
	if !errors.As(err, &genErr) {
		t.Errorf("error type = %T, want *domain.GenerateError", err)
	}
}

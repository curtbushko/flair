package fileio_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/adapters/wrappers"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

const testBgHex = "#1a1b26"

func TestWriteUniversal_Empty(t *testing.T) {
	ts := domain.NewTokenSet()
	var buf bytes.Buffer

	err := fileio.WriteUniversal(&buf, ts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Output should be valid YAML
	var uf ports.UniversalFile
	if err := yaml.Unmarshal(buf.Bytes(), &uf); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	// Tokens map should be empty (or nil, either is acceptable)
	if len(uf.Tokens) != 0 {
		t.Errorf("expected empty tokens map, got %d entries", len(uf.Tokens))
	}
}

func TestWriteUniversal_ColorOnly(t *testing.T) {
	ts := domain.NewTokenSet()
	bg, err := domain.ParseHex(testBgHex)
	if err != nil {
		t.Fatalf("failed to parse hex: %v", err)
	}
	ts.Set("surface.background", domain.Token{Color: bg})

	var buf bytes.Buffer
	if err := fileio.WriteUniversal(&buf, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var uf ports.UniversalFile
	if err := yaml.Unmarshal(buf.Bytes(), &uf); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	tok, ok := uf.Tokens["surface.background"]
	if !ok {
		t.Fatalf("expected 'surface.background' in tokens, got keys: %v", tokenKeys(uf.Tokens))
	}

	if tok.Color != testBgHex {
		t.Errorf("expected color '#1a1b26', got %q", tok.Color)
	}

	// No style flags should be set
	if tok.Bold || tok.Italic || tok.Underline || tok.Undercurl || tok.Strikethrough {
		t.Errorf("expected no style flags, got bold=%v italic=%v underline=%v undercurl=%v strikethrough=%v",
			tok.Bold, tok.Italic, tok.Underline, tok.Undercurl, tok.Strikethrough)
	}
}

func TestWriteUniversal_WithStyles(t *testing.T) {
	ts := domain.NewTokenSet()
	c, err := domain.ParseHex("#565f89")
	if err != nil {
		t.Fatalf("failed to parse hex: %v", err)
	}
	ts.Set("syntax.comment", domain.Token{Color: c, Italic: true})

	var buf bytes.Buffer
	if err := fileio.WriteUniversal(&buf, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var uf ports.UniversalFile
	if err := yaml.Unmarshal(buf.Bytes(), &uf); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	tok, ok := uf.Tokens["syntax.comment"]
	if !ok {
		t.Fatalf("expected 'syntax.comment' in tokens, got keys: %v", tokenKeys(uf.Tokens))
	}

	if tok.Color != "#565f89" {
		t.Errorf("expected color '#565f89', got %q", tok.Color)
	}

	if !tok.Italic {
		t.Error("expected italic to be true")
	}
}

func TestWriteUniversal_NoneColor(t *testing.T) {
	ts := domain.NewTokenSet()
	ts.Set("markup.bold", domain.Token{
		Color: domain.NoneColor(),
		Bold:  true,
	})

	var buf bytes.Buffer
	if err := fileio.WriteUniversal(&buf, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var uf ports.UniversalFile
	if err := yaml.Unmarshal(buf.Bytes(), &uf); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	tok, ok := uf.Tokens["markup.bold"]
	if !ok {
		t.Fatalf("expected 'markup.bold' in tokens, got keys: %v", tokenKeys(uf.Tokens))
	}

	// NoneColor should result in empty color string
	if tok.Color != "" {
		t.Errorf("expected empty color for NoneColor token, got %q", tok.Color)
	}

	if !tok.Bold {
		t.Error("expected bold to be true")
	}
}

func TestWriteUniversal_AllPaths(t *testing.T) {
	ts := domain.NewTokenSet()

	paths := []string{
		"surface.background",
		"text.primary",
		"syntax.keyword",
		"status.error",
		"terminal.red",
	}

	for i, p := range paths {
		c := domain.Color{R: uint8(i * 50), G: uint8(i * 30), B: uint8(i * 20)}
		ts.Set(p, domain.Token{Color: c})
	}

	var buf bytes.Buffer
	if err := fileio.WriteUniversal(&buf, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var uf ports.UniversalFile
	if err := yaml.Unmarshal(buf.Bytes(), &uf); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	if len(uf.Tokens) != len(paths) {
		t.Errorf("expected %d tokens, got %d", len(paths), len(uf.Tokens))
	}

	for _, p := range paths {
		if _, ok := uf.Tokens[p]; !ok {
			t.Errorf("expected path %q in tokens map", p)
		}
	}

	// Verify output has sorted keys by checking raw YAML order
	output := buf.String()
	lines := strings.Split(output, "\n")
	var foundPaths []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Token paths appear as top-level keys under "tokens:" with 4-space indent
		for _, p := range paths {
			if strings.HasPrefix(trimmed, p+":") {
				foundPaths = append(foundPaths, p)
			}
		}
	}

	// Verify sorted order
	for i := 1; i < len(foundPaths); i++ {
		if foundPaths[i] < foundPaths[i-1] {
			t.Errorf("token paths not sorted: %v", foundPaths)
			break
		}
	}
}

func TestWriteUniversal_WithVersionedWriter(t *testing.T) {
	ts := domain.NewTokenSet()
	c, err := domain.ParseHex("#7aa2f7")
	if err != nil {
		t.Fatalf("failed to parse hex: %v", err)
	}
	ts.Set("accent.primary", domain.Token{Color: c})

	var buf bytes.Buffer
	vw := wrappers.NewVersionedWriter(&buf, domain.FileKindUniversal, "tokyonight")

	if err := fileio.WriteUniversal(vw, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Should start with schema_version header
	wantHeader := fmt.Sprintf("schema_version: %d\nkind: universal\ntheme_name: tokyonight\n",
		domain.CurrentVersion(domain.FileKindUniversal))
	if !strings.HasPrefix(output, wantHeader) {
		t.Errorf("output does not start with expected header.\ngot:\n%s\nwant prefix:\n%s", output, wantHeader)
	}

	// The full output (after header) should be valid YAML when parsed as UniversalFile
	// Parse just the part after the header
	afterHeader := strings.TrimPrefix(output, wantHeader)
	var uf ports.UniversalFile
	if err := yaml.Unmarshal([]byte(afterHeader), &uf); err != nil {
		t.Fatalf("post-header content is not valid YAML: %v\ncontent:\n%s", err, afterHeader)
	}

	tok, ok := uf.Tokens["accent.primary"]
	if !ok {
		t.Fatalf("expected 'accent.primary' in tokens")
	}
	if tok.Color != "#7aa2f7" {
		t.Errorf("expected color '#7aa2f7', got %q", tok.Color)
	}
}

// --- ReadUniversal tests ---

func TestReadUniversal_Valid(t *testing.T) {
	yamlData := `tokens:
  surface.background:
    color: "#1a1b26"
  syntax.keyword:
    color: "#bb9af7"
`
	ts, err := fileio.ReadUniversal(bytes.NewReader([]byte(yamlData)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ts.Len() != 2 {
		t.Fatalf("expected 2 tokens, got %d", ts.Len())
	}

	bg, ok := ts.Get("surface.background")
	if !ok {
		t.Fatal("expected 'surface.background' in token set")
	}
	if bg.Color.Hex() != testBgHex {
		t.Errorf("expected color '#1a1b26', got %q", bg.Color.Hex())
	}

	kw, ok := ts.Get("syntax.keyword")
	if !ok {
		t.Fatal("expected 'syntax.keyword' in token set")
	}
	if kw.Color.Hex() != "#bb9af7" {
		t.Errorf("expected color '#bb9af7', got %q", kw.Color.Hex())
	}
}

func TestReadUniversal_StyleFlags(t *testing.T) {
	yamlData := `tokens:
  syntax.comment:
    color: "#565f89"
    italic: true
  markup.bold:
    color: "#c0caf5"
    bold: true
  ui.underline:
    color: "#7aa2f7"
    underline: true
    undercurl: true
    strikethrough: true
`
	ts, err := fileio.ReadUniversal(bytes.NewReader([]byte(yamlData)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	comment, ok := ts.Get("syntax.comment")
	if !ok {
		t.Fatal("expected 'syntax.comment' in token set")
	}
	if !comment.Italic {
		t.Error("expected italic=true for syntax.comment")
	}
	if comment.Bold {
		t.Error("expected bold=false for syntax.comment")
	}

	bold, ok := ts.Get("markup.bold")
	if !ok {
		t.Fatal("expected 'markup.bold' in token set")
	}
	if !bold.Bold {
		t.Error("expected bold=true for markup.bold")
	}

	ul, ok := ts.Get("ui.underline")
	if !ok {
		t.Fatal("expected 'ui.underline' in token set")
	}
	if !ul.Underline {
		t.Error("expected underline=true for ui.underline")
	}
	if !ul.Undercurl {
		t.Error("expected undercurl=true for ui.underline")
	}
	if !ul.Strikethrough {
		t.Error("expected strikethrough=true for ui.underline")
	}
}

func TestUniversal_RoundTrip(t *testing.T) {
	// Build a TokenSet with mixed color-only and styled tokens.
	original := domain.NewTokenSet()

	bg, err := domain.ParseHex(testBgHex)
	if err != nil {
		t.Fatalf("failed to parse hex: %v", err)
	}
	original.Set("surface.background", domain.Token{Color: bg})

	comment, err := domain.ParseHex("#565f89")
	if err != nil {
		t.Fatalf("failed to parse hex: %v", err)
	}
	original.Set("syntax.comment", domain.Token{Color: comment, Italic: true})

	// NoneColor with style only
	original.Set("markup.bold", domain.Token{
		Color: domain.NoneColor(),
		Bold:  true,
	})

	// Write
	var buf bytes.Buffer
	writeErr := fileio.WriteUniversal(&buf, original)
	if writeErr != nil {
		t.Fatalf("WriteUniversal error: %v", writeErr)
	}

	// Read back
	restored, readErr := fileio.ReadUniversal(bytes.NewReader(buf.Bytes()))
	if readErr != nil {
		t.Fatalf("ReadUniversal error: %v", readErr)
	}

	// Compare
	if original.Len() != restored.Len() {
		t.Fatalf("token count mismatch: original=%d restored=%d", original.Len(), restored.Len())
	}

	for _, path := range original.Paths() {
		orig, _ := original.Get(path)
		rest, ok := restored.Get(path)
		if !ok {
			t.Errorf("path %q missing in restored TokenSet", path)
			continue
		}

		if orig.Color.IsNone != rest.Color.IsNone {
			t.Errorf("path %q IsNone mismatch: orig=%v rest=%v", path, orig.Color.IsNone, rest.Color.IsNone)
		}
		if !orig.Color.IsNone && orig.Color.Hex() != rest.Color.Hex() {
			t.Errorf("path %q color mismatch: orig=%s rest=%s", path, orig.Color.Hex(), rest.Color.Hex())
		}
		if orig.Bold != rest.Bold {
			t.Errorf("path %q bold mismatch: orig=%v rest=%v", path, orig.Bold, rest.Bold)
		}
		if orig.Italic != rest.Italic {
			t.Errorf("path %q italic mismatch: orig=%v rest=%v", path, orig.Italic, rest.Italic)
		}
		if orig.Underline != rest.Underline {
			t.Errorf("path %q underline mismatch: orig=%v rest=%v", path, orig.Underline, rest.Underline)
		}
		if orig.Undercurl != rest.Undercurl {
			t.Errorf("path %q undercurl mismatch: orig=%v rest=%v", path, orig.Undercurl, rest.Undercurl)
		}
		if orig.Strikethrough != rest.Strikethrough {
			t.Errorf("path %q strikethrough mismatch: orig=%v rest=%v", path, orig.Strikethrough, rest.Strikethrough)
		}
	}
}

func TestReadUniversal_VersionMismatch(t *testing.T) {
	yamlData := `schema_version: 99
kind: universal
theme_name: test
tokens:
  surface.background:
    color: "#1a1b26"
`
	vr := wrappers.NewValidatingReader(bytes.NewReader([]byte(yamlData)), domain.FileKindUniversal)

	_, err := fileio.ReadUniversal(vr)
	if err == nil {
		t.Fatal("expected error for version mismatch, got nil")
	}

	var schemaErr *domain.SchemaVersionError
	if !errors.As(err, &schemaErr) {
		t.Fatalf("expected SchemaVersionError, got %T: %v", err, err)
	}
	if !schemaErr.NeedsUpgrade {
		t.Error("expected NeedsUpgrade=true for schema_version 99")
	}
}

func TestReadUniversal_EmptyTokens(t *testing.T) {
	yamlData := `tokens: {}
`
	ts, err := fileio.ReadUniversal(bytes.NewReader([]byte(yamlData)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ts.Len() != 0 {
		t.Errorf("expected empty token set, got %d tokens", ts.Len())
	}
}

func TestReadUniversal_InvalidYAML(t *testing.T) {
	malformed := []byte(`{not valid yaml: [[[`)

	_, err := fileio.ReadUniversal(bytes.NewReader(malformed))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// tokenKeys is a test helper that returns the keys of a token map.
func tokenKeys(m map[string]ports.UniversalToken) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

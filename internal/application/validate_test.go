package application

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/curtbushko/flair/internal/adapters/wrappers"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/domain"
)

// mockValidateStore implements ports.ThemeStore for validate tests.
type mockValidateStore struct {
	files map[string][]byte // key: "theme/file"
}

func newMockValidateStore() *mockValidateStore {
	return &mockValidateStore{files: make(map[string][]byte)}
}

func (m *mockValidateStore) ConfigDir() string             { return "/mock" }
func (m *mockValidateStore) ThemeDir(name string) string   { return "/mock/" + name }
func (m *mockValidateStore) EnsureThemeDir(string) error   { return nil }
func (m *mockValidateStore) ListThemes() ([]string, error) { return nil, nil }
func (m *mockValidateStore) SelectedTheme() (string, error) {
	return "", nil
}
func (m *mockValidateStore) Select(string) error { return nil }
func (m *mockValidateStore) FileMtime(string, string) (time.Time, error) {
	return time.Time{}, errors.New("not implemented")
}

func (m *mockValidateStore) putFile(theme, file string, data []byte) {
	m.files[theme+"/"+file] = data
}

func (m *mockValidateStore) FileExists(theme, file string) bool {
	_, ok := m.files[theme+"/"+file]
	return ok
}

func (m *mockValidateStore) OpenReader(theme, file string) (io.ReadCloser, error) {
	data, ok := m.files[theme+"/"+file]
	if !ok {
		return nil, fmt.Errorf("file not found: %s/%s", theme, file)
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *mockValidateStore) OpenWriter(theme, file string) (io.WriteCloser, error) {
	return nil, errors.New("not implemented")
}

// testSchemaValidator uses the real ValidatingReader to check schema versions.
func testSchemaValidator(r io.Reader, kind domain.FileKind) error {
	vr := wrappers.NewValidatingReader(r, kind)
	_, err := io.ReadAll(vr)
	return err
}

// validPaletteYAML returns a well-formed palette YAML with all 24 colors.
func validPaletteYAML() []byte {
	return []byte(`schema_version: 1
kind: palette
theme_name: test-theme
system: base24
author: test
variant: dark
palette:
  base00: "16161e"
  base01: "1a1b26"
  base02: "2f3549"
  base03: "444b6a"
  base04: "787c99"
  base05: "a9b1d6"
  base06: "cbccd1"
  base07: "d5d6db"
  base08: "c0caf5"
  base09: "a9b1d6"
  base0A: "0db9d7"
  base0B: "9ece6a"
  base0C: "b4f9f8"
  base0D: "2ac3de"
  base0E: "bb9af7"
  base0F: "f7768e"
  base10: "1a1b26"
  base11: "1a1b26"
  base12: "c0caf5"
  base13: "0db9d7"
  base14: "9ece6a"
  base15: "b4f9f8"
  base16: "2ac3de"
  base17: "bb9af7"
`)
}

// newParser returns a real YAML palette parser for tests.
func newParser() *yamlparser.Parser {
	return yamlparser.NewParser()
}

// populateAllFiles adds all expected theme files to the mock store.
func populateAllFiles(store *mockValidateStore, theme string) {
	store.putFile(theme, "universal.yaml", []byte("schema_version: 1\nkind: universal\n"))
	for _, f := range outputFiles {
		store.putFile(theme, f, []byte("content"))
	}
	for _, mf := range []string{
		"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml",
		"qss-mapping.yaml", "stylix-mapping.yaml",
	} {
		store.putFile(theme, mf, []byte("content"))
	}
}

func TestValidateThemeUseCase_ValidTheme(t *testing.T) {
	store := newMockValidateStore()
	store.putFile("my-theme", "palette.yaml", validPaletteYAML())
	populateAllFiles(store, "my-theme")

	uc := NewValidateThemeUseCase(store, newParser(), testSchemaValidator)
	violations, err := uc.Execute("my-theme")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestValidateThemeUseCase_MissingPalette(t *testing.T) {
	store := newMockValidateStore()
	// No palette.yaml at all

	uc := NewValidateThemeUseCase(store, newParser(), testSchemaValidator)
	violations, err := uc.Execute("my-theme")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) == 0 {
		t.Fatal("expected violations for missing palette, got none")
	}

	found := false
	for _, v := range violations {
		if strings.Contains(v, "palette.yaml") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation mentioning palette.yaml, got: %v", violations)
	}
}

func TestValidateThemeUseCase_InvalidPalette(t *testing.T) {
	store := newMockValidateStore()
	// Palette with bad luminance ordering: dark variant with white bg (#ffffff), black fg (#000000)
	badPalette := []byte(`schema_version: 1
kind: palette
theme_name: test-theme
system: base24
author: test
variant: dark
palette:
  base00: "ffffff"
  base01: "ffffff"
  base02: "ffffff"
  base03: "ffffff"
  base04: "ffffff"
  base05: "000000"
  base06: "ffffff"
  base07: "ffffff"
  base08: "ffffff"
  base09: "ffffff"
  base0A: "ffffff"
  base0B: "ffffff"
  base0C: "ffffff"
  base0D: "ffffff"
  base0E: "ffffff"
  base0F: "ffffff"
  base10: "ffffff"
  base11: "ffffff"
  base12: "ffffff"
  base13: "ffffff"
  base14: "ffffff"
  base15: "ffffff"
  base16: "ffffff"
  base17: "ffffff"
`)
	store.putFile("my-theme", "palette.yaml", badPalette)
	populateAllFiles(store, "my-theme")

	uc := NewValidateThemeUseCase(store, newParser(), testSchemaValidator)
	violations, err := uc.Execute("my-theme")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) == 0 {
		t.Fatal("expected violations from palette validation, got none")
	}

	// Should have luminance ordering violation (dark theme with white bg, black fg)
	found := false
	for _, v := range violations {
		if strings.Contains(v, "luminance") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected luminance violation, got: %v", violations)
	}
}

func TestValidateThemeUseCase_MissingFiles(t *testing.T) {
	store := newMockValidateStore()
	store.putFile("my-theme", "palette.yaml", validPaletteYAML())
	// Missing universal.yaml and output files

	uc := NewValidateThemeUseCase(store, newParser(), testSchemaValidator)
	violations, err := uc.Execute("my-theme")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) == 0 {
		t.Fatal("expected violations for missing files, got none")
	}

	// Should mention universal.yaml
	found := false
	for _, v := range violations {
		if strings.Contains(v, "universal.yaml") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation mentioning universal.yaml, got: %v", violations)
	}
}

func TestValidateThemeUseCase_SchemaVersionMismatch(t *testing.T) {
	store := newMockValidateStore()
	// Palette with wrong schema version
	badVersionPalette := []byte(`schema_version: 99
kind: palette
theme_name: test-theme
system: base24
author: test
variant: dark
palette:
  base00: "1a1b26"
  base01: "16161e"
  base02: "2f3549"
  base03: "444b6a"
  base04: "787c99"
  base05: "a9b1d6"
  base06: "cbccd1"
  base07: "d5d6db"
  base08: "c0caf5"
  base09: "a9b1d6"
  base0A: "0db9d7"
  base0B: "9ece6a"
  base0C: "b4f9f8"
  base0D: "2ac3de"
  base0E: "bb9af7"
  base0F: "f7768e"
  base10: "1a1b26"
  base11: "1a1b26"
  base12: "c0caf5"
  base13: "0db9d7"
  base14: "9ece6a"
  base15: "b4f9f8"
  base16: "2ac3de"
  base17: "bb9af7"
`)
	store.putFile("my-theme", "palette.yaml", badVersionPalette)
	populateAllFiles(store, "my-theme")

	uc := NewValidateThemeUseCase(store, newParser(), testSchemaValidator)
	violations, err := uc.Execute("my-theme")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) == 0 {
		t.Fatal("expected violations for schema version mismatch, got none")
	}

	found := false
	for _, v := range violations {
		if strings.Contains(v, "schema") || strings.Contains(v, "version") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected schema version violation, got: %v", violations)
	}
}

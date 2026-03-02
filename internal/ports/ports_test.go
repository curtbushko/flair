package ports_test

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// TestCompileCheck_PaletteParser verifies the PaletteParser interface compiles
// with the correct method signature: Parse(io.Reader) (*domain.Palette, error).
func TestCompileCheck_PaletteParser(t *testing.T) {
	var _ ports.PaletteParser = (*mockPaletteParser)(nil)
}

type mockPaletteParser struct{}

func (m *mockPaletteParser) Parse(r io.Reader) (*domain.Palette, error) {
	return nil, nil
}

// TestCompileCheck_PaletteSource verifies the PaletteSource interface compiles
// with List, Get, Has methods.
func TestCompileCheck_PaletteSource(t *testing.T) {
	var _ ports.PaletteSource = (*mockPaletteSource)(nil)
}

type mockPaletteSource struct{}

func (m *mockPaletteSource) List() []string                     { return nil }
func (m *mockPaletteSource) Get(name string) (io.Reader, error) { return nil, nil }
func (m *mockPaletteSource) Has(name string) bool               { return false }

// TestCompileCheck_Tokenizer verifies the Tokenizer interface compiles
// with Tokenize(*domain.Palette) *domain.TokenSet.
func TestCompileCheck_Tokenizer(t *testing.T) {
	var _ ports.Tokenizer = (*mockTokenizer)(nil)
}

type mockTokenizer struct{}

func (m *mockTokenizer) Tokenize(p *domain.Palette) *domain.TokenSet {
	return nil
}

// TestCompileCheck_Mapper verifies the Mapper interface compiles with
// Name() string and Map(*domain.ResolvedTheme) (ports.MappedTheme, error).
func TestCompileCheck_Mapper(t *testing.T) {
	var _ ports.Mapper = (*mockMapper)(nil)
}

type mockMapper struct{}

func (m *mockMapper) Name() string { return "" }
func (m *mockMapper) Map(theme *domain.ResolvedTheme) (ports.MappedTheme, error) {
	return nil, nil
}

// TestCompileCheck_Generator verifies the Generator interface compiles with
// Name, DefaultFilename, Generate methods.
func TestCompileCheck_Generator(t *testing.T) {
	var _ ports.Generator = (*mockGenerator)(nil)
}

type mockGenerator struct{}

func (m *mockGenerator) Name() string            { return "" }
func (m *mockGenerator) DefaultFilename() string { return "" }
func (m *mockGenerator) Generate(w io.Writer, mapped ports.MappedTheme) error {
	return nil
}

// TestCompileCheck_ThemeStore verifies the ThemeStore interface compiles with
// all 10 methods.
func TestCompileCheck_ThemeStore(t *testing.T) {
	var _ ports.ThemeStore = (*mockThemeStore)(nil)
}

type mockThemeStore struct{}

func (m *mockThemeStore) ConfigDir() string                     { return "" }
func (m *mockThemeStore) ThemeDir(themeName string) string      { return "" }
func (m *mockThemeStore) EnsureThemeDir(themeName string) error { return nil }
func (m *mockThemeStore) ListThemes() ([]string, error)         { return nil, nil }
func (m *mockThemeStore) SelectedTheme() (string, error)        { return "", nil }
func (m *mockThemeStore) Select(themeName string) error         { return nil }
func (m *mockThemeStore) OpenReader(themeName, filename string) (io.ReadCloser, error) {
	return nil, nil
}
func (m *mockThemeStore) OpenWriter(themeName, filename string) (io.WriteCloser, error) {
	return nil, nil
}
func (m *mockThemeStore) FileExists(themeName, filename string) bool { return false }
func (m *mockThemeStore) FileMtime(themeName, filename string) (time.Time, error) {
	return time.Time{}, nil
}

// TestCompileCheck_TargetStruct verifies the Target struct pairs Mapper + Generator + MappingFile.
func TestCompileCheck_TargetStruct(t *testing.T) {
	m := ports.Mapper(&mockMapper{})
	gen := ports.Generator(&mockGenerator{})
	target := ports.Target{
		Mapper:          m,
		Generator:       gen,
		MappingFile:     "vim-mapping.yaml",
		MappingFileKind: domain.FileKindVimMapping,
		WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
			return nil
		},
	}

	if target.MappingFile != "vim-mapping.yaml" {
		t.Errorf("expected MappingFile = %q, got %q", "vim-mapping.yaml", target.MappingFile)
	}
	if target.Mapper.Name() != "" {
		t.Errorf("expected Mapper.Name() = %q, got %q", "", target.Mapper.Name())
	}
	if target.Generator.Name() != "" {
		t.Errorf("expected Generator.Name() = %q, got %q", "", target.Generator.Name())
	}
	if target.MappingFileKind != domain.FileKindVimMapping {
		t.Errorf("expected MappingFileKind = %q, got %q", domain.FileKindVimMapping, target.MappingFileKind)
	}
	if target.WriteMappingFile == nil {
		t.Error("expected WriteMappingFile to be set")
	}
}

// assertAllFieldsReadable uses reflect to read every field of a struct, ensuring
// the govet unusedwrite checker sees all fields as used.
func assertAllFieldsReadable(t *testing.T, name string, v any) {
	t.Helper()
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fv := rv.Field(i)
		if !fv.IsValid() {
			t.Errorf("%s.%s is not valid", name, field.Name)
		}
	}
}

// TestCompileCheck_FileStructs verifies all file structs compile with yaml tags.
func TestCompileCheck_FileStructs(t *testing.T) {
	// FileHeader
	h := ports.FileHeader{
		SchemaVersion: 1,
		Kind:          domain.FileKindPalette,
		ThemeName:     "test",
	}
	assertAllFieldsReadable(t, "FileHeader", h)

	// PaletteFile
	pf := ports.PaletteFile{
		FileHeader: ports.FileHeader{SchemaVersion: 1, Kind: domain.FileKindPalette, ThemeName: "t"},
		System:     "base24",
		Author:     "author",
		Variant:    "dark",
		Palette:    map[string]string{"base00": "1a1b26"},
	}
	assertAllFieldsReadable(t, "PaletteFile", pf)

	// TokenEntry
	te := ports.TokenEntry{
		Color:         "#7aa2f7",
		Bold:          true,
		Italic:        true,
		Underline:     true,
		Undercurl:     true,
		Strikethrough: true,
	}
	assertAllFieldsReadable(t, "TokenEntry", te)

	// TokensFile
	tf := ports.TokensFile{
		FileHeader: ports.FileHeader{SchemaVersion: 1, Kind: domain.FileKindTokens, ThemeName: "t"},
		Tokens:     map[string]ports.TokenEntry{"text.primary": te},
	}
	assertAllFieldsReadable(t, "TokensFile", tf)

	// VimMappingHighlight
	vmh := ports.VimMappingHighlight{
		Fg: "#ff0000", Bg: "#000000", Sp: "#00ff00",
		Bold: true, Italic: true, Underline: true, Undercurl: true,
		Strikethrough: true, Reverse: true, Nocombine: true,
		Link: "Normal",
	}
	assertAllFieldsReadable(t, "VimMappingHighlight", vmh)

	// VimMappingFile
	vmf := ports.VimMappingFile{
		FileHeader:     ports.FileHeader{SchemaVersion: 1, Kind: domain.FileKindVimMapping, ThemeName: "t"},
		Highlights:     map[string]ports.VimMappingHighlight{"Normal": vmh},
		TerminalColors: [16]string{"#000000"},
	}
	assertAllFieldsReadable(t, "VimMappingFile", vmf)

	// CSSRuleEntry
	cre := ports.CSSRuleEntry{
		Selector:   "body",
		Properties: map[string]string{"color": "#fff"},
	}
	assertAllFieldsReadable(t, "CSSRuleEntry", cre)

	// CSSMappingFile
	cmf := ports.CSSMappingFile{
		FileHeader:       ports.FileHeader{SchemaVersion: 1, Kind: domain.FileKindCSSMapping, ThemeName: "t"},
		CustomProperties: map[string]string{"--bg": "#000"},
		Rules:            []ports.CSSRuleEntry{cre},
	}
	assertAllFieldsReadable(t, "CSSMappingFile", cmf)

	// GtkMappingFile
	gmf := ports.GtkMappingFile{
		FileHeader: ports.FileHeader{SchemaVersion: 1, Kind: domain.FileKindGtkMapping, ThemeName: "t"},
		Colors:     map[string]string{"window_bg_color": "#1a1b26"},
		Rules:      []ports.CSSRuleEntry{cre},
	}
	assertAllFieldsReadable(t, "GtkMappingFile", gmf)

	// QssMappingFile
	qmf := ports.QssMappingFile{
		FileHeader: ports.FileHeader{SchemaVersion: 1, Kind: domain.FileKindQssMapping, ThemeName: "t"},
		Rules:      []ports.CSSRuleEntry{cre},
	}
	assertAllFieldsReadable(t, "QssMappingFile", qmf)

	// StylixMappingFile
	smf := ports.StylixMappingFile{
		FileHeader: ports.FileHeader{SchemaVersion: 1, Kind: domain.FileKindStylixMapping, ThemeName: "t"},
		Values:     map[string]string{"bg": "#1a1b26"},
	}
	assertAllFieldsReadable(t, "StylixMappingFile", smf)
}

// TestCompileCheck_ThemeStructs verifies theme DTO structs reference domain types correctly.
func TestCompileCheck_ThemeStructs(t *testing.T) {
	red := domain.Color{R: 255, G: 0, B: 0}
	blue := domain.Color{R: 0, G: 0, B: 255}

	// VimHighlight uses *domain.Color
	vh := ports.VimHighlight{
		Fg:            &red,
		Bg:            &blue,
		Sp:            nil,
		Bold:          true,
		Italic:        true,
		Underline:     true,
		Undercurl:     true,
		Strikethrough: true,
		Reverse:       true,
		Nocombine:     true,
		Link:          "Normal",
	}
	if vh.Fg.R != 255 {
		t.Errorf("VimHighlight.Fg.R = %d, want 255", vh.Fg.R)
	}
	if vh.Bg.B != 255 {
		t.Errorf("VimHighlight.Bg.B = %d, want 255", vh.Bg.B)
	}
	assertAllFieldsReadable(t, "VimHighlight", vh)

	// VimTheme
	vt := ports.VimTheme{
		Name:           "test",
		Highlights:     map[string]ports.VimHighlight{"Normal": vh},
		TerminalColors: [16]domain.Color{red},
	}
	assertAllFieldsReadable(t, "VimTheme", vt)
	if vt.Name != "test" {
		t.Errorf("VimTheme.Name = %q, want %q", vt.Name, "test")
	}

	// StylixTheme
	st := ports.StylixTheme{
		Values: map[string]string{"bg": "#1a1b26"},
	}
	assertAllFieldsReadable(t, "StylixTheme", st)

	// CSSProperty
	cp := ports.CSSProperty{Property: "color", Value: "#fff"}
	assertAllFieldsReadable(t, "CSSProperty", cp)

	// CSSRule
	cr := ports.CSSRule{
		Selector:   "body",
		Properties: []ports.CSSProperty{cp},
	}
	assertAllFieldsReadable(t, "CSSRule", cr)

	// CSSTheme
	ct := ports.CSSTheme{
		CustomProperties: map[string]string{"--bg": "#000"},
		Rules:            []ports.CSSRule{cr},
	}
	assertAllFieldsReadable(t, "CSSTheme", ct)

	// GtkColorDef
	gcd := ports.GtkColorDef{Name: "window_bg_color", Value: "#1a1b26"}
	assertAllFieldsReadable(t, "GtkColorDef", gcd)

	// GtkTheme
	gt := ports.GtkTheme{
		Colors: []ports.GtkColorDef{gcd},
		Rules:  []ports.CSSRule{cr},
	}
	assertAllFieldsReadable(t, "GtkTheme", gt)

	// QssTheme
	qt := ports.QssTheme{
		Rules: []ports.CSSRule{cr},
	}
	assertAllFieldsReadable(t, "QssTheme", qt)
}

// TestFileStructsHaveYamlTags verifies all file struct fields have yaml tags.
func TestFileStructsHaveYamlTags(t *testing.T) {
	structsToCheck := []struct {
		name string
		val  any
	}{
		{"FileHeader", ports.FileHeader{}},
		{"PaletteFile", ports.PaletteFile{}},
		{"TokenEntry", ports.TokenEntry{}},
		{"TokensFile", ports.TokensFile{}},
		{"VimMappingHighlight", ports.VimMappingHighlight{}},
		{"VimMappingFile", ports.VimMappingFile{}},
		{"CSSRuleEntry", ports.CSSRuleEntry{}},
		{"CSSMappingFile", ports.CSSMappingFile{}},
		{"GtkMappingFile", ports.GtkMappingFile{}},
		{"QssMappingFile", ports.QssMappingFile{}},
		{"StylixMappingFile", ports.StylixMappingFile{}},
	}

	for _, sc := range structsToCheck {
		t.Run(sc.name, func(t *testing.T) {
			rt := reflect.TypeOf(sc.val)
			for i := 0; i < rt.NumField(); i++ {
				field := rt.Field(i)
				// Skip embedded fields (like FileHeader) - they use inline yaml
				if field.Anonymous {
					continue
				}
				tag := field.Tag.Get("yaml")
				if tag == "" {
					t.Errorf("field %s.%s has no yaml tag", sc.name, field.Name)
				}
			}
		})
	}
}

// TestMappedThemeIsAny verifies MappedTheme is an alias for any.
func TestMappedThemeIsAny(t *testing.T) {
	var m ports.MappedTheme
	// Should be able to assign any type
	m = "string"
	if m != "string" {
		t.Errorf("MappedTheme should accept string")
	}
	m = 42
	if m != 42 {
		t.Errorf("MappedTheme should accept int")
	}
	m = ports.VimTheme{}
	if _, ok := m.(ports.VimTheme); !ok {
		t.Errorf("MappedTheme should accept VimTheme")
	}
}

// TestFileHeaderFieldNames verifies FileHeader has the expected yaml field names.
func TestFileHeaderFieldNames(t *testing.T) {
	rt := reflect.TypeOf(ports.FileHeader{})

	tests := []struct {
		fieldName   string
		expectedTag string
	}{
		{"SchemaVersion", "schema_version"},
		{"Kind", "kind"},
		{"ThemeName", "theme_name"},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			field, ok := rt.FieldByName(tt.fieldName)
			if !ok {
				t.Fatalf("field %s not found", tt.fieldName)
			}
			tag := field.Tag.Get("yaml")
			if !strings.HasPrefix(tag, tt.expectedTag) {
				t.Errorf("field %s yaml tag = %q, want prefix %q", tt.fieldName, tag, tt.expectedTag)
			}
		})
	}
}

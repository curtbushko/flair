package application_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/curtbushko/flair/internal/adapters/deriver"
	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/adapters/generator"
	"github.com/curtbushko/flair/internal/adapters/mapper"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/application"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// --- Stub implementations ---

// stubPaletteSource is a test stub for ports.PaletteSource.
type stubPaletteSource struct {
	palettes map[string][]byte
}

func newStubPaletteSource() *stubPaletteSource {
	return &stubPaletteSource{palettes: make(map[string][]byte)}
}

func (s *stubPaletteSource) List() []string {
	names := make([]string, 0, len(s.palettes))
	for k := range s.palettes {
		names = append(names, k)
	}
	return names
}

func (s *stubPaletteSource) Get(name string) (io.Reader, error) {
	data, ok := s.palettes[name]
	if !ok {
		return nil, fmt.Errorf("palette %q not found", name)
	}
	return bytes.NewReader(data), nil
}

func (s *stubPaletteSource) Has(name string) bool {
	_, ok := s.palettes[name]
	return ok
}

// recordedWrite captures what was written to a specific file.
type recordedWrite struct {
	data []byte
}

// stubThemeStore is an in-memory ThemeStore that records all writes.
type stubThemeStore struct {
	mu             sync.Mutex
	files          map[string]map[string]*recordedWrite // themeName -> filename -> data
	ensureDirCalls []string
}

func newStubThemeStore() *stubThemeStore {
	return &stubThemeStore{
		files: make(map[string]map[string]*recordedWrite),
	}
}

func (s *stubThemeStore) ConfigDir() string { return "/fake/config" }

func (s *stubThemeStore) ThemeDir(themeName string) string {
	return "/fake/config/" + themeName
}

func (s *stubThemeStore) EnsureThemeDir(themeName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureDirCalls = append(s.ensureDirCalls, themeName)
	if s.files[themeName] == nil {
		s.files[themeName] = make(map[string]*recordedWrite)
	}
	return nil
}

func (s *stubThemeStore) ListThemes() ([]string, error)  { return nil, nil }
func (s *stubThemeStore) SelectedTheme() (string, error) { return "", nil }
func (s *stubThemeStore) Select(_ string) error          { return nil }

func (s *stubThemeStore) OpenReader(themeName, filename string) (io.ReadCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if theme, ok := s.files[themeName]; ok {
		if rw, ok := theme[filename]; ok {
			return io.NopCloser(bytes.NewReader(rw.data)), nil
		}
	}
	return nil, fmt.Errorf("file not found: %s/%s", themeName, filename)
}

func (s *stubThemeStore) OpenWriter(themeName, filename string) (io.WriteCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.files[themeName] == nil {
		s.files[themeName] = make(map[string]*recordedWrite)
	}
	rw := &recordedWrite{}
	s.files[themeName][filename] = rw
	return &bufWriteCloser{rw: rw}, nil
}

func (s *stubThemeStore) FileExists(themeName, filename string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if theme, ok := s.files[themeName]; ok {
		_, ok := theme[filename]
		return ok
	}
	return false
}

func (s *stubThemeStore) FileMtime(_, _ string) (time.Time, error) {
	return time.Now(), nil
}

func (s *stubThemeStore) writtenFiles(themeName string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var names []string
	if theme, ok := s.files[themeName]; ok {
		for name := range theme {
			names = append(names, name)
		}
	}
	return names
}

func (s *stubThemeStore) hasFile(themeName, filename string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if theme, ok := s.files[themeName]; ok {
		_, ok := theme[filename]
		return ok
	}
	return false
}

// bufWriteCloser wraps a recordedWrite to implement io.WriteCloser.
type bufWriteCloser struct {
	buf bytes.Buffer
	rw  *recordedWrite
}

func (bwc *bufWriteCloser) Write(p []byte) (int, error) {
	return bwc.buf.Write(p)
}

func (bwc *bufWriteCloser) Close() error {
	bwc.rw.data = bwc.buf.Bytes()
	return nil
}

// stubGenParser is a test stub for ports.PaletteParser used in generate tests.
type stubGenParser struct {
	palette *domain.Palette
	err     error
}

func (s *stubGenParser) Parse(_ io.Reader) (*domain.Palette, error) {
	return s.palette, s.err
}

// stubGenDeriver is a test stub for ports.TokenDeriver used in generate tests.
type stubGenDeriver struct {
	tokenSet *domain.TokenSet
}

func (s *stubGenDeriver) Derive(_ *domain.Palette) *domain.TokenSet {
	return s.tokenSet
}

// stubGenMapper is a test stub for ports.Mapper.
type stubGenMapper struct {
	name   string
	result ports.MappedTheme
	err    error
}

func (s *stubGenMapper) Name() string { return s.name }
func (s *stubGenMapper) Map(_ *domain.ResolvedTheme) (ports.MappedTheme, error) {
	return s.result, s.err
}

// stubGenGenerator is a test stub for ports.Generator.
type stubGenGenerator struct {
	name            string
	defaultFilename string
	err             error
}

func (s *stubGenGenerator) Name() string            { return s.name }
func (s *stubGenGenerator) DefaultFilename() string { return s.defaultFilename }
func (s *stubGenGenerator) Generate(w io.Writer, _ ports.MappedTheme) error {
	if s.err != nil {
		return s.err
	}
	_, err := w.Write([]byte("generated-" + s.name))
	return err
}

// --- Helpers ---

func makeGenPalette(t *testing.T) *domain.Palette {
	t.Helper()
	colors := map[string]string{
		"base00": "1a1b26", "base01": "1f2335", "base02": "292e42", "base03": "565f89",
		"base04": "a9b1d6", "base05": "c0caf5", "base06": "c0caf5", "base07": "c8d3f5",
		"base08": "f7768e", "base09": "ff9e64", "base0A": "e0af68", "base0B": "9ece6a",
		"base0C": "7dcfff", "base0D": "7aa2f7", "base0E": "bb9af7", "base0F": "db4b4b",
		"base10": "16161e", "base11": "101014", "base12": "ff899d", "base13": "e9c582",
		"base14": "afd67a", "base15": "97d8f8", "base16": "8db6fa", "base17": "c8acf8",
	}
	pal, err := domain.NewPalette("Tokyo Night Dark", "test-author", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("makeGenPalette: %v", err)
	}
	return pal
}

func makeGenTokenSet(t *testing.T, pal *domain.Palette) *domain.TokenSet {
	t.Helper()
	ts := domain.NewTokenSet()
	ts.Set("surface.background", domain.Token{Color: pal.Base(0x00)})
	ts.Set("text.primary", domain.Token{Color: pal.Base(0x05)})
	return ts
}

func makeStubTargets() []ports.Target {
	names := []string{"vim", "css", "gtk", "qss", "stylix"}
	mappingFiles := []string{
		"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml",
		"qss-mapping.yaml", "stylix-mapping.yaml",
	}
	outputFiles := []string{
		"style.lua", "style.css", "gtk.css",
		"style.qss", "style.json",
	}
	mappingKinds := []domain.FileKind{
		domain.FileKindVimMapping, domain.FileKindCSSMapping,
		domain.FileKindGtkMapping, domain.FileKindQssMapping,
		domain.FileKindStylixMapping,
	}

	targets := make([]ports.Target, len(names))
	for i, name := range names {
		targets[i] = ports.Target{
			Mapper: &stubGenMapper{
				name:   name,
				result: &ports.StylixTheme{Values: map[string]string{"key": "val"}},
			},
			Generator: &stubGenGenerator{
				name:            name,
				defaultFilename: outputFiles[i],
			},
			MappingFile:     mappingFiles[i],
			MappingFileKind: mappingKinds[i],
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				_, err := w.Write([]byte("mapping-data"))
				return err
			},
		}
	}
	return targets
}

// --- Tests ---

func TestGenerateTheme_FullPipeline(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newStubThemeStore()
	builtins := newStubPaletteSource()

	targets := makeStubTargets()

	uc := application.NewGenerateThemeUseCase(
		&stubGenParser{palette: pal},
		&stubGenDeriver{tokenSet: ts},
		targets,
		store,
		builtins,
	)

	err := uc.Execute(bytes.NewReader([]byte("palette-yaml")), "test-theme", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have written: palette.yaml, universal.yaml, 5 mapping files, 5 output files = 12 total
	written := store.writtenFiles("test-theme")
	if len(written) != 12 {
		t.Errorf("expected 12 files written, got %d: %v", len(written), written)
	}

	expectedFiles := []string{
		"palette.yaml", "universal.yaml",
		"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml",
		"qss-mapping.yaml", "stylix-mapping.yaml",
		"style.lua", "style.css", "gtk.css", "style.qss", "style.json",
	}
	for _, f := range expectedFiles {
		if !store.hasFile("test-theme", f) {
			t.Errorf("expected file %q to be written", f)
		}
	}
}

func TestGenerateTheme_BuiltinName(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newStubThemeStore()

	builtins := newStubPaletteSource()
	builtins.palettes["tokyo-night-dark"] = []byte("palette-yaml")

	targets := makeStubTargets()

	uc := application.NewGenerateThemeUseCase(
		&stubGenParser{palette: pal},
		&stubGenDeriver{tokenSet: ts},
		targets,
		store,
		builtins,
	)

	// When paletteRef is a built-in name and themeName is empty,
	// the use case should resolve the palette from builtins and infer theme name.
	err := uc.ExecuteBuiltin("tokyo-night-dark", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Theme name should be inferred as "tokyo-night-dark"
	if len(store.ensureDirCalls) == 0 {
		t.Fatal("expected EnsureThemeDir to be called")
	}
	if store.ensureDirCalls[0] != "tokyo-night-dark" {
		t.Errorf("EnsureThemeDir called with %q, want %q", store.ensureDirCalls[0], "tokyo-night-dark")
	}

	// Should have written files under "tokyo-night-dark"
	written := store.writtenFiles("tokyo-night-dark")
	if len(written) != 12 {
		t.Errorf("expected 12 files written, got %d: %v", len(written), written)
	}
}

func TestGenerateTheme_TargetFilter(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newStubThemeStore()
	builtins := newStubPaletteSource()

	targets := makeStubTargets()

	uc := application.NewGenerateThemeUseCase(
		&stubGenParser{palette: pal},
		&stubGenDeriver{tokenSet: ts},
		targets,
		store,
		builtins,
	)

	err := uc.Execute(bytes.NewReader([]byte("palette-yaml")), "test-theme", "stylix")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have: palette.yaml, universal.yaml, stylix-mapping.yaml, style.json = 4 files
	written := store.writtenFiles("test-theme")
	if len(written) != 4 {
		t.Errorf("expected 4 files written, got %d: %v", len(written), written)
	}

	mustHave := []string{"palette.yaml", "universal.yaml", "stylix-mapping.yaml", "style.json"}
	for _, f := range mustHave {
		if !store.hasFile("test-theme", f) {
			t.Errorf("expected file %q to be written", f)
		}
	}

	// Should NOT have other target files
	mustNotHave := []string{"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml", "qss-mapping.yaml"}
	for _, f := range mustNotHave {
		if store.hasFile("test-theme", f) {
			t.Errorf("file %q should not be written when filtering to stylix", f)
		}
	}
}

func TestGenerateTheme_CreateDir(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newStubThemeStore()
	builtins := newStubPaletteSource()

	targets := makeStubTargets()

	uc := application.NewGenerateThemeUseCase(
		&stubGenParser{palette: pal},
		&stubGenDeriver{tokenSet: ts},
		targets,
		store,
		builtins,
	)

	err := uc.Execute(bytes.NewReader([]byte("palette-yaml")), "my-theme", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(store.ensureDirCalls) == 0 {
		t.Fatal("expected EnsureThemeDir to be called")
	}
	if store.ensureDirCalls[0] != "my-theme" {
		t.Errorf("EnsureThemeDir called with %q, want %q", store.ensureDirCalls[0], "my-theme")
	}
}

func TestGenerateTheme_ParseError(t *testing.T) {
	parseErr := errors.New("bad palette yaml")
	store := newStubThemeStore()
	builtins := newStubPaletteSource()

	targets := makeStubTargets()

	uc := application.NewGenerateThemeUseCase(
		&stubGenParser{err: parseErr},
		&stubGenDeriver{tokenSet: domain.NewTokenSet()},
		targets,
		store,
		builtins,
	)

	err := uc.Execute(bytes.NewReader([]byte("bad")), "test-theme", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, parseErr) {
		t.Errorf("error = %v, want wrapped %v", err, parseErr)
	}
}

func TestGenerateTheme_MapperError(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newStubThemeStore()
	builtins := newStubPaletteSource()

	// Create targets where one mapper fails
	targets := makeStubTargets()
	targets[1] = ports.Target{
		Mapper: &stubGenMapper{
			name: "css",
			err:  errors.New("css mapper broke"),
		},
		Generator: &stubGenGenerator{
			name:            "css",
			defaultFilename: "style.css",
		},
		MappingFile:     "css-mapping.yaml",
		MappingFileKind: domain.FileKindCSSMapping,
		WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
			_, err := w.Write([]byte("mapping-data"))
			return err
		},
	}

	uc := application.NewGenerateThemeUseCase(
		&stubGenParser{palette: pal},
		&stubGenDeriver{tokenSet: ts},
		targets,
		store,
		builtins,
	)

	err := uc.Execute(bytes.NewReader([]byte("palette-yaml")), "test-theme", "")

	// Should return an error mentioning the failed target
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "css") {
		t.Errorf("error should mention 'css', got: %v", err)
	}

	// Other targets should still have been generated
	if !store.hasFile("test-theme", "palette.yaml") {
		t.Error("expected palette.yaml to be written")
	}
	if !store.hasFile("test-theme", "universal.yaml") {
		t.Error("expected universal.yaml to be written")
	}
	// vim, gtk, qss, stylix should still work
	for _, f := range []string{"vim-mapping.yaml", "style.lua", "gtk-mapping.yaml", "gtk.css", "qss-mapping.yaml", "style.qss", "stylix-mapping.yaml", "style.json"} {
		if !store.hasFile("test-theme", f) {
			t.Errorf("expected file %q to still be written despite css mapper error", f)
		}
	}

	// css files should NOT be written since mapper failed
	if store.hasFile("test-theme", "css-mapping.yaml") {
		t.Error("css-mapping.yaml should not be written when mapper fails")
	}
	if store.hasFile("test-theme", "style.css") {
		t.Error("style.css should not be written when mapper fails")
	}
}

func TestGenerateTheme_Integration(t *testing.T) {
	yamlBytes, err := os.ReadFile("../../testdata/tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("read test fixture: %v", err)
	}

	parser := yamlparser.NewParser()
	deriv := deriver.New()
	store := newStubThemeStore()
	builtins := newStubPaletteSource()

	targets := []ports.Target{
		{
			Mapper:          mapper.NewVim(),
			Generator:       generator.NewVim(),
			MappingFile:     "vim-mapping.yaml",
			MappingFileKind: domain.FileKindVimMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				theme, ok := mapped.(*ports.VimTheme)
				if !ok {
					return fmt.Errorf("expected *ports.VimTheme, got %T", mapped)
				}
				mf := vimThemeToMappingFile(theme)
				return fileio.WriteVimMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewCSS(),
			Generator:       generator.NewCSS(),
			MappingFile:     "css-mapping.yaml",
			MappingFileKind: domain.FileKindCSSMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				theme, ok := mapped.(*ports.CSSTheme)
				if !ok {
					return fmt.Errorf("expected *ports.CSSTheme, got %T", mapped)
				}
				mf := cssThemeToMappingFile(theme)
				return fileio.WriteCSSMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewGtk(),
			Generator:       generator.NewGtk(),
			MappingFile:     "gtk-mapping.yaml",
			MappingFileKind: domain.FileKindGtkMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				theme, ok := mapped.(*ports.GtkTheme)
				if !ok {
					return fmt.Errorf("expected *ports.GtkTheme, got %T", mapped)
				}
				mf := gtkThemeToMappingFile(theme)
				return fileio.WriteGtkMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewQss(),
			Generator:       generator.NewQss(),
			MappingFile:     "qss-mapping.yaml",
			MappingFileKind: domain.FileKindQssMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				theme, ok := mapped.(*ports.QssTheme)
				if !ok {
					return fmt.Errorf("expected *ports.QssTheme, got %T", mapped)
				}
				mf := qssThemeToMappingFile(theme)
				return fileio.WriteQssMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewStylix(),
			Generator:       generator.NewStylix(),
			MappingFile:     "stylix-mapping.yaml",
			MappingFileKind: domain.FileKindStylixMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				theme, ok := mapped.(*ports.StylixTheme)
				if !ok {
					return fmt.Errorf("expected *ports.StylixTheme, got %T", mapped)
				}
				mf := ports.StylixMappingFile{Values: theme.Values}
				return fileio.WriteStylixMapping(w, mf)
			},
		},
	}

	uc := application.NewGenerateThemeUseCase(parser, deriv, targets, store, builtins)

	err = uc.Execute(bytes.NewReader(yamlBytes), "tokyo-night-dark", "")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	// All 12 files should be produced
	expectedFiles := []string{
		"palette.yaml", "universal.yaml",
		"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml",
		"qss-mapping.yaml", "stylix-mapping.yaml",
		"style.lua", "style.css", "gtk.css", "style.qss", "style.json",
	}

	written := store.writtenFiles("tokyo-night-dark")
	if len(written) != 12 {
		t.Errorf("expected 12 files written, got %d: %v", len(written), written)
	}

	for _, f := range expectedFiles {
		if !store.hasFile("tokyo-night-dark", f) {
			t.Errorf("expected file %q to be written", f)
		}
	}

	// Verify mapping files contain valid YAML (non-empty)
	for _, mf := range []string{"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml", "qss-mapping.yaml", "stylix-mapping.yaml"} {
		if !store.hasFile("tokyo-night-dark", mf) {
			continue
		}
		rw := store.files["tokyo-night-dark"][mf]
		if len(rw.data) == 0 {
			t.Errorf("mapping file %q is empty", mf)
		}
	}

	// Verify output files are non-empty
	for _, of := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
		if !store.hasFile("tokyo-night-dark", of) {
			continue
		}
		rw := store.files["tokyo-night-dark"][of]
		if len(rw.data) == 0 {
			t.Errorf("output file %q is empty", of)
		}
	}
}

// --- Integration test helper functions for converting theme DTOs to mapping files ---

func vimThemeToMappingFile(theme *ports.VimTheme) ports.VimMappingFile {
	highlights := make(map[string]ports.VimMappingHighlight, len(theme.Highlights))
	for name, hl := range theme.Highlights {
		mhl := ports.VimMappingHighlight{
			Bold:          hl.Bold,
			Italic:        hl.Italic,
			Underline:     hl.Underline,
			Undercurl:     hl.Undercurl,
			Strikethrough: hl.Strikethrough,
			Reverse:       hl.Reverse,
			Nocombine:     hl.Nocombine,
			Link:          hl.Link,
		}
		if hl.Fg != nil {
			mhl.Fg = hl.Fg.Hex()
		}
		if hl.Bg != nil {
			mhl.Bg = hl.Bg.Hex()
		}
		if hl.Sp != nil {
			mhl.Sp = hl.Sp.Hex()
		}
		highlights[name] = mhl
	}

	var termColors [16]string
	for i, c := range theme.TerminalColors {
		termColors[i] = c.Hex()
	}

	return ports.VimMappingFile{
		Highlights:     highlights,
		TerminalColors: termColors,
	}
}

func cssThemeToMappingFile(theme *ports.CSSTheme) ports.CSSMappingFile {
	rules := make([]ports.CSSRuleEntry, len(theme.Rules))
	for i, r := range theme.Rules {
		props := make(map[string]string, len(r.Properties))
		for _, p := range r.Properties {
			props[p.Property] = p.Value
		}
		rules[i] = ports.CSSRuleEntry{
			Selector:   r.Selector,
			Properties: props,
		}
	}
	return ports.CSSMappingFile{
		CustomProperties: theme.CustomProperties,
		Rules:            rules,
	}
}

func gtkThemeToMappingFile(theme *ports.GtkTheme) ports.GtkMappingFile {
	colors := make(map[string]string, len(theme.Colors))
	for _, c := range theme.Colors {
		colors[c.Name] = c.Value
	}
	rules := make([]ports.CSSRuleEntry, len(theme.Rules))
	for i, r := range theme.Rules {
		props := make(map[string]string, len(r.Properties))
		for _, p := range r.Properties {
			props[p.Property] = p.Value
		}
		rules[i] = ports.CSSRuleEntry{
			Selector:   r.Selector,
			Properties: props,
		}
	}
	return ports.GtkMappingFile{
		Colors: colors,
		Rules:  rules,
	}
}

func qssThemeToMappingFile(theme *ports.QssTheme) ports.QssMappingFile {
	rules := make([]ports.CSSRuleEntry, len(theme.Rules))
	for i, r := range theme.Rules {
		props := make(map[string]string, len(r.Properties))
		for _, p := range r.Properties {
			props[p.Property] = p.Value
		}
		rules[i] = ports.CSSRuleEntry{
			Selector:   r.Selector,
			Properties: props,
		}
	}
	return ports.QssMappingFile{
		Rules: rules,
	}
}

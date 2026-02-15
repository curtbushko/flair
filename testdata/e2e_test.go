package testdata_test

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/curtbushko/flair/internal/adapters/deriver"
	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/adapters/generator"
	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/adapters/palettes"
	"github.com/curtbushko/flair/internal/adapters/store"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/application"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

var update = flag.Bool("update", false, "update golden files")

// goldenFiles maps each output filename to its golden file path under expected/.
var goldenFiles = map[string]string{
	"universal.yaml":      "expected/universal.yaml",
	"vim-mapping.yaml":    "expected/vim-mapping.yaml",
	"css-mapping.yaml":    "expected/css-mapping.yaml",
	"gtk-mapping.yaml":    "expected/gtk-mapping.yaml",
	"qss-mapping.yaml":    "expected/qss-mapping.yaml",
	"stylix-mapping.yaml": "expected/stylix-mapping.yaml",
	"style.lua":           "expected/style.lua",
	"style.css":           "expected/style.css",
	"gtk.css":             "expected/gtk.css",
	"style.qss":           "expected/style.qss",
	"style.json":          "expected/style.json",
}

// runPipeline runs the full flair pipeline for Tokyo Night Dark in the given
// temp directory and returns the path to the generated theme directory.
func runPipeline(t *testing.T, tmpDir string) string {
	t.Helper()

	paletteData, err := os.ReadFile("tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("read palette: %v", err)
	}

	parser := yamlparser.NewParser()
	drv := deriver.New()
	fsStore := store.NewFsStore(tmpDir)

	targets := buildTargets()

	uc := application.NewGenerateThemeUseCase(
		parser, drv, targets, fsStore, nil,
		application.WithPaletteWriter(writePaletteYAML),
		application.WithUniversalWriter(func(w io.Writer, ts *domain.TokenSet) error {
			return fileio.WriteUniversal(w, ts)
		}),
	)

	if err := uc.Execute(bytes.NewReader(paletteData), "tokyo-night-dark", ""); err != nil {
		t.Fatalf("pipeline execute: %v", err)
	}

	return filepath.Join(tmpDir, "tokyo-night-dark")
}

// buildTargets returns the full set of real targets (same wiring as production).
func buildTargets() []ports.Target {
	return []ports.Target{
		{
			Mapper:          mapper.NewVim(),
			Generator:       generator.NewVim(),
			MappingFile:     "vim-mapping.yaml",
			MappingFileKind: domain.FileKindVimMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				vt, ok := mapped.(*ports.VimTheme)
				if !ok {
					return fmt.Errorf("expected *ports.VimTheme, got %T", mapped)
				}
				mf := vimMappingFromTheme(vt)
				return fileio.WriteVimMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewCSS(),
			Generator:       generator.NewCSS(),
			MappingFile:     "css-mapping.yaml",
			MappingFileKind: domain.FileKindCSSMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				ct, ok := mapped.(*ports.CSSTheme)
				if !ok {
					return fmt.Errorf("expected *ports.CSSTheme, got %T", mapped)
				}
				mf := cssMappingFromTheme(ct)
				return fileio.WriteCSSMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewGtk(),
			Generator:       generator.NewGtk(),
			MappingFile:     "gtk-mapping.yaml",
			MappingFileKind: domain.FileKindGtkMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				gt, ok := mapped.(*ports.GtkTheme)
				if !ok {
					return fmt.Errorf("expected *ports.GtkTheme, got %T", mapped)
				}
				mf := gtkMappingFromTheme(gt)
				return fileio.WriteGtkMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewQss(),
			Generator:       generator.NewQss(),
			MappingFile:     "qss-mapping.yaml",
			MappingFileKind: domain.FileKindQssMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				qt, ok := mapped.(*ports.QssTheme)
				if !ok {
					return fmt.Errorf("expected *ports.QssTheme, got %T", mapped)
				}
				mf := qssMappingFromTheme(qt)
				return fileio.WriteQssMapping(w, mf)
			},
		},
		{
			Mapper:          mapper.NewStylix(),
			Generator:       generator.NewStylix(),
			MappingFile:     "stylix-mapping.yaml",
			MappingFileKind: domain.FileKindStylixMapping,
			WriteMappingFile: func(w io.Writer, mapped ports.MappedTheme) error {
				st, ok := mapped.(*ports.StylixTheme)
				if !ok {
					return fmt.Errorf("expected *ports.StylixTheme, got %T", mapped)
				}
				mf := ports.StylixMappingFile{Values: st.Values}
				return fileio.WriteStylixMapping(w, mf)
			},
		},
	}
}

// TestE2E_TokyoNightDark_FullPipeline runs the full pipeline for Tokyo Night
// Dark and compares all generated files against golden files in testdata/expected/.
// Use -update to regenerate golden files.
func TestE2E_TokyoNightDark_FullPipeline(t *testing.T) {
	tmpDir := t.TempDir()
	themeDir := runPipeline(t, tmpDir)

	if *update {
		if err := os.MkdirAll("expected", 0o755); err != nil {
			t.Fatalf("create expected dir: %v", err)
		}
	}

	for outputFile, goldenPath := range goldenFiles {
		t.Run(outputFile, func(t *testing.T) {
			gotData, err := os.ReadFile(filepath.Join(themeDir, outputFile))
			if err != nil {
				t.Fatalf("read generated file %q: %v", outputFile, err)
			}

			if *update {
				if err := os.WriteFile(goldenPath, gotData, 0o644); err != nil {
					t.Fatalf("update golden file %q: %v", goldenPath, err)
				}
				t.Logf("updated golden file: %s", goldenPath)
				return
			}

			wantData, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden file %q (run with -update to generate): %v", goldenPath, err)
			}

			if !bytes.Equal(gotData, wantData) {
				t.Errorf("output %q does not match golden file %q\n--- got (%d bytes) ---\n%s\n--- want (%d bytes) ---\n%s",
					outputFile, goldenPath, len(gotData), truncate(gotData, 500), len(wantData), truncate(wantData, 500))
			}
		})
	}
}

// TestE2E_TokyoNightDark_Deterministic runs the pipeline twice and verifies
// all outputs are byte-identical.
func TestE2E_TokyoNightDark_Deterministic(t *testing.T) {
	tmpDir1 := t.TempDir()
	themeDir1 := runPipeline(t, tmpDir1)

	tmpDir2 := t.TempDir()
	themeDir2 := runPipeline(t, tmpDir2)

	for outputFile := range goldenFiles {
		t.Run(outputFile, func(t *testing.T) {
			data1, err := os.ReadFile(filepath.Join(themeDir1, outputFile))
			if err != nil {
				t.Fatalf("read run1 %q: %v", outputFile, err)
			}

			data2, err := os.ReadFile(filepath.Join(themeDir2, outputFile))
			if err != nil {
				t.Fatalf("read run2 %q: %v", outputFile, err)
			}

			if !bytes.Equal(data1, data2) {
				t.Errorf("non-deterministic output for %q\n--- run1 (%d bytes) ---\n%s\n--- run2 (%d bytes) ---\n%s",
					outputFile, len(data1), truncate(data1, 500), len(data2), truncate(data2, 500))
			}
		})
	}
}

// truncate returns the first n bytes as a string, with "..." appended if truncated.
func truncate(data []byte, n int) string {
	if len(data) <= n {
		return string(data)
	}
	return string(data[:n]) + "..."
}

// --- Mapping file conversion helpers (same as production wire.go) ---

func vimMappingFromTheme(theme *ports.VimTheme) ports.VimMappingFile {
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
		if hl.Fg != nil && !hl.Fg.IsNone {
			mhl.Fg = hl.Fg.Hex()
		}
		if hl.Bg != nil && !hl.Bg.IsNone {
			mhl.Bg = hl.Bg.Hex()
		}
		if hl.Sp != nil && !hl.Sp.IsNone {
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

func cssMappingFromTheme(theme *ports.CSSTheme) ports.CSSMappingFile {
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

func gtkMappingFromTheme(theme *ports.GtkTheme) ports.GtkMappingFile {
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

// runPipelineBuiltin runs the full flair pipeline using the built-in palette
// name (via ExecuteBuiltin) and returns the path to the generated theme directory.
func runPipelineBuiltin(t *testing.T, tmpDir string) string {
	t.Helper()

	parser := yamlparser.NewParser()
	drv := deriver.New()
	fsStore := store.NewFsStore(tmpDir)
	builtinSrc := palettes.NewSource()

	targets := buildTargets()

	uc := application.NewGenerateThemeUseCase(
		parser, drv, targets, fsStore, builtinSrc,
		application.WithPaletteWriter(writePaletteYAML),
		application.WithUniversalWriter(func(w io.Writer, ts *domain.TokenSet) error {
			return fileio.WriteUniversal(w, ts)
		}),
	)

	if err := uc.ExecuteBuiltin("tokyo-night-dark", "tokyo-night-dark", ""); err != nil {
		t.Fatalf("pipeline execute builtin: %v", err)
	}

	return filepath.Join(tmpDir, "tokyo-night-dark")
}

// allGeneratedFiles lists every file the pipeline creates (palette.yaml +
// universal.yaml + 5 mapping files + 5 output files = 12 total).
var allGeneratedFiles = []string{
	"palette.yaml",
	"universal.yaml",
	"vim-mapping.yaml",
	"css-mapping.yaml",
	"gtk-mapping.yaml",
	"qss-mapping.yaml",
	"stylix-mapping.yaml",
	"style.lua",
	"style.css",
	"gtk.css",
	"style.qss",
	"style.json",
}

// TestE2E_BuiltinVsFile_IdenticalOutput generates a theme using the built-in
// palette name 'tokyo-night-dark' and separately from the testdata YAML file,
// then compares all 12 generated files are byte-identical.
func TestE2E_BuiltinVsFile_IdenticalOutput(t *testing.T) {
	// Generate from file
	tmpDirFile := t.TempDir()
	themeDirFile := runPipeline(t, tmpDirFile)

	// Generate from built-in name
	tmpDirBuiltin := t.TempDir()
	themeDirBuiltin := runPipelineBuiltin(t, tmpDirBuiltin)

	for _, outputFile := range allGeneratedFiles {
		t.Run(outputFile, func(t *testing.T) {
			fileData, err := os.ReadFile(filepath.Join(themeDirFile, outputFile))
			if err != nil {
				t.Fatalf("read file-generated %q: %v", outputFile, err)
			}

			builtinData, err := os.ReadFile(filepath.Join(themeDirBuiltin, outputFile))
			if err != nil {
				t.Fatalf("read builtin-generated %q: %v", outputFile, err)
			}

			if !bytes.Equal(fileData, builtinData) {
				t.Errorf("output %q differs between file and builtin generation\n--- file (%d bytes) ---\n%s\n--- builtin (%d bytes) ---\n%s",
					outputFile, len(fileData), truncate(fileData, 500), len(builtinData), truncate(builtinData, 500))
			}
		})
	}
}

// runPipelineForScheme runs the full flair pipeline for the given built-in
// scheme name and returns the path to the generated theme directory.
func runPipelineForScheme(t *testing.T, tmpDir, schemeName string) string {
	t.Helper()

	parser := yamlparser.NewParser()
	drv := deriver.New()
	fsStore := store.NewFsStore(tmpDir)
	builtinSrc := palettes.NewSource()

	targets := buildTargets()

	uc := application.NewGenerateThemeUseCase(
		parser, drv, targets, fsStore, builtinSrc,
		application.WithPaletteWriter(writePaletteYAML),
		application.WithUniversalWriter(func(w io.Writer, ts *domain.TokenSet) error {
			return fileio.WriteUniversal(w, ts)
		}),
	)

	if err := uc.ExecuteBuiltin(schemeName, schemeName, ""); err != nil {
		t.Fatalf("pipeline execute builtin %q: %v", schemeName, err)
	}

	return filepath.Join(tmpDir, schemeName)
}

// TestE2E_AdditionalSchemes runs the full pipeline for gruvbox-dark and
// catppuccin-mocha built-in palettes, verifying all 12 files are generated,
// non-empty, and contain format-specific markers.
func TestE2E_AdditionalSchemes(t *testing.T) {
	schemes := []struct {
		name string
	}{
		{name: "gruvbox-dark"},
		{name: "catppuccin-mocha"},
	}

	for _, sc := range schemes {
		t.Run(sc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			themeDir := runPipelineForScheme(t, tmpDir, sc.name)

			// Verify all 12 files exist and are non-empty
			for _, f := range allGeneratedFiles {
				t.Run(f, func(t *testing.T) {
					data, err := os.ReadFile(filepath.Join(themeDir, f))
					if err != nil {
						t.Fatalf("read generated file %q: %v", f, err)
					}
					if len(data) == 0 {
						t.Errorf("generated file %q is empty", f)
					}
				})
			}

			// Verify style.lua starts with expected Lua preamble
			luaData, err := os.ReadFile(filepath.Join(themeDir, "style.lua"))
			if err != nil {
				t.Fatalf("read style.lua: %v", err)
			}
			if !bytes.Contains(luaData, []byte("vim.cmd")) {
				t.Errorf("style.lua does not contain expected Lua preamble ('vim.cmd')")
			}

			// Verify style.json contains valid JSON (starts with '{')
			jsonData, err := os.ReadFile(filepath.Join(themeDir, "style.json"))
			if err != nil {
				t.Fatalf("read style.json: %v", err)
			}
			trimmedJSON := bytes.TrimSpace(jsonData)
			if len(trimmedJSON) == 0 || trimmedJSON[0] != '{' {
				t.Errorf("style.json does not start with '{', got %q", truncate(trimmedJSON, 20))
			}

			// Verify style.css contains ':root' selector
			cssData, err := os.ReadFile(filepath.Join(themeDir, "style.css"))
			if err != nil {
				t.Fatalf("read style.css: %v", err)
			}
			if !bytes.Contains(cssData, []byte(":root")) {
				t.Errorf("style.css does not contain ':root' selector")
			}

			// Verify gtk.css contains '@define-color' declarations
			gtkData, err := os.ReadFile(filepath.Join(themeDir, "gtk.css"))
			if err != nil {
				t.Fatalf("read gtk.css: %v", err)
			}
			if !bytes.Contains(gtkData, []byte("@define-color")) {
				t.Errorf("gtk.css does not contain '@define-color' declarations")
			}
		})
	}
}

// TestE2E_AdditionalSchemes_Deterministic runs each additional scheme twice
// and verifies all outputs are byte-identical.
func TestE2E_AdditionalSchemes_Deterministic(t *testing.T) {
	schemes := []string{"gruvbox-dark", "catppuccin-mocha"}

	for _, sc := range schemes {
		t.Run(sc, func(t *testing.T) {
			tmpDir1 := t.TempDir()
			themeDir1 := runPipelineForScheme(t, tmpDir1, sc)

			tmpDir2 := t.TempDir()
			themeDir2 := runPipelineForScheme(t, tmpDir2, sc)

			for _, f := range allGeneratedFiles {
				t.Run(f, func(t *testing.T) {
					data1, err := os.ReadFile(filepath.Join(themeDir1, f))
					if err != nil {
						t.Fatalf("read run1 %q: %v", f, err)
					}

					data2, err := os.ReadFile(filepath.Join(themeDir2, f))
					if err != nil {
						t.Fatalf("read run2 %q: %v", f, err)
					}

					if !bytes.Equal(data1, data2) {
						t.Errorf("non-deterministic output for %q\n--- run1 (%d bytes) ---\n%s\n--- run2 (%d bytes) ---\n%s",
							f, len(data1), truncate(data1, 500), len(data2), truncate(data2, 500))
					}
				})
			}
		})
	}
}

func qssMappingFromTheme(theme *ports.QssTheme) ports.QssMappingFile {
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

// writePaletteYAML serializes a domain.Palette back to the tinted-theming
// common YAML format so regenerate can re-parse it.
func writePaletteYAML(w io.Writer, pal *domain.Palette) error {
	if _, err := fmt.Fprintf(w, "system: %q\n", pal.System); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "name: %q\n", pal.Name); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "author: %q\n", pal.Author); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "variant: %q\n", pal.Variant); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "palette:\n"); err != nil {
		return err
	}
	slotNames := []string{
		"base00", "base01", "base02", "base03",
		"base04", "base05", "base06", "base07",
		"base08", "base09", "base0A", "base0B",
		"base0C", "base0D", "base0E", "base0F",
		"base10", "base11", "base12", "base13",
		"base14", "base15", "base16", "base17",
	}
	for i, name := range slotNames {
		c := pal.Base(i)
		if _, err := fmt.Fprintf(w, "  %s: %q\n", name, c.Hex()); err != nil {
			return err
		}
	}
	return nil
}

// --- Regenerate E2E helpers ---

// checksumFile returns the SHA-256 hex digest of a file.
func checksumFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("checksumFile %s: %v", path, err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// checksumAllThemeFiles returns a map of filename -> SHA-256 for all theme files.
func checksumAllThemeFiles(t *testing.T, themeDir string) map[string]string {
	t.Helper()
	sums := make(map[string]string, len(allGeneratedFiles))
	for _, f := range allGeneratedFiles {
		sums[f] = checksumFile(t, filepath.Join(themeDir, f))
	}
	return sums
}

// touchFileWithComment appends a YAML comment to a file, changing its content
// and mtime. This simulates a user editing the file.
func touchFileWithComment(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("touchFileWithComment read %s: %v", path, err)
	}
	data = append(data, []byte("\n# user-edit\n")...)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("touchFileWithComment write %s: %v", path, err)
	}
}

// sleepForMtime sleeps briefly to ensure filesystem mtime granularity is exceeded.
func sleepForMtime() {
	time.Sleep(50 * time.Millisecond)
}

// newRegenUseCase builds a RegenerateThemeUseCase with real adapters.
func newRegenUseCase(fsStore *store.FsStore) *application.RegenerateThemeUseCase {
	parser := yamlparser.NewParser()
	drv := deriver.New()
	targets := buildTargets()

	return application.NewRegenerateThemeUseCase(
		fsStore, parser, drv, targets,
		application.WithRegenUniversalWriter(func(w io.Writer, ts *domain.TokenSet) error {
			return fileio.WriteUniversal(w, ts)
		}),
	)
}

// --- Regenerate E2E Tests ---

// fileMtime returns a file's modification time.
func fileMtime(t *testing.T, path string) time.Time {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return info.ModTime()
}

// recordMtimes returns a map of filename -> mtime for all theme files.
func recordMtimes(t *testing.T, themeDir string) map[string]time.Time {
	t.Helper()
	mtimes := make(map[string]time.Time, len(allGeneratedFiles))
	for _, f := range allGeneratedFiles {
		mtimes[f] = fileMtime(t, filepath.Join(themeDir, f))
	}
	return mtimes
}

// TestE2E_Regenerate_UniversalEdit_DownstreamOnly generates a full theme,
// records mtimes of all files, edits universal.yaml, runs regenerate, and
// verifies: palette.yaml is unchanged, universal.yaml content is preserved
// (not overwritten), and all 5 mapping files + 5 output files are regenerated
// (their mtimes are updated and content is rewritten).
func TestE2E_Regenerate_UniversalEdit_DownstreamOnly(t *testing.T) {
	tmpDir := t.TempDir()
	themeDir := runPipeline(t, tmpDir)
	themeName := "tokyo-night-dark"
	fsStore := store.NewFsStore(tmpDir)

	// Record mtimes and checksum of palette.yaml after initial generation.
	beforeMtimes := recordMtimes(t, themeDir)
	paletteChecksum := checksumFile(t, filepath.Join(themeDir, "palette.yaml"))

	// Sleep to ensure mtime granularity is exceeded.
	sleepForMtime()

	// Edit universal.yaml (simulating user edit).
	universalPath := filepath.Join(themeDir, "universal.yaml")
	touchFileWithComment(t, universalPath)

	// Record the universal.yaml content after edit to verify it is preserved.
	editedUniversalChecksum := checksumFile(t, universalPath)

	// Run regenerate.
	regenUC := newRegenUseCase(fsStore)
	msg, err := regenUC.Execute(themeName, "")
	if err != nil {
		t.Fatalf("regenerate: %v", err)
	}

	if msg == "" || strings.Contains(msg, "nothing to do") {
		t.Fatalf("expected regeneration to occur, got: %q", msg)
	}

	// Record mtimes after regeneration.
	afterMtimes := recordMtimes(t, themeDir)

	// 1. palette.yaml checksum must be unchanged.
	afterPaletteChecksum := checksumFile(t, filepath.Join(themeDir, "palette.yaml"))
	if afterPaletteChecksum != paletteChecksum {
		t.Error("palette.yaml content was modified during regeneration; expected it to remain unchanged")
	}

	// palette.yaml mtime must not have changed (file was not touched).
	if !afterMtimes["palette.yaml"].Equal(beforeMtimes["palette.yaml"]) {
		t.Error("palette.yaml mtime changed during regeneration; expected it to remain untouched")
	}

	// 2. universal.yaml must NOT be overwritten (user edit preserved).
	afterUniversalChecksum := checksumFile(t, filepath.Join(themeDir, "universal.yaml"))
	if afterUniversalChecksum != editedUniversalChecksum {
		t.Error("universal.yaml was overwritten during regeneration; expected user edit to be preserved")
	}

	// 3. All 5 mapping files must have been rewritten (mtime updated).
	mappings := []string{
		"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml",
		"qss-mapping.yaml", "stylix-mapping.yaml",
	}
	for _, mf := range mappings {
		if !afterMtimes[mf].After(beforeMtimes[mf]) {
			t.Errorf("mapping file %q mtime was NOT updated; expected it to be regenerated", mf)
		}
	}

	// 4. All 5 output files must have been rewritten (mtime updated).
	outputs := []string{
		"style.lua", "style.css", "gtk.css", "style.qss", "style.json",
	}
	for _, of := range outputs {
		if !afterMtimes[of].After(beforeMtimes[of]) {
			t.Errorf("output file %q mtime was NOT updated; expected it to be regenerated", of)
		}
	}
}

// TestE2E_AllBuiltins_ParseClean verifies that every built-in palette returned
// by PaletteSource.List() parses successfully via PaletteParser.Parse().
func TestE2E_AllBuiltins_ParseClean(t *testing.T) {
	src := palettes.NewSource()
	parser := yamlparser.NewParser()

	names := src.List()
	if len(names) == 0 {
		t.Fatal("PaletteSource.List() returned zero palettes; expected at least one built-in")
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			r, err := src.Get(name)
			if err != nil {
				t.Fatalf("PaletteSource.Get(%q): %v", name, err)
			}

			pal, err := parser.Parse(r)
			if err != nil {
				t.Fatalf("PaletteParser.Parse(%q): %v", name, err)
			}

			if pal.Name == "" {
				t.Errorf("parsed palette %q has empty Name", name)
			}
		})
	}
}

// TestE2E_AllBuiltins_ValidateClean verifies that every built-in palette
// passes domain.ValidatePalette() with zero hard violations (completeness
// and luminance ordering errors). Soft warnings (monotonicity and bright
// variant) are logged but do not fail the test, since upstream palette
// authors may intentionally deviate from those heuristics.
func TestE2E_AllBuiltins_ValidateClean(t *testing.T) {
	src := palettes.NewSource()
	parser := yamlparser.NewParser()

	names := src.List()
	if len(names) == 0 {
		t.Fatal("PaletteSource.List() returned zero palettes; expected at least one built-in")
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			r, err := src.Get(name)
			if err != nil {
				t.Fatalf("PaletteSource.Get(%q): %v", name, err)
			}

			pal, err := parser.Parse(r)
			if err != nil {
				t.Fatalf("PaletteParser.Parse(%q): %v", name, err)
			}

			violations := domain.ValidatePalette(pal)
			for _, v := range violations {
				if strings.Contains(v, "warning") {
					t.Logf("warning (non-fatal): %s", v)
				} else {
					t.Errorf("violation: %s", v)
				}
			}
		})
	}
}

// TestE2E_Regenerate_MappingEdit_SingleTarget generates a full theme, edits
// vim-mapping.yaml, runs regenerate, and verifies only style.lua is regenerated
// while all other files remain unchanged.
func TestE2E_Regenerate_MappingEdit_SingleTarget(t *testing.T) {
	tmpDir := t.TempDir()
	themeDir := runPipeline(t, tmpDir)
	themeName := "tokyo-night-dark"
	fsStore := store.NewFsStore(tmpDir)

	// Record mtimes after initial generation.
	beforeMtimes := recordMtimes(t, themeDir)
	beforeChecksums := checksumAllThemeFiles(t, themeDir)

	// Sleep to ensure mtime granularity is exceeded.
	sleepForMtime()

	// Edit vim-mapping.yaml (simulating user edit to a single mapping file).
	vimMappingPath := filepath.Join(themeDir, "vim-mapping.yaml")
	touchFileWithComment(t, vimMappingPath)

	// Run regenerate.
	regenUC := newRegenUseCase(fsStore)
	msg, err := regenUC.Execute(themeName, "")
	if err != nil {
		t.Fatalf("regenerate: %v", err)
	}

	if msg == "" || strings.Contains(msg, "nothing to do") {
		t.Fatalf("expected regeneration to occur, got: %q", msg)
	}

	// Record after state.
	afterMtimes := recordMtimes(t, themeDir)
	afterChecksums := checksumAllThemeFiles(t, themeDir)

	// 1. palette.yaml must be unchanged (both content and mtime).
	if afterChecksums["palette.yaml"] != beforeChecksums["palette.yaml"] {
		t.Error("palette.yaml content was modified; expected unchanged")
	}
	if !afterMtimes["palette.yaml"].Equal(beforeMtimes["palette.yaml"]) {
		t.Error("palette.yaml mtime changed; expected unchanged")
	}

	// 2. universal.yaml must be unchanged (both content and mtime).
	if afterChecksums["universal.yaml"] != beforeChecksums["universal.yaml"] {
		t.Error("universal.yaml content was modified; expected unchanged")
	}
	if !afterMtimes["universal.yaml"].Equal(beforeMtimes["universal.yaml"]) {
		t.Error("universal.yaml mtime changed; expected unchanged")
	}

	// 3. Only style.lua (vim output) should be regenerated (mtime updated).
	if !afterMtimes["style.lua"].After(beforeMtimes["style.lua"]) {
		t.Error("style.lua mtime was NOT updated; expected it to be regenerated after vim-mapping.yaml edit")
	}

	// 4. All other mapping and output files must remain unchanged.
	unchangedFiles := []string{
		"css-mapping.yaml", "gtk-mapping.yaml", "qss-mapping.yaml", "stylix-mapping.yaml",
		"style.css", "gtk.css", "style.qss", "style.json",
	}
	for _, f := range unchangedFiles {
		if afterChecksums[f] != beforeChecksums[f] {
			t.Errorf("file %q content was modified; expected it to remain unchanged", f)
		}
		if !afterMtimes[f].Equal(beforeMtimes[f]) {
			t.Errorf("file %q mtime changed; expected it to remain unchanged after vim-mapping.yaml edit only", f)
		}
	}
}

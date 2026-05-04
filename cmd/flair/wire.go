package main

import (
	"fmt"
	"io"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/adapters/generator"
	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/adapters/palettes"
	"github.com/curtbushko/flair/internal/adapters/store"
	"github.com/curtbushko/flair/internal/adapters/tokenizer"
	"github.com/curtbushko/flair/internal/adapters/wrappers"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/application"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// App holds the wired use cases and dependencies for the flair CLI.
type App struct {
	Targets    []ports.Target
	Generate   *application.GenerateThemeUseCase
	Derive     *application.DeriveThemeUseCase
	Select     *application.SelectThemeUseCase
	List       *application.ListThemesUseCase
	Validate   *application.ValidateThemeUseCase
	Preview    *application.PreviewThemeUseCase
	Regenerate *application.RegenerateThemeUseCase
	Store      ports.ThemeStore
	Builtins   ports.PaletteSource
}

// Wire creates all adapters, targets, and use cases, returning a fully
// wired App. configDir is the root configuration directory.
func Wire(configDir string) *App {
	// Adapters
	parser := yamlparser.NewParser()
	tok := tokenizer.New()
	fsStore := store.NewFsStore(configDir)
	builtins := palettes.NewSource()

	// Targets: one per output format.
	targets := []ports.Target{
		{
			Mapper:           mapper.NewVim(),
			Generator:        generator.NewVim(),
			MappingFile:      "vim-mapping.yaml",
			MappingFileKind:  domain.FileKindVimMapping,
			WriteMappingFile: writeVimMapping,
		},
		{
			Mapper:           mapper.NewCSS(),
			Generator:        generator.NewCSS(),
			MappingFile:      "css-mapping.yaml",
			MappingFileKind:  domain.FileKindCSSMapping,
			WriteMappingFile: writeCSSMapping,
		},
		{
			Mapper:           mapper.NewGtk(),
			Generator:        generator.NewGtk(),
			MappingFile:      "gtk-mapping.yaml",
			MappingFileKind:  domain.FileKindGtkMapping,
			WriteMappingFile: writeGtkMapping,
		},
		{
			Mapper:           mapper.NewQss(),
			Generator:        generator.NewQss(),
			MappingFile:      "qss-mapping.yaml",
			MappingFileKind:  domain.FileKindQssMapping,
			WriteMappingFile: writeQssMapping,
		},
		{
			Mapper:           mapper.NewStylix(),
			Generator:        generator.NewStylix(),
			MappingFile:      "stylix-mapping.yaml",
			MappingFileKind:  domain.FileKindStylixMapping,
			WriteMappingFile: writeStylixMapping,
		},
	}

	// Use cases
	deriveUC := application.NewDeriveThemeUseCase(parser, tok)
	generateUC := application.NewGenerateThemeUseCase(
		parser, tok, targets, fsStore, builtins,
		application.WithPaletteWriter(func(w io.Writer, pal *domain.Palette) error {
			return fileio.WritePalette(w, pal)
		}),
		application.WithTokensWriter(func(w io.Writer, ts *domain.TokenSet) error {
			return fileio.WriteTokens(w, ts)
		}),
	)
	regenerateUC := application.NewRegenerateThemeUseCase(
		fsStore, parser, tok, targets,
		application.WithRegenTokensWriter(func(w io.Writer, ts *domain.TokenSet) error {
			return fileio.WriteTokens(w, ts)
		}),
	)
	selectUC := application.NewSelectThemeUseCase(fsStore, builtins, generateUC, regenerateUC)
	listUC := application.NewListThemesUseCase(fsStore, builtins)
	validateUC := application.NewValidateThemeUseCase(fsStore, parser, schemaValidatorFunc())
	previewUC := application.NewPreviewThemeUseCase(fsStore, parser, fileio.ReadTokens, tok, builtins)

	return &App{
		Targets:    targets,
		Generate:   generateUC,
		Derive:     deriveUC,
		Select:     selectUC,
		List:       listUC,
		Validate:   validateUC,
		Preview:    previewUC,
		Regenerate: regenerateUC,
		Store:      fsStore,
		Builtins:   builtins,
	}
}

// schemaValidatorFunc returns a SchemaValidator that uses the ValidatingReader
// adapter to check schema_version headers.
func schemaValidatorFunc() application.SchemaValidator {
	return func(r io.Reader, kind domain.FileKind) error {
		vr := wrappers.NewValidatingReader(r, kind)
		_, err := io.ReadAll(vr)
		return err
	}
}

// writeVimMapping adapts fileio.WriteVimMapping to ports.MappingFileWriter.
func writeVimMapping(w io.Writer, mapped ports.MappedTheme) error {
	vt, ok := mapped.(*ports.VimTheme)
	if !ok {
		return fmt.Errorf("write vim mapping: expected *ports.VimTheme, got %T", mapped)
	}

	mf := ports.VimMappingFile{
		Highlights:     make(map[string]ports.VimMappingHighlight, len(vt.Highlights)),
		TerminalColors: vimTerminalColors(vt.TerminalColors),
	}

	for name, hl := range vt.Highlights {
		mf.Highlights[name] = vimHighlightToMapping(hl)
	}

	return fileio.WriteVimMapping(w, mf)
}

// writeCSSMapping adapts fileio.WriteCSSMapping to ports.MappingFileWriter.
func writeCSSMapping(w io.Writer, mapped ports.MappedTheme) error {
	ct, ok := mapped.(*ports.CSSTheme)
	if !ok {
		return fmt.Errorf("write css mapping: expected *ports.CSSTheme, got %T", mapped)
	}

	mf := ports.CSSMappingFile{
		CustomProperties: ct.CustomProperties,
		Rules:            cssRulesToEntries(ct.Rules),
	}

	return fileio.WriteCSSMapping(w, mf)
}

// writeGtkMapping adapts fileio.WriteGtkMapping to ports.MappingFileWriter.
func writeGtkMapping(w io.Writer, mapped ports.MappedTheme) error {
	gt, ok := mapped.(*ports.GtkTheme)
	if !ok {
		return fmt.Errorf("write gtk mapping: expected *ports.GtkTheme, got %T", mapped)
	}

	colors := make(map[string]string, len(gt.Colors))
	for _, c := range gt.Colors {
		colors[c.Name] = c.Value
	}

	mf := ports.GtkMappingFile{
		Colors: colors,
		Rules:  cssRulesToEntries(gt.Rules),
	}

	return fileio.WriteGtkMapping(w, mf)
}

// writeQssMapping adapts fileio.WriteQssMapping to ports.MappingFileWriter.
func writeQssMapping(w io.Writer, mapped ports.MappedTheme) error {
	qt, ok := mapped.(*ports.QssTheme)
	if !ok {
		return fmt.Errorf("write qss mapping: expected *ports.QssTheme, got %T", mapped)
	}

	mf := ports.QssMappingFile{
		Rules: cssRulesToEntries(qt.Rules),
	}

	return fileio.WriteQssMapping(w, mf)
}

// writeStylixMapping adapts fileio.WriteStylixMapping to ports.MappingFileWriter.
func writeStylixMapping(w io.Writer, mapped ports.MappedTheme) error {
	st, ok := mapped.(*ports.StylixTheme)
	if !ok {
		return fmt.Errorf("write stylix mapping: expected *ports.StylixTheme, got %T", mapped)
	}

	mf := ports.StylixMappingFile{
		Values: st.Values,
	}

	return fileio.WriteStylixMapping(w, mf)
}

// vimHighlightToMapping converts a ports.VimHighlight to a VimMappingHighlight.
func vimHighlightToMapping(hl ports.VimHighlight) ports.VimMappingHighlight {
	m := ports.VimMappingHighlight{
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
		m.Fg = hl.Fg.Hex()
	}

	if hl.Bg != nil && !hl.Bg.IsNone {
		m.Bg = hl.Bg.Hex()
	}

	if hl.Sp != nil && !hl.Sp.IsNone {
		m.Sp = hl.Sp.Hex()
	}

	return m
}

// vimTerminalColors converts a [16]domain.Color to [16]string hex values.
func vimTerminalColors(colors [16]domain.Color) [16]string {
	var result [16]string
	for i, c := range colors {
		result[i] = c.Hex()
	}
	return result
}

// cssRulesToEntries converts []ports.CSSRule to []ports.CSSRuleEntry for YAML serialization.
func cssRulesToEntries(rules []ports.CSSRule) []ports.CSSRuleEntry {
	entries := make([]ports.CSSRuleEntry, 0, len(rules))
	for _, rule := range rules {
		props := make(map[string]string, len(rule.Properties))
		for _, p := range rule.Properties {
			props[p.Property] = p.Value
		}

		entries = append(entries, ports.CSSRuleEntry{
			Selector:   rule.Selector,
			Properties: props,
		})
	}
	return entries
}

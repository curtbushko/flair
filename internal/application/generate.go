package application

import (
	"fmt"
	"io"
	"strings"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// PaletteWriter writes a palette to an io.Writer. The composition root wires
// this to write raw palette YAML (optionally via VersionedWriter).
type PaletteWriter func(w io.Writer, palette *domain.Palette) error

// UniversalWriter writes a token set to an io.Writer. The composition root
// wires this to fileio.WriteUniversal (optionally via VersionedWriter).
type UniversalWriter func(w io.Writer, tokens *domain.TokenSet) error

// GenerateThemeUseCase orchestrates the full theme generation pipeline:
// parse palette -> derive tokens -> for each target: map -> write mapping -> generate output.
// It depends only on port interfaces, keeping the application layer adapter-agnostic.
type GenerateThemeUseCase struct {
	parser          ports.PaletteParser
	deriver         ports.TokenDeriver
	targets         []ports.Target
	store           ports.ThemeStore
	builtins        ports.PaletteSource
	paletteWriter   PaletteWriter
	universalWriter UniversalWriter
}

// NewGenerateThemeUseCase returns a new GenerateThemeUseCase wired to the given
// port implementations. paletteWriter and universalWriter default to simple
// pass-through writers if not provided via options.
func NewGenerateThemeUseCase(
	parser ports.PaletteParser,
	deriver ports.TokenDeriver,
	targets []ports.Target,
	store ports.ThemeStore,
	builtins ports.PaletteSource,
	opts ...GenerateOption,
) *GenerateThemeUseCase {
	uc := &GenerateThemeUseCase{
		parser:   parser,
		deriver:  deriver,
		targets:  targets,
		store:    store,
		builtins: builtins,
	}
	for _, opt := range opts {
		opt(uc)
	}
	return uc
}

// GenerateOption configures optional dependencies on GenerateThemeUseCase.
type GenerateOption func(*GenerateThemeUseCase)

// WithPaletteWriter sets a custom palette writer function.
func WithPaletteWriter(pw PaletteWriter) GenerateOption {
	return func(uc *GenerateThemeUseCase) {
		uc.paletteWriter = pw
	}
}

// WithUniversalWriter sets a custom universal writer function.
func WithUniversalWriter(uw UniversalWriter) GenerateOption {
	return func(uc *GenerateThemeUseCase) {
		uc.universalWriter = uw
	}
}

// Execute runs the full pipeline from an io.Reader palette source.
// If targetFilter is non-empty, only that target is generated (plus
// palette.yaml and universal.yaml).
func (uc *GenerateThemeUseCase) Execute(r io.Reader, themeName, targetFilter string) error {
	// 1. Parse palette
	palette, err := uc.parser.Parse(r)
	if err != nil {
		return fmt.Errorf("parse palette: %w", err)
	}

	return uc.generate(palette, themeName, targetFilter)
}

// ExecuteBuiltin resolves a built-in palette name, optionally infers the theme
// name, and runs the full pipeline.
func (uc *GenerateThemeUseCase) ExecuteBuiltin(builtinName, themeName, targetFilter string) error {
	r, err := uc.builtins.Get(builtinName)
	if err != nil {
		return fmt.Errorf("get built-in palette %q: %w", builtinName, err)
	}

	palette, err := uc.parser.Parse(r)
	if err != nil {
		return fmt.Errorf("parse palette: %w", err)
	}

	// Infer theme name from built-in name if not provided.
	if themeName == "" {
		themeName = builtinName
	}

	return uc.generate(palette, themeName, targetFilter)
}

// generate is the core pipeline: derive tokens, ensure dir, write files, map+generate targets.
func (uc *GenerateThemeUseCase) generate(palette *domain.Palette, themeName, targetFilter string) error {
	// 2. Derive tokens
	tokens := uc.deriver.Derive(palette)

	// 3. Ensure theme directory exists
	if err := uc.store.EnsureThemeDir(themeName); err != nil {
		return fmt.Errorf("ensure theme dir: %w", err)
	}

	// 4. Write palette.yaml
	if err := uc.writePalette(palette, themeName); err != nil {
		return fmt.Errorf("write palette: %w", err)
	}

	// 5. Write universal.yaml
	if err := uc.writeUniversal(tokens, themeName); err != nil {
		return fmt.Errorf("write universal: %w", err)
	}

	// 6. Build resolved theme for mappers
	resolved := &domain.ResolvedTheme{
		Name:    palette.Name,
		Variant: palette.Variant,
		Palette: palette,
		Tokens:  tokens,
	}

	// 7. Process each target
	var targetErrors []string

	for _, target := range uc.targets {
		// Apply target filter
		if targetFilter != "" && target.Mapper.Name() != targetFilter {
			continue
		}

		if err := uc.processTarget(target, resolved, themeName); err != nil {
			targetErrors = append(targetErrors, fmt.Sprintf("%s: %v", target.Mapper.Name(), err))
			continue
		}
	}

	if len(targetErrors) > 0 {
		return fmt.Errorf("target errors: %s", strings.Join(targetErrors, "; "))
	}

	return nil
}

// processTarget runs the map -> write mapping -> generate pipeline for a single target.
func (uc *GenerateThemeUseCase) processTarget(target ports.Target, resolved *domain.ResolvedTheme, themeName string) error {
	// Map
	mapped, err := target.Mapper.Map(resolved)
	if err != nil {
		return fmt.Errorf("map: %w", err)
	}

	// Write mapping file
	if err := uc.writeMappingFile(target, mapped, themeName); err != nil {
		return fmt.Errorf("write mapping: %w", err)
	}

	// Generate output file
	if err := uc.writeOutputFile(target, mapped, themeName); err != nil {
		return fmt.Errorf("generate output: %w", err)
	}

	return nil
}

// writePalette writes palette.yaml to the theme directory.
func (uc *GenerateThemeUseCase) writePalette(palette *domain.Palette, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, "palette.yaml")
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	if uc.paletteWriter != nil {
		return uc.paletteWriter(w, palette)
	}

	// Default: write a minimal palette marker.
	_, writeErr := fmt.Fprintf(w, "# palette for %s\n", themeName)
	return writeErr
}

// writeUniversal writes universal.yaml to the theme directory.
func (uc *GenerateThemeUseCase) writeUniversal(tokens *domain.TokenSet, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, "universal.yaml")
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	if uc.universalWriter != nil {
		return uc.universalWriter(w, tokens)
	}

	// Default: write a minimal universal marker.
	_, writeErr := fmt.Fprintf(w, "# universal for %s\n", themeName)
	return writeErr
}

// writeMappingFile writes the target's mapping YAML to the theme directory.
func (uc *GenerateThemeUseCase) writeMappingFile(target ports.Target, mapped ports.MappedTheme, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, target.MappingFile)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	if target.WriteMappingFile != nil {
		return target.WriteMappingFile(w, mapped)
	}

	// Fallback: write a minimal marker.
	_, writeErr := fmt.Fprintf(w, "# mapping for %s\n", target.Mapper.Name())
	return writeErr
}

// writeOutputFile generates and writes the final output file.
func (uc *GenerateThemeUseCase) writeOutputFile(target ports.Target, mapped ports.MappedTheme, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, target.Generator.DefaultFilename())
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	return target.Generator.Generate(w, mapped)
}

package application

import (
	"fmt"
	"strings"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// RegenerateThemeUseCase regenerates all theme files from the palette.
//
// Dependency chain:
//
//	palette.yaml -> tokens.yaml -> mapping files -> output files
//
// All files in the dependency chain are always regenerated to ensure
// consistency when mapper code or tokenizer logic changes.
type RegenerateThemeUseCase struct {
	store        ports.ThemeStore
	parser       ports.PaletteParser
	tokenizer    ports.Tokenizer
	targets      []ports.Target
	tokensWriter TokensWriter
}

// RegenerateOption configures optional dependencies on RegenerateThemeUseCase.
type RegenerateOption func(*RegenerateThemeUseCase)

// WithRegenTokensWriter sets a custom tokens writer for regeneration.
func WithRegenTokensWriter(tw TokensWriter) RegenerateOption {
	return func(uc *RegenerateThemeUseCase) {
		uc.tokensWriter = tw
	}
}

// NewRegenerateThemeUseCase returns a new RegenerateThemeUseCase.
func NewRegenerateThemeUseCase(
	store ports.ThemeStore,
	parser ports.PaletteParser,
	tokenizer ports.Tokenizer,
	targets []ports.Target,
	opts ...RegenerateOption,
) *RegenerateThemeUseCase {
	uc := &RegenerateThemeUseCase{
		store:     store,
		parser:    parser,
		tokenizer: tokenizer,
		targets:   targets,
	}
	for _, opt := range opts {
		opt(uc)
	}
	return uc
}

// Execute regenerates all theme files from the palette.
// If targetFilter is non-empty, only that target is considered.
// Returns a human-readable message and any error.
func (uc *RegenerateThemeUseCase) Execute(themeName, targetFilter string) (string, error) {
	// Verify the theme exists by checking palette.yaml.
	if !uc.store.FileExists(themeName, "palette.yaml") {
		return "", fmt.Errorf("theme %q not found (missing palette.yaml)", themeName)
	}

	// Read palette for re-derivation.
	palette, err := uc.readPalette(themeName)
	if err != nil {
		return "", fmt.Errorf("read palette: %w", err)
	}

	// Resolve the theme for downstream operations.
	tokens := uc.tokenizer.Tokenize(palette)
	resolved := &domain.ResolvedTheme{
		Name:    palette.Name,
		Variant: palette.Variant,
		Palette: palette,
		Tokens:  tokens,
	}

	// Always regenerate everything from palette.
	regenerated, err := uc.regenFromPalette(resolved, themeName, targetFilter)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("regenerated %d files: %s", len(regenerated), strings.Join(regenerated, ", ")), nil
}


// regenFromPalette re-derives tokens.yaml and all targets.
func (uc *RegenerateThemeUseCase) regenFromPalette(
	resolved *domain.ResolvedTheme, themeName, targetFilter string,
) ([]string, error) {
	if err := uc.writeTokens(resolved.Tokens, themeName); err != nil {
		return nil, fmt.Errorf("write tokens: %w", err)
	}

	regenerated := []string{"tokens.yaml"}

	regen, err := uc.regenerateTargets(resolved, themeName, targetFilter)
	if err != nil {
		return regenerated, err
	}

	return append(regenerated, regen...), nil
}


// regenerateTargets maps and generates all (filtered) targets.
func (uc *RegenerateThemeUseCase) regenerateTargets(resolved *domain.ResolvedTheme, themeName, targetFilter string) ([]string, error) {
	var regenerated []string
	var targetErrors []string

	for _, target := range uc.targets {
		if targetFilter != "" && target.Mapper.Name() != targetFilter {
			continue
		}

		if err := uc.processTarget(target, resolved, themeName); err != nil {
			targetErrors = append(targetErrors, fmt.Sprintf("%s: %v", target.Mapper.Name(), err))
			continue
		}
		regenerated = append(regenerated, target.MappingFile, target.Generator.DefaultFilename())
	}

	if len(targetErrors) > 0 {
		return regenerated, fmt.Errorf("target errors: %s", strings.Join(targetErrors, "; "))
	}

	return regenerated, nil
}


// processTarget runs the map -> write mapping -> generate pipeline for a single target.
func (uc *RegenerateThemeUseCase) processTarget(target ports.Target, resolved *domain.ResolvedTheme, themeName string) error {
	mapped, err := target.Mapper.Map(resolved)
	if err != nil {
		return fmt.Errorf("map: %w", err)
	}

	if err := uc.writeMappingFile(target, mapped, themeName); err != nil {
		return fmt.Errorf("write mapping: %w", err)
	}

	if err := uc.writeOutputFile(target, mapped, themeName); err != nil {
		return fmt.Errorf("generate output: %w", err)
	}

	return nil
}

// readPalette reads and parses palette.yaml from the theme directory.
func (uc *RegenerateThemeUseCase) readPalette(themeName string) (*domain.Palette, error) {
	rc, err := uc.store.OpenReader(themeName, "palette.yaml")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rc.Close() }()

	return uc.parser.Parse(rc)
}

// writeTokens writes tokens.yaml to the theme directory.
func (uc *RegenerateThemeUseCase) writeTokens(tokens *domain.TokenSet, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, "tokens.yaml")
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	if uc.tokensWriter != nil {
		return uc.tokensWriter(w, tokens)
	}

	_, writeErr := fmt.Fprintf(w, "# tokens for %s\n", themeName)
	return writeErr
}

// writeMappingFile writes the target's mapping YAML to the theme directory.
func (uc *RegenerateThemeUseCase) writeMappingFile(target ports.Target, mapped ports.MappedTheme, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, target.MappingFile)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	if target.WriteMappingFile != nil {
		return target.WriteMappingFile(w, mapped)
	}

	_, writeErr := fmt.Fprintf(w, "# mapping for %s\n", target.Mapper.Name())
	return writeErr
}

// writeOutputFile generates and writes the final output file.
func (uc *RegenerateThemeUseCase) writeOutputFile(target ports.Target, mapped ports.MappedTheme, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, target.Generator.DefaultFilename())
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	return target.Generator.Generate(w, mapped)
}

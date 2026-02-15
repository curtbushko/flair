package application

import (
	"fmt"
	"strings"
	"time"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// RegenerateThemeUseCase inspects modification times to determine which
// downstream files need re-derivation, then regenerates only the stale ones.
//
// Dependency chain:
//
//	palette.yaml -> universal.yaml -> mapping files -> output files
//
// If palette is newer than universal, everything downstream is regenerated.
// If universal is newer than mappings, all mappings + outputs are regenerated.
// If a mapping is newer than its output, only that output is regenerated.
type RegenerateThemeUseCase struct {
	store           ports.ThemeStore
	parser          ports.PaletteParser
	deriver         ports.TokenDeriver
	targets         []ports.Target
	universalWriter UniversalWriter
}

// RegenerateOption configures optional dependencies on RegenerateThemeUseCase.
type RegenerateOption func(*RegenerateThemeUseCase)

// WithRegenUniversalWriter sets a custom universal writer for regeneration.
func WithRegenUniversalWriter(uw UniversalWriter) RegenerateOption {
	return func(uc *RegenerateThemeUseCase) {
		uc.universalWriter = uw
	}
}

// NewRegenerateThemeUseCase returns a new RegenerateThemeUseCase.
func NewRegenerateThemeUseCase(
	store ports.ThemeStore,
	parser ports.PaletteParser,
	deriver ports.TokenDeriver,
	targets []ports.Target,
	opts ...RegenerateOption,
) *RegenerateThemeUseCase {
	uc := &RegenerateThemeUseCase{
		store:   store,
		parser:  parser,
		deriver: deriver,
		targets: targets,
	}
	for _, opt := range opts {
		opt(uc)
	}
	return uc
}

// Execute inspects mtimes in the theme directory and regenerates stale files.
// If targetFilter is non-empty, only that target is considered.
// Returns a human-readable message and any error.
func (uc *RegenerateThemeUseCase) Execute(themeName, targetFilter string) (string, error) {
	// Verify the theme exists by checking palette.yaml.
	if !uc.store.FileExists(themeName, "palette.yaml") {
		return "", fmt.Errorf("theme %q not found (missing palette.yaml)", themeName)
	}

	// Read palette for potential re-derivation.
	palette, err := uc.readPalette(themeName)
	if err != nil {
		return "", fmt.Errorf("read palette: %w", err)
	}

	// Get mtimes.
	paletteMtime, err := uc.store.FileMtime(themeName, "palette.yaml")
	if err != nil {
		return "", fmt.Errorf("palette mtime: %w", err)
	}

	universalMtime, universalErr := uc.store.FileMtime(themeName, "universal.yaml")

	// Determine what needs regeneration.
	// If universal.yaml is missing (error from FileMtime), force full regeneration
	// from palette rather than relying on zero-time coincidence.
	paletteEdited := universalErr != nil || paletteMtime.After(universalMtime)
	universalEdited := !paletteEdited && uc.isUniversalNewerThanMappings(themeName, universalMtime, targetFilter)

	// Resolve the theme for downstream operations.
	tokens := uc.deriver.Derive(palette)
	resolved := &domain.ResolvedTheme{
		Name:    palette.Name,
		Variant: palette.Variant,
		Palette: palette,
		Tokens:  tokens,
	}

	regenerated, err := uc.applyRegeneration(resolved, themeName, targetFilter, paletteEdited, universalEdited)
	if err != nil {
		return "", err
	}

	if len(regenerated) == 0 {
		return "nothing to do", nil
	}

	return fmt.Sprintf("regenerated %d files: %s", len(regenerated), strings.Join(regenerated, ", ")), nil
}

// applyRegeneration dispatches to the appropriate regeneration strategy based
// on which upstream file was edited.
func (uc *RegenerateThemeUseCase) applyRegeneration(
	resolved *domain.ResolvedTheme, themeName, targetFilter string,
	paletteEdited, universalEdited bool,
) ([]string, error) {
	switch {
	case paletteEdited:
		return uc.regenFromPalette(resolved, themeName, targetFilter)
	case universalEdited:
		return uc.regenerateTargets(resolved, themeName, targetFilter)
	default:
		return uc.regenerateStaleMappings(resolved, themeName, targetFilter)
	}
}

// regenFromPalette re-derives universal.yaml and all targets.
func (uc *RegenerateThemeUseCase) regenFromPalette(
	resolved *domain.ResolvedTheme, themeName, targetFilter string,
) ([]string, error) {
	if err := uc.writeUniversal(resolved.Tokens, themeName); err != nil {
		return nil, fmt.Errorf("write universal: %w", err)
	}

	regenerated := []string{"universal.yaml"}

	regen, err := uc.regenerateTargets(resolved, themeName, targetFilter)
	if err != nil {
		return regenerated, err
	}

	return append(regenerated, regen...), nil
}

// isUniversalNewerThanMappings checks if universal.yaml is newer than any mapping file.
func (uc *RegenerateThemeUseCase) isUniversalNewerThanMappings(themeName string, universalMtime time.Time, targetFilter string) bool {
	for _, tgt := range uc.targets {
		if targetFilter != "" && tgt.Mapper.Name() != targetFilter {
			continue
		}
		mappingMtime, err := uc.store.FileMtime(themeName, tgt.MappingFile)
		if err != nil {
			// File missing -> needs regen.
			return true
		}
		if universalMtime.After(mappingMtime) {
			return true
		}
	}
	return false
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

// regenerateStaleMappings checks individual mapping->output pairs and only
// regenerates outputs whose mapping is newer.
func (uc *RegenerateThemeUseCase) regenerateStaleMappings(resolved *domain.ResolvedTheme, themeName, targetFilter string) ([]string, error) {
	var regenerated []string
	var targetErrors []string

	for _, target := range uc.targets {
		if targetFilter != "" && target.Mapper.Name() != targetFilter {
			continue
		}

		mappingMtime, err := uc.store.FileMtime(themeName, target.MappingFile)
		if err != nil {
			continue
		}

		outputMtime, err := uc.store.FileMtime(themeName, target.Generator.DefaultFilename())
		if err != nil {
			// Output missing -> needs regen.
			if err := uc.processTarget(target, resolved, themeName); err != nil {
				targetErrors = append(targetErrors, fmt.Sprintf("%s: %v", target.Mapper.Name(), err))
				continue
			}
			regenerated = append(regenerated, target.Generator.DefaultFilename())
			continue
		}

		if mappingMtime.After(outputMtime) {
			if err := uc.processTarget(target, resolved, themeName); err != nil {
				targetErrors = append(targetErrors, fmt.Sprintf("%s: %v", target.Mapper.Name(), err))
				continue
			}
			regenerated = append(regenerated, target.Generator.DefaultFilename())
		}
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

// writeUniversal writes universal.yaml to the theme directory.
func (uc *RegenerateThemeUseCase) writeUniversal(tokens *domain.TokenSet, themeName string) error {
	w, err := uc.store.OpenWriter(themeName, "universal.yaml")
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	if uc.universalWriter != nil {
		return uc.universalWriter(w, tokens)
	}

	_, writeErr := fmt.Fprintf(w, "# universal for %s\n", themeName)
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

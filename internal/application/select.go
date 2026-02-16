package application

import (
	"fmt"
	"strings"

	"github.com/curtbushko/flair/internal/ports"
)

// outputFiles lists the output files that must exist before a theme can be selected.
var outputFiles = []string{
	"style.lua",
	"style.css",
	"gtk.css",
	"style.qss",
	"style.json",
}

// SelectThemeUseCase verifies that a theme directory exists and contains
// all required output files, then activates it via store.Select().
// When a theme is not on disk but matches a built-in palette, it is
// auto-generated before selection.
type SelectThemeUseCase struct {
	store     ports.ThemeStore
	builtins  ports.PaletteSource
	generator *GenerateThemeUseCase
}

// NewSelectThemeUseCase returns a new SelectThemeUseCase wired to the given
// store, builtins source, and generator. builtins and generator may be nil
// if auto-generation of built-in themes is not needed.
func NewSelectThemeUseCase(store ports.ThemeStore, builtins ports.PaletteSource, generator *GenerateThemeUseCase) *SelectThemeUseCase {
	return &SelectThemeUseCase{store: store, builtins: builtins, generator: generator}
}

// Execute verifies the theme is complete and selects it as active.
// If the theme does not exist on disk but matches a built-in palette,
// it is auto-generated first.
func (uc *SelectThemeUseCase) Execute(themeName string) error {
	if !uc.themeExists(themeName) {
		if err := uc.autoGenerate(themeName); err != nil {
			return err
		}
	}

	// Verify all output files exist.
	var missing []string
	for _, f := range outputFiles {
		if !uc.store.FileExists(themeName, f) {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("theme %q is incomplete, missing output files: %s", themeName, strings.Join(missing, ", "))
	}

	// Activate the theme via symlinks.
	if err := uc.store.Select(themeName); err != nil {
		return fmt.Errorf("select theme %q: %w", themeName, err)
	}

	return nil
}

// autoGenerate attempts to generate a built-in theme if the name matches a
// known palette. Returns an error if the theme is not a built-in.
func (uc *SelectThemeUseCase) autoGenerate(themeName string) error {
	if uc.builtins == nil || uc.generator == nil || !uc.builtins.Has(themeName) {
		return fmt.Errorf("theme %q not found", themeName)
	}
	if err := uc.generator.ExecuteBuiltin(themeName, "", ""); err != nil {
		return fmt.Errorf("auto-generate built-in theme %q: %w", themeName, err)
	}
	return nil
}

// themeExists checks whether the theme directory exists by testing whether
// at least one output file is present, or by checking the dir via ListThemes.
func (uc *SelectThemeUseCase) themeExists(themeName string) bool {
	// Check if any output file exists (fast path for a generated theme).
	for _, f := range outputFiles {
		if uc.store.FileExists(themeName, f) {
			return true
		}
	}

	// Fall back to listing themes to check for directory existence.
	themes, err := uc.store.ListThemes()
	if err != nil {
		return false
	}
	for _, name := range themes {
		if name == themeName {
			return true
		}
	}

	return false
}

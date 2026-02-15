package application

import (
	"fmt"

	"github.com/curtbushko/flair/internal/ports"
)

// ThemeInfo describes an installed theme for the list command.
type ThemeInfo struct {
	Name     string
	Selected bool
	Complete bool
}

// ListThemesUseCase lists installed themes and built-in palettes.
type ListThemesUseCase struct {
	store    ports.ThemeStore
	builtins ports.PaletteSource
}

// NewListThemesUseCase returns a new ListThemesUseCase wired to the given
// store and palette source.
func NewListThemesUseCase(store ports.ThemeStore, builtins ports.PaletteSource) *ListThemesUseCase {
	return &ListThemesUseCase{store: store, builtins: builtins}
}

// Execute returns a slice of ThemeInfo for all installed themes, marking which
// one is currently selected and whether each theme is complete (has all output files).
func (uc *ListThemesUseCase) Execute() ([]ThemeInfo, error) {
	names, err := uc.store.ListThemes()
	if err != nil {
		return nil, fmt.Errorf("list themes: %w", err)
	}

	selected, err := uc.store.SelectedTheme()
	if err != nil {
		return nil, fmt.Errorf("selected theme: %w", err)
	}

	themes := make([]ThemeInfo, 0, len(names))
	for _, name := range names {
		themes = append(themes, ThemeInfo{
			Name:     name,
			Selected: name == selected,
			Complete: uc.isComplete(name),
		})
	}

	return themes, nil
}

// ListBuiltins returns the names of all built-in palettes.
func (uc *ListThemesUseCase) ListBuiltins() []string {
	return uc.builtins.List()
}

// isComplete checks whether a theme has all required output files.
func (uc *ListThemesUseCase) isComplete(themeName string) bool {
	for _, f := range outputFiles {
		if !uc.store.FileExists(themeName, f) {
			return false
		}
	}
	return true
}

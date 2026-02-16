package application

import (
	"fmt"

	"github.com/curtbushko/flair/internal/ports"
)

// ThemeInfo describes a theme for the list command.
type ThemeInfo struct {
	Name      string
	Selected  bool
	Complete  bool
	Generated bool
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

// Execute returns a slice of ThemeInfo for generated themes merged with
// available built-in palettes. Generated themes are identified by having a
// palette.yaml file. Built-in palettes that have not been generated appear
// with Generated=false.
func (uc *ListThemesUseCase) Execute() ([]ThemeInfo, error) {
	names, err := uc.store.ListThemes()
	if err != nil {
		return nil, fmt.Errorf("list themes: %w", err)
	}

	selected, err := uc.store.SelectedTheme()
	if err != nil {
		return nil, fmt.Errorf("selected theme: %w", err)
	}

	seen := make(map[string]bool)
	var themes []ThemeInfo

	// Add generated themes (directories that contain palette.yaml).
	for _, name := range names {
		if !uc.store.FileExists(name, "palette.yaml") {
			continue
		}
		seen[name] = true
		themes = append(themes, ThemeInfo{
			Name:      name,
			Selected:  name == selected,
			Complete:  uc.isComplete(name),
			Generated: true,
		})
	}

	// Add built-in palettes that have not been generated.
	for _, name := range uc.builtins.List() {
		if seen[name] {
			continue
		}
		themes = append(themes, ThemeInfo{
			Name:      name,
			Generated: false,
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

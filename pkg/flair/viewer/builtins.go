// Package viewer provides a bubbletea-based TUI for browsing flair themes.
package viewer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/curtbushko/flair/pkg/flair"
	"github.com/curtbushko/flair/pkg/flair/palettes"
)

// BuiltinThemeLoader implements [ThemeLoader] for embedded built-in themes.
//
// It loads themes from the embedded palettes without requiring filesystem access.
// This is the standard loader for use with [RunBuiltins].
type BuiltinThemeLoader struct{}

// NewBuiltinThemeLoader creates a new BuiltinThemeLoader.
//
// The returned loader can load any theme available via [flair.ListBuiltins].
func NewBuiltinThemeLoader() *BuiltinThemeLoader {
	return &BuiltinThemeLoader{}
}

// LoadPalette loads the base24 colors for a built-in theme.
//
// The returned PaletteData contains 24 hex color strings corresponding
// to base00-base17 slots.
func (l *BuiltinThemeLoader) LoadPalette(name string) (PaletteData, error) {
	palette, err := loadRawPalette(name)
	if err != nil {
		return PaletteData{}, fmt.Errorf("load built-in palette %q: %w", name, err)
	}

	var pd PaletteData
	for i := 0; i < 24; i++ {
		if c := palette.Base(i); c != nil {
			pd.Colors[i] = c.Hex()
		}
	}

	return pd, nil
}

// loadRawPalette loads the raw palette (not tokenized) for a theme.
func loadRawPalette(name string) (*flair.Palette, error) {
	filename := name + ".yaml"
	data, err := palettes.EmbeddedFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("palette %q not found", name)
	}

	palette, err := flair.ParsePalette(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse palette %q: %w", name, err)
	}

	return palette, nil
}

// LoadTokens loads the semantic tokens for a built-in theme.
//
// The returned TokenData contains tokens grouped by category (Surface, Text,
// Status, Syntax, Diff, Statusline) with hex color values.
func (l *BuiltinThemeLoader) LoadTokens(name string) (TokenData, error) {
	theme, err := flair.LoadBuiltin(name)
	if err != nil {
		return TokenData{}, fmt.Errorf("load built-in tokens %q: %w", name, err)
	}

	td := TokenData{
		Surface:    make(map[string]string),
		Text:       make(map[string]string),
		Status:     make(map[string]string),
		Syntax:     make(map[string]string),
		Diff:       make(map[string]string),
		Statusline: make(map[string]string),
	}

	colors := theme.Colors()
	for key, color := range colors {
		hex := color.Hex()
		switch {
		case strings.HasPrefix(key, "surface."):
			td.Surface[key] = hex
		case strings.HasPrefix(key, "text."):
			td.Text[key] = hex
		case strings.HasPrefix(key, "status."):
			td.Status[key] = hex
		case strings.HasPrefix(key, "syntax."):
			td.Syntax[key] = hex
		case strings.HasPrefix(key, "diff."):
			td.Diff[key] = hex
		case strings.HasPrefix(key, "statusline."):
			td.Statusline[key] = hex
		}
	}

	return td, nil
}

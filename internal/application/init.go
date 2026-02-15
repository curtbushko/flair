package application

import (
	"fmt"
	"path/filepath"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// scaffoldPalette is the template written by flair init to create a new theme.
const scaffoldPalette = `schema_version: %d
system: "base24"
name: "%s"
author: "Your Name"
variant: "dark"
palette:
  # Background shades (darkest to lightest)
  base00: "000000"  # Default background
  base01: "111111"  # Lighter background (status bars, line numbers)
  base02: "222222"  # Selection background
  base03: "333333"  # Comments, invisibles, line highlighting
  # Foreground shades (darkest to lightest)
  base04: "999999"  # Dark foreground (status bars)
  base05: "cccccc"  # Default foreground, caret, delimiters
  base06: "dddddd"  # Light foreground
  base07: "eeeeee"  # Lightest foreground
  # Accent colors
  base08: "ff0000"  # Red       - Variables, XML tags, markup link text
  base09: "ff8800"  # Orange    - Integers, booleans, constants
  base0A: "ffff00"  # Yellow    - Classes, markup bold, search bg
  base0B: "00ff00"  # Green     - Strings, inherited class, markup code
  base0C: "00ffff"  # Cyan      - Support, regex, escape chars, markup quotes
  base0D: "0088ff"  # Blue      - Functions, methods, attribute IDs, headings
  base0E: "ff00ff"  # Magenta   - Keywords, storage, selector, markup italic
  base0F: "884400"  # Brown     - Deprecated, opening/closing embedded tags
  # Extended base24 slots
  base10: "0a0a0a"  # Deeper background
  base11: "050505"  # Deepest background
  base12: "ff4444"  # Bright red
  base13: "ffcc00"  # Bright yellow
  base14: "44ff44"  # Bright green
  base15: "44ffff"  # Bright cyan
  base16: "4488ff"  # Bright blue
  base17: "ff44ff"  # Bright magenta
`

// InitThemeUseCase creates a new theme directory with a scaffold palette.yaml.
type InitThemeUseCase struct {
	store ports.ThemeStore
}

// NewInitThemeUseCase returns a new InitThemeUseCase wired to the given store.
func NewInitThemeUseCase(store ports.ThemeStore) *InitThemeUseCase {
	return &InitThemeUseCase{store: store}
}

// Execute creates the theme directory and writes a scaffold palette.yaml.
// Returns the path to the created palette.yaml, or an error if it already exists.
func (uc *InitThemeUseCase) Execute(themeName string) (string, error) {
	// Check if palette.yaml already exists (before creating dir).
	if uc.store.FileExists(themeName, "palette.yaml") {
		return "", fmt.Errorf("palette.yaml already exists for theme %q", themeName)
	}

	// Create the theme directory.
	if err := uc.store.EnsureThemeDir(themeName); err != nil {
		return "", fmt.Errorf("init theme %q: %w", themeName, err)
	}

	// Write scaffold palette.yaml.
	w, err := uc.store.OpenWriter(themeName, "palette.yaml")
	if err != nil {
		return "", fmt.Errorf("init theme %q: %w", themeName, err)
	}

	content := fmt.Sprintf(scaffoldPalette, domain.SchemaPalette, themeName)
	if _, writeErr := w.Write([]byte(content)); writeErr != nil {
		_ = w.Close()
		return "", fmt.Errorf("write scaffold palette: %w", writeErr)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("close scaffold palette: %w", err)
	}

	palettePath := filepath.Join(uc.store.ThemeDir(themeName), "palette.yaml")
	return palettePath, nil
}

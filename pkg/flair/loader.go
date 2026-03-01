package flair

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// outputFiles lists the output files that are checked for symlinks.
var outputFiles = []string{
	"style.lua",
	"style.css",
	"gtk.css",
	"style.qss",
	"style.json",
}

// universalToken represents a single semantic token in universal.yaml.
// This is a private type to avoid exposing YAML details.
type universalToken struct {
	Color         string `yaml:"color"`
	Bold          bool   `yaml:"bold,omitempty"`
	Italic        bool   `yaml:"italic,omitempty"`
	Underline     bool   `yaml:"underline,omitempty"`
	Undercurl     bool   `yaml:"undercurl,omitempty"`
	Strikethrough bool   `yaml:"strikethrough,omitempty"`
}

// universalFile represents the structure of universal.yaml.
type universalFile struct {
	Tokens map[string]universalToken `yaml:"tokens"`
}

// ErrNoSelectedTheme is returned when no theme is currently selected.
var ErrNoSelectedTheme = errors.New("no theme selected")

// ErrThemeNotFound is returned when a requested theme does not exist.
var ErrThemeNotFound = errors.New("theme not found")

// Load loads the currently selected theme from the default config directory.
// It respects XDG_CONFIG_HOME if set, otherwise uses ~/.config/flair.
func Load() (*Theme, error) {
	return LoadFrom(defaultConfigDir())
}

// LoadFrom loads the currently selected theme from the specified config directory.
// The selected theme is determined by following symlinks (e.g., style.json -> themename/style.json).
func LoadFrom(configDir string) (*Theme, error) {
	selected, err := SelectedThemeFrom(configDir)
	if err != nil {
		return nil, fmt.Errorf("load theme: %w", err)
	}

	if selected == "" {
		return nil, ErrNoSelectedTheme
	}

	return LoadNamedFrom(configDir, selected)
}

// LoadNamed loads a specific theme by name from the default config directory.
func LoadNamed(name string) (*Theme, error) {
	return LoadNamedFrom(defaultConfigDir(), name)
}

// LoadNamedFrom loads a specific theme by name from the specified config directory.
func LoadNamedFrom(configDir, name string) (*Theme, error) {
	themeDir := filepath.Join(configDir, name)
	universalPath := filepath.Join(themeDir, "universal.yaml")

	data, err := os.ReadFile(universalPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("load theme %q: %w", name, ErrThemeNotFound)
		}
		return nil, fmt.Errorf("load theme %q: %w", name, err)
	}

	var uf universalFile
	if err := yaml.Unmarshal(data, &uf); err != nil {
		return nil, fmt.Errorf("load theme %q: parse universal.yaml: %w", name, err)
	}

	colors := make(map[string]Color, len(uf.Tokens))
	for tokenPath, token := range uf.Tokens {
		if token.Color == "" {
			// Skip tokens without a color value.
			continue
		}

		c, err := ParseHex(token.Color)
		if err != nil {
			return nil, fmt.Errorf("load theme %q: token %q: %w", name, tokenPath, err)
		}
		colors[tokenPath] = c
	}

	// Extract variant from theme name if present (e.g., "tokyo-night-dark" -> "dark").
	variant := extractVariant(name)

	return NewTheme(name, variant, colors), nil
}

// ListThemes returns the names of all available themes in the default config directory.
func ListThemes() ([]string, error) {
	return ListThemesFrom(defaultConfigDir())
}

// ListThemesFrom returns the names of all available themes in the specified config directory.
// Theme directories are identified by the presence of a universal.yaml file.
func ListThemesFrom(configDir string) ([]string, error) {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list themes: %w", err)
	}

	themes := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check for universal.yaml to confirm it's a valid theme directory.
		universalPath := filepath.Join(configDir, entry.Name(), "universal.yaml")
		if _, err := os.Stat(universalPath); err == nil {
			themes = append(themes, entry.Name())
		}
	}

	sort.Strings(themes)
	return themes, nil
}

// SelectedTheme returns the name of the currently selected theme from the default config directory.
// Returns an empty string if no theme is selected.
func SelectedTheme() (string, error) {
	return SelectedThemeFrom(defaultConfigDir())
}

// SelectedThemeFrom returns the name of the currently selected theme by reading symlinks.
// Returns an empty string if no theme is selected.
func SelectedThemeFrom(configDir string) (string, error) {
	for _, f := range outputFiles {
		link := filepath.Join(configDir, f)
		target, err := os.Readlink(link)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			// Not a symlink or other error, skip.
			continue
		}

		// Target is relative: <themeName>/<filename>
		themeName := filepath.Dir(target)
		if themeName != "." && themeName != "" {
			return themeName, nil
		}
	}

	return "", nil
}

// defaultConfigDir returns the default flair config directory.
// It respects XDG_CONFIG_HOME when set, falling back to ~/.config/flair.
func defaultConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "flair")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to relative path if home directory cannot be determined.
		return filepath.Join(".config", "flair")
	}

	return filepath.Join(home, ".config", "flair")
}

// extractVariant attempts to extract a variant name from the theme name.
// Common patterns: "tokyo-night-dark" -> "dark", "gruvbox-light" -> "light".
func extractVariant(themeName string) string {
	variants := []string{"dark", "light", "storm", "moon", "day", "night"}

	for _, v := range variants {
		suffix := "-" + v
		if len(themeName) > len(suffix) && themeName[len(themeName)-len(suffix):] == suffix {
			return v
		}
	}

	return ""
}

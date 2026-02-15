// Package config provides configuration types and defaults for the flair CLI.
package config

import (
	"os"
	"path/filepath"
)

// Config holds runtime configuration for flair.
type Config struct {
	// ConfigDir is the root directory where flair stores themes and output files.
	ConfigDir string
}

// DefaultConfigDir returns the default configuration directory.
// It respects XDG_CONFIG_HOME when set, falling back to ~/.config/flair.
func DefaultConfigDir() string {
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

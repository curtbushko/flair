package ports

import (
	"io"
	"time"
)

// ThemeStore manages theme directories and symlinks on the filesystem.
type ThemeStore interface {
	// ConfigDir returns the root config directory (e.g. ~/.config/flair).
	ConfigDir() string

	// ThemeDir returns the path for a named theme.
	ThemeDir(themeName string) string

	// EnsureThemeDir creates the theme directory if it doesn't exist.
	EnsureThemeDir(themeName string) error

	// ListThemes returns all theme directory names.
	ListThemes() ([]string, error)

	// SelectedTheme returns the currently symlinked theme name, or "" if none.
	SelectedTheme() (string, error)

	// Select creates/updates symlinks at the config root pointing to the
	// given theme's output files.
	Select(themeName string) error

	// OpenReader opens a file within a theme directory for reading.
	// The caller must close the returned reader.
	OpenReader(themeName, filename string) (io.ReadCloser, error)

	// OpenWriter opens (or creates) a file within a theme directory for writing.
	// The caller must close the returned writer.
	OpenWriter(themeName, filename string) (io.WriteCloser, error)

	// FileExists checks whether a file exists in a theme directory.
	FileExists(themeName, filename string) bool

	// FileMtime returns the modification time of a file.
	FileMtime(themeName, filename string) (time.Time, error)
}

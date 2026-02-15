package ports

import (
	"io"
	"time"
)

// ThemeDirProvider provides directory path information.
type ThemeDirProvider interface {
	// ConfigDir returns the root config directory (e.g. ~/.config/flair).
	ConfigDir() string

	// ThemeDir returns the path for a named theme.
	ThemeDir(themeName string) string
}

// ThemeManager manages theme directories and selection.
type ThemeManager interface {
	// EnsureThemeDir creates the theme directory if it doesn't exist.
	EnsureThemeDir(themeName string) error

	// ListThemes returns all theme directory names.
	ListThemes() ([]string, error)

	// SelectedTheme returns the currently symlinked theme name, or "" if none.
	SelectedTheme() (string, error)

	// Select creates/updates symlinks at the config root pointing to the
	// given theme's output files.
	Select(themeName string) error
}

// ThemeFileIO provides file read/write operations within theme directories.
type ThemeFileIO interface {
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

// ThemeStore combines all theme storage operations.
type ThemeStore interface {
	ThemeDirProvider
	ThemeManager
	ThemeFileIO
}

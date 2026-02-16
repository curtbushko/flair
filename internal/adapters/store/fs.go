// Package store provides a filesystem-based implementation of ports.ThemeStore.
// It manages theme directories, symlinks for the selected theme, and file I/O
// via io.ReadCloser / io.WriteCloser.
package store

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// outputFiles lists the output files that are symlinked when a theme is selected.
var outputFiles = []string{
	"style.lua",
	"style.css",
	"gtk.css",
	"style.qss",
	"style.json",
}

// FsStore implements ports.ThemeStore using the local filesystem.
type FsStore struct {
	configDir string
}

// NewFsStore returns a new FsStore rooted at configDir.
func NewFsStore(configDir string) *FsStore {
	return &FsStore{configDir: configDir}
}

// ConfigDir returns the root config directory.
func (s *FsStore) ConfigDir() string {
	return s.configDir
}

// ThemeDir returns the path for a named theme.
func (s *FsStore) ThemeDir(themeName string) string {
	return filepath.Join(s.configDir, themeName)
}

// EnsureThemeDir creates the theme directory if it doesn't exist.
func (s *FsStore) EnsureThemeDir(themeName string) error {
	dir := s.ThemeDir(themeName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("ensure theme dir %q: %w", themeName, err)
	}
	return nil
}

// ListThemes returns all theme directory names, sorted alphabetically.
// Returns an empty slice if the config directory does not exist.
func (s *FsStore) ListThemes() ([]string, error) {
	entries, err := os.ReadDir(s.configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list themes: %w", err)
	}

	themes := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			themes = append(themes, entry.Name())
		}
	}
	sort.Strings(themes)
	return themes, nil
}

// OpenReader opens a file within a theme directory for reading.
// The caller must close the returned reader.
func (s *FsStore) OpenReader(themeName, filename string) (io.ReadCloser, error) {
	path := filepath.Join(s.ThemeDir(themeName), filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open reader %s/%s: %w", themeName, filename, err)
	}
	return f, nil
}

// OpenWriter opens (or creates/truncates) a file within a theme directory for writing.
// The caller must close the returned writer.
func (s *FsStore) OpenWriter(themeName, filename string) (io.WriteCloser, error) {
	path := filepath.Join(s.ThemeDir(themeName), filename)
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("open writer %s/%s: %w", themeName, filename, err)
	}
	return f, nil
}

// FileExists checks whether a file exists in a theme directory.
func (s *FsStore) FileExists(themeName, filename string) bool {
	path := filepath.Join(s.ThemeDir(themeName), filename)
	_, err := os.Stat(path)
	return err == nil
}

// FileMtime returns the modification time of a file in a theme directory.
func (s *FsStore) FileMtime(themeName, filename string) (time.Time, error) {
	path := filepath.Join(s.ThemeDir(themeName), filename)
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("file mtime %s/%s: %w", themeName, filename, err)
	}
	return info.ModTime(), nil
}

// Select creates relative symlinks at the config root pointing to the given
// theme's output files. Existing symlinks are removed before creating new ones.
func (s *FsStore) Select(themeName string) error {
	for _, f := range outputFiles {
		if err := s.createSymlink(themeName, f); err != nil {
			return err
		}
	}
	return nil
}

// createSymlink creates a single relative symlink at the config root for a
// theme output file.
func (s *FsStore) createSymlink(themeName, filename string) error {
	link := filepath.Join(s.configDir, filename)
	target := filepath.Join(themeName, filename)

	// Remove existing symlink or file if present.
	if err := os.Remove(link); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove existing symlink %s: %w", filename, err)
	}

	if err := os.Symlink(target, link); err != nil {
		return fmt.Errorf("create symlink %s -> %s: %w", filename, target, err)
	}
	return nil
}

// SelectedTheme reads the first found symlink target to determine the currently
// selected theme name. Returns empty string and no error if no symlinks exist.
func (s *FsStore) SelectedTheme() (string, error) {
	for _, f := range outputFiles {
		link := filepath.Join(s.configDir, f)
		target, err := os.Readlink(link)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			// Not a symlink, skip.
			continue
		}
		// Target is relative: <themeName>/<filename>
		themeName := filepath.Dir(target)
		return themeName, nil
	}
	return "", nil
}

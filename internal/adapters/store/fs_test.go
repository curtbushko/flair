package store_test

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/curtbushko/flair/internal/adapters/store"
	"github.com/curtbushko/flair/internal/ports"
)

// writeAndClose is a test helper that writes data to a theme file and closes the writer.
func writeAndClose(t *testing.T, s *store.FsStore, themeName, filename string, data []byte) {
	t.Helper()
	w, err := s.OpenWriter(themeName, filename)
	if err != nil {
		t.Fatalf("OpenWriter(%s, %s): %v", themeName, filename, err)
	}
	_, err = w.Write(data)
	if err != nil {
		t.Fatalf("Write(%s/%s): %v", themeName, filename, err)
	}
	err = w.Close()
	if err != nil {
		t.Fatalf("Close(%s/%s): %v", themeName, filename, err)
	}
}

// TestFsStore_Paths verifies ConfigDir and ThemeDir return correct paths.
func TestFsStore_Paths(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	if got := s.ConfigDir(); got != dir {
		t.Errorf("ConfigDir() = %q, want %q", got, dir)
	}

	want := filepath.Join(dir, "test")
	if got := s.ThemeDir("test"); got != want {
		t.Errorf("ThemeDir('test') = %q, want %q", got, want)
	}
}

// TestFsStore_EnsureThemeDir verifies that EnsureThemeDir creates the theme directory.
func TestFsStore_EnsureThemeDir(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	if err := s.EnsureThemeDir("tokyonight"); err != nil {
		t.Fatalf("EnsureThemeDir() error = %v", err)
	}

	themeDir := filepath.Join(dir, "tokyonight")
	info, err := os.Stat(themeDir)
	if err != nil {
		t.Fatalf("theme dir does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected theme dir to be a directory")
	}
}

// TestFsStore_EnsureThemeDir_Idempotent verifies no error on second call.
func TestFsStore_EnsureThemeDir_Idempotent(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	if err := s.EnsureThemeDir("tokyonight"); err != nil {
		t.Fatalf("first EnsureThemeDir() error = %v", err)
	}
	if err := s.EnsureThemeDir("tokyonight"); err != nil {
		t.Fatalf("second EnsureThemeDir() error = %v", err)
	}
}

// TestFsStore_ListThemes verifies sorted directory listing.
func TestFsStore_ListThemes(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	// Create theme dirs
	for _, name := range []string{"tokyonight", "gruvbox"} {
		if err := os.MkdirAll(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", name, err)
		}
	}
	// Create a file (should be excluded from listing)
	if err := os.WriteFile(filepath.Join(dir, "somefile.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("create file: %v", err)
	}

	themes, err := s.ListThemes()
	if err != nil {
		t.Fatalf("ListThemes() error = %v", err)
	}

	want := []string{"gruvbox", "tokyonight"}
	if !sort.StringsAreSorted(themes) {
		t.Errorf("ListThemes() not sorted: %v", themes)
	}
	if len(themes) != len(want) {
		t.Fatalf("ListThemes() = %v, want %v", themes, want)
	}
	for i, w := range want {
		if themes[i] != w {
			t.Errorf("ListThemes()[%d] = %q, want %q", i, themes[i], w)
		}
	}
}

// TestFsStore_ImplementsInterface verifies FsStore implements ports.ThemeStore.
func TestFsStore_ImplementsInterface(t *testing.T) {
	var _ ports.ThemeStore = (*store.FsStore)(nil)
}

// TestFsStore_OpenWriter verifies that OpenWriter creates a writable file.
func TestFsStore_OpenWriter(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	if err := s.EnsureThemeDir("tokyonight"); err != nil {
		t.Fatalf("EnsureThemeDir: %v", err)
	}

	data := []byte("schema_version: 1\n")
	writeAndClose(t, s, "tokyonight", "universal.yaml", data)

	// Verify file exists on disk with correct content.
	got, err := os.ReadFile(filepath.Join(dir, "tokyonight", "universal.yaml"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("file content = %q, want %q", got, data)
	}
}

// TestFsStore_OpenReader verifies that OpenReader reads existing file.
func TestFsStore_OpenReader(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	themeDir := filepath.Join(dir, "tokyonight")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := []byte("tokens:\n  text.primary: '#7aa2f7'\n")
	if err := os.WriteFile(filepath.Join(themeDir, "universal.yaml"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	rc, err := s.OpenReader("tokyonight", "universal.yaml")
	if err != nil {
		t.Fatalf("OpenReader() error = %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("read content = %q, want %q", got, data)
	}
}

// TestFsStore_ReadWriteRoundTrip verifies write-then-read round trip.
func TestFsStore_ReadWriteRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	if err := s.EnsureThemeDir("tokyonight"); err != nil {
		t.Fatal(err)
	}

	data := []byte("round-trip content\n")
	writeAndClose(t, s, "tokyonight", "test.yaml", data)

	rc, err := s.OpenReader("tokyonight", "test.yaml")
	if err != nil {
		t.Fatalf("OpenReader: %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("round-trip content = %q, want %q", got, data)
	}
}

// TestFsStore_FileExistsAndMtime verifies FileExists and FileMtime.
func TestFsStore_FileExistsAndMtime(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	if err := s.EnsureThemeDir("tokyonight"); err != nil {
		t.Fatal(err)
	}

	// File does not exist yet.
	if s.FileExists("tokyonight", "missing.yaml") {
		t.Error("FileExists should return false for nonexistent file")
	}

	// Write a file.
	before := time.Now().Add(-time.Second) // allow 1s tolerance
	writeAndClose(t, s, "tokyonight", "test.yaml", []byte("data"))
	after := time.Now().Add(time.Second)

	if !s.FileExists("tokyonight", "test.yaml") {
		t.Error("FileExists should return true for existing file")
	}

	mtime, err := s.FileMtime("tokyonight", "test.yaml")
	if err != nil {
		t.Fatalf("FileMtime() error = %v", err)
	}
	if mtime.Before(before) || mtime.After(after) {
		t.Errorf("FileMtime() = %v, want between %v and %v", mtime, before, after)
	}
}

// TestFsStore_Select verifies symlink creation for theme output files.
func TestFsStore_Select(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	// Create theme dir with output files.
	themeDir := filepath.Join(dir, "tokyonight")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		if err := os.WriteFile(filepath.Join(themeDir, f), []byte("content-"+f), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := s.Select("tokyonight"); err != nil {
		t.Fatalf("Select() error = %v", err)
	}

	// Verify symlinks exist at config root pointing to theme files.
	for _, f := range outputFiles {
		link := filepath.Join(dir, f)
		target, err := os.Readlink(link)
		if err != nil {
			t.Errorf("Readlink(%s) error = %v", f, err)
			continue
		}
		// Target should be relative: tokyonight/<filename>
		wantTarget := filepath.Join("tokyonight", f)
		if target != wantTarget {
			t.Errorf("symlink %s -> %q, want %q", f, target, wantTarget)
		}
	}
}

// TestFsStore_Select_Replace verifies symlinks are replaced when selecting a new theme.
func TestFsStore_Select_Replace(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}

	// Create two themes.
	for _, theme := range []string{"gruvbox", "tokyonight"} {
		themeDir := filepath.Join(dir, theme)
		if err := os.MkdirAll(themeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, f := range outputFiles {
			if err := os.WriteFile(filepath.Join(themeDir, f), []byte(theme+"-"+f), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Select gruvbox first.
	if err := s.Select("gruvbox"); err != nil {
		t.Fatalf("Select(gruvbox) error = %v", err)
	}

	// Now select tokyonight.
	if err := s.Select("tokyonight"); err != nil {
		t.Fatalf("Select(tokyonight) error = %v", err)
	}

	// Verify symlinks now point to tokyonight.
	for _, f := range outputFiles {
		link := filepath.Join(dir, f)
		target, err := os.Readlink(link)
		if err != nil {
			t.Errorf("Readlink(%s) error = %v", f, err)
			continue
		}
		wantTarget := filepath.Join("tokyonight", f)
		if target != wantTarget {
			t.Errorf("symlink %s -> %q, want %q", f, target, wantTarget)
		}
	}
}

// TestFsStore_SelectedTheme verifies reading the currently selected theme.
func TestFsStore_SelectedTheme(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	// Create theme with output files and select it.
	themeDir := filepath.Join(dir, "tokyonight")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		if err := os.WriteFile(filepath.Join(themeDir, f), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := s.Select("tokyonight"); err != nil {
		t.Fatal(err)
	}

	got, err := s.SelectedTheme()
	if err != nil {
		t.Fatalf("SelectedTheme() error = %v", err)
	}
	if got != "tokyonight" {
		t.Errorf("SelectedTheme() = %q, want %q", got, "tokyonight")
	}
}

// TestFsStore_SelectedTheme_None verifies empty string when no symlinks exist.
func TestFsStore_SelectedTheme_None(t *testing.T) {
	dir := t.TempDir()
	s := store.NewFsStore(dir)

	got, err := s.SelectedTheme()
	if err != nil {
		t.Fatalf("SelectedTheme() error = %v", err)
	}
	if got != "" {
		t.Errorf("SelectedTheme() = %q, want empty string", got)
	}
}

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	// Enable dry-run mode for all tests to avoid TTY requirements.
	viewerDryRun = true
}

// TestSelect_WithThemeName verifies that 'flair select <theme>' applies symlinks.
func TestSelect_WithThemeName(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate the theme first.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runSelect([]string{"flair", "select", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Verify symlinks exist at config root.
	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		link := filepath.Join(dir, f)
		target, err := os.Readlink(link)
		if err != nil {
			t.Errorf("Readlink(%s) error = %v", f, err)
			continue
		}
		wantTarget := filepath.Join("tokyo-night-dark", f)
		if target != wantTarget {
			t.Errorf("symlink %s -> %q, want %q", f, target, wantTarget)
		}
	}
}

// TestSelect_NoArgs_LaunchesViewer verifies that 'flair select' with no args
// launches the viewer (stub message for now until task 10/viewer is implemented).
func TestSelect_NoArgs_LaunchesViewer(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate a theme so we have something to view.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runSelect([]string{"flair", "select", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Should indicate viewer mode (stub for now).
	output := stdout.String()
	if !strings.Contains(output, "viewer") {
		t.Errorf("stdout = %q, want it to mention viewer mode", output)
	}
}

// TestSelect_ViewerFlag verifies that '--viewer' launches viewer with theme pre-selected.
func TestSelect_ViewerFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate the theme first.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runSelect([]string{"flair", "select", "--viewer", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Should indicate viewer mode with theme pre-selected.
	output := stdout.String()
	if !strings.Contains(output, "viewer") {
		t.Errorf("stdout = %q, want it to mention viewer mode", output)
	}
	if !strings.Contains(output, "tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to mention pre-selected theme", output)
	}
}

// TestSelect_ViewerFlag_NoTheme verifies that '--viewer' without theme uses current selection.
func TestSelect_ViewerFlag_NoTheme(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate and select a theme first so there's a current selection.
	generateThemeForTest(t, dir, "tokyo-night-dark")
	_ = runSelect([]string{"flair", "select", "tokyo-night-dark", "--dir", dir}, &bytes.Buffer{}, &bytes.Buffer{})

	code := runSelect([]string{"flair", "select", "--viewer", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Should indicate viewer mode.
	output := stdout.String()
	if !strings.Contains(output, "viewer") {
		t.Errorf("stdout = %q, want it to mention viewer mode", output)
	}
}

func TestSelectCmd_PrintsConfirmation(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// First generate a theme so it has all output files.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runSelect([]string{"flair", "select", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain theme name", output)
	}
}

func TestSelectCmd_ThemeNotFound(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runSelect([]string{"flair", "select", "nonexistent", "--dir", dir}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "nonexistent") {
		t.Errorf("stderr = %q, want it to mention theme name", output)
	}
}

func TestSelectCmd_CreatesSymlinks(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate a theme first.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runSelect([]string{"flair", "select", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Verify symlinks exist at config root.
	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		link := filepath.Join(dir, f)
		target, err := os.Readlink(link)
		if err != nil {
			t.Errorf("Readlink(%s) error = %v", f, err)
			continue
		}
		wantTarget := filepath.Join("tokyo-night-dark", f)
		if target != wantTarget {
			t.Errorf("symlink %s -> %q, want %q", f, target, wantTarget)
		}
	}
}

func TestSelectCmd_BuiltinAutoGenerate(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Select a built-in theme WITHOUT generating first — should auto-generate.
	code := runSelect([]string{"flair", "select", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain theme name", output)
	}

	// Verify symlinks exist at config root.
	outputFiles := []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"}
	for _, f := range outputFiles {
		link := filepath.Join(dir, f)
		target, err := os.Readlink(link)
		if err != nil {
			t.Errorf("Readlink(%s) error = %v", f, err)
			continue
		}
		wantTarget := filepath.Join("tokyo-night-dark", f)
		if target != wantTarget {
			t.Errorf("symlink %s -> %q, want %q", f, target, wantTarget)
		}
	}
}

func TestSelectCmd_DirFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate theme first.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runSelect([]string{"flair", "select", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
}

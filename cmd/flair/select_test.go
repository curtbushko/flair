package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectCmd_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runSelect([]string{"flair", "select"}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "theme") || !strings.Contains(output, "required") {
		t.Errorf("stderr = %q, want it to mention theme argument is required", output)
	}
}

func TestSelectCmd_PrintsConfirmation(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// First generate a theme so it has all output files.
	genCode := runGenerate(
		[]string{"flair", "generate", "tokyo-night-dark", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if genCode != 0 {
		t.Fatalf("generate setup failed with exit code %d", genCode)
	}

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
	genCode := runGenerate(
		[]string{"flair", "generate", "tokyo-night-dark", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if genCode != 0 {
		t.Fatalf("generate setup failed with exit code %d", genCode)
	}

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

func TestSelectCmd_DirFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate theme first.
	genCode := runGenerate(
		[]string{"flair", "generate", "tokyo-night-dark", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if genCode != 0 {
		t.Fatalf("generate setup failed with exit code %d", genCode)
	}

	code := runSelect([]string{"flair", "select", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
}

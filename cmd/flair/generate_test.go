package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCmd_BuiltinName(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runGenerate([]string{"flair", "generate", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Should create theme dir with 12 files.
	themeDir := filepath.Join(dir, "tokyo-night-dark")
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		t.Fatalf("read theme dir: %v", err)
	}

	if len(entries) != 12 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected 12 files, got %d: %v", len(entries), names)
	}
}

func TestGenerateCmd_FilePath(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Use the test palette file from testdata.
	palettePath := filepath.Join("testdata", "tokyo-night-dark.yaml")
	// Resolve to absolute path relative to the project root.
	absPath, err := filepath.Abs(filepath.Join("..", "..", palettePath))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}

	code := runGenerate([]string{"flair", "generate", absPath, "--name", "my-theme", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	themeDir := filepath.Join(dir, "my-theme")
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		t.Fatalf("read theme dir: %v", err)
	}

	if len(entries) != 12 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected 12 files, got %d: %v", len(entries), names)
	}
}

func TestGenerateCmd_DirFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runGenerate([]string{"flair", "generate", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Theme dir should be under the custom dir.
	themeDir := filepath.Join(dir, "tokyo-night-dark")
	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		t.Errorf("expected theme dir %q to exist", themeDir)
	}
}

func TestGenerateCmd_TargetFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runGenerate([]string{"flair", "generate", "tokyo-night-dark", "--dir", dir, "--target", "stylix"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Should produce only: palette.yaml, universal.yaml, stylix-mapping.yaml, style.json = 4 files.
	themeDir := filepath.Join(dir, "tokyo-night-dark")
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		t.Fatalf("read theme dir: %v", err)
	}

	if len(entries) != 4 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected 4 files with --target stylix, got %d: %v", len(entries), names)
	}

	// Verify expected files exist.
	expectedFiles := []string{"palette.yaml", "universal.yaml", "stylix-mapping.yaml", "style.json"}
	for _, f := range expectedFiles {
		path := filepath.Join(themeDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %q to exist", f)
		}
	}
}

func TestGenerateCmd_InfersThemeName(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// When --name is not provided and using a builtin, theme name should be the builtin name.
	code := runGenerate([]string{"flair", "generate", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	themeDir := filepath.Join(dir, "tokyo-night-dark")
	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		t.Errorf("expected theme dir %q to exist (name inferred from builtin)", themeDir)
	}
}

func TestGenerateCmd_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runGenerate([]string{"flair", "generate"}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "palette") || !strings.Contains(output, "required") {
		t.Errorf("stderr = %q, want it to mention palette argument is required", output)
	}
}

func TestGenerateCmd_InvalidPalette(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runGenerate([]string{"flair", "generate", "nonexistent-file.yaml", "--dir", dir}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if output == "" {
		t.Error("expected error message on stderr, got empty")
	}
}

func TestGenerateCmd_PrintsSummary(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runGenerate([]string{"flair", "generate", "tokyo-night-dark", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Summary should mention the theme name and number of files.
	output := stdout.String()
	if !strings.Contains(output, "tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain theme name", output)
	}

	if !strings.Contains(output, "12") {
		t.Errorf("stdout = %q, want it to mention file count", output)
	}
}

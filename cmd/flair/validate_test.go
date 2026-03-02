package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCmd_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runValidate([]string{"flair", "validate"}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "theme") || !strings.Contains(output, "required") {
		t.Errorf("stderr = %q, want it to mention theme argument is required", output)
	}
}

func TestValidateCmd_PrintsViolations(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Create a theme directory with only palette.yaml (missing other files).
	themeName := "incomplete-theme"
	themeDir := filepath.Join(dir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Write a valid palette so parsing works, but leave other files missing.
	paletteData := []byte(`schema_version: 1
kind: palette
theme_name: incomplete-theme
system: base24
author: test
variant: dark
palette:
  base00: "16161e"
  base01: "1a1b26"
  base02: "2f3549"
  base03: "444b6a"
  base04: "787c99"
  base05: "a9b1d6"
  base06: "cbccd1"
  base07: "d5d6db"
  base08: "c0caf5"
  base09: "a9b1d6"
  base0A: "0db9d7"
  base0B: "9ece6a"
  base0C: "b4f9f8"
  base0D: "2ac3de"
  base0E: "bb9af7"
  base0F: "f7768e"
  base10: "16161e"
  base11: "16161e"
  base12: "c0caf5"
  base13: "0db9d7"
  base14: "9ece6a"
  base15: "b4f9f8"
  base16: "2ac3de"
  base17: "bb9af7"
`)
	if err := os.WriteFile(filepath.Join(themeDir, "palette.yaml"), paletteData, 0o644); err != nil {
		t.Fatalf("write palette: %v", err)
	}

	code := runValidate([]string{"flair", "validate", themeName, "--dir", dir}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1 (violations found)", code)
	}

	output := stdout.String()
	// Each violation should be on a separate line.
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		t.Errorf("expected multiple violation lines, got %d: %q", len(lines), output)
	}

	// Should mention missing files.
	if !strings.Contains(output, "missing") {
		t.Errorf("stdout = %q, want it to mention missing files", output)
	}
}

func TestValidateCmd_ExitCodeOnFailure(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Create an empty theme directory (all files missing).
	themeName := "empty-theme"
	themeDir := filepath.Join(dir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	code := runValidate([]string{"flair", "validate", themeName, "--dir", dir}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1 when violations found", code)
	}
}

func TestValidateCmd_ValidTheme(t *testing.T) {
	// Create a theme directory with all required files and proper headers.
	dir := t.TempDir()
	themeName := "valid-theme"
	themeDir := filepath.Join(dir, themeName)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Write a valid palette.yaml with schema_version header.
	paletteData := []byte(`schema_version: 1
kind: palette
theme_name: valid-theme
system: base24
author: test
variant: dark
palette:
  base00: "16161e"
  base01: "1a1b26"
  base02: "2f3549"
  base03: "444b6a"
  base04: "787c99"
  base05: "a9b1d6"
  base06: "cbccd1"
  base07: "d5d6db"
  base08: "c0caf5"
  base09: "a9b1d6"
  base0A: "0db9d7"
  base0B: "9ece6a"
  base0C: "b4f9f8"
  base0D: "2ac3de"
  base0E: "bb9af7"
  base0F: "f7768e"
  base10: "16161e"
  base11: "16161e"
  base12: "c0caf5"
  base13: "0db9d7"
  base14: "9ece6a"
  base15: "b4f9f8"
  base16: "2ac3de"
  base17: "bb9af7"
`)
	if err := os.WriteFile(filepath.Join(themeDir, "palette.yaml"), paletteData, 0o644); err != nil {
		t.Fatalf("write palette: %v", err)
	}

	// Write tokens.yaml with schema_version header.
	tokensData := []byte("schema_version: 1\nkind: tokens\ntokens: {}\n")
	if err := os.WriteFile(filepath.Join(themeDir, "tokens.yaml"), tokensData, 0o644); err != nil {
		t.Fatalf("write universal: %v", err)
	}

	// Write all expected mapping and output files.
	for _, f := range []string{
		"vim-mapping.yaml", "css-mapping.yaml", "gtk-mapping.yaml",
		"qss-mapping.yaml", "stylix-mapping.yaml",
		"style.lua", "style.css", "gtk.css", "style.qss", "style.json",
	} {
		if err := os.WriteFile(filepath.Join(themeDir, f), []byte("content"), 0o644); err != nil {
			t.Fatalf("write %s: %v", f, err)
		}
	}

	var stdout, stderr bytes.Buffer
	code := runValidate([]string{"flair", "validate", themeName, "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0 for valid theme; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "valid") {
		t.Errorf("stdout = %q, want it to confirm theme is valid", output)
	}
}

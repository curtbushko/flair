package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPreviewCmd_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runPreview([]string{"flair", "preview"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected non-zero exit code for missing theme name")
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "usage") && !strings.Contains(errOutput, "Usage") &&
		!strings.Contains(errOutput, "theme name") {
		t.Errorf("stderr = %q, want usage or error message about missing theme name", errOutput)
	}
}

func TestPreviewCmd_ValidTheme(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate a theme first.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runPreview(
		[]string{"flair", "preview", "tokyo-night-dark", "--dir", dir},
		&stdout, &stderr,
	)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()

	// Should contain ANSI escape sequences.
	if !strings.Contains(output, "\x1b[") {
		t.Error("output should contain ANSI escape sequences")
	}

	// Should contain semantic token names from tokens.yaml.
	if !strings.Contains(output, "surface.background") {
		t.Error("output should contain semantic token names like surface.background")
	}
}

func TestPreviewCmd_DirFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate a theme first.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	code := runPreview(
		[]string{"flair", "preview", "tokyo-night-dark", "--dir", dir},
		&stdout, &stderr,
	)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	if stdout.Len() == 0 {
		t.Error("expected output from preview command")
	}
}

func TestPreviewCmd_BuiltinWithoutGenerate(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Preview a built-in theme WITHOUT generating first.
	// This should work by deriving tokens on-the-fly.
	code := runPreview(
		[]string{"flair", "preview", "tokyo-night-dark", "--dir", dir},
		&stdout, &stderr,
	)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()

	// Should contain ANSI escape sequences.
	if !strings.Contains(output, "\x1b[") {
		t.Error("output should contain ANSI escape sequences")
	}

	// Should contain semantic token names (derived on-the-fly).
	if !strings.Contains(output, "surface.background") {
		t.Error("output should contain semantic token names like surface.background")
	}

	// Should contain palette colors.
	if !strings.Contains(output, "base00") {
		t.Error("output should contain palette slot names")
	}
}

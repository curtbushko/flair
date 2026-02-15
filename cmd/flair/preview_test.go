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
	genCode := runGenerate(
		[]string{"flair", "generate", "tokyo-night-dark", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if genCode != 0 {
		t.Fatalf("generate setup failed with exit code %d", genCode)
	}

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

	// Should contain semantic token names from universal.yaml.
	if !strings.Contains(output, "surface.background") {
		t.Error("output should contain semantic token names like surface.background")
	}
}

func TestPreviewCmd_DirFlag(t *testing.T) {
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

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestListCmd_FormatsOutput(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate two themes.
	genCode := runGenerate(
		[]string{"flair", "generate", "tokyo-night-dark", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if genCode != 0 {
		t.Fatalf("generate setup failed with exit code %d", genCode)
	}

	genCode = runGenerate(
		[]string{"flair", "generate", "catppuccin-mocha", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if genCode != 0 {
		t.Fatalf("generate setup failed with exit code %d", genCode)
	}

	// Select one theme.
	selCode := runSelect(
		[]string{"flair", "select", "tokyo-night-dark", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if selCode != 0 {
		t.Fatalf("select setup failed with exit code %d", selCode)
	}

	code := runList([]string{"flair", "list", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()

	// Selected theme should have "* " prefix.
	if !strings.Contains(output, "* tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain '* tokyo-night-dark'", output)
	}

	// Non-selected theme should have "  " prefix.
	if !strings.Contains(output, "  catppuccin-mocha") {
		t.Errorf("stdout = %q, want it to contain '  catppuccin-mocha'", output)
	}
}

func TestListCmd_BuiltinsFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runList([]string{"flair", "list", "--builtins", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()

	// Should contain at least one built-in palette name.
	if !strings.Contains(output, "tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain a built-in palette name", output)
	}

	// Builtins should print plain names, not prefixed with "* " selection markers.
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "* ") {
			t.Errorf("builtin output should not have selection marker, got %q", line)
		}
	}
}

func TestListCmd_NoThemes(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runList([]string{"flair", "list", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "No themes installed") {
		t.Errorf("stdout = %q, want it to contain helpful message about no themes", output)
	}
}

func TestListCmd_DirFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Generate a theme in custom dir.
	genCode := runGenerate(
		[]string{"flair", "generate", "tokyo-night-dark", "--dir", dir},
		&bytes.Buffer{}, &bytes.Buffer{},
	)
	if genCode != 0 {
		t.Fatalf("generate setup failed with exit code %d", genCode)
	}

	code := runList([]string{"flair", "list", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain theme name", output)
	}
}

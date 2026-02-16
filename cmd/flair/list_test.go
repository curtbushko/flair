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

	// Selected theme should have "* " prefix (no [available] since it was generated).
	if !strings.Contains(output, "* tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain '* tokyo-night-dark'", output)
	}

	// Generated non-selected theme should have "  " prefix without [available].
	if !strings.Contains(output, "  catppuccin-mocha") {
		t.Errorf("stdout = %q, want it to contain '  catppuccin-mocha'", output)
	}

	// Generated themes should NOT have [available] marker.
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip lines for ungenerated builtins.
		if strings.Contains(line, "[available]") {
			continue
		}
		// Generated themes should not have [available].
		if strings.Contains(trimmed, "tokyo-night-dark") && strings.Contains(line, "[available]") {
			t.Errorf("generated theme should not have [available] marker: %q", line)
		}
	}
}

func TestListCmd_ShowsAvailableBuiltins(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Do not generate any themes - just list.
	code := runList([]string{"flair", "list", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()

	// Built-in palettes should appear with [available] marker.
	if !strings.Contains(output, "tokyo-night-dark [available]") {
		t.Errorf("stdout = %q, want it to contain 'tokyo-night-dark [available]'", output)
	}
	if !strings.Contains(output, "catppuccin-mocha [available]") {
		t.Errorf("stdout = %q, want it to contain 'catppuccin-mocha [available]'", output)
	}
	if !strings.Contains(output, "gruvbox-dark [available]") {
		t.Errorf("stdout = %q, want it to contain 'gruvbox-dark [available]'", output)
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

func TestListCmd_NoGeneratedThemesShowsBuiltins(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runList([]string{"flair", "list", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()
	// With no generated themes, built-in palettes should still appear.
	if !strings.Contains(output, "[available]") {
		t.Errorf("stdout = %q, want it to contain built-in palettes with [available]", output)
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
	// Generated theme should appear without [available].
	if !strings.Contains(output, "tokyo-night-dark") {
		t.Errorf("stdout = %q, want it to contain theme name", output)
	}
	// The generated theme line should NOT have [available].
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "tokyo-night-dark") && strings.Contains(line, "[available]") {
			t.Errorf("generated theme should not have [available] marker: %q", line)
		}
	}
}

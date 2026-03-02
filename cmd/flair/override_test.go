package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunOverride_AddColor(t *testing.T) {
	dir := t.TempDir()

	// Generate a theme to work with.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	var stdout, stderr bytes.Buffer

	// Run override command to add color override.
	code := runOverride([]string{
		"flair", "override", "tokyo-night-dark", "syntax.keyword", "#ff00ff",
		"--dir", dir,
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Read palette.yaml and verify override was added.
	palettePath := filepath.Join(dir, "tokyo-night-dark", "palette.yaml")
	content, err := os.ReadFile(palettePath)
	if err != nil {
		t.Fatalf("read palette: %v", err)
	}

	if !strings.Contains(string(content), "overrides:") {
		t.Errorf("palette.yaml missing overrides section")
	}
	if !strings.Contains(string(content), "syntax.keyword:") {
		t.Errorf("palette.yaml missing syntax.keyword override")
	}
	if !strings.Contains(string(content), "ff00ff") {
		t.Errorf("palette.yaml missing color value ff00ff")
	}
}

func TestRunOverride_AddStyle(t *testing.T) {
	dir := t.TempDir()

	generateThemeForTest(t, dir, "tokyo-night-dark")

	var stdout, stderr bytes.Buffer

	// Run override command with --bold flag.
	code := runOverride([]string{
		"flair", "override", "tokyo-night-dark", "syntax.keyword", "--bold",
		"--dir", dir,
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Read palette.yaml and verify style override was added.
	palettePath := filepath.Join(dir, "tokyo-night-dark", "palette.yaml")
	content, err := os.ReadFile(palettePath)
	if err != nil {
		t.Fatalf("read palette: %v", err)
	}

	if !strings.Contains(string(content), "syntax.keyword:") {
		t.Errorf("palette.yaml missing syntax.keyword override")
	}
	if !strings.Contains(string(content), "bold: true") {
		t.Errorf("palette.yaml missing bold: true")
	}
}

func TestRunOverride_List(t *testing.T) {
	dir := t.TempDir()

	generateThemeForTest(t, dir, "tokyo-night-dark")

	// First add two overrides.
	var stdout, stderr bytes.Buffer
	code := runOverride([]string{
		"flair", "override", "tokyo-night-dark", "syntax.keyword", "#ff00ff",
		"--dir", dir,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("add first override failed: %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = runOverride([]string{
		"flair", "override", "tokyo-night-dark", "syntax.string", "#00ff00", "--italic",
		"--dir", dir,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("add second override failed: %s", stderr.String())
	}

	// Now list overrides.
	stdout.Reset()
	stderr.Reset()
	code = runOverride([]string{
		"flair", "override", "tokyo-night-dark", "--list",
		"--dir", dir,
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("list failed: %s", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "syntax.keyword") {
		t.Errorf("output = %q, missing syntax.keyword", output)
	}
	if !strings.Contains(output, "syntax.string") {
		t.Errorf("output = %q, missing syntax.string", output)
	}
}

func TestRunOverride_Remove(t *testing.T) {
	dir := t.TempDir()

	generateThemeForTest(t, dir, "tokyo-night-dark")

	// First add an override.
	var stdout, stderr bytes.Buffer
	code := runOverride([]string{
		"flair", "override", "tokyo-night-dark", "syntax.keyword", "#ff00ff",
		"--dir", dir,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("add override failed: %s", stderr.String())
	}

	// Now remove it.
	stdout.Reset()
	stderr.Reset()
	code = runOverride([]string{
		"flair", "override", "tokyo-night-dark", "--remove", "syntax.keyword",
		"--dir", dir,
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("remove failed: %s", stderr.String())
	}

	// Read palette.yaml and verify override was removed.
	palettePath := filepath.Join(dir, "tokyo-night-dark", "palette.yaml")
	content, err := os.ReadFile(palettePath)
	if err != nil {
		t.Fatalf("read palette: %v", err)
	}

	// Should either not have overrides section or not have syntax.keyword.
	contentStr := string(content)
	hasOverrides := strings.Contains(contentStr, "overrides:")
	hasKeyword := strings.Contains(contentStr, "syntax.keyword:")

	if hasKeyword {
		t.Errorf("palette.yaml should not contain syntax.keyword after removal")
	}
	// If it has overrides but no keyword, that's OK (empty overrides).
	_ = hasOverrides
}

func TestRunOverride_Update(t *testing.T) {
	dir := t.TempDir()

	generateThemeForTest(t, dir, "tokyo-night-dark")

	// Add initial override.
	var stdout, stderr bytes.Buffer
	code := runOverride([]string{
		"flair", "override", "tokyo-night-dark", "syntax.keyword", "#ff0000",
		"--dir", dir,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("add override failed: %s", stderr.String())
	}

	// Update to different color.
	stdout.Reset()
	stderr.Reset()
	code = runOverride([]string{
		"flair", "override", "tokyo-night-dark", "syntax.keyword", "#00ff00",
		"--dir", dir,
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("update failed: %s", stderr.String())
	}

	// Read palette.yaml and verify new color.
	palettePath := filepath.Join(dir, "tokyo-night-dark", "palette.yaml")
	content, err := os.ReadFile(palettePath)
	if err != nil {
		t.Fatalf("read palette: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "00ff00") {
		t.Errorf("palette.yaml missing updated color 00ff00")
	}
	if strings.Contains(contentStr, "ff0000") {
		t.Errorf("palette.yaml should not contain old color ff0000")
	}
}

func TestRunOverride_Help(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runOverride([]string{"flair", "override", "--help"}, &stdout, &stderr)

	// Help should return 0.
	if code != 0 {
		t.Errorf("exit code = %d, want 0 for help", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "Usage:") {
		t.Errorf("output = %q, missing Usage", output)
	}
	if !strings.Contains(output, "override") {
		t.Errorf("output = %q, missing 'override'", output)
	}
}

func TestRunOverride_InvalidTheme(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runOverride([]string{
		"flair", "override", "nonexistent", "syntax.keyword", "#ff00ff",
		"--dir", dir,
	}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1 for invalid theme", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "error") || !strings.Contains(output, "nonexistent") {
		t.Errorf("stderr = %q, want error mentioning nonexistent theme", output)
	}
}

func TestRunOverride_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runOverride([]string{"flair", "override"}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1 when no args", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "Usage") || !strings.Contains(output, "theme") {
		t.Errorf("stderr = %q, want usage mentioning theme", output)
	}
}

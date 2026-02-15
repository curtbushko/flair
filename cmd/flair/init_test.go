package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCmd_NameFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runInit([]string{"flair", "init", "--name", "my-theme", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Theme directory should exist.
	themeDir := filepath.Join(dir, "my-theme")
	info, err := os.Stat(themeDir)
	if err != nil {
		t.Fatalf("theme dir not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", themeDir)
	}

	// palette.yaml should exist.
	palettePath := filepath.Join(themeDir, "palette.yaml")
	if _, err := os.Stat(palettePath); err != nil {
		t.Fatalf("palette.yaml not created: %v", err)
	}
}

func TestInitCmd_NoName(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runInit([]string{"flair", "init", "--dir", dir}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "name") || !strings.Contains(output, "required") {
		t.Errorf("stderr = %q, want it to mention --name is required", output)
	}
}

func TestInitCmd_PrintsConfirmation(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runInit([]string{"flair", "init", "--name", "my-theme", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "palette.yaml") {
		t.Errorf("stdout = %q, want it to contain palette.yaml path", output)
	}
}

func TestInitCmd_AlreadyExists(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Create the theme first.
	code := runInit([]string{"flair", "init", "--name", "existing", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("first init failed: exit code = %d; stderr: %s", code, stderr.String())
	}

	// Try again -- should fail.
	stdout.Reset()
	stderr.Reset()
	code = runInit([]string{"flair", "init", "--name", "existing", "--dir", dir}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "already exists") {
		t.Errorf("stderr = %q, want it to mention already exists", output)
	}
}

func TestInitCmd_DirFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runInit([]string{"flair", "init", "--name", "custom-dir-theme", "--dir", dir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	// Verify palette.yaml was created under the custom dir.
	palettePath := filepath.Join(dir, "custom-dir-theme", "palette.yaml")
	if _, err := os.Stat(palettePath); err != nil {
		t.Fatalf("palette.yaml not created at %s: %v", palettePath, err)
	}
}

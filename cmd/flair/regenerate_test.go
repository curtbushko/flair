package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegenerateCmd_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runRegenerate([]string{"flair", "regenerate"}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "theme") || !strings.Contains(output, "required") {
		t.Errorf("stderr = %q, want it to mention theme argument is required", output)
	}
}

func TestRegenerateCmd_TargetFlag(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// First generate a theme to regenerate.
	generateThemeForTest(t, dir, "tokyo-night-dark")

	// Verify theme was created.
	themeDir := filepath.Join(dir, "tokyo-night-dark")
	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		t.Fatalf("theme dir %q not found", themeDir)
	}

	// Copy real palette YAML into the theme dir so regenerate can parse it.
	copyBuiltinPalette(t, themeDir, "tokyo-night-dark")

	// Now regenerate with --target vim.
	code := runRegenerate([]string{"flair", "regenerate", "tokyo-night-dark", "--target", "vim", "--dir", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("regenerate exit code = %d, want 0; stderr: %s", code, stderr.String())
	}

	output := stdout.String()
	if output == "" {
		t.Error("expected output on stdout, got empty")
	}
}

func TestRegenerateCmd_NonexistentTheme(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := runRegenerate([]string{"flair", "regenerate", "does-not-exist", "--dir", dir}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if output == "" {
		t.Error("expected error message on stderr, got empty")
	}
}

// copyBuiltinPalette reads the built-in palette and writes it to palette.yaml
// in the theme directory so regenerate can parse it.
func copyBuiltinPalette(t *testing.T, themeDir, paletteName string) {
	t.Helper()

	app := Wire(filepath.Dir(themeDir))
	r, err := app.Builtins.Get(paletteName)
	if err != nil {
		t.Fatalf("get builtin palette: %v", err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read builtin palette: %v", err)
	}

	if err := os.WriteFile(filepath.Join(themeDir, "palette.yaml"), data, 0o644); err != nil {
		t.Fatalf("write palette.yaml: %v", err)
	}
}

func TestRegenerateCmd_HelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runRegenerate([]string{"flair", "regenerate", "--help"}, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
}

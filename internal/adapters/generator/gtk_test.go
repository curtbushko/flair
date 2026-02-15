package generator_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/generator"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// sampleGtkTheme returns a GtkTheme with color definitions and widget rules.
func sampleGtkTheme() *ports.GtkTheme {
	return &ports.GtkTheme{
		Colors: []ports.GtkColorDef{
			{Name: "window_bg_color", Value: "#1a1b26"},
			{Name: "window_fg_color", Value: "#c0caf5"},
			{Name: "accent_bg_color", Value: "#7aa2f7"},
		},
		Rules: []ports.CSSRule{
			{
				Selector: "window",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "@window_bg_color"},
					{Property: "color", Value: "@window_fg_color"},
				},
			},
			{
				Selector: "headerbar",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "@headerbar_bg_color"},
					{Property: "color", Value: "@headerbar_fg_color"},
				},
			},
		},
	}
}

// TestGtkGenerator_Interface verifies that the GTK generator implements
// ports.Generator and returns the expected Name() and DefaultFilename().
func TestGtkGenerator_Interface(t *testing.T) {
	g := generator.NewGtk()

	// Compile-time interface check.
	var _ ports.Generator = g

	if name := g.Name(); name != "gtk" {
		t.Errorf("Name() = %q, want %q", name, "gtk")
	}

	if filename := g.DefaultFilename(); filename != "gtk.css" {
		t.Errorf("DefaultFilename() = %q, want %q", filename, "gtk.css")
	}
}

// TestGtkGenerator_Output verifies that Generate produces @define-color lines
// first, followed by widget selector rules in valid CSS format.
func TestGtkGenerator_Output(t *testing.T) {
	g := generator.NewGtk()
	theme := sampleGtkTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Must contain @define-color declarations.
	if !strings.Contains(output, "@define-color window_bg_color #1a1b26;") {
		t.Error("output missing '@define-color window_bg_color #1a1b26;' declaration")
	}
	if !strings.Contains(output, "@define-color window_fg_color #c0caf5;") {
		t.Error("output missing '@define-color window_fg_color #c0caf5;' declaration")
	}
	if !strings.Contains(output, "@define-color accent_bg_color #7aa2f7;") {
		t.Error("output missing '@define-color accent_bg_color #7aa2f7;' declaration")
	}

	// Must contain widget selector rules.
	if !strings.Contains(output, "window {") {
		t.Error("output does not contain 'window {' rule")
	}
	if !strings.Contains(output, "headerbar {") {
		t.Error("output does not contain 'headerbar {' rule")
	}

	// Widget rules must contain property declarations.
	if !strings.Contains(output, "background-color: @window_bg_color;") {
		t.Error("output missing 'background-color: @window_bg_color;' declaration")
	}
	if !strings.Contains(output, "color: @window_fg_color;") {
		t.Error("output missing 'color: @window_fg_color;' in window rule")
	}

	// @define-color lines must come before selector rules.
	defineIdx := strings.Index(output, "@define-color")
	windowIdx := strings.Index(output, "window {")
	if defineIdx < 0 || windowIdx < 0 {
		t.Fatal("output missing @define-color or window block")
	}
	if windowIdx <= defineIdx {
		t.Error("widget rules should appear after @define-color declarations")
	}
}

// TestGtkGenerator_WrongType verifies that passing a non-GtkTheme value
// as MappedTheme returns a GenerateError.
func TestGtkGenerator_WrongType(t *testing.T) {
	g := generator.NewGtk()

	var buf bytes.Buffer

	// Pass a string instead of *ports.GtkTheme.
	err := g.Generate(&buf, "not a gtk theme")
	if err == nil {
		t.Fatal("Generate() with wrong type should return error, got nil")
	}

	var genErr *domain.GenerateError
	if !errors.As(err, &genErr) {
		t.Errorf("error type = %T, want *domain.GenerateError", err)
	}
}

// TestGtkGenerator_Deterministic verifies that generating the same GtkTheme
// twice produces byte-identical output.
func TestGtkGenerator_Deterministic(t *testing.T) {
	g := generator.NewGtk()
	theme := sampleGtkTheme()

	var buf1, buf2 bytes.Buffer
	if err := g.Generate(&buf1, theme); err != nil {
		t.Fatalf("first Generate() error: %v", err)
	}
	if err := g.Generate(&buf2, theme); err != nil {
		t.Fatalf("second Generate() error: %v", err)
	}

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Error("output is not deterministic; two runs produced different results")
		t.Logf("run 1:\n%s", buf1.String())
		t.Logf("run 2:\n%s", buf2.String())
	}
}

// TestGtkGenerator_CommentHeader verifies that the output starts with
// a GTK-specific comment header.
func TestGtkGenerator_CommentHeader(t *testing.T) {
	g := generator.NewGtk()
	theme := sampleGtkTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Output should start with a comment.
	if !strings.HasPrefix(output, "/*") {
		t.Error("output should start with a CSS comment header")
	}
}

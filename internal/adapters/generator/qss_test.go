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

// sampleQssTheme returns a QssTheme with widget rules and pseudo-state rules
// using literal hex color values.
func sampleQssTheme() *ports.QssTheme {
	return &ports.QssTheme{
		Rules: []ports.CSSRule{
			{
				Selector: "QWidget",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "#1a1b26"},
					{Property: "color", Value: "#c0caf5"},
				},
			},
			{
				Selector: "QMainWindow",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "#1a1b26"},
					{Property: "color", Value: "#c0caf5"},
				},
			},
			{
				Selector: "QPushButton",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "#1f2335"},
					{Property: "color", Value: "#c0caf5"},
					{Property: "border", Value: "1px solid #3b4261"},
				},
			},
			{
				Selector: "QPushButton:hover",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "#292e42"},
					{Property: "color", Value: "#c0caf5"},
				},
			},
			{
				Selector: "QPushButton:pressed",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "#3b4261"},
					{Property: "color", Value: "#c0caf5"},
				},
			},
			{
				Selector: "QLineEdit",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "#16161e"},
					{Property: "color", Value: "#c0caf5"},
					{Property: "border", Value: "1px solid #3b4261"},
				},
			},
			{
				Selector: "QLineEdit:focus",
				Properties: []ports.CSSProperty{
					{Property: "border", Value: "1px solid #5c77bb"},
				},
			},
		},
	}
}

// TestQssGenerator_Interface verifies that the QSS generator implements
// ports.Generator and returns the expected Name() and DefaultFilename().
func TestQssGenerator_Interface(t *testing.T) {
	g := generator.NewQss()

	// Compile-time interface check.
	var _ ports.Generator = g

	if name := g.Name(); name != "qss" {
		t.Errorf("Name() = %q, want %q", name, "qss")
	}

	if filename := g.DefaultFilename(); filename != "style.qss" {
		t.Errorf("DefaultFilename() = %q, want %q", filename, "style.qss")
	}
}

// TestQssGenerator_LiteralHex verifies that Generate produces output with
// literal hex values and no variable references (no var(), no @define-color).
func TestQssGenerator_LiteralHex(t *testing.T) {
	g := generator.NewQss()
	theme := sampleQssTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Must contain literal hex values.
	if !strings.Contains(output, "#1a1b26") {
		t.Error("output missing literal hex value '#1a1b26'")
	}
	if !strings.Contains(output, "#c0caf5") {
		t.Error("output missing literal hex value '#c0caf5'")
	}
	if !strings.Contains(output, "#1f2335") {
		t.Error("output missing literal hex value '#1f2335'")
	}

	// Must NOT contain variable references.
	if strings.Contains(output, "var(") {
		t.Error("output contains 'var(' -- QSS must use literal hex, not CSS variables")
	}
	if strings.Contains(output, "@define-color") {
		t.Error("output contains '@define-color' -- QSS must use literal hex, not GTK color definitions")
	}

	// Must contain widget selectors.
	if !strings.Contains(output, "QWidget {") {
		t.Error("output missing 'QWidget {' rule")
	}
	if !strings.Contains(output, "QPushButton {") {
		t.Error("output missing 'QPushButton {' rule")
	}
	if !strings.Contains(output, "QLineEdit {") {
		t.Error("output missing 'QLineEdit {' rule")
	}

	// Must contain pseudo-state selectors.
	if !strings.Contains(output, "QPushButton:hover {") {
		t.Error("output missing 'QPushButton:hover {' rule")
	}
	if !strings.Contains(output, "QPushButton:pressed {") {
		t.Error("output missing 'QPushButton:pressed {' rule")
	}
	if !strings.Contains(output, "QLineEdit:focus {") {
		t.Error("output missing 'QLineEdit:focus {' rule")
	}

	// Must contain property declarations with literal hex.
	if !strings.Contains(output, "background-color: #1a1b26;") {
		t.Error("output missing 'background-color: #1a1b26;' declaration")
	}
	if !strings.Contains(output, "color: #c0caf5;") {
		t.Error("output missing 'color: #c0caf5;' declaration")
	}
}

// TestQssGenerator_WrongType verifies that passing a non-QssTheme value
// as MappedTheme returns a GenerateError.
func TestQssGenerator_WrongType(t *testing.T) {
	g := generator.NewQss()

	var buf bytes.Buffer

	// Pass a string instead of *ports.QssTheme.
	err := g.Generate(&buf, "not a qss theme")
	if err == nil {
		t.Fatal("Generate() with wrong type should return error, got nil")
	}

	var genErr *domain.GenerateError
	if !errors.As(err, &genErr) {
		t.Errorf("error type = %T, want *domain.GenerateError", err)
	}
}

// TestQssGenerator_Deterministic verifies that generating the same QssTheme
// twice produces byte-identical output.
func TestQssGenerator_Deterministic(t *testing.T) {
	g := generator.NewQss()
	theme := sampleQssTheme()

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

// TestQssGenerator_CommentHeader verifies that the output starts with
// a QSS-specific comment header.
func TestQssGenerator_CommentHeader(t *testing.T) {
	g := generator.NewQss()
	theme := sampleQssTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Output should start with a comment.
	if !strings.HasPrefix(output, "/*") {
		t.Error("output should start with a QSS comment header")
	}
}

// TestQssGenerator_TrailingNewline verifies that the output ends with
// a trailing newline for POSIX compliance.
func TestQssGenerator_TrailingNewline(t *testing.T) {
	g := generator.NewQss()
	theme := sampleQssTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	if !strings.HasSuffix(output, "\n") {
		t.Error("output should end with a trailing newline")
	}
}

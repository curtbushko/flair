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

// sampleCSSTheme returns a CSSTheme with custom properties and element rules.
func sampleCSSTheme() *ports.CSSTheme {
	return &ports.CSSTheme{
		CustomProperties: map[string]string{
			"--flair-bg":             "#1a1b26",
			"--flair-fg":             "#c0caf5",
			"--flair-accent-primary": "#7aa2f7",
		},
		Rules: []ports.CSSRule{
			{
				Selector: "body",
				Properties: []ports.CSSProperty{
					{Property: "background-color", Value: "var(--flair-bg)"},
					{Property: "color", Value: "var(--flair-fg)"},
				},
			},
			{
				Selector: "a",
				Properties: []ports.CSSProperty{
					{Property: "color", Value: "var(--flair-accent-primary)"},
				},
			},
		},
	}
}

// TestCSSGenerator_Interface verifies that the CSS generator implements
// ports.Generator and returns the expected Name() and DefaultFilename().
func TestCSSGenerator_Interface(t *testing.T) {
	g := generator.NewCSS()

	// Compile-time interface check.
	var _ ports.Generator = g

	if name := g.Name(); name != "css" {
		t.Errorf("Name() = %q, want %q", name, "css")
	}

	if filename := g.DefaultFilename(); filename != "style.css" {
		t.Errorf("DefaultFilename() = %q, want %q", filename, "style.css")
	}
}

// TestCSSGenerator_RootBlock verifies that Generate produces a :root block
// with --flair-* custom properties sorted alphabetically.
func TestCSSGenerator_RootBlock(t *testing.T) {
	g := generator.NewCSS()
	theme := &ports.CSSTheme{
		CustomProperties: map[string]string{
			"--flair-fg":             "#c0caf5",
			"--flair-bg":             "#1a1b26",
			"--flair-accent-primary": "#7aa2f7",
		},
		Rules: []ports.CSSRule{},
	}

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Must contain :root block.
	if !strings.Contains(output, ":root {") {
		t.Fatal("output does not contain ':root {' block")
	}

	// Must contain all custom properties.
	for prop, val := range theme.CustomProperties {
		decl := prop + ": " + val + ";"
		if !strings.Contains(output, decl) {
			t.Errorf("output missing declaration %q", decl)
		}
	}

	// Custom properties must be sorted alphabetically.
	rootStart := strings.Index(output, ":root {")
	rootEnd := strings.Index(output[rootStart:], "}")
	rootBlock := output[rootStart : rootStart+rootEnd+1]
	lines := strings.Split(rootBlock, "\n")

	var props []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--flair-") {
			colonIdx := strings.Index(trimmed, ":")
			if colonIdx > 0 {
				props = append(props, trimmed[:colonIdx])
			}
		}
	}

	if len(props) != 3 {
		t.Fatalf("expected 3 custom properties in :root block, got %d", len(props))
	}

	for i := 1; i < len(props); i++ {
		if props[i-1] > props[i] {
			t.Errorf("custom properties not sorted: %q before %q", props[i-1], props[i])
		}
	}
}

// TestCSSGenerator_ElementRules verifies that Generate produces element
// selector rules with proper CSS formatting after the :root block.
func TestCSSGenerator_ElementRules(t *testing.T) {
	g := generator.NewCSS()
	theme := sampleCSSTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Must contain body rule.
	if !strings.Contains(output, "body {") {
		t.Error("output does not contain 'body {' rule")
	}

	// Must contain a rule.
	if !strings.Contains(output, "a {") {
		t.Error("output does not contain 'a {' rule")
	}

	// Body rule must contain property declarations.
	if !strings.Contains(output, "background-color: var(--flair-bg);") {
		t.Error("output missing 'background-color: var(--flair-bg);' declaration")
	}
	if !strings.Contains(output, "color: var(--flair-fg);") {
		t.Error("output missing 'color: var(--flair-fg);' in body rule")
	}

	// Anchor rule must contain property declaration.
	if !strings.Contains(output, "color: var(--flair-accent-primary);") {
		t.Error("output missing 'color: var(--flair-accent-primary);' in a rule")
	}

	// Element rules must come after :root block.
	rootIdx := strings.Index(output, ":root {")
	bodyIdx := strings.Index(output, "body {")
	if rootIdx < 0 || bodyIdx < 0 {
		t.Fatal("output missing :root or body block")
	}
	if bodyIdx <= rootIdx {
		t.Error("body rule should appear after :root block")
	}
}

// TestCSSGenerator_WrongType verifies that passing a non-CSSTheme value
// as MappedTheme returns a GenerateError.
func TestCSSGenerator_WrongType(t *testing.T) {
	g := generator.NewCSS()

	var buf bytes.Buffer

	// Pass a string instead of *ports.CSSTheme.
	err := g.Generate(&buf, "not a css theme")
	if err == nil {
		t.Fatal("Generate() with wrong type should return error, got nil")
	}

	var genErr *domain.GenerateError
	if !errors.As(err, &genErr) {
		t.Errorf("error type = %T, want *domain.GenerateError", err)
	}
}

// TestCSSGenerator_Deterministic verifies that generating the same CSSTheme
// twice produces byte-identical output.
func TestCSSGenerator_Deterministic(t *testing.T) {
	g := generator.NewCSS()
	theme := sampleCSSTheme()

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

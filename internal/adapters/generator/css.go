package generator

import (
	"fmt"
	"io"
	"sort"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// CSS implements ports.Generator for the CSS target.
// It writes a style.css file with a :root{} block of custom properties
// followed by element selector rules from a CSSTheme.
type CSS struct{}

// NewCSS returns a new CSS generator.
func NewCSS() *CSS {
	return &CSS{}
}

// Name returns the target name for this generator.
func (c *CSS) Name() string {
	return "css"
}

// DefaultFilename returns the default output filename for CSS.
func (c *CSS) DefaultFilename() string {
	return "style.css"
}

// Generate writes the CSSTheme as a CSS stylesheet to w. The mapped argument
// must be a *ports.CSSTheme; a type assertion failure returns a
// *domain.GenerateError. Output consists of a :root{} block with sorted
// custom properties followed by element selector rules.
func (c *CSS) Generate(w io.Writer, mapped ports.MappedTheme) error {
	theme, ok := mapped.(*ports.CSSTheme)
	if !ok {
		return &domain.GenerateError{
			Target:  "css",
			Message: fmt.Sprintf("expected *ports.CSSTheme, got %T", mapped),
		}
	}

	if err := writeRootBlock(w, theme.CustomProperties); err != nil {
		return &domain.GenerateError{
			Target:  "css",
			Message: "failed to write :root block",
			Cause:   err,
		}
	}

	if err := writeElementRules(w, theme.Rules); err != nil {
		return &domain.GenerateError{
			Target:  "css",
			Message: "failed to write element rules",
			Cause:   err,
		}
	}

	return nil
}

// writeRootBlock writes a :root{} block with custom properties sorted
// alphabetically by property name.
func writeRootBlock(w io.Writer, props map[string]string) error {
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if _, err := fmt.Fprint(w, ":root {\n"); err != nil {
		return err
	}

	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "  %s: %s;\n", k, props[k]); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "}\n"); err != nil {
		return err
	}

	return nil
}

// writeElementRules writes CSS element rules with selector { property: value; }
// formatting. A blank line separates each rule from the previous block.
func writeElementRules(w io.Writer, rules []ports.CSSRule) error {
	for _, rule := range rules {
		// Blank line before each rule to separate from previous block.
		if _, err := fmt.Fprintf(w, "\n%s {\n", rule.Selector); err != nil {
			return err
		}

		for _, prop := range rule.Properties {
			if _, err := fmt.Fprintf(w, "  %s: %s;\n", prop.Property, prop.Value); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprint(w, "}\n"); err != nil {
			return err
		}
	}

	return nil
}

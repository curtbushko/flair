package mapper_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/ports"
)

// TestCSSMapper_Interface verifies that the CSS mapper implements
// ports.Mapper and Name() returns "css".
func TestCSSMapper_Interface(t *testing.T) {
	m := mapper.NewCSS()

	// Compile-time interface check.
	var _ ports.Mapper = m

	if name := m.Name(); name != "css" {
		t.Errorf("Name() = %q, want %q", name, "css")
	}
}

// TestCSSMapper_CustomProperties verifies that the CSS mapper produces
// custom properties with --flair-* naming convention from semantic tokens.
func TestCSSMapper_CustomProperties(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewCSS()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	ct, ok := result.(*ports.CSSTheme)
	if !ok {
		t.Fatalf("Map() returned %T, want *ports.CSSTheme", result)
	}

	if ct.CustomProperties == nil {
		t.Fatal("CSSTheme.CustomProperties is nil")
	}

	// All custom property names must follow --flair-* naming convention.
	for name := range ct.CustomProperties {
		if !strings.HasPrefix(name, "--flair-") {
			t.Errorf("custom property %q does not follow --flair-* naming convention", name)
		}
	}

	// Should have a reasonable number of custom properties
	// (surface, text, syntax, status, accent, etc.)
	if len(ct.CustomProperties) < 20 {
		t.Errorf("CSSTheme.CustomProperties has %d entries, want >= 20", len(ct.CustomProperties))
	}

	// Verify specific expected custom properties and their values.
	expectedProps := []struct {
		name    string
		wantHex string
	}{
		{"--flair-bg", "#1a1b26"},
		{"--flair-fg", "#c0caf5"},
		{"--flair-bg-raised", "#1f2335"},
		{"--flair-bg-highlight", "#292e42"},
		{"--flair-text-muted", "#565f89"},
		{"--flair-text-secondary", "#a9b1d6"},
		{"--flair-accent-primary", "#7aa2f7"},
		{"--flair-accent-secondary", "#bb9af7"},
		{"--flair-status-error", "#ff899d"},
		{"--flair-status-warning", "#e9c582"},
		{"--flair-status-success", "#afd67a"},
		{"--flair-syntax-keyword", "#bb9af7"},
		{"--flair-syntax-string", "#9ece6a"},
		{"--flair-syntax-function", "#7aa2f7"},
		{"--flair-syntax-comment", "#565f89"},
	}

	for _, tc := range expectedProps {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ct.CustomProperties[tc.name]
			if !ok {
				keys := make([]string, 0, len(ct.CustomProperties))
				for k := range ct.CustomProperties {
					keys = append(keys, k)
				}
				t.Fatalf("custom property %q not found. Available: %v", tc.name, keys)
			}
			gotColor := mustParseHex(t, got)
			wantColor := mustParseHex(t, tc.wantHex)
			if !gotColor.Equal(wantColor) {
				t.Errorf("property %q = %q, want %q", tc.name, got, tc.wantHex)
			}
		})
	}

	// All values must be valid hex colors.
	for name, val := range ct.CustomProperties {
		if !isValidHex(val) {
			t.Errorf("custom property %q has invalid hex value %q", name, val)
		}
	}
}

// TestCSSMapper_ElementRules verifies that the CSS mapper produces element
// rules for body, a, ::selection, and pre/code with correct properties
// referencing custom properties.
func TestCSSMapper_ElementRules(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewCSS()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	ct := result.(*ports.CSSTheme)

	if len(ct.Rules) == 0 {
		t.Fatal("CSSTheme.Rules is empty")
	}

	// Collect rules by selector for easier assertion.
	rulesBySelector := make(map[string][]ports.CSSProperty)
	for _, rule := range ct.Rules {
		rulesBySelector[rule.Selector] = rule.Properties
	}

	// body rule should exist with background-color and color referencing
	// custom properties.
	requiredSelectors := []string{"body", "a", "::selection"}
	for _, sel := range requiredSelectors {
		if _, ok := rulesBySelector[sel]; !ok {
			selectors := make([]string, 0, len(rulesBySelector))
			for s := range rulesBySelector {
				selectors = append(selectors, s)
			}
			t.Errorf("missing rule for selector %q. Available selectors: %v", sel, selectors)
		}
	}

	// Check for pre/code rule (could be "pre, code" or separate "pre" and "code")
	hasPreCode := false
	for sel := range rulesBySelector {
		if strings.Contains(sel, "pre") || strings.Contains(sel, "code") {
			hasPreCode = true
			break
		}
	}
	if !hasPreCode {
		t.Error("missing rule for pre/code elements")
	}

	// body rule should reference custom properties via var() syntax.
	assertRuleHasVarProperty(t, rulesBySelector, "body", "background-color")
	assertRuleHasVarProperty(t, rulesBySelector, "body", "color")

	// a rule should have a color property.
	assertRuleHasProperty(t, rulesBySelector, "a", "color")

	// ::selection rule should have background-color.
	assertRuleHasProperty(t, rulesBySelector, "::selection", "background-color")
}

// TestCSSMapper_NilTheme verifies that the CSS mapper returns an error
// when given a nil theme.
func TestCSSMapper_NilTheme(t *testing.T) {
	m := mapper.NewCSS()
	_, err := m.Map(nil)
	if err == nil {
		t.Fatal("expected error for nil theme, got nil")
	}
}

// assertRuleHasProperty checks that the given selector has a rule with the
// specified property name.
func assertRuleHasProperty(t *testing.T, rules map[string][]ports.CSSProperty, selector, property string) {
	t.Helper()
	props, ok := rules[selector]
	if !ok {
		return // selector presence is checked separately
	}
	for _, prop := range props {
		if prop.Property == property {
			return
		}
	}
	t.Errorf("%s rule missing %s property", selector, property)
}

// assertRuleHasVarProperty checks that the given selector has a rule with
// the specified property name whose value references a CSS custom property
// via var(--flair-...).
func assertRuleHasVarProperty(t *testing.T, rules map[string][]ports.CSSProperty, selector, property string) {
	t.Helper()
	props, ok := rules[selector]
	if !ok {
		return // selector presence is checked separately
	}
	for _, prop := range props {
		if prop.Property == property {
			if !strings.Contains(prop.Value, "var(--flair-") {
				t.Errorf("%s %s should reference a custom property via var(), got %q", selector, property, prop.Value)
			}
			return
		}
	}
	t.Errorf("%s rule missing %s property", selector, property)
}

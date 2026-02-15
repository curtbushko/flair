package mapper_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/ports"
)

// TestQssMapper_Interface verifies that the QSS mapper implements
// ports.Mapper and Name() returns "qss".
func TestQssMapper_Interface(t *testing.T) {
	m := mapper.NewQss()

	// Compile-time interface check.
	var _ ports.Mapper = m

	if name := m.Name(); name != "qss" {
		t.Errorf("Name() = %q, want %q", name, "qss")
	}
}

// TestQssMapper_WidgetRules verifies that the QSS mapper produces widget
// rules for QWidget, QMainWindow, QPushButton, QLineEdit, QTextEdit,
// QScrollBar, etc. with literal hex color values.
func TestQssMapper_WidgetRules(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewQss()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	qt, ok := result.(*ports.QssTheme)
	if !ok {
		t.Fatalf("Map() returned %T, want *ports.QssTheme", result)
	}

	if qt.Rules == nil {
		t.Fatal("QssTheme.Rules is nil")
	}

	// QSS themes should have a reasonable number of widget rules.
	if len(qt.Rules) < 5 {
		t.Errorf("QssTheme.Rules has %d entries, want >= 5", len(qt.Rules))
	}

	// Collect selectors for assertion.
	selectorSet := make(map[string]bool)
	for _, rule := range qt.Rules {
		selectorSet[rule.Selector] = true
	}

	// Required Qt widget selectors.
	requiredSelectors := []string{
		"QWidget",
		"QMainWindow",
		"QPushButton",
		"QLineEdit",
		"QTextEdit",
		"QScrollBar",
	}

	for _, sel := range requiredSelectors {
		if !selectorSet[sel] {
			selectors := make([]string, 0, len(selectorSet))
			for s := range selectorSet {
				selectors = append(selectors, s)
			}
			t.Errorf("missing rule for selector %q. Available selectors: %v", sel, selectors)
		}
	}

	// Each rule must have at least one property.
	for _, rule := range qt.Rules {
		if len(rule.Properties) == 0 {
			t.Errorf("rule for selector %q has no properties", rule.Selector)
		}
	}

	// All property values that look like colors must be valid literal hex.
	for _, rule := range qt.Rules {
		for _, prop := range rule.Properties {
			if len(prop.Value) > 0 && prop.Value[0] == '#' {
				if !isValidHex(prop.Value) {
					t.Errorf("rule %q property %q has invalid hex value %q",
						rule.Selector, prop.Property, prop.Value)
				}
			}
		}
	}
}

// TestQssMapper_PseudoStates verifies that the QSS mapper produces
// pseudo-state rules for QPushButton:hover, QPushButton:pressed,
// and QLineEdit:focus.
func TestQssMapper_PseudoStates(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewQss()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	qt := result.(*ports.QssTheme)

	// Collect selectors for assertion.
	selectorSet := make(map[string]bool)
	for _, rule := range qt.Rules {
		selectorSet[rule.Selector] = true
	}

	// Required pseudo-state selectors.
	requiredPseudoStates := []string{
		"QPushButton:hover",
		"QPushButton:pressed",
		"QLineEdit:focus",
	}

	for _, sel := range requiredPseudoStates {
		if !selectorSet[sel] {
			selectors := make([]string, 0, len(selectorSet))
			for s := range selectorSet {
				selectors = append(selectors, s)
			}
			t.Errorf("missing pseudo-state rule for %q. Available selectors: %v", sel, selectors)
		}
	}

	// Pseudo-state rules should also use literal hex values.
	for _, rule := range qt.Rules {
		for _, prop := range rule.Properties {
			if len(prop.Value) > 0 && prop.Value[0] == '#' {
				if !isValidHex(prop.Value) {
					t.Errorf("pseudo-state rule %q property %q has invalid hex value %q",
						rule.Selector, prop.Property, prop.Value)
				}
			}
		}
	}
}

// TestQssMapper_LiteralHexOnly verifies that the QSS mapper uses only
// literal hex color values -- no var(), @define-color, or other variable
// references in any property value.
func TestQssMapper_LiteralHexOnly(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewQss()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	qt := result.(*ports.QssTheme)

	for _, rule := range qt.Rules {
		for _, prop := range rule.Properties {
			if len(prop.Value) == 0 {
				continue
			}
			// Must not use CSS var() syntax.
			if len(prop.Value) >= 4 && prop.Value[:4] == "var(" {
				t.Errorf("rule %q property %q uses var() syntax: %q -- must be literal hex",
					rule.Selector, prop.Property, prop.Value)
			}
			// Must not use GTK @name syntax.
			if prop.Value[0] == '@' {
				t.Errorf("rule %q property %q uses @name syntax: %q -- must be literal hex",
					rule.Selector, prop.Property, prop.Value)
			}
		}
	}
}

// TestQssMapper_NilTheme verifies that the QSS mapper returns an error
// when given a nil theme.
func TestQssMapper_NilTheme(t *testing.T) {
	m := mapper.NewQss()
	_, err := m.Map(nil)
	if err == nil {
		t.Fatal("expected error for nil theme, got nil")
	}
}

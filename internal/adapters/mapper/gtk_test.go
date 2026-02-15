package mapper_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/ports"
)

// TestGtkMapper_Interface verifies that the GTK mapper implements
// ports.Mapper and Name() returns "gtk".
func TestGtkMapper_Interface(t *testing.T) {
	m := mapper.NewGtk()

	// Compile-time interface check.
	var _ ports.Mapper = m

	if name := m.Name(); name != "gtk" {
		t.Errorf("Name() = %q, want %q", name, "gtk")
	}
}

// TestGtkMapper_ColorDefinitions verifies that the GTK mapper produces
// @define-color entries with GTK-standard color names from semantic tokens.
func TestGtkMapper_ColorDefinitions(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewGtk()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	gt, ok := result.(*ports.GtkTheme)
	if !ok {
		t.Fatalf("Map() returned %T, want *ports.GtkTheme", result)
	}

	if gt.Colors == nil {
		t.Fatal("GtkTheme.Colors is nil")
	}

	// GTK themes should have a reasonable number of color definitions
	// covering window, header, accent, status, and text colors.
	if len(gt.Colors) < 10 {
		t.Errorf("GtkTheme.Colors has %d entries, want >= 10", len(gt.Colors))
	}

	// Build a lookup map for easier assertions.
	colorMap := make(map[string]string)
	for _, cd := range gt.Colors {
		colorMap[cd.Name] = cd.Value
	}

	// Verify specific expected @define-color entries and their values.
	expectedColors := []struct {
		name    string
		wantHex string
	}{
		{"window_bg_color", "#1a1b26"},
		{"window_fg_color", "#c0caf5"},
		{"view_bg_color", "#16161e"},
		{"view_fg_color", "#c0caf5"},
		{"headerbar_bg_color", "#16161e"},
		{"headerbar_fg_color", "#c0caf5"},
		{"accent_bg_color", "#7aa2f7"},
		{"accent_fg_color", "#1a1b26"},
		{"card_bg_color", "#1f2335"},
		{"card_fg_color", "#c0caf5"},
	}

	for _, tc := range expectedColors {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := colorMap[tc.name]
			if !ok {
				names := make([]string, 0, len(colorMap))
				for n := range colorMap {
					names = append(names, n)
				}
				t.Fatalf("color %q not found. Available: %v", tc.name, names)
			}
			gotColor := mustParseHex(t, got)
			wantColor := mustParseHex(t, tc.wantHex)
			if !gotColor.Equal(wantColor) {
				t.Errorf("color %q = %q, want %q", tc.name, got, tc.wantHex)
			}
		})
	}

	// All values must be valid hex colors.
	for _, cd := range gt.Colors {
		if !isValidHex(cd.Value) {
			t.Errorf("color %q has invalid hex value %q", cd.Name, cd.Value)
		}
	}
}

// TestGtkMapper_WidgetRules verifies that the GTK mapper produces widget
// selector rules for window, headerbar, button, entry, textview, etc.
func TestGtkMapper_WidgetRules(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewGtk()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	gt := result.(*ports.GtkTheme)

	if len(gt.Rules) == 0 {
		t.Fatal("GtkTheme.Rules is empty")
	}

	// Collect selectors for assertion.
	selectorSet := make(map[string]bool)
	for _, rule := range gt.Rules {
		selectorSet[rule.Selector] = true
	}

	// Required widget selectors.
	requiredSelectors := []string{
		"window",
		"headerbar",
		"button",
		"entry",
		"textview",
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
	for _, rule := range gt.Rules {
		if len(rule.Properties) == 0 {
			t.Errorf("rule for selector %q has no properties", rule.Selector)
		}
	}
}

// TestGtkMapper_NilTheme verifies that the GTK mapper returns an error
// when given a nil theme.
func TestGtkMapper_NilTheme(t *testing.T) {
	m := mapper.NewGtk()
	_, err := m.Map(nil)
	if err == nil {
		t.Fatal("expected error for nil theme, got nil")
	}
}

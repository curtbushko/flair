package domain_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestValidateOverridePath_Valid(t *testing.T) {
	// Arrange: Override path 'syntax.keyword' (exists in token inventory)
	path := "syntax.keyword"

	// Act: Call ValidateOverridePath(path)
	err := domain.ValidateOverridePath(path)

	// Assert: Returns nil (no error)
	if err != nil {
		t.Errorf("ValidateOverridePath(%q) = %v, want nil", path, err)
	}
}

func TestValidateOverridePath_Invalid(t *testing.T) {
	// Arrange: Override path 'invalid.nonexistent'
	path := "invalid.nonexistent"

	// Act: Call ValidateOverridePath(path)
	err := domain.ValidateOverridePath(path)

	// Assert: Returns error with unknown path message
	if err == nil {
		t.Errorf("ValidateOverridePath(%q) = nil, want error", path)
	}
}

func TestValidateOverridePath_AllValidPaths(t *testing.T) {
	// Test that all paths in the inventory are valid
	validPaths := []string{
		// Surface tokens
		"surface.background",
		"surface.background.raised",
		"surface.background.sunken",
		"surface.background.darkest",
		"surface.background.highlight",
		"surface.background.selection",
		"surface.background.search",
		"surface.background.overlay",
		"surface.background.popup",
		"surface.background.sidebar",
		"surface.background.statusbar",
		// Text tokens
		"text.primary",
		"text.secondary",
		"text.muted",
		"text.subtle",
		"text.inverse",
		"text.overlay",
		"text.sidebar",
		// Syntax tokens
		"syntax.keyword",
		"syntax.string",
		"syntax.function",
		"syntax.comment",
		"syntax.variable",
		"syntax.constant",
		"syntax.operator",
		"syntax.type",
		"syntax.number",
		"syntax.tag",
		"syntax.property",
		"syntax.parameter",
		"syntax.regexp",
		"syntax.escape",
		"syntax.constructor",
	}

	for _, path := range validPaths {
		t.Run(path, func(t *testing.T) {
			err := domain.ValidateOverridePath(path)
			if err != nil {
				t.Errorf("ValidateOverridePath(%q) = %v, want nil", path, err)
			}
		})
	}
}

func TestValidateOverride_ShadowWarning(t *testing.T) {
	// Arrange: Override for 'syntax.comment' (which has default derivation with italic style)
	path := "syntax.comment"

	// Act: Call ValidateOverrideWithWarnings()
	warnings := domain.ValidateOverrideWithWarnings(path)

	// Assert: Returns warning message about shadowing
	if len(warnings) == 0 {
		t.Errorf("ValidateOverrideWithWarnings(%q) returned no warnings, want at least one", path)
	}
}

func TestTokenInventory_Complete(t *testing.T) {
	// Arrange: Full token inventory
	inventory := domain.TokenInventory()

	// Assert: Inventory contains expected prefixes
	expectedPrefixes := []string{
		"surface.",
		"text.",
		"syntax.",
		"status.",
		"diff.",
		"markup.",
		"comment.",
		"accent.",
		"border.",
		"scrollbar.",
		"state.",
		"git.",
		"terminal.",
		"statusline.",
	}

	// Count paths per prefix
	prefixCounts := make(map[string]int)
	for _, path := range inventory {
		for _, prefix := range expectedPrefixes {
			if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
				prefixCounts[prefix]++
				break
			}
		}
	}

	// Assert each prefix has at least one path
	for _, prefix := range expectedPrefixes {
		if prefixCounts[prefix] == 0 {
			t.Errorf("TokenInventory() missing paths with prefix %q", prefix)
		}
	}
}

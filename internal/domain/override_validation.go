package domain

import (
	"errors"
	"fmt"
)

// tokenInventory is the canonical set of all valid token paths in the system.
// These paths are derived by the tokenizer and represent all valid override targets.
var tokenInventory = map[string]struct{}{
	// Surface tokens (11)
	"surface.background":           {},
	"surface.background.raised":    {},
	"surface.background.sunken":    {},
	"surface.background.darkest":   {},
	"surface.background.highlight": {},
	"surface.background.selection": {},
	"surface.background.search":    {},
	"surface.background.overlay":   {},
	"surface.background.popup":     {},
	"surface.background.sidebar":   {},
	"surface.background.statusbar": {},

	// Text tokens (7)
	"text.primary":   {},
	"text.secondary": {},
	"text.muted":     {},
	"text.subtle":    {},
	"text.inverse":   {},
	"text.overlay":   {},
	"text.sidebar":   {},

	// Status tokens (6)
	"status.error":   {},
	"status.warning": {},
	"status.success": {},
	"status.info":    {},
	"status.hint":    {},
	"status.todo":    {},

	// Diff tokens (9)
	"diff.added.fg":     {},
	"diff.added.bg":     {},
	"diff.added.sign":   {},
	"diff.deleted.fg":   {},
	"diff.deleted.bg":   {},
	"diff.deleted.sign": {},
	"diff.changed.fg":   {},
	"diff.changed.bg":   {},
	"diff.ignored":      {},

	// Syntax tokens (15)
	"syntax.keyword":     {},
	"syntax.string":      {},
	"syntax.function":    {},
	"syntax.comment":     {},
	"syntax.variable":    {},
	"syntax.constant":    {},
	"syntax.operator":    {},
	"syntax.type":        {},
	"syntax.number":      {},
	"syntax.tag":         {},
	"syntax.property":    {},
	"syntax.parameter":   {},
	"syntax.regexp":      {},
	"syntax.escape":      {},
	"syntax.constructor": {},

	// Markup tokens (15)
	"markup.heading":        {},
	"markup.heading.1":      {},
	"markup.heading.2":      {},
	"markup.heading.3":      {},
	"markup.heading.4":      {},
	"markup.heading.5":      {},
	"markup.heading.6":      {},
	"markup.link":           {},
	"markup.code":           {},
	"markup.bold":           {},
	"markup.italic":         {},
	"markup.strikethrough":  {},
	"markup.quote":          {},
	"markup.list.bullet":    {},
	"markup.list.checked":   {},
	"markup.list.unchecked": {},

	// Comment tokens (6)
	"comment.error":   {},
	"comment.warning": {},
	"comment.info":    {},
	"comment.hint":    {},
	"comment.note":    {},
	"comment.todo":    {},

	// Accent tokens (3)
	"accent.primary":    {},
	"accent.secondary":  {},
	"accent.foreground": {},

	// Border tokens (3)
	"border.default": {},
	"border.focus":   {},
	"border.muted":   {},

	// Scrollbar tokens (2)
	"scrollbar.thumb": {},
	"scrollbar.track": {},

	// State tokens (3)
	"state.hover":       {},
	"state.active":      {},
	"state.disabled.fg": {},

	// Git tokens (4)
	"git.added":    {},
	"git.modified": {},
	"git.deleted":  {},
	"git.ignored":  {},

	// Terminal tokens (16)
	"terminal.black":     {},
	"terminal.red":       {},
	"terminal.green":     {},
	"terminal.yellow":    {},
	"terminal.blue":      {},
	"terminal.magenta":   {},
	"terminal.cyan":      {},
	"terminal.white":     {},
	"terminal.brblack":   {},
	"terminal.brred":     {},
	"terminal.brgreen":   {},
	"terminal.bryellow":  {},
	"terminal.brblue":    {},
	"terminal.brmagenta": {},
	"terminal.brcyan":    {},
	"terminal.brwhite":   {},

	// Statusline tokens (6)
	"statusline.a.bg": {},
	"statusline.a.fg": {},
	"statusline.b.bg": {},
	"statusline.b.fg": {},
	"statusline.c.bg": {},
	"statusline.c.fg": {},
}

// styledTokens lists token paths that have default style flags (bold, italic, etc.).
// Overriding these tokens will shadow the default styling, which may be intentional
// but deserves a warning.
var styledTokens = map[string]string{
	"syntax.comment":       "italic",
	"markup.heading":       "bold",
	"markup.heading.1":     "bold",
	"markup.heading.2":     "bold",
	"markup.heading.3":     "bold",
	"markup.heading.4":     "bold",
	"markup.heading.5":     "bold",
	"markup.heading.6":     "bold",
	"markup.bold":          "bold",
	"markup.italic":        "italic",
	"markup.strikethrough": "strikethrough",
	"markup.quote":         "italic",
}

// ErrUnknownOverridePath is returned when an override path is not in the token inventory.
var ErrUnknownOverridePath = errors.New("unknown override path")

// ValidateOverridePath checks if the given path is a valid token path.
// Returns nil if the path exists in the token inventory, or an error if it does not.
func ValidateOverridePath(path string) error {
	if _, ok := tokenInventory[path]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownOverridePath, path)
	}
	return nil
}

// ValidateOverrideWithWarnings checks if overriding the given path would shadow
// a default derived value (e.g., a token with default styling).
// Returns a slice of warning messages (may be empty).
func ValidateOverrideWithWarnings(path string) []string {
	var warnings []string

	// Check if this path has default styling that would be shadowed
	if style, ok := styledTokens[path]; ok {
		warnings = append(warnings, "override shadows default "+style+" style for "+path)
	}

	return warnings
}

// TokenInventory returns a sorted slice of all valid token paths.
// This is the canonical source of truth for what override paths are valid.
func TokenInventory() []string {
	paths := make([]string, 0, len(tokenInventory))
	for path := range tokenInventory {
		paths = append(paths, path)
	}
	return paths
}

// ValidateOverrides validates all override paths in a map against the token inventory.
// Returns a slice of validation error messages for any invalid paths.
func ValidateOverrides(overrides map[string]TokenOverride) []string {
	if overrides == nil {
		return nil
	}

	var violations []string
	for path := range overrides {
		if err := ValidateOverridePath(path); err != nil {
			violations = append(violations, "invalid override path: "+path)
		}
	}
	return violations
}

// ValidateOverridesWithWarnings validates all override paths and collects warnings
// for any that shadow default derived values.
func ValidateOverridesWithWarnings(overrides map[string]TokenOverride) (violations []string, warnings []string) {
	if overrides == nil {
		return nil, nil
	}

	for path := range overrides {
		if err := ValidateOverridePath(path); err != nil {
			violations = append(violations, "invalid override path: "+path)
			continue
		}
		warnings = append(warnings, ValidateOverrideWithWarnings(path)...)
	}
	return violations, warnings
}

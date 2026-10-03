// Package mapper provides adapters that transform a ResolvedTheme into
// target-specific mapped theme structs (ports.StylixTheme, etc.).
package mapper

import (
	"errors"
	"sort"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// Stylix implements ports.Mapper for the Stylix target.
// It maps a ResolvedTheme into a flat key-value StylixTheme containing
// base24 palette passthrough slots and semantic token colors in kebab-case.
type Stylix struct{}

// NewStylix returns a new Stylix mapper.
func NewStylix() *Stylix {
	return &Stylix{}
}

// Name returns the target name for this mapper.
func (s *Stylix) Name() string {
	return "stylix"
}

// Map transforms a ResolvedTheme into a *ports.StylixTheme with at least
// 60 key-value pairs: 24 base palette slots plus semantic token colors.
func (s *Stylix) Map(theme *domain.ResolvedTheme) (ports.MappedTheme, error) {
	if theme == nil {
		return nil, errors.New("stylix mapper: nil theme")
	}
	if theme.Palette == nil {
		return nil, errors.New("stylix mapper: nil palette")
	}
	if theme.Tokens == nil {
		return nil, errors.New("stylix mapper: nil tokens")
	}

	values := make(map[string]string)

	// Pass through all 24 base palette slots.
	mapPaletteSlots(theme.Palette, values)

	// Map semantic token colors.
	mapSemanticTokens(theme.Tokens, values)

	return &ports.StylixTheme{Values: values}, nil
}

// mapPaletteSlots adds base00 through base17 to the values map.
func mapPaletteSlots(pal *domain.Palette, values map[string]string) {
	slotNames := pal.SlotNames()
	for i, name := range slotNames {
		values[name] = pal.Base(i).Hex()
	}
}

// semanticMapping maps a semantic token path to its Stylix key name.
// The key uses kebab-case derived from the dot-separated token path.
type semanticMapping struct {
	tokenPath string
	stylixKey string
}

// semanticMappings defines the semantic token paths and their corresponding
// Stylix key names. This produces the semantic portion of the >=60 key
// requirement (24 base slots + these semantic keys).
//
//nolint:dupl // Each mapper has its own naming convention; structural similarity is intentional.
var semanticMappings = func() []semanticMapping {
	// Build mappings grouped by category, matching the PLAN.md token inventory.
	// Includes all existing color-bearing semantic tokens.
	m := make([]semanticMapping, 0, 115)

	// Surface tokens (11)
	m = append(m,
		semanticMapping{"surface.background", "surface-bg"},
		semanticMapping{"surface.background.raised", "surface-bg-raised"},
		semanticMapping{"surface.background.sunken", "surface-bg-sunken"},
		semanticMapping{"surface.background.darkest", "surface-bg-darkest"},
		semanticMapping{"surface.background.highlight", "surface-bg-highlight"},
		semanticMapping{"surface.background.selection", "surface-bg-selection"},
		semanticMapping{"surface.background.search", "surface-bg-search"},
		semanticMapping{"surface.background.overlay", "surface-bg-overlay"},
		semanticMapping{"surface.background.popup", "surface-bg-popup"},
		semanticMapping{"surface.background.sidebar", "surface-bg-sidebar"},
		semanticMapping{"surface.background.statusbar", "surface-bg-statusbar"},
	)

	// Text tokens (7)
	m = append(m,
		semanticMapping{"text.primary", "text-primary"},
		semanticMapping{"text.secondary", "text-secondary"},
		semanticMapping{"text.muted", "text-muted"},
		semanticMapping{"text.subtle", "text-subtle"},
		semanticMapping{"text.inverse", "text-inverse"},
		semanticMapping{"text.overlay", "text-overlay"},
		semanticMapping{"text.sidebar", "text-sidebar"},
	)

	// Status tokens (6)
	m = append(m,
		semanticMapping{"status.error", "status-error"},
		semanticMapping{"status.warning", "status-warning"},
		semanticMapping{"status.success", "status-success"},
		semanticMapping{"status.info", "status-info"},
		semanticMapping{"status.hint", "status-hint"},
		semanticMapping{"status.todo", "status-todo"},
	)

	// Diff tokens (10)
	m = append(m,
		semanticMapping{"diff.added.fg", "diff-added-fg"},
		semanticMapping{"diff.added.bg", "diff-added-bg"},
		semanticMapping{"diff.added.sign", "diff-added-sign"},
		semanticMapping{"diff.deleted.fg", "diff-deleted-fg"},
		semanticMapping{"diff.deleted.bg", "diff-deleted-bg"},
		semanticMapping{"diff.deleted.sign", "diff-deleted-sign"},
		semanticMapping{"diff.changed.fg", "diff-changed-fg"},
		semanticMapping{"diff.changed.bg", "diff-changed-bg"},
		semanticMapping{"diff.changed.sign", "diff-changed-sign"},
		semanticMapping{"diff.ignored", "diff-ignored"},
	)

	// Syntax tokens (25)
	m = append(m,
		semanticMapping{"syntax.keyword", "syntax-keyword"},
		semanticMapping{"syntax.string", "syntax-string"},
		semanticMapping{"syntax.function", "syntax-function"},
		semanticMapping{"syntax.comment", "syntax-comment"},
		semanticMapping{"syntax.variable", "syntax-variable"},
		semanticMapping{"syntax.constant", "syntax-constant"},
		semanticMapping{"syntax.operator", "syntax-operator"},
		semanticMapping{"syntax.type", "syntax-type"},
		semanticMapping{"syntax.number", "syntax-number"},
		semanticMapping{"syntax.tag", "syntax-tag"},
		semanticMapping{"syntax.property", "syntax-property"},
		semanticMapping{"syntax.parameter", "syntax-parameter"},
		semanticMapping{"syntax.regexp", "syntax-regexp"},
		semanticMapping{"syntax.escape", "syntax-escape"},
		semanticMapping{"syntax.constructor", "syntax-constructor"},
		semanticMapping{"syntax.boolean", "syntax-boolean"},
		semanticMapping{"syntax.function.builtin", "syntax-function-builtin"},
		semanticMapping{"syntax.type.builtin", "syntax-type-builtin"},
		semanticMapping{"syntax.module", "syntax-module"},
		semanticMapping{"syntax.module.builtin", "syntax-module-builtin"},
		semanticMapping{"syntax.string.documentation", "syntax-string-documentation"},
		semanticMapping{"syntax.label", "syntax-label"},
		semanticMapping{"syntax.punctuation", "syntax-punctuation"},
		semanticMapping{"syntax.deprecated", "syntax-deprecated"},
		semanticMapping{"syntax.macro", "syntax-macro"},
	)

	// Markup tokens (13 with color, 3 style-only excluded from Stylix)
	m = append(m,
		semanticMapping{"markup.heading", "markup-heading"},
		semanticMapping{"markup.heading.1", "markup-heading-1"},
		semanticMapping{"markup.heading.2", "markup-heading-2"},
		semanticMapping{"markup.heading.3", "markup-heading-3"},
		semanticMapping{"markup.heading.4", "markup-heading-4"},
		semanticMapping{"markup.heading.5", "markup-heading-5"},
		semanticMapping{"markup.heading.6", "markup-heading-6"},
		semanticMapping{"markup.link", "markup-link"},
		semanticMapping{"markup.code", "markup-code"},
		semanticMapping{"markup.quote", "markup-quote"},
		semanticMapping{"markup.list.bullet", "markup-list-bullet"},
		semanticMapping{"markup.list.checked", "markup-list-checked"},
		semanticMapping{"markup.list.unchecked", "markup-list-unchecked"},
	)

	// Comment annotation tokens (6)
	m = append(m,
		semanticMapping{"comment.error", "comment-error"},
		semanticMapping{"comment.warning", "comment-warning"},
		semanticMapping{"comment.info", "comment-info"},
		semanticMapping{"comment.hint", "comment-hint"},
		semanticMapping{"comment.note", "comment-note"},
		semanticMapping{"comment.todo", "comment-todo"},
	)

	// Accent tokens (3)
	m = append(m,
		semanticMapping{"accent.primary", "accent-primary"},
		semanticMapping{"accent.secondary", "accent-secondary"},
		semanticMapping{"accent.foreground", "accent-foreground"},
	)

	// Border tokens (3)
	m = append(m,
		semanticMapping{"border.default", "border-default"},
		semanticMapping{"border.focus", "border-focus"},
		semanticMapping{"border.muted", "border-muted"},
	)

	// Scrollbar tokens (2)
	m = append(m,
		semanticMapping{"scrollbar.thumb", "scrollbar-thumb"},
		semanticMapping{"scrollbar.track", "scrollbar-track"},
	)

	// State tokens (3)
	m = append(m,
		semanticMapping{"state.hover", "state-hover"},
		semanticMapping{"state.active", "state-active"},
		semanticMapping{"state.disabled.fg", "state-disabled-fg"},
	)

	// Git tokens (4)
	m = append(m,
		semanticMapping{"git.added", "git-added"},
		semanticMapping{"git.modified", "git-modified"},
		semanticMapping{"git.deleted", "git-deleted"},
		semanticMapping{"git.ignored", "git-ignored"},
	)

	// Terminal ANSI tokens (16)
	m = append(m,
		semanticMapping{"terminal.black", "terminal-black"},
		semanticMapping{"terminal.red", "terminal-red"},
		semanticMapping{"terminal.green", "terminal-green"},
		semanticMapping{"terminal.yellow", "terminal-yellow"},
		semanticMapping{"terminal.blue", "terminal-blue"},
		semanticMapping{"terminal.magenta", "terminal-magenta"},
		semanticMapping{"terminal.cyan", "terminal-cyan"},
		semanticMapping{"terminal.white", "terminal-white"},
		semanticMapping{"terminal.brblack", "terminal-brblack"},
		semanticMapping{"terminal.brred", "terminal-brred"},
		semanticMapping{"terminal.brgreen", "terminal-brgreen"},
		semanticMapping{"terminal.bryellow", "terminal-bryellow"},
		semanticMapping{"terminal.brblue", "terminal-brblue"},
		semanticMapping{"terminal.brmagenta", "terminal-brmagenta"},
		semanticMapping{"terminal.brcyan", "terminal-brcyan"},
		semanticMapping{"terminal.brwhite", "terminal-brwhite"},
	)

	// Statusline tokens (6)
	m = append(m,
		semanticMapping{"statusline.a.bg", "statusline-a-bg"},
		semanticMapping{"statusline.a.fg", "statusline-a-fg"},
		semanticMapping{"statusline.b.bg", "statusline-b-bg"},
		semanticMapping{"statusline.b.fg", "statusline-b-fg"},
		semanticMapping{"statusline.c.bg", "statusline-c-bg"},
		semanticMapping{"statusline.c.fg", "statusline-c-fg"},
	)

	return m
}()

// mapSemanticTokens adds semantic token colors to the values map.
// Style-only tokens (NoneColor) are skipped since Stylix only uses colors.
func mapSemanticTokens(ts *domain.TokenSet, values map[string]string) {
	for _, sm := range semanticMappings {
		tok, ok := ts.Get(sm.tokenPath)
		if !ok {
			continue
		}
		// Skip style-only tokens (markup.bold, markup.italic, etc.)
		if tok.Color.IsNone {
			continue
		}
		values[sm.stylixKey] = tok.Color.Hex()
	}
}

// SortedKeys returns the keys of a StylixTheme.Values map in sorted order.
// Used by the file I/O layer for deterministic output.
func SortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

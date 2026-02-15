package mapper

import (
	"errors"
	"fmt"
	"sort"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// CSS implements ports.Mapper for the CSS target.
// It maps a ResolvedTheme into a CSSTheme containing --flair-* custom
// properties derived from semantic tokens and element selector rules
// for body, a, ::selection, and pre/code.
type CSS struct{}

// NewCSS returns a new CSS mapper.
func NewCSS() *CSS {
	return &CSS{}
}

// Name returns the target name for this mapper.
func (c *CSS) Name() string {
	return "css"
}

// Map transforms a ResolvedTheme into a *ports.CSSTheme with --flair-*
// custom properties from semantic tokens and element rules for
// body, a, ::selection, and pre/code.
func (c *CSS) Map(theme *domain.ResolvedTheme) (ports.MappedTheme, error) {
	if theme == nil {
		return nil, errors.New("css mapper: nil theme")
	}
	if theme.Palette == nil {
		return nil, errors.New("css mapper: nil palette")
	}
	if theme.Tokens == nil {
		return nil, errors.New("css mapper: nil tokens")
	}

	props := make(map[string]string)
	mapCSSCustomProperties(theme.Tokens, props)

	rules := buildCSSElementRules()

	return &ports.CSSTheme{
		CustomProperties: props,
		Rules:            rules,
	}, nil
}

// cssPropertyMapping maps a semantic token path to its CSS custom property name.
type cssPropertyMapping struct {
	tokenPath    string
	propertyName string
}

// cssPropertyMappings defines the semantic token paths and their corresponding
// CSS custom property names following the --flair-* naming convention.
//
//nolint:dupl // Each mapper has its own naming convention; structural similarity is intentional.
var cssPropertyMappings = func() []cssPropertyMapping {
	m := make([]cssPropertyMapping, 0, 80)

	// Surface tokens
	m = append(m,
		cssPropertyMapping{"surface.background", "--flair-bg"},
		cssPropertyMapping{"surface.background.raised", "--flair-bg-raised"},
		cssPropertyMapping{"surface.background.sunken", "--flair-bg-sunken"},
		cssPropertyMapping{"surface.background.darkest", "--flair-bg-darkest"},
		cssPropertyMapping{"surface.background.highlight", "--flair-bg-highlight"},
		cssPropertyMapping{"surface.background.selection", "--flair-bg-selection"},
		cssPropertyMapping{"surface.background.search", "--flair-bg-search"},
		cssPropertyMapping{"surface.background.overlay", "--flair-bg-overlay"},
		cssPropertyMapping{"surface.background.popup", "--flair-bg-popup"},
		cssPropertyMapping{"surface.background.sidebar", "--flair-bg-sidebar"},
		cssPropertyMapping{"surface.background.statusbar", "--flair-bg-statusbar"},
	)

	// Text tokens
	m = append(m,
		cssPropertyMapping{"text.primary", "--flair-fg"},
		cssPropertyMapping{"text.secondary", "--flair-text-secondary"},
		cssPropertyMapping{"text.muted", "--flair-text-muted"},
		cssPropertyMapping{"text.subtle", "--flair-text-subtle"},
		cssPropertyMapping{"text.inverse", "--flair-text-inverse"},
		cssPropertyMapping{"text.overlay", "--flair-text-overlay"},
		cssPropertyMapping{"text.sidebar", "--flair-text-sidebar"},
	)

	// Status tokens
	m = append(m,
		cssPropertyMapping{"status.error", "--flair-status-error"},
		cssPropertyMapping{"status.warning", "--flair-status-warning"},
		cssPropertyMapping{"status.success", "--flair-status-success"},
		cssPropertyMapping{"status.info", "--flair-status-info"},
		cssPropertyMapping{"status.hint", "--flair-status-hint"},
		cssPropertyMapping{"status.todo", "--flair-status-todo"},
	)

	// Diff tokens
	m = append(m,
		cssPropertyMapping{"diff.added.fg", "--flair-diff-added-fg"},
		cssPropertyMapping{"diff.added.bg", "--flair-diff-added-bg"},
		cssPropertyMapping{"diff.added.sign", "--flair-diff-added-sign"},
		cssPropertyMapping{"diff.deleted.fg", "--flair-diff-deleted-fg"},
		cssPropertyMapping{"diff.deleted.bg", "--flair-diff-deleted-bg"},
		cssPropertyMapping{"diff.deleted.sign", "--flair-diff-deleted-sign"},
		cssPropertyMapping{"diff.changed.fg", "--flair-diff-changed-fg"},
		cssPropertyMapping{"diff.changed.bg", "--flair-diff-changed-bg"},
		cssPropertyMapping{"diff.ignored", "--flair-diff-ignored"},
	)

	// Syntax tokens
	m = append(m,
		cssPropertyMapping{"syntax.keyword", "--flair-syntax-keyword"},
		cssPropertyMapping{"syntax.string", "--flair-syntax-string"},
		cssPropertyMapping{"syntax.function", "--flair-syntax-function"},
		cssPropertyMapping{"syntax.comment", "--flair-syntax-comment"},
		cssPropertyMapping{"syntax.variable", "--flair-syntax-variable"},
		cssPropertyMapping{"syntax.constant", "--flair-syntax-constant"},
		cssPropertyMapping{"syntax.operator", "--flair-syntax-operator"},
		cssPropertyMapping{"syntax.type", "--flair-syntax-type"},
		cssPropertyMapping{"syntax.number", "--flair-syntax-number"},
		cssPropertyMapping{"syntax.tag", "--flair-syntax-tag"},
		cssPropertyMapping{"syntax.property", "--flair-syntax-property"},
		cssPropertyMapping{"syntax.parameter", "--flair-syntax-parameter"},
		cssPropertyMapping{"syntax.regexp", "--flair-syntax-regexp"},
		cssPropertyMapping{"syntax.escape", "--flair-syntax-escape"},
		cssPropertyMapping{"syntax.constructor", "--flair-syntax-constructor"},
	)

	// Markup tokens (only those with colors)
	m = append(m,
		cssPropertyMapping{"markup.heading", "--flair-markup-heading"},
		cssPropertyMapping{"markup.link", "--flair-markup-link"},
		cssPropertyMapping{"markup.code", "--flair-markup-code"},
		cssPropertyMapping{"markup.quote", "--flair-markup-quote"},
		cssPropertyMapping{"markup.list.bullet", "--flair-markup-list-bullet"},
		cssPropertyMapping{"markup.list.checked", "--flair-markup-list-checked"},
		cssPropertyMapping{"markup.list.unchecked", "--flair-markup-list-unchecked"},
	)

	// Accent tokens
	m = append(m,
		cssPropertyMapping{"accent.primary", "--flair-accent-primary"},
		cssPropertyMapping{"accent.secondary", "--flair-accent-secondary"},
		cssPropertyMapping{"accent.foreground", "--flair-accent-foreground"},
	)

	// Border tokens
	m = append(m,
		cssPropertyMapping{"border.default", "--flair-border-default"},
		cssPropertyMapping{"border.focus", "--flair-border-focus"},
		cssPropertyMapping{"border.muted", "--flair-border-muted"},
	)

	// Scrollbar tokens
	m = append(m,
		cssPropertyMapping{"scrollbar.thumb", "--flair-scrollbar-thumb"},
		cssPropertyMapping{"scrollbar.track", "--flair-scrollbar-track"},
	)

	// State tokens
	m = append(m,
		cssPropertyMapping{"state.hover", "--flair-state-hover"},
		cssPropertyMapping{"state.active", "--flair-state-active"},
		cssPropertyMapping{"state.disabled.fg", "--flair-state-disabled-fg"},
	)

	// Git tokens
	m = append(m,
		cssPropertyMapping{"git.added", "--flair-git-added"},
		cssPropertyMapping{"git.modified", "--flair-git-modified"},
		cssPropertyMapping{"git.deleted", "--flair-git-deleted"},
		cssPropertyMapping{"git.ignored", "--flair-git-ignored"},
	)

	// Terminal ANSI tokens
	m = append(m,
		cssPropertyMapping{"terminal.black", "--flair-terminal-black"},
		cssPropertyMapping{"terminal.red", "--flair-terminal-red"},
		cssPropertyMapping{"terminal.green", "--flair-terminal-green"},
		cssPropertyMapping{"terminal.yellow", "--flair-terminal-yellow"},
		cssPropertyMapping{"terminal.blue", "--flair-terminal-blue"},
		cssPropertyMapping{"terminal.magenta", "--flair-terminal-magenta"},
		cssPropertyMapping{"terminal.cyan", "--flair-terminal-cyan"},
		cssPropertyMapping{"terminal.white", "--flair-terminal-white"},
		cssPropertyMapping{"terminal.brblack", "--flair-terminal-brblack"},
		cssPropertyMapping{"terminal.brred", "--flair-terminal-brred"},
		cssPropertyMapping{"terminal.brgreen", "--flair-terminal-brgreen"},
		cssPropertyMapping{"terminal.bryellow", "--flair-terminal-bryellow"},
		cssPropertyMapping{"terminal.brblue", "--flair-terminal-brblue"},
		cssPropertyMapping{"terminal.brmagenta", "--flair-terminal-brmagenta"},
		cssPropertyMapping{"terminal.brcyan", "--flair-terminal-brcyan"},
		cssPropertyMapping{"terminal.brwhite", "--flair-terminal-brwhite"},
	)

	return m
}()

// mapCSSCustomProperties populates custom properties from semantic tokens.
// Style-only tokens (NoneColor) are skipped since CSS custom properties
// represent color values.
func mapCSSCustomProperties(ts *domain.TokenSet, props map[string]string) {
	for _, pm := range cssPropertyMappings {
		tok, ok := ts.Get(pm.tokenPath)
		if !ok {
			continue
		}
		// Skip style-only tokens (markup.bold, markup.italic, etc.)
		if tok.Color.IsNone {
			continue
		}
		props[pm.propertyName] = tok.Color.Hex()
	}
}

// buildCSSElementRules creates the standard CSS element rules that reference
// custom properties via var() syntax.
func buildCSSElementRules() []ports.CSSRule {
	return []ports.CSSRule{
		{
			Selector: "body",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: varRef("--flair-bg")},
				{Property: "color", Value: varRef("--flair-fg")},
				{Property: "font-family", Value: "system-ui, -apple-system, sans-serif"},
			},
		},
		{
			Selector: "a",
			Properties: []ports.CSSProperty{
				{Property: "color", Value: varRef("--flair-accent-primary")},
			},
		},
		{
			Selector: "a:hover",
			Properties: []ports.CSSProperty{
				{Property: "color", Value: varRef("--flair-accent-secondary")},
			},
		},
		{
			Selector: "::selection",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: varRef("--flair-bg-selection")},
				{Property: "color", Value: varRef("--flair-fg")},
			},
		},
		{
			Selector: "pre, code",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: varRef("--flair-bg-raised")},
				{Property: "color", Value: varRef("--flair-fg")},
				{Property: "border-radius", Value: "4px"},
			},
		},
	}
}

// varRef returns a CSS var() reference for a custom property name.
func varRef(name string) string {
	return fmt.Sprintf("var(%s)", name)
}

// SortedCSSPropertyNames returns the keys of a CSSTheme.CustomProperties map
// in sorted order. Used by the file I/O layer for deterministic output.
func SortedCSSPropertyNames(props map[string]string) []string {
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

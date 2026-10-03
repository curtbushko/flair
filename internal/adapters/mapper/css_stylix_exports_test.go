package mapper_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// TestCSSStylixExportsInventory checks public aliases against independently
// overridden tokens, including legacy aliases and the extended color inventory.
func TestCSSStylixExportsInventory(t *testing.T) {
	paths := strings.Fields(`
		surface.background surface.background.raised surface.background.sunken
		surface.background.darkest surface.background.highlight surface.background.selection
		surface.background.search surface.background.overlay surface.background.popup
		surface.background.sidebar surface.background.statusbar
		text.primary text.secondary text.muted text.subtle text.inverse text.overlay text.sidebar
		status.error status.warning status.success status.info status.hint status.todo
		diff.added.fg diff.added.bg diff.added.sign diff.deleted.fg diff.deleted.bg diff.deleted.sign
		diff.changed.fg diff.changed.bg diff.changed.sign diff.ignored
		syntax.keyword syntax.string syntax.function syntax.comment syntax.variable syntax.constant
		syntax.operator syntax.type syntax.number syntax.tag syntax.property syntax.parameter
		syntax.regexp syntax.escape syntax.constructor syntax.boolean syntax.function.builtin
		syntax.type.builtin syntax.module syntax.module.builtin syntax.string.documentation
		syntax.label syntax.punctuation syntax.deprecated syntax.macro
		markup.heading markup.heading.1 markup.heading.2 markup.heading.3 markup.heading.4
		markup.heading.5 markup.heading.6 markup.link markup.code markup.quote
		markup.list.bullet markup.list.checked markup.list.unchecked
		comment.error comment.warning comment.info comment.hint comment.note comment.todo
		accent.primary accent.secondary accent.foreground border.default border.focus border.muted
		scrollbar.thumb scrollbar.track state.hover state.active state.disabled.fg
		git.added git.modified git.deleted git.ignored
		terminal.black terminal.red terminal.green terminal.yellow terminal.blue terminal.magenta
		terminal.cyan terminal.white terminal.brblack terminal.brred terminal.brgreen terminal.bryellow
		terminal.brblue terminal.brmagenta terminal.brcyan terminal.brwhite
		statusline.a.bg statusline.a.fg statusline.b.bg statusline.b.fg statusline.c.bg statusline.c.fg
	`)
	theme := buildResolvedTheme(t)
	for i, path := range paths {
		tok, ok := theme.Tokens.Get(path)
		require.True(t, ok, "existing token %s", path)
		// Unique colors distinguish even tokens that normally share a palette slot.
		tok.Color = mustParseHex(t, fmt.Sprintf("#%06x", 0x010000+i))
		theme.Tokens.Set(path, tok)
	}

	cssResult, err := mapper.NewCSS().Map(theme)
	require.NoError(t, err)
	css, ok := cssResult.(*ports.CSSTheme)
	require.True(t, ok)
	stylixResult, err := mapper.NewStylix().Map(theme)
	require.NoError(t, err)
	stylix, ok := stylixResult.(*ports.StylixTheme)
	require.True(t, ok)

	for i, path := range paths {
		t.Run(path, func(t *testing.T) {
			alias := strings.ReplaceAll(path, ".", "-")
			cssAlias := alias
			if strings.HasPrefix(path, "surface.background") {
				suffix := strings.TrimPrefix(alias, "surface-background")
				alias = "surface-bg" + suffix
				cssAlias = "bg" + suffix
			} else if path == "text.primary" {
				cssAlias = "fg"
			}
			want := fmt.Sprintf("#%06x", 0x010000+i)
			assert.Equal(t, want, css.CustomProperties["--flair-"+cssAlias])
			assert.Equal(t, want, stylix.Values[alias])
		})
	}
	assert.Len(t, css.CustomProperties, len(paths))
	assert.Len(t, stylix.Values, len(paths)+24)
}

func TestCSSStylixExportsLinkRule(t *testing.T) {
	theme := buildResolvedTheme(t)
	theme.Tokens.Set("markup.link", domain.Token{Color: mustParseHex(t, "#123456")})
	theme.Tokens.Set("accent.primary", domain.Token{Color: mustParseHex(t, "#654321")})
	result, err := mapper.NewCSS().Map(theme)
	require.NoError(t, err)
	css, ok := result.(*ports.CSSTheme)
	require.True(t, ok)
	assert.Equal(t, "#123456", css.CustomProperties["--flair-markup-link"])
	rules := make(map[string][]ports.CSSProperty)
	for _, rule := range css.Rules {
		rules[rule.Selector] = rule.Properties
	}
	assert.Equal(t, []ports.CSSProperty{{Property: "color", Value: "var(--flair-markup-link)"}}, rules["a"])
	assert.Equal(t, []ports.CSSProperty{{Property: "color", Value: "var(--flair-accent-secondary)"}}, rules["a:hover"])
}

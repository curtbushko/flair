package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// Distinct sentinels prevent coincident palette colors from masking token mixups.
func uiTheme(t *testing.T) *domain.ResolvedTheme {
	t.Helper()
	theme := buildResolvedTheme(t)
	for path, hex := range map[string]string{
		"surface.background.selection": "#123456",
		"surface.background.highlight": "#234567",
		"surface.background.statusbar": "#345678",
		"surface.background.sunken":    "#456789",
		"text.primary":                 "#56789a",
		"text.secondary":               "#6789ab",
		"state.hover":                  "#789abc",
		"state.active":                 "#89abcd",
		"state.disabled.fg":            "#9abcde",
		"border.focus":                 "#abcdef",
		"accent.primary":               "#bcdef0",
		"accent.foreground":            "#cdef01",
	} {
		theme.Tokens.Set(path, domain.Token{Color: mustParseHex(t, hex)})
	}
	return theme
}

func uiProperty(t *testing.T, rules []ports.CSSRule, selector, property string) string {
	t.Helper()
	for _, rule := range rules {
		if rule.Selector != selector {
			continue
		}
		for _, prop := range rule.Properties {
			if prop.Property == property {
				return prop.Value
			}
		}
	}
	require.FailNow(t, "missing UI property", "%s: %s", selector, property)
	return ""
}

func TestGtkMapper_SemanticStates(t *testing.T) {
	mapped, err := mapper.NewGtk().Map(uiTheme(t))
	require.NoError(t, err)
	theme, ok := mapped.(*ports.GtkTheme)
	require.True(t, ok)
	colors := make(map[string]string)
	for _, color := range theme.Colors {
		colors[color.Name] = color.Value
	}
	for name, want := range map[string]string{
		"hover_bg_color": "#789abc", "active_bg_color": "#89abcd",
		"disabled_fg_color": "#9abcde", "focus_border_color": "#abcdef",
		"selection_bg_color": "#123456", "selection_fg_color": "#56789a",
		"error_fg_color": "#56789a", "warning_fg_color": "#56789a", "success_fg_color": "#56789a",
	} {
		assert.Equal(t, want, colors[name], name)
	}
	for _, rule := range theme.Rules {
		for _, prop := range rule.Properties {
			if len(prop.Value) > 0 && prop.Value[0] == '@' {
				assert.Contains(t, colors, prop.Value[1:], rule.Selector)
			}
		}
	}
	for _, tc := range []struct{ selector, property, want string }{
		{"button", "color", "@card_fg_color"},
		{"button:hover", "background-color", "@hover_bg_color"},
		{"button:active", "background-color", "@active_bg_color"},
		{"button:disabled", "color", "@disabled_fg_color"},
		{"entry:disabled", "color", "@disabled_fg_color"},
		{"textview:disabled text", "color", "@disabled_fg_color"},
		{"button:focus", "border-color", "@focus_border_color"},
		{"entry:focus", "border-color", "@focus_border_color"},
		{"textview:focus", "border-color", "@focus_border_color"},
		{"entry selection", "background-color", "@selection_bg_color"},
		{"entry selection", "color", "@selection_fg_color"},
		{"textview text selection", "background-color", "@selection_bg_color"},
		{"textview text selection", "color", "@selection_fg_color"},
	} {
		assert.Equal(t, tc.want, uiProperty(t, theme.Rules, tc.selector, tc.property), tc.selector)
	}
}

func TestQssMapper_SemanticStates(t *testing.T) {
	mapped, err := mapper.NewQss().Map(uiTheme(t))
	require.NoError(t, err)
	theme, ok := mapped.(*ports.QssTheme)
	require.True(t, ok)
	assert.Equal(t, "#345678", uiProperty(t, theme.Rules, "QStatusBar", "background-color"))
	assert.Equal(t, "#6789ab", uiProperty(t, theme.Rules, "QStatusBar", "color"))
	assert.Equal(t, "#789abc", uiProperty(t, theme.Rules, "QPushButton:hover", "background-color"))
	assert.Equal(t, "#89abcd", uiProperty(t, theme.Rules, "QPushButton:pressed", "background-color"))
	for _, selector := range []string{"QTreeView::item:selected", "QListView::item:selected", "QTableView::item:selected", "QMenu::item:selected"} {
		assert.Equal(t, "#123456", uiProperty(t, theme.Rules, selector, "background-color"), selector)
		assert.Equal(t, "#56789a", uiProperty(t, theme.Rules, selector, "color"), selector)
	}
	for _, selector := range []string{"QLineEdit:focus", "QTextEdit:focus", "QPlainTextEdit:focus"} {
		assert.Equal(t, "1px solid #abcdef", uiProperty(t, theme.Rules, selector, "border"), selector)
	}
	for _, selector := range []string{"QWidget:disabled", "QPushButton:disabled", "QPushButton:disabled:hover", "QPushButton:disabled:pressed", "QMenu::item:selected:disabled"} {
		assert.Equal(t, "#9abcde", uiProperty(t, theme.Rules, selector, "color"), selector)
	}
	// Equal-specificity button rules must occur after hover/pressed rules.
	disabled, hover, pressed := -1, -1, -1
	for i, rule := range theme.Rules {
		switch rule.Selector {
		case "QPushButton:disabled":
			disabled = i
		case "QPushButton:hover":
			hover = i
		case "QPushButton:pressed":
			pressed = i
		}
	}
	assert.Greater(t, disabled, hover)
	assert.Greater(t, disabled, pressed)
}

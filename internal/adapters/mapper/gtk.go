package mapper

import (
	"errors"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// Gtk implements ports.Mapper for the GTK target.
// It maps a ResolvedTheme into a GtkTheme containing @define-color
// entries derived from semantic tokens and widget CSS selector rules
// for window, headerbar, button, entry, textview, etc.
type Gtk struct{}

// NewGtk returns a new GTK mapper.
func NewGtk() *Gtk {
	return &Gtk{}
}

// Name returns the target name for this mapper.
func (g *Gtk) Name() string {
	return "gtk"
}

// Map transforms a ResolvedTheme into a *ports.GtkTheme with @define-color
// entries from semantic tokens and widget selector rules for standard GTK
// widgets.
func (g *Gtk) Map(theme *domain.ResolvedTheme) (ports.MappedTheme, error) {
	if theme == nil {
		return nil, errors.New("gtk mapper: nil theme")
	}
	if theme.Palette == nil {
		return nil, errors.New("gtk mapper: nil palette")
	}
	if theme.Tokens == nil {
		return nil, errors.New("gtk mapper: nil tokens")
	}

	colors := buildGtkColorDefinitions(theme.Tokens)
	rules := buildGtkWidgetRules()

	return &ports.GtkTheme{
		Colors: colors,
		Rules:  rules,
	}, nil
}

// gtkColorMapping maps a semantic token path to a GTK @define-color name.
type gtkColorMapping struct {
	tokenPath string
	colorName string
}

// gtkColorMappings defines the semantic token paths and their corresponding
// GTK @define-color names following the GNOME/Adwaita naming convention.
//
//nolint:dupl // Each mapper has its own naming convention; structural similarity is intentional.
var gtkColorMappings = func() []gtkColorMapping {
	m := make([]gtkColorMapping, 0, 40)

	// Window colors
	m = append(m,
		gtkColorMapping{"surface.background", "window_bg_color"},
		gtkColorMapping{"text.primary", "window_fg_color"},
	)

	// View colors (for content areas)
	m = append(m,
		gtkColorMapping{"surface.background.sunken", "view_bg_color"},
		gtkColorMapping{"text.primary", "view_fg_color"},
	)

	// Header bar colors
	m = append(m,
		gtkColorMapping{"surface.background.sunken", "headerbar_bg_color"},
		gtkColorMapping{"text.primary", "headerbar_fg_color"},
	)

	// Sidebar colors
	m = append(m,
		gtkColorMapping{"surface.background.sidebar", "sidebar_bg_color"},
		gtkColorMapping{"text.sidebar", "sidebar_fg_color"},
	)

	// Card colors
	m = append(m,
		gtkColorMapping{"surface.background.raised", "card_bg_color"},
		gtkColorMapping{"text.primary", "card_fg_color"},
	)

	// Dialog colors
	m = append(m,
		gtkColorMapping{"surface.background.overlay", "dialog_bg_color"},
		gtkColorMapping{"text.primary", "dialog_fg_color"},
	)

	// Popover colors
	m = append(m,
		gtkColorMapping{"surface.background.popup", "popover_bg_color"},
		gtkColorMapping{"text.primary", "popover_fg_color"},
	)

	// Accent colors
	m = append(m,
		gtkColorMapping{"accent.primary", "accent_bg_color"},
		gtkColorMapping{"accent.foreground", "accent_fg_color"},
		gtkColorMapping{"accent.primary", "accent_color"},
	)

	// Status colors
	m = append(m,
		gtkColorMapping{"status.error", "error_bg_color"},
		gtkColorMapping{"text.primary", "error_fg_color"},
		gtkColorMapping{"status.error", "error_color"},
		gtkColorMapping{"status.warning", "warning_bg_color"},
		gtkColorMapping{"text.primary", "warning_fg_color"},
		gtkColorMapping{"status.warning", "warning_color"},
		gtkColorMapping{"status.success", "success_bg_color"},
		gtkColorMapping{"text.primary", "success_fg_color"},
		gtkColorMapping{"status.success", "success_color"},
	)

	// Border and scrollbar
	m = append(m,
		gtkColorMapping{"border.default", "borders"},
		gtkColorMapping{"scrollbar.thumb", "scrollbar_outline_color"},
	)

	return m
}()

// buildGtkColorDefinitions creates @define-color entries from semantic tokens.
func buildGtkColorDefinitions(ts *domain.TokenSet) []ports.GtkColorDef {
	var colors []ports.GtkColorDef

	for _, cm := range gtkColorMappings {
		tok, ok := ts.Get(cm.tokenPath)
		if !ok {
			continue
		}
		if tok.Color.IsNone {
			continue
		}
		colors = append(colors, ports.GtkColorDef{
			Name:  cm.colorName,
			Value: tok.Color.Hex(),
		})
	}

	return colors
}

// buildGtkWidgetRules creates the standard GTK widget CSS rules that reference
// @define-color names. GTK CSS uses @name syntax to reference defined colors.
func buildGtkWidgetRules() []ports.CSSRule {
	return []ports.CSSRule{
		{
			Selector: "window",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: "@window_bg_color"},
				{Property: "color", Value: "@window_fg_color"},
			},
		},
		{
			Selector: "headerbar",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: "@headerbar_bg_color"},
				{Property: "color", Value: "@headerbar_fg_color"},
			},
		},
		{
			Selector: "button",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: "@card_bg_color"},
				{Property: "color", Value: "@window_fg_color"},
			},
		},
		{
			Selector: "entry",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: "@view_bg_color"},
				{Property: "color", Value: "@view_fg_color"},
			},
		},
		{
			Selector: "textview",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: "@view_bg_color"},
				{Property: "color", Value: "@view_fg_color"},
			},
		},
		{
			Selector: ".sidebar",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: "@sidebar_bg_color"},
				{Property: "color", Value: "@sidebar_fg_color"},
			},
		},
		{
			Selector: "popover",
			Properties: []ports.CSSProperty{
				{Property: "background-color", Value: "@popover_bg_color"},
				{Property: "color", Value: "@popover_fg_color"},
			},
		},
	}
}

package mapper

import (
	"errors"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// Qss implements ports.Mapper for the Qt Style Sheet (QSS) target.
// It maps a ResolvedTheme into a QssTheme containing widget selector
// rules and pseudo-state rules with literal hex color values. Unlike
// GTK and CSS, QSS does not support variables -- all colors are inlined.
type Qss struct{}

// NewQss returns a new QSS mapper.
func NewQss() *Qss {
	return &Qss{}
}

// Name returns the target name for this mapper.
func (q *Qss) Name() string {
	return "qss"
}

// Map transforms a ResolvedTheme into a *ports.QssTheme with widget
// selector rules and pseudo-state rules. All color values are literal
// hex strings -- QSS does not support variable references.
func (q *Qss) Map(theme *domain.ResolvedTheme) (ports.MappedTheme, error) {
	if theme == nil {
		return nil, errors.New("qss mapper: nil theme")
	}
	if theme.Palette == nil {
		return nil, errors.New("qss mapper: nil palette")
	}
	if theme.Tokens == nil {
		return nil, errors.New("qss mapper: nil tokens")
	}

	rules := buildQssWidgetRules(theme.Tokens)
	pseudoRules := buildQssPseudoStateRules(theme.Tokens)
	rules = append(rules, pseudoRules...)

	return &ports.QssTheme{
		Rules: rules,
	}, nil
}

// tokenHex is a helper that extracts the hex color for a token path,
// returning an empty string if the token is not found or has no color.
func tokenHex(ts *domain.TokenSet, path string) string {
	tok, ok := ts.Get(path)
	if !ok {
		return ""
	}
	if tok.Color.IsNone {
		return ""
	}
	return tok.Color.Hex()
}

// buildQssWidgetRules creates the standard Qt widget selector rules
// with literal hex color values derived from semantic tokens.
//
//nolint:funlen // Large widget rule table is intentionally in one function.
func buildQssWidgetRules(ts *domain.TokenSet) []ports.CSSRule {
	rules := make([]ports.CSSRule, 0, 35)

	bg := tokenHex(ts, "surface.background")
	fg := tokenHex(ts, "text.primary")
	bgRaised := tokenHex(ts, "surface.background.raised")
	bgSunken := tokenHex(ts, "surface.background.sunken")
	bgPopup := tokenHex(ts, "surface.background.popup")
	textSecondary := tokenHex(ts, "text.secondary")
	borderDefault := tokenHex(ts, "border.default")
	borderFocus := tokenHex(ts, "border.focus")
	accentPrimary := tokenHex(ts, "accent.primary")
	accentFg := tokenHex(ts, "accent.foreground")
	scrollbarThumb := tokenHex(ts, "scrollbar.thumb")
	scrollbarTrack := tokenHex(ts, "scrollbar.track")
	bgHighlight := tokenHex(ts, "surface.background.highlight")
	bgSidebar := tokenHex(ts, "surface.background.sidebar")
	textSidebar := tokenHex(ts, "text.sidebar")

	// QWidget -- base widget styling
	rules = append(rules, ports.CSSRule{
		Selector: "QWidget",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bg},
			{Property: "color", Value: fg},
		},
	})

	// QMainWindow
	rules = append(rules, ports.CSSRule{
		Selector: "QMainWindow",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bg},
			{Property: "color", Value: fg},
		},
	})

	// QPushButton
	rules = append(rules, ports.CSSRule{
		Selector: "QPushButton",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgRaised},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QLineEdit
	rules = append(rules, ports.CSSRule{
		Selector: "QLineEdit",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QTextEdit
	rules = append(rules, ports.CSSRule{
		Selector: "QTextEdit",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QPlainTextEdit
	rules = append(rules, ports.CSSRule{
		Selector: "QPlainTextEdit",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QComboBox
	rules = append(rules, ports.CSSRule{
		Selector: "QComboBox",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgRaised},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QMenuBar
	rules = append(rules, ports.CSSRule{
		Selector: "QMenuBar",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
		},
	})

	// QMenu
	rules = append(rules, ports.CSSRule{
		Selector: "QMenu",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgPopup},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QMenu::item:selected
	rules = append(rules, ports.CSSRule{
		Selector: "QMenu::item:selected",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgHighlight},
			{Property: "color", Value: fg},
		},
	})

	// QToolBar
	rules = append(rules, ports.CSSRule{
		Selector: "QToolBar",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
			{Property: "border", Value: "none"},
		},
	})

	// QStatusBar
	rules = append(rules, ports.CSSRule{
		Selector: "QStatusBar",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: textSecondary},
		},
	})

	// QTabWidget::pane
	rules = append(rules, ports.CSSRule{
		Selector: "QTabWidget::pane",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QTabBar::tab
	rules = append(rules, ports.CSSRule{
		Selector: "QTabBar::tab",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgRaised},
			{Property: "color", Value: textSecondary},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QTabBar::tab:selected
	rules = append(rules, ports.CSSRule{
		Selector: "QTabBar::tab:selected",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bg},
			{Property: "color", Value: fg},
			{Property: "border-bottom", Value: "2px solid " + accentPrimary},
		},
	})

	// QCheckBox and QRadioButton
	rules = append(rules, ports.CSSRule{
		Selector: "QCheckBox",
		Properties: []ports.CSSProperty{
			{Property: "color", Value: fg},
		},
	})

	rules = append(rules, ports.CSSRule{
		Selector: "QRadioButton",
		Properties: []ports.CSSProperty{
			{Property: "color", Value: fg},
		},
	})

	// QGroupBox
	rules = append(rules, ports.CSSRule{
		Selector: "QGroupBox",
		Properties: []ports.CSSProperty{
			{Property: "border", Value: "1px solid " + borderDefault},
			{Property: "color", Value: fg},
		},
	})

	// QScrollBar (vertical)
	rules = append(rules, ports.CSSRule{
		Selector: "QScrollBar",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: scrollbarTrack},
		},
	})

	// QScrollBar::handle
	rules = append(rules, ports.CSSRule{
		Selector: "QScrollBar::handle",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: scrollbarThumb},
		},
	})

	// QToolTip
	rules = append(rules, ports.CSSRule{
		Selector: "QToolTip",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgPopup},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QProgressBar
	rules = append(rules, ports.CSSRule{
		Selector: "QProgressBar",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgRaised},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QProgressBar::chunk
	rules = append(rules, ports.CSSRule{
		Selector: "QProgressBar::chunk",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: accentPrimary},
		},
	})

	// QDockWidget
	rules = append(rules, ports.CSSRule{
		Selector: "QDockWidget",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSidebar},
			{Property: "color", Value: textSidebar},
		},
	})

	// QHeaderView::section (table/tree header)
	rules = append(rules, ports.CSSRule{
		Selector: "QHeaderView::section",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgRaised},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QTreeView / QListView / QTableView
	rules = append(rules, ports.CSSRule{
		Selector: "QTreeView",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
		},
	})

	rules = append(rules, ports.CSSRule{
		Selector: "QListView",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
		},
	})

	rules = append(rules, ports.CSSRule{
		Selector: "QTableView",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
		},
	})

	// Selected items in list/tree/table views
	rules = append(rules, ports.CSSRule{
		Selector: "QTreeView::item:selected",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: accentPrimary},
			{Property: "color", Value: accentFg},
		},
	})

	// QSpinBox
	rules = append(rules, ports.CSSRule{
		Selector: "QSpinBox",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgSunken},
			{Property: "color", Value: fg},
			{Property: "border", Value: "1px solid " + borderDefault},
		},
	})

	// QSlider::groove
	rules = append(rules, ports.CSSRule{
		Selector: "QSlider::groove",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: bgRaised},
		},
	})

	// QSlider::handle
	rules = append(rules, ports.CSSRule{
		Selector: "QSlider::handle",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: accentPrimary},
		},
	})

	// QLabel
	rules = append(rules, ports.CSSRule{
		Selector: "QLabel",
		Properties: []ports.CSSProperty{
			{Property: "color", Value: fg},
		},
	})

	// Focus border for QLineEdit and QTextEdit
	rules = append(rules, ports.CSSRule{
		Selector: "QLineEdit:focus",
		Properties: []ports.CSSProperty{
			{Property: "border", Value: "1px solid " + borderFocus},
		},
	})

	rules = append(rules, ports.CSSRule{
		Selector: "QTextEdit:focus",
		Properties: []ports.CSSProperty{
			{Property: "border", Value: "1px solid " + borderFocus},
		},
	})

	return rules
}

// buildQssPseudoStateRules creates pseudo-state rules for interactive widgets.
// Qt style sheets use the :pseudo-state syntax for hover, pressed, focus, etc.
func buildQssPseudoStateRules(ts *domain.TokenSet) []ports.CSSRule {
	rules := make([]ports.CSSRule, 0, 6)

	stateHover := tokenHex(ts, "state.hover")
	stateActive := tokenHex(ts, "state.active")
	fg := tokenHex(ts, "text.primary")
	accentPrimary := tokenHex(ts, "accent.primary")
	accentFg := tokenHex(ts, "accent.foreground")

	// QPushButton:hover
	rules = append(rules, ports.CSSRule{
		Selector: "QPushButton:hover",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: stateHover},
			{Property: "color", Value: fg},
		},
	})

	// QPushButton:pressed
	rules = append(rules, ports.CSSRule{
		Selector: "QPushButton:pressed",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: stateActive},
			{Property: "color", Value: fg},
		},
	})

	// QComboBox:hover
	rules = append(rules, ports.CSSRule{
		Selector: "QComboBox:hover",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: stateHover},
		},
	})

	// QTabBar::tab:hover
	rules = append(rules, ports.CSSRule{
		Selector: "QTabBar::tab:hover",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: stateHover},
		},
	})

	// QCheckBox::indicator:checked
	rules = append(rules, ports.CSSRule{
		Selector: "QCheckBox::indicator:checked",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: accentPrimary},
			{Property: "color", Value: accentFg},
		},
	})

	// QRadioButton::indicator:checked
	rules = append(rules, ports.CSSRule{
		Selector: "QRadioButton::indicator:checked",
		Properties: []ports.CSSProperty{
			{Property: "background-color", Value: accentPrimary},
			{Property: "color", Value: accentFg},
		},
	})

	return rules
}

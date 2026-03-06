package viewer

import (
	"strings"
	"testing"
)

// TestView_TextStatusPage verifies the Text & Status page renders realistic content.
func TestView_TextStatusPage(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})

	m.currentPage = PageTextStatus
	m.tokens = TokenData{
		Text: map[string]string{
			"text.primary":   "#c0caf5",
			"text.secondary": "#a9b1d6",
			"text.muted":     "#565f89",
		},
		Status: map[string]string{
			"status.error":   "#f7768e",
			"status.warning": "#e0af68",
			"status.success": "#9ece6a",
			"status.info":    "#7dcfff",
		},
	}

	view := m.View()

	// Should contain the page title.
	if !strings.Contains(view.Content, "Text & Status") {
		t.Error("text status page missing title")
	}

	// Should contain realistic text content (not token names).
	textPhrases := []string{
		"Lorem ipsum",
		"Primary text",
		"Secondary text",
		"Muted text",
	}
	for _, phrase := range textPhrases {
		if !strings.Contains(view.Content, phrase) {
			t.Errorf("text status page missing phrase %q", phrase)
		}
	}

	// Should contain status message labels.
	statusLabels := []string{"Error:", "Warning:", "Success:", "Info:"}
	for _, label := range statusLabels {
		if !strings.Contains(view.Content, label) {
			t.Errorf("text status page missing status label %q", label)
		}
	}
}

// TestView_InteractivePage verifies the Interactive Components page.
func TestView_InteractivePage(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})

	m.currentPage = PageInteractive
	m.tokens = TokenData{
		Text: map[string]string{
			"text.primary": "#c0caf5",
			"text.muted":   "#565f89",
		},
		Surface: map[string]string{
			"surface.background":        "#1a1b26",
			"surface.background.raised": "#24283b",
			"surface.background.sunken": "#16161e",
		},
	}

	view := m.View()

	// Should contain the page title.
	if !strings.Contains(view.Content, "Interactive Components") {
		t.Error("interactive page missing title")
	}

	// Should contain button examples.
	buttonLabels := []string{"Submit", "Cancel", "Disabled"}
	for _, label := range buttonLabels {
		if !strings.Contains(view.Content, label) {
			t.Errorf("interactive page missing button %q", label)
		}
	}

	// Should contain input field section.
	if !strings.Contains(view.Content, "Input Fields") {
		t.Error("interactive page missing Input Fields section")
	}

	// Should contain selection list section.
	if !strings.Contains(view.Content, "Selection List") {
		t.Error("interactive page missing Selection List section")
	}
}

// TestView_DataDisplayPage verifies the Data Display page.
func TestView_DataDisplayPage(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})

	m.currentPage = PageDataDisplay
	m.tokens = TokenData{
		Text: map[string]string{
			"text.primary":   "#c0caf5",
			"text.secondary": "#a9b1d6",
		},
		Syntax: map[string]string{
			"syntax.keyword":  "#bb9af7",
			"syntax.string":   "#9ece6a",
			"syntax.function": "#7aa2f7",
		},
	}

	view := m.View()

	// Should contain the page title.
	if !strings.Contains(view.Content, "Data Display") {
		t.Error("data display page missing title")
	}

	// Should contain table section with headers.
	if !strings.Contains(view.Content, "Table") {
		t.Error("data display page missing Table section")
	}

	// Should contain column headers for sample table.
	tableHeaders := []string{"Name", "Status", "Progress"}
	for _, header := range tableHeaders {
		if !strings.Contains(view.Content, header) {
			t.Errorf("data display page missing table header %q", header)
		}
	}

	// Should contain dialog section.
	if !strings.Contains(view.Content, "Dialog") {
		t.Error("data display page missing Dialog section")
	}

	// Should contain code block section.
	if !strings.Contains(view.Content, "Code") {
		t.Error("data display page missing Code section")
	}
}

// TestView_TwoPanelLayout verifies 2-panel layout with themes on left.
func TestView_TwoPanelLayout(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"tokyo-night-dark", "gruvbox-dark", "catppuccin-mocha"},
	})

	m.currentPage = PageTextStatus
	m.cursor = 1 // gruvbox-dark highlighted
	m.width = 120
	m.height = 40

	view := m.View()

	// Should contain all theme names in the left panel.
	for _, theme := range m.themes {
		if !strings.Contains(view.Content, theme) {
			t.Errorf("view missing theme %q", theme)
		}
	}

	// Should contain "Styles" title for left panel.
	if !strings.Contains(view.Content, "Styles") {
		t.Error("view missing Styles title for left panel")
	}

	// Should contain the content page title on the right.
	if !strings.Contains(view.Content, "Text & Status") {
		t.Error("view missing content page title")
	}
}

// TestView_ThemeListShowsSelection verifies selected theme is marked.
func TestView_ThemeListShowsSelection(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1", "theme2", "theme3"},
	})

	m.selectedTheme = "theme2" //nolint:goconst // test data
	m.cursor = 1
	m.width = 120
	m.height = 40

	view := m.View()

	// The cursor indicator should appear.
	if !strings.Contains(view.Content, ">") {
		t.Error("view missing cursor indicator")
	}
}

// TestView_HelpFooter verifies help footer is rendered.
func TestView_HelpFooter(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})

	view := m.View()

	// Should contain navigation hints.
	hints := []string{"Tab:", "↑/↓/j/k:", "Enter:", "q/Esc:"}
	for _, hint := range hints {
		if !strings.Contains(view.Content, hint) {
			t.Errorf("help footer missing hint containing %q", hint)
		}
	}
}

// TestView_HelpFooterAtBottom verifies help footer is pinned to window bottom.
func TestView_HelpFooterAtBottom(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})
	m.width = 120
	m.height = 30

	view := m.View()

	// Count lines in the view.
	lines := strings.Split(view.Content, "\n")

	// The view should use the full height (minus 1 for the help line).
	// Last non-empty line should contain help hints.
	lastLine := ""
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			lastLine = lines[i]
			break
		}
	}

	if !strings.Contains(lastLine, "Tab:") {
		t.Errorf("help footer should be at bottom, got last line: %q", lastLine)
	}

	// View should have approximately height lines (allow some variance for borders).
	if len(lines) < m.height-2 {
		t.Errorf("view has %d lines, expected at least %d to fill window", len(lines), m.height-2)
	}
}

// TestView_StatusBarSimulation verifies status bar is rendered with powerline style.
func TestView_StatusBarSimulation(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})
	m.currentPage = PageInteractive
	m.width = 120
	m.height = 40
	m.tokens = TokenData{
		Text: map[string]string{
			"text.primary": "#c0caf5",
		},
		Surface: map[string]string{
			"surface.background": "#1a1b26",
		},
		Statusline: map[string]string{
			"statusline.a.bg": "#7aa2f7",
			"statusline.a.fg": "#1a1b26",
			"statusline.b.bg": "#3b4261",
			"statusline.b.fg": "#c0caf5",
			"statusline.c.bg": "#24283b",
			"statusline.c.fg": "#a9b1d6",
		},
	}

	view := m.View()

	// Should contain status bar section.
	if !strings.Contains(view.Content, "Status Bar") {
		t.Error("interactive page missing Status Bar section")
	}

	// Should contain powerline separator characters.
	if !strings.Contains(view.Content, "") {
		t.Error("status bar missing powerline separator")
	}

	// Should contain sample content like mode indicator.
	if !strings.Contains(view.Content, "NORMAL") || !strings.Contains(view.Content, "main") {
		t.Error("status bar missing sample content (NORMAL mode or main branch)")
	}
}

// TestTokenData_HasStatuslineField verifies TokenData includes Statusline map.
func TestTokenData_HasStatuslineField(t *testing.T) {
	td := TokenData{
		Statusline: map[string]string{
			"statusline.a.bg": "#7aa2f7",
		},
	}

	if td.Statusline["statusline.a.bg"] != "#7aa2f7" {
		t.Error("TokenData.Statusline field not working correctly")
	}
}

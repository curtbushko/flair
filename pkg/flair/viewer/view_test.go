package viewer

import (
	"strings"
	"testing"
)

// TestView_PalettePage verifies palette page renders base colors.
func TestView_PalettePage(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})

	// Switch to palette page.
	m.currentPage = PagePalette
	m.palette = PaletteData{
		Colors: [24]string{
			"#1a1b26", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#c0caf5",
			"#414868", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#c0caf5",
			"#15161e", "#101014", "#ff9e64", "#9ece6a",
			"#73daca", "#7dcfff", "#2ac3de", "#ff007c",
		},
	}

	view := m.View()

	// Should contain base slot labels.
	labels := []string{"base00", "base01", "base07", "base0D", "base17"}
	for _, label := range labels {
		if !strings.Contains(view, label) {
			t.Errorf("palette view missing label %q", label)
		}
	}

	// Should contain hex colors (lowercase).
	if !strings.Contains(view, "#1a1b26") && !strings.Contains(view, "1a1b26") {
		t.Error("palette view missing hex color")
	}
}

// TestView_TokensPage verifies tokens page renders category headers.
func TestView_TokensPage(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})

	m.currentPage = PageTokens
	m.tokens = TokenData{
		Surface: map[string]string{"surface.background": "#1a1b26"},
		Text:    map[string]string{"text.primary": "#c0caf5"},
		Status:  map[string]string{"status.error": "#f7768e"},
	}

	view := m.View()

	// Should contain category section headers.
	headers := []string{"Surface", "Text", "Status"}
	for _, header := range headers {
		if !strings.Contains(view, header) {
			t.Errorf("tokens view missing header %q", header)
		}
	}
}

// TestView_ComponentsPage verifies components page renders styled examples.
func TestView_ComponentsPage(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"test-theme"},
	})

	m.currentPage = PageComponents
	m.tokens = TokenData{
		Status: map[string]string{
			"status.error":   "#f7768e",
			"status.warning": "#e0af68",
			"status.success": "#9ece6a",
		},
	}

	view := m.View()

	// Should contain component examples with token names as labels.
	labels := []string{"status.error", "status.warning", "status.success"}
	for _, label := range labels {
		if !strings.Contains(view, label) {
			t.Errorf("components view missing label %q", label)
		}
	}
}

// TestView_SelectorPage verifies selector page renders theme list.
func TestView_SelectorPage(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"tokyo-night-dark", "gruvbox-dark", "catppuccin-mocha"},
	})

	m.currentPage = PageSelector
	m.cursor = 1 // gruvbox-dark selected

	view := m.View()

	// Should contain all theme names.
	for _, theme := range m.themes {
		if !strings.Contains(view, theme) {
			t.Errorf("selector view missing theme %q", theme)
		}
	}
}

package viewer

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestModel_Init verifies model initializes with themes loaded.
func TestModel_Init(t *testing.T) {
	themes := []string{"tokyo-night-dark", "gruvbox-dark", "catppuccin-mocha"}

	opts := Options{
		Themes:       themes,
		InitialTheme: "gruvbox-dark",
	}

	m := NewModel(opts)

	// Model should have themes loaded.
	if len(m.themes) != len(themes) {
		t.Errorf("got %d themes, want %d", len(m.themes), len(themes))
	}

	// Current page should be text status (first content page in 2-panel layout).
	if m.currentPage != PageTextStatus {
		t.Errorf("got page %v, want PageTextStatus", m.currentPage)
	}

	// Initial theme should be highlighted (cursor at that index).
	wantIdx := 1 // gruvbox-dark is at index 1
	if m.cursor != wantIdx {
		t.Errorf("cursor at %d, want %d", m.cursor, wantIdx)
	}

	// Selected theme name should match.
	if m.selectedTheme != "gruvbox-dark" {
		t.Errorf("selectedTheme = %q, want %q", m.selectedTheme, "gruvbox-dark")
	}
}

// TestModel_TabSwitchesPages verifies Tab key cycles through content pages.
func TestModel_TabSwitchesPages(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1", "theme2"},
	})

	// Start on text status page (first content page in 2-panel layout).
	if m.currentPage != PageTextStatus {
		t.Fatalf("initial page = %v, want PageTextStatus", m.currentPage)
	}

	// Tab to interactive page.
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageInteractive {
		t.Errorf("after Tab: page = %v, want PageInteractive", m.currentPage)
	}

	// Tab to data display page.
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageDataDisplay {
		t.Errorf("after Tab: page = %v, want PageDataDisplay", m.currentPage)
	}

	// Tab to bubbletea page.
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageBubbletea {
		t.Errorf("after Tab: page = %v, want PageBubbletea", m.currentPage)
	}

	// Tab to huh page.
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageHuh {
		t.Errorf("after Tab: page = %v, want PageHuh", m.currentPage)
	}

	// Tab to bubbles page.
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageBubbles {
		t.Errorf("after Tab: page = %v, want PageBubbles", m.currentPage)
	}

	// Tab wraps back to text status.
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageTextStatus {
		t.Errorf("after Tab (wrap): page = %v, want PageTextStatus", m.currentPage)
	}
}

// TestUpdate_TabCyclesBubbletea verifies Tab from DataDisplay goes to Bubbletea.
func TestUpdate_TabCyclesBubbletea(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1"},
	})

	// Set to PageDataDisplay.
	m.currentPage = PageDataDisplay

	// Tab should go to PageBubbletea.
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	if m.currentPage != PageBubbletea {
		t.Errorf("after Tab from DataDisplay: page = %v, want PageBubbletea", m.currentPage)
	}
}

// TestUpdate_TabCyclesHuh verifies Tab from Bubbletea goes to Huh.
func TestUpdate_TabCyclesHuh(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1"},
	})

	// Set to PageBubbletea.
	m.currentPage = PageBubbletea

	// Tab should go to PageHuh.
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	if m.currentPage != PageHuh {
		t.Errorf("after Tab from Bubbletea: page = %v, want PageHuh", m.currentPage)
	}
}

// TestUpdate_TabCyclesBubbles verifies Tab from Huh goes to Bubbles.
func TestUpdate_TabCyclesBubbles(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1"},
	})

	// Set to PageHuh.
	m.currentPage = PageHuh

	// Tab should go to PageBubbles.
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	if m.currentPage != PageBubbles {
		t.Errorf("after Tab from Huh: page = %v, want PageBubbles", m.currentPage)
	}
}

// TestUpdate_TabWrapsFromBubbles verifies Tab from Bubbles wraps to TextStatus.
func TestUpdate_TabWrapsFromBubbles(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1"},
	})

	// Set to PageBubbles (last page).
	m.currentPage = PageBubbles

	// Tab should wrap to PageTextStatus.
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	if m.currentPage != PageTextStatus {
		t.Errorf("after Tab from Bubbles: page = %v, want PageTextStatus", m.currentPage)
	}
}

// TestModel_Navigation verifies j/k keys navigate in selector.
func TestModel_Navigation(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1", "theme2", "theme3"},
	})

	// Initial cursor at 0.
	if m.cursor != 0 {
		t.Fatalf("initial cursor = %d, want 0", m.cursor)
	}

	// Press 'j' to move down.
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'j'})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Errorf("after j: cursor = %d, want 1", m.cursor)
	}

	// Press 'j' again.
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'j'})
	m = updated.(Model)
	if m.cursor != 2 {
		t.Errorf("after j: cursor = %d, want 2", m.cursor)
	}

	// Press 'j' at bottom (should stay at 2).
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'j'})
	m = updated.(Model)
	if m.cursor != 2 {
		t.Errorf("after j at bottom: cursor = %d, want 2", m.cursor)
	}

	// Press 'k' to move up.
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'k'})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Errorf("after k: cursor = %d, want 1", m.cursor)
	}

	// Press 'k' again.
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'k'})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Errorf("after k: cursor = %d, want 0", m.cursor)
	}

	// Press 'k' at top (should stay at 0).
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'k'})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Errorf("after k at top: cursor = %d, want 0", m.cursor)
	}
}

// TestModel_NavigationReloadsThemeData verifies j/k keys reload preview data immediately.
func TestModel_NavigationReloadsThemeData(t *testing.T) {
	loader := &mockThemeLoader{
		palettes: map[string]PaletteData{
			"theme1": {Colors: [24]string{"#111111"}},
			"theme2": {Colors: [24]string{"#222222"}},
			"theme3": {Colors: [24]string{"#333333"}},
		},
		tokens: map[string]TokenData{
			"theme1": {Status: map[string]string{"status.error": "#ff0000"}},
			"theme2": {Status: map[string]string{"status.error": "#00ff00"}},
			"theme3": {Status: map[string]string{"status.error": "#0000ff"}},
		},
	}

	m := NewModel(Options{
		Themes:       []string{"theme1", "theme2", "theme3"},
		InitialTheme: "theme1",
		ThemeLoader:  loader,
	})

	// Verify initial data loaded for theme1.
	if m.palette.Colors[0] != "#111111" {
		t.Errorf("initial palette = %q, want #111111", m.palette.Colors[0])
	}

	// Press 'j' to move to theme2 - preview data should reload immediately.
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'j'})
	m = updated.(Model)

	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}
	// selectedTheme should NOT change until Enter is pressed.
	if m.selectedTheme != "theme1" {
		t.Errorf("selectedTheme = %q, want theme1 (unchanged)", m.selectedTheme)
	}
	// But preview data should be loaded for theme2.
	if m.palette.Colors[0] != "#222222" {
		t.Errorf("after j: palette = %q, want #222222", m.palette.Colors[0])
	}
	if m.tokens.Status["status.error"] != "#00ff00" {
		t.Errorf("after j: status.error = %q, want #00ff00", m.tokens.Status["status.error"])
	}

	// Press 'j' again to move to theme3.
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'j'})
	m = updated.(Model)

	if m.palette.Colors[0] != "#333333" {
		t.Errorf("after second j: palette = %q, want #333333", m.palette.Colors[0])
	}

	// Press 'k' to go back to theme2.
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'k'})
	m = updated.(Model)

	if m.palette.Colors[0] != "#222222" {
		t.Errorf("after k: palette = %q, want #222222", m.palette.Colors[0])
	}
}

// TestModel_EnterConfirmsSelection verifies Enter key confirms selection.
func TestModel_EnterConfirmsSelection(t *testing.T) {
	var selectedName string
	m := NewModel(Options{
		Themes: []string{"theme1", "theme2", "theme3"},
		OnSelect: func(name string) {
			selectedName = name
		},
	})

	// Navigate to theme2.
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'j'})
	m = updated.(Model)

	// selectedTheme should still be empty (no initial theme, no Enter pressed).
	if m.selectedTheme != "" {
		t.Errorf("before Enter: selectedTheme = %q, want empty", m.selectedTheme)
	}
	if m.selectedIndex != -1 {
		t.Errorf("before Enter: selectedIndex = %d, want -1", m.selectedIndex)
	}

	// Press Enter to confirm selection.
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if m.selectedTheme != "theme2" {
		t.Errorf("after Enter: selectedTheme = %q, want theme2", m.selectedTheme)
	}
	if m.selectedIndex != 1 {
		t.Errorf("after Enter: selectedIndex = %d, want 1", m.selectedIndex)
	}
	if selectedName != "theme2" {
		t.Errorf("onSelect called with %q, want theme2", selectedName)
	}
}

// TestModel_QuitOnQ verifies 'q' key returns tea.Quit.
func TestModel_QuitOnQ(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1"},
	})

	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})

	// cmd should be tea.Quit.
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}

	// Execute the command and check it returns tea.QuitMsg.
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", msg)
	}
}

// TestModel_ImplementsTeaModel verifies Model implements tea.Model interface.
func TestModel_ImplementsTeaModel(t *testing.T) {
	var _ tea.Model = Model{}
}

// mockThemeLoader implements ThemeLoader for testing.
type mockThemeLoader struct {
	palettes map[string]PaletteData
	tokens   map[string]TokenData
}

func (m *mockThemeLoader) LoadPalette(name string) (PaletteData, error) {
	if pd, ok := m.palettes[name]; ok {
		return pd, nil
	}
	return PaletteData{}, nil
}

func (m *mockThemeLoader) LoadTokens(name string) (TokenData, error) {
	if td, ok := m.tokens[name]; ok {
		return td, nil
	}
	return TokenData{}, nil
}

// TestModel_InitialThemeLoadsData verifies initial theme loads data on model creation.
func TestModel_InitialThemeLoadsData(t *testing.T) {
	loader := &mockThemeLoader{
		palettes: map[string]PaletteData{
			"gruvbox": {Colors: [24]string{"#282828", "#cc241d"}},
		},
		tokens: map[string]TokenData{
			"gruvbox": {Surface: map[string]string{"surface.background": "#282828"}},
		},
	}

	m := NewModel(Options{
		Themes:       []string{"gruvbox"},
		InitialTheme: "gruvbox",
		ThemeLoader:  loader,
	})

	// Verify data was loaded.
	if m.palette.Colors[0] != "#282828" {
		t.Errorf("palette[0] = %q, want #282828", m.palette.Colors[0])
	}
	if m.tokens.Surface["surface.background"] != "#282828" {
		t.Errorf("surface.background = %q, want #282828", m.tokens.Surface["surface.background"])
	}
}

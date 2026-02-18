package viewer

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

	// Current page should be selector.
	if m.currentPage != PageSelector {
		t.Errorf("got page %v, want PageSelector", m.currentPage)
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

// TestModel_TabSwitchesPages verifies Tab key cycles through pages.
func TestModel_TabSwitchesPages(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1", "theme2"},
	})

	// Start on selector page.
	if m.currentPage != PageSelector {
		t.Fatalf("initial page = %v, want PageSelector", m.currentPage)
	}

	// Tab to palette page.
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PagePalette {
		t.Errorf("after Tab: page = %v, want PagePalette", m.currentPage)
	}

	// Tab to tokens page.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageTokens {
		t.Errorf("after Tab: page = %v, want PageTokens", m.currentPage)
	}

	// Tab to components page.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageComponents {
		t.Errorf("after Tab: page = %v, want PageComponents", m.currentPage)
	}

	// Tab wraps back to selector.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.currentPage != PageSelector {
		t.Errorf("after Tab (wrap): page = %v, want PageSelector", m.currentPage)
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
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Errorf("after j: cursor = %d, want 1", m.cursor)
	}

	// Press 'j' again.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if m.cursor != 2 {
		t.Errorf("after j: cursor = %d, want 2", m.cursor)
	}

	// Press 'j' at bottom (should stay at 2).
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if m.cursor != 2 {
		t.Errorf("after j at bottom: cursor = %d, want 2", m.cursor)
	}

	// Press 'k' to move up.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Errorf("after k: cursor = %d, want 1", m.cursor)
	}

	// Press 'k' again.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Errorf("after k: cursor = %d, want 0", m.cursor)
	}

	// Press 'k' at top (should stay at 0).
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Errorf("after k at top: cursor = %d, want 0", m.cursor)
	}
}

// TestModel_EnterAppliesTheme verifies Enter on selector calls OnSelect callback.
func TestModel_EnterAppliesTheme(t *testing.T) {
	var selectedName string
	m := NewModel(Options{
		Themes: []string{"theme1", "theme2"},
		OnSelect: func(name string) {
			selectedName = name
		},
	})

	// Move to theme2.
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)

	// Press Enter.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	if selectedName != "theme2" {
		t.Errorf("OnSelect called with %q, want %q", selectedName, "theme2")
	}

	// Verify selectedTheme is updated.
	if m.selectedTheme != "theme2" {
		t.Errorf("selectedTheme = %q, want %q", m.selectedTheme, "theme2")
	}
}

// TestModel_QuitOnQ verifies 'q' key returns tea.Quit.
func TestModel_QuitOnQ(t *testing.T) {
	m := NewModel(Options{
		Themes: []string{"theme1"},
	})

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

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

// TestModel_EnterReloadsThemeData verifies Enter key reloads palette and token data.
func TestModel_EnterReloadsThemeData(t *testing.T) {
	loader := &mockThemeLoader{
		palettes: map[string]PaletteData{
			"theme1": {Colors: [24]string{"#111111"}},
			"theme2": {Colors: [24]string{"#222222"}},
		},
		tokens: map[string]TokenData{
			"theme1": {Status: map[string]string{"status.error": "#ff0000"}},
			"theme2": {Status: map[string]string{"status.error": "#00ff00"}},
		},
	}

	m := NewModel(Options{
		Themes:       []string{"theme1", "theme2"},
		InitialTheme: "theme1",
		ThemeLoader:  loader,
	})

	// Verify initial data is loaded.
	if m.palette.Colors[0] != "#111111" {
		t.Errorf("initial palette color = %q, want #111111", m.palette.Colors[0])
	}
	if m.tokens.Status["status.error"] != "#ff0000" {
		t.Errorf("initial status.error = %q, want #ff0000", m.tokens.Status["status.error"])
	}

	// Move to theme2 and press Enter.
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	// Verify data is reloaded for theme2.
	if m.palette.Colors[0] != "#222222" {
		t.Errorf("after Enter: palette color = %q, want #222222", m.palette.Colors[0])
	}
	if m.tokens.Status["status.error"] != "#00ff00" {
		t.Errorf("after Enter: status.error = %q, want #00ff00", m.tokens.Status["status.error"])
	}
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

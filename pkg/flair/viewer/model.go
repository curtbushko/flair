// Package viewer provides a bubbletea-based TUI for browsing and selecting
// flair themes. It displays theme palettes, semantic tokens, and styled
// component examples with live theme switching.
package viewer

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Page represents which content page is displayed in the right panel.
type Page int

const (
	// PageTextStatus shows text paragraphs and status messages.
	PageTextStatus Page = iota
	// PageInteractive shows buttons, inputs, and selection lists.
	PageInteractive
	// PageDataDisplay shows tables, dialogs, and code blocks.
	PageDataDisplay
)

// pageCount is the total number of content pages for Tab cycling.
const pageCount = 3

// Options configures the viewer Model.
type Options struct {
	// Themes is the list of available theme names.
	Themes []string

	// InitialTheme is the theme to pre-select on startup.
	InitialTheme string

	// OnSelect is called when the user selects a theme with Enter.
	// If nil, selection only updates the internal state.
	OnSelect func(name string)

	// ThemeLoader loads theme data for display. If nil, palette/token
	// pages will show placeholder content.
	ThemeLoader ThemeLoader
}

// ThemeLoader loads theme data for rendering in the viewer.
type ThemeLoader interface {
	// LoadPalette returns base24 colors for a theme.
	LoadPalette(name string) (PaletteData, error)

	// LoadTokens returns semantic tokens for a theme.
	LoadTokens(name string) (TokenData, error)
}

// PaletteData contains the base24 colors for display.
type PaletteData struct {
	Colors [24]string // Hex colors for base00-base17
}

// TokenData contains semantic tokens grouped by category.
type TokenData struct {
	Surface map[string]string // surface.* tokens
	Text    map[string]string // text.* tokens
	Status  map[string]string // status.* tokens
	Syntax  map[string]string // syntax.* tokens
	Diff    map[string]string // diff.* tokens
}

// Model implements tea.Model for the style viewer TUI.
type Model struct {
	themes        []string
	cursor        int // Current cursor position (for navigation/preview)
	selectedIndex int // Index of confirmed selection (-1 if none)
	currentPage   Page
	selectedTheme string // Name of confirmed selection
	onSelect      func(string)
	themeLoader   ThemeLoader

	// Cached data for previewed theme (at cursor position).
	palette PaletteData
	tokens  TokenData

	// Terminal dimensions.
	width  int
	height int
}

// NewModel creates a new viewer Model with the given options.
func NewModel(opts Options) Model {
	themes := opts.Themes
	if themes == nil {
		themes = []string{}
	}

	m := Model{
		themes:        themes,
		selectedIndex: -1, // No selection until Enter is pressed
		currentPage:   PageTextStatus,
		onSelect:      opts.OnSelect,
		themeLoader:   opts.ThemeLoader,
	}

	// Set initial theme and cursor position.
	if opts.InitialTheme != "" {
		m.selectedTheme = opts.InitialTheme
		for i, name := range opts.Themes {
			if name == opts.InitialTheme {
				m.cursor = i
				m.selectedIndex = i
				break
			}
		}
		// Load data for initial theme.
		m.loadThemeData(opts.InitialTheme)
	} else if len(themes) > 0 {
		// Load first theme for preview.
		m.loadThemeData(themes[0])
	}

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

// handleKey processes keyboard input.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyTab:
		m.currentPage = (m.currentPage + 1) % pageCount
		return m, nil

	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m.previewCurrentTheme()
		}
		return m, nil

	case tea.KeyDown:
		if m.cursor < len(m.themes)-1 {
			m.cursor++
			m.previewCurrentTheme()
		}
		return m, nil

	case tea.KeyEnter:
		m.confirmSelection()
		return m, nil

	case tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
				m.previewCurrentTheme()
			}
			return m, nil
		case "k":
			if m.cursor > 0 {
				m.cursor--
				m.previewCurrentTheme()
			}
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// previewCurrentTheme loads theme data for the cursor position (live preview).
func (m *Model) previewCurrentTheme() {
	if len(m.themes) == 0 || m.cursor >= len(m.themes) {
		return
	}
	m.loadThemeData(m.themes[m.cursor])
}

// confirmSelection marks the current cursor position as selected and calls onSelect.
func (m *Model) confirmSelection() {
	if len(m.themes) == 0 || m.cursor >= len(m.themes) {
		return
	}
	m.selectedIndex = m.cursor
	m.selectedTheme = m.themes[m.cursor]
	if m.onSelect != nil {
		m.onSelect(m.selectedTheme)
	}
}

// loadThemeData loads palette and token data for the given theme using ThemeLoader.
// If ThemeLoader is nil, this is a no-op.
func (m *Model) loadThemeData(themeName string) {
	if m.themeLoader == nil {
		return
	}

	// Load palette data.
	if palette, err := m.themeLoader.LoadPalette(themeName); err == nil {
		m.palette = palette
	}

	// Load token data.
	if tokens, err := m.themeLoader.LoadTokens(themeName); err == nil {
		m.tokens = tokens
	}
}

// View is implemented in view.go

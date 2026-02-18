// Package viewer provides a bubbletea-based TUI for browsing and selecting
// flair themes. It displays theme palettes, semantic tokens, and styled
// component examples with live theme switching.
package viewer

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Page represents which page is currently displayed in the viewer.
type Page int

const (
	// PageSelector shows the theme list for selection.
	PageSelector Page = iota
	// PagePalette shows the base00-base17 palette colors.
	PagePalette
	// PageTokens shows semantic tokens grouped by category.
	PageTokens
	// PageComponents shows styled component examples.
	PageComponents
)

// pageCount is the total number of pages for Tab cycling.
const pageCount = 4

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
	cursor        int
	currentPage   Page
	selectedTheme string
	onSelect      func(string)
	themeLoader   ThemeLoader

	// Cached data for current theme.
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
		themes:      themes,
		currentPage: PageSelector,
		onSelect:    opts.OnSelect,
		themeLoader: opts.ThemeLoader,
	}

	// Set initial theme and cursor position.
	if opts.InitialTheme != "" {
		m.selectedTheme = opts.InitialTheme
		for i, name := range opts.Themes {
			if name == opts.InitialTheme {
				m.cursor = i
				break
			}
		}
		// Load data for initial theme.
		m.loadThemeData(opts.InitialTheme)
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

	case tea.KeyEnter:
		if m.currentPage == PageSelector && len(m.themes) > 0 {
			m.selectedTheme = m.themes[m.cursor]
			if m.onSelect != nil {
				m.onSelect(m.selectedTheme)
			}
			// Reload palette and token data for the newly selected theme.
			m.loadThemeData(m.selectedTheme)
		}
		return m, nil

	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			}
			return m, nil
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
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

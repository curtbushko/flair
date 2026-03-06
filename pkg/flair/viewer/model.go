// Package viewer provides a bubbletea-based TUI for browsing and selecting
// flair themes. It displays theme palettes, semantic tokens, and styled
// component examples with live theme switching.
package viewer

import (
	"sort"

	tea "charm.land/bubbletea/v2"
)

// Page represents which content page is displayed in the right panel.
//
// The viewer has multiple pages to showcase different aspects of a theme.
// Users can cycle through pages using the Tab key.
type Page int

const (
	// PageTextStatus shows text paragraphs and status messages.
	// Demonstrates primary/secondary/muted text and error/warning/success/info colors.
	PageTextStatus Page = iota
	// PageInteractive shows buttons, inputs, and selection lists.
	// Demonstrates interactive component styles like focused buttons and inputs.
	PageInteractive
	// PageDataDisplay shows tables, dialogs, and code blocks.
	// Demonstrates data presentation styles including syntax highlighting.
	PageDataDisplay
	// PageBubbletea shows bubbletea component examples.
	// Demonstrates spinner, progress bar, and text input styles.
	PageBubbletea
	// PageHuh shows huh form component examples.
	// Demonstrates text inputs, selects, and confirm dialogs styled with the theme.
	PageHuh
	// PageBubbles shows bubbles component examples.
	// Demonstrates list, table, and viewport components styled with the theme.
	PageBubbles
)

// pageCount is the total number of content pages for Tab cycling.
const pageCount = 6

// Options configures the viewer [Model].
//
// Options controls which themes are available, callback behavior, and
// how theme data is loaded for preview.
type Options struct {
	// Themes is the list of available theme names to display.
	// Names are sorted alphabetically when creating the Model.
	Themes []string

	// InitialTheme is the theme to pre-select on startup.
	// If empty or not found in Themes, the first theme is selected.
	InitialTheme string

	// OnSelect is called when the user confirms a theme selection with Enter.
	// If nil, selection only updates the internal state without side effects.
	OnSelect func(name string)

	// OnInstall is called when the user confirms theme selection.
	// This can be used to trigger theme installation or other actions.
	// If nil, no installation action is taken.
	OnInstall func(name string) error

	// ThemeLoader loads theme data for display in preview panels.
	// If nil, preview pages will show placeholder content.
	// Use [NewBuiltinThemeLoader] for built-in themes.
	ThemeLoader ThemeLoader
}

// ThemeLoader loads theme data for rendering in the viewer.
//
// Implementations provide palette and token data for theme preview.
// Use [NewBuiltinThemeLoader] for the standard implementation that
// loads from embedded palettes.
type ThemeLoader interface {
	// LoadPalette returns base24 colors for a theme.
	// The returned PaletteData contains 24 hex color strings.
	LoadPalette(name string) (PaletteData, error)

	// LoadTokens returns semantic tokens for a theme.
	// The returned TokenData contains tokens grouped by category.
	LoadTokens(name string) (TokenData, error)
}

// PaletteData contains the base24 colors for display in the palette preview.
type PaletteData struct {
	// Colors contains hex color strings for base00-base17 (indices 0-23).
	Colors [24]string
}

// TokenData contains semantic tokens grouped by category for preview rendering.
//
// Each map contains token path keys (e.g., "surface.background") mapped to
// hex color strings.
type TokenData struct {
	Surface    map[string]string // surface.* tokens (backgrounds)
	Text       map[string]string // text.* tokens (foregrounds)
	Status     map[string]string // status.* tokens (error, warning, etc.)
	Syntax     map[string]string // syntax.* tokens (code highlighting)
	Diff       map[string]string // diff.* tokens (version control)
	Statusline map[string]string // statusline.* tokens (status bar)
}

// Model implements tea.Model for the style viewer TUI.
//
// Model manages the viewer state including theme list navigation, page selection,
// and cached theme data for preview. Create a Model using [NewModel].
//
// Model can be used directly with bubbletea for advanced integration, or
// use the convenience functions [Run], [RunBuiltins], or [RunWithOptions].
type Model struct {
	themes        []string
	cursor        int // Current cursor position (for navigation/preview)
	selectedIndex int // Index of confirmed selection (-1 if none)
	currentPage   Page
	selectedTheme string // Name of confirmed selection
	onSelect      func(string)
	onInstall     func(string) error
	themeLoader   ThemeLoader

	// Cached data for previewed theme (at cursor position).
	palette PaletteData
	tokens  TokenData

	// Terminal dimensions.
	width  int
	height int

	// altScreen controls whether the view uses alternate screen buffer.
	altScreen bool
}

// NewModel creates a new viewer Model with the given options.
//
// The theme list is sorted alphabetically. If InitialTheme is specified
// and found in the list, it is pre-selected; otherwise the first theme
// is selected for preview.
//
// Example:
//
//	model := viewer.NewModel(viewer.Options{
//	    Themes:       flair.ListBuiltins(),
//	    ThemeLoader:  viewer.NewBuiltinThemeLoader(),
//	})
//	p := tea.NewProgram(model)
//	p.Run()
func NewModel(opts Options) Model {
	themes := opts.Themes
	if themes == nil {
		themes = []string{}
	}

	// Sort themes alphabetically.
	sort.Strings(themes)

	m := Model{
		themes:        themes,
		selectedIndex: -1, // No selection until Enter is pressed
		currentPage:   PageTextStatus,
		onSelect:      opts.OnSelect,
		onInstall:     opts.OnInstall,
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

// Init implements [tea.Model] and returns nil (no initial command).
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements [tea.Model] and handles keyboard and window events.
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
	switch msg.String() {
	case "tab":
		m.currentPage = (m.currentPage + 1) % pageCount
		return m, nil

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.previewCurrentTheme()
		}
		return m, nil

	case "down", "j":
		if m.cursor < len(m.themes)-1 {
			m.cursor++
			m.previewCurrentTheme()
		}
		return m, nil

	case "enter":
		m.confirmSelection()
		return m, nil

	case "esc", "q":
		return m, tea.Quit
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

// confirmSelection marks the current cursor position as selected and calls onSelect/onInstall.
func (m *Model) confirmSelection() {
	if len(m.themes) == 0 || m.cursor >= len(m.themes) {
		return
	}
	m.selectedIndex = m.cursor
	m.selectedTheme = m.themes[m.cursor]
	if m.onSelect != nil {
		m.onSelect(m.selectedTheme)
	}
	if m.onInstall != nil {
		// Ignore error for now - could be enhanced to show error in UI.
		_ = m.onInstall(m.selectedTheme)
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

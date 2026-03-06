package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

// NewStyles creates list.Styles from a flair.Theme.
//
// NewStyles applies theme colors to the list component's title, filter prompt,
// status bar, and pagination elements. If theme is nil, it returns list.DefaultStyles().
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	styles := list.NewStyles(theme)
//	myList.Styles = *styles
func NewStyles(theme *flair.Theme) *list.Styles {
	if theme == nil {
		s := list.DefaultStyles()
		return &s
	}

	// Get colors from theme with fallbacks.
	textPrimary := getColorHex(theme, "text.primary", "#c0caf5")
	textSecondary := getColorHex(theme, "text.secondary", "#a9b1d6")
	textMuted := getColorHex(theme, "text.muted", "#565f89")
	accentPrimary := getColorHex(theme, "accent.primary", "#7aa2f7")
	accentSecondary := getColorHex(theme, "accent.secondary", "#bb9af7")
	statusInfo := getColorHex(theme, "status.info", "#7dcfff")
	surfaceBg := getColorHex(theme, "surface.background", "#1a1b26")
	surfaceRaised := getColorHex(theme, "surface.background.raised", "#24283b")

	s := list.DefaultStyles()

	// Title styling.
	s.Title = s.Title.
		Foreground(lipgloss.Color(accentPrimary)).
		Bold(true)

	s.TitleBar = s.TitleBar.
		Foreground(lipgloss.Color(textPrimary)).
		Background(lipgloss.Color(surfaceRaised))

	// Filter styling.
	s.FilterPrompt = s.FilterPrompt.
		Foreground(lipgloss.Color(accentSecondary))

	s.FilterCursor = s.FilterCursor.
		Foreground(lipgloss.Color(accentPrimary))

	s.DefaultFilterCharacterMatch = s.DefaultFilterCharacterMatch.
		Foreground(lipgloss.Color(accentPrimary)).
		Bold(true)

	// Status bar styling.
	s.StatusBar = s.StatusBar.
		Foreground(lipgloss.Color(textSecondary)).
		Background(lipgloss.Color(surfaceBg))

	s.StatusEmpty = s.StatusEmpty.
		Foreground(lipgloss.Color(textMuted))

	s.StatusBarActiveFilter = s.StatusBarActiveFilter.
		Foreground(lipgloss.Color(accentSecondary))

	s.StatusBarFilterCount = s.StatusBarFilterCount.
		Foreground(lipgloss.Color(textMuted))

	// No items message.
	s.NoItems = s.NoItems.
		Foreground(lipgloss.Color(textMuted))

	// Pagination styling.
	s.PaginationStyle = s.PaginationStyle.
		Foreground(lipgloss.Color(textMuted))

	s.ActivePaginationDot = s.ActivePaginationDot.
		Foreground(lipgloss.Color(accentPrimary))

	s.InactivePaginationDot = s.InactivePaginationDot.
		Foreground(lipgloss.Color(textMuted))

	s.ArabicPagination = s.ArabicPagination.
		Foreground(lipgloss.Color(textSecondary))

	s.DividerDot = s.DividerDot.
		Foreground(lipgloss.Color(textMuted))

	// Help styling.
	s.HelpStyle = s.HelpStyle.
		Foreground(lipgloss.Color(textMuted))

	// Spinner styling.
	s.Spinner = s.Spinner.
		Foreground(lipgloss.Color(statusInfo))

	return &s
}

// NewDelegate creates a themed list.DefaultDelegate from a flair.Theme.
//
// NewDelegate configures the delegate's item styles with theme colors,
// styling normal items, selected items, dimmed items, and filter matches.
// If theme is nil, it returns a default delegate.
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	delegate := list.NewDelegate(theme)
//	myList := bubbles_list.New(items, *delegate, width, height)
func NewDelegate(theme *flair.Theme) *list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()

	if theme == nil {
		return &delegate
	}

	delegate.Styles = NewItemStyles(theme)
	return &delegate
}

// NewItemStyles creates list.DefaultItemStyles from a flair.Theme.
//
// NewItemStyles applies theme colors to item title, description, and filter match
// styles for normal, selected, and dimmed states. If theme is nil, it returns
// list.NewDefaultItemStyles().
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	itemStyles := list.NewItemStyles(theme)
//	delegate.Styles = itemStyles
func NewItemStyles(theme *flair.Theme) list.DefaultItemStyles {
	if theme == nil {
		return list.NewDefaultItemStyles()
	}

	// Get colors from theme with fallbacks.
	textPrimary := getColorHex(theme, "text.primary", "#c0caf5")
	textSecondary := getColorHex(theme, "text.secondary", "#a9b1d6")
	textMuted := getColorHex(theme, "text.muted", "#565f89")
	accentPrimary := getColorHex(theme, "accent.primary", "#7aa2f7")
	selectionBg := getColorHex(theme, "surface.background.selection", "#33467c")

	s := list.NewDefaultItemStyles()

	// Normal state.
	s.NormalTitle = s.NormalTitle.
		Foreground(lipgloss.Color(textPrimary))

	s.NormalDesc = s.NormalDesc.
		Foreground(lipgloss.Color(textSecondary))

	// Selected state - accent color with selection background.
	s.SelectedTitle = s.SelectedTitle.
		Foreground(lipgloss.Color(accentPrimary)).
		Background(lipgloss.Color(selectionBg)).
		Bold(true)

	s.SelectedDesc = s.SelectedDesc.
		Foreground(lipgloss.Color(textSecondary)).
		Background(lipgloss.Color(selectionBg))

	// Dimmed state - muted colors for filtering.
	s.DimmedTitle = s.DimmedTitle.
		Foreground(lipgloss.Color(textMuted))

	s.DimmedDesc = s.DimmedDesc.
		Foreground(lipgloss.Color(textMuted))

	// Filter match highlighting.
	s.FilterMatch = s.FilterMatch.
		Foreground(lipgloss.Color(accentPrimary)).
		Bold(true)

	return s
}

// getColorHex retrieves a color from the theme by path, returning the hex string.
// If the color is not found, it returns the fallback value.
func getColorHex(theme *flair.Theme, path, fallback string) string {
	if c, ok := theme.Color(path); ok {
		return c.Hex()
	}
	return fallback
}

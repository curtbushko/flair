package lipgloss

import (
	"github.com/curtbushko/flair/pkg/flair"
)

// NewStyles creates a Styles struct from a flair theme.
// It builds pre-configured lipgloss.Style instances for all style categories.
func NewStyles(theme *flair.Theme) *Styles {
	if theme == nil {
		return nil
	}

	s := &Styles{}

	// Surface styles (via helper functions)
	s.Background = BuildSurfaceBackground(theme)
	s.Raised = BuildSurfaceRaised(theme)
	s.Sunken = BuildSurfaceSunken(theme)
	s.Overlay = BuildSurfaceOverlay(theme)
	s.Popup = BuildSurfacePopup(theme)

	// Text styles (via helper functions)
	s.Text = BuildTextPrimary(theme)
	s.Secondary = BuildTextSecondary(theme)
	s.Muted = BuildTextMuted(theme)
	s.Inverse = BuildTextInverse(theme)

	// Status styles (via helper functions)
	s.Error = BuildStatusError(theme)
	s.Warning = BuildStatusWarning(theme)
	s.Success = BuildStatusSuccess(theme)
	s.Info = BuildStatusInfo(theme)

	// Border styles (via helper functions)
	s.Border = BuildBorderDefault(theme)
	s.BorderFocus = BuildBorderFocus(theme)
	s.BorderMuted = BuildBorderMuted(theme)

	// Component styles (via helper functions)
	s.Button = BuildButton(theme)
	s.ButtonFocused = BuildButtonFocused(theme)
	s.Input = BuildInput(theme)
	s.InputFocused = BuildInputFocused(theme)
	s.ListItem = BuildListItem(theme)
	s.ListSelected = BuildListSelected(theme)
	s.Table = BuildTable(theme)
	s.TableHeader = BuildTableHeader(theme)
	s.Dialog = BuildDialog(theme)

	// State styles (via helper functions)
	s.Hover = BuildStateHover(theme)
	s.Active = BuildStateActive(theme)
	s.Disabled = BuildStateDisabled(theme)
	s.Selected = BuildStateSelected(theme)

	return s
}

// Default returns styles using the currently selected flair theme.
// Returns nil if no theme is selected or loading fails.
func Default() *Styles {
	theme, err := flair.Load()
	if err != nil {
		return nil
	}
	return NewStyles(theme)
}

// DefaultFrom returns styles from the specified config directory.
// Returns nil if no theme is selected or loading fails.
func DefaultFrom(configDir string) *Styles {
	theme, err := flair.LoadFrom(configDir)
	if err != nil {
		return nil
	}
	return NewStyles(theme)
}

// getColor retrieves a color from the theme by path, with a fallback default.
func getColor(theme *flair.Theme, path string, fallback flair.Color) flair.Color {
	if c, ok := theme.Get(path); ok {
		return c
	}
	return fallback
}

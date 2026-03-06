package huh

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/flair"
)

const (
	buttonPaddingHorizontal = 2
	buttonPaddingVertical   = 0
)

// NewTheme creates a [huh.Theme] from a [flair.Theme].
//
// The returned theme applies flair's semantic colors to all huh form components
// including titles, descriptions, errors, buttons, inputs, and selectors.
//
// NewTheme returns nil if the provided theme is nil.
//
// Example:
//
//	flairTheme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	huhTheme := huh.NewTheme(flairTheme)
//	form := huh.NewForm(...).WithTheme(huhTheme)
func NewTheme(theme *flair.Theme) *huh.Theme {
	if theme == nil {
		return nil
	}

	// Extract semantic colors from flair theme.
	text := theme.Text()
	status := theme.Status()

	// Get accent and border colors with fallbacks.
	accentPrimary := getColorHex(theme, "accent.primary", "#7aa2f7")
	accentSecondary := getColorHex(theme, "accent.secondary", "#bb9af7")
	borderDefault := getColorHex(theme, "border.default", "#565f89")
	borderFocus := getColorHex(theme, "border.focus", "#7aa2f7")

	// Start with a base theme.
	t := buildBaseTheme()

	// Apply flair colors to focused styles.
	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color(borderFocus))
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color(accentPrimary))
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lipgloss.Color(accentPrimary)).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(lipgloss.Color(accentPrimary))
	t.Focused.File = t.Focused.File.Foreground(lipgloss.Color(text.Primary.Hex()))
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.Color(text.Secondary.Hex()))
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(lipgloss.Color(status.Error.Hex()))
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(lipgloss.Color(status.Error.Hex()))

	// Selector styles.
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color(accentSecondary))
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(lipgloss.Color(accentSecondary))
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(lipgloss.Color(accentSecondary))
	t.Focused.Option = t.Focused.Option.Foreground(lipgloss.Color(text.Primary.Hex()))

	// Multi-select styles.
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lipgloss.Color(accentSecondary))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color(status.Success.Hex()))
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color(status.Success.Hex()))
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(lipgloss.Color(text.Primary.Hex()))
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(lipgloss.Color(text.Muted.Hex()))

	// Button styles.
	t.Focused.FocusedButton = t.Focused.FocusedButton.
		Foreground(lipgloss.Color(text.Inverse.Hex())).
		Background(lipgloss.Color(accentPrimary))
	t.Focused.BlurredButton = t.Focused.BlurredButton.
		Foreground(lipgloss.Color(text.Primary.Hex())).
		Background(lipgloss.Color(borderDefault))
	t.Focused.Next = t.Focused.FocusedButton

	// Text input styles.
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(lipgloss.Color(accentPrimary))
	t.Focused.TextInput.CursorText = t.Focused.TextInput.CursorText.Foreground(lipgloss.Color(text.Primary.Hex()))
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.Color(text.Muted.Hex()))
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(lipgloss.Color(accentSecondary))
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(lipgloss.Color(text.Primary.Hex()))

	// Apply blurred styles (copy from focused with modifications).
	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.Title = t.Blurred.Title.Foreground(lipgloss.Color(text.Muted.Hex()))
	t.Blurred.Description = t.Blurred.Description.Foreground(lipgloss.Color(text.Muted.Hex()))
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	// Group styles.
	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	// Help styles.
	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(lipgloss.Color(text.Muted.Hex()))
	t.Help.ShortKey = t.Help.ShortKey.Foreground(lipgloss.Color(text.Secondary.Hex()))
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(lipgloss.Color(text.Muted.Hex()))
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(lipgloss.Color(text.Muted.Hex()))
	t.Help.FullKey = t.Help.FullKey.Foreground(lipgloss.Color(text.Secondary.Hex()))
	t.Help.FullDesc = t.Help.FullDesc.Foreground(lipgloss.Color(text.Muted.Hex()))
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(lipgloss.Color(text.Muted.Hex()))

	return t
}

// Default loads a [huh.Theme] from the currently selected flair theme.
//
// Default uses [flair.Default] to load the theme, which falls back to the
// built-in default theme (tokyo-night-dark) if no theme is selected.
//
// Example:
//
//	theme, err := huh.Default()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	form := huh.NewForm(...).WithTheme(theme)
func Default() (*huh.Theme, error) {
	theme, err := flair.Default()
	if err != nil {
		return nil, err
	}
	return NewTheme(theme), nil
}

// DefaultFrom loads a [huh.Theme] from the specified config directory.
//
// DefaultFrom is useful for testing or when using a non-standard config
// location. It loads the currently selected theme from the config directory.
//
// Returns an error if no theme is selected or the theme cannot be loaded.
func DefaultFrom(configDir string) (*huh.Theme, error) {
	theme, err := flair.LoadFrom(configDir)
	if err != nil {
		return nil, err
	}
	return NewTheme(theme), nil
}

// buildBaseTheme creates a base huh theme with structural styles.
func buildBaseTheme() *huh.Theme {
	var t huh.Theme

	t.Form.Base = lipgloss.NewStyle()
	t.Group.Base = lipgloss.NewStyle()
	t.FieldSeparator = lipgloss.NewStyle().SetString("\n\n")

	button := lipgloss.NewStyle().
		Padding(buttonPaddingVertical, buttonPaddingHorizontal).
		MarginRight(1)

	// Focused styles.
	t.Focused.Base = lipgloss.NewStyle().
		PaddingLeft(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeft(true)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = lipgloss.NewStyle()
	t.Focused.NoteTitle = lipgloss.NewStyle()
	t.Focused.Description = lipgloss.NewStyle()
	t.Focused.Directory = lipgloss.NewStyle()
	t.Focused.File = lipgloss.NewStyle()
	t.Focused.ErrorIndicator = lipgloss.NewStyle().SetString(" *")
	t.Focused.ErrorMessage = lipgloss.NewStyle()
	t.Focused.SelectSelector = lipgloss.NewStyle().SetString("> ")
	t.Focused.NextIndicator = lipgloss.NewStyle().MarginLeft(1).SetString("->")
	t.Focused.PrevIndicator = lipgloss.NewStyle().MarginRight(1).SetString("<-")
	t.Focused.Option = lipgloss.NewStyle()
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().SetString("> ")
	t.Focused.SelectedOption = lipgloss.NewStyle()
	t.Focused.SelectedPrefix = lipgloss.NewStyle().SetString("[x] ")
	t.Focused.UnselectedOption = lipgloss.NewStyle()
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ")
	t.Focused.FocusedButton = button
	t.Focused.BlurredButton = button
	t.Focused.Next = button
	t.Focused.TextInput.Cursor = lipgloss.NewStyle()
	t.Focused.TextInput.CursorText = lipgloss.NewStyle()
	t.Focused.TextInput.Placeholder = lipgloss.NewStyle()
	t.Focused.TextInput.Prompt = lipgloss.NewStyle()
	t.Focused.TextInput.Text = lipgloss.NewStyle()

	t.Help = help.New().Styles

	// Blurred styles (copy from focused, hide border).
	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return &t
}

// getColorHex retrieves a color from the theme by path, returning the hex string.
// If the color is not found, it returns the fallback value.
func getColorHex(theme *flair.Theme, path, fallback string) string {
	if c, ok := theme.Color(path); ok {
		return c.Hex()
	}
	return fallback
}

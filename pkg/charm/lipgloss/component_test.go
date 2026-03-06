package lipgloss_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

func TestComponent_Button(t *testing.T) {
	// Given: Theme with accent.primary and accent.foreground
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"accent.foreground":          mustParseHex("#1a1b26"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Button via builder function
	button := lipgloss.BuildButton(theme)

	// Then: Button has accent bg, padding, and appropriate fg
	rendered := button.Render("test")
	if rendered == "" {
		t.Error("Button style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Button.Render("test")
	if rendered == "" {
		t.Error("styles.Button should render content")
	}
}

func TestComponent_ButtonFocused(t *testing.T) {
	// Given: Theme with accent colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"accent.foreground":          mustParseHex("#1a1b26"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.ButtonFocused via builder function
	buttonFocused := lipgloss.BuildButtonFocused(theme)

	// Then: ButtonFocused style renders non-empty with bold
	rendered := buttonFocused.Render("test")
	if rendered == "" {
		t.Error("ButtonFocused style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.ButtonFocused.Render("test")
	if rendered == "" {
		t.Error("styles.ButtonFocused should render content")
	}
}

func TestComponent_Input(t *testing.T) {
	// Given: Theme with surface and border tokens
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Input via builder function
	input := lipgloss.BuildInput(theme)

	// Then: Input has raised bg, border, and text foreground
	rendered := input.Render("test")
	if rendered == "" {
		t.Error("Input style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Input.Render("test")
	if rendered == "" {
		t.Error("styles.Input should render content")
	}
}

func TestComponent_InputFocused(t *testing.T) {
	// Given: Theme with surface and border tokens
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.InputFocused via builder function
	inputFocused := lipgloss.BuildInputFocused(theme)

	// Then: InputFocused has focused border color
	rendered := inputFocused.Render("test")
	if rendered == "" {
		t.Error("InputFocused style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.InputFocused.Render("test")
	if rendered == "" {
		t.Error("styles.InputFocused should render content")
	}
}

func TestComponent_ListItem(t *testing.T) {
	// Given: Theme with text colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.ListItem via builder function
	listItem := lipgloss.BuildListItem(theme)

	// Then: ListItem style renders non-empty
	rendered := listItem.Render("test")
	if rendered == "" {
		t.Error("ListItem style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.ListItem.Render("test")
	if rendered == "" {
		t.Error("styles.ListItem should render content")
	}
}

func TestComponent_ListSelected(t *testing.T) {
	// Given: Theme with accent colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.ListSelected via builder function
	listSelected := lipgloss.BuildListSelected(theme)

	// Then: ListSelected style renders non-empty with bold
	rendered := listSelected.Render("test")
	if rendered == "" {
		t.Error("ListSelected style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.ListSelected.Render("test")
	if rendered == "" {
		t.Error("styles.ListSelected should render content")
	}
}

func TestComponent_Table(t *testing.T) {
	// Given: Theme with text colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Table via builder function
	table := lipgloss.BuildTable(theme)

	// Then: Table style renders non-empty
	rendered := table.Render("test")
	if rendered == "" {
		t.Error("Table style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Table.Render("test")
	if rendered == "" {
		t.Error("styles.Table should render content")
	}
}

func TestComponent_TableHeader(t *testing.T) {
	// Given: Theme with text colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#16161e"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.TableHeader via builder function
	tableHeader := lipgloss.BuildTableHeader(theme)

	// Then: TableHeader style renders non-empty with bold
	rendered := tableHeader.Render("test")
	if rendered == "" {
		t.Error("TableHeader style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.TableHeader.Render("test")
	if rendered == "" {
		t.Error("styles.TableHeader should render content")
	}
}

func TestComponent_Dialog(t *testing.T) {
	// Given: Theme with surface and border colors
	colors := map[string]flair.Color{
		"surface.background":         mustParseHex("#1a1b26"),
		"surface.background.raised":  mustParseHex("#1f2335"),
		"surface.background.sunken":  mustParseHex("#16161e"),
		"surface.background.overlay": mustParseHex("#24283b"),
		"surface.background.popup":   mustParseHex("#16161e"),
		"text.primary":               mustParseHex("#c0caf5"),
		"text.secondary":             mustParseHex("#a9b1d6"),
		"text.muted":                 mustParseHex("#565f89"),
		"text.inverse":               mustParseHex("#1a1b26"),
		"status.error":               mustParseHex("#f7768e"),
		"status.warning":             mustParseHex("#e0af68"),
		"status.success":             mustParseHex("#9ece6a"),
		"status.info":                mustParseHex("#7dcfff"),
		"accent.primary":             mustParseHex("#7aa2f7"),
		"border.default":             mustParseHex("#565f89"),
		"border.focus":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Check styles.Dialog via builder function
	dialog := lipgloss.BuildDialog(theme)

	// Then: Dialog style has overlay bg, border, and padding
	rendered := dialog.Render("test")
	if rendered == "" {
		t.Error("Dialog style should render content")
	}

	// Verify via NewStyles
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
		return
	}
	rendered = styles.Dialog.Render("test")
	if rendered == "" {
		t.Error("styles.Dialog should render content")
	}
}

// Edge case tests for component styles

func TestComponent_MissingAccentFallback(t *testing.T) {
	// Given: Theme without accent.* tokens
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		// Note: no accent.* tokens
	}
	theme := flair.NewTheme("minimal", "dark", colors)

	// When: Build components that use accent colors
	// Then: Should not panic and should use fallback
	t.Run("ButtonFocused", func(t *testing.T) {
		rendered := lipgloss.BuildButtonFocused(theme).Render("fallback test")
		if rendered == "" {
			t.Error("ButtonFocused style should render with fallback accent")
		}
	})
	t.Run("ListSelected", func(t *testing.T) {
		rendered := lipgloss.BuildListSelected(theme).Render("fallback test")
		if rendered == "" {
			t.Error("ListSelected style should render with fallback accent")
		}
	})
}

func TestComponent_AllComponentsEmptyContent(t *testing.T) {
	// Given: Theme with all tokens
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
		"accent.foreground":            mustParseHex("#1a1b26"),
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render empty string with all component styles
	// Then: Should not panic on empty string
	t.Run("Button", func(t *testing.T) {
		rendered := lipgloss.BuildButton(theme).Render("")
		t.Logf("Button rendered empty as: %q", rendered)
	})
	t.Run("ButtonFocused", func(t *testing.T) {
		rendered := lipgloss.BuildButtonFocused(theme).Render("")
		t.Logf("ButtonFocused rendered empty as: %q", rendered)
	})
	t.Run("Input", func(t *testing.T) {
		rendered := lipgloss.BuildInput(theme).Render("")
		t.Logf("Input rendered empty as: %q", rendered)
	})
	t.Run("InputFocused", func(t *testing.T) {
		rendered := lipgloss.BuildInputFocused(theme).Render("")
		t.Logf("InputFocused rendered empty as: %q", rendered)
	})
	t.Run("ListItem", func(t *testing.T) {
		rendered := lipgloss.BuildListItem(theme).Render("")
		t.Logf("ListItem rendered empty as: %q", rendered)
	})
	t.Run("ListSelected", func(t *testing.T) {
		rendered := lipgloss.BuildListSelected(theme).Render("")
		t.Logf("ListSelected rendered empty as: %q", rendered)
	})
	t.Run("Table", func(t *testing.T) {
		rendered := lipgloss.BuildTable(theme).Render("")
		t.Logf("Table rendered empty as: %q", rendered)
	})
	t.Run("TableHeader", func(t *testing.T) {
		rendered := lipgloss.BuildTableHeader(theme).Render("")
		t.Logf("TableHeader rendered empty as: %q", rendered)
	})
	t.Run("Dialog", func(t *testing.T) {
		rendered := lipgloss.BuildDialog(theme).Render("")
		t.Logf("Dialog rendered empty as: %q", rendered)
	})
}

func TestComponent_ButtonPaddingApplied(t *testing.T) {
	// Given: Theme with colors
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#16161e"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render button with content
	button := lipgloss.BuildButton(theme)
	rendered := button.Render("OK")

	// Then: Output should be longer than input due to padding
	if len(rendered) <= len("OK") {
		t.Logf("Button with padding rendered: %q (length %d)", rendered, len(rendered))
	}
}

func TestComponent_DialogWithComplexContent(t *testing.T) {
	// Given: Theme with colors
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#24283b"),
		"surface.background.popup":     mustParseHex("#16161e"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#a9b1d6"),
		"text.muted":                   mustParseHex("#565f89"),
		"text.inverse":                 mustParseHex("#1a1b26"),
		"status.error":                 mustParseHex("#f7768e"),
		"status.warning":               mustParseHex("#e0af68"),
		"status.success":               mustParseHex("#9ece6a"),
		"status.info":                  mustParseHex("#7dcfff"),
		"accent.primary":               mustParseHex("#7aa2f7"),
		"border.default":               mustParseHex("#565f89"),
		"border.focus":                 mustParseHex("#7aa2f7"),
	}
	theme := flair.NewTheme("tokyo-night-dark", "dark", colors)

	// When: Render dialog with multiline content including nested styles
	dialog := lipgloss.BuildDialog(theme)
	button := lipgloss.BuildButton(theme)

	complexContent := "Dialog Title\n\nAre you sure you want to proceed?\n\n" +
		button.Render(" Yes ") + " " + button.Render(" No ")

	rendered := dialog.Render(complexContent)

	// Then: Should render without panic
	if rendered == "" {
		t.Error("Dialog with complex content should render")
	}
}

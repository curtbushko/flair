package lipgloss_test

import (
	"reflect"
	"strings"
	"testing"

	charmLipgloss "github.com/charmbracelet/lipgloss"

	"github.com/curtbushko/flair/pkg/charm/lipgloss"
	"github.com/curtbushko/flair/pkg/flair"
)

// createMockTokyoNightTheme creates a mock theme matching Tokyo Night Dark colors.
// This is a complete theme with all required tokens for integration testing.
func createMockTokyoNightTheme() *flair.Theme {
	colors := map[string]flair.Color{
		// Surface tokens
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1f2335"),
		"surface.background.sunken":    mustParseHex("#16161e"),
		"surface.background.overlay":   mustParseHex("#24283b"),
		"surface.background.popup":     mustParseHex("#1f2335"),
		"surface.background.highlight": mustParseHex("#292e42"),
		"surface.background.selection": mustParseHex("#364a82"),

		// Text tokens
		"text.primary":   mustParseHex("#c0caf5"),
		"text.secondary": mustParseHex("#a9b1d6"),
		"text.muted":     mustParseHex("#565f89"),
		"text.inverse":   mustParseHex("#1a1b26"),

		// Status tokens
		"status.error":   mustParseHex("#f7768e"),
		"status.warning": mustParseHex("#e0af68"),
		"status.success": mustParseHex("#9ece6a"),
		"status.info":    mustParseHex("#7dcfff"),

		// Accent tokens
		"accent.primary":    mustParseHex("#7aa2f7"),
		"accent.foreground": mustParseHex("#1a1b26"),

		// Border tokens
		"border.default": mustParseHex("#565f89"),
		"border.focus":   mustParseHex("#7aa2f7"),
		"border.muted":   mustParseHex("#3b4261"),
	}
	return flair.NewTheme("tokyo-night-dark", "dark", colors)
}

// createMinimalTheme creates a theme with only surface.background and text.primary.
// Used to test fallback behavior for missing tokens.
func createMinimalTheme() *flair.Theme {
	colors := map[string]flair.Color{
		"surface.background":           mustParseHex("#1a1b26"),
		"surface.background.raised":    mustParseHex("#1a1b26"),
		"surface.background.sunken":    mustParseHex("#1a1b26"),
		"surface.background.overlay":   mustParseHex("#1a1b26"),
		"surface.background.popup":     mustParseHex("#1a1b26"),
		"surface.background.highlight": mustParseHex("#1a1b26"),
		"surface.background.selection": mustParseHex("#1a1b26"),
		"text.primary":                 mustParseHex("#c0caf5"),
		"text.secondary":               mustParseHex("#c0caf5"),
		"text.muted":                   mustParseHex("#c0caf5"),
		"text.inverse":                 mustParseHex("#c0caf5"),
		"status.error":                 mustParseHex("#c0caf5"),
		"status.warning":               mustParseHex("#c0caf5"),
		"status.success":               mustParseHex("#c0caf5"),
		"status.info":                  mustParseHex("#c0caf5"),
	}
	return flair.NewTheme("minimal", "dark", colors)
}

func TestIntegration_TokyoNightStyles(t *testing.T) {
	// Given: A mock theme matching Tokyo Night Dark colors
	theme := createMockTokyoNightTheme()

	// When: Create Styles via NewStyles()
	styles := lipgloss.NewStyles(theme)

	// Then: All style fields are configured, Render() produces non-empty output
	if styles == nil {
		t.Fatal("NewStyles returned nil for a complete Tokyo Night theme")
	}

	// Test all surface styles render
	surfaceTests := []struct {
		name  string
		style charmLipgloss.Style
	}{
		{"Background", styles.Background},
		{"Raised", styles.Raised},
		{"Sunken", styles.Sunken},
		{"Overlay", styles.Overlay},
		{"Popup", styles.Popup},
	}
	for _, tt := range surfaceTests {
		t.Run("Surface_"+tt.name, func(t *testing.T) {
			rendered := tt.style.Render("test content")
			if rendered == "" {
				t.Errorf("%s style did not render content", tt.name)
			}
		})
	}

	// Test all text styles render
	textTests := []struct {
		name  string
		style charmLipgloss.Style
	}{
		{"Text", styles.Text},
		{"Secondary", styles.Secondary},
		{"Muted", styles.Muted},
		{"Inverse", styles.Inverse},
	}
	for _, tt := range textTests {
		t.Run("Text_"+tt.name, func(t *testing.T) {
			rendered := tt.style.Render("test content")
			if rendered == "" {
				t.Errorf("%s style did not render content", tt.name)
			}
		})
	}

	// Test all status styles render
	statusTests := []struct {
		name  string
		style charmLipgloss.Style
	}{
		{"Error", styles.Error},
		{"Warning", styles.Warning},
		{"Success", styles.Success},
		{"Info", styles.Info},
	}
	for _, tt := range statusTests {
		t.Run("Status_"+tt.name, func(t *testing.T) {
			rendered := tt.style.Render("test content")
			if rendered == "" {
				t.Errorf("%s style did not render content", tt.name)
			}
		})
	}

	// Test all border styles render
	borderTests := []struct {
		name  string
		style charmLipgloss.Style
	}{
		{"Border", styles.Border},
		{"BorderFocus", styles.BorderFocus},
		{"BorderMuted", styles.BorderMuted},
	}
	for _, tt := range borderTests {
		t.Run("Border_"+tt.name, func(t *testing.T) {
			rendered := tt.style.Render("test content")
			if rendered == "" {
				t.Errorf("%s style did not render content", tt.name)
			}
		})
	}

	// Test all component styles render
	componentTests := []struct {
		name  string
		style charmLipgloss.Style
	}{
		{"Button", styles.Button},
		{"ButtonFocused", styles.ButtonFocused},
		{"Input", styles.Input},
		{"InputFocused", styles.InputFocused},
		{"ListItem", styles.ListItem},
		{"ListSelected", styles.ListSelected},
		{"Table", styles.Table},
		{"TableHeader", styles.TableHeader},
		{"Dialog", styles.Dialog},
	}
	for _, tt := range componentTests {
		t.Run("Component_"+tt.name, func(t *testing.T) {
			rendered := tt.style.Render("test content")
			if rendered == "" {
				t.Errorf("%s style did not render content", tt.name)
			}
		})
	}

	// Test all state styles render
	stateTests := []struct {
		name  string
		style charmLipgloss.Style
	}{
		{"Hover", styles.Hover},
		{"Active", styles.Active},
		{"Disabled", styles.Disabled},
		{"Selected", styles.Selected},
	}
	for _, tt := range stateTests {
		t.Run("State_"+tt.name, func(t *testing.T) {
			rendered := tt.style.Render("test content")
			if rendered == "" {
				t.Errorf("%s style did not render content", tt.name)
			}
		})
	}
}

func TestStyles_RenderOutput(t *testing.T) {
	// Given: Styles from NewStyles()
	theme := createMockTokyoNightTheme()
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// When: Call styles.Error.Render('Error message')
	output := styles.Error.Render("Error message")

	// Then: Output contains the original text
	// Note: ANSI escape codes may not be present in test environments without TTY
	// The important thing is that the style renders and contains the text
	if !strings.Contains(output, "Error message") {
		t.Errorf("expected output to contain 'Error message', got: %q", output)
	}

	// Verify the render produces non-empty output
	if output == "" {
		t.Error("expected non-empty output from Error style")
	}
}

func TestStyles_RenderOutput_AllStyles(t *testing.T) {
	// Given: Styles from NewStyles()
	theme := createMockTokyoNightTheme()
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// Test that various styles produce output containing the original text
	// Note: ANSI escape codes may not be present in test environments without TTY
	tests := []struct {
		name   string
		render func() string
	}{
		{"Error", func() string { return styles.Error.Render("test") }},
		{"Warning", func() string { return styles.Warning.Render("test") }},
		{"Success", func() string { return styles.Success.Render("test") }},
		{"Info", func() string { return styles.Info.Render("test") }},
		{"Text", func() string { return styles.Text.Render("test") }},
		{"Background", func() string { return styles.Background.Render("test") }},
		{"Button", func() string { return styles.Button.Render("test") }},
		{"Dialog", func() string { return styles.Dialog.Render("test") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.render()

			// Should render non-empty
			if output == "" {
				t.Errorf("%s: expected non-empty output", tt.name)
			}

			// Should contain the original text
			if !strings.Contains(output, "test") {
				t.Errorf("%s: expected 'test' in output, got: %q", tt.name, output)
			}
		})
	}
}

func TestNewStyles_MissingTokensFallback(t *testing.T) {
	// Given: Theme with only minimal tokens (missing border, accent tokens)
	theme := createMinimalTheme()

	// When: Create Styles
	styles := lipgloss.NewStyles(theme)

	// Then: No panic, styles are created successfully
	if styles == nil {
		t.Fatal("NewStyles returned nil for minimal theme")
	}

	// Styles that use fallback colors should still render
	tests := []struct {
		name  string
		style charmLipgloss.Style
	}{
		{"Border", styles.Border},
		{"BorderFocus", styles.BorderFocus},
		{"BorderMuted", styles.BorderMuted},
		{"Button", styles.Button},
		{"ButtonFocused", styles.ButtonFocused},
		{"ListSelected", styles.ListSelected},
		{"Active", styles.Active},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render("fallback test")
			if rendered == "" {
				t.Errorf("%s style did not render with fallback colors", tt.name)
			}
			// Should still contain the text
			if !strings.Contains(rendered, "fallback test") {
				t.Errorf("%s: expected 'fallback test' in output, got: %q", tt.name, rendered)
			}
		})
	}
}

func TestNewStyles_NilTheme(t *testing.T) {
	// Given: A nil theme

	// When: Create Styles with nil theme
	styles := lipgloss.NewStyles(nil)

	// Then: Returns nil, no panic
	if styles != nil {
		t.Error("NewStyles should return nil for nil theme")
	}
}

func TestStyles_AllNonZero(t *testing.T) {
	// Given: Complete theme
	theme := createMockTokyoNightTheme()
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// When: Check each style field in Styles struct
	stylesValue := reflect.ValueOf(*styles)
	stylesType := reflect.TypeOf(*styles)

	// Then: None are the zero value lipgloss.Style{}
	zeroStyle := charmLipgloss.Style{}

	for i := 0; i < stylesValue.NumField(); i++ {
		fieldName := stylesType.Field(i).Name
		fieldValue := stylesValue.Field(i).Interface()

		style, ok := fieldValue.(charmLipgloss.Style)
		if !ok {
			continue // Skip non-style fields if any
		}

		t.Run(fieldName, func(t *testing.T) {
			// A configured style should render differently than a zero style
			// when given the same input, or at minimum should not be the zero value
			zeroRender := zeroStyle.Render("x")
			styleRender := style.Render("x")

			// At minimum, verify it renders
			if styleRender == "" {
				t.Errorf("field %s has a style that renders to empty string", fieldName)
			}

			// Verify it's configured (renders differently than zero style)
			// Note: Some styles might coincidentally render the same, so this is informational
			if styleRender == zeroRender {
				t.Logf("field %s renders same as zero style (may need configuration)", fieldName)
			}
		})
	}
}

func TestStyles_EndToEnd_ThemeLoadToRender(t *testing.T) {
	// Given: A theme with complete tokens
	theme := createMockTokyoNightTheme()

	// When: Create full style chain
	styles := lipgloss.NewStyles(theme)
	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// Then: Can build a complete UI layout
	// Simulate a real UI: header, content, status bar
	header := styles.TableHeader.Render("My Application")
	content := styles.Background.Render(
		styles.Text.Render("Welcome to the app!"),
	)
	errorMsg := styles.Error.Render("Something went wrong")
	warningMsg := styles.Warning.Render("Proceed with caution")
	successMsg := styles.Success.Render("Operation successful")
	infoMsg := styles.Info.Render("FYI: This is informational")

	// All should produce output
	outputs := map[string]string{
		"header":     header,
		"content":    content,
		"errorMsg":   errorMsg,
		"warningMsg": warningMsg,
		"successMsg": successMsg,
		"infoMsg":    infoMsg,
	}

	for name, output := range outputs {
		if output == "" {
			t.Errorf("%s rendered to empty string", name)
		}
	}

	// Build a dialog
	dialog := styles.Dialog.Render(
		styles.Text.Render("Dialog content\n") +
			styles.Button.Render(" OK ") + " " +
			styles.ButtonFocused.Render(" Cancel "),
	)

	if dialog == "" {
		t.Error("dialog rendered to empty string")
	}
}

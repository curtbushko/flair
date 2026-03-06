// Package lipgloss provides themed lipgloss styles built from flair themes.
//
// This package creates pre-configured [lipgloss.Style] instances for common
// TUI elements including surfaces, text, status indicators, borders, and
// interactive components. It reads theme data from pkg/flair and translates
// semantic color tokens into ready-to-use lipgloss styles.
//
// The package is fully independent from flair's internal packages.
// External projects can import it without pulling in flair's internal
// implementation details.
//
// # Quick Start
//
// The simplest way to get themed styles is using [Default]:
//
//	styles := lipgloss.Default()
//	if styles != nil {
//	    header := styles.Raised.Render("My Application")
//	    message := styles.Text.Render("Hello, world!")
//	    warning := styles.Warning.Render("This is a warning")
//	    fmt.Println(header)
//	    fmt.Println(message)
//	    fmt.Println(warning)
//	}
//
// # Using a Specific Theme
//
// Load a specific theme and create styles with [NewStyles]:
//
//	theme, err := flair.LoadBuiltin("gruvbox-dark")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	styles := lipgloss.NewStyles(theme)
//
// # Style Categories
//
// The [Styles] struct provides pre-built styles organized by category:
//
//   - Surface: Background, Raised, Sunken, Overlay, Popup
//   - Text: Text (primary), Secondary, Muted, Inverse
//   - Status: Error, Warning, Success, Info
//   - Border: Border, BorderFocus, BorderMuted
//   - Components: Button, ButtonFocused, Input, InputFocused, ListItem,
//     ListSelected, Table, TableHeader, Dialog
//   - State: Hover, Active, Disabled, Selected
//
// # Builder Functions
//
// For fine-grained control, use the individual Build* functions to create
// specific styles:
//
//	theme, _ := flair.Default()
//	errorStyle := lipgloss.BuildStatusError(theme)
//	buttonStyle := lipgloss.BuildButton(theme)
//
// # Dependencies
//
// This package depends on:
//   - github.com/charmbracelet/lipgloss
//   - github.com/curtbushko/flair/pkg/flair
package lipgloss

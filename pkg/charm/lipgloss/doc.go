// Package lipgloss provides themed lipgloss styles built from flair themes.
//
// This package creates pre-configured lipgloss.Style instances for common
// TUI elements including surfaces, text, status indicators, borders, and
// interactive components. It reads theme data from pkg/flair and translates
// semantic color tokens into ready-to-use lipgloss styles.
//
// The package is fully independent from flair's internal packages.
// External projects can import it without pulling in flair's internal
// implementation details.
//
// # Basic Usage
//
// Load the currently selected flair theme and create styles:
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
// Load a named theme and create styles:
//
//	theme, err := flair.LoadNamed("tokyo-night-dark")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	styles := lipgloss.NewStyles(theme)
//
// # Style Categories
//
// The Styles struct provides styles for:
//
//   - Surface: Background, Raised, Sunken, Overlay, Popup
//   - Text: Text (primary), Secondary, Muted, Inverse
//   - Status: Error, Warning, Success, Info
//   - Border: Border, BorderFocus
//   - Components: Button, ButtonFocused, Input, InputFocused, ListItem,
//     ListSelected, Table, TableHeader, Dialog
//
// # Dependencies
//
// This package depends on:
//   - github.com/charmbracelet/lipgloss
//   - github.com/curtbushko/flair/pkg/flair
package lipgloss

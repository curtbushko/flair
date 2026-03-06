// Package huh provides themed huh form components built from flair themes.
//
// This package creates [github.com/charmbracelet/huh.Theme] instances configured
// with colors from a flair theme. It enables building consistent, themed terminal
// forms and prompts using the huh library.
//
// The package is fully independent from flair's internal packages.
// External projects can import it without pulling in flair's internal
// implementation details.
//
// # Quick Start
//
// The simplest way to get a themed huh form is using [Default]:
//
//	theme, err := huh.Default()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	form := huh.NewForm(
//	    huh.NewGroup(
//	        huh.NewInput().Title("Name").Value(&name),
//	    ),
//	).WithTheme(theme)
//
// # Using a Specific Theme
//
// Load a specific flair theme and create a huh theme with [NewTheme]:
//
//	flairTheme, err := flair.LoadBuiltin("gruvbox-dark")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	huhTheme := huh.NewTheme(flairTheme)
//
// # Theme Components
//
// The returned [github.com/charmbracelet/huh.Theme] includes styled components:
//
//   - Form and Group base styles
//   - Focused and Blurred field styles (Title, Description, Errors)
//   - Select, MultiSelect, and FilePicker indicators
//   - TextInput styles (Cursor, Placeholder, Prompt, Text)
//   - Button styles (Focused and Blurred)
//   - Help text styles
//
// # Dependencies
//
// This package depends on:
//   - github.com/charmbracelet/huh
//   - github.com/curtbushko/flair/pkg/flair
package huh

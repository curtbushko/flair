// Package bubbletea provides themed lipgloss styles for bubbletea applications.
//
// This package integrates flair color themes with the Charm ecosystem, allowing
// bubbletea TUI applications to use consistent, theme-aware styling.
//
// # Quick Start
//
// The simplest way to get themed styles is using [Default], which loads styles
// from the currently selected flair theme:
//
//	styles, err := bubbletea.Default()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use surface styles for backgrounds
//	bg := styles.Surface.Background
//	panel := bg.Padding(1, 2).Render("Hello, World!")
//
//	// Use text styles for foregrounds
//	primary := styles.Text.Primary
//	msg := primary.Render("Welcome to the app")
//
//	// Use status styles for messages
//	errStyle := styles.Status.Error.Bold(true)
//	errMsg := errStyle.Render("An error occurred")
//
// # Creating Styles from a Theme
//
// For more control, use [NewStyles] with a [flair.Theme]:
//
//	theme, _ := flair.LoadBuiltin("gruvbox-dark")
//	styles := bubbletea.NewStyles(theme)
//
// # Available Style Categories
//
// The [Styles] struct provides pre-configured lipgloss styles organized by
// semantic category:
//
//   - [Styles.Surface]: Background styles (Background, Raised, Sunken, Selection, etc.)
//   - [Styles.Text]: Foreground text styles (Primary, Secondary, Muted, Subtle)
//   - [Styles.Status]: Status message styles (Error, Warning, Success, Info)
//   - [Styles.Accent]: Accent color styles (Primary, Secondary)
//   - [Styles.Border]: Border color styles (Default, Focus)
//
// Each style is a [lipgloss.Style] that can be further customized using
// lipgloss methods like Padding, Margin, Bold, etc.
package bubbletea

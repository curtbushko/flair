// Package viewport provides themed viewport styles for charmbracelet/bubbles viewport.
//
// This package integrates flair color themes with the bubbles viewport component,
// allowing TUI applications to use consistent, theme-aware viewport styling.
//
// # Quick Start
//
// The simplest way to get a themed viewport is using [NewModel] with a [flair.Theme]:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	vp := viewport.NewModel(theme, 80, 24)
//	vp.SetContent("Hello, world!")
//
// # Styling Details
//
// The viewport styles configure:
//   - Background: Uses surface.background from theme
//   - Foreground: Uses text.primary from theme
//
// # Available Functions
//
//   - [NewStyle]: Creates lipgloss.Style for viewport from a flair.Theme
//   - [NewModel]: Creates a pre-styled viewport.Model
//   - [Default]: Returns an empty lipgloss.Style (convenience function)
//
// All functions handle nil themes gracefully by returning defaults.
package viewport

// Package table provides themed table styles for charmbracelet/bubbles table.
//
// This package integrates flair color themes with the bubbles table component,
// allowing TUI applications to use consistent, theme-aware table styling.
//
// # Quick Start
//
// The simplest way to get themed table styles is using [NewStyles] with a [flair.Theme]:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	styles := table.NewStyles(theme)
//
//	// Apply to a bubbles table model
//	myTable := bubbletable.New(
//	    bubbletable.WithColumns(columns),
//	    bubbletable.WithRows(rows),
//	    bubbletable.WithStyles(styles),
//	)
//
// # Styling Details
//
// The table styles configure:
//   - Header: Bold text with secondary foreground color
//   - Cell: Primary text color for data cells
//   - Selected: Accent foreground with selection background for the focused row
//
// # Available Functions
//
//   - [NewStyles]: Creates table.Styles from a flair.Theme
//   - [Default]: Returns default table.Styles (shorthand for DefaultStyles)
//
// All functions handle nil themes gracefully by returning default styles.
package table

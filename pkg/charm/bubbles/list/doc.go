// Package list provides themed list styles and delegates for charmbracelet/bubbles list.
//
// This package integrates flair color themes with the bubbles list component,
// allowing TUI applications to use consistent, theme-aware list styling.
//
// # Quick Start
//
// The simplest way to get themed list styles is using [NewStyles] with a [flair.Theme]:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	styles := list.NewStyles(theme)
//
//	// Apply to a bubbles list model
//	myList := list.New(items, list.NewDelegate(theme), width, height)
//	myList.Styles = *styles
//
// # Creating a Themed Delegate
//
// Use [NewDelegate] to create a themed list delegate:
//
//	theme, _ := flair.LoadBuiltin("gruvbox-dark")
//	delegate := list.NewDelegate(theme)
//
//	// Create list with themed delegate
//	myList := bubbles_list.New(items, *delegate, width, height)
//
// # Available Functions
//
//   - [NewStyles]: Creates list.Styles from a flair.Theme
//   - [NewDelegate]: Creates a themed list.DefaultDelegate
//   - [NewItemStyles]: Creates list.DefaultItemStyles from a flair.Theme
//
// All functions handle nil themes gracefully by returning default styles.
package list

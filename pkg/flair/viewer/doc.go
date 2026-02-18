// Package viewer provides a bubbletea-based TUI for browsing and selecting
// flair color themes.
//
// The viewer has four pages:
//
//   - Theme Selector: List available themes, select one to apply
//   - Palette: Display base00-base17 colors with swatches
//   - Tokens: Show semantic tokens grouped by category (Surface, Text, Status, etc.)
//   - Components: Showcase styled UI components using theme colors
//
// # Usage
//
// The simplest way to use the viewer is via the Run function:
//
//	err := viewer.Run(viewer.Options{
//	    Themes:       []string{"tokyo-night-dark", "gruvbox-dark"},
//	    InitialTheme: "tokyo-night-dark",
//	    OnSelect: func(name string) {
//	        // Apply theme via symlinks or other mechanism
//	    },
//	})
//
// # Integration
//
// For advanced integration, create a Model directly:
//
//	model := viewer.NewModel(viewer.Options{...})
//	p := tea.NewProgram(model)
//	p.Run()
//
// The Model implements tea.Model and can be embedded in other bubbletea
// applications.
//
// # Key Bindings
//
//   - Tab: Switch pages
//   - j/k: Navigate up/down
//   - Enter: Select theme
//   - q: Quit
package viewer

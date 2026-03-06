// Package flair provides public types for loading and using flair color themes.
//
// This package is standalone and has no dependencies on internal/ packages.
// It is designed to be imported by external consumers who want to work with
// flair themes programmatically.
//
// # Core Types
//
// The main types are:
//
//   - [Color]: An RGB color value with methods for hex conversion and comparison.
//   - [Theme]: A named collection of semantic color tokens representing a color theme.
//   - [Palette]: A base24 color palette with 24 indexed color slots.
//   - [Store]: A manager for installing, selecting, and loading themes from disk.
//
// # Loading Themes
//
// The simplest way to load a theme is using [Default], which loads the currently
// selected theme or falls back to a built-in default:
//
//	theme, err := flair.Default()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Theme:", theme.Name())
//
// For more control, use [Load] to load only the selected theme, or [LoadBuiltin]
// to load a specific built-in theme by name:
//
//	// Load only if a theme is selected
//	theme, err := flair.Load()
//
//	// Load a specific built-in theme
//	theme, err := flair.LoadBuiltin("tokyo-night-dark")
//
// # Using Theme Colors
//
// Themes provide semantic color tokens organized by category. Use the typed
// accessor methods for compile-time safety:
//
//	theme, _ := flair.Default()
//
//	// Surface colors for backgrounds
//	bg := theme.Surface().Background
//	fmt.Println("Background:", bg.Hex())
//
//	// Text colors for foregrounds
//	fg := theme.Text().Primary
//	fmt.Println("Foreground:", fg.Hex())
//
//	// Status colors for messages
//	errColor := theme.Status().Error
//	fmt.Println("Error color:", errColor.Hex())
//
// Alternatively, use [Theme.Color] or [Theme.Get] for dynamic token lookup:
//
//	if c, ok := theme.Color("syntax.keyword"); ok {
//	    fmt.Println("Keyword color:", c.Hex())
//	}
//
// # Theme Management
//
// Use [Store] to manage themes on disk:
//
//	store := flair.NewStore()
//
//	// Install a built-in theme
//	store.Install("gruvbox-dark")
//
//	// Select a theme (creates symlinks)
//	store.Select("gruvbox-dark")
//
//	// List installed themes
//	themes, _ := store.List()
//
// # Built-in Themes
//
// The package embeds several popular color themes. Use [ListBuiltins] to see
// available themes and [LoadBuiltin] to load them without filesystem access:
//
//	for _, name := range flair.ListBuiltins() {
//	    fmt.Println(name)
//	}
//
// # Creating Custom Themes
//
// Create themes programmatically using [NewTheme]:
//
//	colors := map[string]flair.Color{
//	    "surface.background": {R: 26, G: 27, B: 38},
//	    "text.primary":       {R: 192, G: 202, B: 245},
//	}
//	theme := flair.NewTheme("my-theme", "dark", colors)
//
// Or parse a base24 palette and tokenize it:
//
//	palette, err := flair.ParsePalette(reader)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	theme := flair.Tokenize(palette)
package flair

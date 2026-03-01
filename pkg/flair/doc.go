// Package flair provides public types for loading and using flair color themes.
//
// This package is standalone and has no dependencies on internal/ packages.
// It is designed to be imported by external consumers who want to work with
// flair themes programmatically.
//
// The main types are:
//   - Color: An RGB color value with methods for hex conversion and comparison.
//   - Theme: A named collection of colors representing a color theme.
//
// Example usage:
//
//	colors := map[string]flair.Color{
//	    "background": {R: 26, G: 27, B: 38},
//	    "foreground": {R: 192, G: 202, B: 245},
//	}
//	theme := flair.NewTheme("tokyo-night", "storm", colors)
//
//	if bg, ok := theme.Color("background"); ok {
//	    fmt.Println("Background:", bg.Hex())
//	}
package flair

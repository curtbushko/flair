package flair

// Theme represents a color theme with a name, variant, and a collection of semantic colors.
//
// Theme is a read-only type designed for external consumers of flair themes.
// It provides typed accessor methods for common color categories (Surface, Text,
// Status, Syntax, Diff, Terminal) as well as direct token lookup via Color and Get.
//
// Themes are typically created via [Load], [LoadBuiltin], [Tokenize], or [NewTheme].
// All methods are safe to call on a nil Theme (they return zero values).
//
// Example:
//
//	theme, _ := flair.LoadBuiltin("tokyo-night-dark")
//	bg := theme.Surface().Background
//	fg := theme.Text().Primary
//	fmt.Printf("Background: %s, Foreground: %s\n", bg.Hex(), fg.Hex())
type Theme struct {
	name    string
	variant string
	colors  map[string]Color
}

// NewTheme creates a new Theme with the given name, variant, and colors.
//
// The colors map is defensively copied to ensure the Theme is immutable.
// Subsequent modifications to the original map do not affect the Theme.
//
// Parameters:
//   - name: The theme name (e.g., "tokyo-night-dark", "gruvbox-light").
//   - variant: The theme variant (e.g., "dark", "light", "storm").
//   - colors: A map of semantic token paths to colors.
//
// Example:
//
//	colors := map[string]flair.Color{
//	    "surface.background": {R: 26, G: 27, B: 38},
//	    "text.primary":       {R: 192, G: 202, B: 245},
//	}
//	theme := flair.NewTheme("my-theme", "dark", colors)
func NewTheme(name, variant string, colors map[string]Color) *Theme {
	// Copy the colors map to ensure immutability
	colorsCopy := make(map[string]Color, len(colors))
	for k, v := range colors {
		colorsCopy[k] = v
	}

	return &Theme{
		name:    name,
		variant: variant,
		colors:  colorsCopy,
	}
}

// Name returns the theme's name (e.g., "tokyo-night-dark").
func (t *Theme) Name() string {
	return t.name
}

// Variant returns the theme's variant (e.g., "dark", "light", "storm").
//
// The variant is typically extracted from the theme name or specified
// explicitly during theme creation. It can be used to determine if a
// theme is dark or light for UI adaptation.
func (t *Theme) Variant() string {
	return t.variant
}

// HasColors reports whether the theme has at least one color defined.
//
// This is useful for validating that a theme was loaded correctly.
func (t *Theme) HasColors() bool {
	return len(t.colors) > 0
}

// Color retrieves a color by its semantic token path.
//
// Token paths follow a hierarchical naming convention such as
// "surface.background", "text.primary", or "syntax.keyword".
//
// Color returns the color and true if found, or a zero Color and false if not.
//
// Example:
//
//	if c, ok := theme.Color("accent.primary"); ok {
//	    fmt.Println("Accent:", c.Hex())
//	}
func (t *Theme) Color(key string) (Color, bool) {
	c, ok := t.colors[key]
	return c, ok
}

// Colors returns a copy of all colors in the theme as a map.
//
// The returned map is a defensive copy; modifying it does not affect
// the theme. Keys are semantic token paths (e.g., "surface.background").
func (t *Theme) Colors() map[string]Color {
	result := make(map[string]Color, len(t.colors))
	for k, v := range t.colors {
		result[k] = v
	}
	return result
}

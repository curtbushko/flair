package flair

// Theme represents a color theme with a name, variant, and a collection of colors.
// This is a read-only type designed for external consumers of flair themes.
// It has no dependencies on internal/ packages.
type Theme struct {
	name    string
	variant string
	colors  map[string]Color
}

// NewTheme creates a new Theme with the given name, variant, and colors.
// The colors map is copied to ensure the Theme is immutable.
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

// Name returns the theme's name.
func (t *Theme) Name() string {
	return t.name
}

// Variant returns the theme's variant (e.g., "dark", "light", "storm").
func (t *Theme) Variant() string {
	return t.variant
}

// HasColors returns true if the theme has at least one color defined.
func (t *Theme) HasColors() bool {
	return len(t.colors) > 0
}

// Color retrieves a color by its key.
// Returns the color and true if found, or a zero Color and false if not.
func (t *Theme) Color(key string) (Color, bool) {
	c, ok := t.colors[key]
	return c, ok
}

// Colors returns a copy of all colors in the theme.
// Modifying the returned map does not affect the theme.
func (t *Theme) Colors() map[string]Color {
	result := make(map[string]Color, len(t.colors))
	for k, v := range t.colors {
		result[k] = v
	}
	return result
}

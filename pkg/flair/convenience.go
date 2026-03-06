package flair

// DefaultFallbackTheme is the built-in theme used when no theme is selected.
//
// This constant is used by [Default] and [MustLoad] as the fallback theme
// when no theme is explicitly selected in the config directory.
const DefaultFallbackTheme = "tokyo-night-dark"

// Default loads the currently selected theme or falls back to the default
// built-in theme (tokyo-night-dark) if no theme is selected.
//
// This is the recommended way to load a theme for most use cases. It ensures
// that a valid theme is always returned, even on first run before any themes
// are installed or selected.
//
// The fallback theme is [DefaultFallbackTheme] (tokyo-night-dark).
//
// Example:
//
//	theme, err := flair.Default()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Using theme: %s\n", theme.Name())
func Default() (*Theme, error) {
	theme, err := Load()
	if err == nil {
		return theme, nil
	}

	// If no theme is selected or load fails, fall back to built-in.
	return LoadBuiltin(DefaultFallbackTheme)
}

// MustLoad loads the currently selected theme or the default built-in theme.
//
// MustLoad panics if no theme can be loaded. This is useful for init-time
// loading where failure should be fatal, such as in package-level variables.
//
// Example:
//
//	var theme = flair.MustLoad() // panics on failure
//
//	func main() {
//	    fmt.Println("Background:", theme.Surface().Background.Hex())
//	}
func MustLoad() *Theme {
	theme, err := Default()
	if err != nil {
		panic("flair: failed to load theme: " + err.Error())
	}
	return theme
}

// LoadOrDefault tries to load the named theme from the config directory,
// falling back to a built-in theme if not found.
//
// This function is useful when you want to prefer an installed theme but
// guarantee a fallback to a known built-in theme.
//
// Parameters:
//   - name: The name of the theme to try loading from the config directory.
//   - fallback: The name of the built-in theme to use if name is not found.
//
// Returns the loaded theme or an error if both name and fallback fail to load.
//
// Example:
//
//	// Try user's custom theme, fall back to gruvbox-dark
//	theme, err := flair.LoadOrDefault("my-custom-theme", "gruvbox-dark")
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadOrDefault(name, fallback string) (*Theme, error) {
	// Try loading the named theme from the config directory.
	theme, err := LoadNamed(name)
	if err == nil {
		return theme, nil
	}

	// Fall back to the built-in fallback theme.
	return LoadBuiltin(fallback)
}

// EnsureInstalled installs all built-in themes to the config directory
// if the directory is empty (has no themes).
//
// This is useful for first-run setup to populate the config directory
// with all available themes. Call this early in your application's
// initialization to ensure themes are available for selection.
//
// If the config directory already contains at least one theme, this function
// does nothing and returns nil (idempotent behavior).
//
// Example:
//
//	func main() {
//	    // Ensure themes are installed on first run
//	    if err := flair.EnsureInstalled(); err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Now themes are guaranteed to be available
//	    themes, _ := flair.ListThemes()
//	    fmt.Printf("Available themes: %d\n", len(themes))
//	}
func EnsureInstalled() error {
	store := NewStore()

	// Check if any themes are already installed.
	themes, err := store.List()
	if err != nil {
		return err
	}

	// If themes exist, do nothing.
	if len(themes) > 0 {
		return nil
	}

	// Install all built-in themes.
	return store.InstallAll()
}

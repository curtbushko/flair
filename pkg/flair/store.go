package flair

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Store manages themes in a config directory.
//
// Store provides a high-level API for theme lifecycle management:
//   - Install built-in themes to disk
//   - Select a theme (creates symlinks for tools to consume)
//   - Load themes from disk
//   - List installed themes
//
// Use [NewStore] for the default config directory (~/.config/flair or
// $XDG_CONFIG_HOME/flair), or [NewStoreAt] for a custom path.
//
// Example:
//
//	store := flair.NewStore()
//
//	// Install and select a theme
//	if err := store.Install("tokyo-night-dark"); err != nil {
//	    log.Fatal(err)
//	}
//	if err := store.Select("tokyo-night-dark"); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Load the selected theme
//	theme, err := store.Load()
//	if err != nil {
//	    log.Fatal(err)
//	}
type Store struct {
	configDir string
}

// NewStore creates a new Store using the default config directory.
//
// The default config directory is determined as follows:
//   - If XDG_CONFIG_HOME is set: $XDG_CONFIG_HOME/flair
//   - Otherwise: ~/.config/flair
//
// The directory is created automatically when themes are installed.
func NewStore() *Store {
	return &Store{configDir: defaultConfigDir()}
}

// NewStoreAt creates a new Store rooted at the specified directory.
//
// This is useful for testing or when using a non-standard config location.
// The directory is created automatically when themes are installed.
func NewStoreAt(configDir string) *Store {
	return &Store{configDir: configDir}
}

// ConfigDir returns the root config directory path.
//
// This is useful for displaying the config location to users or for
// manual file operations.
func (s *Store) ConfigDir() string {
	return s.configDir
}

// Install copies a built-in theme to the config directory.
//
// Install creates the theme directory (e.g., ~/.config/flair/tokyo-night-dark/)
// and writes:
//   - tokens.yaml: Semantic color tokens
//   - style.lua: Neovim/Lua format
//   - style.css: CSS custom properties
//   - gtk.css: GTK CSS format
//   - style.qss: Qt stylesheet format
//   - style.json: JSON format
//
// Install is idempotent; calling it multiple times overwrites existing files.
// Returns an error if the built-in theme does not exist or files cannot be written.
//
// Example:
//
//	store := flair.NewStore()
//	if err := store.Install("gruvbox-dark"); err != nil {
//	    log.Fatal(err)
//	}
func (s *Store) Install(name string) error {
	if !HasBuiltin(name) {
		return fmt.Errorf("install theme %q: built-in theme not found", name)
	}

	theme, err := LoadBuiltin(name)
	if err != nil {
		return fmt.Errorf("install theme %q: %w", name, err)
	}

	themeDir := filepath.Join(s.configDir, name)
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		return fmt.Errorf("install theme %q: create directory: %w", name, err)
	}

	// Write tokens.yaml.
	if err := s.writeTokensYAML(themeDir, theme); err != nil {
		return fmt.Errorf("install theme %q: %w", name, err)
	}

	// Write output files.
	if err := s.writeOutputFiles(themeDir, name, theme); err != nil {
		return fmt.Errorf("install theme %q: %w", name, err)
	}

	return nil
}

// InstallAll installs all built-in themes to the config directory.
//
// This is a convenience method for populating the config directory with
// all available themes. It calls [Store.Install] for each theme returned
// by [ListBuiltins].
//
// Returns an error if any theme fails to install, stopping at the first failure.
func (s *Store) InstallAll() error {
	builtins := ListBuiltins()
	for _, name := range builtins {
		if err := s.Install(name); err != nil {
			return fmt.Errorf("install all: %w", err)
		}
	}
	return nil
}

// Select activates a theme by creating symlinks at the config root.
//
// Select creates relative symlinks from the config root to the theme's
// output files. For example:
//
//	~/.config/flair/style.lua -> tokyo-night-dark/style.lua
//	~/.config/flair/style.css -> tokyo-night-dark/style.css
//
// The theme must already be installed via [Store.Install]. Existing
// symlinks are removed before creating new ones.
//
// External tools can read these symlinks to apply the selected theme.
//
// Example:
//
//	store := flair.NewStore()
//	store.Install("tokyo-night-dark")
//	if err := store.Select("tokyo-night-dark"); err != nil {
//	    log.Fatal(err)
//	}
func (s *Store) Select(name string) error {
	themeDir := filepath.Join(s.configDir, name)
	tokensPath := filepath.Join(themeDir, "tokens.yaml")

	if _, err := os.Stat(tokensPath); os.IsNotExist(err) {
		return fmt.Errorf("select theme %q: theme not installed", name)
	}

	for _, f := range outputFiles {
		if err := s.createSymlink(name, f); err != nil {
			return fmt.Errorf("select theme %q: %w", name, err)
		}
	}

	return nil
}

// createSymlink creates a single relative symlink at the config root.
func (s *Store) createSymlink(themeName, filename string) error {
	link := filepath.Join(s.configDir, filename)
	target := filepath.Join(themeName, filename)

	// Remove existing symlink or file if present.
	if err := os.Remove(link); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove existing symlink %s: %w", filename, err)
	}

	if err := os.Symlink(target, link); err != nil {
		return fmt.Errorf("create symlink %s -> %s: %w", filename, target, err)
	}

	return nil
}

// Load loads the currently selected theme from the store's config directory.
//
// The selected theme is determined by following symlinks at the config root.
// Returns [ErrNoSelectedTheme] if no theme is selected.
//
// Example:
//
//	store := flair.NewStore()
//	theme, err := store.Load()
//	if errors.Is(err, flair.ErrNoSelectedTheme) {
//	    fmt.Println("No theme selected")
//	}
func (s *Store) Load() (*Theme, error) {
	selected, err := s.Selected()
	if err != nil {
		return nil, fmt.Errorf("load theme: %w", err)
	}

	if selected == "" {
		return nil, ErrNoSelectedTheme
	}

	return s.LoadNamed(selected)
}

// LoadNamed loads a specific theme by name from the store's config directory.
//
// The theme must be installed. Returns [ErrThemeNotFound] if the theme
// does not exist.
func (s *Store) LoadNamed(name string) (*Theme, error) {
	return LoadNamedFrom(s.configDir, name)
}

// Selected returns the name of the currently selected theme.
//
// Returns an empty string (not an error) if no theme is selected.
// The selection is determined by reading symlinks at the config root.
func (s *Store) Selected() (string, error) {
	return SelectedThemeFrom(s.configDir)
}

// List returns the names of all installed themes, sorted alphabetically.
//
// Theme directories are identified by the presence of a tokens.yaml file.
// Returns nil (not an error) if the config directory does not exist.
//
// Example:
//
//	store := flair.NewStore()
//	themes, _ := store.List()
//	for _, name := range themes {
//	    fmt.Println(name)
//	}
func (s *Store) List() ([]string, error) {
	return ListThemesFrom(s.configDir)
}

// writeTokensYAML writes the tokens.yaml file for a theme.
func (s *Store) writeTokensYAML(themeDir string, theme *Theme) error {
	tokensPath := filepath.Join(themeDir, "tokens.yaml")

	// Build the tokens structure.
	tokens := make(map[string]tokensEntry)
	for key, color := range theme.Colors() {
		tokens[key] = tokensEntry{
			Color: color.Hex(),
		}
	}

	tf := tokensFile{Tokens: tokens}

	data, err := yaml.Marshal(tf)
	if err != nil {
		return fmt.Errorf("marshal tokens.yaml: %w", err)
	}

	if err := os.WriteFile(tokensPath, data, 0o644); err != nil {
		return fmt.Errorf("write tokens.yaml: %w", err)
	}

	return nil
}

// writeOutputFiles writes all output files (style.lua, style.css, etc.) for a theme.
func (s *Store) writeOutputFiles(themeDir, name string, theme *Theme) error {
	// Write style.lua (Neovim/Lua format).
	if err := s.writeStyleLua(themeDir, theme); err != nil {
		return err
	}

	// Write style.css (CSS custom properties).
	if err := s.writeStyleCSS(themeDir, theme); err != nil {
		return err
	}

	// Write gtk.css (GTK CSS format).
	if err := s.writeGtkCSS(themeDir, theme); err != nil {
		return err
	}

	// Write style.qss (Qt stylesheet format).
	if err := s.writeStyleQSS(themeDir, theme); err != nil {
		return err
	}

	// Write style.json (JSON format).
	return s.writeStyleJSON(themeDir, name, theme)
}

// writeStyleLua writes the Lua output file.
func (s *Store) writeStyleLua(themeDir string, theme *Theme) error {
	path := filepath.Join(themeDir, "style.lua")

	colors := theme.Colors()
	keys := make([]string, 0, len(colors))
	for k := range colors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString("return {\n")
	for _, k := range keys {
		fmt.Fprintf(&sb, "  [\"%s\"] = \"%s\",\n", k, colors[k].Hex())
	}
	sb.WriteString("}\n")

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// writeStyleCSS writes the CSS output file with custom properties.
func (s *Store) writeStyleCSS(themeDir string, theme *Theme) error {
	path := filepath.Join(themeDir, "style.css")

	colors := theme.Colors()
	keys := make([]string, 0, len(colors))
	for k := range colors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(":root {\n")
	for _, k := range keys {
		fmt.Fprintf(&sb, "  --%s: %s;\n", k, colors[k].Hex())
	}
	sb.WriteString("}\n")

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// writeGtkCSS writes the GTK CSS output file.
func (s *Store) writeGtkCSS(themeDir string, theme *Theme) error {
	path := filepath.Join(themeDir, "gtk.css")

	colors := theme.Colors()
	keys := make([]string, 0, len(colors))
	for k := range colors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString("/* GTK Theme Colors */\n")
	fmt.Fprintf(&sb, "@define-color bg_color %s;\n", colors["surface.background"].Hex())
	if fg, ok := colors["text.primary"]; ok {
		fmt.Fprintf(&sb, "@define-color fg_color %s;\n", fg.Hex())
	}

	// Add all colors as define-color.
	for _, k := range keys {
		fmt.Fprintf(&sb, "@define-color %s %s;\n", k, colors[k].Hex())
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// writeStyleQSS writes the Qt stylesheet output file.
func (s *Store) writeStyleQSS(themeDir string, theme *Theme) error {
	path := filepath.Join(themeDir, "style.qss")

	colors := theme.Colors()

	bg := "#1a1b26"
	fg := "#c0caf5"
	if c, ok := colors["surface.background"]; ok {
		bg = c.Hex()
	}
	if c, ok := colors["text.primary"]; ok {
		fg = c.Hex()
	}

	content := fmt.Sprintf(`/* Qt Stylesheet */
QWidget {
  background-color: %s;
  color: %s;
}
`, bg, fg)

	return os.WriteFile(path, []byte(content), 0o644)
}

// writeStyleJSON writes the JSON output file.
func (s *Store) writeStyleJSON(themeDir, name string, theme *Theme) error {
	path := filepath.Join(themeDir, "style.json")

	colors := theme.Colors()
	colorMap := make(map[string]string, len(colors))
	for k, c := range colors {
		colorMap[k] = c.Hex()
	}

	output := map[string]interface{}{
		"name":    name,
		"variant": theme.Variant(),
		"colors":  colorMap,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal style.json: %w", err)
	}

	return os.WriteFile(path, data, 0o644)
}

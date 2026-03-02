package main

import (
	"flag"
	"io"
	"strings"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/config"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/pkg/flair/viewer"
)

// viewerDryRun is a package-level flag that tests can set to skip the TUI.
// This is used only in test builds to avoid TTY requirements.
var viewerDryRun bool

// runSelect parses flags and executes the select subcommand.
// It writes diagnostic/error output to stderr and confirmation to stdout.
// Returns an exit code.
func runSelect(args []string, stdout, stderr io.Writer) int {
	subArgs := args[2:]
	flagArgs, positional, viewerFlag := splitSelectFlagsAndArgs(subArgs)

	fs := flag.NewFlagSet("select", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")

	if err := fs.Parse(flagArgs); err != nil {
		return 1
	}

	// Resolve config directory.
	configDir := *dirFlag
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	// Wire the app with the config directory.
	app := Wire(configDir)

	// Determine if we should launch viewer mode.
	// Viewer mode is triggered when:
	// 1. --viewer flag is set, OR
	// 2. No theme argument is provided
	if viewerFlag || len(positional) == 0 {
		return runSelectViewer(stdout, stderr, app, positional)
	}

	// Standard select mode: apply theme via symlinks.
	themeName := positional[0]

	if err := app.Select.Execute(themeName); err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	writeStr(stdout, "Selected theme %q\n", themeName)
	return 0
}

// runSelectViewer launches the style viewer TUI.
func runSelectViewer(stdout, stderr io.Writer, app *App, positional []string) int {
	// Get available themes.
	themeInfos, err := app.List.Execute()
	if err != nil {
		writeStr(stderr, "error listing themes: %v\n", err)
		return 1
	}

	// Extract theme names from ThemeInfo.
	themes := make([]string, 0, len(themeInfos))
	for _, info := range themeInfos {
		themes = append(themes, info.Name)
	}

	// Determine initial theme.
	var initialTheme string
	if len(positional) > 0 {
		initialTheme = positional[0]
	} else {
		// Try to get the currently selected theme.
		selected, err := app.Store.SelectedTheme()
		if err == nil && selected != "" {
			initialTheme = selected
		}
	}

	// Print confirmation before launching viewer.
	if initialTheme != "" {
		writeStr(stdout, "Launching style viewer with theme %q pre-selected...\n", initialTheme)
	} else {
		writeStr(stdout, "Launching style viewer...\n")
	}

	// Create the theme loader adapter.
	loader := &appThemeLoader{app: app}

	// Configure viewer options.
	opts := viewer.Options{
		Themes:       themes,
		InitialTheme: initialTheme,
		ThemeLoader:  loader,
		OnSelect: func(name string) {
			// Apply theme via symlinks when selected.
			if err := app.Select.Execute(name); err != nil {
				// Log error but don't interrupt the viewer.
				writeStr(stderr, "warning: failed to select theme: %v\n", err)
			}
		},
	}

	// Launch the viewer TUI.
	runOpts := viewer.RunOptions{
		WithAltScreen: true,
		DryRun:        viewerDryRun,
	}
	if err := viewer.RunWithOptions(opts, runOpts); err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

// appThemeLoader adapts App dependencies to viewer.ThemeLoader interface.
type appThemeLoader struct {
	app *App
}

// LoadPalette returns base24 colors for a theme.
func (l *appThemeLoader) LoadPalette(name string) (viewer.PaletteData, error) {
	pd := viewer.PaletteData{}

	// First check if the theme exists on disk.
	rc, err := l.app.Store.OpenReader(name, "palette.yaml")
	if err == nil {
		defer func() { _ = rc.Close() }()
		palette, err := l.app.Preview.Parser().Parse(rc)
		if err == nil {
			for i := 0; i < 24; i++ {
				pd.Colors[i] = palette.Colors[i].Hex()
			}
			return pd, nil
		}
	}

	// Check built-in palettes.
	if l.app.Builtins.Has(name) {
		r, err := l.app.Builtins.Get(name)
		if err != nil {
			return pd, err
		}
		palette, err := l.app.Preview.Parser().Parse(r)
		if err != nil {
			return pd, err
		}
		for i := 0; i < 24; i++ {
			pd.Colors[i] = palette.Colors[i].Hex()
		}
		return pd, nil
	}

	return pd, nil
}

// LoadTokens returns semantic tokens for a theme.
func (l *appThemeLoader) LoadTokens(name string) (viewer.TokenData, error) {
	td := viewer.TokenData{
		Surface: make(map[string]string),
		Text:    make(map[string]string),
		Status:  make(map[string]string),
		Syntax:  make(map[string]string),
		Diff:    make(map[string]string),
	}

	// Try to read tokens.yaml.
	rc, err := l.app.Store.OpenReader(name, "tokens.yaml")
	if err != nil {
		// No tokens.yaml, try to derive from palette.
		return l.deriveTokensFromPalette(name, td)
	}
	defer func() { _ = rc.Close() }()

	ts, err := fileio.ReadTokens(rc)
	if err != nil {
		return td, err
	}

	// Group tokens by category prefix.
	for _, path := range ts.Paths() {
		tok, _ := ts.Get(path)
		if tok.Color.IsNone {
			continue
		}
		hex := tok.Color.Hex()

		switch {
		case strings.HasPrefix(path, "surface."):
			td.Surface[path] = hex
		case strings.HasPrefix(path, "text."):
			td.Text[path] = hex
		case strings.HasPrefix(path, "status."):
			td.Status[path] = hex
		case strings.HasPrefix(path, "syntax."):
			td.Syntax[path] = hex
		case strings.HasPrefix(path, "diff."):
			td.Diff[path] = hex
		}
	}

	return td, nil
}

// deriveTokensFromPalette derives tokens from a palette when tokens.yaml doesn't exist.
func (l *appThemeLoader) deriveTokensFromPalette(name string, td viewer.TokenData) (viewer.TokenData, error) {
	// Try to parse palette from store.
	rc, err := l.app.Store.OpenReader(name, "palette.yaml")
	if err == nil {
		defer func() { _ = rc.Close() }()
		palette, err := l.app.Preview.Parser().Parse(rc)
		if err == nil {
			ts := l.app.Preview.Tokenizer().Tokenize(palette)
			return l.tokenSetToTokenData(ts), nil
		}
	}

	// Check built-in palettes.
	if l.app.Builtins.Has(name) {
		r, err := l.app.Builtins.Get(name)
		if err != nil {
			return td, nil
		}
		palette, err := l.app.Preview.Parser().Parse(r)
		if err != nil {
			return td, nil
		}
		ts := l.app.Preview.Tokenizer().Tokenize(palette)
		return l.tokenSetToTokenData(ts), nil
	}

	return td, nil
}

// tokenSetToTokenData converts a domain.TokenSet to viewer.TokenData.
func (l *appThemeLoader) tokenSetToTokenData(ts *domain.TokenSet) viewer.TokenData {
	td := viewer.TokenData{
		Surface: make(map[string]string),
		Text:    make(map[string]string),
		Status:  make(map[string]string),
		Syntax:  make(map[string]string),
		Diff:    make(map[string]string),
	}

	for _, path := range ts.Paths() {
		tok, _ := ts.Get(path)
		if tok.Color.IsNone {
			continue
		}
		hex := tok.Color.Hex()

		switch {
		case strings.HasPrefix(path, "surface."):
			td.Surface[path] = hex
		case strings.HasPrefix(path, "text."):
			td.Text[path] = hex
		case strings.HasPrefix(path, "status."):
			td.Status[path] = hex
		case strings.HasPrefix(path, "syntax."):
			td.Syntax[path] = hex
		case strings.HasPrefix(path, "diff."):
			td.Diff[path] = hex
		}
	}

	return td
}

// splitSelectFlagsAndArgs separates flag arguments from positional arguments
// for the select subcommand. Returns flags, positional args, and whether
// --viewer was specified.
func splitSelectFlagsAndArgs(args []string) (flags, positional []string, viewerFlag bool) {
	// Flags that take a value argument.
	flagsWithValue := map[string]bool{
		"--dir": true,
	}

	// Boolean flags (no value).
	boolFlags := map[string]bool{
		"--viewer": true,
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == flagHelpShrt || arg == flagHelp {
			flags = append(flags, arg)
			i++
			continue
		}

		if arg == "--viewer" {
			viewerFlag = true
			i++
			continue
		}

		if boolFlags[arg] {
			i++
			continue
		}

		if flagsWithValue[arg] {
			flags = append(flags, arg)
			i++
			if i < len(args) {
				flags = append(flags, args[i])
				i++
			}
			continue
		}

		positional = append(positional, arg)
		i++
	}

	return flags, positional, viewerFlag
}

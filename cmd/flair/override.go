package main

import (
	"flag"
	"fmt"
	"io"
	"sort"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/config"
	"github.com/curtbushko/flair/internal/domain"
)

const overrideUsage = `Usage: flair override <theme> [options]

Manage token color and style overrides for a theme.

Add/update override:
  flair override <theme> <token> <color>
  flair override <theme> <token> --bold --italic

List overrides:
  flair override <theme> --list

Remove override:
  flair override <theme> --remove <token>

Options:
  --dir <path>       Config directory (default: ~/.config/flair)
  --list             List all overrides for the theme
  --remove <token>   Remove the override for the specified token
  --bold             Set bold style
  --italic           Set italic style
  --underline        Set underline style
  --undercurl        Set undercurl style
  --strikethrough    Set strikethrough style
  -h, --help         Show this help message

Examples:
  flair override mytheme syntax.keyword #ff00ff
  flair override mytheme syntax.keyword --bold --italic
  flair override mytheme syntax.keyword #ff00ff --bold
  flair override mytheme --list
  flair override mytheme --remove syntax.keyword
`

// overrideOpts holds parsed flags for the override command.
type overrideOpts struct {
	configDir     string
	list          bool
	remove        string
	bold          bool
	italic        bool
	underline     bool
	undercurl     bool
	strikethrough bool
}

// hasStyle returns true if any style flag is set.
func (o *overrideOpts) hasStyle() bool {
	return o.bold || o.italic || o.underline || o.undercurl || o.strikethrough
}

// runOverride parses flags and executes the override subcommand.
// It manages token overrides in a theme's palette.yaml file.
// Returns an exit code: 0 on success, 1 on error.
func runOverride(args []string, stdout, stderr io.Writer) int {
	themeName, positional, opts, code := parseOverrideArgs(args, stderr)
	if code >= 0 {
		return code
	}

	// Wire the app.
	app := Wire(opts.configDir)

	// Check if theme exists.
	if !app.Store.FileExists(themeName, "palette.yaml") {
		writeStr(stderr, "error: theme %q does not exist (no palette.yaml found)\n", themeName)
		return 1
	}

	return dispatchOverrideAction(themeName, positional, opts, app, stdout, stderr)
}

// parseOverrideArgs parses command-line arguments for override command.
// Returns themeName, remaining positional args, options, and exit code.
// Exit code < 0 means continue processing; code >= 0 means return that code.
func parseOverrideArgs(args []string, stderr io.Writer) (string, []string, *overrideOpts, int) {
	if len(args) < 3 {
		writeStr(stderr, "error: theme argument is required\n\n")
		writeStr(stderr, overrideUsage)
		return "", nil, nil, 1
	}

	// Check for help flag early.
	if hasHelpFlag(args[2:]) {
		writeStr(stderr, overrideUsage)
		return "", nil, nil, 0
	}

	subArgs := args[2:]
	flagArgs, positional := splitOverrideFlags(subArgs)

	if len(positional) == 0 {
		writeStr(stderr, "error: theme argument is required\n\n")
		writeStr(stderr, overrideUsage)
		return "", nil, nil, 1
	}

	themeName := positional[0]
	positional = positional[1:]

	opts, err := parseOverrideFlags(flagArgs, stderr)
	if err != nil {
		return "", nil, nil, 1
	}

	return themeName, positional, opts, -1
}

// hasHelpFlag checks if any arg is a help flag.
func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == flagHelp || arg == flagHelpShrt {
			return true
		}
	}
	return false
}

// parseOverrideFlags parses the flag arguments and returns options.
func parseOverrideFlags(flagArgs []string, stderr io.Writer) (*overrideOpts, error) {
	fs := flag.NewFlagSet("override", flag.ContinueOnError)
	fs.SetOutput(stderr)

	opts := &overrideOpts{}
	fs.StringVar(&opts.configDir, "dir", "", "config directory (default: ~/.config/flair)")
	fs.BoolVar(&opts.list, "list", false, "list all overrides")
	fs.StringVar(&opts.remove, "remove", "", "token to remove")
	fs.BoolVar(&opts.bold, "bold", false, "set bold style")
	fs.BoolVar(&opts.italic, "italic", false, "set italic style")
	fs.BoolVar(&opts.underline, "underline", false, "set underline style")
	fs.BoolVar(&opts.undercurl, "undercurl", false, "set undercurl style")
	fs.BoolVar(&opts.strikethrough, "strikethrough", false, "set strikethrough style")

	if err := fs.Parse(flagArgs); err != nil {
		return nil, err
	}

	if opts.configDir == "" {
		opts.configDir = config.DefaultConfigDir()
	}

	return opts, nil
}

// dispatchOverrideAction routes to the appropriate handler based on flags.
func dispatchOverrideAction(
	themeName string, positional []string, opts *overrideOpts,
	app *App, stdout, stderr io.Writer,
) int {
	if opts.list {
		return handleList(themeName, app, stdout, stderr)
	}

	if opts.remove != "" {
		return handleRemove(themeName, opts.remove, app, stdout, stderr)
	}

	return handleAddUpdateAction(themeName, positional, opts, app, stdout, stderr)
}

// handleAddUpdateAction validates args and calls handleAddUpdate.
func handleAddUpdateAction(
	themeName string, positional []string, opts *overrideOpts,
	app *App, stdout, stderr io.Writer,
) int {
	if len(positional) == 0 {
		writeStr(stderr, "error: token name is required for add/update\n\n")
		writeStr(stderr, overrideUsage)
		return 1
	}

	tokenName := positional[0]
	var colorHex string
	if len(positional) > 1 {
		colorHex = positional[1]
	}

	if colorHex == "" && !opts.hasStyle() {
		writeStr(stderr, "error: must specify color and/or style flags\n\n")
		writeStr(stderr, overrideUsage)
		return 1
	}

	return handleAddUpdate(
		themeName, tokenName, colorHex,
		opts.bold, opts.italic, opts.underline, opts.undercurl, opts.strikethrough,
		app, stdout, stderr,
	)
}

// handleList lists all overrides for a theme.
func handleList(themeName string, app *App, stdout, stderr io.Writer) int {
	pal, err := readPalette(themeName, app)
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	if len(pal.Overrides) == 0 {
		writeStr(stdout, "No overrides defined for theme %q\n", themeName)
		return 0
	}

	keys := sortedKeys(pal.Overrides)

	writeStr(stdout, "Overrides for theme %q:\n", themeName)
	for _, key := range keys {
		writeOverrideLine(stdout, key, pal.Overrides[key])
	}

	return 0
}

// sortedKeys returns the sorted keys of an override map.
func sortedKeys(overrides map[string]domain.TokenOverride) []string {
	keys := make([]string, 0, len(overrides))
	for k := range overrides {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// writeOverrideLine writes a single override to stdout.
func writeOverrideLine(w io.Writer, key string, override domain.TokenOverride) {
	writeStr(w, "  %s:", key)
	if override.Color != nil {
		writeStr(w, " color=#%s", override.Color.Hex())
	}
	if override.Bold {
		writeStr(w, " bold")
	}
	if override.Italic {
		writeStr(w, " italic")
	}
	if override.Underline {
		writeStr(w, " underline")
	}
	if override.Undercurl {
		writeStr(w, " undercurl")
	}
	if override.Strikethrough {
		writeStr(w, " strikethrough")
	}
	writeStr(w, "\n")
}

// handleRemove removes an override from a theme.
func handleRemove(themeName, tokenName string, app *App, stdout, stderr io.Writer) int {
	pal, err := readPalette(themeName, app)
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	if pal.Overrides == nil {
		writeStr(stdout, "No overrides defined for theme %q\n", themeName)
		return 0
	}

	if _, exists := pal.Overrides[tokenName]; !exists {
		writeStr(stdout, "Override %q not found in theme %q\n", tokenName, themeName)
		return 0
	}

	delete(pal.Overrides, tokenName)

	if err := writePalette(themeName, pal, app); err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	writeStr(stdout, "Removed override %q from theme %q\n", tokenName, themeName)
	return 0
}

// handleAddUpdate adds or updates an override.
func handleAddUpdate(
	themeName, tokenName, colorHex string,
	bold, italic, underline, undercurl, strikethrough bool,
	app *App, stdout, stderr io.Writer,
) int {
	pal, err := readPalette(themeName, app)
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	override, err := domain.NewTokenOverride(colorHex, bold, italic, underline, undercurl, strikethrough)
	if err != nil {
		writeStr(stderr, "error: invalid color %q: %v\n", colorHex, err)
		return 1
	}

	if pal.Overrides == nil {
		pal.Overrides = make(map[string]domain.TokenOverride)
	}

	_, exists := pal.Overrides[tokenName]
	pal.Overrides[tokenName] = *override

	if err := writePalette(themeName, pal, app); err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	if exists {
		writeStr(stdout, "Updated override %q in theme %q\n", tokenName, themeName)
	} else {
		writeStr(stdout, "Added override %q to theme %q\n", tokenName, themeName)
	}
	return 0
}

// readPalette reads the palette.yaml file for a theme.
func readPalette(themeName string, app *App) (*domain.Palette, error) {
	r, err := app.Store.OpenReader(themeName, "palette.yaml")
	if err != nil {
		return nil, fmt.Errorf("open palette: %w", err)
	}
	defer func() { _ = r.Close() }()

	parser := app.Validate.Parser()
	pal, err := parser.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parse palette: %w", err)
	}

	return pal, nil
}

// writePalette writes the palette.yaml file for a theme.
func writePalette(themeName string, pal *domain.Palette, app *App) error {
	w, err := app.Store.OpenWriter(themeName, "palette.yaml")
	if err != nil {
		return fmt.Errorf("open palette for writing: %w", err)
	}
	defer func() { _ = w.Close() }()

	if err := fileio.WritePalette(w, pal); err != nil {
		return fmt.Errorf("write palette: %w", err)
	}

	return nil
}

// splitOverrideFlags separates flag arguments from positional arguments.
func splitOverrideFlags(args []string) (flags, positional []string) {
	knownFlagsWithValue := map[string]bool{
		"--dir":    true,
		"--remove": true,
	}

	knownFlagsNoValue := map[string]bool{
		"--list":          true,
		"--bold":          true,
		"--italic":        true,
		"--underline":     true,
		"--undercurl":     true,
		"--strikethrough": true,
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == flagHelpShrt || arg == flagHelp {
			flags = append(flags, arg)
			i++
			continue
		}

		if knownFlagsWithValue[arg] {
			flags = append(flags, arg)
			i++
			if i < len(args) {
				flags = append(flags, args[i])
				i++
			}
			continue
		}

		if knownFlagsNoValue[arg] {
			flags = append(flags, arg)
			i++
			continue
		}

		positional = append(positional, arg)
		i++
	}

	return flags, positional
}

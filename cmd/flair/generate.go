package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/curtbushko/flair/internal/config"
)

const generateUsage = `Usage: flair generate <palette> [options]

Generate theme files from a palette. The palette argument can be a
built-in name (e.g. "tokyo-night-dark") or a file path to a palette YAML.

Options:
  --dir <path>     Config directory (default: ~/.config/flair)
  --target <name>  Generate only the named target (e.g. "stylix")
  --name <name>    Override the theme name (default: inferred from palette)
  -h, --help       Show this help message
`

// runGenerate parses flags and executes the generate subcommand.
// It writes diagnostic/error output to stderr and summary output to stdout.
// Returns an exit code.
func runGenerate(args []string, stdout, stderr io.Writer) int {
	// Separate flags and positional args so flags can appear in any position.
	// Go's flag package stops at the first non-flag arg by default.
	subArgs := args[2:]
	flagArgs, positional := splitFlagsAndArgs(subArgs)

	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")
	targetFlag := fs.String("target", "", "generate only the named target")
	nameFlag := fs.String("name", "", "override the theme name")

	if err := fs.Parse(flagArgs); err != nil {
		return 1
	}

	if len(positional) == 0 {
		writeStr(stderr, "error: palette argument is required\n\n")
		writeStr(stderr, generateUsage)
		return 1
	}

	paletteRef := positional[0]

	// Resolve config directory.
	configDir := *dirFlag
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	// Wire the app with the config directory.
	app := Wire(configDir)

	// Resolve theme name and execute the pipeline.
	themeName, err := executeGenerate(app, paletteRef, *nameFlag, *targetFlag)
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	// Print summary to stdout.
	themeDir := app.Store.ThemeDir(themeName)
	fileCount := countFiles(themeDir)
	writeStr(stdout, "Generated theme %q: %d files written to %s\n", themeName, fileCount, themeDir)

	return 0
}

// executeGenerate resolves the palette reference and runs the generate pipeline.
// It returns the theme name used and any error.
func executeGenerate(app *App, paletteRef, nameOverride, targetFilter string) (string, error) {
	if app.Builtins.Has(paletteRef) {
		return executeBuiltin(app, paletteRef, nameOverride, targetFilter)
	}

	return executeFile(app, paletteRef, nameOverride, targetFilter)
}

// executeBuiltin runs the generate pipeline for a built-in palette name.
func executeBuiltin(app *App, builtinName, nameOverride, targetFilter string) (string, error) {
	themeName := builtinName
	if nameOverride != "" {
		themeName = nameOverride
	}

	err := app.Generate.ExecuteBuiltin(builtinName, themeName, targetFilter)
	return themeName, err
}

// executeFile runs the generate pipeline from a palette file path.
func executeFile(app *App, path, nameOverride, targetFilter string) (string, error) {
	themeName := inferNameFromPath(path)
	if nameOverride != "" {
		themeName = nameOverride
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	return themeName, app.Generate.Execute(f, themeName, targetFilter)
}

// inferNameFromPath extracts a theme name from a file path by using
// the base filename without extension.
func inferNameFromPath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext != "" {
		return base[:len(base)-len(ext)]
	}
	return base
}

// countFiles counts the number of files in a directory.
func countFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}
	return count
}

// splitFlagsAndArgs separates flag arguments (--flag value) from positional
// arguments so the Go flag package can parse them regardless of order.
func splitFlagsAndArgs(args []string) (flags, positional []string) {
	knownFlags := map[string]bool{
		"--dir":    true,
		"--target": true,
		"--name":   true,
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == flagHelpShrt || arg == flagHelp {
			flags = append(flags, arg)
			i++
			continue
		}

		if knownFlags[arg] {
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

	return flags, positional
}

package main

import (
	"flag"
	"io"

	"github.com/curtbushko/flair/internal/config"
)

// runList parses flags and executes the list subcommand.
// It writes theme list to stdout and errors to stderr.
// Returns an exit code.
func runList(args []string, stdout, stderr io.Writer) int {
	subArgs := args[2:]
	flagArgs := splitListFlags(subArgs)

	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")
	builtinsFlag := fs.Bool("builtins", false, "list built-in palette names only")

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

	if *builtinsFlag {
		return listBuiltins(app, stdout)
	}

	return listThemes(app, stdout, stderr)
}

// listBuiltins prints the names of all built-in palettes.
func listBuiltins(app *App, stdout io.Writer) int {
	names := app.List.ListBuiltins()
	for _, name := range names {
		writeStr(stdout, "%s\n", name)
	}
	return 0
}

// listThemes prints installed themes with a selection marker.
func listThemes(app *App, stdout, stderr io.Writer) int {
	themes, err := app.List.Execute()
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	if len(themes) == 0 {
		writeStr(stdout, "No themes installed. Use 'flair generate <palette>' to create one.\n")
		return 0
	}

	for _, info := range themes {
		marker := "  "
		if info.Selected {
			marker = "* "
		}
		writeStr(stdout, "%s%s\n", marker, info.Name)
	}

	return 0
}

// splitListFlags separates flag arguments for the list subcommand.
func splitListFlags(args []string) []string {
	knownFlags := map[string]bool{
		"--dir": true,
	}
	knownBoolFlags := map[string]bool{
		"--builtins": true,
	}

	var flags []string
	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == flagHelpShrt || arg == flagHelp {
			flags = append(flags, arg)
			i++
			continue
		}

		if knownBoolFlags[arg] {
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

		// Ignore unexpected positional args for list.
		i++
	}

	return flags
}

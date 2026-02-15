package main

import (
	"flag"
	"io"

	"github.com/curtbushko/flair/internal/config"
)

// runPreview parses flags and executes the preview subcommand.
// It writes the ANSI-colored preview to stdout and errors to stderr.
// Returns an exit code.
func runPreview(args []string, stdout, stderr io.Writer) int {
	subArgs := args[2:]
	flagArgs, themeName := splitPreviewArgs(subArgs)

	fs := flag.NewFlagSet("preview", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")

	if err := fs.Parse(flagArgs); err != nil {
		return 1
	}

	if themeName == "" {
		writeStr(stderr, "error: missing theme name\n")
		writeStr(stderr, "Usage: flair preview <theme-name> [--dir <path>]\n")
		return 1
	}

	// Resolve config directory.
	configDir := *dirFlag
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	// Wire the app with the config directory.
	app := Wire(configDir)

	if err := app.Preview.Execute(themeName, stdout); err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

// splitPreviewArgs separates the theme name positional arg from flag arguments.
func splitPreviewArgs(args []string) (flags []string, themeName string) {
	knownFlags := map[string]bool{
		"--dir": true,
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

		// First non-flag argument is the theme name.
		if !isFlag(arg) && themeName == "" {
			themeName = arg
			i++
			continue
		}

		i++
	}

	return flags, themeName
}

// isFlag returns true if the argument looks like a flag (starts with "-").
func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

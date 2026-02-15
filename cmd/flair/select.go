package main

import (
	"flag"
	"io"

	"github.com/curtbushko/flair/internal/config"
)

const selectUsage = `Usage: flair select <theme-name> [options]

Select a theme as active by creating symlinks to its output files.

Options:
  --dir <path>     Config directory (default: ~/.config/flair)
  -h, --help       Show this help message
`

// runSelect parses flags and executes the select subcommand.
// It writes diagnostic/error output to stderr and confirmation to stdout.
// Returns an exit code.
func runSelect(args []string, stdout, stderr io.Writer) int {
	subArgs := args[2:]
	flagArgs, positional := splitSelectFlagsAndArgs(subArgs)

	fs := flag.NewFlagSet("select", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")

	if err := fs.Parse(flagArgs); err != nil {
		return 1
	}

	if len(positional) == 0 {
		writeStr(stderr, "error: theme argument is required\n\n")
		writeStr(stderr, selectUsage)
		return 1
	}

	themeName := positional[0]

	// Resolve config directory.
	configDir := *dirFlag
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	// Wire the app with the config directory.
	app := Wire(configDir)

	if err := app.Select.Execute(themeName); err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	writeStr(stdout, "Selected theme %q\n", themeName)
	return 0
}

// splitSelectFlagsAndArgs separates flag arguments from positional arguments
// for the select subcommand.
func splitSelectFlagsAndArgs(args []string) (flags, positional []string) {
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

		positional = append(positional, arg)
		i++
	}

	return flags, positional
}

package main

import (
	"errors"
	"flag"
	"io"

	"github.com/curtbushko/flair/internal/config"
)

const regenerateUsage = `Usage: flair regenerate <theme-name> [options]

Regenerate all theme files from the palette. This ensures all downstream
files (tokens, mappings, outputs) are consistent with the current theme
definition and mapper code.

Options:
  --dir <path>     Config directory (default: ~/.config/flair)
  --target <name>  Regenerate only the named target (e.g. "vim")
  -h, --help       Show this help message
`

// runRegenerate parses flags and executes the regenerate subcommand.
func runRegenerate(args []string, stdout, stderr io.Writer) int {
	subArgs := args[2:]
	flagArgs, positional := splitRegenFlagsAndArgs(subArgs)

	fs := flag.NewFlagSet("regenerate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")
	targetFlag := fs.String("target", "", "regenerate only the named target")

	if err := fs.Parse(flagArgs); err != nil {
		return handleRegenFlagError(err)
	}

	if len(positional) == 0 {
		writeStr(stderr, "error: theme argument is required\n\n")
		writeStr(stderr, regenerateUsage)
		return 1
	}

	themeName := positional[0]

	configDir := *dirFlag
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	app := Wire(configDir)

	msg, err := app.Regenerate.Execute(themeName, *targetFlag)
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	writeStr(stdout, "%s\n", msg)
	return 0
}

// handleRegenFlagError returns 0 for help flag, 1 for real errors.
func handleRegenFlagError(err error) int {
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	return 1
}

// splitRegenFlagsAndArgs separates flag arguments from positional arguments
// for the regenerate command.
func splitRegenFlagsAndArgs(args []string) (flags, positional []string) {
	knownFlags := map[string]bool{
		"--dir":    true,
		"--target": true,
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

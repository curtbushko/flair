package main

import (
	"flag"
	"io"

	"github.com/curtbushko/flair/internal/config"
)

const validateUsage = `Usage: flair validate <theme-name> [options]

Validate a theme directory for completeness, schema correctness, and palette
validity. Prints any violations found and exits with code 1 if there are any.

Options:
  --dir <path>  Config directory (default: ~/.config/flair)
  -h, --help    Show this help message
`

// runValidate parses flags and executes the validate subcommand.
// It prints violations to stdout and errors/usage to stderr.
// Returns an exit code: 0 if valid, 1 if violations found or error.
func runValidate(args []string, stdout, stderr io.Writer) int {
	subArgs := args[2:]
	flagArgs, positional := splitValidateFlags(subArgs)

	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")

	if err := fs.Parse(flagArgs); err != nil {
		return 1
	}

	if len(positional) == 0 {
		writeStr(stderr, "error: theme argument is required\n\n")
		writeStr(stderr, validateUsage)
		return 1
	}

	themeName := positional[0]

	// Resolve config directory.
	configDir := *dirFlag
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	// Wire the app and create the validate use case.
	app := Wire(configDir)

	violations, err := app.Validate.Execute(themeName)
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	if len(violations) == 0 {
		writeStr(stdout, "theme %q is valid\n", themeName)
		return 0
	}

	// Print each violation on a separate line.
	for _, v := range violations {
		writeStr(stdout, "%s\n", v)
	}

	return 1
}

// splitValidateFlags separates flag arguments from positional arguments.
func splitValidateFlags(args []string) (flags, positional []string) {
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

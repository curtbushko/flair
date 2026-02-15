package main

import (
	"flag"
	"io"

	"github.com/curtbushko/flair/internal/adapters/store"
	"github.com/curtbushko/flair/internal/application"
	"github.com/curtbushko/flair/internal/config"
)

const initUsage = `Usage: flair init --name <theme-name> [options]

Scaffold a new theme directory with a starter palette.yaml.

Options:
  --name <name>  Theme name (required)
  --dir <path>   Config directory (default: ~/.config/flair)
  -h, --help     Show this help message
`

// runInit parses flags and executes the init subcommand.
// It writes confirmation to stdout and errors to stderr.
// Returns an exit code.
func runInit(args []string, stdout, stderr io.Writer) int {
	subArgs := args[2:]
	flagArgs := splitInitFlags(subArgs)

	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)

	nameFlag := fs.String("name", "", "theme name (required)")
	dirFlag := fs.String("dir", "", "config directory (default: ~/.config/flair)")

	if err := fs.Parse(flagArgs); err != nil {
		return 1
	}

	if *nameFlag == "" {
		writeStr(stderr, "error: --name flag is required\n\n")
		writeStr(stderr, initUsage)
		return 1
	}

	// Resolve config directory.
	configDir := *dirFlag
	if configDir == "" {
		configDir = config.DefaultConfigDir()
	}

	// Wire dependencies directly (init doesn't need the full App).
	fsStore := store.NewFsStore(configDir)
	uc := application.NewInitThemeUseCase(fsStore)

	palettePath, err := uc.Execute(*nameFlag)
	if err != nil {
		writeStr(stderr, "error: %v\n", err)
		return 1
	}

	writeStr(stdout, "Created %s\n", palettePath)
	return 0
}

// splitInitFlags separates flag arguments for the init subcommand.
func splitInitFlags(args []string) []string {
	knownFlags := map[string]bool{
		"--name": true,
		"--dir":  true,
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

		if knownFlags[arg] {
			flags = append(flags, arg)
			i++
			if i < len(args) {
				flags = append(flags, args[i])
				i++
			}
			continue
		}

		// Ignore unexpected positional args.
		i++
	}

	return flags
}

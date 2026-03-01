// Package main is the composition root for the flair CLI. It parses
// subcommands from os.Args and dispatches to the appropriate use case.
package main

import (
	"fmt"
	"io"
	"os"
)

const (
	flagHelp     = "--help"
	flagHelpShrt = "-h"
)

const usageText = `Usage: flair <command> [options]

Commands:
  list        List available built-in palettes
  preview     Preview a theme with ANSI colors
  regenerate  Re-derive stale downstream files
  select      Select a theme or launch style viewer
  validate    Validate a theme directory

Flags:
  -h, --help  Show this help message
`

func main() {
	code := run(os.Args, os.Stderr)
	os.Exit(code)
}

// run parses subcommands and dispatches to the appropriate handler.
// It writes diagnostics and usage to stderr and returns an exit code.
func run(args []string, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stderr)
		return 1
	}

	cmd := args[1]

	switch cmd {
	case flagHelp, flagHelpShrt:
		printUsage(stderr)
		return 0
	case "select":
		return runSelect(args, os.Stdout, stderr)
	case "list":
		return runList(args, os.Stdout, stderr)
	case "preview":
		return runPreview(args, os.Stdout, stderr)
	case "validate":
		return runValidate(args, os.Stdout, stderr)
	case "regenerate":
		return runRegenerate(args, os.Stdout, stderr)
	default:
		writeStr(stderr, "unknown command: %s\n\n", cmd)
		printUsage(stderr)
		return 1
	}
}

// printUsage writes the usage text to w.
func printUsage(w io.Writer) {
	_, _ = fmt.Fprint(w, usageText)
}

// writeStr writes a formatted string to w, discarding any write error.
func writeStr(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

package viewer

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
)

// RunOptions provides additional configuration for Run.
type RunOptions struct {
	// Input is the reader for keyboard input. If nil, os.Stdin is used.
	Input io.Reader

	// Output is the writer for TUI output. If nil, os.Stdout is used.
	Output io.Writer

	// WithAltScreen uses the alternate screen buffer (default: true).
	WithAltScreen bool

	// DryRun skips launching the TUI and returns immediately. Useful for
	// testing option validation without starting the interactive viewer.
	DryRun bool
}

// Run starts the style viewer TUI with the given options.
// It blocks until the user quits and returns any error encountered.
func Run(opts Options) error {
	return RunWithOptions(opts, RunOptions{WithAltScreen: true})
}

// RunWithOptions starts the style viewer TUI with additional run configuration.
// This is useful for testing or embedding in other applications.
func RunWithOptions(opts Options, runOpts RunOptions) error {
	model := NewModel(opts)

	var teaOpts []tea.ProgramOption

	if runOpts.WithAltScreen {
		teaOpts = append(teaOpts, tea.WithAltScreen())
	}

	if runOpts.Input != nil {
		teaOpts = append(teaOpts, tea.WithInput(runOpts.Input))
	}

	if runOpts.Output != nil {
		teaOpts = append(teaOpts, tea.WithOutput(runOpts.Output))
	}

	// In dry-run mode, skip launching the TUI.
	if runOpts.DryRun {
		return nil
	}

	p := tea.NewProgram(model, teaOpts...)

	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("viewer: %w", err)
	}

	return nil
}

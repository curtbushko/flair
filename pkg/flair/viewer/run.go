package viewer

import (
	"errors"
	"fmt"
	"io"

	tea "charm.land/bubbletea/v2"

	"github.com/curtbushko/flair/pkg/flair"
)

// RunOptions provides additional configuration for [Run] and [RunWithOptions].
//
// These options control I/O streams and display behavior of the TUI viewer.
type RunOptions struct {
	// Input is the reader for keyboard input. If nil, os.Stdin is used.
	Input io.Reader

	// Output is the writer for TUI output. If nil, os.Stdout is used.
	Output io.Writer

	// WithAltScreen uses the alternate screen buffer (default: true).
	// When true, the TUI runs in a separate screen and restores the
	// original terminal content on exit.
	WithAltScreen bool

	// DryRun skips launching the TUI and returns immediately. Useful for
	// testing option validation without starting the interactive viewer.
	DryRun bool
}

// RunBuiltinsOptions configures the [RunBuiltins] convenience function.
//
// This options struct combines run configuration with callback functions
// for theme selection and installation.
type RunBuiltinsOptions struct {
	// Input is the reader for keyboard input. If nil, os.Stdin is used.
	Input io.Reader

	// Output is the writer for TUI output. If nil, os.Stdout is used.
	Output io.Writer

	// WithAltScreen uses the alternate screen buffer (default: true).
	WithAltScreen bool

	// DryRun skips launching the TUI and returns immediately.
	DryRun bool

	// OnSelect is called when the user confirms a theme selection with Enter.
	// The callback receives the selected theme name.
	OnSelect func(name string)

	// OnInstall is called when the user confirms theme installation.
	// The callback should install the theme and return any error.
	// If nil, no installation action is taken.
	OnInstall func(name string) error

	// InitialTheme is the theme to pre-select on startup.
	// If empty, the first theme in the list is selected.
	InitialTheme string
}

// RunBuiltins starts the style viewer TUI with all built-in themes.
//
// This is a zero-config convenience function that requires no filesystem setup.
// It loads all themes from embedded palettes and does not require ~/.config/flair
// to exist.
//
// Example:
//
//	err := viewer.RunBuiltins(viewer.RunBuiltinsOptions{
//	    OnSelect: func(name string) {
//	        fmt.Printf("Selected: %s\n", name)
//	    },
//	    OnInstall: func(name string) error {
//	        store := flair.NewStore()
//	        return store.Install(name)
//	    },
//	})
func RunBuiltins(opts RunBuiltinsOptions) error {
	themes := flair.ListBuiltins()
	if len(themes) == 0 {
		return errors.New("no built-in themes available")
	}

	viewerOpts := Options{
		Themes:       themes,
		InitialTheme: opts.InitialTheme,
		OnSelect:     opts.OnSelect,
		OnInstall:    opts.OnInstall,
		ThemeLoader:  NewBuiltinThemeLoader(),
	}

	runOpts := RunOptions{
		Input:         opts.Input,
		Output:        opts.Output,
		WithAltScreen: opts.WithAltScreen,
		DryRun:        opts.DryRun,
	}

	// Default to alt screen if not explicitly disabled.
	if opts.Input == nil && opts.Output == nil && !opts.DryRun {
		runOpts.WithAltScreen = true
	}

	return RunWithOptions(viewerOpts, runOpts)
}

// Run starts the style viewer TUI with the given options.
//
// Run blocks until the user quits (via q, Escape, or Ctrl+C) and returns
// any error encountered during execution.
//
// The viewer displays a 2-panel layout with theme list on the left and
// preview content on the right. Use Tab to switch between preview pages.
//
// Example:
//
//	err := viewer.Run(viewer.Options{
//	    Themes:       []string{"tokyo-night-dark", "gruvbox-dark"},
//	    InitialTheme: "tokyo-night-dark",
//	    OnSelect: func(name string) {
//	        fmt.Printf("Selected: %s\n", name)
//	    },
//	})
func Run(opts Options) error {
	return RunWithOptions(opts, RunOptions{WithAltScreen: true})
}

// RunWithOptions starts the style viewer TUI with additional run configuration.
//
// This function provides full control over both viewer options and run-time
// configuration. It is useful for testing (with custom I/O) or embedding
// the viewer in other bubbletea applications.
func RunWithOptions(opts Options, runOpts RunOptions) error {
	model := NewModel(opts)

	// In bubbletea v2, alt screen is controlled via View.AltScreen field.
	model.altScreen = runOpts.WithAltScreen

	var teaOpts []tea.ProgramOption

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

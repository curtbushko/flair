# Bubbletea Example

This example demonstrates using flair with [bubbletea](https://github.com/charmbracelet/bubbletea) for building themed TUI applications.

## What it shows

- Loading a theme with `flair.MustLoad()`
- Creating themed styles with `flairbubble.NewStyles(theme)`
- Using themed spinners and progress bars
- Accessing raw theme colors for custom styles

## Running

```sh
# From the project root
go run ./examples/charm-bubbletea

# Or select a different theme first
flair select catppuccin-mocha
go run ./examples/charm-bubbletea
```

## Code highlights

```go
// Load theme and create styles
theme := flair.MustLoad()
styles := flairbubble.NewStyles(theme)

// Use in your View() method
spinner := styles.Spinner.Render(spinnerFrames[m.frame])
progress := styles.ProgressFilled.Render(strings.Repeat("█", filled))

// Access raw colors for custom styles
colors := theme.Colors()
customStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color(colors["status.success"].Hex()))
```

## Available styles

| Style | Description |
|-------|-------------|
| Title | Bold text for titles |
| Label | Standard label text |
| Spinner | Accent-colored spinner frames |
| ProgressFilled | Filled portion of progress bar |
| ProgressEmpty | Empty portion of progress bar |
| Help | Muted help text |
| Error | Error message styling |
| Success | Success message styling |

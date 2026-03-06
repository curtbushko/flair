# Bubbles Table Example

This example demonstrates using flair theme colors to style table UI elements.

## What it shows

- Loading a theme with `flair.MustLoad()`
- Accessing theme colors with `theme.Colors()`
- Creating styled table headers and cells
- Highlighting selected rows with theme colors

## Running

```sh
# From the project root
go run ./examples/charm-table

# Or select a different theme first
flair select gruvbox-dark
go run ./examples/charm-table
```

## Code highlights

```go
// Load theme and access colors
theme := flair.MustLoad()
colors := theme.Colors()

// Create header style
headerStyle := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color(colors["text.secondary"].Hex())).
    BorderBottom(true).
    BorderForeground(lipgloss.Color(colors["border.default"].Hex()))

// Create selected row style
selectedStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color(colors["text.inverse"].Hex())).
    Background(lipgloss.Color(colors["accent.primary"].Hex())).
    Bold(true)
```

## Useful color tokens for tables

| Token | Usage |
|-------|-------|
| `text.primary` | Cell text |
| `text.secondary` | Header text |
| `text.inverse` | Selected cell text |
| `accent.primary` | Selected row background |
| `border.default` | Table borders |

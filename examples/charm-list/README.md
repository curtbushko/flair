# Bubbles List Example

This example demonstrates using flair theme colors to style list-like UI elements.

## What it shows

- Loading a theme with `flair.MustLoad()`
- Accessing theme colors with `theme.Colors()`
- Creating styled list items with selection highlighting
- Using semantic color tokens for consistent theming

## Running

```sh
# From the project root
go run ./examples/charm-list

# Or select a different theme first
flair select nord
go run ./examples/charm-list
```

## Code highlights

```go
// Load theme and access colors
theme := flair.MustLoad()
colors := theme.Colors()

// Create styles using theme colors
itemStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color(colors["text.primary"].Hex()))

selectedStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color(colors["accent.primary"].Hex())).
    Background(lipgloss.Color(colors["surface.background.selection"].Hex())).
    Bold(true)
```

## Useful color tokens for lists

| Token | Usage |
|-------|-------|
| `text.primary` | Normal item text |
| `text.muted` | Item descriptions |
| `accent.primary` | Selected item text |
| `surface.background.selection` | Selected item background |
| `border.default` | List borders |

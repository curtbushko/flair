# Bubbles Viewport Example

This example demonstrates using flair theme colors to style content areas and panels.

## What it shows

- Loading a theme with `flair.MustLoad()`
- Accessing theme colors with `theme.Colors()`
- Creating bordered content panels with themed colors
- Using syntax colors for code-like content

## Running

```sh
# From the project root
go run ./examples/charm-viewport

# Or select a different theme first
flair select catppuccin-mocha
go run ./examples/charm-viewport
```

## Code highlights

```go
// Load theme and access colors
theme := flair.MustLoad()
colors := theme.Colors()

// Create content panel style
contentStyle := lipgloss.NewStyle().
    Background(lipgloss.Color(colors["surface.background.sunken"].Hex())).
    Foreground(lipgloss.Color(colors["text.primary"].Hex())).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color(colors["border.default"].Hex())).
    Padding(1, 2)

// Access theme metadata
fmt.Printf("Name: %s\n", theme.Name())
fmt.Printf("Variant: %s\n", theme.Variant())
```

## Useful color tokens for viewports

| Token | Usage |
|-------|-------|
| `surface.background.sunken` | Inset content areas |
| `text.primary` | Content text |
| `text.secondary` | Headers and labels |
| `border.default` | Panel borders |
| `syntax.*` | Code highlighting |

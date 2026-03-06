# Lipgloss Example

This example demonstrates using flair with [lipgloss](https://github.com/charmbracelet/lipgloss) for styled terminal output.

## What it shows

- Loading a theme with `flair.MustLoad()`
- Creating pre-configured styles with `flairlip.NewStyles(theme)`
- Using text styles (primary, secondary, muted)
- Using status styles (error, warning, success, info)
- Using surface styles (background, raised, sunken)
- Using component styles (buttons, list items)

## Running

```sh
# From the project root
go run ./examples/charm-lipgloss

# Or select a different theme first
flair select gruvbox-dark
go run ./examples/charm-lipgloss
```

## Code highlights

```go
// Load theme and create styles
theme := flair.MustLoad()
styles := flairlip.NewStyles(theme)

// Use pre-configured styles
fmt.Println(styles.Text.Render("Primary text"))
fmt.Println(styles.Error.Render("Error message"))
fmt.Println(styles.Button.Render(" Submit "))
```

## Available styles

| Category  | Styles |
|-----------|--------|
| Surface   | Background, Raised, Sunken, Overlay, Popup |
| Text      | Text, Secondary, Muted, Inverse |
| Status    | Error, Warning, Success, Info |
| Component | Button, ButtonFocused, ListItem, ListSelected |

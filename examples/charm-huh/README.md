# Huh Example

This example demonstrates using flair with [huh](https://github.com/charmbracelet/huh) for themed forms.

## What it shows

- Loading a theme with `flair.MustLoad()`
- Creating a huh theme with `flairhuh.NewTheme(theme)`
- Using themed inputs, selects, and confirms
- Applying the theme to a form with `.WithTheme()`

## Running

```sh
# From the project root
go run ./examples/charm-huh

# Or select a different theme first
flair select dracula
go run ./examples/charm-huh
```

## Code highlights

```go
// Load theme and create huh theme
theme := flair.MustLoad()
huhTheme := flairhuh.NewTheme(theme)

// Create a form with the theme
form := huh.NewForm(
    huh.NewGroup(
        huh.NewInput().
            Title("What's your name?").
            Value(&name),

        huh.NewSelect[string]().
            Title("Choose an option").
            Options(
                huh.NewOption("Option A", "a"),
                huh.NewOption("Option B", "b"),
            ).
            Value(&choice),
    ),
).WithTheme(huhTheme)

form.Run()
```

## Themed components

The flair huh theme styles all huh components:

- **Input** - Text inputs with themed cursor and placeholder
- **Text** - Multi-line text areas
- **Select** - Single-select dropdowns with themed indicators
- **MultiSelect** - Multi-select with themed checkboxes
- **Confirm** - Yes/No confirmations with themed buttons
- **Note** - Informational notes

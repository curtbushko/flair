# Flair

A directory-based theme pipeline for generating consistent color themes across
multiple targets. Each theme lives in its own directory under
`~/.config/flair/<theme-name>/` and flows through a four-stage pipeline:
palette, universal tokens, target-specific mappings, and final output files.

Flair also provides Go packages for building themed CLI applications with
[Charmbracelet](https://charm.sh) libraries (lipgloss, bubbletea, huh, bubbles).

## Installation

### From source

```sh
go install github.com/curtbushko/flair/cmd/flair@latest
```

### Build locally

```sh
git clone https://github.com/curtbushko/flair.git
cd flair
make build
# binary is in ./bin/flair
```

---

## CLI Usage

The flair CLI generates, manages, and previews color themes.

### Quick Start

```sh
# Select a built-in theme (auto-generates if needed)
flair select tokyo-night-dark

# Preview the theme in your terminal
flair preview tokyo-night-dark

# List available themes
flair list

# Launch interactive style viewer
flair select
```

### Commands

| Command      | Description                                          |
|--------------|------------------------------------------------------|
| `select`     | Select a theme as active (creates symlinks)          |
| `regenerate` | Re-derive downstream files from modified sources     |
| `validate`   | Check theme for completeness and correctness         |
| `preview`    | Preview theme with ANSI colors in terminal           |
| `list`       | List available or built-in themes                    |
| `override`   | Manage token overrides in palette.yaml               |

### select

Select a theme as active by creating symlinks to its output files. If the theme
is a built-in and hasn't been generated yet, it will be auto-generated first.

```sh
# Select a specific theme
flair select gruvbox-dark

# Launch interactive style viewer (no arguments)
flair select

# Force viewer mode with a theme pre-selected
flair select --viewer tokyo-night-dark
```

### regenerate

Re-derive downstream files by inspecting modification times. Only stale files
are regenerated.

```sh
# Regenerate all stale files for a theme
flair regenerate tokyo-night-dark

# Regenerate only the CSS target
flair regenerate tokyo-night-dark --target css
```

### validate

Validate a theme directory for completeness, schema correctness, and palette
validity.

```sh
flair validate tokyo-night-dark
```

### preview

Preview a theme with ANSI colors in the terminal.

```sh
flair preview catppuccin-mocha
```

### list

List available themes or built-in palette names.

```sh
# List installed themes (selected theme marked with *)
flair list

# List built-in palettes
flair list --builtins
```

### override

Manage token overrides without manually editing YAML.

```sh
# Add/update an override
flair override mytheme syntax.keyword "#ff00ff" --bold --italic

# List current overrides
flair override mytheme --list

# Remove an override
flair override mytheme --remove syntax.keyword
```

---

## Using Flair with Charm Packages

Flair provides themed adapters for all major Charmbracelet packages. Load a
theme once and use it across your entire TUI application.

### Supported Packages

| Package | Import Path | Description |
|---------|-------------|-------------|
| lipgloss | `pkg/charm/lipgloss` | Pre-configured styles for text rendering |
| bubbletea | `pkg/charm/bubbletea` | Themed styles for bubbletea components |
| huh | `pkg/charm/huh` | Themed forms (inputs, selects, confirms) |
| bubbles/list | `pkg/charm/bubbles/list` | Themed list component |
| bubbles/table | `pkg/charm/bubbles/table` | Themed table component |
| bubbles/viewport | `pkg/charm/bubbles/viewport` | Themed viewport/scrollable area |

### Basic Usage Pattern

All flair charm packages follow the same pattern:

```go
import (
    "github.com/curtbushko/flair/pkg/flair"
    flairlip "github.com/curtbushko/flair/pkg/charm/lipgloss"
)

func main() {
    // 1. Load a theme
    theme := flair.MustLoad()

    // 2. Create themed styles
    styles := flairlip.NewStyles(theme)

    // 3. Use the styles
    fmt.Println(styles.Text.Render("Hello, themed world!"))
}
```

---

### lipgloss - Text Styling

The lipgloss adapter provides pre-configured styles for common UI patterns.

```go
import (
    "github.com/curtbushko/flair/pkg/flair"
    flairlip "github.com/curtbushko/flair/pkg/charm/lipgloss"
)

func main() {
    theme := flair.MustLoad()
    styles := flairlip.NewStyles(theme)

    // Surface styles for containers
    header := styles.Raised.Render("My Application")
    content := styles.Background.Render("Main content")

    // Text styles
    title := styles.Text.Render("Primary text")
    hint := styles.Muted.Render("Hint text")

    // Status indicators
    fmt.Println(styles.Error.Render("Error: something went wrong"))
    fmt.Println(styles.Success.Render("Success: operation complete"))
    fmt.Println(styles.Warning.Render("Warning: proceed with caution"))
    fmt.Println(styles.Info.Render("Info: FYI"))

    // Component styles
    button := styles.Button.Render("Submit")
    focused := styles.ButtonFocused.Render("Submit")
}
```

**Available Styles:**

| Category  | Styles |
|-----------|--------|
| Surface   | Background, Raised, Sunken, Overlay, Popup |
| Text      | Text, Secondary, Muted, Inverse |
| Status    | Error, Warning, Success, Info |
| Border    | Border, BorderFocus, BorderMuted |
| Component | Button, ButtonFocused, Input, InputFocused, ListItem, ListSelected, Table, TableHeader, Dialog |
| State     | Hover, Active, Disabled, Selected |

See [examples/charm-lipgloss](examples/charm-lipgloss) for a complete example.

---

### bubbletea - TUI Applications

The bubbletea adapter provides themed styles for building TUI applications.

```go
import (
    "github.com/curtbushko/flair/pkg/flair"
    flairbubble "github.com/curtbushko/flair/pkg/charm/bubbletea"
)

func main() {
    theme := flair.MustLoad()
    styles := flairbubble.NewStyles(theme)

    // Use styles in your bubbletea model's View()
    view := styles.Title.Render("My App") + "\n"
    view += styles.Spinner.Render("Loading...")
    view += styles.ProgressFilled.Render("████") +
            styles.ProgressEmpty.Render("░░░░")
}
```

See [examples/charm-bubbletea](examples/charm-bubbletea) for a complete example.

---

### huh - Forms

The huh adapter creates themed form components.

```go
import (
    "github.com/charmbracelet/huh"
    "github.com/curtbushko/flair/pkg/flair"
    flairhuh "github.com/curtbushko/flair/pkg/charm/huh"
)

func main() {
    theme := flair.MustLoad()
    huhTheme := flairhuh.NewTheme(theme)

    var name string
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Name").
                Value(&name),
        ),
    ).WithTheme(huhTheme)

    form.Run()
}
```

See [examples/charm-huh](examples/charm-huh) for a complete example.

---

### bubbles/list - Lists

The list adapter provides themed delegates and styles for the bubbles list component.

```go
import (
    "github.com/charmbracelet/bubbles/list"
    "github.com/curtbushko/flair/pkg/flair"
    flairlist "github.com/curtbushko/flair/pkg/charm/bubbles/list"
)

func main() {
    theme := flair.MustLoad()

    // Create a themed delegate for item rendering
    delegate := flairlist.NewDelegate(theme)

    items := []list.Item{
        item{title: "Item 1", desc: "First item"},
        item{title: "Item 2", desc: "Second item"},
    }

    l := list.New(items, delegate, 40, 20)
    l.Styles = flairlist.NewStyles(theme)
}
```

See [examples/charm-list](examples/charm-list) for a complete example.

---

### bubbles/table - Tables

The table adapter provides themed styles for the bubbles table component.

```go
import (
    "github.com/charmbracelet/bubbles/table"
    "github.com/curtbushko/flair/pkg/flair"
    flairtable "github.com/curtbushko/flair/pkg/charm/bubbles/table"
)

func main() {
    theme := flair.MustLoad()
    styles := flairtable.NewStyles(theme)

    t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithStyles(styles),
    )
}
```

See [examples/charm-table](examples/charm-table) for a complete example.

---

### bubbles/viewport - Scrollable Content

The viewport adapter provides themed styles for scrollable content areas.

```go
import (
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/curtbushko/flair/pkg/flair"
    flairviewport "github.com/curtbushko/flair/pkg/charm/bubbles/viewport"
)

func main() {
    theme := flair.MustLoad()

    // Create a themed viewport
    vp := flairviewport.NewModel(theme, 80, 24)
    vp.SetContent(longContent)
}
```

See [examples/charm-viewport](examples/charm-viewport) for a complete example.

---

## Using Flair as a Library

Flair can be used as a Go library without installing the CLI. All built-in
themes are embedded, so external programs can load themes with zero setup.

### Zero-setup with built-in themes

```go
import "github.com/curtbushko/flair/pkg/flair"

func main() {
    // Load a built-in theme directly - no CLI, no config needed
    theme, err := flair.LoadBuiltin("tokyo-night-dark")
    if err != nil {
        log.Fatal(err)
    }

    // Access theme colors
    colors := theme.Colors()
    fmt.Println("Background:", colors["surface.background"].Hex())
}
```

### Discovering built-in themes

```go
// List all available built-in theme names
names := flair.ListBuiltins()

// Check if a specific built-in exists
if flair.HasBuiltin("gruvbox-dark") {
    theme, _ := flair.LoadBuiltin("gruvbox-dark")
}
```

### Convenience functions

```go
// Default() - loads selected theme or falls back to tokyo-night-dark
theme, err := flair.Default()

// MustLoad() - same as Default() but panics on error
theme := flair.MustLoad()

// LoadOrDefault() - try a specific theme, fall back to a built-in
theme, err := flair.LoadOrDefault("my-custom-theme", "gruvbox-dark")

// EnsureInstalled() - install built-ins to config dir if empty
err := flair.EnsureInstalled()
```

### The Store type

For programmatic theme management:

```go
store := flair.NewStore()

// Install and select themes
store.Install("tokyo-night-dark")
store.InstallAll()
store.Select("tokyo-night-dark")

// Load themes
theme, _ := store.Load()           // currently selected
theme, _ := store.LoadNamed("gruvbox-dark")

// Query
themes, _ := store.List()
selected, _ := store.Selected()
```

---

## Customizing with Overrides

Override specific semantic tokens directly in `palette.yaml`:

```yaml
system: "base24"
name: "Tokyo Night Dark Custom"
palette:
  base00: "1a1b26"
  # ... palette colors ...

overrides:
  syntax.keyword:
    color: "#ff00ff"
    italic: true
  status.error:
    color: "#ff0000"
    bold: true
```

### Supported token paths

| Category   | Token paths |
|------------|-------------|
| Surface    | `surface.*` -- backgrounds, raised, sunken, overlay |
| Text       | `text.*` -- primary, secondary, muted, inverse |
| Status     | `status.*` -- error, warning, success, info, hint |
| Syntax     | `syntax.*` -- keyword, string, comment, function |
| Statusline | `statusline.*` -- a/b/c sections with fg/bg |

### Override properties

| Property | Type | Description |
|----------|------|-------------|
| `color` | string | Hex color (e.g., `"#ff00ff"`) |
| `bold` | boolean | Bold text style |
| `italic` | boolean | Italic text style |
| `underline` | boolean | Underline text style |

---

## Pipeline

Flair processes themes through a four-stage pipeline:

```
palette.yaml --> universal.yaml --> *-mapping.yaml --> style.*
  (input)         (derived)         (per-target)      (final output)
```

1. **Palette** (`palette.yaml`) -- Base24 color palette
2. **Universal tokens** (`universal.yaml`) -- ~87 semantic tokens
3. **Target mappings** (`*-mapping.yaml`) -- Per-target mappings
4. **Output files** (`style.*`) -- Final output files

## Directory Layout

```
~/.config/flair/
  style.lua --------> tokyo-night-dark/style.lua    # symlinks
  style.css --------> tokyo-night-dark/style.css

  tokyo-night-dark/
    palette.yaml           # Layer A: base24 palette
    universal.yaml         # Layer B: semantic tokens
    vim-mapping.yaml       # Layer C: target mappings
    css-mapping.yaml
    style.lua              # Layer D: final outputs
    style.css
```

## Supported Targets

| Target | Output file | Description |
|--------|-------------|-------------|
| Vim/Neovim | `style.lua` | Lua colorscheme |
| CSS | `style.css` | Custom properties and rules |
| GTK | `gtk.css` | GTK @define-color + rules |
| QSS | `style.qss` | Qt Style Sheets |
| Stylix | `style.json` | Key-value pairs for Stylix/NixOS |

## Built-in Palettes

Flair ships with 189 built-in palettes including:

- tokyo-night-dark, tokyo-night-storm, tokyo-night-moon
- gruvbox-dark, gruvbox-light, gruvbox-material
- catppuccin-mocha, catppuccin-macchiato, catppuccin-frappe, catppuccin-latte
- dracula, nord, one-dark, solarized-dark
- and many more...

List them with:

```sh
flair list --builtins
```

## Architecture

Flair follows hexagonal architecture (ports and adapters):

```
internal/
  domain/          # Entities, value objects (no external deps)
  ports/           # Interface definitions
  application/     # Use cases
  adapters/        # External implementations
cmd/flair/         # CLI entry point
pkg/               # Public packages for library use
  flair/           # Theme loading and store
  charm/           # Charmbracelet integrations
```

## License

See [LICENSE](LICENSE) for details.

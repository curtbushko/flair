# Flair

A directory-based theme pipeline for generating consistent color themes across
multiple targets. Each theme lives in its own directory under
`~/.config/flair/<theme-name>/` and flows through a four-stage pipeline:
palette, universal tokens, target-specific mappings, and final output files.

Edit any intermediate YAML file, run regenerate, and everything downstream
updates automatically.

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

## Using Flair in Your CLI

Flair provides Go packages for building themed CLI applications with
[lipgloss](https://github.com/charmbracelet/lipgloss).

### Import the packages

```go
import (
    "github.com/curtbushko/flair/pkg/flair"
    "github.com/curtbushko/flair/pkg/charm/lipgloss"
)
```

### Load a theme and create styles

```go
// Load the currently selected theme (or fall back to tokyo-night-dark)
theme, err := flair.Default()
if err != nil {
    log.Fatal(err)
}

// Create pre-configured lipgloss styles from the theme
styles := lipgloss.NewStyles(theme)
```

### Use the pre-configured styles

```go
// Surface styles for containers
header := styles.Raised.Render("My Application")
content := styles.Background.Render("Main content area")

// Text styles
title := styles.Text.Render("Primary text")
subtitle := styles.Secondary.Render("Secondary text")
hint := styles.Muted.Render("Muted hint")

// Status indicators
errorMsg := styles.Error.Render("Something went wrong")
successMsg := styles.Success.Render("Operation complete")
warnMsg := styles.Warning.Render("Proceed with caution")
infoMsg := styles.Info.Render("FYI")

// Component styles for interactive elements
button := styles.Button.Render("Submit")
focused := styles.ButtonFocused.Render("Submit")
listItem := styles.ListItem.Render("Item 1")
selected := styles.ListSelected.Render("Item 2")
```

### Available style categories

| Category   | Styles                                          |
|------------|-------------------------------------------------|
| Surface    | Background, Raised, Sunken, Overlay, Popup      |
| Text       | Text, Secondary, Muted, Inverse                 |
| Status     | Error, Warning, Success, Info                   |
| Border     | Border, BorderFocus, BorderMuted                |
| Component  | Button, ButtonFocused, Input, InputFocused      |
|            | ListItem, ListSelected, Table, TableHeader, Dialog |
| State      | Hover, Active, Disabled, Selected               |

### Alternative loading methods

```go
// Load a specific theme by name
theme, err := flair.LoadNamed("gruvbox-dark")

// Load from the currently selected theme only (no fallback)
theme, err := flair.Load()

// Quick one-liner for styles (returns nil on error)
styles := lipgloss.Default()
```

## Using Flair as a Library

Flair can be used as a Go library without installing the CLI. All built-in
themes are embedded in the binary, so external programs can load and use
themes with zero filesystem setup.

### Zero-setup usage with built-in themes

```go
import "github.com/curtbushko/flair/pkg/flair"

func main() {
    // Load a built-in theme directly - no CLI, no config directory needed
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
// Returns: ["catppuccin-mocha", "gruvbox-dark", "tokyo-night-dark", ...]

// Check if a specific built-in exists
if flair.HasBuiltin("gruvbox-dark") {
    theme, _ := flair.LoadBuiltin("gruvbox-dark")
}
```

### Convenience functions

```go
// Default() - loads selected theme or falls back to tokyo-night-dark
// This is the recommended way to load a theme for most use cases
theme, err := flair.Default()

// MustLoad() - same as Default() but panics on error
// Useful for init-time loading where failure should be fatal
theme := flair.MustLoad()

// LoadOrDefault() - try a specific theme, fall back to a built-in
theme, err := flair.LoadOrDefault("my-custom-theme", "gruvbox-dark")

// EnsureInstalled() - install built-ins to config dir if empty
// Useful for first-run setup in CLI applications
err := flair.EnsureInstalled()
```

### The Store type for programmatic theme management

The `Store` type provides full control over theme installation and selection:

```go
import "github.com/curtbushko/flair/pkg/flair"

// Create a store using the default config directory (~/.config/flair)
store := flair.NewStore()

// Or use a custom directory
store := flair.NewStoreAt("/path/to/themes")

// Install a built-in theme to the config directory
err := store.Install("tokyo-night-dark")

// Install all built-in themes
err := store.InstallAll()

// Select a theme (creates symlinks at config root)
err := store.Select("tokyo-night-dark")

// Load the currently selected theme
theme, err := store.Load()

// Load a specific installed theme by name
theme, err := store.LoadNamed("gruvbox-dark")

// List all installed themes
themes, err := store.List()

// Get the currently selected theme name
selected, err := store.Selected()
```

### Example: themed CLI without the flair CLI

```go
package main

import (
    "fmt"
    "os"

    "github.com/curtbushko/flair/pkg/flair"
    "github.com/curtbushko/flair/pkg/charm/lipgloss"
)

func main() {
    // Load directly from built-ins - no flair CLI needed
    theme, err := flair.LoadBuiltin("catppuccin-mocha")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to load theme:", err)
        os.Exit(1)
    }

    // Create lipgloss styles from the theme
    styles := lipgloss.NewStyles(theme)

    // Use the styles in your CLI
    fmt.Println(styles.Text.Render("Welcome to my app!"))
    fmt.Println(styles.Success.Render("Operation successful"))
    fmt.Println(styles.Error.Render("Something went wrong"))
}
```

## Customizing with Overrides

Flair allows you to override specific semantic tokens directly in your
`palette.yaml` file. Overrides are applied after default tokenization, giving
you fine-grained control over individual tokens while inheriting the rest from
the palette.

### Override format in palette.yaml

Add an `overrides` section to your theme's `palette.yaml`:

```yaml
system: "base24"
name: "Tokyo Night Dark Custom"
author: "Your Name"
variant: "dark"
palette:
  base00: "1a1b26"
  base01: "1f2335"
  # ... palette colors ...

overrides:
  syntax.keyword:
    color: "#ff00ff"
    italic: true
  syntax.comment:
    color: "#666666"
  surface.background:
    color: "#000000"
  status.error:
    color: "#ff0000"
    bold: true
```

### Supported token paths

Any token path from the token inventory can be overridden:

| Category       | Token paths                                        |
|----------------|----------------------------------------------------|
| Surface        | `surface.*` -- backgrounds, raised, sunken, overlay |
| Text           | `text.*` -- primary, secondary, muted, inverse     |
| Status         | `status.*` -- error, warning, success, info, hint  |
| Syntax         | `syntax.*` -- keyword, string, comment, function, etc. |
| Markup         | `markup.*` -- heading, bold, italic, link, code    |
| Diff           | `diff.*` -- added, removed, changed                |
| Accent         | `accent.*` -- primary, secondary accents           |
| Border         | `border.*` -- default, focus, muted                |
| Terminal       | `terminal.*` -- ANSI colors                        |
| Git            | `git.*` -- added, modified, deleted, ignored       |
| State          | `state.*` -- hover, active, disabled, selected     |
| Scrollbar      | `scrollbar.*` -- track, thumb                      |
| Statusline     | `statusline.*` -- a/b/c sections with fg/bg        |

### Override properties

Each override can specify any combination of these properties:

| Property        | Type    | Description                    |
|-----------------|---------|--------------------------------|
| `color`         | string  | Hex color (e.g., `"#ff00ff"`)  |
| `bold`          | boolean | Bold text style                |
| `italic`        | boolean | Italic text style              |
| `underline`     | boolean | Underline text style           |
| `undercurl`     | boolean | Undercurl text style           |
| `strikethrough` | boolean | Strikethrough text style       |

### Managing overrides with the CLI

The `flair override` command provides a convenient way to manage overrides
without manually editing YAML files.

Add or update an override:

```sh
# Set color only
flair override mytheme syntax.keyword "#ff00ff"

# Set style flags only
flair override mytheme syntax.keyword --bold --italic

# Set both color and styles
flair override mytheme syntax.keyword "#ff00ff" --bold
```

List current overrides:

```sh
flair override mytheme --list
```

Remove an override:

```sh
flair override mytheme --remove syntax.keyword
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)
- `--bold`, `--italic`, `--underline`, `--undercurl`, `--strikethrough` -- Style flags
- `--list` -- List all overrides for the theme
- `--remove <token>` -- Remove an override

After modifying overrides, regenerate the theme to apply changes:

```sh
flair regenerate mytheme
```

## Quick Start

```sh
# Select a built-in theme (auto-generates if needed)
flair select tokyo-night-dark

# Preview the theme in your terminal
flair preview tokyo-night-dark

# List available themes
flair list
```

## Commands

### select

Select a theme as active by creating symlinks to its output files at the
config root. If the theme is a built-in and hasn't been generated yet, it
will be auto-generated first.

```sh
flair select <theme-name> [options]
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)

Example:

```sh
flair select gruvbox-dark
```

### regenerate

Re-derive downstream files by inspecting modification times. Only stale files
are regenerated: palette edits re-derive everything, universal edits re-map
all targets, and mapping edits re-generate only that output.

```sh
flair regenerate <theme-name> [options]
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)
- `--target <name>` -- Regenerate only the named target (e.g. `vim`)

Examples:

```sh
# Regenerate all stale files for a theme
flair regenerate tokyo-night-dark

# Regenerate only the CSS target
flair regenerate tokyo-night-dark --target css
```

### validate

Validate a theme directory for completeness, schema correctness, and palette
validity. Prints any violations found and exits with code 1 if invalid.

```sh
flair validate <theme-name> [options]
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)

Example:

```sh
flair validate tokyo-night-dark
```

### preview

Preview a theme with ANSI colors in the terminal.

```sh
flair preview <theme-name> [options]
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)

Example:

```sh
flair preview catppuccin-mocha
```

### list

List available themes or built-in palette names.

```sh
flair list [options]
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)
- `--builtins` -- List built-in palette names only

Examples:

```sh
# List installed themes (selected theme marked with *)
flair list

# List built-in palettes
flair list --builtins
```

## Pipeline

Flair processes themes through a four-stage pipeline. Each stage produces a
versioned YAML file that can be inspected and edited:

```
palette.yaml --> universal.yaml --> *-mapping.yaml --> style.*
  (input)         (derived)         (per-target)      (final output)
```

1. **Palette** (`palette.yaml`) -- A base24 color palette defining the raw
   colors for a theme.

2. **Universal tokens** (`universal.yaml`) -- Approximately 87 semantic tokens
   derived from the palette. These represent abstract concepts like
   "editor background" or "warning foreground" independent of any target.

3. **Target mappings** (`*-mapping.yaml`) -- Per-target files that map
   universal tokens to target-specific constructs (Vim highlight groups,
   CSS custom properties, GTK colors, etc.).

4. **Output files** (`style.*`) -- Final output files ready for consumption
   by applications.

## Directory Layout

```
~/.config/flair/
  style.lua --------> tokyo-night-dark/style.lua    # symlinks to selected theme
  style.css --------> tokyo-night-dark/style.css
  gtk.css ----------> tokyo-night-dark/gtk.css
  style.qss --------> tokyo-night-dark/style.qss
  style.json -------> tokyo-night-dark/style.json

  tokyo-night-dark/
    palette.yaml           # Layer A: base24 palette (input)
    universal.yaml         # Layer B: semantic tokens (derived)
    vim-mapping.yaml       # Layer C: target mappings
    css-mapping.yaml
    gtk-mapping.yaml
    qss-mapping.yaml
    stylix-mapping.yaml
    style.lua              # Layer D: final output files
    style.css
    gtk.css
    style.qss
    style.json

  gruvbox-dark/
    palette.yaml
    ...
```

## Supported Targets

| Target        | Mapping file           | Output file   | Description                  |
|---------------|------------------------|---------------|------------------------------|
| Vim / Neovim  | `vim-mapping.yaml`     | `style.lua`   | Lua colorscheme              |
| CSS           | `css-mapping.yaml`     | `style.css`   | Custom properties and rules  |
| GTK           | `gtk-mapping.yaml`     | `gtk.css`     | GTK @define-color + rules    |
| QSS           | `qss-mapping.yaml`     | `style.qss`   | Qt Style Sheets              |
| Stylix / NixOS| `stylix-mapping.yaml`  | `style.json`  | Key-value pairs for Stylix   |

## Bufferline Integration

Flair generates a `bufferline_theme` table in `style.lua` that can be used
with [bufferline.nvim](https://github.com/akinsho/bufferline.nvim).

The theme uses the same `statusline.*` tokens as lualine:
- Selected buffer: `statusline.a.*` (brightest)
- Visible buffers: `statusline.b.*`
- Background/hidden: `statusline.c.*`

### Usage

The generated theme is automatically applied if bufferline is installed:

```lua
-- In your Neovim config, load the colorscheme:
require('flair.style')  -- or vim.cmd('colorscheme flair')

-- The bufferline theme is applied automatically via pcall
```

### Manual Setup

If you prefer manual control:

```lua
-- style.lua exposes bufferline_theme as a local
-- You can copy the generated theme or customize it:
require('bufferline').setup({
  highlights = require('flair.bufferline_theme'),  -- if exported
  options = {
    -- your options
  },
})
```

### Customizing Statusline Tokens

Override `statusline.*` tokens in your palette to customize both
lualine and bufferline:

```yaml
tokens:
  statusline:
    a:
      fg: "#000000"
      bg: "#ffffff"
    b:
      fg: "#888888"
      bg: "#333333"
    c:
      fg: "#666666"
      bg: "#1a1a1a"
```

## Built-in Palettes

Flair ships with three built-in palettes:

- **tokyo-night-dark** -- A dark theme inspired by Tokyo Night
- **gruvbox-dark** -- A dark retro groove theme
- **catppuccin-mocha** -- The mocha variant of Catppuccin

List them with:

```sh
flair list --builtins
```

## Customization

Flair uses no override system. Instead, customize themes by editing
intermediate YAML files directly:

1. Select a theme: `flair select tokyo-night-dark`
2. Edit any intermediate file (e.g. `universal.yaml` to change semantic
   token assignments, or a `*-mapping.yaml` to adjust target-specific output)
3. Regenerate downstream files: `flair regenerate tokyo-night-dark`

The regenerate command inspects modification times and only re-derives files
that are stale relative to their upstream source. Edit the palette and
everything regenerates. Edit a mapping file and only that target's output
regenerates.

Every YAML file includes a `schema_version` field for forward compatibility.

## Architecture

Flair follows hexagonal architecture (ports and adapters), enforced by
`go-arch-lint`:

```
internal/
  domain/          # Entities, value objects, domain errors (no external deps)
  ports/           # Interface definitions (PaletteSource, ThemeStore, etc.)
  application/     # Use cases (Generate, Select, List, Validate, Preview, etc.)
  adapters/
    palettes/      # Built-in palette source (go:embed)
    store/         # Filesystem-based theme store
    yaml/          # YAML parser adapter
    tokenizer/     # Palette-to-universal token derivation
    mapper/        # Universal-to-target mapping (vim, css, gtk, qss, stylix)
    generator/     # Mapping-to-output file generation
    fileio/        # YAML file readers and writers
    wrappers/      # Decorators (schema validation, versioning)
  config/          # Runtime configuration
cmd/flair/         # CLI entry point and dependency wiring
```

Dependencies flow inward: adapters depend on ports, ports depend on domain.
The domain layer has no external dependencies.

## License

See [LICENSE](LICENSE) for details.

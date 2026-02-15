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

## Quick Start

```sh
# Generate a theme from a built-in palette
flair generate tokyo-night-dark

# Preview the theme in your terminal
flair preview tokyo-night-dark

# Select it as the active theme (creates symlinks)
flair select tokyo-night-dark
```

## Commands

### generate

Generate theme files from a palette. The palette argument can be a built-in
name or a file path to a palette YAML.

```sh
flair generate <palette> [options]
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)
- `--target <name>` -- Generate only the named target (e.g. `stylix`)
- `--name <name>` -- Override the theme name (default: inferred from palette)

Examples:

```sh
# Generate from a built-in palette
flair generate catppuccin-mocha

# Generate from a custom palette file
flair generate ~/palettes/my-palette.yaml --name my-theme

# Generate only the vim target
flair generate gruvbox-dark --target vim
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

### select

Select a theme as active by creating symlinks to its output files at the
config root.

```sh
flair select <theme-name> [options]
```

Options:

- `--dir <path>` -- Config directory (default: `~/.config/flair`)

Example:

```sh
flair select gruvbox-dark
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

### init

Scaffold a new theme directory with a starter palette.yaml.

```sh
flair init --name <theme-name> [options]
```

Options:

- `--name <name>` -- Theme name (required)
- `--dir <path>` -- Config directory (default: `~/.config/flair`)

Example:

```sh
flair init --name my-custom-theme
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

1. Generate a theme: `flair generate tokyo-night-dark`
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
    deriver/       # Palette-to-universal token derivation
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

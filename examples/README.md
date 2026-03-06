# Flair Examples

This directory contains example programs demonstrating how to use the flair theming library.

## Examples

### minimal-cli

Demonstrates basic CLI theming with flair.

```bash
go run ./examples/minimal-cli
```

Shows how to:
- Load a theme using `flair.Default()`
- Create lipgloss styles from the theme
- Print styled output to the terminal

### embed-viewer

Demonstrates embedding the flair style viewer in a CLI application.

```bash
go run ./examples/embed-viewer
```

Shows how to:
- Use `viewer.Run()` with custom options
- Provide an `OnSelect` callback to respond to theme selection
- Configure the viewer with a theme loader

### zero-setup

Demonstrates CLI theming with no configuration required.

```bash
go run ./examples/zero-setup
```

Shows how to:
- Use `flair.LoadBuiltin()` to load themes directly from embedded palettes
- Work without any `~/.config/flair` directory
- List available built-in themes

### programmatic

Demonstrates programmatic theme generation.

```bash
go run ./examples/programmatic
```

Shows how to:
- Parse a palette YAML file with `flair.ParsePalette()`
- Tokenize the palette to generate semantic tokens with `flair.Tokenize()`
- Access the generated theme's colors

### viewer-install

Demonstrates the viewer with install-on-select behavior.

```bash
go run ./examples/viewer-install
```

Shows how to:
- Use `viewer.RunBuiltins()` for zero-config viewer startup
- Provide `OnSelect` and `OnInstall` callbacks
- Install themes to `~/.config/flair` when selected

## Running Examples

All examples can be run from the repository root:

```bash
# Run any example
go run ./examples/<example-name>

# Build all examples
go build ./examples/...
```

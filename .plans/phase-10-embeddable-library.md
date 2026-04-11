# Phase 10: Embeddable Library (Zero CLI Dependency)

## Tasks

- [x] 10.1 — Built-in Themes in `pkg/flair`
  - [x] 10.1a — Move embedded palette YAML files to `pkg/flair/palettes/`
  - [x] 10.1b — Create `pkg/flair/builtins.go` with `ListBuiltins()` and `LoadBuiltin(name)`
  - [x] 10.1c — `LoadBuiltin()` returns a fully tokenized `*Theme` (no filesystem)
  - [x] 10.1d — Built-in themes work without any `~/.config/flair` setup
  - [x] 10.1e — Unit tests for built-in theme loading

- [x] 10.2 — Public Palette and Tokenizer in `pkg/flair`
  - [x] 10.2a — Create `pkg/flair/palette.go` — public Palette type (base24 colors)
  - [x] 10.2b — Create `pkg/flair/color.go` — public Color type with hex parsing
  - [x] 10.2c — Create `pkg/flair/tokenizer.go` — derive tokens from palette
  - [x] 10.2d — Create `pkg/flair/parse.go` — parse palette YAML from `io.Reader`
  - [x] 10.2e — Unit tests for palette parsing and tokenization

- [x] 10.3 — Public Theme Store in `pkg/flair`
  - [x] 10.3a — Create `pkg/flair/store.go` — read/write themes to `~/.config/flair`
  - [x] 10.3b — `SaveTheme(theme *Theme)` — write theme files to config dir
  - [x] 10.3c — `Select(name string)` — update symlinks to select a theme
  - [x] 10.3d — `Install(name string)` — install a built-in theme to config dir
  - [x] 10.3e — `InstallAll()` — install all built-in themes
  - [x] 10.3f — Unit tests for store operations

- [x] 10.4 — Convenience Functions
  - [x] 10.4a — `flair.Default()` — load selected theme, fallback to built-in if none
  - [x] 10.4b — `flair.MustLoad()` — panic variant for init-time loading
  - [x] 10.4c — `flair.LoadOrDefault(name, fallback string)` — try named, fallback to built-in
  - [x] 10.4d — `flair.EnsureInstalled()` — install built-ins if config dir empty
  - [x] 10.4e — Unit tests for convenience functions

- [x] 10.5 — Viewer Integration
  - [x] 10.5a — Update viewer to work with built-in themes (no filesystem required)
  - [x] 10.5b — `viewer.Run()` shows built-ins even if `~/.config/flair` doesn't exist
  - [x] 10.5c — Option to install theme on selection (`OnSelect` can call `Install`)
  - [x] 10.5d — Unit tests for viewer with built-ins

- [x] 10.6 — Internal Refactoring
  - [x] 10.6a — `internal/adapters/palettes/` imports from `pkg/flair/palettes/` (DRY)
  - [x] 10.6b — `internal/adapters/tokenizer/` delegates to `pkg/flair/tokenizer.go`
  - [x] 10.6c — CLI commands use `pkg/flair` where appropriate
  - [x] 10.6d — Ensure `go-arch-lint check` still passes
  - [x] 10.6e — No breaking changes to existing CLI behavior

- [x] 10.7 — Documentation
  - [x] 10.7a — README section: Using flair as a library (no CLI required)
  - [x] 10.7b — Example: CLI with built-in flair theming (zero setup)
  - [x] 10.7c — Example: programmatic theme generation
  - [x] 10.7d — Example: embedding viewer with "install on select"
  - [x] 10.7e — godoc comments for all new public APIs

## Notes

Makes flair fully usable as a Go library without requiring CLI installation.

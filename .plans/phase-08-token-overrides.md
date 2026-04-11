# Phase 8: Token Overrides

## Tasks

- [x] 8.1 — Domain: Token override types
  - [x] 8.1a — Define `TokenOverride` struct (color, bold, italic, underline, etc.)
  - [x] 8.1b — Add `Overrides map[string]TokenOverride` field to `domain.Palette`
  - [x] 8.1c — Unit tests for override parsing and application

- [x] 8.2 — Adapter: YAML parser updates
  - [x] 8.2a — Update `adapters/yaml/parser.go` to parse `overrides` section
  - [x] 8.2b — Validate override token paths against known token paths
  - [x] 8.2c — Validate override colors (hex format)
  - [x] 8.2d — Unit tests for override YAML parsing

- [x] 8.3 — Adapter: Tokenizer override application
  - [x] 8.3a — Update `Tokenizer.Tokenize()` to accept optional overrides
  - [x] 8.3b — Apply overrides after default derivation
  - [x] 8.3c — Merge override styles with derived styles (override wins)
  - [x] 8.3d — Unit tests for override application

- [x] 8.4 — Adapter: Palette writer updates
  - [x] 8.4a — Update `fileio.WritePalette()` to include overrides section
  - [x] 8.4b — Preserve user overrides during regeneration
  - [x] 8.4c — Unit tests for round-trip (parse → write → parse)

- [x] 8.5 — Application: Use case updates
  - [x] 8.5a — Update `GenerateThemeUseCase` to pass overrides to tokenizer
  - [x] 8.5b — Update `RegenerateThemeUseCase` to preserve overrides
  - [x] 8.5c — Update `PreviewThemeUseCase` to display overridden tokens
  - [x] 8.5d — Integration tests for override flow

- [x] 8.6 — CLI: Override-related commands
  - [x] 8.6a — `flair override <theme> <token> <color>` — Add/update override
  - [x] 8.6b — `flair override <theme> --list` — List current overrides
  - [x] 8.6c — `flair override <theme> --remove <token>` — Remove override
  - [x] 8.6d — Help text and documentation

- [x] 8.7 — Validation
  - [x] 8.7a — Validate override token paths exist in token inventory
  - [x] 8.7b — Warning for overrides that shadow derived values
  - [x] 8.7c — Update `validate` command to check overrides

- [x] 8.8 — Documentation & Polish
  - [x] 8.8a — README section: Customizing with overrides
  - [x] 8.8b — Example: Creating a custom theme with overrides
  - [x] 8.8c — BDD feature files for override scenarios

## Notes

Allows users to override specific semantic tokens directly in palette.yaml for simple customization.

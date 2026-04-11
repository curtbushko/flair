# Phase 2: Layer B — Token Derivation

## Tasks

- [x] 2.1 — Adapter: Default tokenizer (base24 → ~87 semantic tokens)
  - [x] 2.1a — Surface tokens (11 tokens)
  - [x] 2.1b — Text tokens (7 tokens)
  - [x] 2.1c — Status tokens (6 tokens)
  - [x] 2.1d — Diff tokens (9 tokens)
  - [x] 2.1e — Syntax tokens (14 tokens)
  - [x] 2.1f — Markup tokens (10 tokens)
  - [x] 2.1g — Accent, border, scrollbar, state tokens (11 tokens)
  - [x] 2.1h — Git tokens (4 tokens)
  - [x] 2.1i — Terminal ANSI colors (16 tokens)
- [x] 2.2 — Adapter: UniversalFile writer (TokenSet → io.Writer as YAML)
- [x] 2.3 — Adapter: UniversalFile reader (io.Reader → TokenSet)
- [x] 2.4 — Application: DeriveTheme use case (io.Reader → universal.yaml via io.Writer)
- [x] 2.5 — Unit tests for derivation rules against Tokyo Night Dark palette

## Notes

Implements the tokenization layer that converts base24 palette colors into semantic tokens.

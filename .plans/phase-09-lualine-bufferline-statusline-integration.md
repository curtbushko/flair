# Phase 9: Lualine & Bufferline Statusline Integration

## Tasks

- [x] 9.1 — Port types for bufferline
  - [x] 9.1a — Define `BufferlineTheme` struct in `ports/themes.go`
  - [x] 9.1b — Define `BufferlineColors` struct (fg, bg, sp for each state)
  - [x] 9.1c — Add `Bufferline *BufferlineTheme` field to `VimTheme`
  - [x] 9.1d — Unit tests for bufferline port types

- [x] 9.2 — Mapper: bufferline mapping
  - [x] 9.2a — Add `mapBufferline()` function in `internal/adapters/mapper/vim.go`
  - [x] 9.2b — Map selected buffer to `statusline.a.*` tokens (brightest)
  - [x] 9.2c — Map visible buffers to `statusline.b.*` tokens
  - [x] 9.2d — Map background/hidden to `statusline.c.*` tokens
  - [x] 9.2e — Map separator colors from `border.default`
  - [x] 9.2f — Map indicator colors from `accent.primary`
  - [x] 9.2g — Map modified indicator from `status.warning`
  - [x] 9.2h — Map diagnostic counts from `status.error/warning/info/hint`
  - [x] 9.2i — Unit tests for bufferline mapping

- [x] 9.3 — Generator: bufferline output
  - [x] 9.3a — Add `generateBufferline()` function in `internal/adapters/generator/vim.go`
  - [x] 9.3b — Output `local bufferline_theme = {...}` table
  - [x] 9.3c — Output `bufferline.setup({ highlights = bufferline_theme })` via pcall
  - [x] 9.3d — Unit tests for bufferline Lua generation

- [x] 9.4 — Mapping file support
  - [x] 9.4a — Update `VimMappingFile` to include bufferline section
  - [x] 9.4b — Update `fileio.WriteVimMapping()` to serialize bufferline
  - [x] 9.4c — Update `fileio.ReadVimMapping()` to parse bufferline
  - [x] 9.4d — Unit tests for round-trip (write → read)

- [x] 9.5 — Golden file updates
  - [x] 9.5a — Update `testdata/expected/style.lua` with bufferline theme
  - [x] 9.5b — Update `testdata/expected/vim-mapping.yaml` with bufferline section
  - [x] 9.5c — Regenerate golden files with `-update` flag

- [x] 9.6 — Documentation
  - [x] 9.6a — Document bufferline theming in README
  - [x] 9.6b — Add example bufferline.nvim setup snippet
  - [x] 9.6c — Document statusline.* token overrides for customization

## Notes

Ensures lualine and bufferline plugins use statusline.* tokens for consistent theming.

# Phase 3: Layer C+D — Mapping + Generation (per target)

## Tasks

- [x] 3.1 — Adapter: Stylix mapper + generator (simplest, validates full pipeline)
  - [x] 3.1a — Mapper: ResolvedTheme → ports.StylixTheme (≥60 keys)
  - [x] 3.1b — MappingFile writer: StylixTheme → stylix-mapping.yaml
  - [x] 3.1c — MappingFile reader: stylix-mapping.yaml → StylixTheme
  - [x] 3.1d — Generator: StylixTheme → style.json (sorted, 2-space indent)
- [x] 3.2 — Adapter: CSS mapper + generator
  - [x] 3.2a — Mapper: custom properties + element rules
  - [x] 3.2b — MappingFile writer/reader
  - [x] 3.2c — Generator: :root{} + element selectors → style.css
- [x] 3.3 — Adapter: Vim mapper + generator (most complex)
  - [x] 3.3a — Mapper: base highlights (Normal, Comment, Visual, CursorLine, etc.)
  - [x] 3.3b — Mapper: treesitter highlights (@keyword, @string, @function, etc.)
  - [x] 3.3c — Mapper: LSP semantic token links
  - [x] 3.3d — Mapper: diagnostic highlights (virtual text, underlines)
  - [x] 3.3e — Mapper: plugin highlights (telescope, gitsigns, etc.)
  - [x] 3.3f — Mapper: markup highlights
  - [x] 3.3g — Mapper: terminal ANSI colors (16)
  - [x] 3.3h — MappingFile writer/reader
  - [x] 3.3i — Generator: .lua output (hi clear, nvim_set_hl, links, terminal)
- [x] 3.4 — Adapter: GTK mapper + generator
  - [x] 3.4a — Mapper: @define-color definitions + widget selector rules
  - [x] 3.4b — MappingFile writer/reader
  - [x] 3.4c — Generator: CSS output (@define-color then selectors) → gtk.css
- [x] 3.5 — Adapter: QSS mapper + generator
  - [x] 3.5a — Mapper: widget + pseudo-state rules
  - [x] 3.5b — MappingFile writer/reader
  - [x] 3.5c — Generator: literal hex, no variables → style.qss
- [x] 3.6 — Application: GenerateTheme use case (full pipeline or partial regeneration)

## Notes

Implements mappers and generators for all target formats: Stylix, CSS, Vim, GTK, and QSS.

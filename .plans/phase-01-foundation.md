# Phase 1: Foundation

## Tasks

- [x] 1.1 — Project scaffolding (go module, directory structure, Makefile, `.go-arch-lint.yml`)
- [x] 1.2 — Domain: Color value object and parsing (hex ↔ RGB ↔ HSL)
- [x] 1.3 — Domain: Color operations (Blend, BlendBg, Lighten, Darken, Desaturate, ShiftHue)
- [x] 1.4 — Domain: Palette entity (base24 struct, slot access, base16 fallbacks)
- [x] 1.5 — Domain: Token value object and TokenSet aggregate
- [x] 1.6 — Domain: ResolvedTheme aggregate
- [x] 1.7 — Domain: Palette validation rules (luminance ordering, completeness)
- [x] 1.8 — Domain: Error types (ParseError, ValidationError, GenerateError, SchemaVersionError)
- [x] 1.9 — Domain: Schema version constants and file type registry
- [x] 1.10 — Port interfaces (PaletteParser, PaletteSource, Tokenizer, Mapper, Generator, ThemeStore)
- [x] 1.11 — Port file structs (PaletteFile, UniversalFile, VimMappingFile, etc.)
- [x] 1.12 — Port theme structs (VimTheme, GtkTheme, QssTheme, CssTheme, StylixTheme)
- [x] 1.13 — Adapter: YAML palette parser (io.Reader → domain.Palette, common tinted-theming format only)
- [x] 1.14 — Adapter: ThemeStore (filesystem — read/write theme dirs, symlink management)
- [x] 1.15 — Adapter: Built-in palettes (//go:embed, PaletteSource impl, Get returns io.Reader)
  - [x] 1.15a — tokyo-night-dark.yaml
  - [x] 1.15b — gruvbox-dark.yaml
  - [x] 1.15c — catppuccin-mocha.yaml
  - [x] 1.15d — andromeda.yaml
  - [x] 1.15e — everforest.yaml
  - [x] 1.15f — gruvbox-material.yaml
  - [x] 1.15g — rebel-scum.yaml
  - [x] 1.15h — tokyo-night-neon.yaml
- [x] 1.16 — Adapter: VersionedWriter (wraps io.Writer, prepends schema_version + kind header)
- [x] 1.17 — Adapter: ValidatingReader (wraps io.Reader, peeks schema version, returns SchemaVersionError)
- [x] 1.18 — Testdata: Tokyo Night Dark reference palette YAML

## Notes

Foundation phase establishes the hexagonal architecture, domain models, and core infrastructure.

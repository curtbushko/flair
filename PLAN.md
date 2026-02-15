# Flair — Implementation Plan

**Language:** Go
**Working title:** `flair`
**Architecture:** Hexagonal (Ports & Adapters), enforced by `go-arch-lint`
**Module path:** `github.com/curtbushko/flair`

---

## Core Concept

Flair is a directory-based theme pipeline. Each theme lives in its own directory
under `~/.config/flair/<theme-name>/`. The pipeline produces a chain of
intermediate files — each one a versioned Go struct serialized to YAML — from
palette through semantic tokens through target-specific mappings to final output
files. Users can customize at any layer by editing the intermediate YAML, and
the tool regenerates everything downstream.

The currently selected theme is exposed via symlinks at the flair config root,
so applications only need to point at a single stable path.

### Pipeline

```
palette.yaml ──► universal.yaml ──► *-mapping.yaml ──► style.*
  (input)         (derived)          (per-target)       (final output)
```

### Directory Layout (runtime)

```
~/.config/flair/
├── style.lua ──────► tokyonight/style.lua       # symlinks to selected theme
├── style.css ──────► tokyonight/style.css
├── gtk.css ────────► tokyonight/gtk.css
├── style.qss ─────► tokyonight/style.qss
├── style.json ────► tokyonight/style.json
│
├── tokyonight/                                   # one directory per theme
│   ├── palette.yaml           # Layer A: base24 palette (input)
│   ├── universal.yaml         # Layer B: ~87 semantic tokens (derived)
│   ├── vim-mapping.yaml       # Layer C: Vim highlight groups (mapped)
│   ├── css-mapping.yaml       #          CSS custom props + rules
│   ├── gtk-mapping.yaml       #          GTK @define-color + selectors
│   ├── qss-mapping.yaml       #          QSS widget rules
│   ├── stylix-mapping.yaml    #          Stylix key-value pairs
│   ├── style.lua              # Layer D: final Neovim colorscheme
│   ├── style.css              #          final CSS
│   ├── gtk.css                #          final GTK CSS
│   ├── style.qss              #          final QSS
│   └── style.json             #          final Stylix JSON
│
├── gruvbox/
│   ├── palette.yaml
│   └── ...
```

### Schema Versioning

Every YAML file produced by flair includes a `schema_version` field. This
allows the tool to detect stale files and re-derive them when derivation rules,
mappings, or output formats change across flair releases.

```yaml
schema_version: 1
# ... rest of file
```

When flair encounters a file with a schema version older than its current
version for that file type, it logs a notice and regenerates. When it encounters
a version newer than it understands, it errors with a "please upgrade flair"
message.

### Customization Model

No override system. Instead, users copy a theme directory and edit any
intermediate file directly. Running `flair generate` on a directory
re-derives downstream files from the furthest-upstream file that was
modified. The intermediate YAML files are the customization surface.

### Go Patterns: Reader/Writer + Embedding for Composition

Two idiomatic Go patterns shape how code is written within the hex layers:

**1. io.Reader / io.Writer as the universal seam.** Every function that
consumes or produces bytes takes an `io.Reader` or `io.Writer` — never a
file path, never `[]byte`. This decouples all transformations from their
data source: a `PaletteParser` works identically on a file, an embedded
built-in, a test literal, or stdin. The composition root (cmd) is the only
place that opens files and creates readers.

**2. Embedding for composition (decorator/wrapper pattern).** Behavior is
layered by wrapping readers and writers rather than modifying them. A
`VersionedWriter` wraps any `io.Writer` and prepends the schema version
header — so generators never know about versioning. A `ValidatingReader`
wraps any `io.Reader`, peeks at the schema version, and returns a
`SchemaVersionError` if incompatible — so every file reader gets version
checking for free. Each wrapper satisfies the same interface it wraps, so
they stack: `VersionedWriter(bufio.Writer(file))`.

These two patterns work together: `io.Reader`/`io.Writer` provide the
uniform seam, and wrappers layer cross-cutting concerns onto that seam
without modifying any existing code. One `VersionedWriter` gives schema
versioning to all five generators at once. One `ValidatingReader` gives
version checking to every file reader.

**Practical consequence:** `PaletteParserBytes` doesn't exist — a parser
that takes `io.Reader` already handles files, embedded bytes, and test
buffers. The composition root opens the file or wraps embedded bytes in a
`bytes.Reader` and passes the reader in. The parser never knows the
difference.

---

## Implementation Checklist

Track progress by marking items `[x]` as completed.

### Phase 1: Foundation

- [x] 1.1 — Project scaffolding (go module, directory structure, Makefile, `.go-arch-lint.yml`)
- [x] 1.2 — Domain: Color value object and parsing (hex ↔ RGB ↔ HSL)
- [x] 1.3 — Domain: Color operations (Blend, BlendBg, Lighten, Darken, Desaturate, ShiftHue)
- [x] 1.4 — Domain: Palette entity (base24 struct, slot access, base16 fallbacks)
- [x] 1.5 — Domain: Token value object and TokenSet aggregate
- [x] 1.6 — Domain: ResolvedTheme aggregate
- [x] 1.7 — Domain: Palette validation rules (luminance ordering, completeness)
- [x] 1.8 — Domain: Error types (ParseError, ValidationError, GenerateError, SchemaVersionError)
- [x] 1.9 — Domain: Schema version constants and file type registry
- [x] 1.10 — Port interfaces (PaletteParser, PaletteSource, TokenDeriver, Mapper, Generator, ThemeStore)
- [x] 1.11 — Port file structs (PaletteFile, UniversalFile, VimMappingFile, etc.)
- [x] 1.12 — Port theme structs (VimTheme, GtkTheme, QssTheme, CssTheme, StylixTheme)
- [x] 1.13 — Adapter: YAML palette parser (io.Reader → domain.Palette, common tinted-theming format only)
- [x] 1.14 — Adapter: ThemeStore (filesystem — read/write theme dirs, symlink management)
- [x] 1.15 — Adapter: Built-in palettes (//go:embed, PaletteSource impl, Get returns io.Reader)
  - [x] 1.15a — tokyo-night-dark.yaml
  - [x] 1.15b — gruvbox-dark.yaml
  - [x] 1.15c — catppuccin-mocha.yaml
- [x] 1.16 — Adapter: VersionedWriter (wraps io.Writer, prepends schema_version + kind header)
- [x] 1.17 — Adapter: ValidatingReader (wraps io.Reader, peeks schema version, returns SchemaVersionError)
- [x] 1.18 — Testdata: Tokyo Night Dark reference palette YAML

### Phase 2: Layer B — Token Derivation

- [ ] 2.1 — Adapter: Default token deriver (base24 → ~87 semantic tokens)
  - [ ] 2.1a — Surface tokens (11 tokens)
  - [ ] 2.1b — Text tokens (7 tokens)
  - [ ] 2.1c — Status tokens (6 tokens)
  - [ ] 2.1d — Diff tokens (9 tokens)
  - [ ] 2.1e — Syntax tokens (14 tokens)
  - [ ] 2.1f — Markup tokens (10 tokens)
  - [ ] 2.1g — Accent, border, scrollbar, state tokens (11 tokens)
  - [ ] 2.1h — Git tokens (4 tokens)
  - [ ] 2.1i — Terminal ANSI colors (16 tokens)
- [ ] 2.2 — Adapter: UniversalFile writer (TokenSet → io.Writer as YAML)
- [ ] 2.3 — Adapter: UniversalFile reader (io.Reader → TokenSet)
- [ ] 2.4 — Application: DeriveTheme use case (io.Reader → universal.yaml via io.Writer)
- [ ] 2.5 — Unit tests for derivation rules against Tokyo Night Dark palette

### Phase 3: Layer C+D — Mapping + Generation (per target)

- [ ] 3.1 — Adapter: Stylix mapper + generator (simplest, validates full pipeline)
  - [ ] 3.1a — Mapper: ResolvedTheme → ports.StylixTheme (≥60 keys)
  - [ ] 3.1b — MappingFile writer: StylixTheme → stylix-mapping.yaml
  - [ ] 3.1c — MappingFile reader: stylix-mapping.yaml → StylixTheme
  - [ ] 3.1d — Generator: StylixTheme → style.json (sorted, 2-space indent)
- [ ] 3.2 — Adapter: CSS mapper + generator
  - [ ] 3.2a — Mapper: custom properties + element rules
  - [ ] 3.2b — MappingFile writer/reader
  - [ ] 3.2c — Generator: :root{} + element selectors → style.css
- [ ] 3.3 — Adapter: Vim mapper + generator (most complex)
  - [ ] 3.3a — Mapper: base highlights (Normal, Comment, Visual, CursorLine, etc.)
  - [ ] 3.3b — Mapper: treesitter highlights (@keyword, @string, @function, etc.)
  - [ ] 3.3c — Mapper: LSP semantic token links
  - [ ] 3.3d — Mapper: diagnostic highlights (virtual text, underlines)
  - [ ] 3.3e — Mapper: plugin highlights (telescope, gitsigns, etc.)
  - [ ] 3.3f — Mapper: markup highlights
  - [ ] 3.3g — Mapper: terminal ANSI colors (16)
  - [ ] 3.3h — MappingFile writer/reader
  - [ ] 3.3i — Generator: .lua output (hi clear, nvim_set_hl, links, terminal)
- [ ] 3.4 — Adapter: GTK mapper + generator
  - [ ] 3.4a — Mapper: @define-color definitions + widget selector rules
  - [ ] 3.4b — MappingFile writer/reader
  - [ ] 3.4c — Generator: CSS output (@define-color then selectors) → gtk.css
- [ ] 3.5 — Adapter: QSS mapper + generator
  - [ ] 3.5a — Mapper: widget + pseudo-state rules
  - [ ] 3.5b — MappingFile writer/reader
  - [ ] 3.5c — Generator: literal hex, no variables → style.qss
- [ ] 3.6 — Application: GenerateTheme use case (full pipeline or partial regeneration)

### Phase 4: CLI (Composition Root)

- [ ] 4.1 — CLI scaffolding and dependency wiring in cmd
- [ ] 4.2 — `generate` command (palette file → full theme directory)
  - [ ] 4.2a — `--dir` flag (custom config dir, default `~/.config/flair`)
  - [ ] 4.2b — `--target` flag (single target output)
- [ ] 4.3 — `select` command (switch active theme via symlinks)
- [ ] 4.4 — `validate` command (lint all files in a theme directory)
- [ ] 4.5 — `preview` command (ANSI color preview in terminal)
- [ ] 4.6 — `init` command (scaffold a new theme directory with palette template)
- [ ] 4.7 — `list` command (show available themes, mark selected)
- [ ] 4.8 — `regenerate` command (re-derive downstream from edited intermediate files)

### Phase 5: Polish

- [ ] 5.1 — End-to-end test: Tokyo Night Dark base24 → full theme directory, all outputs (golden files)
- [ ] 5.2 — End-to-end test: Generate from built-in name produces identical output to file
- [ ] 5.3 — End-to-end test: 2 additional base24 schemes from tinted-theming
- [ ] 5.4 — End-to-end test: Edit universal.yaml, regenerate → only downstream files change
- [ ] 5.5 — End-to-end test: All built-in palettes parse and validate cleanly
- [ ] 5.6 — `go-arch-lint check` passes clean
- [ ] 5.7 — README and usage documentation
- [ ] 5.8 — Makefile targets (build, test, lint, arch-lint, install)

---

## Hexagonal Architecture Mapping

### How the theme layers map to hex layers

```
Theme Concept          Hex Layer         Why
─────────────────────  ────────────────  ──────────────────────────────────
Color, Palette,        Domain            Pure types + value object logic.
Token, ResolvedTheme,                    No deps. Color ops are domain math.
color ops, validation,                   Schema version constants live here.
error types, schema
versions

PaletteParser,         Ports             Interfaces defining boundaries.
TokenDeriver,                            Depend only on domain types.
Mapper, Generator,                       Also: file structs (PaletteFile,
ThemeStore                               UniversalFile, *MappingFile) and
                                         mapped theme structs (DTOs) shared
                                         between mapper and generator
                                         adapters without cross-adapter deps.

DeriveTheme,           Application       Use cases orchestrating domain
GenerateTheme,                           logic through port interfaces.
RegenerateTheme,                         Depends on ports, NOT adapters.
ValidateTheme,
PreviewTheme,
InitTheme,
SelectTheme,
ListThemes

YAML parser,           Adapters          Concrete implementations of ports.
Default deriver,                         Depend on ports + domain.
ThemeStore (fs),                         Adapters do NOT depend on each
file readers/writers,                    other.
Built-in palettes,                       Reader/Writer wrappers (Versioned
VersionedWriter,                         Writer, ValidatingReader) layer
ValidatingReader,                        cross-cutting concerns via
Vim/GTK/QSS/CSS/                         embedding/composition.
Stylix mappers,
Vim/GTK/QSS/CSS/
Stylix generators

CLI, DI wiring         Cmd               Composition root. Wires adapters
                                         to application via ports.
                                         May depend on everything.
```

### Dependency flow

```
    ┌──────────────────────────────────────────┐
    │                  cmd                      │
    │  (composition root — wires everything)    │
    └────────┬──────────────┬──────────────────┘
             │              │
     ┌───────▼──────┐  ┌───▼──────────────────┐
     │ application  │  │      adapters         │
     │ (use cases)  │  │  ┌─────┐ ┌─────────┐ │
     │              │  │  │yaml │ │ deriver  │ │
     │              │  │  └─────┘ └─────────┘ │
     └───────┬──────┘  │  ┌─────┐ ┌─────────┐ │
             │         │  │store│ │palettes │ │
             │         │  └─────┘ └─────────┘ │
             │         │  ┌──────────────────┐ │
             │         │  │     fileio       │ │
             │         │  └──────────────────┘ │
             │         │  ┌──────────────────┐ │
             │         │  │    wrappers      │ │
             │         │  │VersionedWriter   │ │
             │         │  │ValidatingReader  │ │
             │         │  └──────────────────┘ │
             │         │  ┌─────┐ ┌─────────┐ │
             │         │  │ map │ │   gen    │ │
             │         │  │ vim │ │   vim    │ │
             │         │  │ gtk │ │   gtk    │ │
             │         │  │ qss │ │   qss    │ │
             │         │  │ css │ │   css    │ │
             │         │  │stlx │ │  stlx    │ │
             │         │  └─────┘ └─────────┘ │
             │         └──────────┬────────────┘
             │                    │
     ┌───────▼────────────────────▼─┐
     │            ports             │
     │  (interfaces + file structs  │
     │   + theme DTOs)              │
     └───────────────┬──────────────┘
                     │
     ┌───────────────▼──────────────┐
     │           domain             │
     │  (pure types + color math    │
     │   + schema versions)         │
     └──────────────────────────────┘
```

All arrows point inward. Adapters never import other adapters.

---

## Directory Structure (source)

```
flair/
├── cmd/
│   └── flair/
│       └── main.go                     # Composition root: wires adapters → ports → app
│
├── internal/
│   ├── domain/
│   │   ├── color.go                    # Color value object (RGB, HSL, parsing, formatting)
│   │   ├── color_test.go
│   │   ├── ops.go                      # Blend, BlendBg, Lighten, Darken, Desaturate, ShiftHue
│   │   ├── ops_test.go
│   │   ├── palette.go                  # Palette entity (24 slots, base16 fallbacks, slot access)
│   │   ├── palette_test.go
│   │   ├── token.go                    # Token value object (Color + style flags)
│   │   ├── tokenset.go                 # TokenSet aggregate (map[string]Token, merge, access)
│   │   ├── tokenset_test.go
│   │   ├── theme.go                    # ResolvedTheme aggregate (Palette + TokenSet)
│   │   ├── theme_test.go
│   │   ├── errors.go                   # ParseError, ValidationError, GenerateError, SchemaVersionError
│   │   ├── schema.go                   # Schema version constants per file type
│   │   ├── validation.go              # Palette validation rules (luminance, completeness)
│   │   └── validation_test.go
│   │
│   ├── ports/
│   │   ├── parser.go                   # PaletteParser interface (io.Reader based)
│   │   ├── palettes.go                 # PaletteSource interface
│   │   ├── deriver.go                  # TokenDeriver interface
│   │   ├── mapper.go                   # Mapper interface
│   │   ├── generator.go               # Generator interface + Target struct
│   │   ├── store.go                    # ThemeStore interface (dirs, symlinks, file I/O)
│   │   ├── files.go                    # PaletteFile, UniversalFile, *MappingFile structs
│   │   └── themes.go                   # VimTheme, GtkTheme, QssTheme, CssTheme, StylixTheme
│   │
│   ├── application/
│   │   ├── derive.go                   # DeriveTheme use case (palette → universal)
│   │   ├── derive_test.go
│   │   ├── generate.go                 # GenerateTheme use case (full pipeline)
│   │   ├── generate_test.go
│   │   ├── regenerate.go              # RegenerateTheme use case (partial rebuild)
│   │   ├── regenerate_test.go
│   │   ├── validate.go                 # ValidateTheme use case
│   │   ├── validate_test.go
│   │   ├── select.go                   # SelectTheme use case (symlink management)
│   │   ├── select_test.go
│   │   ├── list.go                     # ListThemes use case
│   │   ├── preview.go                  # PreviewTheme use case
│   │   └── init.go                     # InitTheme use case (scaffold directory)
│   │
│   ├── adapters/
│   │   ├── yaml/
│   │   │   ├── parser.go              # PaletteParser impl (io.Reader → domain.Palette)
│   │   │   └── parser_test.go
│   │   │
│   │   ├── store/
│   │   │   ├── fs.go                  # ThemeStore impl (filesystem, dirs, symlinks)
│   │   │   └── fs_test.go
│   │   │
│   │   ├── fileio/
│   │   │   ├── universal.go           # Read/write universal.yaml (io.Reader/io.Writer)
│   │   │   ├── universal_test.go
│   │   │   ├── mapping.go            # Read/write *-mapping.yaml (io.Reader/io.Writer)
│   │   │   └── mapping_test.go
│   │   │
│   │   ├── wrappers/
│   │   │   ├── versioned.go          # VersionedWriter: wraps io.Writer, prepends schema header
│   │   │   ├── validating.go         # ValidatingReader: wraps io.Reader, checks schema version
│   │   │   └── wrappers_test.go
│   │   │
│   │   ├── palettes/
│   │   │   ├── palettes.go           # //go:embed *.yaml + List(), Get() functions
│   │   │   ├── palettes_test.go
│   │   │   ├── tokyo-night-dark.yaml # Built-in base24 palette
│   │   │   ├── gruvbox-dark.yaml
│   │   │   └── catppuccin-mocha.yaml
│   │   │
│   │   ├── deriver/
│   │   │   ├── deriver.go             # TokenDeriver impl (default derivation rules)
│   │   │   └── deriver_test.go
│   │   │
│   │   ├── mapper/
│   │   │   ├── vim.go                 # Vim Mapper: ResolvedTheme → ports.VimTheme
│   │   │   ├── vim_test.go
│   │   │   ├── vim_plugins.go         # Vim plugin highlight sub-mappings
│   │   │   ├── gtk.go
│   │   │   ├── gtk_test.go
│   │   │   ├── qss.go
│   │   │   ├── qss_test.go
│   │   │   ├── css.go
│   │   │   ├── css_test.go
│   │   │   ├── stylix.go
│   │   │   └── stylix_test.go
│   │   │
│   │   └── generator/
│   │       ├── vim.go                 # Vim Generator: ports.VimTheme → style.lua
│   │       ├── vim_test.go
│   │       ├── gtk.go                 # GTK Generator: ports.GtkTheme → gtk.css
│   │       ├── gtk_test.go
│   │       ├── qss.go                # QSS Generator: ports.QssTheme → style.qss
│   │       ├── qss_test.go
│   │       ├── css.go                 # CSS Generator: ports.CssTheme → style.css
│   │       ├── css_test.go
│   │       ├── stylix.go             # Stylix Generator: ports.StylixTheme → style.json
│   │       └── stylix_test.go
│   │
│   └── config/
│       └── config.go                  # App configuration (config dir path, defaults)
│
├── testdata/
│   ├── tokyo-night-dark.yaml          # Reference base24 palette
│   ├── tokyo-night-dark-base16.yaml   # base16-only variant for fallback testing
│   ├── invalid-palette.yaml           # Malformed input for error testing
│   └── expected/                      # Golden files for each generated output
│       ├── universal.yaml
│       ├── vim-mapping.yaml
│       ├── css-mapping.yaml
│       ├── gtk-mapping.yaml
│       ├── qss-mapping.yaml
│       ├── stylix-mapping.yaml
│       ├── style.lua
│       ├── style.css
│       ├── gtk.css
│       ├── style.qss
│       └── style.json
│
├── .go-arch-lint.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Architecture Linter Configuration

```yaml
# .go-arch-lint.yml
version: 3

allow:
  depOnAnyVendor: true
  deepScan: false

components:
  domain:
    in: internal/domain/**

  ports:
    in: internal/ports/**

  application:
    in: internal/application/**

  adapters:
    in: internal/adapters/**

  config:
    in: internal/config/**

  cmd:
    in: cmd/**

commonComponents:
  - domain
  - ports

deps:
  domain:
    anyVendorDeps: true

  ports:
    anyVendorDeps: true

  application:
    mayDependOn:
      - ports

  adapters:
    mayDependOn:
      - ports

  config:
    anyVendorDeps: true

  cmd:
    anyProjectDeps: true
    anyVendorDeps: true

excludeFiles:
  - _test\.go$
  - mock_.*\.go$
  - /mocks/
  - /testdata/
```

---

## Schema Version Constants

```go
// internal/domain/schema.go
package domain

// Schema versions for each file type. Bump when the file format changes.
const (
    SchemaPalette        = 1
    SchemaUniversal      = 1
    SchemaVimMapping     = 1
    SchemaCssMapping     = 1
    SchemaGtkMapping     = 1
    SchemaQssMapping     = 1
    SchemaStylixMapping  = 1
)

// FileKind identifies a file type in the pipeline.
type FileKind string

const (
    FileKindPalette       FileKind = "palette"
    FileKindUniversal     FileKind = "universal"
    FileKindVimMapping    FileKind = "vim-mapping"
    FileKindCssMapping    FileKind = "css-mapping"
    FileKindGtkMapping    FileKind = "gtk-mapping"
    FileKindQssMapping    FileKind = "qss-mapping"
    FileKindStylixMapping FileKind = "stylix-mapping"
)

// CurrentVersion returns the current schema version for a file kind.
func CurrentVersion(kind FileKind) int { /* switch on kind */ }
```

---

## File Structs (Port Layer)

Every intermediate YAML file is a serialized Go struct. These live in
`ports/files.go` so both the file reader/writer adapter and the use cases
can reference them without cross-adapter imports.

```go
// internal/ports/files.go
package ports

import "github.com/curtbushko/flair/internal/domain"

// FileHeader is embedded in every YAML file produced by flair.
type FileHeader struct {
    SchemaVersion int            `yaml:"schema_version"`
    Kind          domain.FileKind `yaml:"kind"`
    ThemeName     string         `yaml:"theme_name"`
}

// PaletteFile is the input palette (base24).
type PaletteFile struct {
    FileHeader `yaml:",inline"`
    System     string            `yaml:"system"`
    Author     string            `yaml:"author"`
    Variant    string            `yaml:"variant"`
    Palette    map[string]string `yaml:"palette"` // "base00": "1a1b26"
}

// UniversalToken is a single semantic token in universal.yaml.
type UniversalToken struct {
    Color         string `yaml:"color"`                    // "#7aa2f7"
    Bold          bool   `yaml:"bold,omitempty"`
    Italic        bool   `yaml:"italic,omitempty"`
    Underline     bool   `yaml:"underline,omitempty"`
    Undercurl     bool   `yaml:"undercurl,omitempty"`
    Strikethrough bool   `yaml:"strikethrough,omitempty"`
}

// UniversalFile is the derived semantic token set.
type UniversalFile struct {
    FileHeader `yaml:",inline"`
    Tokens     map[string]UniversalToken `yaml:"tokens"`
}

// VimMappingHighlight is a single Vim highlight group in the mapping file.
type VimMappingHighlight struct {
    Fg            string `yaml:"fg,omitempty"`
    Bg            string `yaml:"bg,omitempty"`
    Sp            string `yaml:"sp,omitempty"`
    Bold          bool   `yaml:"bold,omitempty"`
    Italic        bool   `yaml:"italic,omitempty"`
    Underline     bool   `yaml:"underline,omitempty"`
    Undercurl     bool   `yaml:"undercurl,omitempty"`
    Strikethrough bool   `yaml:"strikethrough,omitempty"`
    Reverse       bool   `yaml:"reverse,omitempty"`
    Nocombine     bool   `yaml:"nocombine,omitempty"`
    Link          string `yaml:"link,omitempty"`
}

// VimMappingFile is the Vim-specific mapping.
type VimMappingFile struct {
    FileHeader     `yaml:",inline"`
    Highlights     map[string]VimMappingHighlight `yaml:"highlights"`
    TerminalColors [16]string                     `yaml:"terminal_colors"`
}

// CssRuleEntry is a CSS selector with its properties.
type CssRuleEntry struct {
    Selector   string            `yaml:"selector"`
    Properties map[string]string `yaml:"properties"`
}

// CssMappingFile is the CSS-specific mapping.
type CssMappingFile struct {
    FileHeader       `yaml:",inline"`
    CustomProperties map[string]string `yaml:"custom_properties"`
    Rules            []CssRuleEntry    `yaml:"rules"`
}

// GtkMappingFile is the GTK-specific mapping.
type GtkMappingFile struct {
    FileHeader `yaml:",inline"`
    Colors     map[string]string `yaml:"colors"` // "window_bg_color": "#1a1b26"
    Rules      []CssRuleEntry    `yaml:"rules"`
}

// QssMappingFile is the QSS-specific mapping.
type QssMappingFile struct {
    FileHeader `yaml:",inline"`
    Rules      []CssRuleEntry `yaml:"rules"`
}

// StylixMappingFile is the Stylix-specific mapping.
type StylixMappingFile struct {
    FileHeader `yaml:",inline"`
    Values     map[string]string `yaml:"values"` // "bg": "#1a1b26"
}
```

---

## Cross-Adapter Type Sharing

The key arch constraint: **adapters cannot depend on each other.**

A Vim mapper produces a `VimTheme`. A Vim generator consumes a `VimTheme`.
If `VimTheme` lived in the mapper package, the generator would import the
mapper — violating the rule.

**Solution: Mapped theme structs live in `ports/themes.go`.**

Both mapper and generator adapters import from `ports`. They share the
theme types without depending on each other. The application layer passes
`MappedTheme` (typed as `any`) between them; the generator type-asserts
to the expected struct.

Similarly, the `ports/files.go` file structs are shared between the
`fileio` adapter (which reads/writes them) and the application use cases
(which orchestrate the pipeline).

```go
// internal/ports/themes.go
package ports

import "github.com/curtbushko/flair/internal/domain"

type VimHighlight struct {
    Fg, Bg, Sp    *domain.Color
    Bold          bool
    Italic        bool
    Underline     bool
    Undercurl     bool
    Strikethrough bool
    Reverse       bool
    Nocombine     bool
    Link          string
}

type VimTheme struct {
    Name           string
    Highlights     map[string]VimHighlight
    TerminalColors [16]domain.Color
}

type StylixTheme struct {
    Values map[string]string
}

type CssProperty struct {
    Property string
    Value    string
}

type CssRule struct {
    Selector   string
    Properties []CssProperty
}

type CssTheme struct {
    CustomProperties map[string]string
    Rules            []CssRule
}

type GtkColorDef struct {
    Name  string
    Value string
}

type GtkTheme struct {
    Colors []GtkColorDef
    Rules  []CssRule
}

type QssTheme struct {
    Rules []CssRule
}
```

---

## Port Interfaces

```go
// internal/ports/parser.go
package ports

import (
    "io"
    "github.com/curtbushko/flair/internal/domain"
)

// PaletteParser reads palette YAML from a reader and returns a domain Palette.
// The caller is responsible for opening/closing the underlying source.
// Works identically on files, embedded built-ins, test buffers, or stdin.
type PaletteParser interface {
    Parse(r io.Reader) (domain.Palette, error)
}
```

```go
// internal/ports/palettes.go
package ports

import "io"

// PaletteSource provides access to built-in palettes shipped with flair.
type PaletteSource interface {
    // List returns the names of all built-in palettes (e.g. "tokyo-night-dark").
    List() []string

    // Get returns a reader for the named built-in palette's YAML.
    // Returns an error if the name is not found.
    Get(name string) (io.Reader, error)

    // Has returns true if the named palette exists as a built-in.
    Has(name string) bool
}
```

```go
// internal/ports/deriver.go
package ports

import "github.com/curtbushko/flair/internal/domain"

// TokenDeriver derives the full semantic token set from a base24 palette.
type TokenDeriver interface {
    Derive(p domain.Palette) domain.TokenSet
}
```

```go
// internal/ports/mapper.go
package ports

import "github.com/curtbushko/flair/internal/domain"

type MappedTheme any

// Mapper transforms a ResolvedTheme into a target-specific theme struct.
type Mapper interface {
    Name() string
    Map(theme domain.ResolvedTheme) (MappedTheme, error)
}
```

```go
// internal/ports/generator.go
package ports

import "io"

// Generator writes the final output file from a mapped theme.
type Generator interface {
    Name() string
    DefaultFilename() string  // e.g. "style.lua", "gtk.css"
    Generate(w io.Writer, mapped MappedTheme) error
}

// Target pairs a mapper with its generator and file I/O.
type Target struct {
    Mapper       Mapper
    Generator    Generator
    MappingFile  string // filename in theme dir, e.g. "vim-mapping.yaml"
}
```

```go
// internal/ports/store.go
package ports

import (
    "io"
    "time"
)

// ThemeStore manages theme directories and symlinks on the filesystem.
type ThemeStore interface {
    // ConfigDir returns the root config directory (e.g. ~/.config/flair).
    ConfigDir() string

    // ThemeDir returns the path for a named theme.
    ThemeDir(themeName string) string

    // EnsureThemeDir creates the theme directory if it doesn't exist.
    EnsureThemeDir(themeName string) error

    // ListThemes returns all theme directory names.
    ListThemes() ([]string, error)

    // SelectedTheme returns the currently symlinked theme name, or "" if none.
    SelectedTheme() (string, error)

    // Select creates/updates symlinks at the config root pointing to the
    // given theme's output files.
    Select(themeName string) error

    // OpenReader opens a file within a theme directory for reading.
    // The caller must close the returned reader.
    OpenReader(themeName, filename string) (io.ReadCloser, error)

    // OpenWriter opens (or creates) a file within a theme directory for writing.
    // The caller must close the returned writer.
    OpenWriter(themeName, filename string) (io.WriteCloser, error)

    // FileExists checks whether a file exists in a theme directory.
    FileExists(themeName, filename string) bool

    // FileMtime returns the modification time of a file.
    FileMtime(themeName, filename string) (time.Time, error)
}
```

---

## Built-in Palettes Adapter

```go
// internal/adapters/palettes/palettes.go
package palettes

import (
    "bytes"
    "embed"
    "io"
)

//go:embed *.yaml
var fs embed.FS

// Source implements ports.PaletteSource using embedded YAML files.
type Source struct{}

func NewSource() *Source { return &Source{} }

// List returns the names of all built-in palettes (filename without .yaml).
func (s *Source) List() []string { /* read dir from fs, strip .yaml */ }

// Get returns a reader for the named built-in palette's YAML.
// Wraps the embedded bytes in a bytes.Reader — no file I/O.
func (s *Source) Get(name string) (io.Reader, error) {
    data, err := fs.ReadFile(name + ".yaml")
    if err != nil { return nil, err }
    return bytes.NewReader(data), nil
}

// Has returns true if the named palette exists as a built-in.
func (s *Source) Has(name string) bool { /* attempt fs.ReadFile, return err == nil */ }
```

The composition root in `cmd/flair/main.go` uses `PaletteSource.Has()` to
determine whether a `generate` argument is a built-in name or a file path.
If built-in, `Get()` returns an `io.Reader` that is passed directly to
`PaletteParser.Parse()`. If a file path, the composition root opens the
file (getting an `io.ReadCloser`) and passes it to the same `Parse()`.
The parser never knows the difference — it just reads from a reader.

Built-in palettes ship as `.yaml` files in the same directory as the Go
source. Adding a new built-in palette is just adding a YAML file — no code
changes required.

---

## Reader/Writer Wrappers (Decorator Pattern)

Wrappers live in `adapters/wrappers/` and layer cross-cutting concerns onto
`io.Reader` and `io.Writer` via embedding. Each wrapper satisfies the same
interface it wraps, so they compose: `VersionedWriter(bufio.Writer(file))`.

```go
// internal/adapters/wrappers/versioned.go
package wrappers

import (
    "io"
    "github.com/curtbushko/flair/internal/domain"
)

// VersionedWriter wraps an io.Writer and prepends a YAML schema_version
// + kind header before the first write. Generators write their content
// without knowing about versioning; the wrapper adds it.
type VersionedWriter struct {
    inner       io.Writer
    kind        domain.FileKind
    themeName   string
    headerDone  bool
}

func NewVersionedWriter(w io.Writer, kind domain.FileKind, themeName string) *VersionedWriter {
    return &VersionedWriter{inner: w, kind: kind, themeName: themeName}
}

// Write prepends the schema header on first call, then delegates to inner.
func (vw *VersionedWriter) Write(p []byte) (int, error) { /* ... */ }
```

```go
// internal/adapters/wrappers/validating.go
package wrappers

import (
    "io"
    "github.com/curtbushko/flair/internal/domain"
)

// ValidatingReader wraps an io.Reader, peeks at the first few bytes to
// extract the schema_version field, and returns a SchemaVersionError if
// the version is incompatible. If valid, the full content (including the
// peeked bytes) is available for reading.
type ValidatingReader struct {
    inner    io.Reader
    kind     domain.FileKind
    validated bool
}

func NewValidatingReader(r io.Reader, kind domain.FileKind) *ValidatingReader {
    return &ValidatingReader{inner: r, kind: kind}
}

// Read validates the schema version on first call (using io.MultiReader
// to replay peeked bytes), then delegates to inner for all subsequent reads.
func (vr *ValidatingReader) Read(p []byte) (int, error) { /* ... */ }
```

### How wrappers compose in the pipeline

The composition root (cmd) wires wrappers when opening readers/writers:

```go
// Writing universal.yaml:
w, _ := store.OpenWriter(themeName, "universal.yaml")
defer w.Close()
vw := wrappers.NewVersionedWriter(w, domain.FileKindUniversal, themeName)
fileio.WriteUniversal(vw, tokenSet)  // writer doesn't know about versioning

// Reading universal.yaml:
r, _ := store.OpenReader(themeName, "universal.yaml")
defer r.Close()
vr := wrappers.NewValidatingReader(r, domain.FileKindUniversal)
tokenSet, err := fileio.ReadUniversal(vr)  // reader doesn't know about validation
// err may be *domain.SchemaVersionError → triggers regeneration

// Generating style.lua:
w, _ := store.OpenWriter(themeName, "style.lua")
defer w.Close()
bw := bufio.NewWriter(w)         // buffer for performance
generator.Generate(bw, vimTheme) // generator just writes to io.Writer
bw.Flush()
// (no VersionedWriter here — final output files don't have schema headers)
```

This means:
- **Generators** never know about schema versions — they just write to `io.Writer`
- **File readers** never know about version validation — they just read from `io.Reader`
- **The fileio adapter** reads/writes structs via `io.Reader`/`io.Writer` — no file paths
- **Schema version logic** lives in exactly one place (the wrappers), applied at composition time

---

## Domain Types

```go
// internal/domain/color.go
package domain

type Color struct {
    R, G, B uint8
    IsNone  bool
}

type HSL struct {
    H float64 // 0-360
    S float64 // 0-1
    L float64 // 0-1
}

func ParseHex(s string) (Color, error)  { /* ... */ }
func (c Color) Hex() string             { /* ... */ }
func (c Color) ToHSL() HSL              { /* ... */ }
func (h HSL) ToRGB() Color              { /* ... */ }
func (c Color) Luminance() float64      { /* WCAG relative luminance */ }
func (c Color) Equal(other Color) bool  { /* ... */ }
func NoneColor() Color                  { return Color{IsNone: true} }
```

```go
// internal/domain/ops.go
package domain

// Blend performs linear RGB interpolation: result = a + t*(b - a)
// t=0.0 returns a, t=1.0 returns b.
func Blend(a, b Color, t float64) Color { /* ... */ }

// BlendBg blends fg into bg at the given ratio.
// Equivalent to Blend(bg, fg, amount) — the bg is the base, fg is mixed in.
// This matches TokyoNight semantics where BlendBg(fg, bg, 0.25) means
// "25% of fg mixed into bg".
func BlendBg(fg, bg Color, amount float64) Color { /* ... */ }

func Lighten(c Color, amount float64) Color    { /* HSL: L = min(L+amount, 1.0) */ }
func Darken(c Color, amount float64) Color     { /* HSL: L = max(L-amount, 0.0) */ }
func Desaturate(c Color, amount float64) Color { /* HSL: S = S * (1 - amount) */ }
func ShiftHue(c Color, degrees float64) Color  { /* HSL: H = (H + degrees) mod 360 */ }
```

```go
// internal/domain/token.go
package domain

type Token struct {
    Color         Color
    Bold          bool
    Italic        bool
    Underline     bool
    Undercurl     bool
    Strikethrough bool
}

func (t Token) HasStyle() bool { /* ... */ }
```

```go
// internal/domain/tokenset.go
package domain

type TokenSet struct {
    tokens map[string]Token
}

func NewTokenSet() TokenSet                              { /* ... */ }
func (ts *TokenSet) Set(path string, t Token)            { /* ... */ }
func (ts TokenSet) Get(path string) (Token, bool)        { /* ... */ }
func (ts TokenSet) MustGet(path string) Token            { /* panics if missing */ }
func (ts TokenSet) Paths() []string                      { /* sorted */ }
func (ts TokenSet) Len() int                             { /* ... */ }
```

```go
// internal/domain/theme.go
package domain

type ResolvedTheme struct {
    Name    string
    Variant string  // "dark" or "light"
    Palette Palette
    Tokens  TokenSet
}

func (rt ResolvedTheme) Token(path string) (Token, bool) { /* ... */ }
```

```go
// internal/domain/palette.go
package domain

type Palette struct {
    Name    string
    Author  string
    Variant string  // "dark" or "light"
    System  string  // "base24" or "base16"
    Slug    string  // machine-readable name (e.g. "tokyonight")
    Colors  [24]Color
}

// NewPalette constructs a Palette. If fewer than 24 colors are provided,
// base16 fallback rules are applied for base10–base17.
func NewPalette(name, author, variant, system string, colors map[string]Color) (Palette, error)

func (p Palette) Base(n int) Color
func (p Palette) Slot(name string) (Color, error)
func (p Palette) SlotNames() [24]string
```

---

## Error Types

```go
// internal/domain/errors.go
package domain

import "fmt"

// ParseError indicates a failure to parse input data.
type ParseError struct {
    Field   string
    Message string
    Cause   error
}

func (e *ParseError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("parse error in %s: %s: %v", e.Field, e.Message, e.Cause)
    }
    return fmt.Sprintf("parse error in %s: %s", e.Field, e.Message)
}
func (e *ParseError) Unwrap() error { return e.Cause }

// ValidationError indicates a palette or file fails validation.
type ValidationError struct {
    Violations []string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed: %d violations", len(e.Violations))
}

// GenerateError indicates a target generation failure.
type GenerateError struct {
    Target  string
    Message string
    Cause   error
}

func (e *GenerateError) Error() string {
    return fmt.Sprintf("generate %s: %s", e.Target, e.Message)
}
func (e *GenerateError) Unwrap() error { return e.Cause }

// SchemaVersionError indicates a file has an incompatible schema version.
type SchemaVersionError struct {
    File           string
    Found          int
    Expected       int
    NeedsUpgrade   bool  // true if Found > Expected (user needs newer flair)
}

func (e *SchemaVersionError) Error() string {
    if e.NeedsUpgrade {
        return fmt.Sprintf("%s: schema version %d is newer than supported %d — please upgrade flair",
            e.File, e.Found, e.Expected)
    }
    return fmt.Sprintf("%s: schema version %d is outdated (current: %d) — will regenerate",
        e.File, e.Found, e.Expected)
}
```

---

## Configuration

```go
// internal/config/config.go
package config

import (
    "os"
    "path/filepath"
)

type Config struct {
    ConfigDir string   // Root dir, default: ~/.config/flair
}

// DefaultConfigDir returns ~/.config/flair, respecting XDG_CONFIG_HOME.
func DefaultConfigDir() string {
    if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
        return filepath.Join(xdg, "flair")
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "flair")
}
```

---

## CLI Design

Go stdlib only — no external CLI framework. Subcommand dispatch via `switch`
on `os.Args[1]`, each subcommand uses its own `flag.FlagSet`.

```
flair generate <palette> [--name <theme-name>] [--dir <config-dir>] [--target <n>]
flair regenerate <theme-name> [--dir <config-dir>] [--target <n>]
flair select <theme-name> [--dir <config-dir>]
flair validate <theme-name> [--dir <config-dir>]
flair preview <theme-name> [--dir <config-dir>]
flair list [--dir <config-dir>]
flair list --builtins
flair init [--name <theme-name>] [--dir <config-dir>]
```

`<palette>` can be either a **built-in name** (e.g. `tokyo-night-dark`) or a
**file path** (e.g. `./my-palette.yaml`). Flair checks built-in names first.

| Command      | Description                                                        |
|--------------|--------------------------------------------------------------------|
| `generate`   | Import a palette (built-in or file), create theme dir, run pipeline|
| `regenerate` | Re-derive downstream files from the furthest-upstream edit         |
| `select`     | Set the active theme (update symlinks at config root)              |
| `validate`   | Check all files in a theme dir for correctness and version skew    |
| `preview`    | Print ANSI-colored palette + token preview to terminal             |
| `list`       | Show available themes; `--builtins` shows built-in palette names   |
| `init`       | Scaffold a new theme directory with a template palette.yaml        |

### `generate` flow

1. Resolve `<palette>`: check `PaletteSource.Has(arg)` — if true, `Get()` returns `io.Reader`; otherwise `os.Open(arg)` returns `io.ReadCloser`
2. Pass reader to `PaletteParser.Parse(r)` → `domain.Palette`
3. Infer `theme-name` from palette name (or `--name` flag) → `store.EnsureThemeDir()`
4. Write `palette.yaml` via `store.OpenWriter()` → `VersionedWriter` → fileio
5. Derive → write `universal.yaml` via `VersionedWriter`
6. For each target: map → write `*-mapping.yaml` via `VersionedWriter` → generate → write final output via `store.OpenWriter()`
7. Print summary of files written

### `regenerate` flow

1. Inspect theme dir for file mtimes
2. Read upstream files via `ValidatingReader` — detect stale schema versions
3. Determine the furthest-upstream file that was manually edited (newer than its downstream)
4. Re-derive from that point forward:
   - `palette.yaml` edited → re-derive everything
   - `universal.yaml` edited → re-map + re-generate all targets
   - `vim-mapping.yaml` edited → re-generate only `style.lua`
5. Write regenerated files via `VersionedWriter`
6. Leave upstream files untouched

### `select` flow

1. Verify theme dir exists and has output files
2. Remove existing symlinks at config root
3. Create new symlinks: `~/.config/flair/style.lua → <theme>/style.lua`, etc.
4. Print confirmation

---

## Application Use Case Signatures

```go
// internal/application/derive.go
package application

import (
    "io"
    "github.com/curtbushko/flair/internal/domain"
    "github.com/curtbushko/flair/internal/ports"
)

type DeriveThemeUseCase struct {
    parser  ports.PaletteParser
    deriver ports.TokenDeriver
}

func NewDeriveThemeUseCase(p ports.PaletteParser, d ports.TokenDeriver) *DeriveThemeUseCase

// Execute parses a palette from a reader and derives the full semantic token set.
// The caller provides the reader (file, built-in bytes.Reader, etc.).
func (uc *DeriveThemeUseCase) Execute(r io.Reader) (domain.ResolvedTheme, error)
```

```go
// internal/application/generate.go
package application

import "github.com/curtbushko/flair/internal/ports"

type GenerateThemeUseCase struct {
    parser   ports.PaletteParser
    deriver  ports.TokenDeriver
    targets  []ports.Target
    store    ports.ThemeStore
    builtins ports.PaletteSource
}

func NewGenerateThemeUseCase(
    p ports.PaletteParser,
    d ports.TokenDeriver,
    targets []ports.Target,
    store ports.ThemeStore,
    builtins ports.PaletteSource,
) *GenerateThemeUseCase

// Execute runs the full pipeline: palette → universal → mappings → outputs.
// paletteRef is either a built-in name or a file path. If built-in,
// builtins.Get() provides the io.Reader; if file, the use case opens it.
// If targetFilter is non-empty, only that target is generated.
func (uc *GenerateThemeUseCase) Execute(paletteRef, themeName, targetFilter string) error
```

```go
// internal/application/regenerate.go
package application

import "github.com/curtbushko/flair/internal/ports"

type RegenerateThemeUseCase struct {
    deriver ports.TokenDeriver
    targets []ports.Target
    store   ports.ThemeStore
}

func NewRegenerateThemeUseCase(
    d ports.TokenDeriver,
    targets []ports.Target,
    store ports.ThemeStore,
) *RegenerateThemeUseCase

// Execute inspects a theme directory and regenerates downstream files
// from the furthest-upstream edit.
func (uc *RegenerateThemeUseCase) Execute(themeName, targetFilter string) error
```

```go
// internal/application/validate.go
package application

import "github.com/curtbushko/flair/internal/ports"

type ValidateThemeUseCase struct {
    parser ports.PaletteParser
    store  ports.ThemeStore
}

func NewValidateThemeUseCase(p ports.PaletteParser, s ports.ThemeStore) *ValidateThemeUseCase

// Execute validates all files in a theme directory: schema versions,
// palette correctness, token completeness.
func (uc *ValidateThemeUseCase) Execute(themeName string) ([]string, error)
```

```go
// internal/application/select.go
package application

import "github.com/curtbushko/flair/internal/ports"

type SelectThemeUseCase struct {
    store   ports.ThemeStore
    targets []ports.Target
}

func NewSelectThemeUseCase(s ports.ThemeStore, targets []ports.Target) *SelectThemeUseCase

// Execute creates/updates symlinks at the config root.
func (uc *SelectThemeUseCase) Execute(themeName string) error
```

```go
// internal/application/list.go
package application

import "github.com/curtbushko/flair/internal/ports"

type ListThemesUseCase struct {
    store    ports.ThemeStore
    builtins ports.PaletteSource
}

func NewListThemesUseCase(s ports.ThemeStore, b ports.PaletteSource) *ListThemesUseCase

type ThemeInfo struct {
    Name     string
    Variant  string
    Selected bool
    Complete bool
}

// Execute returns all installed themes.
func (uc *ListThemesUseCase) Execute() ([]ThemeInfo, error)

// ListBuiltins returns the names of all built-in palettes.
func (uc *ListThemesUseCase) ListBuiltins() []string
```

```go
// internal/application/preview.go
package application

import (
    "io"
    "github.com/curtbushko/flair/internal/ports"
)

type PreviewThemeUseCase struct {
    parser  ports.PaletteParser
    deriver ports.TokenDeriver
    store   ports.ThemeStore
}

func NewPreviewThemeUseCase(p ports.PaletteParser, d ports.TokenDeriver, s ports.ThemeStore) *PreviewThemeUseCase

// Execute writes an ANSI-colored preview to w.
func (uc *PreviewThemeUseCase) Execute(themeName string, w io.Writer) error
```

```go
// internal/application/init.go
package application

import "github.com/curtbushko/flair/internal/ports"

type InitThemeUseCase struct {
    store ports.ThemeStore
}

func NewInitThemeUseCase(s ports.ThemeStore) *InitThemeUseCase

// Execute creates a new theme directory with a scaffold palette.yaml.
func (uc *InitThemeUseCase) Execute(themeName string) error
```

---

## Input YAML Format

Flair accepts the **common format** from tinted-theming (spec 0.11+).

```yaml
system: "base24"
name: "Tokyo Night Dark"
slug: "tokyo-night-dark"        # optional, inferred from name if absent
author: "Michaël Ball"
variant: "dark"
palette:
  base00: "1a1b26"
  base01: "1f2335"
  # ... all 24 slots (or 16 with base16 fallbacks)
```

When imported via `flair generate`, the file is rewritten as `palette.yaml`
inside the theme directory with the `schema_version` and `kind` header added.

---

## Reference Palette: Tokyo Night Dark (base24)

```yaml
# testdata/tokyo-night-dark.yaml
system: "base24"
name: "Tokyo Night Dark"
author: "Michaël Ball (base24 by curtbushko)"
variant: "dark"
palette:
  base00: "1a1b26"    # Default Background
  base01: "1f2335"    # Lighter Background (status bars, line numbers)
  base02: "292e42"    # Selection Background
  base03: "565f89"    # Comments, Invisibles, Line Highlighting
  base04: "a9b1d6"    # Dark Foreground (status bars)
  base05: "c0caf5"    # Default Foreground
  base06: "c0caf5"    # Light Foreground
  base07: "c8d3f5"    # Lightest Foreground
  base08: "f7768e"    # Red
  base09: "ff9e64"    # Orange
  base0A: "e0af68"    # Yellow
  base0B: "9ece6a"    # Green
  base0C: "7dcfff"    # Cyan
  base0D: "7aa2f7"    # Blue
  base0E: "bb9af7"    # Magenta
  base0F: "db4b4b"    # Brown/Dark Red
  base10: "16161e"    # Darker Background
  base11: "101014"    # Darkest Background
  base12: "ff899d"    # Bright Red
  base13: "e9c582"    # Bright Yellow
  base14: "afd67a"    # Bright Green
  base15: "97d8f8"    # Bright Cyan
  base16: "8db6fa"    # Bright Blue
  base17: "c8acf8"    # Bright Magenta
```

### Base16 Fallback Rules

When only base00–base0F are provided, base10–base17 are derived:

| Slot   | Fallback | Reason                      |
|--------|----------|-----------------------------|
| base10 | base00   | Darker bg ← same as bg      |
| base11 | base00   | Darkest bg ← same as bg     |
| base12 | base08   | Bright red ← red            |
| base13 | base0A   | Bright yellow ← yellow      |
| base14 | base0B   | Bright green ← green        |
| base15 | base0C   | Bright cyan ← cyan          |
| base16 | base0D   | Bright blue ← blue          |
| base17 | base0E   | Bright magenta ← magenta    |

---

## Init Scaffold Template

`flair init --name my-theme` creates `~/.config/flair/my-theme/palette.yaml`:

```yaml
schema_version: 1
kind: palette
theme_name: "my-theme"
system: "base24"
name: "My Theme"
author: "Your Name"
variant: "dark"
palette:
  # Neutral ramp (backgrounds → foregrounds, darkest → lightest)
  base00: "1a1b26"    # Default Background
  base01: "1f2335"    # Lighter Background (status bars, line numbers)
  base02: "292e42"    # Selection Background
  base03: "565f89"    # Comments, Invisibles
  base04: "a9b1d6"    # Dark Foreground
  base05: "c0caf5"    # Default Foreground
  base06: "c0caf5"    # Light Foreground
  base07: "c8d3f5"    # Lightest Foreground

  # Accent colors
  base08: "f7768e"    # Red
  base09: "ff9e64"    # Orange
  base0A: "e0af68"    # Yellow
  base0B: "9ece6a"    # Green
  base0C: "7dcfff"    # Cyan
  base0D: "7aa2f7"    # Blue
  base0E: "bb9af7"    # Magenta
  base0F: "db4b4b"    # Brown/Dark Red

  # Extended base24 (brighter variants for terminal / status)
  base10: "16161e"    # Darker Background
  base11: "101014"    # Darkest Background
  base12: "ff899d"    # Bright Red
  base13: "e9c582"    # Bright Yellow
  base14: "afd67a"    # Bright Green
  base15: "97d8f8"    # Bright Cyan
  base16: "8db6fa"    # Bright Blue
  base17: "c8acf8"    # Bright Magenta
```

---

## Validation Rules

```go
// internal/domain/validation.go

// Luminance computes WCAG 2.1 relative luminance.
// For each channel c in [0,1]:
//   if c <= 0.04045: c' = c / 12.92
//   else:            c' = ((c + 0.055) / 1.055) ^ 2.4
// L = 0.2126*R' + 0.7152*G' + 0.0722*B'
func (c Color) Luminance() float64 { /* ... */ }

func ValidatePalette(p Palette) []string { /* ... */ }
```

**Rules:**

1. **Completeness**: All 24 slots must have non-zero colors (unless IsNone).
2. **Luminance ordering (dark)**: base00.Luminance < base05.Luminance.
3. **Luminance ordering (light)**: base00.Luminance > base05.Luminance.
4. **Neutral ramp monotonicity**: Luminance should generally increase base00 → base07. Warning.
5. **Bright variants brighter**: base12–base17 luminance ≥ base08–base0E. Warning.

---

## Complete Token Inventory (~88 tokens)

### Surface tokens (11)

| Token Path                        | Derivation                            |
|-----------------------------------|---------------------------------------|
| `surface.background`              | `base00`                              |
| `surface.background.raised`       | `base01`                              |
| `surface.background.sunken`       | `base10`                              |
| `surface.background.darkest`      | `base11`                              |
| `surface.background.highlight`    | `base02`                              |
| `surface.background.selection`    | `BlendBg(base0D, base00, 0.30)`       |
| `surface.background.search`       | `BlendBg(base0A, base00, 0.30)`       |
| `surface.background.overlay`      | `base10`                              |
| `surface.background.popup`        | `base10`                              |
| `surface.background.sidebar`      | `base10`                              |
| `surface.background.statusbar`    | `base10`                              |

### Text tokens (7)

| Token Path          | Derivation                          |
|---------------------|-------------------------------------|
| `text.primary`      | `base05`                            |
| `text.secondary`    | `base04`                            |
| `text.muted`        | `base03`                            |
| `text.subtle`       | `BlendBg(base03, base00, 0.50)`     |
| `text.inverse`      | `base00`                            |
| `text.overlay`      | `base06`                            |
| `text.sidebar`      | `base04`                            |

### Status tokens (6)

| Token Path        | Derivation  |
|-------------------|-------------|
| `status.error`    | `base12`    |
| `status.warning`  | `base13`    |
| `status.success`  | `base14`    |
| `status.info`     | `base15`    |
| `status.hint`     | `base15`    |
| `status.todo`     | `base0D`    |

### Diff tokens (9)

| Token Path            | Derivation                            |
|-----------------------|---------------------------------------|
| `diff.added.fg`       | `base14`                              |
| `diff.added.bg`       | `BlendBg(base0B, base00, 0.25)`       |
| `diff.added.sign`     | `base14`                              |
| `diff.deleted.fg`     | `base12`                              |
| `diff.deleted.bg`     | `BlendBg(base08, base00, 0.25)`       |
| `diff.deleted.sign`   | `base12`                              |
| `diff.changed.fg`     | `base16`                              |
| `diff.changed.bg`     | `BlendBg(base0D, base00, 0.15)`       |
| `diff.ignored`        | `base03`                              |

### Syntax tokens (14)

| Token Path              | Derivation              |
|-------------------------|-------------------------|
| `syntax.keyword`        | `base0E`                |
| `syntax.string`         | `base0B`                |
| `syntax.function`       | `base0D`                |
| `syntax.comment`        | `base03` + italic       |
| `syntax.variable`       | `base05`                |
| `syntax.constant`       | `base09`                |
| `syntax.operator`       | `base16`                |
| `syntax.type`           | `base0A`                |
| `syntax.number`         | `base09`                |
| `syntax.tag`            | `base08`                |
| `syntax.property`       | `base14`                |
| `syntax.parameter`      | `base13`                |
| `syntax.regexp`         | `base0C`                |
| `syntax.escape`         | `base0E`                |
| `syntax.constructor`    | `base17`                |

### Markup tokens (10)

| Token Path                  | Derivation             |
|-----------------------------|------------------------|
| `markup.heading`            | `base0D` + bold        |
| `markup.link`               | `base0C`               |
| `markup.code`               | `base0B`               |
| `markup.bold`               | bold (fg inherited)    |
| `markup.italic`             | italic (fg inherited)  |
| `markup.strikethrough`      | strikethrough          |
| `markup.quote`              | `base03` + italic      |
| `markup.list.bullet`        | `base09`               |
| `markup.list.checked`       | `base0B`               |
| `markup.list.unchecked`     | `base0D`               |

### Accent, border, scrollbar tokens (11)

| Token Path              | Derivation                                 |
|-------------------------|--------------------------------------------|
| `accent.primary`        | `base0D`                                   |
| `accent.secondary`      | `base0E`                                   |
| `accent.foreground`     | `base00`                                   |
| `border.default`        | `BlendBg(base03, base00, 0.40)`            |
| `border.focus`          | `BlendBg(base0D, base00, 0.70)`            |
| `border.muted`          | `base01`                                   |
| `scrollbar.thumb`       | `BlendBg(base03, base00, 0.40)`            |
| `scrollbar.track`       | `base01`                                   |
| `state.hover`           | `surface.background.highlight` (alias)     |
| `state.active`          | `BlendBg(base0D, base00, 0.20)`            |
| `state.disabled.fg`     | `text.muted` (alias)                       |

### Git tokens (4)

| Token Path         | Derivation  |
|--------------------|-------------|
| `git.added`        | `base0B`    |
| `git.modified`     | `base0D`    |
| `git.deleted`      | `base08`    |
| `git.ignored`      | `base03`    |

### Terminal ANSI colors (16)

| Token Path           | Derivation | ANSI Index |
|----------------------|------------|------------|
| `terminal.black`     | `base01`   | 0          |
| `terminal.red`       | `base08`   | 1          |
| `terminal.green`     | `base0B`   | 2          |
| `terminal.yellow`    | `base0A`   | 3          |
| `terminal.blue`      | `base0D`   | 4          |
| `terminal.magenta`   | `base0E`   | 5          |
| `terminal.cyan`      | `base0C`   | 6          |
| `terminal.white`     | `base05`   | 7          |
| `terminal.brblack`   | `base03`   | 8          |
| `terminal.brred`     | `base12`   | 9          |
| `terminal.brgreen`   | `base14`   | 10         |
| `terminal.bryellow`  | `base13`   | 11         |
| `terminal.brblue`    | `base16`   | 12         |
| `terminal.brmagenta` | `base17`   | 13         |
| `terminal.brcyan`    | `base15`   | 14         |
| `terminal.brwhite`   | `base07`   | 15         |

---

## Output File Names (per target)

| Target  | Mapping File            | Output File  |
|---------|-------------------------|--------------|
| vim     | `vim-mapping.yaml`      | `style.lua`  |
| css     | `css-mapping.yaml`      | `style.css`  |
| gtk     | `gtk-mapping.yaml`      | `gtk.css`    |
| qss     | `qss-mapping.yaml`      | `style.qss`  |
| stylix  | `stylix-mapping.yaml`   | `style.json` |

### Symlinks (at `~/.config/flair/`)

| Symlink          | Target                         |
|------------------|--------------------------------|
| `style.lua`      | `<theme>/style.lua`            |
| `style.css`      | `<theme>/style.css`            |
| `gtk.css`        | `<theme>/gtk.css`              |
| `style.qss`      | `<theme>/style.qss`            |
| `style.json`     | `<theme>/style.json`           |

---

## Feature Specifications (BDD)

### Feature 1.2: Domain — Color value object

```gherkin
Feature: Color representation and parsing

  Scenario: Parse a 6-digit hex color
    Given the hex string "#7aa2f7"
    When I parse it as a Color
    Then the RGB values should be R=122 G=162 B=247

  Scenario: Parse without hash prefix
    Given the hex string "7aa2f7"
    Then the RGB values should be R=122 G=162 B=247

  Scenario: Parse 3-digit shorthand "#f00" → R=255 G=0 B=0
  Scenario: Reject invalid hex "#zzzzzz" → ParseError
  Scenario: Reject wrong length "#12345" → ParseError
  Scenario: Format as hex → "#7aa2f7"
  Scenario: RGB to HSL (R=255 G=0 B=0 → H=0 S=1 L=0.5)
  Scenario: Round-trip "#bb9af7" → HSL → RGB → "#bb9af7"
  Scenario: NONE sentinel → IsNone=true
  Scenario: Luminance("#ffffff") ≈ 1.0, Luminance("#000000") ≈ 0.0
```

### Feature 1.3: Domain — Color operations

```gherkin
Feature: Color blending and manipulation

  Scenario: Blend 50% "#000000"+"#ffffff" → "#808080"
  Scenario: Blend 0.0 returns source, 1.0 returns target
  Scenario: BlendBg(fg, bg, 0.25) = Blend(bg, fg, 0.25)
  Scenario: Lighten, Darken, Desaturate, ShiftHue edge cases
```

### Feature 1.4: Domain — Palette entity

```gherkin
Feature: Base24 palette entity

  Scenario: Construct full palette (24 colors)
  Scenario: Access by name ("base0D") and index (13)
  Scenario: Base16 fallbacks (only 16 → 24 via table)
  Scenario: Missing required slot → ParseError
```

### Feature 1.9: Domain — Schema versions

```gherkin
Feature: Schema version management

  Scenario: CurrentVersion returns correct version for each FileKind
  Scenario: All FileKind constants have version > 0
```

### Feature 1.14: Adapter — ThemeStore (filesystem)

```gherkin
Feature: ThemeStore filesystem operations

  Scenario: EnsureThemeDir creates directory
  Scenario: ListThemes returns sorted directory names
  Scenario: Select creates symlinks to theme output files
  Scenario: Select replaces existing symlinks
  Scenario: SelectedTheme reads symlink target

  Scenario: OpenWriter creates writable file
    Given theme "tokyonight" exists
    When OpenWriter("tokyonight", "universal.yaml")
    Then returns io.WriteCloser
    And writing and closing produces the file on disk

  Scenario: OpenReader reads existing file
    Given theme "tokyonight" has "universal.yaml"
    When OpenReader("tokyonight", "universal.yaml")
    Then returns io.ReadCloser with file contents

  Scenario: OpenReader/OpenWriter round-trip
    Given data written via OpenWriter
    When read via OpenReader
    Then bytes match

  Scenario: FileExists and FileMtime
```

### Feature 1.15: Adapter — Built-in palettes

```gherkin
Feature: Built-in palette source

  Scenario: List returns all embedded palette names
    When List()
    Then returns ["catppuccin-mocha", "gruvbox-dark", "tokyo-night-dark"] (sorted)

  Scenario: Get returns YAML bytes
    When Get("tokyo-night-dark")
    Then returns valid base24 YAML bytes

  Scenario: Get unknown name → error
    When Get("nonexistent")
    Then returns error

  Scenario: Has returns true for built-in
    When Has("tokyo-night-dark")
    Then returns true

  Scenario: Has returns false for unknown
    When Has("my-custom-theme")
    Then returns false
```

### Feature 1.16–1.17: Adapter — Reader/Writer wrappers

```gherkin
Feature: VersionedWriter

  Scenario: Prepends schema header on first write
    Given a VersionedWriter wrapping a bytes.Buffer for FileKindUniversal
    When I write "tokens:\n  ..."
    Then the buffer contains "schema_version: 1\nkind: universal\n..." before the content

  Scenario: Header written only once
    Given a VersionedWriter
    When I write twice
    Then the header appears exactly once at the start

  Scenario: Correct version per file kind
    Given VersionedWriters for each FileKind
    Then each prepends the correct schema_version from domain.CurrentVersion()

Feature: ValidatingReader

  Scenario: Valid schema version — passes through
    Given YAML with schema_version: 1 for FileKindUniversal
    When wrapped in ValidatingReader and read
    Then all bytes are available (including the header)

  Scenario: Outdated schema version — SchemaVersionError
    Given YAML with schema_version: 0 for FileKindUniversal (current is 1)
    When wrapped in ValidatingReader and read
    Then returns SchemaVersionError with NeedsUpgrade=false

  Scenario: Future schema version — SchemaVersionError
    Given YAML with schema_version: 99
    When wrapped in ValidatingReader and read
    Then returns SchemaVersionError with NeedsUpgrade=true

  Scenario: Composable with other readers
    Given ValidatingReader wrapping bufio.Reader wrapping bytes.Reader
    Then reads validate and pass through correctly
```

### Feature 2.1–2.5: Token derivation + universal.yaml

```gherkin
Feature: Token derivation and universal file

  Scenario: Full derivation → ≥87 tokens, all paths have non-nil Color
  Scenario: Surface, text, status, diff, syntax, markup tokens match tables

  Scenario: Write universal.yaml via io.Writer
    Given a TokenSet and an io.Writer (bytes.Buffer)
    When fileio.WriteUniversal(w, tokenSet)
    Then output is valid YAML with all tokens

  Scenario: VersionedWriter prepends header
    Given a VersionedWriter wrapping a buffer
    When fileio.WriteUniversal(vw, tokenSet)
    Then output starts with schema_version: 1, kind: universal

  Scenario: Read universal.yaml via io.Reader round-trip
    Given written universal.yaml bytes
    When fileio.ReadUniversal(bytes.NewReader(data))
    Then TokenSet matches original

  Scenario: ValidatingReader catches version mismatch
    Given universal.yaml with schema_version: 99
    When wrapped in ValidatingReader and read
    Then SchemaVersionError with NeedsUpgrade=true
```

### Feature 3.1: Adapter — Stylix target

```gherkin
Feature: Stylix (simplest, validates pipeline)

  Scenario: Mapper → ≥60 keys, raw palette passthrough

  Scenario: Mapping file write via io.Writer + read via io.Reader round-trip
    Given a StylixTheme
    When written to bytes.Buffer and read back from bytes.Reader
    Then StylixTheme matches, and VersionedWriter adds schema_version

  Scenario: Generator writes to io.Writer
    Given a StylixTheme and a bytes.Buffer
    When generator.Generate(buf, theme)
    Then buffer contains valid JSON, sorted keys, 2-space indent, hex colors
```

### Feature 3.3: Adapter — Vim target

```gherkin
Feature: Vim (most complex)

  Scenario: ≥200 highlight groups (base, treesitter, LSP, diagnostic, plugin, markup)
  Scenario: 16 terminal ANSI colors
  Scenario: Mapping file write/read via io.Writer/io.Reader round-trip
  Scenario: Generator writes to io.Writer → hi clear, colors_name, nvim_set_hl, links, terminal_color_N
```

### Feature 3.2, 3.4, 3.5: CSS, GTK, QSS

```gherkin
Feature: CSS — custom properties + element rules, mapping write/read via io.Reader/io.Writer
Feature: GTK — @define-color + widget rules, mapping write/read via io.Reader/io.Writer
Feature: QSS — widget + pseudo-state rules, all literal hex, mapping write/read via io.Reader/io.Writer
```

### Feature 3.6: Application — GenerateTheme

```gherkin
Feature: GenerateTheme orchestration

  Scenario: Full pipeline from file → 12 files
    Given a palette file path
    When Execute("./my-palette.yaml", "my-theme", "")
    Then theme dir has: palette.yaml, universal.yaml,
         5 mapping files, 5 output files (12 files total)

  Scenario: Full pipeline from built-in name
    When Execute("tokyo-night-dark", "", "")
    Then theme dir "tokyo-night-dark" created with all 12 files

  Scenario: Built-in name infers theme name
    When Execute("tokyo-night-dark", "", "")
    Then theme name is "tokyo-night-dark"

  Scenario: Single target filter → 4 files
  Scenario: Theme dir created if missing
  Scenario: All files have correct schema versions
  Scenario: One target failure doesn't block others
```

### Feature 4.3: CLI — select

```gherkin
Feature: Theme selection

  Scenario: Creates symlinks to all 5 output files
  Scenario: Non-existent theme → error
  Scenario: Incomplete theme → error listing missing files
```

### Feature 4.7: CLI — list

```gherkin
Feature: Theme listing

  Scenario: Shows installed themes with selected marker
    Given themes "gruvbox" and "tokyonight", "tokyonight" selected
    When flair list
    Then output shows:
      "  gruvbox (dark)"
      "* tokyonight (dark)"

  Scenario: No themes → helpful message

  Scenario: List built-in palettes
    When flair list --builtins
    Then output shows:
      "catppuccin-mocha"
      "gruvbox-dark"
      "tokyo-night-dark"
```

### Feature 4.8: CLI — regenerate

```gherkin
Feature: Partial regeneration

  Scenario: Edit palette.yaml → everything downstream regenerated
  Scenario: Edit universal.yaml → all mappings + outputs regenerated
  Scenario: Edit vim-mapping.yaml → only style.lua regenerated
  Scenario: No edits detected → "nothing to do"
```

### Feature 5.1–5.6: End-to-end

```gherkin
Feature: End-to-end

  Scenario: Tokyo Night Dark → full theme directory, all golden files match
  Scenario: Generate from built-in name → identical to generate from file
  Scenario: 2 additional schemes → outputs valid
  Scenario: Regenerate after universal.yaml edit → correct partial rebuild
  Scenario: Deterministic (byte-identical on repeated runs)
  Scenario: go-arch-lint check → zero violations
  Scenario: All built-in palettes parse and validate cleanly
```

---

## Dependency Summary

```
Layer          Imports From           Never Imports From
─────────────  ─────────────────────  ──────────────────
domain         (stdlib only)          ports, application, adapters, cmd
ports          domain                 application, adapters, cmd
application    ports, domain          adapters, cmd
adapters/*     ports, domain          other adapters/*, application, cmd
config         (stdlib, vendor)       domain, ports, application, adapters
cmd            everything             (composition root)
```

External dependency: `gopkg.in/yaml.v3` — used by `adapters/yaml` and `adapters/fileio` only.

---

## Implementation Order Rationale

1. **Phase 1** builds the full inner core: domain types, error types, schema versions, port interfaces (all io.Reader/io.Writer based), port file structs, port theme DTOs, YAML palette parser, ThemeStore adapter, built-in palettes adapter, VersionedWriter and ValidatingReader wrappers. After Phase 1 you can parse a palette (from file or built-in via io.Reader) and manage theme directories. `go-arch-lint check` passes.

2. **Phase 2** adds the deriver adapter, universal.yaml read/write, and the DeriveTheme use case. After Phase 2 you can go from palette → universal.yaml.

3. **Phase 3** builds targets one at a time. Stylix first (simplest), then CSS, Vim, GTK, QSS. Each target includes mapper, mapping file read/write, and generator. GenerateTheme use case added once first target works.

4. **Phase 4** builds the CLI. `generate` and `select` first (core workflow), then `validate`, `preview`, `list`, `init`, `regenerate`.

5. **Phase 5** is golden files, regeneration tests, arch lint, docs.

`go-arch-lint check` passes at every commit.

# Phase 6: Charm Integration & Style Viewer

## Tasks

#### 6.1 — Public Package Foundation (`pkg/`)

- [x] 6.1a — Create `pkg/` directory structure (independent from `/internal`)
- [x] 6.1b — Define `pkg/flair/theme.go` — public Theme type (read-only, no internal deps)
- [x] 6.1c — Define `pkg/flair/loader.go` — load selected theme from `~/.config/flair`
- [x] 6.1d — Define `pkg/flair/colors.go` — public color accessors (surface, text, status, etc.)
- [x] 6.1e — Unit tests for theme loading and color accessors

#### 6.2 — Lipgloss Adapter (`pkg/charm/lipgloss`)

- [x] 6.2a — `pkg/charm/lipgloss/styles.go` — LipglossStyles struct with pre-configured styles
- [x] 6.2b — `pkg/charm/lipgloss/builder.go` — NewStyles(theme) → LipglossStyles
- [x] 6.2c — Surface styles (Background, Raised, Sunken, Overlay, Popup)
- [x] 6.2d — Text styles (Primary, Secondary, Muted, Inverse)
- [x] 6.2e — Status styles (Error, Warning, Success, Info)
- [x] 6.2f — Border styles (Default, Focus, Muted)
- [x] 6.2g — Component styles (Button, Input, List, Table, Dialog)
- [x] 6.2h — State styles (Hover, Active, Disabled, Selected)
- [x] 6.2i — Unit tests for all lipgloss style builders

#### 6.3 — Style Viewer (`pkg/flair/viewer`)

- [x] 6.3a — `pkg/flair/viewer/model.go` — Bubbletea model for style viewer
- [x] 6.3b — `pkg/flair/viewer/view.go` — Render style showcase pages
- [x] 6.3c — Theme selector component (list available themes, live preview)
- [x] 6.3d — Palette display page (base00–base17 with color swatches)
- [x] 6.3e — Token display page (semantic tokens grouped by category)
- [x] 6.3f — Lipgloss component showcase page (buttons, inputs, tables, etc.)
- [x] 6.3g — Dynamic theme switching (select theme → update all styles live)
- [x] 6.3h — Use flair token names as example labels in component showcase
- [x] 6.3i — `pkg/flair/viewer/run.go` — Public Run() function for embedding in other CLIs
- [x] 6.3j — Keyboard navigation (j/k scroll, Enter select, q quit, Tab switch pages)
- [x] 6.3k — Unit tests for viewer model and view rendering

#### 6.4 — CLI Integration

- [x] 6.4a — Update `select` command: no args launches style viewer
- [x] 6.4b — `select <theme-name>` retains existing symlink behavior
- [x] 6.4c — Add `--viewer` flag to force viewer mode with theme pre-selected
- [x] 6.4d — Integration tests for select command variants

#### 6.5 — Documentation & Polish

- [x] 6.5a — README section: Using flair in your CLI
- [x] 6.5b — Example: minimal CLI with flair theming
- [x] 6.5c — Example: embedding style viewer in a CLI
- [x] 6.5d — godoc comments for all public APIs in `pkg/`
- [x] 6.5e — BDD feature files for Phase 6 features

## Notes

Integrates with Charmbracelet's lipgloss library for themed TUI components.

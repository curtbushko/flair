# Phase 4: CLI (Composition Root)

## Tasks

- [x] 4.1 — CLI scaffolding and dependency wiring in cmd
- [x] 4.2 — `generate` command (palette file → full theme directory)
  - [x] 4.2a — `--dir` flag (custom config dir, default `~/.config/flair`)
  - [x] 4.2b — `--target` flag (single target output)
- [x] 4.3 — `select` command (switch active theme via symlinks)
- [x] 4.4 — `validate` command (lint all files in a theme directory)
- [x] 4.5 — `preview` command (ANSI color preview in terminal)
- [x] 4.6 — `init` command (scaffold a new theme directory with palette template)
- [x] 4.7 — `list` command (show available themes, mark selected)
- [x] 4.8 — `regenerate` command (re-derive downstream from edited intermediate files)

## Notes

CLI composition root wires together all adapters and use cases.

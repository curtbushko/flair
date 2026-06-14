# Flair Project Instructions

## Overview

Flair is a Go CLI tool for managing color palettes and themes. It uses base24 color palettes stored in YAML files.

## Palette Structure

Palettes are stored in `pkg/flair/palettes/*.yaml` and follow the base24 color system:

```yaml
system: "base24"
name: "Theme Name"
author: "Author"
variant: "dark"  # or "light"
palette:
  base00: "XXXXXX"  # main background
  base01: "XXXXXX"  # darker background
  base02: "XXXXXX"  # selection background
  base03: "XXXXXX"  # comments/muted text
  base04: "XXXXXX"  # dark foreground
  base05: "XXXXXX"  # main foreground
  base06: "XXXXXX"  # light foreground
  base07: "XXXXXX"  # lightest foreground
  base08: "XXXXXX"  # red
  base09: "XXXXXX"  # orange
  base0A: "XXXXXX"  # yellow
  base0B: "XXXXXX"  # green
  base0C: "XXXXXX"  # cyan
  base0D: "XXXXXX"  # blue
  base0E: "XXXXXX"  # purple/magenta
  base0F: "XXXXXX"  # brown/accent
  base10: "XXXXXX"  # sunken surface (darker than base00)
  base11: "XXXXXX"  # darkest background
  base12: "XXXXXX"  # bright red
  base13: "XXXXXX"  # bright yellow
  base14: "XXXXXX"  # bright green
  base15: "XXXXXX"  # bright cyan
  base16: "XXXXXX"  # bright blue
  base17: "XXXXXX"  # bright magenta
```

## Statusline Override Pattern

All themes MUST use this consistent statusline pattern in their `overrides` section:

```yaml
overrides:
  statusline_a_bg: # base03
    color: "XXXXXX"
  statusline_a_fg: # base00
    color: "XXXXXX"
  statusline_b_bg: # base10
    color: "XXXXXX"
  statusline_b_fg: # base05
    color: "XXXXXX"
  statusline_c_bg: # base01
    color: "XXXXXX"
  statusline_c_fg: # base05
    color: "XXXXXX"
```

### Statusline Pattern Explanation

| Entry | Base | Purpose |
|-------|------|---------|
| statusline_a_bg | base03 | Muted text color used as accent background |
| statusline_a_fg | base00 | Main background used as contrasting text |
| statusline_b_bg | base10 | Sunken surface for middle section |
| statusline_b_fg | base05 | Main foreground text |
| statusline_c_bg | base01 | Darker background, blends with editor |
| statusline_c_fg | base05 | Main foreground text |

This creates a three-tier statusline:
- **Section A**: Most prominent (inverted colors for mode indicator)
- **Section B**: Middle layer (file info, branch, etc.)
- **Section C**: Blends with editor background (position, encoding)

**Exception**: statusline_c entries (statusline_c_bg and statusline_c_fg) can deviate from the base pattern if the comment contains the word 'custom'. This allows themes to use custom colors for section C when the standard pattern doesn't work well.

### Comment Format

Every statusline entry MUST have a comment indicating which base color it maps to:

```yaml
  statusline_a_bg: # base03
    color: "504945"
```

## Adding New Themes

1. Create a new YAML file in `pkg/flair/palettes/`
2. Define all base00-base17 colors
3. Add the statusline overrides using the pattern above
4. Ensure all hex values are UPPERCASE (e.g., `"FF5555"` not `"ff5555"`)
5. Add comments to palette entries describing the color purpose

## Code Style

- All hex color values must be UPPERCASE
- Use Nerd Font icons instead of emojis
- Follow hexagonal architecture patterns

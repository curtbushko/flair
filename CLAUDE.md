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
  base00: "xxxxxx"  # main background
  base01: "xxxxxx"  # darker background
  base02: "xxxxxx"  # selection background
  base03: "xxxxxx"  # comments/muted text
  base04: "xxxxxx"  # dark foreground
  base05: "xxxxxx"  # main foreground
  base06: "xxxxxx"  # light foreground
  base07: "xxxxxx"  # lightest foreground
  base08: "xxxxxx"  # red
  base09: "xxxxxx"  # orange
  base0A: "xxxxxx"  # yellow
  base0B: "xxxxxx"  # green
  base0C: "xxxxxx"  # cyan
  base0D: "xxxxxx"  # blue
  base0E: "xxxxxx"  # purple/magenta
  base0F: "xxxxxx"  # brown/accent
  base10: "xxxxxx"  # sunken surface (darker than base00)
  base11: "xxxxxx"  # darkest background
  base12: "xxxxxx"  # bright red
  base13: "xxxxxx"  # bright yellow
  base14: "xxxxxx"  # bright green
  base15: "xxxxxx"  # bright cyan
  base16: "xxxxxx"  # bright blue
  base17: "xxxxxx"  # bright magenta
```

## Statusline Override Pattern

All themes MUST use this consistent statusline pattern in their `overrides` section:

```yaml
overrides:
  statusline_a_bg: # base0F
    color: "xxxxxx"
  statusline_a_fg: # base00
    color: "xxxxxx"
  statusline_b_bg: # base10
    color: "xxxxxx"
  statusline_b_fg: # base0D
    color: "xxxxxx"
  statusline_c_bg: # base01
    color: "xxxxxx"
  statusline_c_fg: # base0D
    color: "xxxxxx"
```

### Statusline Pattern Explanation

| Entry | Base | Purpose |
|-------|------|---------|
| statusline_a_bg | base0F | Brown/accent color for prominent mode indicator background |
| statusline_a_fg | base00 | Main background used as contrasting text |
| statusline_b_bg | base10 | Sunken surface for middle section |
| statusline_b_fg | base0D | Blue accent for readable foreground text |
| statusline_c_bg | base01 | Darker background, blends with editor |
| statusline_c_fg | base0D | Blue accent for consistent foreground text |

This creates a three-tier statusline:
- **Section A**: Most prominent (accent color background with inverted text for mode indicator)
- **Section B**: Middle layer (sunken background with blue accent text for file info, branch, etc.)
- **Section C**: Blends with editor background (darker background with blue accent text for position, encoding)

### Comment Format

Every statusline entry MUST have a comment indicating which base color it maps to:

```yaml
  statusline_a_bg: # base0F
    color: "a7c080"
  statusline_a_fg: # base00
    color: "272e33"
  statusline_b_fg: # base0D
    color: "7fbbb3"
```

## Adding New Themes

1. Create a new YAML file in `pkg/flair/palettes/`
2. Define all base00-base17 colors with lowercase hex values
3. Add the statusline overrides using the pattern above with lowercase hex values
4. Add comments indicating which base color each statusline entry uses
5. Add comments to palette entries describing the color purpose

## Code Style

- All hex color values must be lowercase (e.g., `"ff5555"` not `"FF5555"`)
- Use Nerd Font icons instead of emojis
- Follow hexagonal architecture patterns

package flair

// Tokenize derives semantic color tokens from a base24 palette.
//
// Tokenize transforms a raw 24-color palette into a Theme containing
// approximately 88 semantic tokens organized into categories:
//
//   - Surface: Background colors for UI elements
//   - Text: Foreground colors for text content
//   - Status: Colors for error, warning, success, info states
//   - Diff: Colors for version control diffs
//   - Syntax: Colors for code syntax highlighting
//   - Markup: Colors for documentation and markdown
//   - Accent/Border: UI accent and border colors
//   - Git: Git status indicator colors
//   - Terminal: 16 ANSI terminal colors
//
// The token derivation follows consistent mapping rules from base24 indices
// to semantic meanings, with some tokens using alpha blending for
// subtle variations (e.g., selection backgrounds).
//
// Tokenize returns nil if the palette is nil.
//
// Example:
//
//	palette, _ := flair.ParsePalette(reader)
//	theme := flair.Tokenize(palette)
//	fmt.Printf("Created theme: %s\n", theme.Name())
//	fmt.Printf("Background: %s\n", theme.Surface().Background.Hex())
func Tokenize(p *Palette) *Theme {
	if p == nil {
		return nil
	}

	colors := make(map[string]Color, 100)

	// Derive all token categories.
	deriveSurface(p, colors)
	deriveText(p, colors)
	deriveStatus(p, colors)
	deriveDiff(p, colors)
	deriveSyntax(p, colors)
	deriveMarkup(p, colors)
	deriveAccentBorder(p, colors)
	deriveGit(p, colors)
	deriveTerminal(p, colors)

	return NewTheme(p.Name(), p.Variant(), colors)
}

// deriveSurface derives surface tokens from the palette.
func deriveSurface(p *Palette, colors map[string]Color) {
	base00 := *p.Base(0x00)
	base01 := *p.Base(0x01)
	base02 := *p.Base(0x02)
	base0A := *p.Base(0x0A)
	base0D := *p.Base(0x0D)
	base10 := *p.Base(0x10)
	base11 := *p.Base(0x11)

	colors["surface.background"] = base00
	colors["surface.background.raised"] = base01
	colors["surface.background.sunken"] = base10
	colors["surface.background.darkest"] = base11
	colors["surface.background.highlight"] = base02
	colors["surface.background.selection"] = BlendBg(base0D, base00, 0.30)
	colors["surface.background.search"] = BlendBg(base0A, base00, 0.30)
	colors["surface.background.overlay"] = base10
	colors["surface.background.popup"] = base10
	colors["surface.background.sidebar"] = base10
	colors["surface.background.statusbar"] = base10
}

// deriveText derives text tokens from the palette.
func deriveText(p *Palette, colors map[string]Color) {
	base00 := *p.Base(0x00)
	base03 := *p.Base(0x03)
	base04 := *p.Base(0x04)
	base05 := *p.Base(0x05)
	base06 := *p.Base(0x06)

	colors["text.primary"] = base05
	colors["text.secondary"] = base04
	colors["text.muted"] = base03
	colors["text.subtle"] = BlendBg(base03, base00, 0.50)
	colors["text.inverse"] = base00
	colors["text.overlay"] = base06
	colors["text.sidebar"] = base04
}

// deriveStatus derives status tokens from the palette.
func deriveStatus(p *Palette, colors map[string]Color) {
	base0D := *p.Base(0x0D)
	base12 := *p.Base(0x12)
	base13 := *p.Base(0x13)
	base14 := *p.Base(0x14)
	base15 := *p.Base(0x15)

	colors["status.error"] = base12
	colors["status.warning"] = base13
	colors["status.success"] = base14
	colors["status.info"] = base15
	colors["status.hint"] = base15
	colors["status.todo"] = base0D
}

// deriveDiff derives diff tokens from the palette.
func deriveDiff(p *Palette, colors map[string]Color) {
	base00 := *p.Base(0x00)
	base03 := *p.Base(0x03)
	base08 := *p.Base(0x08)
	base0B := *p.Base(0x0B)
	base0D := *p.Base(0x0D)
	base12 := *p.Base(0x12)
	base14 := *p.Base(0x14)
	base16 := *p.Base(0x16)

	colors["diff.added.fg"] = base14
	colors["diff.added.bg"] = BlendBg(base0B, base00, 0.25)
	colors["diff.added.sign"] = base14
	colors["diff.deleted.fg"] = base12
	colors["diff.deleted.bg"] = BlendBg(base08, base00, 0.25)
	colors["diff.deleted.sign"] = base12
	colors["diff.changed.fg"] = base16
	colors["diff.changed.bg"] = BlendBg(base0D, base00, 0.15)
	colors["diff.ignored"] = base03
}

// deriveSyntax derives syntax highlighting tokens from the palette.
func deriveSyntax(p *Palette, colors map[string]Color) {
	base03 := *p.Base(0x03)
	base05 := *p.Base(0x05)
	base08 := *p.Base(0x08)
	base09 := *p.Base(0x09)
	base0A := *p.Base(0x0A)
	base0B := *p.Base(0x0B)
	base0C := *p.Base(0x0C)
	base0D := *p.Base(0x0D)
	base0E := *p.Base(0x0E)
	base13 := *p.Base(0x13)
	base14 := *p.Base(0x14)
	base16 := *p.Base(0x16)
	base17 := *p.Base(0x17)

	colors["syntax.keyword"] = base0E
	colors["syntax.string"] = base0B
	colors["syntax.function"] = base0D
	colors["syntax.comment"] = base03
	colors["syntax.variable"] = base05
	colors["syntax.constant"] = base09
	colors["syntax.operator"] = base16
	colors["syntax.type"] = base0A
	colors["syntax.number"] = base09
	colors["syntax.tag"] = base08
	colors["syntax.property"] = base14
	colors["syntax.parameter"] = base13
	colors["syntax.regexp"] = base0C
	colors["syntax.escape"] = base0E
	colors["syntax.constructor"] = base17
}

// deriveMarkup derives markup/documentation tokens from the palette.
func deriveMarkup(p *Palette, colors map[string]Color) {
	base03 := *p.Base(0x03)
	base09 := *p.Base(0x09)
	base0B := *p.Base(0x0B)
	base0C := *p.Base(0x0C)
	base0D := *p.Base(0x0D)

	colors["markup.heading"] = base0D
	colors["markup.link"] = base0C
	colors["markup.code"] = base0B
	colors["markup.quote"] = base03
	colors["markup.list.bullet"] = base09
	colors["markup.list.checked"] = base0B
	colors["markup.list.unchecked"] = base0D
}

// deriveAccentBorder derives accent, border, scrollbar, and state tokens from the palette.
func deriveAccentBorder(p *Palette, colors map[string]Color) {
	base00 := *p.Base(0x00)
	base01 := *p.Base(0x01)
	base02 := *p.Base(0x02)
	base03 := *p.Base(0x03)
	base0D := *p.Base(0x0D)
	base0E := *p.Base(0x0E)

	colors["accent.primary"] = base0D
	colors["accent.secondary"] = base0E
	colors["accent.foreground"] = base00
	colors["border.default"] = BlendBg(base03, base00, 0.40)
	colors["border.focus"] = BlendBg(base0D, base00, 0.70)
	colors["border.muted"] = base01
	colors["scrollbar.thumb"] = BlendBg(base03, base00, 0.40)
	colors["scrollbar.track"] = base01
	colors["state.hover"] = base02 // alias for surface.background.highlight
	colors["state.active"] = BlendBg(base0D, base00, 0.20)
	colors["state.disabled.fg"] = base03 // alias for text.muted
}

// deriveGit derives git-related tokens from the palette.
func deriveGit(p *Palette, colors map[string]Color) {
	base03 := *p.Base(0x03)
	base08 := *p.Base(0x08)
	base0B := *p.Base(0x0B)
	base0D := *p.Base(0x0D)

	colors["git.added"] = base0B
	colors["git.modified"] = base0D
	colors["git.deleted"] = base08
	colors["git.ignored"] = base03
}

// deriveTerminal derives terminal ANSI color tokens from the palette.
func deriveTerminal(p *Palette, colors map[string]Color) {
	colors["terminal.black"] = *p.Base(0x01)
	colors["terminal.red"] = *p.Base(0x08)
	colors["terminal.green"] = *p.Base(0x0B)
	colors["terminal.yellow"] = *p.Base(0x0A)
	colors["terminal.blue"] = *p.Base(0x0D)
	colors["terminal.magenta"] = *p.Base(0x0E)
	colors["terminal.cyan"] = *p.Base(0x0C)
	colors["terminal.white"] = *p.Base(0x05)
	colors["terminal.brblack"] = *p.Base(0x03)
	colors["terminal.brred"] = *p.Base(0x12)
	colors["terminal.brgreen"] = *p.Base(0x14)
	colors["terminal.bryellow"] = *p.Base(0x13)
	colors["terminal.brblue"] = *p.Base(0x16)
	colors["terminal.brmagenta"] = *p.Base(0x17)
	colors["terminal.brcyan"] = *p.Base(0x15)
	colors["terminal.brwhite"] = *p.Base(0x07)
}

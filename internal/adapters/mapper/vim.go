package mapper

import (
	"errors"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// Vim implements ports.Mapper for the Vim/Neovim target.
// It maps a ResolvedTheme into a VimTheme containing base editor highlights,
// treesitter highlights, LSP semantic token links, and diagnostic highlights.
type Vim struct{}

// NewVim returns a new Vim mapper.
func NewVim() *Vim {
	return &Vim{}
}

// Name returns the target name for this mapper.
func (v *Vim) Name() string {
	return "vim"
}

// Map transforms a ResolvedTheme into a *ports.VimTheme with base editor
// highlights, treesitter groups, LSP semantic links, diagnostic groups,
// plugin highlights, markup highlights, and terminal ANSI colors.
func (v *Vim) Map(theme *domain.ResolvedTheme) (ports.MappedTheme, error) {
	if theme == nil {
		return nil, errors.New("vim mapper: nil theme")
	}
	if theme.Palette == nil {
		return nil, errors.New("vim mapper: nil palette")
	}
	if theme.Tokens == nil {
		return nil, errors.New("vim mapper: nil tokens")
	}

	highlights := make(map[string]ports.VimHighlight)

	mapBase(theme, highlights)
	mapTreesitter(theme, highlights)
	mapLSP(highlights)
	mapDiagnostic(theme, highlights)
	mapPlugins(theme, highlights)
	mapMarkup(theme, highlights)

	termColors := mapTerminal(theme)

	return &ports.VimTheme{
		Name:           theme.Name,
		Highlights:     highlights,
		TerminalColors: termColors,
	}, nil
}

// colorOf retrieves a token color as a *domain.Color pointer.
// Returns nil if the token is not found or has IsNone.
func colorOf(theme *domain.ResolvedTheme, path string) *domain.Color {
	tok, ok := theme.Tokens.Get(path)
	if !ok || tok.Color.IsNone {
		return nil
	}
	c := tok.Color
	return &c
}

// mapBase adds standard Vim editor highlight groups to the highlights map.
//
//nolint:funlen // Large mapping table is intentionally in one function for clarity.
func mapBase(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	// Shorthand helpers for common token lookups.
	fg := func(path string) *domain.Color { return colorOf(theme, path) }
	bg := func(path string) *domain.Color { return colorOf(theme, path) }

	// --- Core editor groups ---
	hl["Normal"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background")}
	hl["NormalFloat"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.popup")}
	hl["NormalNC"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background")}
	hl["Comment"] = ports.VimHighlight{Fg: fg("syntax.comment"), Italic: true}
	hl["Cursor"] = ports.VimHighlight{Reverse: true}
	hl["lCursor"] = ports.VimHighlight{Reverse: true}
	hl["CursorIM"] = ports.VimHighlight{Reverse: true}

	// --- UI groups ---
	hl["CursorLine"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["CursorColumn"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["CursorLineNr"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["ColorColumn"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["LineNr"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["SignColumn"] = ports.VimHighlight{Bg: bg("surface.background")}
	hl["FoldColumn"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background")}
	hl["Folded"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.raised")}
	hl["VertSplit"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["WinSeparator"] = ports.VimHighlight{Fg: fg("border.default")}

	// --- Visual / Selection ---
	hl["Visual"] = ports.VimHighlight{Bg: bg("surface.background.selection")}
	hl["VisualNOS"] = ports.VimHighlight{Bg: bg("surface.background.selection")}

	// --- Search ---
	hl["Search"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.search")}
	hl["IncSearch"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary")}
	hl["CurSearch"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary")}
	hl["Substitute"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("status.error")}

	// --- Popup / completion menu ---
	hl["Pmenu"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.popup")}
	hl["PmenuSel"] = ports.VimHighlight{Bg: bg("surface.background.selection")}
	hl["PmenuSbar"] = ports.VimHighlight{Bg: bg("surface.background.raised")}
	hl["PmenuThumb"] = ports.VimHighlight{Bg: bg("scrollbar.thumb")}
	hl["FloatBorder"] = ports.VimHighlight{Fg: fg("border.default"), Bg: bg("surface.background.popup")}
	hl["FloatTitle"] = ports.VimHighlight{Fg: fg("accent.primary"), Bg: bg("surface.background.popup"), Bold: true}

	// --- Tab line ---
	hl["TabLine"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.raised")}
	hl["TabLineSel"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background"), Bold: true}
	hl["TabLineFill"] = ports.VimHighlight{Bg: bg("surface.background.sunken")}

	// --- Status line ---
	hl["StatusLine"] = ports.VimHighlight{Fg: fg("text.secondary"), Bg: bg("surface.background.statusbar")}
	hl["StatusLineNC"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.statusbar")}
	hl["WildMenu"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary")}

	// --- Messages ---
	hl["ErrorMsg"] = ports.VimHighlight{Fg: fg("status.error")}
	hl["WarningMsg"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["ModeMsg"] = ports.VimHighlight{Fg: fg("text.primary"), Bold: true}
	hl["MoreMsg"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["Question"] = ports.VimHighlight{Fg: fg("accent.primary")}

	// --- Diff ---
	hl["DiffAdd"] = ports.VimHighlight{Bg: bg("diff.added.bg")}
	hl["DiffChange"] = ports.VimHighlight{Bg: bg("diff.changed.bg")}
	hl["DiffDelete"] = ports.VimHighlight{Fg: fg("diff.deleted.fg"), Bg: bg("diff.deleted.bg")}
	hl["DiffText"] = ports.VimHighlight{Bg: bg("diff.changed.bg"), Bold: true}

	// --- Spelling ---
	hl["SpellBad"] = ports.VimHighlight{Sp: fg("status.error"), Undercurl: true}
	hl["SpellCap"] = ports.VimHighlight{Sp: fg("status.warning"), Undercurl: true}
	hl["SpellRare"] = ports.VimHighlight{Sp: fg("accent.secondary"), Undercurl: true}
	hl["SpellLocal"] = ports.VimHighlight{Sp: fg("status.info"), Undercurl: true}

	// --- Syntax base groups ---
	hl["Constant"] = ports.VimHighlight{Fg: fg("syntax.constant")}
	hl["String"] = ports.VimHighlight{Fg: fg("syntax.string")}
	hl["Character"] = ports.VimHighlight{Fg: fg("syntax.string")}
	hl["Number"] = ports.VimHighlight{Fg: fg("syntax.number")}
	hl["Boolean"] = ports.VimHighlight{Fg: fg("syntax.constant")}
	hl["Float"] = ports.VimHighlight{Fg: fg("syntax.number")}
	hl["Identifier"] = ports.VimHighlight{Fg: fg("syntax.variable")}
	hl["Function"] = ports.VimHighlight{Fg: fg("syntax.function")}
	hl["Statement"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Conditional"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Repeat"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Label"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Operator"] = ports.VimHighlight{Fg: fg("syntax.operator")}
	hl["Keyword"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Exception"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["PreProc"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Include"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Define"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Macro"] = ports.VimHighlight{Fg: fg("syntax.constant")}
	hl["PreCondit"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Type"] = ports.VimHighlight{Fg: fg("syntax.type")}
	hl["StorageClass"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["Structure"] = ports.VimHighlight{Fg: fg("syntax.type")}
	hl["Typedef"] = ports.VimHighlight{Fg: fg("syntax.type")}
	hl["Special"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["SpecialChar"] = ports.VimHighlight{Fg: fg("syntax.escape")}
	hl["Tag"] = ports.VimHighlight{Fg: fg("syntax.tag")}
	hl["Delimiter"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["SpecialComment"] = ports.VimHighlight{Fg: fg("syntax.comment"), Italic: true}
	hl["Debug"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["Underlined"] = ports.VimHighlight{Underline: true}
	hl["Ignore"] = ports.VimHighlight{}
	hl["Error"] = ports.VimHighlight{Fg: fg("status.error")}
	hl["Todo"] = ports.VimHighlight{Fg: fg("status.todo"), Bold: true}

	// --- Miscellaneous ---
	hl["MatchParen"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["NonText"] = ports.VimHighlight{Fg: fg("text.subtle")}
	hl["SpecialKey"] = ports.VimHighlight{Fg: fg("text.subtle")}
	hl["Conceal"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["Directory"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["Title"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["EndOfBuffer"] = ports.VimHighlight{Fg: fg("surface.background")}

	// --- Markup (from semantic tokens) ---
	hl["markdownH1"] = ports.VimHighlight{Fg: fg("markup.heading"), Bold: true}
	hl["markdownH2"] = ports.VimHighlight{Fg: fg("markup.heading"), Bold: true}
	hl["markdownH3"] = ports.VimHighlight{Fg: fg("markup.heading"), Bold: true}
	hl["markdownH4"] = ports.VimHighlight{Fg: fg("markup.heading"), Bold: true}
	hl["markdownH5"] = ports.VimHighlight{Fg: fg("markup.heading"), Bold: true}
	hl["markdownH6"] = ports.VimHighlight{Fg: fg("markup.heading"), Bold: true}
	hl["markdownUrl"] = ports.VimHighlight{Fg: fg("markup.link"), Underline: true}
	hl["markdownCode"] = ports.VimHighlight{Fg: fg("markup.code")}
	hl["markdownCodeBlock"] = ports.VimHighlight{Fg: fg("markup.code")}
	hl["markdownBold"] = ports.VimHighlight{Bold: true}
	hl["markdownItalic"] = ports.VimHighlight{Italic: true}
	hl["markdownListMarker"] = ports.VimHighlight{Fg: fg("markup.list.bullet")}
	hl["markdownBlockquote"] = ports.VimHighlight{Fg: fg("markup.quote"), Italic: true}
}

// treesitterMapping maps a treesitter highlight group name to a semantic token
// path. An empty tokenPath means the group uses a Link instead.
type treesitterMapping struct {
	group     string
	tokenPath string
	link      string
	italic    bool
	bold      bool
}

// treesitterMappings defines the treesitter highlight groups and their semantic
// token sources. Uses table-driven approach for maintainability.
//
//nolint:dupl // Intentional structural similarity between mapping tables.
var treesitterMappings = []treesitterMapping{
	// Core treesitter groups mapped directly to semantic tokens.
	{group: "@keyword", tokenPath: "syntax.keyword"},
	{group: "@keyword.function", tokenPath: "syntax.keyword"},
	{group: "@keyword.operator", tokenPath: "syntax.operator"},
	{group: "@keyword.return", tokenPath: "syntax.keyword"},
	{group: "@keyword.coroutine", tokenPath: "syntax.keyword"},
	{group: "@keyword.exception", tokenPath: "syntax.keyword"},
	{group: "@keyword.conditional", tokenPath: "syntax.keyword"},
	{group: "@keyword.repeat", tokenPath: "syntax.keyword"},
	{group: "@keyword.import", tokenPath: "syntax.keyword"},
	{group: "@string", tokenPath: "syntax.string"},
	{group: "@string.escape", tokenPath: "syntax.escape"},
	{group: "@string.regex", tokenPath: "syntax.regexp"},
	{group: "@string.special", tokenPath: "syntax.escape"},
	{group: "@character", tokenPath: "syntax.string"},
	{group: "@character.special", tokenPath: "syntax.escape"},
	{group: "@function", tokenPath: "syntax.function"},
	{group: "@function.builtin", tokenPath: "syntax.function"},
	{group: "@function.call", tokenPath: "syntax.function"},
	{group: "@function.macro", tokenPath: "syntax.constant"},
	{group: "@method", tokenPath: "syntax.function"},
	{group: "@method.call", tokenPath: "syntax.function"},
	{group: "@variable", tokenPath: "syntax.variable"},
	{group: "@variable.builtin", tokenPath: "syntax.variable"},
	{group: "@variable.parameter", tokenPath: "syntax.parameter"},
	{group: "@variable.member", tokenPath: "syntax.property"},
	{group: "@type", tokenPath: "syntax.type"},
	{group: "@type.builtin", tokenPath: "syntax.type"},
	{group: "@type.definition", tokenPath: "syntax.type"},
	{group: "@type.qualifier", tokenPath: "syntax.keyword"},
	{group: "@constant", tokenPath: "syntax.constant"},
	{group: "@constant.builtin", tokenPath: "syntax.constant"},
	{group: "@constant.macro", tokenPath: "syntax.constant"},
	{group: "@number", tokenPath: "syntax.number"},
	{group: "@number.float", tokenPath: "syntax.number"},
	{group: "@boolean", tokenPath: "syntax.constant"},
	{group: "@operator", tokenPath: "syntax.operator"},
	{group: "@punctuation.bracket", tokenPath: "text.secondary"},
	{group: "@punctuation.delimiter", tokenPath: "text.primary"},
	{group: "@punctuation.special", tokenPath: "accent.primary"},
	{group: "@tag", tokenPath: "syntax.tag"},
	{group: "@tag.attribute", tokenPath: "syntax.property"},
	{group: "@tag.delimiter", tokenPath: "text.secondary"},
	{group: "@property", tokenPath: "syntax.property"},
	{group: "@parameter", tokenPath: "syntax.parameter"},
	{group: "@constructor", tokenPath: "syntax.constructor"},
	{group: "@namespace", tokenPath: "syntax.type"},
	{group: "@module", tokenPath: "syntax.type"},
	{group: "@label", tokenPath: "syntax.keyword"},
	{group: "@include", tokenPath: "syntax.keyword"},
	{group: "@exception", tokenPath: "syntax.keyword"},
	{group: "@define", tokenPath: "syntax.keyword"},
	{group: "@preproc", tokenPath: "syntax.keyword"},

	// Link groups (reference other highlight groups).
	{group: "@comment", link: "Comment"},
	{group: "@text", link: "Normal"},
	{group: "@text.strong", bold: true},
	{group: "@text.emphasis", italic: true},
	{group: "@text.underline", link: "Underlined"},
	{group: "@text.strike", link: "Underlined"},
	{group: "@text.title", link: "Title"},
	{group: "@text.uri", link: "Underlined"},
	{group: "@text.todo", link: "Todo"},
	{group: "@text.note", link: "Todo"},
	{group: "@text.warning", link: "WarningMsg"},
	{group: "@text.danger", link: "ErrorMsg"},

	// Markup treesitter groups.
	{group: "@markup.heading", tokenPath: "markup.heading", bold: true},
	{group: "@markup.link", tokenPath: "markup.link"},
	{group: "@markup.link.url", tokenPath: "markup.link"},
	{group: "@markup.raw", tokenPath: "markup.code"},
	{group: "@markup.list", tokenPath: "markup.list.bullet"},
	{group: "@markup.strong", bold: true},
	{group: "@markup.italic", italic: true},
	{group: "@markup.strikethrough", tokenPath: "markup.quote"},
	{group: "@markup.quote", tokenPath: "markup.quote", italic: true},
}

// mapTreesitter adds treesitter highlight groups to the highlights map.
func mapTreesitter(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	for _, m := range treesitterMappings {
		if m.link != "" {
			hl[m.group] = ports.VimHighlight{Link: m.link}
			continue
		}

		h := ports.VimHighlight{
			Italic: m.italic,
			Bold:   m.bold,
		}

		if m.tokenPath != "" {
			h.Fg = colorOf(theme, m.tokenPath)
		}

		hl[m.group] = h
	}
}

// lspMapping maps an LSP semantic token type to its linked treesitter group.
type lspMapping struct {
	group string
	link  string
}

// lspMappings defines the LSP semantic token link groups.
var lspMappings = []lspMapping{
	{"@lsp.type.function", "@function"},
	{"@lsp.type.method", "@function"},
	{"@lsp.type.variable", "@variable"},
	{"@lsp.type.parameter", "@parameter"},
	{"@lsp.type.property", "@property"},
	{"@lsp.type.keyword", "@keyword"},
	{"@lsp.type.type", "@type"},
	{"@lsp.type.namespace", "@type"},
	{"@lsp.type.enum", "@type"},
	{"@lsp.type.enumMember", "@constant"},
	{"@lsp.type.struct", "@type"},
	{"@lsp.type.class", "@type"},
	{"@lsp.type.interface", "@type"},
	{"@lsp.type.string", "@string"},
	{"@lsp.type.number", "@number"},
	{"@lsp.type.operator", "@operator"},
	{"@lsp.type.comment", "@comment"},
	{"@lsp.type.decorator", "@function"},
	{"@lsp.type.macro", "@constant"},
	{"@lsp.type.typeParameter", "@type"},
	{"@lsp.type.event", "@type"},
	{"@lsp.type.modifier", "@keyword"},
	{"@lsp.type.regexp", "@string.regex"},
}

// mapLSP adds LSP semantic token link groups to the highlights map.
func mapLSP(hl map[string]ports.VimHighlight) {
	for _, m := range lspMappings {
		hl[m.group] = ports.VimHighlight{Link: m.link}
	}
}

// mapDiagnostic adds diagnostic highlight groups to the highlights map.
// Includes text, underline, sign, and virtual text variants.
func mapDiagnostic(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	type diagLevel struct {
		suffix    string
		tokenPath string
	}

	levels := []diagLevel{
		{"Error", "status.error"},
		{"Warn", "status.warning"},
		{"Info", "status.info"},
		{"Hint", "status.hint"},
	}

	for _, level := range levels {
		c := colorOf(theme, level.tokenPath)

		// DiagnosticError/Warn/Info/Hint: foreground text color
		hl["Diagnostic"+level.suffix] = ports.VimHighlight{Fg: c}

		// DiagnosticUnderlineError/etc: undercurl with sp color
		hl["DiagnosticUnderline"+level.suffix] = ports.VimHighlight{
			Sp:        c,
			Undercurl: true,
		}

		// DiagnosticSignError/etc: sign column indicators
		hl["DiagnosticSign"+level.suffix] = ports.VimHighlight{Fg: c}

		// DiagnosticVirtualTextError/etc: virtual text
		hl["DiagnosticVirtualText"+level.suffix] = ports.VimHighlight{Fg: c, Italic: true}

		// DiagnosticFloatingError/etc: floating window diagnostics
		hl["DiagnosticFloating"+level.suffix] = ports.VimHighlight{Fg: c}
	}
}

// mapMarkup adds additional markup-related highlight groups beyond the
// @markup.* treesitter groups (which are handled in treesitterMappings).
// This covers legacy markdown and help file groups.
func mapMarkup(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	fg := func(path string) *domain.Color { return colorOf(theme, path) }

	// --- Help file groups ---
	hl["helpCommand"] = ports.VimHighlight{Fg: fg("syntax.string")}
	hl["helpExample"] = ports.VimHighlight{Fg: fg("markup.code")}
	hl["helpHyperTextEntry"] = ports.VimHighlight{Fg: fg("markup.link"), Underline: true}
	hl["helpHyperTextJump"] = ports.VimHighlight{Fg: fg("markup.link"), Underline: true}
	hl["helpSectionDelim"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["helpHeader"] = ports.VimHighlight{Fg: fg("markup.heading"), Bold: true}
}

// terminalTokenOrder defines the ANSI terminal color order (0-15) mapped to
// their semantic token paths.
var terminalTokenOrder = [16]string{
	"terminal.black",     // 0
	"terminal.red",       // 1
	"terminal.green",     // 2
	"terminal.yellow",    // 3
	"terminal.blue",      // 4
	"terminal.magenta",   // 5
	"terminal.cyan",      // 6
	"terminal.white",     // 7
	"terminal.brblack",   // 8
	"terminal.brred",     // 9
	"terminal.brgreen",   // 10
	"terminal.bryellow",  // 11
	"terminal.brblue",    // 12
	"terminal.brmagenta", // 13
	"terminal.brcyan",    // 14
	"terminal.brwhite",   // 15
}

// mapTerminal builds the 16-entry terminal ANSI color array from terminal.*
// semantic tokens.
func mapTerminal(theme *domain.ResolvedTheme) [16]domain.Color {
	var colors [16]domain.Color

	for i, tokenPath := range terminalTokenOrder {
		tok, ok := theme.Tokens.Get(tokenPath)
		if ok && !tok.Color.IsNone {
			colors[i] = tok.Color
		}
	}

	return colors
}

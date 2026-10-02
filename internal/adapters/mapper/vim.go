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
	mapLSP(theme, highlights)
	mapDiagnostic(theme, highlights)
	mapPlugins(theme, highlights)
	mapMarkup(theme, highlights)

	termColors := mapTerminal(theme)
	lualineTheme := mapLualine(theme)
	bufferlineTheme := mapBufferline(theme)

	return &ports.VimTheme{
		Name:           theme.Name,
		Highlights:     highlights,
		TerminalColors: termColors,
		Lualine:        lualineTheme,
		Bufferline:     bufferlineTheme,
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

// noneColor returns a pointer to a Color with IsNone set to true.
// This produces bg = 'none' or fg = 'none' in the output.
func noneColor() *domain.Color {
	return &domain.Color{IsNone: true}
}

// mapBase adds standard Vim editor highlight groups to the highlights map.
//
//nolint:funlen // Large mapping table is intentionally in one function for clarity.
func mapBase(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	fg := func(path string) *domain.Color { return colorOf(theme, path) }
	bg := func(path string) *domain.Color { return colorOf(theme, path) }

	// --- Core editor groups ---
	hl["Normal"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["NormalFloat"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["NormalNC"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["Comment"] = ports.VimHighlight{
		Fg:     fg("syntax.comment"),
		Italic: true,
	}

	hl["Cursor"] = ports.VimHighlight{Reverse: true}
	hl["lCursor"] = ports.VimHighlight{Reverse: true}
	hl["CursorIM"] = ports.VimHighlight{Reverse: true}

	// --- UI groups ---
	hl["CursorLine"] = ports.VimHighlight{
		Bg: bg("surface.background.highlight"),
	}
	hl["CursorColumn"] = ports.VimHighlight{
		Bg: bg("surface.background.highlight"),
	}
	hl["CursorLineNr"] = ports.VimHighlight{
		Fg:   fg("accent.primary"),
		Bold: true,
	}
	hl["ColorColumn"] = ports.VimHighlight{
		Bg: bg("surface.background.highlight"),
	}
	hl["LineNr"] = ports.VimHighlight{
		Fg: fg("text.muted"),
	}
	hl["SignColumn"] = ports.VimHighlight{
		Bg: bg("surface.background.raised"),
	}
	hl["FoldColumn"] = ports.VimHighlight{
		Fg: fg("text.muted"),
		Bg: bg("surface.background.raised"),
	}
	hl["Folded"] = ports.VimHighlight{
		Fg: fg("text.muted"),
		Bg: bg("surface.background.raised"),
	}
	hl["VertSplit"] = ports.VimHighlight{
		Fg: fg("border.default"),
	}
	hl["WinSeparator"] = ports.VimHighlight{
		Fg: fg("border.default"),
	}

	// --- Visual / Selection ---
	hl["Visual"] = ports.VimHighlight{
		Bg: bg("surface.background.selection"),
	}
	hl["VisualNOS"] = ports.VimHighlight{
		Bg: bg("surface.background.selection"),
	}

	// --- Search ---
	hl["Search"] = ports.VimHighlight{
		Fg: fg("text.primary"),
		Bg: bg("surface.background.search"),
	}
	hl["IncSearch"] = ports.VimHighlight{
		Fg: fg("text.inverse"),
		Bg: bg("accent.primary"),
	}
	hl["CurSearch"] = ports.VimHighlight{
		Fg: fg("text.inverse"),
		Bg: bg("accent.primary"),
	}
	hl["Substitute"] = ports.VimHighlight{
		Fg: fg("text.inverse"),
		Bg: bg("status.error"),
	}

	// --- Popup / completion menu ---
	hl["Pmenu"] = ports.VimHighlight{
		Fg: fg("text.primary"),
	}
	hl["PmenuSel"] = ports.VimHighlight{
		Bg: bg("surface.background.selection"),
	}
	hl["PmenuSbar"] = ports.VimHighlight{}
	hl["PmenuThumb"] = ports.VimHighlight{
		Bg: bg("scrollbar.thumb"),
	}
	hl["FloatBorder"] = ports.VimHighlight{
		Fg: fg("syntax.property"),
	}
	hl["FloatTitle"] = ports.VimHighlight{
		Fg: fg("status.hint"),
	}

	// --- Tab line ---
	hl["TabLine"] = ports.VimHighlight{
		Fg: fg("text.muted"),
		Bg: noneColor(),
	}
	hl["TabLineSel"] = ports.VimHighlight{
		Fg:   fg("text.primary"),
		Bg:   bg("surface.background"),
		Bold: true,
	}
	hl["TabLineFill"] = ports.VimHighlight{
		Bg: noneColor(),
	}
	hl["TabLineFile"] = ports.VimHighlight{
		Fg: fg("text.secondary"),
		Bg: noneColor(),
	}

	// --- Status line ---
	hl["StatusLine"] = ports.VimHighlight{
		Fg: fg("text.secondary"),
		Bg: noneColor(),
	}
	hl["StatusLineNC"] = ports.VimHighlight{
		Fg: fg("text.muted"),
		Bg: noneColor(),
	}
	hl["WildMenu"] = ports.VimHighlight{
		Fg: fg("text.inverse"),
		Bg: bg("accent.primary"),
	}

	// --- Messages ---
	hl["ErrorMsg"] = ports.VimHighlight{
		Fg: fg("status.error"),
	}
	hl["WarningMsg"] = ports.VimHighlight{
		Fg: fg("status.warning"),
	}
	hl["ModeMsg"] = ports.VimHighlight{
		Fg:   fg("text.primary"),
		Bold: true,
	}
	hl["MoreMsg"] = ports.VimHighlight{
		Fg: fg("accent.primary"),
	}
	hl["Question"] = ports.VimHighlight{
		Fg: fg("accent.primary"),
	}

	// --- Diff ---
	hl["DiffAdd"] = ports.VimHighlight{
		Bg: bg("diff.added.bg"),
	}
	hl["DiffChange"] = ports.VimHighlight{
		Bg: bg("diff.changed.bg"),
	}
	hl["DiffDelete"] = ports.VimHighlight{
		Fg: fg("diff.deleted.fg"),
		Bg: bg("diff.deleted.bg"),
	}
	hl["DiffText"] = ports.VimHighlight{
		Bg:   bg("diff.changed.bg"),
		Bold: true,
	}

	// --- Spelling ---
	hl["SpellBad"] = ports.VimHighlight{
		Sp:        fg("status.error"),
		Undercurl: true,
	}
	hl["SpellCap"] = ports.VimHighlight{
		Sp:        fg("status.warning"),
		Undercurl: true,
	}
	hl["SpellRare"] = ports.VimHighlight{
		Sp:        fg("accent.secondary"),
		Undercurl: true,
	}
	hl["SpellLocal"] = ports.VimHighlight{
		Sp:        fg("status.info"),
		Undercurl: true,
	}

	// --- Syntax base groups ---
	hl["Constant"] = ports.VimHighlight{
		Fg: fg("syntax.constant"),
	}
	hl["String"] = ports.VimHighlight{
		Fg: fg("syntax.string"),
	}
	hl["Character"] = ports.VimHighlight{
		Fg: fg("syntax.string"),
	}
	hl["Number"] = ports.VimHighlight{
		Fg: fg("syntax.number"),
	}
	hl["Boolean"] = ports.VimHighlight{
		Fg: fg("syntax.boolean"),
	}
	hl["Float"] = ports.VimHighlight{
		Fg: fg("syntax.number"),
	}

	hl["Identifier"] = ports.VimHighlight{
		Fg: fg("syntax.variable"),
	}
	hl["Function"] = ports.VimHighlight{
		Fg: fg("syntax.function"),
	}

	hl["Statement"] = ports.VimHighlight{
		Fg: fg("syntax.keyword"),
	}
	hl["Conditional"] = ports.VimHighlight{
		Fg: fg("syntax.keyword"),
	}
	hl["Repeat"] = ports.VimHighlight{
		Fg: fg("syntax.keyword"),
	}
	hl["Keyword"] = ports.VimHighlight{
		Fg: fg("syntax.keyword"),
	}
	hl["Exception"] = ports.VimHighlight{
		Fg: fg("syntax.keyword"),
	}

	hl["Label"] = ports.VimHighlight{
		Fg: fg("syntax.label"),
	}
	hl["Operator"] = ports.VimHighlight{
		Fg: fg("syntax.operator"),
	}

	// Preprocessor concepts are macros/directives, not ordinary control flow.
	hl["PreProc"] = ports.VimHighlight{
		Fg: fg("syntax.macro"),
	}
	hl["Define"] = ports.VimHighlight{
		Fg: fg("syntax.macro"),
	}
	hl["Macro"] = ports.VimHighlight{
		Fg: fg("syntax.macro"),
	}
	hl["PreCondit"] = ports.VimHighlight{
		Fg: fg("syntax.macro"),
	}

	// Includes and imports describe module relationships.
	hl["Include"] = ports.VimHighlight{
		Fg: fg("syntax.module"),
	}

	// Type-related concepts stay within the type family.
	hl["Type"] = ports.VimHighlight{
		Fg: fg("syntax.type"),
	}
	hl["StorageClass"] = ports.VimHighlight{
		Fg: fg("syntax.type"),
	}
	hl["Structure"] = ports.VimHighlight{
		Fg: fg("syntax.type"),
	}
	hl["Typedef"] = ports.VimHighlight{
		Fg: fg("syntax.type"),
	}

	hl["Special"] = ports.VimHighlight{
		Fg: fg("accent.primary"),
	}
	hl["SpecialChar"] = ports.VimHighlight{
		Fg: fg("syntax.escape"),
	}
	hl["Tag"] = ports.VimHighlight{
		Fg: fg("syntax.tag"),
	}
	hl["Delimiter"] = ports.VimHighlight{
		Fg: fg("syntax.punctuation"),
	}
	hl["SpecialComment"] = ports.VimHighlight{
		Fg:     fg("syntax.comment"),
		Italic: true,
	}
	hl["Debug"] = ports.VimHighlight{
		Fg: fg("status.warning"),
	}
	hl["Underlined"] = ports.VimHighlight{
		Underline: true,
	}
	hl["Ignore"] = ports.VimHighlight{}
	hl["Error"] = ports.VimHighlight{
		Fg: fg("status.error"),
	}
	hl["Todo"] = ports.VimHighlight{
		Fg:   fg("status.todo"),
		Bold: true,
	}

	// --- Miscellaneous ---
	hl["MatchParen"] = ports.VimHighlight{
		Fg:   fg("accent.primary"),
		Bold: true,
	}
	hl["NonText"] = ports.VimHighlight{
		Fg: fg("text.subtle"),
	}
	hl["SpecialKey"] = ports.VimHighlight{
		Fg: fg("text.subtle"),
	}
	hl["Whitespace"] = ports.VimHighlight{
		Fg: fg("text.subtle"),
	}
	hl["Conceal"] = ports.VimHighlight{
		Fg: fg("text.muted"),
	}
	hl["Directory"] = ports.VimHighlight{
		Fg: fg("accent.primary"),
	}
	hl["Title"] = ports.VimHighlight{
		Fg:   fg("accent.primary"),
		Bold: true,
	}
	hl["EndOfBuffer"] = ports.VimHighlight{
		Fg: fg("surface.background"),
	}

	// --- Markup ---
	hl["markdownH1"] = ports.VimHighlight{
		Fg:   fg("markup.heading"),
		Bold: true,
	}
	hl["markdownH2"] = ports.VimHighlight{
		Fg:   fg("markup.heading"),
		Bold: true,
	}
	hl["markdownH3"] = ports.VimHighlight{
		Fg:   fg("markup.heading"),
		Bold: true,
	}
	hl["markdownH4"] = ports.VimHighlight{
		Fg:   fg("markup.heading"),
		Bold: true,
	}
	hl["markdownH5"] = ports.VimHighlight{
		Fg:   fg("markup.heading"),
		Bold: true,
	}
	hl["markdownH6"] = ports.VimHighlight{
		Fg:   fg("markup.heading"),
		Bold: true,
	}
	hl["markdownUrl"] = ports.VimHighlight{
		Fg:        fg("markup.link"),
		Underline: true,
	}
	hl["markdownCode"] = ports.VimHighlight{
		Fg: fg("markup.code"),
	}
	hl["markdownCodeBlock"] = ports.VimHighlight{
		Fg: fg("markup.code"),
	}
	hl["markdownBold"] = ports.VimHighlight{
		Bold: true,
	}
	hl["markdownItalic"] = ports.VimHighlight{
		Italic: true,
	}
	hl["markdownListMarker"] = ports.VimHighlight{
		Fg: fg("markup.list.bullet"),
	}
	hl["markdownBlockquote"] = ports.VimHighlight{
		Fg:     fg("markup.quote"),
		Italic: true,
	}
}

// treesitterMapping maps a treesitter highlight group name to a semantic token
// path. An empty tokenPath means the group uses a Link instead.
type treesitterMapping struct {
	group         string
	tokenPath     string
	bgPath        string
	link          string
	italic        bool
	bold          bool
	strikethrough bool
}

// treesitterMappings defines the treesitter highlight groups and their semantic
// token sources.
//
//nolint:dupl // Intentional structural similarity between mapping tables.
var treesitterMappings = []treesitterMapping{
	// Keywords.
	{group: "@keyword", tokenPath: "syntax.keyword", italic: true},
	{group: "@keyword.function", tokenPath: "syntax.function"},
	{group: "@keyword.operator", tokenPath: "syntax.operator"},
	{group: "@keyword.return", tokenPath: "syntax.keyword"},
	{group: "@keyword.coroutine", tokenPath: "syntax.function.builtin"},
	{group: "@keyword.exception", tokenPath: "syntax.keyword"},
	{group: "@keyword.conditional", tokenPath: "syntax.keyword"},
	{group: "@keyword.repeat", tokenPath: "syntax.keyword"},
	{group: "@keyword.import", tokenPath: "syntax.module"},
	{group: "@keyword.debug", link: "Debug"},
	{group: "@keyword.directive", tokenPath: "syntax.macro"},
	{group: "@keyword.directive.define", tokenPath: "syntax.macro"},
	{group: "@keyword.storage", tokenPath: "syntax.type"},

	// Strings.
	{group: "@string", tokenPath: "syntax.string"},
	{group: "@string.escape", tokenPath: "syntax.escape"},
	{group: "@string.regex", tokenPath: "syntax.regexp"},
	{group: "@string.regexp", tokenPath: "syntax.regexp"},
	{group: "@string.special", tokenPath: "syntax.escape"},
	{group: "@string.documentation", tokenPath: "syntax.string.documentation"},

	{group: "@character", tokenPath: "syntax.string"},
	{group: "@character.special", tokenPath: "syntax.escape"},
	{group: "@character.printf", tokenPath: "syntax.escape"},

	// Functions.
	{group: "@function", tokenPath: "syntax.function"},
	{group: "@function.builtin", tokenPath: "syntax.function.builtin"},
	{group: "@function.call", tokenPath: "syntax.function"},
	{group: "@function.macro", tokenPath: "syntax.macro"},
	{group: "@function.method", tokenPath: "syntax.function"},
	{group: "@function.method.call", tokenPath: "syntax.function"},

	// Legacy Tree-sitter method captures.
	{group: "@method", tokenPath: "syntax.function"},
	{group: "@method.call", tokenPath: "syntax.function"},

	// Constructors.
	{group: "@constructor", tokenPath: "syntax.constructor"},
	{group: "@constructor.tsx", tokenPath: "syntax.constructor"},

	// Variables.
	{group: "@variable", tokenPath: "syntax.variable"},
	{group: "@variable.builtin", tokenPath: "syntax.constant"},
	{group: "@variable.parameter", tokenPath: "syntax.parameter"},
	{group: "@variable.parameter.builtin", tokenPath: "syntax.parameter"},
	{group: "@variable.member", tokenPath: "syntax.property"},

	// Legacy captures.
	{group: "@property", tokenPath: "syntax.property"},
	{group: "@parameter", tokenPath: "syntax.parameter"},

	// Types.
	{group: "@type", tokenPath: "syntax.type"},
	{group: "@type.builtin", tokenPath: "syntax.type.builtin"},
	{group: "@type.definition", tokenPath: "syntax.type"},
	{group: "@type.qualifier", tokenPath: "syntax.type"},

	// Constants and literals.
	{group: "@constant", tokenPath: "syntax.constant"},
	{group: "@constant.builtin", tokenPath: "syntax.constant"},
	{group: "@constant.macro", tokenPath: "syntax.macro"},
	{group: "@number", tokenPath: "syntax.number"},
	{group: "@number.float", tokenPath: "syntax.number"},
	{group: "@boolean", tokenPath: "syntax.boolean"},

	// Operators and punctuation.
	{group: "@operator", tokenPath: "syntax.operator"},
	{group: "@punctuation.bracket", tokenPath: "syntax.punctuation"},
	{group: "@punctuation.delimiter", tokenPath: "syntax.punctuation"},
	{group: "@punctuation.special", tokenPath: "syntax.punctuation"},

	// Tags.
	{group: "@tag", tokenPath: "syntax.tag"},
	{group: "@tag.attribute", tokenPath: "syntax.property"},
	{group: "@tag.delimiter", tokenPath: "syntax.punctuation"},
	{group: "@tag.builtin", tokenPath: "syntax.tag"},
	{group: "@tag.javascript", tokenPath: "syntax.tag"},
	{group: "@tag.tsx", tokenPath: "syntax.tag"},
	{group: "@tag.delimiter.tsx", tokenPath: "syntax.punctuation"},

	// Modules and namespaces.
	{group: "@namespace", tokenPath: "syntax.module"},
	{group: "@namespace.builtin", tokenPath: "syntax.module.builtin"},
	{group: "@module", tokenPath: "syntax.module"},
	{group: "@module.builtin", tokenPath: "syntax.module.builtin"},

	// Misc language constructs.
	{group: "@label", tokenPath: "syntax.label"},
	{group: "@include", tokenPath: "syntax.module"},
	{group: "@exception", tokenPath: "syntax.keyword"},
	{group: "@define", tokenPath: "syntax.macro"},
	{group: "@preproc", tokenPath: "syntax.macro"},
	{group: "@annotation", tokenPath: "syntax.macro"},
	{group: "@attribute", tokenPath: "syntax.macro"},
	{group: "@none", link: "Normal"},

	// Comments.
	{group: "@comment", link: "Comment"},
	{group: "@comment.error", tokenPath: "comment.error"},
	{group: "@comment.warning", tokenPath: "comment.warning"},
	{group: "@comment.info", tokenPath: "comment.info"},
	{group: "@comment.hint", tokenPath: "comment.hint"},
	{group: "@comment.note", tokenPath: "comment.note"},
	{group: "@comment.todo", tokenPath: "comment.todo"},

	// Diff.
	{group: "@diff.plus", link: "DiffAdd"},
	{group: "@diff.minus", link: "DiffDelete"},
	{group: "@diff.delta", link: "DiffChange"},

	// Legacy text captures.
	{group: "@text", link: "Normal"},
	{group: "@text.strong", bold: true},
	{group: "@text.emphasis", italic: true},
	{group: "@text.underline", link: "Underlined"},
	{group: "@text.strike", strikethrough: true},
	{group: "@text.title", link: "Title"},
	{group: "@text.uri", link: "Underlined"},
	{group: "@text.todo", link: "Todo"},
	{group: "@text.note", link: "Todo"},
	{group: "@text.warning", link: "WarningMsg"},
	{group: "@text.danger", link: "ErrorMsg"},

	// Markup.
	{group: "@markup", link: "@none"},
	{group: "@markup.heading", tokenPath: "markup.heading", bold: true},
	{group: "@markup.heading.1.markdown", tokenPath: "markup.heading.1", bold: true},
	{group: "@markup.heading.2.markdown", tokenPath: "markup.heading.2", bold: true},
	{group: "@markup.heading.3.markdown", tokenPath: "markup.heading.3", bold: true},
	{group: "@markup.heading.4.markdown", tokenPath: "markup.heading.4", bold: true},
	{group: "@markup.heading.5.markdown", tokenPath: "markup.heading.5", bold: true},
	{group: "@markup.heading.6.markdown", tokenPath: "markup.heading.6", bold: true},

	{group: "@markup.link", tokenPath: "markup.link"},
	{group: "@markup.link.url", tokenPath: "markup.link"},
	{group: "@markup.link.label", tokenPath: "syntax.escape"},
	{group: "@markup.link.label.symbol", tokenPath: "syntax.variable"},

	{group: "@markup.raw", tokenPath: "markup.code"},
	{group: "@markup.raw.markdown_inline", tokenPath: "markup.code"},

	{group: "@markup.list", tokenPath: "markup.list.bullet"},
	{group: "@markup.list.markdown", tokenPath: "markup.list.bullet", bold: true},
	{group: "@markup.list.checked", tokenPath: "markup.list.checked"},
	{group: "@markup.list.unchecked", tokenPath: "markup.list.unchecked"},

	{group: "@markup.strong", bold: true},
	{group: "@markup.italic", italic: true},
	{group: "@markup.emphasis", italic: true},
	{group: "@markup.underline", link: "Underlined"},
	{group: "@markup.strikethrough", strikethrough: true},
	{group: "@markup.quote", tokenPath: "markup.quote", italic: true},
	{group: "@markup.math", link: "Special"},
	{group: "@markup.environment", tokenPath: "syntax.macro"},
	{group: "@markup.environment.name", tokenPath: "syntax.type"},
}

// mapTreesitter adds treesitter highlight groups to the highlights map.
func mapTreesitter(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	for _, m := range treesitterMappings {
		if m.link != "" {
			hl[m.group] = ports.VimHighlight{Link: m.link}
			continue
		}

		h := ports.VimHighlight{
			Italic:        m.italic,
			Bold:          m.bold,
			Strikethrough: m.strikethrough,
		}

		if m.tokenPath != "" {
			h.Fg = colorOf(theme, m.tokenPath)
		}

		if m.bgPath != "" {
			h.Bg = colorOf(theme, m.bgPath)
		}

		hl[m.group] = h
	}
}

// lspMapping maps an LSP semantic token onto its closest Tree-sitter capture.
// Tree-sitter defines the visual vocabulary; LSP adds semantic precision.
type lspMapping struct {
	group string
	link  string
}

// lspMappings defines LSP semantic-token aliases into the Tree-sitter
// highlight hierarchy.
var lspMappings = []lspMapping{
	// Functions.
	{"@lsp.type.function", "@function"},
	{"@lsp.type.method", "@function.method"},
	{"@lsp.type.macro", "@function.macro"},

	// Variables and members.
	{"@lsp.type.variable", "@variable"},
	{"@lsp.type.parameter", "@variable.parameter"},
	{"@lsp.type.property", "@variable.member"},
	{"@lsp.type.generic", "@variable"},
	{"@lsp.type.selfKeyword", "@variable.builtin"},
	{"@lsp.type.selfTypeKeyword", "@variable.builtin"},

	// Types.
	{"@lsp.type.type", "@type"},
	{"@lsp.type.builtinType", "@type.builtin"},
	{"@lsp.type.typeParameter", "@type"},
	{"@lsp.type.typeAlias", "@type.definition"},
	{"@lsp.type.enum", "@type"},
	{"@lsp.type.struct", "@type"},
	{"@lsp.type.class", "@type"},
	{"@lsp.type.interface", "@type"},

	// Modules.
	{"@lsp.type.namespace", "@module"},
	{"@lsp.type.namespace.python", "@module"},

	// Constants and literals.
	{"@lsp.type.enumMember", "@constant"},
	{"@lsp.type.string", "@string"},
	{"@lsp.type.number", "@number"},
	{"@lsp.type.boolean", "@boolean"},
	{"@lsp.type.regexp", "@string.regex"},
	{"@lsp.type.escapeSequence", "@string.escape"},
	{"@lsp.type.formatSpecifier", "@string.special"},

	// Language constructs.
	{"@lsp.type.keyword", "@keyword"},
	{"@lsp.type.modifier", "@type.qualifier"},
	{"@lsp.type.lifetime", "@keyword.storage"},
	{"@lsp.type.operator", "@operator"},

	// Attributes / annotations.
	{"@lsp.type.decorator", "@attribute"},
	{"@lsp.type.deriveHelper", "@attribute"},

	// Comments.
	{"@lsp.type.comment", "@comment"},

	// No exact Tree-sitter equivalent.
	{"@lsp.type.event", "@type"},

	// Default-library symbols.
	{"@lsp.typemod.class.defaultLibrary", "@type.builtin"},
	{"@lsp.typemod.enum.defaultLibrary", "@type.builtin"},
	{"@lsp.typemod.enumMember.defaultLibrary", "@constant.builtin"},
	{"@lsp.typemod.function.defaultLibrary", "@function.builtin"},
	{"@lsp.typemod.macro.defaultLibrary", "@function.macro"},
	{"@lsp.typemod.method.defaultLibrary", "@function.builtin"},
	{"@lsp.typemod.struct.defaultLibrary", "@type.builtin"},
	{"@lsp.typemod.variable.defaultLibrary", "@variable.builtin"},

	// Modifiers.
	{"@lsp.typemod.keyword.async", "@keyword.coroutine"},
	{"@lsp.typemod.keyword.injected", "@keyword"},
	{"@lsp.typemod.operator.injected", "@operator"},
	{"@lsp.typemod.string.injected", "@string"},
	{"@lsp.typemod.variable.injected", "@variable"},
	{"@lsp.typemod.variable.callable", "@function"},
	{"@lsp.typemod.variable.static", "@constant"},
}

// mapLSP adds LSP semantic-token link groups to the highlights map.
func mapLSP(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	for _, m := range lspMappings {
		hl[m.group] = ports.VimHighlight{Link: m.link}
	}

	// Default-library types are built-in types semantically.
	hl["@lsp.typemod.type.defaultLibrary"] = ports.VimHighlight{
		Link: "@type.builtin",
	}
	hl["@lsp.typemod.typeAlias.defaultLibrary"] = ports.VimHighlight{
		Link: "@type.builtin",
	}

	// Unresolved references are an LSP-only diagnostic state.
	hl["@lsp.type.unresolvedReference"] = ports.VimHighlight{
		Sp:        colorOf(theme, "status.error"),
		Undercurl: true,
	}

	// Deprecation is semantic state rather than syntax.
	hl["@lsp.mod.deprecated"] = ports.VimHighlight{
		Fg:            colorOf(theme, "syntax.deprecated"),
		Strikethrough: true,
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

		hl["Diagnostic"+level.suffix] = ports.VimHighlight{
			Fg: c,
		}

		hl["DiagnosticUnderline"+level.suffix] = ports.VimHighlight{
			Sp:        c,
			Undercurl: true,
		}

		hl["DiagnosticSign"+level.suffix] = ports.VimHighlight{
			Fg: c,
		}

		hl["DiagnosticVirtualText"+level.suffix] = ports.VimHighlight{
			Fg:     c,
			Italic: true,
		}

		hl["DiagnosticFloating"+level.suffix] = ports.VimHighlight{
			Fg: c,
		}
	}
}

// mapMarkup adds additional markup-related highlight groups beyond the
// @markup.* treesitter groups.
func mapMarkup(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	fg := func(path string) *domain.Color { return colorOf(theme, path) }

	hl["helpCommand"] = ports.VimHighlight{
		Fg: fg("syntax.string"),
	}
	hl["helpExample"] = ports.VimHighlight{
		Fg: fg("markup.code"),
	}
	hl["helpHyperTextEntry"] = ports.VimHighlight{
		Fg:        fg("markup.link"),
		Underline: true,
	}
	hl["helpHyperTextJump"] = ports.VimHighlight{
		Fg:        fg("markup.link"),
		Underline: true,
	}
	hl["helpSectionDelim"] = ports.VimHighlight{
		Fg: fg("text.muted"),
	}
	hl["helpHeader"] = ports.VimHighlight{
		Fg:   fg("markup.heading"),
		Bold: true,
	}
}

// terminalTokenOrder defines the ANSI terminal color order (0-15) mapped to
// their semantic token paths.
var terminalTokenOrder = [16]string{
	"terminal.black",
	"terminal.red",
	"terminal.green",
	"terminal.yellow",
	"terminal.blue",
	"terminal.magenta",
	"terminal.cyan",
	"terminal.white",
	"terminal.brblack",
	"terminal.brred",
	"terminal.brgreen",
	"terminal.bryellow",
	"terminal.brblue",
	"terminal.brmagenta",
	"terminal.brcyan",
	"terminal.brwhite",
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

// mapLualine builds a lualine theme from the statusline semantic tokens.
func mapLualine(theme *domain.ResolvedTheme) *ports.LualineTheme {
	fg := func(path string) *domain.Color { return colorOf(theme, path) }
	bg := func(path string) *domain.Color { return colorOf(theme, path) }

	baseMode := ports.LualineMode{
		A: ports.LualineModeColors{
			Fg: fg("statusline.a.fg"),
			Bg: bg("statusline.a.bg"),
		},
		B: ports.LualineModeColors{
			Fg: fg("statusline.b.fg"),
			Bg: bg("statusline.b.bg"),
		},
		C: ports.LualineModeColors{
			Fg: fg("statusline.c.fg"),
			Bg: bg("statusline.c.bg"),
		},
	}

	return &ports.LualineTheme{
		Normal:   baseMode,
		Insert:   baseMode,
		Visual:   baseMode,
		Replace:  baseMode,
		Command:  baseMode,
		Inactive: baseMode,
	}
}

// mapBufferline builds a bufferline theme from semantic tokens.
func mapBufferline(theme *domain.ResolvedTheme) *ports.BufferlineTheme {
	fg := func(path string) *domain.Color { return colorOf(theme, path) }
	bg := func(path string) *domain.Color { return colorOf(theme, path) }

	bgBase := bg("surface.background")
	bgRaised := bg("surface.background.raised")
	bgSelected := bg("statusline.a.bg")

	return &ports.BufferlineTheme{
		Background:                ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgRaised},
		Fill:                      ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgBase},
		BufferSelected:            ports.BufferlineColors{Fg: fg("accent.primary"), Bg: bgSelected, Bold: true},
		BufferVisible:             ports.BufferlineColors{Fg: fg("text.secondary"), Bg: bgRaised},
		CloseButton:               ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgRaised},
		CloseButtonSelected:       ports.BufferlineColors{Fg: fg("accent.primary"), Bg: bgSelected},
		CloseButtonVisible:        ports.BufferlineColors{Fg: fg("text.secondary"), Bg: bgRaised},
		Diagnostic:                ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgRaised},
		DiagnosticSelected:        ports.BufferlineColors{Fg: fg("accent.primary"), Bg: bgSelected},
		DiagnosticVisible:         ports.BufferlineColors{Fg: fg("text.secondary"), Bg: bgRaised},
		Duplicate:                 ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgRaised, Italic: true},
		DuplicateSelected:         ports.BufferlineColors{Fg: fg("accent.primary"), Bg: bgSelected, Italic: true},
		DuplicateVisible:          ports.BufferlineColors{Fg: fg("text.secondary"), Bg: bgRaised, Italic: true},
		Error:                     ports.BufferlineColors{Fg: fg("status.error"), Bg: bgRaised},
		ErrorSelected:             ports.BufferlineColors{Fg: fg("status.error"), Bg: bgSelected},
		ErrorVisible:              ports.BufferlineColors{Fg: fg("status.error"), Bg: bgRaised},
		ErrorDiagnostic:           ports.BufferlineColors{Fg: fg("status.error"), Bg: bgRaised},
		ErrorDiagnosticSelected:   ports.BufferlineColors{Fg: fg("status.error"), Bg: bgSelected},
		ErrorDiagnosticVisible:    ports.BufferlineColors{Fg: fg("status.error"), Bg: bgRaised},
		Hint:                      ports.BufferlineColors{Fg: fg("status.hint"), Bg: bgRaised},
		HintSelected:              ports.BufferlineColors{Fg: fg("status.hint"), Bg: bgSelected},
		HintVisible:               ports.BufferlineColors{Fg: fg("status.hint"), Bg: bgRaised},
		HintDiagnostic:            ports.BufferlineColors{Fg: fg("status.hint"), Bg: bgRaised},
		HintDiagnosticSelected:    ports.BufferlineColors{Fg: fg("status.hint"), Bg: bgSelected},
		HintDiagnosticVisible:     ports.BufferlineColors{Fg: fg("status.hint"), Bg: bgRaised},
		IndicatorSelected:         ports.BufferlineColors{Fg: fg("accent.primary"), Bg: bgSelected},
		IndicatorVisible:          ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgRaised},
		Info:                      ports.BufferlineColors{Fg: fg("status.info"), Bg: bgRaised},
		InfoSelected:              ports.BufferlineColors{Fg: fg("status.info"), Bg: bgSelected},
		InfoVisible:               ports.BufferlineColors{Fg: fg("status.info"), Bg: bgRaised},
		InfoDiagnostic:            ports.BufferlineColors{Fg: fg("status.info"), Bg: bgRaised},
		InfoDiagnosticSelected:    ports.BufferlineColors{Fg: fg("status.info"), Bg: bgSelected},
		InfoDiagnosticVisible:     ports.BufferlineColors{Fg: fg("status.info"), Bg: bgRaised},
		Modified:                  ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgBase},
		ModifiedSelected:          ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgSelected},
		ModifiedVisible:           ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgBase},
		Numbers:                   ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgRaised},
		NumbersSelected:           ports.BufferlineColors{Fg: fg("accent.primary"), Bg: bgSelected},
		NumbersVisible:            ports.BufferlineColors{Fg: fg("text.secondary"), Bg: bgRaised},
		OffsetSeparator:           ports.BufferlineColors{Fg: bgBase, Bg: bgBase},
		Pick:                      ports.BufferlineColors{Fg: fg("accent.secondary"), Bg: bgRaised, Bold: true},
		PickSelected:              ports.BufferlineColors{Fg: fg("accent.secondary"), Bg: bgSelected, Bold: true},
		PickVisible:               ports.BufferlineColors{Fg: fg("accent.secondary"), Bg: bgRaised, Bold: true},
		Separator:                 ports.BufferlineColors{Fg: bgRaised, Bg: bgBase},
		SeparatorSelected:         ports.BufferlineColors{Fg: bgSelected, Bg: bgBase},
		SeparatorVisible:          ports.BufferlineColors{Fg: bgRaised, Bg: bgBase},
		Tab:                       ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgRaised},
		TabClose:                  ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgBase},
		TabSelected:               ports.BufferlineColors{Fg: fg("accent.primary"), Bg: bgSelected, Bold: true},
		TabSeparator:              ports.BufferlineColors{Fg: bgRaised, Bg: bgBase},
		TabSeparatorSelected:      ports.BufferlineColors{Fg: bgSelected, Bg: bgBase},
		TruncMarker:               ports.BufferlineColors{Fg: fg("text.muted"), Bg: bgBase},
		Warning:                   ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgRaised},
		WarningSelected:           ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgSelected},
		WarningVisible:            ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgRaised},
		WarningDiagnostic:         ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgRaised},
		WarningDiagnosticSelected: ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgSelected},
		WarningDiagnosticVisible:  ports.BufferlineColors{Fg: fg("status.warning"), Bg: bgRaised},
	}
}

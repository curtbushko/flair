package mapper_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// TestVimMapper_Interface verifies that the Vim mapper implements
// ports.Mapper and Name() returns "vim".
func TestVimMapper_Interface(t *testing.T) {
	m := mapper.NewVim()

	// Compile-time interface check.
	var _ ports.Mapper = m

	if name := m.Name(); name != "vim" {
		t.Errorf("Name() = %q, want %q", name, "vim")
	}
}

// TestVimMapper_BaseHighlights verifies that the Vim mapper produces
// standard editor highlight groups (Normal, Comment, Visual, CursorLine,
// LineNr, SignColumn, etc.) with correct colors from the ResolvedTheme.
func TestVimMapper_BaseHighlights(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt, ok := result.(*ports.VimTheme)
	if !ok {
		t.Fatalf("Map() returned %T, want *ports.VimTheme", result)
	}

	if vt.Highlights == nil {
		t.Fatal("VimTheme.Highlights is nil")
	}

	// colorPtr is a helper to get a *domain.Color for comparison.
	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	tests := []struct {
		group  string
		wantFg *domain.Color
		wantBg *domain.Color
		italic bool
	}{
		// Normal: fg=text.primary (#c0caf5), no bg (transparency)
		{"Normal", colorPtr("#c0caf5"), nil, false},
		// Comment: fg=syntax.comment (#565f89), italic
		{"Comment", colorPtr("#565f89"), nil, true},
		// Visual: bg=surface.background.selection (blended)
		{"Visual", nil, nil, false},
		// CursorLine: bg=surface.background.highlight (#292e42)
		{"CursorLine", nil, colorPtr("#292e42"), false},
		// LineNr: fg=text.muted (#565f89)
		{"LineNr", colorPtr("#565f89"), nil, false},
		// SignColumn: bg=surface.background.raised (#1f2335)
		{"SignColumn", nil, colorPtr("#1f2335"), false},
	}

	for _, tc := range tests {
		t.Run(tc.group, func(t *testing.T) {
			hl, ok := vt.Highlights[tc.group]
			if !ok {
				t.Fatalf("highlight group %q not found", tc.group)
			}

			if tc.wantFg != nil {
				if hl.Fg == nil {
					t.Errorf("%s: Fg is nil, want %s", tc.group, tc.wantFg.Hex())
				} else if !hl.Fg.Equal(*tc.wantFg) {
					t.Errorf("%s: Fg = %s, want %s", tc.group, hl.Fg.Hex(), tc.wantFg.Hex())
				}
			}

			if tc.wantBg != nil {
				if hl.Bg == nil {
					t.Errorf("%s: Bg is nil, want %s", tc.group, tc.wantBg.Hex())
				} else if !hl.Bg.Equal(*tc.wantBg) {
					t.Errorf("%s: Bg = %s, want %s", tc.group, hl.Bg.Hex(), tc.wantBg.Hex())
				}
			}

			if tc.italic && !hl.Italic {
				t.Errorf("%s: Italic = false, want true", tc.group)
			}
		})
	}

	// Additionally verify these groups exist.
	requiredBase := []string{
		"Normal", "Comment", "Visual", "CursorLine", "LineNr",
		"SignColumn", "StatusLine", "StatusLineNC", "Pmenu",
		"PmenuSel", "PmenuSbar", "PmenuThumb", "TabLine",
		"TabLineSel", "TabLineFill", "TabLineFile", "Search", "IncSearch",
		"MatchParen", "Folded", "FoldColumn", "DiffAdd",
		"DiffChange", "DiffDelete", "DiffText", "ErrorMsg",
		"WarningMsg", "VertSplit", "WinSeparator", "NonText",
		"SpecialKey", "Directory", "Title", "Question",
		"MoreMsg", "ModeMsg", "WildMenu", "Conceal",
		"SpellBad", "SpellCap", "SpellRare", "SpellLocal",
		"ColorColumn", "CursorColumn", "CursorLineNr",
	}

	for _, group := range requiredBase {
		if _, ok := vt.Highlights[group]; !ok {
			t.Errorf("missing required base highlight group %q", group)
		}
	}
}

// TestVimMapper_TabLineFileBgNone verifies that TabLineFile has bg = 'none'.
func TestVimMapper_TabLineFileBgNone(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	hl, ok := vt.Highlights["TabLineFile"]
	if !ok {
		t.Fatal("TabLineFile highlight group not found")
	}

	if hl.Bg == nil {
		t.Fatal("TabLineFile Bg is nil, expected IsNone color")
	}

	if !hl.Bg.IsNone {
		t.Errorf("TabLineFile Bg.IsNone = %v, want true", hl.Bg.IsNone)
	}
}

// TestVimMapper_TreesitterHighlights verifies that the Vim mapper produces
// treesitter highlight groups (@keyword, @string, @function, @variable, etc.)
// with correct colors from semantic tokens.
func TestVimMapper_TreesitterHighlights(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	// colorPtr is a helper to get a *domain.Color for comparison.
	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	// Treesitter groups with expected foreground colors.
	tsTests := []struct {
		group  string
		wantFg *domain.Color
		italic bool
		bold   bool
		link   string
	}{
		{"@keyword", colorPtr("#bb9af7"), false, false, ""},
		{"@string", colorPtr("#9ece6a"), false, false, ""},
		{"@function", colorPtr("#7aa2f7"), false, false, ""},
		{"@comment", nil, false, false, "Comment"},
		{"@variable", colorPtr("#c0caf5"), false, false, ""},
		{"@type", colorPtr("#e0af68"), false, false, ""},
		{"@constant", colorPtr("#ff9e64"), false, false, ""},
		{"@operator", colorPtr("#8db6fa"), false, false, ""},
		{"@number", colorPtr("#ff9e64"), false, false, ""},
		{"@string.escape", colorPtr("#c8acf8"), false, false, ""}, // syntax.escape = base17
		{"@string.regex", colorPtr("#7dcfff"), false, false, ""},
		{"@tag", colorPtr("#f7768e"), false, false, ""},
		{"@property", colorPtr("#97d8f8"), false, false, ""}, // syntax.property = base15
		{"@parameter", colorPtr("#e0af68"), false, false, ""},
		{"@constructor", colorPtr("#c8acf8"), false, false, ""},
		{"@function.builtin", colorPtr("#7dcfff"), false, false, ""}, // syntax.function.builtin = base0C
		{"@type.builtin", colorPtr("#7dcfff"), false, false, ""},     // syntax.type.builtin = base0C
		{"@variable.builtin", colorPtr("#f7768e"), false, false, ""},
		{"@keyword.return", colorPtr("#bb9af7"), false, false, ""},
		{"@keyword.function", colorPtr("#bb9af7"), false, false, ""},
	}

	for _, tc := range tsTests {
		t.Run(tc.group, func(t *testing.T) {
			hl, ok := vt.Highlights[tc.group]
			if !ok {
				t.Fatalf("treesitter highlight group %q not found", tc.group)
			}

			if tc.link != "" {
				if hl.Link != tc.link {
					t.Errorf("%s: Link = %q, want %q", tc.group, hl.Link, tc.link)
				}
				return // Link groups don't need color checks
			}

			if tc.wantFg != nil {
				if hl.Fg == nil {
					t.Errorf("%s: Fg is nil, want %s", tc.group, tc.wantFg.Hex())
				} else if !hl.Fg.Equal(*tc.wantFg) {
					t.Errorf("%s: Fg = %s, want %s", tc.group, hl.Fg.Hex(), tc.wantFg.Hex())
				}
			}

			if tc.italic && !hl.Italic {
				t.Errorf("%s: Italic = false, want true", tc.group)
			}

			if tc.bold && !hl.Bold {
				t.Errorf("%s: Bold = false, want true", tc.group)
			}
		})
	}
}

// TestVimMapper_LSPHighlights verifies that the Vim mapper produces LSP
// semantic token link groups (@lsp.type.function, @lsp.type.variable, etc.)
// that link to their corresponding treesitter groups.
func TestVimMapper_LSPHighlights(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	// LSP semantic token links.
	lspTests := []struct {
		group    string
		wantLink string
	}{
		{"@lsp.type.function", "@function"},
		{"@lsp.type.variable", "@variable"},
		{"@lsp.type.keyword", "@keyword"},
		{"@lsp.type.type", "@type"},
		{"@lsp.type.property", "@property"},
		{"@lsp.type.parameter", "@parameter"},
		{"@lsp.type.method", "@function"},
		{"@lsp.type.string", "@string"},
		{"@lsp.type.number", "@number"},
		{"@lsp.type.operator", "@operator"},
		{"@lsp.type.comment", "@comment"},
		{"@lsp.type.namespace", "@type"},
		{"@lsp.type.enum", "@type"},
		{"@lsp.type.enumMember", "@constant"},
		{"@lsp.type.struct", "@type"},
		{"@lsp.type.class", "@type"},
		{"@lsp.type.interface", "@type"},
		{"@lsp.type.decorator", "@function"},
		{"@lsp.type.macro", "@function.macro"},
	}

	for _, tc := range lspTests {
		t.Run(tc.group, func(t *testing.T) {
			hl, ok := vt.Highlights[tc.group]
			if !ok {
				t.Fatalf("LSP highlight group %q not found", tc.group)
			}
			if hl.Link != tc.wantLink {
				t.Errorf("%s: Link = %q, want %q", tc.group, hl.Link, tc.wantLink)
			}
		})
	}
}

// TestVimMapper_DiagnosticHighlights verifies that the Vim mapper produces
// diagnostic highlight groups (DiagnosticError, DiagnosticWarn, DiagnosticInfo,
// DiagnosticHint) and their underline variants with correct colors.
func TestVimMapper_DiagnosticHighlights(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	// colorPtr is a helper to get a *domain.Color for comparison.
	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	// Diagnostic text groups: fg = status color.
	diagTextTests := []struct {
		group  string
		wantFg *domain.Color
	}{
		{"DiagnosticError", colorPtr("#ff899d")}, // status.error
		{"DiagnosticWarn", colorPtr("#e9c582")},  // status.warning
		{"DiagnosticInfo", colorPtr("#afd67a")},  // status.info (base14)
		{"DiagnosticHint", colorPtr("#ff9e64")},  // status.hint (base09)
	}

	for _, tc := range diagTextTests {
		t.Run(tc.group, func(t *testing.T) {
			hl, ok := vt.Highlights[tc.group]
			if !ok {
				t.Fatalf("diagnostic highlight group %q not found", tc.group)
			}
			if hl.Fg == nil {
				t.Fatalf("%s: Fg is nil, want %s", tc.group, tc.wantFg.Hex())
			}
			if !hl.Fg.Equal(*tc.wantFg) {
				t.Errorf("%s: Fg = %s, want %s", tc.group, hl.Fg.Hex(), tc.wantFg.Hex())
			}
		})
	}

	// Diagnostic underline groups: undercurl + sp = status color.
	diagUnderlineTests := []struct {
		group  string
		wantSp *domain.Color
	}{
		{"DiagnosticUnderlineError", colorPtr("#ff899d")},
		{"DiagnosticUnderlineWarn", colorPtr("#e9c582")},
		{"DiagnosticUnderlineInfo", colorPtr("#afd67a")},
		{"DiagnosticUnderlineHint", colorPtr("#ff9e64")},
	}

	for _, tc := range diagUnderlineTests {
		t.Run(tc.group, func(t *testing.T) {
			hl, ok := vt.Highlights[tc.group]
			if !ok {
				t.Fatalf("diagnostic underline group %q not found", tc.group)
			}
			if !hl.Undercurl {
				t.Errorf("%s: Undercurl = false, want true", tc.group)
			}
			if hl.Sp == nil {
				t.Fatalf("%s: Sp is nil, want %s", tc.group, tc.wantSp.Hex())
			}
			if !hl.Sp.Equal(*tc.wantSp) {
				t.Errorf("%s: Sp = %s, want %s", tc.group, hl.Sp.Hex(), tc.wantSp.Hex())
			}
		})
	}

	// Diagnostic sign groups should exist.
	signGroups := []string{
		"DiagnosticSignError", "DiagnosticSignWarn",
		"DiagnosticSignInfo", "DiagnosticSignHint",
	}
	for _, group := range signGroups {
		if _, ok := vt.Highlights[group]; !ok {
			t.Errorf("missing diagnostic sign group %q", group)
		}
	}

	// DiagnosticVirtualText groups should exist.
	vtGroups := []string{
		"DiagnosticVirtualTextError", "DiagnosticVirtualTextWarn",
		"DiagnosticVirtualTextInfo", "DiagnosticVirtualTextHint",
	}
	for _, group := range vtGroups {
		if _, ok := vt.Highlights[group]; !ok {
			t.Errorf("missing diagnostic virtual text group %q", group)
		}
	}
}

// TestVimMapper_MinimumHighlightCount verifies that the Vim mapper produces
// at least 100 total highlight groups (base + treesitter + LSP + diagnostic).
func TestVimMapper_MinimumHighlightCount(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	if len(vt.Highlights) < 100 {
		t.Errorf("VimTheme.Highlights has %d groups, want >= 100", len(vt.Highlights))
		// Print count breakdown for debugging.
		base, ts, lsp, diag, other := 0, 0, 0, 0, 0
		for name := range vt.Highlights {
			switch {
			case len(name) > 0 && name[0] == '@' && len(name) > 5 && name[:5] == "@lsp.":
				lsp++
			case len(name) > 0 && name[0] == '@':
				ts++
			case len(name) >= 10 && name[:10] == "Diagnostic":
				diag++
			default:
				if isStandardVimGroup(name) {
					base++
				} else {
					other++
				}
			}
		}
		t.Logf("Breakdown: base=%d ts=%d lsp=%d diag=%d other=%d total=%d",
			base, ts, lsp, diag, other, len(vt.Highlights))
	}

	t.Logf("Total Vim highlight groups: %d", len(vt.Highlights))
}

// TestVimMapper_NilTheme verifies that the Vim mapper returns an error
// when given a nil theme.
func TestVimMapper_NilTheme(t *testing.T) {
	m := mapper.NewVim()
	_, err := m.Map(nil)
	if err == nil {
		t.Fatal("expected error for nil theme, got nil")
	}
}

// TestVimMapper_PluginHighlights verifies that the Vim mapper produces
// highlight groups for common Neovim plugins: telescope, gitsigns, nvim-tree,
// indent-blankline, dashboard, lazy, mason, cmp, and notify.
func TestVimMapper_PluginHighlights(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	// Plugin highlight groups that must exist.
	requiredPluginGroups := []string{
		// Telescope
		"TelescopeNormal",
		"TelescopeBorder",
		"TelescopePromptNormal",
		"TelescopePromptBorder",
		"TelescopePromptTitle",
		"TelescopePreviewTitle",
		"TelescopeResultsTitle",
		"TelescopeSelection",
		"TelescopeMatching",

		// GitSigns
		"GitSignsAdd",
		"GitSignsChange",
		"GitSignsDelete",

		// NvimTree
		"NvimTreeNormal",
		"NvimTreeFolderIcon",
		"NvimTreeFolderName",
		"NvimTreeOpenedFolderName",
		"NvimTreeRootFolder",
		"NvimTreeGitDirty",
		"NvimTreeGitNew",
		"NvimTreeGitDeleted",

		// IndentBlankline
		"IndentBlanklineChar",
		"IndentBlanklineContextChar",

		// Dashboard
		"DashboardHeader",
		"DashboardFooter",

		// Lazy
		"LazyButton",
		"LazyButtonActive",

		// Mason
		"MasonHeader",
		"MasonHighlight",

		// Cmp
		"CmpItemAbbr",
		"CmpItemAbbrMatch",
		"CmpItemKind",
		"CmpItemMenu",

		// Notify
		"NotifyERRORBorder",
		"NotifyWARNBorder",
		"NotifyINFOBorder",
	}

	for _, group := range requiredPluginGroups {
		t.Run(group, func(t *testing.T) {
			if _, ok := vt.Highlights[group]; !ok {
				t.Errorf("missing plugin highlight group %q", group)
			}
		})
	}

	// Verify specific colors for key plugin groups.
	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	// GitSignsAdd should use diff.added token color.
	if hl, ok := vt.Highlights["GitSignsAdd"]; ok {
		wantFg := colorPtr("#afd67a") // diff.added.fg = base14
		if hl.Fg == nil {
			t.Errorf("GitSignsAdd: Fg is nil, want %s", wantFg.Hex())
		} else if !hl.Fg.Equal(*wantFg) {
			t.Errorf("GitSignsAdd: Fg = %s, want %s", hl.Fg.Hex(), wantFg.Hex())
		}
	}

	// TelescopeMatching should use accent.primary color.
	if hl, ok := vt.Highlights["TelescopeMatching"]; ok {
		wantFg := colorPtr("#7aa2f7") // accent.primary = base0D
		if hl.Fg == nil {
			t.Errorf("TelescopeMatching: Fg is nil, want %s", wantFg.Hex())
		} else if !hl.Fg.Equal(*wantFg) {
			t.Errorf("TelescopeMatching: Fg = %s, want %s", hl.Fg.Hex(), wantFg.Hex())
		}
	}

	// TelescopeNormal should have main background (surface.background = base00).
	if hl, ok := vt.Highlights["TelescopeNormal"]; ok {
		wantBg := colorPtr("#1a1b26") // surface.background = base00
		if hl.Bg == nil {
			t.Errorf("TelescopeNormal: Bg is nil, want %s", wantBg.Hex())
		} else if !hl.Bg.Equal(*wantBg) {
			t.Errorf("TelescopeNormal: Bg = %s, want %s", hl.Bg.Hex(), wantBg.Hex())
		}
	}

	// TelescopeBorder should match background color for both fg and bg.
	if hl, ok := vt.Highlights["TelescopeBorder"]; ok {
		wantColor := colorPtr("#1a1b26") // surface.background = base00
		if hl.Fg == nil {
			t.Errorf("TelescopeBorder: Fg is nil, want %s", wantColor.Hex())
		} else if !hl.Fg.Equal(*wantColor) {
			t.Errorf("TelescopeBorder: Fg = %s, want %s", hl.Fg.Hex(), wantColor.Hex())
		}
		if hl.Bg == nil {
			t.Errorf("TelescopeBorder: Bg is nil, want %s", wantColor.Hex())
		} else if !hl.Bg.Equal(*wantColor) {
			t.Errorf("TelescopeBorder: Bg = %s, want %s", hl.Bg.Hex(), wantColor.Hex())
		}
	}
}

// TestVimMapper_MarkupHighlights verifies that the Vim mapper produces
// treesitter @markup.* highlight groups with correct colors and styles.
func TestVimMapper_MarkupHighlights(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	markupTests := []struct {
		group  string
		wantFg *domain.Color
		bold   bool
		italic bool
	}{
		{"@markup.heading", colorPtr("#7aa2f7"), true, false}, // markup.heading = base0D, bold
		{"@markup.link", colorPtr("#7dcfff"), false, false},   // markup.link = base0C
		{"@markup.raw", colorPtr("#9ece6a"), false, false},    // markup.code = base0B
		{"@markup.strong", nil, true, false},                  // bold, no specific fg
		{"@markup.italic", nil, false, true},                  // italic, no specific fg
		{"@markup.quote", colorPtr("#565f89"), false, true},   // markup.quote = base03, italic
		{"@markup.list", colorPtr("#ff9e64"), false, false},   // markup.list.bullet = base09
	}

	for _, tc := range markupTests {
		t.Run(tc.group, func(t *testing.T) {
			hl, ok := vt.Highlights[tc.group]
			if !ok {
				t.Fatalf("markup highlight group %q not found", tc.group)
			}

			if tc.wantFg != nil {
				if hl.Fg == nil {
					t.Errorf("%s: Fg is nil, want %s", tc.group, tc.wantFg.Hex())
				} else if !hl.Fg.Equal(*tc.wantFg) {
					t.Errorf("%s: Fg = %s, want %s", tc.group, hl.Fg.Hex(), tc.wantFg.Hex())
				}
			}

			if tc.bold && !hl.Bold {
				t.Errorf("%s: Bold = false, want true", tc.group)
			}

			if tc.italic && !hl.Italic {
				t.Errorf("%s: Italic = false, want true", tc.group)
			}
		})
	}
}

// TestVimMapper_TerminalColors verifies that the Vim mapper produces
// a 16-entry terminal colors array in ANSI order from terminal.* tokens.
func TestVimMapper_TerminalColors(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	// Expected terminal ANSI colors from the plan (Tokyo Night Dark palette).
	// Index: token path → expected hex.
	expectedTerminal := []struct {
		index   int
		name    string
		wantHex string
	}{
		{0, "terminal.black", "#1f2335"},      // base01
		{1, "terminal.red", "#f7768e"},        // base08
		{2, "terminal.green", "#9ece6a"},      // base0B
		{3, "terminal.yellow", "#e0af68"},     // base0A
		{4, "terminal.blue", "#7aa2f7"},       // base0D
		{5, "terminal.magenta", "#bb9af7"},    // base0E
		{6, "terminal.cyan", "#7dcfff"},       // base0C
		{7, "terminal.white", "#c0caf5"},      // base05
		{8, "terminal.brblack", "#565f89"},    // base03
		{9, "terminal.brred", "#ff899d"},      // base12
		{10, "terminal.brgreen", "#afd67a"},   // base14
		{11, "terminal.bryellow", "#e9c582"},  // base13
		{12, "terminal.brblue", "#8db6fa"},    // base16
		{13, "terminal.brmagenta", "#c8acf8"}, // base17
		{14, "terminal.brcyan", "#97d8f8"},    // base15
		{15, "terminal.brwhite", "#c8d3f5"},   // base07
	}

	for _, tc := range expectedTerminal {
		t.Run(tc.name, func(t *testing.T) {
			want := mustParseHex(t, tc.wantHex)
			got := vt.TerminalColors[tc.index]

			if got.IsNone {
				t.Fatalf("TerminalColors[%d] (%s) is IsNone, want %s", tc.index, tc.name, tc.wantHex)
			}

			if !got.Equal(want) {
				t.Errorf("TerminalColors[%d] (%s) = %s, want %s", tc.index, tc.name, got.Hex(), tc.wantHex)
			}
		})
	}
}

// TestVimMapper_TotalHighlightCount verifies that the Vim mapper produces
// at least 200 total highlight groups (base + treesitter + LSP + diagnostic +
// plugin + markup).
func TestVimMapper_TotalHighlightCount(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)

	if len(vt.Highlights) < 200 {
		t.Errorf("VimTheme.Highlights has %d groups, want >= 200", len(vt.Highlights))

		// Count breakdown for debugging.
		base, ts, lsp, diag, plugin, other := 0, 0, 0, 0, 0, 0
		for name := range vt.Highlights {
			switch {
			case len(name) > 5 && name[:5] == "@lsp.":
				lsp++
			case len(name) > 0 && name[0] == '@':
				ts++
			case len(name) >= 10 && name[:10] == "Diagnostic":
				diag++
			case isStandardVimGroup(name):
				base++
			default:
				other++
			}
		}

		// Estimate plugin groups (non-standard capitalized that aren't base vim).
		_ = plugin
		t.Logf("Breakdown: base=%d ts=%d lsp=%d diag=%d other=%d total=%d",
			base, ts, lsp, diag, other, len(vt.Highlights))
	}

	t.Logf("Total Vim highlight groups: %d (want >= 200)", len(vt.Highlights))
}

// isStandardVimGroup checks whether a name looks like a standard Vim highlight
// group (capitalized, no @ prefix, not starting with "Diagnostic").
func isStandardVimGroup(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}

// TestVimMapper_BufferlineTheme verifies that the Vim mapper produces a
// BufferlineTheme with all 15 highlight groups mapped to correct statusline tokens.
func TestVimMapper_BufferlineTheme(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt, ok := result.(*ports.VimTheme)
	if !ok {
		t.Fatalf("Map() returned %T, want *ports.VimTheme", result)
	}

	if vt.Bufferline == nil {
		t.Fatal("VimTheme.Bufferline is nil")
	}

	bl := vt.Bufferline

	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	// Verify BufferSelected uses accent.primary fg with statusline.a.bg and bold
	t.Run("BufferSelected", func(t *testing.T) {
		// From tokenizer: accent.primary = base0D (#7aa2f7), statusline.a.bg = base04 (#a9b1d6)
		wantFg := colorPtr("#7aa2f7")
		wantBg := colorPtr("#a9b1d6")

		if bl.BufferSelected.Fg == nil {
			t.Errorf("BufferSelected.Fg is nil, want %s", wantFg.Hex())
		} else if !bl.BufferSelected.Fg.Equal(*wantFg) {
			t.Errorf("BufferSelected.Fg = %s, want %s", bl.BufferSelected.Fg.Hex(), wantFg.Hex())
		}

		if bl.BufferSelected.Bg == nil {
			t.Errorf("BufferSelected.Bg is nil, want %s", wantBg.Hex())
		} else if !bl.BufferSelected.Bg.Equal(*wantBg) {
			t.Errorf("BufferSelected.Bg = %s, want %s", bl.BufferSelected.Bg.Hex(), wantBg.Hex())
		}

		if !bl.BufferSelected.Bold {
			t.Error("BufferSelected.Bold = false, want true")
		}
	})

	// Verify BufferVisible uses text.secondary with surface.background.raised bg
	// (raised bg provides contrast for triangle separators)
	t.Run("BufferVisible", func(t *testing.T) {
		// From tokenizer: text.secondary = base04 (#a9b1d6), surface.background.raised = base01 (#1f2335)
		wantFg := colorPtr("#a9b1d6")
		wantBg := colorPtr("#1f2335")

		if bl.BufferVisible.Fg == nil {
			t.Errorf("BufferVisible.Fg is nil, want %s", wantFg.Hex())
		} else if !bl.BufferVisible.Fg.Equal(*wantFg) {
			t.Errorf("BufferVisible.Fg = %s, want %s", bl.BufferVisible.Fg.Hex(), wantFg.Hex())
		}

		if bl.BufferVisible.Bg == nil {
			t.Errorf("BufferVisible.Bg is nil, want %s", wantBg.Hex())
		} else if !bl.BufferVisible.Bg.Equal(*wantBg) {
			t.Errorf("BufferVisible.Bg = %s, want %s", bl.BufferVisible.Bg.Hex(), wantBg.Hex())
		}
	})

	// Verify Background uses text.muted with surface.background.raised bg
	// (raised bg provides contrast for triangle separators)
	t.Run("Background", func(t *testing.T) {
		// From tokenizer: text.muted = base03 (#565f89), surface.background.raised = base01 (#1f2335)
		wantFg := colorPtr("#565f89")
		wantBg := colorPtr("#1f2335")

		if bl.Background.Fg == nil {
			t.Errorf("Background.Fg is nil, want %s", wantFg.Hex())
		} else if !bl.Background.Fg.Equal(*wantFg) {
			t.Errorf("Background.Fg = %s, want %s", bl.Background.Fg.Hex(), wantFg.Hex())
		}

		if bl.Background.Bg == nil {
			t.Errorf("Background.Bg is nil, want %s", wantBg.Hex())
		} else if !bl.Background.Bg.Equal(*wantBg) {
			t.Errorf("Background.Bg = %s, want %s", bl.Background.Bg.Hex(), wantBg.Hex())
		}
	})

	// Verify Error uses status.error token
	t.Run("Error", func(t *testing.T) {
		// status.error = base12 (#ff899d)
		wantFg := colorPtr("#ff899d")

		if bl.Error.Fg == nil {
			t.Errorf("Error.Fg is nil, want %s", wantFg.Hex())
		} else if !bl.Error.Fg.Equal(*wantFg) {
			t.Errorf("Error.Fg = %s, want %s", bl.Error.Fg.Hex(), wantFg.Hex())
		}
	})
}

// TestVimMapper_BufferlineSeparators verifies that separator groups use
// surface.background for both fg and bg (invisible separators).
func TestVimMapper_BufferlineSeparators(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)
	if vt.Bufferline == nil {
		t.Fatal("VimTheme.Bufferline is nil")
	}

	bl := vt.Bufferline

	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	// surface.background = base00 (#1a1b26)
	wantSurfaceBg := colorPtr("#1a1b26")
	// surface.background.raised = base01 (#1f2335)
	wantSurfaceRaised := colorPtr("#1f2335")
	// statusline.a.bg = base04 (#a9b1d6)
	wantStatuslineABg := colorPtr("#a9b1d6")

	// Separator: fg = surface.background.raised, bg = surface.background (for triangle)
	t.Run("Separator", func(t *testing.T) {
		if bl.Separator.Fg == nil {
			t.Errorf("Separator.Fg is nil, want %s", wantSurfaceRaised.Hex())
		} else if !bl.Separator.Fg.Equal(*wantSurfaceRaised) {
			t.Errorf("Separator.Fg = %s, want %s", bl.Separator.Fg.Hex(), wantSurfaceRaised.Hex())
		}

		if bl.Separator.Bg == nil {
			t.Errorf("Separator.Bg is nil, want %s", wantSurfaceBg.Hex())
		} else if !bl.Separator.Bg.Equal(*wantSurfaceBg) {
			t.Errorf("Separator.Bg = %s, want %s", bl.Separator.Bg.Hex(), wantSurfaceBg.Hex())
		}
	})

	// SeparatorVisible: fg = surface.background.raised, bg = surface.background
	t.Run("SeparatorVisible", func(t *testing.T) {
		if bl.SeparatorVisible.Fg == nil {
			t.Errorf("SeparatorVisible.Fg is nil, want %s", wantSurfaceRaised.Hex())
		} else if !bl.SeparatorVisible.Fg.Equal(*wantSurfaceRaised) {
			t.Errorf("SeparatorVisible.Fg = %s, want %s", bl.SeparatorVisible.Fg.Hex(), wantSurfaceRaised.Hex())
		}

		if bl.SeparatorVisible.Bg == nil {
			t.Errorf("SeparatorVisible.Bg is nil, want %s", wantSurfaceBg.Hex())
		} else if !bl.SeparatorVisible.Bg.Equal(*wantSurfaceBg) {
			t.Errorf("SeparatorVisible.Bg = %s, want %s", bl.SeparatorVisible.Bg.Hex(), wantSurfaceBg.Hex())
		}
	})

	// SeparatorSelected: fg = statusline.a.bg, bg = surface.background (triangle from selected)
	t.Run("SeparatorSelected", func(t *testing.T) {
		if bl.SeparatorSelected.Fg == nil {
			t.Errorf("SeparatorSelected.Fg is nil, want %s", wantStatuslineABg.Hex())
		} else if !bl.SeparatorSelected.Fg.Equal(*wantStatuslineABg) {
			t.Errorf("SeparatorSelected.Fg = %s, want %s", bl.SeparatorSelected.Fg.Hex(), wantStatuslineABg.Hex())
		}

		if bl.SeparatorSelected.Bg == nil {
			t.Errorf("SeparatorSelected.Bg is nil, want %s", wantSurfaceBg.Hex())
		} else if !bl.SeparatorSelected.Bg.Equal(*wantSurfaceBg) {
			t.Errorf("SeparatorSelected.Bg = %s, want %s", bl.SeparatorSelected.Bg.Hex(), wantSurfaceBg.Hex())
		}
	})
}

// TestVimMapper_BufferlineIndicator verifies that IndicatorSelected uses
// accent.primary for fg with statusline.a.bg for bg.
func TestVimMapper_BufferlineIndicator(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)
	if vt.Bufferline == nil {
		t.Fatal("VimTheme.Bufferline is nil")
	}

	bl := vt.Bufferline

	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	// accent.primary = base0D (#7aa2f7), statusline.a.bg = base04 (#a9b1d6)
	wantFg := colorPtr("#7aa2f7")
	wantBg := colorPtr("#a9b1d6")

	if bl.IndicatorSelected.Fg == nil {
		t.Errorf("IndicatorSelected.Fg is nil, want %s", wantFg.Hex())
	} else if !bl.IndicatorSelected.Fg.Equal(*wantFg) {
		t.Errorf("IndicatorSelected.Fg = %s, want %s", bl.IndicatorSelected.Fg.Hex(), wantFg.Hex())
	}

	if bl.IndicatorSelected.Bg == nil {
		t.Errorf("IndicatorSelected.Bg is nil, want %s", wantBg.Hex())
	} else if !bl.IndicatorSelected.Bg.Equal(*wantBg) {
		t.Errorf("IndicatorSelected.Bg = %s, want %s", bl.IndicatorSelected.Bg.Hex(), wantBg.Hex())
	}
}

// TestVimMapper_BufferlineModified verifies that all Modified states use
// status.warning for fg, with surface.background or statusline.a.bg for bg.
func TestVimMapper_BufferlineModified(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)
	if vt.Bufferline == nil {
		t.Fatal("VimTheme.Bufferline is nil")
	}

	bl := vt.Bufferline

	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	// status.warning = base13 (#e9c582)
	wantWarningFg := colorPtr("#e9c582")

	tests := []struct {
		name   string
		colors ports.BufferlineColors
		wantBg *domain.Color
	}{
		{"Modified", bl.Modified, colorPtr("#1a1b26")},                 // surface.background = base00
		{"ModifiedVisible", bl.ModifiedVisible, colorPtr("#1a1b26")},   // surface.background = base00
		{"ModifiedSelected", bl.ModifiedSelected, colorPtr("#a9b1d6")}, // statusline.a.bg = base04
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.colors.Fg == nil {
				t.Errorf("%s.Fg is nil, want %s", tc.name, wantWarningFg.Hex())
			} else if !tc.colors.Fg.Equal(*wantWarningFg) {
				t.Errorf("%s.Fg = %s, want %s", tc.name, tc.colors.Fg.Hex(), wantWarningFg.Hex())
			}

			if tc.colors.Bg == nil {
				t.Errorf("%s.Bg is nil, want %s", tc.name, tc.wantBg.Hex())
			} else if !tc.colors.Bg.Equal(*tc.wantBg) {
				t.Errorf("%s.Bg = %s, want %s", tc.name, tc.colors.Bg.Hex(), tc.wantBg.Hex())
			}
		})
	}
}

// TestVimMapper_BufferlineDiagnostics verifies that diagnostic colors map to
// status.error/warning/info/hint tokens.
func TestVimMapper_BufferlineDiagnostics(t *testing.T) {
	theme := buildResolvedTheme(t)
	m := mapper.NewVim()

	result, err := m.Map(theme)
	if err != nil {
		t.Fatalf("Map() error: %v", err)
	}

	vt := result.(*ports.VimTheme)
	if vt.Bufferline == nil {
		t.Fatal("VimTheme.Bufferline is nil")
	}

	bl := vt.Bufferline

	colorPtr := func(hex string) *domain.Color {
		t.Helper()
		c := mustParseHex(t, hex)
		return &c
	}

	tests := []struct {
		name   string
		colors ports.BufferlineColors
		wantFg *domain.Color
	}{
		{"Error", bl.Error, colorPtr("#ff899d")},     // status.error = base12
		{"Warning", bl.Warning, colorPtr("#e9c582")}, // status.warning = base13
		{"Info", bl.Info, colorPtr("#afd67a")},       // status.info = base14
		{"Hint", bl.Hint, colorPtr("#ff9e64")},       // status.hint = base09
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.colors.Fg == nil {
				t.Errorf("%s.Fg is nil, want %s", tc.name, tc.wantFg.Hex())
			} else if !tc.colors.Fg.Equal(*tc.wantFg) {
				t.Errorf("%s.Fg = %s, want %s", tc.name, tc.colors.Fg.Hex(), tc.wantFg.Hex())
			}
		})
	}
}

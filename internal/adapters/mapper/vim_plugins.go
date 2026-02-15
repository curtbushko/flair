package mapper

import (
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// mapPlugins adds highlight groups for common Neovim plugins to the highlights
// map. Plugins are organized alphabetically with clear section headers.
//
//nolint:funlen // Large mapping table is intentionally in one function for clarity.
func mapPlugins(theme *domain.ResolvedTheme, hl map[string]ports.VimHighlight) {
	fg := func(path string) *domain.Color { return colorOf(theme, path) }
	bg := func(path string) *domain.Color { return colorOf(theme, path) }

	// --- alpha-nvim / dashboard-nvim ---
	hl["DashboardHeader"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["DashboardCenter"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["DashboardFooter"] = ports.VimHighlight{Fg: fg("text.muted"), Italic: true}
	hl["DashboardShortCut"] = ports.VimHighlight{Fg: fg("accent.primary")}

	// --- bufferline.nvim ---
	hl["BufferLineFill"] = ports.VimHighlight{Bg: bg("surface.background.sunken")}
	hl["BufferLineBackground"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.sunken")}
	hl["BufferLineBufferSelected"] = ports.VimHighlight{Fg: fg("text.primary"), Bold: true}
	hl["BufferLineBufferVisible"] = ports.VimHighlight{Fg: fg("text.muted")}

	// --- nvim-cmp ---
	hl["CmpItemAbbr"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["CmpItemAbbrDeprecated"] = ports.VimHighlight{Fg: fg("text.muted"), Strikethrough: true}
	hl["CmpItemAbbrMatch"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["CmpItemAbbrMatchFuzzy"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["CmpItemKind"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["CmpItemKindClass"] = ports.VimHighlight{Fg: fg("syntax.type")}
	hl["CmpItemKindConstant"] = ports.VimHighlight{Fg: fg("syntax.constant")}
	hl["CmpItemKindFunction"] = ports.VimHighlight{Fg: fg("syntax.function")}
	hl["CmpItemKindKeyword"] = ports.VimHighlight{Fg: fg("syntax.keyword")}
	hl["CmpItemKindMethod"] = ports.VimHighlight{Fg: fg("syntax.function")}
	hl["CmpItemKindProperty"] = ports.VimHighlight{Fg: fg("syntax.property")}
	hl["CmpItemKindSnippet"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["CmpItemKindVariable"] = ports.VimHighlight{Fg: fg("syntax.variable")}
	hl["CmpItemMenu"] = ports.VimHighlight{Fg: fg("text.muted")}

	// --- gitsigns.nvim ---
	hl["GitSignsAdd"] = ports.VimHighlight{Fg: fg("diff.added.fg")}
	hl["GitSignsChange"] = ports.VimHighlight{Fg: fg("diff.changed.fg")}
	hl["GitSignsDelete"] = ports.VimHighlight{Fg: fg("diff.deleted.fg")}
	hl["GitSignsCurrentLineBlame"] = ports.VimHighlight{Fg: fg("text.subtle"), Italic: true}
	hl["GitSignsAddNr"] = ports.VimHighlight{Fg: fg("diff.added.fg")}
	hl["GitSignsChangeNr"] = ports.VimHighlight{Fg: fg("diff.changed.fg")}
	hl["GitSignsDeleteNr"] = ports.VimHighlight{Fg: fg("diff.deleted.fg")}
	hl["GitSignsAddLn"] = ports.VimHighlight{Bg: bg("diff.added.bg")}
	hl["GitSignsChangeLn"] = ports.VimHighlight{Bg: bg("diff.changed.bg")}
	hl["GitSignsDeleteLn"] = ports.VimHighlight{Bg: bg("diff.deleted.bg")}

	// --- indent-blankline.nvim ---
	hl["IndentBlanklineChar"] = ports.VimHighlight{Fg: fg("border.muted")}
	hl["IndentBlanklineContextChar"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["IndentBlanklineContextStart"] = ports.VimHighlight{Underline: true}
	hl["IblIndent"] = ports.VimHighlight{Fg: fg("border.muted")}
	hl["IblScope"] = ports.VimHighlight{Fg: fg("accent.primary")}

	// --- lazy.nvim ---
	hl["LazyButton"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.raised")}
	hl["LazyButtonActive"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary"), Bold: true}
	hl["LazyH1"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary"), Bold: true}
	hl["LazyH2"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["LazyComment"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["LazyDimmed"] = ports.VimHighlight{Fg: fg("text.subtle")}
	hl["LazySpecial"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["LazyProgressDone"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["LazyProgressTodo"] = ports.VimHighlight{Fg: fg("text.subtle")}

	// --- mason.nvim ---
	hl["MasonHeader"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary"), Bold: true}
	hl["MasonHighlight"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["MasonHighlightBlock"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary")}
	hl["MasonHighlightBlockBold"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary"), Bold: true}
	hl["MasonMuted"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["MasonMutedBlock"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.raised")}

	// --- neo-tree.nvim ---
	hl["NeoTreeNormal"] = ports.VimHighlight{Fg: fg("text.sidebar"), Bg: bg("surface.background.sidebar")}
	hl["NeoTreeNormalNC"] = ports.VimHighlight{Fg: fg("text.sidebar"), Bg: bg("surface.background.sidebar")}
	hl["NeoTreeDimText"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["NeoTreeFileName"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["NeoTreeGitAdded"] = ports.VimHighlight{Fg: fg("git.added")}
	hl["NeoTreeGitDeleted"] = ports.VimHighlight{Fg: fg("git.deleted")}
	hl["NeoTreeGitModified"] = ports.VimHighlight{Fg: fg("git.modified")}
	hl["NeoTreeRootName"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["NeoTreeIndentMarker"] = ports.VimHighlight{Fg: fg("border.muted")}

	// --- nvim-notify ---
	hl["NotifyERRORBorder"] = ports.VimHighlight{Fg: fg("status.error")}
	hl["NotifyERRORIcon"] = ports.VimHighlight{Fg: fg("status.error")}
	hl["NotifyERRORTitle"] = ports.VimHighlight{Fg: fg("status.error"), Bold: true}
	hl["NotifyWARNBorder"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["NotifyWARNIcon"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["NotifyWARNTitle"] = ports.VimHighlight{Fg: fg("status.warning"), Bold: true}
	hl["NotifyINFOBorder"] = ports.VimHighlight{Fg: fg("status.info")}
	hl["NotifyINFOIcon"] = ports.VimHighlight{Fg: fg("status.info")}
	hl["NotifyINFOTitle"] = ports.VimHighlight{Fg: fg("status.info"), Bold: true}
	hl["NotifyDEBUGBorder"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["NotifyDEBUGIcon"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["NotifyDEBUGTitle"] = ports.VimHighlight{Fg: fg("text.muted"), Bold: true}
	hl["NotifyTRACEBorder"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["NotifyTRACEIcon"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["NotifyTRACETitle"] = ports.VimHighlight{Fg: fg("accent.secondary"), Bold: true}

	// --- nvim-tree.lua ---
	hl["NvimTreeNormal"] = ports.VimHighlight{Fg: fg("text.sidebar"), Bg: bg("surface.background.sidebar")}
	hl["NvimTreeNormalNC"] = ports.VimHighlight{Fg: fg("text.sidebar"), Bg: bg("surface.background.sidebar")}
	hl["NvimTreeFolderIcon"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["NvimTreeFolderName"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["NvimTreeOpenedFolderName"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["NvimTreeEmptyFolderName"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["NvimTreeRootFolder"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["NvimTreeGitDirty"] = ports.VimHighlight{Fg: fg("git.modified")}
	hl["NvimTreeGitNew"] = ports.VimHighlight{Fg: fg("git.added")}
	hl["NvimTreeGitDeleted"] = ports.VimHighlight{Fg: fg("git.deleted")}
	hl["NvimTreeGitIgnored"] = ports.VimHighlight{Fg: fg("git.ignored")}
	hl["NvimTreeIndentMarker"] = ports.VimHighlight{Fg: fg("border.muted")}
	hl["NvimTreeSpecialFile"] = ports.VimHighlight{Fg: fg("accent.secondary"), Underline: true}
	hl["NvimTreeImageFile"] = ports.VimHighlight{Fg: fg("text.secondary")}
	hl["NvimTreeSymlink"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["NvimTreeWinSeparator"] = ports.VimHighlight{Fg: fg("border.default"), Bg: bg("surface.background.sidebar")}

	// --- telescope.nvim ---
	hl["TelescopeNormal"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.popup")}
	hl["TelescopeBorder"] = ports.VimHighlight{Fg: fg("border.default"), Bg: bg("surface.background.popup")}
	hl["TelescopePromptNormal"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.raised")}
	hl["TelescopePromptBorder"] = ports.VimHighlight{Fg: fg("border.default"), Bg: bg("surface.background.raised")}
	hl["TelescopePromptTitle"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary"), Bold: true}
	hl["TelescopePromptPrefix"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["TelescopePreviewTitle"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("status.success"), Bold: true}
	hl["TelescopePreviewNormal"] = ports.VimHighlight{Bg: bg("surface.background.popup")}
	hl["TelescopePreviewBorder"] = ports.VimHighlight{Fg: fg("border.default"), Bg: bg("surface.background.popup")}
	hl["TelescopeResultsTitle"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.secondary"), Bold: true}
	hl["TelescopeSelection"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.selection")}
	hl["TelescopeSelectionCaret"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["TelescopeMatching"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["TelescopeMultiSelection"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["TelescopeMultiIcon"] = ports.VimHighlight{Fg: fg("accent.primary")}

	// --- trouble.nvim ---
	hl["TroubleNormal"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.sidebar")}
	hl["TroubleText"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["TroubleCount"] = ports.VimHighlight{Fg: fg("accent.secondary"), Bold: true}

	// --- which-key.nvim ---
	hl["WhichKey"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["WhichKeyGroup"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["WhichKeyDesc"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["WhichKeySeperator"] = ports.VimHighlight{Fg: fg("text.subtle")}
	hl["WhichKeySeparator"] = ports.VimHighlight{Fg: fg("text.subtle")}
	hl["WhichKeyFloat"] = ports.VimHighlight{Bg: bg("surface.background.popup")}
	hl["WhichKeyValue"] = ports.VimHighlight{Fg: fg("text.muted")}
}

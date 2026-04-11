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

	// --- lualine.nvim ---
	// All modes use the same statusline.a/b/c tokens
	// Sections: a, b, c (left side) and x, y, z (right side, mirrors c, b, a)
	// Normal mode
	hl["lualine_a_normal"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	hl["lualine_b_normal"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_c_normal"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_x_normal"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_y_normal"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_z_normal"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	// Insert mode
	hl["lualine_a_insert"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	hl["lualine_b_insert"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_c_insert"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_x_insert"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_y_insert"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_z_insert"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	// Visual mode
	hl["lualine_a_visual"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	hl["lualine_b_visual"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_c_visual"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_x_visual"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_y_visual"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_z_visual"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	// Replace mode
	hl["lualine_a_replace"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	hl["lualine_b_replace"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_c_replace"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_x_replace"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_y_replace"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_z_replace"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	// Command mode
	hl["lualine_a_command"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	hl["lualine_b_command"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_c_command"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_x_command"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_y_command"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_z_command"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	// Terminal mode
	hl["lualine_a_terminal"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	hl["lualine_b_terminal"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_c_terminal"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_x_terminal"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_y_terminal"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_z_terminal"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	// Inactive mode
	hl["lualine_a_inactive"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	hl["lualine_b_inactive"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_c_inactive"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_x_inactive"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["lualine_y_inactive"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["lualine_z_inactive"] = ports.VimHighlight{Fg: fg("statusline.a.fg"), Bg: bg("statusline.a.bg")}
	// Transitional highlights (separators between sections)
	hl["lualine_transitional_lualine_a_normal_to_lualine_b_normal"] = ports.VimHighlight{Fg: fg("statusline.a.bg"), Bg: bg("statusline.b.bg")}
	hl["lualine_transitional_lualine_b_normal_to_lualine_c_normal"] = ports.VimHighlight{Fg: fg("statusline.b.bg"), Bg: bg("statusline.c.bg")}

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
	hl["GitSignsAdd"] = ports.VimHighlight{Fg: fg("diff.added.sign")}
	hl["GitSignsChange"] = ports.VimHighlight{Fg: fg("diff.changed.sign")}
	hl["GitSignsDelete"] = ports.VimHighlight{Fg: fg("diff.deleted.sign")}
	hl["GitSignsCurrentLineBlame"] = ports.VimHighlight{Fg: fg("text.subtle"), Italic: true}
	hl["GitSignsAddNr"] = ports.VimHighlight{Fg: fg("diff.added.sign")}
	hl["GitSignsChangeNr"] = ports.VimHighlight{Fg: fg("diff.changed.sign")}
	hl["GitSignsDeleteNr"] = ports.VimHighlight{Fg: fg("diff.deleted.sign")}
	hl["GitSignsAddLn"] = ports.VimHighlight{Bg: bg("diff.added.bg")}
	hl["GitSignsChangeLn"] = ports.VimHighlight{Bg: bg("diff.changed.bg")}
	hl["GitSignsDeleteLn"] = ports.VimHighlight{Bg: bg("diff.deleted.bg")}

	// --- indent-blankline.nvim ---
	hl["IndentBlanklineChar"] = ports.VimHighlight{Fg: fg("text.subtle"), Nocombine: true}
	hl["IndentBlanklineContextChar"] = ports.VimHighlight{Fg: fg("accent.primary"), Nocombine: true}
	hl["IndentBlanklineContextStart"] = ports.VimHighlight{Underline: true}
	hl["IblIndent"] = ports.VimHighlight{Fg: fg("text.subtle"), Nocombine: true}
	hl["IblScope"] = ports.VimHighlight{Fg: fg("accent.primary"), Nocombine: true}
	hl["IblWhitespace"] = ports.VimHighlight{Fg: fg("text.subtle"), Nocombine: true}
	// indent-blankline v3 internal highlights (defined with colors, not links)
	hl["@ibl.indent.char.1"] = ports.VimHighlight{Fg: fg("text.subtle"), Nocombine: true}
	hl["@ibl.whitespace.char.1"] = ports.VimHighlight{Fg: fg("text.subtle"), Nocombine: true}
	hl["@ibl.scope.char.1"] = ports.VimHighlight{Fg: fg("accent.primary"), Nocombine: true}
	hl["@ibl.scope.underline.1"] = ports.VimHighlight{Fg: fg("accent.primary"), Underline: true}

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
	hl["NeoTreeNormal"] = ports.VimHighlight{Fg: fg("text.overlay")}
	hl["NeoTreeNormalNC"] = ports.VimHighlight{Fg: fg("text.overlay")}
	hl["NeoTreeDimText"] = ports.VimHighlight{Fg: fg("text.overlay")}
	hl["NeoTreeFileName"] = ports.VimHighlight{Fg: fg("text.overlay")}
	hl["NeoTreeGitAdded"] = ports.VimHighlight{Fg: fg("git.added")}
	hl["NeoTreeGitDeleted"] = ports.VimHighlight{Fg: fg("git.deleted")}
	hl["NeoTreeGitModified"] = ports.VimHighlight{Fg: fg("status.hint")}
	hl["NeoTreeGitStaged"] = ports.VimHighlight{Fg: fg("diff.added.fg")}
	hl["NeoTreeGitUntracked"] = ports.VimHighlight{Fg: fg("syntax.constructor")}
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
	hl["NvimTreeNormal"] = ports.VimHighlight{Fg: fg("text.overlay")}
	hl["NvimTreeNormalNC"] = ports.VimHighlight{Fg: fg("text.overlay")}
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
	hl["NvimTreeWinSeparator"] = ports.VimHighlight{Fg: fg("border.default")}

	// --- telescope.nvim ---
	hl["TelescopeNormal"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["TelescopeBorder"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["TelescopePromptNormal"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["TelescopePromptBorder"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["TelescopePromptTitle"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["TelescopePromptPrefix"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["TelescopePreviewTitle"] = ports.VimHighlight{Fg: fg("status.success"), Bold: true}
	hl["TelescopePreviewNormal"] = ports.VimHighlight{}
	hl["TelescopePreviewBorder"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["TelescopeResultsTitle"] = ports.VimHighlight{Fg: fg("accent.secondary"), Bold: true}
	hl["TelescopeSelection"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.selection")}
	hl["TelescopeSelectionCaret"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["TelescopeMatching"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["TelescopeMultiSelection"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["TelescopeMultiIcon"] = ports.VimHighlight{Fg: fg("accent.primary")}

	// --- trouble.nvim ---
	hl["TroubleNormal"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["TroubleText"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["TroubleCount"] = ports.VimHighlight{Fg: fg("accent.secondary"), Bold: true}

	// --- which-key.nvim ---
	hl["WhichKey"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["WhichKeyGroup"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["WhichKeyDesc"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["WhichKeySeperator"] = ports.VimHighlight{Fg: fg("text.subtle")}
	hl["WhichKeySeparator"] = ports.VimHighlight{Fg: fg("text.subtle")}
	hl["WhichKeyFloat"] = ports.VimHighlight{}
	hl["WhichKeyValue"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["WhichKeyNormal"] = ports.VimHighlight{}

	// --- aerial.nvim ---
	hl["AerialNormal"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["AerialGuide"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["AerialLine"] = ports.VimHighlight{Link: "LspInlayHint"}
	hl["AerialArrayIcon"] = ports.VimHighlight{Link: "LspKindArray"}
	hl["AerialBooleanIcon"] = ports.VimHighlight{Link: "LspKindBoolean"}
	hl["AerialClassIcon"] = ports.VimHighlight{Link: "LspKindClass"}
	hl["AerialColorIcon"] = ports.VimHighlight{Link: "LspKindColor"}
	hl["AerialConstantIcon"] = ports.VimHighlight{Link: "LspKindConstant"}
	hl["AerialConstructorIcon"] = ports.VimHighlight{Link: "LspKindConstructor"}
	hl["AerialEnumIcon"] = ports.VimHighlight{Link: "LspKindEnum"}
	hl["AerialEnumMemberIcon"] = ports.VimHighlight{Link: "LspKindEnumMember"}
	hl["AerialEventIcon"] = ports.VimHighlight{Link: "LspKindEvent"}
	hl["AerialFieldIcon"] = ports.VimHighlight{Link: "LspKindField"}
	hl["AerialFileIcon"] = ports.VimHighlight{Link: "LspKindFile"}
	hl["AerialFolderIcon"] = ports.VimHighlight{Link: "LspKindFolder"}
	hl["AerialFunctionIcon"] = ports.VimHighlight{Link: "LspKindFunction"}
	hl["AerialInterfaceIcon"] = ports.VimHighlight{Link: "LspKindInterface"}
	hl["AerialKeyIcon"] = ports.VimHighlight{Link: "LspKindKey"}
	hl["AerialKeywordIcon"] = ports.VimHighlight{Link: "LspKindKeyword"}
	hl["AerialMethodIcon"] = ports.VimHighlight{Link: "LspKindMethod"}
	hl["AerialModuleIcon"] = ports.VimHighlight{Link: "LspKindModule"}
	hl["AerialNamespaceIcon"] = ports.VimHighlight{Link: "LspKindNamespace"}
	hl["AerialNullIcon"] = ports.VimHighlight{Link: "LspKindNull"}
	hl["AerialNumberIcon"] = ports.VimHighlight{Link: "LspKindNumber"}
	hl["AerialObjectIcon"] = ports.VimHighlight{Link: "LspKindObject"}
	hl["AerialOperatorIcon"] = ports.VimHighlight{Link: "LspKindOperator"}
	hl["AerialPropertyIcon"] = ports.VimHighlight{Link: "LspKindProperty"}
	hl["AerialReferenceIcon"] = ports.VimHighlight{Link: "LspKindReference"}
	hl["AerialSnippetIcon"] = ports.VimHighlight{Link: "LspKindSnippet"}
	hl["AerialStringIcon"] = ports.VimHighlight{Link: "LspKindString"}
	hl["AerialStructIcon"] = ports.VimHighlight{Link: "LspKindStruct"}
	hl["AerialTextIcon"] = ports.VimHighlight{Link: "LspKindText"}
	hl["AerialTypeParameterIcon"] = ports.VimHighlight{Link: "LspKindTypeParameter"}
	hl["AerialUnitIcon"] = ports.VimHighlight{Link: "LspKindUnit"}
	hl["AerialValueIcon"] = ports.VimHighlight{Link: "LspKindValue"}
	hl["AerialVariableIcon"] = ports.VimHighlight{Link: "LspKindVariable"}

	// --- alpha-nvim ---
	hl["AlphaButtons"] = ports.VimHighlight{Fg: fg("syntax.regexp")}
	hl["AlphaFooter"] = ports.VimHighlight{Fg: fg("syntax.function")}
	hl["AlphaHeader"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["AlphaHeaderLabel"] = ports.VimHighlight{Fg: fg("syntax.constant")}
	hl["AlphaShortcut"] = ports.VimHighlight{Fg: fg("syntax.constant")}

	// --- flash.nvim ---
	hl["FlashBackdrop"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["FlashLabel"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary"), Bold: true}

	// --- fzf-lua ---
	hl["FzfLuaBorder"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["FzfLuaCursor"] = ports.VimHighlight{Link: "IncSearch"}
	hl["FzfLuaDirPart"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["FzfLuaFilePart"] = ports.VimHighlight{Link: "FzfLuaFzfNormal"}
	hl["FzfLuaFzfCursorLine"] = ports.VimHighlight{Link: "Visual"}
	hl["FzfLuaFzfNormal"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["FzfLuaFzfPointer"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["FzfLuaFzfSeparator"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["FzfLuaHeaderBind"] = ports.VimHighlight{Link: "@punctuation.special"}
	hl["FzfLuaHeaderText"] = ports.VimHighlight{Link: "Title"}
	hl["FzfLuaNormal"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["FzfLuaPath"] = ports.VimHighlight{Link: "Directory"}
	hl["FzfLuaPreviewTitle"] = ports.VimHighlight{Fg: fg("status.success"), Bold: true}
	hl["FzfLuaTitle"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}

	// --- git-gutter ---
	hl["GitGutterAdd"] = ports.VimHighlight{Fg: fg("diff.added.fg")}
	hl["GitGutterAddLineNr"] = ports.VimHighlight{Fg: fg("diff.added.fg")}
	hl["GitGutterChange"] = ports.VimHighlight{Fg: fg("diff.changed.fg")}
	hl["GitGutterChangeLineNr"] = ports.VimHighlight{Fg: fg("diff.changed.fg")}
	hl["GitGutterDelete"] = ports.VimHighlight{Fg: fg("diff.deleted.fg")}
	hl["GitGutterDeleteLineNr"] = ports.VimHighlight{Fg: fg("diff.deleted.fg")}

	// --- grug-far.nvim ---
	hl["GrugFarHelpHeader"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["GrugFarHelpHeaderKey"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["GrugFarInputLabel"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["GrugFarInputPlaceholder"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["GrugFarResultsChangeIndicator"] = ports.VimHighlight{Fg: fg("diff.changed.fg")}
	hl["GrugFarResultsHeader"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["GrugFarResultsLineColumn"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["GrugFarResultsLineNo"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["GrugFarResultsMatch"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary")}
	hl["GrugFarResultsStats"] = ports.VimHighlight{Fg: fg("accent.primary")}

	// --- hop.nvim ---
	hl["HopNextKey"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["HopUnmatched"] = ports.VimHighlight{Fg: fg("text.muted")}

	// --- illuminate.vim ---
	hl["IlluminatedWordRead"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["IlluminatedWordText"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["IlluminatedWordWrite"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["illuminatedWord"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["illuminatedCurWord"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}

	// --- LspKind (used by many plugins) ---
	hl["LspKindArray"] = ports.VimHighlight{Link: "@punctuation.bracket"}
	hl["LspKindBoolean"] = ports.VimHighlight{Link: "@boolean"}
	hl["LspKindClass"] = ports.VimHighlight{Link: "@type"}
	hl["LspKindColor"] = ports.VimHighlight{Link: "Special"}
	hl["LspKindConstant"] = ports.VimHighlight{Link: "@constant"}
	hl["LspKindConstructor"] = ports.VimHighlight{Link: "@constructor"}
	hl["LspKindEnum"] = ports.VimHighlight{Link: "@lsp.type.enum"}
	hl["LspKindEnumMember"] = ports.VimHighlight{Link: "@lsp.type.enumMember"}
	hl["LspKindEvent"] = ports.VimHighlight{Link: "Special"}
	hl["LspKindField"] = ports.VimHighlight{Link: "@variable.member"}
	hl["LspKindFile"] = ports.VimHighlight{Link: "Normal"}
	hl["LspKindFolder"] = ports.VimHighlight{Link: "Directory"}
	hl["LspKindFunction"] = ports.VimHighlight{Link: "@function"}
	hl["LspKindInterface"] = ports.VimHighlight{Link: "@lsp.type.interface"}
	hl["LspKindKey"] = ports.VimHighlight{Link: "@variable.member"}
	hl["LspKindKeyword"] = ports.VimHighlight{Link: "@lsp.type.keyword"}
	hl["LspKindMethod"] = ports.VimHighlight{Link: "@function.method"}
	hl["LspKindModule"] = ports.VimHighlight{Link: "@module"}
	hl["LspKindNamespace"] = ports.VimHighlight{Link: "@module"}
	hl["LspKindNull"] = ports.VimHighlight{Link: "@constant.builtin"}
	hl["LspKindNumber"] = ports.VimHighlight{Link: "@number"}
	hl["LspKindObject"] = ports.VimHighlight{Link: "@constant"}
	hl["LspKindOperator"] = ports.VimHighlight{Link: "@operator"}
	hl["LspKindPackage"] = ports.VimHighlight{Link: "@module"}
	hl["LspKindProperty"] = ports.VimHighlight{Link: "@property"}
	hl["LspKindReference"] = ports.VimHighlight{Link: "@markup.link"}
	hl["LspKindSnippet"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["LspKindString"] = ports.VimHighlight{Link: "@string"}
	hl["LspKindStruct"] = ports.VimHighlight{Link: "@lsp.type.struct"}
	hl["LspKindText"] = ports.VimHighlight{Link: "@markup"}
	hl["LspKindTypeParameter"] = ports.VimHighlight{Link: "@lsp.type.typeParameter"}
	hl["LspKindUnit"] = ports.VimHighlight{Link: "@lsp.type.struct"}
	hl["LspKindValue"] = ports.VimHighlight{Link: "@string"}
	hl["LspKindVariable"] = ports.VimHighlight{Link: "@variable"}

	// --- lsp-saga ---
	hl["LspSagaCodeActionBorder"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["LspSagaCodeActionContent"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["LspSagaCodeActionTitle"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}

	// --- mini.nvim ---
	hl["MiniClueNextKey"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["MiniClueNextKeyWithPostkeys"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["MiniFilesBorder"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["MiniFilesBorderModified"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["MiniIndentscopeSymbol"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["MiniIndentscopeSymbolOff"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["MiniPickBorder"] = ports.VimHighlight{Link: "FloatBorder"}
	hl["MiniPickPreviewLine"] = ports.VimHighlight{Link: "CursorLine"}
	hl["MiniPickPreviewRegion"] = ports.VimHighlight{Link: "IncSearch"}
	hl["MiniStarterCurrent"] = ports.VimHighlight{Link: "CursorLine"}
	hl["MiniStarterHeader"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["MiniStarterFooter"] = ports.VimHighlight{Fg: fg("text.muted"), Italic: true}
	hl["MiniStarterItem"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["MiniStarterItemBullet"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["MiniStarterItemPrefix"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["MiniStarterQuery"] = ports.VimHighlight{Fg: fg("status.info")}
	hl["MiniStarterSection"] = ports.VimHighlight{Fg: fg("accent.secondary")}
	hl["MiniStatuslineDevinfo"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.highlight")}
	hl["MiniStatuslineFileinfo"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.highlight")}
	hl["MiniStatuslineFilename"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.statusbar")}
	hl["MiniStatuslineInactive"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.statusbar")}
	hl["MiniStatuslineModeCommand"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("status.warning"), Bold: true}
	hl["MiniStatuslineModeInsert"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("status.success"), Bold: true}
	hl["MiniStatuslineModeNormal"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary"), Bold: true}
	hl["MiniStatuslineModeOther"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.secondary"), Bold: true}
	hl["MiniStatuslineModeReplace"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("status.error"), Bold: true}
	hl["MiniStatuslineModeVisual"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.secondary"), Bold: true}
	hl["MiniTablineCurrent"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background"), Bold: true}
	hl["MiniTablineFill"] = ports.VimHighlight{Bg: bg("surface.background.sunken")}
	hl["MiniTablineHidden"] = ports.VimHighlight{Fg: fg("text.muted"), Bg: bg("surface.background.sunken")}
	hl["MiniTablineModifiedCurrent"] = ports.VimHighlight{Fg: fg("status.warning"), Bg: bg("surface.background"), Bold: true}
	hl["MiniTablineModifiedHidden"] = ports.VimHighlight{Fg: fg("status.warning"), Bg: bg("surface.background.sunken")}
	hl["MiniTablineModifiedVisible"] = ports.VimHighlight{Fg: fg("status.warning"), Bg: bg("surface.background.statusbar")}
	hl["MiniTablineTabpagesection"] = ports.VimHighlight{Fg: fg("text.inverse"), Bg: bg("accent.primary")}
	hl["MiniTablineVisible"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.statusbar")}

	// --- noice.nvim ---
	hl["NoiceCmdlineIcon"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["NoiceCmdlineIconSearch"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["NoiceCmdlinePopup"] = ports.VimHighlight{}
	hl["NoiceCmdlinePopupBorder"] = ports.VimHighlight{Fg: fg("border.default")}
	hl["NoiceCmdlinePopupBorderSearch"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["NoiceCmdlinePopupTitle"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["NoiceCompletionItemKindClass"] = ports.VimHighlight{Link: "LspKindClass"}
	hl["NoiceCompletionItemKindColor"] = ports.VimHighlight{Link: "LspKindColor"}
	hl["NoiceCompletionItemKindConstant"] = ports.VimHighlight{Link: "LspKindConstant"}
	hl["NoiceCompletionItemKindConstructor"] = ports.VimHighlight{Link: "LspKindConstructor"}
	hl["NoiceCompletionItemKindDefault"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["NoiceCompletionItemKindEnum"] = ports.VimHighlight{Link: "LspKindEnum"}
	hl["NoiceCompletionItemKindEnumMember"] = ports.VimHighlight{Link: "LspKindEnumMember"}
	hl["NoiceCompletionItemKindEvent"] = ports.VimHighlight{Link: "LspKindEvent"}
	hl["NoiceCompletionItemKindField"] = ports.VimHighlight{Link: "LspKindField"}
	hl["NoiceCompletionItemKindFile"] = ports.VimHighlight{Link: "LspKindFile"}
	hl["NoiceCompletionItemKindFolder"] = ports.VimHighlight{Link: "LspKindFolder"}
	hl["NoiceCompletionItemKindFunction"] = ports.VimHighlight{Link: "LspKindFunction"}
	hl["NoiceCompletionItemKindInterface"] = ports.VimHighlight{Link: "LspKindInterface"}
	hl["NoiceCompletionItemKindKeyword"] = ports.VimHighlight{Link: "LspKindKeyword"}
	hl["NoiceCompletionItemKindMethod"] = ports.VimHighlight{Link: "LspKindMethod"}
	hl["NoiceCompletionItemKindModule"] = ports.VimHighlight{Link: "LspKindModule"}
	hl["NoiceCompletionItemKindOperator"] = ports.VimHighlight{Link: "LspKindOperator"}
	hl["NoiceCompletionItemKindProperty"] = ports.VimHighlight{Link: "LspKindProperty"}
	hl["NoiceCompletionItemKindReference"] = ports.VimHighlight{Link: "LspKindReference"}
	hl["NoiceCompletionItemKindSnippet"] = ports.VimHighlight{Link: "LspKindSnippet"}
	hl["NoiceCompletionItemKindStruct"] = ports.VimHighlight{Link: "LspKindStruct"}
	hl["NoiceCompletionItemKindText"] = ports.VimHighlight{Link: "LspKindText"}
	hl["NoiceCompletionItemKindTypeParameter"] = ports.VimHighlight{Link: "LspKindTypeParameter"}
	hl["NoiceCompletionItemKindUnit"] = ports.VimHighlight{Link: "LspKindUnit"}
	hl["NoiceCompletionItemKindValue"] = ports.VimHighlight{Link: "LspKindValue"}
	hl["NoiceCompletionItemKindVariable"] = ports.VimHighlight{Link: "LspKindVariable"}

	// --- rainbow-delimiters ---
	hl["RainbowDelimiterBlue"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["RainbowDelimiterCyan"] = ports.VimHighlight{Fg: fg("syntax.regexp")}
	hl["RainbowDelimiterGreen"] = ports.VimHighlight{Fg: fg("syntax.string")}
	hl["RainbowDelimiterOrange"] = ports.VimHighlight{Fg: fg("syntax.constant")}
	hl["RainbowDelimiterRed"] = ports.VimHighlight{Fg: fg("syntax.tag")}
	hl["RainbowDelimiterViolet"] = ports.VimHighlight{Fg: fg("syntax.constructor")}
	hl["RainbowDelimiterYellow"] = ports.VimHighlight{Fg: fg("syntax.type")}

	// --- render-markdown.nvim ---
	hl["RenderMarkdownH1"] = ports.VimHighlight{Fg: fg("markup.heading.1"), Bold: true}
	hl["RenderMarkdownH2"] = ports.VimHighlight{Fg: fg("markup.heading.2"), Bold: true}
	hl["RenderMarkdownH3"] = ports.VimHighlight{Fg: fg("markup.heading.3"), Bold: true}
	hl["RenderMarkdownH4"] = ports.VimHighlight{Fg: fg("markup.heading.4"), Bold: true}
	hl["RenderMarkdownH5"] = ports.VimHighlight{Fg: fg("markup.heading.5"), Bold: true}
	hl["RenderMarkdownH6"] = ports.VimHighlight{Fg: fg("markup.heading.6"), Bold: true}
	hl["RenderMarkdownH1Bg"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["RenderMarkdownH2Bg"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["RenderMarkdownH3Bg"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["RenderMarkdownH4Bg"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["RenderMarkdownH5Bg"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["RenderMarkdownH6Bg"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["RenderMarkdownCode"] = ports.VimHighlight{Bg: bg("surface.background.raised")}
	hl["RenderMarkdownCodeInline"] = ports.VimHighlight{Bg: bg("surface.background.raised")}
	hl["RenderMarkdownTableHead"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["RenderMarkdownTableRow"] = ports.VimHighlight{Fg: fg("text.primary")}

	// --- scrollbar.nvim ---
	hl["ScrollbarHandle"] = ports.VimHighlight{Bg: bg("scrollbar.thumb")}
	hl["ScrollbarSearchHandle"] = ports.VimHighlight{Fg: fg("status.warning"), Bg: bg("scrollbar.thumb")}
	hl["ScrollbarSearch"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["ScrollbarErrorHandle"] = ports.VimHighlight{Fg: fg("status.error"), Bg: bg("scrollbar.thumb")}
	hl["ScrollbarError"] = ports.VimHighlight{Fg: fg("status.error")}
	hl["ScrollbarWarnHandle"] = ports.VimHighlight{Fg: fg("status.warning"), Bg: bg("scrollbar.thumb")}
	hl["ScrollbarWarn"] = ports.VimHighlight{Fg: fg("status.warning")}
	hl["ScrollbarInfoHandle"] = ports.VimHighlight{Fg: fg("status.info"), Bg: bg("scrollbar.thumb")}
	hl["ScrollbarInfo"] = ports.VimHighlight{Fg: fg("status.info")}
	hl["ScrollbarHintHandle"] = ports.VimHighlight{Fg: fg("status.hint"), Bg: bg("scrollbar.thumb")}
	hl["ScrollbarHint"] = ports.VimHighlight{Fg: fg("status.hint")}

	// --- WinBar ---
	hl["WinBar"] = ports.VimHighlight{Fg: fg("text.primary")}
	hl["WinBarNC"] = ports.VimHighlight{Fg: fg("text.muted")}

	// --- yanky.nvim ---
	hl["YankyPut"] = ports.VimHighlight{Link: "IncSearch"}
	hl["YankyYanked"] = ports.VimHighlight{Link: "IncSearch"}

	// --- Additional CmpItemKind links ---
	hl["CmpItemKindArray"] = ports.VimHighlight{Link: "LspKindArray"}
	hl["CmpItemKindBoolean"] = ports.VimHighlight{Link: "LspKindBoolean"}
	hl["CmpItemKindCodeium"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["CmpItemKindColor"] = ports.VimHighlight{Link: "LspKindColor"}
	hl["CmpItemKindConstructor"] = ports.VimHighlight{Link: "LspKindConstructor"}
	hl["CmpItemKindCopilot"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["CmpItemKindDefault"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["CmpItemKindEnum"] = ports.VimHighlight{Link: "LspKindEnum"}
	hl["CmpItemKindEnumMember"] = ports.VimHighlight{Link: "LspKindEnumMember"}
	hl["CmpItemKindEvent"] = ports.VimHighlight{Link: "LspKindEvent"}
	hl["CmpItemKindField"] = ports.VimHighlight{Link: "LspKindField"}
	hl["CmpItemKindFile"] = ports.VimHighlight{Link: "LspKindFile"}
	hl["CmpItemKindFolder"] = ports.VimHighlight{Link: "LspKindFolder"}
	hl["CmpItemKindInterface"] = ports.VimHighlight{Link: "LspKindInterface"}
	hl["CmpItemKindKey"] = ports.VimHighlight{Link: "LspKindKey"}
	hl["CmpItemKindModule"] = ports.VimHighlight{Link: "LspKindModule"}
	hl["CmpItemKindNamespace"] = ports.VimHighlight{Link: "LspKindNamespace"}
	hl["CmpItemKindNull"] = ports.VimHighlight{Link: "LspKindNull"}
	hl["CmpItemKindNumber"] = ports.VimHighlight{Link: "LspKindNumber"}
	hl["CmpItemKindObject"] = ports.VimHighlight{Link: "LspKindObject"}
	hl["CmpItemKindOperator"] = ports.VimHighlight{Link: "LspKindOperator"}
	hl["CmpItemKindPackage"] = ports.VimHighlight{Link: "LspKindPackage"}
	hl["CmpItemKindReference"] = ports.VimHighlight{Link: "LspKindReference"}
	hl["CmpItemKindString"] = ports.VimHighlight{Link: "LspKindString"}
	hl["CmpItemKindStruct"] = ports.VimHighlight{Link: "LspKindStruct"}
	hl["CmpItemKindSupermaven"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["CmpItemKindTabNine"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["CmpItemKindText"] = ports.VimHighlight{Link: "LspKindText"}
	hl["CmpItemKindTypeParameter"] = ports.VimHighlight{Link: "LspKindTypeParameter"}
	hl["CmpItemKindUnit"] = ports.VimHighlight{Link: "LspKindUnit"}
	hl["CmpItemKindValue"] = ports.VimHighlight{Link: "LspKindValue"}
	hl["CmpDocumentation"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.popup")}
	hl["CmpDocumentationBorder"] = ports.VimHighlight{Fg: fg("border.default"), Bg: bg("surface.background.popup")}
	hl["CmpGhostText"] = ports.VimHighlight{Fg: fg("text.muted")}

	// --- Neogit ---
	hl["NeogitHunkHeader"] = ports.VimHighlight{Fg: fg("text.primary"), Bg: bg("surface.background.highlight")}
	hl["NeogitHunkHeaderHighlight"] = ports.VimHighlight{Fg: fg("accent.primary"), Bg: bg("surface.background.selection")}
	hl["NeogitDiffContext"] = ports.VimHighlight{Bg: bg("surface.background")}
	hl["NeogitDiffContextHighlight"] = ports.VimHighlight{Bg: bg("surface.background.highlight")}
	hl["NeogitDiffAdd"] = ports.VimHighlight{Fg: fg("diff.added.fg"), Bg: bg("diff.added.bg")}
	hl["NeogitDiffAddHighlight"] = ports.VimHighlight{Fg: fg("diff.added.fg"), Bg: bg("diff.added.bg")}
	hl["NeogitDiffDelete"] = ports.VimHighlight{Fg: fg("diff.deleted.fg"), Bg: bg("diff.deleted.bg")}
	hl["NeogitDiffDeleteHighlight"] = ports.VimHighlight{Fg: fg("diff.deleted.fg"), Bg: bg("diff.deleted.bg")}

	// --- navic ---
	hl["NavicIconsArray"] = ports.VimHighlight{Link: "LspKindArray"}
	hl["NavicIconsBoolean"] = ports.VimHighlight{Link: "LspKindBoolean"}
	hl["NavicIconsClass"] = ports.VimHighlight{Link: "LspKindClass"}
	hl["NavicIconsConstant"] = ports.VimHighlight{Link: "LspKindConstant"}
	hl["NavicIconsConstructor"] = ports.VimHighlight{Link: "LspKindConstructor"}
	hl["NavicIconsEnum"] = ports.VimHighlight{Link: "LspKindEnum"}
	hl["NavicIconsEnumMember"] = ports.VimHighlight{Link: "LspKindEnumMember"}
	hl["NavicIconsEvent"] = ports.VimHighlight{Link: "LspKindEvent"}
	hl["NavicIconsField"] = ports.VimHighlight{Link: "LspKindField"}
	hl["NavicIconsFile"] = ports.VimHighlight{Link: "LspKindFile"}
	hl["NavicIconsFunction"] = ports.VimHighlight{Link: "LspKindFunction"}
	hl["NavicIconsInterface"] = ports.VimHighlight{Link: "LspKindInterface"}
	hl["NavicIconsKey"] = ports.VimHighlight{Link: "LspKindKey"}
	hl["NavicIconsKeyword"] = ports.VimHighlight{Link: "LspKindKeyword"}
	hl["NavicIconsMethod"] = ports.VimHighlight{Link: "LspKindMethod"}
	hl["NavicIconsModule"] = ports.VimHighlight{Link: "LspKindModule"}
	hl["NavicIconsNamespace"] = ports.VimHighlight{Link: "LspKindNamespace"}
	hl["NavicIconsNull"] = ports.VimHighlight{Link: "LspKindNull"}
	hl["NavicIconsNumber"] = ports.VimHighlight{Link: "LspKindNumber"}
	hl["NavicIconsObject"] = ports.VimHighlight{Link: "LspKindObject"}
	hl["NavicIconsOperator"] = ports.VimHighlight{Link: "LspKindOperator"}
	hl["NavicIconsPackage"] = ports.VimHighlight{Link: "LspKindPackage"}
	hl["NavicIconsProperty"] = ports.VimHighlight{Link: "LspKindProperty"}
	hl["NavicIconsString"] = ports.VimHighlight{Link: "LspKindString"}
	hl["NavicIconsStruct"] = ports.VimHighlight{Link: "LspKindStruct"}
	hl["NavicIconsTypeParameter"] = ports.VimHighlight{Link: "LspKindTypeParameter"}
	hl["NavicIconsVariable"] = ports.VimHighlight{Link: "LspKindVariable"}
	hl["NavicSeparator"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["NavicText"] = ports.VimHighlight{Fg: fg("text.primary")}

	// --- NeoTree additional ---
	hl["NeoTreeTabActive"] = ports.VimHighlight{Fg: fg("accent.primary"), Bold: true}
	hl["NeoTreeTabInactive"] = ports.VimHighlight{Fg: fg("text.muted")}
	hl["NeoTreeTabSeparatorActive"] = ports.VimHighlight{Fg: fg("accent.primary")}
	hl["NeoTreeTabSeparatorInactive"] = ports.VimHighlight{Fg: fg("text.muted")}

	// --- BufferLine additional ---
	// Uses statusline.c for inactive/fill, statusline.b for selected
	// For slope separators: fg draws the slope, bg is behind it
	hl["BufferLineBackground"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: noneColor()}
	hl["BufferLineBuffer"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: noneColor()}
	hl["BufferLineBufferVisible"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["BufferLineBufferSelected"] = ports.VimHighlight{Fg: fg("accent.primary"), Bg: bg("statusline.b.bg"), Bold: true}
	hl["BufferLineTab"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: noneColor()}
	hl["BufferLineTabSelected"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["BufferLineFill"] = ports.VimHighlight{Fg: fg("statusline.c.bg"), Bg: noneColor()}
	hl["BufferLineSeparator"] = ports.VimHighlight{Fg: fg("statusline.c.bg"), Bg: noneColor()}
	hl["BufferLineSeparatorSelected"] = ports.VimHighlight{Fg: fg("statusline.c.bg"), Bg: bg("statusline.b.bg")}
	hl["BufferLineSeparatorVisible"] = ports.VimHighlight{Fg: fg("statusline.c.bg"), Bg: noneColor()}
	hl["BufferLineIndicatorSelected"] = ports.VimHighlight{Fg: fg("statusline.b.bg"), Bg: bg("statusline.b.bg")}
	hl["BufferLineModified"] = ports.VimHighlight{Fg: fg("status.warning"), Bg: bg("statusline.c.bg")}
	hl["BufferLineModifiedSelected"] = ports.VimHighlight{Fg: fg("statusline.b.bg"), Bg: bg("statusline.b.bg")}
	hl["BufferLineModifiedVisible"] = ports.VimHighlight{Fg: fg("status.warning"), Bg: bg("statusline.c.bg")}
	hl["BufferLineOffsetSeparator"] = ports.VimHighlight{Fg: fg("statusline.c.bg"), Bg: noneColor()}
	hl["BufferLineDuplicate"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["BufferLineDuplicateSelected"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
	hl["BufferLineDuplicateVisible"] = ports.VimHighlight{Fg: fg("statusline.c.fg"), Bg: bg("statusline.c.bg")}
	hl["BufferLineTabSeparator"] = ports.VimHighlight{Fg: fg("statusline.c.bg"), Bg: noneColor()}
	hl["BufferLineTabSeparatorSelected"] = ports.VimHighlight{Fg: fg("statusline.b.fg"), Bg: bg("statusline.b.bg")}
}

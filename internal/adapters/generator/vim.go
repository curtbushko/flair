package generator

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// Vim implements ports.Generator for the Vim/Neovim target.
// It writes a style.lua file with hi clear, vim.g.colors_name,
// nvim_set_hl calls (sorted alphabetically), and terminal_color_N
// assignments from a VimTheme.
type Vim struct{}

// NewVim returns a new Vim generator.
func NewVim() *Vim {
	return &Vim{}
}

// Name returns the target name for this generator.
func (v *Vim) Name() string {
	return "vim"
}

// DefaultFilename returns the default output filename for Vim.
func (v *Vim) DefaultFilename() string {
	return "style.lua"
}

// Generate writes the VimTheme as a Neovim Lua colorscheme to w. The mapped
// argument must be a *ports.VimTheme; a type assertion failure returns a
// *domain.GenerateError. Output consists of:
//  1. hi clear and colors_name header
//  2. Sorted nvim_set_hl calls for all highlight groups
//  3. Terminal color assignments (vim.g.terminal_color_0..15)
//
//nolint:funlen // Generator output logic is intentionally consolidated.
func (v *Vim) Generate(w io.Writer, mapped ports.MappedTheme) error {
	theme, ok := mapped.(*ports.VimTheme)
	if !ok {
		return &domain.GenerateError{
			Target:  "vim",
			Message: fmt.Sprintf("expected *ports.VimTheme, got %T", mapped),
		}
	}

	if err := writeHeader(w, theme.Name); err != nil {
		return &domain.GenerateError{
			Target:  "vim",
			Message: "failed to write header",
			Cause:   err,
		}
	}

	if err := writeHighlights(w, theme.Highlights); err != nil {
		return &domain.GenerateError{
			Target:  "vim",
			Message: "failed to write highlights",
			Cause:   err,
		}
	}

	if err := writeTerminalColors(w, theme.TerminalColors); err != nil {
		return &domain.GenerateError{
			Target:  "vim",
			Message: "failed to write terminal colors",
			Cause:   err,
		}
	}

	// Write VimEnter autocmd to re-apply plugin highlights after plugins load.
	// This ensures highlights aren't overwritten by plugins like bufferline.
	if err := writePluginHighlightAutocmd(w, theme.Highlights); err != nil {
		return &domain.GenerateError{
			Target:  "vim",
			Message: "failed to write plugin highlight autocmd",
			Cause:   err,
		}
	}

	// Write lualine theme setup
	if theme.Lualine != nil {
		if err := writeLualineSetup(w, theme.Lualine); err != nil {
			return &domain.GenerateError{
				Target:  "vim",
				Message: "failed to write lualine setup",
				Cause:   err,
			}
		}
	}

	// Write bufferline theme setup
	if theme.Bufferline != nil {
		if err := writeBufferlineSetup(w, theme.Bufferline); err != nil {
			return &domain.GenerateError{
				Target:  "vim",
				Message: "failed to write bufferline setup",
				Cause:   err,
			}
		}
	}

	return nil
}

// writeHeader writes the Lua colorscheme header: hi clear and colors_name.
func writeHeader(w io.Writer, name string) error {
	if _, err := fmt.Fprint(w, "vim.cmd('hi clear')\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "vim.g.colors_name = '%s'\n", name); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return err
	}
	return nil
}

// writeHighlights writes sorted nvim_set_hl calls for all highlight groups.
// Link groups use { link = 'Target' } syntax. Regular groups emit fg, bg, sp,
// and style attributes.
func writeHighlights(w io.Writer, highlights map[string]ports.VimHighlight) error {
	names := make([]string, 0, len(highlights))
	for name := range highlights {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		hl := highlights[name]
		line := formatSetHl(name, hl)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}

	return nil
}

// formatSetHl formats a single vim.api.nvim_set_hl call for the given group.
func formatSetHl(name string, hl ports.VimHighlight) string {
	if hl.Link != "" {
		return fmt.Sprintf("vim.api.nvim_set_hl(0, '%s', { link = '%s' })", name, hl.Link)
	}

	var attrs []string

	if hl.Fg != nil {
		if hl.Fg.IsNone {
			attrs = append(attrs, "fg = 'none'")
		} else {
			attrs = append(attrs, fmt.Sprintf("fg = '%s'", hl.Fg.Hex()))
		}
	}
	if hl.Bg != nil {
		if hl.Bg.IsNone {
			attrs = append(attrs, "bg = 'none'")
		} else {
			attrs = append(attrs, fmt.Sprintf("bg = '%s'", hl.Bg.Hex()))
		}
	}
	if hl.Sp != nil {
		if hl.Sp.IsNone {
			attrs = append(attrs, "sp = 'none'")
		} else {
			attrs = append(attrs, fmt.Sprintf("sp = '%s'", hl.Sp.Hex()))
		}
	}
	if hl.Bold {
		attrs = append(attrs, "bold = true")
	}
	if hl.Italic {
		attrs = append(attrs, "italic = true")
	}
	if hl.Underline {
		attrs = append(attrs, "underline = true")
	}
	if hl.Undercurl {
		attrs = append(attrs, "undercurl = true")
	}
	if hl.Strikethrough {
		attrs = append(attrs, "strikethrough = true")
	}
	if hl.Reverse {
		attrs = append(attrs, "reverse = true")
	}
	if hl.Nocombine {
		attrs = append(attrs, "nocombine = true")
	}

	return fmt.Sprintf("vim.api.nvim_set_hl(0, '%s', { %s })", name, strings.Join(attrs, ", "))
}

// writeTerminalColors writes vim.g.terminal_color_N assignments for the
// 16 ANSI terminal colors.
func writeTerminalColors(w io.Writer, colors [16]domain.Color) error {
	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return err
	}

	for i, c := range colors {
		if _, err := fmt.Fprintf(w, "vim.g.terminal_color_%d = '%s'\n", i, c.Hex()); err != nil {
			return err
		}
	}

	return nil
}

// pluginHighlightPrefixes are highlight group prefixes that belong to plugins
// which may overwrite our highlights when they load. These need to be re-applied
// via VimEnter autocmd.
var pluginHighlightPrefixes = []string{
	"BufferLine",
	"lualine",
	"Lualine",
}

// isPluginHighlight returns true if the highlight group name belongs to a plugin.
func isPluginHighlight(name string) bool {
	for _, prefix := range pluginHighlightPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// writePluginHighlightAutocmd writes autocmds that re-apply plugin highlights
// after plugins have loaded. Uses multiple events and delays to ensure highlights
// are applied after plugins like bufferline initialize.
func writePluginHighlightAutocmd(w io.Writer, highlights map[string]ports.VimHighlight) error {
	// Collect plugin highlights
	var pluginNames []string
	for name := range highlights {
		if isPluginHighlight(name) {
			pluginNames = append(pluginNames, name)
		}
	}

	if len(pluginNames) == 0 {
		return nil
	}

	sort.Strings(pluginNames)

	// Write a local function to apply plugin highlights
	if _, err := fmt.Fprint(w, "\n-- Re-apply plugin highlights after plugins load\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "local function apply_plugin_highlights()\n"); err != nil {
		return err
	}

	for _, name := range pluginNames {
		hl := highlights[name]
		line := formatSetHl(name, hl)
		if _, err := fmt.Fprintf(w, "  %s\n", line); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "end\n\n"); err != nil {
		return err
	}

	// Apply immediately (catches early loading)
	if _, err := fmt.Fprint(w, "-- Apply immediately\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "apply_plugin_highlights()\n\n"); err != nil {
		return err
	}

	// Apply on UIEnter with longer delay for lazy-loaded plugins
	if _, err := fmt.Fprint(w, "-- Apply on UIEnter with delay for lazy-loaded plugins\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "vim.api.nvim_create_autocmd('UIEnter', {\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  once = true,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  callback = function()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    apply_plugin_highlights()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    vim.defer_fn(apply_plugin_highlights, 100)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    vim.defer_fn(apply_plugin_highlights, 500)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  end,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "})\n\n"); err != nil {
		return err
	}

	// Re-apply on ColorScheme change
	if _, err := fmt.Fprint(w, "-- Re-apply when colorscheme changes\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "vim.api.nvim_create_autocmd('ColorScheme', {\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  callback = function()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    apply_plugin_highlights()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    vim.defer_fn(apply_plugin_highlights, 50)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  end,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "})\n\n"); err != nil {
		return err
	}

	// Re-apply on User BufferlineRender event (bufferline fires this)
	if _, err := fmt.Fprint(w, "-- Re-apply when bufferline renders\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "vim.api.nvim_create_autocmd('User', {\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  pattern = 'BufferlineRender',\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  callback = apply_plugin_highlights,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "})\n"); err != nil {
		return err
	}

	return nil
}

// writeLualineSetup writes the lualine theme configuration.
// This sets up lualine with the theme colors so it uses the colorscheme's
// statusline tokens instead of its default colors.
func writeLualineSetup(w io.Writer, theme *ports.LualineTheme) error {
	if _, err := fmt.Fprint(w, "\n-- Lualine theme configuration\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "local lualine_theme = {\n"); err != nil {
		return err
	}

	// Write each mode
	modes := []struct {
		name string
		mode ports.LualineMode
	}{
		{"normal", theme.Normal},
		{"insert", theme.Insert},
		{"visual", theme.Visual},
		{"replace", theme.Replace},
		{"command", theme.Command},
		{"inactive", theme.Inactive},
	}

	for _, m := range modes {
		if err := writeLualineMode(w, m.name, m.mode); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "}\n\n"); err != nil {
		return err
	}

	// Register the theme globally so user can reference it
	if _, err := fmt.Fprint(w, "-- Register lualine theme globally\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "_G.flair_lualine_theme = lualine_theme\n\n"); err != nil {
		return err
	}

	// Write the setup call that applies the theme if lualine is available
	if _, err := fmt.Fprint(w, "-- Apply lualine theme\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "local function apply_lualine_theme()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  local ok, lualine = pcall(require, 'lualine')\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  if ok then\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    local cfg = require('lualine').get_config()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    cfg.options.theme = lualine_theme\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    require('lualine').setup(cfg)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  end\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "end\n\n"); err != nil {
		return err
	}

	// Apply immediately if lualine is already loaded
	if _, err := fmt.Fprint(w, "-- Apply immediately if lualine is loaded\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "if package.loaded['lualine'] then\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  apply_lualine_theme()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "end\n\n"); err != nil {
		return err
	}

	// Also apply after UIEnter with delay
	if _, err := fmt.Fprint(w, "-- Apply after UI loads\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "vim.api.nvim_create_autocmd('UIEnter', {\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  once = true,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  callback = function()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    vim.defer_fn(apply_lualine_theme, 100)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  end,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "})\n"); err != nil {
		return err
	}

	return nil
}

// writeLualineMode writes a single lualine mode configuration.
func writeLualineMode(w io.Writer, name string, mode ports.LualineMode) error {
	if _, err := fmt.Fprintf(w, "  %s = {\n", name); err != nil {
		return err
	}

	sections := []struct {
		name   string
		colors ports.LualineModeColors
	}{
		{"a", mode.A},
		{"b", mode.B},
		{"c", mode.C},
	}

	for _, s := range sections {
		if err := writeLualineSection(w, s.name, s.colors); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "  },\n"); err != nil {
		return err
	}

	return nil
}

// writeLualineSection writes a single lualine section (a, b, or c).
func writeLualineSection(w io.Writer, name string, colors ports.LualineModeColors) error {
	if _, err := fmt.Fprintf(w, "    %s = { ", name); err != nil {
		return err
	}

	var parts []string
	if colors.Fg != nil {
		parts = append(parts, fmt.Sprintf("fg = '%s'", colors.Fg.Hex()))
	}
	if colors.Bg != nil {
		parts = append(parts, fmt.Sprintf("bg = '%s'", colors.Bg.Hex()))
	}

	if _, err := fmt.Fprint(w, strings.Join(parts, ", ")); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, " },\n"); err != nil {
		return err
	}

	return nil
}

// writeBufferlineSetup writes the bufferline.setup() configuration with all highlights.
func writeBufferlineSetup(w io.Writer, theme *ports.BufferlineTheme) error {
	if _, err := fmt.Fprint(w, "\n-- Bufferline theme configuration\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "local bufferline_highlights = {\n"); err != nil {
		return err
	}

	// Helper to format a highlight group
	writeGroup := func(name string, colors ports.BufferlineColors) error {
		if _, err := fmt.Fprintf(w, "  %s = { ", name); err != nil {
			return err
		}

		var parts []string
		if colors.Bg != nil {
			parts = append(parts, fmt.Sprintf("bg = '%s'", colors.Bg.Hex()))
		}
		if colors.Fg != nil {
			parts = append(parts, fmt.Sprintf("fg = '%s'", colors.Fg.Hex()))
		}
		if colors.Bold {
			parts = append(parts, "bold = true")
		}
		if colors.Italic {
			parts = append(parts, "italic = true")
		}

		if _, err := fmt.Fprint(w, strings.Join(parts, ", ")); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, " },\n"); err != nil {
			return err
		}
		return nil
	}

	// Write all highlight groups in alphabetical order
	groups := []struct {
		name   string
		colors ports.BufferlineColors
	}{
		{"background", theme.Background},
		{"buffer_selected", theme.BufferSelected},
		{"buffer_visible", theme.BufferVisible},
		{"close_button", theme.CloseButton},
		{"close_button_selected", theme.CloseButtonSelected},
		{"close_button_visible", theme.CloseButtonVisible},
		{"diagnostic", theme.Diagnostic},
		{"diagnostic_selected", theme.DiagnosticSelected},
		{"diagnostic_visible", theme.DiagnosticVisible},
		{"duplicate", theme.Duplicate},
		{"duplicate_selected", theme.DuplicateSelected},
		{"duplicate_visible", theme.DuplicateVisible},
		{"error", theme.Error},
		{"error_diagnostic", theme.ErrorDiagnostic},
		{"error_diagnostic_selected", theme.ErrorDiagnosticSelected},
		{"error_diagnostic_visible", theme.ErrorDiagnosticVisible},
		{"error_selected", theme.ErrorSelected},
		{"error_visible", theme.ErrorVisible},
		{"fill", theme.Fill},
		{"hint", theme.Hint},
		{"hint_diagnostic", theme.HintDiagnostic},
		{"hint_diagnostic_selected", theme.HintDiagnosticSelected},
		{"hint_diagnostic_visible", theme.HintDiagnosticVisible},
		{"hint_selected", theme.HintSelected},
		{"hint_visible", theme.HintVisible},
		{"indicator_selected", theme.IndicatorSelected},
		{"indicator_visible", theme.IndicatorVisible},
		{"info", theme.Info},
		{"info_diagnostic", theme.InfoDiagnostic},
		{"info_diagnostic_selected", theme.InfoDiagnosticSelected},
		{"info_diagnostic_visible", theme.InfoDiagnosticVisible},
		{"info_selected", theme.InfoSelected},
		{"info_visible", theme.InfoVisible},
		{"modified", theme.Modified},
		{"modified_selected", theme.ModifiedSelected},
		{"modified_visible", theme.ModifiedVisible},
		{"numbers", theme.Numbers},
		{"numbers_selected", theme.NumbersSelected},
		{"numbers_visible", theme.NumbersVisible},
		{"offset_separator", theme.OffsetSeparator},
		{"pick", theme.Pick},
		{"pick_selected", theme.PickSelected},
		{"pick_visible", theme.PickVisible},
		{"separator", theme.Separator},
		{"separator_selected", theme.SeparatorSelected},
		{"separator_visible", theme.SeparatorVisible},
		{"tab", theme.Tab},
		{"tab_close", theme.TabClose},
		{"tab_selected", theme.TabSelected},
		{"tab_separator", theme.TabSeparator},
		{"tab_separator_selected", theme.TabSeparatorSelected},
		{"trunc_marker", theme.TruncMarker},
		{"warning", theme.Warning},
		{"warning_diagnostic", theme.WarningDiagnostic},
		{"warning_diagnostic_selected", theme.WarningDiagnosticSelected},
		{"warning_diagnostic_visible", theme.WarningDiagnosticVisible},
		{"warning_selected", theme.WarningSelected},
		{"warning_visible", theme.WarningVisible},
	}

	for _, g := range groups {
		if err := writeGroup(g.name, g.colors); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "}\n\n"); err != nil {
		return err
	}

	// Register globally for user access
	if _, err := fmt.Fprint(w, "-- Register bufferline highlights globally\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "_G.flair_bufferline_highlights = bufferline_highlights\n\n"); err != nil {
		return err
	}

	// Apply highlights by merging with existing bufferline config
	if _, err := fmt.Fprint(w, "-- Apply bufferline highlights (merges with your existing config)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "local function apply_bufferline_highlights()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  local ok, bufferline = pcall(require, 'bufferline')\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  if not ok then return end\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  local state = require('bufferline.state')\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  if state and state.current_element_hl then\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    -- Bufferline already configured, merge highlights\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    local cfg = bufferline.get_config and bufferline.get_config() or {}\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    cfg.highlights = vim.tbl_deep_extend('force', cfg.highlights or {}, bufferline_highlights)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    bufferline.setup(cfg)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  end\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "end\n\n"); err != nil {
		return err
	}

	// Apply after a delay to let user's bufferline config load first
	if _, err := fmt.Fprint(w, "vim.api.nvim_create_autocmd('User', {\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  pattern = 'LazyDone',\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  callback = function()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    vim.defer_fn(apply_bufferline_highlights, 50)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  end,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "})\n\n"); err != nil {
		return err
	}

	// Also try on VimEnter for non-lazy setups
	if _, err := fmt.Fprint(w, "vim.api.nvim_create_autocmd('VimEnter', {\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  callback = function()\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "    vim.defer_fn(apply_bufferline_highlights, 100)\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  end,\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "})\n"); err != nil {
		return err
	}

	return nil
}

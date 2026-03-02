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

	if theme.Lualine != nil {
		if err := writeLualineTheme(w, theme.Lualine); err != nil {
			return &domain.GenerateError{
				Target:  "vim",
				Message: "failed to write lualine theme",
				Cause:   err,
			}
		}
	}

	if theme.Bufferline != nil {
		if err := writeBufferlineTheme(w, theme.Bufferline); err != nil {
			return &domain.GenerateError{
				Target:  "vim",
				Message: "failed to write bufferline theme",
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
		attrs = append(attrs, fmt.Sprintf("fg = '%s'", hl.Fg.Hex()))
	}
	if hl.Bg != nil {
		attrs = append(attrs, fmt.Sprintf("bg = '%s'", hl.Bg.Hex()))
	}
	if hl.Sp != nil {
		attrs = append(attrs, fmt.Sprintf("sp = '%s'", hl.Sp.Hex()))
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

// writeLualineTheme writes the lualine theme setup to the Lua file.
// It defines a local theme table and calls require("lualine").setup().
func writeLualineTheme(w io.Writer, theme *ports.LualineTheme) error {
	if _, err := fmt.Fprint(w, "\n-- Lualine theme\n"); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, "local lualine_theme = {\n"); err != nil {
		return err
	}

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

	// Write the lualine setup call
	if _, err := fmt.Fprint(w, "-- Apply lualine theme if lualine is available\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "local ok, lualine = pcall(require, 'lualine')\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "if ok then\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  lualine.setup({ options = { theme = lualine_theme } })\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "end\n"); err != nil {
		return err
	}

	return nil
}

// writeLualineMode writes a single lualine mode (normal, insert, etc.) to the output.
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

// writeLualineSection writes a single section (a, b, or c) of a lualine mode.
func writeLualineSection(w io.Writer, name string, colors ports.LualineModeColors) error {
	var parts []string

	if colors.Fg != nil {
		parts = append(parts, fmt.Sprintf("fg = '%s'", colors.Fg.Hex()))
	}
	if colors.Bg != nil {
		parts = append(parts, fmt.Sprintf("bg = '%s'", colors.Bg.Hex()))
	}

	if _, err := fmt.Fprintf(w, "    %s = { %s },\n", name, strings.Join(parts, ", ")); err != nil {
		return err
	}

	return nil
}

// writeBufferlineTheme writes the bufferline theme setup to the Lua file.
// It defines a local bufferline_theme table and applies it via pcall to
// protect against missing bufferline plugin.
func writeBufferlineTheme(w io.Writer, theme *ports.BufferlineTheme) error {
	if _, err := fmt.Fprint(w, "\n-- Bufferline theme\n"); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, "local bufferline_theme = {\n"); err != nil {
		return err
	}

	// Define groups in a fixed order for deterministic output.
	groups := []struct {
		name   string
		colors ports.BufferlineColors
	}{
		{"fill", theme.Fill},
		{"background", theme.Background},
		{"buffer_visible", theme.BufferVisible},
		{"buffer_selected", theme.BufferSelected},
		{"separator", theme.Separator},
		{"separator_visible", theme.SeparatorVisible},
		{"separator_selected", theme.SeparatorSelected},
		{"indicator_selected", theme.IndicatorSelected},
		{"modified", theme.Modified},
		{"modified_visible", theme.ModifiedVisible},
		{"modified_selected", theme.ModifiedSelected},
		{"error", theme.Error},
		{"warning", theme.Warning},
		{"info", theme.Info},
		{"hint", theme.Hint},
	}

	for _, g := range groups {
		if err := writeBufferlineGroup(w, g.name, g.colors); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "}\n\n"); err != nil {
		return err
	}

	// Write pcall wrapper to apply bufferline theme if plugin is available.
	if _, err := fmt.Fprint(w, "-- Apply bufferline theme if bufferline is available\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "local bl_ok, bufferline = pcall(require, 'bufferline')\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "if bl_ok then\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "  bufferline.setup({ highlights = bufferline_theme })\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "end\n"); err != nil {
		return err
	}

	return nil
}

// writeBufferlineGroup writes a single bufferline highlight group to the output.
func writeBufferlineGroup(w io.Writer, name string, colors ports.BufferlineColors) error {
	var parts []string

	if colors.Fg != nil {
		parts = append(parts, fmt.Sprintf("fg = '%s'", colors.Fg.Hex()))
	}
	if colors.Bg != nil {
		parts = append(parts, fmt.Sprintf("bg = '%s'", colors.Bg.Hex()))
	}
	if colors.Bold {
		parts = append(parts, "bold = true")
	}
	if colors.Italic {
		parts = append(parts, "italic = true")
	}

	if _, err := fmt.Fprintf(w, "  %s = { %s },\n", name, strings.Join(parts, ", ")); err != nil {
		return err
	}

	return nil
}

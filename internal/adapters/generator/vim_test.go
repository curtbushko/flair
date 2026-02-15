package generator_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/generator"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// sampleVimTheme returns a VimTheme with a few highlights and terminal colors.
func sampleVimTheme() *ports.VimTheme {
	fgColor := domain.Color{R: 0xc0, G: 0xca, B: 0xf5} // #c0caf5
	bgColor := domain.Color{R: 0x1a, G: 0x1b, B: 0x26} // #1a1b26

	return &ports.VimTheme{
		Name: "tokyonight",
		Highlights: map[string]ports.VimHighlight{
			"Normal": {Fg: &fgColor, Bg: &bgColor},
		},
		TerminalColors: [16]domain.Color{
			{R: 0x1a, G: 0x1b, B: 0x26}, // 0
			{R: 0xf7, G: 0x76, B: 0x8e}, // 1
			{R: 0x9e, G: 0xce, B: 0x6a}, // 2
			{R: 0xe0, G: 0xaf, B: 0x68}, // 3
			{R: 0x7a, G: 0xa2, B: 0xf7}, // 4
			{R: 0xbb, G: 0x9a, B: 0xf7}, // 5
			{R: 0x7d, G: 0xcf, B: 0xff}, // 6
			{R: 0xa9, G: 0xb1, B: 0xd6}, // 7
			{R: 0x41, G: 0x48, B: 0x68}, // 8
			{R: 0xf7, G: 0x76, B: 0x8e}, // 9
			{R: 0x9e, G: 0xce, B: 0x6a}, // 10
			{R: 0xe0, G: 0xaf, B: 0x68}, // 11
			{R: 0x7a, G: 0xa2, B: 0xf7}, // 12
			{R: 0xbb, G: 0x9a, B: 0xf7}, // 13
			{R: 0x7d, G: 0xcf, B: 0xff}, // 14
			{R: 0xc0, G: 0xca, B: 0xf5}, // 15
		},
	}
}

// TestVimGenerator_Interface verifies that the Vim generator implements
// ports.Generator and returns the expected Name() and DefaultFilename().
func TestVimGenerator_Interface(t *testing.T) {
	g := generator.NewVim()

	// Compile-time interface check.
	var _ ports.Generator = g

	if name := g.Name(); name != "vim" {
		t.Errorf("Name() = %q, want %q", name, "vim")
	}

	if filename := g.DefaultFilename(); filename != "style.lua" {
		t.Errorf("DefaultFilename() = %q, want %q", filename, "style.lua")
	}
}

// TestVimGenerator_HiClear verifies that Generate produces output starting
// with vim.cmd('hi clear') and vim.g.colors_name assignment.
func TestVimGenerator_HiClear(t *testing.T) {
	g := generator.NewVim()
	theme := sampleVimTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "vim.cmd('hi clear')") {
		t.Error("output does not contain vim.cmd('hi clear')")
	}

	expected := fmt.Sprintf("vim.g.colors_name = '%s'", theme.Name)
	if !strings.Contains(output, expected) {
		t.Errorf("output does not contain %q", expected)
	}

	// hi clear should appear before colors_name.
	hiIdx := strings.Index(output, "vim.cmd('hi clear')")
	nameIdx := strings.Index(output, "vim.g.colors_name")
	if hiIdx >= nameIdx {
		t.Error("vim.cmd('hi clear') should appear before vim.g.colors_name")
	}
}

// TestVimGenerator_SetHl verifies that Generate produces nvim_set_hl calls
// with correct fg and bg values.
func TestVimGenerator_SetHl(t *testing.T) {
	g := generator.NewVim()

	fgColor := domain.Color{R: 0xc0, G: 0xca, B: 0xf5}
	bgColor := domain.Color{R: 0x1a, G: 0x1b, B: 0x26}

	theme := &ports.VimTheme{
		Name: "test",
		Highlights: map[string]ports.VimHighlight{
			"Normal": {Fg: &fgColor, Bg: &bgColor},
		},
		TerminalColors: [16]domain.Color{},
	}

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Must contain nvim_set_hl call for Normal group.
	if !strings.Contains(output, "vim.api.nvim_set_hl(0, 'Normal'") {
		t.Error("output does not contain nvim_set_hl call for Normal")
	}

	// Must contain fg and bg values.
	if !strings.Contains(output, "fg = '#c0caf5'") {
		t.Error("output does not contain fg = '#c0caf5'")
	}
	if !strings.Contains(output, "bg = '#1a1b26'") {
		t.Error("output does not contain bg = '#1a1b26'")
	}
}

// TestVimGenerator_Links verifies that Generate produces link highlight groups
// using { link = 'Target' } syntax.
func TestVimGenerator_Links(t *testing.T) {
	g := generator.NewVim()

	theme := &ports.VimTheme{
		Name: "test",
		Highlights: map[string]ports.VimHighlight{
			"@comment": {Link: "Comment"},
		},
		TerminalColors: [16]domain.Color{},
	}

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "vim.api.nvim_set_hl(0, '@comment', { link = 'Comment' })") {
		t.Errorf("output does not contain expected link call.\noutput:\n%s", output)
	}
}

// TestVimGenerator_TerminalColors verifies that Generate produces
// vim.g.terminal_color_0 through vim.g.terminal_color_15 assignments.
func TestVimGenerator_TerminalColors(t *testing.T) {
	g := generator.NewVim()
	theme := sampleVimTheme()

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	for i := 0; i < 16; i++ {
		prefix := fmt.Sprintf("vim.g.terminal_color_%d", i)
		if !strings.Contains(output, prefix) {
			t.Errorf("output does not contain %q", prefix)
		}
	}

	// Verify a specific terminal color value.
	if !strings.Contains(output, "vim.g.terminal_color_0 = '#1a1b26'") {
		t.Error("terminal_color_0 should be '#1a1b26'")
	}
	if !strings.Contains(output, "vim.g.terminal_color_15 = '#c0caf5'") {
		t.Error("terminal_color_15 should be '#c0caf5'")
	}
}

// TestVimGenerator_StyleAttributes verifies that Generate produces
// bold, italic, and other style attributes in nvim_set_hl calls.
func TestVimGenerator_StyleAttributes(t *testing.T) {
	g := generator.NewVim()

	fgColor := domain.Color{R: 0xbb, G: 0x9a, B: 0xf7}

	theme := &ports.VimTheme{
		Name: "test",
		Highlights: map[string]ports.VimHighlight{
			"Keyword": {Fg: &fgColor, Bold: true, Italic: true},
		},
		TerminalColors: [16]domain.Color{},
	}

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "bold = true") {
		t.Error("output does not contain 'bold = true'")
	}
	if !strings.Contains(output, "italic = true") {
		t.Error("output does not contain 'italic = true'")
	}
}

// TestVimGenerator_WrongType verifies that passing a non-VimTheme value
// as MappedTheme returns a GenerateError.
func TestVimGenerator_WrongType(t *testing.T) {
	g := generator.NewVim()

	var buf bytes.Buffer

	err := g.Generate(&buf, "not a vim theme")
	if err == nil {
		t.Fatal("Generate() with wrong type should return error, got nil")
	}

	var genErr *domain.GenerateError
	if !errors.As(err, &genErr) {
		t.Errorf("error type = %T, want *domain.GenerateError", err)
	}
}

// TestVimGenerator_Deterministic verifies that generating the same VimTheme
// twice produces byte-identical output.
func TestVimGenerator_Deterministic(t *testing.T) {
	g := generator.NewVim()
	theme := sampleVimTheme()

	var buf1, buf2 bytes.Buffer
	if err := g.Generate(&buf1, theme); err != nil {
		t.Fatalf("first Generate() error: %v", err)
	}
	if err := g.Generate(&buf2, theme); err != nil {
		t.Fatalf("second Generate() error: %v", err)
	}

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Error("output is not deterministic; two runs produced different results")
		t.Logf("run 1:\n%s", buf1.String())
		t.Logf("run 2:\n%s", buf2.String())
	}
}

// TestVimGenerator_SortedHighlights verifies that highlight groups are
// emitted in sorted alphabetical order for deterministic output.
func TestVimGenerator_SortedHighlights(t *testing.T) {
	g := generator.NewVim()

	fgColor := domain.Color{R: 0xc0, G: 0xca, B: 0xf5}

	theme := &ports.VimTheme{
		Name: "test",
		Highlights: map[string]ports.VimHighlight{
			"Zebra":   {Fg: &fgColor},
			"Alpha":   {Fg: &fgColor},
			"Middle":  {Fg: &fgColor},
			"Comment": {Fg: &fgColor},
		},
		TerminalColors: [16]domain.Color{},
	}

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	// Extract the order of nvim_set_hl group names.
	lines := strings.Split(output, "\n")
	var groups []string
	for _, line := range lines {
		if strings.Contains(line, "vim.api.nvim_set_hl(0, '") {
			start := strings.Index(line, "nvim_set_hl(0, '") + len("nvim_set_hl(0, '")
			end := strings.Index(line[start:], "'")
			if end > 0 {
				groups = append(groups, line[start:start+end])
			}
		}
	}

	if len(groups) != 4 {
		t.Fatalf("expected 4 highlight groups, got %d: %v", len(groups), groups)
	}

	for i := 1; i < len(groups); i++ {
		if groups[i-1] > groups[i] {
			t.Errorf("highlights not sorted: %q before %q", groups[i-1], groups[i])
		}
	}
}

// TestVimGenerator_SpColor verifies that the sp attribute is included
// in nvim_set_hl calls for highlights with special colors.
func TestVimGenerator_SpColor(t *testing.T) {
	g := generator.NewVim()

	spColor := domain.Color{R: 0xf7, G: 0x76, B: 0x8e}

	theme := &ports.VimTheme{
		Name: "test",
		Highlights: map[string]ports.VimHighlight{
			"SpellBad": {Sp: &spColor, Undercurl: true},
		},
		TerminalColors: [16]domain.Color{},
	}

	var buf bytes.Buffer
	if err := g.Generate(&buf, theme); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "sp = '#f7768e'") {
		t.Error("output does not contain sp = '#f7768e'")
	}
	if !strings.Contains(output, "undercurl = true") {
		t.Error("output does not contain 'undercurl = true'")
	}
}

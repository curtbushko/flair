package flair_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

const (
	// testVariantDark is the expected variant for our test palette.
	testVariantDark = "dark"
)

// tokyoNightPaletteYAML is a complete base24 palette YAML for tokenizer tests.
const tokyoNightPaletteYAML = `system: "base24"
name: "Tokyo Night Dark"
author: "Test Author"
variant: "dark"
palette:
  base00: "1a1b26"
  base01: "1f2335"
  base02: "292e42"
  base03: "565f89"
  base04: "a9b1d6"
  base05: "c0caf5"
  base06: "c0caf5"
  base07: "c8d3f5"
  base08: "f7768e"
  base09: "ff9e64"
  base0A: "e0af68"
  base0B: "9ece6a"
  base0C: "7dcfff"
  base0D: "7aa2f7"
  base0E: "bb9af7"
  base0F: "db4b4b"
  base10: "16161e"
  base11: "101014"
  base12: "ff899d"
  base13: "e9c582"
  base14: "afd67a"
  base15: "97d8f8"
  base16: "8db6fa"
  base17: "c8acf8"
`

// parsePalette is a test helper that parses the Tokyo Night palette.
func parsePalette(t *testing.T) *flair.Palette {
	t.Helper()
	r := strings.NewReader(tokyoNightPaletteYAML)
	pal, err := flair.ParsePalette(r)
	if err != nil {
		t.Fatalf("ParsePalette() unexpected error: %v", err)
	}
	return pal
}

// TestTokenize_SurfaceTokens tests that Tokenize derives surface tokens correctly.
func TestTokenize_SurfaceTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	// Verify theme metadata.
	if got := theme.Name(); got != "Tokyo Night Dark" {
		t.Errorf("Name() = %q, want %q", got, "Tokyo Night Dark")
	}
	if got := theme.Variant(); got != testVariantDark {
		t.Errorf("Variant() = %q, want %q", got, testVariantDark)
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// surface.background = base00 = "1a1b26"
		{"surface.background", "#1a1b26"},
		// surface.background.raised = base01 = "1f2335"
		{"surface.background.raised", "#1f2335"},
		// surface.background.sunken = base10 = "16161e"
		{"surface.background.sunken", "#16161e"},
		// surface.background.darkest = base11 = "101014"
		{"surface.background.darkest", "#101014"},
		// surface.background.highlight = base02 = "292e42"
		{"surface.background.highlight", "#292e42"},
		// surface.background.overlay = base10 = "16161e"
		{"surface.background.overlay", "#16161e"},
		// surface.background.popup = base10 = "16161e"
		{"surface.background.popup", "#16161e"},
		// surface.background.sidebar = base10 = "16161e"
		{"surface.background.sidebar", "#16161e"},
		// surface.background.statusbar = base10 = "16161e"
		{"surface.background.statusbar", "#16161e"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}

	// Test blended colors separately (selection, search).
	// surface.background.selection = BlendBg(base0D, base00, 0.30)
	// base0D = "7aa2f7", base00 = "1a1b26"
	t.Run("surface.background.selection", func(t *testing.T) {
		c, ok := theme.Color("surface.background.selection")
		if !ok {
			t.Fatal("Color(surface.background.selection) not found")
		}
		// Just verify it exists and is non-zero; exact blend value tested in BlendBg test.
		if c.R == 0 && c.G == 0 && c.B == 0 {
			t.Error("Color(surface.background.selection) is zero, expected blended color")
		}
	})

	t.Run("surface.background.search", func(t *testing.T) {
		c, ok := theme.Color("surface.background.search")
		if !ok {
			t.Fatal("Color(surface.background.search) not found")
		}
		if c.R == 0 && c.G == 0 && c.B == 0 {
			t.Error("Color(surface.background.search) is zero, expected blended color")
		}
	})
}

// TestTokenize_TextTokens tests that Tokenize derives text tokens correctly.
func TestTokenize_TextTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// text.primary = base05 = "c0caf5"
		{"text.primary", "#c0caf5"},
		// text.secondary = base04 = "a9b1d6"
		{"text.secondary", "#a9b1d6"},
		// text.muted = base03 = "565f89"
		{"text.muted", "#565f89"},
		// text.inverse = base00 = "1a1b26"
		{"text.inverse", "#1a1b26"},
		// text.overlay = base06 = "c0caf5"
		{"text.overlay", "#c0caf5"},
		// text.sidebar = base04 = "a9b1d6"
		{"text.sidebar", "#a9b1d6"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}

	// text.subtle = BlendBg(base03, base00, 0.50) - verify exists.
	t.Run("text.subtle", func(t *testing.T) {
		c, ok := theme.Color("text.subtle")
		if !ok {
			t.Fatal("Color(text.subtle) not found")
		}
		if c.R == 0 && c.G == 0 && c.B == 0 {
			t.Error("Color(text.subtle) is zero, expected blended color")
		}
	})
}

// TestTokenize_SyntaxTokens tests that Tokenize derives syntax tokens correctly.
func TestTokenize_SyntaxTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// syntax.keyword = base0E = "bb9af7"
		{"syntax.keyword", "#bb9af7"},
		// syntax.string = base0B = "9ece6a"
		{"syntax.string", "#9ece6a"},
		// syntax.function = base0D = "7aa2f7"
		{"syntax.function", "#7aa2f7"},
		// syntax.comment = base03 = "565f89"
		{"syntax.comment", "#565f89"},
		// syntax.variable = base05 = "c0caf5"
		{"syntax.variable", "#c0caf5"},
		// syntax.constant = base09 = "ff9e64"
		{"syntax.constant", "#ff9e64"},
		// syntax.operator = base16 = "8db6fa"
		{"syntax.operator", "#8db6fa"},
		// syntax.type = base0A = "e0af68"
		{"syntax.type", "#e0af68"},
		// syntax.number = base09 = "ff9e64"
		{"syntax.number", "#ff9e64"},
		// syntax.tag = base08 = "f7768e"
		{"syntax.tag", "#f7768e"},
		// syntax.property = base14 = "afd67a"
		{"syntax.property", "#afd67a"},
		// syntax.parameter = base13 = "e9c582"
		{"syntax.parameter", "#e9c582"},
		// syntax.regexp = base0C = "7dcfff"
		{"syntax.regexp", "#7dcfff"},
		// syntax.escape = base0E = "bb9af7"
		{"syntax.escape", "#bb9af7"},
		// syntax.constructor = base17 = "c8acf8"
		{"syntax.constructor", "#c8acf8"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}
}

// TestTokenize_StatusTokens tests that Tokenize derives status tokens correctly.
func TestTokenize_StatusTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// status.error = base12 = "ff899d"
		{"status.error", "#ff899d"},
		// status.warning = base13 = "e9c582"
		{"status.warning", "#e9c582"},
		// status.success = base14 = "afd67a"
		{"status.success", "#afd67a"},
		// status.info = base15 = "97d8f8"
		{"status.info", "#97d8f8"},
		// status.hint = base15 = "97d8f8"
		{"status.hint", "#97d8f8"},
		// status.todo = base0D = "7aa2f7"
		{"status.todo", "#7aa2f7"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}
}

// TestTokenize_DiffTokens tests that Tokenize derives diff tokens correctly.
func TestTokenize_DiffTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// diff.added.fg = base14 = "afd67a"
		{"diff.added.fg", "#afd67a"},
		// diff.added.sign = base14 = "afd67a"
		{"diff.added.sign", "#afd67a"},
		// diff.deleted.fg = base12 = "ff899d"
		{"diff.deleted.fg", "#ff899d"},
		// diff.deleted.sign = base12 = "ff899d"
		{"diff.deleted.sign", "#ff899d"},
		// diff.changed.fg = base16 = "8db6fa"
		{"diff.changed.fg", "#8db6fa"},
		// diff.ignored = base03 = "565f89"
		{"diff.ignored", "#565f89"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}

	// Blended diff backgrounds - just verify they exist.
	blendedTokens := []string{"diff.added.bg", "diff.deleted.bg", "diff.changed.bg"}
	for _, token := range blendedTokens {
		t.Run(token, func(t *testing.T) {
			_, ok := theme.Color(token)
			if !ok {
				t.Fatalf("Color(%q) not found", token)
			}
		})
	}
}

// TestTokenize_MarkupTokens tests that Tokenize derives markup tokens correctly.
func TestTokenize_MarkupTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// markup.heading = base0D = "7aa2f7"
		{"markup.heading", "#7aa2f7"},
		// markup.link = base0C = "7dcfff"
		{"markup.link", "#7dcfff"},
		// markup.code = base0B = "9ece6a"
		{"markup.code", "#9ece6a"},
		// markup.quote = base03 = "565f89"
		{"markup.quote", "#565f89"},
		// markup.list.bullet = base09 = "ff9e64"
		{"markup.list.bullet", "#ff9e64"},
		// markup.list.checked = base0B = "9ece6a"
		{"markup.list.checked", "#9ece6a"},
		// markup.list.unchecked = base0D = "7aa2f7"
		{"markup.list.unchecked", "#7aa2f7"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}
}

// TestTokenize_AccentBorderTokens tests that Tokenize derives accent/border tokens correctly.
func TestTokenize_AccentBorderTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// accent.primary = base0D = "7aa2f7"
		{"accent.primary", "#7aa2f7"},
		// accent.secondary = base0E = "bb9af7"
		{"accent.secondary", "#bb9af7"},
		// accent.foreground = base00 = "1a1b26"
		{"accent.foreground", "#1a1b26"},
		// border.muted = base01 = "1f2335"
		{"border.muted", "#1f2335"},
		// scrollbar.track = base01 = "1f2335"
		{"scrollbar.track", "#1f2335"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}

	// Blended tokens - verify they exist.
	blendedTokens := []string{"border.default", "border.focus", "scrollbar.thumb", "state.active"}
	for _, token := range blendedTokens {
		t.Run(token, func(t *testing.T) {
			_, ok := theme.Color(token)
			if !ok {
				t.Fatalf("Color(%q) not found", token)
			}
		})
	}

	// Alias tokens - verify state.hover = surface.background.highlight.
	t.Run("state.hover", func(t *testing.T) {
		hover, ok := theme.Color("state.hover")
		if !ok {
			t.Fatal("Color(state.hover) not found")
		}
		highlight, ok := theme.Color("surface.background.highlight")
		if !ok {
			t.Fatal("Color(surface.background.highlight) not found")
		}
		if !hover.Equal(highlight) {
			t.Errorf("state.hover = %s, want %s (surface.background.highlight)", hover.Hex(), highlight.Hex())
		}
	})

	// state.disabled.fg = text.muted.
	t.Run("state.disabled.fg", func(t *testing.T) {
		disabled, ok := theme.Color("state.disabled.fg")
		if !ok {
			t.Fatal("Color(state.disabled.fg) not found")
		}
		muted, ok := theme.Color("text.muted")
		if !ok {
			t.Fatal("Color(text.muted) not found")
		}
		if !disabled.Equal(muted) {
			t.Errorf("state.disabled.fg = %s, want %s (text.muted)", disabled.Hex(), muted.Hex())
		}
	})
}

// TestTokenize_GitTokens tests that Tokenize derives git tokens correctly.
func TestTokenize_GitTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// git.added = base0B = "9ece6a"
		{"git.added", "#9ece6a"},
		// git.modified = base0D = "7aa2f7"
		{"git.modified", "#7aa2f7"},
		// git.deleted = base08 = "f7768e"
		{"git.deleted", "#f7768e"},
		// git.ignored = base03 = "565f89"
		{"git.ignored", "#565f89"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}
}

// TestTokenize_TerminalTokens tests that Tokenize derives terminal ANSI tokens correctly.
func TestTokenize_TerminalTokens(t *testing.T) {
	pal := parsePalette(t)
	theme := flair.Tokenize(pal)

	if theme == nil {
		t.Fatal("Tokenize() returned nil")
	}

	tests := []struct {
		token   string
		wantHex string
	}{
		// terminal.black = base01 = "1f2335"
		{"terminal.black", "#1f2335"},
		// terminal.red = base08 = "f7768e"
		{"terminal.red", "#f7768e"},
		// terminal.green = base0B = "9ece6a"
		{"terminal.green", "#9ece6a"},
		// terminal.yellow = base0A = "e0af68"
		{"terminal.yellow", "#e0af68"},
		// terminal.blue = base0D = "7aa2f7"
		{"terminal.blue", "#7aa2f7"},
		// terminal.magenta = base0E = "bb9af7"
		{"terminal.magenta", "#bb9af7"},
		// terminal.cyan = base0C = "7dcfff"
		{"terminal.cyan", "#7dcfff"},
		// terminal.white = base05 = "c0caf5"
		{"terminal.white", "#c0caf5"},
		// terminal.brblack = base03 = "565f89"
		{"terminal.brblack", "#565f89"},
		// terminal.brred = base12 = "ff899d"
		{"terminal.brred", "#ff899d"},
		// terminal.brgreen = base14 = "afd67a"
		{"terminal.brgreen", "#afd67a"},
		// terminal.bryellow = base13 = "e9c582"
		{"terminal.bryellow", "#e9c582"},
		// terminal.brblue = base16 = "8db6fa"
		{"terminal.brblue", "#8db6fa"},
		// terminal.brmagenta = base17 = "c8acf8"
		{"terminal.brmagenta", "#c8acf8"},
		// terminal.brcyan = base15 = "97d8f8"
		{"terminal.brcyan", "#97d8f8"},
		// terminal.brwhite = base07 = "c8d3f5"
		{"terminal.brwhite", "#c8d3f5"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			c, ok := theme.Color(tt.token)
			if !ok {
				t.Fatalf("Color(%q) not found", tt.token)
			}
			if got := c.Hex(); got != tt.wantHex {
				t.Errorf("Color(%q).Hex() = %q, want %q", tt.token, got, tt.wantHex)
			}
		})
	}
}

// TestBlendBg tests the BlendBg function for alpha blending colors.
func TestBlendBg(t *testing.T) {
	tests := []struct {
		name    string
		fg      flair.Color
		bg      flair.Color
		alpha   float64
		wantHex string
	}{
		{
			name:    "50% blend",
			fg:      flair.Color{R: 255, G: 0, B: 0}, // red
			bg:      flair.Color{R: 0, G: 0, B: 255}, // blue
			alpha:   0.5,
			wantHex: "#800080", // 50% red + 50% blue = purple (128, 0, 128)
		},
		{
			name:    "0% blend (all background)",
			fg:      flair.Color{R: 255, G: 0, B: 0},
			bg:      flair.Color{R: 0, G: 255, B: 0},
			alpha:   0.0,
			wantHex: "#00ff00", // pure green (background)
		},
		{
			name:    "100% blend (all foreground)",
			fg:      flair.Color{R: 255, G: 0, B: 0},
			bg:      flair.Color{R: 0, G: 255, B: 0},
			alpha:   1.0,
			wantHex: "#ff0000", // pure red (foreground)
		},
		{
			name:    "30% blend for selection-like color",
			fg:      flair.Color{R: 0x7a, G: 0xa2, B: 0xf7}, // base0D
			bg:      flair.Color{R: 0x1a, G: 0x1b, B: 0x26}, // base00
			alpha:   0.30,
			wantHex: "#374465", // 0.3*122 + 0.7*26 = 55, 0.3*162 + 0.7*27 = 68, 0.3*247 + 0.7*38 = 101
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flair.BlendBg(tt.fg, tt.bg, tt.alpha)
			if got.Hex() != tt.wantHex {
				t.Errorf("BlendBg() = %s, want %s", got.Hex(), tt.wantHex)
			}
		})
	}
}

// TestTokenize_NilPalette tests that Tokenize handles nil palette gracefully.
func TestTokenize_NilPalette(t *testing.T) {
	theme := flair.Tokenize(nil)
	if theme != nil {
		t.Error("Tokenize(nil) should return nil")
	}
}

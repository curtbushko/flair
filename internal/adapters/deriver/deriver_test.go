package deriver_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/deriver"
	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// tokyoNightDarkPalette returns the reference Tokyo Night Dark base24 palette
// used as the canonical test fixture throughout the deriver tests.
func tokyoNightDarkPalette(t *testing.T) *domain.Palette {
	t.Helper()

	colors := map[string]string{
		"base00": "1a1b26",
		"base01": "1f2335",
		"base02": "292e42",
		"base03": "565f89",
		"base04": "a9b1d6",
		"base05": "c0caf5",
		"base06": "c0caf5",
		"base07": "c8d3f5",
		"base08": "f7768e",
		"base09": "ff9e64",
		"base0A": "e0af68",
		"base0B": "9ece6a",
		"base0C": "7dcfff",
		"base0D": "7aa2f7",
		"base0E": "bb9af7",
		"base0F": "db4b4b",
		"base10": "16161e",
		"base11": "101014",
		"base12": "ff899d",
		"base13": "e9c582",
		"base14": "afd67a",
		"base15": "97d8f8",
		"base16": "8db6fa",
		"base17": "c8acf8",
	}

	pal, err := domain.NewPalette("Tokyo Night Dark", "Michael Ball", "dark", "base24", colors)
	if err != nil {
		t.Fatalf("failed to create test palette: %v", err)
	}
	return pal
}

// mustParseHex is a test helper that parses a hex color or fails the test.
func mustParseHex(t *testing.T, hex string) domain.Color {
	t.Helper()
	c, err := domain.ParseHex(hex)
	if err != nil {
		t.Fatalf("failed to parse hex %q: %v", hex, err)
	}
	return c
}

func TestDefaultDeriver_ImplementsInterface(t *testing.T) {
	var _ ports.TokenDeriver = &deriver.DefaultDeriver{}
}

func TestDeriveSurface_Background(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("surface.background")
	if !ok {
		t.Fatal("surface.background not found in token set")
	}

	want := mustParseHex(t, "#1a1b26")
	if !tok.Color.Equal(want) {
		t.Errorf("surface.background = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSurface_BackgroundRaised(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("surface.background.raised")
	if !ok {
		t.Fatal("surface.background.raised not found in token set")
	}

	want := mustParseHex(t, "#1f2335")
	if !tok.Color.Equal(want) {
		t.Errorf("surface.background.raised = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSurface_BackgroundSunken(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("surface.background.sunken")
	if !ok {
		t.Fatal("surface.background.sunken not found in token set")
	}

	want := mustParseHex(t, "#16161e")
	if !tok.Color.Equal(want) {
		t.Errorf("surface.background.sunken = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSurface_BackgroundDarkest(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("surface.background.darkest")
	if !ok {
		t.Fatal("surface.background.darkest not found in token set")
	}

	want := mustParseHex(t, "#101014")
	if !tok.Color.Equal(want) {
		t.Errorf("surface.background.darkest = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSurface_BackgroundHighlight(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("surface.background.highlight")
	if !ok {
		t.Fatal("surface.background.highlight not found in token set")
	}

	want := mustParseHex(t, "#292e42")
	if !tok.Color.Equal(want) {
		t.Errorf("surface.background.highlight = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSurface_BackgroundSelection(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("surface.background.selection")
	if !ok {
		t.Fatal("surface.background.selection not found in token set")
	}

	// BlendBg(base0D, base00, 0.30) = Blend(base00, base0D, 0.30)
	want := domain.BlendBg(pal.Base(0x0D), pal.Base(0x00), 0.30)
	if !tok.Color.Equal(want) {
		t.Errorf("surface.background.selection = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSurface_BackgroundSearch(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("surface.background.search")
	if !ok {
		t.Fatal("surface.background.search not found in token set")
	}

	// BlendBg(base0A, base00, 0.30) = Blend(base00, base0A, 0.30)
	want := domain.BlendBg(pal.Base(0x0A), pal.Base(0x00), 0.30)
	if !tok.Color.Equal(want) {
		t.Errorf("surface.background.search = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSurface_Base10Aliases(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	want := mustParseHex(t, "#16161e") // base10
	aliases := []string{
		"surface.background.overlay",
		"surface.background.popup",
		"surface.background.sidebar",
		"surface.background.statusbar",
	}

	for _, path := range aliases {
		tok, ok := ts.Get(path)
		if !ok {
			t.Errorf("%s not found in token set", path)
			continue
		}
		if !tok.Color.Equal(want) {
			t.Errorf("%s = %s, want %s", path, tok.Color.Hex(), want.Hex())
		}
	}
}

func TestDeriveSurface_AllPresent(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	surfacePaths := []string{
		"surface.background",
		"surface.background.raised",
		"surface.background.sunken",
		"surface.background.darkest",
		"surface.background.highlight",
		"surface.background.selection",
		"surface.background.search",
		"surface.background.overlay",
		"surface.background.popup",
		"surface.background.sidebar",
		"surface.background.statusbar",
	}

	for _, path := range surfacePaths {
		if _, ok := ts.Get(path); !ok {
			t.Errorf("missing surface token: %s", path)
		}
	}

	// Verify exactly 11 surface tokens
	count := 0
	for _, p := range ts.Paths() {
		if len(p) >= 7 && p[:8] == "surface." {
			count++
		}
	}
	if count != 11 {
		t.Errorf("expected 11 surface tokens, got %d", count)
	}
}

// --- Text token tests ---

func TestDeriveText_Primary(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("text.primary")
	if !ok {
		t.Fatal("text.primary not found in token set")
	}

	want := mustParseHex(t, "#c0caf5") // base05
	if !tok.Color.Equal(want) {
		t.Errorf("text.primary = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveText_Secondary(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("text.secondary")
	if !ok {
		t.Fatal("text.secondary not found in token set")
	}

	want := mustParseHex(t, "#a9b1d6") // base04
	if !tok.Color.Equal(want) {
		t.Errorf("text.secondary = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveText_Muted(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("text.muted")
	if !ok {
		t.Fatal("text.muted not found in token set")
	}

	want := mustParseHex(t, "#565f89") // base03
	if !tok.Color.Equal(want) {
		t.Errorf("text.muted = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveText_Subtle(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("text.subtle")
	if !ok {
		t.Fatal("text.subtle not found in token set")
	}

	// BlendBg(base03, base00, 0.50) = Blend(base00, base03, 0.50)
	want := domain.BlendBg(pal.Base(0x03), pal.Base(0x00), 0.50)
	if !tok.Color.Equal(want) {
		t.Errorf("text.subtle = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveText_Inverse(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("text.inverse")
	if !ok {
		t.Fatal("text.inverse not found in token set")
	}

	want := mustParseHex(t, "#1a1b26") // base00
	if !tok.Color.Equal(want) {
		t.Errorf("text.inverse = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveText_Overlay(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("text.overlay")
	if !ok {
		t.Fatal("text.overlay not found in token set")
	}

	want := mustParseHex(t, "#c0caf5") // base06
	if !tok.Color.Equal(want) {
		t.Errorf("text.overlay = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveText_Sidebar(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("text.sidebar")
	if !ok {
		t.Fatal("text.sidebar not found in token set")
	}

	want := mustParseHex(t, "#a9b1d6") // base04
	if !tok.Color.Equal(want) {
		t.Errorf("text.sidebar = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveText_AllPresent(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	textPaths := []string{
		"text.primary",
		"text.secondary",
		"text.muted",
		"text.subtle",
		"text.inverse",
		"text.overlay",
		"text.sidebar",
	}

	for _, path := range textPaths {
		if _, ok := ts.Get(path); !ok {
			t.Errorf("missing text token: %s", path)
		}
	}

	// Verify exactly 7 text tokens
	count := 0
	for _, p := range ts.Paths() {
		if len(p) >= 5 && p[:5] == "text." {
			count++
		}
	}
	if count != 7 {
		t.Errorf("expected 7 text tokens, got %d", count)
	}
}

// --- Status token tests ---

func TestDeriveStatus_Error(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("status.error")
	if !ok {
		t.Fatal("status.error not found in token set")
	}

	want := mustParseHex(t, "#ff899d") // base12
	if !tok.Color.Equal(want) {
		t.Errorf("status.error = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveStatus_Warning(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("status.warning")
	if !ok {
		t.Fatal("status.warning not found in token set")
	}

	want := mustParseHex(t, "#e9c582") // base13
	if !tok.Color.Equal(want) {
		t.Errorf("status.warning = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveStatus_Success(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("status.success")
	if !ok {
		t.Fatal("status.success not found in token set")
	}

	want := mustParseHex(t, "#afd67a") // base14
	if !tok.Color.Equal(want) {
		t.Errorf("status.success = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveStatus_InfoHint(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	want := mustParseHex(t, "#97d8f8") // base15

	tok, ok := ts.Get("status.info")
	if !ok {
		t.Fatal("status.info not found in token set")
	}
	if !tok.Color.Equal(want) {
		t.Errorf("status.info = %s, want %s", tok.Color.Hex(), want.Hex())
	}

	tok, ok = ts.Get("status.hint")
	if !ok {
		t.Fatal("status.hint not found in token set")
	}
	if !tok.Color.Equal(want) {
		t.Errorf("status.hint = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveStatus_Todo(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("status.todo")
	if !ok {
		t.Fatal("status.todo not found in token set")
	}

	want := mustParseHex(t, "#7aa2f7") // base0D
	if !tok.Color.Equal(want) {
		t.Errorf("status.todo = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

// --- Diff token tests ---

func TestDeriveDiff_AddedFg(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.added.fg")
	if !ok {
		t.Fatal("diff.added.fg not found in token set")
	}

	want := mustParseHex(t, "#afd67a") // base14
	if !tok.Color.Equal(want) {
		t.Errorf("diff.added.fg = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_AddedBg(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.added.bg")
	if !ok {
		t.Fatal("diff.added.bg not found in token set")
	}

	// BlendBg(base0B, base00, 0.25)
	want := domain.BlendBg(pal.Base(0x0B), pal.Base(0x00), 0.25)
	if !tok.Color.Equal(want) {
		t.Errorf("diff.added.bg = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_AddedSign(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.added.sign")
	if !ok {
		t.Fatal("diff.added.sign not found in token set")
	}

	want := mustParseHex(t, "#afd67a") // base14
	if !tok.Color.Equal(want) {
		t.Errorf("diff.added.sign = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_DeletedFg(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.deleted.fg")
	if !ok {
		t.Fatal("diff.deleted.fg not found in token set")
	}

	want := mustParseHex(t, "#ff899d") // base12
	if !tok.Color.Equal(want) {
		t.Errorf("diff.deleted.fg = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_DeletedBg(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.deleted.bg")
	if !ok {
		t.Fatal("diff.deleted.bg not found in token set")
	}

	// BlendBg(base08, base00, 0.25)
	want := domain.BlendBg(pal.Base(0x08), pal.Base(0x00), 0.25)
	if !tok.Color.Equal(want) {
		t.Errorf("diff.deleted.bg = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_DeletedSign(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.deleted.sign")
	if !ok {
		t.Fatal("diff.deleted.sign not found in token set")
	}

	want := mustParseHex(t, "#ff899d") // base12
	if !tok.Color.Equal(want) {
		t.Errorf("diff.deleted.sign = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_ChangedFg(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.changed.fg")
	if !ok {
		t.Fatal("diff.changed.fg not found in token set")
	}

	want := mustParseHex(t, "#8db6fa") // base16
	if !tok.Color.Equal(want) {
		t.Errorf("diff.changed.fg = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_ChangedBg(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.changed.bg")
	if !ok {
		t.Fatal("diff.changed.bg not found in token set")
	}

	// BlendBg(base0D, base00, 0.15)
	want := domain.BlendBg(pal.Base(0x0D), pal.Base(0x00), 0.15)
	if !tok.Color.Equal(want) {
		t.Errorf("diff.changed.bg = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveDiff_Ignored(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("diff.ignored")
	if !ok {
		t.Fatal("diff.ignored not found in token set")
	}

	want := mustParseHex(t, "#565f89") // base03
	if !tok.Color.Equal(want) {
		t.Errorf("diff.ignored = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveStatusDiff_AllPresent(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	statusPaths := []string{
		"status.error",
		"status.warning",
		"status.success",
		"status.info",
		"status.hint",
		"status.todo",
	}
	for _, path := range statusPaths {
		if _, ok := ts.Get(path); !ok {
			t.Errorf("missing status token: %s", path)
		}
	}

	diffPaths := []string{
		"diff.added.fg",
		"diff.added.bg",
		"diff.added.sign",
		"diff.deleted.fg",
		"diff.deleted.bg",
		"diff.deleted.sign",
		"diff.changed.fg",
		"diff.changed.bg",
		"diff.ignored",
	}
	for _, path := range diffPaths {
		if _, ok := ts.Get(path); !ok {
			t.Errorf("missing diff token: %s", path)
		}
	}

	// Verify counts
	statusCount := 0
	diffCount := 0
	for _, p := range ts.Paths() {
		if len(p) >= 7 && p[:7] == "status." {
			statusCount++
		}
		if len(p) >= 5 && p[:5] == "diff." {
			diffCount++
		}
	}
	if statusCount != 6 {
		t.Errorf("expected 6 status tokens, got %d", statusCount)
	}
	if diffCount != 9 {
		t.Errorf("expected 9 diff tokens, got %d", diffCount)
	}
}

// --- Syntax token tests ---

func TestDeriveSyntax_Keyword(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.keyword")
	if !ok {
		t.Fatal("syntax.keyword not found in token set")
	}

	want := mustParseHex(t, "#bb9af7") // base0E
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.keyword = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_String(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.string")
	if !ok {
		t.Fatal("syntax.string not found in token set")
	}

	want := mustParseHex(t, "#9ece6a") // base0B
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.string = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Function(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.function")
	if !ok {
		t.Fatal("syntax.function not found in token set")
	}

	want := mustParseHex(t, "#7aa2f7") // base0D
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.function = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Comment(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.comment")
	if !ok {
		t.Fatal("syntax.comment not found in token set")
	}

	want := mustParseHex(t, "#565f89") // base03
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.comment color = %s, want %s", tok.Color.Hex(), want.Hex())
	}
	if !tok.Italic {
		t.Error("syntax.comment should have Italic=true")
	}
}

func TestDeriveSyntax_Variable(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.variable")
	if !ok {
		t.Fatal("syntax.variable not found in token set")
	}

	want := mustParseHex(t, "#c0caf5") // base05
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.variable = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Constant(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.constant")
	if !ok {
		t.Fatal("syntax.constant not found in token set")
	}

	want := mustParseHex(t, "#ff9e64") // base09
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.constant = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Operator(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.operator")
	if !ok {
		t.Fatal("syntax.operator not found in token set")
	}

	want := mustParseHex(t, "#8db6fa") // base16
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.operator = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Type(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.type")
	if !ok {
		t.Fatal("syntax.type not found in token set")
	}

	want := mustParseHex(t, "#e0af68") // base0A
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.type = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Number(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.number")
	if !ok {
		t.Fatal("syntax.number not found in token set")
	}

	want := mustParseHex(t, "#ff9e64") // base09
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.number = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Tag(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.tag")
	if !ok {
		t.Fatal("syntax.tag not found in token set")
	}

	want := mustParseHex(t, "#f7768e") // base08
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.tag = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Property(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.property")
	if !ok {
		t.Fatal("syntax.property not found in token set")
	}

	want := mustParseHex(t, "#afd67a") // base14
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.property = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Parameter(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.parameter")
	if !ok {
		t.Fatal("syntax.parameter not found in token set")
	}

	want := mustParseHex(t, "#e9c582") // base13
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.parameter = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Regexp(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.regexp")
	if !ok {
		t.Fatal("syntax.regexp not found in token set")
	}

	want := mustParseHex(t, "#7dcfff") // base0C
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.regexp = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Escape(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.escape")
	if !ok {
		t.Fatal("syntax.escape not found in token set")
	}

	want := mustParseHex(t, "#bb9af7") // base0E
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.escape = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_Constructor(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.constructor")
	if !ok {
		t.Fatal("syntax.constructor not found in token set")
	}

	want := mustParseHex(t, "#c8acf8") // base17
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.constructor = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveSyntax_AllPresent(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	syntaxPaths := []string{
		"syntax.keyword",
		"syntax.string",
		"syntax.function",
		"syntax.comment",
		"syntax.variable",
		"syntax.constant",
		"syntax.operator",
		"syntax.type",
		"syntax.number",
		"syntax.tag",
		"syntax.property",
		"syntax.parameter",
		"syntax.regexp",
		"syntax.escape",
		"syntax.constructor",
	}

	for _, path := range syntaxPaths {
		if _, ok := ts.Get(path); !ok {
			t.Errorf("missing syntax token: %s", path)
		}
	}

	// Verify exactly 15 syntax tokens (14 from PLAN + constructor)
	count := 0
	for _, p := range ts.Paths() {
		if len(p) >= 7 && p[:7] == "syntax." {
			count++
		}
	}
	if count != 15 {
		t.Errorf("expected 15 syntax tokens, got %d", count)
	}
}

// --- Markup token tests ---

func TestDeriveMarkup_Heading(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.heading")
	if !ok {
		t.Fatal("markup.heading not found in token set")
	}

	want := mustParseHex(t, "#7aa2f7") // base0D
	if !tok.Color.Equal(want) {
		t.Errorf("markup.heading color = %s, want %s", tok.Color.Hex(), want.Hex())
	}
	if !tok.Bold {
		t.Error("markup.heading should have Bold=true")
	}
}

func TestDeriveMarkup_Link(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.link")
	if !ok {
		t.Fatal("markup.link not found in token set")
	}

	want := mustParseHex(t, "#7dcfff") // base0C
	if !tok.Color.Equal(want) {
		t.Errorf("markup.link = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_Code(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.code")
	if !ok {
		t.Fatal("markup.code not found in token set")
	}

	want := mustParseHex(t, "#9ece6a") // base0B
	if !tok.Color.Equal(want) {
		t.Errorf("markup.code = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_Bold(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.bold")
	if !ok {
		t.Fatal("markup.bold not found in token set")
	}

	if !tok.Bold {
		t.Error("markup.bold should have Bold=true")
	}
	want := pal.Base(0x05)
	if tok.Color.Hex() != want.Hex() {
		t.Errorf("markup.bold color = %s, want %s (base05)", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_Italic(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.italic")
	if !ok {
		t.Fatal("markup.italic not found in token set")
	}

	if !tok.Italic {
		t.Error("markup.italic should have Italic=true")
	}
	want := pal.Base(0x05)
	if tok.Color.Hex() != want.Hex() {
		t.Errorf("markup.italic color = %s, want %s (base05)", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_Strikethrough(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.strikethrough")
	if !ok {
		t.Fatal("markup.strikethrough not found in token set")
	}

	if !tok.Strikethrough {
		t.Error("markup.strikethrough should have Strikethrough=true")
	}
	want := pal.Base(0x03)
	if tok.Color.Hex() != want.Hex() {
		t.Errorf("markup.strikethrough color = %s, want %s (base03)", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_Quote(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.quote")
	if !ok {
		t.Fatal("markup.quote not found in token set")
	}

	want := mustParseHex(t, "#565f89") // base03
	if !tok.Color.Equal(want) {
		t.Errorf("markup.quote color = %s, want %s", tok.Color.Hex(), want.Hex())
	}
	if !tok.Italic {
		t.Error("markup.quote should have Italic=true")
	}
}

func TestDeriveMarkup_ListBullet(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.list.bullet")
	if !ok {
		t.Fatal("markup.list.bullet not found in token set")
	}

	want := mustParseHex(t, "#ff9e64") // base09
	if !tok.Color.Equal(want) {
		t.Errorf("markup.list.bullet = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_ListChecked(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.list.checked")
	if !ok {
		t.Fatal("markup.list.checked not found in token set")
	}

	want := mustParseHex(t, "#9ece6a") // base0B
	if !tok.Color.Equal(want) {
		t.Errorf("markup.list.checked = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_ListUnchecked(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.list.unchecked")
	if !ok {
		t.Fatal("markup.list.unchecked not found in token set")
	}

	want := mustParseHex(t, "#7aa2f7") // base0D
	if !tok.Color.Equal(want) {
		t.Errorf("markup.list.unchecked = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveMarkup_AllPresent(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	markupPaths := []string{
		"markup.heading",
		"markup.link",
		"markup.code",
		"markup.bold",
		"markup.italic",
		"markup.strikethrough",
		"markup.quote",
		"markup.list.bullet",
		"markup.list.checked",
		"markup.list.unchecked",
	}

	for _, path := range markupPaths {
		if _, ok := ts.Get(path); !ok {
			t.Errorf("missing markup token: %s", path)
		}
	}

	// Verify exactly 10 markup tokens
	count := 0
	for _, p := range ts.Paths() {
		if len(p) >= 7 && p[:7] == "markup." {
			count++
		}
	}
	if count != 10 {
		t.Errorf("expected 10 markup tokens, got %d", count)
	}
}

// --- Accent token tests ---

func TestDeriveAccent_Primary(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("accent.primary")
	if !ok {
		t.Fatal("accent.primary not found in token set")
	}

	want := mustParseHex(t, "#7aa2f7") // base0D
	if !tok.Color.Equal(want) {
		t.Errorf("accent.primary = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveAccent_Secondary(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("accent.secondary")
	if !ok {
		t.Fatal("accent.secondary not found in token set")
	}

	want := mustParseHex(t, "#bb9af7") // base0E
	if !tok.Color.Equal(want) {
		t.Errorf("accent.secondary = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveAccent_Foreground(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("accent.foreground")
	if !ok {
		t.Fatal("accent.foreground not found in token set")
	}

	want := mustParseHex(t, "#1a1b26") // base00
	if !tok.Color.Equal(want) {
		t.Errorf("accent.foreground = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

// --- Border token tests ---

func TestDeriveBorder_Default(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("border.default")
	if !ok {
		t.Fatal("border.default not found in token set")
	}

	// BlendBg(base03, base00, 0.40)
	want := domain.BlendBg(pal.Base(0x03), pal.Base(0x00), 0.40)
	if !tok.Color.Equal(want) {
		t.Errorf("border.default = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveBorder_Focus(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("border.focus")
	if !ok {
		t.Fatal("border.focus not found in token set")
	}

	// BlendBg(base0D, base00, 0.70)
	want := domain.BlendBg(pal.Base(0x0D), pal.Base(0x00), 0.70)
	if !tok.Color.Equal(want) {
		t.Errorf("border.focus = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveBorder_Muted(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("border.muted")
	if !ok {
		t.Fatal("border.muted not found in token set")
	}

	want := mustParseHex(t, "#1f2335") // base01
	if !tok.Color.Equal(want) {
		t.Errorf("border.muted = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

// --- Scrollbar token tests ---

func TestDeriveScrollbar_Thumb(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("scrollbar.thumb")
	if !ok {
		t.Fatal("scrollbar.thumb not found in token set")
	}

	// BlendBg(base03, base00, 0.40)
	want := domain.BlendBg(pal.Base(0x03), pal.Base(0x00), 0.40)
	if !tok.Color.Equal(want) {
		t.Errorf("scrollbar.thumb = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveScrollbar_Track(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("scrollbar.track")
	if !ok {
		t.Fatal("scrollbar.track not found in token set")
	}

	want := mustParseHex(t, "#1f2335") // base01
	if !tok.Color.Equal(want) {
		t.Errorf("scrollbar.track = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

// --- State token tests ---

func TestDeriveState_Hover(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("state.hover")
	if !ok {
		t.Fatal("state.hover not found in token set")
	}

	// Aliases surface.background.highlight (base02)
	want := mustParseHex(t, "#292e42")
	if !tok.Color.Equal(want) {
		t.Errorf("state.hover = %s, want %s", tok.Color.Hex(), want.Hex())
	}

	// Verify it matches the surface.background.highlight token
	surfTok, ok := ts.Get("surface.background.highlight")
	if !ok {
		t.Fatal("surface.background.highlight not found for alias check")
	}
	if !tok.Color.Equal(surfTok.Color) {
		t.Errorf("state.hover (%s) should match surface.background.highlight (%s)",
			tok.Color.Hex(), surfTok.Color.Hex())
	}
}

func TestDeriveState_Active(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("state.active")
	if !ok {
		t.Fatal("state.active not found in token set")
	}

	// BlendBg(base0D, base00, 0.20)
	want := domain.BlendBg(pal.Base(0x0D), pal.Base(0x00), 0.20)
	if !tok.Color.Equal(want) {
		t.Errorf("state.active = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveState_DisabledFg(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("state.disabled.fg")
	if !ok {
		t.Fatal("state.disabled.fg not found in token set")
	}

	// Aliases text.muted (base03)
	want := mustParseHex(t, "#565f89")
	if !tok.Color.Equal(want) {
		t.Errorf("state.disabled.fg = %s, want %s", tok.Color.Hex(), want.Hex())
	}

	// Verify it matches the text.muted token
	textTok, ok := ts.Get("text.muted")
	if !ok {
		t.Fatal("text.muted not found for alias check")
	}
	if !tok.Color.Equal(textTok.Color) {
		t.Errorf("state.disabled.fg (%s) should match text.muted (%s)",
			tok.Color.Hex(), textTok.Color.Hex())
	}
}

// --- Git token tests ---

func TestDeriveGit_Added(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("git.added")
	if !ok {
		t.Fatal("git.added not found in token set")
	}

	want := mustParseHex(t, "#9ece6a") // base0B
	if !tok.Color.Equal(want) {
		t.Errorf("git.added = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveGit_Modified(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("git.modified")
	if !ok {
		t.Fatal("git.modified not found in token set")
	}

	want := mustParseHex(t, "#7aa2f7") // base0D
	if !tok.Color.Equal(want) {
		t.Errorf("git.modified = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveGit_Deleted(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("git.deleted")
	if !ok {
		t.Fatal("git.deleted not found in token set")
	}

	want := mustParseHex(t, "#f7768e") // base08
	if !tok.Color.Equal(want) {
		t.Errorf("git.deleted = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

func TestDeriveGit_Ignored(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("git.ignored")
	if !ok {
		t.Fatal("git.ignored not found in token set")
	}

	want := mustParseHex(t, "#565f89") // base03
	if !tok.Color.Equal(want) {
		t.Errorf("git.ignored = %s, want %s", tok.Color.Hex(), want.Hex())
	}
}

// --- Terminal ANSI token tests ---

func TestDeriveTerminal_AllColors(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tests := []struct {
		path string
		base int
	}{
		{"terminal.black", 0x01},
		{"terminal.red", 0x08},
		{"terminal.green", 0x0B},
		{"terminal.yellow", 0x0A},
		{"terminal.blue", 0x0D},
		{"terminal.magenta", 0x0E},
		{"terminal.cyan", 0x0C},
		{"terminal.white", 0x05},
		{"terminal.brblack", 0x03},
		{"terminal.brred", 0x12},
		{"terminal.brgreen", 0x14},
		{"terminal.bryellow", 0x13},
		{"terminal.brblue", 0x16},
		{"terminal.brmagenta", 0x17},
		{"terminal.brcyan", 0x15},
		{"terminal.brwhite", 0x07},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			tok, ok := ts.Get(tt.path)
			if !ok {
				t.Fatalf("%s not found in token set", tt.path)
			}

			want := pal.Base(tt.base)
			if !tok.Color.Equal(want) {
				t.Errorf("%s = %s, want %s", tt.path, tok.Color.Hex(), want.Hex())
			}
		})
	}
}

// --- All 31 accent/border/scrollbar/state/git/terminal tokens present ---

func TestDeriveAccentGitTerminal_AllPresent(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	allPaths := []string{
		// Accent (3)
		"accent.primary",
		"accent.secondary",
		"accent.foreground",
		// Border (3)
		"border.default",
		"border.focus",
		"border.muted",
		// Scrollbar (2)
		"scrollbar.thumb",
		"scrollbar.track",
		// State (3)
		"state.hover",
		"state.active",
		"state.disabled.fg",
		// Git (4)
		"git.added",
		"git.modified",
		"git.deleted",
		"git.ignored",
		// Terminal (16)
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

	for _, path := range allPaths {
		if _, ok := ts.Get(path); !ok {
			t.Errorf("missing token: %s", path)
		}
	}

	if len(allPaths) != 31 {
		t.Fatalf("expected 31 paths in test, got %d", len(allPaths))
	}
}

// =============================================================================
// Comprehensive Integration Tests (Task 10 — Full Derivation Validation)
// =============================================================================

// TestFullDerivation_TokenCount verifies that full derivation produces at least
// 87 tokens from the Tokyo Night Dark palette.
func TestFullDerivation_TokenCount(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	if ts.Len() < 87 {
		t.Errorf("expected at least 87 tokens, got %d", ts.Len())
	}
}

// TestFullDerivation_AllColorsPresent verifies that every token path has a
// valid RGB color.
func TestFullDerivation_AllColorsPresent(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	for _, path := range ts.Paths() {
		tok, ok := ts.Get(path)
		if !ok {
			t.Errorf("path %q listed but not retrievable", path)
			continue
		}

		if tok.Color.IsNone {
			t.Errorf("token %q should have a valid RGB color, got NoneColor", path)
		}
	}
}

// TestFullDerivation_SurfaceValues verifies exact hex values for all surface
// tokens against the Tokyo Night Dark palette.
func TestFullDerivation_SurfaceValues(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tests := []struct {
		path string
		want string
	}{
		{"surface.background", "#1a1b26"},
		{"surface.background.raised", "#1f2335"},
		{"surface.background.sunken", "#16161e"},
		{"surface.background.darkest", "#101014"},
		{"surface.background.highlight", "#292e42"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			tok, ok := ts.Get(tt.path)
			if !ok {
				t.Fatalf("%s not found in token set", tt.path)
			}
			want := mustParseHex(t, tt.want)
			if !tok.Color.Equal(want) {
				t.Errorf("%s = %s, want %s", tt.path, tok.Color.Hex(), want.Hex())
			}
		})
	}
}

// TestFullDerivation_BlendedValues verifies that all blended tokens match
// independently computed BlendBg results. This is the key cross-check:
// the test independently computes the same blend operations and compares.
func TestFullDerivation_BlendedValues(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tests := []struct {
		path   string
		fg     int // base slot for fg
		bg     int // base slot for bg
		amount float64
	}{
		{"surface.background.selection", 0x0D, 0x00, 0.30},
		{"surface.background.search", 0x0A, 0x00, 0.30},
		{"text.subtle", 0x03, 0x00, 0.50},
		{"border.default", 0x03, 0x00, 0.40},
		{"border.focus", 0x0D, 0x00, 0.70},
		{"scrollbar.thumb", 0x03, 0x00, 0.40},
		{"state.active", 0x0D, 0x00, 0.20},
		{"diff.added.bg", 0x0B, 0x00, 0.25},
		{"diff.deleted.bg", 0x08, 0x00, 0.25},
		{"diff.changed.bg", 0x0D, 0x00, 0.15},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			tok, ok := ts.Get(tt.path)
			if !ok {
				t.Fatalf("%s not found in token set", tt.path)
			}

			// Independently compute the expected value.
			want := domain.BlendBg(pal.Base(tt.fg), pal.Base(tt.bg), tt.amount)
			if !tok.Color.Equal(want) {
				t.Errorf("%s = %s, want %s (BlendBg(base%02X, base%02X, %.2f))",
					tt.path, tok.Color.Hex(), want.Hex(), tt.fg, tt.bg, tt.amount)
			}
		})
	}
}

// TestFullDerivation_SyntaxCommentStyle verifies that syntax.comment has
// Italic=true and the correct color (#565f89 = base03).
func TestFullDerivation_SyntaxCommentStyle(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("syntax.comment")
	if !ok {
		t.Fatal("syntax.comment not found in token set")
	}

	want := mustParseHex(t, "#565f89")
	if !tok.Color.Equal(want) {
		t.Errorf("syntax.comment color = %s, want %s", tok.Color.Hex(), want.Hex())
	}
	if !tok.Italic {
		t.Error("syntax.comment should have Italic=true")
	}
	if tok.Bold {
		t.Error("syntax.comment should not have Bold=true")
	}
}

// TestFullDerivation_MarkupHeadingStyle verifies that markup.heading has
// Bold=true and the correct color (#7aa2f7 = base0D).
func TestFullDerivation_MarkupHeadingStyle(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tok, ok := ts.Get("markup.heading")
	if !ok {
		t.Fatal("markup.heading not found in token set")
	}

	want := mustParseHex(t, "#7aa2f7")
	if !tok.Color.Equal(want) {
		t.Errorf("markup.heading color = %s, want %s", tok.Color.Hex(), want.Hex())
	}
	if !tok.Bold {
		t.Error("markup.heading should have Bold=true")
	}
	if tok.Italic {
		t.Error("markup.heading should not have Italic=true")
	}
}

// TestFullDerivation_TerminalColors verifies that all 16 terminal tokens
// map to their correct base slots from the Tokyo Night Dark palette.
func TestFullDerivation_TerminalColors(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	ts := d.Derive(pal)

	tests := []struct {
		path string
		base int
		hex  string
	}{
		{"terminal.black", 0x01, "#1f2335"},
		{"terminal.red", 0x08, "#f7768e"},
		{"terminal.green", 0x0B, "#9ece6a"},
		{"terminal.yellow", 0x0A, "#e0af68"},
		{"terminal.blue", 0x0D, "#7aa2f7"},
		{"terminal.magenta", 0x0E, "#bb9af7"},
		{"terminal.cyan", 0x0C, "#7dcfff"},
		{"terminal.white", 0x05, "#c0caf5"},
		{"terminal.brblack", 0x03, "#565f89"},
		{"terminal.brred", 0x12, "#ff899d"},
		{"terminal.brgreen", 0x14, "#afd67a"},
		{"terminal.bryellow", 0x13, "#e9c582"},
		{"terminal.brblue", 0x16, "#8db6fa"},
		{"terminal.brmagenta", 0x17, "#c8acf8"},
		{"terminal.brcyan", 0x15, "#97d8f8"},
		{"terminal.brwhite", 0x07, "#c8d3f5"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			tok, ok := ts.Get(tt.path)
			if !ok {
				t.Fatalf("%s not found in token set", tt.path)
			}

			// Verify against base slot
			baseColor := pal.Base(tt.base)
			if !tok.Color.Equal(baseColor) {
				t.Errorf("%s = %s, want %s (base%02X)", tt.path, tok.Color.Hex(), baseColor.Hex(), tt.base)
			}

			// Also verify against expected hex for double-confirmation
			wantHex := mustParseHex(t, tt.hex)
			if !tok.Color.Equal(wantHex) {
				t.Errorf("%s = %s, want hex %s", tt.path, tok.Color.Hex(), tt.hex)
			}
		})
	}
}

// TestFullDerivation_WriteReadRoundTrip verifies that serializing the full
// derived TokenSet to YAML via WriteUniversal and reading it back via
// ReadUniversal produces an identical TokenSet.
func TestFullDerivation_WriteReadRoundTrip(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()
	original := d.Derive(pal)

	// Write to buffer.
	var buf bytes.Buffer
	if err := fileio.WriteUniversal(&buf, original); err != nil {
		t.Fatalf("WriteUniversal failed: %v", err)
	}

	// Read back from buffer.
	roundTripped, err := fileio.ReadUniversal(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadUniversal failed: %v", err)
	}

	// Compare counts.
	if original.Len() != roundTripped.Len() {
		t.Fatalf("token count mismatch: original=%d, round-tripped=%d",
			original.Len(), roundTripped.Len())
	}

	// Compare all paths and token values.
	origPaths := original.Paths()
	rtPaths := roundTripped.Paths()

	for i, path := range origPaths {
		if i >= len(rtPaths) || rtPaths[i] != path {
			t.Errorf("path mismatch at index %d: original=%q, round-tripped=%q",
				i, path, safeIndex(rtPaths, i))
			continue
		}

		origTok, _ := original.Get(path)
		rtTok, _ := roundTripped.Get(path)

		// Compare color.
		if origTok.Color.IsNone != rtTok.Color.IsNone {
			t.Errorf("token %q: IsNone mismatch: original=%v, round-tripped=%v",
				path, origTok.Color.IsNone, rtTok.Color.IsNone)
		} else if !origTok.Color.IsNone && !origTok.Color.Equal(rtTok.Color) {
			t.Errorf("token %q: color mismatch: original=%s, round-tripped=%s",
				path, origTok.Color.Hex(), rtTok.Color.Hex())
		}

		// Compare style flags.
		if origTok.Bold != rtTok.Bold {
			t.Errorf("token %q: Bold mismatch: %v vs %v", path, origTok.Bold, rtTok.Bold)
		}
		if origTok.Italic != rtTok.Italic {
			t.Errorf("token %q: Italic mismatch: %v vs %v", path, origTok.Italic, rtTok.Italic)
		}
		if origTok.Underline != rtTok.Underline {
			t.Errorf("token %q: Underline mismatch: %v vs %v", path, origTok.Underline, rtTok.Underline)
		}
		if origTok.Undercurl != rtTok.Undercurl {
			t.Errorf("token %q: Undercurl mismatch: %v vs %v", path, origTok.Undercurl, rtTok.Undercurl)
		}
		if origTok.Strikethrough != rtTok.Strikethrough {
			t.Errorf("token %q: Strikethrough mismatch: %v vs %v", path, origTok.Strikethrough, rtTok.Strikethrough)
		}
	}
}

// safeIndex returns the string at index i or "<missing>" if out of bounds.
func safeIndex(s []string, i int) string {
	if i >= len(s) {
		return "<missing>"
	}
	return s[i]
}

// TestFullDerivation_Deterministic verifies that two independent derivations
// from the same palette produce byte-identical TokenSets.
func TestFullDerivation_Deterministic(t *testing.T) {
	pal := tokyoNightDarkPalette(t)
	d := deriver.New()

	ts1 := d.Derive(pal)
	ts2 := d.Derive(pal)

	// Same count.
	if ts1.Len() != ts2.Len() {
		t.Fatalf("token count differs: %d vs %d", ts1.Len(), ts2.Len())
	}

	// Same paths.
	paths1 := ts1.Paths()
	paths2 := ts2.Paths()
	for i := range paths1 {
		if paths1[i] != paths2[i] {
			t.Fatalf("path mismatch at index %d: %q vs %q", i, paths1[i], paths2[i])
		}
	}

	// Same tokens.
	for _, path := range paths1 {
		tok1, _ := ts1.Get(path)
		tok2, _ := ts2.Get(path)

		if tok1.Color.IsNone != tok2.Color.IsNone {
			t.Errorf("token %q: IsNone differs: %v vs %v", path, tok1.Color.IsNone, tok2.Color.IsNone)
		}
		if !tok1.Color.IsNone && !tok1.Color.Equal(tok2.Color) {
			t.Errorf("token %q: color differs: %s vs %s", path, tok1.Color.Hex(), tok2.Color.Hex())
		}
		if tok1.Bold != tok2.Bold {
			t.Errorf("token %q: Bold differs: %v vs %v", path, tok1.Bold, tok2.Bold)
		}
		if tok1.Italic != tok2.Italic {
			t.Errorf("token %q: Italic differs: %v vs %v", path, tok1.Italic, tok2.Italic)
		}
		if tok1.Underline != tok2.Underline {
			t.Errorf("token %q: Underline differs: %v vs %v", path, tok1.Underline, tok2.Underline)
		}
		if tok1.Undercurl != tok2.Undercurl {
			t.Errorf("token %q: Undercurl differs: %v vs %v", path, tok1.Undercurl, tok2.Undercurl)
		}
		if tok1.Strikethrough != tok2.Strikethrough {
			t.Errorf("token %q: Strikethrough differs: %v vs %v", path, tok1.Strikethrough, tok2.Strikethrough)
		}
	}

	// Additionally verify byte-identical serialization.
	var buf1, buf2 bytes.Buffer
	if err := fileio.WriteUniversal(&buf1, ts1); err != nil {
		t.Fatalf("WriteUniversal ts1 failed: %v", err)
	}
	if err := fileio.WriteUniversal(&buf2, ts2); err != nil {
		t.Fatalf("WriteUniversal ts2 failed: %v", err)
	}
	if buf1.String() != buf2.String() {
		t.Error("serialized YAML output differs between two derivations")
	}
}

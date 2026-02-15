package application_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/deriver"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/application"
	"github.com/curtbushko/flair/internal/domain"
)

// --- Stub implementations of ports ---

// stubParser is a test stub for ports.PaletteParser.
type stubParser struct {
	palette *domain.Palette
	err     error
}

func (s *stubParser) Parse(_ io.Reader) (*domain.Palette, error) {
	return s.palette, s.err
}

// stubDeriver is a test stub for ports.TokenDeriver.
type stubDeriver struct {
	tokenSet *domain.TokenSet
}

func (s *stubDeriver) Derive(_ *domain.Palette) *domain.TokenSet {
	return s.tokenSet
}

// --- Helpers ---

func makePalette(t *testing.T, name, variant string) *domain.Palette {
	t.Helper()
	colors := map[string]string{
		"base00": "1a1b26", "base01": "1f2335", "base02": "292e42", "base03": "565f89",
		"base04": "a9b1d6", "base05": "c0caf5", "base06": "c0caf5", "base07": "c8d3f5",
		"base08": "f7768e", "base09": "ff9e64", "base0A": "e0af68", "base0B": "9ece6a",
		"base0C": "7dcfff", "base0D": "7aa2f7", "base0E": "bb9af7", "base0F": "db4b4b",
		"base10": "16161e", "base11": "101014", "base12": "ff899d", "base13": "e9c582",
		"base14": "afd67a", "base15": "97d8f8", "base16": "8db6fa", "base17": "c8acf8",
	}
	pal, err := domain.NewPalette(name, "test-author", variant, "base24", colors)
	if err != nil {
		t.Fatalf("makePalette: %v", err)
	}
	return pal
}

// --- Tests ---

func TestDeriveThemeUseCase_Execute(t *testing.T) {
	pal := makePalette(t, "Test Theme", "dark")

	ts := domain.NewTokenSet()
	ts.Set("surface.background", domain.Token{Color: pal.Base(0x00)})
	ts.Set("text.primary", domain.Token{Color: pal.Base(0x05)})

	uc := application.NewDeriveThemeUseCase(
		&stubParser{palette: pal},
		&stubDeriver{tokenSet: ts},
	)

	theme, err := uc.Execute(bytes.NewReader([]byte("dummy")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if theme.Name != "Test Theme" {
		t.Errorf("Name = %q, want %q", theme.Name, "Test Theme")
	}
	if theme.Variant != "dark" {
		t.Errorf("Variant = %q, want %q", theme.Variant, "dark")
	}
	if theme.Palette != pal {
		t.Error("Palette pointer does not match")
	}
	if theme.Tokens != ts {
		t.Error("Tokens pointer does not match")
	}

	tok, ok := theme.Token("surface.background")
	if !ok {
		t.Fatal("expected surface.background token to exist")
	}
	if tok.Color != pal.Base(0x00) {
		t.Errorf("surface.background color = %v, want %v", tok.Color, pal.Base(0x00))
	}
}

func TestDeriveThemeUseCase_ParseError(t *testing.T) {
	parseErr := errors.New("bad yaml")
	uc := application.NewDeriveThemeUseCase(
		&stubParser{err: parseErr},
		&stubDeriver{tokenSet: domain.NewTokenSet()},
	)

	_, err := uc.Execute(bytes.NewReader([]byte("bad")))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, parseErr) {
		t.Errorf("error = %v, want wrapped %v", err, parseErr)
	}
}

func TestDeriveThemeUseCase_Metadata(t *testing.T) {
	pal := makePalette(t, "Tokyo Night Dark", "dark")
	ts := domain.NewTokenSet()

	uc := application.NewDeriveThemeUseCase(
		&stubParser{palette: pal},
		&stubDeriver{tokenSet: ts},
	)

	theme, err := uc.Execute(bytes.NewReader([]byte("dummy")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if theme.Name != "Tokyo Night Dark" {
		t.Errorf("Name = %q, want %q", theme.Name, "Tokyo Night Dark")
	}
	if theme.Variant != "dark" {
		t.Errorf("Variant = %q, want %q", theme.Variant, "dark")
	}
}

func TestDeriveThemeUseCase_Integration(t *testing.T) {
	yamlBytes, err := os.ReadFile("../../testdata/tokyo-night-dark.yaml")
	if err != nil {
		t.Fatalf("read test fixture: %v", err)
	}

	parser := yamlparser.NewParser()
	deriv := deriver.New()

	uc := application.NewDeriveThemeUseCase(parser, deriv)

	theme, err := uc.Execute(bytes.NewReader(yamlBytes))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if theme.Name != "Tokyo Night Dark" {
		t.Errorf("Name = %q, want %q", theme.Name, "Tokyo Night Dark")
	}
	if theme.Tokens.Len() < 87 {
		t.Errorf("Tokens.Len() = %d, want >= 87", theme.Tokens.Len())
	}
}

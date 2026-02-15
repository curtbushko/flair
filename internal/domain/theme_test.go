package domain_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestResolvedTheme_Fields(t *testing.T) {
	pal := &domain.Palette{
		Name:    "tokyonight",
		Variant: "storm",
	}
	ts := domain.NewTokenSet()
	ts.Set("surface.background", domain.Token{
		Color: domain.Color{R: 26, G: 27, B: 38},
	})

	rt := domain.ResolvedTheme{
		Name:    "tokyonight",
		Variant: "storm",
		Palette: pal,
		Tokens:  ts,
	}

	if rt.Name != "tokyonight" {
		t.Errorf("Name = %q, want %q", rt.Name, "tokyonight")
	}
	if rt.Variant != "storm" {
		t.Errorf("Variant = %q, want %q", rt.Variant, "storm")
	}
	if rt.Palette != pal {
		t.Error("Palette does not match construction value")
	}
	if rt.Tokens != ts {
		t.Error("Tokens does not match construction value")
	}
}

func TestResolvedTheme_Token_Found(t *testing.T) {
	ts := domain.NewTokenSet()
	expected := domain.Token{
		Color: domain.Color{R: 26, G: 27, B: 38},
		Bold:  true,
	}
	ts.Set("surface.background", expected)

	rt := domain.ResolvedTheme{
		Name:    "tokyonight",
		Variant: "storm",
		Palette: &domain.Palette{},
		Tokens:  ts,
	}

	got, ok := rt.Token("surface.background")
	if !ok {
		t.Fatal("Token() ok = false, want true for existing path")
	}
	if got.Color.R != expected.Color.R || got.Color.G != expected.Color.G || got.Color.B != expected.Color.B {
		t.Errorf("Token() Color = {R:%d, G:%d, B:%d}, want {R:%d, G:%d, B:%d}",
			got.Color.R, got.Color.G, got.Color.B,
			expected.Color.R, expected.Color.G, expected.Color.B)
	}
	if got.Bold != expected.Bold {
		t.Errorf("Token() Bold = %v, want %v", got.Bold, expected.Bold)
	}
}

func TestResolvedTheme_Token_NotFound(t *testing.T) {
	ts := domain.NewTokenSet()

	rt := domain.ResolvedTheme{
		Name:    "tokyonight",
		Variant: "storm",
		Palette: &domain.Palette{},
		Tokens:  ts,
	}

	got, ok := rt.Token("nonexistent")
	if ok {
		t.Error("Token() ok = true, want false for missing path")
	}
	if (got != domain.Token{}) {
		t.Errorf("Token() = %+v, want zero Token for missing path", got)
	}
}

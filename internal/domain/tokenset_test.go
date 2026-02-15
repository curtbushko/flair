package domain_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestToken_HasStyle_NoFlags(t *testing.T) {
	tok := domain.Token{
		Color: domain.Color{R: 255, G: 0, B: 0},
	}

	if tok.HasStyle() {
		t.Error("HasStyle() = true, want false when no flags set")
	}
}

func TestToken_HasStyle_Bold(t *testing.T) {
	tok := domain.Token{
		Bold: true,
	}

	if !tok.HasStyle() {
		t.Error("HasStyle() = false, want true when Bold is set")
	}
}

func TestToken_HasStyle_AllFlags(t *testing.T) {
	tests := []struct {
		name  string
		token domain.Token
	}{
		{"Bold", domain.Token{Bold: true}},
		{"Italic", domain.Token{Italic: true}},
		{"Underline", domain.Token{Underline: true}},
		{"Undercurl", domain.Token{Undercurl: true}},
		{"Strikethrough", domain.Token{Strikethrough: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.token.HasStyle() {
				t.Errorf("HasStyle() = false, want true when %s is set", tt.name)
			}
		})
	}
}

func TestNewTokenSet_Empty(t *testing.T) {
	ts := domain.NewTokenSet()

	if ts.Len() != 0 {
		t.Errorf("Len() = %d, want 0 for new empty TokenSet", ts.Len())
	}
}

func TestTokenSet_SetGet(t *testing.T) {
	ts := domain.NewTokenSet()
	tok := domain.Token{
		Color: domain.Color{R: 26, G: 27, B: 38},
		Bold:  true,
	}

	ts.Set("surface.background", tok)
	got, ok := ts.Get("surface.background")

	if !ok {
		t.Fatal("Get() ok = false, want true after Set")
	}
	if got.Color.R != tok.Color.R || got.Color.G != tok.Color.G || got.Color.B != tok.Color.B {
		t.Errorf("Get() Color = {R:%d, G:%d, B:%d}, want {R:%d, G:%d, B:%d}",
			got.Color.R, got.Color.G, got.Color.B,
			tok.Color.R, tok.Color.G, tok.Color.B)
	}
	if got.Bold != tok.Bold {
		t.Errorf("Get() Bold = %v, want %v", got.Bold, tok.Bold)
	}
}

func TestTokenSet_Get_Missing(t *testing.T) {
	ts := domain.NewTokenSet()

	got, ok := ts.Get("nonexistent")

	if ok {
		t.Error("Get() ok = true, want false for missing path")
	}
	if (got != domain.Token{}) {
		t.Errorf("Get() = %+v, want zero Token for missing path", got)
	}
}

func TestTokenSet_MustGet_Panics(t *testing.T) {
	ts := domain.NewTokenSet()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet() did not panic for missing path")
		}
	}()

	ts.MustGet("nonexistent")
}

func TestTokenSet_Paths_Sorted(t *testing.T) {
	ts := domain.NewTokenSet()
	ts.Set("z.b", domain.Token{})
	ts.Set("a.c", domain.Token{})
	ts.Set("m.d", domain.Token{})

	paths := ts.Paths()

	expected := []string{"a.c", "m.d", "z.b"}
	if len(paths) != len(expected) {
		t.Fatalf("Paths() len = %d, want %d", len(paths), len(expected))
	}
	for i, p := range paths {
		if p != expected[i] {
			t.Errorf("Paths()[%d] = %q, want %q", i, p, expected[i])
		}
	}
}

func TestTokenSet_Len(t *testing.T) {
	ts := domain.NewTokenSet()
	ts.Set("a", domain.Token{})
	ts.Set("b", domain.Token{})
	ts.Set("c", domain.Token{})

	if ts.Len() != 3 {
		t.Errorf("Len() = %d, want 3", ts.Len())
	}
}

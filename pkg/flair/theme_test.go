package flair_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

func TestTheme_Fields(t *testing.T) {
	// Create a theme with specific fields
	colors := map[string]flair.Color{
		"background": {R: 26, G: 27, B: 38},
		"foreground": {R: 192, G: 202, B: 245},
	}

	theme := flair.NewTheme("tokyo-night", "storm", colors)

	// Test Name field
	if theme.Name() != "tokyo-night" {
		t.Errorf("Theme.Name() = %v, want %v", theme.Name(), "tokyo-night")
	}

	// Test Variant field
	if theme.Variant() != "storm" {
		t.Errorf("Theme.Variant() = %v, want %v", theme.Variant(), "storm")
	}
}

func TestTheme_HasColors(t *testing.T) {
	tests := []struct {
		name   string
		colors map[string]flair.Color
		want   bool
	}{
		{
			name:   "with colors",
			colors: map[string]flair.Color{"bg": {R: 0, G: 0, B: 0}},
			want:   true,
		},
		{
			name:   "empty map",
			colors: map[string]flair.Color{},
			want:   false,
		},
		{
			name:   "nil map",
			colors: nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := flair.NewTheme("test", "dark", tt.colors)
			got := theme.HasColors()
			if got != tt.want {
				t.Errorf("Theme.HasColors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTheme_Color(t *testing.T) {
	colors := map[string]flair.Color{
		"background": {R: 26, G: 27, B: 38},
		"foreground": {R: 192, G: 202, B: 245},
	}
	theme := flair.NewTheme("tokyo-night", "storm", colors)

	tests := []struct {
		name      string
		key       string
		wantColor flair.Color
		wantOK    bool
	}{
		{
			name:      "existing key",
			key:       "background",
			wantColor: flair.Color{R: 26, G: 27, B: 38},
			wantOK:    true,
		},
		{
			name:      "another existing key",
			key:       "foreground",
			wantColor: flair.Color{R: 192, G: 202, B: 245},
			wantOK:    true,
		},
		{
			name:      "non-existent key",
			key:       "accent",
			wantColor: flair.Color{},
			wantOK:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := theme.Color(tt.key)
			if ok != tt.wantOK {
				t.Errorf("Theme.Color(%q) ok = %v, want %v", tt.key, ok, tt.wantOK)
			}
			if !got.Equal(tt.wantColor) {
				t.Errorf("Theme.Color(%q) = %v, want %v", tt.key, got, tt.wantColor)
			}
		})
	}
}

func TestTheme_Colors(t *testing.T) {
	colors := map[string]flair.Color{
		"background": {R: 26, G: 27, B: 38},
		"foreground": {R: 192, G: 202, B: 245},
	}
	theme := flair.NewTheme("tokyo-night", "storm", colors)

	// Get all colors
	allColors := theme.Colors()

	// Verify count
	if len(allColors) != 2 {
		t.Errorf("Theme.Colors() returned %d colors, want 2", len(allColors))
	}

	// Verify contents
	if bg, ok := allColors["background"]; !ok || !bg.Equal(colors["background"]) {
		t.Errorf("Theme.Colors() missing or wrong background color")
	}
	if fg, ok := allColors["foreground"]; !ok || !fg.Equal(colors["foreground"]) {
		t.Errorf("Theme.Colors() missing or wrong foreground color")
	}

	// Verify it's a copy (modifying returned map doesn't affect theme)
	allColors["newkey"] = flair.Color{R: 1, G: 2, B: 3}
	_, ok := theme.Color("newkey")
	if ok {
		t.Error("Theme.Colors() should return a copy, but modification affected original")
	}
}

package domain_test

import (
	"errors"
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestTokenOverride_ColorOnly(t *testing.T) {
	// Arrange: Create a TokenOverride with only Color field set
	color, err := domain.ParseHex("#ff0000")
	if err != nil {
		t.Fatalf("ParseHex() unexpected error: %v", err)
	}
	override := domain.TokenOverride{
		Color: &color,
	}

	// Act & Assert: HasColor() should return true, HasStyle() should return false
	if !override.HasColor() {
		t.Error("HasColor() = false, want true")
	}
	if override.HasStyle() {
		t.Error("HasStyle() = true, want false")
	}
}

func TestTokenOverride_WithStyles(t *testing.T) {
	// Arrange: Create a TokenOverride with Bold=true, Italic=true
	override := domain.TokenOverride{
		Bold:   true,
		Italic: true,
	}

	// Act & Assert: HasStyle() should return true
	if !override.HasStyle() {
		t.Error("HasStyle() = false, want true")
	}
	// HasColor should return false since no color is set
	if override.HasColor() {
		t.Error("HasColor() = true, want false")
	}
}

func TestTokenOverride_AllStyleFlags(t *testing.T) {
	tests := []struct {
		name     string
		override domain.TokenOverride
		want     bool
	}{
		{"bold only", domain.TokenOverride{Bold: true}, true},
		{"italic only", domain.TokenOverride{Italic: true}, true},
		{"underline only", domain.TokenOverride{Underline: true}, true},
		{"undercurl only", domain.TokenOverride{Undercurl: true}, true},
		{"strikethrough only", domain.TokenOverride{Strikethrough: true}, true},
		{"no styles", domain.TokenOverride{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.override.HasStyle()
			if got != tt.want {
				t.Errorf("HasStyle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenOverride_ValidateColor(t *testing.T) {
	// This test verifies that Validate() returns a *domain.ParseError for invalid hex colors
	// Note: Since Color is already a *domain.Color (parsed), validation happens at parse time.
	// We test the NewTokenOverride constructor or a Validate method if one exists.

	tests := []struct {
		name      string
		colorHex  string
		wantError bool
	}{
		{"valid 6-digit hex", "#ff0000", false},
		{"valid 3-digit hex", "#f00", false},
		{"invalid hex chars", "#gggggg", true},
		{"invalid length", "#ff00", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			override, err := domain.NewTokenOverride(tt.colorHex, false, false, false, false, false)
			if (err != nil) != tt.wantError {
				t.Errorf("NewTokenOverride() error = %v, wantError %v", err, tt.wantError)
			}
			if tt.wantError {
				var pe *domain.ParseError
				if !errors.As(err, &pe) {
					t.Errorf("error type = %T, want *domain.ParseError", err)
				}
			} else if override == nil {
				t.Error("NewTokenOverride() returned nil, want non-nil TokenOverride")
			}
		})
	}
}

func TestTokenOverride_ApplyToToken(t *testing.T) {
	// Arrange: A base Token and a TokenOverride with color and italic
	baseColor, _ := domain.ParseHex("#000000")
	baseToken := domain.Token{
		Color:     baseColor,
		Bold:      true,
		Italic:    false,
		Underline: true,
	}

	overrideColor, _ := domain.ParseHex("#ff0000")
	override := domain.TokenOverride{
		Color:  &overrideColor,
		Italic: true,
	}

	// Act: Apply the override to the token
	result := override.Apply(baseToken)

	// Assert: Returns new Token with overridden color and italic, preserving other styles
	// Color should be overridden
	if !result.Color.Equal(overrideColor) {
		t.Errorf("result.Color = %s, want %s", result.Color.Hex(), overrideColor.Hex())
	}
	// Italic should be overridden to true
	if !result.Italic {
		t.Error("result.Italic = false, want true")
	}
	// Bold should be preserved from base token
	if !result.Bold {
		t.Error("result.Bold = false, want true (preserved from base)")
	}
	// Underline should be preserved from base token
	if !result.Underline {
		t.Error("result.Underline = false, want true (preserved from base)")
	}
}

func TestTokenOverride_ApplyPreservesBaseWhenNoOverride(t *testing.T) {
	// Arrange: A base Token with all styles set, and an empty override
	baseColor, _ := domain.ParseHex("#123456")
	baseToken := domain.Token{
		Color:         baseColor,
		Bold:          true,
		Italic:        true,
		Underline:     true,
		Undercurl:     true,
		Strikethrough: true,
	}

	override := domain.TokenOverride{} // Empty override

	// Act: Apply the empty override
	result := override.Apply(baseToken)

	// Assert: All values should be preserved from base
	if !result.Color.Equal(baseColor) {
		t.Errorf("result.Color = %s, want %s", result.Color.Hex(), baseColor.Hex())
	}
	if !result.Bold {
		t.Error("result.Bold = false, want true")
	}
	if !result.Italic {
		t.Error("result.Italic = false, want true")
	}
	if !result.Underline {
		t.Error("result.Underline = false, want true")
	}
	if !result.Undercurl {
		t.Error("result.Undercurl = false, want true")
	}
	if !result.Strikethrough {
		t.Error("result.Strikethrough = false, want true")
	}
}

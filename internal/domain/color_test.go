package domain_test

import (
	"errors"
	"math"
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestParseHex_SixDigitWithHash(t *testing.T) {
	c, err := domain.ParseHex("#7aa2f7")
	if err != nil {
		t.Fatalf("ParseHex(#7aa2f7) unexpected error: %v", err)
	}
	if c.R != 122 || c.G != 162 || c.B != 247 {
		t.Errorf("ParseHex(#7aa2f7) = {R:%d, G:%d, B:%d}, want {R:122, G:162, B:247}", c.R, c.G, c.B)
	}
}

func TestParseHex_WithoutHash(t *testing.T) {
	c, err := domain.ParseHex("7aa2f7")
	if err != nil {
		t.Fatalf("ParseHex(7aa2f7) unexpected error: %v", err)
	}
	if c.R != 122 || c.G != 162 || c.B != 247 {
		t.Errorf("ParseHex(7aa2f7) = {R:%d, G:%d, B:%d}, want {R:122, G:162, B:247}", c.R, c.G, c.B)
	}
}

func TestParseHex_ThreeDigit(t *testing.T) {
	c, err := domain.ParseHex("#f00")
	if err != nil {
		t.Fatalf("ParseHex(#f00) unexpected error: %v", err)
	}
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("ParseHex(#f00) = {R:%d, G:%d, B:%d}, want {R:255, G:0, B:0}", c.R, c.G, c.B)
	}
}

func TestParseHex_ThreeDigitWithoutHash(t *testing.T) {
	c, err := domain.ParseHex("f00")
	if err != nil {
		t.Fatalf("ParseHex(f00) unexpected error: %v", err)
	}
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("ParseHex(f00) = {R:%d, G:%d, B:%d}, want {R:255, G:0, B:0}", c.R, c.G, c.B)
	}
}

func TestParseHex_InvalidChars(t *testing.T) {
	_, err := domain.ParseHex("#zzzzzz")
	if err == nil {
		t.Fatal("ParseHex(#zzzzzz) expected error, got nil")
	}
	var pe *domain.ParseError
	if !errors.As(err, &pe) {
		t.Errorf("ParseHex(#zzzzzz) error type = %T, want *domain.ParseError", err)
	}
}

func TestParseHex_WrongLength(t *testing.T) {
	_, err := domain.ParseHex("#12345")
	if err == nil {
		t.Fatal("ParseHex(#12345) expected error, got nil")
	}
	var pe *domain.ParseError
	if !errors.As(err, &pe) {
		t.Errorf("ParseHex(#12345) error type = %T, want *domain.ParseError", err)
	}
}

func TestParseHex_EmptyString(t *testing.T) {
	_, err := domain.ParseHex("")
	if err == nil {
		t.Fatal("ParseHex('') expected error, got nil")
	}
	var pe *domain.ParseError
	if !errors.As(err, &pe) {
		t.Errorf("ParseHex('') error type = %T, want *domain.ParseError", err)
	}
}

func TestParseHex_CaseInsensitive(t *testing.T) {
	c, err := domain.ParseHex("#7AA2F7")
	if err != nil {
		t.Fatalf("ParseHex(#7AA2F7) unexpected error: %v", err)
	}
	if c.R != 122 || c.G != 162 || c.B != 247 {
		t.Errorf("ParseHex(#7AA2F7) = {R:%d, G:%d, B:%d}, want {R:122, G:162, B:247}", c.R, c.G, c.B)
	}
}

func TestColor_Hex(t *testing.T) {
	c := domain.Color{R: 122, G: 162, B: 247}
	got := c.Hex()
	want := "#7aa2f7"
	if got != want {
		t.Errorf("Color{122,162,247}.Hex() = %q, want %q", got, want)
	}
}

func TestColor_Hex_Black(t *testing.T) {
	c := domain.Color{R: 0, G: 0, B: 0}
	got := c.Hex()
	want := "#000000"
	if got != want {
		t.Errorf("Color{0,0,0}.Hex() = %q, want %q", got, want)
	}
}

func TestColor_Hex_White(t *testing.T) {
	c := domain.Color{R: 255, G: 255, B: 255}
	got := c.Hex()
	want := "#ffffff"
	if got != want {
		t.Errorf("Color{255,255,255}.Hex() = %q, want %q", got, want)
	}
}

func TestNoneColor(t *testing.T) {
	c := domain.NoneColor()
	if !c.IsNone {
		t.Error("NoneColor().IsNone = false, want true")
	}
}

func TestNoneColor_RegularColorIsNotNone(t *testing.T) {
	c, err := domain.ParseHex("#646464")
	if err != nil {
		t.Fatalf("ParseHex(#646464) unexpected error: %v", err)
	}
	if c.IsNone {
		t.Error("Regular color should have IsNone = false")
	}
}

func TestColor_Equal(t *testing.T) {
	tests := []struct {
		name string
		a    domain.Color
		b    domain.Color
		want bool
	}{
		{
			name: "identical colors",
			a:    domain.Color{R: 122, G: 162, B: 247},
			b:    domain.Color{R: 122, G: 162, B: 247},
			want: true,
		},
		{
			name: "different red",
			a:    domain.Color{R: 100, G: 162, B: 247},
			b:    domain.Color{R: 122, G: 162, B: 247},
			want: false,
		},
		{
			name: "different green",
			a:    domain.Color{R: 122, G: 100, B: 247},
			b:    domain.Color{R: 122, G: 162, B: 247},
			want: false,
		},
		{
			name: "different blue",
			a:    domain.Color{R: 122, G: 162, B: 200},
			b:    domain.Color{R: 122, G: 162, B: 247},
			want: false,
		},
		{
			name: "ignores IsNone for equality",
			a:    domain.Color{R: 0, G: 0, B: 0, IsNone: true},
			b:    domain.Color{R: 0, G: 0, B: 0, IsNone: false},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Equal(tt.b)
			if got != tt.want {
				t.Errorf("Color.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_ToHSL_Red(t *testing.T) {
	c := domain.Color{R: 255, G: 0, B: 0}
	hsl := c.ToHSL()
	if math.Abs(hsl.H-0) > 1.0 {
		t.Errorf("Red ToHSL().H = %f, want approx 0", hsl.H)
	}
	if math.Abs(hsl.S-1.0) > 0.01 {
		t.Errorf("Red ToHSL().S = %f, want approx 1.0", hsl.S)
	}
	if math.Abs(hsl.L-0.5) > 0.01 {
		t.Errorf("Red ToHSL().L = %f, want approx 0.5", hsl.L)
	}
}

func TestColor_ToHSL_Green(t *testing.T) {
	c := domain.Color{R: 0, G: 255, B: 0}
	hsl := c.ToHSL()
	if math.Abs(hsl.H-120) > 1.0 {
		t.Errorf("Green ToHSL().H = %f, want approx 120", hsl.H)
	}
	if math.Abs(hsl.S-1.0) > 0.01 {
		t.Errorf("Green ToHSL().S = %f, want approx 1.0", hsl.S)
	}
	if math.Abs(hsl.L-0.5) > 0.01 {
		t.Errorf("Green ToHSL().L = %f, want approx 0.5", hsl.L)
	}
}

func TestColor_ToHSL_Blue(t *testing.T) {
	c := domain.Color{R: 0, G: 0, B: 255}
	hsl := c.ToHSL()
	if math.Abs(hsl.H-240) > 1.0 {
		t.Errorf("Blue ToHSL().H = %f, want approx 240", hsl.H)
	}
	if math.Abs(hsl.S-1.0) > 0.01 {
		t.Errorf("Blue ToHSL().S = %f, want approx 1.0", hsl.S)
	}
	if math.Abs(hsl.L-0.5) > 0.01 {
		t.Errorf("Blue ToHSL().L = %f, want approx 0.5", hsl.L)
	}
}

func TestColor_ToHSL_White(t *testing.T) {
	c := domain.Color{R: 255, G: 255, B: 255}
	hsl := c.ToHSL()
	if math.Abs(hsl.S) > 0.01 {
		t.Errorf("White ToHSL().S = %f, want approx 0", hsl.S)
	}
	if math.Abs(hsl.L-1.0) > 0.01 {
		t.Errorf("White ToHSL().L = %f, want approx 1.0", hsl.L)
	}
}

func TestColor_ToHSL_Black(t *testing.T) {
	c := domain.Color{R: 0, G: 0, B: 0}
	hsl := c.ToHSL()
	if math.Abs(hsl.S) > 0.01 {
		t.Errorf("Black ToHSL().S = %f, want approx 0", hsl.S)
	}
	if math.Abs(hsl.L) > 0.01 {
		t.Errorf("Black ToHSL().L = %f, want approx 0", hsl.L)
	}
}

func TestColor_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		hex  string
	}{
		{"bb9af7", "#bb9af7"},
		{"7aa2f7", "#7aa2f7"},
		{"pure red", "#ff0000"},
		{"pure green", "#00ff00"},
		{"pure blue", "#0000ff"},
		{"white", "#ffffff"},
		{"black", "#000000"},
		{"mid gray", "#808080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original, err := domain.ParseHex(tt.hex)
			if err != nil {
				t.Fatalf("ParseHex(%s) unexpected error: %v", tt.hex, err)
			}
			roundTripped := original.ToHSL().ToRGB()
			got := roundTripped.Hex()
			if got != tt.hex {
				t.Errorf("Round-trip %s -> ToHSL -> ToRGB -> Hex = %q, want %q", tt.hex, got, tt.hex)
			}
		})
	}
}

func TestHSL_ToRGB_PureRed(t *testing.T) {
	hsl := domain.HSL{H: 0, S: 1.0, L: 0.5}
	c := hsl.ToRGB()
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("HSL{0,1,0.5}.ToRGB() = {R:%d, G:%d, B:%d}, want {R:255, G:0, B:0}", c.R, c.G, c.B)
	}
}

func TestHSL_ToRGB_PureGreen(t *testing.T) {
	hsl := domain.HSL{H: 120, S: 1.0, L: 0.5}
	c := hsl.ToRGB()
	if c.R != 0 || c.G != 255 || c.B != 0 {
		t.Errorf("HSL{120,1,0.5}.ToRGB() = {R:%d, G:%d, B:%d}, want {R:0, G:255, B:0}", c.R, c.G, c.B)
	}
}

func TestHSL_ToRGB_PureBlue(t *testing.T) {
	hsl := domain.HSL{H: 240, S: 1.0, L: 0.5}
	c := hsl.ToRGB()
	if c.R != 0 || c.G != 0 || c.B != 255 {
		t.Errorf("HSL{240,1,0.5}.ToRGB() = {R:%d, G:%d, B:%d}, want {R:0, G:0, B:255}", c.R, c.G, c.B)
	}
}

func TestColor_Luminance_White(t *testing.T) {
	c := domain.Color{R: 255, G: 255, B: 255}
	got := c.Luminance()
	if math.Abs(got-1.0) > 0.01 {
		t.Errorf("White Luminance() = %f, want approx 1.0", got)
	}
}

func TestColor_Luminance_Black(t *testing.T) {
	c := domain.Color{R: 0, G: 0, B: 0}
	got := c.Luminance()
	if math.Abs(got) > 0.01 {
		t.Errorf("Black Luminance() = %f, want approx 0.0", got)
	}
}

func TestColor_Luminance_MidRange(t *testing.T) {
	c := domain.Color{R: 128, G: 128, B: 128}
	got := c.Luminance()
	// Mid-gray luminance should be between 0 and 1
	if got <= 0 || got >= 1.0 {
		t.Errorf("Mid-gray Luminance() = %f, want between 0 and 1", got)
	}
	// WCAG relative luminance for #808080 is approximately 0.2159
	if math.Abs(got-0.2159) > 0.01 {
		t.Errorf("Mid-gray Luminance() = %f, want approx 0.2159", got)
	}
}

func TestColor_Luminance_WCAG_Formula(t *testing.T) {
	// Test that luminance follows the WCAG 2.1 formula:
	// L = 0.2126 * R + 0.7152 * G + 0.0722 * B
	// where each channel is linearized from sRGB.
	// Pure red should have luminance ~ 0.2126
	red := domain.Color{R: 255, G: 0, B: 0}
	got := red.Luminance()
	if math.Abs(got-0.2126) > 0.01 {
		t.Errorf("Pure red Luminance() = %f, want approx 0.2126", got)
	}
}

package domain_test

import (
	"math"
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestBlend_50Percent(t *testing.T) {
	black := domain.Color{R: 0, G: 0, B: 0}
	white := domain.Color{R: 255, G: 255, B: 255}

	got := domain.Blend(black, white, 0.5)

	// Linear interpolation: 0 + 0.5*(255-0) = 127.5, rounded to 128
	if got.R != 128 || got.G != 128 || got.B != 128 {
		t.Errorf("Blend(black, white, 0.5) = {R:%d, G:%d, B:%d}, want {R:128, G:128, B:128}",
			got.R, got.G, got.B)
	}
}

func TestBlend_Zero(t *testing.T) {
	a := domain.Color{R: 100, G: 150, B: 200}
	b := domain.Color{R: 50, G: 75, B: 100}

	got := domain.Blend(a, b, 0.0)

	if !got.Equal(a) {
		t.Errorf("Blend(a, b, 0.0) = {R:%d, G:%d, B:%d}, want {R:%d, G:%d, B:%d}",
			got.R, got.G, got.B, a.R, a.G, a.B)
	}
}

func TestBlend_One(t *testing.T) {
	a := domain.Color{R: 100, G: 150, B: 200}
	b := domain.Color{R: 50, G: 75, B: 100}

	got := domain.Blend(a, b, 1.0)

	if !got.Equal(b) {
		t.Errorf("Blend(a, b, 1.0) = {R:%d, G:%d, B:%d}, want {R:%d, G:%d, B:%d}",
			got.R, got.G, got.B, b.R, b.G, b.B)
	}
}

func TestBlendBg_EquivalentToBlend(t *testing.T) {
	fg, _ := domain.ParseHex("#7aa2f7")
	bg, _ := domain.ParseHex("#1a1b26")

	blendBgResult := domain.BlendBg(fg, bg, 0.25)
	blendResult := domain.Blend(bg, fg, 0.25)

	if !blendBgResult.Equal(blendResult) {
		t.Errorf("BlendBg(fg, bg, 0.25) = %s, Blend(bg, fg, 0.25) = %s; expected equal",
			blendBgResult.Hex(), blendResult.Hex())
	}
}

func TestLighten(t *testing.T) {
	// Mid-range color
	color := domain.Color{R: 122, G: 162, B: 247}
	originalLum := color.Luminance()

	got := domain.Lighten(color, 0.1)
	gotLum := got.Luminance()

	if gotLum <= originalLum {
		t.Errorf("Lighten: result luminance %f should be higher than original %f", gotLum, originalLum)
	}
}

func TestLighten_ClampMax(t *testing.T) {
	// White should stay white when lightened
	white := domain.Color{R: 255, G: 255, B: 255}
	got := domain.Lighten(white, 0.5)

	if got.R != 255 || got.G != 255 || got.B != 255 {
		t.Errorf("Lighten(white, 0.5) = {R:%d, G:%d, B:%d}, want white",
			got.R, got.G, got.B)
	}
}

func TestDarken(t *testing.T) {
	// Mid-range color
	color := domain.Color{R: 122, G: 162, B: 247}
	originalLum := color.Luminance()

	got := domain.Darken(color, 0.1)
	gotLum := got.Luminance()

	if gotLum >= originalLum {
		t.Errorf("Darken: result luminance %f should be lower than original %f", gotLum, originalLum)
	}
}

func TestDarken_ClampMin(t *testing.T) {
	// Black should stay black when darkened
	black := domain.Color{R: 0, G: 0, B: 0}
	got := domain.Darken(black, 0.5)

	if got.R != 0 || got.G != 0 || got.B != 0 {
		t.Errorf("Darken(black, 0.5) = {R:%d, G:%d, B:%d}, want black",
			got.R, got.G, got.B)
	}
}

func TestDesaturate(t *testing.T) {
	// Saturated red color
	color := domain.Color{R: 255, G: 0, B: 0}
	originalHSL := color.ToHSL()

	got := domain.Desaturate(color, 0.5)
	gotHSL := got.ToHSL()

	// Desaturate(c, 0.5) should produce S = original_S * (1 - 0.5) = original_S * 0.5
	wantS := originalHSL.S * 0.5
	if math.Abs(gotHSL.S-wantS) > 0.02 {
		t.Errorf("Desaturate: got S=%f, want approx %f", gotHSL.S, wantS)
	}
}

func TestDesaturate_Full(t *testing.T) {
	// Full desaturation should produce a grayscale color
	color := domain.Color{R: 255, G: 0, B: 0}

	got := domain.Desaturate(color, 1.0)
	gotHSL := got.ToHSL()

	if math.Abs(gotHSL.S) > 0.02 {
		t.Errorf("Desaturate(color, 1.0): got S=%f, want approx 0", gotHSL.S)
	}
}

func TestShiftHue(t *testing.T) {
	// Pure red (H=0)
	color := domain.Color{R: 255, G: 0, B: 0}

	got := domain.ShiftHue(color, 120.0)
	gotHSL := got.ToHSL()

	// Shifting red by 120 degrees should give approximately green (H=120)
	if math.Abs(gotHSL.H-120.0) > 1.0 {
		t.Errorf("ShiftHue(red, 120): got H=%f, want approx 120", gotHSL.H)
	}
}

func TestShiftHue_Wrap(t *testing.T) {
	// Color with H=300 (magenta-ish)
	hsl := domain.HSL{H: 300, S: 1.0, L: 0.5}
	color := hsl.ToRGB()

	got := domain.ShiftHue(color, 120.0)
	gotHSL := got.ToHSL()

	// (300 + 120) % 360 = 60
	wantH := 60.0
	if math.Abs(gotHSL.H-wantH) > 1.0 {
		t.Errorf("ShiftHue(H=300, 120): got H=%f, want approx %f", gotHSL.H, wantH)
	}
}

func TestShiftHue_Negative(t *testing.T) {
	// Pure green (H=120)
	color := domain.Color{R: 0, G: 255, B: 0}

	got := domain.ShiftHue(color, -120.0)
	gotHSL := got.ToHSL()

	// Shifting green by -120 should give red (H=0)
	if gotHSL.H > 1.0 && gotHSL.H < 359.0 {
		t.Errorf("ShiftHue(green, -120): got H=%f, want approx 0 or 360", gotHSL.H)
	}
}

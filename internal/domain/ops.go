package domain

import "math"

// Blend performs linear RGB interpolation between two colors.
// t=0.0 returns a, t=1.0 returns b.
func Blend(a, b Color, t float64) Color {
	return Color{
		R: clampUint8(lerp(float64(a.R), float64(b.R), t)),
		G: clampUint8(lerp(float64(a.G), float64(b.G), t)),
		B: clampUint8(lerp(float64(a.B), float64(b.B), t)),
	}
}

// BlendBg blends a foreground color into a background color by the given amount.
// BlendBg(fg, bg, amount) is equivalent to Blend(bg, fg, amount).
func BlendBg(fg, bg Color, amount float64) Color {
	return Blend(bg, fg, amount)
}

// Lighten increases the HSL lightness of a color by the given amount.
// The lightness is clamped to a maximum of 1.0.
func Lighten(c Color, amount float64) Color {
	hsl := c.ToHSL()
	hsl.L = math.Min(hsl.L+amount, 1.0)
	return hsl.ToRGB()
}

// Darken decreases the HSL lightness of a color by the given amount.
// The lightness is clamped to a minimum of 0.0.
func Darken(c Color, amount float64) Color {
	hsl := c.ToHSL()
	hsl.L = math.Max(hsl.L-amount, 0.0)
	return hsl.ToRGB()
}

// Desaturate reduces the HSL saturation of a color by a factor.
// The new saturation is S * (1 - amount). An amount of 1.0 produces grayscale.
func Desaturate(c Color, amount float64) Color {
	hsl := c.ToHSL()
	hsl.S *= (1.0 - amount)
	return hsl.ToRGB()
}

// ShiftHue rotates the hue of a color by the given number of degrees.
// The result wraps around modulo 360.
func ShiftHue(c Color, degrees float64) Color {
	hsl := c.ToHSL()
	hsl.H = math.Mod(hsl.H+degrees, 360.0)
	if hsl.H < 0 {
		hsl.H += 360.0
	}
	return hsl.ToRGB()
}

// lerp performs linear interpolation between two float64 values.
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// clampUint8 rounds and clamps a float64 to a uint8 value [0, 255].
func clampUint8(v float64) uint8 {
	rounded := math.Round(v)
	if rounded < 0 {
		return 0
	}
	if rounded > 255 {
		return 255
	}
	return uint8(rounded)
}

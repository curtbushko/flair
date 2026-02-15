package domain

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color represents an RGB color value.
// When IsNone is true, the color is a sentinel indicating "no color".
type Color struct {
	R      uint8
	G      uint8
	B      uint8
	IsNone bool
}

// HSL represents a color in the HSL (Hue, Saturation, Lightness) color space.
// H is in degrees [0, 360), S and L are in the range [0, 1].
type HSL struct {
	H float64
	S float64
	L float64
}

// ParseHex parses a hex color string and returns a Color.
// It accepts 3-digit and 6-digit hex strings, with or without a leading '#'.
// Returns a *ParseError for invalid input.
func ParseHex(hex string) (Color, error) {
	s := strings.TrimPrefix(hex, "#")

	switch len(s) {
	case 3:
		// Expand shorthand: "f00" -> "ff0000"
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	case 6:
		// Valid length, proceed.
	default:
		return Color{}, &ParseError{
			Field:   "hex",
			Message: fmt.Sprintf("invalid length %d, expected 3 or 6 hex digits", len(s)),
		}
	}

	val, err := strconv.ParseUint(s, 16, 24)
	if err != nil {
		return Color{}, &ParseError{
			Field:   "hex",
			Message: fmt.Sprintf("invalid hex characters in %q", s),
			Cause:   err,
		}
	}

	return Color{
		R: uint8(val >> 16),
		G: uint8(val >> 8),
		B: uint8(val),
	}, nil
}

// Hex returns the color formatted as a lowercase 6-digit hex string with '#' prefix.
func (c Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// NoneColor returns a sentinel Color with IsNone set to true.
func NoneColor() Color {
	return Color{IsNone: true}
}

// Equal compares two colors by their RGB values, ignoring the IsNone flag.
func (c Color) Equal(other Color) bool {
	return c.R == other.R && c.G == other.G && c.B == other.B
}

// ToHSL converts an RGB Color to the HSL color space.
func (c Color) ToHSL() HSL {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	maxVal := math.Max(r, math.Max(g, b))
	minVal := math.Min(r, math.Min(g, b))
	delta := maxVal - minVal

	l := (maxVal + minVal) / 2.0

	if delta == 0 {
		return HSL{H: 0, S: 0, L: l}
	}

	var s float64
	if l < 0.5 {
		s = delta / (maxVal + minVal)
	} else {
		s = delta / (2.0 - maxVal - minVal)
	}

	var h float64
	switch maxVal {
	case r:
		h = (g - b) / delta
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/delta + 2
	case b:
		h = (r-g)/delta + 4
	}
	h *= 60

	return HSL{H: h, S: s, L: l}
}

// ToRGB converts an HSL color back to an RGB Color.
func (hsl HSL) ToRGB() Color {
	if hsl.S == 0 {
		v := uint8(math.Round(hsl.L * 255))
		return Color{R: v, G: v, B: v}
	}

	var q float64
	if hsl.L < 0.5 {
		q = hsl.L * (1 + hsl.S)
	} else {
		q = hsl.L + hsl.S - hsl.L*hsl.S
	}
	p := 2*hsl.L - q

	h := hsl.H / 360.0

	r := hueToRGB(p, q, h+1.0/3.0)
	g := hueToRGB(p, q, h)
	b := hueToRGB(p, q, h-1.0/3.0)

	return Color{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
	}
}

// hueToRGB is a helper for HSL-to-RGB conversion.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6.0:
		return p + (q-p)*6*t
	case t < 1.0/2.0:
		return q
	case t < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-t)*6
	default:
		return p
	}
}

// Luminance calculates the WCAG 2.1 relative luminance of the color.
// Returns a value between 0.0 (black) and 1.0 (white).
// Formula: L = 0.2126 * R + 0.7152 * G + 0.0722 * B
// where each channel is linearized from sRGB.
func (c Color) Luminance() float64 {
	r := linearize(float64(c.R) / 255.0)
	g := linearize(float64(c.G) / 255.0)
	b := linearize(float64(c.B) / 255.0)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// linearize converts an sRGB channel value to linear RGB.
func linearize(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

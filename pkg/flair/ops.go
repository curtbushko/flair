package flair

import "math"

// BlendBg blends a foreground color onto a background color at the given alpha.
//
// This function performs standard alpha compositing using linear interpolation:
//
//	result = fg * alpha + bg * (1 - alpha)
//
// Parameters:
//   - fg: The foreground color to blend.
//   - bg: The background color to blend onto.
//   - alpha: The blend factor (0.0 = pure background, 1.0 = pure foreground).
//
// Alpha values are clamped to the range [0.0, 1.0].
//
// BlendBg is used internally by [Tokenize] to create subtle color variations
// such as selection highlights and search backgrounds.
//
// Example:
//
//	blue := flair.Color{R: 122, G: 162, B: 247}
//	black := flair.Color{R: 0, G: 0, B: 0}
//	subtle := flair.BlendBg(blue, black, 0.3) // 30% blue on black
func BlendBg(fg, bg Color, alpha float64) Color {
	// Clamp alpha to [0.0, 1.0].
	if alpha < 0.0 {
		alpha = 0.0
	}
	if alpha > 1.0 {
		alpha = 1.0
	}

	// Linear interpolation: result = fg * alpha + bg * (1 - alpha)
	r := float64(fg.R)*alpha + float64(bg.R)*(1.0-alpha)
	g := float64(fg.G)*alpha + float64(bg.G)*(1.0-alpha)
	b := float64(fg.B)*alpha + float64(bg.B)*(1.0-alpha)

	return Color{
		R: uint8(math.Round(r)),
		G: uint8(math.Round(g)),
		B: uint8(math.Round(b)),
	}
}

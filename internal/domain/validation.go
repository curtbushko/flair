package domain

import "fmt"

// ValidatePalette checks a palette against validation rules and returns
// a slice of violation/warning messages. An empty slice means the palette
// is valid.
//
// Rules:
//  1. Completeness: all 24 slots must have non-IsNone colors.
//  2. Luminance ordering (dark): base00.Luminance < base05.Luminance.
//  3. Luminance ordering (light): base00.Luminance > base05.Luminance.
//  4. Neutral ramp monotonicity: luminance should increase base00-base07 (warning).
//  5. Bright variants: base12-base17 luminance >= base08-base0E (warning).
func ValidatePalette(p *Palette) []string {
	var violations []string

	if v := checkCompleteness(p); len(v) > 0 {
		violations = append(violations, v...)
	}
	if v := checkLuminanceOrdering(p); len(v) > 0 {
		violations = append(violations, v...)
	}
	if v := checkMonotonicity(p); len(v) > 0 {
		violations = append(violations, v...)
	}
	if v := checkBrightVariants(p); len(v) > 0 {
		violations = append(violations, v...)
	}

	return violations
}

// checkCompleteness verifies that all 24 color slots are populated (non-IsNone).
func checkCompleteness(p *Palette) []string {
	var violations []string
	for i := 0; i < 24; i++ {
		if p.Colors[i].IsNone {
			violations = append(violations, fmt.Sprintf("completeness: %s is missing (IsNone)", slotNames[i]))
		}
	}
	return violations
}

// checkLuminanceOrdering verifies that backgrounds are darker than foregrounds
// for dark themes, and lighter than foregrounds for light themes.
func checkLuminanceOrdering(p *Palette) []string {
	var violations []string

	bg := p.Colors[0] // base00
	fg := p.Colors[5] // base05
	bgLum := bg.Luminance()
	fgLum := fg.Luminance()

	switch p.Variant {
	case "dark":
		// Dark theme: background should be darker (lower luminance) than foreground
		if bgLum >= fgLum {
			violations = append(violations,
				fmt.Sprintf("luminance ordering (dark): base00 luminance (%.4f) must be less than base05 luminance (%.4f)",
					bgLum, fgLum))
		}
	case "light":
		// Light theme: background should be lighter (higher luminance) than foreground
		if bgLum <= fgLum {
			violations = append(violations,
				fmt.Sprintf("luminance ordering (light): base00 luminance (%.4f) must be greater than base05 luminance (%.4f)",
					bgLum, fgLum))
		}
	}

	return violations
}

// checkMonotonicity warns if the neutral ramp (base00-base07) does not have
// monotonically increasing luminance.
func checkMonotonicity(p *Palette) []string {
	var violations []string

	for i := 1; i <= 7; i++ {
		prevLum := p.Colors[i-1].Luminance()
		currLum := p.Colors[i].Luminance()
		if currLum < prevLum {
			violations = append(violations,
				fmt.Sprintf("monotonicity warning: %s luminance (%.4f) is less than %s luminance (%.4f)",
					slotNames[i], currLum, slotNames[i-1], prevLum))
		}
	}

	return violations
}

// checkBrightVariants warns if any bright variant (base12-base17) has lower
// luminance than its corresponding base accent (base08-base0E).
//
// Mapping: base12<->base08, base13<->base0A, base14<->base0B,
//
//	base15<->base0C, base16<->base0D, base17<->base0E
func checkBrightVariants(p *Palette) []string {
	var violations []string

	// Bright variant index -> base accent index
	pairs := [][2]int{
		{18, 8},  // base12 <-> base08
		{19, 10}, // base13 <-> base0A
		{20, 11}, // base14 <-> base0B
		{21, 12}, // base15 <-> base0C
		{22, 13}, // base16 <-> base0D
		{23, 14}, // base17 <-> base0E
	}

	for _, pair := range pairs {
		brightIdx := pair[0]
		baseIdx := pair[1]
		brightLum := p.Colors[brightIdx].Luminance()
		baseLum := p.Colors[baseIdx].Luminance()

		if brightLum < baseLum {
			violations = append(violations,
				fmt.Sprintf("bright variant warning: %s luminance (%.4f) is less than %s luminance (%.4f)",
					slotNames[brightIdx], brightLum, slotNames[baseIdx], baseLum))
		}
	}

	return violations
}

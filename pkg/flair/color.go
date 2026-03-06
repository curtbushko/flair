package flair

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents an RGB color value with 8 bits per channel.
//
// Color is a simple, immutable value type designed for external consumers.
// It provides methods for hex string conversion and equality comparison.
// The zero value represents black (#000000).
//
// Example:
//
//	c := flair.Color{R: 122, G: 162, B: 247}
//	fmt.Println(c.Hex()) // Output: #7aa2f7
type Color struct {
	// R is the red channel (0-255).
	R uint8
	// G is the green channel (0-255).
	G uint8
	// B is the blue channel (0-255).
	B uint8
}

// Hex returns the color formatted as a lowercase 6-digit hex string with '#' prefix.
//
// Example:
//
//	c := flair.Color{R: 255, G: 128, B: 0}
//	fmt.Println(c.Hex()) // Output: #ff8000
func (c Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// Equal reports whether c and other have identical RGB values.
func (c Color) Equal(other Color) bool {
	return c.R == other.R && c.G == other.G && c.B == other.B
}

// ParseHex parses a hex color string and returns a Color.
//
// ParseHex accepts 3-digit and 6-digit hex strings, with or without a leading '#'.
// 3-digit hex strings are expanded (e.g., "f00" becomes "ff0000").
//
// Examples of valid input:
//
//	"#ff8000"  // 6-digit with hash
//	"ff8000"   // 6-digit without hash
//	"#f80"     // 3-digit with hash (expands to ff8800)
//	"f80"      // 3-digit without hash
//
// ParseHex returns an error for invalid input such as wrong length or
// non-hexadecimal characters.
func ParseHex(hex string) (Color, error) {
	s := strings.TrimPrefix(hex, "#")

	switch len(s) {
	case 3:
		// Expand shorthand: "f00" -> "ff0000"
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	case 6:
		// Valid length, proceed.
	default:
		return Color{}, fmt.Errorf("invalid hex color %q: expected 3 or 6 hex digits", hex)
	}

	val, err := strconv.ParseUint(s, 16, 24)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}

	return Color{
		R: uint8(val >> 16),
		G: uint8(val >> 8),
		B: uint8(val),
	}, nil
}

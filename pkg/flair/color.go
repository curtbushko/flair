package flair

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents an RGB color value with 8 bits per channel.
// This is a simple, read-only value type for external consumers.
type Color struct {
	R uint8
	G uint8
	B uint8
}

// Hex returns the color formatted as a lowercase 6-digit hex string with '#' prefix.
func (c Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// Equal compares two colors by their RGB values.
func (c Color) Equal(other Color) bool {
	return c.R == other.R && c.G == other.G && c.B == other.B
}

// ParseHex parses a hex color string and returns a Color.
// It accepts 3-digit and 6-digit hex strings, with or without a leading '#'.
// Returns an error for invalid input.
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

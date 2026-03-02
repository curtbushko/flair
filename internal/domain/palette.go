package domain

import "fmt"

// slotNames lists the 24 base24 slot names in index order.
var slotNames = [24]string{
	"base00", "base01", "base02", "base03",
	"base04", "base05", "base06", "base07",
	"base08", "base09", "base0A", "base0B",
	"base0C", "base0D", "base0E", "base0F",
	"base10", "base11", "base12", "base13",
	"base14", "base15", "base16", "base17",
}

// slotIndex maps a slot name to its index in the Colors array.
var slotIndex = func() map[string]int {
	idx := make(map[string]int, len(slotNames))
	for i, name := range slotNames {
		idx[name] = i
	}
	return idx
}()

// base16Fallbacks maps base24-extended slots (indices 16-23) to their
// fallback base16 slot index when only 16 colors are provided.
//
//	base10 -> base00, base11 -> base00, base12 -> base08,
//	base13 -> base0A, base14 -> base0B, base15 -> base0C,
//	base16 -> base0D, base17 -> base0E
var base16Fallbacks = map[int]int{
	16: 0,  // base10 -> base00
	17: 0,  // base11 -> base00
	18: 8,  // base12 -> base08
	19: 10, // base13 -> base0A
	20: 11, // base14 -> base0B
	21: 12, // base15 -> base0C
	22: 13, // base16 -> base0D
	23: 14, // base17 -> base0E
}

// Palette represents a complete base24 color palette.
type Palette struct {
	Name      string
	Author    string
	Variant   string
	System    string
	Slug      string
	Colors    [24]Color
	Overrides map[string]TokenOverride // Optional token overrides by path (nil when none)
}

// NewPalette constructs a Palette from a map of slot names to hex color strings.
// When system is "base16" and only 16 colors are provided, the remaining 8 slots
// are filled using the standard base16-to-base24 fallback rules.
// Returns a *ParseError if required slots are missing or a hex value is invalid.
func NewPalette(name, author, variant, system string, colors map[string]string) (*Palette, error) {
	pal := &Palette{
		Name:    name,
		Author:  author,
		Variant: variant,
		System:  system,
		Slug:    name + "-" + variant,
	}

	if err := parseRequiredSlots(pal, system, colors); err != nil {
		return nil, err
	}

	if err := fillExtendedSlots(pal, system); err != nil {
		return nil, err
	}

	return pal, nil
}

// parseRequiredSlots parses the first 16 (base16) or 24 (base24) required slots.
func parseRequiredSlots(pal *Palette, system string, colors map[string]string) error {
	requiredCount := 24
	if system == "base16" {
		requiredCount = 16
	}

	for i := 0; i < requiredCount; i++ {
		slotName := slotNames[i]
		hex, ok := colors[slotName]
		if !ok {
			return &ParseError{
				Field:   slotName,
				Message: "missing required slot " + slotName,
			}
		}
		c, err := ParseHex(hex)
		if err != nil {
			return &ParseError{
				Field:   slotName,
				Message: fmt.Sprintf("invalid hex value %q for slot %s", hex, slotName),
				Cause:   err,
			}
		}
		pal.Colors[i] = c
	}

	return nil
}

// fillExtendedSlots populates slots 16-23 from fallback rules when the
// system is base16. For base24, slots 16-23 are already parsed in parseRequiredSlots.
func fillExtendedSlots(pal *Palette, system string) error {
	if system == "base16" {
		return applyBase16Fallbacks(pal)
	}
	// base24 slots 16-23 are already parsed in parseRequiredSlots.
	return nil
}

// applyBase16Fallbacks copies colors from base16 slots to the extended base24 slots.
func applyBase16Fallbacks(pal *Palette) error {
	for i := 16; i < 24; i++ {
		fallbackIdx, ok := base16Fallbacks[i]
		if !ok {
			return &ParseError{
				Field:   slotNames[i],
				Message: fmt.Sprintf("no fallback defined for slot index %d", i),
			}
		}
		pal.Colors[i] = pal.Colors[fallbackIdx]
	}
	return nil
}

// Base returns the color at the given index (0-23).
// Returns NoneColor for out-of-range indices.
func (p *Palette) Base(n int) Color {
	if n < 0 || n >= len(p.Colors) {
		return NoneColor()
	}
	return p.Colors[n]
}

// Slot returns the color for the given slot name (e.g., "base0D").
// Returns an error if the slot name is not recognized.
func (p *Palette) Slot(name string) (Color, error) {
	idx, ok := slotIndex[name]
	if !ok {
		return Color{}, &ParseError{
			Field:   name,
			Message: "unknown slot name " + name,
		}
	}
	return p.Colors[idx], nil
}

// SlotNames returns the ordered list of all 24 slot name strings.
func (p *Palette) SlotNames() []string {
	names := make([]string, len(slotNames))
	copy(names, slotNames[:])
	return names
}

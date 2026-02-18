package viewer

// This file contains helper types and functions for page rendering.
// The actual rendering logic is in view.go.

// String returns a human-readable name for the page.
func (p Page) String() string {
	switch p {
	case PageSelector:
		return "Theme Selector"
	case PagePalette:
		return "Palette"
	case PageTokens:
		return "Tokens"
	case PageComponents:
		return "Components"
	default:
		return "Unknown"
	}
}

// TokenCategories returns the ordered list of token category names.
func TokenCategories() []string {
	return []string{"Surface", "Text", "Status", "Syntax", "Diff"}
}

// SlotNames returns the 24 base24 slot name strings.
func SlotNames() [24]string {
	return slotNames
}

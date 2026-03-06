package viewer

// This file contains helper types and functions for page rendering.
// The actual rendering logic is in view.go.

// String returns a human-readable name for the page.
//
// This is useful for displaying page titles in the UI.
func (p Page) String() string {
	switch p {
	case PageTextStatus:
		return "Text & Status"
	case PageInteractive:
		return "Interactive Components"
	case PageDataDisplay:
		return "Data Display"
	default:
		return "Unknown"
	}
}

// TokenCategories returns the ordered list of token category names.
//
// These categories correspond to the groupings in [TokenData].
func TokenCategories() []string {
	return []string{"Surface", "Text", "Status", "Syntax", "Diff"}
}

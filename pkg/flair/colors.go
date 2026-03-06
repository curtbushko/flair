package flair

// SurfaceColors provides typed access to surface/background color tokens.
//
// These colors are used for UI backgrounds and container elements.
// Each field corresponds to a semantic token in the "surface.*" namespace.
type SurfaceColors struct {
	// Background is the primary background color (surface.background).
	Background Color
	// Raised is for elevated surfaces like cards (surface.background.raised).
	Raised Color
	// Sunken is for inset areas like input fields (surface.background.sunken).
	Sunken Color
	// Darkest is the darkest background shade (surface.background.darkest).
	Darkest Color
	// Overlay is for modal overlays (surface.background.overlay).
	Overlay Color
	// Popup is for popup menus and tooltips (surface.background.popup).
	Popup Color
	// Highlight is for hover states (surface.background.highlight).
	Highlight Color
	// Selection is for selected text backgrounds (surface.background.selection).
	Selection Color
	// Search is for search match highlights (surface.background.search).
	Search Color
	// Sidebar is for sidebar backgrounds (surface.background.sidebar).
	Sidebar Color
	// Statusbar is for status bar backgrounds (surface.background.statusbar).
	Statusbar Color
}

// TextColors provides typed access to text/foreground color tokens.
//
// These colors are used for text and foreground elements.
// Each field corresponds to a semantic token in the "text.*" namespace.
type TextColors struct {
	// Primary is the main text color for body content (text.primary).
	Primary Color
	// Secondary is for less prominent text (text.secondary).
	Secondary Color
	// Muted is for disabled or placeholder text (text.muted).
	Muted Color
	// Subtle is for very low-contrast text (text.subtle).
	Subtle Color
	// Inverse is text on accent backgrounds (text.inverse).
	Inverse Color
	// Overlay is text in overlay contexts (text.overlay).
	Overlay Color
	// Sidebar is text in sidebars (text.sidebar).
	Sidebar Color
}

// StatusColors provides typed access to semantic status color tokens.
//
// These colors indicate states and severity levels in the UI.
// Each field corresponds to a semantic token in the "status.*" namespace.
type StatusColors struct {
	// Error is for error messages and indicators (status.error).
	Error Color
	// Warning is for warning messages (status.warning).
	Warning Color
	// Success is for success messages (status.success).
	Success Color
	// Info is for informational messages (status.info).
	Info Color
	// Hint is for hints and suggestions (status.hint).
	Hint Color
	// Todo is for TODO comments and markers (status.todo).
	Todo Color
}

// SyntaxColors provides typed access to syntax highlighting color tokens.
//
// These colors are used for code syntax highlighting in editors.
// Each field corresponds to a semantic token in the "syntax.*" namespace.
type SyntaxColors struct {
	// Keyword is for language keywords like if, for, return (syntax.keyword).
	Keyword Color
	// String is for string literals (syntax.string).
	String Color
	// Function is for function names and calls (syntax.function).
	Function Color
	// Comment is for code comments (syntax.comment).
	Comment Color
	// Variable is for variable names (syntax.variable).
	Variable Color
	// Constant is for constant values (syntax.constant).
	Constant Color
	// Operator is for operators like +, -, = (syntax.operator).
	Operator Color
	// Type is for type names and annotations (syntax.type).
	Type Color
	// Number is for numeric literals (syntax.number).
	Number Color
	// Tag is for HTML/XML tags (syntax.tag).
	Tag Color
	// Property is for object properties (syntax.property).
	Property Color
	// Parameter is for function parameters (syntax.parameter).
	Parameter Color
	// Regexp is for regular expressions (syntax.regexp).
	Regexp Color
	// Escape is for escape sequences (syntax.escape).
	Escape Color
	// Constructor is for constructor functions (syntax.constructor).
	Constructor Color
}

// DiffColors provides typed access to diff/version control color tokens.
//
// These colors are used for displaying code diffs and version control status.
// Each field corresponds to a semantic token in the "diff.*" namespace.
type DiffColors struct {
	// AddedFg is the foreground color for added lines (diff.added.fg).
	AddedFg Color
	// AddedBg is the background color for added lines (diff.added.bg).
	AddedBg Color
	// AddedSign is the sign/gutter color for added lines (diff.added.sign).
	AddedSign Color
	// DeletedFg is the foreground color for deleted lines (diff.deleted.fg).
	DeletedFg Color
	// DeletedBg is the background color for deleted lines (diff.deleted.bg).
	DeletedBg Color
	// DeletedSign is the sign/gutter color for deleted lines (diff.deleted.sign).
	DeletedSign Color
	// ChangedFg is the foreground color for changed lines (diff.changed.fg).
	ChangedFg Color
	// ChangedBg is the background color for changed lines (diff.changed.bg).
	ChangedBg Color
	// Ignored is the color for ignored files (diff.ignored).
	Ignored Color
}

// terminalTokenPaths maps ANSI color indices (0-15) to token paths.
var terminalTokenPaths = [16]string{
	"terminal.black",     // 0
	"terminal.red",       // 1
	"terminal.green",     // 2
	"terminal.yellow",    // 3
	"terminal.blue",      // 4
	"terminal.magenta",   // 5
	"terminal.cyan",      // 6
	"terminal.white",     // 7
	"terminal.brblack",   // 8
	"terminal.brred",     // 9
	"terminal.brgreen",   // 10
	"terminal.bryellow",  // 11
	"terminal.brblue",    // 12
	"terminal.brmagenta", // 13
	"terminal.brcyan",    // 14
	"terminal.brwhite",   // 15
}

// Surface returns typed access to surface/background colors from the theme.
//
// This method provides compile-time safe access to surface tokens. Missing
// tokens return zero Color values (black).
//
// Example:
//
//	surface := theme.Surface()
//	fmt.Printf("Background: %s\n", surface.Background.Hex())
//	fmt.Printf("Raised: %s\n", surface.Raised.Hex())
func (t *Theme) Surface() SurfaceColors {
	return SurfaceColors{
		Background: t.getColor("surface.background"),
		Raised:     t.getColor("surface.background.raised"),
		Sunken:     t.getColor("surface.background.sunken"),
		Darkest:    t.getColor("surface.background.darkest"),
		Overlay:    t.getColor("surface.background.overlay"),
		Popup:      t.getColor("surface.background.popup"),
		Highlight:  t.getColor("surface.background.highlight"),
		Selection:  t.getColor("surface.background.selection"),
		Search:     t.getColor("surface.background.search"),
		Sidebar:    t.getColor("surface.background.sidebar"),
		Statusbar:  t.getColor("surface.background.statusbar"),
	}
}

// Text returns typed access to text/foreground colors from the theme.
//
// This method provides compile-time safe access to text tokens. Missing
// tokens return zero Color values (black).
//
// Example:
//
//	text := theme.Text()
//	fmt.Printf("Primary: %s\n", text.Primary.Hex())
//	fmt.Printf("Muted: %s\n", text.Muted.Hex())
func (t *Theme) Text() TextColors {
	return TextColors{
		Primary:   t.getColor("text.primary"),
		Secondary: t.getColor("text.secondary"),
		Muted:     t.getColor("text.muted"),
		Subtle:    t.getColor("text.subtle"),
		Inverse:   t.getColor("text.inverse"),
		Overlay:   t.getColor("text.overlay"),
		Sidebar:   t.getColor("text.sidebar"),
	}
}

// Status returns typed access to semantic status colors from the theme.
//
// This method provides compile-time safe access to status tokens. Missing
// tokens return zero Color values (black).
//
// Example:
//
//	status := theme.Status()
//	fmt.Printf("Error: %s\n", status.Error.Hex())
//	fmt.Printf("Success: %s\n", status.Success.Hex())
func (t *Theme) Status() StatusColors {
	return StatusColors{
		Error:   t.getColor("status.error"),
		Warning: t.getColor("status.warning"),
		Success: t.getColor("status.success"),
		Info:    t.getColor("status.info"),
		Hint:    t.getColor("status.hint"),
		Todo:    t.getColor("status.todo"),
	}
}

// Syntax returns typed access to syntax highlighting colors from the theme.
//
// This method provides compile-time safe access to syntax tokens. Missing
// tokens return zero Color values (black).
//
// Example:
//
//	syntax := theme.Syntax()
//	fmt.Printf("Keyword: %s\n", syntax.Keyword.Hex())
//	fmt.Printf("String: %s\n", syntax.String.Hex())
func (t *Theme) Syntax() SyntaxColors {
	return SyntaxColors{
		Keyword:     t.getColor("syntax.keyword"),
		String:      t.getColor("syntax.string"),
		Function:    t.getColor("syntax.function"),
		Comment:     t.getColor("syntax.comment"),
		Variable:    t.getColor("syntax.variable"),
		Constant:    t.getColor("syntax.constant"),
		Operator:    t.getColor("syntax.operator"),
		Type:        t.getColor("syntax.type"),
		Number:      t.getColor("syntax.number"),
		Tag:         t.getColor("syntax.tag"),
		Property:    t.getColor("syntax.property"),
		Parameter:   t.getColor("syntax.parameter"),
		Regexp:      t.getColor("syntax.regexp"),
		Escape:      t.getColor("syntax.escape"),
		Constructor: t.getColor("syntax.constructor"),
	}
}

// Diff returns typed access to diff/version control colors from the theme.
//
// This method provides compile-time safe access to diff tokens. Missing
// tokens return zero Color values (black).
//
// Example:
//
//	diff := theme.Diff()
//	fmt.Printf("Added: %s\n", diff.AddedFg.Hex())
//	fmt.Printf("Deleted: %s\n", diff.DeletedFg.Hex())
func (t *Theme) Diff() DiffColors {
	return DiffColors{
		AddedFg:     t.getColor("diff.added.fg"),
		AddedBg:     t.getColor("diff.added.bg"),
		AddedSign:   t.getColor("diff.added.sign"),
		DeletedFg:   t.getColor("diff.deleted.fg"),
		DeletedBg:   t.getColor("diff.deleted.bg"),
		DeletedSign: t.getColor("diff.deleted.sign"),
		ChangedFg:   t.getColor("diff.changed.fg"),
		ChangedBg:   t.getColor("diff.changed.bg"),
		Ignored:     t.getColor("diff.ignored"),
	}
}

// Terminal returns the 16 ANSI terminal colors (indices 0-15).
//
// The colors are ordered according to the standard ANSI color scheme:
//
//	0: black,   1: red,       2: green,   3: yellow
//	4: blue,    5: magenta,   6: cyan,    7: white
//	8: brblack, 9: brred,    10: brgreen, 11: bryellow
//	12: brblue, 13: brmagenta, 14: brcyan, 15: brwhite
//
// Missing tokens return zero Color values (black) at their respective indices.
//
// Example:
//
//	colors := theme.Terminal()
//	fmt.Printf("Red: %s\n", colors[1].Hex())
//	fmt.Printf("Bright Blue: %s\n", colors[12].Hex())
func (t *Theme) Terminal() [16]Color {
	var colors [16]Color
	for i, path := range terminalTokenPaths {
		colors[i] = t.getColor(path)
	}
	return colors
}

// Get retrieves a color by its semantic token path.
//
// Get is an alias for [Theme.Color] provided for convenience.
// Token paths follow a hierarchical naming convention such as
// "accent.primary", "border.default", or "scrollbar.thumb".
//
// Returns the color and true if found, or a zero Color and false if not.
func (t *Theme) Get(path string) (Color, bool) {
	return t.Color(path)
}

// getColor is a helper that returns the color for a token path,
// or a zero Color if the token is not found.
func (t *Theme) getColor(path string) Color {
	return t.colors[path]
}

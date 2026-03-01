package flair

// SurfaceColors provides typed access to surface/background color tokens.
type SurfaceColors struct {
	Background Color
	Raised     Color
	Sunken     Color
	Darkest    Color
	Overlay    Color
	Popup      Color
	Highlight  Color
	Selection  Color
	Search     Color
	Sidebar    Color
	Statusbar  Color
}

// TextColors provides typed access to text/foreground color tokens.
type TextColors struct {
	Primary   Color
	Secondary Color
	Muted     Color
	Subtle    Color
	Inverse   Color
	Overlay   Color
	Sidebar   Color
}

// StatusColors provides typed access to semantic status color tokens.
type StatusColors struct {
	Error   Color
	Warning Color
	Success Color
	Info    Color
	Hint    Color
	Todo    Color
}

// SyntaxColors provides typed access to syntax highlighting color tokens.
type SyntaxColors struct {
	Keyword     Color
	String      Color
	Function    Color
	Comment     Color
	Variable    Color
	Constant    Color
	Operator    Color
	Type        Color
	Number      Color
	Tag         Color
	Property    Color
	Parameter   Color
	Regexp      Color
	Escape      Color
	Constructor Color
}

// DiffColors provides typed access to diff/version control color tokens.
type DiffColors struct {
	AddedFg     Color
	AddedBg     Color
	AddedSign   Color
	DeletedFg   Color
	DeletedBg   Color
	DeletedSign Color
	ChangedFg   Color
	ChangedBg   Color
	Ignored     Color
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

// Surface returns typed access to surface/background colors.
// Missing tokens return zero Color values.
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

// Text returns typed access to text/foreground colors.
// Missing tokens return zero Color values.
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

// Status returns typed access to semantic status colors.
// Missing tokens return zero Color values.
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

// Syntax returns typed access to syntax highlighting colors.
// Missing tokens return zero Color values.
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

// Diff returns typed access to diff/version control colors.
// Missing tokens return zero Color values.
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
// Missing tokens return zero Color values at their respective indices.
func (t *Theme) Terminal() [16]Color {
	var colors [16]Color
	for i, path := range terminalTokenPaths {
		colors[i] = t.getColor(path)
	}
	return colors
}

// Get retrieves a color by its token path (e.g., "accent.primary").
// Returns the color and true if found, or a zero Color and false if not.
func (t *Theme) Get(path string) (Color, bool) {
	return t.Color(path)
}

// getColor is a helper that returns the color for a token path,
// or a zero Color if the token is not found.
func (t *Theme) getColor(path string) Color {
	return t.colors[path]
}

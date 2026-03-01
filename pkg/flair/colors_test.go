package flair_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

// Helper to create a theme with all token types for testing.
func newTestTheme() *flair.Theme {
	colors := map[string]flair.Color{
		// Surface colors
		"surface.background":           {R: 26, G: 27, B: 38},
		"surface.background.raised":    {R: 31, G: 35, B: 53},
		"surface.background.sunken":    {R: 22, G: 22, B: 30},
		"surface.background.darkest":   {R: 16, G: 16, B: 20},
		"surface.background.overlay":   {R: 30, G: 31, B: 42},
		"surface.background.popup":     {R: 30, G: 31, B: 42},
		"surface.background.highlight": {R: 41, G: 46, B: 66},
		"surface.background.selection": {R: 50, G: 55, B: 75},
		"surface.background.search":    {R: 60, G: 65, B: 85},
		"surface.background.sidebar":   {R: 28, G: 29, B: 40},
		"surface.background.statusbar": {R: 25, G: 26, B: 37},

		// Text colors
		"text.primary":   {R: 192, G: 202, B: 245},
		"text.secondary": {R: 169, G: 177, B: 214},
		"text.muted":     {R: 86, G: 95, B: 137},
		"text.subtle":    {R: 70, G: 80, B: 120},
		"text.inverse":   {R: 26, G: 27, B: 38},
		"text.overlay":   {R: 200, G: 210, B: 250},
		"text.sidebar":   {R: 169, G: 177, B: 214},

		// Status colors
		"status.error":   {R: 247, G: 118, B: 142},
		"status.warning": {R: 224, G: 175, B: 104},
		"status.success": {R: 158, G: 206, B: 106},
		"status.info":    {R: 125, G: 207, B: 255},
		"status.hint":    {R: 125, G: 207, B: 255},
		"status.todo":    {R: 122, G: 162, B: 247},

		// Syntax colors
		"syntax.keyword":     {R: 187, G: 154, B: 247},
		"syntax.string":      {R: 158, G: 206, B: 106},
		"syntax.function":    {R: 122, G: 162, B: 247},
		"syntax.comment":     {R: 86, G: 95, B: 137},
		"syntax.variable":    {R: 192, G: 202, B: 245},
		"syntax.constant":    {R: 255, G: 158, B: 100},
		"syntax.operator":    {R: 137, G: 220, B: 235},
		"syntax.type":        {R: 224, G: 175, B: 104},
		"syntax.number":      {R: 255, G: 158, B: 100},
		"syntax.tag":         {R: 247, G: 118, B: 142},
		"syntax.property":    {R: 158, G: 206, B: 106},
		"syntax.parameter":   {R: 224, G: 175, B: 104},
		"syntax.regexp":      {R: 125, G: 207, B: 255},
		"syntax.escape":      {R: 187, G: 154, B: 247},
		"syntax.constructor": {R: 200, G: 172, B: 248},

		// Diff colors
		"diff.added.fg":     {R: 158, G: 206, B: 106},
		"diff.added.bg":     {R: 40, G: 60, B: 40},
		"diff.added.sign":   {R: 158, G: 206, B: 106},
		"diff.deleted.fg":   {R: 247, G: 118, B: 142},
		"diff.deleted.bg":   {R: 60, G: 40, B: 40},
		"diff.deleted.sign": {R: 247, G: 118, B: 142},
		"diff.changed.fg":   {R: 137, G: 220, B: 235},
		"diff.changed.bg":   {R: 40, G: 50, B: 60},
		"diff.ignored":      {R: 86, G: 95, B: 137},

		// Terminal colors (ANSI 0-15)
		"terminal.black":     {R: 31, G: 35, B: 53},
		"terminal.red":       {R: 247, G: 118, B: 142},
		"terminal.green":     {R: 158, G: 206, B: 106},
		"terminal.yellow":    {R: 224, G: 175, B: 104},
		"terminal.blue":      {R: 122, G: 162, B: 247},
		"terminal.magenta":   {R: 187, G: 154, B: 247},
		"terminal.cyan":      {R: 125, G: 207, B: 255},
		"terminal.white":     {R: 192, G: 202, B: 245},
		"terminal.brblack":   {R: 86, G: 95, B: 137},
		"terminal.brred":     {R: 255, G: 137, B: 157},
		"terminal.brgreen":   {R: 175, G: 214, B: 122},
		"terminal.bryellow":  {R: 233, G: 197, B: 130},
		"terminal.brblue":    {R: 141, G: 182, B: 250},
		"terminal.brmagenta": {R: 200, G: 172, B: 248},
		"terminal.brcyan":    {R: 151, G: 216, B: 248},
		"terminal.brwhite":   {R: 200, G: 211, B: 245},

		// Accent color for Get() test
		"accent.primary": {R: 122, G: 162, B: 247},
	}
	return flair.NewTheme("test-theme", "dark", colors)
}

func TestTheme_Surface(t *testing.T) {
	theme := newTestTheme()

	surface := theme.Surface()

	tests := []struct {
		name string
		got  flair.Color
		want flair.Color
	}{
		{"Background", surface.Background, flair.Color{R: 26, G: 27, B: 38}},
		{"Raised", surface.Raised, flair.Color{R: 31, G: 35, B: 53}},
		{"Sunken", surface.Sunken, flair.Color{R: 22, G: 22, B: 30}},
		{"Darkest", surface.Darkest, flair.Color{R: 16, G: 16, B: 20}},
		{"Overlay", surface.Overlay, flair.Color{R: 30, G: 31, B: 42}},
		{"Popup", surface.Popup, flair.Color{R: 30, G: 31, B: 42}},
		{"Highlight", surface.Highlight, flair.Color{R: 41, G: 46, B: 66}},
		{"Selection", surface.Selection, flair.Color{R: 50, G: 55, B: 75}},
		{"Search", surface.Search, flair.Color{R: 60, G: 65, B: 85}},
		{"Sidebar", surface.Sidebar, flair.Color{R: 28, G: 29, B: 40}},
		{"Statusbar", surface.Statusbar, flair.Color{R: 25, G: 26, B: 37}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.want) {
				t.Errorf("Surface.%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestTheme_Text(t *testing.T) {
	theme := newTestTheme()

	text := theme.Text()

	tests := []struct {
		name string
		got  flair.Color
		want flair.Color
	}{
		{"Primary", text.Primary, flair.Color{R: 192, G: 202, B: 245}},
		{"Secondary", text.Secondary, flair.Color{R: 169, G: 177, B: 214}},
		{"Muted", text.Muted, flair.Color{R: 86, G: 95, B: 137}},
		{"Subtle", text.Subtle, flair.Color{R: 70, G: 80, B: 120}},
		{"Inverse", text.Inverse, flair.Color{R: 26, G: 27, B: 38}},
		{"Overlay", text.Overlay, flair.Color{R: 200, G: 210, B: 250}},
		{"Sidebar", text.Sidebar, flair.Color{R: 169, G: 177, B: 214}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.want) {
				t.Errorf("Text.%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestTheme_Status(t *testing.T) {
	theme := newTestTheme()

	status := theme.Status()

	tests := []struct {
		name string
		got  flair.Color
		want flair.Color
	}{
		{"Error", status.Error, flair.Color{R: 247, G: 118, B: 142}},
		{"Warning", status.Warning, flair.Color{R: 224, G: 175, B: 104}},
		{"Success", status.Success, flair.Color{R: 158, G: 206, B: 106}},
		{"Info", status.Info, flair.Color{R: 125, G: 207, B: 255}},
		{"Hint", status.Hint, flair.Color{R: 125, G: 207, B: 255}},
		{"Todo", status.Todo, flair.Color{R: 122, G: 162, B: 247}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.want) {
				t.Errorf("Status.%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestTheme_Syntax(t *testing.T) {
	theme := newTestTheme()

	syntax := theme.Syntax()

	tests := []struct {
		name string
		got  flair.Color
		want flair.Color
	}{
		{"Keyword", syntax.Keyword, flair.Color{R: 187, G: 154, B: 247}},
		{"String", syntax.String, flair.Color{R: 158, G: 206, B: 106}},
		{"Function", syntax.Function, flair.Color{R: 122, G: 162, B: 247}},
		{"Comment", syntax.Comment, flair.Color{R: 86, G: 95, B: 137}},
		{"Variable", syntax.Variable, flair.Color{R: 192, G: 202, B: 245}},
		{"Constant", syntax.Constant, flair.Color{R: 255, G: 158, B: 100}},
		{"Operator", syntax.Operator, flair.Color{R: 137, G: 220, B: 235}},
		{"Type", syntax.Type, flair.Color{R: 224, G: 175, B: 104}},
		{"Number", syntax.Number, flair.Color{R: 255, G: 158, B: 100}},
		{"Tag", syntax.Tag, flair.Color{R: 247, G: 118, B: 142}},
		{"Property", syntax.Property, flair.Color{R: 158, G: 206, B: 106}},
		{"Parameter", syntax.Parameter, flair.Color{R: 224, G: 175, B: 104}},
		{"Regexp", syntax.Regexp, flair.Color{R: 125, G: 207, B: 255}},
		{"Escape", syntax.Escape, flair.Color{R: 187, G: 154, B: 247}},
		{"Constructor", syntax.Constructor, flair.Color{R: 200, G: 172, B: 248}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.want) {
				t.Errorf("Syntax.%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestTheme_Diff(t *testing.T) {
	theme := newTestTheme()

	diff := theme.Diff()

	tests := []struct {
		name string
		got  flair.Color
		want flair.Color
	}{
		{"AddedFg", diff.AddedFg, flair.Color{R: 158, G: 206, B: 106}},
		{"AddedBg", diff.AddedBg, flair.Color{R: 40, G: 60, B: 40}},
		{"AddedSign", diff.AddedSign, flair.Color{R: 158, G: 206, B: 106}},
		{"DeletedFg", diff.DeletedFg, flair.Color{R: 247, G: 118, B: 142}},
		{"DeletedBg", diff.DeletedBg, flair.Color{R: 60, G: 40, B: 40}},
		{"DeletedSign", diff.DeletedSign, flair.Color{R: 247, G: 118, B: 142}},
		{"ChangedFg", diff.ChangedFg, flair.Color{R: 137, G: 220, B: 235}},
		{"ChangedBg", diff.ChangedBg, flair.Color{R: 40, G: 50, B: 60}},
		{"Ignored", diff.Ignored, flair.Color{R: 86, G: 95, B: 137}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.want) {
				t.Errorf("Diff.%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestTheme_Terminal(t *testing.T) {
	theme := newTestTheme()

	terminal := theme.Terminal()

	expected := [16]flair.Color{
		{R: 31, G: 35, B: 53},    // 0: black
		{R: 247, G: 118, B: 142}, // 1: red
		{R: 158, G: 206, B: 106}, // 2: green
		{R: 224, G: 175, B: 104}, // 3: yellow
		{R: 122, G: 162, B: 247}, // 4: blue
		{R: 187, G: 154, B: 247}, // 5: magenta
		{R: 125, G: 207, B: 255}, // 6: cyan
		{R: 192, G: 202, B: 245}, // 7: white
		{R: 86, G: 95, B: 137},   // 8: bright black
		{R: 255, G: 137, B: 157}, // 9: bright red
		{R: 175, G: 214, B: 122}, // 10: bright green
		{R: 233, G: 197, B: 130}, // 11: bright yellow
		{R: 141, G: 182, B: 250}, // 12: bright blue
		{R: 200, G: 172, B: 248}, // 13: bright magenta
		{R: 151, G: 216, B: 248}, // 14: bright cyan
		{R: 200, G: 211, B: 245}, // 15: bright white
	}

	for i := 0; i < 16; i++ {
		if !terminal[i].Equal(expected[i]) {
			t.Errorf("Terminal[%d] = %v, want %v", i, terminal[i], expected[i])
		}
	}
}

func TestTheme_Get(t *testing.T) {
	theme := newTestTheme()

	tests := []struct {
		name      string
		path      string
		wantColor flair.Color
		wantOK    bool
	}{
		{
			name:      "existing accent.primary token",
			path:      "accent.primary",
			wantColor: flair.Color{R: 122, G: 162, B: 247},
			wantOK:    true,
		},
		{
			name:      "existing surface.background token",
			path:      "surface.background",
			wantColor: flair.Color{R: 26, G: 27, B: 38},
			wantOK:    true,
		},
		{
			name:      "non-existent token",
			path:      "does.not.exist",
			wantColor: flair.Color{},
			wantOK:    false,
		},
		{
			name:      "empty path",
			path:      "",
			wantColor: flair.Color{},
			wantOK:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := theme.Get(tt.path)
			if ok != tt.wantOK {
				t.Errorf("Theme.Get(%q) ok = %v, want %v", tt.path, ok, tt.wantOK)
			}
			if !got.Equal(tt.wantColor) {
				t.Errorf("Theme.Get(%q) = %v, want %v", tt.path, got, tt.wantColor)
			}
		})
	}
}

func TestTheme_Surface_MissingTokens(t *testing.T) {
	// Theme with minimal tokens should return zero colors for missing fields
	colors := map[string]flair.Color{
		"surface.background": {R: 26, G: 27, B: 38},
	}
	theme := flair.NewTheme("minimal", "dark", colors)

	surface := theme.Surface()

	// Background should be populated
	if !surface.Background.Equal(flair.Color{R: 26, G: 27, B: 38}) {
		t.Errorf("Surface.Background = %v, want {26 27 38}", surface.Background)
	}

	// Raised should be zero value since token is missing
	if !surface.Raised.Equal(flair.Color{}) {
		t.Errorf("Surface.Raised = %v, want zero value", surface.Raised)
	}
}

func TestTheme_Terminal_MissingTokens(t *testing.T) {
	// Theme with partial terminal tokens
	colors := map[string]flair.Color{
		"terminal.black": {R: 31, G: 35, B: 53},
		"terminal.red":   {R: 247, G: 118, B: 142},
		// Missing other terminal colors
	}
	theme := flair.NewTheme("partial", "dark", colors)

	terminal := theme.Terminal()

	// Index 0 (black) should be populated
	if !terminal[0].Equal(flair.Color{R: 31, G: 35, B: 53}) {
		t.Errorf("Terminal[0] = %v, want {31 35 53}", terminal[0])
	}

	// Index 1 (red) should be populated
	if !terminal[1].Equal(flair.Color{R: 247, G: 118, B: 142}) {
		t.Errorf("Terminal[1] = %v, want {247 118 142}", terminal[1])
	}

	// Index 2 (green) should be zero value since token is missing
	if !terminal[2].Equal(flair.Color{}) {
		t.Errorf("Terminal[2] = %v, want zero value", terminal[2])
	}
}

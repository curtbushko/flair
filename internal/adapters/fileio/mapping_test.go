package fileio_test

import (
	"bytes"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/ports"
)

// malformedYAML is a shared constant for tests that verify invalid YAML handling.
const malformedYAML = `{not valid yaml: [[[`

// TestWriteStylixMapping_ValidYAML verifies that WriteStylixMapping produces
// valid YAML that can be unmarshaled back to a StylixMappingFile.
func TestWriteStylixMapping_ValidYAML(t *testing.T) {
	mf := ports.StylixMappingFile{
		Values: map[string]string{
			"base00":         "#1a1b26",
			"base0D":         "#7aa2f7",
			"surface-bg":     "#1a1b26",
			"text-primary":   "#c0caf5",
			"syntax-keyword": "#bb9af7",
		},
	}

	var buf bytes.Buffer
	err := fileio.WriteStylixMapping(&buf, mf)
	if err != nil {
		t.Fatalf("WriteStylixMapping error: %v", err)
	}

	output := buf.String()

	// Output should be valid YAML
	var parsed ports.StylixMappingFile
	if err := yaml.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, output)
	}

	// Should contain all 5 values
	if len(parsed.Values) != 5 {
		t.Errorf("expected 5 values, got %d", len(parsed.Values))
	}

	// Verify specific values
	for key, want := range mf.Values {
		got, ok := parsed.Values[key]
		if !ok {
			t.Errorf("key %q missing in parsed output", key)
			continue
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", key, got, want)
		}
	}

	// Verify keys are sorted in the YAML output
	lines := strings.Split(output, "\n")
	var foundKeys []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Value keys appear as 4-space indented under "values:"
		if strings.HasPrefix(line, "    ") && strings.Contains(trimmed, ":") {
			key := strings.SplitN(trimmed, ":", 2)[0]
			foundKeys = append(foundKeys, key)
		}
	}

	for i := 1; i < len(foundKeys); i++ {
		if foundKeys[i] < foundKeys[i-1] {
			t.Errorf("keys not sorted: %v", foundKeys)
			break
		}
	}
}

// TestStylixMappingFile_RoundTrip verifies that writing a StylixMappingFile
// via WriteStylixMapping and reading it back via ReadStylixMapping produces
// an identical StylixMappingFile.
func TestStylixMappingFile_RoundTrip(t *testing.T) {
	original := ports.StylixMappingFile{
		Values: map[string]string{
			"base00":         "#1a1b26",
			"base01":         "#1f2335",
			"base02":         "#292e42",
			"base0D":         "#7aa2f7",
			"surface-bg":     "#1a1b26",
			"text-primary":   "#c0caf5",
			"syntax-keyword": "#bb9af7",
			"syntax-string":  "#9ece6a",
			"terminal-red":   "#f7768e",
			"accent-primary": "#7aa2f7",
		},
	}

	// Write
	var buf bytes.Buffer
	if err := fileio.WriteStylixMapping(&buf, original); err != nil {
		t.Fatalf("WriteStylixMapping error: %v", err)
	}

	// Read back
	restored, err := fileio.ReadStylixMapping(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadStylixMapping error: %v", err)
	}

	// Compare
	if len(original.Values) != len(restored.Values) {
		t.Fatalf("value count mismatch: original=%d, restored=%d",
			len(original.Values), len(restored.Values))
	}

	for key, want := range original.Values {
		got, ok := restored.Values[key]
		if !ok {
			t.Errorf("key %q missing in restored mapping", key)
			continue
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", key, got, want)
		}
	}
}

// TestReadStylixMapping_Valid verifies that ReadStylixMapping can parse
// a valid YAML input with values.
func TestReadStylixMapping_Valid(t *testing.T) {
	yamlData := `values:
  base00: "#1a1b26"
  surface-bg: "#1a1b26"
  text-primary: "#c0caf5"
`
	mf, err := fileio.ReadStylixMapping(strings.NewReader(yamlData))
	if err != nil {
		t.Fatalf("ReadStylixMapping error: %v", err)
	}

	if len(mf.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(mf.Values))
	}

	if got := mf.Values["base00"]; got != testBgHex {
		t.Errorf("base00 = %q, want %q", got, testBgHex)
	}
	if got := mf.Values["surface-bg"]; got != testBgHex {
		t.Errorf("surface-bg = %q, want %q", got, testBgHex)
	}
	if got := mf.Values["text-primary"]; got != "#c0caf5" {
		t.Errorf("text-primary = %q, want %q", got, "#c0caf5")
	}
}

// TestReadStylixMapping_InvalidYAML verifies that invalid YAML returns an error.
func TestReadStylixMapping_InvalidYAML(t *testing.T) {
	_, err := fileio.ReadStylixMapping(strings.NewReader(malformedYAML))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// TestWriteStylixMapping_Empty verifies that an empty values map
// still produces valid YAML output.
func TestWriteStylixMapping_Empty(t *testing.T) {
	mf := ports.StylixMappingFile{
		Values: map[string]string{},
	}

	var buf bytes.Buffer
	if err := fileio.WriteStylixMapping(&buf, mf); err != nil {
		t.Fatalf("WriteStylixMapping error: %v", err)
	}

	var parsed ports.StylixMappingFile
	if err := yaml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	if len(parsed.Values) != 0 {
		t.Errorf("expected empty values map, got %d entries", len(parsed.Values))
	}
}

// --------------------------------------------------------------------------
// CSS mapping file tests
// --------------------------------------------------------------------------

// TestWriteCSSMapping_ValidYAML verifies that WriteCSSMapping produces
// valid YAML that can be unmarshaled back to a CSSMappingFile.
func TestWriteCSSMapping_ValidYAML(t *testing.T) {
	mf := ports.CSSMappingFile{
		CustomProperties: map[string]string{
			"--flair-bg":             "#1a1b26",
			"--flair-fg":             "#c0caf5",
			"--flair-accent-primary": "#7aa2f7",
			"--flair-syntax-keyword": "#bb9af7",
		},
		Rules: []ports.CSSRuleEntry{
			{
				Selector: "body",
				Properties: map[string]string{
					"background-color": "var(--flair-bg)",
					"color":            "var(--flair-fg)",
				},
			},
			{
				Selector: "a",
				Properties: map[string]string{
					"color": "var(--flair-accent-primary)",
				},
			},
		},
	}

	var buf bytes.Buffer
	err := fileio.WriteCSSMapping(&buf, mf)
	if err != nil {
		t.Fatalf("WriteCSSMapping error: %v", err)
	}

	output := buf.String()

	// Output should be valid YAML.
	var parsed ports.CSSMappingFile
	if err := yaml.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, output)
	}

	// Should contain all 4 custom properties.
	if len(parsed.CustomProperties) != 4 {
		t.Errorf("expected 4 custom properties, got %d", len(parsed.CustomProperties))
	}

	// Verify specific values.
	for key, want := range mf.CustomProperties {
		got, ok := parsed.CustomProperties[key]
		if !ok {
			t.Errorf("custom property %q missing in parsed output", key)
			continue
		}
		if got != want {
			t.Errorf("custom property %q: got %q, want %q", key, got, want)
		}
	}

	// Should contain 2 rules.
	if len(parsed.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(parsed.Rules))
	}

	// Verify custom property keys are sorted in the YAML output.
	lines := strings.Split(output, "\n")
	var foundKeys []string
	inCustomProps := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "custom_properties:" {
			inCustomProps = true
			continue
		}
		if inCustomProps && strings.HasPrefix(line, "    ") && strings.Contains(trimmed, ":") {
			key := strings.SplitN(trimmed, ":", 2)[0]
			foundKeys = append(foundKeys, key)
		}
		if inCustomProps && !strings.HasPrefix(line, "    ") && trimmed != "" {
			inCustomProps = false
		}
	}

	for i := 1; i < len(foundKeys); i++ {
		if foundKeys[i] < foundKeys[i-1] {
			t.Errorf("custom property keys not sorted: %v", foundKeys)
			break
		}
	}
}

// TestCSSMappingFile_RoundTrip verifies that writing a CSSMappingFile
// via WriteCSSMapping and reading it back via ReadCSSMapping produces
// an identical CSSMappingFile.
func TestCSSMappingFile_RoundTrip(t *testing.T) {
	original := ports.CSSMappingFile{
		CustomProperties: map[string]string{
			"--flair-bg":               "#1a1b26",
			"--flair-fg":               "#c0caf5",
			"--flair-accent-primary":   "#7aa2f7",
			"--flair-syntax-keyword":   "#bb9af7",
			"--flair-syntax-string":    "#9ece6a",
			"--flair-status-error":     "#ff899d",
			"--flair-bg-selection":     "#3b4261",
			"--flair-bg-raised":        "#1f2335",
			"--flair-accent-secondary": "#bb9af7",
		},
		Rules: []ports.CSSRuleEntry{
			{
				Selector: "body",
				Properties: map[string]string{
					"background-color": "var(--flair-bg)",
					"color":            "var(--flair-fg)",
					"font-family":      "system-ui, -apple-system, sans-serif",
				},
			},
			{
				Selector: "a",
				Properties: map[string]string{
					"color": "var(--flair-accent-primary)",
				},
			},
			{
				Selector: "::selection",
				Properties: map[string]string{
					"background-color": "var(--flair-bg-selection)",
					"color":            "var(--flair-fg)",
				},
			},
			{
				Selector: "pre, code",
				Properties: map[string]string{
					"background-color": "var(--flair-bg-raised)",
					"color":            "var(--flair-fg)",
					"border-radius":    "4px",
				},
			},
		},
	}

	// Write.
	var buf bytes.Buffer
	if err := fileio.WriteCSSMapping(&buf, original); err != nil {
		t.Fatalf("WriteCSSMapping error: %v", err)
	}

	// Read back.
	restored, err := fileio.ReadCSSMapping(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadCSSMapping error: %v\nwritten:\n%s", err, buf.String())
	}

	// Compare custom properties.
	if len(original.CustomProperties) != len(restored.CustomProperties) {
		t.Fatalf("custom property count mismatch: original=%d, restored=%d",
			len(original.CustomProperties), len(restored.CustomProperties))
	}

	for key, want := range original.CustomProperties {
		got, ok := restored.CustomProperties[key]
		if !ok {
			t.Errorf("custom property %q missing in restored mapping", key)
			continue
		}
		if got != want {
			t.Errorf("custom property %q: got %q, want %q", key, got, want)
		}
	}

	// Compare rules.
	if len(original.Rules) != len(restored.Rules) {
		t.Fatalf("rule count mismatch: original=%d, restored=%d",
			len(original.Rules), len(restored.Rules))
	}

	for i, origRule := range original.Rules {
		restoredRule := restored.Rules[i]
		if origRule.Selector != restoredRule.Selector {
			t.Errorf("rule %d selector: got %q, want %q", i, restoredRule.Selector, origRule.Selector)
		}
		if len(origRule.Properties) != len(restoredRule.Properties) {
			t.Errorf("rule %d property count mismatch: original=%d, restored=%d",
				i, len(origRule.Properties), len(restoredRule.Properties))
			continue
		}
		for propKey, wantVal := range origRule.Properties {
			gotVal, ok := restoredRule.Properties[propKey]
			if !ok {
				t.Errorf("rule %d (%q) missing property %q", i, origRule.Selector, propKey)
				continue
			}
			if gotVal != wantVal {
				t.Errorf("rule %d (%q) property %q: got %q, want %q",
					i, origRule.Selector, propKey, gotVal, wantVal)
			}
		}
	}
}

// TestReadCSSMapping_Valid verifies that ReadCSSMapping can parse
// a valid YAML input with custom_properties and rules.
func TestReadCSSMapping_Valid(t *testing.T) {
	yamlData := `custom_properties:
  --flair-bg: "#1a1b26"
  --flair-fg: "#c0caf5"
rules:
  - selector: "body"
    properties:
      background-color: "var(--flair-bg)"
      color: "var(--flair-fg)"
`
	mf, err := fileio.ReadCSSMapping(strings.NewReader(yamlData))
	if err != nil {
		t.Fatalf("ReadCSSMapping error: %v", err)
	}

	if len(mf.CustomProperties) != 2 {
		t.Fatalf("expected 2 custom properties, got %d", len(mf.CustomProperties))
	}

	if got := mf.CustomProperties["--flair-bg"]; got != "#1a1b26" {
		t.Errorf("--flair-bg = %q, want %q", got, "#1a1b26")
	}

	if len(mf.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mf.Rules))
	}

	if mf.Rules[0].Selector != "body" {
		t.Errorf("rule selector = %q, want %q", mf.Rules[0].Selector, "body")
	}
}

// TestReadCSSMapping_InvalidYAML verifies that invalid YAML returns an error.
func TestReadCSSMapping_InvalidYAML(t *testing.T) {
	_, err := fileio.ReadCSSMapping(strings.NewReader(malformedYAML))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// TestWriteCSSMapping_Empty verifies that empty custom properties and rules
// still produce valid YAML output.
func TestWriteCSSMapping_Empty(t *testing.T) {
	mf := ports.CSSMappingFile{
		CustomProperties: map[string]string{},
		Rules:            []ports.CSSRuleEntry{},
	}

	var buf bytes.Buffer
	if err := fileio.WriteCSSMapping(&buf, mf); err != nil {
		t.Fatalf("WriteCSSMapping error: %v", err)
	}

	var parsed ports.CSSMappingFile
	if err := yaml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	if len(parsed.CustomProperties) != 0 {
		t.Errorf("expected empty custom properties, got %d entries", len(parsed.CustomProperties))
	}
	if len(parsed.Rules) != 0 {
		t.Errorf("expected empty rules, got %d entries", len(parsed.Rules))
	}
}

// --------------------------------------------------------------------------
// Vim mapping file tests
// --------------------------------------------------------------------------

// TestVimMappingFile_RoundTrip verifies that writing a VimMappingFile via
// WriteVimMapping and reading it back via ReadVimMapping produces an
// identical VimMappingFile with all highlights and terminal colors preserved.
func TestVimMappingFile_RoundTrip(t *testing.T) {
	original := ports.VimMappingFile{
		Highlights: map[string]ports.VimMappingHighlight{
			"Normal":   {Fg: "text.primary", Bg: "surface.background"},
			"Comment":  {Fg: "syntax.comment", Italic: true},
			"Keyword":  {Fg: "syntax.keyword", Bold: true},
			"String":   {Fg: "syntax.string"},
			"Function": {Fg: "syntax.function"},
			"Type":     {Fg: "syntax.type"},
			"SpellBad": {Sp: "status.error", Undercurl: true},
			"Visual":   {Bg: "surface.background.selection"},
			"ErrorMsg": {Fg: "status.error", Bold: true, Strikethrough: true},
			"@comment": {Link: "Comment"},
		},
		TerminalColors: [16]string{
			"#1a1b26", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#a9b1d6",
			"#414868", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#c0caf5",
		},
	}

	// Write.
	var buf bytes.Buffer
	if err := fileio.WriteVimMapping(&buf, original); err != nil {
		t.Fatalf("WriteVimMapping error: %v", err)
	}

	// Read back.
	restored, err := fileio.ReadVimMapping(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadVimMapping error: %v\nwritten:\n%s", err, buf.String())
	}

	// Compare highlights.
	if len(original.Highlights) != len(restored.Highlights) {
		t.Fatalf("highlight count mismatch: original=%d, restored=%d",
			len(original.Highlights), len(restored.Highlights))
	}

	for name, origHL := range original.Highlights {
		restoredHL, ok := restored.Highlights[name]
		if !ok {
			t.Errorf("highlight %q missing in restored mapping", name)
			continue
		}
		if origHL != restoredHL {
			t.Errorf("highlight %q mismatch:\n  got  %+v\n  want %+v", name, restoredHL, origHL)
		}
	}

	// Compare terminal colors.
	if original.TerminalColors != restored.TerminalColors {
		t.Errorf("terminal colors mismatch:\n  got  %v\n  want %v",
			restored.TerminalColors, original.TerminalColors)
	}

	// Verify highlight keys are sorted in the YAML output.
	output := buf.String()
	lines := strings.Split(output, "\n")
	var foundKeys []string
	inHighlights := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "highlights:" {
			inHighlights = true
			continue
		}
		if inHighlights && strings.HasPrefix(line, "    ") && !strings.HasPrefix(line, "        ") {
			key := strings.SplitN(trimmed, ":", 2)[0]
			foundKeys = append(foundKeys, key)
		}
		if inHighlights && !strings.HasPrefix(line, "    ") && trimmed != "" {
			inHighlights = false
		}
	}

	for i := 1; i < len(foundKeys); i++ {
		if foundKeys[i] < foundKeys[i-1] {
			t.Errorf("highlight keys not sorted: %v", foundKeys)
			break
		}
	}
}

// TestReadVimMapping_InvalidYAML verifies that invalid YAML returns an error.
func TestReadVimMapping_InvalidYAML(t *testing.T) {
	_, err := fileio.ReadVimMapping(strings.NewReader(malformedYAML))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// TestWriteVimMapping_Empty verifies that empty highlights and zero terminal
// colors still produce valid YAML output.
func TestWriteVimMapping_Empty(t *testing.T) {
	mf := ports.VimMappingFile{
		Highlights:     map[string]ports.VimMappingHighlight{},
		TerminalColors: [16]string{},
	}

	var buf bytes.Buffer
	if err := fileio.WriteVimMapping(&buf, mf); err != nil {
		t.Fatalf("WriteVimMapping error: %v", err)
	}

	var parsed ports.VimMappingFile
	if err := yaml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	if len(parsed.Highlights) != 0 {
		t.Errorf("expected empty highlights, got %d entries", len(parsed.Highlights))
	}
}

// --------------------------------------------------------------------------
// GTK mapping file tests
// --------------------------------------------------------------------------

// TestGtkMappingFile_RoundTrip verifies that writing a GtkMappingFile via
// WriteGtkMapping and reading it back via ReadGtkMapping produces an
// identical GtkMappingFile with all colors and rules preserved.
func TestGtkMappingFile_RoundTrip(t *testing.T) {
	original := ports.GtkMappingFile{
		Colors: map[string]string{
			"window_bg_color":    "#1a1b26",
			"window_fg_color":    "#c0caf5",
			"headerbar_bg_color": "#16161e",
			"headerbar_fg_color": "#c0caf5",
			"accent_bg_color":    "#7aa2f7",
			"accent_fg_color":    "#1a1b26",
			"card_bg_color":      "#1f2335",
			"card_fg_color":      "#c0caf5",
			"error_color":        "#ff899d",
		},
		Rules: []ports.CSSRuleEntry{
			{
				Selector: "window",
				Properties: map[string]string{
					"background-color": "@window_bg_color",
					"color":            "@window_fg_color",
				},
			},
			{
				Selector: "headerbar",
				Properties: map[string]string{
					"background-color": "@headerbar_bg_color",
					"color":            "@headerbar_fg_color",
				},
			},
			{
				Selector: "button",
				Properties: map[string]string{
					"background-color": "@card_bg_color",
					"color":            "@window_fg_color",
				},
			},
		},
	}

	// Write.
	var buf bytes.Buffer
	if err := fileio.WriteGtkMapping(&buf, original); err != nil {
		t.Fatalf("WriteGtkMapping error: %v", err)
	}

	// Read back.
	restored, err := fileio.ReadGtkMapping(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadGtkMapping error: %v\nwritten:\n%s", err, buf.String())
	}

	// Compare colors.
	if len(original.Colors) != len(restored.Colors) {
		t.Fatalf("color count mismatch: original=%d, restored=%d",
			len(original.Colors), len(restored.Colors))
	}

	for key, want := range original.Colors {
		got, ok := restored.Colors[key]
		if !ok {
			t.Errorf("color %q missing in restored mapping", key)
			continue
		}
		if got != want {
			t.Errorf("color %q: got %q, want %q", key, got, want)
		}
	}

	// Compare rules.
	if len(original.Rules) != len(restored.Rules) {
		t.Fatalf("rule count mismatch: original=%d, restored=%d",
			len(original.Rules), len(restored.Rules))
	}

	for i, origRule := range original.Rules {
		restoredRule := restored.Rules[i]
		if origRule.Selector != restoredRule.Selector {
			t.Errorf("rule %d selector: got %q, want %q", i, restoredRule.Selector, origRule.Selector)
		}
		if len(origRule.Properties) != len(restoredRule.Properties) {
			t.Errorf("rule %d property count mismatch: original=%d, restored=%d",
				i, len(origRule.Properties), len(restoredRule.Properties))
			continue
		}
		for propKey, wantVal := range origRule.Properties {
			gotVal, ok := restoredRule.Properties[propKey]
			if !ok {
				t.Errorf("rule %d (%q) missing property %q", i, origRule.Selector, propKey)
				continue
			}
			if gotVal != wantVal {
				t.Errorf("rule %d (%q) property %q: got %q, want %q",
					i, origRule.Selector, propKey, gotVal, wantVal)
			}
		}
	}

	// Verify color keys are sorted in the YAML output.
	output := buf.String()
	lines := strings.Split(output, "\n")
	var foundKeys []string
	inColors := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "colors:" {
			inColors = true
			continue
		}
		if inColors && strings.HasPrefix(line, "    ") && strings.Contains(trimmed, ":") {
			key := strings.SplitN(trimmed, ":", 2)[0]
			foundKeys = append(foundKeys, key)
		}
		if inColors && !strings.HasPrefix(line, "    ") && trimmed != "" {
			inColors = false
		}
	}

	for i := 1; i < len(foundKeys); i++ {
		if foundKeys[i] < foundKeys[i-1] {
			t.Errorf("color keys not sorted: %v", foundKeys)
			break
		}
	}
}

// TestReadGtkMapping_Valid verifies that ReadGtkMapping can parse
// a valid YAML input with colors and rules.
func TestReadGtkMapping_Valid(t *testing.T) {
	yamlData := `colors:
  window_bg_color: "#1a1b26"
  window_fg_color: "#c0caf5"
rules:
  - selector: "window"
    properties:
      background-color: "@window_bg_color"
      color: "@window_fg_color"
`
	mf, err := fileio.ReadGtkMapping(strings.NewReader(yamlData))
	if err != nil {
		t.Fatalf("ReadGtkMapping error: %v", err)
	}

	if len(mf.Colors) != 2 {
		t.Fatalf("expected 2 colors, got %d", len(mf.Colors))
	}

	if got := mf.Colors["window_bg_color"]; got != "#1a1b26" {
		t.Errorf("window_bg_color = %q, want %q", got, "#1a1b26")
	}

	if len(mf.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mf.Rules))
	}

	if mf.Rules[0].Selector != "window" {
		t.Errorf("rule selector = %q, want %q", mf.Rules[0].Selector, "window")
	}
}

// TestReadGtkMapping_InvalidYAML verifies that invalid YAML returns an error.
func TestReadGtkMapping_InvalidYAML(t *testing.T) {
	_, err := fileio.ReadGtkMapping(strings.NewReader(malformedYAML))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// TestWriteGtkMapping_Empty verifies that empty colors and rules
// still produce valid YAML output.
func TestWriteGtkMapping_Empty(t *testing.T) {
	mf := ports.GtkMappingFile{
		Colors: map[string]string{},
		Rules:  []ports.CSSRuleEntry{},
	}

	var buf bytes.Buffer
	if err := fileio.WriteGtkMapping(&buf, mf); err != nil {
		t.Fatalf("WriteGtkMapping error: %v", err)
	}

	var parsed ports.GtkMappingFile
	if err := yaml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	if len(parsed.Colors) != 0 {
		t.Errorf("expected empty colors, got %d entries", len(parsed.Colors))
	}
	if len(parsed.Rules) != 0 {
		t.Errorf("expected empty rules, got %d entries", len(parsed.Rules))
	}
}

// --------------------------------------------------------------------------
// QSS mapping file tests
// --------------------------------------------------------------------------

// TestQssMappingFile_RoundTrip verifies that writing a QssMappingFile via
// WriteQssMapping and reading it back via ReadQssMapping produces an
// identical QssMappingFile with all rules preserved.
func TestQssMappingFile_RoundTrip(t *testing.T) {
	original := ports.QssMappingFile{
		Rules: []ports.CSSRuleEntry{
			{
				Selector: "QWidget",
				Properties: map[string]string{
					"background-color": "#1a1b26",
					"color":            "#c0caf5",
				},
			},
			{
				Selector: "QMainWindow",
				Properties: map[string]string{
					"background-color": "#1a1b26",
					"color":            "#c0caf5",
				},
			},
			{
				Selector: "QPushButton",
				Properties: map[string]string{
					"background-color": "#1f2335",
					"color":            "#c0caf5",
					"border":           "1px solid #3b4261",
				},
			},
			{
				Selector: "QPushButton:hover",
				Properties: map[string]string{
					"background-color": "#292e42",
					"color":            "#c0caf5",
				},
			},
			{
				Selector: "QPushButton:pressed",
				Properties: map[string]string{
					"background-color": "#3b4261",
					"color":            "#c0caf5",
				},
			},
			{
				Selector: "QLineEdit",
				Properties: map[string]string{
					"background-color": "#16161e",
					"color":            "#c0caf5",
					"border":           "1px solid #3b4261",
				},
			},
			{
				Selector: "QLineEdit:focus",
				Properties: map[string]string{
					"border": "1px solid #5c77bb",
				},
			},
		},
	}

	// Write.
	var buf bytes.Buffer
	if err := fileio.WriteQssMapping(&buf, original); err != nil {
		t.Fatalf("WriteQssMapping error: %v", err)
	}

	// Read back.
	restored, err := fileio.ReadQssMapping(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadQssMapping error: %v\nwritten:\n%s", err, buf.String())
	}

	// Compare rules.
	if len(original.Rules) != len(restored.Rules) {
		t.Fatalf("rule count mismatch: original=%d, restored=%d",
			len(original.Rules), len(restored.Rules))
	}

	for i, origRule := range original.Rules {
		restoredRule := restored.Rules[i]
		if origRule.Selector != restoredRule.Selector {
			t.Errorf("rule %d selector: got %q, want %q", i, restoredRule.Selector, origRule.Selector)
		}
		if len(origRule.Properties) != len(restoredRule.Properties) {
			t.Errorf("rule %d property count mismatch: original=%d, restored=%d",
				i, len(origRule.Properties), len(restoredRule.Properties))
			continue
		}
		for propKey, wantVal := range origRule.Properties {
			gotVal, ok := restoredRule.Properties[propKey]
			if !ok {
				t.Errorf("rule %d (%q) missing property %q", i, origRule.Selector, propKey)
				continue
			}
			if gotVal != wantVal {
				t.Errorf("rule %d (%q) property %q: got %q, want %q",
					i, origRule.Selector, propKey, gotVal, wantVal)
			}
		}
	}
}

// TestReadQssMapping_Valid verifies that ReadQssMapping can parse
// a valid YAML input with rules.
func TestReadQssMapping_Valid(t *testing.T) {
	yamlData := `rules:
  - selector: "QWidget"
    properties:
      background-color: "#1a1b26"
      color: "#c0caf5"
  - selector: "QPushButton:hover"
    properties:
      background-color: "#292e42"
`
	mf, err := fileio.ReadQssMapping(strings.NewReader(yamlData))
	if err != nil {
		t.Fatalf("ReadQssMapping error: %v", err)
	}

	if len(mf.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(mf.Rules))
	}

	if mf.Rules[0].Selector != "QWidget" {
		t.Errorf("rule 0 selector = %q, want %q", mf.Rules[0].Selector, "QWidget")
	}

	if mf.Rules[1].Selector != "QPushButton:hover" {
		t.Errorf("rule 1 selector = %q, want %q", mf.Rules[1].Selector, "QPushButton:hover")
	}
}

// TestReadQssMapping_InvalidYAML verifies that invalid YAML returns an error.
func TestReadQssMapping_InvalidYAML(t *testing.T) {
	_, err := fileio.ReadQssMapping(strings.NewReader(malformedYAML))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// TestWriteQssMapping_Empty verifies that empty rules
// still produce valid YAML output.
func TestWriteQssMapping_Empty(t *testing.T) {
	mf := ports.QssMappingFile{
		Rules: []ports.CSSRuleEntry{},
	}

	var buf bytes.Buffer
	if err := fileio.WriteQssMapping(&buf, mf); err != nil {
		t.Fatalf("WriteQssMapping error: %v", err)
	}

	var parsed ports.QssMappingFile
	if err := yaml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, buf.String())
	}

	if len(parsed.Rules) != 0 {
		t.Errorf("expected empty rules, got %d entries", len(parsed.Rules))
	}
}

// --------------------------------------------------------------------------
// VimMappingFile bufferline tests
// --------------------------------------------------------------------------

// TestVimMappingFile_BufferlineRoundTrip verifies that writing a VimMappingFile
// with Bufferline populated via WriteVimMapping and reading it back via
// ReadVimMapping produces an identical VimMappingFile with all bufferline
// highlight groups preserved.
func TestVimMappingFile_BufferlineRoundTrip(t *testing.T) {
	original := ports.VimMappingFile{
		Highlights: map[string]ports.VimMappingHighlight{
			"Normal":  {Fg: "#c0caf5", Bg: "#1a1b26"},
			"Comment": {Fg: "#565f89", Italic: true},
		},
		TerminalColors: [16]string{
			"#1a1b26", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#a9b1d6",
			"#414868", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#c0caf5",
		},
		Bufferline: &ports.BufferlineMappingTheme{
			Fill:              ports.BufferlineMappingColors{Fg: "#c0caf5", Bg: "#16161e"},
			Background:        ports.BufferlineMappingColors{Fg: "#565f89", Bg: "#16161e"},
			BufferVisible:     ports.BufferlineMappingColors{Fg: "#565f89", Bg: "#16161e"},
			BufferSelected:    ports.BufferlineMappingColors{Fg: "#c0caf5", Bg: "#1a1b26", Bold: true},
			Separator:         ports.BufferlineMappingColors{Fg: "#16161e", Bg: "#16161e"},
			SeparatorVisible:  ports.BufferlineMappingColors{Fg: "#16161e", Bg: "#16161e"},
			SeparatorSelected: ports.BufferlineMappingColors{Fg: "#16161e", Bg: "#1a1b26"},
			IndicatorSelected: ports.BufferlineMappingColors{Fg: "#7aa2f7", Bg: "#1a1b26"},
			Modified:          ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#16161e"},
			ModifiedVisible:   ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#16161e"},
			ModifiedSelected:  ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#1a1b26"},
			Error:             ports.BufferlineMappingColors{Fg: "#f7768e", Bg: "#16161e"},
			Warning:           ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#16161e"},
			Info:              ports.BufferlineMappingColors{Fg: "#7aa2f7", Bg: "#16161e"},
			Hint:              ports.BufferlineMappingColors{Fg: "#1abc9c", Bg: "#16161e"},
		},
	}

	// Write.
	var buf bytes.Buffer
	if err := fileio.WriteVimMapping(&buf, original); err != nil {
		t.Fatalf("WriteVimMapping error: %v", err)
	}

	// Read back.
	restored, err := fileio.ReadVimMapping(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadVimMapping error: %v\nwritten:\n%s", err, buf.String())
	}

	// Verify bufferline is not nil.
	if restored.Bufferline == nil {
		t.Fatal("restored.Bufferline is nil, expected non-nil")
	}

	// Compare all bufferline groups.
	bl := original.Bufferline
	rbl := restored.Bufferline

	assertBufferlineColors := func(name string, orig, rest ports.BufferlineMappingColors) {
		t.Helper()
		if orig.Fg != rest.Fg {
			t.Errorf("%s.Fg: got %q, want %q", name, rest.Fg, orig.Fg)
		}
		if orig.Bg != rest.Bg {
			t.Errorf("%s.Bg: got %q, want %q", name, rest.Bg, orig.Bg)
		}
		if orig.Bold != rest.Bold {
			t.Errorf("%s.Bold: got %v, want %v", name, rest.Bold, orig.Bold)
		}
		if orig.Italic != rest.Italic {
			t.Errorf("%s.Italic: got %v, want %v", name, rest.Italic, orig.Italic)
		}
	}

	assertBufferlineColors("Fill", bl.Fill, rbl.Fill)
	assertBufferlineColors("Background", bl.Background, rbl.Background)
	assertBufferlineColors("BufferVisible", bl.BufferVisible, rbl.BufferVisible)
	assertBufferlineColors("BufferSelected", bl.BufferSelected, rbl.BufferSelected)
	assertBufferlineColors("Separator", bl.Separator, rbl.Separator)
	assertBufferlineColors("SeparatorVisible", bl.SeparatorVisible, rbl.SeparatorVisible)
	assertBufferlineColors("SeparatorSelected", bl.SeparatorSelected, rbl.SeparatorSelected)
	assertBufferlineColors("IndicatorSelected", bl.IndicatorSelected, rbl.IndicatorSelected)
	assertBufferlineColors("Modified", bl.Modified, rbl.Modified)
	assertBufferlineColors("ModifiedVisible", bl.ModifiedVisible, rbl.ModifiedVisible)
	assertBufferlineColors("ModifiedSelected", bl.ModifiedSelected, rbl.ModifiedSelected)
	assertBufferlineColors("Error", bl.Error, rbl.Error)
	assertBufferlineColors("Warning", bl.Warning, rbl.Warning)
	assertBufferlineColors("Info", bl.Info, rbl.Info)
	assertBufferlineColors("Hint", bl.Hint, rbl.Hint)
}

// TestVimMappingFile_NilBufferline verifies that a VimMappingFile with nil
// Bufferline still works correctly (no bufferline section in YAML output).
func TestVimMappingFile_NilBufferline(t *testing.T) {
	original := ports.VimMappingFile{
		Highlights: map[string]ports.VimMappingHighlight{
			"Normal": {Fg: "#c0caf5", Bg: "#1a1b26"},
		},
		TerminalColors: [16]string{
			"#1a1b26", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#a9b1d6",
			"#414868", "#f7768e", "#9ece6a", "#e0af68",
			"#7aa2f7", "#bb9af7", "#7dcfff", "#c0caf5",
		},
		Bufferline: nil,
	}

	// Write.
	var buf bytes.Buffer
	if err := fileio.WriteVimMapping(&buf, original); err != nil {
		t.Fatalf("WriteVimMapping error: %v", err)
	}

	output := buf.String()

	// Verify no bufferline section in output.
	if strings.Contains(output, "bufferline:") {
		t.Errorf("YAML output should not contain bufferline section when nil:\n%s", output)
	}

	// Read back.
	restored, err := fileio.ReadVimMapping(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadVimMapping error: %v", err)
	}

	// Verify bufferline remains nil.
	if restored.Bufferline != nil {
		t.Errorf("restored.Bufferline should be nil, got %+v", restored.Bufferline)
	}
}

// TestVimMappingFile_BufferlineSorted verifies that bufferline keys are
// serialized in a deterministic order.
func TestVimMappingFile_BufferlineSorted(t *testing.T) {
	mf := ports.VimMappingFile{
		Highlights: map[string]ports.VimMappingHighlight{
			"Normal": {Fg: "#c0caf5", Bg: "#1a1b26"},
		},
		TerminalColors: [16]string{},
		Bufferline: &ports.BufferlineMappingTheme{
			Fill:              ports.BufferlineMappingColors{Fg: "#c0caf5", Bg: "#16161e"},
			Background:        ports.BufferlineMappingColors{Fg: "#565f89", Bg: "#16161e"},
			BufferVisible:     ports.BufferlineMappingColors{Fg: "#565f89", Bg: "#16161e"},
			BufferSelected:    ports.BufferlineMappingColors{Fg: "#c0caf5", Bg: "#1a1b26", Bold: true},
			Separator:         ports.BufferlineMappingColors{Fg: "#16161e", Bg: "#16161e"},
			SeparatorVisible:  ports.BufferlineMappingColors{Fg: "#16161e", Bg: "#16161e"},
			SeparatorSelected: ports.BufferlineMappingColors{Fg: "#16161e", Bg: "#1a1b26"},
			IndicatorSelected: ports.BufferlineMappingColors{Fg: "#7aa2f7", Bg: "#1a1b26"},
			Modified:          ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#16161e"},
			ModifiedVisible:   ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#16161e"},
			ModifiedSelected:  ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#1a1b26"},
			Error:             ports.BufferlineMappingColors{Fg: "#f7768e", Bg: "#16161e"},
			Warning:           ports.BufferlineMappingColors{Fg: "#e0af68", Bg: "#16161e"},
			Info:              ports.BufferlineMappingColors{Fg: "#7aa2f7", Bg: "#16161e"},
			Hint:              ports.BufferlineMappingColors{Fg: "#1abc9c", Bg: "#16161e"},
		},
	}

	// Write twice and verify output is identical (deterministic).
	var buf1, buf2 bytes.Buffer
	if err := fileio.WriteVimMapping(&buf1, mf); err != nil {
		t.Fatalf("WriteVimMapping (1) error: %v", err)
	}
	if err := fileio.WriteVimMapping(&buf2, mf); err != nil {
		t.Fatalf("WriteVimMapping (2) error: %v", err)
	}

	if buf1.String() != buf2.String() {
		t.Errorf("output is not deterministic:\nfirst:\n%s\nsecond:\n%s", buf1.String(), buf2.String())
	}

	// Verify bufferline section appears after terminal_colors.
	output := buf1.String()
	tcIdx := strings.Index(output, "terminal_colors:")
	blIdx := strings.Index(output, "bufferline:")

	if tcIdx == -1 {
		t.Error("terminal_colors section not found")
	}
	if blIdx == -1 {
		t.Error("bufferline section not found")
	}
	if tcIdx != -1 && blIdx != -1 && blIdx < tcIdx {
		t.Error("bufferline should appear after terminal_colors")
	}
}

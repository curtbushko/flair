package fileio

import (
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/curtbushko/flair/internal/ports"
)

// WriteStylixMapping serializes a StylixMappingFile as YAML to w with sorted
// keys for deterministic output. The caller is responsible for wrapping w with
// a VersionedWriter if schema version headers are desired.
func WriteStylixMapping(w io.Writer, mf ports.StylixMappingFile) error {
	out, err := marshalSortedStylixMapping(mf)
	if err != nil {
		return fmt.Errorf("write stylix mapping: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("write stylix mapping: %w", err)
	}

	return nil
}

// ReadStylixMapping reads YAML from r and returns a StylixMappingFile.
// The caller may wrap r with a ValidatingReader to enforce schema version
// checking before the YAML is decoded.
func ReadStylixMapping(r io.Reader) (ports.StylixMappingFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return ports.StylixMappingFile{}, fmt.Errorf("read stylix mapping: %w", err)
	}

	var mf ports.StylixMappingFile
	if err := yaml.Unmarshal(data, &mf); err != nil {
		return ports.StylixMappingFile{}, fmt.Errorf("read stylix mapping: %w", err)
	}

	return mf, nil
}

// WriteCSSMapping serializes a CSSMappingFile as YAML to w with sorted
// custom property keys for deterministic output. The caller is responsible
// for wrapping w with a VersionedWriter if schema version headers are desired.
func WriteCSSMapping(w io.Writer, mf ports.CSSMappingFile) error {
	out, err := marshalSortedCSSMapping(mf)
	if err != nil {
		return fmt.Errorf("write css mapping: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("write css mapping: %w", err)
	}

	return nil
}

// ReadCSSMapping reads YAML from r and returns a CSSMappingFile.
// The caller may wrap r with a ValidatingReader to enforce schema version
// checking before the YAML is decoded.
func ReadCSSMapping(r io.Reader) (ports.CSSMappingFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return ports.CSSMappingFile{}, fmt.Errorf("read css mapping: %w", err)
	}

	var mf ports.CSSMappingFile
	if err := yaml.Unmarshal(data, &mf); err != nil {
		return ports.CSSMappingFile{}, fmt.Errorf("read css mapping: %w", err)
	}

	return mf, nil
}

// WriteVimMapping serializes a VimMappingFile as YAML to w with sorted
// highlight keys for deterministic output. The caller is responsible for
// wrapping w with a VersionedWriter if schema version headers are desired.
func WriteVimMapping(w io.Writer, mf ports.VimMappingFile) error {
	out, err := marshalSortedVimMapping(mf)
	if err != nil {
		return fmt.Errorf("write vim mapping: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("write vim mapping: %w", err)
	}

	return nil
}

// ReadVimMapping reads YAML from r and returns a VimMappingFile.
// The caller may wrap r with a ValidatingReader to enforce schema version
// checking before the YAML is decoded.
func ReadVimMapping(r io.Reader) (ports.VimMappingFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return ports.VimMappingFile{}, fmt.Errorf("read vim mapping: %w", err)
	}

	var mf ports.VimMappingFile
	if err := yaml.Unmarshal(data, &mf); err != nil {
		return ports.VimMappingFile{}, fmt.Errorf("read vim mapping: %w", err)
	}

	return mf, nil
}

// WriteGtkMapping serializes a GtkMappingFile as YAML to w with sorted
// color keys for deterministic output. The caller is responsible for wrapping
// w with a VersionedWriter if schema version headers are desired.
func WriteGtkMapping(w io.Writer, mf ports.GtkMappingFile) error {
	out, err := marshalSortedGtkMapping(mf)
	if err != nil {
		return fmt.Errorf("write gtk mapping: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("write gtk mapping: %w", err)
	}

	return nil
}

// ReadGtkMapping reads YAML from r and returns a GtkMappingFile.
// The caller may wrap r with a ValidatingReader to enforce schema version
// checking before the YAML is decoded.
func ReadGtkMapping(r io.Reader) (ports.GtkMappingFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return ports.GtkMappingFile{}, fmt.Errorf("read gtk mapping: %w", err)
	}

	var mf ports.GtkMappingFile
	if err := yaml.Unmarshal(data, &mf); err != nil {
		return ports.GtkMappingFile{}, fmt.Errorf("read gtk mapping: %w", err)
	}

	return mf, nil
}

// marshalSortedStylixMapping encodes a StylixMappingFile to YAML with sorted
// value keys. We build the YAML node tree manually to guarantee key ordering.
func marshalSortedStylixMapping(mf ports.StylixMappingFile) ([]byte, error) {
	// Build the values mapping node with sorted keys.
	valuesNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	keys := make([]string, 0, len(mf.Values))
	for k := range mf.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		valuesNode.Content = append(valuesNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: mf.Values[key], Tag: "!!str"},
		)
	}

	// Build the root document node.
	root := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "values", Tag: "!!str"},
		valuesNode,
	)

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}

	return yaml.Marshal(doc)
}

// marshalSortedCSSMapping encodes a CSSMappingFile to YAML with sorted
// custom property keys for deterministic output. We build the YAML node
// tree manually to guarantee key ordering in the custom_properties section.
func marshalSortedCSSMapping(mf ports.CSSMappingFile) ([]byte, error) {
	root := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	propsNode := marshalSortedMap(mf.CustomProperties)

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "custom_properties", Tag: "!!str"},
		propsNode,
	)

	rulesNode := marshalCSSRuleEntries(mf.Rules)
	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "rules", Tag: "!!str"},
		rulesNode,
	)

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}

	return yaml.Marshal(doc)
}

// marshalSortedGtkMapping encodes a GtkMappingFile to YAML with sorted
// color keys for deterministic output. The structure mirrors the CSS mapping
// format: a colors mapping node followed by a rules sequence node.
func marshalSortedGtkMapping(mf ports.GtkMappingFile) ([]byte, error) {
	root := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	colorsNode := marshalSortedMap(mf.Colors)

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "colors", Tag: "!!str"},
		colorsNode,
	)

	rulesNode := marshalCSSRuleEntries(mf.Rules)
	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "rules", Tag: "!!str"},
		rulesNode,
	)

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}

	return yaml.Marshal(doc)
}

// marshalSortedMap builds a YAML mapping node from a map[string]string with
// keys sorted alphabetically for deterministic output.
func marshalSortedMap(m map[string]string) *yaml.Node {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: m[key], Tag: "!!str"},
		)
	}

	return node
}

// marshalCSSRuleEntries builds a YAML sequence node from a slice of
// CSSRuleEntry, with properties sorted by key for deterministic output.
func marshalCSSRuleEntries(rules []ports.CSSRuleEntry) *yaml.Node {
	rulesNode := &yaml.Node{
		Kind: yaml.SequenceNode,
		Tag:  "!!seq",
	}

	for _, rule := range rules {
		ruleMapNode := &yaml.Node{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
		}

		// selector field.
		ruleMapNode.Content = append(ruleMapNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "selector", Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: rule.Selector, Tag: "!!str"},
		)

		// properties mapping with sorted keys.
		rulePropsNode := marshalSortedMap(rule.Properties)

		ruleMapNode.Content = append(ruleMapNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "properties", Tag: "!!str"},
			rulePropsNode,
		)

		rulesNode.Content = append(rulesNode.Content, ruleMapNode)
	}

	return rulesNode
}

// WriteQssMapping serializes a QssMappingFile as YAML to w with sorted
// rule property keys for deterministic output. The caller is responsible for
// wrapping w with a VersionedWriter if schema version headers are desired.
func WriteQssMapping(w io.Writer, mf ports.QssMappingFile) error {
	out, err := marshalQssMapping(mf)
	if err != nil {
		return fmt.Errorf("write qss mapping: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("write qss mapping: %w", err)
	}

	return nil
}

// ReadQssMapping reads YAML from r and returns a QssMappingFile.
// The caller may wrap r with a ValidatingReader to enforce schema version
// checking before the YAML is decoded.
func ReadQssMapping(r io.Reader) (ports.QssMappingFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return ports.QssMappingFile{}, fmt.Errorf("read qss mapping: %w", err)
	}

	var mf ports.QssMappingFile
	if err := yaml.Unmarshal(data, &mf); err != nil {
		return ports.QssMappingFile{}, fmt.Errorf("read qss mapping: %w", err)
	}

	return mf, nil
}

// marshalQssMapping encodes a QssMappingFile to YAML with sorted rule
// property keys for deterministic output. The QSS mapping contains only
// rules (no color definitions or custom properties).
func marshalQssMapping(mf ports.QssMappingFile) ([]byte, error) {
	root := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	rulesNode := marshalCSSRuleEntries(mf.Rules)
	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "rules", Tag: "!!str"},
		rulesNode,
	)

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}

	return yaml.Marshal(doc)
}

// marshalSortedVimMapping encodes a VimMappingFile to YAML with sorted
// highlight keys for deterministic output. We build the YAML node tree
// manually to guarantee key ordering in the highlights section.
//
//nolint:funlen // Large mapping serialization is intentionally in one function.
func marshalSortedVimMapping(mf ports.VimMappingFile) ([]byte, error) {
	root := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	// Build highlights mapping node with sorted keys.
	hlNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	hlKeys := make([]string, 0, len(mf.Highlights))
	for k := range mf.Highlights {
		hlKeys = append(hlKeys, k)
	}
	sort.Strings(hlKeys)

	for _, name := range hlKeys {
		hl := mf.Highlights[name]
		hlValueNode := marshalVimHighlight(hl)
		hlNode.Content = append(hlNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: name, Tag: "!!str"},
			hlValueNode,
		)
	}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "highlights", Tag: "!!str"},
		hlNode,
	)

	// Build terminal_colors sequence node.
	tcNode := &yaml.Node{
		Kind: yaml.SequenceNode,
		Tag:  "!!seq",
	}

	for _, color := range mf.TerminalColors {
		tcNode.Content = append(tcNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: color, Tag: "!!str"},
		)
	}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "terminal_colors", Tag: "!!str"},
		tcNode,
	)

	// Build bufferline section if present.
	if mf.Bufferline != nil {
		blNode := marshalBufferlineTheme(mf.Bufferline)
		root.Content = append(root.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "bufferline", Tag: "!!str"},
			blNode,
		)
	}

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}

	return yaml.Marshal(doc)
}

// marshalBufferlineTheme builds a YAML mapping node for a BufferlineMappingTheme.
// Groups are serialized in a deterministic order matching the struct field order.
func marshalBufferlineTheme(bl *ports.BufferlineMappingTheme) *yaml.Node {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	addGroup := func(name string, colors ports.BufferlineMappingColors) {
		groupNode := marshalBufferlineColors(colors)
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: name, Tag: "!!str"},
			groupNode,
		)
	}

	// Add groups in deterministic order.
	addGroup("fill", bl.Fill)
	addGroup("background", bl.Background)
	addGroup("buffer_visible", bl.BufferVisible)
	addGroup("buffer_selected", bl.BufferSelected)
	addGroup("separator", bl.Separator)
	addGroup("separator_visible", bl.SeparatorVisible)
	addGroup("separator_selected", bl.SeparatorSelected)
	addGroup("indicator_selected", bl.IndicatorSelected)
	addGroup("modified", bl.Modified)
	addGroup("modified_visible", bl.ModifiedVisible)
	addGroup("modified_selected", bl.ModifiedSelected)
	addGroup("error", bl.Error)
	addGroup("warning", bl.Warning)
	addGroup("info", bl.Info)
	addGroup("hint", bl.Hint)

	return node
}

// marshalBufferlineColors builds a YAML mapping node for BufferlineMappingColors.
// Only non-zero fields are included (matching the omitempty YAML tags).
func marshalBufferlineColors(c ports.BufferlineMappingColors) *yaml.Node {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	if c.Fg != "" {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "fg", Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: c.Fg, Tag: "!!str"},
		)
	}
	if c.Bg != "" {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "bg", Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: c.Bg, Tag: "!!str"},
		)
	}
	if c.Bold {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "bold", Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
		)
	}
	if c.Italic {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "italic", Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
		)
	}

	return node
}

// marshalVimHighlight builds a YAML mapping node for a single VimMappingHighlight.
// Only non-zero fields are included (matching the omitempty YAML tags).
func marshalVimHighlight(hl ports.VimMappingHighlight) *yaml.Node {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	addStr := func(key, val string) {
		if val == "" {
			return
		}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: val, Tag: "!!str"},
		)
	}

	addBool := func(key string, val bool) {
		if !val {
			return
		}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
		)
	}

	addStr("fg", hl.Fg)
	addStr("bg", hl.Bg)
	addStr("sp", hl.Sp)
	addBool("bold", hl.Bold)
	addBool("italic", hl.Italic)
	addBool("underline", hl.Underline)
	addBool("undercurl", hl.Undercurl)
	addBool("strikethrough", hl.Strikethrough)
	addBool("reverse", hl.Reverse)
	addBool("nocombine", hl.Nocombine)
	addStr("link", hl.Link)

	return node
}

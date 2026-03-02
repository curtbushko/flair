// Package fileio provides adapters for reading and writing flair pipeline
// files (tokens.yaml, mapping files) via io.Reader and io.Writer.
package fileio

import (
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// ReadTokens reads YAML from r and returns a TokenSet. The YAML is
// expected to conform to the ports.TokensFile structure. Each
// TokenEntry is converted back to a domain.Token: an empty color
// string produces a NoneColor, while non-empty strings are parsed as hex.
//
// The caller may wrap r with a ValidatingReader to enforce schema version
// checking before the YAML is decoded.
func ReadTokens(r io.Reader) (*domain.TokenSet, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read tokens: %w", err)
	}

	var tf ports.TokensFile
	if err := yaml.Unmarshal(data, &tf); err != nil {
		return nil, fmt.Errorf("read tokens: %w", err)
	}

	ts := domain.NewTokenSet()
	for path, te := range tf.Tokens {
		tok := domain.Token{
			Bold:          te.Bold,
			Italic:        te.Italic,
			Underline:     te.Underline,
			Undercurl:     te.Undercurl,
			Strikethrough: te.Strikethrough,
		}

		if te.Color == "" {
			tok.Color = domain.NoneColor()
		} else {
			c, err := domain.ParseHex(te.Color)
			if err != nil {
				return nil, fmt.Errorf("read tokens: token %q: %w", path, err)
			}
			tok.Color = c
		}

		ts.Set(path, tok)
	}

	return ts, nil
}

// WriteTokens serializes a TokenSet as YAML to w using the
// ports.TokensFile structure. Domain tokens are converted to
// ports.TokenEntry (color as hex string, style flags preserved).
// Token paths are sorted for deterministic output.
//
// The caller is responsible for wrapping w with a VersionedWriter if
// schema version headers are desired.
func WriteTokens(w io.Writer, ts *domain.TokenSet) error {
	tf := ports.TokensFile{
		Tokens: make(map[string]ports.TokenEntry, ts.Len()),
	}

	paths := ts.Paths() // already sorted
	for _, path := range paths {
		tok, _ := ts.Get(path)
		te := ports.TokenEntry{
			Bold:          tok.Bold,
			Italic:        tok.Italic,
			Underline:     tok.Underline,
			Undercurl:     tok.Undercurl,
			Strikethrough: tok.Strikethrough,
		}

		if !tok.Color.IsNone {
			te.Color = tok.Color.Hex()
		}

		tf.Tokens[path] = te
	}

	// Encode using a yaml.Encoder with sorted map keys for determinism.
	out, err := marshalSortedTokens(tf)
	if err != nil {
		return fmt.Errorf("write tokens: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("write tokens: %w", err)
	}

	return nil
}

// marshalSortedTokens encodes a TokensFile to YAML with sorted
// token keys. The standard yaml.Marshal does not guarantee map key order,
// so we build the YAML node tree manually.
func marshalSortedTokens(tf ports.TokensFile) ([]byte, error) {
	// Build the tokens mapping node with sorted keys.
	tokensNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	paths := make([]string, 0, len(tf.Tokens))
	for p := range tf.Tokens {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, path := range paths {
		tok := tf.Tokens[path]

		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: path,
			Tag:   "!!str",
		}

		// Build the token value as a mapping node.
		valNode := &yaml.Node{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
		}

		// Always include color field.
		valNode.Content = append(valNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "color", Tag: "!!str"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: tok.Color, Tag: "!!str"},
		)

		// Add style flags only when set (matching omitempty behavior).
		if tok.Bold {
			valNode.Content = append(valNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "bold", Tag: "!!str"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
			)
		}
		if tok.Italic {
			valNode.Content = append(valNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "italic", Tag: "!!str"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
			)
		}
		if tok.Underline {
			valNode.Content = append(valNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "underline", Tag: "!!str"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
			)
		}
		if tok.Undercurl {
			valNode.Content = append(valNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "undercurl", Tag: "!!str"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
			)
		}
		if tok.Strikethrough {
			valNode.Content = append(valNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "strikethrough", Tag: "!!str"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
			)
		}

		tokensNode.Content = append(tokensNode.Content, keyNode, valNode)
	}

	// Build the root document node.
	root := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "tokens", Tag: "!!str"},
		tokensNode,
	)

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}

	out, err := yaml.Marshal(doc)
	if err != nil {
		return nil, err
	}

	return out, nil
}

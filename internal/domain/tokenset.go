package domain

import (
	"fmt"
	"sort"
)

// TokenSet is an aggregate that maps semantic path strings to Token values.
// It provides encapsulated access to an underlying map.
type TokenSet struct {
	tokens map[string]Token
}

// NewTokenSet returns an empty TokenSet ready for use.
func NewTokenSet() *TokenSet {
	return &TokenSet{
		tokens: make(map[string]Token),
	}
}

// Set stores a Token at the given path, overwriting any existing value.
func (ts *TokenSet) Set(path string, tok Token) {
	ts.tokens[path] = tok
}

// Get retrieves the Token at the given path.
// Returns the token and true if found, or a zero Token and false if not.
func (ts *TokenSet) Get(path string) (Token, bool) {
	tok, ok := ts.tokens[path]
	return tok, ok
}

// MustGet retrieves the Token at the given path, panicking if the path
// is not found. Use this only when the caller is certain the path exists.
func (ts *TokenSet) MustGet(path string) Token {
	tok, ok := ts.tokens[path]
	if !ok {
		panic(fmt.Sprintf("tokenset: path %q not found", path))
	}
	return tok
}

// Paths returns a sorted slice of all token path strings in the set.
func (ts *TokenSet) Paths() []string {
	paths := make([]string, 0, len(ts.tokens))
	for p := range ts.tokens {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// Len returns the number of tokens in the set.
func (ts *TokenSet) Len() int {
	return len(ts.tokens)
}

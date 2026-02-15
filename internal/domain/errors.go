// Package domain contains the core domain types, value objects, and error
// types for the flair theme pipeline. It has no external dependencies.
package domain

import "fmt"

// ParseError indicates a failure to parse input data.
type ParseError struct {
	Field   string
	Message string
	Cause   error
}

func (e *ParseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("parse error in %s: %s: %v", e.Field, e.Message, e.Cause)
	}
	return fmt.Sprintf("parse error in %s: %s", e.Field, e.Message)
}

func (e *ParseError) Unwrap() error { return e.Cause }

// ValidationError indicates a palette or file fails validation.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %d violations", len(e.Violations))
}

// GenerateError indicates a target generation failure.
type GenerateError struct {
	Target  string
	Message string
	Cause   error
}

func (e *GenerateError) Error() string {
	return fmt.Sprintf("generate %s: %s", e.Target, e.Message)
}

func (e *GenerateError) Unwrap() error { return e.Cause }

// SchemaVersionError indicates a file has an incompatible schema version.
type SchemaVersionError struct {
	File         string
	Found        int
	Expected     int
	NeedsUpgrade bool // true if Found > Expected (user needs newer flair)
}

func (e *SchemaVersionError) Error() string {
	if e.NeedsUpgrade {
		return fmt.Sprintf("%s: schema version %d is newer than supported %d — please upgrade flair",
			e.File, e.Found, e.Expected)
	}
	return fmt.Sprintf("%s: schema version %d is outdated (current: %d) — will regenerate",
		e.File, e.Found, e.Expected)
}

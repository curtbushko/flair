package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestParseError_Error(t *testing.T) {
	e := &domain.ParseError{
		Field:   "hex",
		Message: "invalid format",
	}
	got := e.Error()
	want := "parse error in hex: invalid format"
	if got != want {
		t.Errorf("ParseError.Error() = %q, want %q", got, want)
	}
}

func TestParseError_ErrorWithCause(t *testing.T) {
	cause := errors.New("underlying issue")
	e := &domain.ParseError{
		Field:   "hex",
		Message: "invalid",
		Cause:   cause,
	}
	got := e.Error()
	want := "parse error in hex: invalid: underlying issue"
	if got != want {
		t.Errorf("ParseError.Error() = %q, want %q", got, want)
	}
}

func TestParseError_Unwrap(t *testing.T) {
	cause := errors.New("wrapped cause")
	e := &domain.ParseError{
		Field:   "hex",
		Message: "invalid",
		Cause:   cause,
	}
	unwrapped := errors.Unwrap(e)
	if !errors.Is(unwrapped, cause) {
		t.Errorf("ParseError.Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestParseError_UnwrapNil(t *testing.T) {
	e := &domain.ParseError{
		Field:   "hex",
		Message: "invalid",
	}
	unwrapped := errors.Unwrap(e)
	if unwrapped != nil {
		t.Errorf("ParseError.Unwrap() = %v, want nil", unwrapped)
	}
}

func TestValidationError_Error(t *testing.T) {
	e := &domain.ValidationError{
		Violations: []string{"a", "b"},
	}
	got := e.Error()
	want := "validation failed: 2 violations"
	if got != want {
		t.Errorf("ValidationError.Error() = %q, want %q", got, want)
	}
}

func TestValidationError_ErrorSingle(t *testing.T) {
	e := &domain.ValidationError{
		Violations: []string{"only one"},
	}
	got := e.Error()
	want := "validation failed: 1 violations"
	if got != want {
		t.Errorf("ValidationError.Error() = %q, want %q", got, want)
	}
}

func TestGenerateError_Error(t *testing.T) {
	e := &domain.GenerateError{
		Target:  "vim",
		Message: "write failed",
	}
	got := e.Error()
	want := "generate vim: write failed"
	if got != want {
		t.Errorf("GenerateError.Error() = %q, want %q", got, want)
	}
}

func TestGenerateError_Unwrap(t *testing.T) {
	cause := errors.New("io error")
	e := &domain.GenerateError{
		Target:  "vim",
		Message: "write failed",
		Cause:   cause,
	}
	unwrapped := errors.Unwrap(e)
	if !errors.Is(unwrapped, cause) {
		t.Errorf("GenerateError.Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestGenerateError_UnwrapNil(t *testing.T) {
	e := &domain.GenerateError{
		Target:  "vim",
		Message: "write failed",
	}
	unwrapped := errors.Unwrap(e)
	if unwrapped != nil {
		t.Errorf("GenerateError.Unwrap() = %v, want nil", unwrapped)
	}
}

func TestSchemaVersionError_Outdated(t *testing.T) {
	e := &domain.SchemaVersionError{
		File:         "u.yaml",
		Found:        0,
		Expected:     1,
		NeedsUpgrade: false,
	}
	got := e.Error()
	if !strings.Contains(got, "outdated") {
		t.Errorf("SchemaVersionError.Error() = %q, want to contain 'outdated'", got)
	}
	if !strings.Contains(got, "will regenerate") {
		t.Errorf("SchemaVersionError.Error() = %q, want to contain 'will regenerate'", got)
	}
	if !strings.Contains(got, "u.yaml") {
		t.Errorf("SchemaVersionError.Error() = %q, want to contain 'u.yaml'", got)
	}
}

func TestSchemaVersionError_NeedsUpgrade(t *testing.T) {
	e := &domain.SchemaVersionError{
		File:         "u.yaml",
		Found:        99,
		Expected:     1,
		NeedsUpgrade: true,
	}
	got := e.Error()
	if !strings.Contains(got, "newer than supported") {
		t.Errorf("SchemaVersionError.Error() = %q, want to contain 'newer than supported'", got)
	}
	if !strings.Contains(got, "please upgrade flair") {
		t.Errorf("SchemaVersionError.Error() = %q, want to contain 'please upgrade flair'", got)
	}
	if !strings.Contains(got, "u.yaml") {
		t.Errorf("SchemaVersionError.Error() = %q, want to contain 'u.yaml'", got)
	}
}

func TestErrorTypes_ImplementError(t *testing.T) {
	// Compile-time check: all error types implement the error interface.
	var _ error = &domain.ParseError{}
	var _ error = &domain.ValidationError{}
	var _ error = &domain.GenerateError{}
	var _ error = &domain.SchemaVersionError{}
}

func TestParseError_WorksWithErrorsIs(t *testing.T) {
	cause := errors.New("root cause")
	e := &domain.ParseError{
		Field:   "hex",
		Message: "invalid",
		Cause:   cause,
	}
	if !errors.Is(e, cause) {
		t.Error("errors.Is should find the wrapped cause in ParseError")
	}
}

func TestGenerateError_WorksWithErrorsAs(t *testing.T) {
	inner := &domain.ParseError{Field: "color", Message: "bad"}
	e := &domain.GenerateError{
		Target:  "vim",
		Message: "generation failed",
		Cause:   inner,
	}
	var target *domain.ParseError
	if !errors.As(e, &target) {
		t.Error("errors.As should find the wrapped ParseError in GenerateError")
	}
	if target.Field != "color" {
		t.Errorf("errors.As target.Field = %q, want 'color'", target.Field)
	}
}

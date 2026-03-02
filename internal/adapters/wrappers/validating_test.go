package wrappers_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/wrappers"
	"github.com/curtbushko/flair/internal/domain"
)

func TestValidatingReader_ValidVersion(t *testing.T) {
	content := fmt.Sprintf("schema_version: %d\nkind: tokens\ntheme_name: tokyonight\ntokens:\n  fg: '#c0caf5'\n",
		domain.CurrentVersion(domain.FileKindTokens))

	vr := wrappers.NewValidatingReader(bytes.NewReader([]byte(content)), domain.FileKindTokens)

	got, err := io.ReadAll(vr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != content {
		t.Errorf("content mismatch.\ngot:\n%s\nwant:\n%s", string(got), content)
	}
}

func TestValidatingReader_OutdatedVersion(t *testing.T) {
	content := "schema_version: 0\nkind: tokens\ntheme_name: old\ntokens:\n  fg: '#000000'\n"

	vr := wrappers.NewValidatingReader(bytes.NewReader([]byte(content)), domain.FileKindTokens)

	_, err := io.ReadAll(vr)
	if err == nil {
		t.Fatal("expected error for outdated version, got nil")
	}

	var schemaErr *domain.SchemaVersionError
	if !errors.As(err, &schemaErr) {
		t.Fatalf("expected *domain.SchemaVersionError, got %T: %v", err, err)
	}

	if schemaErr.NeedsUpgrade {
		t.Error("NeedsUpgrade should be false for outdated version")
	}
	if schemaErr.Found != 0 {
		t.Errorf("Found = %d, want 0", schemaErr.Found)
	}
	if schemaErr.Expected != domain.CurrentVersion(domain.FileKindTokens) {
		t.Errorf("Expected = %d, want %d", schemaErr.Expected, domain.CurrentVersion(domain.FileKindTokens))
	}
}

func TestValidatingReader_FutureVersion(t *testing.T) {
	content := "schema_version: 99\nkind: tokens\ntheme_name: future\ntokens:\n  fg: '#ffffff'\n"

	vr := wrappers.NewValidatingReader(bytes.NewReader([]byte(content)), domain.FileKindTokens)

	_, err := io.ReadAll(vr)
	if err == nil {
		t.Fatal("expected error for future version, got nil")
	}

	var schemaErr *domain.SchemaVersionError
	if !errors.As(err, &schemaErr) {
		t.Fatalf("expected *domain.SchemaVersionError, got %T: %v", err, err)
	}

	if !schemaErr.NeedsUpgrade {
		t.Error("NeedsUpgrade should be true for future version")
	}
	if schemaErr.Found != 99 {
		t.Errorf("Found = %d, want 99", schemaErr.Found)
	}
	if schemaErr.Expected != domain.CurrentVersion(domain.FileKindTokens) {
		t.Errorf("Expected = %d, want %d", schemaErr.Expected, domain.CurrentVersion(domain.FileKindTokens))
	}
}

func TestValidatingReader_MissingVersion(t *testing.T) {
	content := "kind: tokens\ntheme_name: no-version\ntokens:\n  fg: '#aabbcc'\n"

	vr := wrappers.NewValidatingReader(bytes.NewReader([]byte(content)), domain.FileKindTokens)

	_, err := io.ReadAll(vr)
	if err == nil {
		t.Fatal("expected error for missing schema_version, got nil")
	}
}

func TestValidatingReader_Composable(t *testing.T) {
	content := fmt.Sprintf("schema_version: %d\nkind: tokens\ntheme_name: compose\ntokens:\n  fg: '#112233'\n",
		domain.CurrentVersion(domain.FileKindTokens))

	// Wrap a bytes.Reader in ValidatingReader
	inner := bytes.NewReader([]byte(content))
	vr := wrappers.NewValidatingReader(inner, domain.FileKindTokens)

	got, err := io.ReadAll(vr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != content {
		t.Errorf("composable read mismatch.\ngot:\n%s\nwant:\n%s", string(got), content)
	}
}

func TestValidatingReader_MultipleReads(t *testing.T) {
	content := fmt.Sprintf("schema_version: %d\nkind: tokens\ntheme_name: multi\ntokens:\n  fg: '#c0caf5'\n  bg: '#1a1b26'\n",
		domain.CurrentVersion(domain.FileKindTokens))

	vr := wrappers.NewValidatingReader(bytes.NewReader([]byte(content)), domain.FileKindTokens)

	// Read in small chunks
	var result bytes.Buffer
	buf := make([]byte, 8)
	for {
		n, err := vr.Read(buf)
		if n > 0 {
			result.Write(buf[:n])
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error during chunked read: %v", err)
		}
	}

	if result.String() != content {
		t.Errorf("chunked read mismatch.\ngot:\n%s\nwant:\n%s", result.String(), content)
	}
}

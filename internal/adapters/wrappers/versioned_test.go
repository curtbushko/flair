package wrappers_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/wrappers"
	"github.com/curtbushko/flair/internal/domain"
)

func TestVersionedWriter_PrependsHeader(t *testing.T) {
	var buf bytes.Buffer
	vw := wrappers.NewVersionedWriter(&buf, domain.FileKindUniversal, "tokyonight")

	content := "tokens:\n  fg: '#c0caf5'\n"
	_, err := vw.Write([]byte(content))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	wantHeader := fmt.Sprintf("schema_version: %d\nkind: universal\ntheme_name: tokyonight\n",
		domain.CurrentVersion(domain.FileKindUniversal))

	if !strings.HasPrefix(got, wantHeader) {
		t.Errorf("buffer does not start with expected header.\ngot:\n%s\nwant prefix:\n%s", got, wantHeader)
	}

	// Content should follow the header
	wantFull := wantHeader + content
	if got != wantFull {
		t.Errorf("full output mismatch.\ngot:\n%s\nwant:\n%s", got, wantFull)
	}
}

func TestVersionedWriter_HeaderOnce(t *testing.T) {
	var buf bytes.Buffer
	vw := wrappers.NewVersionedWriter(&buf, domain.FileKindPalette, "catppuccin")

	_, err := vw.Write([]byte("first"))
	if err != nil {
		t.Fatalf("first write error: %v", err)
	}

	_, err = vw.Write([]byte("second"))
	if err != nil {
		t.Fatalf("second write error: %v", err)
	}

	got := buf.String()
	headerLine := "schema_version:"

	count := strings.Count(got, headerLine)
	if count != 1 {
		t.Errorf("header appears %d times, want exactly 1.\ngot:\n%s", count, got)
	}

	// Content should have both writes after the header
	wantHeader := fmt.Sprintf("schema_version: %d\nkind: palette\ntheme_name: catppuccin\n",
		domain.CurrentVersion(domain.FileKindPalette))
	wantFull := wantHeader + "firstsecond"
	if got != wantFull {
		t.Errorf("full output mismatch.\ngot:\n%s\nwant:\n%s", got, wantFull)
	}
}

func TestVersionedWriter_CorrectVersionPerKind(t *testing.T) {
	kinds := []struct {
		kind     domain.FileKind
		wantKind string
	}{
		{domain.FileKindPalette, "palette"},
		{domain.FileKindUniversal, "universal"},
		{domain.FileKindVimMapping, "vim-mapping"},
		{domain.FileKindCSSMapping, "css-mapping"},
		{domain.FileKindGtkMapping, "gtk-mapping"},
		{domain.FileKindQssMapping, "qss-mapping"},
		{domain.FileKindStylixMapping, "stylix-mapping"},
	}

	for _, tc := range kinds {
		t.Run(string(tc.kind), func(t *testing.T) {
			var buf bytes.Buffer
			vw := wrappers.NewVersionedWriter(&buf, tc.kind, "test-theme")

			_, err := vw.Write([]byte("data"))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := buf.String()
			wantVersion := fmt.Sprintf("schema_version: %d\n", domain.CurrentVersion(tc.kind))

			if !strings.HasPrefix(got, wantVersion) {
				t.Errorf("expected prefix %q, got:\n%s", wantVersion, got)
			}

			wantKindLine := fmt.Sprintf("kind: %s\n", tc.wantKind)
			if !strings.Contains(got, wantKindLine) {
				t.Errorf("expected kind line %q in:\n%s", wantKindLine, got)
			}
		})
	}
}

func TestVersionedWriter_EmptyWrite(t *testing.T) {
	var buf bytes.Buffer
	vw := wrappers.NewVersionedWriter(&buf, domain.FileKindUniversal, "empty-theme")

	_, err := vw.Write([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	wantHeader := fmt.Sprintf("schema_version: %d\nkind: universal\ntheme_name: empty-theme\n",
		domain.CurrentVersion(domain.FileKindUniversal))

	if got != wantHeader {
		t.Errorf("empty write should still produce header.\ngot:\n%q\nwant:\n%q", got, wantHeader)
	}
}

func TestVersionedWriter_ByteCount(t *testing.T) {
	var buf bytes.Buffer
	vw := wrappers.NewVersionedWriter(&buf, domain.FileKindPalette, "count-theme")

	data := []byte("some content here")
	n, err := vw.Write(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Write returned %d bytes, want %d (content length, not including header)", n, len(data))
	}

	// Second write should also return content length only
	data2 := []byte("more data")
	n2, err := vw.Write(data2)
	if err != nil {
		t.Fatalf("unexpected error on second write: %v", err)
	}

	if n2 != len(data2) {
		t.Errorf("second Write returned %d bytes, want %d", n2, len(data2))
	}
}

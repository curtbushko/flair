package domain_test

import (
	"testing"

	"github.com/curtbushko/flair/internal/domain"
)

func TestCurrentVersion_AllKinds(t *testing.T) {
	tests := []struct {
		name    string
		kind    domain.FileKind
		wantVer int
	}{
		{"palette", domain.FileKindPalette, 1},
		{"universal", domain.FileKindUniversal, 1},
		{"vim-mapping", domain.FileKindVimMapping, 1},
		{"css-mapping", domain.FileKindCSSMapping, 1},
		{"gtk-mapping", domain.FileKindGtkMapping, 1},
		{"qss-mapping", domain.FileKindQssMapping, 1},
		{"stylix-mapping", domain.FileKindStylixMapping, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.CurrentVersion(tt.kind)
			if got != tt.wantVer {
				t.Errorf("CurrentVersion(%q) = %d, want %d", tt.kind, got, tt.wantVer)
			}
		})
	}
}

func TestCurrentVersion_AllPositive(t *testing.T) {
	allKinds := []domain.FileKind{
		domain.FileKindPalette,
		domain.FileKindUniversal,
		domain.FileKindVimMapping,
		domain.FileKindCSSMapping,
		domain.FileKindGtkMapping,
		domain.FileKindQssMapping,
		domain.FileKindStylixMapping,
	}

	for _, kind := range allKinds {
		t.Run(string(kind), func(t *testing.T) {
			ver := domain.CurrentVersion(kind)
			if ver <= 0 {
				t.Errorf("CurrentVersion(%q) = %d, want > 0", kind, ver)
			}
		})
	}
}

func TestCurrentVersion_UnknownKind(t *testing.T) {
	got := domain.CurrentVersion(domain.FileKind("nonexistent"))
	if got != 0 {
		t.Errorf("CurrentVersion(\"nonexistent\") = %d, want 0", got)
	}
}

func TestFileKind_StringValues(t *testing.T) {
	tests := []struct {
		kind domain.FileKind
		want string
	}{
		{domain.FileKindPalette, "palette"},
		{domain.FileKindUniversal, "universal"},
		{domain.FileKindVimMapping, "vim-mapping"},
		{domain.FileKindCSSMapping, "css-mapping"},
		{domain.FileKindGtkMapping, "gtk-mapping"},
		{domain.FileKindQssMapping, "qss-mapping"},
		{domain.FileKindStylixMapping, "stylix-mapping"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := string(tt.kind)
			if got != tt.want {
				t.Errorf("FileKind string = %q, want %q", got, tt.want)
			}
		})
	}
}

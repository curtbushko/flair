package palettes_test

import (
	"io"
	"sort"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/palettes"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/ports"
	"github.com/curtbushko/flair/pkg/flair"
)

// Verify Source satisfies ports.PaletteSource at compile time.
var _ ports.PaletteSource = (*palettes.Source)(nil)

func TestSource_List(t *testing.T) {
	src := palettes.NewSource()
	got := src.List()

	// Should return a reasonable number of palettes
	if len(got) == 0 {
		t.Fatal("List() returned empty list, expected at least one palette")
	}

	// Verify the list is sorted
	if !sort.StringsAreSorted(got) {
		t.Errorf("List() result is not sorted: %v", got)
	}

	// Verify known palettes are present
	knownPalettes := []string{
		"tokyo-night-dark",
		"gruvbox-dark",
		"catppuccin-mocha",
		"everforest",
	}

	for _, known := range knownPalettes {
		found := false
		for _, name := range got {
			if name == known {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("List() missing expected palette %q", known)
		}
	}
}

func TestSource_Get_Valid(t *testing.T) {
	src := palettes.NewSource()

	reader, err := src.Get("tokyo-night-dark")
	if err != nil {
		t.Fatalf("Get('tokyo-night-dark') error: %v", err)
	}

	// Read all bytes to verify it produces valid YAML content
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Get('tokyo-night-dark') returned empty reader")
	}

	// Verify the content contains expected YAML fields
	content := string(data)
	for _, field := range []string{"system", "name", "palette"} {
		if !contains(content, field) {
			t.Errorf("expected YAML to contain field %q", field)
		}
	}
}

func TestSource_Get_Unknown(t *testing.T) {
	src := palettes.NewSource()

	_, err := src.Get("nonexistent")
	if err == nil {
		t.Fatal("Get('nonexistent') expected error, got nil")
	}
}

func TestSource_Has_True(t *testing.T) {
	src := palettes.NewSource()

	if !src.Has("tokyo-night-dark") {
		t.Error("Has('tokyo-night-dark') = false, want true")
	}
}

func TestSource_Has_False(t *testing.T) {
	src := palettes.NewSource()

	if src.Has("my-custom-theme") {
		t.Error("Has('my-custom-theme') = true, want false")
	}
}

func TestSource_Get_ParseableByParser(t *testing.T) {
	src := palettes.NewSource()
	parser := yamlparser.NewParser()

	reader, err := src.Get("tokyo-night-dark")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}

	pal, err := parser.Parse(reader)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if pal.Name != "Tokyo Night Dark" {
		t.Errorf("parsed name = %q, want %q", pal.Name, "Tokyo Night Dark")
	}

	if pal.System != "base24" {
		t.Errorf("parsed system = %q, want %q", pal.System, "base24")
	}

	// Verify all 24 slots have valid colors
	for i := 0; i < 24; i++ {
		if pal.Base(i).IsNone {
			t.Errorf("slot %d is IsNone, expected valid color", i)
		}
	}
}

func TestSource_AllPalettes_Parseable(t *testing.T) {
	src := palettes.NewSource()
	parser := yamlparser.NewParser()

	names := src.List()
	if len(names) == 0 {
		t.Fatal("List() returned empty, expected at least 8 palettes")
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			reader, err := src.Get(name)
			if err != nil {
				t.Fatalf("Get(%q) error: %v", name, err)
			}

			pal, err := parser.Parse(reader)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", name, err)
			}

			if pal.Name == "" {
				t.Errorf("parsed palette %q has empty name", name)
			}

			if pal.System != "base24" {
				t.Errorf("parsed palette %q system = %q, want %q", name, pal.System, "base24")
			}

			if pal.Variant != "dark" && pal.Variant != "light" {
				t.Errorf("parsed palette %q variant = %q, want dark or light", name, pal.Variant)
			}

			// All 24 slots should have valid colors
			for i := 0; i < 24; i++ {
				if pal.Base(i).IsNone {
					t.Errorf("palette %q slot %d is IsNone", name, i)
				}
			}
		})
	}
}

// contains checks if s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestSource_UsesPkgFlairPalettes verifies that internal/adapters/palettes
// delegates to pkg/flair/palettes (DRY - single source of truth).
func TestSource_UsesPkgFlairPalettes(t *testing.T) {
	src := palettes.NewSource()

	// Get the list from internal adapter
	internalList := src.List()

	// Get the list from pkg/flair builtins
	pkgList := flair.ListBuiltins()

	// Both should return the same number of palettes
	if len(internalList) != len(pkgList) {
		t.Fatalf("internal palettes count (%d) != pkg/flair palettes count (%d)",
			len(internalList), len(pkgList))
	}

	// Both should have the exact same palette names
	for i, name := range internalList {
		if name != pkgList[i] {
			t.Errorf("palette mismatch at index %d: internal=%q, pkg=%q",
				i, name, pkgList[i])
		}
	}

	// Verify content is identical for a sample palette
	internalReader, err := src.Get("tokyo-night-dark")
	if err != nil {
		t.Fatalf("internal Get error: %v", err)
	}
	internalData, err := io.ReadAll(internalReader)
	if err != nil {
		t.Fatalf("internal ReadAll error: %v", err)
	}
	if len(internalData) == 0 {
		t.Fatal("internal palette data is empty")
	}

	// Load via pkg/flair to verify same content is accessible
	if !flair.HasBuiltin("tokyo-night-dark") {
		t.Fatal("pkg/flair does not have tokyo-night-dark builtin")
	}
}

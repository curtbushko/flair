package palettes_test

import (
	"io"
	"sort"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/palettes"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/ports"
)

// Verify Source satisfies ports.PaletteSource at compile time.
var _ ports.PaletteSource = (*palettes.Source)(nil)

func TestSource_List(t *testing.T) {
	src := palettes.NewSource()
	got := src.List()

	want := []string{
		"0x96f",
		"3024-day",
		"3024-night",
		"adventure-time",
		"alien-blood",
		"andromeda",
		"argonaut",
		"arthur",
		"atelier-sulphurpool",
		"ayu-dark",
		"ayu-light",
		"ayu-mirage",
		"banana-blueberry",
		"batman",
		"birds-of-paradise",
		"blazer",
		"blue-berry-pie",
		"blue-matrix",
		"bluloco-dark",
		"bluloco-light",
		"borland",
		"breeze",
		"broadcast",
		"brogrammer",
		"builtin-dark",
		"builtin-light",
		"builtin-pastel-dark",
		"builtin-solarized-dark",
		"builtin-solarized-light",
		"builtin-tango-dark",
		"builtin-tango-light",
		"catppuccin-frappe",
		"catppuccin-latte",
		"catppuccin-macchiato",
		"catppuccin-mocha",
		"chalk",
		"chalkboard",
		"challenger-deep",
		"ciapre",
		"clrs",
		"cobalt-neon",
		"cobalt2",
		"crayon-pony-fish",
		"cyberdyne",
		"dark-plus",
		"deep",
		"deep-oceanic-next",
		"desert",
		"dimmed-monokai",
		"dracula",
		"earthsong",
		"eldritch",
		"elemental",
		"elementary",
		"embarcadero",
		"encom",
		"espresso",
		"espresso-libre",
		"everforest",
		"fideloper",
		"firefox-dev",
		"fish-tank",
		"flat",
		"flatland",
		"flexoki-dark",
		"flexoki-light",
		"floraverse",
		"forest-blue",
		"framer",
		"front-end-delight",
		"fun-forrest",
		"galaxy",
		"github",
		"github-dark",
		"grape",
		"gruvbox-dark",
		"gruvbox-light",
		"gruvbox-material",
		"hacktober",
		"hardcore",
		"hipster-green",
		"hivacruz",
		"homebrew",
		"hopscotch",
		"hurtado",
		"hybrid",
		"ic-green-ppl",
		"ic-orange-ppl",
		"idea",
		"idle-toes",
		"jackie-brown",
		"japanesque",
		"jellybeans",
		"jet-brains-darcula",
		"kanagawa-dragon",
		"kibble",
		"lab-fox",
		"laser",
		"later-this-evening",
		"lavandula",
		"lovelace",
		"man-page",
		"material",
		"material-dark",
		"mathias",
		"medallion",
		"mission-brogue",
		"misterioso",
		"molokai",
		"mona-lisa",
		"monokai-vivid",
		"mountain",
		"neofusion",
		"night-lion-v1",
		"night-lion-v2",
		"night-owlish-light",
		"nocturnal-winter",
		"obsidian",
		"ocean",
		"oceanic-material",
		"ollie",
		"one-black",
		"one-dark",
		"one-half-light",
		"one-light",
		"operator-mono-dark",
		"orng",
		"pandora",
		"papercolor-dark",
		"papercolor-light",
		"paul-millr",
		"pencil-dark",
		"pencil-light",
		"piatto-light",
		"pnevma",
		"pro",
		"pro-light",
		"purple-rain",
		"purplepeter",
		"rebecca",
		"rebel-scum",
		"red-alert",
		"red-planet",
		"red-sands",
		"rippedcasts",
		"royal",
		"scarlet-protocol",
		"sea-shells",
		"seafoam-pastel",
		"shades-of-purple",
		"shaman",
		"slate",
		"sleepy-hollow",
		"smyck",
		"solarized-dark-patched",
		"space-gray-eighties",
		"space-gray-eighties-dull",
		"spacedust",
		"sparky",
		"spiderman",
		"square",
		"sundried",
		"tango-adapted",
		"tango-half-adapted",
		"terminal-basic",
		"thayer-bright",
		"the-hulk",
		"tokyo-night-dark",
		"tokyo-night-light",
		"tokyo-night-moon",
		"tokyo-night-neon",
		"tokyo-night-storm",
		"tomorrow-night",
		"toy-chest",
		"treehouse",
		"twilight",
		"ubuntu",
		"ultra-violet",
		"under-the-sea",
		"unikitty",
		"vibrant-ink",
		"violet-dark",
		"violet-light",
		"warm-neon",
		"wez",
		"wild-cherry",
		"wombat",
		"wryan",
		"zenburn",
	}

	if len(got) != len(want) {
		t.Fatalf("List() returned %d items, want %d: %v", len(got), len(want), got)
	}

	// Verify the list is sorted
	if !sort.StringsAreSorted(got) {
		t.Errorf("List() result is not sorted: %v", got)
	}

	for i, name := range want {
		if got[i] != name {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], name)
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

package fileio_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/adapters/fileio"
	"github.com/curtbushko/flair/internal/domain"
)

func TestWritePalette_ContainsAllFields(t *testing.T) {
	// Create a palette with known values.
	var colors [24]domain.Color
	for i := range colors {
		// Use a simple pattern: #00XXYY where XX=i*10, YY=i*11
		colors[i] = domain.Color{R: uint8(i * 10), G: uint8(i * 11), B: uint8(i * 12)}
	}
	pal := &domain.Palette{
		System:  "base24",
		Name:    "Test Theme",
		Author:  "Test Author",
		Variant: "dark",
		Colors:  colors,
	}

	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// Verify metadata fields are present.
	if !strings.Contains(output, `system: "base24"`) {
		t.Error("output should contain system field")
	}
	if !strings.Contains(output, `name: "Test Theme"`) {
		t.Error("output should contain name field")
	}
	if !strings.Contains(output, `author: "Test Author"`) {
		t.Error("output should contain author field")
	}
	if !strings.Contains(output, `variant: "dark"`) {
		t.Error("output should contain variant field")
	}
	if !strings.Contains(output, "palette:") {
		t.Error("output should contain palette section")
	}

	// Verify all 24 base slots are present.
	expectedSlots := []string{
		"base00:", "base01:", "base02:", "base03:",
		"base04:", "base05:", "base06:", "base07:",
		"base08:", "base09:", "base0A:", "base0B:",
		"base0C:", "base0D:", "base0E:", "base0F:",
		"base10:", "base11:", "base12:", "base13:",
		"base14:", "base15:", "base16:", "base17:",
	}
	for _, slot := range expectedSlots {
		if !strings.Contains(output, slot) {
			t.Errorf("output should contain %s", slot)
		}
	}
}

func TestWritePalette_HexFormat(t *testing.T) {
	// Create a palette with a specific color to verify hex formatting.
	var colors [24]domain.Color
	colors[0] = domain.Color{R: 0x1a, G: 0x1b, B: 0x26}
	for i := 1; i < 24; i++ {
		colors[i] = domain.Color{R: 0xff, G: 0xff, B: 0xff}
	}
	pal := &domain.Palette{
		System:  "base24",
		Name:    "Test",
		Author:  "Test",
		Variant: "dark",
		Colors:  colors,
	}

	var buf bytes.Buffer
	err := fileio.WritePalette(&buf, pal)
	if err != nil {
		t.Fatalf("WritePalette() error = %v", err)
	}

	output := buf.String()

	// base00 should have the specific color.
	if !strings.Contains(output, `base00: "#1a1b26"`) {
		t.Errorf("output should contain base00 with correct hex value, got: %s", output)
	}
}

// Example: programmatic demonstrates programmatic theme generation.
//
// This example shows how to:
// - Parse a palette YAML file with flair.ParsePalette()
// - Tokenize the palette to generate semantic tokens with flair.Tokenize()
// - Access the generated theme's colors
//
// This is useful for building theme generation tools or processing
// custom palettes not included in the built-in set.
//
// Run with: go run ./examples/programmatic
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/curtbushko/flair/pkg/flair"
)

// Sample palette YAML embedded for demonstration.
// In production, you would read this from a file.
const samplePalette = `
name: "Custom Theme"
author: "Example Author"
variant: "dark"
palette:
  base00: "#1e1e2e"
  base01: "#181825"
  base02: "#313244"
  base03: "#45475a"
  base04: "#585b70"
  base05: "#cdd6f4"
  base06: "#f5e0dc"
  base07: "#b4befe"
  base08: "#f38ba8"
  base09: "#fab387"
  base0A: "#f9e2af"
  base0B: "#a6e3a1"
  base0C: "#94e2d5"
  base0D: "#89b4fa"
  base0E: "#cba6f7"
  base0F: "#f2cdcd"
  base10: "#11111b"
  base11: "#0c0c14"
  base12: "#f38ba8"
  base13: "#f9e2af"
  base14: "#a6e3a1"
  base15: "#89dceb"
  base16: "#74c7ec"
  base17: "#b4befe"
`

func main() {
	fmt.Println("Programmatic Theme Generation")
	fmt.Println("==============================")
	fmt.Println()

	// Parse the palette YAML.
	// ParsePalette reads a base24 palette from any io.Reader.
	palette, err := flair.ParsePalette(strings.NewReader(samplePalette))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed palette: %s\n", palette.Name())
	fmt.Printf("Author: %s\n", palette.Author())
	fmt.Printf("Variant: %s\n\n", palette.Variant())

	// Access base24 colors directly from the palette.
	fmt.Println("Base24 Colors:")
	fmt.Printf("  base00 (background): %s\n", palette.Base(0x00).Hex())
	fmt.Printf("  base05 (foreground): %s\n", palette.Base(0x05).Hex())
	fmt.Printf("  base08 (red):        %s\n", palette.Base(0x08).Hex())
	fmt.Printf("  base0B (green):      %s\n", palette.Base(0x0B).Hex())
	fmt.Printf("  base0D (blue):       %s\n", palette.Base(0x0D).Hex())
	fmt.Println()

	// Tokenize the palette to generate semantic tokens.
	// This transforms the 24-color palette into ~88 semantic tokens.
	theme := flair.Tokenize(palette)

	fmt.Printf("Generated theme: %s\n\n", theme.Name())

	// Access semantic tokens via typed accessors.
	fmt.Println("Surface Tokens:")
	surface := theme.Surface()
	fmt.Printf("  background:   %s\n", surface.Background.Hex())
	fmt.Printf("  raised:       %s\n", surface.Raised.Hex())
	fmt.Printf("  selection:    %s\n", surface.Selection.Hex())
	fmt.Println()

	fmt.Println("Status Tokens:")
	status := theme.Status()
	fmt.Printf("  error:   %s\n", status.Error.Hex())
	fmt.Printf("  warning: %s\n", status.Warning.Hex())
	fmt.Printf("  success: %s\n", status.Success.Hex())
	fmt.Printf("  info:    %s\n", status.Info.Hex())
	fmt.Println()

	// Access all colors as a map.
	colors := theme.Colors()
	fmt.Printf("Total semantic tokens generated: %d\n", len(colors))
}

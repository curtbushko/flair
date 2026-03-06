// Example: custom-overrides demonstrates creating a custom theme with token overrides.
//
// This example shows how to:
// - Create a custom palette YAML file with token overrides
// - Load and parse a palette with overrides using the internal adapter
// - Apply overrides during tokenization
// - Access the customized theme colors
//
// Token overrides allow you to customize specific semantic tokens without
// modifying the underlying base24 palette or derivation rules. This is useful
// for fine-tuning themes to your preferences.
//
// See the accompanying my-theme/palette.yaml for the override format.
//
// Run with: go run ./examples/custom-overrides
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/curtbushko/flair/internal/adapters/tokenizer"
	"github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/domain"
)

func main() {
	fmt.Println("Custom Theme with Token Overrides")
	fmt.Println("==================================")
	fmt.Println()

	// Locate the palette file relative to this example.
	// In a real application, you might use a config path or embed the palette.
	exampleDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Try to find the palette file
	palettePath := filepath.Join(exampleDir, "examples", "custom-overrides", "my-theme", "palette.yaml")
	if _, statErr := os.Stat(palettePath); os.IsNotExist(statErr) {
		// Try running from the example directory itself
		palettePath = filepath.Join(exampleDir, "my-theme", "palette.yaml")
	}

	// Open the palette file.
	file, err := os.Open(palettePath)
	if err != nil {
		log.Fatalf("Failed to open palette: %v\n\nMake sure to run from the project root:\n  go run ./examples/custom-overrides", err)
	}
	defer func() { _ = file.Close() }()

	// Parse the palette using the YAML parser.
	// This parser supports the overrides section.
	parser := yaml.NewParser()
	palette, err := parser.Parse(file)
	if err != nil {
		log.Fatalf("Failed to parse palette: %v", err)
	}

	fmt.Printf("Loaded palette: %s\n", palette.Name)
	fmt.Printf("Author: %s\n", palette.Author)
	fmt.Printf("Variant: %s\n\n", palette.Variant)

	// Show the overrides that were parsed.
	printOverrides(palette)

	// Tokenize the palette with overrides applied.
	// The tokenizer automatically applies any overrides defined in the palette.
	tok := tokenizer.New()
	tokenSet := tok.Tokenize(palette)

	fmt.Println("Tokenization complete!")
	fmt.Printf("Total tokens generated: %d\n\n", len(tokenSet.Paths()))

	// Demonstrate that overrides were applied by comparing to expected values.
	fmt.Println("Sample overridden tokens:")

	// Show syntax.keyword - should be magenta (#ff00ff) and bold
	if token, ok := tokenSet.Get("syntax.keyword"); ok {
		fmt.Printf("  syntax.keyword:  %s", token.Color.Hex())
		if token.Bold {
			fmt.Print(" (bold)")
		}
		fmt.Println()
	}

	// Show syntax.string - should be the custom green (#98c379)
	if token, ok := tokenSet.Get("syntax.string"); ok {
		fmt.Printf("  syntax.string:   %s\n", token.Color.Hex())
	}

	// Show syntax.comment - should be #6a737d with italic
	if token, ok := tokenSet.Get("syntax.comment"); ok {
		fmt.Printf("  syntax.comment:  %s", token.Color.Hex())
		if token.Italic {
			fmt.Print(" (italic)")
		}
		fmt.Println()
	}

	// Show status.error - should be #e06c75 with bold
	if token, ok := tokenSet.Get("status.error"); ok {
		fmt.Printf("  status.error:    %s", token.Color.Hex())
		if token.Bold {
			fmt.Print(" (bold)")
		}
		fmt.Println()
	}

	// Show statusline overrides
	if token, ok := tokenSet.Get("statusline.a.bg"); ok {
		fmt.Printf("  statusline.a.bg: %s\n", token.Color.Hex())
	}

	fmt.Println()

	// Compare an overridden token with a non-overridden one.
	fmt.Println("Non-overridden tokens (derived from palette):")

	if token, ok := tokenSet.Get("syntax.type"); ok {
		fmt.Printf("  syntax.type:     %s (derived from base0A)\n", token.Color.Hex())
	}

	if token, ok := tokenSet.Get("syntax.constant"); ok {
		fmt.Printf("  syntax.constant: %s (derived from base09)\n", token.Color.Hex())
	}

	fmt.Println()
	fmt.Println("Override precedence:")
	fmt.Println("  1. Base24 palette colors define the foundation")
	fmt.Println("  2. Tokenizer derives ~88 semantic tokens from the palette")
	fmt.Println("  3. Overrides are applied last, replacing specific tokens")
	fmt.Println()
	fmt.Println("This allows you to customize individual tokens while")
	fmt.Println("keeping the rest of the theme's color harmony intact.")
}

func printOverrides(palette *domain.Palette) {
	if len(palette.Overrides) == 0 {
		return
	}

	fmt.Printf("Token overrides defined: %d\n", len(palette.Overrides))
	fmt.Println()
	fmt.Println("Override Summary:")

	for path, override := range palette.Overrides {
		colorStr := "(no color change)"
		if override.HasColor() {
			colorStr = override.Color.Hex()
		}

		styleStr := formatStyles(override)
		fmt.Printf("  %-35s %s%s\n", path+":", colorStr, styleStr)
	}
	fmt.Println()
}

func formatStyles(override domain.TokenOverride) string {
	styles := []string{}
	if override.Bold {
		styles = append(styles, "bold")
	}
	if override.Italic {
		styles = append(styles, "italic")
	}
	if override.Underline {
		styles = append(styles, "underline")
	}
	if override.Strikethrough {
		styles = append(styles, "strikethrough")
	}

	if len(styles) == 0 {
		return ""
	}
	return fmt.Sprintf(" [%v]", styles)
}

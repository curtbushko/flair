package fileio

import (
	"fmt"
	"io"
	"sort"

	"github.com/curtbushko/flair/internal/domain"
)

// slotNames lists the base24 palette slot names in order.
var slotNames = []string{
	"base00", "base01", "base02", "base03",
	"base04", "base05", "base06", "base07",
	"base08", "base09", "base0A", "base0B",
	"base0C", "base0D", "base0E", "base0F",
	"base10", "base11", "base12", "base13",
	"base14", "base15", "base16", "base17",
}

// WritePalette serializes a domain.Palette to YAML format compatible with
// the tinted-theming palette format. This allows the palette to be re-parsed
// by the regenerate use case.
//
// The caller is responsible for wrapping w with a VersionedWriter if
// schema version headers are desired.
func WritePalette(w io.Writer, pal *domain.Palette) error {
	if _, err := fmt.Fprintf(w, "system: %q\n", pal.System); err != nil {
		return fmt.Errorf("write palette: %w", err)
	}
	if _, err := fmt.Fprintf(w, "name: %q\n", pal.Name); err != nil {
		return fmt.Errorf("write palette: %w", err)
	}
	if _, err := fmt.Fprintf(w, "author: %q\n", pal.Author); err != nil {
		return fmt.Errorf("write palette: %w", err)
	}
	if _, err := fmt.Fprintf(w, "variant: %q\n", pal.Variant); err != nil {
		return fmt.Errorf("write palette: %w", err)
	}
	if _, err := fmt.Fprintf(w, "palette:\n"); err != nil {
		return fmt.Errorf("write palette: %w", err)
	}
	for i, name := range slotNames {
		c := pal.Base(i)
		if _, err := fmt.Fprintf(w, "  %s: %q\n", name, c.Hex()); err != nil {
			return fmt.Errorf("write palette: %w", err)
		}
	}

	// Write overrides section if present.
	return writeOverrides(w, pal.Overrides)
}

// writeOverrides serializes the overrides map to YAML format.
// Keys are sorted alphabetically for deterministic output.
// Skipped when overrides is nil or empty.
func writeOverrides(w io.Writer, overrides map[string]domain.TokenOverride) error {
	if len(overrides) == 0 {
		return nil
	}

	if _, err := fmt.Fprintf(w, "overrides:\n"); err != nil {
		return fmt.Errorf("write overrides header: %w", err)
	}

	// Sort keys for deterministic output.
	keys := make([]string, 0, len(overrides))
	for k := range overrides {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		override := overrides[key]
		if _, err := fmt.Fprintf(w, "  %s:\n", key); err != nil {
			return fmt.Errorf("write override key: %w", err)
		}

		if err := writeOverrideFields(w, override); err != nil {
			return err
		}
	}

	return nil
}

// writeOverrideFields serializes individual override fields with omitempty behavior.
func writeOverrideFields(w io.Writer, o domain.TokenOverride) error {
	// Write color if present.
	if o.Color != nil {
		if _, err := fmt.Fprintf(w, "    color: %q\n", o.Color.Hex()); err != nil {
			return fmt.Errorf("write override color: %w", err)
		}
	}

	// Write style flags only when true (omitempty behavior).
	if o.Bold {
		if _, err := fmt.Fprintf(w, "    bold: true\n"); err != nil {
			return fmt.Errorf("write override bold: %w", err)
		}
	}
	if o.Italic {
		if _, err := fmt.Fprintf(w, "    italic: true\n"); err != nil {
			return fmt.Errorf("write override italic: %w", err)
		}
	}
	if o.Underline {
		if _, err := fmt.Fprintf(w, "    underline: true\n"); err != nil {
			return fmt.Errorf("write override underline: %w", err)
		}
	}
	if o.Undercurl {
		if _, err := fmt.Fprintf(w, "    undercurl: true\n"); err != nil {
			return fmt.Errorf("write override undercurl: %w", err)
		}
	}
	if o.Strikethrough {
		if _, err := fmt.Fprintf(w, "    strikethrough: true\n"); err != nil {
			return fmt.Errorf("write override strikethrough: %w", err)
		}
	}

	return nil
}

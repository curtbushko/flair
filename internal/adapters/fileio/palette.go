package fileio

import (
	"fmt"
	"io"

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
	return nil
}

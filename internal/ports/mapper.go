package ports

import "github.com/curtbushko/flair/internal/domain"

// MappedTheme is the output of a Mapper, consumed by a Generator.
// Each target defines its own concrete type (VimTheme, CssTheme, etc.)
// but the pipeline passes them through as any.
type MappedTheme = any

// Mapper transforms a ResolvedTheme into a target-specific theme struct.
type Mapper interface {
	// Name returns the target name (e.g. "vim", "css", "gtk").
	Name() string

	// Map transforms a ResolvedTheme into a target-specific mapped theme.
	Map(theme *domain.ResolvedTheme) (MappedTheme, error)
}

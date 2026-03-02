package application

import (
	"fmt"
	"io"

	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// expectedFiles lists the files that a complete theme directory should contain.
var expectedFiles = []string{
	"palette.yaml",
	"tokens.yaml",
	"vim-mapping.yaml",
	"css-mapping.yaml",
	"gtk-mapping.yaml",
	"qss-mapping.yaml",
	"stylix-mapping.yaml",
	"style.lua",
	"style.css",
	"gtk.css",
	"style.qss",
	"style.json",
}

// SchemaValidator reads from a reader and returns an error if the schema
// version is invalid or missing. The composition root wires this to the
// ValidatingReader adapter.
type SchemaValidator func(r io.Reader, kind domain.FileKind) error

// ValidateThemeUseCase checks a theme directory for completeness,
// schema correctness, and palette validity.
type ValidateThemeUseCase struct {
	store           ports.ThemeStore
	parser          ports.PaletteParser
	schemaValidator SchemaValidator
}

// NewValidateThemeUseCase returns a new ValidateThemeUseCase.
func NewValidateThemeUseCase(store ports.ThemeStore, parser ports.PaletteParser, sv SchemaValidator) *ValidateThemeUseCase {
	return &ValidateThemeUseCase{store: store, parser: parser, schemaValidator: sv}
}

// Execute validates the named theme directory and returns all violations found.
// A nil error with an empty slice means the theme is valid.
func (uc *ValidateThemeUseCase) Execute(themeName string) ([]string, error) {
	var violations []string

	// 1. Check for missing files.
	violations = append(violations, uc.checkMissingFiles(themeName)...)

	// 2. Validate palette.yaml schema version and content.
	paletteViolations, palette := uc.validatePalette(themeName)
	violations = append(violations, paletteViolations...)

	// 3. Run domain palette validation if we got a valid palette.
	if palette != nil {
		violations = append(violations, domain.ValidatePalette(palette)...)
	}

	// 4. Validate tokens.yaml schema version.
	violations = append(violations, uc.validateSchemaVersion(themeName, "tokens.yaml", domain.FileKindTokens)...)

	return violations, nil
}

// checkMissingFiles returns violations for each expected file that is absent.
func (uc *ValidateThemeUseCase) checkMissingFiles(themeName string) []string {
	var violations []string
	for _, f := range expectedFiles {
		if !uc.store.FileExists(themeName, f) {
			violations = append(violations, "missing file: "+f)
		}
	}
	return violations
}

// validatePalette checks palette.yaml schema version and parses its content.
// Schema check and content parse are done in separate reads so a schema issue
// does not prevent content validation.
func (uc *ValidateThemeUseCase) validatePalette(themeName string) ([]string, *domain.Palette) {
	var violations []string

	if !uc.store.FileExists(themeName, "palette.yaml") {
		return violations, nil
	}

	// Check schema version (separate read).
	violations = append(violations, uc.validateSchemaVersion(themeName, "palette.yaml", domain.FileKindPalette)...)

	// Parse palette content (separate read).
	rc, err := uc.store.OpenReader(themeName, "palette.yaml")
	if err != nil {
		violations = append(violations, fmt.Sprintf("palette.yaml: cannot open: %v", err))
		return violations, nil
	}
	defer func() { _ = rc.Close() }()

	palette, err := uc.parser.Parse(rc)
	if err != nil {
		violations = append(violations, fmt.Sprintf("palette.yaml: parse error: %v", err))
		return violations, nil
	}

	return violations, palette
}

// validateSchemaVersion checks the schema version of a file via the injected
// SchemaValidator. Returns violations for mismatches or missing headers.
func (uc *ValidateThemeUseCase) validateSchemaVersion(themeName, filename string, kind domain.FileKind) []string {
	if !uc.store.FileExists(themeName, filename) {
		return nil // Already reported by checkMissingFiles.
	}

	if uc.schemaValidator == nil {
		return nil
	}

	rc, err := uc.store.OpenReader(themeName, filename)
	if err != nil {
		return []string{fmt.Sprintf("%s: cannot open: %v", filename, err)}
	}
	defer func() { _ = rc.Close() }()

	if err := uc.schemaValidator(rc, kind); err != nil {
		return []string{fmt.Sprintf("%s: %v", filename, err)}
	}

	return nil
}

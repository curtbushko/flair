# Flair BDD Test Specifications

This document contains all BDD (Behavior-Driven Development) specifications for the Flair theme pipeline.
Tests are implemented using godog (cucumber for Go) and located in `features/`.

---

## Feature Checklist

Track testing progress by marking items as completed.

### Domain Layer

| Feature | Status | Location |
|---------|--------|----------|
| Color representation and parsing | TESTED | `features/domain/color.feature` |
| Color blending and manipulation | TESTED | `features/domain/color_ops.feature` |
| Base24 palette entity | TESTED | `features/domain/palette.feature` |
| Schema version management | TESTED | `features/domain/schema.feature` |
| Palette validation rules | TESTED | `features/domain/validation.feature` |

### Adapters Layer

| Feature | Status | Location |
|---------|--------|----------|
| ThemeStore filesystem operations | TESTED | `features/adapters/store.feature` |
| Built-in palette source | TESTED | `features/adapters/builtins.feature` |
| Reader/Writer wrappers | TESTED | `features/adapters/wrappers.feature` |
| CSS mapper + generator | TESTED | `features/adapters/css.feature` |
| GTK mapper + generator | TESTED | `features/adapters/gtk.feature` |
| QSS mapper + generator | TESTED | `features/adapters/qss.feature` |
| Stylix mapper + generator | TESTED | `features/application/generate.feature` |
| Vim mapper + generator | TESTED | `features/application/generate.feature` |

### Application Layer

| Feature | Status | Location |
|---------|--------|----------|
| Token derivation | TESTED | `features/application/generate.feature` |
| GenerateTheme orchestration | PARTIAL | See UNTESTED section |
| RegenerateTheme use case | UNTESTED | See UNTESTED section |
| ValidateTheme use case | UNTESTED | See UNTESTED section |
| SelectTheme use case | UNTESTED | See UNTESTED section |
| ListThemes use case | UNTESTED | See UNTESTED section |

### End-to-End

| Feature | Status | Location |
|---------|--------|----------|
| Full pipeline (all targets) | TESTED | `features/e2e/pipeline.feature` |
| Deterministic output | TESTED | `features/e2e/pipeline.feature` |
| Partial regeneration | UNTESTED | See UNTESTED section |

---

## TESTED Features

The following features have complete test coverage. Source files are in `features/`.

### Feature: Color representation and parsing

```gherkin
Feature: Color representation and parsing
  As a theme developer
  I need reliable color parsing and formatting
  So that color values are handled consistently throughout the pipeline

  Scenario: Parse a 6-digit hex color with hash
    Given the hex string "#7aa2f7"
    When I parse it as a Color
    Then the RGB values should be R=122 G=162 B=247

  Scenario: Parse a 6-digit hex color without hash
    Given the hex string "7aa2f7"
    When I parse it as a Color
    Then the RGB values should be R=122 G=162 B=247

  Scenario: Parse 3-digit shorthand
    Given the hex string "#f00"
    When I parse it as a Color
    Then the RGB values should be R=255 G=0 B=0

  Scenario: Parse white color
    Given the hex string "#ffffff"
    When I parse it as a Color
    Then the RGB values should be R=255 G=255 B=255

  Scenario: Parse black color
    Given the hex string "#000000"
    When I parse it as a Color
    Then the RGB values should be R=0 G=0 B=0

  Scenario: Reject invalid hex characters
    Given the hex string "#zzzzzz"
    When I parse it as a Color
    Then parsing should fail with a ParseError

  Scenario: Reject wrong length
    Given the hex string "#12345"
    When I parse it as a Color
    Then parsing should fail with a ParseError

  Scenario: Format color as hex
    Given the hex string "#7aa2f7"
    When I parse it as a Color
    Then the color formatted as hex should be "#7aa2f7"

  Scenario: RGB to HSL conversion for red
    Given the hex string "#ff0000"
    When I parse it as a Color
    And I convert it to HSL
    Then the HSL values should be H=0 S=1.0 L=0.5

  Scenario: HSL round-trip
    Given the hex string "#bb9af7"
    When I parse it as a Color
    And I convert it to HSL
    And I convert the HSL back to RGB
    Then the color formatted as hex should be "#bb9af7"

  Scenario: NONE color sentinel
    Given a NONE color
    Then IsNone should be true

  Scenario: Luminance of white
    Given the hex string "#ffffff"
    When I parse it as a Color
    Then the luminance should be approximately 1.0

  Scenario: Luminance of black
    Given the hex string "#000000"
    When I parse it as a Color
    Then the luminance should be approximately 0.0
```

### Feature: Color blending and manipulation

```gherkin
Feature: Color blending and manipulation
  As a theme developer
  I need color operations like blending and lightening
  So that I can derive consistent color variations

  Scenario: Blend 50% black and white
    Given two colors "#000000" and "#ffffff"
    When I blend them with ratio 0.5
    Then the result should be approximately "#808080"

  Scenario: Blend 0.0 returns source
    Given two colors "#ff0000" and "#00ff00"
    When I blend them with ratio 0.0
    Then the result should be approximately "#ff0000"

  Scenario: Blend 1.0 returns target
    Given two colors "#ff0000" and "#00ff00"
    When I blend them with ratio 1.0
    Then the result should be approximately "#00ff00"

  Scenario: Blend 25% creates subtle mix
    Given two colors "#1a1b26" and "#7aa2f7"
    When I blend them with ratio 0.25
    Then the result should be approximately "#323d5a"

  Scenario: Lighten a dark color
    Given the color "#1a1b26"
    When I lighten it by 0.1
    Then the luminance should be approximately 0.032

  Scenario: Darken a light color
    Given the color "#c0caf5"
    When I darken it by 0.1
    Then the luminance should be approximately 0.394

  Scenario: Desaturate removes color
    Given the color "#ff0000"
    When I desaturate it by 1.0
    Then I convert it to HSL
    And the HSL values should be H=0 S=0.0 L=0.5

  Scenario: Shift hue by 120 degrees
    Given the color "#ff0000"
    When I shift hue by 120 degrees
    Then the result should be approximately "#00ff00"

  Scenario: Shift hue wraps around
    Given the color "#ff0000"
    When I shift hue by 360 degrees
    Then the result should be approximately "#ff0000"
```

### Feature: Base24 palette entity

```gherkin
Feature: Base24 palette entity
  As a theme developer
  I need a palette entity that holds 24 color slots
  So that I can access colors by name or index

  Scenario: Full palette has 24 colors
    Given the Tokyo Night Dark palette from testdata
    Then the palette should have 24 colors

  Scenario: Access base00 slot
    Given the Tokyo Night Dark palette from testdata
    Then slot "base00" should be "#1a1b26"

  Scenario: Access base0D slot (blue)
    Given the Tokyo Night Dark palette from testdata
    Then slot "base0D" should be "#7aa2f7"

  Scenario: Access base0E slot (magenta)
    Given the Tokyo Night Dark palette from testdata
    Then slot "base0E" should be "#bb9af7"

  Scenario: Base index matches slot name
    Given the Tokyo Night Dark palette from testdata
    Then Base(0) should return the same as Slot("base00")
    And Base(13) should return the same as Slot("base0D")

  Scenario: Base16 fallback for base10 slot
    Given a base16 palette with only 16 colors
    Then base16 should be a fallback from base0

  Scenario: Base16 fallback for base12 slot
    Given a base16 palette with only 16 colors
    Then base18 should be a fallback from base8
```

### Feature: Schema version management

```gherkin
Feature: Schema version management
  As a theme developer
  I need schema versioning for file compatibility
  So that file format changes can be detected and handled

  Scenario: CurrentVersion for palette
    Given file kind "palette"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for universal
    Given file kind "universal"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for vim-mapping
    Given file kind "vim-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for css-mapping
    Given file kind "css-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for gtk-mapping
    Given file kind "gtk-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for qss-mapping
    Given file kind "qss-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for stylix-mapping
    Given file kind "stylix-mapping"
    Then CurrentVersion should return 1

  Scenario: All FileKind constants have version greater than 0
    Given all FileKind constants
    Then each should have a version greater than 0
```

### Feature: Palette validation rules

```gherkin
Feature: Palette validation rules
  As a theme developer
  I need palette validation
  So that color schemes meet quality requirements

  Scenario: Valid dark palette passes validation
    Given the Tokyo Night Dark palette from testdata
    When I validate the palette
    Then validation should pass with no errors

  Scenario: Validate luminance ordering for dark theme
    Given a dark palette where base00 is lighter than base05
    When I validate the palette
    Then validation should fail with luminance ordering error

  Scenario: Validate luminance ordering for light theme
    Given a light palette where base00 is darker than base05
    When I validate the palette
    Then validation should fail with luminance ordering error

  Scenario: Validate neutral ramp monotonicity warning
    Given a palette where base01 luminance is less than base00
    When I validate the palette
    Then validation should warn about monotonicity

  Scenario: Validate bright variants are brighter
    Given a palette where base12 is darker than base08
    When I validate the palette
    Then validation should warn about bright variant luminance
```

### Feature: ThemeStore filesystem operations

```gherkin
Feature: ThemeStore filesystem operations
  As the flair CLI
  I need to manage theme directories and symlinks
  So that themes can be stored and selected

  Scenario: EnsureThemeDir creates directory
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    Then the theme directory should exist

  Scenario: OpenWriter creates writable file
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    And I call OpenWriter("tokyonight", "test.yaml") and write "hello world"
    Then FileExists("tokyonight", "test.yaml") should return true

  Scenario: OpenReader reads existing file
    Given theme "tokyonight" exists with file "universal.yaml"
    When I call OpenWriter("tokyonight", "universal.yaml") and write "test content"
    And I call OpenReader("tokyonight", "universal.yaml")
    Then the content should be "test content"

  Scenario: Select creates symlinks
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    And I call Select("tokyonight")
    Then symlink "style.lua" should point to "tokyonight/style.lua"
    And symlink "style.json" should point to "tokyonight/style.json"

  Scenario: SelectedTheme reads symlink target
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    And I call Select("tokyonight")
    Then SelectedTheme should return "tokyonight"

  Scenario: FileExists returns false for missing file
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    Then FileExists("tokyonight", "nonexistent.yaml") should return false
```

### Feature: Built-in palette source

```gherkin
Feature: Built-in palette source
  As a theme developer
  I need access to built-in palettes
  So that I can generate themes without external files

  Scenario: List returns all embedded palette names
    When I call List() on the built-in source
    Then the result should contain "tokyo-night-dark"
    And the result should contain "gruvbox-dark"
    And the result should contain "catppuccin-mocha"
    And the result should be sorted alphabetically

  Scenario: Get returns YAML bytes for tokyo-night-dark
    When I call Get("tokyo-night-dark") on the built-in source
    Then I should receive valid YAML bytes

  Scenario: Get returns YAML bytes for gruvbox-dark
    When I call Get("gruvbox-dark") on the built-in source
    Then I should receive valid YAML bytes

  Scenario: Get returns YAML bytes for catppuccin-mocha
    When I call Get("catppuccin-mocha") on the built-in source
    Then I should receive valid YAML bytes

  Scenario: Get unknown name returns error
    When I call Get("nonexistent-theme") on the built-in source
    Then Get should return an error

  Scenario: Has returns true for tokyo-night-dark
    When I call Has("tokyo-night-dark") on the built-in source
    Then Has should return true

  Scenario: Has returns true for gruvbox-dark
    When I call Has("gruvbox-dark") on the built-in source
    Then Has should return true

  Scenario: Has returns false for unknown
    When I call Has("my-custom-theme") on the built-in source
    Then Has should return false
```

### Feature: Reader/Writer wrappers

```gherkin
Feature: Reader/Writer wrappers
  As the flair pipeline
  I need versioned writers and validating readers
  So that schema versions are handled consistently

  Scenario: VersionedWriter prepends schema header
    Given a VersionedWriter for kind "universal" and theme "tokyonight"
    When I write "tokens:\n  syntax.keyword: '#bb9af7'"
    Then the output should start with "schema_version: 1"
    And the output should contain "kind: universal"
    And the output should contain "theme_name: tokyonight"

  Scenario: VersionedWriter includes correct version for palette
    Given a VersionedWriter for kind "palette" and theme "gruvbox"
    When I write "palette data"
    Then the output should start with "schema_version: 1"
    And the output should contain "kind: palette"

  Scenario: VersionedWriter includes correct version for vim-mapping
    Given a VersionedWriter for kind "vim-mapping" and theme "catppuccin"
    When I write "highlights:"
    Then the output should start with "schema_version: 1"
    And the output should contain "kind: vim-mapping"

  Scenario: ValidatingReader passes valid schema version
    Given YAML with schema_version 1 for kind "universal"
    When I wrap it in ValidatingReader and read
    Then reading should succeed

  Scenario: ValidatingReader rejects outdated schema version
    Given YAML with schema_version 0 for kind "universal"
    When I wrap it in ValidatingReader and read
    Then reading should fail with SchemaVersionError
    And NeedsUpgrade should be false

  Scenario: ValidatingReader rejects future schema version
    Given YAML with schema_version 99 for kind "universal"
    When I wrap it in ValidatingReader and read
    Then reading should fail with SchemaVersionError
    And NeedsUpgrade should be true

  Scenario: ValidatingReader works with vim-mapping
    Given YAML with schema_version 1 for kind "vim-mapping"
    When I wrap it in ValidatingReader and read
    Then reading should succeed

  Scenario: ValidatingReader works with stylix-mapping
    Given YAML with schema_version 1 for kind "stylix-mapping"
    When I wrap it in ValidatingReader and read
    Then reading should succeed
```

### Feature: Token derivation and theme generation

```gherkin
Feature: Token derivation and theme generation
  As a theme developer
  I need to derive semantic tokens and generate output files
  So that I can create consistent themes across targets

  Scenario: Derive tokens from Tokyo Night Dark
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then the TokenSet should have at least 87 tokens

  Scenario: Surface tokens have correct values
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "surface.background" should have color "#1a1b26"

  Scenario: Syntax keyword token
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "syntax.keyword" should have color "#bb9af7"

  Scenario: Syntax comment is italic
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "syntax.comment" should be italic

  Scenario: Markup heading is bold
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "markup.heading" should be bold

  Scenario: Stylix mapper produces at least 60 values
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the Stylix mapper
    Then the StylixTheme should have at least 60 values

  Scenario: Vim mapper produces at least 200 highlight groups
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the Vim mapper
    Then the VimTheme should have at least 200 highlight groups
    And the VimTheme should have 16 terminal colors

  Scenario: Stylix generator produces valid JSON
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the Stylix mapper
    When I generate Stylix output
    Then the output should be valid JSON
    And the JSON should contain key "base00"

  Scenario: Vim generator produces Lua output
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the Vim mapper
    When I generate Vim output
    Then the generated output should contain "nvim_set_hl"
    And the generated output should contain "vim.cmd"
```

### Feature: CSS mapper and generator

```gherkin
Feature: CSS mapper and generator
  As a theme developer
  I need CSS output generation
  So that themes can be used in web applications

  Scenario: CSS mapper produces custom properties
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the CSS mapper
    Then the CssTheme should have custom properties

  Scenario: CSS mapper produces element rules
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the CSS mapper
    Then the CssTheme should have element rules

  Scenario: CSS generator produces valid CSS output
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the CSS mapper
    When I generate CSS output
    Then the generated output should contain ":root"
    And the generated output should contain "--"
```

### Feature: GTK mapper and generator

```gherkin
Feature: GTK mapper and generator
  As a theme developer
  I need GTK CSS output generation
  So that themes can be used in GTK applications

  Scenario: GTK mapper produces color definitions
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the GTK mapper
    Then the GtkTheme should have color definitions

  Scenario: GTK mapper produces widget rules
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the GTK mapper
    Then the GtkTheme should have widget rules

  Scenario: GTK generator produces valid GTK CSS
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the GTK mapper
    When I generate GTK output
    Then the generated output should contain "@define-color"
```

### Feature: QSS mapper and generator

```gherkin
Feature: QSS mapper and generator
  As a theme developer
  I need QSS output generation
  So that themes can be used in Qt applications

  Scenario: QSS mapper produces widget rules
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the QSS mapper
    Then the QssTheme should have widget rules

  Scenario: QSS generator produces valid QSS output
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the QSS mapper
    When I generate QSS output
    Then the generated output should contain "QWidget"
    And the generated output should contain "#"
```

### Feature: End-to-end pipeline

```gherkin
Feature: End-to-end pipeline
  As a theme developer
  I need the full pipeline to work correctly
  So that I can generate complete themes from palettes

  Scenario: Tokyo Night Dark full pipeline produces all 12 files
    When I run the full pipeline for "tokyo-night-dark"
    Then all 12 files should be created

  Scenario: Gruvbox Dark full pipeline produces all 12 files
    When I run the full pipeline for "gruvbox-dark"
    Then all 12 files should be created

  Scenario: Catppuccin Mocha full pipeline produces all 12 files
    When I run the full pipeline for "catppuccin-mocha"
    Then all 12 files should be created

  Scenario: Pipeline produces deterministic output
    When I run the full pipeline for "tokyo-night-dark"
    Then running the pipeline again should produce identical output
```

---

## UNTESTED Features

The following features need test implementation. These scenarios are extracted from
the original PLAN.md specifications but do not yet have step definitions or .feature files.

### Feature: GenerateTheme orchestration

**Status:** Needs CLI integration scenarios

```gherkin
Feature: GenerateTheme orchestration
  As a user
  I need the generate command to work end-to-end
  So that I can create themes from palettes

  Scenario: Full pipeline from file path
    Given a palette file at "./my-palette.yaml"
    When I execute GenerateTheme("./my-palette.yaml", "my-theme", "")
    Then theme dir should have palette.yaml
    And theme dir should have universal.yaml
    And theme dir should have 5 mapping files
    And theme dir should have 5 output files

  Scenario: Full pipeline from built-in name
    When I execute GenerateTheme("tokyo-night-dark", "", "")
    Then theme dir "tokyo-night-dark" should be created
    And all 12 files should exist

  Scenario: Built-in name infers theme name
    When I execute GenerateTheme("tokyo-night-dark", "", "")
    Then the theme name should be "tokyo-night-dark"

  Scenario: Single target filter generates 4 files
    When I execute GenerateTheme("tokyo-night-dark", "", "vim")
    Then only palette.yaml, universal.yaml, vim-mapping.yaml, style.lua should exist

  Scenario: Theme dir created if missing
    Given no theme directory exists
    When I execute GenerateTheme("tokyo-night-dark", "", "")
    Then the theme directory should be created

  Scenario: All files have correct schema versions
    When I execute GenerateTheme("tokyo-night-dark", "", "")
    Then all intermediate files should have schema_version: 1

  Scenario: One target failure does not block others
    Given a mapper that fails for vim target
    When I execute GenerateTheme with all targets
    Then other targets should still generate successfully
    And an error should be reported for vim target
```

### Feature: CLI select command

**Status:** Needs CLI integration scenarios

```gherkin
Feature: CLI select command
  As a user
  I need to select the active theme
  So that applications use my chosen theme

  Scenario: Creates symlinks to all 5 output files
    Given theme "tokyonight" exists with all output files
    When I run "flair select tokyonight"
    Then symlink "style.lua" should point to "tokyonight/style.lua"
    And symlink "style.css" should point to "tokyonight/style.css"
    And symlink "gtk.css" should point to "tokyonight/gtk.css"
    And symlink "style.qss" should point to "tokyonight/style.qss"
    And symlink "style.json" should point to "tokyonight/style.json"

  Scenario: Non-existent theme returns error
    Given theme "nonexistent" does not exist
    When I run "flair select nonexistent"
    Then the command should fail
    And the error should mention "theme not found"

  Scenario: Incomplete theme returns error
    Given theme "partial" exists without style.lua
    When I run "flair select partial"
    Then the command should fail
    And the error should list missing files
```

### Feature: CLI list command

**Status:** Needs CLI integration scenarios

```gherkin
Feature: CLI list command
  As a user
  I need to see available themes
  So that I can choose which theme to use

  Scenario: Shows installed themes with selected marker
    Given themes "gruvbox" and "tokyonight" exist
    And "tokyonight" is selected
    When I run "flair list"
    Then output should show "  gruvbox (dark)"
    And output should show "* tokyonight (dark)"

  Scenario: No themes shows helpful message
    Given no themes exist
    When I run "flair list"
    Then output should show "No themes installed"

  Scenario: List built-in palettes
    When I run "flair list --builtins"
    Then output should show "catppuccin-mocha"
    And output should show "gruvbox-dark"
    And output should show "tokyo-night-dark"
```

### Feature: CLI regenerate command

**Status:** Needs CLI integration scenarios

```gherkin
Feature: CLI regenerate command
  As a user
  I need to regenerate themes after editing intermediate files
  So that my customizations are applied to output files

  Scenario: Edit palette.yaml regenerates everything
    Given theme "tokyonight" exists
    And I modify palette.yaml to change base00
    When I run "flair regenerate tokyonight"
    Then universal.yaml should be regenerated
    And all mapping files should be regenerated
    And all output files should be regenerated
    And base00 change should be reflected in outputs

  Scenario: Edit universal.yaml regenerates mappings and outputs
    Given theme "tokyonight" exists
    And I modify universal.yaml to change syntax.keyword color
    When I run "flair regenerate tokyonight"
    Then palette.yaml should NOT be modified
    And all mapping files should be regenerated
    And all output files should be regenerated

  Scenario: Edit vim-mapping.yaml regenerates only style.lua
    Given theme "tokyonight" exists
    And I modify vim-mapping.yaml to add a highlight group
    When I run "flair regenerate tokyonight"
    Then only style.lua should be regenerated
    And other output files should NOT be modified

  Scenario: No edits detected shows nothing to do
    Given theme "tokyonight" exists with no modifications
    When I run "flair regenerate tokyonight"
    Then output should show "nothing to regenerate"
```

### Feature: Advanced end-to-end scenarios

**Status:** Needs additional E2E coverage

```gherkin
Feature: Advanced end-to-end scenarios
  As a developer
  I need comprehensive E2E coverage
  So that the full system is validated

  Scenario: Generate from built-in produces identical output to file
    Given I generate theme from built-in "tokyo-night-dark"
    And I generate theme from file "tokyo-night-dark.yaml"
    Then both themes should have identical output files

  Scenario: go-arch-lint check passes
    When I run "go-arch-lint check"
    Then the command should exit with code 0
    And no architecture violations should be reported

  Scenario: All built-in palettes parse and validate cleanly
    Given all built-in palettes
    When I parse and validate each
    Then no errors should occur
    And no warnings should be reported
```

---

## Running Tests

```bash
# Run all BDD tests
go test ./features/...

# Run with verbose output
go test ./features/... -v

# Run specific feature file
go test ./features/... -godog.paths=features/domain/color.feature

# Run scenarios matching a pattern
go test ./features/... -godog.tags="@wip"
```

---

## Implementation Notes

### Step Definitions

All step definitions are in `features/steps/common.go`. When adding new scenarios:

1. Add step definitions in `registerXxxSteps()` functions
2. Use `TestContext` to share state between steps
3. Perform real validation against domain types and adapters
4. Return concrete errors, not stub implementations

### Test Data

- `testdata/tokyo-night-dark.yaml` - Reference base24 palette
- Built-in palettes in `internal/adapters/palettes/`
- Temp directories created per-scenario via Before/After hooks

### Adding New Features

1. Add feature file to appropriate directory (`domain/`, `adapters/`, `application/`, `e2e/`)
2. Update godog_test.go `Paths` if adding new directory
3. Add step definitions in common.go
4. Update Feature Checklist in this document
5. Move from UNTESTED to TESTED section once complete

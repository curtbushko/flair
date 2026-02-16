# BDD Test Implementation Plan

## Context

The PLAN.md contains detailed Gherkin feature specifications (Features 1.2-5.6) that describe the expected behavior of flair's components. Currently these serve as documentation but aren't executable. By implementing godog tests, we codify these specifications as living documentation that validates the implementation.

---

## Critical: Real Validation, Not Stubs

Every step definition MUST perform actual validation:

- **Given steps**: Set up real state (not mocks that bypass logic)
- **When steps**: Execute real code paths (not stubs)
- **Then steps**: Assert concrete values with specific expectations

**Anti-patterns to avoid:**
```go
// BAD: Always passes
func theRGBValuesShouldBe(r, g, b int) error {
    return nil  // NO! Must actually check values
}

// GOOD: Real assertion
func theRGBValuesShouldBe(r, g, b int) error {
    if ctx.color.R != uint8(r) || ctx.color.G != uint8(g) || ctx.color.B != uint8(b) {
        return fmt.Errorf("RGB = (%d,%d,%d), want (%d,%d,%d)",
            ctx.color.R, ctx.color.G, ctx.color.B, r, g, b)
    }
    return nil
}
```

**Validation requirements:**
1. Parse results checked against expected values (not just "no error")
2. Generated output compared byte-for-byte against golden files
3. Error types verified with `errors.As()` for specific error kinds
4. File existence/content verified by reading actual filesystem
5. Color math validated with tolerance for floating-point (HSL conversions)

---

## Directory Structure

```
features/
├── domain/
│   ├── color.feature           # Feature 1.2: Color parsing
│   ├── color_ops.feature       # Feature 1.3: Color operations
│   ├── palette.feature         # Feature 1.4: Palette entity
│   └── schema.feature          # Feature 1.9: Schema versions
├── adapters/
│   ├── store.feature           # Feature 1.14: ThemeStore
│   ├── builtins.feature        # Feature 1.15: Built-in palettes
│   └── wrappers.feature        # Feature 1.16-1.17: Versioned/Validating
├── derivation/
│   └── tokens.feature          # Feature 2.1-2.5: Token derivation
├── targets/
│   ├── stylix.feature          # Feature 3.1: Stylix
│   └── vim.feature             # Feature 3.3: Vim
├── application/
│   └── generate.feature        # Feature 3.6: GenerateTheme
├── cli/
│   ├── select.feature          # Feature 4.3: select command
│   ├── list.feature            # Feature 4.7: list command
│   └── regenerate.feature      # Feature 4.8: regenerate command
├── e2e/
│   └── pipeline.feature        # Feature 5.1-5.6: End-to-end
├── steps/
│   ├── common_test.go          # Shared context, helpers
│   ├── domain_test.go          # Domain step definitions
│   ├── adapter_test.go         # Adapter step definitions
│   ├── application_test.go     # Use case step definitions
│   ├── cli_test.go             # CLI step definitions
│   └── e2e_test.go             # E2E step definitions
└── godog_test.go               # Test runner with TestMain
```

---

## Validation Strategy Per Feature

| Feature | What We Assert | How |
|---------|----------------|-----|
| Color parsing | Exact R,G,B values | `color.R == 122 && color.G == 162 && color.B == 247` |
| Color ops (Blend) | Calculated RGB within tolerance | `math.Abs(float64(got.R) - float64(want.R)) < 2` |
| HSL round-trip | Original hex recovered | `color.Hex() == originalHex` |
| Luminance | WCAG formula result | `math.Abs(lum - 1.0) < 0.001` for white |
| Palette slots | All 24 colors accessible | Loop through `base00`-`base17`, verify non-zero |
| Schema versions | Correct version per FileKind | `domain.CurrentVersion(kind) == expectedVersion` |
| ThemeStore | Files exist on disk | `os.Stat(path)` returns no error |
| Symlinks | Target matches expected | `os.Readlink()` returns correct path |
| VersionedWriter | Header in output | `bytes.HasPrefix(buf.Bytes(), expectedHeader)` |
| ValidatingReader | Error type for bad version | `errors.As(err, &schemaVersionError)` |
| Token derivation | Token count >= 87 | `tokenSet.Len() >= 87` |
| Token values | Specific tokens have expected colors | `tokenSet.MustGet("syntax.keyword").Color.Hex() == "#bb9af7"` |
| Stylix output | Valid JSON, sorted keys | `json.Valid()` + key order check |
| Vim output | Contains `nvim_set_hl` calls | `strings.Contains(output, "nvim_set_hl")` |
| Golden files | Byte-identical match | `bytes.Equal(got, golden)` |
| Deterministic | Two runs produce same output | Generate twice, compare bytes |

---

## Implementation Checklist

### Phase 0: Setup

- [ ] 0.1 — Add godog dependency (`go get github.com/cucumber/godog`)
- [ ] 0.2 — Create `features/` directory structure
- [ ] 0.3 — Create test runner (`features/godog_test.go`)
- [ ] 0.4 — Create shared context (`features/steps/common_test.go`)
- [ ] 0.5 — Add `test:bdd` task to Taskfile.yml
- [ ] 0.6 — Verify `go test ./features/...` runs (no scenarios yet)

### Phase 1: Domain Features

- [ ] 1.1 — `features/domain/color.feature` (Feature 1.2: Color parsing)
  - [ ] Scenario: Parse 6-digit hex with hash
  - [ ] Scenario: Parse 6-digit hex without hash
  - [ ] Scenario: Parse 3-digit shorthand
  - [ ] Scenario: Reject invalid hex
  - [ ] Scenario: Reject wrong length
  - [ ] Scenario: Format as hex
  - [ ] Scenario: RGB to HSL conversion
  - [ ] Scenario: HSL round-trip
  - [ ] Scenario: NONE sentinel
  - [ ] Scenario: Luminance calculations

- [ ] 1.2 — `features/steps/domain_test.go` — Color step definitions
  - [ ] Step: `the hex string {string}`
  - [ ] Step: `I parse it as a Color`
  - [ ] Step: `the RGB values should be R={int} G={int} B={int}`
  - [ ] Step: `parsing should fail with ParseError`
  - [ ] Step: `the hex output should be {string}`
  - [ ] Step: `the luminance should be approximately {float}`

- [ ] 1.3 — `features/domain/color_ops.feature` (Feature 1.3: Color operations)
  - [ ] Scenario: Blend 50% black + white
  - [ ] Scenario: Blend 0.0 returns source
  - [ ] Scenario: Blend 1.0 returns target
  - [ ] Scenario: BlendBg equivalence
  - [ ] Scenario: Lighten edge cases
  - [ ] Scenario: Darken edge cases
  - [ ] Scenario: Desaturate edge cases
  - [ ] Scenario: ShiftHue wrap-around

- [ ] 1.4 — Step definitions for color operations
  - [ ] Step: `I blend {string} and {string} at {float}`
  - [ ] Step: `the result should be {string}`
  - [ ] Step: `I lighten {string} by {float}`
  - [ ] Step: `I darken {string} by {float}`

- [ ] 1.5 — `features/domain/palette.feature` (Feature 1.4: Palette entity)
  - [ ] Scenario: Construct full 24-color palette
  - [ ] Scenario: Access by name (base0D)
  - [ ] Scenario: Access by index (13)
  - [ ] Scenario: Base16 fallbacks applied
  - [ ] Scenario: Missing required slot error

- [ ] 1.6 — Step definitions for palette
  - [ ] Step: `a palette with {int} colors`
  - [ ] Step: `I access slot {string}`
  - [ ] Step: `the color should be {string}`
  - [ ] Step: `palette construction should fail`

- [ ] 1.7 — `features/domain/schema.feature` (Feature 1.9: Schema versions)
  - [ ] Scenario: CurrentVersion returns correct version per FileKind
  - [ ] Scenario: All FileKind constants have version > 0

- [ ] 1.8 — Step definitions for schema
  - [ ] Step: `FileKind {string}`
  - [ ] Step: `CurrentVersion should return {int}`

### Phase 2: Adapter Features

- [ ] 2.1 — `features/adapters/store.feature` (Feature 1.14: ThemeStore)
  - [ ] Scenario: EnsureThemeDir creates directory
  - [ ] Scenario: ListThemes returns sorted names
  - [ ] Scenario: Select creates symlinks
  - [ ] Scenario: Select replaces existing symlinks
  - [ ] Scenario: SelectedTheme reads symlink target
  - [ ] Scenario: OpenWriter creates writable file
  - [ ] Scenario: OpenReader reads existing file
  - [ ] Scenario: Read/Write round-trip
  - [ ] Scenario: FileExists and FileMtime

- [ ] 2.2 — `features/steps/adapter_test.go` — Store step definitions
  - [ ] Step: `a temporary config directory`
  - [ ] Step: `I ensure theme dir {string}`
  - [ ] Step: `directory {string} should exist`
  - [ ] Step: `I select theme {string}`
  - [ ] Step: `symlink {string} should point to {string}`

- [ ] 2.3 — `features/adapters/builtins.feature` (Feature 1.15: Built-in palettes)
  - [ ] Scenario: List returns all embedded names (sorted)
  - [ ] Scenario: Get returns valid YAML bytes
  - [ ] Scenario: Get unknown name returns error
  - [ ] Scenario: Has returns true for built-in
  - [ ] Scenario: Has returns false for unknown

- [ ] 2.4 — Step definitions for builtins
  - [ ] Step: `I call List on PaletteSource`
  - [ ] Step: `the result should contain {string}`
  - [ ] Step: `I call Get for {string}`
  - [ ] Step: `I should receive valid base24 YAML`
  - [ ] Step: `I should receive an error`

- [ ] 2.5 — `features/adapters/wrappers.feature` (Features 1.16-1.17)
  - [ ] Scenario: VersionedWriter prepends header on first write
  - [ ] Scenario: Header written only once
  - [ ] Scenario: Correct version per file kind
  - [ ] Scenario: ValidatingReader passes valid version
  - [ ] Scenario: ValidatingReader rejects outdated version
  - [ ] Scenario: ValidatingReader rejects future version
  - [ ] Scenario: Wrappers composable with other readers

- [ ] 2.6 — Step definitions for wrappers
  - [ ] Step: `a VersionedWriter for {string}`
  - [ ] Step: `I write {string}`
  - [ ] Step: `the output should start with schema_version: {int}`
  - [ ] Step: `a ValidatingReader with schema_version: {int}`
  - [ ] Step: `reading should succeed`
  - [ ] Step: `reading should fail with SchemaVersionError`

### Phase 3: Derivation Features

- [ ] 3.1 — `features/derivation/tokens.feature` (Features 2.1-2.5)
  - [ ] Scenario: Full derivation produces >= 87 tokens
  - [ ] Scenario: Surface tokens match derivation table
  - [ ] Scenario: Text tokens match derivation table
  - [ ] Scenario: Status tokens match derivation table
  - [ ] Scenario: Syntax tokens match derivation table
  - [ ] Scenario: Terminal ANSI colors correct
  - [ ] Scenario: Write universal.yaml via io.Writer
  - [ ] Scenario: Read universal.yaml round-trip
  - [ ] Scenario: ValidatingReader catches version mismatch

- [ ] 3.2 — Step definitions for derivation
  - [ ] Step: `the Tokyo Night Dark palette`
  - [ ] Step: `I derive tokens`
  - [ ] Step: `the token count should be at least {int}`
  - [ ] Step: `token {string} should have color {string}`
  - [ ] Step: `token {string} should have italic style`

### Phase 4: Target Features

- [ ] 4.1 — `features/targets/stylix.feature` (Feature 3.1)
  - [ ] Scenario: Mapper produces >= 60 keys
  - [ ] Scenario: Mapping file write/read round-trip
  - [ ] Scenario: Generator writes valid JSON
  - [ ] Scenario: JSON has sorted keys
  - [ ] Scenario: JSON has 2-space indent

- [ ] 4.2 — `features/targets/vim.feature` (Feature 3.3)
  - [ ] Scenario: >= 200 highlight groups
  - [ ] Scenario: 16 terminal ANSI colors
  - [ ] Scenario: Mapping file round-trip
  - [ ] Scenario: Generator writes hi clear
  - [ ] Scenario: Generator writes nvim_set_hl calls
  - [ ] Scenario: Generator writes terminal_color_N

- [ ] 4.3 — Step definitions for targets
  - [ ] Step: `I map a ResolvedTheme to Stylix`
  - [ ] Step: `the StylixTheme should have at least {int} keys`
  - [ ] Step: `I generate Stylix output`
  - [ ] Step: `the output should be valid JSON`
  - [ ] Step: `the VimTheme should have at least {int} highlights`

### Phase 5: Application Features

- [ ] 5.1 — `features/application/generate.feature` (Feature 3.6)
  - [ ] Scenario: Full pipeline from file creates 12 files
  - [ ] Scenario: Full pipeline from built-in name
  - [ ] Scenario: Built-in name infers theme name
  - [ ] Scenario: Single target filter creates 4 files
  - [ ] Scenario: Theme dir created if missing
  - [ ] Scenario: All files have correct schema versions

- [ ] 5.2 — `features/steps/application_test.go`
  - [ ] Step: `I run GenerateTheme with palette {string}`
  - [ ] Step: `theme directory {string} should have {int} files`
  - [ ] Step: `file {string} should exist in theme {string}`
  - [ ] Step: `file {string} should have schema_version {int}`

### Phase 6: CLI Features

- [ ] 6.1 — `features/cli/select.feature` (Feature 4.3)
  - [ ] Scenario: Creates symlinks to all 5 output files
  - [ ] Scenario: Non-existent theme returns error
  - [ ] Scenario: Incomplete theme lists missing files

- [ ] 6.2 — `features/cli/list.feature` (Feature 4.7)
  - [ ] Scenario: Shows installed themes with selected marker
  - [ ] Scenario: No themes shows helpful message
  - [ ] Scenario: --builtins lists built-in palettes

- [ ] 6.3 — `features/cli/regenerate.feature` (Feature 4.8)
  - [ ] Scenario: Edit palette.yaml regenerates everything
  - [ ] Scenario: Edit universal.yaml regenerates mappings + outputs
  - [ ] Scenario: Edit vim-mapping.yaml regenerates only style.lua
  - [ ] Scenario: No edits shows "nothing to do"

- [ ] 6.4 — `features/steps/cli_test.go`
  - [ ] Step: `I run flair {string}`
  - [ ] Step: `the exit code should be {int}`
  - [ ] Step: `stdout should contain {string}`
  - [ ] Step: `stderr should contain {string}`

### Phase 7: End-to-End Features

- [ ] 7.1 — `features/e2e/pipeline.feature` (Features 5.1-5.6)
  - [ ] Scenario: Tokyo Night Dark full pipeline matches golden files
  - [ ] Scenario: Built-in produces identical output to file
  - [ ] Scenario: Additional schemes produce valid output
  - [ ] Scenario: Regenerate after edit produces correct partial rebuild
  - [ ] Scenario: Pipeline is deterministic (byte-identical on reruns)
  - [ ] Scenario: All built-in palettes parse and validate

- [ ] 7.2 — `features/steps/e2e_test.go`
  - [ ] Step: `I generate theme from {string}`
  - [ ] Step: `output file {string} should match golden file`
  - [ ] Step: `I generate theme twice`
  - [ ] Step: `both outputs should be byte-identical`
  - [ ] Step: `all built-in palettes should parse without error`

### Phase 8: Verification

- [ ] 8.1 — All BDD tests pass: `go test -v ./features/...`
- [ ] 8.2 — Taskfile integration: `task test:bdd` works
- [ ] 8.3 — Mutation test: break implementation, verify tests fail
- [ ] 8.4 — CI integration (if applicable)

---

## Dependencies

```
github.com/cucumber/godog v0.14.1
```

---

## Taskfile Addition

```yaml
test:bdd:
  desc: Run BDD acceptance tests
  cmds:
    - go test -v ./features/...
```

---

## Reusable Assets

| Asset | Location | Purpose |
|-------|----------|---------|
| Tokyo Night palette | `testdata/tokyo-night-dark.yaml` | Canonical test fixture |
| Golden files | `testdata/expected/*` | Reference outputs for comparison |
| Stub implementations | `internal/application/*_test.go` | Pattern for in-memory stores |

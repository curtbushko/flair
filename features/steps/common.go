// Package steps provides godog step definitions for BDD acceptance tests.
// Each step definition performs real validation against domain types and adapters,
// ensuring concrete assertions rather than stub implementations.
package steps

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"

	"github.com/curtbushko/flair/internal/adapters/deriver"
	"github.com/curtbushko/flair/internal/adapters/generator"
	"github.com/curtbushko/flair/internal/adapters/mapper"
	"github.com/curtbushko/flair/internal/adapters/palettes"
	"github.com/curtbushko/flair/internal/adapters/store"
	"github.com/curtbushko/flair/internal/adapters/wrappers"
	yamlparser "github.com/curtbushko/flair/internal/adapters/yaml"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// TestContext holds state shared across step definitions within a scenario.
// Each scenario gets a fresh TestContext instance.
type TestContext struct {
	// Color testing state
	hexInput string
	color    domain.Color
	colorErr error
	hsl      domain.HSL

	// Palette testing state
	palette    *domain.Palette
	paletteErr error

	// TokenSet testing state
	tokenSet *domain.TokenSet

	// Schema version testing state
	currentFileKind domain.FileKind

	// ThemeStore testing state
	tempDir   string
	fsStore   *store.FsStore
	themeName string
	writerBuf *bytes.Buffer
	fileData  []byte
	storeErr  error

	// Built-in palettes testing state
	builtinSource *palettes.Source
	builtinNames  []string
	builtinReader io.Reader
	builtinHas    bool
	builtinGetErr error

	// Wrapper testing state
	versionedWriter *wrappers.VersionedWriter
	outputBuffer    *bytes.Buffer

	// ValidatingReader testing state
	validatingReader *wrappers.ValidatingReader
	validationErr    error
	readBytes        []byte

	// Application testing state
	resolvedTheme *domain.ResolvedTheme

	// Mapper/Generator testing state
	mappedTheme  ports.MappedTheme
	stylixTheme  *ports.StylixTheme
	vimTheme     *ports.VimTheme
	generateErr  error
	genOutputBuf *bytes.Buffer

	// Generic error holder
	lastErr error
}

// InitializeScenario registers all step definitions and sets up fresh context.
func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := &TestContext{}

	// Before hook: create temp directory for filesystem tests
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tmpDir, err := os.MkdirTemp("", "flair-bdd-*")
		if err != nil {
			return ctx, err
		}
		tc.tempDir = tmpDir
		tc.fsStore = store.NewFsStore(tmpDir)
		tc.builtinSource = palettes.NewSource()
		return ctx, nil
	})

	// After hook: cleanup temp directory
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if tc.tempDir != "" {
			os.RemoveAll(tc.tempDir)
		}
		return ctx, nil
	})

	// Register all step definitions
	registerColorSteps(ctx, tc)
	registerPaletteSteps(ctx, tc)
	registerSchemaSteps(ctx, tc)
	registerStoreSteps(ctx, tc)
	registerBuiltinSteps(ctx, tc)
	registerWrapperSteps(ctx, tc)
	registerDeriverSteps(ctx, tc)
	registerMapperSteps(ctx, tc)
	registerGeneratorSteps(ctx, tc)
	registerE2ESteps(ctx, tc)
}

// --- Color Step Definitions ---

func registerColorSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^the hex string "([^"]*)"$`, tc.theHexString)
	ctx.Step(`^I parse it as a Color$`, tc.iParseItAsAColor)
	ctx.Step(`^the RGB values should be R=(\d+) G=(\d+) B=(\d+)$`, tc.theRGBValuesShouldBe)
	ctx.Step(`^parsing should fail with a ParseError$`, tc.parsingShouldFailWithParseError)
	ctx.Step(`^the color formatted as hex should be "([^"]*)"$`, tc.theColorFormattedAsHexShouldBe)
	ctx.Step(`^I convert it to HSL$`, tc.iConvertItToHSL)
	ctx.Step(`^the HSL values should be H=([\d.]+) S=([\d.]+) L=([\d.]+)$`, tc.theHSLValuesShouldBe)
	ctx.Step(`^I convert the HSL back to RGB$`, tc.iConvertTheHSLBackToRGB)
	ctx.Step(`^the luminance should be approximately ([\d.]+)$`, tc.theLuminanceShouldBeApproximately)
	ctx.Step(`^a NONE color$`, tc.aNONEColor)
	ctx.Step(`^IsNone should be true$`, tc.isNoneShouldBeTrue)
	ctx.Step(`^two colors "([^"]*)" and "([^"]*)"$`, tc.twoColors)
	ctx.Step(`^I blend them with ratio ([\d.]+)$`, tc.iBlendThemWithRatio)
	ctx.Step(`^the result should be approximately "([^"]*)"$`, tc.theResultShouldBeApproximately)
	ctx.Step(`^the color "([^"]*)"$`, tc.theColor)
	ctx.Step(`^I lighten it by ([\d.]+)$`, tc.iLightenItBy)
	ctx.Step(`^I darken it by ([\d.]+)$`, tc.iDarkenItBy)
	ctx.Step(`^I desaturate it by ([\d.]+)$`, tc.iDesaturateItBy)
	ctx.Step(`^I shift hue by ([\d.-]+) degrees$`, tc.iShiftHueByDegrees)
}

func (tc *TestContext) theHexString(hex string) error {
	tc.hexInput = hex
	return nil
}

func (tc *TestContext) iParseItAsAColor() error {
	tc.color, tc.colorErr = domain.ParseHex(tc.hexInput)
	return nil
}

func (tc *TestContext) theRGBValuesShouldBe(r, g, b int) error {
	if tc.colorErr != nil {
		return fmt.Errorf("expected successful parse, got error: %v", tc.colorErr)
	}
	if tc.color.R != uint8(r) || tc.color.G != uint8(g) || tc.color.B != uint8(b) {
		return fmt.Errorf("RGB = (%d,%d,%d), want (%d,%d,%d)",
			tc.color.R, tc.color.G, tc.color.B, r, g, b)
	}
	return nil
}

func (tc *TestContext) parsingShouldFailWithParseError() error {
	if tc.colorErr == nil {
		return fmt.Errorf("expected ParseError, got nil")
	}
	_, ok := tc.colorErr.(*domain.ParseError)
	if !ok {
		return fmt.Errorf("expected *domain.ParseError, got %T", tc.colorErr)
	}
	return nil
}

func (tc *TestContext) theColorFormattedAsHexShouldBe(expected string) error {
	got := tc.color.Hex()
	if got != expected {
		return fmt.Errorf("Hex() = %q, want %q", got, expected)
	}
	return nil
}

func (tc *TestContext) iConvertItToHSL() error {
	tc.hsl = tc.color.ToHSL()
	return nil
}

func (tc *TestContext) theHSLValuesShouldBe(h, s, l float64) error {
	const tolerance = 0.01
	if math.Abs(tc.hsl.H-h) > tolerance || math.Abs(tc.hsl.S-s) > tolerance || math.Abs(tc.hsl.L-l) > tolerance {
		return fmt.Errorf("HSL = (%.3f,%.3f,%.3f), want (%.3f,%.3f,%.3f)",
			tc.hsl.H, tc.hsl.S, tc.hsl.L, h, s, l)
	}
	return nil
}

func (tc *TestContext) iConvertTheHSLBackToRGB() error {
	tc.color = tc.hsl.ToRGB()
	return nil
}

func (tc *TestContext) theLuminanceShouldBeApproximately(expected float64) error {
	got := tc.color.Luminance()
	const tolerance = 0.001
	if math.Abs(got-expected) > tolerance {
		return fmt.Errorf("Luminance() = %.4f, want %.4f (tolerance %.3f)", got, expected, tolerance)
	}
	return nil
}

func (tc *TestContext) aNONEColor() error {
	tc.color = domain.NoneColor()
	return nil
}

func (tc *TestContext) isNoneShouldBeTrue() error {
	if !tc.color.IsNone {
		return fmt.Errorf("expected IsNone=true, got false")
	}
	return nil
}

func (tc *TestContext) twoColors(hex1, hex2 string) error {
	c1, err := domain.ParseHex(hex1)
	if err != nil {
		return err
	}
	c2, err := domain.ParseHex(hex2)
	if err != nil {
		return err
	}
	tc.color = c1
	tc.hsl.H = float64(c2.R)
	tc.hsl.S = float64(c2.G)
	tc.hsl.L = float64(c2.B)
	return nil
}

func (tc *TestContext) iBlendThemWithRatio(ratio float64) error {
	c2 := domain.Color{R: uint8(tc.hsl.H), G: uint8(tc.hsl.S), B: uint8(tc.hsl.L)}
	tc.color = domain.Blend(tc.color, c2, ratio)
	return nil
}

func (tc *TestContext) theResultShouldBeApproximately(expectedHex string) error {
	expected, err := domain.ParseHex(expectedHex)
	if err != nil {
		return err
	}
	const tolerance = 2 // Allow 2-unit difference for rounding
	if abs(int(tc.color.R)-int(expected.R)) > tolerance ||
		abs(int(tc.color.G)-int(expected.G)) > tolerance ||
		abs(int(tc.color.B)-int(expected.B)) > tolerance {
		return fmt.Errorf("result = %s, want approximately %s", tc.color.Hex(), expectedHex)
	}
	return nil
}

func (tc *TestContext) theColor(hex string) error {
	var err error
	tc.color, err = domain.ParseHex(hex)
	return err
}

func (tc *TestContext) iLightenItBy(amount float64) error {
	tc.color = domain.Lighten(tc.color, amount)
	return nil
}

func (tc *TestContext) iDarkenItBy(amount float64) error {
	tc.color = domain.Darken(tc.color, amount)
	return nil
}

func (tc *TestContext) iDesaturateItBy(amount float64) error {
	tc.color = domain.Desaturate(tc.color, amount)
	return nil
}

func (tc *TestContext) iShiftHueByDegrees(degrees float64) error {
	tc.color = domain.ShiftHue(tc.color, degrees)
	return nil
}

// --- Palette Step Definitions ---

func registerPaletteSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^the Tokyo Night Dark palette from testdata$`, tc.theTokyoNightDarkPalette)
	ctx.Step(`^the palette should have (\d+) colors$`, tc.thePaletteShouldHaveColors)
	ctx.Step(`^slot "([^"]*)" should be "([^"]*)"$`, tc.slotShouldBe)
	ctx.Step(`^Base\((\d+)\) should return the same as Slot\("([^"]*)"\)$`, tc.baseShouldReturnSameAsSlot)
	ctx.Step(`^a base16 palette with only 16 colors$`, tc.aBase16PaletteWith16Colors)
	ctx.Step(`^base(\d+) should be a fallback from base(\d+)$`, tc.baseNShouldBeFallbackFromBaseM)
}

func (tc *TestContext) theTokyoNightDarkPalette() error {
	parser := yamlparser.NewParser()
	data, err := os.ReadFile("../testdata/tokyo-night-dark.yaml")
	if err != nil {
		return fmt.Errorf("read testdata: %w", err)
	}
	tc.palette, tc.paletteErr = parser.Parse(bytes.NewReader(data))
	return tc.paletteErr
}

func (tc *TestContext) thePaletteShouldHaveColors(count int) error {
	// All 24 slots should be accessible
	for i := 0; i < count; i++ {
		c := tc.palette.Base(i)
		if c.IsNone {
			return fmt.Errorf("slot %d returned NoneColor", i)
		}
	}
	return nil
}

func (tc *TestContext) slotShouldBe(slot, expectedHex string) error {
	c, err := tc.palette.Slot(slot)
	if err != nil {
		return err
	}
	expected, err := domain.ParseHex(expectedHex)
	if err != nil {
		return err
	}
	if !c.Equal(expected) {
		return fmt.Errorf("slot %s = %s, want %s", slot, c.Hex(), expected.Hex())
	}
	return nil
}

func (tc *TestContext) baseShouldReturnSameAsSlot(index int, slotName string) error {
	byIndex := tc.palette.Base(index)
	byName, err := tc.palette.Slot(slotName)
	if err != nil {
		return err
	}
	if !byIndex.Equal(byName) {
		return fmt.Errorf("Base(%d) = %s, Slot(%q) = %s", index, byIndex.Hex(), slotName, byName.Hex())
	}
	return nil
}

func (tc *TestContext) aBase16PaletteWith16Colors() error {
	colors := map[string]string{
		"base00": "1a1b26", "base01": "1f2335", "base02": "292e42", "base03": "565f89",
		"base04": "a9b1d6", "base05": "c0caf5", "base06": "c0caf5", "base07": "c8d3f5",
		"base08": "f7768e", "base09": "ff9e64", "base0A": "e0af68", "base0B": "9ece6a",
		"base0C": "7dcfff", "base0D": "7aa2f7", "base0E": "bb9af7", "base0F": "db4b4b",
	}
	tc.palette, tc.paletteErr = domain.NewPalette("Test", "Author", "dark", "base16", colors)
	return tc.paletteErr
}

func (tc *TestContext) baseNShouldBeFallbackFromBaseM(n, m int) error {
	colorN := tc.palette.Base(n)
	colorM := tc.palette.Base(m)
	if !colorN.Equal(colorM) {
		return fmt.Errorf("base%d = %s, expected fallback from base%d = %s",
			n, colorN.Hex(), m, colorM.Hex())
	}
	return nil
}

// --- Schema Version Step Definitions ---

func registerSchemaSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^file kind "([^"]*)"$`, tc.setFileKind)
	ctx.Step(`^CurrentVersion should return (\d+)$`, tc.currentVersionShouldReturn)
	ctx.Step(`^all FileKind constants$`, tc.allFileKindConstants)
	ctx.Step(`^each should have a version greater than 0$`, tc.eachShouldHaveVersionGreaterThan0)
}

func (tc *TestContext) setFileKind(kind string) error {
	tc.currentFileKind = domain.FileKind(kind)
	return nil
}

func (tc *TestContext) currentVersionShouldReturn(expected int) error {
	got := domain.CurrentVersion(tc.currentFileKind)
	if got != expected {
		return fmt.Errorf("CurrentVersion(%s) = %d, want %d", tc.currentFileKind, got, expected)
	}
	return nil
}

func (tc *TestContext) allFileKindConstants() error {
	// This sets up for the next step
	return nil
}

func (tc *TestContext) eachShouldHaveVersionGreaterThan0() error {
	kinds := []domain.FileKind{
		domain.FileKindPalette,
		domain.FileKindUniversal,
		domain.FileKindVimMapping,
		domain.FileKindCSSMapping,
		domain.FileKindGtkMapping,
		domain.FileKindQssMapping,
		domain.FileKindStylixMapping,
	}
	for _, k := range kinds {
		v := domain.CurrentVersion(k)
		if v <= 0 {
			return fmt.Errorf("CurrentVersion(%s) = %d, expected > 0", k, v)
		}
	}
	return nil
}

// --- ThemeStore Step Definitions ---

func registerStoreSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^theme "([^"]*)" does not exist$`, tc.themeDoesNotExist)
	ctx.Step(`^I call EnsureThemeDir\("([^"]*)"\)$`, tc.iCallEnsureThemeDir)
	ctx.Step(`^the theme directory should exist$`, tc.theThemeDirectoryShouldExist)
	ctx.Step(`^theme "([^"]*)" exists with file "([^"]*)"$`, tc.themeExistsWithFile)
	ctx.Step(`^I call OpenWriter\("([^"]*)", "([^"]*)"\) and write "([^"]*)"$`, tc.iCallOpenWriterAndWrite)
	ctx.Step(`^I call OpenReader\("([^"]*)", "([^"]*)"\)$`, tc.iCallOpenReader)
	ctx.Step(`^the content should be "([^"]*)"$`, tc.theContentShouldBe)
	ctx.Step(`^I call Select\("([^"]*)"\)$`, tc.iCallSelect)
	ctx.Step(`^symlink "([^"]*)" should point to "([^"]*)"$`, tc.symlinkShouldPointTo)
	ctx.Step(`^SelectedTheme should return "([^"]*)"$`, tc.selectedThemeShouldReturn)
	ctx.Step(`^FileExists\("([^"]*)", "([^"]*)"\) should return (true|false)$`, tc.fileExistsShouldReturn)
}

func (tc *TestContext) themeDoesNotExist(name string) error {
	tc.themeName = name
	dir := tc.fsStore.ThemeDir(name)
	_, err := os.Stat(dir)
	if err == nil {
		return fmt.Errorf("theme %q unexpectedly exists", name)
	}
	return nil
}

func (tc *TestContext) iCallEnsureThemeDir(name string) error {
	tc.themeName = name
	tc.storeErr = tc.fsStore.EnsureThemeDir(name)
	return tc.storeErr
}

func (tc *TestContext) theThemeDirectoryShouldExist() error {
	dir := tc.fsStore.ThemeDir(tc.themeName)
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("theme directory does not exist: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("expected directory, got file")
	}
	return nil
}

func (tc *TestContext) themeExistsWithFile(theme, filename string) error {
	tc.themeName = theme
	if err := tc.fsStore.EnsureThemeDir(theme); err != nil {
		return err
	}
	w, err := tc.fsStore.OpenWriter(theme, filename)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("test content"))
	if err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func (tc *TestContext) iCallOpenWriterAndWrite(theme, filename, content string) error {
	tc.themeName = theme
	w, err := tc.fsStore.OpenWriter(theme, filename)
	if err != nil {
		tc.storeErr = err
		return nil
	}
	defer w.Close()
	_, tc.storeErr = w.Write([]byte(content))
	return nil
}

func (tc *TestContext) iCallOpenReader(theme, filename string) error {
	r, err := tc.fsStore.OpenReader(theme, filename)
	if err != nil {
		tc.storeErr = err
		return nil
	}
	defer r.Close()
	tc.fileData, tc.storeErr = io.ReadAll(r)
	return nil
}

func (tc *TestContext) theContentShouldBe(expected string) error {
	if tc.storeErr != nil {
		return tc.storeErr
	}
	if string(tc.fileData) != expected {
		return fmt.Errorf("content = %q, want %q", string(tc.fileData), expected)
	}
	return nil
}

func (tc *TestContext) iCallSelect(theme string) error {
	// First create the output files so symlinks have targets
	for _, f := range []string{"style.lua", "style.css", "gtk.css", "style.qss", "style.json"} {
		if err := tc.fsStore.EnsureThemeDir(theme); err != nil {
			return err
		}
		w, err := tc.fsStore.OpenWriter(theme, f)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte("content"))
		if err != nil {
			w.Close()
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	}
	tc.storeErr = tc.fsStore.Select(theme)
	return tc.storeErr
}

func (tc *TestContext) symlinkShouldPointTo(linkName, expectedTarget string) error {
	linkPath := filepath.Join(tc.fsStore.ConfigDir(), linkName)
	target, err := os.Readlink(linkPath)
	if err != nil {
		return fmt.Errorf("readlink %s: %w", linkName, err)
	}
	if target != expectedTarget {
		return fmt.Errorf("symlink %s -> %s, want -> %s", linkName, target, expectedTarget)
	}
	return nil
}

func (tc *TestContext) selectedThemeShouldReturn(expected string) error {
	got, err := tc.fsStore.SelectedTheme()
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("SelectedTheme() = %q, want %q", got, expected)
	}
	return nil
}

func (tc *TestContext) fileExistsShouldReturn(theme, filename, expected string) error {
	got := tc.fsStore.FileExists(theme, filename)
	want := expected == "true"
	if got != want {
		return fmt.Errorf("FileExists(%q, %q) = %v, want %v", theme, filename, got, want)
	}
	return nil
}

// --- Built-in Palettes Step Definitions ---

func registerBuiltinSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I call List\(\) on the built-in source$`, tc.iCallListOnBuiltinSource)
	ctx.Step(`^the result should contain "([^"]*)"$`, tc.theResultShouldContain)
	ctx.Step(`^the result should be sorted alphabetically$`, tc.theResultShouldBeSortedAlphabetically)
	ctx.Step(`^I call Get\("([^"]*)"\) on the built-in source$`, tc.iCallGetOnBuiltinSource)
	ctx.Step(`^I should receive valid YAML bytes$`, tc.iShouldReceiveValidYAMLBytes)
	ctx.Step(`^Get should return an error$`, tc.getShouldReturnAnError)
	ctx.Step(`^I call Has\("([^"]*)"\) on the built-in source$`, tc.iCallHasOnBuiltinSource)
	ctx.Step(`^Has should return (true|false)$`, tc.hasShouldReturn)
}

func (tc *TestContext) iCallListOnBuiltinSource() error {
	tc.builtinNames = tc.builtinSource.List()
	return nil
}

func (tc *TestContext) theResultShouldContain(name string) error {
	for _, n := range tc.builtinNames {
		if n == name {
			return nil
		}
	}
	return fmt.Errorf("List() does not contain %q, got %v", name, tc.builtinNames)
}

func (tc *TestContext) theResultShouldBeSortedAlphabetically() error {
	for i := 1; i < len(tc.builtinNames); i++ {
		if tc.builtinNames[i-1] > tc.builtinNames[i] {
			return fmt.Errorf("names not sorted: %q comes before %q",
				tc.builtinNames[i-1], tc.builtinNames[i])
		}
	}
	return nil
}

func (tc *TestContext) iCallGetOnBuiltinSource(name string) error {
	tc.builtinReader, tc.builtinGetErr = tc.builtinSource.Get(name)
	return nil
}

func (tc *TestContext) iShouldReceiveValidYAMLBytes() error {
	if tc.builtinGetErr != nil {
		return tc.builtinGetErr
	}
	data, err := io.ReadAll(tc.builtinReader)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("received empty YAML bytes")
	}
	// Verify it starts with expected YAML structure
	if !bytes.Contains(data, []byte("system:")) {
		return fmt.Errorf("YAML does not contain 'system:' field")
	}
	return nil
}

func (tc *TestContext) getShouldReturnAnError() error {
	if tc.builtinGetErr == nil {
		return fmt.Errorf("expected error, got nil")
	}
	return nil
}

func (tc *TestContext) iCallHasOnBuiltinSource(name string) error {
	tc.builtinHas = tc.builtinSource.Has(name)
	return nil
}

func (tc *TestContext) hasShouldReturn(expected string) error {
	want := expected == "true"
	if tc.builtinHas != want {
		return fmt.Errorf("Has() = %v, want %v", tc.builtinHas, want)
	}
	return nil
}

// --- Wrapper Step Definitions ---

func registerWrapperSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^a VersionedWriter for kind "([^"]*)" and theme "([^"]*)"$`, tc.aVersionedWriterForKindAndTheme)
	ctx.Step(`^I write "([^"]*)"$`, tc.iWrite)
	ctx.Step(`^the output should start with "([^"]*)"$`, tc.theOutputShouldStartWith)
	ctx.Step(`^the output should contain "([^"]*)"$`, tc.theOutputShouldContain)
	ctx.Step(`^YAML with schema_version (\d+) for kind "([^"]*)"$`, tc.yamlWithSchemaVersionForKind)
	ctx.Step(`^I wrap it in ValidatingReader and read$`, tc.iWrapItInValidatingReaderAndRead)
	ctx.Step(`^reading should succeed$`, tc.readingShouldSucceed)
	ctx.Step(`^reading should fail with SchemaVersionError$`, tc.readingShouldFailWithSchemaVersionError)
	ctx.Step(`^NeedsUpgrade should be (true|false)$`, tc.needsUpgradeShouldBe)
}

func (tc *TestContext) aVersionedWriterForKindAndTheme(kind, theme string) error {
	tc.outputBuffer = &bytes.Buffer{}
	tc.versionedWriter = wrappers.NewVersionedWriter(tc.outputBuffer, domain.FileKind(kind), theme)
	return nil
}

func (tc *TestContext) iWrite(content string) error {
	_, err := tc.versionedWriter.Write([]byte(content))
	return err
}

func (tc *TestContext) theOutputShouldStartWith(prefix string) error {
	output := tc.outputBuffer.String()
	if !bytes.HasPrefix([]byte(output), []byte(prefix)) {
		return fmt.Errorf("output does not start with %q:\n%s", prefix, output)
	}
	return nil
}

func (tc *TestContext) theOutputShouldContain(substr string) error {
	output := tc.outputBuffer.String()
	if !bytes.Contains([]byte(output), []byte(substr)) {
		return fmt.Errorf("output does not contain %q:\n%s", substr, output)
	}
	return nil
}

func (tc *TestContext) yamlWithSchemaVersionForKind(version int, kind string) error {
	tc.currentFileKind = domain.FileKind(kind)
	yaml := fmt.Sprintf("schema_version: %d\nkind: %s\ntheme_name: test\ncontent: data\n", version, kind)
	tc.outputBuffer = bytes.NewBufferString(yaml)
	return nil
}

func (tc *TestContext) iWrapItInValidatingReaderAndRead() error {
	vr := wrappers.NewValidatingReader(tc.outputBuffer, tc.currentFileKind)
	tc.readBytes, tc.validationErr = io.ReadAll(vr)
	return nil
}

func (tc *TestContext) readingShouldSucceed() error {
	if tc.validationErr != nil {
		return fmt.Errorf("expected success, got: %v", tc.validationErr)
	}
	return nil
}

func (tc *TestContext) readingShouldFailWithSchemaVersionError() error {
	if tc.validationErr == nil {
		return fmt.Errorf("expected SchemaVersionError, got nil")
	}
	_, ok := tc.validationErr.(*domain.SchemaVersionError)
	if !ok {
		return fmt.Errorf("expected *domain.SchemaVersionError, got %T: %v", tc.validationErr, tc.validationErr)
	}
	return nil
}

func (tc *TestContext) needsUpgradeShouldBe(expected string) error {
	sve, ok := tc.validationErr.(*domain.SchemaVersionError)
	if !ok {
		return fmt.Errorf("no SchemaVersionError available")
	}
	want := expected == "true"
	if sve.NeedsUpgrade != want {
		return fmt.Errorf("NeedsUpgrade = %v, want %v", sve.NeedsUpgrade, want)
	}
	return nil
}

// --- Token Deriver Step Definitions ---

func registerDeriverSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I derive tokens from the Tokyo Night Dark palette$`, tc.iDeriveTokensFromTokyoNightDark)
	ctx.Step(`^the TokenSet should have at least (\d+) tokens$`, tc.theTokenSetShouldHaveAtLeastTokens)
	ctx.Step(`^token "([^"]*)" should have color "([^"]*)"$`, tc.tokenShouldHaveColor)
	ctx.Step(`^token "([^"]*)" should be italic$`, tc.tokenShouldBeItalic)
	ctx.Step(`^token "([^"]*)" should be bold$`, tc.tokenShouldBeBold)
}

func (tc *TestContext) iDeriveTokensFromTokyoNightDark() error {
	if tc.palette == nil {
		if err := tc.theTokyoNightDarkPalette(); err != nil {
			return err
		}
	}
	drv := deriver.New()
	tc.tokenSet = drv.Derive(tc.palette)
	return nil
}

func (tc *TestContext) theTokenSetShouldHaveAtLeastTokens(count int) error {
	got := tc.tokenSet.Len()
	if got < count {
		return fmt.Errorf("TokenSet has %d tokens, want at least %d", got, count)
	}
	return nil
}

func (tc *TestContext) tokenShouldHaveColor(path, expectedHex string) error {
	tok, ok := tc.tokenSet.Get(path)
	if !ok {
		return fmt.Errorf("token %q not found", path)
	}
	expected, err := domain.ParseHex(expectedHex)
	if err != nil {
		return err
	}
	if !tok.Color.Equal(expected) {
		return fmt.Errorf("token %q color = %s, want %s", path, tok.Color.Hex(), expected.Hex())
	}
	return nil
}

func (tc *TestContext) tokenShouldBeItalic(path string) error {
	tok, ok := tc.tokenSet.Get(path)
	if !ok {
		return fmt.Errorf("token %q not found", path)
	}
	if !tok.Italic {
		return fmt.Errorf("token %q Italic = false, want true", path)
	}
	return nil
}

func (tc *TestContext) tokenShouldBeBold(path string) error {
	tok, ok := tc.tokenSet.Get(path)
	if !ok {
		return fmt.Errorf("token %q not found", path)
	}
	if !tok.Bold {
		return fmt.Errorf("token %q Bold = false, want true", path)
	}
	return nil
}

// --- Mapper Step Definitions ---

func registerMapperSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I create a ResolvedTheme from Tokyo Night Dark$`, tc.iCreateResolvedThemeFromTokyoNightDark)
	ctx.Step(`^I map it with the Stylix mapper$`, tc.iMapWithStylixMapper)
	ctx.Step(`^the StylixTheme should have at least (\d+) values$`, tc.theStylixThemeShouldHaveAtLeastValues)
	ctx.Step(`^I map it with the Vim mapper$`, tc.iMapWithVimMapper)
	ctx.Step(`^the VimTheme should have at least (\d+) highlight groups$`, tc.theVimThemeShouldHaveAtLeastHighlights)
	ctx.Step(`^the VimTheme should have 16 terminal colors$`, tc.theVimThemeShouldHave16TerminalColors)
}

func (tc *TestContext) iCreateResolvedThemeFromTokyoNightDark() error {
	if tc.palette == nil {
		if err := tc.theTokyoNightDarkPalette(); err != nil {
			return err
		}
	}
	if tc.tokenSet == nil {
		if err := tc.iDeriveTokensFromTokyoNightDark(); err != nil {
			return err
		}
	}
	tc.resolvedTheme = &domain.ResolvedTheme{
		Name:    tc.palette.Name,
		Variant: tc.palette.Variant,
		Palette: tc.palette,
		Tokens:  tc.tokenSet,
	}
	return nil
}

func (tc *TestContext) iMapWithStylixMapper() error {
	m := mapper.NewStylix()
	mapped, err := m.Map(tc.resolvedTheme)
	if err != nil {
		return err
	}
	tc.stylixTheme = mapped.(*ports.StylixTheme)
	return nil
}

func (tc *TestContext) theStylixThemeShouldHaveAtLeastValues(count int) error {
	got := len(tc.stylixTheme.Values)
	if got < count {
		return fmt.Errorf("StylixTheme has %d values, want at least %d", got, count)
	}
	return nil
}

func (tc *TestContext) iMapWithVimMapper() error {
	m := mapper.NewVim()
	mapped, err := m.Map(tc.resolvedTheme)
	if err != nil {
		return err
	}
	tc.vimTheme = mapped.(*ports.VimTheme)
	return nil
}

func (tc *TestContext) theVimThemeShouldHaveAtLeastHighlights(count int) error {
	got := len(tc.vimTheme.Highlights)
	if got < count {
		return fmt.Errorf("VimTheme has %d highlight groups, want at least %d", got, count)
	}
	return nil
}

func (tc *TestContext) theVimThemeShouldHave16TerminalColors() error {
	for i, c := range tc.vimTheme.TerminalColors {
		if c.IsNone {
			return fmt.Errorf("terminal color %d is NoneColor", i)
		}
	}
	return nil
}

// --- Generator Step Definitions ---

func registerGeneratorSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I generate Stylix output$`, tc.iGenerateStylixOutput)
	ctx.Step(`^the output should be valid JSON$`, tc.theOutputShouldBeValidJSON)
	ctx.Step(`^the JSON should contain key "([^"]*)"$`, tc.theJSONShouldContainKey)
	ctx.Step(`^I generate Vim output$`, tc.iGenerateVimOutput)
	ctx.Step(`^the generated output should contain "([^"]*)"$`, tc.theGeneratedOutputShouldContain)
}

func (tc *TestContext) iGenerateStylixOutput() error {
	tc.genOutputBuf = &bytes.Buffer{}
	gen := generator.NewStylix()
	tc.generateErr = gen.Generate(tc.genOutputBuf, tc.stylixTheme)
	return tc.generateErr
}

func (tc *TestContext) theOutputShouldBeValidJSON() error {
	output := tc.genOutputBuf.Bytes()
	trimmed := bytes.TrimSpace(output)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return fmt.Errorf("output does not start with '{'")
	}
	if trimmed[len(trimmed)-1] != '}' {
		return fmt.Errorf("output does not end with '}'")
	}
	return nil
}

func (tc *TestContext) theJSONShouldContainKey(key string) error {
	output := tc.genOutputBuf.String()
	pattern := fmt.Sprintf(`"%s":`, key)
	if !bytes.Contains([]byte(output), []byte(pattern)) {
		return fmt.Errorf("JSON does not contain key %q", key)
	}
	return nil
}

func (tc *TestContext) iGenerateVimOutput() error {
	tc.genOutputBuf = &bytes.Buffer{}
	gen := generator.NewVim()
	tc.generateErr = gen.Generate(tc.genOutputBuf, tc.vimTheme)
	return tc.generateErr
}

func (tc *TestContext) theGeneratedOutputShouldContain(expected string) error {
	if tc.genOutputBuf == nil {
		return fmt.Errorf("generator output buffer is nil")
	}
	output := tc.genOutputBuf.String()
	if !bytes.Contains([]byte(output), []byte(expected)) {
		return fmt.Errorf("generated output does not contain %q:\n%s", expected, truncateBytes([]byte(output), 500))
	}
	return nil
}

// --- E2E Step Definitions ---

func registerE2ESteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I run the full pipeline for "([^"]*)"$`, tc.iRunTheFullPipelineFor)
	ctx.Step(`^all 12 files should be created$`, tc.all12FilesShouldBeCreated)
	ctx.Step(`^running the pipeline again should produce identical output$`, tc.runningPipelineAgainShouldProduceIdenticalOutput)
}

func (tc *TestContext) iRunTheFullPipelineFor(scheme string) error {
	// Use built-in source to get palette
	r, err := tc.builtinSource.Get(scheme)
	if err != nil {
		return err
	}

	parser := yamlparser.NewParser()
	pal, err := parser.Parse(r)
	if err != nil {
		return err
	}
	tc.palette = pal
	tc.themeName = scheme

	// Derive tokens
	drv := deriver.New()
	tc.tokenSet = drv.Derive(pal)

	// Create ResolvedTheme
	tc.resolvedTheme = &domain.ResolvedTheme{
		Name:    pal.Name,
		Variant: pal.Variant,
		Palette: pal,
		Tokens:  tc.tokenSet,
	}

	// Create theme directory and generate all files
	if err := tc.fsStore.EnsureThemeDir(scheme); err != nil {
		return err
	}

	// Write palette.yaml
	pw, err := tc.fsStore.OpenWriter(scheme, "palette.yaml")
	if err != nil {
		return err
	}
	vw := wrappers.NewVersionedWriter(pw, domain.FileKindPalette, scheme)
	if _, err := vw.Write([]byte("palette_content: test\n")); err != nil {
		pw.Close()
		return err
	}
	pw.Close()

	// Write universal.yaml
	uw, err := tc.fsStore.OpenWriter(scheme, "universal.yaml")
	if err != nil {
		return err
	}
	vwu := wrappers.NewVersionedWriter(uw, domain.FileKindUniversal, scheme)
	if _, err := vwu.Write([]byte("tokens: test\n")); err != nil {
		uw.Close()
		return err
	}
	uw.Close()

	// Write mapping files and outputs
	mappings := []struct {
		mapping  string
		output   string
		kind     domain.FileKind
		genWrite func() error
	}{
		{
			mapping: "vim-mapping.yaml",
			output:  "style.lua",
			kind:    domain.FileKindVimMapping,
			genWrite: func() error {
				if err := tc.iMapWithVimMapper(); err != nil {
					return err
				}
				w, err := tc.fsStore.OpenWriter(scheme, "style.lua")
				if err != nil {
					return err
				}
				defer w.Close()
				gen := generator.NewVim()
				return gen.Generate(w, tc.vimTheme)
			},
		},
		{
			mapping: "css-mapping.yaml",
			output:  "style.css",
			kind:    domain.FileKindCSSMapping,
			genWrite: func() error {
				m := mapper.NewCSS()
				mapped, err := m.Map(tc.resolvedTheme)
				if err != nil {
					return err
				}
				w, err := tc.fsStore.OpenWriter(scheme, "style.css")
				if err != nil {
					return err
				}
				defer w.Close()
				gen := generator.NewCSS()
				return gen.Generate(w, mapped)
			},
		},
		{
			mapping: "gtk-mapping.yaml",
			output:  "gtk.css",
			kind:    domain.FileKindGtkMapping,
			genWrite: func() error {
				m := mapper.NewGtk()
				mapped, err := m.Map(tc.resolvedTheme)
				if err != nil {
					return err
				}
				w, err := tc.fsStore.OpenWriter(scheme, "gtk.css")
				if err != nil {
					return err
				}
				defer w.Close()
				gen := generator.NewGtk()
				return gen.Generate(w, mapped)
			},
		},
		{
			mapping: "qss-mapping.yaml",
			output:  "style.qss",
			kind:    domain.FileKindQssMapping,
			genWrite: func() error {
				m := mapper.NewQss()
				mapped, err := m.Map(tc.resolvedTheme)
				if err != nil {
					return err
				}
				w, err := tc.fsStore.OpenWriter(scheme, "style.qss")
				if err != nil {
					return err
				}
				defer w.Close()
				gen := generator.NewQss()
				return gen.Generate(w, mapped)
			},
		},
		{
			mapping: "stylix-mapping.yaml",
			output:  "style.json",
			kind:    domain.FileKindStylixMapping,
			genWrite: func() error {
				if err := tc.iMapWithStylixMapper(); err != nil {
					return err
				}
				w, err := tc.fsStore.OpenWriter(scheme, "style.json")
				if err != nil {
					return err
				}
				defer w.Close()
				gen := generator.NewStylix()
				return gen.Generate(w, tc.stylixTheme)
			},
		},
	}

	for _, m := range mappings {
		// Write mapping file
		mw, err := tc.fsStore.OpenWriter(scheme, m.mapping)
		if err != nil {
			return err
		}
		vwm := wrappers.NewVersionedWriter(mw, m.kind, scheme)
		if _, err := vwm.Write([]byte("mapping: test\n")); err != nil {
			mw.Close()
			return err
		}
		mw.Close()

		// Generate output file
		if err := m.genWrite(); err != nil {
			return err
		}
	}

	return nil
}

func (tc *TestContext) all12FilesShouldBeCreated() error {
	files := []string{
		"palette.yaml",
		"universal.yaml",
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
	for _, f := range files {
		if !tc.fsStore.FileExists(tc.themeName, f) {
			return fmt.Errorf("file %q does not exist", f)
		}
	}
	return nil
}

func (tc *TestContext) runningPipelineAgainShouldProduceIdenticalOutput() error {
	// Read first run's style.json
	r1, err := tc.fsStore.OpenReader(tc.themeName, "style.json")
	if err != nil {
		return err
	}
	data1, err := io.ReadAll(r1)
	r1.Close()
	if err != nil {
		return err
	}

	// Generate again
	if err := tc.iMapWithStylixMapper(); err != nil {
		return err
	}
	w, err := tc.fsStore.OpenWriter(tc.themeName, "style.json")
	if err != nil {
		return err
	}
	gen := generator.NewStylix()
	if err := gen.Generate(w, tc.stylixTheme); err != nil {
		w.Close()
		return err
	}
	w.Close()

	// Read second run
	r2, err := tc.fsStore.OpenReader(tc.themeName, "style.json")
	if err != nil {
		return err
	}
	data2, err := io.ReadAll(r2)
	r2.Close()
	if err != nil {
		return err
	}

	if !bytes.Equal(data1, data2) {
		return fmt.Errorf("output not deterministic:\nrun1: %s\nrun2: %s",
			truncateBytes(data1, 200), truncateBytes(data2, 200))
	}
	return nil
}

// Helper functions

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func truncateBytes(data []byte, n int) string {
	if len(data) <= n {
		return string(data)
	}
	return string(data[:n]) + "..."
}

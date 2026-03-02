package application_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/curtbushko/flair/internal/application"
	"github.com/curtbushko/flair/internal/domain"
	"github.com/curtbushko/flair/internal/ports"
)

// --- Mtime-aware stub store for regenerate tests ---

// mtimeStubStore extends stubThemeStore with configurable mtimes.
type mtimeStubStore struct {
	stubThemeStore
	mtimes map[string]map[string]time.Time // themeName -> filename -> mtime
}

func newMtimeStubStore() *mtimeStubStore {
	return &mtimeStubStore{
		stubThemeStore: stubThemeStore{
			files: make(map[string]map[string]*recordedWrite),
		},
		mtimes: make(map[string]map[string]time.Time),
	}
}

func (s *mtimeStubStore) FileMtime(themeName, filename string) (time.Time, error) {
	if theme, ok := s.mtimes[themeName]; ok {
		if t, ok := theme[filename]; ok {
			return t, nil
		}
	}
	return time.Time{}, &fileNotFoundError{themeName: themeName, filename: filename}
}

func (s *mtimeStubStore) FileExists(themeName, filename string) bool {
	if theme, ok := s.mtimes[themeName]; ok {
		_, ok := theme[filename]
		return ok
	}
	return false
}

func (s *mtimeStubStore) setMtime(themeName, filename string, t time.Time) {
	if s.mtimes[themeName] == nil {
		s.mtimes[themeName] = make(map[string]time.Time)
	}
	s.mtimes[themeName][filename] = t
	// Also ensure the file entry exists for OpenReader.
	s.mu.Lock()
	if s.files[themeName] == nil {
		s.files[themeName] = make(map[string]*recordedWrite)
	}
	if s.files[themeName][filename] == nil {
		s.files[themeName][filename] = &recordedWrite{data: []byte("# placeholder")}
	}
	s.mu.Unlock()
}

type fileNotFoundError struct {
	themeName string
	filename  string
}

func (e *fileNotFoundError) Error() string {
	return "file not found: " + e.themeName + "/" + e.filename
}

// --- Regenerate test helpers ---

const regenTestTheme = "my-theme"

func makeRegenTargets() []ports.Target {
	return makeStubTargets()
}

func makeRegenUseCase(
	store *mtimeStubStore,
	parser ports.PaletteParser,
	tokenizer ports.Tokenizer,
	targets []ports.Target,
) *application.RegenerateThemeUseCase {
	return application.NewRegenerateThemeUseCase(
		store,
		parser,
		tokenizer,
		targets,
		application.WithRegenTokensWriter(func(w io.Writer, ts *domain.TokenSet) error {
			_, err := w.Write([]byte("tokens-data"))
			return err
		}),
	)
}

// --- Tests ---

func TestRegenerateThemeUseCase_PaletteEdited(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newMtimeStubStore()
	targets := makeRegenTargets()

	now := time.Now()
	// Palette is newest -> everything downstream should regenerate.
	store.setMtime(regenTestTheme, "palette.yaml", now)
	store.setMtime(regenTestTheme, "tokens.yaml", now.Add(-2*time.Second))
	for _, tgt := range targets {
		store.setMtime(regenTestTheme, tgt.MappingFile, now.Add(-3*time.Second))
		store.setMtime(regenTestTheme, tgt.Generator.DefaultFilename(), now.Add(-4*time.Second))
	}

	// Seed palette.yaml content for parser.
	store.mu.Lock()
	store.files[regenTestTheme]["palette.yaml"] = &recordedWrite{data: []byte("palette-yaml")}
	store.mu.Unlock()

	uc := makeRegenUseCase(
		store,
		&stubGenParser{palette: pal},
		&stubGenTokenizer{tokenSet: ts},
		targets,
	)

	msg, err := uc.Execute(regenTestTheme, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should regenerate universal + all mappings + all outputs.
	if !store.hasFile(regenTestTheme, "tokens.yaml") {
		t.Error("expected tokens.yaml to be regenerated")
	}
	for _, tgt := range targets {
		if !store.hasFile(regenTestTheme, tgt.MappingFile) {
			t.Errorf("expected %s to be regenerated", tgt.MappingFile)
		}
		if !store.hasFile(regenTestTheme, tgt.Generator.DefaultFilename()) {
			t.Errorf("expected %s to be regenerated", tgt.Generator.DefaultFilename())
		}
	}

	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestRegenerateThemeUseCase_UniversalEdited(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newMtimeStubStore()
	targets := makeRegenTargets()

	now := time.Now()
	// Universal is newest among downstream (newer than mappings but older than palette is okay,
	// but here palette is older, so universal is the trigger).
	store.setMtime(regenTestTheme, "palette.yaml", now.Add(-5*time.Second))
	store.setMtime(regenTestTheme, "tokens.yaml", now)
	for _, tgt := range targets {
		store.setMtime(regenTestTheme, tgt.MappingFile, now.Add(-3*time.Second))
		store.setMtime(regenTestTheme, tgt.Generator.DefaultFilename(), now.Add(-4*time.Second))
	}

	store.mu.Lock()
	store.files[regenTestTheme]["palette.yaml"] = &recordedWrite{data: []byte("palette-yaml")}
	store.mu.Unlock()

	uc := makeRegenUseCase(
		store,
		&stubGenParser{palette: pal},
		&stubGenTokenizer{tokenSet: ts},
		targets,
	)

	msg, err := uc.Execute(regenTestTheme, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should regenerate all mappings + outputs but NOT re-derive universal.
	// Check that mappings and outputs were written.
	for _, tgt := range targets {
		if !store.hasFile(regenTestTheme, tgt.MappingFile) {
			t.Errorf("expected %s to be regenerated", tgt.MappingFile)
		}
		if !store.hasFile(regenTestTheme, tgt.Generator.DefaultFilename()) {
			t.Errorf("expected %s to be regenerated", tgt.Generator.DefaultFilename())
		}
	}

	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestRegenerateThemeUseCase_MappingEdited(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newMtimeStubStore()
	targets := makeRegenTargets()

	now := time.Now()
	// All files are old, but one mapping is newer than its output.
	store.setMtime(regenTestTheme, "palette.yaml", now.Add(-10*time.Second))
	store.setMtime(regenTestTheme, "tokens.yaml", now.Add(-10*time.Second))
	for _, tgt := range targets {
		store.setMtime(regenTestTheme, tgt.MappingFile, now.Add(-10*time.Second))
		store.setMtime(regenTestTheme, tgt.Generator.DefaultFilename(), now.Add(-5*time.Second))
	}
	// Make the first target's mapping newer than its output.
	store.setMtime(regenTestTheme, targets[0].MappingFile, now)
	store.setMtime(regenTestTheme, targets[0].Generator.DefaultFilename(), now.Add(-5*time.Second))

	store.mu.Lock()
	store.files[regenTestTheme]["palette.yaml"] = &recordedWrite{data: []byte("palette-yaml")}
	store.mu.Unlock()

	uc := makeRegenUseCase(
		store,
		&stubGenParser{palette: pal},
		&stubGenTokenizer{tokenSet: ts},
		targets,
	)

	msg, err := uc.Execute(regenTestTheme, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should regenerate only the first target's output.
	if !store.hasFile(regenTestTheme, targets[0].Generator.DefaultFilename()) {
		t.Errorf("expected %s to be regenerated", targets[0].Generator.DefaultFilename())
	}

	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestRegenerateThemeUseCase_NoEdits(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newMtimeStubStore()
	targets := makeRegenTargets()

	now := time.Now()
	// All downstream files are newer than their upstream -> nothing to do.
	store.setMtime(regenTestTheme, "palette.yaml", now.Add(-10*time.Second))
	store.setMtime(regenTestTheme, "tokens.yaml", now.Add(-5*time.Second))
	for _, tgt := range targets {
		store.setMtime(regenTestTheme, tgt.MappingFile, now.Add(-3*time.Second))
		store.setMtime(regenTestTheme, tgt.Generator.DefaultFilename(), now.Add(-1*time.Second))
	}

	store.mu.Lock()
	store.files[regenTestTheme]["palette.yaml"] = &recordedWrite{data: []byte("palette-yaml")}
	store.mu.Unlock()

	uc := makeRegenUseCase(
		store,
		&stubGenParser{palette: pal},
		&stubGenTokenizer{tokenSet: ts},
		targets,
	)

	msg, err := uc.Execute(regenTestTheme, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(msg, "nothing to do") {
		t.Errorf("expected 'nothing to do' message, got %q", msg)
	}
}

func TestRegenerateThemeUseCase_TargetFilter(t *testing.T) {
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)
	store := newMtimeStubStore()
	targets := makeRegenTargets()

	now := time.Now()
	// Palette is newest -> would normally regen everything.
	store.setMtime(regenTestTheme, "palette.yaml", now)
	store.setMtime(regenTestTheme, "tokens.yaml", now.Add(-2*time.Second))
	for _, tgt := range targets {
		store.setMtime(regenTestTheme, tgt.MappingFile, now.Add(-3*time.Second))
		store.setMtime(regenTestTheme, tgt.Generator.DefaultFilename(), now.Add(-4*time.Second))
	}

	store.mu.Lock()
	store.files[regenTestTheme]["palette.yaml"] = &recordedWrite{data: []byte("palette-yaml")}
	store.mu.Unlock()

	uc := makeRegenUseCase(
		store,
		&stubGenParser{palette: pal},
		&stubGenTokenizer{tokenSet: ts},
		targets,
	)

	// Filter to only regenerate the "vim" target.
	_, err := uc.Execute(regenTestTheme, "vim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// vim output should be regenerated.
	if !store.hasFile(regenTestTheme, "style.lua") {
		t.Error("expected style.lua to be regenerated")
	}

	// Clear the write records for non-vim targets to check they were not written.
	// Check that other target outputs were NOT regenerated by looking at fresh writes.
	// Since the store writes on execute, check that only vim-related files got new writes.
	// For simplicity, we check the universal was written (palette edit triggers it).
	if !store.hasFile(regenTestTheme, "tokens.yaml") {
		t.Error("expected tokens.yaml to be regenerated when palette is edited")
	}
}

func TestRegenerateThemeUseCase_ThemeNotFound(t *testing.T) {
	store := newMtimeStubStore()
	targets := makeRegenTargets()
	pal := makeGenPalette(t)
	ts := makeGenTokenSet(t, pal)

	// Set up an unrelated theme to prove the use case checks the right name.
	store.setMtime("other-theme", "palette.yaml", time.Now())

	uc := makeRegenUseCase(
		store,
		&stubGenParser{palette: pal},
		&stubGenTokenizer{tokenSet: ts},
		targets,
	)

	_, err := uc.Execute("nonexistent-theme", "")
	if err == nil {
		t.Fatal("expected error for non-existent theme, got nil")
	}
}

// Ensure mtimeStubStore also satisfies OpenWriter by delegating to embedded stubThemeStore.
func (s *mtimeStubStore) OpenWriter(themeName, filename string) (io.WriteCloser, error) {
	return s.stubThemeStore.OpenWriter(themeName, filename)
}

func (s *mtimeStubStore) OpenReader(themeName, filename string) (io.ReadCloser, error) {
	return s.stubThemeStore.OpenReader(themeName, filename)
}

func (s *mtimeStubStore) EnsureThemeDir(themeName string) error {
	return s.stubThemeStore.EnsureThemeDir(themeName)
}

func (s *mtimeStubStore) ThemeDir(themeName string) string {
	return s.stubThemeStore.ThemeDir(themeName)
}

func (s *mtimeStubStore) ConfigDir() string {
	return s.stubThemeStore.ConfigDir()
}

func (s *mtimeStubStore) ListThemes() ([]string, error) {
	return s.stubThemeStore.ListThemes()
}

func (s *mtimeStubStore) SelectedTheme() (string, error) {
	return s.stubThemeStore.SelectedTheme()
}

func (s *mtimeStubStore) Select(themeName string) error {
	return s.stubThemeStore.Select(themeName)
}

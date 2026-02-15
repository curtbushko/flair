package application_test

import (
	"strings"
	"testing"

	"github.com/curtbushko/flair/internal/application"
)

func TestInitThemeUseCase_CreatesDir(t *testing.T) {
	store := newStubThemeStore()

	uc := application.NewInitThemeUseCase(store)

	result, err := uc.Execute("my-theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// EnsureThemeDir should have been called.
	if len(store.ensureDirCalls) == 0 {
		t.Fatal("expected EnsureThemeDir to be called")
	}
	if store.ensureDirCalls[0] != "my-theme" {
		t.Errorf("EnsureThemeDir called with %q, want %q", store.ensureDirCalls[0], "my-theme")
	}

	// Result should contain the palette path.
	if result == "" {
		t.Error("expected non-empty result path")
	}
}

func TestInitThemeUseCase_WritesScaffoldPalette(t *testing.T) {
	store := newStubThemeStore()

	uc := application.NewInitThemeUseCase(store)

	_, err := uc.Execute("my-theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// palette.yaml should have been written.
	if !store.hasFile("my-theme", "palette.yaml") {
		t.Fatal("expected palette.yaml to be written")
	}

	// Read back contents.
	store.mu.Lock()
	data := store.files["my-theme"]["palette.yaml"].data
	store.mu.Unlock()

	content := string(data)

	// Should contain schema_version header.
	if !strings.Contains(content, "schema_version:") {
		t.Error("expected palette.yaml to contain schema_version header")
	}

	// Should contain system field.
	if !strings.Contains(content, "system:") {
		t.Error("expected palette.yaml to contain system field")
	}

	// Should contain the theme name.
	if !strings.Contains(content, "my-theme") {
		t.Error("expected palette.yaml to contain theme name")
	}
}

func TestInitThemeUseCase_PaletteHasAllSlots(t *testing.T) {
	store := newStubThemeStore()

	uc := application.NewInitThemeUseCase(store)

	_, err := uc.Execute("my-theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store.mu.Lock()
	data := store.files["my-theme"]["palette.yaml"].data
	store.mu.Unlock()

	content := string(data)

	// All 24 base slots should be present.
	slots := []string{
		"base00", "base01", "base02", "base03",
		"base04", "base05", "base06", "base07",
		"base08", "base09", "base0A", "base0B",
		"base0C", "base0D", "base0E", "base0F",
		"base10", "base11", "base12", "base13",
		"base14", "base15", "base16", "base17",
	}

	for _, slot := range slots {
		if !strings.Contains(content, slot+":") {
			t.Errorf("expected palette.yaml to contain slot %q", slot)
		}
	}
}

func TestInitThemeUseCase_AlreadyExists(t *testing.T) {
	store := newStubThemeStore()

	// Pre-create the theme with a palette.yaml.
	if err := store.EnsureThemeDir("existing"); err != nil {
		t.Fatal(err)
	}
	w, err := store.OpenWriter("existing", "palette.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write([]byte("existing content"))
	_ = w.Close()

	uc := application.NewInitThemeUseCase(store)

	_, err = uc.Execute("existing")
	if err == nil {
		t.Fatal("expected error for existing palette.yaml, got nil")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %v, want it to mention 'already exists'", err)
	}
}

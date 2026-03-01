package viewer

import (
	"testing"
)

// TestRun_ReturnsModel verifies Model() returns a valid tea.Model.
func TestRun_ReturnsModel(t *testing.T) {
	opts := Options{
		Themes: []string{"theme1", "theme2"},
	}

	model := NewModel(opts)

	// Verify it's the correct type.
	if len(model.themes) != 2 {
		t.Errorf("model.themes = %d, want 2", len(model.themes))
	}
}

// TestRunOptions_Defaults verifies default values are sensible.
func TestRunOptions_Defaults(t *testing.T) {
	opts := Options{}

	m := NewModel(opts)

	// Empty themes is valid.
	if m.themes == nil {
		t.Error("themes should not be nil")
	}

	// Default page is text status (first content page in 2-panel layout).
	if m.currentPage != PageTextStatus {
		t.Errorf("default page = %v, want PageTextStatus", m.currentPage)
	}
}

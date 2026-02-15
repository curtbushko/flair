package main

import (
	"testing"
)

func TestWireTargets(t *testing.T) {
	app := Wire("/tmp/flair-test")

	// Verify exactly 5 targets.
	if got := len(app.Targets); got != 5 {
		t.Fatalf("len(Targets) = %d, want 5", got)
	}

	// Expected target names and mapping files.
	type want struct {
		name        string
		mappingFile string
	}

	expected := []want{
		{"vim", "vim-mapping.yaml"},
		{"css", "css-mapping.yaml"},
		{"gtk", "gtk-mapping.yaml"},
		{"qss", "qss-mapping.yaml"},
		{"stylix", "stylix-mapping.yaml"},
	}

	// Build a lookup from actual targets.
	found := make(map[string]string, len(app.Targets))
	for _, tgt := range app.Targets {
		found[tgt.Mapper.Name()] = tgt.MappingFile
	}

	for _, exp := range expected {
		t.Run(exp.name, func(t *testing.T) {
			mf, ok := found[exp.name]
			if !ok {
				t.Fatalf("target %q not found in Targets slice", exp.name)
			}

			if mf != exp.mappingFile {
				t.Errorf("target %q MappingFile = %q, want %q", exp.name, mf, exp.mappingFile)
			}
		})
	}

	// Verify each target has a WriteMappingFile function wired.
	for _, tgt := range app.Targets {
		if tgt.WriteMappingFile == nil {
			t.Errorf("target %q has nil WriteMappingFile", tgt.Mapper.Name())
		}
	}

	// Verify each target has a non-nil Generator.
	for _, tgt := range app.Targets {
		if tgt.Generator == nil {
			t.Errorf("target %q has nil Generator", tgt.Mapper.Name())
		}
	}
}

package main

import (
	"testing"
)

// generateThemeForTest generates a built-in theme in the given directory
// using the internal GenerateThemeUseCase. This is used by tests that need
// a theme to exist before running other commands.
func generateThemeForTest(t *testing.T, dir, themeName string) {
	t.Helper()

	app := Wire(dir)
	if err := app.Generate.ExecuteBuiltin(themeName, themeName, ""); err != nil {
		t.Fatalf("generateThemeForTest(%s): %v", themeName, err)
	}
}

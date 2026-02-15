package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/curtbushko/flair/internal/config"
)

func TestDefaultConfigDir_XDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/custom-xdg")

	got := config.DefaultConfigDir()
	want := filepath.Join("/tmp/custom-xdg", "flair")

	if got != want {
		t.Errorf("DefaultConfigDir() = %q, want %q", got, want)
	}
}

func TestDefaultConfigDir_Fallback(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error: %v", err)
	}

	got := config.DefaultConfigDir()
	want := filepath.Join(home, ".config", "flair")

	if got != want {
		t.Errorf("DefaultConfigDir() = %q, want %q", got, want)
	}
}

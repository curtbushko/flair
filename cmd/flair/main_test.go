package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLI_UnknownCommand(t *testing.T) {
	var stderr bytes.Buffer
	code := run([]string{"flair", "bogus"}, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	if !strings.Contains(stderr.String(), "unknown command") {
		t.Errorf("stderr = %q, want it to contain %q", stderr.String(), "unknown command")
	}
}

func TestCLI_NoArgs(t *testing.T) {
	var stderr bytes.Buffer
	code := run([]string{"flair"}, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}

	output := stderr.String()
	if !strings.Contains(output, "Usage:") {
		t.Errorf("stderr = %q, want it to contain %q", output, "Usage:")
	}
}

func TestCLI_HelpFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"--help", []string{"flair", "--help"}},
		{"-h", []string{"flair", "-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			code := run(tt.args, &stderr)

			if code != 0 {
				t.Errorf("exit code = %d, want 0", code)
			}

			output := stderr.String()
			if !strings.Contains(output, "Usage:") {
				t.Errorf("stderr = %q, want it to contain %q", output, "Usage:")
			}
		})
	}
}

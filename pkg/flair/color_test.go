package flair_test

import (
	"testing"

	"github.com/curtbushko/flair/pkg/flair"
)

func TestColor_Hex(t *testing.T) {
	tests := []struct {
		name  string
		color flair.Color
		want  string
	}{
		{
			name:  "tokyo night background",
			color: flair.Color{R: 26, G: 27, B: 38},
			want:  "#1a1b26",
		},
		{
			name:  "black",
			color: flair.Color{R: 0, G: 0, B: 0},
			want:  "#000000",
		},
		{
			name:  "white",
			color: flair.Color{R: 255, G: 255, B: 255},
			want:  "#ffffff",
		},
		{
			name:  "red",
			color: flair.Color{R: 255, G: 0, B: 0},
			want:  "#ff0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.Hex()
			if got != tt.want {
				t.Errorf("Color.Hex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_Equal(t *testing.T) {
	tests := []struct {
		name  string
		c1    flair.Color
		c2    flair.Color
		equal bool
	}{
		{
			name:  "identical colors",
			c1:    flair.Color{R: 26, G: 27, B: 38},
			c2:    flair.Color{R: 26, G: 27, B: 38},
			equal: true,
		},
		{
			name:  "different red",
			c1:    flair.Color{R: 26, G: 27, B: 38},
			c2:    flair.Color{R: 27, G: 27, B: 38},
			equal: false,
		},
		{
			name:  "different green",
			c1:    flair.Color{R: 26, G: 27, B: 38},
			c2:    flair.Color{R: 26, G: 28, B: 38},
			equal: false,
		},
		{
			name:  "different blue",
			c1:    flair.Color{R: 26, G: 27, B: 38},
			c2:    flair.Color{R: 26, G: 27, B: 39},
			equal: false,
		},
		{
			name:  "completely different",
			c1:    flair.Color{R: 0, G: 0, B: 0},
			c2:    flair.Color{R: 255, G: 255, B: 255},
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c1.Equal(tt.c2)
			if got != tt.equal {
				t.Errorf("Color.Equal() = %v, want %v", got, tt.equal)
			}
		})
	}
}

func TestParseHex_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  flair.Color
	}{
		{
			name:  "6-digit with hash",
			input: "#1a1b26",
			want:  flair.Color{R: 26, G: 27, B: 38},
		},
		{
			name:  "6-digit without hash",
			input: "1a1b26",
			want:  flair.Color{R: 26, G: 27, B: 38},
		},
		{
			name:  "3-digit with hash",
			input: "#fff",
			want:  flair.Color{R: 255, G: 255, B: 255},
		},
		{
			name:  "3-digit without hash",
			input: "f00",
			want:  flair.Color{R: 255, G: 0, B: 0},
		},
		{
			name:  "uppercase hex",
			input: "#AABBCC",
			want:  flair.Color{R: 170, G: 187, B: 204},
		},
		{
			name:  "mixed case hex",
			input: "#AaBbCc",
			want:  flair.Color{R: 170, G: 187, B: 204},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flair.ParseHex(tt.input)
			if err != nil {
				t.Fatalf("ParseHex(%q) returned error: %v", tt.input, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseHex(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseHex_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "too short",
			input: "#ff",
		},
		{
			name:  "too long",
			input: "#fffffff",
		},
		{
			name:  "invalid characters",
			input: "#gggggg",
		},
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "only hash",
			input: "#",
		},
		{
			name:  "4 digits",
			input: "#1234",
		},
		{
			name:  "5 digits",
			input: "#12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := flair.ParseHex(tt.input)
			if err == nil {
				t.Errorf("ParseHex(%q) expected error, got nil", tt.input)
			}
		})
	}
}

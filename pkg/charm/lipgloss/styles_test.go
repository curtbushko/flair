package lipgloss_test

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/lipgloss"

	flairlip "github.com/curtbushko/flair/pkg/charm/lipgloss"
)

func TestStyles_HasAllCategories(t *testing.T) {
	// Given: A Styles struct
	stylesType := reflect.TypeOf(flairlip.Styles{})

	// Then: It should have all expected fields
	expectedFields := []struct {
		name     string
		fieldTyp reflect.Type
	}{
		// Surface styles
		{"Background", reflect.TypeOf(lipgloss.Style{})},
		{"Raised", reflect.TypeOf(lipgloss.Style{})},
		{"Sunken", reflect.TypeOf(lipgloss.Style{})},
		{"Overlay", reflect.TypeOf(lipgloss.Style{})},
		{"Popup", reflect.TypeOf(lipgloss.Style{})},

		// Text styles
		{"Text", reflect.TypeOf(lipgloss.Style{})},
		{"Secondary", reflect.TypeOf(lipgloss.Style{})},
		{"Muted", reflect.TypeOf(lipgloss.Style{})},
		{"Inverse", reflect.TypeOf(lipgloss.Style{})},

		// Status styles
		{"Error", reflect.TypeOf(lipgloss.Style{})},
		{"Warning", reflect.TypeOf(lipgloss.Style{})},
		{"Success", reflect.TypeOf(lipgloss.Style{})},
		{"Info", reflect.TypeOf(lipgloss.Style{})},

		// Border styles
		{"Border", reflect.TypeOf(lipgloss.Style{})},
		{"BorderFocus", reflect.TypeOf(lipgloss.Style{})},

		// Component styles
		{"Button", reflect.TypeOf(lipgloss.Style{})},
		{"ButtonFocused", reflect.TypeOf(lipgloss.Style{})},
		{"Input", reflect.TypeOf(lipgloss.Style{})},
		{"InputFocused", reflect.TypeOf(lipgloss.Style{})},
		{"ListItem", reflect.TypeOf(lipgloss.Style{})},
		{"ListSelected", reflect.TypeOf(lipgloss.Style{})},
		{"Table", reflect.TypeOf(lipgloss.Style{})},
		{"TableHeader", reflect.TypeOf(lipgloss.Style{})},
		{"Dialog", reflect.TypeOf(lipgloss.Style{})},
	}

	for _, expected := range expectedFields {
		t.Run(expected.name, func(t *testing.T) {
			field, ok := stylesType.FieldByName(expected.name)
			if !ok {
				t.Errorf("Styles struct is missing field %q", expected.name)
				return
			}
			if field.Type != expected.fieldTyp {
				t.Errorf("field %q has type %v, want %v", expected.name, field.Type, expected.fieldTyp)
			}
		})
	}
}

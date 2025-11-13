package colors_test

import (
	"testing"

	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestValidHexColor(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectsValid bool
	}{
		{"valid 6-digit hex", "#FF5733", true},
		{"valid lowercase", "#ff5733", true},
		{"valid mixed case", "#Ff5733", true},
		{"empty string", "", false},
		{"no hash", "FF5733", false},
		{"3-digit hex", "#FFF", false},
		{"7-digit hex", "#FF57331", false},
		{"invalid chars", "#GG5733", false},
		{"spaces", "#FF 733", false},
		{"none value", "none", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := colors.GetColor(tt.input)

			if tt.expectsValid {
				_, isColor := result.(lipgloss.Color)
				assert.True(t, isColor, "Expected valid color for input: %s", tt.input)
			} else {
				_, isNoColor := result.(lipgloss.NoColor)
				assert.True(t, isNoColor, "Expected NoColor for invalid input: %s", tt.input)
			}
		})
	}
}

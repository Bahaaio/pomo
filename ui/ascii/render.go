// Package ascii provides ASCII art style.
package ascii

import (
	"github.com/charmbracelet/lipgloss"
)

func RenderNumber(number string, font Font) string {
	digits := make([]string, 0, len(number))

	for _, digit := range number {
		digits = append(digits, renderDigit(digit, font))
	}

	asciiDigits := lipgloss.JoinHorizontal(lipgloss.Top, digits...)
	return asciiDigits
}

func GetFont(fontName string) Font {
	if font, exists := fonts[fontName]; exists {
		return font
	}

	return fonts[DefaultFont]
}

func renderDigit(digit rune, font Font) string {
	if digit == ':' {
		return font[len(font)-1]
	}

	if digit < '0' || digit > '9' {
		return ""
	}

	return font[digit-'0']
}

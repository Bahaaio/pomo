// Package colors defines color constants for UI elements.
package colors

import (
	"log"
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

const (
	InactiveButtonFg = lipgloss.Color("#FFF7DB")
	InactiveButtonBg = lipgloss.Color("#888B7E")
	ActiveButtonFg   = InactiveButtonFg
	ActiveButtonBg   = lipgloss.Color("#F25D94")
	BorderFg         = lipgloss.Color("#5A56E0")
	SuccessMessageFg = lipgloss.Color("#198754")
	TimerFg          = lipgloss.Color("#5A56E0")
)

var validColorRegex *regexp.Regexp = nil

func init() {
	var err error
	validColorRegex, err = regexp.Compile("^#[0-9a-fA-F]{6}$")
	if err != nil {
		log.Println("failed to compile isHex regex:", err)
	}
}

func GetColor(color string) lipgloss.TerminalColor {
	if validColorRegex == nil || !validColorRegex.MatchString(color) {
		log.Println("using no color")
		return lipgloss.NoColor{}
	}

	log.Println("using color:", color)
	return lipgloss.Color(color)
}

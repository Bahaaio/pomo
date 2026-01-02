// Package colors defines color constants for UI elements.
package colors

import (
	"log"
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

// Palette
const (
	Purple  = lipgloss.Color("#5A56E0")
	Pink    = lipgloss.Color("#F25D94")
	Cream   = lipgloss.Color("#FFF7DB")
	Gray    = lipgloss.Color("#888B7E")
	Green   = lipgloss.Color("#198754")
	Blue    = lipgloss.Color("#4A9EFF")
	NoColor = lipgloss.Color("default")
)

const (
	// Timer & primary UI
	TimerFg  = Purple
	BorderFg = Purple

	// Session types
	WorkSessionFg  = Purple
	BreakSessionFg = NoColor

	// Buttons
	InactiveButtonFg = Cream
	InactiveButtonBg = Gray
	ActiveButtonFg   = Cream
	ActiveButtonBg   = Pink

	// Messages
	SuccessMessageFg = Green
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

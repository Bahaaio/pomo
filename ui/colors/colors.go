// Package colors defines color constants for UI elements.
package colors

import (
	"log"
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

// Palette
const (
	Purple      = lipgloss.Color("#8860FF")
	PurpleDark  = lipgloss.Color("#8840FF")
	PurpleLight = lipgloss.Color("#A070FF")
	PurplePale  = lipgloss.Color("#C3A1FF")
	Pink        = lipgloss.Color("#F25D94")
	Red         = lipgloss.Color("#FF4C4C")
	Cream       = lipgloss.Color("#FFF7DB")
	Gray        = lipgloss.Color("#888B7E")
	Green       = lipgloss.Color("#198754")
	Blue        = lipgloss.Color("#4A9EFF")
	DimGray     = lipgloss.Color("#606060")
	NoColor     = lipgloss.Color("default")
)

const (
	// Timer & primary UI
	TimerFg  = Purple
	BorderFg = Purple
	PauseFg  = DimGray

	// heat map
	HeatMapFg0 = DimGray
	HeatMapFg1 = PurpleDark
	HeatMapFg2 = Purple
	HeatMapFg3 = PurpleLight
	HeatMapFg4 = PurplePale

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
	ErrorMessageFg   = Red
)

var validColorRegex *regexp.Regexp = nil

func init() {
	var err error
	validColorRegex, err = regexp.Compile("^#[0-9a-fA-F]{6}$")
	if err != nil {
		log.Println("failed to compile isHex regex:", err)
	}
}

// GetColor returns a [lipgloss.TerminalColor] based on the provided color string.
// If the color string is not a valid hex color code, it returns [lipgloss.NoColor].
func GetColor(color string) lipgloss.TerminalColor {
	if validColorRegex == nil || !validColorRegex.MatchString(color) {
		log.Println("using no color")
		return lipgloss.NoColor{}
	}

	log.Println("using color:", color)
	return lipgloss.Color(color)
}

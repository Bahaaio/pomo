package ui

import (
	"fmt"

	"github.com/Bahaaio/pomo/ui/ascii"
)

const (
	maxWidth           = 80
	margin             = 4
	padding            = 2
	separator          = " — "
	pausedIndicator    = "(paused)"
	completedIndicator = "done!"
)

func (m *Model) buildMainContent() string {
	timeLeft := m.buildTimeLeft()

	if m.useTimerArt {
		return timeLeft + "\n\n" + m.title
	}

	content := m.title
	if !m.timer.Timedout() {
		content += separator + timeLeft
	}

	return content
}

func (m *Model) buildStatusIndicators() string {
	if m.timer.Timedout() {
		return separator + completedIndicator
	}

	if m.paused {
		return " " + pausedIndicator
	}

	return ""
}

func (m *Model) buildProgressBar() string {
	return "\n\n" + m.progressBar.View() + "\n"
}

// returns time left as a string in HH:MM:SS format
func (m *Model) buildTimeLeft() string {
	left := m.timer.Timeout
	hours := int(left.Hours())
	minutes := int(left.Minutes()) % 60
	seconds := int(left.Seconds()) % 60

	time := ""

	// only show hours if they are non-zero
	if hours > 0 {
		time += fmt.Sprintf("%02d:", hours)
	}
	time += fmt.Sprintf("%02d:%02d", minutes, seconds)

	if m.useTimerArt {
		time = ascii.RenderNumber(time, m.timerFont)
		return m.asciiTimerStyle.Render(time)
	}

	return time
}

func (m *Model) buildHelpView() string {
	return m.help.View(keyMap)
}

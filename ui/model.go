package ui

import (
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui/ascii"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// components
	progressBar progress.Model
	help        help.Model

	// timer
	timer           timer.Model
	duration        time.Duration
	initialDuration time.Duration
	passed          time.Duration

	// state
	title         string
	width, height int // window dimensions
	paused        bool
	quitting      bool
	exitStatus    ExitStatus

	// ASCII art
	useTimerArt     bool
	timerFont       ascii.Font
	asciiTimerStyle lipgloss.Style
}

func NewModel(task config.Task, asciiArt config.ASCIIArt) Model {
	var timerFont ascii.Font
	timerStyle := lipgloss.NewStyle()

	if asciiArt.Enabled {
		timerFont = ascii.GetFont(asciiArt.Font)

		timerColor := colors.GetColor(asciiArt.Color)
		timerStyle = timerStyle.Foreground(timerColor)
	}

	return Model{
		title:           task.Title,
		timer:           timer.New(task.Duration),
		progressBar:     progress.New(progress.WithDefaultGradient()),
		duration:        task.Duration,
		initialDuration: task.Duration,
		useTimerArt:     asciiArt.Enabled,
		timerFont:       timerFont,
		asciiTimerStyle: timerStyle,
		help:            help.New(),
		exitStatus:      Quit,
	}
}

type ExitStatus byte

const (
	Completed ExitStatus = iota
	Skipped
	Quit
)

func (e ExitStatus) String() string {
	switch e {
	case Completed:
		return "Completed"
	case Skipped:
		return "Skipped"
	case Quit:
		return "Quit"
	default:
		return "Unknown"
	}
}

func (m Model) ExitStatus() ExitStatus {
	return m.exitStatus
}

func (m Model) Elapsed() time.Duration {
	return m.passed
}

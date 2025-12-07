package ui

import (
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui/ascii"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/Bahaaio/pomo/ui/confirm"
	"github.com/Bahaaio/pomo/ui/summary"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// components
	progressBar   progress.Model
	confirmDialog confirm.Model
	help          help.Model

	// timer
	timer    timer.Model
	duration time.Duration
	elapsed  time.Duration

	// state
	width, height       int // window dimensions
	shouldAskToContinue bool
	sessionState        SessionState
	currentTaskType     config.TaskType
	currentTask         config.Task
	sessionSummary      summary.SessionSummary
	isShortSession      bool

	// ASCII art
	useTimerArt     bool
	timerFont       ascii.Font
	asciiTimerStyle lipgloss.Style
}

func NewModel(taskType config.TaskType, asciiArt config.ASCIIArt, askToContinue bool) Model {
	task := taskType.GetTask()

	var timerFont ascii.Font
	timerStyle := lipgloss.NewStyle()

	if asciiArt.Enabled {
		timerFont = ascii.GetFont(asciiArt.Font)

		timerColor := colors.GetColor(asciiArt.Color)
		timerStyle = timerStyle.Foreground(timerColor)
	}

	return Model{
		progressBar:   progress.New(progress.WithDefaultGradient()),
		confirmDialog: confirm.New(),
		help:          help.New(),

		timer:    timer.New(task.Duration),
		duration: task.Duration,

		shouldAskToContinue: askToContinue,
		sessionState:        Running,
		currentTaskType:     taskType,
		currentTask:         *task,
		sessionSummary:      summary.SessionSummary{},

		useTimerArt:     asciiArt.Enabled,
		timerFont:       timerFont,
		asciiTimerStyle: timerStyle,
	}
}

type SessionState byte

const (
	Running SessionState = iota
	Paused
	ShowingConfirm
	Quitting
)

func (m Model) GetSessionSummary() summary.SessionSummary {
	return m.sessionSummary
}

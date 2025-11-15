// Package ui provides the terminal user interface for pomodoro sessions.
package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui/ascii"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxWidth           = 80
	margin             = 4
	padding            = 2
	interval           = time.Second
	separator          = " — "
	pausedIndicator    = "(paused)"
	completedIndicator = "done!"
)

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

type Model struct {
	title           string
	timer           timer.Model
	progress        progress.Model
	duration        time.Duration
	initialDuration time.Duration
	passed          time.Duration
	width, height   int
	paused          bool
	useTimerArt     bool
	timerFont       ascii.Font
	ASCIITimerStyle lipgloss.Style
	help            help.Model
	quitting        bool
	exitStatus      ExitStatus
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
		timer:           timer.NewWithInterval(task.Duration, interval),
		progress:        progress.New(progress.WithDefaultGradient()),
		duration:        task.Duration,
		initialDuration: task.Duration,
		useTimerArt:     asciiArt.Enabled,
		timerFont:       timerFont,
		ASCIITimerStyle: timerStyle,
		help:            help.New(),
		exitStatus:      Quit,
	}
}

func (m Model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Increase):
			m.duration += time.Minute
			return m, m.resetTimer()

		case key.Matches(msg, Keys.Pause):
			m.paused = !m.paused

			if !m.paused {
				return m, m.timer.Start()
			}
			return m, nil

		case key.Matches(msg, Keys.Reset):
			m.passed = 0
			m.duration = m.initialDuration
			return m, m.resetTimer()

		case key.Matches(msg, Keys.Skip):
			m.exitStatus = Skipped
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, Keys.Quit):
			m.exitStatus = Quit
			m.quitting = true
			return m, tea.Quit

		default:
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = min(m.width-2*padding-margin, maxWidth)
		return m, nil

	case timer.TickMsg:
		if m.paused {
			return m, nil
		}

		var cmds []tea.Cmd

		m.passed += m.timer.Interval

		percent := m.getPercent()
		cmds = append(cmds, m.progress.SetPercent(percent))

		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)

		if m.progress.Percent() >= 1.0 && !m.progress.IsAnimating() {
			log.Println("timer completed")

			m.exitStatus = Completed
			m.quitting = true
			return m, tea.Quit
		}

		return m, cmd

	default:
		return m, nil
	}
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	content := m.buildMainContent()
	content += m.buildStatusIndicators()
	content += "\n\n" + m.progress.View() + "\n"

	help := m.help.View(Keys)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, content, help),
	)
}

func (m Model) buildMainContent() string {
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

func (m Model) buildStatusIndicators() string {
	if m.timer.Timedout() {
		return separator + completedIndicator
	}

	if m.paused {
		return " " + pausedIndicator
	}

	return ""
}

// returns time left as a string in HH:MM:SS format
func (m Model) buildTimeLeft() string {
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
		time = ascii.ToASCIIArt(time, m.timerFont)
		return m.ASCIITimerStyle.Render(time)
	}

	return time
}

func (m Model) ExitStatus() ExitStatus {
	return m.exitStatus
}

func (m *Model) resetTimer() tea.Cmd {
	// reset timer with new duration minus passed time
	m.timer.Timeout = m.duration - m.passed

	// update progress bar
	return m.progress.SetPercent(m.getPercent())
}

func (m Model) getPercent() float64 {
	passed := float64(m.passed.Milliseconds())
	duration := float64(m.duration.Milliseconds())

	return passed / duration
}

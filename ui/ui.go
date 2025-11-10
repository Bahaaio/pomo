// Package ui provides the terminal user interface for pomodoro sessions.
package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxWidth = 80
	margin   = 4
	padding  = 2
	interval = time.Second
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
	exitStatus      ExitStatus
	help            help.Model
	quitting        bool
}

func NewModel(task config.Task) Model {
	return Model{
		title:           task.Title,
		timer:           timer.NewWithInterval(task.Duration, interval),
		progress:        progress.New(progress.WithDefaultGradient()),
		duration:        task.Duration,
		initialDuration: task.Duration,
		exitStatus:      Quit,
		help:            help.New(),
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
			if !config.C.AskToContinue {
				return m, nil
			}

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

	s := m.title + " — "
	if m.timer.Timedout() {
		s += "done!"
	} else {
		left := m.timer.Timeout
		hours := int(left.Hours())
		minutes := int(left.Minutes()) % 60
		seconds := int(left.Seconds()) % 60

		// HH:MM:SS format
		// only show hours if they are non-zero
		if hours > 0 {
			s += fmt.Sprintf("%02d:", hours)
		}
		s += fmt.Sprintf("%02d:%02d", minutes, seconds)

		// Show pause indicator
		if m.paused {
			s += " (paused)"
		}
	}

	s += "\n\n" +
		m.progress.View() + "\n"

	help := m.help.View(Keys)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			s,
			help,
		),
	)
}

func (m Model) ExitStatus() ExitStatus {
	return m.exitStatus
}

func (m *Model) resetTimer() tea.Cmd {
	m.timer = timer.NewWithInterval(
		m.duration-m.passed,
		interval,
	)

	return tea.Batch(
		m.progress.SetPercent(m.getPercent()),
		m.timer.Start(),
	)
}

func (m Model) getPercent() float64 {
	passed := float64(m.passed.Milliseconds())
	duration := float64(m.duration.Milliseconds())

	return passed / duration
}

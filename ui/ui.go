// Package ui provides the terminal user interface for pomodoro sessions.
package ui

import (
	"time"

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
)

type model struct {
	name            string
	timer           timer.Model
	progress        progress.Model
	duration        time.Duration
	initialDuration time.Duration
	passed          time.Duration
	width           int
	height          int
	altScreen       bool
	help            help.Model
	quitting        bool
}

const interval = time.Second

func NewModel(duration time.Duration, taskName string, altScreen bool) model {
	return model{
		name:            taskName,
		timer:           timer.NewWithInterval(duration, interval),
		progress:        progress.New(progress.WithDefaultGradient()),
		duration:        duration,
		initialDuration: duration,
		altScreen:       altScreen,
		help:            help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Increase):
			m.duration += time.Minute
			return m, m.resetTimer()

		case key.Matches(msg, Keys.Decrease):
			m.duration -= time.Minute
			return m, m.resetTimer()

		case key.Matches(msg, Keys.Reset):
			m.passed = 0
			m.duration = m.initialDuration
			return m, m.resetTimer()

		case key.Matches(msg, Keys.Quit):
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
		var cmds []tea.Cmd

		m.passed += m.timer.Interval

		percent := float64(m.passed.Milliseconds()) / float64(m.duration.Milliseconds())
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

		if m.progress.Percent() == 1.0 && !m.progress.IsAnimating() {
			m.quitting = true
			return m, tea.Quit
		}

		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := m.name + " session — "
	if m.timer.Timedout() {
		s = "done!"
	} else {
		s += m.timer.View()
	}

	s += "\n\n" +
		m.progress.View() + "\n"

	help := m.help.View(Keys)

	if m.altScreen {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				s,
				help,
			),
		)
	}

	return lipgloss.Place(
		m.width,
		1,
		lipgloss.Center,
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"\n",
			s,
			help,
		),
	)
}

func (m *model) resetTimer() tea.Cmd {
	m.timer = timer.NewWithInterval(
		m.duration-m.passed,
		interval,
	)

	return m.timer.Init()
}

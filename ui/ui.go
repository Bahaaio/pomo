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
)

type Model struct {
	title           string
	timer           timer.Model
	progress        progress.Model
	duration        time.Duration
	initialDuration time.Duration
	passed          time.Duration
	width           int
	height          int
	altScreen       bool
	paused          bool
	help            help.Model
	quitting        bool
}

const interval = time.Second

func NewModel(task config.Task, altScreen bool) Model {
	return Model{
		title:           task.Title,
		timer:           timer.NewWithInterval(task.Duration, interval),
		progress:        progress.New(progress.WithDefaultGradient()),
		duration:        task.Duration,
		initialDuration: task.Duration,
		altScreen:       altScreen,
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

		case key.Matches(msg, Keys.Decrease):
			if m.duration > time.Minute {
				m.duration -= time.Minute
				return m, m.resetTimer()
			}
			return m, nil

		case key.Matches(msg, Keys.Pause):
			m.paused = !m.paused

			if !m.paused {
				return m, m.resetTimer()
			}
			return m, nil

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
		if m.paused {
			return m, nil
		}

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
			log.Println("timer completed")
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

		// HH:MM:SS format
		// only show hours if they are non-zero
		if int(left.Hours()) > 0 {
			s += fmt.Sprintf("%02d:", int(left.Hours()))
		}
		s += fmt.Sprintf("%02d:%02d", int(left.Minutes())%60, int(left.Seconds())%60)

		// Show pause indicator
		if m.paused {
			s += " (paused)"
		}
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

func (m *Model) resetTimer() tea.Cmd {
	m.timer = timer.NewWithInterval(
		m.duration-m.passed,
		interval,
	)

	return m.timer.Init()
}

func (m Model) TimerCompleted() bool {
	return m.timer.Timedout()
}

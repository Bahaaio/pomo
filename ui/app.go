// Package ui provides the terminal user interface for pomodoro sessions.
package ui

import (
	"github.com/Bahaaio/pomo/ui/confirm"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.handleKeys(msg)

	case tea.WindowSizeMsg:
		return m, m.handleWindowResize(msg)

	case timer.TickMsg:
		return m, m.handleTimerTick(msg)

	case timer.StartStopMsg:
		return m, m.handleTimerStartStop(msg)

	case confirm.ChoiceMsg:
		return m, m.handleConfirmChoice(msg)

	case progress.FrameMsg:
		return m, m.handleProgressBarFrame(msg)

	default:
		return m, nil
	}
}

func (m Model) View() string {
	if m.sessionState == Quitting {
		return ""
	}

	// show confirmation dialog
	if m.sessionState == ShowingConfirm {
		title := m.currentTaskType.Opposite().GetTask().Title
		return m.confirmDialog.View("start " + title + "?")
	}

	content := m.buildMainContent()
	content += m.buildStatusIndicators()
	content += m.buildProgressBar()

	help := m.buildHelpView()

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, content, help),
	)
}

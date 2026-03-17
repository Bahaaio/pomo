// Package ui provides the terminal user interface for pomodoro sessions.
package ui

import (
	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/sound"
	"github.com/Bahaaio/pomo/ui/confirm"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// initSession runs the onSessionStart and duringSession hooks for the initial session
func (m Model) initSession() (tea.Model, tea.Cmd) {
	// determine which hooks to use (per-task overrides global)
	var onStartCmds, duringCmds [][]string
	if len(m.currentTask.OnStart) > 0 {
		onStartCmds = m.currentTask.OnStart
	} else {
		onStartCmds = config.C.OnSessionStart
	}
	if len(m.currentTask.During) > 0 {
		duringCmds = m.currentTask.During
	} else {
		duringCmds = config.C.DuringSession
	}

	// run start actions (fire and forget)
	for _, cmd := range onStartCmds {
		sound.PlayOnce(cmd[0])
	}

	// run during actions (ambient sounds)
	if len(duringCmds) > 0 {
		// Use the first sound file for looping
		m.duringSoundPlayer.PlayLoop(duringCmds[0][0])
	}

	return m, m.timer.Init()
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return initSessionMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initSessionMsg:
		return m.initSession()

	case tea.KeyMsg:
		return m, m.handleKeys(msg)

	case tea.WindowSizeMsg:
		return m, m.handleWindowResize(msg)

	case timer.TickMsg:
		return m, m.handleTimerTick(msg)

	case confirmTickMsg:
		return m, m.handleConfirmTick()

	case timer.StartStopMsg:
		return m, m.handleTimerStartStop(msg)

	case progress.FrameMsg:
		return m, m.handleProgressBarFrame(msg)

	case confirm.ChoiceMsg:
		return m, m.handleConfirmChoice(msg)

	case commandsDoneMsg:
		return m, m.handleCommandsDone()

	default:
		return m, nil
	}
}

func (m Model) View() string {
	if m.sessionState == Quitting {
		return ""
	}

	if m.sessionState == WaitingForCommands {
		return m.buildWaitingForCommandsView()
	}

	// show confirmation dialog
	if m.sessionState == ShowingConfirm {
		return m.buildConfirmDialogView()
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

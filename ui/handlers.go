package ui

import (
	"log"
	"time"

	"github.com/Bahaaio/pomo/actions"
	"github.com/Bahaaio/pomo/ui/confirm"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	if m.sessionState == ShowingConfirm {
		return m.confirmDialog.HandleKeys(msg)
	}

	switch {
	case key.Matches(msg, keyMap.Increase):
		m.duration += time.Minute
		return m.updateProgressBar()

	case key.Matches(msg, keyMap.Pause):
		if m.sessionState == Paused {
			m.sessionState = Running
		} else {
			m.sessionState = Paused
		}

		if m.sessionState == Running {
			return m.timer.Start()
		}

		return nil

	case key.Matches(msg, keyMap.Reset):
		m.elapsed = 0
		m.duration = m.currentTask.Duration
		return m.updateProgressBar()

	case key.Matches(msg, keyMap.Skip):
		m.recordSession()
		return m.nextSession()

	case key.Matches(msg, keyMap.Quit):
		m.recordSession()
		return m.Quit()

	default:
		return nil
	}
}

func (m *Model) handleConfirmChoice(msg confirm.ChoiceMsg) tea.Cmd {
	switch msg.Choice {
	case confirm.Confirm:
		return m.nextSession()
	case confirm.Cancel:
		return m.Quit()
	}

	return nil
}

func (m *Model) handleWindowResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.confirmDialog.HandleWindowResize(msg) // always update it

	m.width = msg.Width
	m.height = msg.Height
	m.progressBar.Width = min(m.width-2*padding-margin, maxWidth)

	return nil
}

func (m *Model) handleTimerTick(msg timer.TickMsg) tea.Cmd {
	if m.sessionState == Paused {
		return nil
	}

	var cmds []tea.Cmd

	m.elapsed += m.timer.Interval

	percent := m.getPercent()
	cmds = append(cmds, m.progressBar.SetPercent(percent))

	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) handleTimerStartStop(msg timer.StartStopMsg) tea.Cmd {
	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)

	return cmd
}

func (m *Model) handleProgressBarFrame(msg progress.FrameMsg) tea.Cmd {
	if m.progressBar.Percent() >= 1.0 && !m.progressBar.IsAnimating() && m.sessionState == Running {
		return m.handleCompletion()
	}

	progressModel, cmd := m.progressBar.Update(msg)
	m.progressBar = progressModel.(progress.Model)

	return cmd
}

func (m *Model) updateProgressBar() tea.Cmd {
	// reset timer with new duration minus passed time
	m.timer.Timeout = m.duration - m.elapsed

	// update progress bar
	return m.progressBar.SetPercent(m.getPercent())
}

func (m Model) getPercent() float64 {
	passed := float64(m.elapsed.Milliseconds())
	duration := float64(m.duration.Milliseconds())

	return passed / duration
}

func (m *Model) handleCompletion() tea.Cmd {
	log.Println("timer completed")

	m.recordSession()
	actions.RunPostActions(&m.currentTask).Wait()

	if m.shouldAskToContinue {
		m.sessionState = ShowingConfirm
		return nil
	}

	return m.Quit()
}

// updates model with next task and starts the timer
func (m *Model) nextSession() tea.Cmd {
	m.currentTaskType = m.currentTaskType.Opposite()
	m.currentTask = *m.currentTaskType.GetTask()

	m.elapsed = 0
	m.duration = m.currentTask.Duration
	m.timer = timer.New(m.currentTask.Duration)

	m.sessionState = Running
	return tea.Batch(
		m.progressBar.SetPercent(0.0),
		m.timer.Start(),
	)
}

// records the current session into the session summary
func (m *Model) recordSession() {
	m.sessionSummary.AddSession(m.currentTaskType, m.elapsed)
}

func (m *Model) Quit() tea.Cmd {
	m.sessionState = Quitting
	return tea.Quit
}

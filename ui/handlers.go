package ui

import (
	"log"
	"time"

	"github.com/Bahaaio/pomo/actions"
	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/db"
	"github.com/Bahaaio/pomo/ui/confirm"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type confirmTickMsg struct{}

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
	case confirm.ShortSession:
		return m.shortSession()
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

func (m *Model) handleConfirmTick() tea.Cmd {
	// send tick every second to update idle time
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return confirmTickMsg{}
	})
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

	// continue after the completion according to config
	switch m.onSessionEnd {
	case "ask":
		m.sessionState = ShowingConfirm
		m.confirmStartTime = time.Now()

		// send first confirm tick
		return func() tea.Msg {
			return confirmTickMsg{}
		}
	case "skip":
		return m.nextSession()
	case "quit":
		return m.Quit()
	default:
		log.Printf("unknown onSessionEnd value %q, defaulting to quit", m.onSessionEnd)
		return m.Quit()
	}
}

// starts session with the opposite task type (work <-> break)
func (m *Model) nextSession() tea.Cmd {
	nextTaskType := m.currentTaskType.Opposite()
	return m.startSession(nextTaskType, *nextTaskType.GetTask(), false)
}

// starts a short session of the current task type
func (m *Model) shortSession() tea.Cmd {
	shortTask := m.currentTask
	shortTask.Duration = 2 * time.Minute // TODO: make configurable
	shortTask.Title = "short " + m.currentTaskType.GetTask().Title

	return m.startSession(m.currentTaskType, shortTask, true)
}

// initializes and starts a new session with the given task
func (m *Model) startSession(taskType config.TaskType, task config.Task, isShortSession bool) tea.Cmd {
	m.isShortSession = isShortSession
	m.currentTaskType = taskType
	m.currentTask = task

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
	// ignore very short or zero duration sessions
	if m.elapsed < time.Second {
		return
	}

	// short sessions extend the current session without incrementing the count
	if m.isShortSession {
		m.sessionSummary.AddDuration(m.currentTaskType, m.elapsed)
		return
	}

	m.sessionSummary.AddSession(m.currentTaskType, m.elapsed)

	// return if no database is configured
	if m.repo == nil {
		return
	}

	if err := m.repo.CreateSession(
		time.Now(),
		m.elapsed,
		db.GetSessionType(m.currentTaskType),
	); err != nil {
		log.Printf("failed to record session: %v", err)
	}
}

func (m *Model) Quit() tea.Cmd {
	m.sessionState = Quitting
	return tea.Quit
}

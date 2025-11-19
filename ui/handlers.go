package ui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, keyMap.Increase):
		m.duration += time.Minute
		return m.updateProgressBar()

	case key.Matches(msg, keyMap.Pause):
		m.paused = !m.paused

		if !m.paused {
			return m.timer.Start()
		}
		return nil

	case key.Matches(msg, keyMap.Reset):
		m.passed = 0
		m.duration = m.initialDuration
		return m.updateProgressBar()

	case key.Matches(msg, keyMap.Skip):
		m.exitStatus = Skipped
		m.quitting = true
		return tea.Quit

	case key.Matches(msg, keyMap.Quit):
		m.exitStatus = Quit
		m.quitting = true
		return tea.Quit

	default:
		return nil
	}
}

func (m *Model) handleWindowResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.width = msg.Width
	m.height = msg.Height
	m.progressBar.Width = min(m.width-2*padding-margin, maxWidth)

	return nil
}

func (m *Model) handleTimerTick(msg timer.TickMsg) tea.Cmd {
	if m.paused {
		return nil
	}

	var cmds []tea.Cmd

	m.passed += m.timer.Interval

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
	progressModel, cmd := m.progressBar.Update(msg)
	m.progressBar = progressModel.(progress.Model)

	if m.progressBar.Percent() >= 1.0 && !m.progressBar.IsAnimating() {
		log.Println("timer completed")

		m.exitStatus = Completed
		m.quitting = true
		return tea.Quit
	}

	return cmd
}

func (m *Model) updateProgressBar() tea.Cmd {
	// reset timer with new duration minus passed time
	m.timer.Timeout = m.duration - m.passed

	// update progress bar
	return m.progressBar.SetPercent(m.getPercent())
}

func (m Model) getPercent() float64 {
	passed := float64(m.passed.Milliseconds())
	duration := float64(m.duration.Milliseconds())

	return passed / duration
}

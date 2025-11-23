// Package summary tracks pomodoro sessions and renders visual summary with progress bar.
package summary

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/lipgloss"
)

var messageStyle = lipgloss.NewStyle().Foreground(colors.SuccessMessageFg)

type SessionSummary struct {
	totalWorkSessions int
	totalWorkDuration time.Duration

	totalBreakSessions int
	totalBreakDuration time.Duration
}

// AddSession adds a session to the summary based on the task type and elapsed time.
func (t *SessionSummary) AddSession(taskType config.TaskType, elapsed time.Duration) {
	if taskType == config.WorkTask {
		t.totalWorkSessions++
		t.totalWorkDuration += elapsed
	} else {
		t.totalBreakSessions++
		t.totalBreakDuration += elapsed
	}
}

// Print prints the session summary to the console.
func (t SessionSummary) Print() {
	if t.totalWorkDuration == 0 && t.totalBreakDuration == 0 {
		return
	}

	sessionIndicator := "sessions"
	if t.totalWorkSessions == 1 {
		sessionIndicator = "session"
	}

	fmt.Println(messageStyle.Render("Session Summary:"))

	if t.totalWorkDuration > 0 {
		fmt.Printf(" Work : %v (%d %s)\n", t.totalWorkDuration, t.totalWorkSessions, sessionIndicator)
	}

	if t.totalBreakDuration > 0 {
		fmt.Printf(" Break: %v (%d %s)\n", t.totalBreakDuration, t.totalBreakSessions, sessionIndicator)
	}

	if t.totalBreakDuration > 0 && t.totalWorkDuration > 0 {
		fmt.Println(" Total:", t.totalWorkDuration+t.totalBreakDuration)
	}

	if t.totalWorkDuration > 0 {
		t.printProgressBar()
	}
}

// prints a progress bar showing the ratio of work to total time.
func (t SessionSummary) printProgressBar() {
	const barWidth = 30

	totalDuration := t.totalWorkDuration + t.totalBreakDuration
	workRatio := float64(t.totalWorkDuration.Milliseconds()) / float64(totalDuration.Milliseconds())

	filledWidth := int(workRatio * barWidth)
	emptyWidth := barWidth - filledWidth

	bar := lipgloss.NewStyle().Foreground(colors.TimerFg).
		Render(strings.Repeat("█", filledWidth)) +
		strings.Repeat("░", emptyWidth)

	fmt.Printf("\n [%s] %.0f%% work\n", bar, workRatio*100)
}

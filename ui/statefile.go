package ui

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Bahaaio/pomo/config"
)

const stateFilePath = "/tmp/pomo.json"

type sessionStateFile struct {
	State            string  `json:"state"`
	Type             string  `json:"type"`
	ElapsedSeconds   int     `json:"elapsed_seconds"`
	DurationSeconds  int     `json:"duration_seconds"`
	RemainingSeconds int     `json:"remaining_seconds"`
	Percent          float64 `json:"percent"`
	SessionsDone     int     `json:"sessions_done"`
	CyclePosition    int     `json:"cycle_position"`
}

func (m Model) writeState() {
	if !m.writeStateFile {
		return
	}

	s := sessionStateFile{
		State:            m.stateString(),
		Type:             m.typeString(),
		ElapsedSeconds:   int(m.elapsed.Seconds()),
		DurationSeconds:  int(m.duration.Seconds()),
		RemainingSeconds: int((m.duration - m.elapsed).Seconds()),
		Percent:          m.getPercent(),
		SessionsDone:     m.sessionSummary.WorkSessions(),
		CyclePosition:    m.cyclePosition,
	}

	data, err := json.Marshal(s)
	if err != nil {
		log.Printf("failed to marshal state file: %v", err)
		return
	}

	if err := os.WriteFile(stateFilePath, data, 0644); err != nil {
		log.Printf("failed to write state file: %v", err)
	}
}

func (m Model) removeStateFile() {
	if !m.writeStateFile {
		return
	}

	if err := os.Remove(stateFilePath); err != nil && !os.IsNotExist(err) {
		log.Printf("failed to remove state file: %v", err)
	}
}

func (m Model) stateString() string {
	switch m.sessionState {
	case Paused:
		return "paused"
	case ShowingConfirm:
		return "completed"
	default:
		return "running"
	}
}

func (m Model) typeString() string {
	if m.currentTaskType == config.WorkTask {
		return "work"
	}
	return "break"
}

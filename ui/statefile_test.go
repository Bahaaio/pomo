package ui

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui/summary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestModel(writeStateFile bool) Model {
	return Model{
		writeStateFile:  writeStateFile,
		sessionState:    Running,
		currentTaskType: config.WorkTask,
		elapsed:         5 * time.Minute,
		duration:        25 * time.Minute,
		sessionSummary:  summary.SessionSummary{},
		cyclePosition:   2,
	}
}

func readStateFile(t *testing.T) sessionStateFile {
	t.Helper()
	data, err := os.ReadFile(stateFilePath)
	require.NoError(t, err)

	var s sessionStateFile
	require.NoError(t, json.Unmarshal(data, &s))
	return s
}

func TestWriteState(t *testing.T) {
	t.Cleanup(func() { os.Remove(stateFilePath) })

	m := newTestModel(true)
	m.writeState()

	s := readStateFile(t)
	assert.Equal(t, "running", s.State)
	assert.Equal(t, "work", s.Type)
	assert.Equal(t, 300, s.ElapsedSeconds)
	assert.Equal(t, 1500, s.DurationSeconds)
	assert.Equal(t, 1200, s.RemainingSeconds)
	assert.Equal(t, 0.2, s.Percent)
	assert.Equal(t, 0, s.SessionsDone)
	assert.Equal(t, 2, s.CyclePosition)
}

func TestWriteStateDisabled(t *testing.T) {
	os.Remove(stateFilePath)

	m := newTestModel(false)
	m.writeState()

	_, err := os.Stat(stateFilePath)
	assert.True(t, os.IsNotExist(err), "state file should not be created when disabled")
}

func TestWriteStatePaused(t *testing.T) {
	t.Cleanup(func() { os.Remove(stateFilePath) })

	m := newTestModel(true)
	m.sessionState = Paused
	m.writeState()

	s := readStateFile(t)
	assert.Equal(t, "paused", s.State)
}

func TestWriteStateBreak(t *testing.T) {
	t.Cleanup(func() { os.Remove(stateFilePath) })

	m := newTestModel(true)
	m.currentTaskType = config.BreakTask
	m.writeState()

	s := readStateFile(t)
	assert.Equal(t, "break", s.Type)
}

func TestWriteStateCompleted(t *testing.T) {
	t.Cleanup(func() { os.Remove(stateFilePath) })

	m := newTestModel(true)
	m.sessionState = ShowingConfirm
	m.writeState()

	s := readStateFile(t)
	assert.Equal(t, "completed", s.State)
}

func TestRemoveStateFile(t *testing.T) {
	m := newTestModel(true)
	m.writeState()

	_, err := os.Stat(stateFilePath)
	require.NoError(t, err, "state file should exist before removal")

	m.removeStateFile()

	_, err = os.Stat(stateFilePath)
	assert.True(t, os.IsNotExist(err), "state file should be removed")
}

func TestRemoveStateFileDisabled(t *testing.T) {
	// create the file manually to ensure removeStateFile respects the flag
	require.NoError(t, os.WriteFile(stateFilePath, []byte("{}"), 0644))
	t.Cleanup(func() { os.Remove(stateFilePath) })

	m := newTestModel(false)
	m.removeStateFile()

	_, err := os.Stat(stateFilePath)
	assert.NoError(t, err, "state file should not be removed when disabled")
}

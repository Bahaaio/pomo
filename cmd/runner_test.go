package cmd

import (
	"testing"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/stretchr/testify/assert"
)

func TestParseArguments(t *testing.T) {
	testCases := []struct {
		name                  string
		args                  []string
		expectedError         bool
		expectedWorkDuration  time.Duration
		expectedBreakDuration time.Duration
	}{
		{
			name:                 "no arguments",
			args:                 []string{},
			expectedError:        false,
			expectedWorkDuration: 0, // should remain unchanged
		},
		{
			name:                 "valid work duration only",
			args:                 []string{"25m"},
			expectedError:        false,
			expectedWorkDuration: 25 * time.Minute,
		},
		{
			name:                  "valid work and break duration",
			args:                  []string{"45m", "15m"},
			expectedError:         false,
			expectedWorkDuration:  45 * time.Minute,
			expectedBreakDuration: 15 * time.Minute,
		},
		{
			name:          "invalid work duration",
			args:          []string{"invalid"},
			expectedError: true,
		},
		{
			name:          "valid work, invalid break",
			args:          []string{"25m", "invalid"},
			expectedError: true,
		},
		{
			name:                  "complex duration formats",
			args:                  []string{"1h30m", "10m30s"},
			expectedError:         false,
			expectedWorkDuration:  90 * time.Minute,
			expectedBreakDuration: 10*time.Minute + 30*time.Second,
		},
	}

	for _, tt := range testCases {
		task := &config.Task{}
		breakTask := &config.Task{}

		result := parseArguments(tt.args, task, breakTask)
		resultIsError := result != nil

		assert.Equal(t, tt.expectedError, resultIsError)

		if tt.expectedError {
			continue
		}

		if len(tt.args) >= 1 {
			assert.Equal(t, tt.expectedWorkDuration, task.Duration)
		}

		if len(tt.args) == 2 {
			assert.Equal(t, tt.expectedBreakDuration, breakTask.Duration)
		}
	}
}

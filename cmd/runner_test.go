package cmd

import (
	"testing"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/spf13/cobra"
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

func TestParseFlags(t *testing.T) {
	defaultTitle := "default"

	testCases := []struct {
		name          string
		title         string
		expectedTitle string
		expectedError bool
	}{
		{
			name:          "empty title",
			title:         "",
			expectedTitle: defaultTitle, // should remain unchanged
			expectedError: false,
		},
		{
			name:          "valid title",
			title:         "Review PR",
			expectedTitle: "Review PR",
			expectedError: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a root command and set the title flag
			cmd := &cobra.Command{}
			cmd.Flags().StringP("title", "t", "", "work session title")
			_ = cmd.Flags().Set("title", tt.title)

			workTask := &config.Task{Title: defaultTitle}
			err := parseFlags(cmd, workTask)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTitle, workTask.Title)
			}
		})
	}
}

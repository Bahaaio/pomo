package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func runTask(taskType config.TaskType, cmd *cobra.Command) {
	task := taskType.GetTask()

	if err := parseArguments(cmd.Flags().Args(), task, &config.C.Break); err != nil {
		_ = cmd.Usage()
		die(err)
	}

	if err := parseFlags(cmd, &config.C.Work); err != nil {
		die(err)
	}

	log.Printf("starting %v session: %v", taskType.GetTask().Title, taskType.GetTask().Duration)

	m := ui.NewModel(taskType, config.C)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		die(err)
	}

	// print session summary
	finalModel.(ui.Model).GetSessionSummary().Print()
}

// parses the arguments and sets the duration
// returns an error if the duration is invalid
func parseArguments(args []string, task *config.Task, breakTask *config.Task) error {
	if len(args) > 0 {
		var err error
		task.Duration, err = time.ParseDuration(args[0])
		if err != nil {
			return fmt.Errorf("invalid duration: '%v'", args[0])
		}

		if len(args) > 1 {
			breakTask.Duration, err = time.ParseDuration(args[1])
			if err != nil {
				return fmt.Errorf("invalid break duration: '%v'", args[1])
			}
		}
	}

	return nil
}

// parses the flags and sets the title and audio URL
func parseFlags(cmd *cobra.Command, workTask *config.Task) error {
	title, _ := cmd.Flags().GetString("title")
	url, _ := cmd.Flags().GetString("url")

	// discard empty title
	if title != "" {
		workTask.Title = title
	}

	// override audio URL if provided
	if url != "" {
		// set URL for duringSession (global config)
		config.C.DuringSession = [][]string{
			{"mpv", "--no-video", "--loop", url},
		}
	}

	return nil
}

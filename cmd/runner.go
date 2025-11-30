package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func runTask(taskType config.TaskType, cmd *cobra.Command) {
	task := taskType.GetTask()

	if !parseArguments(cmd.Flags().Args(), task, &config.C.Break) {
		_ = cmd.Usage()
		die(nil)
	}

	log.Printf("starting %v session: %v", taskType.GetTask().Title, taskType.GetTask().Duration)

	m := ui.NewModel(taskType, config.C.ASCIIArt, config.C.AskToContinue)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		die(err)
	}

	// print session summary
	finalModel.(ui.Model).GetSessionSummary().Print()
}

// parses the arguments and sets the duration
// returns false if the duration could not be parsed
func parseArguments(args []string, task *config.Task, breakTask *config.Task) bool {
	if len(args) > 0 {
		var err error
		task.Duration, err = time.ParseDuration(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "\ninvalid duration: '%v'\n\n", args[0])
			return false
		}

		if len(args) > 1 {
			breakTask.Duration, err = time.ParseDuration(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "\ninvalid break duration: '%v'\n\n", args[1])
				return false
			}
		}
	}

	return true
}

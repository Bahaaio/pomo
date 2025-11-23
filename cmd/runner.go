package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Bahaaio/pomo/actions"
	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	"github.com/Bahaaio/pomo/ui/confirm"
	"github.com/Bahaaio/pomo/ui/summary"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var sessionSummary = summary.SessionSummary{}

func runTask(taskType config.TaskType, cmd *cobra.Command) {
	task := taskType.GetTask()

	if !parseArguments(cmd.Flags().Args(), task, &config.C.Break) {
		_ = cmd.Usage()
		die(nil)
	}

	exitStatus, elapsedTime := runTimer(task)
	log.Println("session exit status:", exitStatus)

	sessionSummary.AddSession(taskType, elapsedTime)
	wg := &sync.WaitGroup{}

	switch exitStatus {
	case ui.Quit:
		sessionSummary.Print()
		return
	case ui.Skipped:
		// skip to next task directly
	case ui.Completed:
		wg = actions.RunPostActions(task)

		if !config.C.AskToContinue || !promptToContinue(taskType) {
			wg.Wait() // wait for notification and post commands
			sessionSummary.Print()
			return
		}
	}

	wg.Wait()
	runTask(taskType.Opposite(), &cobra.Command{}) // run the next task
}

// runs the timer UI for the given task
// returns the exit status and elapsed time
func runTimer(task *config.Task) (ui.ExitStatus, time.Duration) {
	log.Printf("starting %v session: %v", task.Title, task.Duration)

	m := ui.NewModel(*task, config.C.ASCIIArt)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		die(err)
	}

	return finalModel.(ui.Model).ExitStatus(), finalModel.(ui.Model).Elapsed()
}

// prompts the user to continue to the next task
// returns true if the user confirmed
func promptToContinue(taskType config.TaskType) bool {
	prompt := fmt.Sprintf("start %s?", taskType.Opposite())

	m := confirm.New(prompt)
	p := tea.NewProgram(m, tea.WithAltScreen())

	confirmModel, err := p.Run()
	if err != nil {
		die(err)
	}

	return confirmModel.(confirm.Model).Confirmed && confirmModel.(confirm.Model).Submitted
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

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/Bahaaio/pomo/ui/confirm"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

var (
	messageStyle = lipgloss.NewStyle().Foreground(colors.SuccessMessageFg)

	totalWorkDuration time.Duration
	totalWorkSessions int

	totalBreakDuration time.Duration
	totalBreakSessions int
)

func runTask(taskType config.TaskType, cmd *cobra.Command) {
	task := taskType.GetTask()

	if !parseArguments(cmd.Flags().Args(), task, &config.C.Break) {
		_ = cmd.Usage()
		die(nil)
	}

	exitStatus, ElapsedTime := runTimer(task)
	log.Println("session exit status:", exitStatus)

	if taskType == config.WorkTask {
		totalWorkDuration += ElapsedTime
		totalWorkSessions++
	} else {
		totalBreakDuration += ElapsedTime
		totalBreakSessions++
	}

	wg := &sync.WaitGroup{}

	switch exitStatus {
	case ui.Quit:
		if totalWorkDuration > 0 || totalBreakDuration > 0 {
			printSummary()
		}
		return
	case ui.Skipped:
		// skip to next task directly
	case ui.Completed:
		wg = runPostActions(task)

		if !config.C.AskToContinue || !promptToContinue(taskType) {
			wg.Wait() // wait for notification and post commands
			printSummary()
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

// sends task notification and runs post commands using goroutines
// returns a wait group to wait for their completion
func runPostActions(task *config.Task) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sendNotification(task.Notification)
	}()

	go func() {
		defer wg.Done()
		runPostCommands(task.Then)
	}()

	return &wg
}

func sendNotification(notification config.Notification) {
	if !notification.Enabled {
		log.Println("notification disabled")
		return
	}

	log.Println("sending notification")

	// use the embeded icon
	var icon any = config.Icon

	// if the user has specified an icon
	// use that instead
	if len(notification.Icon) > 0 {
		icon = notification.Icon
	}

	var err error
	if notification.Urgent {
		err = beeep.Alert(notification.Title, notification.Message, icon)
	} else {
		err = beeep.Notify(notification.Title, notification.Message, icon)
	}

	if err != nil {
		log.Println("failed to send notification:", err)
	}
}

func runPostCommands(cmds [][]string) {
	log.Println("running post commands")

	for _, cmd := range cmds {
		c := exec.Command(cmd[0], cmd[1:]...)

		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run command '%q': %v\n", cmd, err)
		}

		// wait some time before running the next command
		time.Sleep(50 * time.Millisecond)
	}
}

func printSummary() {
	sessionIndicator := "sessions"
	if totalWorkSessions == 1 {
		sessionIndicator = "session"
	}

	fmt.Println(messageStyle.Render("Session Summary:"))

	if totalWorkDuration > 0 {
		fmt.Printf(" Work : %v (%d %s)\n", totalWorkDuration, totalWorkSessions, sessionIndicator)
	}

	if totalBreakDuration > 0 {
		fmt.Printf(" Break: %v (%d %s)\n", totalBreakDuration, totalBreakSessions, sessionIndicator)
	}

	if totalBreakDuration > 0 && totalWorkDuration > 0 {
		fmt.Println(" Total:", totalWorkDuration+totalBreakDuration)
	}

	totalDuration := totalWorkDuration + totalBreakDuration
	workRatio := float64(totalWorkDuration.Milliseconds()) / float64(totalDuration.Milliseconds())

	if totalWorkDuration > 0 {
		printProgressBar(workRatio)
	}
}

func printProgressBar(workRatio float64) {
	const barWidth = 30

	filledWidth := int(workRatio * barWidth)
	emptyWidth := barWidth - filledWidth

	bar := lipgloss.NewStyle().Foreground(colors.TimerFg).
		Render(strings.Repeat("█", filledWidth)) +
		strings.Repeat("░", emptyWidth)

	fmt.Printf("\n [%s] %.0f%% work\n", bar, workRatio*100)
}

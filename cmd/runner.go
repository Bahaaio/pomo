package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	"github.com/Bahaaio/pomo/ui/confirm"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

var messageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#198754")) // green

func runTask(taskType config.TaskType, cmd *cobra.Command) {
	task := taskType.GetTask()

	if !runTimer(task, cmd) {
		return // return if the timer is cancelled
	}

	wg := runPostActions(task)

	if !config.C.AskToContinue || !promptToContinue(taskType) {
		wg.Wait() // wait for notification and post commands
		fmt.Println(messageStyle.Render(task.Title, "finished"))
		return
	}

	wg.Wait()
	runTask(taskType.Opposite(), cmd) // run the next task
}

func runTimer(task *config.Task, cmd *cobra.Command) bool {
	if !parseArguments(cmd.Flags().Args(), task) {
		_ = cmd.Usage()
		die(nil)
	}

	log.Printf("starting %v session: %v", task.Title, task.Duration)

	m := ui.NewModel(*task)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		die(err)
	}

	if !finalModel.(ui.Model).TimerCompleted() {
		log.Println("timer did not complete")
		return false
	}

	return true
}

func promptToContinue(taskType config.TaskType) bool {
	prompt := fmt.Sprintf("start %s?", taskType.Opposite().GetTask().Title)

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
func parseArguments(args []string, task *config.Task) bool {
	if len(args) > 0 {
		var err error
		task.Duration, err = time.ParseDuration(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "\ninvalid duration: '%v'\n\n", args[0])
			return false
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

	err := beeep.Notify(notification.Title, notification.Message, icon)
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

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	os.Exit(1)
}

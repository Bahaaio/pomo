package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
)

func runTask(task config.Task) {
	log.Printf("starting %v session: %v", task.Title, task.Duration)
	notification := task.Notification

	m := ui.NewModel(task, config.C.FullScreen)
	p := tea.NewProgram(m)

	var finalModel tea.Model
	var err error

	if finalModel, err = p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// when timer is completed,
	// send the notification and run post commands
	if finalModel.(ui.Model).TimerCompleted() {
		sendNotification(notification)
		runPostCommands(task.Then)
	}
}

func sendNotification(notification config.Notification) {
	if !notification.Enabled {
		log.Println("notification disabled")
		return
	}

	log.Println("sending notification")

	err := beeep.Notify(notification.Title, notification.Message, notification.Icon)
	if err != nil {
		log.Println("failed to send notification:", err)
	}
}

func runPostCommands(cmds []string) {
	log.Println("running post commands")

	runCommand, arg := getRunCommand()

	for _, cmd := range cmds {
		c := exec.Command(runCommand, arg, cmd)

		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run command '%q': %v\n", cmd, err)
		}

		// wait some time before running the next command
		time.Sleep(50 * time.Millisecond)
	}
}

func getRunCommand() (command, arg string) {
	command, arg = "sh", "-c"

	if runtime.GOOS == "windows" {
		command, arg = "cmd.exe", "/c"
	}

	return
}

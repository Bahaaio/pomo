package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func runTask(task config.Task, taskName string) {
	log.Printf("starting %v session: %v", taskName, task.Duration)

	m := ui.NewModel(task.Duration, taskName, config.C.FullScreen)
	p := tea.NewProgram(m)

	var finalModel tea.Model
	var err error

	if finalModel, err = p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if finalModel.(ui.Model).TimerCompleted() {
		runPostCommands(task.Then)
	}
}

func runPostCommands(cmds []string) {
	log.Println("running post commands")

	for _, cmd := range cmds {
		c := exec.Command("sh", "-c", cmd)

		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run command '%q': %v\n", cmd, err)
		}
	}
}

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

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err := config.Save()
	if err != nil {
		fmt.Println("failed to save config:", err)
	}
	log.Println("saved config")

	log.Println("running post commands")
	runPostCommands(task.Then)

	log.Printf("completed %v session: %v", taskName, task.Duration)
}

func runPostCommands(cmds []string) {
	for _, cmd := range cmds {
		c := exec.Command("sh", "-c", cmd)

		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run command '%q': %v\n", cmd, err)
		}
	}
}

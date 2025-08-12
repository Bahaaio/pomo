package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func runTask(task config.Task) {
	log.Printf("starting %v session: %v", task.Title, task.Duration)

	m := ui.NewModel(task, config.C.FullScreen)
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

		// wait some time before running the next command
		time.Sleep(50 * time.Millisecond)
	}
}

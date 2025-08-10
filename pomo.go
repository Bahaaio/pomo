package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Bahaaio/pomo/config"
	tea "github.com/charmbracelet/bubbletea"
)

func runPostCommands() {
	for _, cmd := range config.C.Work.Then {
		c := exec.Command("sh", "-c", cmd)

		if err := c.Run(); err != nil {
			fmt.Printf("failed to run command '%q': %v\n", cmd, err)
		}
	}
}

func main() {
	m := NewModel(config.C.Work.Duration)
	p := tea.NewProgram(m, config.ProgramOptions()...)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err := config.Save()
	if err != nil {
		fmt.Println("failed to write config:", err)
	}
}

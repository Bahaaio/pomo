package main

import (
	"fmt"
	"io"
	"log"
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

func initializeLogging() {
	if len(os.Getenv("DEBUG")) > 0 {
		_, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("failed to setup logging:", err)
			os.Exit(1)
		}
		log.SetFlags(log.Ltime)
		log.Println("debug mode")
	} else {
		log.SetOutput(io.Discard)
	}
}

func main() {
	initializeLogging()

	err := config.LoadConfig()
	if err != nil {
		panic("could not load config")
	}

	m := NewModel(config.C.Work.Duration)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err = config.Save()
	if err != nil {
		fmt.Println("failed to write config:", err)
	}
	log.Println("saved config")
}

package config

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func ProgramOptions() []tea.ProgramOption {
	var opts []tea.ProgramOption
	if C.AltScreen {
		opts = append(opts, tea.WithAltScreen())
	}
	log.Println("AltScreen:", C.AltScreen)

	return opts
}

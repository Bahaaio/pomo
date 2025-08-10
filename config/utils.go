package config

import tea "github.com/charmbracelet/bubbletea"

func ProgramOptions() []tea.ProgramOption {
	var opts []tea.ProgramOption
	if C.AltScreen {
		opts = append(opts, tea.WithAltScreen())
	}

	return opts
}

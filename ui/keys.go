package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Increase key.Binding
	Decrease key.Binding
	Adjust   key.Binding // Combined help for increase/decrease
	Reset    key.Binding
	Pause    key.Binding
	Quit     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Adjust,
		k.Pause,
		k.Reset,
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var Keys = KeyMap{
	Increase: key.NewBinding(
		key.WithKeys("k", "up"),
	),
	Decrease: key.NewBinding(
		key.WithKeys("j", "down"),
	),
	Adjust: key.NewBinding(
		key.WithKeys("k", "up", "j", "down"),
		key.WithHelp("↑/↓", "±1 minute"),
	),
	Reset: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("", "reset"),
	),
	Pause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("󱁐", "pause/resume"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

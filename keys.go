package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Increase key.Binding
	Decrease key.Binding
	Reset    key.Binding
	Quit     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Increase,
		k.Decrease,
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
		key.WithHelp("", "+1 minute"),
	),
	Decrease: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("", "-1 minute"),
	),
	Reset: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("", "reset"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Increase key.Binding
	Reset    key.Binding
	Pause    key.Binding
	Skip     key.Binding
	Quit     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Increase,
		k.Pause,
		k.Reset,
		k.Skip,
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var Keys = KeyMap{
	Increase: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑", "+1 minute"),
	),
	Reset: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("←", "reset"),
	),
	Pause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "pause/resume"),
	),
	Skip: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "skip"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

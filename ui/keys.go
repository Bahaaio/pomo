package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Increase key.Binding
	Decrease key.Binding
	Reset    key.Binding
	Pause    key.Binding
	Skip     key.Binding
	Quit     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Increase,
		k.Decrease,
		k.Pause,
		k.Reset,
		k.Skip,
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var keyMap = KeyMap{
	Increase: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "+1 min"),
	),
	Decrease: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "-1 min"),
	),
	Reset: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("←/h", "reset"),
	),
	Pause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "pause/resume"),
	),
	Skip: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("→/l", "skip"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

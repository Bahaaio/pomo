package confirm

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Toggle  key.Binding
	Confirm key.Binding
	Cancel  key.Binding
	Submit  key.Binding
	Quit    key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Toggle,
		k.Submit,
		k.Confirm,
		k.Cancel,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var Keys = KeyMap{
	Toggle: key.NewBinding(
		key.WithKeys("tab", "h", "l", "left", "right"),
		key.WithHelp("", "toggle"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("y", "Y"),
		key.WithHelp("y", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("n", "N"),
		key.WithHelp("n", "cancel"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

package stats

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

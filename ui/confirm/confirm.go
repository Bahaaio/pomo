// Package confirm provides a simple confirmation dialog component
package confirm

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	confirmText = "Yes"
	cancelText  = "No"
)

var (
	promptStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Bold(true)

	border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5A56E0")).
		Padding(2, 7).
		BorderTop(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			Margin(0, 2)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				Padding(0, 3).
				Margin(0, 2)
)

type Model struct {
	Prompt        string
	Confirmed     bool
	Submitted     bool
	width, height int
	help          help.Model
	quitting      bool
}

func New(prompt string) Model {
	return Model{
		Prompt:    prompt,
		Confirmed: true,
		help:      help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Confirm):
			m.Confirmed = true
			m.Submitted = true
			return m, tea.Quit
		case key.Matches(msg, Keys.Cancel):
			m.Confirmed = false
			m.Submitted = true
			return m, tea.Quit
		case key.Matches(msg, Keys.Toggle):
			m.Confirmed = !m.Confirmed
			return m, nil
		case key.Matches(msg, Keys.Submit):
			m.Submitted = true
			return m, tea.Quit
		case key.Matches(msg, Keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	prompt := promptStyle.Render(m.Prompt)

	var confirmButton, cancelButton string

	if m.Confirmed {
		confirmButton = activeButtonStyle.Render(confirmText)
		cancelButton = buttonStyle.Render(cancelText)
	} else {
		confirmButton = buttonStyle.Render(confirmText)
		cancelButton = activeButtonStyle.Render(cancelText)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Right, confirmButton, cancelButton)
	dialog := lipgloss.JoinVertical(lipgloss.Center, prompt, "\n", buttons)
	ui := border.Render(dialog)
	help := m.help.View(Keys)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			ui,
			"",
			help,
		),
	)
}

// Package confirm provides a simple confirmation dialog component
package confirm

import (
	"github.com/Bahaaio/pomo/ui/colors"
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
	buttonPadding = []int{0, 3}
	buttonMargin  = []int{0, 2}
	borderPadding = []int{2, 6}
)

var (
	promptStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderFg).
			Padding(borderPadding...).
			BorderTop(true)

	InactiveButtonStyle = lipgloss.NewStyle().
				Foreground(colors.InactiveButtonFg).
				Background(colors.InactiveButtonBg).
				Padding(buttonPadding...).
				Margin(buttonMargin...)

	activeButtonStyle = InactiveButtonStyle.
				Foreground(colors.ActiveButtonFg).
				Background(colors.ActiveButtonBg).
				Padding(buttonPadding...).
				Margin(buttonMargin...)
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
		cancelButton = InactiveButtonStyle.Render(cancelText)
	} else {
		confirmButton = InactiveButtonStyle.Render(confirmText)
		cancelButton = activeButtonStyle.Render(cancelText)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Right, confirmButton, cancelButton)
	dialog := lipgloss.JoinVertical(lipgloss.Center, prompt, "\n", buttons)
	ui := borderStyle.Render(dialog)
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

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

type ConfirmChoice int

const (
	Confirm ConfirmChoice = iota
	Cancel
)

type ChoiceMsg struct {
	Choice ConfirmChoice
}

type Model struct {
	confirmed     bool
	width, height int
	help          help.Model
	quitting      bool
}

func New() Model {
	return Model{
		confirmed: true,
		help:      help.New(),
	}
}

func (m Model) View(prompt string) string {
	if m.quitting {
		return ""
	}

	prompt = promptStyle.Render(prompt)

	var confirmButton, cancelButton string

	if m.confirmed {
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

func (m *Model) HandleKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Keys.Confirm):
		return m.Choice(Confirm)

	case key.Matches(msg, Keys.Cancel):
		return m.Choice(Cancel)

	case key.Matches(msg, Keys.Toggle):
		m.confirmed = !m.confirmed
		return nil

	case key.Matches(msg, Keys.Submit):
		if m.confirmed {
			return m.Choice(Confirm)
		}
		return m.Choice(Cancel)

	case key.Matches(msg, Keys.Quit):
		m.quitting = true
		return m.Choice(Cancel)

	default:
		return nil
	}
}

func (m *Model) HandleWindowResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.width = msg.Width
	m.height = msg.Height
	return nil
}

func (m Model) Choice(choice ConfirmChoice) tea.Cmd {
	return func() tea.Msg {
		return ChoiceMsg{Choice: choice}
	}
}

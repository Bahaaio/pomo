// Package stats implements the statistics view for pomo.
package stats

import (
	"fmt"

	"github.com/Bahaaio/pomo/db"
	"github.com/Bahaaio/pomo/ui/stats/components"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const barChartHeight = 12

type Model struct {
	// components
	barChart components.BarChart

	// stats
	allTimeStats db.AllTimeStats
	weeklyStats  []db.DailyStat
	monthlyStats []db.DailyStat

	// state
	width, height int
	help          help.Model
	quitting      bool
}

func New() Model {
	return Model{
		barChart: components.NewBarChart(barChartHeight),
		help:     help.New(),
	}
}

type statsMsg struct {
	allTimeStats db.AllTimeStats
	weeklyStats  []db.DailyStat
	monthlyStats []db.DailyStat
}

func fetchStats() tea.Msg {
	// TODO: remove panics and replace with error message

	database, err := db.Init()
	if err != nil {
		panic(err)
	}

	repo := db.NewSessionRepo(database)

	stats, err := repo.GetAllTimeStats()
	if err != nil {
		panic(err)
	}

	weeklyStats, err := repo.GetWeeklyStats()
	if err != nil {
		panic(err)
	}

	monthlyStats, err := repo.GetMonthlyStats()
	if err != nil {
		panic(err)
	}

	return statsMsg{
		allTimeStats: stats,
		weeklyStats:  weeklyStats,
		monthlyStats: monthlyStats,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchStats
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	title := "Pomodoro statistics\n\n"

	content := "All-Time stats:\n"
	content += fmt.Sprintln("  Total Sessions:", m.allTimeStats.TotalSessions)
	content += fmt.Sprintln("  Work Time:", m.allTimeStats.TotalWorkDuration)
	content += fmt.Sprintln("  Break Time:", m.allTimeStats.TotalBreakDuration)

	// left align
	content = lipgloss.JoinVertical(lipgloss.Left, content, "")

	chart := m.barChart.View(m.weeklyStats)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			content,
			chart,
			"",
			m.help.View(Keys),
		),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsMsg:
		m.allTimeStats = msg.allTimeStats
		m.weeklyStats = msg.weeklyStats
		m.monthlyStats = msg.monthlyStats
		return m, nil
	case tea.KeyMsg:
		return m, handleKeys(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	default:
		return m, nil
	}
}

func handleKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Keys.Quit):
		return tea.Quit
	}
	return nil
}

// Package stats implements the statistics view for pomo.
package stats

import (
	"errors"

	"github.com/Bahaaio/pomo/db"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/Bahaaio/pomo/ui/stats/components"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	barChartHeight     = 12
	durationRatioWidth = 30
)

var errStyle = lipgloss.NewStyle().
	Foreground(colors.ErrorMessageFg).
	AlignHorizontal(lipgloss.Center)

type Model struct {
	// components
	durationRatio components.DurationRatio
	barChart      components.BarChart
	heatMap       components.HeatMap
	streak        components.Streak

	// error message
	err error

	// stats
	allTimeStats db.AllTimeStats
	weeklyStats  []db.DailyStat
	monthlyStats []db.DailyStat
	streakStats  db.StreakStats

	// state
	width, height int
	help          help.Model
	quitting      bool
}

func New() Model {
	return Model{
		durationRatio: components.NewDurationRatio(durationRatioWidth),
		barChart:      components.NewBarChart(barChartHeight),
		heatMap:       components.NewHeatMap(),
		streak:        components.NewStreak(),
		help:          help.New(),
	}
}

type statsMsg struct {
	allTimeStats db.AllTimeStats
	weeklyStats  []db.DailyStat
	monthlyStats []db.DailyStat
	streakStats  db.StreakStats
}

type errMsg struct {
	err error
}

// fetchStats retrieves statistics from the database and returns them as a statsMsg.
// If an error occurs, it returns an errMsg instead.
func fetchStats() tea.Msg {
	database, err := db.Connect()
	if err != nil {
		return errMsg{err: errors.New("failed to connect to the database")}
	}

	repo := db.NewSessionRepo(database)

	stats, err := repo.GetAllTimeStats()
	if err != nil {
		return errMsg{err: errors.New("failed to fetch all-time stats")}
	}

	weeklyStats, err := repo.GetWeeklyStats()
	if err != nil {
		return errMsg{err: errors.New("failed to fetch weekly stats")}
	}

	monthlyStats, err := repo.GetLastMonthsStats(components.NumberOfMonths)
	if err != nil {
		return errMsg{err: errors.New("failed to fetch heatmap stats")}
	}

	streakStats, err := repo.GetStreakStats()
	if err != nil {
		return errMsg{err: errors.New("failed to fetch streak stats")}
	}

	return statsMsg{
		allTimeStats: stats,
		weeklyStats:  weeklyStats,
		monthlyStats: monthlyStats,
		streakStats:  streakStats,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchStats
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.err != nil {
		return m.buildErrorMessage()
	}

	title := "Pomodoro statistics"

	durationRatio := m.durationRatio.View(
		m.allTimeStats.TotalWorkDuration,
		m.allTimeStats.TotalBreakDuration,
	)

	streak := m.streak.View(m.streakStats)

	chart := m.barChart.View(m.weeklyStats)
	hMap := m.heatMap.View(m.monthlyStats)

	charts := lipgloss.JoinHorizontal(lipgloss.Bottom, chart, "   ", hMap)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			"\n\n",
			durationRatio,
			"",
			streak,
			"\n",
			charts,
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
		m.streakStats = msg.streakStats
		return m, nil
	case errMsg:
		m.err = msg.err
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

func (m *Model) buildErrorMessage() string {
	title := "An error occurred while fetching statistics."
	message := m.err.Error()

	help := m.help.View(Keys)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		message,
		"",
		help,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		errStyle.Render(content),
	)
}

func handleKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Keys.Quit):
		return tea.Quit
	}
	return nil
}

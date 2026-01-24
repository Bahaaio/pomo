package components

import (
	"strings"
	"time"

	"github.com/Bahaaio/pomo/db"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/lipgloss"
)

const (
	NumberOfMonths = 4

	cellChar            = "󱓻 "
	cellWidth           = 2
	verticalSeparator   = "─"
	horizontalSeparator = "│"

	labelWidth        = 3
	weekDayLabelWidth = labelWidth + 1 + 1 + cellWidth // including separator, space, and cellWidth*space "Sun │  "
	maxMonthWeeks     = 5
	monthGridWidth    = cellWidth * maxMonthWeeks
	gridWidth         = monthGridWidth*NumberOfMonths + (NumberOfMonths)*cellWidth
)

var (
	style0 = lipgloss.NewStyle().Foreground(colors.HeatMapFg0)
	style1 = lipgloss.NewStyle().Foreground(colors.HeatMapFg1)
	style2 = lipgloss.NewStyle().Foreground(colors.HeatMapFg2)
	style3 = lipgloss.NewStyle().Foreground(colors.HeatMapFg3)
	style4 = lipgloss.NewStyle().Foreground(colors.HeatMapFg4)
	styles = []lipgloss.Style{style0, style1, style2, style3, style4}

	paddingStyle   = lipgloss.NewStyle().Padding(1)
	leftAlignStyle = lipgloss.NewStyle().Align(lipgloss.Left)
)

type monthGrid struct {
	label    string
	numWeeks int
	cells    [7][]string
}

type HeatMap struct{}

func NewHeatMap() HeatMap {
	return HeatMap{}
}

func (h *HeatMap) View(stats []db.DailyStat) string {
	statsMap := buildStatsMap(stats)
	grids := h.makeMonthGrids(statsMap)

	// left align month labels
	monthLabels := h.buildMonthLabels(grids)
	monthLabels = leftAlignStyle.Render(monthLabels)

	dayLabels := h.buildWeekDayLabels()
	grid := h.buildGrids(grids)
	legend := h.buildLegend()

	center := lipgloss.JoinHorizontal(lipgloss.Left, dayLabels, grid)
	body := lipgloss.JoinVertical(lipgloss.Left, monthLabels, center)
	heatMap := lipgloss.JoinVertical(lipgloss.Center, body, "", legend)

	return paddingStyle.Render(heatMap)
}

func (h *HeatMap) buildLegend() string {
	builder := strings.Builder{}
	builder.WriteString("Less ")

	for _, style := range styles {
		builder.WriteString(style.Render(cellChar))
	}

	builder.WriteString(" More")
	return builder.String()
}

func buildStatsMap(stats []db.DailyStat) map[string]time.Duration {
	m := make(map[string]time.Duration)
	for _, stat := range stats {
		m[stat.Date] = stat.WorkDuration
	}
	return m
}

func (h *HeatMap) buildGrids(grids []monthGrid) string {
	var rows [7]strings.Builder

	for monthIdx, grid := range grids {
		for weekday := range 7 {
			for _, cell := range grid.cells[weekday] {
				rows[weekday].WriteString(cell)
			}
			// add separator between months (except last)
			if monthIdx < len(grids)-1 {
				rows[weekday].WriteString(strings.Repeat(" ", cellWidth))
			}
		}
	}

	// convert to string slice and join with newlines
	var result []string
	for i := range rows {
		result = append(result, rows[i].String())
	}

	return strings.Join(result, "\n")
}

func (h *HeatMap) makeMonthGrids(statsMap map[string]time.Duration) []monthGrid {
	now := time.Now()
	var grids []monthGrid

	// build grid for each of the last N months
	for i := NumberOfMonths - 1; i >= 0; i-- {
		monthTime := now.AddDate(0, -i, 0)
		grid := h.makeMonthGrid(monthTime.Year(), monthTime.Month(), now, statsMap)
		grids = append(grids, grid)
	}

	return grids
}

func (h *HeatMap) makeMonthGrid(year int, month time.Month, today time.Time, statsMap map[string]time.Duration) monthGrid {
	// get first and last day of month
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1) // last day of month

	// calculate number of weeks this month spans
	startWeekday := int(firstDay.Weekday()) // 0=Sun, 6=Sat
	daysInMonth := lastDay.Day()
	numWeeks := (startWeekday + daysInMonth + 6) / 7

	// initialize grid with empty cells
	var cells [7][]string
	for i := range cells {
		cells[i] = make([]string, numWeeks)
		for j := range cells[i] {
			cells[i][j] = strings.Repeat(" ", cellWidth) // empty
		}
	}

	// fill in actual days
	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

		// skip future dates
		if date.After(today) {
			continue
		}

		weekday := int(date.Weekday())
		week := (startWeekday + day - 1) / 7

		key := date.Format(db.DateFormat)
		duration := statsMap[key]
		cells[weekday][week] = renderCell(duration)
	}

	return monthGrid{
		label:    month.String()[:3],
		numWeeks: numWeeks,
		cells:    cells,
	}
}

func (h *HeatMap) buildMonthLabels(grids []monthGrid) string {
	var builder strings.Builder

	// add left padding to align with day labels
	builder.WriteString(strings.Repeat(" ", weekDayLabelWidth))

	for i, grid := range grids {
		// center the label over the grid
		leftPad := (monthGridWidth - labelWidth) / 2
		rightPad := monthGridWidth - labelWidth - leftPad + 1

		builder.WriteString(strings.Repeat(" ", leftPad))
		builder.WriteString(grid.label)
		builder.WriteString(strings.Repeat(" ", rightPad))

		// add separator spacing between months (except last)
		if i < len(grids)-1 {
			builder.WriteString(strings.Repeat(" ", cellWidth))
		}
	}

	// add separator line
	separator := strings.Repeat(verticalSeparator, weekDayLabelWidth+gridWidth)

	builder.WriteString("\n")
	builder.WriteString(separator)

	return builder.String()
}

func (h *HeatMap) buildWeekDayLabels() string {
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	builder := strings.Builder{}

	for _, day := range days {
		builder.WriteString(day + " " + horizontalSeparator + strings.Repeat(" ", cellWidth))
		if day != "Sat" {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func renderCell(duration time.Duration) string {
	style := getCellStyle(duration)
	return style.Render(cellChar)
}

func getCellStyle(duration time.Duration) lipgloss.Style {
	if duration < time.Second {
		return style0
	} else if duration <= time.Minute*30 {
		return style1
	} else if duration <= time.Hour {
		return style2
	} else if duration <= time.Hour*2 {
		return style3
	} else {
		return style4
	}
}

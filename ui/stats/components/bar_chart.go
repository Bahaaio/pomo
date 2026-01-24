// Package components implements UI components for stats.
package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bahaaio/pomo/db"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/lipgloss"
)

var barStyle = lipgloss.NewStyle().Foreground(colors.WorkSessionFg)

const (
	barChar     = "█"
	axisChar    = "│"
	tickChar    = "┤"
	cornerChar  = "└"
	lineChar    = "─"
	paddingChar = " "

	barThickness = 3
	spacing      = 2
	daysInWeek   = 7
)

var spacer = strings.Repeat(paddingChar, spacing)

type chartLayout struct {
	barHeight       int
	yAxisLabelWidth int
	yAxisWidth      int // label + space + tick char
	barAreaWidth    int
	totalWidth      int
}

type BarChart struct {
	chartLayout
}

func NewBarChart(height int) BarChart {
	return BarChart{
		chartLayout: chartLayout{
			barHeight: height - 1 - 1, // leave space for x-axis and labels
		},
	}
}

func (b *BarChart) calculateLayout(maxDuration, scale time.Duration) chartLayout {
	longestLabel := 0

	for duration := maxDuration; duration > 0; duration -= scale {
		label := formatDuration(duration)
		longestLabel = max(longestLabel, len(label))
	}

	yAxisLabelWidth := longestLabel
	yAxisWidth := yAxisLabelWidth + 1 + 1 // length of label + space + tick char

	barAreaWidth := spacing + (barThickness+spacing)*daysInWeek

	return chartLayout{
		barHeight:       b.barHeight,
		yAxisLabelWidth: yAxisLabelWidth,
		yAxisWidth:      yAxisWidth,
		barAreaWidth:    barAreaWidth,
		totalWidth:      yAxisWidth + barAreaWidth,
	}
}

func (b *BarChart) View(stats []db.DailyStat) string {
	if len(stats) == 0 {
		return ""
	}

	maxDuration := getMaxDuration(stats)

	// dividing by half of max height to leave space for tick chars
	targetTicks := b.barHeight / 2
	scale := calculateScale(maxDuration, targetTicks)

	// fallback for empty stats
	if maxDuration == 0 {
		maxDuration = time.Hour
		scale = time.Minute * 10
	}

	b.chartLayout = b.calculateLayout(maxDuration, scale)

	yAxis := b.buildYAxis(maxDuration, scale)
	bars := b.buildBars(stats, maxDuration)

	top := lipgloss.JoinHorizontal(lipgloss.Left, yAxis, spacer, bars)
	xAxis := b.buildXAxis()
	labels := b.buildLabels(stats)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		xAxis,
		labels,
	)
}

func (b *BarChart) buildBars(stats []db.DailyStat, maxDuration time.Duration) string {
	bars := make([]string, 0, len(stats))

	for _, stat := range stats {
		if stat.WorkDuration == 0 {
			bars = append(bars, renderBar(0), spacer)
			continue
		}

		barHeight := int((float64(stat.WorkDuration) / float64(maxDuration)) * float64(b.barHeight))
		bar := renderBar(barHeight)
		bars = append(bars, bar, spacer)
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, bars...)
}

func (b *BarChart) buildYAxis(maxDuration, scale time.Duration) string {
	builder := strings.Builder{}

	// smallest duration to print
	epsilon := time.Millisecond * 500

	// print all rows with same width of longest duration label
	for duration := maxDuration; duration >= epsilon; duration -= scale {
		tick := fmt.Sprintf("%-*s %s\n", b.yAxisLabelWidth, formatDuration(duration), tickChar)
		builder.WriteString(tick)

		axis := strings.Repeat(paddingChar, b.yAxisLabelWidth) + paddingChar + axisChar
		builder.WriteString(axis)

		// don't print the last new line
		if duration-scale >= epsilon {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (b *BarChart) buildXAxis() string {
	zeroLabel := fmt.Sprintf("%-*s", b.yAxisLabelWidth, "0")
	return zeroLabel + " " + cornerChar + strings.Repeat(lineChar, b.barAreaWidth)
}

func (b *BarChart) buildLabels(stats []db.DailyStat) string {
	var labels strings.Builder

	for _, stat := range stats {
		day := getDayLabel(stat.Date)
		labels.WriteString(day)
		labels.WriteString(spacer)
	}

	// yaxis width + spacing between yaxis and bars
	paddingLength := b.yAxisWidth + spacing
	padding := strings.Repeat(paddingChar, paddingLength)

	return padding + labels.String()
}

func getDayLabel(day string) string {
	t, err := time.Parse(db.DateFormat, day)
	if err != nil {
		return strings.Repeat(paddingChar, barThickness)
	}

	// get first three letters of weekday
	return t.Weekday().String()[:barThickness]
}

func renderBar(height int) string {
	if height == 0 {
		// return an empty bar for alignment
		return strings.Repeat(paddingChar, barThickness)
	}

	bar := strings.Repeat(barChar, barThickness)

	// don't render last new line char
	return barStyle.Render(strings.Repeat(bar+"\n", height-1) + bar)
}

func getMaxDuration(stats []db.DailyStat) time.Duration {
	var maxDuration time.Duration

	for _, stat := range stats {
		if stat.WorkDuration > maxDuration {
			maxDuration = stat.WorkDuration
		}
	}

	return maxDuration
}

func calculateScale(maxDuration time.Duration, targetTicks int) time.Duration {
	scale := time.Duration(
		float64(maxDuration.Milliseconds())/float64(targetTicks),
	) * time.Millisecond

	// minimum scale of 100 ms to avoid too many ticks
	scale = max(scale, time.Millisecond*100)

	return scale
}

func formatDuration(d time.Duration) string {
	seconds := d.Seconds()
	if seconds < 60 {
		if seconds == float64(int(seconds)) {
			return fmt.Sprintf("%ds", int(seconds))
		}

		// show one decimal place for seconds less than 1
		return fmt.Sprintf("%0.1fs", seconds)
	}

	minutes := int(d.Minutes())
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	mins := minutes % 60

	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dh%dm", hours, mins)
}

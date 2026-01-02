package components

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/lipgloss"
)

type DurationRatio struct {
	width int
}

func NewDurationRatio(width int) DurationRatio {
	return DurationRatio{
		width: width,
	}
}

func (d *DurationRatio) View(workDuration, breakDuration time.Duration) string {
	totalDuration := workDuration + breakDuration

	if totalDuration == 0 {
		return ""
	}

	workRatio := float64(workDuration.Milliseconds()) / float64(totalDuration.Milliseconds())

	workPercentage := int(math.Round(workRatio * 100))
	breakPercentage := 100 - workPercentage

	top := d.buildTop(workDuration, breakDuration)
	bar := d.buildBar(workPercentage)
	bottom := d.buildBottom(workPercentage, breakPercentage)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		bar,
		bottom,
	)
}

func (d *DurationRatio) buildTop(workDuration, breakDuration time.Duration) string {
	workLabel := workDuration.String()
	breakLabel := breakDuration.String()

	paddingLength := max(d.width-len(workLabel)-len(breakLabel), 0)
	padding := strings.Repeat(" ", paddingLength)

	return workLabel + padding + breakLabel
}

func (d *DurationRatio) buildBar(workPercentage int) string {
	filledWidth := int(float64(d.width) * (float64(workPercentage) / 100.0))
	emptyWidth := d.width - filledWidth

	bar := lipgloss.NewStyle().Foreground(colors.WorkSessionFG).
		Render(strings.Repeat("█", filledWidth)) +
		strings.Repeat("░", emptyWidth)

	return bar
}

func (d *DurationRatio) buildBottom(workPercentage, breakPercentage int) string {
	workLabel := fmt.Sprintf("%d%%", workPercentage)
	breakLabel := fmt.Sprintf("%d%%", breakPercentage)

	paddingLength := max(d.width-len(workLabel)-len(breakLabel), 0)
	padding := strings.Repeat(" ", paddingLength)

	return workLabel + padding + breakLabel
}

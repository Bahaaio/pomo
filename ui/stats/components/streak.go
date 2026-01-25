package components

import (
	"fmt"

	"github.com/Bahaaio/pomo/db"
)

type Streak struct{}

func NewStreak() Streak {
	return Streak{}
}

func (s Streak) View(streak db.StreakStats) string {
	return fmt.Sprintf("󱐋 streak %vd · best %vd", streak.Current, streak.Best)
}

package db

import (
	"time"
)

// calculateStreak calculates the current and best streaks from a list of dates (formatted as "YYYY-MM-DD").
// the dates should be sorted in descending order (most recent first).
// a streak is consecutive days with at least one 'work' session.
func calculateStreak(dates []string) StreakStats {
	now := time.Now()
	today := now.Format(DateFormat)
	yesterday := now.AddDate(0, 0, -1).Format(DateFormat)

	currentStreak, bestStreak, tempStreak := 0, 0, 0
	currentStreakBroken := false

	for i, date := range dates {
		if i == 0 {
			// if the most recent date is today or yesterday, start the streak
			if date == today || date == yesterday {
				tempStreak = 1
				currentStreak = 1
				bestStreak = 1
			} else {
				// streak is already broken
				currentStreakBroken = true
				tempStreak = 1
			}
			continue
		}

		if isConsecutiveDate(dates[i-1], date) {
			tempStreak++
			bestStreak = max(bestStreak, tempStreak)

			if !currentStreakBroken {
				currentStreak++
			}
		} else {
			// streak broken
			currentStreakBroken = true
			tempStreak = 1
		}
	}

	return StreakStats{Current: currentStreak, Best: bestStreak}
}

// checks if date1 is the day after date2
func isConsecutiveDate(date1, date2 string) bool {
	t1, err1 := time.Parse(DateFormat, date1)
	t2, err2 := time.Parse(DateFormat, date2)

	if err1 != nil || err2 != nil {
		return false
	}

	return t1.AddDate(0, 0, -1).Equal(t2)
}

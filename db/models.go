package db

import "time"

var schema = `
CREATE TABLE IF NOT EXISTS sessions(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type TEXT NOT NULL,
	duration INTEGER NOT NULL,
	started_at TEXT NOT NULL
);
`

type Session struct {
	ID        int           `db:"id"`
	Type      string        `db:"type"`
	Duration  time.Duration `db:"duration"`
	StartedAt time.Time     `db:"started_at"`
}

type AllTimeStats struct {
	TotalSessions      int           `db:"total_sessions"`
	TotalWorkDuration  time.Duration `db:"total_work_duration"`
	TotalBreakDuration time.Duration `db:"total_break_duration"`
}

type DailyStat struct {
	Date         string        `db:"day"`
	WorkDuration time.Duration `db:"work_duration"`
}

type SessionType string

const (
	WorkSession  SessionType = "work"
	BreakSession SessionType = "break"
)

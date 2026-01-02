package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const DateFormat = "2006-01-02"

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

// CreateSession inserts a new session record into the database.
func (r *SessionRepo) CreateSession(startedAt time.Time, duration time.Duration, sessionType SessionType) error {
	startedAtStr := startedAt.Format(time.RFC3339)

	if _, err := r.db.Exec(
		"insert into sessions (started_at, duration, type) values (?, ?, ?);",
		startedAtStr,
		duration,
		sessionType,
	); err != nil {
		return err
	}

	return nil
}

// GetAllTimeStats retrieves aggregate statistics across all sessions.
func (r *SessionRepo) GetAllTimeStats() (AllTimeStats, error) {
	var totalStats AllTimeStats

	// sqlite treats (type = 'work') as 1 or 0
	if err := r.db.Get(
		&totalStats,
		`
		SELECT
			COUNT(*) AS total_sessions,
			COALESCE(SUM(duration * (type = 'work')), 0)  AS total_work_duration,
			COALESCE(SUM(duration * (type = 'break')), 0) AS total_break_duration
		FROM sessions;
		`,
	); err != nil {
		return AllTimeStats{}, err
	}

	return totalStats, nil
}

func (r *SessionRepo) GetWeeklyStats() ([]DailyStat, error) {
	today := time.Now()
	firstDay := today.AddDate(0, 0, -6)

	return r.getDailyStats(firstDay, today)
}

func (r *SessionRepo) GetMonthlyStats() ([]DailyStat, error) {
	today := time.Now()
	firstDay := today.AddDate(0, 0, -29)

	return r.getDailyStats(firstDay, today)
}

func (r *SessionRepo) getDailyStats(from, to time.Time) ([]DailyStat, error) {
	fromStr := from.Format(DateFormat)
	toStr := to.Format(DateFormat)

	var stats []DailyStat

	if err := r.db.Select(
		&stats,
		`
		SELECT
			date(started_at) AS day,
			COALESCE(SUM(duration * (type = 'work')), 0) AS work_duration
		FROM sessions
		WHERE date(started_at) BETWEEN ? AND ?
		GROUP BY day
		ORDER BY day;
		`,
		fromStr, toStr,
	); err != nil {
		return nil, err
	}

	return r.normalizeStats(from, to, stats), nil
}

func (r *SessionRepo) normalizeStats(from, to time.Time, stats []DailyStat) []DailyStat {
	m := make(map[string]DailyStat)

	for _, stat := range stats {
		m[stat.Date] = stat
	}

	var normalized []DailyStat
	current := from
	for !current.After(to) {
		day := current.Format(DateFormat)

		normalized = append(normalized, DailyStat{
			Date:         day,
			WorkDuration: m[day].WorkDuration,
		})

		current = current.AddDate(0, 0, 1) // next day
	}

	return normalized
}

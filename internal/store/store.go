// Package store keeps the user's submissions in a local SQLite database.
package store

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

type Submission struct {
	ID          int64
	Exercise    string // exercise key, e.g. "1001-extremamente-basico"
	Language    string
	Verdict     string
	Passed      int
	Total       int
	TimeMs      int64
	SamplesOnly bool // "testar exemplos": shown to the user, but never counts as solved
	CreatedAt   time.Time
}

type ExerciseStat struct {
	Solved      bool
	Attempts    int
	LastVerdict string
	LastAt      time.Time
}

const schema = `
CREATE TABLE IF NOT EXISTS submissions (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	exercise     TEXT NOT NULL,
	language     TEXT NOT NULL,
	verdict      TEXT NOT NULL,
	passed       INTEGER NOT NULL,
	total        INTEGER NOT NULL,
	time_ms      INTEGER NOT NULL,
	samples_only INTEGER NOT NULL DEFAULT 0,
	created_at   TEXT NOT NULL,
	code         TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS submissions_exercise ON submissions (exercise);
`

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Add(sub Submission, code string) (int64, error) {
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = time.Now()
	}
	res, err := s.db.Exec(`INSERT INTO submissions
		(exercise, language, verdict, passed, total, time_ms, samples_only, created_at, code)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sub.Exercise, sub.Language, sub.Verdict, sub.Passed, sub.Total, sub.TimeMs,
		sub.SamplesOnly, sub.CreatedAt.UTC().Format(time.RFC3339), code)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Submissions returns real submissions (not sample runs), newest first.
func (s *Store) Submissions(limit int) ([]Submission, error) {
	rows, err := s.db.Query(`SELECT id, exercise, language, verdict, passed, total, time_ms, created_at
		FROM submissions WHERE samples_only = 0 ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []Submission
	for rows.Next() {
		var sub Submission
		var at string
		if err := rows.Scan(&sub.ID, &sub.Exercise, &sub.Language, &sub.Verdict, &sub.Passed, &sub.Total, &sub.TimeMs, &at); err != nil {
			return nil, err
		}
		sub.CreatedAt, _ = time.Parse(time.RFC3339, at)
		subs = append(subs, sub)
	}
	return subs, rows.Err()
}

// Stats summarizes each exercise the user has submitted.
func (s *Store) Stats() (map[string]ExerciseStat, error) {
	subs, err := s.Submissions(1 << 30)
	if err != nil {
		return nil, err
	}
	stats := map[string]ExerciseStat{}
	for i := len(subs) - 1; i >= 0; i-- { // oldest first, so the last one wins
		sub := subs[i]
		st := stats[sub.Exercise]
		st.Attempts++
		st.LastVerdict, st.LastAt = sub.Verdict, sub.CreatedAt
		st.Solved = st.Solved || sub.Verdict == "AC"
		stats[sub.Exercise] = st
	}
	return stats, nil
}

// Activity counts submissions per local day ("2006-01-02").
func (s *Store) Activity() (map[string]int, error) {
	subs, err := s.Submissions(1 << 30)
	if err != nil {
		return nil, err
	}
	days := map[string]int{}
	for _, sub := range subs {
		days[sub.CreatedAt.Local().Format("2006-01-02")]++
	}
	return days, nil
}

// BestStreak is the longest run of consecutive days with submissions.
func BestStreak(activity map[string]int) int {
	best := 0
	for day := range activity {
		d, err := time.ParseInLocation("2006-01-02", day, time.Local)
		if err != nil || activity[d.AddDate(0, 0, -1).Format("2006-01-02")] > 0 {
			continue // not the first day of a run
		}
		n := 0
		for activity[d.Format("2006-01-02")] > 0 {
			n++
			d = d.AddDate(0, 0, 1)
		}
		best = max(best, n)
	}
	return best
}

// Streak counts consecutive days with submissions ending today (or
// yesterday, so the streak survives until the end of today).
func Streak(activity map[string]int, now time.Time) int {
	day := now
	if activity[day.Format("2006-01-02")] == 0 {
		day = day.AddDate(0, 0, -1)
	}
	n := 0
	for activity[day.Format("2006-01-02")] > 0 {
		n++
		day = day.AddDate(0, 0, -1)
	}
	return n
}

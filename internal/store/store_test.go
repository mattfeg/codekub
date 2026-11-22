package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStatsAndActivity(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	day := func(d int) time.Time { return time.Date(2026, 11, d, 12, 0, 0, 0, time.Local) }
	add := func(ex, v string, d int, samples bool) {
		if _, err := s.Add(Submission{Exercise: ex, Language: "ruby", Verdict: v, CreatedAt: day(d), SamplesOnly: samples}, "code"); err != nil {
			t.Fatal(err)
		}
	}
	add("001-hello-world", "WA", 20, false)
	add("001-hello-world", "AC", 21, false)
	add("002-business-card", "AC", 22, true) // sample run: never counts
	add("002-business-card", "PE", 22, false)

	stats, _ := s.Stats()
	if st := stats["001-hello-world"]; !st.Solved || st.Attempts != 2 || st.LastVerdict != "AC" {
		t.Errorf("hello world: %+v", st)
	}
	if st := stats["002-business-card"]; st.Solved || st.Attempts != 1 || st.LastVerdict != "PE" {
		t.Errorf("business card: %+v", st)
	}
	act, _ := s.Activity()
	if act["2026-11-22"] != 1 || act["2026-11-20"] != 1 {
		t.Errorf("activity %v", act)
	}
	if n := Streak(act, day(22)); n != 3 {
		t.Errorf("streak %d, want 3", n)
	}
	if n := Streak(act, day(23)); n != 3 {
		t.Errorf("streak should survive until the end of the next day, got %d", n)
	}
	if n := BestStreak(act); n != 3 {
		t.Errorf("best streak %d, want 3", n)
	}
	if n := Streak(act, day(25)); n != 0 {
		t.Errorf("broken streak should be 0, got %d", n)
	}
}

// Package core wires the track, the user's workspace, the judge and the
// progress store together. Both the CLI commands and the TUI use it.
package core

import (
	"context"
	"os"
	"os/user"
	"strings"

	"github.com/mattfeg/codekub/internal/judge"
	"github.com/mattfeg/codekub/internal/store"
	"github.com/mattfeg/codekub/internal/track"
	"github.com/mattfeg/codekub/internal/workspace"
)

// UnlockRatio is the share of a phase that must be solved to open the next.
const UnlockRatio = 0.7

type Env struct {
	Track  *track.Track
	WS     *workspace.Workspace
	Store  *store.Store
	Lang   judge.Language
	Runner judge.Runner
}

func Open() (*Env, error) {
	fsys, source, err := track.Find()
	if err != nil {
		return nil, err
	}
	tr, err := track.Load(fsys, source)
	if err != nil {
		return nil, err
	}
	ws, err := workspace.Default()
	if err != nil {
		return nil, err
	}
	lang, err := ws.Language()
	if err != nil {
		return nil, err
	}
	st, err := store.Open(ws.DBPath())
	if err != nil {
		return nil, err
	}
	return &Env{Track: tr, WS: ws, Store: st, Lang: lang, Runner: judge.LocalRunner{}}, nil
}

func (e *Env) Close() { e.Store.Close() }

// Stats never fails: a broken database shows as no progress.
func (e *Env) Stats() map[string]store.ExerciseStat {
	stats, err := e.Store.Stats()
	if err != nil {
		return map[string]store.ExerciseStat{}
	}
	return stats
}

func Solved(p *track.Phase, stats map[string]store.ExerciseStat) int {
	n := 0
	for _, ex := range p.Exercises {
		if stats[ex.Key()].Solved {
			n++
		}
	}
	return n
}

// Unlocked reports whether phase i is open: the first always is, the others
// once UnlockRatio of the previous phase is solved.
func (e *Env) Unlocked(i int, stats map[string]store.ExerciseStat) bool {
	if i == 0 {
		return true
	}
	prev := e.Track.Phases[i-1]
	return float64(Solved(prev, stats)) >= UnlockRatio*float64(len(prev.Exercises))
}

// Next is the first unsolved exercise in an unlocked phase, or nil when the
// open part of the track is complete.
func (e *Env) Next(stats map[string]store.ExerciseStat) *track.Exercise {
	for i, p := range e.Track.Phases {
		if !e.Unlocked(i, stats) {
			return nil
		}
		for _, ex := range p.Exercises {
			if !stats[ex.Key()].Solved {
				return ex
			}
		}
	}
	return nil
}

// Continue picks what the home screen offers: the latest attempted exercise
// if it is still unsolved, otherwise the next one.
func (e *Env) Continue(stats map[string]store.ExerciseStat) *track.Exercise {
	var latest *track.Exercise
	for _, ex := range e.Track.Exercises() {
		st, ok := stats[ex.Key()]
		if ok && !st.Solved && (latest == nil || st.LastAt.After(stats[latest.Key()].LastAt)) {
			latest = ex
		}
	}
	if latest != nil {
		return latest
	}
	return e.Next(stats)
}

// LanguageFor picks the language of an explicit solution file, or the
// configured one.
func (e *Env) LanguageFor(file string) (judge.Language, error) {
	if file == "" {
		return e.Lang, nil
	}
	return judge.LanguageForFile(file)
}

type Submission struct {
	Result judge.Result
	ID     int64 // 0 for sample runs and runs that were not recorded
}

// Submit judges a solution and records it. Sample runs only use the
// statement's examples and are recorded as such.
func (e *Env) Submit(ctx context.Context, ex *track.Exercise, file string, samplesOnly bool, emit func(judge.Event)) (Submission, error) {
	lang, err := e.LanguageFor(file)
	if err != nil {
		return Submission{}, err
	}
	if file == "" {
		file = e.WS.SolutionPath(ex, lang)
	}
	code, err := os.ReadFile(file)
	if err != nil {
		return Submission{}, err
	}
	tests, err := ex.Tests()
	if err != nil {
		return Submission{}, err
	}
	if samplesOnly {
		tests = tests[:min(len(tests), max(1, ex.Samples))]
	}
	res := judge.Run(ctx, e.Runner, judge.Request{
		Lang: lang, Source: file, Tests: tests, TimeLimit: ex.Limit(), StopOnFail: true,
	}, emit)
	sub := Submission{Result: res}
	if res.Canceled || res.Verdict == judge.IE {
		return sub, nil
	}
	id, err := e.Store.Add(store.Submission{
		Exercise: ex.Key(), Language: lang.ID, Verdict: string(res.Verdict), Passed: res.Passed(),
		Total: res.Total, TimeMs: res.MaxTime().Milliseconds(), SamplesOnly: samplesOnly,
	}, string(code))
	if !samplesOnly {
		sub.ID = id
	}
	return sub, err
}

// FirstName is used to greet the user.
func FirstName() string {
	if u, err := user.Current(); err == nil {
		name := u.Name
		if name == "" {
			name = u.Username
		}
		if f := strings.Fields(name); len(f) > 0 {
			return f[0]
		}
	}
	return "dev"
}

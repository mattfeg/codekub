package judge

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type TestCase struct {
	InputPath, OutputPath string
	Sample                bool
}

type Request struct {
	Lang       Language
	Source     string
	Tests      []TestCase
	TimeLimit  time.Duration // base limit; the language factor is applied here
	StopOnFail bool
}

type TestResult struct {
	Index   int
	Sample  bool
	Verdict Verdict
	Time    time.Duration
	Diff    *Diff
	Input   string // only filled for samples, which the statement already shows
	Detail  string // runtime error output, or why the judge gave up
}

type Result struct {
	Verdict  Verdict
	Tests    []TestResult
	Total    int
	Limit    time.Duration
	Detail   string // compilation error or internal error message
	Canceled bool
}

func (r Result) Passed() int {
	n := 0
	for _, t := range r.Tests {
		if t.Verdict == AC {
			n++
		}
	}
	return n
}

func (r Result) MaxTime() time.Duration {
	var m time.Duration
	for _, t := range r.Tests {
		m = max(m, t.Time)
	}
	return m
}

type EventKind int

const (
	EventChecking EventKind = iota
	EventTestStart
	EventTestDone
	EventDone
)

type Event struct {
	Kind   EventKind
	Index  int
	Test   *TestResult
	Result *Result
}

// Run judges a solution test by test, reporting progress through emit (which
// may be nil). It stops at the first failure when StopOnFail is set.
func Run(ctx context.Context, runner Runner, req Request, emit func(Event)) Result {
	if emit == nil {
		emit = func(Event) {}
	}
	limit := time.Duration(float64(req.TimeLimit) * req.Lang.TimeFactor)
	res := Result{Verdict: AC, Total: len(req.Tests), Limit: limit}
	finish := func() Result {
		emit(Event{Kind: EventDone, Result: &res})
		return res
	}

	ws, err := runner.Prepare(req.Lang, req.Source)
	if err != nil {
		res.Verdict, res.Detail = IE, err.Error()
		return finish()
	}
	defer ws.Close()

	if len(req.Lang.Check) > 0 {
		emit(Event{Kind: EventChecking})
		ok, out, err := ws.Check(ctx)
		switch {
		case ctx.Err() != nil:
			res.Canceled = true
			return finish()
		case err != nil:
			res.Verdict, res.Detail = IE, err.Error()
			return finish()
		case !ok:
			res.Verdict, res.Detail = CE, out
			return finish()
		}
	}

	for i, tc := range req.Tests {
		emit(Event{Kind: EventTestStart, Index: i})
		tr, err := runTest(ctx, ws, i, tc, limit)
		if err != nil {
			res.Verdict, res.Detail = IE, err.Error()
			return finish()
		}
		if tr == nil {
			res.Canceled = true
			return finish()
		}
		res.Tests = append(res.Tests, *tr)
		emit(Event{Kind: EventTestDone, Index: i, Test: tr})
		if tr.Verdict != AC {
			if res.Verdict == AC {
				res.Verdict = tr.Verdict
			}
			if req.StopOnFail {
				break
			}
		}
	}
	return finish()
}

// runTest returns nil when the run was canceled.
func runTest(ctx context.Context, ws Workspace, i int, tc TestCase, limit time.Duration) (*TestResult, error) {
	input, err := os.ReadFile(tc.InputPath)
	if err != nil {
		return nil, err
	}
	want, err := os.ReadFile(tc.OutputPath)
	if err != nil {
		return nil, err
	}
	ex := ws.Exec(ctx, strings.NewReader(string(input)), limit)
	tr := &TestResult{Index: i, Sample: tc.Sample, Time: ex.Elapsed}
	if tc.Sample {
		tr.Input = string(input)
	}
	switch {
	case ex.Canceled:
		return nil, nil
	case ex.Err != nil:
		return nil, ex.Err
	case ex.TimedOut:
		tr.Verdict = TLE
		tr.Detail = fmt.Sprintf("passou de %s", fmtSeconds(limit))
	case ex.OutputExceeded:
		tr.Verdict = RE
		tr.Detail = fmt.Sprintf("a saída passou de %d MB", OutputLimit>>20)
	case ex.Signal != "" || ex.ExitCode != 0:
		tr.Verdict = RE
		tr.Detail = runtimeDetail(ex)
	default:
		tr.Verdict, tr.Diff = Compare(ex.Stdout, string(want))
	}
	return tr, nil
}

func runtimeDetail(ex Execution) string {
	reason := fmt.Sprintf("terminou com código %d", ex.ExitCode)
	if ex.Signal != "" {
		reason = "encerrado pelo sinal " + ex.Signal
	}
	lines := strings.Split(strings.TrimSpace(ex.Stderr), "\n")
	if len(lines) > 8 {
		lines = lines[len(lines)-8:]
	}
	if len(lines) == 1 && lines[0] == "" {
		return reason
	}
	return reason + "\n" + strings.Join(lines, "\n")
}

// fmtSeconds writes durations the Brazilian way: 3 s, 1,5 s.
func fmtSeconds(d time.Duration) string {
	return strings.Replace(strconv.FormatFloat(d.Seconds(), 'f', -1, 64), ".", ",", 1) + " s"
}

// FormatTime shows a test's running time with millisecond precision: 0,042 s.
func FormatTime(d time.Duration) string {
	return strings.Replace(fmt.Sprintf("%.3f s", d.Seconds()), ".", ",", 1)
}

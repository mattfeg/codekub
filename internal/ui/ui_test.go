package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mattfeg/codekub/internal/core"
)

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// harness drives the model like the Bubble Tea runtime would: it executes
// commands and feeds their messages back, except animation ticks, which the
// test sends itself with frames().
type harness struct {
	t *testing.T
	m *model
}

func (h *harness) send(msg tea.Msg) {
	_, cmd := h.m.Update(msg)
	h.run(cmd)
}

func (h *harness) run(cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	ch := make(chan tea.Msg, 1)
	go func() { ch <- cmd() }()
	var msg tea.Msg
	select {
	case msg = <-ch:
	case <-time.After(2 * time.Second):
		return // a watcher waiting for a save that never comes
	}
	switch msg := msg.(type) {
	case nil, frameMsg:
	case tea.BatchMsg:
		for _, c := range msg {
			h.run(c)
		}
	default:
		h.send(msg)
	}
}

func (h *harness) key(k string) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
	switch k {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		msg = tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		msg = tea.KeyMsg{Type: tea.KeyTab}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	}
	h.send(msg)
}

func (h *harness) frames(n int) {
	for i := 0; i < n; i++ {
		h.m.Update(frameMsg{})
	}
}

func (h *harness) view() string { return ansi.ReplaceAllString(h.m.View(), "") }

func (h *harness) expect(substr string) {
	h.t.Helper()
	if v := h.view(); !strings.Contains(v, substr) {
		h.t.Fatalf("screen should contain %q:\n%s", substr, v)
	}
}

func newEnv(t *testing.T) (*core.Env, string) {
	t.Helper()
	solutions := t.TempDir()
	t.Setenv("CODEKUB_HOME", t.TempDir())
	t.Setenv("CODEKUB_SOLUTIONS", solutions)
	t.Setenv("CODEKUB_TRACK", filepath.Join("..", "..", "track"))
	env, err := core.Open()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(env.Close)
	return env, solutions
}

func TestSolveThroughTheTUI(t *testing.T) {
	snapDir := os.Getenv("SNAP_DIR")
	for _, size := range [][2]int{{80, 24}, {120, 40}, {160, 50}} {
		env, solutions := newEnv(t) // fresh progress for each terminal size
		h := &harness{t: t, m: newModel(env, Options{})}
		h.send(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		snap := func(name string) {
			if snapDir != "" {
				os.WriteFile(filepath.Join(snapDir, fmt.Sprintf("%dx%d-%s.txt", size[0], size[1], name)), []byte(h.view()), 0o644)
			}
		}
		h.frames(60)
		snap("01-splash")
		h.key("x")
		h.frames(40)
		snap("02-home")
		h.key("tab")
		h.frames(20)
		snap("03-track")
		h.key("enter") // focus the exercise list, on the next exercise
		h.key("enter") // open it
		h.frames(20)
		snap("04-problem")
		if h.m.screen != scrProblem || h.m.prob.ex.ID != "001" {
			t.Fatalf("expected problem 001, got screen %d", h.m.screen)
		}

		path := filepath.Join(solutions, "001-hello-world.rb")
		os.WriteFile(path, []byte("puts 'Ola'\n"), 0o644)
		h.key("s")
		h.frames(30)
		snap("05-result-WA")
		if r := h.m.jdg.result; r == nil || r.Verdict != "WA" {
			t.Fatalf("expected WA, got %+v", r)
		}

		os.WriteFile(path, []byte("puts 'Olá, mundo!'\n"), 0o644)
		h.key("r")
		h.frames(8)
		snap("06-result-AC")
		if r := h.m.jdg.result; r == nil || r.Verdict != "AC" || h.m.jdg.subID == 0 {
			t.Fatalf("expected a recorded AC, got %+v", r)
		}
		if !h.m.stats["001-hello-world"].Solved {
			t.Fatal("progress was not updated after AC")
		}

		// Watch: a save triggers a new judgement on its own.
		runs := h.m.jdg.runs
		go func() { // the user saves while the watcher is waiting
			time.Sleep(500 * time.Millisecond)
			os.WriteFile(path, []byte("puts 'Olá, mundo!'\n# salvo de novo\n"), 0o644)
		}()
		h.key("w")
		h.frames(8)
		snap("07-watch")
		if h.m.jdg.runs != runs+1 || h.m.jdg.result == nil || h.m.jdg.result.Verdict != "AC" {
			t.Fatalf("watch should have rejudged once, runs %d→%d", runs, h.m.jdg.runs)
		}

		h.key("esc")
		h.key("esc")
		h.key("tab")
		h.frames(90)
		snap("08-progress")
		if h.m.screen != scrProgress {
			t.Fatalf("expected the progress tab, got screen %d", h.m.screen)
		}
		h.expect("resolvidos")
		h.key("?")
		h.frames(20)
		snap("09-help")
		h.key("esc")
		if h.m.helpOpen {
			t.Fatal("esc should close help")
		}
	}
}

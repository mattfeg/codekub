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
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case "backspace":
		msg = tea.KeyMsg{Type: tea.KeyBackspace}
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
		if h.m.screen != scrProblem || h.m.prob.ex.ID != "1000" {
			t.Fatalf("expected problem 1000, got screen %d", h.m.screen)
		}

		path := filepath.Join(solutions, "1000-hello-world.rb")
		os.WriteFile(path, []byte("puts 'Ola'\n"), 0o644)
		h.key("s")
		h.frames(30)
		snap("05-result-WA")
		if r := h.m.jdg.result; r == nil || r.Verdict != "WA" {
			t.Fatalf("expected WA, got %+v", r)
		}

		os.WriteFile(path, []byte("puts 'Hello World!'\n"), 0o644)
		h.key("r")
		h.frames(8)
		snap("06-result-AC")
		if r := h.m.jdg.result; r == nil || r.Verdict != "AC" || h.m.jdg.subID == 0 {
			t.Fatalf("expected a recorded AC, got %+v", r)
		}
		if !h.m.stats["1000-hello-world"].Solved {
			t.Fatal("progress was not updated after AC")
		}

		// Watch: a save triggers a new judgement on its own.
		runs := h.m.jdg.runs
		go func() { // the user saves while the watcher is waiting
			time.Sleep(500 * time.Millisecond)
			os.WriteFile(path, []byte("puts 'Hello World!'\n# salvo de novo\n"), 0o644)
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

func TestSettingsThroughTheTUI(t *testing.T) {
	env, _ := newEnv(t)
	for _, key := range []string{"CODEKUB_EDITOR", "CODEKUB_LANG"} {
		t.Setenv(key, "")
	}
	h := &harness{t: t, m: newModel(env, Options{})}
	h.send(tea.WindowSizeMsg{Width: 140, Height: 40})
	h.key("x")
	for h.m.screen != scrSettings {
		h.key("tab")
	}
	h.frames(20)
	h.expect("AJUSTES")

	// Pick the first installed editor from the list.
	h.key("enter")
	if len(h.m.set.picker) < 2 {
		t.Fatalf("editor picker should list the default plus installed editors: %+v", h.m.set.picker)
	}
	h.key("down")
	want := h.m.set.picker[1].value
	h.key("enter")
	if got := env.WS.SavedConfig().Editor; got != want {
		t.Fatalf("editor saved as %q, want %q", got, want)
	}
	if snap := os.Getenv("SNAP_DIR"); snap != "" {
		h.key("enter")
		os.WriteFile(filepath.Join(snap, "140x40-10-settings.txt"), []byte(h.view()), 0o644)
		h.key("esc")
	}

	openCustom := func() {
		h.key("enter")
		for !h.m.set.picker[h.m.set.pick].custom {
			h.key("down")
		}
		h.key("enter")
	}

	// Custom command: global keys like q must become text while typing. The
	// command typed is an editor installed here (Windows has no vi).
	installed := strings.Fields(h.m.editorOptions()[1].value)[0]
	openCustom()
	for _, r := range installed + "q" {
		h.key(string(r))
	}
	if h.m.screen != scrSettings || h.m.set.input != installed+"q" {
		t.Fatalf("typing should capture keys, input %q", h.m.set.input)
	}
	h.key("backspace")
	h.key("enter")
	if got := env.WS.SavedConfig().Editor; got != installed {
		t.Fatalf("custom editor saved as %q, want %q", got, installed)
	}

	// A command that does not exist is refused.
	openCustom()
	for _, r := range "nao-existe" {
		h.key(string(r))
	}
	h.key("enter")
	if got := env.WS.SavedConfig().Editor; got != installed {
		t.Fatalf("a missing command must not be saved, got %q", got)
	}
	h.key("esc")
	h.key("esc")

	// Language: choose Python when it is installed.
	h.key("down")
	h.key("enter")
	for i, o := range h.m.set.picker {
		if o.value == "python" && !o.disabled {
			for h.m.set.pick != i {
				h.key("down")
			}
			h.key("enter")
			if env.Lang.ID != "python" || env.WS.SavedConfig().Language != "python" {
				t.Fatalf("language not switched: %s / %s", env.Lang.ID, env.WS.SavedConfig().Language)
			}
		}
	}
}

func TestSearchThroughTheTUI(t *testing.T) {
	env, _ := newEnv(t)
	h := &harness{t: t, m: newModel(env, Options{})}
	h.send(tea.WindowSizeMsg{Width: 120, Height: 40})
	h.key("x")

	// By number: "q" must be text while searching, not quit.
	h.key("/")
	for _, r := range "2q" {
		h.key(string(r))
	}
	if !h.m.srch.open || h.m.srch.query != "2q" {
		t.Fatalf("search should capture keys, query %q", h.m.srch.query)
	}
	h.key("backspace")
	for _, r := range "165" {
		h.key(string(r))
	}
	h.frames(20)
	h.expect("BUSCAR PROBLEMA")
	if snap := os.Getenv("SNAP_DIR"); snap != "" {
		os.WriteFile(filepath.Join(snap, "120x40-11-search.txt"), []byte(h.view()), 0o644)
	}
	h.key("enter")
	if h.m.srch.open || h.m.screen != scrProblem || h.m.prob.ex.ID != "2165" {
		t.Fatalf("enter should open problem 2165, got screen %d ex %v", h.m.screen, h.m.prob.ex)
	}
	if p := env.Track.Phases[h.m.trk.phase]; p.Exercises[h.m.trk.ex].ID != "2165" {
		t.Fatal("the track screen should point at the problem found")
	}

	// By title, ignoring accents, from the problem screen: it replaces the
	// open problem instead of stacking another one.
	depth := len(h.m.stack)
	h.key("/")
	for _, r := range "area do circulo" {
		h.key(string(r))
	}
	h.key("enter")
	if h.m.prob.ex.ID != "1002" || len(h.m.stack) != depth {
		t.Fatalf("expected 1002 at the same depth, got %s (stack %d→%d)", h.m.prob.ex.ID, depth, len(h.m.stack))
	}

	// No results: enter does nothing, esc closes.
	h.key("/")
	h.key("z")
	h.key("z")
	h.key("z")
	h.key("z")
	h.key("enter")
	if !h.m.srch.open {
		t.Fatal("enter without results should keep the search open")
	}
	h.key("esc")
	if h.m.srch.open || h.m.prob.ex.ID != "1002" {
		t.Fatal("esc should close the search and keep the problem")
	}
}

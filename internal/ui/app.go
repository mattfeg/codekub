package ui

import (
	"context"
	"math"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"

	"github.com/mattfeg/codekub/internal/core"
	"github.com/mattfeg/codekub/internal/judge"
	"github.com/mattfeg/codekub/internal/store"
	"github.com/mattfeg/codekub/internal/track"
)

type screen int

const (
	scrSplash screen = iota
	scrHome
	scrTrack
	scrProgress
	scrProblem
	scrJudge
	scrSettings
)

var tabs = []struct {
	s     screen
	label string
}{{scrHome, "Início"}, {scrTrack, "Trilha"}, {scrProgress, "Progresso"}, {scrSettings, "Ajustes"}}

// Options open the TUI somewhere other than the splash screen.
type Options struct {
	Exercise *track.Exercise
	File     string // explicit solution file; empty means the workspace file
	Watch    bool
}

type (
	frameMsg      struct{}
	judgeEventMsg struct {
		gen int
		ev  judge.Event
	}
	judgeDoneMsg struct {
		gen int
		sub core.Submission
		err error
	}
	fileChangedMsg  struct{ gen int }
	debounceMsg     struct{ seq int }
	editorClosedMsg struct{ err error }
)

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(time.Time) tea.Msg { return frameMsg{} })
}

type model struct {
	env        *core.Env
	w, h       int
	frame      int
	screen     screen
	stack      []screen
	wipe       float64
	helpOpen   bool
	toast      string
	toastUntil int
	stats      map[string]store.ExerciseStat
	subs       []store.Submission
	activity   map[string]int
	name       string

	home homeState
	trk  trackState
	prob problemState
	jdg  judgeState
	prog progressState
	set  settingsState
}

func Run(env *core.Env, opts Options) error {
	m := newModel(env, opts)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	m.jdg.stopRun()
	m.jdg.stopWatch()
	return err
}

func newModel(env *core.Env, opts Options) *model {
	m := &model{env: env, w: 110, h: 36, wipe: 1, name: core.FirstName()}
	m.refreshStats()
	m.trk.focusNext(m)
	if opts.Exercise != nil {
		m.openExercise(opts.Exercise, opts.File)
		m.stack = []screen{scrHome, scrProblem}
		m.screen = scrJudge
		m.jdg.watch = opts.Watch
	}
	return m
}

// refreshStats reloads progress from the database; views read the cache.
func (m *model) refreshStats() {
	m.stats = m.env.Stats()
	m.subs, _ = m.env.Store.Submissions(1 << 30)
	m.activity, _ = m.env.Store.Activity()
}

func (m *model) Init() tea.Cmd {
	cmds := []tea.Cmd{tick()}
	if m.screen == scrJudge && m.jdg.watch {
		cmds = append(cmds, m.startWatch())
	} else if m.screen == scrJudge {
		cmds = append(cmds, m.startJudge(false))
	}
	return tea.Batch(cmds...)
}

func (m *model) openExercise(ex *track.Exercise, file string) {
	if m.jdg.ex != ex || m.jdg.file != file {
		m.jdg.stopRun()
		m.jdg.stopWatch()
		m.jdg = judgeState{ex: ex, file: file}
	}
	m.prob.ex = ex
	if file == "" {
		m.env.WS.EnsureSolution(ex, m.env.Lang)
	}
}

func (m *model) solutionPath() string {
	if m.jdg.file != "" {
		return m.jdg.file
	}
	return m.env.WS.SolutionPath(m.jdg.ex, m.env.Lang)
}

func (m *model) navigate(s screen) {
	m.stack = append(m.stack, m.screen)
	m.switchTo(s)
}

func (m *model) switchTo(s screen) {
	m.screen = s
	m.wipe = 0
	switch s {
	case scrHome:
		m.home.enter(m)
	case scrProgress:
		m.prog.enter(m)
	case scrProblem:
		m.prob.offset = 0
	case scrSettings:
		m.set.enter()
	}
}

func (m *model) back() {
	if len(m.stack) == 0 {
		return
	}
	s := m.stack[len(m.stack)-1]
	m.stack = m.stack[:len(m.stack)-1]
	m.switchTo(s)
}

func (m *model) notify(msg string) {
	m.toast = msg
	m.toastUntil = m.frame + 90
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case frameMsg:
		m.frame++
		if m.wipe < 1 {
			m.wipe = math.Min(1, m.wipe+0.085)
		}
		m.home.step()
		m.prog.step()
		m.jdg.step()
		return m, tick()
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, m.handleKey(msg.String())
	case judgeEventMsg:
		if msg.gen != m.jdg.gen {
			return m, nil
		}
		m.jdg.onEvent(msg.ev)
		return m, waitMsg(m.jdg.ch)
	case judgeDoneMsg:
		if msg.gen != m.jdg.gen {
			return m, nil
		}
		m.jdg.onDone(msg)
		m.refreshStats()
		return m, nil
	case fileChangedMsg:
		if msg.gen != m.jdg.watchGen {
			return m, nil
		}
		m.jdg.saveSeq++
		seq := m.jdg.saveSeq
		return m, tea.Batch(waitFile(m.jdg.watchCh, msg.gen),
			tea.Tick(200*time.Millisecond, func(time.Time) tea.Msg { return debounceMsg{seq} }))
	case debounceMsg:
		if msg.seq == m.jdg.saveSeq && m.jdg.watch {
			m.notify(filepath.Base(m.solutionPath()) + " salvo · julgando")
			return m, m.startJudge(false)
		}
	case editorClosedMsg:
		if msg.err != nil {
			m.notify("não consegui abrir o editor")
		}
	}
	return m, nil
}

func (m *model) handleKey(key string) tea.Cmd {
	if m.screen == scrSplash {
		m.switchTo(scrHome)
		return nil
	}
	if m.screen == scrSettings && m.set.typing {
		return m.settingsKey(key) // every key is text while typing a command
	}
	if key == "?" || (m.helpOpen && key == "esc") {
		m.helpOpen = !m.helpOpen
		m.wipe = 0
		return nil
	}
	if m.helpOpen {
		if key == "q" {
			return tea.Quit
		}
		return nil
	}
	switch key {
	case "q":
		return tea.Quit
	case "tab":
		m.cycleTab(1)
		return nil
	case "shift+tab":
		m.cycleTab(-1)
		return nil
	}
	switch m.screen {
	case scrHome:
		return m.homeKey(key)
	case scrTrack:
		return m.trackKey(key)
	case scrProblem:
		return m.problemKey(key)
	case scrJudge:
		return m.judgeKey(key)
	case scrSettings:
		return m.settingsKey(key)
	case scrProgress:
		if key == "r" {
			m.prog.enter(m)
			m.wipe = 0
		}
	}
	return nil
}

func (m *model) currentTab() screen {
	if m.screen == scrProblem || m.screen == scrJudge {
		return scrTrack
	}
	return m.screen
}

func (m *model) cycleTab(dir int) {
	cur := 0
	for i, t := range tabs {
		if t.s == m.currentTab() {
			cur = i
		}
	}
	m.stack = nil
	m.switchTo(tabs[(cur+dir+len(tabs))%len(tabs)].s)
}

// --- judging and watching -------------------------------------------------

func waitMsg(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		return msg
	}
}

func waitFile(ch <-chan struct{}, gen int) tea.Cmd {
	return func() tea.Msg {
		if _, ok := <-ch; !ok {
			return nil
		}
		return fileChangedMsg{gen}
	}
}

func (m *model) startJudge(samples bool) tea.Cmd {
	j := &m.jdg
	j.stopRun()
	j.gen++
	gen := j.gen
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan tea.Msg, 256)
	total := 0
	if tests, err := j.ex.Tests(); err == nil {
		total = len(tests)
		if samples {
			total = min(total, max(1, j.ex.Samples))
		}
	}
	solvedBefore := m.stats[j.ex.Key()].Solved
	*j = judgeState{
		ex: j.ex, file: j.file, watch: j.watch, watcher: j.watcher, watchCh: j.watchCh, watchGen: j.watchGen,
		saveSeq: j.saveSeq, gen: gen, cancel: cancel, ch: ch, running: true, samples: samples,
		total: total, solvedBefore: solvedBefore, runs: j.runs + 1,
	}
	env, ex, file := m.env, j.ex, j.file
	go func() {
		sub, err := env.Submit(ctx, ex, file, samples, func(e judge.Event) {
			if e.Kind != judge.EventDone {
				ch <- judgeEventMsg{gen, e}
			}
		})
		ch <- judgeDoneMsg{gen, sub, err}
		close(ch)
	}()
	return waitMsg(ch)
}

func (m *model) startWatch() tea.Cmd {
	j := &m.jdg
	j.stopWatch()
	j.watch = true
	path := m.solutionPath()
	w, err := fsnotify.NewWatcher()
	if err == nil {
		err = w.Add(filepath.Dir(path))
	}
	if err != nil {
		m.notify("não consegui observar o arquivo")
		j.watch = false
		return nil
	}
	j.watchGen++
	ch := make(chan struct{}, 1)
	j.watcher, j.watchCh = w, ch
	base := filepath.Base(path)
	go func() {
		defer close(ch)
		for {
			select {
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				if filepath.Base(ev.Name) == base && ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
			}
		}
	}()
	return waitFile(ch, j.watchGen)
}

func (m *model) openEditor() tea.Cmd {
	path := m.solutionPath()
	return tea.ExecProcess(m.env.WS.EditorCommand(path), func(err error) tea.Msg { return editorClosedMsg{err} })
}

// --- frame ------------------------------------------------------------------

func (m *model) View() string {
	if m.screen == scrSplash {
		return m.splashView()
	}
	header := m.headerView()
	footer := m.footerView()
	bodyH := m.h - lipgloss.Height(header) - lipgloss.Height(footer)
	var body string
	switch {
	case m.helpOpen:
		body = m.helpView(bodyH)
	case m.screen == scrHome:
		body = m.homeView()
	case m.screen == scrTrack:
		body = m.trackView(bodyH)
	case m.screen == scrProblem:
		body = m.problemView(bodyH)
	case m.screen == scrJudge:
		body = m.judgeView()
	case m.screen == scrProgress:
		body = m.progressView()
	case m.screen == scrSettings:
		body = m.settingsView()
	}
	return header + "\n" + m.fitBody(body, bodyH) + "\n" + footer
}

// fitBody pads/crops the body and applies the top-down wipe transition.
func (m *model) fitBody(body string, h int) string {
	lines := strings.Split(body, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	for len(lines) < h {
		lines = append(lines, "")
	}
	eased := 1 - math.Pow(1-m.wipe, 3)
	visible := int(eased*float64(h) + 0.5)
	clip := lipgloss.NewStyle().MaxWidth(m.w)
	for i := range lines {
		if i >= visible {
			lines[i] = ""
			continue
		}
		lines[i] = clip.Render("  " + lines[i])
	}
	return strings.Join(lines, "\n")
}

func (m *model) headerView() string {
	left := " " + mascot("idle", m.frame) + "  " + gradient("codekub", hexPrimary, hexMint, float64(m.frame)/10)
	var ts []string
	for _, t := range tabs {
		if m.currentTab() == t.s {
			ts = append(ts, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(hexInk)).
				Background(lipgloss.Color(hexPrimary)).Padding(0, 1).Render(t.label))
		} else {
			ts = append(ts, sMuted.Padding(0, 1).Render(t.label))
		}
	}
	center := strings.Join(ts, " ")
	right := sMuted.Render(m.env.Lang.Name) + sFaint.Render("  ·  ") + sMuted.Render(m.name+" ")
	gap := m.w - lipgloss.Width(left) - lipgloss.Width(center) - lipgloss.Width(right) - 6
	return left + "      " + center + strings.Repeat(" ", max(1, gap)) + right + "\n" + sFaint.Render(strings.Repeat("─", m.w))
}

func (m *model) footerView() string {
	var help string
	switch m.screen {
	case scrHome:
		help = helpLine("enter", "continuar", "w", "watch", "tab", "abas", "q", "sair")
	case scrTrack:
		help = helpLine("↑↓", "navegar", "→ enter", "abrir", "←", "voltar", "tab", "abas", "q", "sair")
	case scrProblem:
		help = helpLine("↑↓", "rolar", "e", "editar", "s", "enviar", "x", "exemplos", "w", "watch", "esc", "voltar")
	case scrJudge:
		help = helpLine("e", "editar", "r", "julgar de novo", "x", "exemplos", "w", "watch", "esc", "voltar")
	case scrProgress:
		help = helpLine("r", "repetir animação", "tab", "abas", "q", "sair")
	case scrSettings:
		help = helpLine("↑↓", "navegar", "enter", "mudar", "esc", "fechar lista", "tab", "abas", "q", "sair")
		if m.set.typing {
			help = helpLine("enter", "salvar", "esc", "cancelar")
		}
	}
	if m.helpOpen {
		help = helpLine("esc ?", "fechar ajuda", "q", "sair")
	} else {
		help = helpLine("?", "ajuda") + sFaint.Render("  ·  ") + help
	}
	right := ""
	if m.toast != "" && m.frame < m.toastUntil {
		right = lipgloss.NewStyle().Foreground(lipgloss.Color(hexInk)).Background(lipgloss.Color(hexMint)).
			Bold(true).Padding(0, 1).Render(m.toast)
	}
	gap := m.w - lipgloss.Width(help) - lipgloss.Width(right) - 3
	return sFaint.Render(strings.Repeat("─", m.w)) + "\n " + help + strings.Repeat(" ", max(1, gap)) + right
}

var glyphs = map[rune][5]string{
	'c': {" ████", "██   ", "██   ", "██   ", " ████"},
	'o': {" ███ ", "██ ██", "██ ██", "██ ██", " ███ "},
	'd': {"████ ", "██ ██", "██ ██", "██ ██", "████ "},
	'e': {"█████", "██   ", "████ ", "██   ", "█████"},
	'k': {"██ ██", "██ █ ", "███  ", "██ █ ", "██ ██"},
	'u': {"██ ██", "██ ██", "██ ██", "██ ██", " ███ "},
	'b': {"████ ", "██ ██", "████ ", "██ ██", "████ "},
}

func (m *model) splashView() string {
	lines := make([]string, 5)
	for _, ch := range "codekub" {
		for i := range lines {
			lines[i] += glyphs[ch][i] + " "
		}
	}
	reveal := m.frame * 2
	phase := float64(m.frame) / 7
	for i, l := range lines {
		r := []rune(l)
		if reveal < len(r) {
			r = append(r[:reveal:reveal], []rune(strings.Repeat(" ", len(r)-reveal))...)
		}
		lines[i] = gradient(string(r), hexPrimary, hexMint, phase)
	}
	tag := []rune("by Matheus Ferreira")
	typed := max(0, min(len(tag), m.frame-28))
	tagline := sMuted.Render(string(tag[:typed]))
	if typed < len(tag) && m.frame%10 < 5 {
		tagline += fg(hexMint).Render("▌")
	}
	mood := "sleep"
	if m.frame > 24 {
		mood = "idle"
	}
	hint := " "
	if m.frame > 60 && (m.frame/15)%2 == 0 {
		hint = sFaint.Render("pressione qualquer tecla")
	}
	content := lipgloss.JoinVertical(lipgloss.Center, strings.Join(lines, "\n"), "", mascot(mood, m.frame)+"  "+tagline, "", hint)
	return lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, content)
}

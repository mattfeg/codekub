package ui

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"

	"github.com/mattfeg/codekub/internal/judge"
	"github.com/mattfeg/codekub/internal/track"
)

type particle struct {
	x, y, vx, vy float64
	life         int
	ch, hex      string
}

type judgeState struct {
	ex           *track.Exercise
	file         string
	samples      bool
	running      bool
	checking     bool
	tests        []judge.TestResult
	current      int
	total        int
	result       *judge.Result
	subID        int64
	solvedBefore bool
	runs         int

	gen    int
	cancel context.CancelFunc
	ch     chan tea.Msg

	watch    bool
	watcher  *fsnotify.Watcher
	watchCh  chan struct{}
	watchGen int
	saveSeq  int

	bar       anim
	particles []particle
}

func (j *judgeState) stopRun() {
	if j.cancel != nil {
		j.cancel()
	}
}

func (j *judgeState) stopWatch() {
	if j.watcher != nil {
		j.watcher.Close()
		j.watcher = nil
	}
}

func (j *judgeState) onEvent(ev judge.Event) {
	switch ev.Kind {
	case judge.EventChecking:
		j.checking = true
	case judge.EventTestStart:
		j.checking, j.current = false, ev.Index
	case judge.EventTestDone:
		j.tests = append(j.tests, *ev.Test)
		j.bar.target = float64(len(j.tests)) / float64(max(1, j.total))
	}
}

func (j *judgeState) onDone(msg judgeDoneMsg) {
	j.running, j.checking = false, false
	if msg.err != nil {
		j.result = &judge.Result{Verdict: judge.IE, Detail: msg.err.Error()}
		return
	}
	res := msg.sub.Result
	if res.Canceled {
		return
	}
	j.result, j.subID = &res, msg.sub.ID
	if res.Verdict == judge.AC {
		j.bar.target = 1
		if !j.samples {
			j.burst()
		}
	}
}

func (j *judgeState) step() {
	j.bar.step()
	alive := j.particles[:0]
	for _, p := range j.particles {
		p.x += p.vx
		p.y += p.vy
		p.vy += 0.025
		p.vx *= 0.96
		p.life--
		if p.life > 0 {
			alive = append(alive, p)
		}
	}
	j.particles = alive
}

// burst throws sparks out of the bottom edge of the verdict card.
func (j *judgeState) burst() {
	chars := []string{"✦", "✧", "•", "*", "+", "·"}
	colors := []string{hexMint, hexPrimary, hexWarm, verdictHex["AC"], "#F472B6"}
	for i := 0; i < 60; i++ {
		j.particles = append(j.particles, particle{
			x: 4 + rand.Float64()*48, y: 0,
			vx: (rand.Float64()*2 - 1) * 0.9, vy: -0.35 + rand.Float64()*0.5,
			life: 30 + rand.Intn(30),
			ch:   chars[rand.Intn(len(chars))], hex: colors[rand.Intn(len(colors))],
		})
	}
}

func (m *model) judgeKey(key string) tea.Cmd {
	j := &m.jdg
	switch key {
	case "s", "r":
		return m.startJudge(false)
	case "x":
		return m.startJudge(true)
	case "e":
		return m.openEditor()
	case "o":
		return m.openSource()
	case "w":
		if j.watch {
			j.stopWatch()
			j.watch = false
			m.notify("watch desligado")
			return nil
		}
		m.notify("watch ligado · salve o arquivo para julgar")
		return m.startWatch()
	case "esc":
		j.stopWatch()
		j.watch = false
		m.back()
	}
	return nil
}

func (m *model) judgeView() string {
	j := &m.jdg
	title := sBold.Render(j.ex.ID+" · "+j.ex.Title) + sMuted.Render("  ·  "+m.env.Lang.Name+"  ·  ")
	switch {
	case j.samples:
		title += fg(hexWarm).Render("só exemplos · não conta como submissão")
	case j.subID > 0:
		title += sMuted.Render(fmt.Sprintf("submissão #%d", j.subID))
	default:
		title += sMuted.Render(shortPath(m.solutionPath()))
	}

	watchLine := ""
	if j.watch {
		state := sMuted.Render("observando " + shortPath(m.solutionPath()))
		if j.running {
			state = fg(hexWarm).Render("alteração detectada · julgando")
		}
		watchLine = fg(pulse(m.frame, hexMint, hexFaint)).Render("●") + " " + fg(hexMint).Bold(true).Render("WATCH") + "  " + state
	}

	if j.runs == 0 {
		lines := []string{"", mascot("sleep", m.frame) + fg(hexFaint).Render([]string{" z", " zᶻ", " zᶻz"}[(m.frame/12)%3]), ""}
		if j.watch {
			lines = append(lines, sBold.Render("Esperando você salvar o arquivo…"),
				sMuted.Render("Edite no seu editor (tecla e). A cada salvamento,"),
				sMuted.Render("o codekub roda todos os testes de novo."))
		} else {
			lines = append(lines, sBold.Render("Pronto para julgar."),
				sMuted.Render("s envia a solução  ·  x testa só os exemplos  ·  w liga o watch"))
		}
		lines = append(lines, "", sFaint.Render(shortPath(m.solutionPath())), "")
		return strings.Join([]string{"", title, watchLine, "", panel(hexMint).Width(66).Render(strings.Join(lines, "\n"))}, "\n")
	}

	var status string
	switch {
	case j.running && j.checking:
		status = spinner(m.frame) + " " + sText.Render("verificando sintaxe…")
	case j.running:
		status = spinner(m.frame) + " " + sText.Render(fmt.Sprintf("julgando · teste %d de %d", j.current+1, j.total))
	case j.result != nil:
		status = fg(verdictHex[string(j.result.Verdict)]).Bold(true).Render("■") + " " + sText.Render("julgamento concluído")
	}
	barLine := gradientBar(44, j.bar.pos, hexPrimary, hexMint) + "  " + sMuted.Render(fmt.Sprintf("%3.0f%%", math.Min(100, math.Max(0, j.bar.pos*100))))
	left := strings.Join(append([]string{status, barLine, ""}, m.testRows()...), "\n")

	body := left
	if j.result != nil && !j.running {
		card := m.resultCard()
		if m.w >= 120 {
			body = lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.NewStyle().Width(56).Render(left), "  ", card)
		} else {
			body = left + "\n\n" + card
		}
	}
	return strings.Join([]string{"", title, watchLine, "", body}, "\n")
}

// testRows lists each test, or draws numbered tiles when there are many.
func (m *model) testRows() []string {
	j := &m.jdg
	state := func(i int) (icon, label, time string, done *judge.TestResult) {
		if i < len(j.tests) {
			t := j.tests[i]
			return verdictIcon(t.Verdict), badge(string(t.Verdict)), judge.FormatTime(t.Time), &t
		}
		switch {
		case j.running && i == j.current && !j.checking:
			return spinner(m.frame), fg(pulse(m.frame, hexWarm, "#7A5A10")).Render("rodando…"), "", nil
		case j.running:
			return sFaint.Render("·"), sFaint.Render("na fila"), "", nil
		}
		return sFaint.Render("–"), sFaint.Render("não executado"), "", nil
	}
	if j.total <= 14 {
		var rows []string
		for i := 0; i < j.total; i++ {
			icon, label, time, _ := state(i)
			kind := "oculto "
			if i < j.ex.Samples {
				kind = "exemplo"
			}
			rows = append(rows, fmt.Sprintf("%s  %s  %s  %s  %s", icon, sText.Render(fmt.Sprintf("Teste %2d", i+1)),
				sFaint.Render(kind), lipgloss.NewStyle().Width(9).Align(lipgloss.Right).Render(sMuted.Render(time)), label))
		}
		return rows
	}
	var rows []string
	line := ""
	for i := 0; i < j.total; i++ {
		_, _, _, t := state(i)
		bg, fgc := hexTrack, hexFaint
		switch {
		case t != nil:
			bg, fgc = verdictHex[string(t.Verdict)], hexInk
		case j.running && i == j.current:
			bg, fgc = pulse(m.frame, hexWarm, "#7A5A10"), hexInk
		}
		line += lipgloss.NewStyle().Background(lipgloss.Color(bg)).Foreground(lipgloss.Color(fgc)).Bold(true).
			Render(fmt.Sprintf(" %02d ", i+1)) + " "
		if (i+1)%10 == 0 {
			rows = append(rows, line, "")
			line = ""
		}
	}
	if line != "" {
		rows = append(rows, line)
	}
	return rows
}

func (m *model) resultCard() string {
	j := &m.jdg
	const width = 52
	card := resultBanner(j.ex, m.env.Next(m.stats), *j.result, j.samples, j.solvedBefore, m.frame, width)
	if j.result.Verdict != judge.AC || len(j.particles) == 0 {
		return card
	}
	const fh = 6
	grid := make([][]string, fh)
	for y := range grid {
		grid[y] = strings.Split(strings.Repeat(" ", width+4), "")
	}
	for _, p := range j.particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < width+4 && y >= 0 && y < fh {
			grid[y][x] = fg(p.hex).Render(p.ch)
		}
	}
	var field []string
	for _, row := range grid {
		field = append(field, strings.Join(row, ""))
	}
	return card + "\n" + strings.Join(field, "\n")
}

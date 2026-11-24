package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mattfeg/codekub/internal/core"
	"github.com/mattfeg/codekub/internal/store"
	"github.com/mattfeg/codekub/internal/track"
)

type homeState struct {
	bars []anim // one per phase, then one per topic of the current phase
}

func (h *homeState) enter(m *model) {
	h.bars = nil
	for _, p := range m.env.Track.Phases {
		h.bars = append(h.bars, anim{target: float64(core.Solved(p, m.stats)) / float64(max(1, len(p.Exercises)))})
	}
	if p := m.currentPhase(); p != nil {
		for _, tp := range p.Topics {
			s, n := topicProgress(p, tp.Slug, m.stats)
			if n > 0 {
				h.bars = append(h.bars, anim{target: float64(s) / float64(n)})
			}
		}
	}
}

func (h *homeState) step() {
	for i := range h.bars {
		h.bars[i].step()
	}
}

func (h *homeState) pos(i int) float64 {
	if i < len(h.bars) {
		return h.bars[i].pos
	}
	return 0
}

func topicProgress(p *track.Phase, topic string, stats map[string]store.ExerciseStat) (solved, total int) {
	for _, ex := range p.Exercises {
		if ex.Topic == topic {
			total++
			if stats[ex.Key()].Solved {
				solved++
			}
		}
	}
	return solved, total
}

// currentPhase is the first unlocked phase that is not complete yet.
func (m *model) currentPhase() *track.Phase {
	for i, p := range m.env.Track.Phases {
		if !m.env.Unlocked(i, m.stats) {
			break
		}
		if core.Solved(p, m.stats) < len(p.Exercises) {
			return p
		}
	}
	if n := len(m.env.Track.Phases); n > 0 {
		return m.env.Track.Phases[n-1]
	}
	return nil
}

func (m *model) homeKey(key string) tea.Cmd {
	ex := m.env.Continue(m.stats)
	if ex == nil {
		return nil
	}
	switch key {
	case "enter":
		m.openExercise(ex, "")
		m.navigate(scrProblem)
	case "w":
		m.openExercise(ex, "")
		m.navigate(scrJudge)
		m.notify("watch ligado · salve o arquivo para julgar")
		return m.startWatch()
	}
	return nil
}

func (m *model) homeView() string {
	greet := mascot("idle", m.frame) + "  " + sBold.Render("Olá, "+m.name+"!") + sMuted.Render("  Bora para o próximo?")

	var cont string
	if ex := m.env.Continue(m.stats); ex == nil {
		cont = panel(hexMint).Width(58).Render(mascot("happy", m.frame) + "  " + sBold.Render("Você resolveu tudo o que está liberado!"))
	} else {
		lines := []string{sTitle.Render("CONTINUE DE ONDE PAROU"), "", exerciseLine(ex), sMuted.Render(placeLine(ex))}
		if st, ok := m.stats[ex.Key()]; ok {
			lines = append(lines, "", badge(st.LastVerdict)+sMuted.Render(fmt.Sprintf("  na última tentativa  ·  %d tentativa(s)", st.Attempts)))
		} else if len(ex.Concepts) > 0 {
			lines = append(lines, "", sMuted.Render("conceitos  ")+sText.Render(strings.Join(ex.Concepts, ", ")))
		}
		lines = append(lines, "", helpLine("enter", "abrir", "w", "modo watch"))
		cont = panel(hexPrimary).Width(58).Render(strings.Join(lines, "\n"))
	}

	subs, act := m.subs, m.activity
	solved, total := 0, len(m.env.Track.Exercises())
	for _, st := range m.stats {
		if st.Solved {
			solved++
		}
	}
	counts := verdictCounts(subs)
	stat := func(value, label string) string {
		return fg(hexMint).Bold(true).Render(fmt.Sprintf("%-7s", value)) + sMuted.Render(label)
	}
	stats := panel(hexBorder).Width(32).Render(strings.Join([]string{
		sTitle.Render("RESUMO"), "",
		stat(fmt.Sprintf("%d/%d", solved, total), "resolvidos"),
		stat(fmt.Sprintf("%d%%", percent(counts[0].N, len(subs))), "aceitação"),
		stat(fmt.Sprint(store.Streak(act, time.Now())), streakLabel(store.Streak(act, time.Now()))),
		stat(fmt.Sprint(len(subs)), "submissões"),
	}, "\n"))

	top := lipgloss.JoinHorizontal(lipgloss.Top, cont, "  ", stats)
	if m.w < 100 {
		top = cont
	}

	rows := []string{sTitle.Render("SUA TRILHA"), ""}
	here := m.currentPhase()
	if ex := m.env.Continue(m.stats); ex != nil {
		here = ex.Phase
	}
	for i, p := range m.env.Track.Phases {
		s := core.Solved(p, m.stats)
		status := ""
		switch {
		case s == len(p.Exercises):
			status = fg(verdictHex["AC"]).Render("✔ concluída")
		case p == here:
			status = fg(pulse(m.frame, hexWarm, "#8A6A1A")).Bold(true).Render("◀ você está aqui")
		case !m.env.Unlocked(i, m.stats):
			status = sFaint.Render("· travada")
		}
		rows = append(rows, fmt.Sprintf("%s  %s  %s  %s  %s", sFaint.Render(p.ID),
			lipgloss.NewStyle().Width(16).Render(sText.Render(p.Title)), gradientBar(32, m.home.pos(i), hexPrimary, hexMint),
			sMuted.Render(fmt.Sprintf("%3d/%-3d", s, len(p.Exercises))), status))
	}
	if p := here; p != nil && len(p.Topics) > 1 {
		rows = append(rows, "", sMuted.Render("Assuntos de "+p.Title))
		bar := len(m.env.Track.Phases)
		for _, tp := range p.Topics {
			s, n := topicProgress(p, tp.Slug, m.stats)
			if n == 0 {
				continue
			}
			rows = append(rows, fmt.Sprintf("    %s  %s  %s", lipgloss.NewStyle().Width(26).Render(sText.Render(tp.Title)),
				gradientBar(24, m.home.pos(bar), hexPrimary, hexMint), sMuted.Render(fmt.Sprintf("%d/%d", s, n))))
			bar++
		}
	}
	return strings.Join([]string{"", greet, "", top, "", strings.Join(rows, "\n")}, "\n")
}

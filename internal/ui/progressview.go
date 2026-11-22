package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/mattfeg/codekub/internal/core"
	"github.com/mattfeg/codekub/internal/store"
	"github.com/mattfeg/codekub/internal/track"
)

type progressState struct {
	verdictBars []anim
	topicBars   []anim
	started     int
}

func (p *progressState) enter(m *model) {
	counts := verdictCounts(m.subs)
	top := max(1, maxCount(counts))
	p.verdictBars = make([]anim, len(counts))
	for i, c := range counts {
		p.verdictBars[i] = anim{target: float64(c.N) / float64(top)}
	}
	p.topicBars = nil
	for _, ph := range m.env.Track.Phases {
		for _, tp := range ph.Topics {
			if s, n := topicProgress(ph, tp.Slug, m.stats); n > 0 {
				p.topicBars = append(p.topicBars, anim{target: float64(s) / float64(n)})
			}
		}
	}
	p.started = 0
}

func (p *progressState) step() {
	p.started++
	for i := range p.verdictBars {
		if p.started > i*3 {
			p.verdictBars[i].step()
		}
	}
	for i := range p.topicBars {
		if p.started > 10+i*2 {
			p.topicBars[i].step()
		}
	}
}

func (m *model) progressView() string {
	counts := verdictCounts(m.subs)
	solved, total := 0, len(m.env.Track.Exercises())
	for _, st := range m.stats {
		if st.Solved {
			solved++
		}
	}
	cur := m.currentPhase()
	card := func(value, label, sub string) string {
		return panel(hexBorder).Width(25).Render(
			sMuted.Render(label) + "\n" + gradient(value, hexPrimary, hexMint, float64(m.frame)/14) + "\n" + sFaint.Render(sub))
	}
	phaseSub := ""
	if cur != nil {
		phaseSub = fmt.Sprintf("%s · %d/%d", cur.Title, core.Solved(cur, m.stats), len(cur.Exercises))
	}
	stats := lipgloss.JoinHorizontal(lipgloss.Top,
		card(fmt.Sprint(solved), "resolvidos", fmt.Sprintf("de %d na trilha", total)), " ",
		card(fmt.Sprintf("%d%%", percent(counts[0].N, len(m.subs))), "aceitação", fmt.Sprintf("%d AC em %d envios", counts[0].N, len(m.subs))), " ",
		card(plural(store.Streak(m.activity, time.Now()), "dia", "dias"), "sequência", "recorde: "+plural(store.BestStreak(m.activity), "dia", "dias")), " ",
		card(phaseID(cur), "fase atual", phaseSub))

	heat := append([]string{sTitle.Render("ATIVIDADE") + sMuted.Render("  · últimas 16 semanas"), ""}, heatmap(m.activity, time.Now(), 16, m.frame)...)
	heatPanel := panel(hexBorder).Width(46).Render(strings.Join(heat, "\n"))

	vrows := []string{sTitle.Render("VEREDITOS") + sMuted.Render(fmt.Sprintf("  · %d envios", len(m.subs))), ""}
	for i, c := range counts {
		pos := 0.0
		if i < len(m.prog.verdictBars) {
			pos = m.prog.verdictBars[i].pos
		}
		bar := fg(verdictHex[c.Code]).Render(strings.Repeat("█", max(0, int(pos*30+0.5))))
		vrows = append(vrows, fmt.Sprintf("%s %s %s %s",
			lipgloss.NewStyle().Width(4).Render(fg(verdictHex[c.Code]).Bold(true).Render(c.Code)),
			lipgloss.NewStyle().Width(31).Render(bar), sText.Render(fmt.Sprintf("%3d", c.N)),
			sFaint.Render(fmt.Sprintf("%3d%%", percent(c.N, len(m.subs))))))
	}
	if len(m.subs) == 0 {
		vrows = append(vrows, "", sMuted.Render("Nenhum envio ainda. Que tal começar?"))
	}
	verdictPanel := panel(hexBorder).Width(56).Render(strings.Join(vrows, "\n"))

	trows := []string{sTitle.Render("POR ASSUNTO"), ""}
	bar := 0
	for _, ph := range m.env.Track.Phases {
		trows = append(trows, sMuted.Render(ph.Title))
		for _, tp := range ph.Topics {
			s, n := topicProgress(ph, tp.Slug, m.stats)
			if n == 0 {
				continue
			}
			pos := 0.0
			if bar < len(m.prog.topicBars) {
				pos = m.prog.topicBars[bar].pos
			}
			bar++
			trows = append(trows, fmt.Sprintf("  %s %s %s", lipgloss.NewStyle().Width(26).Render(sText.Render(tp.Title)),
				gradientBar(16, pos, hexPrimary, hexMint), sMuted.Render(fmt.Sprintf("%d/%d", s, n))))
		}
	}
	topicPanel := panel(hexBorder).Width(56).Render(strings.Join(trows, "\n"))

	done, doing, next := conceptGroups(m.env.Track.Exercises(), m.stats)
	crows := []string{sTitle.Render("CONCEITOS"), ""}
	for _, g := range []struct {
		title, hex string
		items      []string
	}{{"DOMINADOS", verdictHex["AC"], done}, {"EM ANDAMENTO", hexWarm, doing}, {"A SEGUIR", hexMuted, next}} {
		if len(g.items) == 0 {
			continue
		}
		items := g.items
		if len(items) > 8 {
			items = append(items[:8:8], fmt.Sprintf("+%d", len(g.items)-8))
		}
		crows = append(crows, fg(g.hex).Bold(true).Render(g.title))
		crows = append(crows, chips(items, g.hex, 42)...)
		crows = append(crows, "")
	}
	conceptPanel := panel(hexBorder).Width(46).Render(strings.TrimRight(strings.Join(crows, "\n"), "\n"))

	row2 := lipgloss.JoinHorizontal(lipgloss.Top, heatPanel, " ", verdictPanel)
	row3 := lipgloss.JoinHorizontal(lipgloss.Top, conceptPanel, " ", topicPanel)
	return strings.Join([]string{"", stats, row2, row3}, "\n")
}

func phaseID(p *track.Phase) string {
	if p == nil {
		return "—"
	}
	return p.ID
}

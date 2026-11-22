package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mattfeg/codekub/internal/core"
)

type trackState struct {
	phase, ex, focus int
}

// focusNext points the track screen at the exercise the user should do next.
func (t *trackState) focusNext(m *model) {
	next := m.env.Continue(m.stats)
	for pi, p := range m.env.Track.Phases {
		for ei, ex := range p.Exercises {
			if ex == next {
				t.phase, t.ex = pi, ei
				return
			}
		}
	}
}

func (m *model) trackKey(key string) tea.Cmd {
	t := &m.trk
	phases := m.env.Track.Phases
	p := phases[t.phase]
	if m.w < 100 && m.env.Unlocked(t.phase, m.stats) {
		t.focus = 1 // the phase panel is hidden on narrow terminals
	}
	switch key {
	case "up", "k":
		if t.focus == 0 {
			t.phase = (t.phase - 1 + len(phases)) % len(phases)
			t.ex = 0
		} else if len(p.Exercises) > 0 {
			t.ex = (t.ex - 1 + len(p.Exercises)) % len(p.Exercises)
		}
	case "down", "j":
		if t.focus == 0 {
			t.phase = (t.phase + 1) % len(phases)
			t.ex = 0
		} else if len(p.Exercises) > 0 {
			t.ex = (t.ex + 1) % len(p.Exercises)
		}
	case "right", "l", "enter":
		switch {
		case !m.env.Unlocked(t.phase, m.stats):
			m.notify(fmt.Sprintf("fase travada · resolva %d%% da anterior", int(core.UnlockRatio*100)))
		case t.focus == 1 && len(p.Exercises) > 0:
			m.openExercise(p.Exercises[t.ex], "")
			m.navigate(scrProblem)
		default:
			t.focus = 1
		}
	case "left", "h", "esc":
		t.focus = 0
	}
	return nil
}

func (m *model) trackView(h int) string {
	t := m.trk
	leftBorder, rightBorder := hexPrimary, hexBorder
	if t.focus == 1 {
		leftBorder, rightBorder = hexBorder, hexPrimary
	}
	panelH := max(8, h-3)

	left := []string{sTitle.Render("FASES"), ""}
	for i, p := range m.env.Track.Phases {
		cursor := "  "
		if i == t.phase {
			cursor = fg(hexMint).Bold(true).Render("› ")
		}
		s := core.Solved(p, m.stats)
		icon, style := fg(hexWarm).Render("●"), sText
		switch {
		case s == len(p.Exercises):
			icon = fg(verdictHex["AC"]).Render("✔")
		case !m.env.Unlocked(i, m.stats):
			icon, style = sFaint.Render("·"), sFaint
		}
		if i == t.phase {
			style = style.Bold(true)
		}
		left = append(left, fmt.Sprintf("%s%s %s %s %s", cursor, icon, sFaint.Render(p.ID),
			lipgloss.NewStyle().Width(14).Render(style.Render(p.Title)), sMuted.Render(fmt.Sprintf("%3d/%d", s, len(p.Exercises)))))
	}
	leftPanel := panel(leftBorder).Width(36).Height(panelH).Render(strings.Join(left, "\n"))

	p := m.env.Track.Phases[t.phase]
	s := core.Solved(p, m.stats)
	head := []string{
		sTitle.Render("FASE "+p.ID) + "  " + sBold.Render(p.Title),
		sMuted.Render(p.Description), "",
		gradientBar(40, float64(s)/float64(max(1, len(p.Exercises))), hexPrimary, hexMint) + "  " + sMuted.Render(fmt.Sprintf("%d de %d", s, len(p.Exercises))),
		"",
	}
	var body []string
	cursorLine := 0
	if !m.env.Unlocked(t.phase, m.stats) {
		prev := m.env.Track.Phases[t.phase-1]
		need := int(core.UnlockRatio*float64(len(prev.Exercises))+0.999) - core.Solved(prev, m.stats)
		body = append(body, "", mascot("sleep", m.frame)+"  "+sMuted.Render("Esta fase ainda está travada."), "",
			sMuted.Render("Faltam ")+fg(hexWarm).Bold(true).Render(fmt.Sprintf("%d exercício(s)", need))+sMuted.Render(" de ")+sText.Render(prev.Title)+sMuted.Render("."))
	}
	topic := ""
	next := m.env.Continue(m.stats)
	for i, ex := range p.Exercises {
		if ex.Topic != topic {
			topic = ex.Topic
			if len(body) > 0 {
				body = append(body, "")
			}
			body = append(body, sMuted.Bold(true).Render(strings.ToUpper(ex.TopicTitle())))
		}
		selected := t.focus == 1 && i == t.ex
		cursor := "  "
		if selected {
			cursor = fg(hexMint).Bold(true).Render("› ")
			cursorLine = len(body)
		}
		st, ok := m.stats[ex.Key()]
		title := sText.Render(ex.Title)
		if selected {
			title = sBold.Render(ex.Title)
		}
		if ex == next {
			title += " " + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(hexInk)).
				Background(lipgloss.Color(pulse(m.frame, hexWarm, "#B8860B"))).Padding(0, 1).Render("PRÓXIMO")
		}
		v := "    "
		if ok {
			v = badge(st.LastVerdict)
		}
		body = append(body, fmt.Sprintf("%s%s  %s  %s %s %s", cursor, statusIcon(st, ok), sFaint.Render(ex.ID),
			lipgloss.NewStyle().Width(38).Render(title), lipgloss.NewStyle().Width(5).Render(stars(ex.Level)), v))
	}

	// Scroll the list so the cursor stays visible.
	visible := panelH - len(head)
	start := 0
	if cursorLine >= visible-2 {
		start = min(cursorLine-visible+3, max(0, len(body)-visible))
	}
	end := min(len(body), start+visible)
	window := body[start:end]
	if start > 0 {
		window[0] = sFaint.Render("  ↑ mais acima")
	}
	if end < len(body) {
		window[len(window)-1] = sFaint.Render("  ↓ mais abaixo")
	}
	if m.w < 100 { // narrow terminal: the exercise list gets the whole width
		return "\n" + panel(hexPrimary).Width(m.w-6).Height(panelH).Render(strings.Join(append(head, window...), "\n"))
	}
	rightPanel := panel(rightBorder).Width(m.w - 46).Height(panelH).Render(strings.Join(append(head, window...), "\n"))
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, " ", rightPanel)
}

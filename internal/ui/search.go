package ui

import (
	"fmt"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"golang.org/x/text/unicode/norm"

	"github.com/mattfeg/codekub/internal/track"
)

// searchState is the "/" box that finds a problem by number or title.
type searchState struct {
	open  bool
	query string
	sel   int
}

const searchMax = 8

// fold lowercases and drops accents, so "area" finds "Área".
func fold(s string) string {
	var b strings.Builder
	for _, r := range norm.NFD.String(strings.ToLower(s)) {
		if !unicode.Is(unicode.Mn, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// searchResults ranks an exact number first, then numbers starting with the
// query, then titles containing it.
func (m *model) searchResults() []*track.Exercise {
	q := strings.TrimSpace(fold(m.srch.query))
	if q == "" {
		return nil
	}
	var exact, prefix, title []*track.Exercise
	for _, ex := range m.env.Track.Exercises() {
		id := strings.TrimLeft(ex.ID, "0")
		switch {
		case id == strings.TrimLeft(q, "0") || ex.ID == q:
			exact = append(exact, ex)
		case strings.HasPrefix(ex.ID, q) || strings.HasPrefix(id, q):
			prefix = append(prefix, ex)
		case strings.Contains(fold(ex.Title), q):
			title = append(title, ex)
		}
	}
	return append(append(exact, prefix...), title...)
}

func (m *model) openSearch() {
	m.srch = searchState{open: true}
}

func (m *model) searchKey(key string) tea.Cmd {
	s := &m.srch
	results := m.searchResults()
	switch key {
	case "esc":
		s.open = false
	case "enter":
		if len(results) == 0 {
			return nil
		}
		s.open = false
		m.goToExercise(results[min(s.sel, len(results)-1)])
	case "up", "ctrl+p":
		s.sel = max(0, s.sel-1)
	case "down", "ctrl+n":
		s.sel = min(min(len(results), searchMax)-1, s.sel+1)
	case "backspace":
		if r := []rune(s.query); len(r) > 0 {
			s.query, s.sel = string(r[:len(r)-1]), 0
		}
	case "space", " ":
		s.query += " "
	default:
		if len([]rune(key)) == 1 {
			s.query, s.sel = s.query+key, 0
		}
	}
	return nil
}

// goToExercise opens ex and points the track screen at it, so going back
// lands on the problem in its phase.
func (m *model) goToExercise(ex *track.Exercise) {
	for pi, p := range m.env.Track.Phases {
		for ei, e := range p.Exercises {
			if e == ex {
				m.trk.phase, m.trk.ex, m.trk.focus = pi, ei, 1
			}
		}
	}
	m.openExercise(ex, "")
	if m.screen == scrProblem || m.screen == scrJudge {
		m.switchTo(scrProblem)
	} else {
		m.navigate(scrProblem)
	}
}

func (m *model) searchView(h int) string {
	const w = 60
	caret := " "
	if m.frame/15%2 == 0 {
		caret = fg(hexMint).Render("▏")
	}
	lines := []string{
		sTitle.Render("BUSCAR PROBLEMA"), "",
		sMuted.Render("número ou título: ") + sBold.Render(m.srch.query) + caret, "",
	}
	results := m.searchResults()
	switch {
	case m.srch.query == "":
		lines = append(lines, sFaint.Render("ex.: 1001, 2165 ou soma"))
	case len(results) == 0:
		lines = append(lines, sMuted.Render("nenhum problema encontrado"))
	}
	for i, ex := range results[:min(len(results), searchMax)] {
		cursor, title := "  ", sText.Render(xansi.Truncate(ex.Title, w-24, "…"))
		if i == m.srch.sel {
			cursor, title = fg(hexMint).Bold(true).Render("› "), sBold.Render(xansi.Truncate(ex.Title, w-24, "…"))
		}
		st, ok := m.stats[ex.Key()]
		lines = append(lines, cursor+statusIcon(st, ok)+"  "+sFaint.Render(ex.ID)+"  "+title+"  "+sFaint.Render(ex.Phase.Title))
	}
	if n := len(results) - searchMax; n > 0 {
		lines = append(lines, sFaint.Render(fmt.Sprintf("   … e mais %d", n)))
	}
	lines = append(lines, "", sFaint.Render("↑↓ escolher  ·  enter abrir  ·  esc fechar"))
	box := panel(hexPrimary).Width(w).Padding(1, 3).Render(strings.Join(lines, "\n"))
	return lipgloss.Place(m.w-4, h, lipgloss.Center, lipgloss.Center, box)
}

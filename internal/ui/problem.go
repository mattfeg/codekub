package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mattfeg/codekub/internal/track"
)

type problemState struct {
	ex       *track.Exercise
	rendered []string
	key      string // exercise and width the cache was rendered for
	offset   int
}

const sideW = 42

func (p *problemState) lines(width int) []string {
	key := fmt.Sprintf("%s@%d", p.ex.Key(), width)
	if p.key != key {
		md, err := p.ex.Statement()
		if err != nil {
			md = "Não consegui ler o enunciado: " + err.Error()
		}
		p.rendered, p.key = strings.Split(renderStatement(md, width, false), "\n"), key
	}
	return p.rendered
}

func (m *model) problemKey(key string) tea.Cmd {
	switch key {
	case "up", "k":
		m.prob.offset = max(0, m.prob.offset-1)
	case "down", "j":
		m.prob.offset++
	case "pgdown", " ":
		m.prob.offset += 10
	case "pgup":
		m.prob.offset = max(0, m.prob.offset-10)
	case "e":
		return m.openEditor()
	case "s":
		m.navigate(scrJudge)
		return m.startJudge(false)
	case "x":
		m.navigate(scrJudge)
		return m.startJudge(true)
	case "w":
		m.navigate(scrJudge)
		m.notify("watch ligado · salve o arquivo para julgar")
		return m.startWatch()
	case "esc":
		m.back()
	}
	return nil
}

func (m *model) problemView(h int) string {
	ex := m.prob.ex
	mdW := max(50, m.w-sideW-12)
	lines := m.prob.lines(mdW)
	viewH := h - 2
	m.prob.offset = min(m.prob.offset, max(0, len(lines)-viewH))
	end := min(len(lines), m.prob.offset+viewH)
	visible := append([]string{}, lines[m.prob.offset:end]...)
	if end < len(lines) {
		visible = append(visible, sFaint.Render("   ↓ mais abaixo (↓ ou espaço)"))
	}
	statement := lipgloss.NewStyle().Width(mdW).Render(strings.Join(visible, "\n"))

	side := []string{
		sTitle.Render("SUA SOLUÇÃO"), "",
		sText.Render(shortPath(m.solutionPath())), "",
		sMuted.Render("limite  ") + sText.Render(limitText(ex, m.env.Lang)),
		sMuted.Render("nível   ") + stars(ex.Level), "",
	}
	side = append(side, chips(ex.Concepts, hexText, sideW-4)...)
	side = append(side, "")
	if st, ok := m.stats[ex.Key()]; ok {
		state := sMuted.Render("ainda não aceito")
		if st.Solved {
			state = fg(verdictHex["AC"]).Bold(true).Render("✔ resolvido")
		}
		side = append(side, state, sMuted.Render("última  ")+badge(st.LastVerdict)+sMuted.Render(fmt.Sprintf("  ·  tentativas %d", st.Attempts)))
	} else {
		side = append(side, sMuted.Render("ainda não enviado"))
	}
	side = append(side, "", sFaint.Render(strings.Repeat("─", sideW-4)), "",
		sKey.Render("e")+"  "+sText.Render("abrir no editor"),
		sKey.Render("s")+"  "+sText.Render("enviar a solução"),
		sKey.Render("x")+"  "+sText.Render("testar só os exemplos"),
		sKey.Render("w")+"  "+sText.Render("modo watch"))
	return lipgloss.JoinHorizontal(lipgloss.Top, statement, "  ", panel(hexBorder).Width(sideW).Render(strings.Join(side, "\n")))
}

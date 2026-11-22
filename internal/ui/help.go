package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var verdictHelp = []struct{ code, text string }{
	{"AC", "passou em todos os testes"},
	{"PE", "saída certa, mas espaços ou quebras de linha diferentes"},
	{"WA", "alguma saída não confere com a esperada"},
	{"TLE", "passou do limite de tempo (já com o fator da linguagem)"},
	{"RE", "o programa quebrou: exceção ou código de saída ≠ 0"},
	{"CE", "erro de sintaxe: a solução nem chegou a rodar"},
}

func keySection(title string, rows [][2]string) string {
	lines := []string{sTitle.Render(title)}
	for _, r := range rows {
		lines = append(lines, lipgloss.NewStyle().Width(12).Render(sKey.Render(r[0]))+sText.Render(r[1]))
	}
	return strings.Join(lines, "\n")
}

// helpView is the overlay opened with "?" from any screen.
func (m *model) helpView(h int) string {
	keys := strings.Join([]string{
		keySection("NAVEGAÇÃO", [][2]string{
			{"tab", "próxima aba"}, {"shift+tab", "aba anterior"}, {"↑↓  j k", "mover ou rolar"},
			{"→  enter", "abrir"}, {"←  esc", "voltar"}, {"?", "abrir ou fechar a ajuda"}, {"q", "sair"},
		}),
		"",
		keySection("NO EXERCÍCIO", [][2]string{
			{"e", "abrir a solução no editor"}, {"s", "enviar a solução"},
			{"x", "testar só os exemplos"}, {"w", "modo watch"},
		}),
		"",
		keySection("AJUSTES", [][2]string{
			{"tab", "até a aba Ajustes"}, {"enter", "escolher editor ou linguagem"},
		}),
		"",
		keySection("NO JULGAMENTO", [][2]string{
			{"r  s", "julgar de novo"}, {"x", "testar só os exemplos"},
			{"w", "ligar ou desligar o watch"}, {"e", "abrir no editor"},
		}),
	}, "\n")

	wrap := lipgloss.NewStyle().Width(44)
	vlines := []string{sTitle.Render("VEREDITOS")}
	for _, v := range verdictHelp {
		vlines = append(vlines, lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(7).Render(badge(v.code)),
			lipgloss.NewStyle().Width(37).Render(sText.Render(v.text))))
	}
	vlines = append(vlines, "",
		sTitle.Render("COMO FUNCIONA"),
		wrap.Render(sMuted.Render("Sua solução lê da entrada padrão e escreve na saída padrão, como no beecrowd.")),
		wrap.Render(sMuted.Render("Os primeiros testes são os exemplos do enunciado; os demais ficam ocultos.")),
		wrap.Render(sMuted.Render("Espaços no fim das linhas e linhas vazias no fim da saída são ignorados.")),
		wrap.Render(sMuted.Render("O julgamento para no primeiro teste que falhar.")),
	)
	verdicts := strings.Join(vlines, "\n")

	cols := lipgloss.JoinHorizontal(lipgloss.Top, keys, "      ", verdicts)
	if m.w < 110 {
		cols = lipgloss.JoinVertical(lipgloss.Left, keys, "", verdicts)
	}
	title := mascot("idle", m.frame) + "  " + sBold.Render("Ajuda do ") + gradient("codekub", hexPrimary, hexMint, float64(m.frame)/10)
	box := panel(hexPrimary).Padding(1, 3).Render(lipgloss.JoinVertical(lipgloss.Left,
		title, "", cols, "", sFaint.Render("esc ou ? para fechar  ·  codekub help mostra os comandos de terminal")))
	return lipgloss.Place(m.w-4, h, lipgloss.Center, lipgloss.Center, box)
}

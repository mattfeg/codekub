package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mattfeg/codekub/internal/judge"
)

// Settings rows the cursor can stop on.
const (
	setEditor = iota
	setLanguage
	setRows
)

type option struct {
	label, detail string
	value         string // editor command or language id; "" means the default editor
	disabled      bool
	custom        bool // asks for a typed command
}

type settingsState struct {
	row     int
	picker  []option // open picker, if any
	pick    int
	typing  bool // typing a custom editor command
	input   string
	langVer map[string]string // language id → installed version, "" when missing
}

var knownEditors = []struct{ cmd, detail string }{
	{"nvim", "terminal"}, {"vim", "terminal"}, {"nano", "terminal"}, {"micro", "terminal"},
	{"hx", "terminal · Helix"}, {"emacs -nw", "terminal"}, {"vi", "terminal"},
	{"code", "gráfico · abre em outra janela"}, {"code -w", "gráfico · espera fechar a aba"},
	{"cursor", "gráfico · abre em outra janela"}, {"subl", "gráfico · abre em outra janela"},
	{"zed", "gráfico · abre em outra janela"},
	{"notepad", "Bloco de Notas · Windows"}, {"notepad++", "gráfico · Windows"},
}

func (s *settingsState) enter() {
	s.picker, s.typing = nil, false
	if s.langVer == nil {
		s.langVer = map[string]string{}
		for _, l := range judge.Languages() {
			v, err := exec.Command(l.Version[0], l.Version[1:]...).CombinedOutput()
			if err == nil {
				s.langVer[l.ID] = strings.TrimSpace(strings.Split(string(v), "\n")[0])
			}
		}
	}
}

func (m *model) editorOptions() []option {
	def := "editor de texto do sistema"
	for _, key := range []string{"VISUAL", "EDITOR"} {
		if v := os.Getenv(key); v != "" {
			def = "$" + key + " = " + v
			break
		}
	}
	opts := []option{{label: "padrão", detail: def}}
	for _, e := range knownEditors {
		if _, err := exec.LookPath(strings.Fields(e.cmd)[0]); err == nil {
			opts = append(opts, option{label: e.cmd, detail: e.detail, value: e.cmd})
		}
	}
	return append(opts, option{label: "personalizado…", detail: "digite qualquer comando", custom: true})
}

func (m *model) languageOptions() []option {
	var opts []option
	for _, l := range judge.Languages() {
		v := m.set.langVer[l.ID]
		o := option{label: l.Name, detail: v, value: l.ID}
		if v == "" {
			o.detail, o.disabled = "não instalado", true
		}
		opts = append(opts, o)
	}
	return opts
}

func (m *model) settingsKey(key string) tea.Cmd {
	s := &m.set
	if s.typing {
		switch key {
		case "esc":
			s.typing = false
		case "enter":
			m.saveEditor(strings.TrimSpace(s.input))
		case "backspace":
			if r := []rune(s.input); len(r) > 0 {
				s.input = string(r[:len(r)-1])
			}
		case "space":
			s.input += " "
		default:
			if len([]rune(key)) == 1 {
				s.input += key
			}
		}
		return nil
	}
	if s.picker != nil {
		switch key {
		case "up", "k":
			s.pick = (s.pick - 1 + len(s.picker)) % len(s.picker)
		case "down", "j":
			s.pick = (s.pick + 1) % len(s.picker)
		case "esc", "left", "h":
			s.picker = nil
		case "enter":
			o := s.picker[s.pick]
			switch {
			case o.disabled:
				m.notify(o.label + " não está instalado")
			case o.custom:
				s.typing, s.input = true, ""
			case s.row == setEditor:
				m.saveEditor(o.value)
			default:
				m.saveLanguage(o.value)
			}
		}
		return nil
	}
	switch key {
	case "up", "k":
		s.row = (s.row - 1 + setRows) % setRows
	case "down", "j":
		s.row = (s.row + 1) % setRows
	case "enter", "right", "l":
		if s.row == setEditor {
			s.picker = m.editorOptions()
		} else {
			s.picker = m.languageOptions()
		}
		s.pick = 0
		for i, o := range s.picker {
			if m.isCurrent(o) {
				s.pick = i
			}
		}
	}
	return nil
}

func (m *model) isCurrent(o option) bool {
	if m.set.row == setEditor {
		return !o.custom && o.value == m.env.WS.SavedConfig().Editor
	}
	return o.value == m.env.Lang.ID
}

func (m *model) saveEditor(command string) {
	if f := strings.Fields(command); len(f) > 0 {
		if _, err := exec.LookPath(f[0]); err != nil {
			m.notify(fmt.Sprintf("não encontrei %q no PATH", f[0]))
			return
		}
	}
	cfg := m.env.WS.SavedConfig()
	cfg.Editor = command
	if err := m.env.WS.SaveConfig(cfg); err != nil {
		m.notify("não consegui salvar: " + err.Error())
		return
	}
	m.set.picker, m.set.typing = nil, false
	ed, source := m.env.WS.Editor()
	if source == "CODEKUB_EDITOR" {
		m.notify("salvo, mas CODEKUB_EDITOR tem prioridade")
		return
	}
	m.notify("editor: " + editorLabel(ed))
}

func (m *model) saveLanguage(id string) {
	lang, err := judge.LanguageByID(id)
	if err != nil {
		m.notify(err.Error())
		return
	}
	cfg := m.env.WS.SavedConfig()
	cfg.Language = lang.ID
	if err := m.env.WS.SaveConfig(cfg); err != nil {
		m.notify("não consegui salvar: " + err.Error())
		return
	}
	m.set.picker = nil
	if os.Getenv("CODEKUB_LANG") != "" {
		m.notify("salvo, mas CODEKUB_LANG tem prioridade")
		return
	}
	// Solutions are per language, so the open exercise state starts over.
	m.env.Lang = lang
	m.jdg.stopRun()
	m.jdg.stopWatch()
	m.jdg = judgeState{}
	m.notify("linguagem: " + lang.Name)
}

func (m *model) settingsView() string {
	s := &m.set
	ws := m.env.WS
	ed, edSource := ws.Editor()
	langSource := "config.json"
	if os.Getenv("CODEKUB_LANG") != "" {
		langSource = "CODEKUB_LANG"
	}

	row := func(i int, label, value, source string) string {
		cursor, style := "  ", sText
		if i == s.row {
			cursor, style = fg(hexMint).Bold(true).Render("› "), sBold
		}
		return cursor + lipgloss.NewStyle().Width(12).Render(sMuted.Render(label)) +
			lipgloss.NewStyle().Width(34).Render(style.Render(value)) + sFaint.Render(source)
	}
	info := func(label, value string) string {
		return "  " + lipgloss.NewStyle().Width(12).Render(sMuted.Render(label)) + sText.Render(value)
	}
	isolation := "sem sandbox (limites de tempo e saída)"
	if judge.SandboxAvailable() {
		isolation = judge.SandboxName()
	}
	border := hexPrimary
	if s.picker != nil {
		border = hexBorder
	}
	main := panel(border).Width(66).Render(strings.Join([]string{
		sTitle.Render("AJUSTES") + "  " + sFaint.Render(shortPath(ws.Data)+"/config.json"), "",
		row(setEditor, "editor", editorLabel(ed), edSource),
		row(setLanguage, "linguagem", m.env.Lang.Name, langSource), "",
		sFaint.Render(strings.Repeat("─", 62)), "",
		info("soluções", shortPath(ws.Solutions)),
		info("progresso", shortPath(ws.DBPath())),
		info("isolamento", isolation), "",
		sMuted.Render("enter muda a opção selecionada"),
	}, "\n"))
	if s.picker == nil {
		return "\n" + main
	}

	title := "ESCOLHA O EDITOR"
	if s.row == setLanguage {
		title = "ESCOLHA A LINGUAGEM"
	}
	lines := []string{sTitle.Render(title), ""}
	for i, o := range s.picker {
		cursor, style := "  ", sText
		if i == s.pick {
			cursor, style = fg(hexMint).Bold(true).Render("› "), sBold
		}
		if o.disabled {
			style = sFaint
		}
		mark := " "
		if m.isCurrent(o) {
			mark = fg(verdictHex["AC"]).Render("✔")
		}
		lines = append(lines, cursor+mark+" "+lipgloss.NewStyle().Width(16).Render(style.Render(o.label))+sFaint.Render(o.detail))
	}
	if s.typing {
		caret := " "
		if m.frame%20 < 10 {
			caret = fg(hexMint).Render("▌")
		}
		lines = append(lines, "", sMuted.Render("comando: ")+sText.Render(s.input)+caret,
			sFaint.Render("enter salva · esc cancela · o caminho do arquivo é acrescentado no fim"))
	} else if tip := editorTip(s.picker[s.pick].value); s.row == setEditor && tip != "" {
		lines = append(lines, "", lipgloss.NewStyle().Width(52).Render(sMuted.Render(tip)))
	}
	picker := panel(hexPrimary).Width(56).Render(strings.Join(lines, "\n"))
	if m.w < 130 {
		return "\n" + main + "\n" + picker
	}
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, main, " ", picker)
}

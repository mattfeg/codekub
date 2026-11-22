package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/mattfeg/codekub/internal/core"
	"github.com/mattfeg/codekub/internal/judge"
	"github.com/mattfeg/codekub/internal/store"
	"github.com/mattfeg/codekub/internal/track"
)

func out(lines ...string) { fmt.Println(strings.Join(lines, "\n")) }

// isTerminal tells whether stdout is a terminal, where live line updates work.
func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, fg(verdictHex["WA"]).Render("erro: ")+err.Error())
	return 2
}

// shortPath shows paths under the home directory as ~/...
func shortPath(p string) string {
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(p, home+string(os.PathSeparator)) {
		return "~" + p[len(home):]
	}
	return p
}

func exerciseLine(ex *track.Exercise) string {
	return sFaint.Render(ex.ID) + "  " + sBold.Render(ex.Title) + "  " + stars(ex.Level)
}

// CmdTest judges a solution, printing each test as it finishes.
func CmdTest(env *core.Env, ref, file string, samplesOnly bool) int {
	ex, err := env.Track.Lookup(ref)
	if err != nil {
		return fail(err)
	}
	if file == "" {
		path := env.WS.SolutionPath(ex, env.Lang)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			env.WS.EnsureSolution(ex, env.Lang)
			out("", " "+mascot("idle", 10)+"  "+sText.Render("Criei sua solução em ")+sKey.Render(shortPath(path)),
				"    "+sMuted.Render("Escreva o código e rode o teste de novo."), "")
			return 1
		}
	}
	lang, err := env.LanguageFor(file)
	if err != nil {
		return fail(err)
	}
	before := env.Stats()[ex.Key()].Solved

	mode := "submissão"
	if samplesOnly {
		mode = "só exemplos"
	}
	out("", " "+exerciseLine(ex)+sMuted.Render("  ·  "+lang.Name+"  ·  "+mode), "")
	clear := ""
	if isTerminal() {
		clear = "\033[2K"
	}
	sub, err := env.Submit(context.Background(), ex, file, samplesOnly, func(e judge.Event) {
		switch e.Kind {
		case judge.EventChecking:
			if clear != "" {
				fmt.Print(sMuted.Render("   verificando sintaxe…") + "\r")
			}
		case judge.EventTestDone:
			t := e.Test
			fmt.Printf("%s   %s  %s  %s  %s  %s\n", clear, verdictIcon(t.Verdict),
				sText.Render(fmt.Sprintf("teste %2d", t.Index+1)), sFaint.Render(fmt.Sprintf("%-7s", testKind(*t))),
				sMuted.Render(fmt.Sprintf("%9s", judge.FormatTime(t.Time))), badge(string(t.Verdict)))
		}
	})
	if err != nil {
		return fail(err)
	}
	fmt.Print(clear)
	res := sub.Result
	out("", resultBanner(ex, env.Next(env.Stats()), res, samplesOnly, before, 0, 60), "")
	if res.Verdict == judge.AC {
		return 0
	}
	return 1
}

// resultBanner is the verdict card shared by `codekub test` and the TUI.
func resultBanner(ex, next *track.Exercise, res judge.Result, samplesOnly, solvedBefore bool, frame, width int) string {
	v := string(res.Verdict)
	color := verdictHex[v]
	icon, mood := "✘", "sad"
	if res.Verdict == judge.AC {
		icon, mood = "✔", "happy"
	}
	lines := []string{mascot(mood, frame) + "   " + fg(color).Bold(true).Render(icon+"  "+strings.ToUpper(verdictName[v])), ""}
	if res.Verdict == judge.CE || res.Verdict == judge.IE {
		lines = append(lines, sMuted.Render("nenhum teste executado"))
	} else {
		lines = append(lines, sMuted.Render(fmt.Sprintf("%d/%d testes  ·  pico %s", res.Passed(), res.Total, judge.FormatTime(res.MaxTime()))))
	}
	if fl := failureLines(res); len(fl) > 0 {
		box := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color(color)).PaddingLeft(2).Width(width - 7).Render(strings.Join(fl, "\n"))
		lines = append(lines, "", box)
	}
	if res.Verdict == judge.AC {
		lines = append(lines, "")
		switch {
		case samplesOnly:
			lines = append(lines, sText.Render("Os exemplos passaram. Agora envie a solução completa."))
		case !solvedBefore:
			lines = append(lines, fg(hexMint).Bold(true).Render("+1 resolvido")+sMuted.Render("  ·  "+ex.TopicTitle()))
		default:
			lines = append(lines, sMuted.Render("Você já tinha resolvido este, mas mandou bem de novo."))
		}
		if next != nil && next != ex && !samplesOnly {
			lines = append(lines, sMuted.Render("próximo  ")+sText.Render(next.ID+" · "+next.Title))
		}
	}
	return panel(color).Width(width).Padding(1, 2).Render(strings.Join(lines, "\n"))
}

func CmdNext(env *core.Env) int {
	stats := env.Stats()
	ex := env.Next(stats)
	if ex == nil {
		out("", " "+mascot("happy", 10)+"  "+sBold.Render("Você resolveu tudo o que está liberado. Parabéns!"), "")
		return 0
	}
	path, err := env.WS.EnsureSolution(ex, env.Lang)
	if err != nil {
		return fail(err)
	}
	card := panel(hexPrimary).Width(62).Render(strings.Join([]string{
		sTitle.Render("PRÓXIMO EXERCÍCIO"), "",
		exerciseLine(ex),
		sMuted.Render(ex.Phase.Title + " · " + ex.TopicTitle()), "",
		sMuted.Render("solução  ") + sText.Render(shortPath(path)), "",
		sKey.Render("codekub show "+strings.TrimLeft(ex.ID, "0")) + sMuted.Render("   ler o enunciado"),
		sKey.Render("codekub watch "+strings.TrimLeft(ex.ID, "0")) + sMuted.Render("  julgar a cada salvamento"),
	}, "\n"))
	out("", card, "")
	return 0
}

func statusIcon(st store.ExerciseStat, ok bool) string {
	switch {
	case st.Solved:
		return fg(verdictHex["AC"]).Render("✔")
	case ok:
		return fg(verdictHex["WA"]).Render("✘")
	}
	return sFaint.Render("○")
}

func CmdList(env *core.Env) int {
	stats := env.Stats()
	for i, p := range env.Track.Phases {
		solved := core.Solved(p, stats)
		head := sTitle.Render("FASE "+p.ID) + "  " + sBold.Render(p.Title) + "  " +
			gradientBar(20, float64(solved)/float64(max(1, len(p.Exercises))), hexPrimary, hexMint) +
			sMuted.Render(fmt.Sprintf("  %d/%d", solved, len(p.Exercises)))
		if !env.Unlocked(i, stats) {
			head += sFaint.Render("  · travada")
		}
		out("", head)
		topic := ""
		for _, ex := range p.Exercises {
			if ex.Topic != topic {
				topic = ex.Topic
				out("", "  "+sMuted.Render(ex.TopicTitle()))
			}
			st, ok := stats[ex.Key()]
			last := ""
			if ok {
				last = badge(st.LastVerdict)
			}
			out(fmt.Sprintf("   %s %s  %s %s %s", statusIcon(st, ok), sFaint.Render(ex.ID),
				lipgloss.NewStyle().Width(34).Render(sText.Render(ex.Title)), lipgloss.NewStyle().Width(5).Render(stars(ex.Level)), last))
		}
	}
	out("")
	return 0
}

func CmdShow(env *core.Env, ref string) int {
	ex, err := env.Track.Lookup(ref)
	if err != nil {
		return fail(err)
	}
	md, err := ex.Statement()
	if err != nil {
		return fail(err)
	}
	path, _ := env.WS.EnsureSolution(ex, env.Lang)
	out(renderStatement(md, 80, !isTerminal()), "",
		"  "+sMuted.Render("solução  ")+sText.Render(shortPath(path)),
		"  "+sMuted.Render("limite   ")+sText.Render(limitText(ex, env.Lang))+"   "+stars(ex.Level), "")
	return 0
}

func limitText(ex *track.Exercise, lang judge.Language) string {
	base := strings.Replace(fmt.Sprintf("%g s", ex.TimeLimit), ".", ",", 1)
	eff := strings.Replace(fmt.Sprintf("%g s", ex.TimeLimit*lang.TimeFactor), ".", ",", 1)
	return fmt.Sprintf("%s (%s: %s)", base, lang.Name, eff)
}

func CmdProgress(env *core.Env) int {
	stats := env.Stats()
	subs, _ := env.Store.Submissions(1 << 30)
	act, _ := env.Store.Activity()
	solved, total := 0, len(env.Track.Exercises())
	for _, st := range stats {
		if st.Solved {
			solved++
		}
	}
	counts := verdictCounts(subs)
	out("", " "+sTitle.Render("PROGRESSO")+"  "+sMuted.Render(env.Lang.Name), "",
		fmt.Sprintf("  %s %s    %s %s    %s %s",
			fg(hexMint).Bold(true).Render(fmt.Sprintf("%d/%d", solved, total)), sMuted.Render("resolvidos"),
			fg(hexMint).Bold(true).Render(fmt.Sprintf("%d%%", percent(counts[0].N, len(subs)))), sMuted.Render("aceitação"),
			fg(hexMint).Bold(true).Render(fmt.Sprint(store.Streak(act, time.Now()))), sMuted.Render(streakLabel(store.Streak(act, time.Now())))), "")
	for _, p := range env.Track.Phases {
		out("  " + sBold.Render(p.Title))
		for _, tp := range p.Topics {
			n, s := 0, 0
			for _, ex := range p.Exercises {
				if ex.Topic == tp.Slug {
					n++
					if stats[ex.Key()].Solved {
						s++
					}
				}
			}
			if n > 0 {
				out(fmt.Sprintf("    %s %s %s", lipgloss.NewStyle().Width(28).Render(sText.Render(tp.Title)),
					gradientBar(24, float64(s)/float64(n), hexPrimary, hexMint), sMuted.Render(fmt.Sprintf("%d/%d", s, n))))
			}
		}
	}
	out("", "  "+sBold.Render("Vereditos")+sMuted.Render(fmt.Sprintf("  · %d envios", len(subs))))
	for _, c := range counts {
		bar := fg(verdictHex[c.Code]).Render(strings.Repeat("█", c.N*30/max(1, maxCount(counts))))
		out(fmt.Sprintf("    %s %s %s", lipgloss.NewStyle().Width(4).Render(fg(verdictHex[c.Code]).Bold(true).Render(c.Code)),
			lipgloss.NewStyle().Width(31).Render(bar), sMuted.Render(fmt.Sprintf("%3d", c.N))))
	}
	out("", "  "+sBold.Render("Atividade")+sMuted.Render("  · últimas 16 semanas"), "")
	for _, l := range heatmap(act, time.Now(), 16, 0) {
		out("    " + l)
	}
	out("")
	return 0
}

func maxCount(cs []verdictCount) int {
	m := 0
	for _, c := range cs {
		m = max(m, c.N)
	}
	return m
}

// CmdDoctor checks everything codekub needs to run solutions.
func CmdDoctor() int {
	ok := fg(verdictHex["AC"]).Render("✔")
	bad := fg(verdictHex["WA"]).Render("✘")
	warn := fg(hexWarm).Render("!")
	out("", " "+mascot("idle", 10)+"  "+sTitle.Render("DIAGNÓSTICO"), "")
	if root, err := track.Find(); err != nil {
		out("  " + bad + " trilha: " + err.Error())
	} else if tr, err := track.Load(root); err != nil {
		out("  " + bad + " trilha: " + err.Error())
	} else {
		out(fmt.Sprintf("  %s trilha      %s %s", ok, sText.Render(shortPath(root)), sMuted.Render(fmt.Sprintf("(%d exercícios)", len(tr.Exercises())))))
	}
	if env, err := core.Open(); err == nil {
		out("  "+ok+" soluções    "+sText.Render(shortPath(env.WS.Solutions)),
			"  "+ok+" progresso   "+sText.Render(shortPath(env.WS.DBPath())),
			"  "+ok+" linguagem   "+sText.Render(env.Lang.Name)+sMuted.Render("  (mude com CODEKUB_LANG ou ~/.codekub/config.json)"))
		env.Close()
	} else {
		out("  " + bad + " ambiente: " + err.Error())
	}
	out("", "  "+sBold.Render("Linguagens"))
	for _, l := range judge.Languages() {
		v, err := exec.Command(l.Version[0], l.Version[1:]...).CombinedOutput()
		if err != nil {
			out(fmt.Sprintf("  %s %-10s %s", bad, l.Name, sMuted.Render("não instalado (ex.: mise use -g "+l.ID+")")))
			continue
		}
		out(fmt.Sprintf("  %s %-10s %s", ok, l.Name, sMuted.Render(strings.TrimSpace(strings.Split(string(v), "\n")[0]))))
	}
	out("", "  "+sBold.Render("Isolamento"))
	if judge.SandboxAvailable() {
		out("  " + ok + " " + judge.SandboxName() + sMuted.Render(": sem rede, escrita só na pasta temporária"))
	} else {
		out("  " + warn + " " + sMuted.Render("sem sandbox: soluções rodam com limite de tempo e de saída, numa pasta temporária"))
	}
	out("  "+warn+" "+sMuted.Render("limite de memória não é aplicado no modo local"), "")
	return 0
}

func Help() string {
	arg := fg(hexMint)
	cmd := func(name, args, desc string) string {
		return "  " + lipgloss.NewStyle().Width(24).Render(sBold.Render(name)+" "+arg.Render(args)) + sMuted.Render(desc)
	}
	ex := func(s string) string { return "  " + sFaint.Render("$ ") + sText.Render(s) }
	return strings.Join([]string{
		"",
		" " + mascot("idle", 10) + "  " + gradient("codekub", hexPrimary, hexMint, 0) + sMuted.Render("  ·  by Matheus Ferreira"),
		"",
		sTitle.Render(" USO"),
		"  " + sBold.Render("codekub") + " " + arg.Render("[comando] [argumentos]"),
		"",
		sTitle.Render(" COMANDOS"),
		cmd("(nenhum)", "", "abre o painel interativo"),
		cmd("next", "", "mostra o próximo exercício da trilha"),
		cmd("list", "", "lista as fases e os exercícios"),
		cmd("show", "<n>", "mostra o enunciado do exercício n"),
		cmd("test", "<n> [arquivo]", "julga sua solução do exercício n"),
		cmd("test", "<n> --exemplos", "roda só os exemplos do enunciado"),
		cmd("watch", "<n> [arquivo]", "julga de novo a cada vez que você salvar"),
		cmd("progress", "", "mostra seu progresso"),
		cmd("doctor", "", "confere linguagens, isolamento e pastas"),
		cmd("help", "", "mostra esta ajuda"),
		"",
		sTitle.Render(" EXEMPLOS"),
		ex("codekub test 13"),
		ex("codekub watch 13 ~/meus-codigos/media.rb"),
		"",
		sTitle.Render(" ONDE FICAM AS COISAS"),
		"  " + sMuted.Render("soluções   ") + sText.Render("~/codekub/solucoes/<número>-<nome>.rb") + sFaint.Render("   (CODEKUB_SOLUTIONS)"),
		"  " + sMuted.Render("progresso  ") + sText.Render("~/.codekub/codekub.db") + sFaint.Render("                   (CODEKUB_HOME)"),
		"",
		sTitle.Render(" NO PAINEL"),
		"  " + helpLine("?", "ajuda", "tab", "abas", "enter", "abrir", "esc", "voltar", "q", "sair"),
		"",
	}, "\n")
}

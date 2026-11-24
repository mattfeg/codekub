package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/mattfeg/codekub/internal/judge"
	"github.com/mattfeg/codekub/internal/store"
	"github.com/mattfeg/codekub/internal/track"
)

// renderStatement renders markdown for the terminal; plain=true drops colors
// for output that is piped or redirected.
func renderStatement(md string, width int, plain bool) string {
	style := "dark"
	if plain {
		style = "notty"
	}
	r, err := glamour.NewTermRenderer(glamour.WithStandardStyle(style), glamour.WithWordWrap(width))
	if err != nil {
		return md
	}
	out, err := r.Render(md)
	if err != nil {
		return md
	}
	return strings.TrimRight(out, "\n")
}

func streakLabel(n int) string {
	if n == 1 {
		return "dia seguido"
	}
	return "dias seguidos"
}

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

// stars shows the level from 1 to 4; level 0 means unknown and shows nothing.
func stars(level int) string {
	if level <= 0 {
		return ""
	}
	return fg(hexWarm).Render(strings.Repeat("★", level)) + sFaint.Render(strings.Repeat("☆", max(0, 4-level)))
}

func verdictIcon(v judge.Verdict) string {
	if v == judge.AC {
		return fg(verdictHex["AC"]).Render("✔")
	}
	return fg(verdictHex[string(v)]).Render("✘")
}

func testKind(t judge.TestResult) string {
	if t.Sample {
		return "exemplo"
	}
	return "oculto"
}

// failureLines explains why a result is not AC, in the same words in the
// CLI and in the TUI. Hidden tests never reveal their input or output.
func failureLines(r judge.Result) []string {
	color := fg(verdictHex[string(r.Verdict)])
	switch r.Verdict {
	case judge.AC:
		return nil
	case judge.CE:
		return append([]string{sMuted.Render("a solução nem chegou a rodar:")}, colorLines(color, r.Detail, 8)...)
	case judge.IE:
		return []string{color.Render(r.Detail), sMuted.Render("rode ") + sKey.Render("codekub doctor") + sMuted.Render(" para conferir o ambiente")}
	}
	if len(r.Tests) == 0 {
		return nil
	}
	t := r.Tests[len(r.Tests)-1]
	head := sMuted.Render(fmt.Sprintf("teste %d · %s", t.Index+1, testKind(t)))
	switch t.Verdict {
	case judge.WA, judge.PE:
		lines := []string{head}
		if t.Verdict == judge.PE {
			lines = append(lines, sText.Render("O conteúdo está certo, mas os espaços ou as quebras de linha não."))
		}
		if t.Diff == nil {
			return lines
		}
		if !t.Sample {
			return append(lines, sText.Render(fmt.Sprintf("A saída não confere na linha %d.", t.Diff.Line)),
				sFaint.Render("Teste oculto: a entrada não é mostrada."))
		}
		show := func(s *string) string {
			if s == nil {
				return sFaint.Render("<fim da saída>")
			}
			return fmt.Sprintf("%q", *s)
		}
		input := strings.Split(strings.TrimRight(t.Input, "\n"), "\n")
		if strings.TrimSpace(t.Input) == "" {
			input = []string{sFaint.Render("(vazia)")}
		}
		if len(input) > 4 {
			input = append(input[:4], "…")
		}
		lines = append(lines, sMuted.Render("entrada   ")+sText.Render(input[0]))
		for _, l := range input[1:] {
			lines = append(lines, "          "+sText.Render(l))
		}
		return append(lines,
			sMuted.Render(fmt.Sprintf("linha %d", t.Diff.Line)),
			sMuted.Render("esperado  ")+fg(verdictHex["AC"]).Bold(true).Render(show(t.Diff.Want)),
			sMuted.Render("obtido    ")+color.Bold(true).Render(show(t.Diff.Got)))
	default:
		return append([]string{head}, colorLines(color, t.Detail, 8)...)
	}
}

func colorLines(style lipgloss.Style, text string, maxLines int) []string {
	var out []string
	for i, l := range strings.Split(strings.TrimSpace(text), "\n") {
		if i == maxLines {
			out = append(out, sFaint.Render("…"))
			break
		}
		out = append(out, style.Render(l))
	}
	return out
}

var heatLevels = []string{hexTrack, "#3F3670", "#6B5BC4", hexPrimary, hexMint}

func heatColor(n int) string {
	switch {
	case n == 0:
		return heatLevels[0]
	case n <= 2:
		return heatLevels[1]
	case n <= 4:
		return heatLevels[2]
	case n <= 7:
		return heatLevels[3]
	}
	return heatLevels[4]
}

// heatmap draws the last `weeks` weeks of activity, Monday to Sunday, with
// today pulsing when animated (frame > 0).
func heatmap(activity map[string]int, now time.Time, weeks, frame int) []string {
	names := []string{"seg", "   ", "qua", "   ", "sex", "   ", "dom"}
	todayRow := (int(now.Weekday()) + 6) % 7 // Monday = 0
	var rows []string
	for d := 0; d < 7; d++ {
		row := sFaint.Render(names[d]) + " "
		for w := 0; w < weeks; w++ {
			offset := (weeks-1-w)*7 + (todayRow - d)
			switch {
			case offset < 0:
				row += "  "
			case offset == 0 && frame > 0:
				row += fg(pulse(frame, hexMint, hexPrimary)).Render("■") + " "
			default:
				day := now.AddDate(0, 0, -offset).Format("2006-01-02")
				row += fg(heatColor(activity[day])).Render("■") + " "
			}
		}
		rows = append(rows, row)
	}
	legend := sFaint.Render("menos ")
	for _, c := range heatLevels {
		legend += fg(c).Render("■") + " "
	}
	return append(rows, "", "    "+legend+sFaint.Render("mais"))
}

type verdictCount struct {
	Code string
	N    int
}

func verdictCounts(subs []store.Submission) []verdictCount {
	counts := map[string]int{}
	for _, s := range subs {
		counts[s.Verdict]++
	}
	var out []verdictCount
	for _, v := range []string{"AC", "WA", "PE", "TLE", "RE", "CE"} {
		out = append(out, verdictCount{v, counts[v]})
	}
	return out
}

// conceptGroups splits the concepts of the open track into mastered (from
// solved exercises), in progress (attempted, unsolved) and up next.
func conceptGroups(exercises []*track.Exercise, stats map[string]store.ExerciseStat) (done, doing, next []string) {
	state := map[string]int{} // 3 done, 2 doing, 1 next
	order := []string{}
	for _, ex := range exercises {
		st, ok := stats[ex.Key()]
		level := 1
		if st.Solved {
			level = 3
		} else if ok {
			level = 2
		}
		for _, c := range ex.Concepts {
			if _, seen := state[c]; !seen {
				order = append(order, c)
			}
			state[c] = max(state[c], level)
		}
	}
	for _, c := range order {
		switch state[c] {
		case 3:
			done = append(done, c)
		case 2:
			doing = append(doing, c)
		default:
			next = append(next, c)
		}
	}
	return done, doing, next
}

func percent(part, total int) int {
	if total == 0 {
		return 0
	}
	return int(math.Round(float64(part) * 100 / float64(total)))
}

func chips(items []string, hex string, width int) []string {
	var lines []string
	line, w := "", 0
	for _, item := range items {
		chip := lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Background(lipgloss.Color(hexTrack)).Padding(0, 1).Render(item)
		if w > 0 && w+lipgloss.Width(chip)+1 > width {
			lines = append(lines, line)
			line, w = "", 0
		}
		line += chip + " "
		w += lipgloss.Width(chip) + 1
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

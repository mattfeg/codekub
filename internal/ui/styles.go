package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	hexPrimary = "#A78BFA"
	hexMint    = "#5EEAD4"
	hexText    = "#E6E6EF"
	hexMuted   = "#9494AD"
	hexFaint   = "#55556E"
	hexTrack   = "#2A2A40"
	hexBorder  = "#3A3A55"
	hexWarm    = "#FBBF24"
	hexBear    = "#D9A877"
	hexInk     = "#14141F"
)

var verdictHex = map[string]string{
	"AC": "#4ADE80", "WA": "#F87171", "PE": "#FACC15",
	"TLE": "#60A5FA", "RE": "#E879F9", "CE": "#FB923C", "IE": "#9494AD",
}

var verdictName = map[string]string{
	"AC": "Aceito", "WA": "Resposta errada", "PE": "Erro de apresentação",
	"TLE": "Tempo limite excedido", "RE": "Erro em execução", "CE": "Erro de sintaxe",
	"IE": "Erro interno",
}

func fg(hex string) lipgloss.Style { return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)) }

var (
	sText  = fg(hexText)
	sMuted = fg(hexMuted)
	sFaint = fg(hexFaint)
	sBold  = fg(hexText).Bold(true)
	sTitle = fg(hexPrimary).Bold(true)
	sKey   = fg(hexMint).Bold(true)
)

func panel(borderHex string) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderHex)).
		Padding(0, 2)
}

func badge(v string) string {
	return lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color(hexInk)).
		Background(lipgloss.Color(verdictHex[v])).
		Padding(0, 1).Render(v)
}

func blend(a, b string, t float64) string {
	ca, _ := colorful.Hex(a)
	cb, _ := colorful.Hex(b)
	return ca.BlendLuv(cb, math.Max(0, math.Min(1, t))).Clamped().Hex()
}

// gradient colors each rune of s between a and b; moving phase makes it shimmer.
func gradient(s, a, b string, phase float64) string {
	var sb strings.Builder
	for i, r := range []rune(s) {
		if r == ' ' {
			sb.WriteRune(r)
			continue
		}
		t := (math.Sin(float64(i)/7-phase) + 1) / 2
		sb.WriteString(fg(blend(a, b, t)).Bold(true).Render(string(r)))
	}
	return sb.String()
}

var partials = []rune("▏▎▍▌▋▊▉")

// gradientBar draws a smooth bar with sub-cell precision.
func gradientBar(width int, pct float64, a, b string) string {
	pct = math.Max(0, math.Min(1, pct))
	filled := pct * float64(width)
	var sb strings.Builder
	for i := 0; i < width; i++ {
		c := blend(a, b, float64(i)/math.Max(1, float64(width-1)))
		switch {
		case float64(i+1) <= filled:
			sb.WriteString(fg(c).Render("█"))
		case float64(i) < filled:
			idx := int((filled-float64(i))*float64(len(partials))) - 1
			if idx < 0 {
				idx = 0
			}
			sb.WriteString(fg(c).Background(lipgloss.Color(hexTrack)).Render(string(partials[idx])))
		default:
			sb.WriteString(fg(hexTrack).Render("█"))
		}
	}
	return sb.String()
}

func helpLine(pairs ...string) string {
	var parts []string
	for i := 0; i+1 < len(pairs); i += 2 {
		parts = append(parts, sKey.Render(pairs[i])+" "+sMuted.Render(pairs[i+1]))
	}
	return strings.Join(parts, sFaint.Render("  ·  "))
}

// anim is a value that chases its target with spring physics.
type anim struct{ pos, vel, target float64 }

var spring = harmonica.NewSpring(harmonica.FPS(30), 5.5, 0.62)

func (a *anim) step() { a.pos, a.vel = spring.Update(a.pos, a.vel, a.target) }

// mascot is the kub: a bear cub that reacts to what is happening.
func mascot(mood string, frame int) string {
	face := map[string]string{
		"idle":  "ʕ•ᴥ•ʔ",
		"happy": "ʕ^ᴥ^ʔ",
		"sad":   "ʕ;ᴥ;ʔ",
		"sleep": "ʕ-ᴥ-ʔ",
		"think": "ʕ•ᴥ•ʔ",
	}[mood]
	if mood == "idle" && frame%96 < 4 {
		face = "ʕ-ᴥ-ʔ"
	}
	if mood == "think" {
		face = []string{"ʕ•ᴥ•ʔ", "ʕ•ᴥ•ʔ", "ʕ°ᴥ°ʔ", "ʕ•ᴥ•ʔ"}[(frame/8)%4]
	}
	return fg(hexBear).Bold(true).Render(face)
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func spinner(frame int) string {
	return fg(hexWarm).Render(spinnerFrames[(frame/2)%len(spinnerFrames)])
}

func pulse(frame int, a, b string) string {
	return blend(a, b, (math.Sin(float64(frame)/5)+1)/2)
}

package judge

import "strings"

type Verdict string

const (
	AC  Verdict = "AC" // accepted
	PE  Verdict = "PE" // right tokens, wrong spacing or line breaks
	WA  Verdict = "WA" // wrong answer
	TLE Verdict = "TLE"
	RE  Verdict = "RE" // runtime error
	CE  Verdict = "CE" // syntax or compilation error
	IE  Verdict = "IE" // internal error: the judge itself could not run the solution
)

// Diff points at the first line where the output differs from the expected one.
// A nil Want or Got means that side had already ended.
type Diff struct {
	Line      int
	Want, Got *string
}

// Compare ignores trailing spaces on each line and blank lines at the end.
// When only the whitespace layout differs, the verdict is PE.
func Compare(got, want string) (Verdict, *Diff) {
	g, w := lines(got), lines(want)
	if equal(g, w) {
		return AC, nil
	}
	if equal(strings.Fields(got), strings.Fields(want)) {
		return PE, firstDiff(g, w)
	}
	return WA, firstDiff(g, w)
}

func lines(s string) []string {
	ls := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for i := range ls {
		ls[i] = strings.TrimRight(ls[i], " \t\r")
	}
	for len(ls) > 0 && ls[len(ls)-1] == "" {
		ls = ls[:len(ls)-1]
	}
	return ls
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func firstDiff(got, want []string) *Diff {
	for i := 0; i < max(len(got), len(want)); i++ {
		d := &Diff{Line: i + 1}
		if i < len(want) {
			d.Want = &want[i]
		}
		if i < len(got) {
			d.Got = &got[i]
		}
		if d.Want == nil || d.Got == nil || *d.Want != *d.Got {
			return d
		}
	}
	return nil
}

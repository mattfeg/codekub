package judge

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	cases := []struct {
		got, want string
		verdict   Verdict
	}{
		{"8\n", "8\n", AC},
		{"8", "8\n", AC},
		{"8   \n\n\n", "8\n", AC},
		{"a b\r\nc\r\n", "a b\nc\n", AC},
		{" 8\n", "8\n", PE},
		{"a  b\n", "a b\n", PE},
		{"a\nb\n", "a b\n", PE},
		{"9\n", "8\n", WA},
		{"", "8\n", WA},
		{"8\n9\n", "8\n", WA},
	}
	for _, c := range cases {
		if v, _ := Compare(c.got, c.want); v != c.verdict {
			t.Errorf("Compare(%q, %q) = %s, want %s", c.got, c.want, v, c.verdict)
		}
	}
	_, d := Compare("1\n2\n", "1\n3\n")
	if d == nil || d.Line != 2 || *d.Want != "3" || *d.Got != "2" {
		t.Errorf("unexpected diff %+v", d)
	}
	_, d = Compare("1\n", "1\n2\n")
	if d == nil || d.Line != 2 || d.Got != nil {
		t.Errorf("missing line should have nil Got, got %+v", d)
	}
}

// sumTests writes a tiny "add two numbers" exercise and returns its tests.
func sumTests(t *testing.T) []TestCase {
	dir := t.TempDir()
	var tests []TestCase
	for i, c := range [][2]string{{"3 5\n", "8\n"}, {"-1 1\n", "0\n"}, {"10 20\n", "30\n"}} {
		in := filepath.Join(dir, string(rune('a'+i))+".in")
		out := filepath.Join(dir, string(rune('a'+i))+".out")
		os.WriteFile(in, []byte(c[0]), 0o644)
		os.WriteFile(out, []byte(c[1]), 0o644)
		tests = append(tests, TestCase{InputPath: in, OutputPath: out, Sample: i == 0})
	}
	return tests
}

func judgeRuby(t *testing.T, code string, limit time.Duration) Result {
	t.Helper()
	src := filepath.Join(t.TempDir(), "solucao.rb")
	os.WriteFile(src, []byte(code), 0o644)
	ruby, _ := LanguageByID("ruby")
	ruby.TimeFactor = 1
	var events []EventKind
	res := Run(context.Background(), LocalRunner{}, Request{
		Lang: ruby, Source: src, Tests: sumTests(t), TimeLimit: limit, StopOnFail: true,
	}, func(e Event) { events = append(events, e.Kind) })
	if len(events) == 0 || events[len(events)-1] != EventDone {
		t.Fatalf("events must end with EventDone, got %v", events)
	}
	return res
}

func TestVerdicts(t *testing.T) {
	cases := []struct {
		name, code string
		verdict    Verdict
		passed     int
	}{
		{"accepted", "a, b = gets.split.map(&:to_i)\nputs a + b\n", AC, 3},
		{"wrong answer", "a, b = gets.split.map(&:to_i)\nputs a - b\n", WA, 0},
		{"presentation", "a, b = gets.split.map(&:to_i)\nputs \"  #{a + b}\"\n", PE, 0},
		{"runtime error", "raise 'boom'\n", RE, 0},
		{"syntax error", "puts (1 + \n", CE, 0},
		{"time limit", "a, b = gets.split.map(&:to_i)\nsleep 5 if a == -1\nputs a + b\n", TLE, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := judgeRuby(t, c.code, 2*time.Second)
			if res.Verdict != c.verdict || res.Passed() != c.passed {
				t.Fatalf("verdict %s passed %d (detail %q), want %s passed %d", res.Verdict, res.Passed(), res.Detail, c.verdict, c.passed)
			}
		})
	}
}

func TestRuntimeErrorShowsMessage(t *testing.T) {
	res := judgeRuby(t, "raise 'boom'\n", 2*time.Second)
	if !strings.Contains(res.Tests[0].Detail, "boom") || strings.Contains(res.Tests[0].Detail, os.TempDir()) {
		t.Fatalf("detail should show the error without temp paths: %q", res.Tests[0].Detail)
	}
}

func TestTimeLimitKillsChildren(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "still-alive")
	code := "pid = spawn('sleep 3; touch " + marker + "')\nsleep 10\n"
	start := time.Now()
	res := judgeRuby(t, code, 500*time.Millisecond)
	if res.Verdict != TLE {
		t.Fatalf("got %s", res.Verdict)
	}
	if time.Since(start) > 3*time.Second {
		t.Fatalf("judge waited too long: %s", time.Since(start))
	}
	time.Sleep(3500 * time.Millisecond)
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("child process survived the time limit")
	}
}

func TestOutputLimit(t *testing.T) {
	res := judgeRuby(t, "loop { puts 'x' * 1000 }\n", 10*time.Second)
	if res.Verdict != RE || !strings.Contains(res.Tests[0].Detail, "MB") {
		t.Fatalf("got %s %q", res.Verdict, res.Tests[0].Detail)
	}
}

func TestSandbox(t *testing.T) {
	if !SandboxAvailable() {
		t.Skip("sandbox not available on this machine")
	}
	outside := filepath.Join(os.Getenv("HOME"), ".codekub-sandbox-probe")
	defer os.Remove(outside)
	code := `begin; File.write("` + outside + `", "x"); rescue; end
begin; require "socket"; TCPSocket.new("1.1.1.1", 80); rescue; end
File.write("inside.txt", "ok")
a, b = gets.split.map(&:to_i)
puts a + b
`
	if res := judgeRuby(t, code, 3*time.Second); res.Verdict != AC {
		t.Fatalf("solution writing in its own dir should pass, got %s %q", res.Verdict, res.Detail)
	}
	if _, err := os.Stat(outside); err == nil {
		t.Fatal("sandbox allowed a write outside the temp dir")
	}
}

func TestCancel(t *testing.T) {
	src := filepath.Join(t.TempDir(), "s.rb")
	os.WriteFile(src, []byte("sleep 10\n"), 0o644)
	ruby, _ := LanguageByID("ruby")
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(700 * time.Millisecond); cancel() }()
	start := time.Now()
	res := Run(ctx, LocalRunner{}, Request{Lang: ruby, Source: src, Tests: sumTests(t), TimeLimit: 5 * time.Second}, nil)
	if !res.Canceled || time.Since(start) > 3*time.Second {
		t.Fatalf("canceled=%v after %s", res.Canceled, time.Since(start))
	}
}

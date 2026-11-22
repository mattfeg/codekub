// Package track loads the exercise track from disk: phases (phase.json) and
// their exercises (meta.json, statement.md, tests/).
package track

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mattfeg/codekub/internal/judge"
)

type Track struct {
	Root   string
	Phases []*Phase
}

type Phase struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Topics      []Topic `json:"topics"`
	Exercises   []*Exercise
	Dir         string
}

type Topic struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type Exercise struct {
	ID        string   `json:"id"`
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	Topic     string   `json:"topic"`
	Level     int      `json:"level"`
	Concepts  []string `json:"concepts"`
	TimeLimit float64  `json:"time_limit"`
	Samples   int      `json:"samples"`
	Rating    int      `json:"rating"`
	Source    Source   `json:"source"`
	Dir       string
	Phase     *Phase
}

type Source struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	License string `json:"license"`
	URL     string `json:"url"`
	Dataset string `json:"dataset"`
	Changes string `json:"changes"`
}

// Find locates the track directory: $CODEKUB_TRACK, then a track/ folder in
// the current directory or any parent, then one next to the executable.
func Find() (string, error) {
	if dir := os.Getenv("CODEKUB_TRACK"); dir != "" {
		return dir, nil
	}
	var starts []string
	if wd, err := os.Getwd(); err == nil {
		starts = append(starts, wd)
	}
	if exe, err := os.Executable(); err == nil {
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			starts = append(starts, filepath.Dir(real))
		}
	}
	for _, start := range starts {
		for dir := start; ; dir = filepath.Dir(dir) {
			candidate := filepath.Join(dir, "track")
			if matches, _ := filepath.Glob(filepath.Join(candidate, "*", "phase.json")); len(matches) > 0 {
				return candidate, nil
			}
			if filepath.Dir(dir) == dir {
				break
			}
		}
	}
	return "", errors.New("não encontrei a pasta track/ (rode dentro do projeto codekub ou defina CODEKUB_TRACK)")
}

func Load(root string) (*Track, error) {
	phaseFiles, err := filepath.Glob(filepath.Join(root, "*", "phase.json"))
	if err != nil || len(phaseFiles) == 0 {
		return nil, fmt.Errorf("nenhuma fase encontrada em %s", root)
	}
	sort.Strings(phaseFiles)
	t := &Track{Root: root}
	for _, pf := range phaseFiles {
		p := &Phase{Dir: filepath.Dir(pf)}
		if err := readJSON(pf, p); err != nil {
			return nil, err
		}
		metas, _ := filepath.Glob(filepath.Join(p.Dir, "*", "meta.json"))
		sort.Strings(metas)
		for _, mf := range metas {
			e := &Exercise{Dir: filepath.Dir(mf), Phase: p, TimeLimit: 1, Samples: 1}
			if err := readJSON(mf, e); err != nil {
				return nil, err
			}
			p.Exercises = append(p.Exercises, e)
		}
		t.Phases = append(t.Phases, p)
	}
	return t, nil
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

// Exercises lists every exercise in track order.
func (t *Track) Exercises() []*Exercise {
	var all []*Exercise
	for _, p := range t.Phases {
		all = append(all, p.Exercises...)
	}
	return all
}

// Lookup accepts "13", "013", "integer-average" or "015-integer-average".
func (t *Track) Lookup(ref string) (*Exercise, error) {
	ref = strings.TrimSpace(ref)
	n, numErr := strconv.Atoi(ref)
	for _, e := range t.Exercises() {
		id, _ := strconv.Atoi(e.ID)
		if (numErr == nil && id == n) || ref == e.Slug || ref == e.Key() {
			return e, nil
		}
	}
	return nil, fmt.Errorf("exercício %q não existe na trilha", ref)
}

func (e *Exercise) Key() string { return e.ID + "-" + e.Slug }

func (e *Exercise) Limit() time.Duration { return time.Duration(e.TimeLimit * float64(time.Second)) }

func (e *Exercise) TopicTitle() string {
	for _, tp := range e.Phase.Topics {
		if tp.Slug == e.Topic {
			return tp.Title
		}
	}
	return e.Topic
}

func (e *Exercise) Statement() (string, error) {
	data, err := os.ReadFile(filepath.Join(e.Dir, "statement.md"))
	return string(data), err
}

var digits = regexp.MustCompile(`\d+`)

func (e *Exercise) Tests() ([]judge.TestCase, error) {
	ins, _ := filepath.Glob(filepath.Join(e.Dir, "tests", "*.in"))
	if len(ins) == 0 {
		return nil, fmt.Errorf("%s não tem testes", e.Key())
	}
	sort.Slice(ins, func(i, j int) bool {
		a, _ := strconv.Atoi(digits.FindString(filepath.Base(ins[i])))
		b, _ := strconv.Atoi(digits.FindString(filepath.Base(ins[j])))
		return a < b
	})
	tests := make([]judge.TestCase, len(ins))
	for i, in := range ins {
		tests[i] = judge.TestCase{InputPath: in, OutputPath: strings.TrimSuffix(in, ".in") + ".out", Sample: i < e.Samples}
	}
	return tests, nil
}

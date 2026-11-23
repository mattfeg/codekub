// Package track loads the exercise track: phases (phase.json) and their
// exercises (meta.json, statement.md, tests/). The track usually comes
// embedded in the binary; CODEKUB_TRACK points at a folder on disk instead,
// which is how exercises are developed.
package track

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mattfeg/codekub/internal/judge"
)

// Embedded is the track compiled into the binary, set by package main.
var Embedded fs.FS

type Track struct {
	Source string // where the track came from, for display
	Phases []*Phase
}

type Phase struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Topics      []Topic `json:"topics"`
	Exercises   []*Exercise
	Dir         string // slash-separated path inside the track
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
	Dir       string   // slash-separated path inside the track
	Phase     *Phase
	fsys      fs.FS
}

type Source struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	License string `json:"license"`
	URL     string `json:"url"`
	Dataset string `json:"dataset"`
	Changes string `json:"changes"`
}

// Find picks the track: the folder in $CODEKUB_TRACK, or the embedded one.
func Find() (fs.FS, string, error) {
	if dir := os.Getenv("CODEKUB_TRACK"); dir != "" {
		if _, err := os.Stat(dir); err != nil {
			return nil, "", fmt.Errorf("CODEKUB_TRACK: %w", err)
		}
		return os.DirFS(dir), dir, nil
	}
	if Embedded != nil {
		return Embedded, "embutida no codekub", nil
	}
	return nil, "", errors.New("nenhuma trilha disponível (defina CODEKUB_TRACK)")
}

func Load(fsys fs.FS, source string) (*Track, error) {
	phaseFiles, err := fs.Glob(fsys, "*/phase.json")
	if err != nil || len(phaseFiles) == 0 {
		return nil, fmt.Errorf("nenhuma fase encontrada em %s", source)
	}
	sort.Strings(phaseFiles)
	t := &Track{Source: source}
	for _, pf := range phaseFiles {
		p := &Phase{Dir: path.Dir(pf)}
		if err := readJSON(fsys, pf, p); err != nil {
			return nil, err
		}
		metas, _ := fs.Glob(fsys, path.Join(p.Dir, "*", "meta.json"))
		sort.Strings(metas)
		for _, mf := range metas {
			e := &Exercise{Dir: path.Dir(mf), Phase: p, TimeLimit: 1, Samples: 1, fsys: fsys}
			if err := readJSON(fsys, mf, e); err != nil {
				return nil, err
			}
			p.Exercises = append(p.Exercises, e)
		}
		t.Phases = append(t.Phases, p)
	}
	return t, nil
}

func readJSON(fsys fs.FS, name string, v any) error {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%s: %w", name, err)
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
	data, err := fs.ReadFile(e.fsys, path.Join(e.Dir, "statement.md"))
	return string(data), err
}

var digits = regexp.MustCompile(`\d+`)

func (e *Exercise) Tests() ([]judge.TestCase, error) {
	ins, _ := fs.Glob(e.fsys, path.Join(e.Dir, "tests", "*.in"))
	if len(ins) == 0 {
		return nil, fmt.Errorf("%s não tem testes", e.Key())
	}
	sort.Slice(ins, func(i, j int) bool {
		a, _ := strconv.Atoi(digits.FindString(path.Base(ins[i])))
		b, _ := strconv.Atoi(digits.FindString(path.Base(ins[j])))
		return a < b
	})
	tests := make([]judge.TestCase, len(ins))
	for i, in := range ins {
		input, err := fs.ReadFile(e.fsys, in)
		if err != nil {
			return nil, err
		}
		output, err := fs.ReadFile(e.fsys, strings.TrimSuffix(in, ".in")+".out")
		if err != nil {
			return nil, err
		}
		tests[i] = judge.TestCase{Input: input, Output: output, Sample: i < e.Samples}
	}
	return tests, nil
}

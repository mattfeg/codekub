package track

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mattfeg/codekub/internal/judge"
)

const repoTrack = "../../track"

func loadRepoTrack(t *testing.T) *Track {
	t.Helper()
	tr, err := Load(os.DirFS(repoTrack), repoTrack)
	if err != nil {
		t.Fatal(err)
	}
	return tr
}

func TestLoadAndLookup(t *testing.T) {
	tr := loadRepoTrack(t)
	if len(tr.Exercises()) < 800 {
		t.Fatalf("expected the beecrowd problems, got %d exercises", len(tr.Exercises()))
	}
	for _, ref := range []string{"1001", "extremamente-basico", "1001-extremamente-basico"} {
		e, err := tr.Lookup(ref)
		if err != nil || e.Key() != "1001-extremamente-basico" {
			t.Errorf("Lookup(%q) = %v, %v", ref, e, err)
		}
	}
	if _, err := tr.Lookup("999"); err == nil {
		t.Error("unknown exercise should fail")
	}
	e, _ := tr.Lookup("1001")
	tests, err := e.Tests()
	if err != nil || len(tests) != 3 || !tests[2].Sample {
		t.Fatalf("tests: %d, err %v", len(tests), err)
	}
	if e.TopicTitle() != "Iniciante" || e.Source.URL != "https://resources.beecrowd.com/repository/UOJ_1001.html" {
		t.Errorf("topic %q, url %q", e.TopicTitle(), e.Source.URL)
	}
}

// Every exercise must be usable: unique id, a statement, tests, and either a
// reference solution the judge accepts or a page to read the problem on.
func TestEveryExerciseIsUsable(t *testing.T) {
	ruby, _ := judge.LanguageByID("ruby")
	seen := map[string]bool{}
	for _, e := range loadRepoTrack(t).Exercises() {
		if seen[e.ID] {
			t.Errorf("duplicate id %s", e.ID)
		}
		seen[e.ID] = true
		e := e
		t.Run(e.Key(), func(t *testing.T) {
			t.Parallel()
			if _, err := e.Statement(); err != nil {
				t.Fatal(err)
			}
			tests, err := e.Tests()
			if err != nil {
				t.Fatal(err)
			}
			reference := filepath.Join(repoTrack, filepath.FromSlash(e.Dir), "reference.rb")
			if _, err := os.Stat(reference); err != nil {
				if e.Source.URL == "" {
					t.Fatal("no reference solution and no source page")
				}
				return
			}
			res := judge.Run(context.Background(), judge.LocalRunner{}, judge.Request{
				Lang: ruby, Source: reference, Tests: tests, TimeLimit: e.Limit(),
			}, nil)
			if res.Verdict != judge.AC {
				t.Fatalf("%s: %s (%s) %+v", e.Key(), res.Verdict, res.Detail, res.Tests)
			}
		})
	}
}

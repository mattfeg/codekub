package track

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mattfeg/codekub/internal/judge"
)

func loadRepoTrack(t *testing.T) *Track {
	t.Helper()
	tr, err := Load(filepath.Join("..", "..", "track"))
	if err != nil {
		t.Fatal(err)
	}
	return tr
}

func TestLoadAndLookup(t *testing.T) {
	tr := loadRepoTrack(t)
	if len(tr.Exercises()) < 36 {
		t.Fatalf("expected the full basics phase, got %d exercises", len(tr.Exercises()))
	}
	for _, ref := range []string{"15", "015", "integer-average", "015-integer-average"} {
		e, err := tr.Lookup(ref)
		if err != nil || e.Key() != "015-integer-average" {
			t.Errorf("Lookup(%q) = %v, %v", ref, e, err)
		}
	}
	if _, err := tr.Lookup("999"); err == nil {
		t.Error("unknown exercise should fail")
	}
	e, _ := tr.Lookup("15")
	tests, err := e.Tests()
	if err != nil || len(tests) < 10 || !tests[0].Sample || !tests[1].Sample || tests[2].Sample {
		t.Fatalf("tests: %d, err %v", len(tests), err)
	}
	if e.TopicTitle() != "Operadores aritméticos" {
		t.Errorf("topic title %q", e.TopicTitle())
	}
}

// Every reference solution must be accepted by the real judge.
func TestReferenceSolutionsAreAccepted(t *testing.T) {
	ruby, _ := judge.LanguageByID("ruby")
	for _, e := range loadRepoTrack(t).Exercises() {
		e := e
		t.Run(e.Key(), func(t *testing.T) {
			t.Parallel()
			tests, err := e.Tests()
			if err != nil {
				t.Fatal(err)
			}
			res := judge.Run(context.Background(), judge.LocalRunner{}, judge.Request{
				Lang: ruby, Source: filepath.Join(e.Dir, "reference.rb"), Tests: tests, TimeLimit: e.Limit(),
			}, nil)
			if res.Verdict != judge.AC {
				t.Fatalf("%s: %s (%s) %+v", e.Key(), res.Verdict, res.Detail, res.Tests)
			}
		})
	}
}

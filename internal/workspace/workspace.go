// Package workspace manages the user's side of codekub: where solutions live,
// solution templates, the editor, and settings.
package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mattfeg/codekub/internal/judge"
	"github.com/mattfeg/codekub/internal/track"
)

type Workspace struct {
	Data      string // ~/.codekub: database and settings
	Solutions string // ~/codekub/solucoes: the files the user edits
}

type Config struct {
	Language string `json:"language"`
}

// Default honors CODEKUB_HOME (data) and CODEKUB_SOLUTIONS (solutions).
func Default() (*Workspace, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	w := &Workspace{
		Data:      envOr("CODEKUB_HOME", filepath.Join(home, ".codekub")),
		Solutions: envOr("CODEKUB_SOLUTIONS", filepath.Join(home, "codekub", "solucoes")),
	}
	for _, dir := range []string{w.Data, w.Solutions} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (w *Workspace) DBPath() string { return filepath.Join(w.Data, "codekub.db") }

func (w *Workspace) Config() Config {
	cfg := Config{Language: "ruby"}
	if data, err := os.ReadFile(filepath.Join(w.Data, "config.json")); err == nil {
		json.Unmarshal(data, &cfg)
	}
	if v := os.Getenv("CODEKUB_LANG"); v != "" {
		cfg.Language = v
	}
	return cfg
}

func (w *Workspace) Language() (judge.Language, error) {
	return judge.LanguageByID(w.Config().Language)
}

func (w *Workspace) SolutionPath(e *track.Exercise, lang judge.Language) string {
	return filepath.Join(w.Solutions, e.Key()+lang.Ext)
}

var hints = map[string]string{
	"ruby":   "Leia a entrada com gets e imprima a resposta com puts.",
	"python": "Leia a entrada com input() e imprima a resposta com print().",
}

// EnsureSolution creates the solution file from a template when missing.
func (w *Workspace) EnsureSolution(e *track.Exercise, lang judge.Language) (string, error) {
	path := w.SolutionPath(e, lang)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	c := lang.Comment
	body := fmt.Sprintf("%s %s · %s\n%s %s\n\n", c, e.ID, e.Title, c, hints[lang.ID])
	return path, os.WriteFile(path, []byte(body), 0o644)
}

// EditorCommand opens a file in $VISUAL or $EDITOR, falling back to the
// system's default text editor.
func EditorCommand(path string) *exec.Cmd {
	for _, key := range []string{"VISUAL", "EDITOR"} {
		if ed := os.Getenv(key); ed != "" {
			return exec.Command("sh", "-c", ed+` "$1"`, "editor", path)
		}
	}
	if runtime.GOOS == "darwin" {
		return exec.Command("open", "-t", path)
	}
	return exec.Command("xdg-open", path)
}

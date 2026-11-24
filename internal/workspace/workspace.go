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
	"strings"

	"github.com/mattfeg/codekub/internal/judge"
	"github.com/mattfeg/codekub/internal/track"
)

type Workspace struct {
	Data      string // ~/.codekub: database and settings
	Solutions string // ~/codekub/solucoes: the files the user edits
}

type Config struct {
	Language string `json:"language"`
	Editor   string `json:"editor,omitempty"` // shell command; the file path is appended
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

func (w *Workspace) configPath() string { return filepath.Join(w.Data, "config.json") }

// SavedConfig is what config.json holds, without environment overrides.
func (w *Workspace) SavedConfig() Config {
	cfg := Config{Language: "ruby"}
	if data, err := os.ReadFile(w.configPath()); err == nil {
		json.Unmarshal(data, &cfg)
	}
	return cfg
}

func (w *Workspace) SaveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(w.configPath(), append(data, '\n'), 0o644)
}

// Config is the effective configuration: config.json plus CODEKUB_LANG.
func (w *Workspace) Config() Config {
	cfg := w.SavedConfig()
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

// Editor returns the editor command and where it came from, in order of
// priority: CODEKUB_EDITOR, config.json, $VISUAL, $EDITOR. An empty command
// means the system's default text editor.
func (w *Workspace) Editor() (command, source string) {
	if v := os.Getenv("CODEKUB_EDITOR"); v != "" {
		return v, "CODEKUB_EDITOR"
	}
	if v := w.SavedConfig().Editor; v != "" {
		return v, "config.json"
	}
	for _, key := range []string{"VISUAL", "EDITOR"} {
		if v := os.Getenv(key); v != "" {
			return v, "$" + key
		}
	}
	return "", "padrão do sistema"
}

// EditorCommand builds the command that opens path in the chosen editor.
// On Unix the editor runs through the shell, so "code -w" or "nvim -p" work;
// Windows has no sh, so the command is split on spaces and run directly.
// Either way, path reaches the editor as a single argument.
func (w *Workspace) EditorCommand(path string) *exec.Cmd {
	ed, _ := w.Editor()
	switch {
	case ed != "" && runtime.GOOS == "windows":
		f := strings.Fields(ed)
		return exec.Command(f[0], append(f[1:], path)...)
	case ed != "":
		return exec.Command("sh", "-c", ed+` "$1"`, "editor", path)
	case runtime.GOOS == "windows":
		return exec.Command("notepad", path)
	case runtime.GOOS == "darwin":
		return exec.Command("open", "-t", path)
	}
	return exec.Command("xdg-open", path)
}

// BrowserCommand builds the command that opens url in the default browser.
func BrowserCommand(url string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		return exec.Command("open", url)
	}
	return exec.Command("xdg-open", url)
}

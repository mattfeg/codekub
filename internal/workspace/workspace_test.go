package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newWorkspace(t *testing.T) *Workspace {
	t.Helper()
	for _, key := range []string{"CODEKUB_EDITOR", "VISUAL", "EDITOR", "CODEKUB_LANG"} {
		t.Setenv(key, "")
	}
	t.Setenv("CODEKUB_HOME", t.TempDir())
	t.Setenv("CODEKUB_SOLUTIONS", t.TempDir())
	w, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func TestEditorPriority(t *testing.T) {
	w := newWorkspace(t)
	check := func(wantCmd, wantSource string) {
		t.Helper()
		if cmd, source := w.Editor(); cmd != wantCmd || source != wantSource {
			t.Fatalf("Editor() = %q from %q, want %q from %q", cmd, source, wantCmd, wantSource)
		}
	}
	check("", "padrão do sistema")
	t.Setenv("EDITOR", "vi")
	check("vi", "$EDITOR")
	t.Setenv("VISUAL", "nano")
	check("nano", "$VISUAL")
	if err := w.SaveConfig(Config{Language: "ruby", Editor: "code -w"}); err != nil {
		t.Fatal(err)
	}
	check("code -w", "config.json")
	t.Setenv("CODEKUB_EDITOR", "nvim")
	check("nvim", "CODEKUB_EDITOR")
}

func TestSaveConfigKeepsEnvOverridesOut(t *testing.T) {
	w := newWorkspace(t)
	t.Setenv("CODEKUB_LANG", "python")
	cfg := w.SavedConfig()
	cfg.Editor = "micro"
	if err := w.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(w.Data, "config.json"))
	if !strings.Contains(string(data), `"language": "ruby"`) || !strings.Contains(string(data), `"editor": "micro"`) {
		t.Fatalf("config.json should hold only saved values: %s", data)
	}
	if w.Config().Language != "python" {
		t.Fatal("CODEKUB_LANG should still override the effective language")
	}
}

func TestEditorCommandPassesPathSafely(t *testing.T) {
	w := newWorkspace(t)
	w.SaveConfig(Config{Language: "ruby", Editor: "printf '%s'"})
	path := filepath.Join(t.TempDir(), "nome com espaço; e $(perigo).rb")
	out, err := w.EditorCommand(path).Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != path {
		t.Fatalf("the path must reach the editor as one argument, got %q", out)
	}
}

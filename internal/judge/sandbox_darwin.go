package judge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const sandboxExec = "/usr/bin/sandbox-exec"

var (
	sandboxOnce sync.Once
	sandboxOK   bool
)

// SandboxAvailable reports whether solutions run under sandbox-exec: no
// network, and writes only inside their temp dir (plus /dev and the caches
// version managers such as mise need). CODEKUB_NO_SANDBOX=1 turns it off.
func SandboxAvailable() bool {
	sandboxOnce.Do(func() {
		if os.Getenv("CODEKUB_NO_SANDBOX") == "1" {
			return
		}
		sandboxOK = exec.Command(sandboxExec, "-p", "(version 1)(allow default)", "/usr/bin/true").Run() == nil
	})
	return sandboxOK
}

func SandboxName() string { return "sandbox-exec (macOS)" }

func sandboxProfile(dir string) string {
	home, _ := os.UserHomeDir()
	allow := fmt.Sprintf("(subpath %q) (subpath \"/dev\")", dir)
	for _, cache := range []string{"Library/Caches", ".cache", ".local/state", ".local/share/mise"} {
		allow += fmt.Sprintf(" (subpath %q)", filepath.Join(home, cache))
	}
	return "(version 1)(allow default)(deny network*)(deny file-write*)(allow file-write* " + allow + ")"
}

//go:build !windows

package judge

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// isolateProcess puts the solution in its own process group, so it can be
// killed together with any process it starts.
func isolateProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killTree(cmd *exec.Cmd) {
	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

// cleanupAfterExit kills whatever the solution left running in its group.
func cleanupAfterExit(cmd *exec.Cmd) { killTree(cmd) }

func exitStatus(ps *os.ProcessState) (code int, signal string) {
	if ws, ok := ps.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
		return 0, ws.Signal().String()
	}
	return ps.ExitCode(), ""
}

// childEnv passes only what interpreters and version managers need.
func childEnv(dir string) []string {
	env := []string{"TMPDIR=" + dir}
	for _, kv := range os.Environ() {
		key := strings.SplitN(kv, "=", 2)[0]
		switch {
		case key == "PATH", key == "HOME", key == "USER",
			strings.HasPrefix(key, "MISE_"), strings.HasPrefix(key, "RBENV_"),
			strings.HasPrefix(key, "ASDF_"), strings.HasPrefix(key, "PYENV_"):
			env = append(env, kv)
		}
	}
	locale := "C.UTF-8"
	if runtime.GOOS == "darwin" {
		locale = "en_US.UTF-8"
	}
	return append(env, "LANG="+locale, "LC_ALL="+locale)
}

//go:build windows

package judge

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// isolateProcess starts the solution in a new process group, detached from
// the console's Ctrl+C.
func isolateProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

// killTree ends the solution and every process it started.
func killTree(cmd *exec.Cmd) {
	if exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)).Run() != nil {
		cmd.Process.Kill()
	}
}

// cleanupAfterExit does nothing: after the process exits its PID may be
// reused by an unrelated program, so it is not safe to kill by PID.
func cleanupAfterExit(*exec.Cmd) {}

func exitStatus(ps *os.ProcessState) (code int, signal string) { return ps.ExitCode(), "" }

// childEnv keeps the environment (Windows programs need variables such as
// SystemRoot to start) and points the temp folders at the work dir.
func childEnv(dir string) []string {
	var env []string
	for _, kv := range os.Environ() {
		key := strings.ToUpper(strings.SplitN(kv, "=", 2)[0])
		if key != "TEMP" && key != "TMP" {
			env = append(env, kv)
		}
	}
	// Without UTF-8 mode, Python writes accented output in the console code
	// page (cp1252) when stdout is a pipe.
	return append(env, "TEMP="+dir, "TMP="+dir, "PYTHONUTF8=1", "PYTHONIOENCODING=utf-8")
}

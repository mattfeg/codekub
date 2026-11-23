package judge

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// OutputLimit caps what a solution may print per test, so an infinite print
// loop cannot exhaust memory.
const OutputLimit = 16 << 20

const checkTimeout = 30 * time.Second

// Runner prepares an isolated place to run one solution. LocalRunner uses the
// toolchains installed on this machine; a DockerRunner can implement the same
// interface later.
type Runner interface {
	Prepare(lang Language, source string) (Workspace, error)
}

type Workspace interface {
	// Check runs the language's syntax check or compilation.
	Check(ctx context.Context) (ok bool, output string, err error)
	Exec(ctx context.Context, stdin io.Reader, limit time.Duration) Execution
	Close() error
}

type Execution struct {
	Stdout, Stderr string
	ExitCode       int
	Signal         string
	TimedOut       bool
	OutputExceeded bool
	Canceled       bool
	Elapsed        time.Duration
	Err            error // the process could not be started at all
}

type LocalRunner struct{}

type localWorkspace struct {
	lang    Language
	dir     string
	src     string
	sandbox bool
}

func (LocalRunner) Prepare(lang Language, source string) (Workspace, error) {
	code, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	tmp, err := os.MkdirTemp("", "codekub-")
	if err != nil {
		return nil, err
	}
	// Sandbox rules match real paths, and macOS temp dirs live behind a symlink.
	dir, err := filepath.EvalSymlinks(tmp)
	if err != nil {
		os.RemoveAll(tmp)
		return nil, err
	}
	src := filepath.Join(dir, filepath.Base(source))
	if err := os.WriteFile(src, code, 0o644); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return &localWorkspace{lang: lang, dir: dir, src: src, sandbox: SandboxAvailable()}, nil
}

func (w *localWorkspace) Close() error { return os.RemoveAll(w.dir) }

func (w *localWorkspace) Check(ctx context.Context) (bool, string, error) {
	if len(w.lang.Check) == 0 {
		return true, "", nil
	}
	ex := w.exec(ctx, expand(w.lang.Check, w.src, w.dir), nil, checkTimeout)
	if ex.Err != nil {
		return false, "", ex.Err
	}
	if ex.TimedOut {
		return false, "a verificação de sintaxe passou do tempo", nil
	}
	out := strings.TrimSpace(ex.Stderr + "\n" + ex.Stdout)
	return ex.ExitCode == 0 && ex.Signal == "", w.relative(out), nil
}

func (w *localWorkspace) Exec(ctx context.Context, stdin io.Reader, limit time.Duration) Execution {
	ex := w.exec(ctx, expand(w.lang.Run, w.src, w.dir), stdin, limit)
	ex.Stderr = w.relative(ex.Stderr)
	return ex
}

// relative hides the temp directory from messages shown to the user.
func (w *localWorkspace) relative(s string) string {
	return strings.ReplaceAll(s, w.dir+string(os.PathSeparator), "")
}

func (w *localWorkspace) exec(ctx context.Context, argv []string, stdin io.Reader, limit time.Duration) Execution {
	if w.sandbox {
		argv = append([]string{sandboxExec, "-p", sandboxProfile(w.dir)}, argv...)
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = w.dir
	cmd.Env = childEnv(w.dir)
	cmd.Stdin = stdin
	isolateProcess(cmd)

	var killOnce sync.Once
	kill := func() {
		killOnce.Do(func() {
			if cmd.Process != nil {
				killTree(cmd)
			}
		})
	}
	stdout := &limitedBuffer{limit: OutputLimit, onExceed: kill}
	stderr := &limitedBuffer{limit: 64 << 10}
	cmd.Stdout, cmd.Stderr = stdout, stderr

	start := time.Now()
	if err := cmd.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			err = fmt.Errorf("%s não está instalado ou não está no PATH", argv[0])
		}
		return Execution{Err: err}
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	var ex Execution
	timer := time.NewTimer(limit)
	defer timer.Stop()
	select {
	case <-done:
		ex.Elapsed = time.Since(start)
	case <-timer.C:
		ex.TimedOut, ex.Elapsed = true, limit
		kill()
		<-done
	case <-ctx.Done():
		ex.Canceled, ex.Elapsed = true, time.Since(start)
		kill()
		<-done
	}
	cleanupAfterExit(cmd) // leftovers the solution may have spawned

	ex.Stdout, ex.Stderr = stdout.String(), stderr.String()
	ex.OutputExceeded = stdout.Exceeded()
	ex.ExitCode, ex.Signal = exitStatus(cmd.ProcessState)
	return ex
}

type limitedBuffer struct {
	mu       sync.Mutex
	buf      bytes.Buffer
	limit    int
	exceeded bool
	onExceed func()
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.buf.Len()+len(p) > b.limit {
		b.buf.Write(p[:max(0, b.limit-b.buf.Len())])
		if !b.exceeded {
			b.exceeded = true
			if b.onExceed != nil {
				go b.onExceed()
			}
		}
		return len(p), nil // keep draining so the process never blocks on a full pipe
	}
	return b.buf.Write(p)
}

func (b *limitedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func (b *limitedBuffer) Exceeded() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.exceeded
}

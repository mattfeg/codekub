//go:build !darwin

package judge

const sandboxExec = ""

// SandboxAvailable is false outside macOS for now; solutions still get a
// temp dir, time and output limits, and a minimal environment.
func SandboxAvailable() bool { return false }

func SandboxName() string { return "" }

func sandboxProfile(string) string { return "" }

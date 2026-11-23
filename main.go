// codekub: a local, progressive programming practice track with an online-judge
// style runner, in the terminal.
package main

import (
	"fmt"
	"os"

	"github.com/mattfeg/codekub/internal/core"
	"github.com/mattfeg/codekub/internal/track"
	"github.com/mattfeg/codekub/internal/ui"
	"github.com/mattfeg/codekub/internal/workspace"
)

// version is set at release time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	ui.Version = version
	track.Embedded = embeddedTrack()
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	cmd := ""
	if len(args) > 0 {
		cmd = args[0]
	}
	switch cmd {
	case "help", "-h", "--help":
		fmt.Println(ui.Help())
		return 0
	case "version", "--version", "-v":
		fmt.Println("codekub", version)
		return 0
	case "doctor":
		return ui.CmdDoctor()
	case "config": // works even when the saved settings are broken
		ws, err := workspace.Default()
		if err != nil {
			fmt.Fprintln(os.Stderr, "codekub:", err)
			return 2
		}
		return ui.CmdConfig(ws, args[1:])
	}

	env, err := core.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "codekub:", err)
		return 2
	}
	defer env.Close()

	switch cmd {
	case "":
		return tui(env, ui.Options{})
	case "next":
		return ui.CmdNext(env)
	case "list":
		return ui.CmdList(env)
	case "progress":
		return ui.CmdProgress(env)
	case "show":
		if len(args) < 2 {
			return usage("codekub show <n>")
		}
		return ui.CmdShow(env, args[1])
	case "test":
		ref, file, samples, ok := exerciseArgs(args[1:])
		if !ok {
			return usage("codekub test <n> [arquivo] [--exemplos]")
		}
		return ui.CmdTest(env, ref, file, samples)
	case "watch":
		ref, file, _, ok := exerciseArgs(args[1:])
		if !ok {
			return usage("codekub watch <n> [arquivo]")
		}
		ex, err := env.Track.Lookup(ref)
		if err != nil {
			fmt.Fprintln(os.Stderr, "codekub:", err)
			return 2
		}
		return tui(env, ui.Options{Exercise: ex, File: file, Watch: true})
	}
	fmt.Fprintf(os.Stderr, "comando desconhecido: %s  (veja codekub help)\n", cmd)
	return 2
}

// exerciseArgs parses "<n> [file] [--exemplos]" in any order.
func exerciseArgs(args []string) (ref, file string, samples, ok bool) {
	var rest []string
	for _, a := range args {
		if a == "--exemplos" || a == "-x" {
			samples = true
		} else {
			rest = append(rest, a)
		}
	}
	if len(rest) == 0 || len(rest) > 2 {
		return "", "", false, false
	}
	if len(rest) == 2 {
		file = rest[1]
	}
	return rest[0], file, samples, true
}

func usage(s string) int {
	fmt.Fprintln(os.Stderr, "uso:", s)
	return 2
}

func tui(env *core.Env, opts ui.Options) int {
	if err := ui.Run(env, opts); err != nil {
		fmt.Fprintln(os.Stderr, "codekub:", err)
		return 1
	}
	return 0
}

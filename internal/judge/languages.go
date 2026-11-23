package judge

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Language describes how to check and run solutions written in one language.
// Command parts may use the placeholders {src} (the solution file) and {dir}
// (the isolated working directory).
type Language struct {
	ID         string
	Name       string
	Ext        string
	Check      []string // syntax check or compilation; empty when not needed
	Run        []string
	TimeFactor float64 // multiplies the exercise's base time limit
	Version    []string
	Comment    string // line comment prefix, used by solution templates
}

// python is "python" on Windows, where "python3" often opens the Microsoft
// Store instead of an interpreter.
var python = pythonCommand()

func pythonCommand() string {
	if runtime.GOOS == "windows" {
		return "python"
	}
	return "python3"
}

var languages = []Language{
	{
		ID: "ruby", Name: "Ruby", Ext: ".rb",
		Check: []string{"ruby", "-c", "{src}"}, Run: []string{"ruby", "{src}"},
		TimeFactor: 3, Version: []string{"ruby", "-v"}, Comment: "#",
	},
	{
		ID: "python", Name: "Python", Ext: ".py",
		Check: []string{python, "-m", "py_compile", "{src}"}, Run: []string{python, "{src}"},
		TimeFactor: 3, Version: []string{python, "--version"}, Comment: "#",
	},
}

func Languages() []Language { return languages }

func LanguageByID(id string) (Language, error) {
	for _, l := range languages {
		if l.ID == id {
			return l, nil
		}
	}
	return Language{}, fmt.Errorf("linguagem não suportada: %s", id)
}

func LanguageForFile(path string) (Language, error) {
	ext := strings.ToLower(filepath.Ext(path))
	for _, l := range languages {
		if l.Ext == ext {
			return l, nil
		}
	}
	return Language{}, fmt.Errorf("não sei rodar arquivos %q", ext)
}

func expand(parts []string, src, dir string) []string {
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = strings.NewReplacer("{src}", src, "{dir}", dir).Replace(p)
	}
	return out
}

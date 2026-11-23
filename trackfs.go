package main

import (
	"embed"
	"io/fs"
)

// The track ships inside the binary: statements, metadata and tests only.
// Reference solutions and generators stay in the repository.
//
//go:embed track/*/phase.json track/*/*/meta.json track/*/*/statement.md track/*/*/tests
var trackFiles embed.FS

func embeddedTrack() fs.FS {
	sub, err := fs.Sub(trackFiles, "track")
	if err != nil {
		panic(err)
	}
	return sub
}

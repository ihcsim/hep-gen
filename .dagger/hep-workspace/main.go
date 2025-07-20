// A generated module for Workspace functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"dagger/hep-workspace/internal/dagger"
)

const (
	defaultImage   = "dannyben/madness:1.2.3"
	defaultPort    = 3000
	defaultWorkdir = "/docs"
	defaultPath    = "/docs"
)

type HepWorkspace struct {
	// the workspace's container state
	// +internal-use-only
	Container *dagger.Container
}

func New(
	source *dagger.Directory,
	tmplDownloadURL string,
	outputFile string,
) HepWorkspace {
	tmpl := dag.HTTP(tmplDownloadURL)
	return HepWorkspace{
		Container: dag.Container().
			From(defaultImage).
			WithMountedDirectory(defaultPath, source).
			WithFile(outputFile, tmpl).
			WithWorkdir(defaultWorkdir).
			WithExposedPort(defaultPort).
			WithDefaultTerminalCmd([]string{"/bin/bash"}),
	}
}

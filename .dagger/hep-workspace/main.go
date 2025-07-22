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
	"context"
	"dagger/hep-workspace/internal/dagger"
)

const defaultImage = "dannyben/madness:1.2.3"

// HepWorkspace represents a containerized workspace containing the tools and
// files needed to write a HEP.
type HepWorkspace struct {
	// the workspace's container state
	// +internal-use-only
	Container *dagger.Container
}

// New returns a new instance of the HEP workspace
func New(
	source *dagger.Directory,
	tmplDownloadURL string,
	outputFile string,
	workdir string,
	exposePort int,
) HepWorkspace {
	tmpl := dag.HTTP(tmplDownloadURL)
	return HepWorkspace{
		Container: dag.Container().
			From(defaultImage).
			WithMountedDirectory(workdir, source).
			WithWorkdir(workdir).
			WithFile(outputFile, tmpl).
			WithExposedPort(exposePort),
	}
}

// ReadFile reads file in the workspace
func (h *HepWorkspace) ReadFile(
	ctx context.Context,
	// the path to the file in the workspace
	path string,
) (string, error) {
	return h.Container.File(path).Contents(ctx)
}

// WriteFile writes a file to the workspace
func (h *HepWorkspace) WriteFile(
	// the path to the file in the workspace
	path string,
	// the new contents of the file
	contents string,
) *HepWorkspace {
	h.Container = h.Container.WithNewFile(path, contents)
	return h
}

// ListFiles list all the files in the workspace rooted at /docs
func (h *HepWorkspace) ListFiles(ctx context.Context) (string, error) {
	workdir, err := h.Container.Workdir(ctx)
	if err != nil {
		return "", err
	}
	return h.Container.
		WithExec([]string{"tree", workdir}).
		Stdout(ctx)
}

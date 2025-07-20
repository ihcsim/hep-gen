// A generated module for HepGen functions
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
	"dagger/hepgen/internal/dagger"
	"dagger/hepgen/pkg/prompt"
	"io"
	"path/filepath"
)

type HepGen struct{}

const (
	docsSite               = "https://docs.harvesterhci.io/v1.5/"
	tmplDownloadURL        = "https://raw.githubusercontent.com/harvester/harvester/refs/heads/master/enhancements/YYYYMMDD-template.md"
	fileHEP                = "index.md"
	fileProblemDescription = "problem.txt"
	workDir                = "work"
	exposePort             = 3000
)

var (
	filepathHEP                = filepath.Join(workDir, fileHEP)
	filepathProblemDescription = filepath.Join(workDir, fileProblemDescription)
)

// Hep generates a HEP draft with the given title. The generated content is output to stdout.
//
// The task is performed in a containerized sandbox workspace.
// The workspace has a bind mount to the host 'source' directory with the following files:
// * problem.txt - the HEP problem statement filled by the HEP author
// * template.md - the HEP template downloaded from tmplDownloadURL
// * index.md - the final draft of the HEP following the sections outlined in template.md
func (m *HepGen) Hep(
	ctx context.Context,
	// the KEP title
	title string,
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) (string, error) {
	w, err := m.workspace(title, source)
	if err != nil {
		return "", err
	}
	return w.Stdout(ctx)
}

// Preview publishes the generated HEP draft to localhost:3000.
// To port-forward to the container, use `dagger -c /bin/sh -c 'preview|up'`
// The markdown server is managed by 'madness' (https://madness.dannyb.co).
func (m *HepGen) Preview(
	ctx context.Context,
	// the KEP title
	title string,
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) (*dagger.Service, error) {
	w, err := m.workspace(title, source)
	if err != nil {
		return nil, err
	}

	serviceOpts := dagger.ContainerAsServiceOpts{
		Args: []string{"madness", "server"},
	}
	return w.AsService(serviceOpts), nil
}

func (m *HepGen) workspace(
	// the KEP title
	title string,
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) (*dagger.Container, error) {
	promptInputs := &prompt.PromptInputs{
		Title:                      title,
		DocsSite:                   docsSite,
		FilepathHEP:                filepathHEP,
		FilepathProblemDescription: filepathProblemDescription,
	}

	out, err := prompt.ExecTmpl(promptInputs)
	if err != nil {
		return nil, err
	}
	prompt, err := io.ReadAll(out)
	if err != nil {
		return nil, err
	}

	ws := dag.HepWorkspace(
		source,
		tmplDownloadURL,
		filepathHEP,
		exposePort,
	)
	env := dag.Env().
		WithHepWorkspaceInput("workspace", ws, "the workspace for this task").
		WithHepWorkspaceOutput("workspace", "the workspace with the generated HEP draft")
	return dag.LLM().
		WithEnv(env).
		WithPrompt(string(prompt)).
		Env().
		Output("workspace").
		AsHepWorkspace().
		Container(), nil
}

// Sandbox returns a sandbox container representing the workspace with a bind mount to the host 'source' directory.
// The sandbox container is exposed at port 3000.
// To port-forward to the container, use `dagger -c /bin/sh -c 'sandbox|up'`
func (m *HepGen) Sandbox(
	ctx context.Context,
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) *dagger.Service {
	serviceOpts := dagger.ContainerAsServiceOpts{
		Args: []string{"madness", "server"},
	}
	return dag.HepWorkspace(
		source,
		tmplDownloadURL,
		fileHEP,
		exposePort,
	).Container().
		AsService(serviceOpts)
}

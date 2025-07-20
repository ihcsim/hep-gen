// A generated module for HepWriter functions
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
	"dagger/hep-writer/internal/dagger"
	"dagger/hep-writer/pkg/prompt"
	"io"
	"path/filepath"
)

type HepWriter struct{}

const (
	docsSite               = "https://docs.harvesterhci.io/v1.5/"
	tmplDownloadURL        = "https://raw.githubusercontent.com/harvester/harvester/refs/heads/master/enhancements/YYYYMMDD-template.md"
	fileHEP                = "index.md"
	fileProblemDescription = "problem.txt"
	workDir                = "work"
)

var (
	filepathHEP                = filepath.Join(workDir, fileHEP)
	filepathProblemDescription = filepath.Join(workDir, fileProblemDescription)
)

// Hep generates a HEP draft with the given title. A sandbox workspace is created with a bind mount to the
// host 'source' directory. The workspace contains the following files:
// * problem.txt - the HEP problem statement
// * index.md - the final draft of the HEP
func (m *HepWriter) Hep(
	ctx context.Context,
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
		filepathHEP)
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

// Workspace returns a sandbox container representing the workspace with a bind mount to the host 'source' directory.
// The sandbox container is exposed at port 3000.
func (m *HepWriter) Workspace(
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) *dagger.Container {
	return dag.HepWorkspace(
		source,
		tmplDownloadURL,
		fileHEP,
	).Container()
}

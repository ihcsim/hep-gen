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
	"fmt"
	"path/filepath"
)

type HepWriter struct{}

const (
	docsSite               = "https://docs.harvesterhci.io/v1.5/"
	workDir                = "work"
	problemDescriptionFile = "problem.txt"
	templateFile           = "template.md"
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
) *dagger.Container {
	prompt := fmt.Sprintf(`Harvester is a modern, open, interoperable, hyperconverged infrastructure (HCI) solution built on Kubernetes. It is an open source project maintained by SUSE.
	Harvester depends on KubeVirt to provide the API to run virtual machines in Kubernetes. Longhorn serves as its main storage provider using the Container Storage Interface (CSI) API.
	More information on Harvester can be found at %s.

	You are a software engineer tasked with the responsbility to write a new Harvester Enhancement Proposal (aka HEP).

	A HEP is a proposal documentation that describes new technical enhancement to address specific user problems. It is expressed in the markdown language.

	The title of your HEP is going to be "%s". The problem description can be found in the %s file in your workspace.

	The enhancement can either propose new features or changes to improve existing features.
	The readers of this HEP are software engineers familiar with technologies like Kubernetes, Golang, Python, YAML etc.
	Please use appropriate formatting such as headers, tables, bullet points and embedded .PNG images to improve readability.
	Including Golang and YAML code samples to illustrate certain changes is encouraged, but not necessary.
	Make sure that the generated HEP contains only valid markdown syntax.

	You can find the template of the HEP at %s. Make a copy of the template file to store the draft content. Name the file after the HEP title, using kebab style.

	You have access to a workspace. The workspace contains the problem description and template files.
`, docsSite, title, filepath.Join(workDir, problemDescriptionFile), filepath.Join(workDir, templateFile))

	ws := dag.HepWorkspace(source)
	env := dag.Env().
		WithHepWorkspaceInput("workspace", ws, "the workspace for this task").
		WithHepWorkspaceOutput("workspace", "the workspace with the generated HEP draft")
	return dag.LLM().
		WithEnv(env).
		WithPrompt(prompt).
		Env().
		Output("workspace").
		AsHepWorkspace().
		Container()
}

// Workspace returns a sandbox container representing the workspace with a bind mount to the host 'source' directory.
// The sandbox container is exposed at port 3000.
func (m *HepWriter) Workspace(
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) *dagger.Container {
	return dag.HepWorkspace(source).Container()
}

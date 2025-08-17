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
	"dagger/hep-gen/internal/dagger"
	prompthep "dagger/hep-gen/pkg/prompt/hep"
	promptreview "dagger/hep-gen/pkg/prompt/review"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type HepGen struct{}

const (
	bindWorkdir      = "/docs"
	srcWorkdir       = "./work"
	docSiteHarvester = "https://docs.harvesterhci.io/v1.5/"
	docSiteKubeVirt  = "https://kubevirt.io/user-guide/"
	docSiteLonghorn  = "https://longhorn.io/docs/1.9.0/"
	docSiteMadness   = "https://madness.dannyb.co/"
	exposePort       = 3000
	fileHEP          = "index.md"
	fileSummary      = "summary.md"
	fileReview       = "review.md"
	tmplDownloadURL  = "https://raw.githubusercontent.com/harvester/harvester/refs/heads/master/enhancements/YYYYMMDD-template.md"
)

var caser = cases.Title(language.English)

// Hep generates a HEP draft with the given title.
// The task is performed in a containerized sandbox workspace.
// The workspace has a bind mount to the host 'source' directory with the following files:
// * problem.txt - the HEP problem statement filled by the HEP author
// * template.md - the HEP template downloaded from tmplDownloadURL
// * index.md - the final draft of the HEP following the sections outlined in template.md
func (m *HepGen) Hep(
	ctx context.Context,
	// the HEP title
	title string,
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) (*dagger.Container, error) {
	promptInputs := &prompthep.PromptInputs{
		HEPTitle:         title,
		DocSiteHarvester: docSiteHarvester,
		DocSiteKubeVirt:  docSiteKubeVirt,
		DocSiteLonghorn:  docSiteLonghorn,
		DocSiteMadness:   docSiteMadness,
		FilepathHEP:      fileHEP,
		FilepathSummary:  fileSummary,
	}

	out, err := prompthep.ExecTmpl(promptInputs)
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
		fileHEP,
		bindWorkdir,
		exposePort,
	)
	env := dag.Env().
		WithHepWorkspaceInput("workspace", ws, "the workspace for this task").
		WithHepWorkspaceOutput("completed", "the workspace with the generated HEP draft")
	return dag.LLM().
		WithEnv(env).
		WithPrompt(string(prompt)).
		Env().
		Output("completed").
		AsHepWorkspace().
		Container(), nil
}

// Prleview publishes the generated HEP draft to localhost:3000.
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
	w, err := m.Hep(ctx, title, source)
	if err != nil {
		return nil, err
	}

	serviceOpts := dagger.ContainerAsServiceOpts{
		Args: []string{"madness", "server"},
	}
	return w.AsService(serviceOpts), nil
}

// Review reviews the content of the HEP found in the local `./work/index.md` file.
func (m *HepGen) Review(
	ctx context.Context,
	// the local directory where the HEP index.md is
	// +defaultPath="./work"
	source *dagger.Directory) (
	*dagger.Directory, error,
) {
	promptInputs := &promptreview.PromptInputs{
		DocSiteHarvester: docSiteHarvester,
		DocSiteKubeVirt:  docSiteKubeVirt,
		DocSiteLonghorn:  docSiteLonghorn,
		FilepathHEP:      fileHEP,
		FilepathReview:   fileReview,
	}
	out, err := promptreview.ExecTmpl(promptInputs)
	if err != nil {
		return nil, err
	}
	prompt, err := io.ReadAll(out)
	if err != nil {
		return nil, err
	}

	env := dag.Env().
		WithDirectoryInput("review", source, "the local directory where the HEP index.md file is").
		WithDirectoryOutput("completed", "the local directory where the review feedback file is").
		WithFileOutput(fileReview, "the file containing the review feedback")

	return dag.LLM().
		WithEnv(env).
		WithPrompt(string(prompt)).
		Env().
		Output("completed").
		AsDirectory(), nil
}

// PullRequest submits a branch and a pull request containing the changes
// in the HEP work/index.md file to the upstream repository. PullRequest only
// works against the 'main' branch.
// Additional files to be included in the changes can be specified using the
// 'includes' argument.
func (m *HepGen) PullRequest(
	ctx context.Context,
	// the upstream repository to submit the pull request to
	upstream string,
	// the GitHub issue number associated with the HEP. hep-gen creates a branch named after this issue number
	issue string,
	// github token to use to submit the pull request
	token *dagger.Secret,
	// name of the fork if the fork option is set to true
	// +default=""
	forkName string,
	// if true, create a fork of the upstream repository
	// +default=false
	fork bool,
	// the source directory with the changes to commit
	// +defaultPath="./work"
	changes *dagger.Directory,
	// additional files to be included in the pull request in addition to the index.md file
	// +default=[]string
	includes []string,
) (string, error) {
	branchOpts := dagger.FeatureBranchOpts{
		ForkName: forkName,
		Fork:     fork,
	}
	timestamp := time.Now().UnixMilli()
	branchName := fmt.Sprintf("gh%s-hep-%d", issue, timestamp)
	branch := dag.FeatureBranch(token, upstream, branchName, branchOpts)

	renamed := fmt.Sprintf("gh%s-hep.md", issue)
	changes = changes.WithFile(".", changes.File(fileHEP).WithName(renamed))
	commitDirectoryOpts := dagger.DirectoryWithDirectoryOpts{
		Include: append(includes, renamed),
	}

	changesOpts := dagger.FeatureBranchWithChangesOpts{
		KeepGit: true,
	}
	staging := branch.Branch().WithDirectory("enhancements", changes, commitDirectoryOpts)
	branch = branch.WithChanges(staging, changesOpts)

	outStatus, err := branch.Env().WithExec([]string{"git", "status"}).Stdout(ctx)
	if err != nil {
		return "", err
	}

	title, err := m.hepTitle(ctx, changes.File(renamed))
	if err != nil {
		return "", err
	}
	title = "[HEP] " + title

	body := fmt.Sprintf("This PR includes a HEP for issue %s", issue)
	outPR, err := branch.PullRequest(ctx, title, body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s", outStatus, outPR), nil
}

// Sandbox returns a sandbox container representing the workspace with a bind mount to the host 'source' directory.
// The sandbox container is exposed at port 3000.
// To start an interactive session, use `dagger -c /bin/sh -c 'sandbox|terminal'`
// To port-forward to the container, use `dagger -c /bin/sh -c 'sandbox|as-service|up'`
func (m *HepGen) Sandbox(
	ctx context.Context,
	// the source directory to mount into the workspace
	// +defaultPath="./work"
	source *dagger.Directory,
) *dagger.Container {
	args := []string{"madness", "server"}
	return dag.HepWorkspace(
		source,
		tmplDownloadURL,
		fileHEP,
		bindWorkdir,
		exposePort,
	).
		Container().
		WithDefaultArgs(args)
}

func (m *HepGen) hepTitle(ctx context.Context, f *dagger.File) (string, error) {
	out, err := f.Contents(ctx)
	if err != nil {
		return "", err
	}

	lines := strings.Split(out, "\n")
	if !strings.HasPrefix(lines[0], "# ") {
		return "", fmt.Errorf("failed to find the HEP title. The first line of %s must begin with a h1 markdown title", fileHEP)
	}

	title := strings.TrimSpace(strings.TrimLeft(lines[0], "#"))
	return caser.String(title), nil
}

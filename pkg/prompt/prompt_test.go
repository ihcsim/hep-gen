package prompt

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPromptTmpl(t *testing.T) {
	var (
		title            = "Increase test coverage via AI-assisted fuzzy testing"
		docSiteHarvester = "https://docs.harvesterhci.io/v1.5/"
		docSiteKubeVirt  = "https://kubevirt.io/user-guide/"
		docSiteLonghorn  = "https://longhorn.io/docs/1.9.0/"
		docSiteMadness   = "https://madness.dannyb.co/"
		filepathHEP      = filepath.Join("work", "index.md")
		filepathSummary  = filepath.Join("work", "summary.md")

		expected = fmt.Sprintf(`Harvester is a modern, open, interoperable, hyperconverged infrastructure (HCI) solution built on Kubernetes. It is an open source project maintained by SUSE.
Harvester depends on KubeVirt (%s) to run virtual machines in Kubernetes. Longhorn (%s) serves as its main storage provider using the Container Storage Interface (CSI) API. More information on Harvester can be found at %s.

You are a software engineer tasked with the responsbility to write a new Harvester Enhancement Proposal (HEP). A HEP is a proposal documentation that describes new technical enhancement to address specific user problems. It is expressed in the markdown language.

The title of your HEP is going to be "%s".

The enhancement can propose either new features or improvements to existing features. The readers of this HEP are software engineers familiar with technologies like Kubernetes, Golang, Python, YAML etc.

Please use appropriate formatting such as headers, tables, bullet points and embedded .PNG images to improve readability. Including Golang and YAML code samples to illustrate certain changes is encouraged, but not necessary. Make sure that the generated HEP contains only valid markdown syntax.

You have access to a workspace to do the work. The workspace is uses 'madness' (%s), an instant markdown server, to render the markdown documents.

The workspace contains a summary file located at '%s' and an index file located at '%s'. Both files are expressed in the markdown language.

The index file contains two sections namely, 'Problem' and 'Solution'. The 'Problem' section describes the problem the HEP is attempting to solve. The 'Solution' section provides a preliminary description of the solution to address the problem. Use this solution as the starting point for the HEP.

The index file is when the HEP should be written to. Maintain its original sections layout. Fill in all the sections. Sections with the '[optional]' label in their title are optional.`, docSiteKubeVirt, docSiteLonghorn, docSiteHarvester, title, docSiteMadness, filepathSummary, filepathHEP)

		inputs = &PromptInputs{
			HEPTitle:         title,
			DocSiteHarvester: docSiteHarvester,
			DocSiteKubeVirt:  docSiteKubeVirt,
			DocSiteLonghorn:  docSiteLonghorn,
			DocSiteMadness:   docSiteMadness,
			FilepathHEP:      filepathHEP,
			FilepathSummary:  filepathSummary,
		}
	)
	out, err := ExecTmpl(inputs)
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	actual, err := io.ReadAll(out)
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	if !reflect.DeepEqual(expected, string(actual)) {
		t.Errorf("mismatch output:\n-> expected: '%s'\n\n-> actual: '%s'", expected, actual)
	}
}

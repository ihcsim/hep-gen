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
		title                      = "Increase test coverage via AI-assisted fuzzy testing"
		docsSite                   = "https://docs.harvesterhci.io/v1.5/"
		filepathHEP                = filepath.Join("work", "index.md")
		filepathProblemDescription = filepath.Join("work", "problem.txt")

		expected = fmt.Sprintf(`Harvester is a modern, open, interoperable, hyperconverged infrastructure (HCI) solution built on Kubernetes. It is an open source project maintained by SUSE.
Harvester depends on KubeVirt to provide the API to run virtual machines in Kubernetes. Longhorn serves as its main storage provider using the Container Storage Interface (CSI) API.
More information on Harvester can be found at %s.

You are a software engineer tasked with the responsbility to write a new Harvester Enhancement Proposal (aka HEP). The title of your HEP is going to be "%s".

A HEP is a proposal documentation that describes new technical enhancement to address specific user problems. It is expressed in the markdown language.

You have access to a workspace to do the work. The workspace contains a problem description file named '%s' and an '%s' file.
Write the HEP in the %s file. Maintain its original sections layout. Fill in all the sections. Sections with the '[optional]' label in their title are optional.

The enhancement can either propose new features or changes to improve existing features.
The readers of this HEP are software engineers familiar with technologies like Kubernetes, Golang, Python, YAML etc.
Please use appropriate formatting such as headers, tables, bullet points and embedded .PNG images to improve readability.
Including Golang and YAML code samples to illustrate certain changes is encouraged, but not necessary.
Make sure that the generated HEP contains only valid markdown syntax.`, docsSite, title, filepathProblemDescription, filepathHEP, filepathHEP)

		inputs = &PromptInputs{
			Title:                      title,
			DocsSite:                   docsSite,
			FilepathHEP:                filepathHEP,
			FilepathProblemDescription: filepathProblemDescription,
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

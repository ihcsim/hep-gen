package prompt

import (
	"bytes"
	"html/template"
	"io"
)

type PromptInputs struct {
	HEPTitle         string
	DocSiteHarvester string
	DocSiteKubeVirt  string
	DocSiteLonghorn  string
	DocSiteMadness   string
	FilepathHEP      string
	FilepathSummary  string
	WordLimit        int
}

const prompt = `Harvester is a modern, open, interoperable, hyperconverged infrastructure (HCI) solution built on Kubernetes. It is an open source project maintained by SUSE.
Harvester depends on KubeVirt ({{ .DocSiteKubeVirt }}) to run virtual machines in Kubernetes. Longhorn ({{ .DocSiteLonghorn }}) serves as its main storage provider using the Container Storage Interface (CSI) API. More information on Harvester can be found at {{ .DocSiteHarvester }}.

You are a software engineer tasked with the responsbility to write a new Harvester Enhancement Proposal (HEP). A HEP is a proposal documentation that describes new technical enhancement to address specific user problems. It is expressed in the markdown language.

The title of your HEP is going to be "{{ .HEPTitle }}".

The enhancement can propose either new features or improvements to existing features. The readers of this HEP are software engineers familiar with technologies like Kubernetes, Golang, Python, YAML etc.

Please use appropriate formatting such as headers, tables, bullet points and embedded .PNG images to improve readability. Including Golang and YAML code samples to illustrate certain changes is encouraged, but not necessary. Make sure that the generated HEP contains only valid markdown syntax.

You have access to a workspace to do the work. The workspace is uses 'madness' ({{ .DocSiteMadness }}), an instant markdown server, to render the markdown documents.

The workspace contains a summary file located at '{{ .FilepathSummary }}' and an index file located at '{{ .FilepathHEP }}'. Both files are expressed in the markdown language.

Please read the summary file. It contains two sections namely, 'Problem' and 'Solution'. The 'Problem' section describes the problem the HEP is attempting to solve. The 'Solution' section provides a preliminary description of the solution to address the problem. Use this solution as the starting point for the HEP.

Please write the HEP to the index file is. Fill in all the sections. Sections with the '[optional]' label in their title are optional.`

func ExecTmpl(inputs *PromptInputs) (io.Reader, error) {
	tmpl, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if err := tmpl.Execute(out, inputs); err != nil {
		return nil, err
	}

	return out, nil
}

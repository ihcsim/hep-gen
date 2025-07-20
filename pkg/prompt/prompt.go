package prompt

import (
	"bytes"
	"html/template"
	"io"
)

type PromptInputs struct {
	Title                      string
	DocsSite                   string
	FilepathHEP                string
	FilepathProblemDescription string
}

const prompt = `Harvester is a modern, open, interoperable, hyperconverged infrastructure (HCI) solution built on Kubernetes. It is an open source project maintained by SUSE.
Harvester depends on KubeVirt to provide the API to run virtual machines in Kubernetes. Longhorn serves as its main storage provider using the Container Storage Interface (CSI) API.
More information on Harvester can be found at {{ .DocsSite }}.

You are a software engineer tasked with the responsbility to write a new Harvester Enhancement Proposal (aka HEP). The title of your HEP is going to be "{{ .Title }}".

A HEP is a proposal documentation that describes new technical enhancement to address specific user problems. It is expressed in the markdown language.

You have access to a workspace to do the work. The workspace contains a problem description file named '{{ .FilepathProblemDescription }}' and an '{{ .FilepathHEP }}' file.
Write the HEP in the {{ .FilepathHEP }} file. Maintain its original sections layout. Fill in all the sections. Sections with the '[optional]' label in their title are optional.

The enhancement can either propose new features or changes to improve existing features.
The readers of this HEP are software engineers familiar with technologies like Kubernetes, Golang, Python, YAML etc.
Please use appropriate formatting such as headers, tables, bullet points and embedded .PNG images to improve readability.
Including Golang and YAML code samples to illustrate certain changes is encouraged, but not necessary.
Make sure that the generated HEP contains only valid markdown syntax.`

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

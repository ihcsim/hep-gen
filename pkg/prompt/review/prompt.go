// Package review contains the LLM prompt for the Review function
package review

import (
	"bytes"
	"html/template"
	"io"
)

type PromptInputs struct {
	DocSiteHarvester string
	DocSiteKubeVirt  string
	DocSiteLonghorn  string
	FilepathHEP      string
	FilepathReview   string
}

const prompt = `Harvester is a modern, open, interoperable, hyperconverged infrastructure (HCI) solution built on Kubernetes. It is an open source project maintained by SUSE.

Harvester depends on KubeVirt ({{ .DocSiteKubeVirt }}) to run virtual machines in Kubernetes. Longhorn ({{ .DocSiteLonghorn }}) serves as its main storage provider using the Container Storage Interface (CSI) API. More information on Harvester can be found at {{ .DocSiteHarvester }}.

You are a software engineer tasked with the responsbility to review a new Harvester Enhancement Proposal (HEP). A HEP is a proposal documentation that describes new technical enhancement to address specific user problems. It is expressed in the markdown language.

Please provide feedback on the HEP by affirming good design decisions, pointing out best practices for cloud-native design, raising concerns and asking clarifying questions about the proposed solution. Suggestions to improve the writing style are also welcome. The enhancement can propose either new features or improvements to existing features. The readers of this HEP are software engineers familiar with technologies like Kubernetes, Golang, Python, YAML etc.

You will be given a directory to work with. Write all your review feedback to {{ .FilepathReview }}. If the {{ .FilepathReview }} doesn't exist, create it.

The HEP to be reviewed can be found at {{ .FilepathHEP }}. If this file doesn't exist, do not proceed. Inform the user that you can't find the file.

The {{ .FilepathHEP }} must have the following main sections:

* Title
* Summary
* Motivation
* Proposal

The 'Title' section provides a clear and descriptive title of the HEP.

The 'Summary' section defines the problem that the HEP is intended to solve. It provides a brief description of the intended solution, without getting caught up in the weed. Readers should be able to grasp the general idea of the HEP by reading this section. 

The 'Motivation' section defines the goals of the HEP. It describes a list of measurable acceptance criteria. Optionally, non-goals are included to limit the scope of the HEP.

The 'Proposal' section provides more insights into how the HEP affects and benefits the users. It outlines user stories and user experience in more details. Changes to Harvester's external API that may create backward incompatibility are called out.

The 'Design' section is the heart of the HEP. This is where the solution is discussed in details, including trade-offs, pitfalls, security and performance consideration, test plans and code snippets. Other general consideration for upgrading cloud-native distributed systems like Harvester can also be listed here.

Please ensure the HEP uses appropriate formatting such as headers, tables, bullet points and embedded .PNG images to improve readability. Make sure that the generated HEP contains only valid markdown syntax.`

func ExecTmpl(inputs *PromptInputs) (io.Reader, error) {
	tmpl, err := template.New("prompt-review").Parse(prompt)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if err := tmpl.Execute(out, inputs); err != nil {
		return nil, err
	}

	return out, nil
}

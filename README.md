# HEP Generator

The `hep-gen` Dagger module is an AI agent that can be used to generate a Harvester
Enhancement Proposal (HEP) draft. The purpose of the tool is to give Harvester
engineers a head start with HEP authoring by:

* Providing initial content and word suggestions for various sections like user
stories and test cases
* Suggesting research topics and additional notes for unfamiliar technologies
* Ensuring the HEP template layout and formatting are used

⚠️ The HEP draft must be reviewed and approved by Harvester maintainers, prior to
acceptance into the `harvester/harvester` repository.

## Prerequisites

* [Dagger](https://dagger.io/) v0.18.2

## Description

The `hep-gen` module implements 3 functions:

```sh
$ dagger functions
Name      Description
hep       Hep generates a HEP draft with the given title. The generated content is output to stdout.
preview   Preview publishes the generated HEP draft to localhost:3000.
review    Review reviews the content of the HEP found in the local `./work/index.md` file.
sandbox   Sandbox returns a sandbox container representing the workspace with a bind mount to the host 'source' directory.
```

The agent performs its task in a containerized sandbox, with a bind mount to the
`./work` folder in this project.

The `hep` function is the main function for writing a HEP draft. The title of the
HEP must be provided as an argument to the function. Use the `./work/summary.md`
file to provide additional context about the HEP to the LLM. In particular, the
LLM expects a problem definition and a brief description of the desired solution.
The generated HEP is written to a file named `index.md` in the workspace. It can
be exported to the `./work` folder using Dagger's built-in `export` function. See
example below.

The `preview` function is built on top of the `hep` function where it publishes
the HEP to localhost:3000, for better readability.

For debugging purposes, the `sandbox` function can be used to start an
interactive session with an empty workspace with a bind mount to the `./work`
folder.

Use the Dagger `.help` command to learn more about each function. E.g.:

```sh
dagger -c '.help hep'
```

## Usage

Follow the Dagger
[LLM endpoint configuration instructions](https://docs.dagger.io/configuration/llm/)
to configure an LLM for use with Dagger.

Think of a good HEP title to be used as the argument to the `hep` function.

Update the `./work/summary.md` file with information on the problem that HEP is
attempting to address. Provide a brief description of the proposed solution.

Launch the Dagger shell:

```sh
dagger
```

Assign the task to create the HEP to a `task` variable:

```sh
task=$(hep "<hep-title>")
```

Generate the HEP and publish it to localhost:3000:

```sh
$task | as-service --args "madness" --args "server" | up
```

To exit the long-running service process from with Dagger shell, switch to the
`navigate` mode and quit the process. See <https://github.com/dagger/dagger/issues/10069>.

To edit the generated `index.md` HEP file, start an interactive session in the
workspace:

```sh
$task | terminal
```

To export the generated `index.md` HEP file to the local `./work` directory:

```sh
$task | directory . | export ./work
```

To get the LLM to review the exported `index.md` HEP file and write its
feedback to the `work/review.md`:

```sh
review | export ./work
```

Once the `index.md` file is exported, a new interactive workspace can be started
using the `sandbox` function with the `./work` folder remounted:

```sh
sandbox | terminal
```

To submit a pull request containing the HEP to the Harvester repository:

```sh
pull-request https://github.com/<fork-org>/harvester.git <github-issue-number> <gh_token>
```

## Development

To build and test the module and its dependencies:

```sh
make
```

The source of module's dependencies can be found in the `.dagger` folder.

To build the `hep-workspace` submodule only:

```sh
make hep-workspace
```

## LICENSE

Apache License Version 2.0.

See [LICENSE](LICENSE) file.

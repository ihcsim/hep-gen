# HEP Generator

The `hep-gen` Dagger module is an AI agent that can be used to generate a Harvester
Enhancement Proposal (HEP) draft. The purpose of the tool is to give Harvester
engineers a head start with HEP authoring. It ensures that the HEP template
layout and formatting are enforced.

⚠️ The HEP draft must be reviewed and approved by Harvester maintainers, prior to
acceptance into the `harvester/harvester` repository.

## Prerequisites

* [Dagger](https://dagger.io/) v0.18.2 - Follow the Dagger [LLM
endpoint configuration instructions](https://docs.dagger.io/configuration/llm/)
to configure an LLM for use with Dagger.

## Description

The `hep-gen` module implements 3 functions:

```sh
$ dagger functions
Name      Description
hep       Hep generates a HEP draft with the given title. The generated content is output to stdout.
preview   Preview publishes the generated HEP draft to localhost:3000.
sandbox   Sandbox returns a sandbox container representing the workspace with a bind mount to the host 'source' directory.
```

The agent performs its task in a containerized sandbox, with a bind mount to the
`./work` folder in this project

The `hep` function is the main function for writing a HEP draft. The title of the
HEP must be provided as an argument to the function. The generated HEP is written
to a file named `index.md`. The HEP author is expected to provide a good problem
statement that the HEP addresses in the `problem.txt` file.

The `preview` function built on top of the `hep` function where it publishes
the HEP to localhost:3000, for better readability.

For debugging purposes, the `workspace` function can be used to start and
interactive session with an empty workspace with a bind mount to the `./work`
folder.

Use the Dagger `.help` command to learn more about each function. E.g.:

```sh
dagger -c '.help hep'
```

## Usages

To launch an interactive sandbox:

```sh
dagger -c /bin/sh -c 'sandbox|terminal'
```

To preview the markdown in the sandbox at <http://localhost:3000>:

```sh
dagger -c /bin/sh -c 'sandbox|as-service|up'
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

See [LICENSE](#license) file.

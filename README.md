# dupl

**dupl** is a tool for finding code clones written in Go. So far it can find clones only
in the Go source files. The method uses suffix tree for serialized ASTs. It ignores values
of AST nodes, it cares just about the types.

## Installation

```bash
go get -u github.com/mibk/dupl
```

## Usage

It searches in a current directory by default (can be changed using the first argument).
Using `-files` flag it is possible to specify coveted files which are read from *stdin*.
The `-t` flag is used to set a minimal clone size measured in syntax units (default is 15).
The default output is just text information about duplicate line ranges in files. Using `-html`
flag the output could turn into HTML page with duplicate code fragments.

## Example

The reduced output of this command with the following parameters for the [Docker](https://www.docker.com) source code
looks like [this](http://htmlpreview.github.io/?https://github.com/mibk/dupl/blob/master/_output_example/docker.html).

```bash
$ dupl -t 200 -html >docker.html
```

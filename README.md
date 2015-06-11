# dupl

**dupl** is a tool for finding code clones written in Go. So far it can find clones only
in the Go source files. The method uses suffix tree for serialized ASTs. It ignores values
of AST nodes, it cares just about the types.

## Installation

```bash
go get -u github.com/mibk/dupl
```

## Usage

`dupl` searches in a current directory by default (can be changed using the first argument).
The output is plaintext about duplicate line ranges in files.

### Flags
- `-files`: read input from *stdin* instead of the current directory.
- `-t size`: set the clone size threshold, measured in syntax units. Default is equivalent to `-t 15`.
- `-html`: output the results as HTML, including duplicate code fragments.
- `-plumbing`: output in an easy-to-parse format for scripts or tools.

## Example

The reduced output of this command with the following parameters for the [Docker](https://www.docker.com) source code
looks like [this](http://htmlpreview.github.io/?https://github.com/mibk/dupl/blob/master/_output_example/docker.html).

```bash
$ dupl -t 200 -html >docker.html
```

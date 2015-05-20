# dupl

**dupl** is a tool for finding code clones written in Go. So far it can find clones only
in the Go source files. The method uses suffix tree for serialized ASTs.

## Usage

It searches in a current directory by default (can be changed using the first argument).
Using `-files` flag it is possible to specify coveted files which are read from *stdin*.
The `-t` flag is used to set a minimal clone size measured in syntax units (default is 15).

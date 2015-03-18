package text

import (
	"fmt"
	"io"
	"io/ioutil"

	"fm.tul.cz/dupl/syntax"
)

type Printer struct {
	writer io.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{
		writer: w,
	}
}

func (p *Printer) Print(dups [][]*syntax.Node) {

	fmt.Fprintf(p.writer, "found %d clones:\n", len(dups))
	for i, dup := range dups {
		if len(dup) == 0 {
			panic("zero length dup")
		}
		nstart := dup[0]
		nend := dup[len(dup)-1]

		// TODO: Duplication could possibly be over several files.
		file, err := ioutil.ReadFile(nstart.Filename)
		if err != nil {
			panic(err)
		}

		lstart, lend := blockLines(file, nstart.Pos, nend.End)
		fmt.Fprintf(p.writer, "  loc %d: %s, line %d-%d,\n", i+1, nstart.Filename, lstart, lend)
	}
}

func blockLines(file []byte, from, to int) (int, int) {
	line := 1
	lineStart, lineEnd := 0, 0
	for offset, b := range file {
		if b == '\n' {
			line++
		}
		if offset == from {
			lineStart = line
		}
		if offset == to-1 {
			lineEnd = line
			break
		}
	}
	return lineStart, lineEnd
}

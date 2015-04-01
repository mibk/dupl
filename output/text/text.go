package text

import (
	"fmt"
	"io"

	"fm.tul.cz/dupl/syntax"
)

type FileReader interface {
	ReadFile(node *syntax.Node) ([]byte, error)
}

type Printer struct {
	writer  io.Writer
	freader FileReader
}

func NewPrinter(w io.Writer, fr FileReader) *Printer {
	return &Printer{
		writer:  w,
		freader: fr,
	}
}

func (p *Printer) Print(dups []*syntax.Seq) {
	fmt.Fprintf(p.writer, "found %d clones:\n", len(dups))
	for i, dup := range dups {
		cnt := len(dup.Nodes)
		if cnt == 0 {
			panic("zero length dup")
		}
		nstart := dup.Nodes[0]
		nend := dup.Nodes[cnt-1]

		file, err := p.freader.ReadFile(nstart)
		if err != nil {
			panic(err)
		}

		lstart, lend := blockLines(file, nstart.Pos, nend.End)
		filename := nstart.Filename
		if nstart.Addr != "" {
			filename = nstart.Addr + "@" + filename
		}
		fmt.Fprintf(p.writer, "  loc %d: %s, line %d-%d,\n", i+1, filename, lstart, lend)
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

package output

import (
	"fmt"
	"io"

	"fm.tul.cz/dupl/syntax"
)

type FileReader interface {
	ReadFile(node *syntax.Node) ([]byte, error)
}

type Printer interface {
	Print(dups [][]*syntax.Node)
	Finish()
}

type TextPrinter struct {
	writer  io.Writer
	freader FileReader
	cnt     int
}

func NewTextPrinter(w io.Writer, fr FileReader) *TextPrinter {
	return &TextPrinter{
		writer:  w,
		freader: fr,
	}
}

func (p *TextPrinter) Print(dups [][]*syntax.Node) {
	p.cnt++
	fmt.Fprintf(p.writer, "found %d clones:\n", len(dups))
	for i, dup := range dups {
		cnt := len(dup)
		if cnt == 0 {
			panic("zero length dup")
		}
		nstart := dup[0]
		nend := dup[cnt-1]

		file, err := p.freader.ReadFile(nstart)
		if err != nil {
			panic(err)
		}

		lstart, lend := blockLines(file, nstart.Pos, nend.End)
		fmt.Fprintf(p.writer, "  loc %d: %s, line %d-%d,\n", i+1, nstart.Filename, lstart, lend)
	}
}

func (p *TextPrinter) Finish() {
	fmt.Fprintf(p.writer, "\nFound total %d clone groups.\n", p.cnt)
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

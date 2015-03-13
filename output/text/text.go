package text

import (
	"fmt"
	"io"
	"io/ioutil"

	"fm.tul.cz/dupl/suffixtree"
	"fm.tul.cz/dupl/syntax"
)

type Printer struct {
	writer io.Writer
	stree  *suffixtree.STree
}

func NewPrinter(w io.Writer, t *suffixtree.STree) *Printer {
	return &Printer{
		writer: w,
		stree:  t,
	}
}

func (p *Printer) Print(m suffixtree.Match) {
	fmt.Fprintf(p.writer, "found match of length %d:\n", m.Len)

	// TODO: Match may not form a whole syntax unit. It may even be comprised
	// of more than one file.
	nstart1 := getNode(p.stree.At(m.P1))
	nend1 := getNode(p.stree.At(m.P1 + m.Len - 1))
	nstart2 := getNode(p.stree.At(m.P2))
	nend2 := getNode(p.stree.At(m.P2 + m.Len - 1))

	file1, err := ioutil.ReadFile(nstart1.Filename)
	if err != nil {
		panic(err)
	}
	file2 := file1
	if nstart1.Filename == nstart2.Filename {
		file2, err = ioutil.ReadFile(nstart2.Filename)
		if err != nil {
			panic(err)
		}
	}

	lstart1, lend1 := blockLines(file1, nstart1.Pos, nend1.End)
	lstart2, lend2 := blockLines(file2, nstart2.Pos, nend2.End)

	fmt.Fprintf(p.writer, "  loc 1: %s, line %d-%d,\n  loc 2: %s, line %d-%d.\n",
		nstart1.Filename, lstart1, lend1, nstart2.Filename, lstart2, lend2)
}

func getNode(tok suffixtree.Token) *syntax.Node {
	if n, ok := tok.(*syntax.Node); ok {
		return n
	}
	panic(fmt.Sprintf("tok (type %T)  is not type *syntax.Node", tok))
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
		if offset == to {
			lineEnd = line
		}
	}
	return lineStart, lineEnd
}

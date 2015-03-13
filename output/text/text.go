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

func (p *Printer) Print(dup1, dup2 []*syntax.Node) {
	if len(dup1) == 0 || len(dup2) == 0 {
		return
	}

	nstart1 := dup1[0]
	nend1 := dup1[len(dup1)-1]
	nstart2 := dup2[0]
	nend2 := dup2[len(dup2)-1]

	// TODO: Duplication could possibly be over several files.
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

	fmt.Fprintf(p.writer, "found clone spanning %d lines:\n", lend1-lstart1+1)
	fmt.Fprintf(p.writer, "  loc 1: %s, line %d-%d,\n  loc 2: %s, line %d-%d.\n",
		nstart1.Filename, lstart1, lend1, nstart2.Filename, lstart2, lend2)
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

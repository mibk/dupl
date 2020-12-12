package clones

import (
	"io/ioutil"

	"github.com/mibk/dupl/internal/syntax"
)

func CreateClones(dups [][]*syntax.Node) ([]Clone, error) {
	clones := make([]Clone, len(dups))
	for i, dup := range dups {
		cnt := len(dup)
		if cnt == 0 {
			panic("zero length dup")
		}
		nstart := dup[0]
		nend := dup[cnt-1]

		file, err := ioutil.ReadFile(nstart.Filename)
		if err != nil {
			return nil, err
		}

		cl := Clone{Filename: nstart.Filename}
		cl.LineStart, cl.LineEnd = blockLines(file, nstart.Pos, nend.End)
		clones[i] = cl
	}
	return clones, nil
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

type ByNameAndLine []Clone

func (c ByNameAndLine) Len() int { return len(c) }

func (c ByNameAndLine) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c ByNameAndLine) Less(i, j int) bool {
	if c[i].Filename == c[j].Filename {
		return c[i].LineStart < c[j].LineStart
	}
	return c[i].Filename < c[j].Filename
}

type Clone struct {
	Filename  string
	LineStart int
	LineEnd   int
	Fragment  []byte
}

type Issue struct {
	From, To Clone
}

package printer

import (
	"fmt"
	"io"
	"sort"

	"github.com/mibk/dupl/syntax"
)

type plumbing struct {
	w       io.Writer
	freader FileReader
}

func NewPlumbing(w io.Writer, fr FileReader) Printer {
	return &plumbing{w, fr}
}

func (p *plumbing) PrintHeader() error { return nil }

func (p *plumbing) PrintClones(dups [][]*syntax.Node) error {
	clones, err := prepareClonesInfo(p.freader, dups)
	if err != nil {
		return err
	}
	sort.Sort(byNameAndLine(clones))
	for i, cl := range clones {
		nextCl := clones[(i+1)%len(clones)]
		fmt.Fprintf(p.w, "%s:%d-%d: duplicate of %s:%d-%d\n", cl.filename, cl.lineStart, cl.lineEnd,
			nextCl.filename, nextCl.lineStart, nextCl.lineEnd)
	}
	return nil
}

func (p *plumbing) PrintFooter() error { return nil }

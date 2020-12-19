package printer

import (
	"fmt"
	"io"
	"sort"

	"github.com/mibk/dupl/internal/clones"
	"github.com/mibk/dupl/internal/syntax"
)

type plumbing struct {
	w io.Writer
}

func NewPlumbing(w io.Writer) Printer {
	return &plumbing{w}
}

func (p *plumbing) PrintHeader() error { return nil }

func (p *plumbing) PrintClones(dups [][]*syntax.Node) error {
	duplicates, err := clones.CreateClones(dups)
	if err != nil {
		return err
	}
	sort.Sort(clones.ByNameAndLine(duplicates))
	for i, cl := range duplicates {
		nextCl := duplicates[(i+1)%len(duplicates)]
		fmt.Fprintf(p.w, "%s:%d-%d: duplicate of %s:%d-%d\n", cl.Filename, cl.LineStart, cl.LineEnd,
			nextCl.Filename, nextCl.LineStart, nextCl.LineEnd)
	}
	return nil
}

func (p *plumbing) PrintFooter() error { return nil }



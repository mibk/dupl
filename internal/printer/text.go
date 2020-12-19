package printer

import (
	"fmt"
	"io"
	"sort"

	"github.com/mibk/dupl/internal/clones"
	"github.com/mibk/dupl/internal/syntax"
)

type text struct {
	cnt int
	w   io.Writer
}

func NewText(w io.Writer) Printer {
	return &text{w: w}
}

func (p *text) PrintHeader() error { return nil }

func (p *text) PrintClones(dups [][]*syntax.Node) error {
	p.cnt++
	fmt.Fprintf(p.w, "found %d clones:\n", len(dups))
	duplicates, err := clones.CreateClones(dups)
	if err != nil {
		return err
	}
	sort.Sort(clones.ByNameAndLine(duplicates))
	for _, cl := range duplicates {
		fmt.Fprintf(p.w, "  %s:%d,%d\n", cl.Filename, cl.LineStart, cl.LineEnd)
	}
	return nil
}

func (p *text) PrintFooter() error {
	_, err := fmt.Fprintf(p.w, "\nFound total %d clone groups.\n", p.cnt)
	return err
}

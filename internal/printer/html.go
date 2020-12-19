package printer

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"regexp"
	"sort"

	"github.com/mibk/dupl/internal/clones"
	"github.com/mibk/dupl/internal/syntax"
)

type htmlprinter struct {
	iota int
	w    io.Writer
}

func NewHTML(w io.Writer) Printer {
	return &htmlprinter{w: w}
}

func (p *htmlprinter) PrintHeader() error {
	_, err := fmt.Fprint(p.w, `<!DOCTYPE html>
<meta charset="utf-8"/>
<title>Duplicates</title>
<style>
	pre {
		background-color: #FFD;
		border: 1px solid #E2E2E2;
		padding: 1ex;
	}
</style>
`)
	return err
}

func (p *htmlprinter) PrintClones(dups [][]*syntax.Node) error {
	p.iota++
	fmt.Fprintf(p.w, "<h1>#%d found %d clones</h1>\n", p.iota, len(dups))

	duplicates, err := clones.CreateClones(dups)
	if err != nil {
		return err
	}

	sort.Sort(clones.ByNameAndLine(duplicates))
	for _, cl := range duplicates {
		fmt.Fprintf(p.w, "<h2>%s:%d</h2>\n<pre>%s</pre>\n", cl.Filename, cl.LineStart,
			html.EscapeString(string(cl.Fragment)))
	}
	return nil
}

func (*htmlprinter) PrintFooter() error { return nil }

func toWhitespace(str []byte) []byte {
	var out []byte
	for _, c := range bytes.Runes(str) {
		if c == '\t' {
			out = append(out, '\t')
		} else {
			out = append(out, ' ')
		}
	}
	return out
}

func deindent(block []byte) []byte {
	const maxVal = 99
	min := maxVal
	re := regexp.MustCompile(`(^|\n)(\t*)\S`)
	for _, line := range re.FindAllSubmatch(block, -1) {
		indent := line[2]
		if len(indent) < min {
			min = len(indent)
		}
	}
	if min == 0 || min == maxVal {
		return block
	}
	block = block[min:]
Loop:
	for i := 0; i < len(block); i++ {
		if block[i] == '\n' && i != len(block)-1 {
			for j := 0; j < min; j++ {
				if block[i+j+1] != '\t' {
					continue Loop
				}
			}
			block = append(block[:i+1], block[i+1+min:]...)
		}
	}
	return block
}

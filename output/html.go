package output

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"fm.tul.cz/dupl/syntax"
)

type HtmlPrinter struct {
	iota int
	*TextPrinter
}

func NewHtmlPrinter(w io.Writer, fr FileReader) *HtmlPrinter {
	fmt.Fprint(w, `<!DOCTYPE html>
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
	return &HtmlPrinter{
		TextPrinter: NewTextPrinter(w, fr),
	}
}

func (p *HtmlPrinter) Print(dups []*syntax.Seq) {
	p.iota++
	fmt.Fprintf(p.writer, "<h1>#%d found %d clones</h1>\n", p.iota, len(dups))
	for _, dup := range dups {
		cnt := len(dup.Nodes)
		if cnt == 0 {
			panic("zero length dup")
		}
		nstart := dup.Nodes[0]
		nend := dup.Nodes[cnt-1]

		filename := nstart.Filename
		if nstart.Addr != "" {
			filename = nstart.Addr + "@" + filename
		}
		file, err := p.freader.ReadFile(nstart)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(p.writer, "<h2>%s</h2>\n", filename)
		start := findLineBeg(file, nstart.Pos)
		content := append(toWhitespace(file[start:nstart.Pos]), file[nstart.Pos:nend.End]...)
		fmt.Fprintf(p.writer, "<pre>%s</pre>\n", deindent(content))
	}
}

func (p *HtmlPrinter) Finish() {}

func findLineBeg(file []byte, index int) int {
	for i := index; i >= 0; i-- {
		if file[i] == '\n' {
			return i + 1
		}
	}
	return 0
}

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
	min := 99
	re := regexp.MustCompile(`(^|\n)(\t*)\S`)
	for _, line := range re.FindAllSubmatch(block, -1) {
		indent := line[2]
		if len(indent) < min {
			min = len(indent)
		}
	}
	if min == 0 {
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"fm.tul.cz/dupl/job"
	"fm.tul.cz/dupl/output"
	"fm.tul.cz/dupl/remote"
	"fm.tul.cz/dupl/suffixtree"
	"fm.tul.cz/dupl/syntax"
)

var (
	dir        = "."
	threshold  = flag.Int("t", 15, "minimum token sequence as a clone")
	serverPort = flag.String("serve", "", "run server at port")
	addrs      AddrList
	html       = flag.Bool("html", false, "html output")
)

type AddrList []string

func (l *AddrList) String() string {
	return fmt.Sprintf("%v", *l)
}

func (l *AddrList) Set(val string) error {
	*l = append(*l, val)
	return nil
}

func init() {
	flag.Var(&addrs, "c", "connect to the given 'addr:port'")
}

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	if len(addrs) != 0 {
		t, clients := remote.RunClient(addrs)
		printDupls(t, remote.NewFileReader(clients))
	} else if *serverPort != "" {
		remote.RunServer(*serverPort, dir)
	} else {
		schan := job.CrawlDir(dir)
		bchan := make(chan *job.Batch)
		go func() {
			for seq := range schan {
				bchan <- job.NewBatch("", seq)
			}
			close(bchan)
		}()
		t, done := job.BuildTree(bchan)
		<-done
		printDupls(t, new(LocalFileReader))
	}
}

type char int

func (c char) Val() int {
	return int(c)
}

type LocalFileReader struct{}

func (r *LocalFileReader) ReadFile(node *syntax.Node) ([]byte, error) {
	return ioutil.ReadFile(node.Filename)
}

func printDupls(t *suffixtree.STree, fr output.FileReader) {
	// finish stream
	t.Update(char(-1))

	// print clones
	var p output.Printer
	if *html {
		p = output.NewHtmlPrinter(os.Stdout, fr)
	} else {
		p = output.NewTextPrinter(os.Stdout, fr)
	}
	mchan := t.FindDuplOver(*threshold)
	for m := range mchan {
		if dups := syntax.FindSyntaxUnits(t, m, *threshold); len(dups) != 0 {
			p.Print(dups)
		}
	}
	p.Finish()
}

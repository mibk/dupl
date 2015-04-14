package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"fm.tul.cz/dupl/job"
	"fm.tul.cz/dupl/output"
	"fm.tul.cz/dupl/remote"
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
		nodesChan := remote.RunClient(addrs, *threshold, dir)
		printDupls(nodesChan)
	} else if *serverPort != "" {
		remote.RunServer(*serverPort)
	} else {
		schan := job.CrawlDir(dir)
		t, done := job.BuildTree(schan)
		<-done

		// finish stream
		t.Update(&syntax.Node{Type: -1})

		mchan := t.FindDuplOver(*threshold)
		nodesChan := make(chan [][]*syntax.Node)
		go func() {
			for m := range mchan {
				nodesChan <- syntax.GetNodes(t, m)
			}
			close(nodesChan)
		}()
		printDupls(nodesChan)
	}
}

type LocalFileReader struct{}

func (r *LocalFileReader) ReadFile(node *syntax.Node) ([]byte, error) {
	return ioutil.ReadFile(node.Filename)
}

func printDupls(nodesChan <-chan [][]*syntax.Node) {
	groups := make(map[string][]*syntax.Seq)
	for seqs := range nodesChan {
		if dups, hash := syntax.FindSyntaxUnits(seqs, *threshold); len(dups) != 0 {
			if _, ok := groups[hash]; ok {
				groups[hash] = JoinGroups(groups[hash], dups)
			} else {
				groups[hash] = dups
			}
		}
	}

	p := getPrinter()
	for _, group := range groups {
		p.Print(group)
	}
	p.Finish()
}

func getPrinter() output.Printer {
	fr := new(LocalFileReader)
	if *html {
		return output.NewHtmlPrinter(os.Stdout, fr)
	}
	return output.NewTextPrinter(os.Stdout, fr)
}

func JoinGroups(grp1, grp2 []*syntax.Seq) []*syntax.Seq {
	// TODO: rm redundant fragments
	return append(grp1, grp2...)
}

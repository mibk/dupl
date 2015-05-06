package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"fm.tul.cz/dupl/job"
	"fm.tul.cz/dupl/output"
	"fm.tul.cz/dupl/remote"
	"fm.tul.cz/dupl/syntax"
)

const DefaultThreshold = 15

var (
	dir        = "."
	verbose    = flag.Bool("verbose", false, "explain what is being done")
	threshold  = flag.Int("threshold", DefaultThreshold, "minimum token sequence as a clone")
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
	flag.BoolVar(verbose, "v", false, "alias for -verbose")
	flag.IntVar(threshold, "t", DefaultThreshold, "alias for -threshold")
	flag.Var(&addrs, "connect", "connect to the given 'addr:port'")
	flag.Var(&addrs, "c", "alias for -connect")
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
		if *verbose {
			log.Println("Building suffix tree")
		}
		schan := job.CrawlDir(dir)
		t, data, done := job.BuildTree(schan)
		<-done

		// finish stream
		t.Update(&syntax.Node{Type: -1})

		if *verbose {
			log.Println("Searching for clones")
		}
		mchan := t.FindDuplOver(*threshold)
		duplChan := make(chan syntax.Match)
		go func() {
			for m := range mchan {
				match := syntax.FindSyntaxUnits(*data, m, *threshold)
				if len(match.Frags) > 0 {
					duplChan <- match
				}
			}
			close(duplChan)
		}()
		printDupls(duplChan)
	}
}

type LocalFileReader struct{}

func (r *LocalFileReader) ReadFile(node *syntax.Node) ([]byte, error) {
	return ioutil.ReadFile(node.Filename)
}

func printDupls(duplChan <-chan syntax.Match) {
	groups := make(map[string][][]*syntax.Node)
	for dupl := range duplChan {
		hash := dupl.Hash
		if _, ok := groups[hash]; ok {
			groups[hash] = append(groups[hash], dupl.Frags...)
		} else {
			groups[hash] = dupl.Frags
		}
	}

	p := getPrinter()
	for _, group := range groups {
		uniq := Unique(group)
		if len(uniq) != 1 {
			p.Print(uniq)
		}
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

func Unique(group [][]*syntax.Node) [][]*syntax.Node {
	fileMap := make(map[string]map[int]bool)

	newGroup := make([][]*syntax.Node, 0)
	for _, seq := range group {
		node := seq[0]
		file, ok := fileMap[node.Filename]
		if !ok {
			file = make(map[int]bool)
			fileMap[node.Filename] = file
		}
		if _, ok = file[node.Pos]; !ok {
			file[node.Pos] = true
			newGroup = append(newGroup, seq)
		}
	}
	return newGroup
}

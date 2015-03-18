package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fm.tul.cz/dupl/output/text"
	"fm.tul.cz/dupl/suffixtree"
	"fm.tul.cz/dupl/syntax"
	"fm.tul.cz/dupl/syntax/golang"
)

type char int

func (c char) Val() int {
	return int(c)
}

var (
	dir       = "."
	threshold = flag.Int("t", 15, "minimum token sequence as a clone")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		dir = args[0]
	}

	// collecting files
	fchan := make(chan string)
	go func() {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				fchan <- path
			}
			return nil
		})
		close(fchan)
	}()

	// AST parsing
	achan := make(chan *syntax.Node)
	go func() {
		for {
			file, ok := <-fchan
			if !ok {
				break
			}
			ast, err := golang.Parse(file)
			if err != nil {
				log.Println(err)
				continue
			}
			achan <- ast
		}
		close(achan)
	}()

	// serialization
	schan := make(chan []*syntax.Node)
	go func() {
		for {
			ast, ok := <-achan
			if !ok {
				break
			}
			seq := syntax.Serialize(ast)
			schan <- seq
		}
		close(schan)
	}()

	// suffix tree
	t := suffixtree.New()
	for {
		seq, ok := <-schan
		if !ok {
			break
		}
		for _, item := range seq {
			t.Update(item)
		}
	}

	// finish stream
	t.Update(char(-1))

	// printing the clones
	printer := text.NewPrinter(os.Stdout)
	mchan := t.FindDuplOver(*threshold)
	cnt := 0
	for {
		m, ok := <-mchan
		if !ok {
			break
		}
		if dups := syntax.FindSyntaxUnits(t, m, *threshold); len(dups) != 0 {
			printer.Print(dups)
			cnt++
		}
	}
	fmt.Printf("\nFound total %d clone groups.\n", cnt)
}

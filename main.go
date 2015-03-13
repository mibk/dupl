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

	t := suffixtree.New()

	for {
		file, ok := <-fchan
		if !ok {
			break
		}
		syn, err := golang.Parse(file)
		if err != nil {
			log.Println(err)
			continue
		}
		stream := syntax.Serialize(syn)
		for _, item := range stream {
			t.Update(item)
		}
	}

	// finish stream
	t.Update(char(-1))

	printer := text.NewPrinter(os.Stdout)

	mchan := t.FindDuplOver(*threshold)
	cnt := 0
	for {
		m, ok := <-mchan
		if !ok {
			break
		}
		dup1, dup2 := syntax.FindSyntaxUnits(t, m)
		printer.Print(dup1, dup2)
		cnt++
	}
	fmt.Printf("\nFound total %d clones.\n", cnt)
}

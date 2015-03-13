package main

import (
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

func main() {
	fchan := make(chan string)
	go func() {
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				fchan <- path
			}
			return nil
		})
		close(fchan)
	}()

	t := suffixtree.New()
	printer := text.NewPrinter(os.Stdout, t)

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

	mchan := t.FindDuplOver(25)
	cnt := 0
	for {
		m, ok := <-mchan
		if !ok {
			break
		}
		printer.Print(m)
		cnt++
	}
	fmt.Printf("\nFound total %d clones.\n", cnt)
}

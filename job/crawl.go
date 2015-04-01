package job

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"fm.tul.cz/dupl/syntax"
	"fm.tul.cz/dupl/syntax/golang"
)

func CrawlDir(dir string) chan []*syntax.Node {

	// collect files
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

	// parse AST
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

	// serialize
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
	return schan
}

package printer

import "github.com/mibk/dupl/syntax"

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

type Printer interface {
	Print(dups [][]*syntax.Node) error
	Finish() error
}

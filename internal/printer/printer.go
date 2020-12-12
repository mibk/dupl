package printer

import "github.com/mibk/dupl/internal/syntax"

type Printer interface {
	PrintHeader() error
	PrintClones(dups [][]*syntax.Node) error
	PrintFooter() error
}

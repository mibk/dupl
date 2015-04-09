package job

import (
	"fm.tul.cz/dupl/suffixtree"
	"fm.tul.cz/dupl/syntax"
)

func BuildTree(schan chan []*syntax.Node) (t *suffixtree.STree, done chan bool) {
	t = suffixtree.New()
	done = make(chan bool)
	go func() {
		for seq := range schan {
			for _, node := range seq {
				t.Update(node)
			}
		}
		done <- true
	}()
	return t, done
}

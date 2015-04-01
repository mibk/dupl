package job

import (
	"fm.tul.cz/dupl/suffixtree"
	"fm.tul.cz/dupl/syntax"
)

type Batch struct {
	addr string
	seq  []*syntax.Node
}

func NewBatch(addr string, seq []*syntax.Node) *Batch {
	return &Batch{addr, seq}
}

func BuildTree(bchan chan *Batch) (t *suffixtree.STree, done chan bool) {
	t = suffixtree.New()
	done = make(chan bool)
	go func() {
		for {
			batch, ok := <-bchan
			if !ok {
				break
			}
			for _, item := range batch.seq {
				item.Addr = batch.addr
				t.Update(item)
			}
		}
		done <- true
	}()
	return t, done
}

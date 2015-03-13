package suffixtree

type Match struct {
	P1, P2 pos
	Len    int
}

type clist struct {
	lists map[int]*plist
}

func newClist() *clist {
	return &clist{make(map[int]*plist)}
}

type plist struct {
	positions []pos
}

func newPlist() *plist {
	return &plist{make([]pos, 0)}
}

func (p *plist) append(p2 *plist) {
	p.positions = append(p.positions, p2.positions...)
}

func (p *plist) add(pos pos) {
	p.positions = append(p.positions, pos)
}

func (c *clist) combine(c2 *clist, length, threshold int, ch chan<- Match) {
	if length < threshold {
		return
	}
	for lc1, pl1 := range c.lists {
		for lc2, pl2 := range c2.lists {
			if lc1 != lc2 {
				for _, p1 := range pl1.positions {
					for _, p2 := range pl2.positions {
						ch <- Match{p1, p2, length}
					}
				}
			}
		}
	}
	c.append(c2)
}

func (c *clist) append(c2 *clist) {
	for lc, pl := range c2.lists {
		if _, ok := c.lists[lc]; ok {
			c.lists[lc].append(pl)
		} else {
			c.lists[lc] = pl
		}
	}
}

// FindDuplOver find pairs of maximal duplicities over a threshold
// length.
func (t *STree) FindDuplOver(threshold int) <-chan Match {
	auxTran := newTran(0, 0, t.root)
	ch := make(chan Match)
	go func() {
		walkTrans(auxTran, 0, threshold, ch)
		close(ch)
	}()
	return ch
}

func walkTrans(parent *tran, length, threshold int, ch chan<- Match) *clist {
	s := parent.state

	cl := newClist()

	if len(s.trans) == 0 {
		pl := newPlist()
		start := parent.end + 1 - pos(length)
		pl.add(start)
		ch := 0
		if start > 0 {
			ch = s.t.data[start-1].Val()
		}
		cl.lists[ch] = pl
		return cl
	}

	for _, t := range s.trans {
		cl.combine(walkTrans(t, length+t.len(), threshold, ch), length, threshold, ch)
	}
	return cl
}

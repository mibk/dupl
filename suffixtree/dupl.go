package suffixtree

type Match struct {
	P1, P2 Pos
	Len    Pos
}

type contextList struct {
	lists map[int]*posList
}

func newContextList() *contextList {
	return &contextList{make(map[int]*posList)}
}

type posList struct {
	positions []Pos
}

func newPosList() *posList {
	return &posList{make([]Pos, 0)}
}

func (p *posList) append(p2 *posList) {
	p.positions = append(p.positions, p2.positions...)
}

func (p *posList) add(pos Pos) {
	p.positions = append(p.positions, pos)
}

func (c *contextList) combine(c2 *contextList, length, threshold int, ch chan<- Match) {
	if length < threshold {
		return
	}
	for lc1, pl1 := range c.lists {
		for lc2, pl2 := range c2.lists {
			if lc1 != lc2 {
				for _, p1 := range pl1.positions {
					for _, p2 := range pl2.positions {
						ch <- Match{p1, p2, Pos(length)}
					}
				}
			}
		}
	}
	c.append(c2)
}

func (c *contextList) append(c2 *contextList) {
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

func walkTrans(parent *tran, length, threshold int, ch chan<- Match) *contextList {
	s := parent.state

	cl := newContextList()

	if len(s.trans) == 0 {
		pl := newPosList()
		start := parent.end + 1 - Pos(length)
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

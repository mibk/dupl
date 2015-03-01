package suffixtree

import "testing"

func TestNew(t *testing.T) {
	str := "cacao"
	s := genStates(8, str)
	// s[0] is root
	s[0].addTran(0, 1, s[1]) // ca
	s[0].addTran(1, 1, s[2]) // a
	s[0].addTran(4, 4, s[3]) // o

	s[1].addTran(2, 4, s[4]) // cao
	s[1].addTran(4, 4, s[5]) // o

	s[2].addTran(2, 4, s[4]) // cao
	s[2].addTran(4, 4, s[5]) // o

	cacao := New(str)
	compareTrees(t, s[0], cacao.root)

	str2 := "banana$"
	r := genStates(11, str2)
	// r[0] is root
	r[0].addTran(0, 6, r[1]) // banana$
	r[0].addTran(1, 1, r[2]) // a
	r[0].addTran(2, 3, r[3]) // na
	r[0].addTran(6, 6, r[4]) // $

	r[2].addTran(2, 3, r[5]) // na
	r[2].addTran(6, 6, r[6]) // $

	r[3].addTran(4, 6, r[7]) // na$
	r[3].addTran(6, 6, r[8]) // $

	r[5].addTran(4, 6, r[9])  // na$
	r[5].addTran(6, 6, r[10]) // $

	banana := New(str2)
	compareTrees(t, r[0], banana.root)
}

func compareTrees(t *testing.T, expected, actual *state) {
	ch1, ch2 := walker(expected), walker(actual)
	for {
		etran, ok1 := <-ch1
		atran, ok2 := <-ch2
		if !ok1 || !ok2 {
			if ok1 {
				t.Error("expected tree is longer")
			} else if ok2 {
				t.Error("actual tree is longer")
			}
			break
		}
		if etran.start != atran.start || etran.ActEnd() != atran.ActEnd() {
			t.Errorf("got transition (%d, %d) '%s', want (%d, %d) '%s'",
				atran.start, atran.ActEnd(), actual.data[atran.start:atran.ActEnd()],
				etran.start, etran.ActEnd(), expected.data[etran.start:etran.ActEnd()+1],
			)
		}
	}
}

func walker(s *state) <-chan *tran {
	ch := make(chan *tran)
	go func() {
		walk(s, ch)
		close(ch)
	}()
	return ch
}

func walk(s *state, ch chan<- *tran) {
	for _, tr := range s.trans {
		ch <- tr
		walk(tr.state, ch)
	}
}

func genStates(count int, data string) []*state {
	states := make([]*state, count)
	for i := range states {
		states[i] = newState(data)
	}
	return states
}

type refPair struct {
	s          *state
	start, end pos
}

func TestCanonize(t *testing.T) {
	s := genStates(4, "somebanana")
	s[0].addTran(0, 3, s[1])
	s[1].addTran(4, 6, s[2])
	s[2].addTran(7, Inf, s[3])

	find := func(needle *state) int {
		for i, state := range s {
			if state == needle {
				return i
			}
		}
		return -1
	}

	var testCases = []struct {
		origin, expected refPair
	}{
		{refPair{s[0], 0, 0}, refPair{s[0], 0, 0}},
		{refPair{s[0], 0, 2}, refPair{s[0], 0, 0}},
		{refPair{s[0], 0, 3}, refPair{s[1], 4, 0}},
		{refPair{s[0], 0, 8}, refPair{s[2], 7, 0}},
		{refPair{s[0], 0, 6}, refPair{s[2], 7, 0}},
		{refPair{s[0], 0, 100}, refPair{s[2], 7, 0}},
	}

	for _, tc := range testCases {
		s, start := canonize(tc.origin.s, tc.origin.start, tc.origin.end)
		if s != tc.expected.s || start != tc.expected.start {
			t.Errorf("for origin ref. pair (%d, (%d, %d)) got (%d, %d), want (%d, %d)",
				find(tc.origin.s), tc.origin.start, tc.origin.end,
				find(s), start,
				find(tc.expected.s), tc.expected.start,
			)
		}
	}
}

func TestSplitting(t *testing.T) {
	data := "banana"
	s1 := newState(data)
	s2 := newState(data)
	s1.addTran(0, 3, s2)

	// active point is (s1, 0, -1), an explicit state
	rets, end := testAndSplit(s1, 0, -1, 'c')
	if rets != s1 {
		t.Errorf("got state %p, want %p", rets, s1)
	}
	if end {
		t.Error("should not be an end-point")
	}
	_, end = testAndSplit(s1, 0, -1, 'b')
	if !end {
		t.Error("should be an end-point")
	}

	// active point is (s1, 0, 2), an implicit state
	rets, end = testAndSplit(s1, 0, 2, 'a')
	if rets != s1 {
		t.Error("returned state should be unchanged")
	}
	if !end {
		t.Error("should be an end-point")
	}

	// [s1]-banana->[s2] => [s1]-ban->[rets]-ana->[s2]
	rets, end = testAndSplit(s1, 0, 2, 'o')
	tr := s1.findTran('b')
	if tr == nil {
		t.Error("should have a b-transition")
	} else if tr.state != rets {
		t.Errorf("got state %p, want %p", tr.state, rets)
	}
	tr2 := rets.findTran('a')
	if tr2 == nil {
		t.Error("should have an a-transition")
	} else if tr2.state != s2 {
		t.Errorf("got state %p, want %p", tr2.state, s2)
	}
	if end {
		t.Error("should not be an end-point")
	}
}

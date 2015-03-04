package suffixtree

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

const Infinity = math.MaxInt32

// pos denotes position in data string.
type pos int32

// STree is a struct representing a suffix tree.
type STree struct {
	data     string
	root     *state
	auxState *state // auxiliary state

	// active point
	s          *state
	start, end pos
}

// New creates new suffix tree.
func New() *STree {
	t := new(STree)
	t.root = newState(t)
	t.auxState = newState(t)
	t.root.linkState = t.auxState
	t.s = t.root
	return t
}

// Update refreshes the suffix tree to by new data.
func (t *STree) Update(data string) {
	spos := pos(len(t.data))
	t.data += data
	for i := range data {
		t.update()
		t.s, t.start = t.canonize(t.s, t.start, spos+pos(i))
	}
}

// update transforms suffix tree T(n) to T(n+1).
func (t *STree) update() {
	oldr := t.root

	// (s, (start, end)) is the canonical reference pair for the active point
	s := t.s
	start, end := t.start, t.end
	r, endPoint := t.testAndSplit(s, start, end-1)
	for !endPoint {
		r.fork(end)
		if oldr != t.root {
			oldr.linkState = r
		}
		oldr = r
		s, start = t.canonize(s.linkState, start, end-1)
		r, endPoint = t.testAndSplit(s, start, end-1)
	}
	if oldr != t.root {
		oldr.linkState = r
	}

	// update active point
	t.s = s
	t.start = start
	t.end++
}

// testAndSplit tests whether a state with canonical ref. pair
// (s, (start, end)) is the end point, that is, a state that have
// a c-transition. If not, then state (exs, (start, end)) is made
// explicit (if not already so).
func (t *STree) testAndSplit(s *state, start, end pos) (exs *state, endPoint bool) {
	c := t.data[t.end]
	if start <= end {
		tr := s.findTran(s.t.data[start])
		splitPoint := tr.start + end - start + 1
		if s.t.data[splitPoint] == c {
			return s, true
		}
		// make the (s, (start, end)) state explicit
		newSt := newState(s.t)
		newSt.addTran(splitPoint, tr.end, tr.state)
		tr.end = splitPoint - 1
		tr.state = newSt
		return newSt, false
	}
	if s == t.auxState || s.findTran(c) != nil {
		return s, true
	}
	return s, false
}

// canonize returns updated state and start position for ref. pair
// (s, (start, end)) of state r so the new ref. pair is canonical,
// that is, referenced from the closest explicit ancestor of r.
func (t *STree) canonize(s *state, start, end pos) (*state, pos) {
	if start > end {
		return s, start
	} else if s == t.auxState {
		return t.root, start + 1
	}

	var tr *tran
	for {
		if start <= end {
			tr = s.findTran(s.t.data[start])
			if tr == nil {
				panic(fmt.Sprintf("there should be some transition for '%c' at %d", s.t.data[start], start))
			}
		}
		if tr.end-tr.start > end-start {
			break
		}
		start += tr.end - tr.start + 1
		s = tr.state
	}
	if s == nil {
		panic("there should always be some suffix link resolution")
	}
	return s, start
}

func (t *STree) String() string {
	buf := new(bytes.Buffer)
	printState(buf, t.root, 0)
	return buf.String()
}

func printState(buf *bytes.Buffer, s *state, ident int) {
	for _, tr := range s.trans {
		fmt.Fprint(buf, strings.Repeat("  ", ident))
		fmt.Fprintf(buf, "* (%d, %d) '%s'\n", tr.start, tr.ActEnd(), s.t.data[tr.start:tr.ActEnd()+1])
		printState(buf, tr.state, ident+1)
	}
}

// state is an explicit state of the suffix tree.
type state struct {
	t         *STree
	trans     []*tran
	linkState *state
}

func newState(t *STree) *state {
	return &state{
		t:         t,
		trans:     make([]*tran, 0),
		linkState: nil,
	}
}

func (s *state) addTran(start, end pos, r *state) {
	s.trans = append(s.trans, newTran(start, end, r))
}

// fork creates a new branch from the state s.
func (s *state) fork(i pos) *state {
	r := newState(s.t)
	s.addTran(i, Infinity, r)
	return r
}

// findTran finds c-transition.
func (s *state) findTran(c byte) *tran {
	for _, tran := range s.trans {
		if s.t.data[tran.start] == c {
			return tran
		}
	}
	return nil
}

// tran represents a state's transition.
type tran struct {
	start, end pos
	state      *state
}

func newTran(start, end pos, s *state) *tran {
	return &tran{start, end, s}
}

// ActEnd returns actual end position as consistent with
// the actual length of the data in the STree.
func (t *tran) ActEnd() pos {
	if t.end == Infinity {
		return pos(len(t.state.t.data)) - 1
	}
	return t.end
}

package suffixtree

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

const Inf = math.MaxInt8

// pos denotes position in data string
type pos int8

// STree is a struct representing a suffix tree
type STree struct {
	data string
	root *state
}

var (
	ai             = 0
	root           *state
	auxiliaryState *state
)

// New creates new suffix tree
func New(data string) *STree {
	tree := &STree{data: data, root: newState(data)}

	root = tree.root
	auxiliaryState = newState(data)
	auxiliaryState.id = 0
	ai--
	tree.root.linkState = auxiliaryState
	s := tree.root
	k := pos(0)
	for i := range data {
		s, k = update(s, k, pos(i))
		s, k = canonize(s, k, pos(i))
	}
	return tree
}

func (t *STree) String() string {
	buf := new(bytes.Buffer)
	printState(buf, t.root, 0)
	return buf.String()
}

func printState(buf *bytes.Buffer, s *state, ident int) {
	fmt.Fprint(buf, strings.Repeat("  ", ident))
	fmt.Fprintln(buf, "id:", s.id)
	for _, tr := range s.trans {
		fmt.Fprint(buf, strings.Repeat("  ", ident))
		fmt.Fprintf(buf, "- tran: %d, %d;  '%s'\n", tr.start, tr.ActEnd(), s.data[tr.start:tr.ActEnd()+1])
		printState(buf, tr.state, ident+1)
	}
}

// state is an explicit state of the suffix tree
type state struct {
	id        int
	data      string
	trans     []*tran
	linkState *state
}

func newState(data string) *state {
	ai++
	return &state{
		id:        ai,
		data:      data,
		trans:     make([]*tran, 0),
		linkState: nil,
	}
}

func (s *state) addTran(start, end pos, r *state) {
	s.trans = append(s.trans, newTran(start, end, r))
}

func (s *state) fork(i pos) *state {
	r := newState(s.data)
	s.addTran(i, Inf, r)
	return r
}

func (s *state) findTran(c byte) *tran {
	for _, tran := range s.trans {
		if s.data[tran.start] == c {
			return tran
		}
	}
	return nil
}

// tran represents a state's transition
type tran struct {
	start, end pos
	state      *state
}

func newTran(start, end pos, s *state) *tran {
	return &tran{start, end, s}
}

func (t *tran) ActEnd() pos {
	if t.end == Inf {
		return pos(len(t.state.data)) - 1
	}
	return t.end
}

func update(s *state, start, end pos) (*state, pos) {
	// (s, (start, end-1)) is the canonical reference pair for the active point
	var oldr *state = root

	r, endPoint := testAndSplit(s, start, end-1, s.data[end])
	for !endPoint {
		r.fork(end)
		if oldr != root {
			oldr.linkState = r
		}
		oldr = r
		s, start = canonize(s.linkState, start, end-1)
		r, endPoint = testAndSplit(s, start, end-1, s.data[end])
	}
	if oldr != root {
		oldr.linkState = r
	}
	return s, start
}

// testAndSplit tests whether a state with canonical ref. pair
// (s, (start, end)) is the end point, that is, a state that have
// a c-transition. If not, then state (exs, (start, end)) is made
// explicit (if not already so).
func testAndSplit(s *state, start, end pos, c byte) (exs *state, endPoint bool) {
	if start <= end {
		tr := s.findTran(s.data[start])
		splitPoint := tr.start + end - start + 1
		if s.data[splitPoint] == c {
			return s, true
		}
		// make the (s, (start, end)) state explicit
		newSt := newState(s.data)
		newSt.addTran(splitPoint, tr.end, tr.state)
		tr.end = splitPoint - 1
		tr.state = newSt
		return newSt, false
	}
	if s != auxiliaryState && s.findTran(c) == nil {
		return s, false
	}
	return s, true
}

// canonize returns updated state and start position for ref. pair
// (s, (start, end)) of state r so the new ref. pair is canonical,
// that is, referenced from the closest explicit ancestor of r.
func canonize(s *state, start, end pos) (*state, pos) {
	if start > end {
		return s, start
	} else if s == auxiliaryState {
		return root, start + 1
	}

	tr := s.findTran(s.data[start])
	if tr == nil {
		panic(fmt.Sprintf("there should be some transition for '%c' at %d", s.data[start], start))
	}
	for tr.end-tr.start <= end-start {
		start += tr.end - tr.start + 1
		s = tr.state
		if start <= end {
			tr = s.findTran(s.data[start])
			if tr == nil {
				panic(fmt.Sprintf("there should be some transition for '%c' at %d", s.data[start], start))
			}
		}
	}
	if s == nil {
		panic("there should always be some suffix link resolution")
	}
	return s, start
}

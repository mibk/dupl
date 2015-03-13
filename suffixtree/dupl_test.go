package suffixtree

import (
	"fmt"
	"testing"
)

func (m Match) String() string {
	return fmt.Sprintf("(%d, %d, %d)", m.P1, m.P2, m.Len)
}

func TestFindingDupl(t *testing.T) {
	testCases := []struct {
		s         string
		threshold int
		matches   []Match
	}{
		{"abab$", 3, []Match{}},
		{"abab$", 2, []Match{{0, 2, 2}}},
		{"abcbcabc$", 3, []Match{{0, 5, 3}}},
		{"abcbcabc$", 2, []Match{{0, 5, 3}, {1, 3, 2}, {3, 6, 2}}},
		{`All work and no play makes Jack a dull boy
All work and no play makes Jack a dull boy$`, 4, []Match{{0, 43, 42}}},
	}

	for _, tc := range testCases {
		tree := New()
		tree.Update(str2tok(tc.s)...)
		ch := tree.FindDuplOver(tc.threshold)
		for _, exp := range tc.matches {
			act, ok := <-ch
			if !ok {
				t.Errorf("missing match %v for '%s'", exp, tc.s)
			} else if exp.P1 != act.P1 || exp.P2 != act.P2 || exp.Len != act.Len {
				t.Errorf("got %v, want %v", act, exp)
			}
		}
		for {
			act, ok := <-ch
			if !ok {
				break
			}
			t.Errorf("beyond expected match %v for '%s'", act, tc.s)
		}
	}
}

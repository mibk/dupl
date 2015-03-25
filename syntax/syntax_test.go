package syntax

import (
	"testing"

	"fm.tul.cz/dupl/suffixtree"
)

func TestSerialization(t *testing.T) {
	n := genNodes(7)
	n[0].AddChildren(n[1], n[2], n[3])
	n[1].AddChildren(n[4], n[5])
	n[2].AddChildren(n[6])
	m := genNodes(6)
	m[0].AddChildren(m[1], m[2], m[3], m[4], m[5])
	testCases := []struct {
		t        *Node
		expected []int
	}{
		{n[0], []int{6, 2, 0, 0, 1, 0, 0}},
		{m[0], []int{5, 0, 0, 0, 0, 0}},
	}

	for _, tc := range testCases {
		compareSeries(t, Serialize(tc.t), tc.expected)
	}
}

func genNodes(cnt int) []*Node {
	nodes := make([]*Node, cnt)
	for i := range nodes {
		nodes[i] = NewNode()
	}
	return nodes
}

func compareSeries(t *testing.T, stream []*Node, owns []int) {
	if len(stream) != len(owns) {
		t.Errorf("series aren't the same length; got %d, want %d", len(stream), len(owns))
		return
	}
	for i, item := range stream {
		if item.Owns != owns[i] {
			t.Errorf("got %d, want %d", item.Owns, owns[i])
		}
	}
}

func TestCyclicDupl(t *testing.T) {
	testCases := []struct {
		seq      string
		indexes  []suffixtree.Pos
		expected bool
	}{
		{"a1b0a2b0", []suffixtree.Pos{0, 2}, false},
		{"a1b0a1b0", []suffixtree.Pos{0, 2}, true},
		{"a0a0", []suffixtree.Pos{0, 1}, true},
		{"a1b0c1b0a1b0c1b0", []suffixtree.Pos{0, 2, 4, 6}, true},
		{"a1b0c1b0a1b0", []suffixtree.Pos{0, 2, 4}, false},
	}

	for _, tc := range testCases {
		stree := suffixtree.New()
		for _, n := range str2nodes(tc.seq) {
			stree.Update(n)
		}
		if tc.expected != isCyclic(stree, tc.indexes, suffixtree.Pos(0)) {
			t.Errorf("for seq %v, indexes %v got %t, want %t", tc.seq, tc.indexes, !tc.expected, tc.expected)
		}
	}
}

// str2nodes converts strint to a sequence of *Node by following principle:
//   - node is represented by 2 characters
//   - first character is node type
//   - second character is the number for Node.Owns.
func str2nodes(str string) []*Node {
	chars := []rune(str)
	nodes := make([]*Node, len(chars)/2)
	for i := 0; i < len(chars)-1; i += 2 {
		nodes[i/2] = &Node{Type: int(chars[i]), Owns: int(chars[i+1] - '0')}
	}
	return nodes
}

package syntax

import "testing"

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

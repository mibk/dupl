package syntax

import (
	"fmt"

	"fm.tul.cz/dupl/suffixtree"
)

type Node struct {
	Type     int
	Filename string
	Pos, End int
	Children []*Node
	Owns     int
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) AddChildren(children ...*Node) {
	n.Children = append(n.Children, children...)
}

func (n *Node) Val() int {
	return n.Type
}

func Serialize(n *Node) []*Node {
	stream := make([]*Node, 0, 10)
	serial(n, &stream)
	return stream
}

func serial(n *Node, stream *[]*Node) int {
	*stream = append(*stream, n)
	var count int
	for _, child := range n.Children {
		count += serial(child, stream)
	}
	n.Owns = count
	return count + 1
}

// FindSyntaxUnits finds all complete syntax units in the match pair and returns them.
func FindSyntaxUnits(stree *suffixtree.STree, m suffixtree.Match) ([]*Node, []*Node) {
	i := 0
	list1 := make([]*Node, 0)
	list2 := make([]*Node, 0)
	for i < int(m.Len) {
		n1 := getNode(stree.At(m.P1 + suffixtree.Pos(i)))
		n2 := getNode(stree.At(m.P2 + suffixtree.Pos(i)))
		if n1.Owns == n2.Owns {
			if n1.Owns >= int(m.Len)-i {
				// not complete syntax unit
				i++
				continue
			}
			list1 = append(list1, n1)
			list2 = append(list2, n2)
			i += n1.Owns + 1

		} else if n1.Owns > n2.Owns {
			i += n1.Owns
		} else {
			i += n2.Owns
		}
	}
	return list1, list2
}

func getNode(tok suffixtree.Token) *Node {
	if n, ok := tok.(*Node); ok {
		return n
	}
	panic(fmt.Sprintf("tok (type %T)  is not type *Node", tok))
}

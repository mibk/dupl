package syntax

import (
	"fmt"

	"fm.tul.cz/dupl/suffixtree"
)

type Seq struct {
	Nodes []*Node
}

func newSeq(cnt int) *Seq {
	return &Seq{make([]*Node, cnt)}
}

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
func FindSyntaxUnits(stree *suffixtree.STree, m suffixtree.Match, threshold int) []*Seq {
	i := 0
	indexes := make([]suffixtree.Pos, 0)
	for i < int(m.Len) {
		n := getNode(stree.At(m.Ps[0] + suffixtree.Pos(i)))
		if n.Owns >= int(m.Len)-i {
			// not complete syntax unit
			i++
			continue
		} else if n.Owns >= threshold {
			indexes = append(indexes, suffixtree.Pos(i))
		}
		i += n.Owns + 1
	}
	if len(indexes) == 0 {
		return make([]*Seq, 0)
	}
	groups := make([]*Seq, len(m.Ps))
	for i, pos := range m.Ps {
		groups[i] = newSeq(len(indexes))
		for j, index := range indexes {
			groups[i].Nodes[j] = getNode(stree.At(pos + index))
		}
	}
	return groups
}

func getNode(tok suffixtree.Token) *Node {
	if n, ok := tok.(*Node); ok {
		return n
	}
	panic(fmt.Sprintf("tok (type %T)  is not type *Node", tok))
}

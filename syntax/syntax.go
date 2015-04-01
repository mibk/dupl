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
	Addr     string
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

	indexCnt := len(indexes)
	if indexCnt > 0 {
		lasti := indexes[indexCnt-1]
		firstn := getNode(stree.At(m.Ps[0] + lasti))
		for i := 1; i < len(m.Ps); i++ {
			n := getNode(stree.At(m.Ps[i] + lasti))
			if firstn.Owns != n.Owns {
				indexes = indexes[:indexCnt-1]
				break
			}
		}
	}
	if len(indexes) == 0 || isCyclic(stree, indexes, m.Ps[0]) {
		return make([]*Seq, 0)
	}
	seqs := make([]*Seq, len(m.Ps))
	for i, pos := range m.Ps {
		seqs[i] = newSeq(len(indexes))
		for j, index := range indexes {
			seqs[i].Nodes[j] = getNode(stree.At(pos + index))
		}
	}
	return seqs
}

// isCyclic finds out whether there is a repetive pattern in the found clone. If positive,
// it return false to point out that the clone would be redundant.
func isCyclic(stree *suffixtree.STree, indexes []suffixtree.Pos, startPos suffixtree.Pos) bool {
	cnt := len(indexes)
	if cnt <= 1 {
		return false
	}

	alts := make(map[int]bool)
	for i := 1; i <= cnt/2; i++ {
		alts[i] = true
	}

	for i := startPos; i < startPos+indexes[cnt/2]; i++ {
		nstart := getNode(stree.At(i + indexes[0]))
		for alt := range alts {
			for j := alt; j < cnt; j += alt {
				nalt := getNode(stree.At(i + indexes[j]))
				if nstart.Owns != nalt.Owns || nstart.Type != nalt.Type {
					delete(alts, alt)
				}
			}
		}
		if len(alts) == 0 {
			return false
		}
	}
	return true
}

func getNode(tok suffixtree.Token) *Node {
	if n, ok := tok.(*Node); ok {
		return n
	}
	panic(fmt.Sprintf("tok (type %T)  is not type *Node", tok))
}

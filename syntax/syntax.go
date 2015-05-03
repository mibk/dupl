package syntax

import (
	"crypto/sha1"
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

// FindSyntaxUnits finds all complete syntax units in the match group and returns them
// with the corresponding hash.
func FindSyntaxUnits(nodeSeqs [][]*Node, threshold int) ([]*Seq, string) {
	indexes := getUnitsIndexes(nodeSeqs[0], threshold)

	// TODO: is this really working?
	indexCnt := len(indexes)
	if indexCnt > 0 {
		lasti := indexes[indexCnt-1]
		firstn := nodeSeqs[0][lasti]
		for i := 1; i < len(nodeSeqs); i++ {
			n := nodeSeqs[i][lasti]
			if firstn.Owns != n.Owns {
				indexes = indexes[:indexCnt-1]
				break
			}
		}
	}
	if len(indexes) == 0 || isCyclic(indexes, nodeSeqs[0]) {
		return make([]*Seq, 0), ""
	}
	seqs := make([]*Seq, len(nodeSeqs))
	for i, nodes := range nodeSeqs {
		seqs[i] = newSeq(len(indexes))
		for j, index := range indexes {
			seqs[i].Nodes[j] = nodes[index]
		}
	}

	lastIndex := indexes[len(indexes)-1]
	firstSeq := nodeSeqs[0]
	hash := hashSeq(firstSeq[indexes[0] : lastIndex+firstSeq[lastIndex].Owns])
	return seqs, hash
}

func getUnitsIndexes(nodeSeq []*Node, threshold int) []int {
	indexes := make([]int, 0)
	for i := 0; i < len(nodeSeq); {
		n := nodeSeq[i]
		if n.Owns >= len(nodeSeq)-i {
			// not complete syntax unit
			i++
			continue
		} else if n.Owns+1 >= threshold {
			indexes = append(indexes, i)
		}
		i += n.Owns + 1
	}
	return indexes
}

// isCyclic finds out whether there is a repetive pattern in the found clone. If positive,
// it return false to point out that the clone would be redundant.
func isCyclic(indexes []int, nodes []*Node) bool {
	cnt := len(indexes)
	if cnt <= 1 {
		return false
	}

	alts := make(map[int]bool)
	for i := 1; i <= cnt/2; i++ {
		if cnt%i == 0 {
			alts[i] = true
		}
	}

	for i := 0; i < indexes[cnt/2]; i++ {
		nstart := nodes[i+indexes[0]]
	AltLoop:
		for alt := range alts {
			for j := alt; j < cnt; j += alt {
				index := i + indexes[j]
				if index < len(nodes) {
					nalt := nodes[index]
					if nstart.Owns == nalt.Owns && nstart.Type == nalt.Type {
						continue
					}
				} else if i >= indexes[alt] {
					return true
				}
				delete(alts, alt)
				continue AltLoop
			}
		}
		if len(alts) == 0 {
			return false
		}
	}
	return true
}

func GetNodes(stree *suffixtree.STree, m suffixtree.Match) [][]*Node {
	seqs := make([][]*Node, len(m.Ps))
	for i, pos := range m.Ps {
		seq := make([]*Node, m.Len)
		for j := suffixtree.Pos(0); j < m.Len; j++ {
			seq[j] = getNode(stree.At(pos + j))
		}
		seqs[i] = seq
	}
	return seqs
}

func getNode(tok suffixtree.Token) *Node {
	if n, ok := tok.(*Node); ok {
		return n
	}
	panic(fmt.Sprintf("tok (type %T)  is not type *Node", tok))
}

func hashSeq(nodes []*Node) string {
	h := sha1.New()
	bytes := make([]byte, len(nodes))
	for i, node := range nodes {
		bytes[i] = byte(node.Type)
	}
	h.Write(bytes)
	return string(h.Sum(nil))
}

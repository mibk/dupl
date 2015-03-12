package syntax

import "go/token"

type Node struct {
	Type     int
	Pos, End token.Pos
	Children []*Node
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) AddChildren(children ...*Node) {
	n.Children = append(n.Children, children...)
}

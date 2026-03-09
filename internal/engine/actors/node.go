package actors

import "mask_of_the_tomb/internal/engine"

// Do we need this? It's just an empty node at this point...

// how do we give this access to the nodeTree info?
// for instance it should be possible to get the parent node from this component
type Node struct {
	treeID string
}

func (n *Node) Init()                        {}
func (n *Node) Update(tree *engine.NodeTree) {}
func (n *Node) SetID(ID string) {
	n.treeID = ID
}

func NewNode() *Node {
	return &Node{}
}

package actors

import (
	"mask_of_the_tomb/internal/engine"

	"github.com/ebitengine/debugui"
)

// Do we need this? It's just an empty node at this point...

// how do we give this access to the nodeTree info?
// for instance it should be possible to get the parent node from this component
type Node struct {
	treeID string `debug:"auto"`
}

func (n *Node) Init()                        {}
func (n *Node) Update(tree *engine.NodeTree) {}
func (n *Node) SetID(ID string) {
	n.treeID = ID
}

func (n *Node) DrawInspector(ctx *debugui.Context) {
	ctx.SetGridLayout(make([]int, 1), make([]int, 1))
	ctx.Text("Node")
	ctx.SetGridLayout([]int{-1, -2}, []int{0, 0})
	RenderFieldsAuto(ctx, n)
}

func NewNode() *Node {
	return &Node{}
}

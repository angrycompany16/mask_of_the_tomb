package nodeactor

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
)

type Node struct {
	treeID string `debug:"auto"`
	node   *engine.Node
}

func (n *Node) Init(cmd *engine.Commands)      {}
func (n *Node) Update(cmd *engine.Commands)    {}
func (n *Node) DrawGizmo(cmd *engine.Commands) {}
func (n *Node) OnTreeAdd(node *engine.Node, cmd *engine.Commands) {
	n.treeID = node.GetID()
	n.node = node
}

func (n *Node) DrawInspector(ctx *debugui.Context) {
	ctx.SetGridLayout(make([]int, 1), make([]int, 1))
	ctx.Header("Node", false, func() {
		utils.RenderFieldsAuto(ctx, n)
	})
}

func (n *Node) GetTreeID() string {
	return n.treeID
}

func (n *Node) GetNode() *engine.Node {
	return n.node
}

func NewNode() *Node {
	return &Node{}
}

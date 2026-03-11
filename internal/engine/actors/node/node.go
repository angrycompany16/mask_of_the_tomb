package node

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/servers"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
)

type Node struct {
	treeID string `debug:"auto"`
	node   *engine.Node
}

//	func (n *Node) Instantiate() *Node {
//		return &Node{
//			treeID: n.treeID,
//			// BRUH this is so stupid
//		}
//	}
func (n *Node) Init()                           {}
func (n *Node) Update(servers *servers.Servers) {}
func (n *Node) OnTreeAdd(node *engine.Node, servers *servers.Servers) {
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

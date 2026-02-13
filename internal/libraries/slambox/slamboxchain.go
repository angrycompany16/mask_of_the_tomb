package slambox

import "mask_of_the_tomb/internal/core/maths"

const nodeRectSize = 8

type ChainNode struct {
	rect *maths.Rect // Mostly for debugging / raycasting

	nextNode    *ChainNode
	nextNodeDir maths.Direction

	prevNode    *ChainNode
	prevNodeDir maths.Direction
}

func (cn *ChainNode) GetRect() *maths.Rect {
	return cn.rect
}

// A chain that can connect multiple slamboxes / slambox groups.
// What if we just simply used the array order as an ordering
// of the chain nodes? Branching is not allowed anyways...
type SlamboxChain struct {
	nodes         []*ChainNode
	slamboxes     []*Slambox
	slamboxGroups []*SlamboxGroup
}

func (sc *SlamboxChain) SlamSlamboxChain() {

}

func (sc *SlamboxChain) GetNodes() []*ChainNode {
	return sc.nodes
}

func (sc *SlamboxChain) GetSlamboxes() []*Slambox {
	return sc.slamboxes
}

func (sc *SlamboxChain) GetSlamboxGroups() []*SlamboxGroup {
	return sc.slamboxGroups
}

func NewSlamboxChain(positionsX, positionsY []float64, slamboxes []*Slambox, slamboxGroups []*SlamboxGroup) *SlamboxChain {
	newSlamboxChain := SlamboxChain{}
	nodes := make([]*ChainNode, 0)
	for i := range positionsX {
		nodes = append(nodes, &ChainNode{
			rect: maths.NewRect(positionsX[i], positionsY[i], nodeRectSize, nodeRectSize),
		})
	}

	for i := 0; i < len(positionsX)-1; i++ {
		nodes[i].nextNode = nodes[i+1]
		dx := positionsX[i+1] - positionsX[i]
		dy := positionsY[i+1] - positionsY[i]
		nodes[i].nextNodeDir = maths.DirFromVector(dx, dy)
	}

	for i := 1; i < len(positionsX); i++ {
		nodes[i].prevNode = nodes[i-1]
		dx := positionsX[i-1] - positionsX[i]
		dy := positionsY[i-1] - positionsY[i]
		nodes[i].prevNodeDir = maths.DirFromVector(dx, dy)
	}

	newSlamboxChain.nodes = nodes
	newSlamboxChain.slamboxes = slamboxes
	newSlamboxChain.slamboxGroups = slamboxGroups
	return &newSlamboxChain
}

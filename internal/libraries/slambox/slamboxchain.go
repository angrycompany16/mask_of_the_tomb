package slambox

import (
	"fmt"
	"mask_of_the_tomb/internal/core/maths"
	"math"
	"slices"
)

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

func (sc *SlamboxChain) Update() {
	for _, slambox := range sc.slamboxes {
		slambox.Update()
	}
	for _, slamboxGroup := range sc.slamboxGroups {
		slamboxGroup.Update()
	}
}

func (sc *SlamboxChain) GetAllSlamboxRects() []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for _, slambox := range sc.slamboxes {
		rects = append(rects, slambox.GetRect())
	}
	for _, slamboxGroup := range sc.slamboxGroups {
		rects = slices.Concat(slamboxGroup.GetSlamboxRects())
	}
	return rects
}

func (sc *SlamboxChain) GetNextDir(i int) maths.Direction {
	if i < 0 || i >= len(sc.nodes)-1 {
		fmt.Println("Invalid index")
		return maths.DirNone
	}

	return sc.nodes[i].nextNodeDir
}

func (sc *SlamboxChain) GetPrevDir(i int) maths.Direction {
	if i < 1 || i >= len(sc.nodes) {
		fmt.Println("Invalid index")
		return maths.DirNone
	}

	return sc.nodes[i].prevNodeDir
}

func (sc *SlamboxChain) DistFromNode(rect maths.Rect, dir maths.Direction, i int) float64 {
	cX0, cY0 := rect.Center()
	nodeRect := sc.nodes[i].GetRect()
	cX1, cY1 := nodeRect.Center()
	hitNode, _, _ := nodeRect.RaycastDirectional(cX0, cY0, dir)

	if hitNode {
		return maths.Norm(1, cX0-cX1, cY0-cY1)
	} else {
		return math.Inf(1)
	}
}

func (sc *SlamboxChain) FindClosestNode(x, y float64) (int, float64) {
	dist := math.Inf(1)
	var closestNodeID int
	for i, node := range sc.nodes {
		nodeRect := node.GetRect()
		cX, cY := nodeRect.Center()
		nodeDist := maths.Norm(1, cX-x, cY-y)
		if dist < nodeDist {
			continue
		}
		dist = nodeDist
		closestNodeID = i
	}
	return closestNodeID, dist
}

// func (sc *SlamboxChain) FindClosestNode(x, y float64, dir maths.Direction) (bool, int, float64) {
// 	dist := math.Inf(1)
// 	foundNode := false
// 	var closestNodeID int
// 	for i, node := range sc.nodes {
// 		nodeRect := node.GetRect()
// 		cX, cY := nodeRect.Center()
// 		nodeDist := maths.Norm(1, cX-x, cY-y)
// 		hitNode, _, _ := nodeRect.RaycastDirectional(x, y, dir)

// 		if nodeRect.IsInDirection(x, y, maths.Opposite(dir)) {
// 			foundNode = true
// 			dist = nodeDist
// 			closestNodeID = i
// 		}

// 		if !hitNode || dist < nodeDist {
// 			continue
// 		}

// 		foundNode = true
// 		dist = nodeDist
// 		closestNodeID = i
// 	}
// 	return foundNode, closestNodeID, dist
// }

func NewSlamboxChain(positionsX, positionsY []float64, slamboxes []*Slambox, slamboxGroups []*SlamboxGroup) *SlamboxChain {
	newSlamboxChain := SlamboxChain{}
	nodes := make([]*ChainNode, len(positionsX))
	for i := range positionsX {
		nodes[i] = &ChainNode{rect: maths.NewRect(positionsX[i], positionsY[i], nodeRectSize, nodeRectSize)}
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

package slambox

import (
	"fmt"
	"mask_of_the_tomb/internal/core/maths"
	"math"
	"slices"
)

const nodeRectSize = 2

type SlamCtx struct {
	againstChain  bool
	closestNodeID int
	dist          float64
}

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
type SlamboxChain struct {
	nodes         []*ChainNode
	slamboxes     []*Slambox
	slamboxGroups []*SlamboxGroup
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

func (sc *SlamboxChain) DistFromNode(x, y float64, i int) float64 {
	cX, cY := sc.nodes[i].GetRect().Center()
	return maths.Norm(1, x-cX, y-cY)
}

// Check whether (x, y) is in the bounding box spanned out by the i,j nodes.
func (sc *SlamboxChain) IsBetween(i, j int, x, y float64) bool {
	BB := maths.BB([]maths.Rect{
		*sc.nodes[i].GetRect(),
		*sc.nodes[j].GetRect(),
	})
	return BB.Contains(x, y)
}

func (sc *SlamboxChain) SortNodesByDist(x, y float64) ([]int, []float64) {
	indices := make([]int, len(sc.nodes))
	dists := make([]float64, len(sc.nodes))
	for i, node := range sc.nodes {
		cX, cY := node.GetRect().Center()
		dists[i] = maths.Norm(1, cX-x, cY-y)
		indices[i] = i
	}

	// MM my favourite insertion sort
	i := 1
	for i < len(sc.nodes) {
		j := i
		for j > 0 && dists[j-1] > dists[j] {
			var tmp float64
			tmp = dists[j]
			dists[j] = dists[j-1]
			dists[j-1] = tmp

			var _tmp int
			_tmp = indices[j]
			indices[j] = indices[j-1]
			indices[j-1] = _tmp

			j--
		}
		i++
	}
	return indices, dists
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

func (sc *SlamboxChain) GetSlamDirection(rect maths.Rect, dir maths.Direction) (bool, bool) {
	cX, cY := rect.Center()
	indices, dists := sc.SortNodesByDist(cX, cY)
	var closestNodeID int
	var dist float64
	var againstChain bool
	// NOTE: This is not completely bug-free. If an invalid closest node overlaps
	// with the slambox, that node will be picked as the closest node.
	// However, this is an unlikely edge case i think
	// I also don't know how to handle it O_o
	for i, index := range indices {
		if sc.nodes[index].GetRect().Overlapping(&rect) {
			closestNodeID = index
			dist = dists[i]
			break
		}
		if hit, _, _ := sc.nodes[index].GetRect().RaycastDirectional(cX, cY, dir); hit {
			closestNodeID = index
			dist = dists[i]
			break
		}
		if hit, _, _ := sc.nodes[index].GetRect().RaycastDirectional(cX, cY, maths.Opposite(dir)); hit {
			closestNodeID = index
			dist = dists[i]
			break
		}
	}

	thisNode := sc.nodes[closestNodeID]

	if dist > 0 {
		if closestNodeID == 0 {
			if dir == thisNode.nextNodeDir {
				againstChain = false
			} else if dir == maths.Opposite(thisNode.nextNodeDir) {
				againstChain = true
			} else {
				return false, false
			}
		} else if closestNodeID == len(sc.nodes)-1 {
			if dir == thisNode.prevNodeDir {
				againstChain = true
			} else if dir == maths.Opposite(thisNode.prevNodeDir) {
				againstChain = false
			} else {
				return false, false
			}
		} else {
			if sc.IsBetween(closestNodeID-1, closestNodeID, cX, cY) {
				if dir == thisNode.prevNodeDir {
					againstChain = true
				} else if dir == maths.Opposite(thisNode.prevNodeDir) {
					againstChain = false
				} else {
					return false, false
				}
			} else if sc.IsBetween(closestNodeID, closestNodeID+1, cX, cY) {
				if dir == thisNode.nextNodeDir {
					againstChain = false
				} else if dir == maths.Opposite(thisNode.nextNodeDir) {
					againstChain = true
				} else {
					return false, false
				}
			} else {
				return false, false
			}
		}
	} else {
		thisNode := sc.nodes[closestNodeID]
		if dir == thisNode.nextNodeDir {
			againstChain = false
		} else if dir == thisNode.prevNodeDir {
			againstChain = true
		} else {
			return false, false
		}
	}
	return true, againstChain
}

func NewSlamboxChain(positionsX, positionsY []float64, slamboxes []*Slambox, slamboxGroups []*SlamboxGroup) *SlamboxChain {
	newSlamboxChain := SlamboxChain{}
	nodes := make([]*ChainNode, len(positionsX))
	for i := range positionsX {
		nodes[i] = &ChainNode{rect: maths.NewRect(positionsX[i], positionsY[i], nodeRectSize, nodeRectSize), nextNodeDir: maths.DirNone, prevNodeDir: maths.DirNone}
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

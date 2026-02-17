package slambox

import (
	"fmt"
	"mask_of_the_tomb/internal/core/maths"
	"math"
	"slices"
)

type extensionType int

const (
	EXTEND_NONE extensionType = iota
	EXTEND_X
	EXTEND_Y
	EXTEND_XY
)

type hitKind int

const (
	NONE hitKind = iota
	SLAMBOX
	SLAMBOX_GROUP
	SLAMBOX_CHAIN
)

type QueryResult struct {
	HitKind hitKind
	Index   int
}

// Note: So far we have kept the grid tiles as simply integers
// HOWEVER this presents a huge opportunity for adding interactivity
// with the environment. If the data type is not bool, we can
// associate metadata (like forces and stuff) in each cell, and
// then apply that metadata to sprites on screen, so that (for
// instance) slamboxes can interact with grass and similar things.
// Just a thought!

// Represents an environment which can contain moving and static boxes.
type SlamboxEnvironment struct {
	TileSize      float64 // Should only ever be a whole float. Data type for convenience
	gridTiles     [][]bool
	slamboxes     []*Slambox
	slamboxGroups []*SlamboxGroup
	slamboxChains []*SlamboxChain
}

func (se *SlamboxEnvironment) Update() {
	for _, slambox := range se.slamboxes {
		slambox.Update()
	}

	for _, slamboxGroup := range se.slamboxGroups {
		slamboxGroup.Update()
	}

	for _, slamboxChain := range se.slamboxChains {
		slamboxChain.Update()
	}
}

// Constructs a list representing the tiles in the SlamboxEnvironment, using a 2D voxel meshing algorithm.
// Essentially it grows a rect as much as possible, then when the limit is reached it finds a new starting points and starts
// growing a new rect. Upon termination it returns a list of rects that cover the same space as the tilemap.
// If the tilemap is empty, a list of length 0 is returned.
func (se *SlamboxEnvironment) Rectify() []*maths.Rect {
	rectList := make([]*maths.Rect, 0)
	cornerX, cornerY := se.findNewRectCorner(rectList)
	done := false
	for !done {
		newRect := maths.NewRect(float64(cornerX)*se.TileSize, float64(cornerY)*se.TileSize, se.TileSize, se.TileSize)
		maximal := false
		if !se.validateRect(newRect, rectList) {
			return rectList
		}
		for !maximal {
			extensionType := EXTEND_NONE

			extendedX := newRect.Extended(maths.DirRight, se.TileSize)
			extendedY := newRect.Extended(maths.DirDown, se.TileSize)
			extendedXY := extendedY.Extended(maths.DirRight, se.TileSize)

			if int(newRect.Right()/se.TileSize) <= len(se.gridTiles[0])-1 {
				if se.validateRect(extendedX, rectList) {
					extensionType = EXTEND_X
				}
			}

			if int(newRect.Bottom()/se.TileSize) <= len(se.gridTiles)-1 {
				if extensionType == EXTEND_X && se.validateRect(extendedXY, rectList) {
					extensionType = EXTEND_XY
				} else if se.validateRect(extendedY, rectList) {
					extensionType = EXTEND_Y
				}
			}

			switch extensionType {
			case EXTEND_NONE:
				rectList = append(rectList, newRect)
				cornerX, cornerY = se.findNewRectCorner(rectList)
				if cornerX == 0 && cornerY == 0 {
					done = true
				}
				// What
				if !done && len(rectList) != 0 {
				}
				maximal = true
			case EXTEND_X:
				newRect = extendedX
			case EXTEND_Y:
				newRect = extendedY
			case EXTEND_XY:
				newRect = extendedXY
			}
		}
	}

	return rectList
}

// Tests whether the rect passed in overlaps with:
//   - An empty tile
//   - Another rect in otherRects
func (se *SlamboxEnvironment) validateRect(rect *maths.Rect, otherRects []*maths.Rect) bool {
	for y := range se.gridTiles {
		for x := range se.gridTiles[y] {
			cX, cY := se.getCenterPos(x, y)
			if rect.Contains(float64(cX), float64(cY)) && !se.gridTiles[y][x] {
				return false
			}
		}
	}
	for _, otherRect := range otherRects {
		if rect.Overlapping(otherRect) {
			return false
		}
	}
	return true
}

// Gets centered position of a tile coord.
func (se *SlamboxEnvironment) getCenterPos(x, y int) (float64, float64) {
	return float64(x)*se.TileSize + se.TileSize/2, float64(y)*se.TileSize + se.TileSize/2
}

// Loops through the environment tiles and checks for tiles that
// are not contained in any of the elements in rectList.
// If none exist, returns 0, 0.
func (se *SlamboxEnvironment) findNewRectCorner(rectList []*maths.Rect) (int, int) {
	for y := range se.gridTiles {
		for x := range se.gridTiles[y] {
			if !se.gridTiles[y][x] {
				continue
			}
			cX, cY := se.getCenterPos(x, y)
			free := true
			for _, rect := range rectList {
				if rect.Contains(cX, cY) {
					free = false
				}
			}

			if free {
				return x, y
			}

			if len(rectList) == 0 {
				return x, y
			}
		}
	}
	return 0, 0
}

// Slams the slambox at index i in the array through the environment.
// Assumes that the index is within the array.
func (se *SlamboxEnvironment) SlamSlambox(index int, dir maths.Direction) {
	slambox := se.slamboxes[index]
	otherRects := slices.Concat(se.GetSlamboxRects(index), se.GetSlamboxGroupRects(-1), se.GetSlamboxChainRects(-1))
	projRect, _ := se.ProjectRect(*slambox.GetRect(), dir, math.Inf(1), otherRects)
	slambox.Slam(projRect.Left(), projRect.Top())
}

// Slams the slambox group at index i in the array through the environment.
// Assumes that the index is within the array.
func (se *SlamboxEnvironment) SlamSlamboxGroup(i int, dir maths.Direction) {
	slamboxGroup := se.slamboxGroups[i]
	rects := slamboxGroup.GetSlamboxRects()
	otherRects := slices.Concat(se.GetSlamboxRects(i), se.GetSlamboxGroupRects(-1), se.GetSlamboxChainRects(-1))
	newRects, _ := se.ProjectRects(rects, dir, math.Inf(1), otherRects)
	slamboxGroup.Slam(newRects)
}

// Slams the slambox chain at index i in the array through the environment.
// objectID should be an index in the array of slamboxes / slambox
// groups. To indicate which one, the objectID is used.
// Assumes that the index is within the array.
// Slightly complicated.
func (se *SlamboxEnvironment) SlamSlamboxChain(i int, objectID int, isGroup bool, dir maths.Direction) {
	chain := se.slamboxChains[i]
	nodes := chain.GetNodes()
	slamboxes := chain.GetSlamboxes()
	slamboxGroups := chain.GetSlamboxGroups()
	otherRects := slices.Concat(se.GetSlamboxRects(-1), se.GetSlamboxGroupRects(-1), se.GetSlamboxChainRects(i))

	// var boxList
	var targetSlambox *Slambox
	if isGroup {
		targetSlamboxGroup := chain.GetSlamboxGroups()[objectID]
		targetSlambox = targetSlamboxGroup.GetCenterSlambox()
	} else {
		targetSlambox = chain.GetSlamboxes()[objectID]
	}

	validDirection, againstChain := chain.GetSlamDirection(*targetSlambox.GetRect(), dir)
	if !validDirection {
		fmt.Println("Invalid direction")
		return
	}

	// Determine loop params based on slam directions
	var i0 int
	var pred func(i int) bool
	var inc func(i int) int
	if againstChain {
		i0 = len(nodes) - 1
		pred = func(i int) bool { return i > 0 }
		inc = func(i int) int { return i - 1 }
	} else {
		i0 = 0
		pred = func(i int) bool { return i < len(nodes)-1 }
		inc = func(i int) int { return i + 1 }
	}

	// Find minimal distance
	minDist := math.Inf(1)
	for i := i0; pred(i); {
		var nodeDir maths.Direction
		var nextIndex int
		if againstChain {
			nodeDir = chain.GetPrevDir(i)
			nextIndex = i - 1
		} else {
			nodeDir = chain.GetNextDir(i)
			nextIndex = i + 1
		}

		for _, slambox := range slamboxes {
			if chain.IsBetween(i, nextIndex, slambox.GetRect().Cx(), slambox.GetRect().Cy()) {
				skip, nodeDist := se.checkSkip(chain, slambox, nextIndex, nodeDir)
				if skip {
					continue
				}
				_, dist := se.ProjectRect(*slambox.GetRect(), nodeDir, nodeDist, otherRects)
				if dist < minDist {
					minDist = dist
				}
			}
		}

		for _, slamboxGroup := range slamboxGroups {
			slambox := slamboxGroup.GetCenterSlambox()
			if chain.IsBetween(i, nextIndex, slambox.GetRect().Cx(), slambox.GetRect().Cy()) {
				skip, nodeDist := se.checkSkip(chain, slambox, nextIndex, nodeDir)
				if skip {
					continue
				}
				_, dist := se.ProjectRects(slamboxGroup.GetSlamboxRects(), nodeDir, nodeDist, otherRects)
				if dist < minDist {
					minDist = dist
				}
			}
		}

		i = inc(i)
	}

	// Slam slamboxes, constrained to minimal distance
	for i := i0; pred(i); {
		var nodeDir maths.Direction
		var nextIndex int
		if againstChain {
			nodeDir = chain.GetPrevDir(i)
			nextIndex = i - 1
		} else {
			nodeDir = chain.GetNextDir(i)
			nextIndex = i + 1
		}
		for _, slambox := range slamboxes {
			if chain.IsBetween(i, nextIndex, slambox.GetRect().Cx(), slambox.GetRect().Cy()) {
				if skip, _ := se.checkSkip(chain, slambox, nextIndex, nodeDir); skip {
					continue
				}
				projRect, _ := se.ProjectRect(*slambox.GetRect(), nodeDir, minDist, otherRects)
				slambox.Slam(projRect.Left(), projRect.Top())
			}
		}

		for _, slamboxGroup := range slamboxGroups {
			slambox := slamboxGroup.GetCenterSlambox()
			if chain.IsBetween(i, nextIndex, slambox.GetRect().Cx(), slambox.GetRect().Cy()) {
				if skip, _ := se.checkSkip(chain, slambox, nextIndex, nodeDir); skip {
					continue
				}
				projRects, _ := se.ProjectRects(slamboxGroup.GetSlamboxRects(), nodeDir, minDist, otherRects)
				slamboxGroup.Slam(projRects)
			}
		}
		i = inc(i)
	}
}

func (se *SlamboxEnvironment) checkSkip(chain *SlamboxChain, slambox *Slambox, nextIndex int, nodeDir maths.Direction) (bool, float64) {
	nodeDist := chain.DistFromNode(slambox.GetRect().Cx(), slambox.GetRect().Cy(), nextIndex)
	if nodeDist == 0 {
		t1 := nextIndex != 0 || nodeDir == chain.GetNextDir(nextIndex)
		t2 := nextIndex != len(chain.GetNodes())-1 || nodeDir == chain.GetPrevDir(nextIndex)
		if t1 && t2 {
			return true, nodeDist
		}
	}
	return false, nodeDist
}

// Projects a group of rects through the environment. Returns
// a list of rects with the same length as the incoming one,
// but projected in the specified direction.
func (se *SlamboxEnvironment) ProjectRects(rects []*maths.Rect, dir maths.Direction, maxDist float64, otherRects []*maths.Rect) ([]maths.Rect, float64) {
	shortestDist := math.Inf(1)
	var closestRect maths.Rect
	var closestID int
	for i, rect := range rects {
		projRect, dist := se.ProjectRect(*rect, dir, maxDist, otherRects)
		if dist < shortestDist {
			closestRect = projRect
			closestID = i
			shortestDist = dist
		}
	}

	projectedList := make([]maths.Rect, 0)
	anchor := rects[closestID]
	anchorX, anchorY := anchor.TopLeft()
	offsetX := closestRect.Left() - anchorX
	offsetY := closestRect.Top() - anchorY
	for _, otherRect := range rects {
		translatedRect := otherRect.Translated(offsetX, offsetY)
		projectedList = append(projectedList, translatedRect)
	}
	return projectedList, shortestDist
}

// Projects a rect (moves it as far as possible) through the slambox
// environment. Returns the projected rect and the distance that it was
// moved.
// Also takes in a distance constraint. To ignore this, pass inmath.Inf(1).
func (se *SlamboxEnvironment) ProjectRect(rect maths.Rect, dir maths.Direction, maxDist float64, otherRects []*maths.Rect) (maths.Rect, float64) {
	rects := slices.Concat(se.Rectify(), otherRects)

	var closestObstruction *maths.Rect
	var closestDist = math.Inf(1)
	for _, otherRect := range rects {
		hrzWithin := maths.IsBetween(rect.Left(), rect.Right()-1, otherRect.Left()) ||
			maths.IsBetween(rect.Left(), rect.Right()-1, otherRect.Right()-1) ||
			maths.IsBetween(otherRect.Left(), otherRect.Right()-1, rect.Left()) ||
			maths.IsBetween(otherRect.Left(), otherRect.Right()-1, rect.Right()-1)
		vrtWithin := maths.IsBetween(rect.Top(), rect.Bottom()-1, otherRect.Top()) ||
			maths.IsBetween(rect.Top(), rect.Bottom()-1, otherRect.Bottom()-1) ||
			maths.IsBetween(otherRect.Top(), otherRect.Bottom()-1, rect.Top()) ||
			maths.IsBetween(otherRect.Top(), otherRect.Bottom()-1, rect.Bottom()-1)
		distY := math.Abs(rect.Cy() - otherRect.Cy())
		distX := math.Abs(rect.Cx() - otherRect.Cx())
		isCloserY := distY < closestDist
		isCloserX := distX < closestDist

		switch dir {
		case maths.DirUp:
			isAbove := otherRect.Bottom() <= rect.Top()

			if !(hrzWithin && isCloserY && isAbove) {
				continue
			}
			closestObstruction = otherRect
			closestDist = distY
		case maths.DirDown:
			isBelow := otherRect.Top() >= rect.Bottom()

			if !(hrzWithin && isCloserY && isBelow) {
				continue
			}
			closestDist = distY
			closestObstruction = otherRect
		case maths.DirRight:
			isRight := otherRect.Left() >= rect.Right()

			if !(vrtWithin && isCloserX && isRight) {
				continue
			}
			closestDist = distX
			closestObstruction = otherRect
		case maths.DirLeft:
			isLeft := otherRect.Right() <= rect.Left()

			if !(vrtWithin && isCloserX && isLeft) {
				continue
			}

			closestDist = distX
			closestObstruction = otherRect
		}
	}

	if closestObstruction == nil {
		return rect, 0
	}

	var dx, dy float64
	switch dir {
	case maths.DirUp:
		dy = closestObstruction.Bottom() - rect.Top()
		dy = maths.Clamp(dy, -maxDist, 0)
	case maths.DirDown:
		dy = closestObstruction.Top() - rect.Bottom()
		dy = maths.Clamp(dy, 0, maxDist)
	case maths.DirRight:
		dx = closestObstruction.Left() - rect.Right()
		dx = maths.Clamp(dx, 0, maxDist)
	case maths.DirLeft:
		dx = closestObstruction.Right() - rect.Left()
		dx = maths.Clamp(dx, -maxDist, 0)
	}

	rect.Translate(dx, dy)
	return rect, maths.Norm(1, dx, dy)
}

// Returns the maths.Rect belonging to each slambox.
func (se *SlamboxEnvironment) GetSlamboxRects(except int) []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for i, slambox := range se.slamboxes {
		if except == i {
			continue
		}
		rects = append(rects, slambox.GetRect())
	}
	return rects
}

func (se *SlamboxEnvironment) GetSlamboxGroupRects(except int) []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for i, slamboxGroup := range se.slamboxGroups {
		if except == i {
			continue
		}
		rects = slices.Concat(rects, slamboxGroup.GetSlamboxRects())
	}
	return rects
}

func (se *SlamboxEnvironment) GetSlamboxChainRects(except int) []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for i, slamboxChain := range se.slamboxChains {
		if except == i {
			continue
		}
		rects = slices.Concat(rects, slamboxChain.GetAllSlamboxRects())
	}
	return rects
}

// Returns a list of IDs of slamboxes overlapping with the rect.
func (se *SlamboxEnvironment) CheckSlamboxOverlap(rect *maths.Rect) []int {
	overlaps := make([]int, 0)
	for i, slambox := range se.slamboxes {
		if slambox.GetRect().Overlapping(rect) {
			overlaps = append(overlaps, i)
		}
	}
	return overlaps
}

func (se *SlamboxEnvironment) GetSlamboxes() []*Slambox {
	return se.slamboxes
}

// Returns a list of IDs of slambox groups overlapping with the rect.
func (se *SlamboxEnvironment) CheckSlamboxGroupOverlap(rect *maths.Rect) []int {
	overlaps := make([]int, 0)
outer:
	for i, slamboxGroup := range se.slamboxGroups {
		for _, slambox := range slamboxGroup.GetSlamboxes() {
			if slambox.GetRect().Overlapping(rect) {
				overlaps = append(overlaps, i)
				continue outer
			}
		}
	}
	return overlaps
}

func (se *SlamboxEnvironment) GetSlamboxGroups() []*SlamboxGroup {
	return se.slamboxGroups
}

// Returns a list of IDs of slambox groups overlapping with the rect.
func (se *SlamboxEnvironment) CheckSlamboxChainOverlap(rect *maths.Rect) []int {
	overlaps := make([]int, 0)
outer:
	for i, slamboxChain := range se.slamboxChains {
		for _, slambox := range slamboxChain.GetSlamboxes() {
			if slambox.GetRect().Overlapping(rect) {
				overlaps = append(overlaps, i)
				continue outer
			}
		}

		for _, slamboxGroup := range slamboxChain.GetSlamboxGroups() {
			for _, slambox := range slamboxGroup.GetSlamboxes() {
				if slambox.GetRect().Overlapping(rect) {
					overlaps = append(overlaps, i)
					continue outer
				}
			}
		}
	}
	return overlaps
}

func (se *SlamboxEnvironment) GetSlamboxChains() []*SlamboxChain {
	return se.slamboxChains
}

func (se *SlamboxEnvironment) SetTiles(tilemap [][]int) {
	se.gridTiles = make([][]bool, 0)
	for i := range se.gridTiles {
		for j := range se.gridTiles {
			se.gridTiles[i] = append(se.gridTiles[i], tilemap[i][j] != 0)
		}
		se.gridTiles = append(se.gridTiles, make([]bool, 0))
	}
}

func (se *SlamboxEnvironment) SetTileSize(tileSize float64) {
	se.TileSize = tileSize
}

func (se *SlamboxEnvironment) AddSlambox(slambox *Slambox) {
	se.slamboxes = append(se.slamboxes, slambox)
}

func (se *SlamboxEnvironment) AddSlamboxGroup(slamboxGroup *SlamboxGroup) {
	se.slamboxGroups = append(se.slamboxGroups, slamboxGroup)
}

func (se *SlamboxEnvironment) AddSlamboxChain(slamboxChain *SlamboxChain) {
	se.slamboxChains = append(se.slamboxChains, slamboxChain)
}

func (se *SlamboxEnvironment) QuerySlamboxes(rect *maths.Rect) QueryResult {
	for i, slamboxRect := range se.GetSlamboxRects(-1) {
		if slamboxRect.Overlapping(rect) {
			return QueryResult{HitKind: SLAMBOX, Index: i}
		}
	}
	for i, slamboxGroup := range se.GetSlamboxGroups() {
		for _, slamboxGroupRect := range slamboxGroup.GetSlamboxRects() {
			if slamboxGroupRect.Overlapping(rect) {
				return QueryResult{HitKind: SLAMBOX_GROUP, Index: i}
			}
		}
	}
	for i, slamboxChain := range se.GetSlamboxChains() {
		for _, slamboxChainRect := range slamboxChain.GetAllSlamboxRects() {
			if slamboxChainRect.Overlapping(rect) {
				return QueryResult{HitKind: SLAMBOX_CHAIN, Index: i}
			}
		}
	}
	return QueryResult{}
}

func NewSlamboxEnvironment(tileSize float64, gridTiles [][]bool, slamboxes []*Slambox, slamboxGroups []*SlamboxGroup, slamboxChains []*SlamboxChain) *SlamboxEnvironment {
	newSlamboxEnvironment := SlamboxEnvironment{
		TileSize:      tileSize,
		gridTiles:     gridTiles,
		slamboxes:     slamboxes,
		slamboxGroups: slamboxGroups,
		slamboxChains: slamboxChains,
	}
	return &newSlamboxEnvironment
}

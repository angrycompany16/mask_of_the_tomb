package slambox

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/maths"
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

type QueryFilter struct {
	IgnoreSlamboxIndex int
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
	TileSize         float64 // Should only ever be a whole float. Data type for convenience
	gridTiles        [][]bool
	environmentRects []*maths.Rect
	slamboxes        []*maths.Rect
	slamboxGroups    []*SlamboxGroup
	slamboxChains    []*SlamboxChain
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

	return slices.Concat(rectList, se.environmentRects)
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

func (se *SlamboxEnvironment) SlamSlambox(index int, dir maths.Direction) (float64, float64) {
	slambox := se.slamboxes[index]
	otherRects := slices.Concat(se.GetSlamboxRects(index), se.GetSlamboxGroupRects(-1), se.GetSlamboxChainRects(-1))
	projRect, _ := se.ProjectRect(*slambox, dir, math.Inf(1), otherRects)
	return projRect.Left(), projRect.Top()
}

func (se *SlamboxEnvironment) SlamSlamboxGroup(i int, dir maths.Direction) []maths.Rect {
	slamboxGroup := se.slamboxGroups[i]
	rects := slamboxGroup.slamboxes
	otherRects := slices.Concat(se.GetSlamboxRects(-1), se.GetSlamboxGroupRects(i), se.GetSlamboxChainRects(-1))
	newRects, _ := se.ProjectRects(rects, dir, math.Inf(1), otherRects)
	return newRects
}

func (se *SlamboxEnvironment) SlamSlamboxChain(i int, objectID int, isGroup bool, dir maths.Direction) (maths.Rect, []maths.Rect) {
	chain := se.slamboxChains[i]
	otherRects := slices.Concat(se.GetSlamboxRects(-1), se.GetSlamboxGroupRects(-1), se.GetSlamboxChainRects(i))

	// var boxList
	var targetSlambox *maths.Rect
	if isGroup {
		targetSlamboxGroup := chain.slamboxGroups[objectID]
		targetSlambox = targetSlamboxGroup.slamboxes[targetSlamboxGroup.centerIndex]
	} else {
		targetSlambox = chain.slamboxes[objectID]
	}

	validDirection, againstChain := chain.GetSlamDirection(*targetSlambox, dir)
	if !validDirection {
		fmt.Println("Invalid direction")
		return maths.Rect{}, nil
	}

	// Determine loop params based on slam directions
	var i0 int
	var pred func(i int) bool
	var inc func(i int) int
	if againstChain {
		i0 = len(chain.nodes) - 1
		pred = func(i int) bool { return i > 0 }
		inc = func(i int) int { return i - 1 }
	} else {
		i0 = 0
		pred = func(i int) bool { return i < len(chain.nodes)-1 }
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

		for _, slambox := range chain.slamboxes {
			if chain.IsBetween(i, nextIndex, slambox.Cx(), slambox.Cy()) {
				skip, nodeDist := se.checkSkip(chain, slambox, nextIndex, nodeDir)
				if skip {
					continue
				}
				_, dist := se.ProjectRect(*slambox, nodeDir, nodeDist, otherRects)
				if dist < minDist {
					minDist = dist
				}
			}
		}

		for _, slamboxGroup := range chain.slamboxGroups {
			slambox := slamboxGroup.slamboxes[slamboxGroup.centerIndex]
			if chain.IsBetween(i, nextIndex, slambox.Cx(), slambox.Cy()) {
				skip, nodeDist := se.checkSkip(chain, slambox, nextIndex, nodeDir)
				if skip {
					continue
				}
				_, dist := se.ProjectRects(slamboxGroup.slamboxes, nodeDir, nodeDist, otherRects)
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
		for _, slambox := range chain.slamboxes {
			if chain.IsBetween(i, nextIndex, slambox.Cx(), slambox.Cy()) {
				if skip, _ := se.checkSkip(chain, slambox, nextIndex, nodeDir); skip {
					continue
				}
				projRect, _ := se.ProjectRect(*slambox, nodeDir, minDist, otherRects)

				return projRect, nil
			}
		}

		for _, slamboxGroup := range chain.slamboxGroups {
			slambox := slamboxGroup.slamboxes[slamboxGroup.centerIndex]
			if chain.IsBetween(i, nextIndex, slambox.Cx(), slambox.Cy()) {
				if skip, _ := se.checkSkip(chain, slambox, nextIndex, nodeDir); skip {
					continue
				}
				projRects, _ := se.ProjectRects(slamboxGroup.slamboxes, nodeDir, minDist, otherRects)

				return maths.Rect{}, projRects
			}
		}
		i = inc(i)
	}
	return maths.Rect{}, nil
}

func (se *SlamboxEnvironment) checkSkip(chain *SlamboxChain, slambox *maths.Rect, nextIndex int, nodeDir maths.Direction) (bool, float64) {
	nodeDist := chain.DistFromNode(slambox.Cx(), slambox.Cy(), nextIndex)
	if nodeDist == 0 {
		t1 := nextIndex != 0 || nodeDir == chain.GetNextDir(nextIndex)
		t2 := nextIndex != len(chain.nodes)-1 || nodeDir == chain.GetPrevDir(nextIndex)
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
// Also takes in a distance constraint. To ignore this, pass in math.Inf(1).
func (se *SlamboxEnvironment) ProjectRect(rect maths.Rect, dir maths.Direction, maxDist float64, otherRects []*maths.Rect) (maths.Rect, float64) {
	rects := slices.Concat(se.Rectify(), otherRects)

	var closestObstruction *maths.Rect
	var closestDist = math.Inf(1)
	for _, otherRect := range rects {
		// alternative method
		projXRect := rect.Translated(0, -rect.Top())
		projXOtherRect := otherRect.Translated(0, -otherRect.Top())
		projYRect := rect.Translated(-rect.Left(), 0)
		projYOtherRect := otherRect.Translated(-otherRect.Left(), 0)
		hrzWithin := projXRect.Overlapping(&projXOtherRect)
		vrtWithin := projYRect.Overlapping(&projYOtherRect)

		switch dir {
		case maths.DirUp:
			dist := math.Abs(otherRect.Bottom() - rect.Top())
			isAbove := otherRect.Bottom() <= rect.Top()

			if !(hrzWithin && isAbove && dist < closestDist) {
				continue
			}
			closestDist = dist
			closestObstruction = otherRect
		case maths.DirDown:
			dist := math.Abs(otherRect.Top() - rect.Bottom())
			isBelow := otherRect.Top() >= rect.Bottom()

			if !(hrzWithin && isBelow && dist < closestDist) {
				continue
			}
			closestDist = dist
			closestObstruction = otherRect
		case maths.DirRight:
			dist := math.Abs(otherRect.Left() - rect.Right())
			isRight := otherRect.Left() >= rect.Right()

			if !(vrtWithin && isRight && dist < closestDist) {
				continue
			}
			closestDist = dist
			closestObstruction = otherRect
		case maths.DirLeft:
			dist := math.Abs(otherRect.Right() - rect.Left())
			isLeft := otherRect.Right() <= rect.Left()

			if !(vrtWithin && isLeft && dist < closestDist) {
				continue
			}
			closestDist = dist
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

// Returns a list of IDs of slamboxes overlapping with the rect.
func (se *SlamboxEnvironment) CheckSlamboxOverlap(rect *maths.Rect) []int {
	overlaps := make([]int, 0)
	for i, slambox := range se.slamboxes {
		if slambox.Overlapping(rect) {
			overlaps = append(overlaps, i)
		}
	}
	return overlaps
}

// Returns a list of IDs of slambox groups overlapping with the rect.
func (se *SlamboxEnvironment) CheckSlamboxGroupOverlap(rect *maths.Rect) []int {
	overlaps := make([]int, 0)
outer:
	for i, slamboxGroup := range se.slamboxGroups {
		for _, slambox := range slamboxGroup.slamboxes {
			if slambox.Overlapping(rect) {
				overlaps = append(overlaps, i)
				continue outer
			}
		}
	}
	return overlaps
}

// Returns a list of IDs of slambox groups overlapping with the rect.
func (se *SlamboxEnvironment) CheckSlamboxChainOverlap(rect *maths.Rect) []int {
	overlaps := make([]int, 0)
outer:
	for i, slamboxChain := range se.slamboxChains {
		for _, slambox := range slamboxChain.slamboxes {
			if slambox.Overlapping(rect) {
				overlaps = append(overlaps, i)
				continue outer
			}
		}

		for _, slamboxGroup := range slamboxChain.slamboxGroups {
			for _, slambox := range slamboxGroup.slamboxes {
				if slambox.Overlapping(rect) {
					overlaps = append(overlaps, i)
					continue outer
				}
			}
		}
	}
	return overlaps
}

func (se *SlamboxEnvironment) QuerySlamboxes(rect *maths.Rect, filter QueryFilter) QueryResult {
	for i, slamboxRect := range se.GetSlamboxRects(-1) {
		if i == filter.IgnoreSlamboxIndex {
			continue
		}
		if slamboxRect.Overlapping(rect) {
			return QueryResult{HitKind: SLAMBOX, Index: i}
		}
	}
	for i, slamboxGroup := range se.GetSlamboxGroups() {
		for _, slamboxGroupRect := range slamboxGroup.slamboxes {
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

func (se *SlamboxEnvironment) CheckTileOverlap(rect *maths.Rect) bool {
	for _, envrect := range se.Rectify() {
		if envrect.Overlapping(rect) {
			return true
		}
	}
	return false
}

func (se *SlamboxEnvironment) CenteredToSlambox(x, y float64) (float64, float64) {
	return x - float64(len(se.gridTiles[0]))*se.TileSize/2,
		y - float64(len(se.gridTiles[0]))*se.TileSize/2
}

// ----- GETTERS -----
// Returns the maths.Rect belonging to each slambox.
func (se *SlamboxEnvironment) GetSlamboxRects(except int) []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for i, slambox := range se.slamboxes {
		if except == i {
			continue
		}
		rects = append(rects, slambox)
	}
	return rects
}

func (se *SlamboxEnvironment) GetSlamboxGroupRects(except int) []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for i, slamboxGroup := range se.slamboxGroups {
		if except == i {
			continue
		}
		rects = slices.Concat(rects, slamboxGroup.slamboxes)
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

func (se *SlamboxEnvironment) GetSlamboxes() []*maths.Rect {
	return se.slamboxes
}

func (se *SlamboxEnvironment) GetSlamboxGroups() []*SlamboxGroup {
	return se.slamboxGroups
}

func (se *SlamboxEnvironment) GetSlamboxChains() []*SlamboxChain {
	return se.slamboxChains
}

// ----- SETTERS -----
func (se *SlamboxEnvironment) SetTileSize(tileSize float64) {
	se.TileSize = tileSize
}

func (se *SlamboxEnvironment) SetTiles(tilemap [][]int) {
	se.gridTiles = make([][]bool, 0)
	for i := range tilemap {
		se.gridTiles = append(se.gridTiles, make([]bool, 0))
		for j := range tilemap[i] {
			se.gridTiles[i] = append(se.gridTiles[i], tilemap[i][j] != 0)
		}
	}
}

func (se *SlamboxEnvironment) AddEnvironmentRect(rect *maths.Rect) {
	se.environmentRects = append(se.environmentRects, rect)
}

func (se *SlamboxEnvironment) ClearEnvironmentRects(rect *maths.Rect) {
	se.environmentRects = make([]*maths.Rect, 0)
}

func (se *SlamboxEnvironment) AddSlambox(slambox *maths.Rect) int {
	se.slamboxes = append(se.slamboxes, slambox)
	return len(se.slamboxes) - 1
}

func (se *SlamboxEnvironment) AddSlamboxGroup(slamboxGroup *SlamboxGroup) int {
	se.slamboxGroups = append(se.slamboxGroups, slamboxGroup)
	return len(se.slamboxGroups) - 1
}

func (se *SlamboxEnvironment) AddSlamboxChain(slamboxChain *SlamboxChain) {
	se.slamboxChains = append(se.slamboxChains, slamboxChain)
}

func NewSlamboxEnvironment(tileSize int) *SlamboxEnvironment {
	return &SlamboxEnvironment{
		TileSize:         float64(tileSize),
		gridTiles:        make([][]bool, 0),
		environmentRects: make([]*maths.Rect, 0),
		slamboxes:        make([]*maths.Rect, 0),
		slamboxGroups:    make([]*SlamboxGroup, 0),
		slamboxChains:    make([]*SlamboxChain, 0),
	}
}

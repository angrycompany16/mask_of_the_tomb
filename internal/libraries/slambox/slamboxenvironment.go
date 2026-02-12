package slambox

import (
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
}

func (se *SlamboxEnvironment) Update() {
	for _, slambox := range se.slamboxes {
		slambox.Update()
	}

	for _, slamboxGroup := range se.slamboxGroups {
		slamboxGroup.Update()
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
				if se.validateRect(&extendedX, rectList) {
					extensionType = EXTEND_X
				}
			}

			if int(newRect.Bottom()/se.TileSize) <= len(se.gridTiles)-1 {
				if extensionType == EXTEND_X && se.validateRect(&extendedXY, rectList) {
					extensionType = EXTEND_XY
				} else if se.validateRect(&extendedY, rectList) {
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
				if !done && len(rectList) != 0 {
				}
				maximal = true
			case EXTEND_X:
				newRect = &extendedX
			case EXTEND_Y:
				newRect = &extendedY
			case EXTEND_XY:
				newRect = &extendedXY
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
func (se *SlamboxEnvironment) SlamSlambox(i int, dir maths.Direction) {
	slambox := se.slamboxes[i]
	projRect, _ := se.ProjectRect(*slambox.GetRect(), dir)
	slambox.Slam(projRect.Left(), projRect.Top())
}

// Slams the slambox group at index i in the array through the environment.
// Assumes that the index is within the array.
func (se *SlamboxEnvironment) SlamSlamboxGroup(i int, dir maths.Direction) {
	slamboxGroup := se.slamboxGroups[i]
	shortestDist := math.Inf(1)
	var closestRect maths.Rect
	var closestID int
	for i, slambox := range slamboxGroup.GetSlamboxes() {
		projRect, dist := se.ProjectRect(*slambox.GetRect(), dir)
		if dist < shortestDist {
			closestRect = projRect
			closestID = i
			shortestDist = dist
		}
	}

	slamboxGroup.Slam(closestRect.Left(), closestRect.Top(), closestID)
}

// Projects a rect (moves it as far as possible) through the slambox
// environment. Returns the projected rect and the distance that it was
// moved.
func (se *SlamboxEnvironment) ProjectRect(rect maths.Rect, dir maths.Direction) (maths.Rect, float64) {
	rects := slices.Concat(se.Rectify(), se.GetSlamboxRects())

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
	case maths.DirDown:
		dy = closestObstruction.Top() - rect.Bottom()
	case maths.DirRight:
		dx = closestObstruction.Left() - rect.Right()
	case maths.DirLeft:
		dx = closestObstruction.Right() - rect.Left()
	}

	rect.Translate(dx, dy)
	return rect, math.Abs(dx + dy)
}

// Returns the maths.Rect belonging to each slambox.
func (se *SlamboxEnvironment) GetSlamboxRects() []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for _, slambox := range se.slamboxes {
		rects = append(rects, slambox.GetRect())
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

func NewSlamboxEnvironment(tileSize float64, gridTiles [][]bool, slamboxes []*Slambox, slamboxGroups []*SlamboxGroup) *SlamboxEnvironment {
	newSlamboxEnvironment := SlamboxEnvironment{
		TileSize:      tileSize,
		gridTiles:     gridTiles,
		slamboxes:     slamboxes,
		slamboxGroups: slamboxGroups,
	}
	return &newSlamboxEnvironment
}

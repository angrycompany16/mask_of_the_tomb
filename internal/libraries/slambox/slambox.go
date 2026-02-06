package slambox

import (
	"mask_of_the_tomb/internal/core/maths"
	"math"
	"slices"
)

// Will need basic stuff like maths, hence a library

type MoveBox struct{}

type Slambox struct{}

// Note: So far we have kept the grid tiles as simply integers
// HOWEVER this presents a huge opportunity for adding interactivity
// with the environment. If the data type is not bool, we can
// associate metadata (like forces and stuff) in each cell, and
// then apply that metadata to sprites on screen, so that (for
// instance) slamboxes can interact with grass and similar things.
// Just a thought!
type SlamboxEnvironment struct {
	TileSize  float64 // Should only ever be a whole float. Data type for convenience
	GridTiles [][]bool
	Slamboxes []maths.Rect
}

// Constructs a list of rects from the GridTiles info
func (se *SlamboxEnvironment) Rectify() []maths.Rect {
	// Algorithm that converts a tilemap into rects. Greedy algorithm?
	// Simple idea:
	// Start rect at first tile
	// 1. let rLeft := If possible, extend one row left
	// 2. let rRight := If possible, extend one row up
	// If only one exists, pick that one
	// If both exists, pick both
	// If none exist, start new rect (where? This will require keeping
	// track of which tiles are covered and not)
	return nil
}

// Projects a rect (moves it as far as possible) through the slambox
// environment. Returns the projected rect and the distance that it was
// moved.
func (se *SlamboxEnvironment) ProjectRect(rect maths.Rect, dir maths.Direction) (maths.Rect, float64) {
	rects := slices.Concat(se.Rectify(), se.Slamboxes)

	var closestObstruction maths.Rect
	var closestDist = math.Inf(1)
	for _, otherRect := range rects {
		hrzWithin := maths.IsBetween(rect.Left(), rect.Right(), otherRect.Left()) ||
			maths.IsBetween(rect.Left(), rect.Right(), otherRect.Right())
		vrtWithin := maths.IsBetween(rect.Top(), rect.Bottom(), otherRect.Top()) ||
			maths.IsBetween(rect.Top(), rect.Bottom(), otherRect.Bottom())
		isCloserY := math.Abs(rect.Cy()-otherRect.Cy()) < closestDist
		isCloserX := math.Abs(rect.Cx()-otherRect.Cx()) < closestDist

		switch dir {
		case maths.DirUp:
			isBelow := otherRect.Bottom() > rect.Top()

			if !hrzWithin || !isCloserY || isBelow {
				continue
			}
			closestObstruction = otherRect
		case maths.DirDown:
			isAbove := otherRect.Bottom() < rect.Top()

			if !hrzWithin || !isCloserY || isAbove {
				continue
			}
			closestObstruction = otherRect
		case maths.DirRight:
			isLeft := otherRect.Right() < rect.Left()

			if !vrtWithin || !isCloserX || isLeft {
				continue
			}
			closestObstruction = otherRect
		case maths.DirLeft:
			isRight := otherRect.Left() > rect.Right()

			if !vrtWithin || !isCloserX || isRight {
				continue
			}
			closestObstruction = otherRect
		}
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

type SlamboxCollisionInfo struct{}

type SlamboxDestructible struct{}

// type SlamboxChain struct{}?

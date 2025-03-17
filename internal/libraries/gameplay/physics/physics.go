package physics

import (
	"mask_of_the_tomb/internal/libraries/gameplay/grid"
	"mask_of_the_tomb/internal/libraries/maths"
	"math"
)

// TODO: STOP THE DOT IMPORTS
type RectCollider struct {
	maths.Rect
}

func NewRectCollider(r maths.Rect) RectCollider {
	return RectCollider{r}
}

// For now this just is a simple wrapper for the grid struct
type TilemapCollider struct {
	grid.Tilemap[int]
}

// Problem: We need to avoid actually colliding with the rects we are connected to

// TODO: Add the ability to link rects together and then we are gucci
// Round the size up to the nearest size fitting the grid
// Project that sized-up rect
// Finally find how far it actually is from hitting the griddy and move it that far
func (tc *TilemapCollider) ProjectRect(collisionRect *maths.Rect, direction maths.Direction, otherRects []*RectCollider) (maths.Rect, float64) {
	enlargedRect := maths.NewRect(
		math.Floor(collisionRect.Left()/tc.TileSize)*tc.TileSize,
		math.Floor(collisionRect.Top()/tc.TileSize)*tc.TileSize,
		math.Ceil(collisionRect.Width()/tc.TileSize)*tc.TileSize,
		math.Ceil(collisionRect.Height()/tc.TileSize)*tc.TileSize,
	)

	gridX, gridY := tc.WorldPosToGrid(collisionRect.TopLeft())

	x := gridX
	y := gridY

	switch direction {
	case maths.DirUp:
		y = -1000
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j >= 0; j-- {
				if tc.Tiles[j][i] == 1 {
					y = maths.MaxInt(y, j+1)
					break
				}
			}
		}
	case maths.DirDown:
		y = 1000
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j <= len(tc.Tiles); j++ {
				if tc.Tiles[j][i] == 1 {
					y = maths.MinInt(y, j-int(enlargedRect.Height()/tc.TileSize))
					break
				}
			}
		}
	case maths.DirLeft:
		x = -1000
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i >= 0; i-- {
				if tc.Tiles[j][i] == 1 {
					x = maths.MaxInt(x, i+1)
					break
				}
			}
		}
	case maths.DirRight:
		x = 1000
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i <= len(tc.Tiles[0]); i++ {
				if tc.Tiles[j][i] == 1 {
					x = maths.MinInt(x, i-int(enlargedRect.Width()/tc.TileSize))
					break
				}
			}
		}
	}
	newX, newY := tc.GridPosToWorld(x, y)

	// Detect collision against rects
	// If there is a collision, change the position based on movedir
	var moveRect *maths.Rect
	switch direction {
	case maths.DirUp:
		moveRect = maths.NewRect(
			newX,
			newY,
			collisionRect.Width(),
			collisionRect.Bottom()-newY,
		)

		// Find the closes rect we are colliding with
		dist := collisionRect.Top() - newY
		for _, rect := range otherRects {
			if rect.Rect.Overlapping(moveRect) {
				if collisionRect.Top()-rect.Rect.Bottom() < dist {
					dist = collisionRect.Top() - rect.Rect.Bottom()
				}
			}
		}

		newY = collisionRect.Top() - dist

		return *maths.NewRect(
			collisionRect.Left(),
			newY,
			collisionRect.Width(),
			collisionRect.Height(),
		), dist
	case maths.DirDown:
		moveRect = maths.NewRect(
			newX,
			collisionRect.Top(),
			collisionRect.Width(),
			newY+2*collisionRect.Height()-collisionRect.Bottom(),
		)

		// Find the closes rect we are colliding with
		dist := newY - collisionRect.Top()
		for _, rect := range otherRects {

			if rect.Rect.Overlapping(moveRect) {
				if rect.Rect.Top()-collisionRect.Bottom() < dist {
					dist = rect.Rect.Top() - collisionRect.Bottom()
				}
			}
		}
		newY = dist + collisionRect.Top()

		return *maths.NewRect(
			collisionRect.Left(),
			newY+enlargedRect.Height()-collisionRect.Height(),
			collisionRect.Width(),
			collisionRect.Height(),
		), dist
	case maths.DirRight:
		moveRect = maths.NewRect(
			collisionRect.Left(),
			newY,
			newX+collisionRect.Width()-collisionRect.Left(),
			collisionRect.Height(),
		)

		dist := newX - collisionRect.Left()
		for _, rect := range otherRects {
			if rect.Rect.Overlapping(moveRect) {
				if rect.Rect.Left()-collisionRect.Right() < dist {
					dist = rect.Rect.Left() - collisionRect.Right()
				}
			}
		}
		newX = dist + collisionRect.Left()

		newPos := newX + enlargedRect.Width() - collisionRect.Width()
		return *maths.NewRect(
			newPos,
			collisionRect.Top(),
			collisionRect.Width(),
			collisionRect.Height(),
		), math.Abs(dist)
	case maths.DirLeft:
		moveRect = maths.NewRect(
			newX,
			newY,
			collisionRect.Right()-newX,
			collisionRect.Height(),
		)

		dist := collisionRect.Left() - newX
		for _, rect := range otherRects {
			if rect.Rect.Overlapping(moveRect) {
				if collisionRect.Left()-rect.Rect.Right() < dist {
					dist = collisionRect.Left() - rect.Rect.Right()
				}
			}
		}
		newX = collisionRect.Left() - dist

		return *maths.NewRect(
			newX,
			collisionRect.Top(),
			collisionRect.Width(),
			collisionRect.Height(),
		), dist
	}

	// Resize the rect again based on move dir
	return *collisionRect, 0
}

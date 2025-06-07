package physics

import (
	"mask_of_the_tomb/internal/core/maths"
	"math"
)

type TilemapCollider struct {
	tilemap[int]
}

// Round the size up to the nearest size fitting the grid
// Project that sized-up rect
// Finally find how far it actually is from hitting the grid(dy) and move it that far
func (tc *TilemapCollider) ProjectRect(collisionRect *maths.Rect, direction maths.Direction, otherRects []*maths.Rect) (maths.Rect, float64) {
	enlargedRect := maths.NewRect(
		math.Floor(collisionRect.Left()/tc.TileSize)*tc.TileSize,
		math.Floor(collisionRect.Top()/tc.TileSize)*tc.TileSize,
		math.Ceil(collisionRect.Width()/tc.TileSize)*tc.TileSize,
		math.Ceil(collisionRect.Height()/tc.TileSize)*tc.TileSize,
	)

	gridX, gridY := tc.worldPosToGrid(collisionRect.TopLeft())

	x := gridX
	y := gridY

	switch direction {
	case maths.DirUp:
		y = -1000
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j >= 0; j-- {
				if tc.Tiles[j][i] > 0 {
					y = maths.MaxInt(y, j+1)
					break
				}
			}
		}
	case maths.DirDown:
		y = 1000
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j <= len(tc.Tiles); j++ {
				if tc.Tiles[j][i] > 0 {
					y = maths.MinInt(y, j-int(enlargedRect.Height()/tc.TileSize))
					break
				}
			}
		}
	case maths.DirLeft:
		x = -1000
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i >= 0; i-- {
				if tc.Tiles[j][i] > 0 {
					x = maths.MaxInt(x, i+1)
					break
				}
			}
		}
	case maths.DirRight:
		x = 1000
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i <= len(tc.Tiles[0]); i++ {
				if tc.Tiles[j][i] > 0 {
					x = maths.MinInt(x, i-int(enlargedRect.Width()/tc.TileSize))
					break
				}
			}
		}
	}
	newX, newY := tc.gridPosToWorld(x, y)

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
			if rect.Overlapping(moveRect) {
				if collisionRect.Top()-rect.Bottom() < dist {
					dist = collisionRect.Top() - rect.Bottom()
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

			if rect.Overlapping(moveRect) {
				if rect.Top()-collisionRect.Bottom() < dist {
					dist = rect.Top() - collisionRect.Bottom()
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
			if rect.Overlapping(moveRect) {
				if rect.Left()-collisionRect.Right() < dist {
					dist = rect.Left() - collisionRect.Right()
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
			if rect.Overlapping(moveRect) {
				if collisionRect.Left()-rect.Right() < dist {
					dist = collisionRect.Left() - rect.Right()
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

func (tc *TilemapCollider) Raycast(posX, posY float64, direction maths.Direction, otherRects []*maths.Rect) (bool, float64, float64) {
	gridX, gridY := tc.worldPosToGrid(posX, posY)

	x := gridX
	y := gridY

	hit := false
	switch direction {
	case maths.DirUp:
		for i := gridY; i > 0; i-- {
			if tc.Tiles[i][gridX] > 0 {
				y = i + 1
				x -= 1
				hit = true
				break
			}
		}
	case maths.DirDown:
		for i := gridY; i < len(tc.Tiles[0]); i++ {
			if tc.Tiles[i][gridX] > 0 {
				y = i - 1
				x -= 1
				hit = true
				break
			}
		}
	case maths.DirLeft:
		for i := gridX; i > 0; i-- {
			// fmt.Println(tc.Tiles[gridY][i])
			if tc.Tiles[gridY][i] > 0 {
				x = i + 1
				y += 1
				hit = true
				break
			}
		}
	case maths.DirRight:
		for i := gridY; i < len(tc.Tiles); i++ {
			if tc.Tiles[gridY][i] > 0 {
				x = i - 1
				y += 1
				hit = true
				break
			}
		}
	}

	worldX, worldY := tc.gridPosToWorld(x, y)
	return hit, worldX, worldY
}

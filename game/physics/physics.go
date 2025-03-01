package physics

import (
	"fmt"
	. "mask_of_the_tomb/utils"
	"mask_of_the_tomb/utils/rect"
	"math"
	"strconv"
)

// TODO: STOP THE DOT IMPORTS
// TODO: rect should not be in utils
var CollisionMatrix = make([]uint32, 32)

// TODO: ignore?
// It would be nice if this was a little easier to work with, I.E. If you could simply
// treat rects as rectcolliders and vice versa
type RectCollider struct {
	rect.Rect
}

func NewRectCollider(r rect.Rect) RectCollider {
	return RectCollider{r}
}

type TilemapCollider struct {
	Tiles    [][]int
	TileSize float64
	// collisionLayer int // 0 - 32
}

func SetupCollisionMatrix() {
	for i := 0; i < len(CollisionMatrix); i++ {
		CollisionMatrix[i] |= 1 << i
	}
}

// Sets collisions between layer and other to be ignored
func SetCollision(layer, other int) {
	CollisionMatrix[layer] |= 1 << other
	CollisionMatrix[other] |= 1 << layer
}

// Sets collisions between layer and other to not be ignored
func IgnoreCollision(layer, other int) {
	CollisionMatrix[layer] &= ^(1 << other)
	CollisionMatrix[other] &= ^(1 << layer)
}

func CheckCollide(layer, other int) bool {
	return CollisionMatrix[layer]&(1<<other) == 1<<other
}

func PrintCollisionMatrixRow(row int) {
	fmt.Printf("%032s\n", strconv.FormatUint(uint64(CollisionMatrix[row]), 2))
}

// TODO: Need to figure out where worldToGrid should go
func (tc *TilemapCollider) gridToWorld(x, y int) (float64, float64) {
	return F64(x) * tc.TileSize, F64(y) * tc.TileSize
}

func (tc *TilemapCollider) worldToGrid(x, y float64) (int, int) {
	return int(x / tc.TileSize), int(y / tc.TileSize)
}

// TODO: rewrite entire function (it ugly)
// Return a rect representing the new rect for collision object
func (tc *TilemapCollider) ProjectRect(collisionRect *rect.Rect, direction Direction, otherRects []*RectCollider) rect.Rect {
	// Goal: be able to specify which objects should be ignored in this function
	// Loop through collisionobjects
	// If collision layers overlap
	// Detect collision

	// gridX, gridY := l.worldToGrid(rect.TopLeft())
	enlargedRect := rect.NewRect(
		math.Floor(collisionRect.Left()/tc.TileSize)*tc.TileSize,
		math.Floor(collisionRect.Top()/tc.TileSize)*tc.TileSize,
		math.Ceil(collisionRect.Width()/tc.TileSize)*tc.TileSize,
		math.Ceil(collisionRect.Height()/tc.TileSize)*tc.TileSize,
	)

	gridX, gridY := tc.worldToGrid(collisionRect.TopLeft())
	// Round the size up to the nearest size fitting the grid
	// Project that sized-up rect
	// Finally find how far it actually is from hitting the griddy and move it that far
	x := gridX
	y := gridY

	switch direction {
	case DirUp:
		y = -1000
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j >= 0; j-- {
				if tc.Tiles[j][i] == 1 {
					y = MaxInt(y, j+1)
					break
				}
			}
		}
	case DirDown:
		y = 1000
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j <= len(tc.Tiles); j++ {
				if tc.Tiles[j][i] == 1 {
					y = MinInt(y, j-int(enlargedRect.Height()/tc.TileSize))
					break
				}
			}
		}
	case DirLeft:
		x = -1000
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i >= 0; i-- {
				if tc.Tiles[j][i] == 1 {
					x = MaxInt(x, i+1)
					break
				}
			}
		}
	case DirRight:
		x = 1000
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i <= len(tc.Tiles[0]); i++ {
				if tc.Tiles[j][i] == 1 {
					x = MinInt(x, i-int(enlargedRect.Width()/tc.TileSize))
					break
				}
			}
		}
	}
	newX, newY := tc.gridToWorld(x, y)

	// Detect collision against rects
	// If there is a collision, change the position based on movedir
	var moveRect *rect.Rect
	switch direction {
	case DirUp:
		moveRect = rect.NewRect(
			newX,
			newY,
			collisionRect.Width(),
			collisionRect.Bottom()-newY,
		)

		// Find the closes rect we are colliding with
		dist := collisionRect.Top() - newY
		for _, rect := range otherRects {
			// fmt.Println(moveRect)
			// fmt.Println(rect)
			if rect.Rect.Overlapping(moveRect) {
				/*
					fmt.Printf("r.Left: %f\n", rect.Rect.Left())
					fmt.Printf("r.Right: %f\n", rect.Rect.Right())
					fmt.Printf("r.Top: %f\n", rect.Rect.Top())
					fmt.Printf("r.Bottom: %f\n", rect.Rect.Bottom())
					fmt.Println()
					fmt.Printf("other.Left: %f\n", moveRect.Left())
					fmt.Printf("other.Right: %f\n", moveRect.Right())
					fmt.Printf("other.Top: %f\n", moveRect.Top())
					fmt.Printf("other.Bottom: %f\n", moveRect.Bottom())
				*/
				if collisionRect.Top()-rect.Rect.Bottom() < dist {
					dist = collisionRect.Top() - rect.Rect.Bottom()
				}
			}
		}

		newY = collisionRect.Top() - dist
	case DirDown:
		moveRect = rect.NewRect(
			newX,
			collisionRect.Top(),
			collisionRect.Width(),
			newY+collisionRect.Height()-collisionRect.Bottom(),
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
	case DirRight:
		moveRect = rect.NewRect(
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
	case DirLeft:
		moveRect = rect.NewRect(
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
	}

	// Resize the rect again based on move dir
	switch direction {
	case DirUp:
		return *rect.NewRect(
			collisionRect.Left(),
			newY,
			collisionRect.Width(),
			collisionRect.Height(),
		)
	case DirDown:
		return *rect.NewRect(
			collisionRect.Left(),
			newY+enlargedRect.Height()-collisionRect.Height(),
			collisionRect.Width(),
			collisionRect.Height(),
		)
	case DirRight:
		return *rect.NewRect(
			newX+enlargedRect.Width()-collisionRect.Width(),
			collisionRect.Top(),
			collisionRect.Width(),
			collisionRect.Height(),
		)
	case DirLeft:
		return *rect.NewRect(
			newX,
			collisionRect.Top(),
			collisionRect.Width(),
			collisionRect.Height(),
		)
	}

	return *collisionRect
}

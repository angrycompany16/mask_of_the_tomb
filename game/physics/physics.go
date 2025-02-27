package physics

import (
	"fmt"
	"mask_of_the_tomb/utils"
	. "mask_of_the_tomb/utils"
	"mask_of_the_tomb/utils/rect"
	"math"
	"strconv"
)

// TODO: rect should not be in utils
var CollisionMatrix = make([]uint32, 32)

// TODO: ignore?
type RectCollider struct {
	r rect.Rect
	// collisionLayer int // 0 - 32
}

func NewRectCollider(r rect.Rect) RectCollider {
	return RectCollider{r: r}
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

func PrintCollisionMatrix() {
	// for _, i := range CollisionMatrix {
	// fmt.Printf("%b", i)
	// fmt.Println(strconv.FormatInt(int64(i), 2))
	// }
}

// Problem: Need to figure out where worldToGrid should go
func (tc *TilemapCollider) gridToWorld(x, y int) (float64, float64) {
	return F64(x) * tc.TileSize, F64(y) * tc.TileSize
}

func (tc *TilemapCollider) worldToGrid(x, y float64) (int, int) {
	return int(x / tc.TileSize), int(y / tc.TileSize)
}

// Return a rect representing the new rect for collision object
func (tc *TilemapCollider) ProjectRect(collisionRect *rect.Rect, direction utils.Direction, otherRects []RectCollider) rect.Rect {

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
	// Finally find how far it actually is from hitting the grid and move it that far
	x := gridX
	y := gridY
	switch direction {
	case utils.DirUp:
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j >= 0; j-- {
				if tc.Tiles[j][i] == 1 {
					y = j + 1
					break
				}
			}
		}
	case utils.DirDown:
		for i := gridX; i < gridX+int(enlargedRect.Width()/tc.TileSize); i++ {
			for j := gridY; j <= len(tc.Tiles); j++ {
				if tc.Tiles[j][i] == 1 {
					y = j - int(enlargedRect.Height()/tc.TileSize)
					break
				}
			}
		}
	case utils.DirLeft:
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i >= 0; i-- {
				if tc.Tiles[j][i] == 1 {
					x = i + 1
					break
				}
			}
		}
	case utils.DirRight:
		for j := gridY; j < gridY+int(enlargedRect.Height()/tc.TileSize); j++ {
			for i := gridX; i < len(tc.Tiles[0]); i++ {
				if tc.Tiles[j][i] == 1 {
					x = i - int(enlargedRect.Width()/tc.TileSize)
					break
				}
			}
		}
	}
	newX, newY := tc.gridToWorld(x, y)

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

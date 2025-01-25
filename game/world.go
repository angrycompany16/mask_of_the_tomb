package game

import (
	"errors"
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/player"
	. "mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// 0 - air, 1 - wall, 2 - spawn

var (
	gameTileMap = [][]int{
		{1, 1, 1, 1, 1, 1},
		{1, 2, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 0, 0, 0, 1, 1},
		{1, 1, 1, 1, 1, 1},
	}
	tileSprite = files.LazyImage(files.TileSpritePath)
	tileSize   = tileSprite.Bounds().Size()
)

type World struct {
	tiles [][]int
}

func (w *World) Update() {
	// Anything...
}

func (w *World) Draw(surf *ebiten.Image) {
	for i, row := range w.tiles {
		for j, col := range row {
			if col == 1 {
				DrawAt(tileSprite, surf, F64(j*tileSize.X), F64(i*tileSize.Y))
			}
		}
	}
}

func (w *World) GetSpawnPoint() (float64, float64) {
	for i, row := range w.tiles {
		for j, col := range row {
			if col == 2 {
				return F64(j * tileSize.X), F64(i * tileSize.Y)
			}
		}
	}
	return 0, 0
}

func (w *World) getCollision(moveDir player.MoveDirection, x, y float64) (float64, float64, error) {
	gridX, gridY := WorldToGrid(x, y)
	switch moveDir {
	case player.DirUp:
		for i := gridY; i >= 0; i-- {
			if w.tiles[i][gridX] == 1 {
				newX, newY := GridToWorld(gridX, i+1)
				return newX, newY, nil
			}
		}
		return x, y, errors.New("failure: Tried to move up and out of bounds")
	case player.DirDown:
		for i := gridY; i < len(w.tiles); i++ {
			if w.tiles[i][gridX] == 1 {
				newX, newY := GridToWorld(gridX, i-1)
				return newX, newY, nil
			}
		}
		return x, y, errors.New("failure: Tried to move down and out of bounds")
	case player.DirLeft:
		for i := gridX; i >= 0; i-- {
			if w.tiles[gridY][i] == 1 {
				newX, newY := GridToWorld(i+1, gridY)
				return newX, newY, nil
			}
		}
		return x, y, errors.New("failure: Tried to move left and out of bounds")
	case player.DirRight:
		for i := gridX; i < len(w.tiles[0]); i++ {
			if w.tiles[gridY][i] == 1 {
				newX, newY := GridToWorld(i-1, gridY)
				return newX, newY, nil
			}
		}
		return x, y, errors.New("failure: Tried to move right and out of bounds")
	default:
		return x, y, nil
	}
}

func GridToWorld(x, y int) (float64, float64) {
	return F64(x) * F64(tileSize.X), F64(y) * F64(tileSize.Y)
}

func WorldToGrid(x, y float64) (int, int) {
	return int(x / F64(tileSize.X)), int(y / F64(tileSize.Y))

}

func NewWorld() *World {
	return &World{tiles: gameTileMap}
}

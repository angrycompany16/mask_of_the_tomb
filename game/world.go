package game

import (
	"image"
	"mask_of_the_tomb/ebitenLDTK"
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/player"
	. "mask_of_the_tomb/utils"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
// tileSprite = files.LazyImage(files.TileSpritePath)
// tileSize = tileSprite.Bounds().Size()
)

type World struct {
	worldLDTK ebitenLDTK.LDTKWorld
	tilemap   *ebiten.Image
	tileSize  float64
	tiles     [][]int
}

func (w *World) Init() {
	w.worldLDTK = *files.LazyLDTK(files.LDTKMapPath)
	w.tiles = w.worldLDTK.Levels[0].MakeBitmap(&w.worldLDTK.Defs.Layers[0], &w.worldLDTK.Levels[0].LayerInstances[0])

	// One folder back to access LDTK folder
	LDTKpath := path.Clean(path.Join(files.LDTKMapPath, ".."))
	tilemapPath := path.Join(LDTKpath, w.worldLDTK.Defs.Tilesets[0].RelPath)
	w.tilemap = files.LazyImage(tilemapPath)
	w.tileSize = float64(w.worldLDTK.Defs.Tilesets[0].TileGridSize)
}

func (w *World) Update() {
	// Anything...
}

func (w *World) Draw(surf *ebiten.Image) {
	tileSize := w.worldLDTK.Defs.Tilesets[0].TileGridSize
	for _, tile := range w.worldLDTK.Levels[0].LayerInstances[0].GridTiles {
		DrawAt(w.tilemap.SubImage(
			image.Rect(
				tile.Src[0],
				tile.Src[1],
				tile.Src[0]+tileSize,
				tile.Src[1]+tileSize,
			),
		).(*ebiten.Image), surf, F64(tile.Px[0]), F64(tile.Px[1]))
	}
}

func (w *World) GetSpawnPoint() (float64, float64) {
	return 0, 0
}

func (w *World) getCollision(moveDir player.MoveDirection, x, y float64) (float64, float64) {
	gridX, gridY := w.worldToGrid(x, y)
	switch moveDir {
	case player.DirUp:
		for i := gridY; i >= 0; i-- {
			if w.tiles[i][gridX] == 1 {
				newX, newY := w.gridToWorld(gridX, i+1)
				return newX, newY
			}
		}
		return x, y
	case player.DirDown:
		for i := gridY; i < len(w.tiles); i++ {
			if w.tiles[i][gridX] == 1 {
				newX, newY := w.gridToWorld(gridX, i-1)
				return newX, newY
			}
		}
		return x, y
	case player.DirLeft:
		for i := gridX; i >= 0; i-- {
			if w.tiles[gridY][i] == 1 {
				newX, newY := w.gridToWorld(i+1, gridY)
				return newX, newY
			}
		}
		return x, y
	case player.DirRight:
		for i := gridX; i < len(w.tiles[0]); i++ {
			if w.tiles[gridY][i] == 1 {
				newX, newY := w.gridToWorld(i-1, gridY)
				return newX, newY
			}
		}
		return x, y
	default:
		return x, y
	}
}

func (w *World) gridToWorld(x, y int) (float64, float64) {
	return F64(x) * w.tileSize, F64(y) * w.tileSize
}

func (w *World) worldToGrid(x, y float64) (int, int) {
	return int(x / w.tileSize), int(y / w.tileSize)

}

func NewWorld() *World {
	return &World{tiles: [][]int{{0}}}
}

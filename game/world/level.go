package world

import (
	"image"
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game/player"
	. "mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	defs         *ebitenLDTK.Defs
	levelLDTK    *ebitenLDTK.Level
	tiles        [][]int
	tileSize     float64
	collectibles []Collectible
}

func (l *Level) draw(surf *ebiten.Image, camX, camY float64) {

	for i := len(l.levelLDTK.LayerInstances) - 1; i >= 0; i-- {
		layerInstance := l.levelLDTK.LayerInstances[i]

		layer, err := l.defs.GetLayerByUid(layerInstance.LayerDefUid)
		HandleLazy(err)

		// IMPROVEMENT: maybe split into separate functions
		if layer.Type == ebitenLDTK.LayerTypeEntities {
			for _, entityInstance := range layerInstance.EntityInstances {
				entity, err := l.defs.GetEntityByUid(entityInstance.Uid)
				HandleLazy(err)
				if entity.RenderMode == ebitenLDTK.RenderModeTile {
					tileset, err := l.defs.GetTilesetByUid(entityInstance.Tile.TilesetUid)
					HandleLazy(err)

					srcX := entityInstance.Tile.X
					srcY := entityInstance.Tile.Y
					tileSize := tileset.TileGridSize
					DrawAt(tileset.Image.SubImage(
						image.Rect(
							int(srcX),
							int(srcY),
							int(srcX+tileSize),
							int(srcY+tileSize),
						),
					).(*ebiten.Image), surf, entityInstance.Px[0]-camX, entityInstance.Px[1]-camY)

				}
			}
		} else if layer.Type == ebitenLDTK.LayerTypeTiles {
			tileset, err := l.defs.GetTilesetByUid(layer.TilesetUid)
			HandleLazy(err)

			tileSize := tileset.TileGridSize
			for _, tile := range layerInstance.GridTiles {
				scaleX, scaleY := 1.0, 1.0
				switch tile.TileOrientation {
				case ebitenLDTK.OrientationFlipX:
					scaleX = -1
				case ebitenLDTK.OrientationFlipY:
					scaleY = -1
				case ebitenLDTK.OrientationFlipXY:
					scaleX, scaleY = -1, -1
				}
				DrawAtScaled(tileset.Image.SubImage(
					image.Rect(
						int(tile.Src[0]),
						int(tile.Src[1]),
						int(tile.Src[0]+tileSize),
						int(tile.Src[1]+tileSize),
					),
				).(*ebiten.Image), surf, tile.Px[0]-camX, tile.Px[1]-camY, scaleX, scaleY, 0.5, 0.5)
			}
		}
	}
}

func (l *Level) GetSpawnPoint() (float64, float64) {
	for _, layerInstance := range l.levelLDTK.LayerInstances {
		for _, entityInstance := range layerInstance.EntityInstances {
			entity, err := l.defs.GetEntityByUid(entityInstance.Uid)
			HandleLazy(err)
			if entity.Name != spawnPosEntityName {
				continue
			}
			return entityInstance.Px[0], entityInstance.Px[1]
		}
	}
	return 0, 0
}

func (l *Level) GetLevelBounds() (float64, float64) {
	return l.levelLDTK.PxWid, l.levelLDTK.PxHei
}

func (l *Level) GetCollision(moveDir player.MoveDirection, x, y float64) (float64, float64) {
	gridX, gridY := l.worldToGrid(x, y)
	switch moveDir {
	case player.DirUp:
		for i := gridY; i >= 0; i-- {
			if l.tiles[i][gridX] == 1 {
				newX, newY := l.gridToWorld(gridX, i+1)
				return newX, newY
			}
		}
		return x, y
	case player.DirDown:
		for i := gridY; i < len(l.tiles); i++ {
			if l.tiles[i][gridX] == 1 {
				newX, newY := l.gridToWorld(gridX, i-1)
				return newX, newY
			}
		}
		return x, y
	case player.DirLeft:
		for i := gridX; i >= 0; i-- {
			if l.tiles[gridY][i] == 1 {
				newX, newY := l.gridToWorld(i+1, gridY)
				return newX, newY
			}
		}
		return x, y
	case player.DirRight:
		for i := gridX; i < len(l.tiles[0]); i++ {
			if l.tiles[gridY][i] == 1 {
				newX, newY := l.gridToWorld(i-1, gridY)
				return newX, newY
			}
		}
		return x, y
	default:
		return x, y
	}
}

func (l *Level) TryCollectibleOverlap(posX, posY, distX, distY float64) int {
	collected := 0
	for _, layerInstance := range l.levelLDTK.LayerInstances {
		layer, err := l.defs.GetLayerByUid(layerInstance.LayerDefUid)
		HandleLazy(err)

		if layer.Name != collectibleLayerName {
			continue
		}

		for _, entityInstance := range layerInstance.EntityInstances {
			itemX, itemY := l.worldToGrid(entityInstance.Px[0], entityInstance.Px[1])
			playerX, playerY := l.worldToGrid(posX, posY)

			if itemY == playerY {
				if entityInstance.Px[0] > posX && entityInstance.Px[0] <= posX+distX {
					collected++
				}

				if entityInstance.Px[0] < posX && entityInstance.Px[0] >= posX+distX {
					collected++
				}
			}

			if itemX == playerX {
				if entityInstance.Px[1] > posY && entityInstance.Px[1] <= posY+distY {
					collected++
				}

				if entityInstance.Px[1] < posY && entityInstance.Px[1] >= posY+distY {
					collected++
				}
			}

		}
	}
	return collected
}

func (l *Level) TryDoorOverlap(x, y float64) (bool, ebitenLDTK.EntityInstance) {
	for _, layerInstance := range l.levelLDTK.LayerInstances {
		for _, entityInstance := range layerInstance.EntityInstances {
			entity, err := l.defs.GetEntityByUid(entityInstance.Uid)
			HandleLazy(err)
			if entity.Name != doorEntityName {
				continue
			}
			if entityInstance.Px[0] == x && entityInstance.Px[1] == y {
				return true, entityInstance
			}
		}
	}
	return false, ebitenLDTK.EntityInstance{}
}

func (l *Level) gridToWorld(x, y int) (float64, float64) {
	return F64(x) * l.tileSize, F64(y) * l.tileSize
}

func (l *Level) worldToGrid(x, y float64) (int, int) {
	return int(x / l.tileSize), int(y / l.tileSize)
}

func newLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (Level, error) {
	newLevel := Level{}
	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs
	newLevel.tiles = levelLDTK.MakeBitmapFromLayer(defs, playerSpaceLayerName)

	playerspace, err := levelLDTK.GetLayerInstanceByName(playerSpaceLayerName)
	if err != nil {
		newLevel.tileSize = 1
		return newLevel, nil
	}

	layer, err := defs.GetLayerByUid(playerspace.LayerDefUid)
	HandleLazy(err)
	newLevel.tileSize = layer.GridSize

	// TODO: get and set up all collectibles

	return newLevel, nil
}

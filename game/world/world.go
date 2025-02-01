package world

import (
	"fmt"
	"image"
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/game/player"
	. "mask_of_the_tomb/utils"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpaceLayerName   = "Playerspace"
	spawnPosEntityName     = "SpawnPosition"
	doorEntityName         = "Door"
	doorOtherSideFieldName = "OtherSide"
	collectibleLayerName   = "Collectibles"
)

type World struct {
	worldLDTK       ebitenLDTK.World
	activeLevel     *ebitenLDTK.Level
	currentTileSize float64
	tiles           [][]int
}

func (w *World) Init() {
	w.worldLDTK = *files.LazyLDTK(LDTKMapPath)

	w.activeLevel = &w.worldLDTK.Levels[0]
	w.tiles = w.worldLDTK.MakeBitmapFromLayer(w.activeLevel, playerSpaceLayerName)

	// One folder back to access LDTK folder
	LDTKpath := path.Clean(path.Join(LDTKMapPath, ".."))
	for i := 0; i < len(w.worldLDTK.Defs.Tilesets); i++ {
		tileset := &w.worldLDTK.Defs.Tilesets[i]
		tilesetPath := path.Join(LDTKpath, tileset.RelPath)
		tileset.Image = files.LazyImage(tilesetPath)
	}

	changeActiveLevel(w, 0)
}

func (w *World) Update() {
	// Anything...
}

func (w *World) Draw(surf *ebiten.Image, camX, camY float64) {

	for i := len(w.activeLevel.LayerInstances) - 1; i >= 0; i-- {
		layerInstance := w.activeLevel.LayerInstances[i]

		layer, err := w.worldLDTK.GetLayerByUid(layerInstance.LayerDefUid)
		HandleLazy(err)

		// IMPROVEMENT: maybe split into separate functions
		if layer.Type == ebitenLDTK.LayerTypeEntities {
			// for _, entity := range w.worldLDTK.Defs.Entities {
			// 	if entity.RenderMode == ebitenLDTK.LayerTypeTiles {
			// 		tile, err := w.worldLDTK.GetTilesetByUid()
			// 	}
			// }

			for _, entityInstance := range layerInstance.EntityInstances {
				entity, err := w.worldLDTK.GetEntityByUid(entityInstance.Uid)
				HandleLazy(err)
				if entity.RenderMode == ebitenLDTK.RenderModeTile {
					tileset, err := w.worldLDTK.GetTilesetByUid(entityInstance.Tile.TilesetUid)
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
			tileset, err := w.worldLDTK.GetTilesetByUid(layer.TilesetUid)
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

func (w *World) GetSpawnPoint() (float64, float64) {
	for _, layerInstance := range w.activeLevel.LayerInstances {
		for _, entityInstance := range layerInstance.EntityInstances {
			entity, err := w.worldLDTK.GetEntityByUid(entityInstance.Uid)
			HandleLazy(err)
			if entity.Name != spawnPosEntityName {
				continue
			}
			return entityInstance.Px[0], entityInstance.Px[1]
		}
	}
	return 0, 0
}

func (w *World) GetActiveLevelBounds() (float64, float64) {
	return w.activeLevel.PxWid, w.activeLevel.PxHei
}

// TODO: replace player.MoveDirection??
func (w *World) GetCollision(moveDir player.MoveDirection, x, y float64) (float64, float64) {
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

func (w *World) TryCollectibleOverlap(posX, posY, distX, distY float64) int {
	collected := 0
	for _, layerInstance := range w.activeLevel.LayerInstances {
		layer, err := w.worldLDTK.GetLayerByUid(layerInstance.LayerDefUid)
		HandleLazy(err)

		if layer.Name != collectibleLayerName {
			continue
		}

		for _, entityInstance := range layerInstance.EntityInstances {
			itemX, itemY := w.worldToGrid(entityInstance.Px[0], entityInstance.Px[1])
			playerX, playerY := w.worldToGrid(posX, posY)

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

func (w *World) TryDoorOverlap(x, y float64) (bool, ebitenLDTK.EntityInstance) {
	for _, layerInstance := range w.activeLevel.LayerInstances {
		for _, entityInstance := range layerInstance.EntityInstances {
			entity, err := w.worldLDTK.GetEntityByUid(entityInstance.Uid)
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

func (w *World) ExitByDoor(doorEntity ebitenLDTK.EntityInstance) (float64, float64) {
	entity, err := w.worldLDTK.GetEntityByUid(doorEntity.Uid)
	HandleLazy(err)
	if entity.Name != doorEntityName {
		return 0, 0 // Should honestly be error handled
	}
	for _, fieldInstance := range doorEntity.FieldInstances {
		if fieldInstance.Name != doorOtherSideFieldName {
			continue
		}

		nextLevel, err := w.worldLDTK.GetLevelByIid(fieldInstance.EntityRefValue.LevelIid)
		HandleLazy(err)
		// Change level
		changeActiveLevel(w, nextLevel.Uid)

		otherSideEntityRef, err := w.activeLevel.GetEntityInstanceByIid(fieldInstance.EntityRefValue.EntityIid)
		HandleLazy(err)
		return otherSideEntityRef.Px[0], otherSideEntityRef.Px[1]
	}
	return 0, 0
}

func (w *World) gridToWorld(x, y int) (float64, float64) {
	return F64(x) * w.currentTileSize, F64(y) * w.currentTileSize
}

func (w *World) worldToGrid(x, y float64) (int, int) {
	return int(x / w.currentTileSize), int(y / w.currentTileSize)
}

func changeActiveLevel[T string | int](world *World, id T) error {
	var newLevel ebitenLDTK.Level
	var err error

	switch v := any(id).(type) {
	case string:
		newLevel, err = world.worldLDTK.GetLevelByName(v)
		if err != nil {
			fmt.Println("Error when switching levels (id string)")
			return err
		}
	case int:
		newLevel, err = world.worldLDTK.GetLevelByUid(v)
		if err != nil {
			fmt.Println("Error when switching levels (id int)")
			return err
		}
	}

	world.activeLevel = &newLevel

	world.tiles = world.worldLDTK.MakeBitmapFromLayer(world.activeLevel, playerSpaceLayerName)

	playerspace, err := world.activeLevel.GetLayerInstanceByName(playerSpaceLayerName)
	if err != nil {
		world.currentTileSize = 1
		return nil
	}

	layer, err := world.worldLDTK.GetLayerByUid(playerspace.LayerDefUid)
	HandleLazy(err)
	world.currentTileSize = layer.GridSize
	return nil
}

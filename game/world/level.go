package world

import (
	"image"
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/player"
	"mask_of_the_tomb/game/save"
	"mask_of_the_tomb/rect"
	"mask_of_the_tomb/rendering"
	. "mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	defs         *ebitenLDTK.Defs
	levelLDTK    *ebitenLDTK.Level
	tiles        [][]int
	tileSize     float64
	collectibles []Collectible
	hazards      []Hazard
	doors        []Door
}

// This REALLY needs a rewrite
func (l *Level) Draw() {
	camX, camY := camera.GlobalCamera.GetPos()
	for i := len(l.levelLDTK.LayerInstances) - 1; i >= 0; i-- {
		layerInstance := l.levelLDTK.LayerInstances[i]

		layer, err := l.defs.GetLayerByUid(layerInstance.LayerDefUid)
		HandleLazy(err)

		if layer.Name == collectibleLayerName {
			for _, collectible := range l.collectibles {
				collectible.draw(rendering.RenderLayers.Playerspace, camX, camY)
			}
		} else if layer.Type == ebitenLDTK.LayerTypeTiles {
			targetLayer := rendering.RenderLayers.Background
			if layer.Name == playerSpaceLayerName {
				targetLayer = rendering.RenderLayers.Playerspace
			} else if layer.Name == foreGroundLayerName {
				targetLayer = rendering.RenderLayers.Foreground
			}

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
				).(*ebiten.Image),
					targetLayer,
					tile.Px[0]-camX, tile.Px[1]-camY, scaleX, scaleY, 0.5, 0.5)
			}
		}
	}
}

func (l *Level) GetSpawnPoint() (float64, float64) {
	for _, layerInstance := range l.levelLDTK.LayerInstances {
		for _, entityInstance := range layerInstance.EntityInstances {
			entity, err := l.defs.GetEntityByUid(entityInstance.DefUid)
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

func (l *Level) GetCollectibleHit(posX, posY, distX, distY float64) int {
	collected := 0
	collect := func(i int) {
		collected++
		l.collectibles[i].collected = true
		save.GlobalSave.GameData.CollectedEntityUids[l.collectibles[i].iid] = true
	}

	for i := 0; i < len(l.collectibles); i++ {
		collectible := l.collectibles[i]
		if collectible.collected {
			continue
		}
		itemX, itemY := l.worldToGrid(collectible.posX, collectible.posY)
		playerX, playerY := l.worldToGrid(posX, posY)

		if itemY == playerY {
			if (collectible.posX > posX && collectible.posX <= posX+distX) ||
				(collectible.posX < posX && collectible.posX >= posX+distX) {
				collect(i)
			}
		}

		if itemX == playerX {
			if (collectible.posY > posY && collectible.posY <= posY+distY) ||
				(collectible.posY < posY && collectible.posY >= posY+distY) {
				collect(i)
			}
		}
	}
	return collected
}

// Should return (hit bool, levelUid int, posX, posY float64)
func (l *Level) GetDoorHit(playerHitbox *rect.Rect) (hit bool, levelIid string, posX, posY float64) {
	for _, door := range l.doors {
		// fmt.Println("Door")
		if door.hitbox.Overlapping(playerHitbox) {
			hit = true
			levelIid = door.levelIid
			posX, posY = playerHitbox.TopLeft()
		}
	}
	return
}

// func (l *Level) GetDoorHit(x, y float64) (bool, ebitenLDTK.EntityInstance) {
// 	for _, layerInstance := range l.levelLDTK.LayerInstances {
// 		for _, entityInstance := range layerInstance.EntityInstances {
// 			entity, err := l.defs.GetEntityByUid(entityInstance.DefUid)
// 			HandleLazy(err)
// 			if entity.Name != doorEntityName {
// 				continue
// 			}
// 			if entityInstance.Px[0] == x && entityInstance.Px[1] == y {
// 				return true, entityInstance
// 			}
// 		}
// 	}
// 	return false, ebitenLDTK.EntityInstance{}
// }

func (l *Level) GetHazardHit(x, y float64) float64 {
	for _, hazard := range l.hazards {
		if hazard.posX+hazard.width > x && x >= hazard.posX {
			if hazard.posY+hazard.height > y && y >= hazard.posY {
				return hazard.damage
			}
		}
	}
	return 0
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

	for _, layerInstance := range levelLDTK.LayerInstances {
		layer, err := defs.GetLayerByUid(layerInstance.LayerDefUid)
		HandleLazy(err)
		if layer.Name == collectibleLayerName {
			for _, entityInstance := range layerInstance.EntityInstances {
				collected := save.GlobalSave.GameData.CollectedEntityUids[entityInstance.Iid]
				newLevel.collectibles = append(newLevel.collectibles, newCollectible(collected, &entityInstance, defs))
			}
		} else if layer.Name == hazardLayerName {
			for _, entityInstance := range layerInstance.EntityInstances {
				newLevel.hazards = append(newLevel.hazards, newHazard(&entityInstance, defs))
			}
		} else if layer.Name == roomTransitionLayerName {
			for _, entityInstance := range layerInstance.EntityInstances {
				newLevel.doors = append(newLevel.doors, newDoor(&entityInstance))
			}
		}
	}

	return newLevel, nil
}

package world

import (
	"fmt"
	"image"
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/player"
	"mask_of_the_tomb/game/save"
	"mask_of_the_tomb/rendering"
	. "mask_of_the_tomb/utils"
	"mask_of_the_tomb/utils/rect"

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
// Now it REALLY REALLY needs a rewrite
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
		} else if layer.Type == ebitenLDTK.LayerTypeIntGrid {
			targetLayer := rendering.RenderLayers.Playerspace

			tileset, err := l.defs.GetTilesetByUid(layer.TilesetUid)
			if err != nil {
				fmt.Println("Failed to get tileset by uid", layer.TilesetUid)
				return
			}
			// fmt.Println("IntGrid")

			tileSize := tileset.TileGridSize
			for _, tile := range layerInstance.AutoLayerTiles {
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

func (l *Level) GetCollision(moveDir player.MoveDirection, rect *rect.Rect) (posX, posY float64) {
	gridX, gridY := l.worldToGrid(rect.TopLeft())
	x := gridX
	y := gridY
	switch moveDir {
	case player.DirUp:
		for i := gridX; i < gridX+int(rect.Width()/l.tileSize); i++ {
			for j := gridY; j >= 0; j-- {
				if l.tiles[j][i] == 1 {
					y = j + 1
					break
				}
			}
		}
	case player.DirDown:
		for i := gridX; i < gridX+int(rect.Width()/l.tileSize); i++ {
			for j := gridY; j <= len(l.tiles); j++ {
				if l.tiles[j][i] == 1 {
					y = j
					break
				}
			}
		}
	case player.DirLeft:
		for j := gridY; j < gridY+int(rect.Height()/l.tileSize); j++ {
			for i := gridX; i >= 0; i-- {
				if l.tiles[j][i] == 1 {
					x = i + 1
					break
				}
			}
		}
	case player.DirRight:
		for j := gridY; j < gridY+int(rect.Height()/l.tileSize); j++ {
			for i := gridX; i < len(l.tiles[0]); i++ {
				if l.tiles[j][i] == 1 {
					x = i
					break
				}
			}
		}
	}
	posX, posY = l.gridToWorld(x, y)
	return
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

func (l *Level) GetDoorHit(playerHitbox *rect.Rect) (hit bool, levelIid, entityIid string) {
	for _, door := range l.doors {
		if door.hitbox.Overlapping(playerHitbox) {
			hit = true
			levelIid = door.levelIid
			entityIid = door.entityIid
		}
	}
	return
}

// TODO: rewrite with rect
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

func (l *Level) GetEntityInstanceByIid(iid string) (ebitenLDTK.EntityInstance, error) {
	return l.levelLDTK.GetEntityInstanceByIid(iid)
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
				newLevel.doors = append(newLevel.doors, newDoor(&entityInstance, levelLDTK))
			}
		}
	}

	return newLevel, nil
}

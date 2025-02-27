package world

import (
	"image"
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/physics"
	"mask_of_the_tomb/game/save"
	"mask_of_the_tomb/rendering"
	. "mask_of_the_tomb/utils"
	"mask_of_the_tomb/utils/rect"

	"github.com/hajimehoshi/ebiten/v2"
)

// Collision detection system structure
// - Note: Player should collide with static level geometry, movable blocks, maybe
//   interactables and foliage etc...
// - Need some way to specify collision masks or something, a way to determine what collides
//   with what
// - This calls for a new package collision

type Level struct {
	defs            *ebitenLDTK.Defs
	levelLDTK       *ebitenLDTK.Level
	TilemapCollider physics.TilemapCollider
	// tiles           [][]int
	tileSize        float64
	collectibles    []Collectible
	hazards         []Hazard
	doors           []Door
	breakableblocks []BreakableBlock
	slamboxes       []slambox
}

func (l *Level) Update() {
	for _, block := range l.breakableblocks {
		block.Update()
	}
}

// This REALLY needs a rewrite
// Now it REALLY REALLY needs a rewrite
func (l *Level) Draw() {
	for _, block := range l.breakableblocks {
		block.Draw()
	}

	camX, camY := camera.GlobalCamera.GetPos()
	for i := len(l.levelLDTK.Layers) - 1; i >= 0; i-- {
		layer := l.levelLDTK.Layers[i]

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
			for _, tile := range layer.GridTiles {
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
			HandleLazy(err)

			tileSize := tileset.TileGridSize
			for _, tile := range layer.AutoLayerTiles {
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
	for _, layer := range l.levelLDTK.Layers {
		for _, entity := range layer.Entities {
			if entity.Name != spawnPosEntityName {
				continue
			}
			return entity.Px[0], entity.Px[1]
		}
	}
	return 0, 0
}

func (l *Level) GetLevelBounds() (float64, float64) {
	return l.levelLDTK.PxWid, l.levelLDTK.PxHei
}

// Projects a Rect through the map in a certain direction
// func (l *Level) GetCollision(moveDir utils.Direction, rect *rect.Rect) (posX, posY float64) {
// 	gridX, gridY := l.worldToGrid(rect.TopLeft())
// 	x := gridX
// 	y := gridY
// 	switch moveDir {
// 	case utils.DirUp:
// 		for i := gridX; i < gridX+int(rect.Width()/l.tileSize); i++ {
// 			for j := gridY; j >= 0; j-- {
// 				if l.tiles[j][i] == 1 {
// 					y = j + 1
// 					break
// 				}
// 			}
// 		}
// 	case utils.DirDown:
// 		for i := gridX; i < gridX+int(rect.Width()/l.tileSize); i++ {
// 			for j := gridY; j <= len(l.tiles); j++ {
// 				if l.tiles[j][i] == 1 {
// 					y = j
// 					break
// 				}
// 			}
// 		}
// 	case utils.DirLeft:
// 		for j := gridY; j < gridY+int(rect.Height()/l.tileSize); j++ {
// 			for i := gridX; i >= 0; i-- {
// 				if l.tiles[j][i] == 1 {
// 					x = i + 1
// 					break
// 				}
// 			}
// 		}
// 	case utils.DirRight:
// 		for j := gridY; j < gridY+int(rect.Height()/l.tileSize); j++ {
// 			for i := gridX; i < len(l.tiles[0]); i++ {
// 				if l.tiles[j][i] == 1 {
// 					x = i
// 					break
// 				}
// 			}
// 		}
// 	}
// 	posX, posY = l.gridToWorld(x, y)
// 	return
// }

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

func (l *Level) GetEntityByIid(iid string) (ebitenLDTK.Entity, error) {
	return l.levelLDTK.GetEntityByIid(iid)
}

func (l *Level) gridToWorld(x, y int) (float64, float64) {
	return F64(x) * l.tileSize, F64(y) * l.tileSize
}

func (l *Level) worldToGrid(x, y float64) (int, int) {
	return int(x / l.tileSize), int(y / l.tileSize)
}

// TODO: Make layer collision bitmap dynamic
func newLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (Level, error) {
	newLevel := Level{}
	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs

	newLevel.TilemapCollider.Tiles = levelLDTK.MakeBitmapFromLayer(defs, playerSpaceLayerName)
	// newLevel.tiles = levelLDTK.MakeBitmapFromLayer(defs, playerSpaceLayerName)

	playerspace, err := levelLDTK.GetLayerByName(playerSpaceLayerName)
	if err != nil {
		newLevel.tileSize = 1
		newLevel.TilemapCollider.TileSize = 1
		return newLevel, nil
	}

	newLevel.tileSize = float64(playerspace.GridSize)
	newLevel.TilemapCollider.TileSize = float64(playerspace.GridSize)

	for _, layer := range levelLDTK.Layers {
		if layer.Name == collectibleLayerName {
			for _, entity := range layer.Entities {
				collected := save.GlobalSave.GameData.CollectedEntityUids[entity.Iid]
				newLevel.collectibles = append(newLevel.collectibles, newCollectible(collected, &entity, defs))
			}
		} else if layer.Name == hazardLayerName {
			for _, entity := range layer.Entities {
				newLevel.hazards = append(newLevel.hazards, newHazard(&entity))
			}
		} else if layer.Name == roomTransitionLayerName {
			for _, entity := range layer.Entities {
				newLevel.doors = append(newLevel.doors, newDoor(&entity))
			}
		} else if layer.Name == breakableBlockLayerName {
			for _, entity := range layer.Entities {
				// TODO: remove the pointer thingy here
				newLevel.breakableblocks = append(newLevel.breakableblocks, *NewBreakableBlock(&entity))
			}
		} else if layer.Name == slamboxLayerName {
			for _, entity := range layer.Entities {
				newLevel.slamboxes = append(newLevel.slamboxes, newSlambox(&entity))
			}
		}
	}

	return newLevel, nil
}

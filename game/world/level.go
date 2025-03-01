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
	// ActiveColliders []physics.RectCollider
	tileSize        float64
	collectibles    []Collectible
	hazards         []Hazard
	doors           []Door
	breakableblocks []BreakableBlock
	slamboxes       []*slambox
}

func (l *Level) Update() {
	for _, block := range l.breakableblocks {
		block.Update()
	}

	for _, slambox := range l.slamboxes {
		slambox.Update()
	}
}

// This REALLY needs a rewrite
// Now it REALLY REALLY needs a rewrite
func (l *Level) Draw() {
	for _, block := range l.breakableblocks {
		block.Draw()
	}

	for _, box := range l.slamboxes {
		box.Draw()
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

// TODO: this should not be here, very spaghetti
func (l *Level) Without(exSlambox *slambox) []*slambox {
	slamboxes := make([]*slambox, 0)
	for _, _slambox := range l.slamboxes {
		if _slambox != exSlambox {
			slamboxes = append(slamboxes, _slambox)
		}
	}
	return slamboxes
}

// For now we assume that we will only ever be slamming one box at a time, though
// this may change later
func (l *Level) GetSlamboxHit(playerCollider *rect.Rect, dir Direction) *slambox {
	extendedRect := playerCollider.Extended(dir, 1)
	for _, slambox := range l.slamboxes {
		if extendedRect.Overlapping(&slambox.collider.Rect) {
			return slambox
		}
	}
	return nil
	// Increase the size of the collider by 1 in the moveDirection
	// If there is an overlap with a slambox, return true
	// Else, return false
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

func (l *Level) GetSlamboxColliders() []*physics.RectCollider {
	colliders := make([]*physics.RectCollider, 0)
	for _, box := range l.slamboxes {
		colliders = append(colliders, &box.collider)
	}
	return colliders
}

func newLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (Level, error) {
	newLevel := Level{}
	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs

	newLevel.TilemapCollider.Tiles = levelLDTK.MakeBitmapFromLayer(defs, playerSpaceLayerName)

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
				newSlambox := newSlambox(&entity)
				newLevel.slamboxes = append(newLevel.slamboxes, newSlambox)
			}
		}
	}

	return newLevel, nil
}

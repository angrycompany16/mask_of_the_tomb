package level

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/engine/advertisers"
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/entities/camera/pubcamera"
	pubgame "mask_of_the_tomb/internal/entities/game/pub"
	"mask_of_the_tomb/internal/libraries/assets/ebitenrenderutil"
	"mask_of_the_tomb/internal/libraries/assets/ldtknames"
	"mask_of_the_tomb/internal/libraries/errs"
	"mask_of_the_tomb/internal/libraries/gameplay/physics"
	"mask_of_the_tomb/internal/libraries/rendering"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

var layerMap = map[string]*ebiten.Image{
	"Foreground":       rendering.RenderLayers.Foreground,
	"PlayerspaceAlt":   rendering.RenderLayers.Playerspace,
	"Playerspace":      rendering.RenderLayers.Playerspace,
	"MidgroundSprites": rendering.RenderLayers.Midground,
	"BackgroundTiles":  rendering.RenderLayers.Background,
}

type InitLevelInfo struct {
	LevelWidth, LevelHeight float64
	SpawnX, SpawnY          float64
}

type Level struct {
	*entities.Entity
	defs            *ebitenLDTK.Defs
	levelLDTK       *ebitenLDTK.Level
	TilemapCollider physics.TilemapCollider
	tileSize        float64
	// hazards         []hazard
	// doors           []door
	// Slamboxes       []*Slambox
}

func (l *Level) Update() {
	// for _, slambox := range l.Slamboxes {
	// 	slambox.Update()
	// }
}

// TODO: Fill the background with backgroundColor from LDTK
func (l *Level) Draw() {
	gameAdv := advertisers.Get(pubgame.GameEntityName)
	gameState := gameAdv.Read().(pubgame.GameAdvertiser)
	if gameState.State == pubgame.StateMainMenu {
		return
	}

	camAdv := advertisers.Get(pubcamera.CameraEntityName)
	camPos := camAdv.Read().(pubcamera.CameraAdvertiser)
	camX, camY := camPos.PosX, camPos.PosY

	// for _, box := range l.Slamboxes {
	// 	box.Draw()
	// }

	// NOTE: we *need* to loop in reverse
	for i := len(l.levelLDTK.Layers) - 1; i >= 0; i-- {
		layer := l.levelLDTK.Layers[i]

		targetRenderLayer, ok := layerMap[layer.Name]
		if !ok && (layer.Type == ebitenLDTK.LayerTypeIntGrid || layer.Type == ebitenLDTK.LayerTypeTiles) {
			fmt.Printf("Layer with name %s does not have a rendering layer\n", layer.Name)
			continue
		}

		if layer.Type == ebitenLDTK.LayerTypeTiles {
			tileset := errs.Must(l.defs.GetTilesetByUid(layer.TilesetUid))
			drawTile(layer.GridTiles, &tileset, targetRenderLayer, tileset.TileGridSize, camX, camY)
		} else if layer.Type == ebitenLDTK.LayerTypeIntGrid {
			tileset := errs.Must(l.defs.GetTilesetByUid(layer.TilesetUid))
			drawTile(layer.AutoLayerTiles, &tileset, targetRenderLayer, tileset.TileGridSize, camX, camY)
		}
	}
}

func drawTile(
	tiles []ebitenLDTK.Tile,
	tileset *ebitenLDTK.Tileset,
	targetLayer *ebiten.Image,
	tileSize, camX, camY float64,
) {
	for _, tile := range tiles {
		scaleX, scaleY := 1.0, 1.0
		switch tile.TileOrientation {
		case ebitenLDTK.OrientationFlipX:
			scaleX = -1
		case ebitenLDTK.OrientationFlipY:
			scaleY = -1
		case ebitenLDTK.OrientationFlipXY:
			scaleX, scaleY = -1, -1
		}
		ebitenrenderutil.DrawAtScaled(tileset.Image.SubImage(
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

func (l *Level) GetSpawnPoint() (float64, float64) {
	for _, layer := range l.levelLDTK.Layers {
		for _, entity := range layer.Entities {
			if entity.Name != ldtknames.SpawnPosEntityName {
				continue
			}
			return entity.Px[0], entity.Px[1]
		}
	}
	return 0, 0
}

// func (l *Level) GetDoorHit(playerHitbox *maths.Rect) (hit bool, levelIid, entityIid string) {
// 	for _, door := range l.doors {
// 		if door.hitbox.Overlapping(playerHitbox) {
// 			hit = true
// 			levelIid = door.levelIid
// 			entityIid = door.entityIid
// 		}
// 	}
// 	return
// }

// func (l *Level) GetHazardHit(playerHitbox *maths.Rect) float64 {
// 	for _, hazard := range l.hazards {
// 		if hazard.hitbox.Overlapping(playerHitbox) {
// 			return hazard.damage
// 		}
// 	}
// 	return 0
// }

// Get all the rect colliders that are not connected to slambox
// func (l *Level) DisconnectedColliders(slambox *Slambox) []*physics.RectCollider {
// 	// I love writing unreadable code
// 	return arrays.MapSlice(
// 		arrays.Filter(
// 			l.Slamboxes, func(s *Slambox) bool { return !slices.Contains(slambox.ConnectedBoxes, s) && s != slambox },
// 		),
// 		func(s *Slambox) *physics.RectCollider { return &s.Collider },
// 	)

// }

// func (l *Level) GetSlamboxColliders() []*physics.RectCollider {
// 	return arrays.MapSlice(l.Slamboxes, func(s *Slambox) *physics.RectCollider { return &s.Collider })
// }

// For now we assume that we will only ever be slamming one box at a time, though
// this may change later
// func (l *Level) GetSlamboxHit(playerCollider *maths.Rect, dir maths.Direction) *Slambox {
// 	extendedRect := playerCollider.Extended(dir, 1)
// 	for _, slambox := range l.Slamboxes {
// 		if extendedRect.Overlapping(&slambox.Collider.Rect) {
// 			return slambox
// 		}
// 	}
// 	return nil
// }

func (l *Level) GetEntityByIid(iid string) (ebitenLDTK.Entity, error) {
	return l.levelLDTK.GetEntityByIid(iid)
}

// TODO: Refactor because this is very big
func NewLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (*Level, InitLevelInfo, error) {
	newLevel := Level{}
	newLevel.Entity = entities.RegisterEntity(&newLevel, "Level")

	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs

	newLevel.TilemapCollider.Tiles = levelLDTK.MakeBitmapFromLayer(defs, ldtknames.PlayerSpaceLayerName)

	playerspace, err := levelLDTK.GetLayerByName(ldtknames.PlayerSpaceLayerName)
	if err != nil {
		newLevel.tileSize = 1
		newLevel.TilemapCollider.TileSize = 1
		return &newLevel, InitLevelInfo{}, nil
	}

	newLevel.tileSize = float64(playerspace.GridSize)
	newLevel.TilemapCollider.TileSize = float64(playerspace.GridSize)

	// for _, layer := range levelLDTK.Layers {
	// 	for _, entity := range layer.Entities {
	// 		switch entity.Name {
	// 		case hazardEntityName:
	// 			newLevel.hazards = append(newLevel.hazards, newHazard(&entity))
	// 		case doorEntityName:
	// 			newLevel.doors = append(newLevel.doors, newDoor(&entity))
	// 		case slamboxEntityName:
	// 			newLevel.Slamboxes = append(newLevel.Slamboxes, newSlambox(&entity))
	// 		}
	// 	}
	// }

	// NOTE: We need to loop twice to ensure that all slamboxes have been added
	// before we link them together
	// for _, slambox := range newLevel.Slamboxes {
	// 	for _, otherSlambox := range newLevel.Slamboxes {
	// 		if slices.Contains(slambox.otherLinkIDs, otherSlambox.LinkID) {
	// 			slambox.ConnectedBoxes = append(slambox.ConnectedBoxes, otherSlambox)
	// 		}
	// 	}
	// }

	spawnX, spawnY := newLevel.GetSpawnPoint()
	newLevelInfo := InitLevelInfo{
		SpawnX:      spawnX,
		SpawnY:      spawnY,
		LevelWidth:  newLevel.levelLDTK.PxWid,
		LevelHeight: newLevel.levelLDTK.PxHei,
	}

	return &newLevel, newLevelInfo, nil
}

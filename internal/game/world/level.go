package world

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/arrays"
	ebitenrenderutil "mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/camera"
	"mask_of_the_tomb/internal/game/physics"
	"mask_of_the_tomb/internal/game/rendering"
	"mask_of_the_tomb/internal/maths"

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

type Level struct {
	defs            *ebitenLDTK.Defs
	levelLDTK       *ebitenLDTK.Level
	TilemapCollider physics.TilemapCollider
	tileSize        float64
	hazards         []hazard
	doors           []door
	Slamboxes       []*Slambox
}

func (l *Level) Update() {
	for _, slambox := range l.Slamboxes {
		slambox.Update()
	}
}

func (l *Level) Draw() {
	for _, box := range l.Slamboxes {
		box.Draw()
	}

	camX, camY := camera.GlobalCamera.GetPos()
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

func (l *Level) GetDoorHit(playerHitbox *maths.Rect) (hit bool, levelIid, entityIid string) {
	for _, door := range l.doors {
		if door.hitbox.Overlapping(playerHitbox) {
			hit = true
			levelIid = door.levelIid
			entityIid = door.entityIid
		}
	}
	return
}

func (l *Level) GetHazardHit(playerHitbox *maths.Rect) float64 {
	for _, hazard := range l.hazards {
		if hazard.hitbox.Overlapping(playerHitbox) {
			return hazard.damage
		}
	}
	return 0
}

// TODO: this should not be here, very spaghetti
func (l *Level) Without(exSlambox *Slambox) []*Slambox {
	slamboxes := make([]*Slambox, 0)
	for _, _slambox := range l.Slamboxes {
		if _slambox != exSlambox {
			slamboxes = append(slamboxes, _slambox)
		}
	}
	return slamboxes
}

func (l *Level) GetSlamboxColliders() []*physics.RectCollider {
	return arrays.MapSlice(l.Slamboxes, func(s *Slambox) *physics.RectCollider { return &s.Collider })
}

// For now we assume that we will only ever be slamming one box at a time, though
// this may change later
func (l *Level) GetSlamboxHit(playerCollider *maths.Rect, dir maths.Direction) *Slambox {
	extendedRect := playerCollider.Extended(dir, 1)
	for _, slambox := range l.Slamboxes {
		if extendedRect.Overlapping(&slambox.Collider.Rect) {
			return slambox
		}
	}
	return nil
}

func (l *Level) GetEntityByIid(iid string) (ebitenLDTK.Entity, error) {
	return l.levelLDTK.GetEntityByIid(iid)
}

func newLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (Level, error) {
	newLevel := Level{}
	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs

	// fmt.Println(playerSpaceLayerName)
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
		if layer.Name == hazardLayerName {
			for _, entity := range layer.Entities {
				newLevel.hazards = append(newLevel.hazards, newHazard(&entity))
			}
		} else if layer.Name == roomTransitionLayerName {
			for _, entity := range layer.Entities {
				newLevel.doors = append(newLevel.doors, newDoor(&entity))
			}
		} else if layer.Name == slamboxLayerName {
			for _, entity := range layer.Entities {
				newSlambox := newSlambox(&entity)
				newLevel.Slamboxes = append(newLevel.Slamboxes, newSlambox)
			}
		}
	}

	return newLevel, nil
}

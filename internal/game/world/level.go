package world

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/arrays"
	ebitenrenderutil "mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	"mask_of_the_tomb/internal/game/physics"
	"mask_of_the_tomb/internal/game/world/entities"
	"mask_of_the_tomb/internal/maths"
	"slices"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpaceLayerName = "Playerspace"
	spawnPosEntityName   = "SpawnPosition"
	doorEntityName       = "Door"
	spawnPointEntityName = "SpawnPoint"
	slamboxEntityName    = "Slambox"
	SpikeIntGridName     = "Spikes"
)

var layerMap = map[string]*ebiten.Image{
	"Foreground":       rendering.RenderLayers.Foreground,
	"PlayerspaceAlt":   rendering.RenderLayers.Playerspace,
	"Playerspace":      rendering.RenderLayers.Playerspace,
	"MidgroundSprites": rendering.RenderLayers.Midground,
	"BackgroundTiles":  rendering.RenderLayers.Background,
}

type Level struct {
	name            string
	defs            *ebitenLDTK.Defs
	levelLDTK       *ebitenLDTK.Level
	TilemapCollider physics.TilemapCollider
	tileSize        float64
	hazards         []entities.Hazard
	doors           []entities.Door
	Slamboxes       []*entities.Slambox
}

func (l *Level) Update() {
	for _, slambox := range l.Slamboxes {
		slambox.Update()
	}
}

// TODO: Fill the background with backgroundColor from LDTK
func (l *Level) Draw() {
	for _, box := range l.Slamboxes {
		box.Draw()
	}

	camX, camY := camera.GetPos()
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
		if door.Hitbox.Overlapping(playerHitbox) {
			hit = true
			levelIid = door.LevelIid
			entityIid = door.EntityIid
		}
	}
	return
}

func (l *Level) GetHazardHit(playerHitbox *maths.Rect) bool {
	for _, hazard := range l.hazards {
		if hazard.Hitbox.Overlapping(playerHitbox) {
			return true
		}
	}
	return false
}

// Get all the rect colliders that are not connected to slambox
func (l *Level) DisconnectedColliders(slambox *entities.Slambox) []*physics.RectCollider {
	// I love writing unreadable code
	return arrays.MapSlice(
		arrays.Filter(
			l.Slamboxes, func(s *entities.Slambox) bool { return !slices.Contains(slambox.ConnectedBoxes, s) && s != slambox },
		),
		func(s *entities.Slambox) *physics.RectCollider { return &s.Collider },
	)

}

func (l *Level) GetSlamboxColliders() []*physics.RectCollider {
	return arrays.MapSlice(l.Slamboxes, func(s *entities.Slambox) *physics.RectCollider { return &s.Collider })
}

// For now we assume that we will only ever be slamming one box at a time, though
// this may change later
func (l *Level) GetSlamboxHit(playerCollider *maths.Rect, dir maths.Direction) *entities.Slambox {
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

func (l *Level) GetBiome() string {
	// fmt.Println(l.levelLDTK)
	field := errs.Must(l.levelLDTK.GetFieldByName("Biome"))
	return field.Biome
}

// TODO: Refactor because this is very big
func newLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (*Level, error) {
	newLevel := &Level{}
	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs
	newLevel.name = levelLDTK.Name

	playerspace, err := levelLDTK.GetLayerByName(playerSpaceLayerName)
	if err != nil {
		newLevel.tileSize = 1
		newLevel.TilemapCollider.TileSize = 1
		return newLevel, nil
	}

	var spikeIntGridID int
	for _, layerDef := range defs.LayerDefs {
		if layerDef.Name == playerSpaceLayerName {
			spikeIntGridID = layerDef.GetIntGridValue(SpikeIntGridName)
		}
	}

	// TODO: Why the hell is this running twice?
	// fmt.Println("hello") // use this to see the point
	intGridCSV := playerspace.ExtractLayerCSV([]int{spikeIntGridID})
	newLevel.TilemapCollider.Tiles = intGridCSV

	newLevel.tileSize = float64(playerspace.GridSize)
	newLevel.TilemapCollider.TileSize = float64(playerspace.GridSize)

	for _, layer := range levelLDTK.Layers {
		for _, entity := range layer.Entities {
			switch entity.Name {
			case entities.HazardEntityName:
				newLevel.hazards = append(newLevel.hazards, entities.NewHazard(&entity))
			case doorEntityName:
				newLevel.doors = append(newLevel.doors, entities.NewDoor(&entity))
			case slamboxEntityName:
				newLevel.Slamboxes = append(newLevel.Slamboxes, entities.NewSlambox(&entity))
			}
		}
	}

	// NOTE: We need to loop twice to ensure that all slamboxes have been added
	// before we link them together
	for _, slambox := range newLevel.Slamboxes {
		for _, otherSlambox := range newLevel.Slamboxes {
			if slices.Contains(slambox.OtherLinkIDs, otherSlambox.LinkID) {
				slambox.ConnectedBoxes = append(slambox.ConnectedBoxes, otherSlambox)
			}
		}
	}

	return newLevel, nil
}

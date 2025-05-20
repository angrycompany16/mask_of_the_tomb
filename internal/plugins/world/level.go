package world

import (
	"fmt"
	"image"
	"image/color"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/arrays"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/libraries/colors"
	"mask_of_the_tomb/internal/libraries/particles"
	"mask_of_the_tomb/internal/libraries/physics"
	"mask_of_the_tomb/internal/libraries/rendering"
	"path/filepath"
	"slices"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	entityLayerName        = "Entities"
	playerSpaceLayerName   = "Playerspace"
	spawnPosEntityName     = "DefaultSpawnPos"
	doorEntityName         = "Door"
	spawnPointEntityName   = "SpawnPoint"
	slamboxEntityName      = "Slambox"
	SpikeIntGridName       = "Spikes"
	gameEntryPosEntityName = "GameEntryPos"
)

var (
	lerpBlend = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorSourceAlpha,
		BlendFactorSourceAlpha:      ebiten.BlendFactorZero,
		BlendFactorDestinationRGB:   ebiten.BlendFactorOneMinusSourceAlpha,
		BlendFactorDestinationAlpha: ebiten.BlendFactorZero,
		BlendOperationRGB:           ebiten.BlendOperationAdd,
		BlendOperationAlpha:         ebiten.BlendOperationAdd,
	}
	particleSysPath              = filepath.Join("assets", "particlesystems", "environment", "basement.yaml")
	playerSpaceNormalTilemapPath = filepath.Join("assets", "sprites", "environment", "tilemaps", "export", "playerspace_tilemap_normal.png")
	// midgroundNormalTilemapPath   = filepath.Join("assets", "sprites", "environment", "tilemaps", "export", "midground_tilemap_normal.png")
	// TODO: *very* temporary solution
	playerspaceNormalTilemap = errs.MustNewImageFromFile(playerSpaceNormalTilemapPath)
	// midgroundNormalTilemap   = errs.MustNewImageFromFile(midgroundNormalTilemapPath)
)

type Level struct {
	name                     string
	defs                     *ebitenLDTK.Defs
	levelLDTK                *ebitenLDTK.Level
	TilemapCollider          physics.TilemapCollider
	tileSize                 float64
	bgColor                  color.Color
	playerspaceNormalTilemap *ebiten.Image
	midgroundNormalTilemap   *ebiten.Image

	// Probably convert into struct in rendering.go
	colorLayers  rendering.LayerList
	normalLayers rendering.LayerList
	// foreground  *ebiten.Image
	// playerspace *ebiten.Image
	// midground   *ebiten.Image
	// background  *ebiten.Image

	fogShader      *ebiten.Shader
	vignetteShader *ebiten.Shader
	// lightsAdditiveShader    *ebiten.Shader
	// lightsSubtractiveShader *ebiten.Shader
	lightsPixelShader *ebiten.Shader
	pixelLightShader  *ebiten.Shader
	resetX, resetY    float64
	ambientParticles  *particles.ParticleSystem
	hazards           []Hazard
	doors             []Door
	slamboxes         []*Slambox
}

// ------ CONSTRUCTOR ------
// TODO: Refactor because this is very big
func newLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (*Level, error) {
	newLevel := &Level{}
	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs
	newLevel.name = levelLDTK.Name
	newLevel.bgColor = colors.HexToRGB(levelLDTK.BgColorHex)
	newLevel.fogShader = errs.Must(ebiten.NewShader(assets.Fog_kage))
	newLevel.vignetteShader = errs.Must(ebiten.NewShader(assets.Vignette_kage))
	newLevel.lightsPixelShader = errs.Must(ebiten.NewShader(assets.Pixel_lights_kage))
	// TODO: set particle system bounds based on level size
	newLevel.ambientParticles = errs.Must(particles.FromFile(particleSysPath, rendering.ScreenLayers.Foreground))

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
			case HazardEntityName:
				newLevel.hazards = append(newLevel.hazards, NewHazard(&entity))
			case doorEntityName:
				newLevel.doors = append(newLevel.doors, NewDoor(&entity))
			case slamboxEntityName:
				newLevel.slamboxes = append(newLevel.slamboxes, NewSlambox(&entity))
			}
		}
	}

	// NOTE: We need to loop twice to ensure that all slamboxes have been added
	// before we link them together
	for _, slambox := range newLevel.slamboxes {
		for _, otherSlambox := range newLevel.slamboxes {
			if slices.Contains(slambox.OtherLinkIDs, otherSlambox.LinkID) {
				// TODO: Change to LinkBox() method
				slambox.ConnectedBoxes = append(slambox.ConnectedBoxes, otherSlambox)
			}
		}
		slambox.CreateSprite()
	}

	newLevel.colorLayers = rendering.NewLayerList(int(levelLDTK.PxWid), int(levelLDTK.PxHei))
	newLevel.normalLayers = rendering.NewLayerList(int(levelLDTK.PxWid), int(levelLDTK.PxHei))

	layerMap := map[string]*ebiten.Image{
		"Foreground":       newLevel.colorLayers.Foreground,
		"PlayerspaceAlt":   newLevel.colorLayers.Playerspace,
		"Playerspace":      newLevel.colorLayers.Playerspace,
		"Props":            newLevel.colorLayers.Midground,
		"MidgroundSprites": newLevel.colorLayers.Midground,
		"BackgroundTiles":  newLevel.colorLayers.Background,
	}

	for i := len(newLevel.levelLDTK.Layers) - 1; i >= 0; i-- {
		layer := newLevel.levelLDTK.Layers[i]

		targetRenderLayer, ok := layerMap[layer.Name]
		if !ok && (layer.Type == ebitenLDTK.LayerTypeIntGrid || layer.Type == ebitenLDTK.LayerTypeTiles) {
			fmt.Printf("Layer with name %s does not have a rendering layer\n", layer.Name)
			continue
		}

		if layer.Type == ebitenLDTK.LayerTypeEntities {
			continue
		}

		var tiles []ebitenLDTK.Tile
		if layer.Type == ebitenLDTK.LayerTypeTiles {
			tiles = layer.GridTiles
		} else if layer.Type == ebitenLDTK.LayerTypeIntGrid {
			tiles = layer.AutoLayerTiles
		}

		// for _, tileset := range defs.Tilesets {
		// 	fmt.Println(tileset.Uid)
		// }
		// // fmt.Println(newLevel.defs.Tilesets)
		// fmt.Println(layer.TilesetUid)
		// fmt.Println(layer.Type)
		// fmt.Println("Done")
		tileset := errs.Must(newLevel.defs.GetTilesetByUid(layer.TilesetUid))
		tilesize := tileset.TileGridSize
		tilesetImage := tileset.Image

		if targetRenderLayer == newLevel.colorLayers.Midground {
			// tilesetImage = midgroundNormalTilemap
			// drawTiles(tiles, tilesetImage, newLevel.normalLayers.Midground, tilesize)
		} else if targetRenderLayer == newLevel.colorLayers.Playerspace {
			tilesetImage = playerspaceNormalTilemap
			drawTiles(tiles, tilesetImage, newLevel.normalLayers.Playerspace, tilesize)
		}

		drawTiles(tiles, tilesetImage, targetRenderLayer, tilesize)
	}

	return newLevel, nil
}

// ------ ENTITY ------
func (l *Level) Update() {
	for _, slambox := range l.slamboxes {
		slambox.Update()
	}
	l.ambientParticles.Update()
}

// Draw the level only once
// Note: Only a few things need to be rendered multiple times:
// Dynamic objects
// Shaders
// Basically we can draw the level only once and then draw that texture onto
// the screen instead of drawing each tile every frame
// Optimization yeah
func (l *Level) Draw(playerX, playerY, camX, camY, time float64) {
	for _, box := range l.slamboxes {
		box.Draw(camX, camY)
	}

	rendering.ScreenLayers.Background2.Fill(l.bgColor)

	// Render fog layer
	shaderOp := ebiten.DrawRectShaderOptions{}
	shaderOp.Uniforms = map[string]any{
		"Time":       time / 5,
		"Amplitude":  1.0,
		"Frequency":  0.025,
		"Strength":   0.7,
		"Threshold":  0.4,
		"Color":      [4]float64{37.0 / 255, 49.0 / 255, 94.0 / 255, 1.0},
		"Center":     [2]float64{0.5, 0.5},
		"Resolution": [2]float64{rendering.GameWidth, rendering.GameHeight},
	}
	shaderOp.Blend = ebiten.BlendSourceOver
	// TODO: Move fog with camera position
	rendering.ScreenLayers.Background2.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.fogShader, &shaderOp)

	shaderOp.Uniforms = map[string]any{
		"PlayerPos":  [2]float64{playerX, playerY},
		"Resolution": [2]float64{rendering.GameWidth, rendering.GameHeight},
	}
	// shaderOp.Blend = lerpBlend
	// rendering.RenderLayers.Foreground.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.lightsAdditiveShader, &shaderOp)

	l.ambientParticles.Draw(camX, camY)

	// l.layers.DrawOnto(&rendering.ScreenLayers, -camX, -camY)
	// rendering.ScreenLayers.Background.

	// shaderOp.Blend = ebiten.BlendSourceOver
	// rendering.RenderLayers.Foreground.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.lightsSubtractiveShader, &shaderOp)

	shaderOp.Blend = ebiten.BlendSourceOver
	rendering.ScreenLayers.Foreground.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.vignetteShader, &shaderOp)

	// Draw lighting on both midground and playerspace
	// This creates problems with camera shake
	trueCamX, trueCamY := rendering.GetStablePos()

	shaderOp.Images = [4]*ebiten.Image{
		l.colorLayers.Playerspace.SubImage(image.Rect(int(trueCamX), int(trueCamY), int(trueCamX+rendering.GameWidth), int(trueCamY+rendering.GameHeight))).(*ebiten.Image),
		nil,
		nil,
		nil,
	}

	playerXrel := playerX - camX
	playerYrel := playerY - camY

	shakeX, shakeY := rendering.GetShake()
	shaderOp.Uniforms = map[string]any{
		"CamShake":     [2]float64{shakeX, shakeY},
		"AmbientLight": [4]float64{0.0, 0.0, 0.0, 1.0},
		"LightPos0":    [3]float64{playerXrel, playerYrel, 0.0},
		"LightColor0":  [3]float64{1.0, 0.0, 0.0},
	}

	// rendering.RenderLayers.Midground.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.lightsPixelShader, &shaderOp)
	// rendering.ScreenLayers.Playerspace.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.lightsPixelShader, &shaderOp)

	l.normalLayers.DrawOnto(&rendering.ScreenLayers, -camX, -camY)

	// pXrel := playerX - camX
	// pYrel := playerY - camY
	// shaderOp.Uniforms = map[string]any{
	// 	"PlayerPos":  [2]float64{pXrel, pYrel},
	// 	"Resolution": [2]float64{rendering.GameWidth, rendering.GameHeight},
	// }
	// shaderOp.Blend = lerpBlend
	// rendering.RenderLayers.Foreground.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.lightsAdditiveShader, &shaderOp)

	// l.ambientParticles.Draw(camX, camY)

	// shaderOp.Blend = ebiten.BlendSourceOver
	// rendering.RenderLayers.Foreground.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.lightsSubtractiveShader, &shaderOp)

	// shaderOp.Blend = ebiten.BlendSourceOver
	// rendering.RenderLayers.Foreground.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.vignetteShader, &shaderOp)

	// rendering.RenderLayers.Playerspace.DrawRectShader(rendering.GameWidth, rendering.GameHeight, l.vignetteShader, &shaderOp)
}

// ------ GETTERS ------
func (l *Level) CheckDoorOverlap(playerHitbox *maths.Rect) (hit bool, levelIid, entityIid string) {
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
func (l *Level) GetDisconnectedColliders(slambox *Slambox) []*physics.RectCollider {
	// I love writing unreadable code
	return arrays.MapSlice(
		arrays.Filter(
			l.slamboxes, func(s *Slambox) bool { return !slices.Contains(slambox.ConnectedBoxes, s) && s != slambox },
		),
		func(s *Slambox) *physics.RectCollider { return &s.Collider },
	)

}

func (l *Level) GetSlamboxColliders() []*physics.RectCollider {
	return arrays.MapSlice(l.slamboxes, func(s *Slambox) *physics.RectCollider { return &s.Collider })
}

func (l *Level) GetSlamboxPositions() []SlamboxPosition {
	return arrays.MapSlice(l.slamboxes, func(s *Slambox) SlamboxPosition {
		return SlamboxPosition{X: s.Collider.Left(), Y: s.Collider.Top()}
	})
}

// For now we assume that we will only ever be slamming one box at a time, though
// this may change later
func (l *Level) GetSlamboxHit(playerCollider *maths.Rect, dir maths.Direction) *Slambox {
	extendedRect := playerCollider.Extended(dir, 1)
	for _, slambox := range l.slamboxes {
		if extendedRect.Overlapping(&slambox.Collider.Rect) {
			return slambox
		}
	}
	return nil
}

func (l *Level) GetBiome() string {
	field, err := l.levelLDTK.GetFieldByName("Biome")
	if err != nil {
		return ""
	}
	return field.Biome
}

func (l *Level) GetBounds() (float64, float64) {
	return l.levelLDTK.PxWid, l.levelLDTK.PxHei
}

func (l *Level) GetDefaultSpawnPoint() (float64, float64) {
	entityLayer := errs.Must(l.levelLDTK.GetLayerByName(entityLayerName))
	for _, entity := range entityLayer.Entities {
		if entity.Name != spawnPosEntityName {
			continue
		}
		return entity.Px[0], entity.Px[1]
	}
	for _, entity := range entityLayer.Entities {
		if entity.Name != doorEntityName {
			continue
		}
		return entity.Px[0], entity.Px[1]
	}
	return 0, 0
}

func (l *Level) GetResetPoint() (float64, float64) {
	return l.resetX, l.resetY
}

func (l *Level) GetName() string {
	return l.levelLDTK.Name
}

func (l *Level) GetGameEntryPos() (float64, float64) {
	entityLayer := errs.Must(l.levelLDTK.GetLayerByName(entityLayerName))
	for _, entity := range entityLayer.Entities {
		if entity.Name == gameEntryPosEntityName {
			return entity.Px[0], entity.Px[1]
		}
	}
	return 0, 0
}

// ------ INTERNAL ------
func drawTiles(
	tiles []ebitenLDTK.Tile,
	tileset *ebiten.Image,
	targetLayer *ebiten.Image,
	tileSize float64,
	// camX, camY float64,
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

		ebitenrenderutil.DrawAtScaled(tileset.SubImage(
			image.Rect(
				int(tile.Src[0]),
				int(tile.Src[1]),
				int(tile.Src[0]+tileSize),
				int(tile.Src[1]+tileSize),
			),
		).(*ebiten.Image),
			targetLayer,
			tile.Px[0], tile.Px[1], scaleX, scaleY, 0.5, 0.5)

		// ebitenrenderutil.DrawAtScaled(tileset.Image.SubImage(
		// 	image.Rect(
		// 		int(tile.Src[0]),
		// 		int(tile.Src[1]),
		// 		int(tile.Src[0]+tileSize),
		// 		int(tile.Src[1]+tileSize),
		// 	),
		// ).(*ebiten.Image),
		// 	targetLayer,
		// 	tile.Px[0]-camX, tile.Px[1]-camY, scaleX, scaleY, 0.5, 0.5)
	}
}

func (l *Level) restoreFromMemory(levelMemory *LevelMemory) {
	entityLayer := errs.Must(l.levelLDTK.GetLayerByName(entityLayerName))
	for _, entity := range entityLayer.Entities {
		if entity.Name != slamboxEntityName {
			continue
		}

		for _, slambox := range l.slamboxes {
			if slambox.LinkID != entity.Iid {
				continue
			}
			slambox.SetPos(entity.Px[0], entity.Px[1])
		}
	}
}

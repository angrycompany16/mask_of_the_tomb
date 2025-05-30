package world

import (
	"fmt"
	"image"
	"image/color"
	"mask_of_the_tomb/internal/core/arrays"
	"mask_of_the_tomb/internal/core/colors"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"
	threads "mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/assettypes"
	"mask_of_the_tomb/internal/libraries/camera"
	"mask_of_the_tomb/internal/libraries/entities/door"
	"mask_of_the_tomb/internal/libraries/entities/grass"
	"mask_of_the_tomb/internal/libraries/entities/hazard"
	"mask_of_the_tomb/internal/libraries/entities/turret"
	"mask_of_the_tomb/internal/libraries/particles"
	"mask_of_the_tomb/internal/libraries/physics"
	"math"
	"path/filepath"
	"slices"
	"time"

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
	grassEntityName        = "Grass"
	hazardEntityName       = "Hazard"
	turretEntityName       = "TurretEnemy"
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
	// TODO: *very* temporary solution
	playerspaceNormalTilemap = errs.MustNewImageFromFile(playerSpaceNormalTilemapPath)
	PlayerLightRadius        = 200.0
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
	tileLayers               rendering.LayerList
	normalLayers             rendering.LayerList
	frameLayers              rendering.LayerList
	fogShader                *ebiten.Shader
	vignetteShader           *ebiten.Shader
	pixelLightShader         *ebiten.Shader
	playerLightBreatheTicker time.Ticker
	resetX, resetY           float64
	ambientParticles         *particles.ParticleSystem
	slamboxes                []*Slambox
	hazards                  []*hazard.Hazard
	doors                    []door.Door
	grassEntities            []grass.Grass
	turrets                  []*turret.Turret
}

// ------ CONSTRUCTOR ------
// TODO: Refactor because this is very big
func newLevel(levelLDTK *ebitenLDTK.Level, defs *ebitenLDTK.Defs) (*Level, error) {
	newLevel := &Level{}
	newLevel.levelLDTK = levelLDTK
	newLevel.defs = defs
	newLevel.name = levelLDTK.Name
	newLevel.bgColor = errs.Must(colors.HexToRGB(levelLDTK.BgColorHex))

	// TODO: set particle system bounds based on level size
	newLevel.fogShader = errs.Must(assettypes.GetShaderAsset("fogShader"))
	newLevel.vignetteShader = errs.Must(assettypes.GetShaderAsset("vignetteShader"))
	newLevel.pixelLightShader = errs.Must(assettypes.GetShaderAsset("pixelLightsShader"))
	newLevel.ambientParticles = errs.Must(particles.GetParticleSystemAsset("ambientParticles"))

	newLevel.playerLightBreatheTicker = *time.NewTicker(time.Millisecond * 560)

	// Most of these layers are completely unnecessary, so maybe it would be a good idea to
	// delete a few of them to save performance
	// After all we do have some pretty dramatic frame drops when switching levels
	newLevel.tileLayers = rendering.NewLayerList(int(levelLDTK.PxWid), int(levelLDTK.PxHei))
	newLevel.normalLayers = rendering.NewLayerList(int(levelLDTK.PxWid), int(levelLDTK.PxHei))
	newLevel.frameLayers = rendering.NewLayerList(int(levelLDTK.PxWid), int(levelLDTK.PxHei))

	layerMap := map[string]*ebiten.Image{
		"Foreground":       newLevel.tileLayers.Foreground,
		"PlayerspaceAlt":   newLevel.tileLayers.Playerspace,
		"Playerspace":      newLevel.tileLayers.Playerspace,
		"Props":            newLevel.tileLayers.Midground,
		"MidgroundSprites": newLevel.tileLayers.Midground,
		"BackgroundTiles":  newLevel.tileLayers.Background,
	}

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

	entityLayer := errs.Must(levelLDTK.GetLayerByName(entityLayerName))
	for _, entity := range entityLayer.Entities {
		switch entity.Name {
		case hazardEntityName:
			newLevel.hazards = append(newLevel.hazards, hazard.NewHazard(&entity))
		case doorEntityName:
			newLevel.doors = append(newLevel.doors, door.NewDoor(&entity))
		case slamboxEntityName:
			newLevel.slamboxes = append(newLevel.slamboxes, NewSlambox(&entity))
		case grassEntityName:
			newLevel.grassEntities = append(newLevel.grassEntities, grass.NewGrass(&entity, 16, rendering.ScreenLayers.Playerspace))
		case turretEntityName:
			newLevel.turrets = append(newLevel.turrets, turret.NewTurret(&entity, entityLayer.GridSize))
		}
	}

	// NOTE: We need to loop twice to ensure that all slamboxes have been added
	// before we link them together
	for _, slambox := range newLevel.slamboxes {
		for _, hazard := range newLevel.hazards {
			if slices.Contains(slambox.attachedHazardIDs, hazard.LinkID) {
				slambox.attachedHazards = append(slambox.attachedHazards, hazard)
				hazard.PosOffsetX = hazard.Hitbox.Left() - slambox.Collider.Left()
				hazard.PosOffsetY = hazard.Hitbox.Top() - slambox.Collider.Top()
			}
		}

		for _, otherSlambox := range newLevel.slamboxes {
			if slices.Contains(slambox.OtherLinkIDs, otherSlambox.LinkID) {
				slambox.ConnectedBoxes = append(slambox.ConnectedBoxes, otherSlambox)
			}
		}
		slambox.CreateSprite(errs.Must(assettypes.GetImageAsset("slamboxTilemap")))
		// slambox.CreateSprite(assetloader.GetAsset("slamboxTilemap").(*assettypes.ImageAsset).Image)
	}

	// Optimization yeah
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

		tileset := errs.Must(newLevel.defs.GetTilesetByUid(layer.TilesetUid))
		tilesize := tileset.TileGridSize
		tilesetImage := tileset.Image

		if targetRenderLayer == newLevel.tileLayers.Midground {
			// tilesetImage = midgroundNormalTilemap
			// drawTiles(tiles, tilesetImage, newLevel.normalLayers.Midground, tilesize)
		} else if targetRenderLayer == newLevel.tileLayers.Playerspace {
			drawTiles(tiles, playerspaceNormalTilemap, newLevel.normalLayers.Playerspace, tilesize)
		}

		drawTiles(tiles, tilesetImage, targetRenderLayer, tilesize)
	}

	return newLevel, nil
}

// ------ ENTITY ------
func (l *Level) Update(playerX, playerY, playerVelX, playerVelY float64) {
	for _, slambox := range l.slamboxes {
		slambox.Update()
	}
	for _, grassEntity := range l.grassEntities {
		grassEntity.Update(playerX, playerY, playerVelX, playerVelY)
	}
	for _, turret := range l.turrets {
		hit, x, y := l.TilemapCollider.Raycast(turret.Hitbox.Left(), turret.Hitbox.Top(), turret.GetAimDir(), l.GetSlamboxColliders())
		if hit {
			turret.RayEndX = x
			turret.RayEndY = y
		}
		for _, rect := range l.GetSlamboxColliders() {
			if rect.Rect.Overlapping(&turret.Hitbox) {
				turret.Die()
			}
		}
	}

	if _, tick := threads.Poll(l.playerLightBreatheTicker.C); tick {
		PlayerLightRadius = 205.0 - math.Copysign(5, PlayerLightRadius-210.0)
	}

	l.ambientParticles.Update()
}

func (l *Level) Draw(ctx rendering.Ctx) {
	l.frameLayers.Playerspace.Clear()
	rendering.ScreenLayers.Background2.Fill(l.bgColor)

	// Render fog layer
	shaderOp := ebiten.DrawRectShaderOptions{}
	shaderOp.Uniforms = map[string]any{
		"Time":       resources.Time / 5,
		"Amplitude":  1.0,
		"Frequency":  0.025,
		"Strength":   0.7,
		"Threshold":  0.4,
		"Color":      [4]float64{37.0 / 255, 49.0 / 255, 94.0 / 255, 1.0},
		"Center":     [2]float64{0.5, 0.5},
		"Resolution": [2]float64{rendering.GAME_WIDTH, rendering.GAME_HEIGHT},
		"PlayerPos":  [2]float64{ctx.PlayerX, ctx.PlayerY},
	}
	shaderOp.Blend = ebiten.BlendSourceOver
	// TODO: Move fog with camera position
	rendering.ScreenLayers.Background2.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.fogShader, &shaderOp)

	rendering.ScreenLayers.Foreground.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.vignetteShader, &shaderOp)

	for _, turretEntity := range l.turrets {
		turretEntity.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	for _, box := range l.slamboxes {
		box.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	for _, hazard := range l.hazards {
		hazard.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	for _, grassEntity := range l.grassEntities {
		grassEntity.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	ebitenrenderutil.DrawAt(l.tileLayers.Playerspace, l.frameLayers.Playerspace, 0, 0)

	// TODO: Give a light to each turret guy
	shaderOp = l.GetShaderOp(ctx, l.frameLayers.Playerspace)
	shaderOp.Uniforms["AmbientLight"] = [3]float64{0.3, 0.3, 0.3}
	rendering.ScreenLayers.Playerspace.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.pixelLightShader, &shaderOp)

	shaderOp = l.GetShaderOp(ctx, l.tileLayers.Midground)
	rendering.ScreenLayers.Midground.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.pixelLightShader, &shaderOp)

	shaderOp = l.GetShaderOp(ctx, l.tileLayers.Background)
	rendering.ScreenLayers.Background.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.pixelLightShader, &shaderOp)

	l.ambientParticles.Draw(rendering.WithLayer(ctx, rendering.ScreenLayers.Foreground))
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

func (l *Level) CheckTurretHit(playerHitBox *maths.Rect) bool {
	for _, turret := range l.turrets {
		if turret.ShouldFire(playerHitBox) {
			return true
		}
	}
	return false
}

// Get all the rect colliders that are not connected to slambox
func (l *Level) GetDisconnectedColliders(_slambox *Slambox) []*physics.RectCollider {
	// I love writing unreadable code
	return arrays.MapSlice(
		arrays.Filter(
			l.slamboxes, func(s *Slambox) bool { return !slices.Contains(_slambox.ConnectedBoxes, s) && s != _slambox },
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

func (l *Level) GetShaderOp(ctx rendering.Ctx, src *ebiten.Image) ebiten.DrawRectShaderOptions {
	trueCamX, trueCamY := camera.GetStablePos()

	shaderOp := ebiten.DrawRectShaderOptions{}
	shaderOp.Images = [4]*ebiten.Image{
		// NEVER touch the first texture argument. EVER.
		nil,
		src.SubImage(image.Rect(int(trueCamX), int(trueCamY), int(trueCamX+rendering.GAME_WIDTH), int(trueCamY+rendering.GAME_HEIGHT))).(*ebiten.Image),
		nil,
		nil,
	}

	shakeX, shakeY := camera.GetShake()

	shaderOp.Uniforms = map[string]any{
		"CamShake":     [2]float64{shakeX, shakeY},
		"Time":         resources.Time / 5,
		"PositionsX":   [10]float64{ctx.PlayerX - ctx.CamX},
		"PositionsY":   [10]float64{ctx.PlayerY - ctx.CamY},
		"InnerRadii":   [10]float64{0.0},
		"OuterRadii":   [10]float64{PlayerLightRadius},
		"ZOffsets":     [10]float64{0.2},
		"Intensities":  [10]float64{0.6},
		"ColorsR":      [10]float64{1.0},
		"ColorsG":      [10]float64{1.0},
		"ColorsB":      [10]float64{1.0},
		"AmbientLight": [3]float64{0.6, 0.6, 0.6},
	}

	return shaderOp
}

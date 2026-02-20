package world

import (
	"fmt"
	"image"
	"image/color"
	"mask_of_the_tomb/internal/core/arrays"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/colors"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/shaders"
	"mask_of_the_tomb/internal/libraries/camera"
	"mask_of_the_tomb/internal/libraries/entities"
	"mask_of_the_tomb/internal/libraries/particles"
	"mask_of_the_tomb/internal/libraries/slambox"
	"mask_of_the_tomb/internal/libraries/speechbubble"
	"math"
	"slices"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

var names = resources.LDTKNames

const (
// entityLayerName            = "Entities"
// playerSpaceLayerName       = "Playerspace"
// spawnPosEntityName         = "DefaultSpawnPos"
// doorEntityName             = "Door"
// spawnPointEntityName       = "SpawnPoint"
// slamboxEntityName          = "Slambox"
// SpikeIntGridName           = "Spikes"
// gameEntryPosEntityName     = "GameEntryPos"
// grassEntityName            = "Grass"
// hazardEntityName           = "Hazard"
// turretEntityName           = "TurretEnemy"
// catcherEntityName          = "Catcher"
// platformEntityName         = "OneWayPlatform"
// lanternEntityName          = "Lantern"
// chainNodeEntityName        = "SlamboxChainNode"
// slamboxChainEntityName     = "SlamboxChain"
// testSpeechBubbleEntityName = "TestSpeechBubble"
// levelTitleFieldName        = "Title"
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
)

// newLevel.chainNodes = append(newLevel.chainNodes, entities.NewChainNode(&entity))
type SlamboxPosition struct {
	X, Y float64
}

// Sort of an everything-container
type Level struct {
	name                     string
	defs                     *ebitenLDTK.Defs
	levelLDTK                *ebitenLDTK.Level
	tileSize                 float64
	bgColor                  color.Color
	playerspaceNormalTilemap *ebiten.Image
	midgroundNormalTilemap   *ebiten.Image
	tileLayers               rendering.LayerList
	frameLayers              rendering.LayerList
	fogShader                *ebiten.Shader
	vignetteShader           *ebiten.Shader
	pixelLightShader         *ebiten.Shader
	resetX, resetY           float64
	ambientParticles         *particles.ParticleSystem
	grassTilemap             *ebiten.Image
	slamboxEnvironment       *slambox.SlamboxEnvironment
	// In a sense these are all game objects
	slamboxEntities  []*SlamboxEntity
	hazards          []*entities.Hazard
	doors            []entities.Door
	grassEntities    []entities.Grass
	turrets          []*entities.Turret
	catchers         []*entities.Catcher
	platforms        []*entities.Platform
	lanterns         []*entities.Lantern
	chainNodes       []*entities.ChainNode
	testSpeechBubble *speechbubble.SpeechBubble
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
	newLevel.ambientParticles = errs.Must(assettypes.GetYamlAsset("ambientParticles")).(*particles.ParticleSystem)
	newLevel.ambientParticles.Init()
	newLevel.grassTilemap = errs.Must(assettypes.GetImageAsset("grassTilemap"))

	// Empty constructor for later setup
	newLevel.slamboxEnvironment = slambox.NewSlamboxEnvironment(1, make([][]bool, 0), make([]*slambox.Slambox, 0), make([]*slambox.SlamboxGroup, 0), make([]*slambox.SlamboxChain, 0))

	// Most of these layers are completely unnecessary, so maybe it would be a good idea to
	// delete a few of them to save performance
	// After all we do have some pretty dramatic frame drops when switching levels
	newLevel.tileLayers = rendering.NewLayerList(int(levelLDTK.PxWid), int(levelLDTK.PxHei))
	newLevel.frameLayers = rendering.NewLayerList(int(levelLDTK.PxWid), int(levelLDTK.PxHei))

	layerMap := map[string]*ebiten.Image{
		"Foreground":       newLevel.tileLayers.Foreground,
		"PlayerspaceAlt":   newLevel.tileLayers.Playerspace,
		"Playerspace":      newLevel.tileLayers.Playerspace,
		"Props":            newLevel.tileLayers.Midground,
		"MidgroundSprites": newLevel.tileLayers.Midground,
		"BackgroundTiles":  newLevel.tileLayers.Background,
	}

	playerspace, err := levelLDTK.GetLayerByName(names.PlayerSpaceLayer)
	if err != nil {
		newLevel.tileSize = 1
		return newLevel, nil
	}

	var spikeIntGridID int
	for _, layerDef := range defs.LayerDefs {
		if layerDef.Name == names.PlayerSpaceLayer {
			spikeIntGridID = layerDef.GetIntGridValue(names.SpikeIntGrid)
		}
	}

	intGridCSV := playerspace.ExtractLayerCSV([]int{spikeIntGridID})
	newLevel.slamboxEnvironment.SetTiles(intGridCSV)

	newLevel.tileSize = float64(playerspace.GridSize)
	newLevel.slamboxEnvironment.SetTileSize(playerspace.GridSize)

	entityLayer := errs.Must(levelLDTK.GetLayerByName(names.EntityLayer))
	for _, entity := range entityLayer.Entities {
		switch entity.Name {
		case names.HazardEntity:
			newLevel.hazards = append(newLevel.hazards, entities.NewHazard(&entity))
		case names.DoorEntity:
			newLevel.doors = append(newLevel.doors, entities.NewDoor(&entity))
		case names.SlamboxEntity:
			newLevel.slamboxEntities = append(newLevel.slamboxEntities, NewSlamboxEntity(&entity, newLevel.slamboxEnvironment, levelLDTK))
		case names.GrassEntity:
			newLevel.grassEntities = append(newLevel.grassEntities, entities.NewGrass(&entity, 16, newLevel.grassTilemap, rendering.ScreenLayers.Playerspace))
		case names.TurretEntity:
			newLevel.turrets = append(newLevel.turrets, entities.NewTurret(&entity, entityLayer.GridSize))
		case names.CatcherEntity:
			newLevel.catchers = append(newLevel.catchers, entities.NewCatcher(&entity))
		case names.PlatformEntity:
			newLevel.platforms = append(newLevel.platforms, entities.NewPlatform(&entity, entityLayer.GridSize))
		case names.LanternEntity:
			newLevel.lanterns = append(newLevel.lanterns, entities.NewLantern(&entity, entityLayer.GridSize))
		// case chainNodeEntityName:
		// 	newLevel.chainNodes = append(newLevel.chainNodes, entities.NewChainNode(&entity))
		case names.TestSpeechBubbleEntity:
			fmt.Println(entity.Px[0], entity.Px[1])
			newLevel.testSpeechBubble = speechbubble.NewSpeechBubble(
				entity.Px[0], entity.Px[1], entity.Width, entity.Height,
			)
		}
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

		drawTiles(tiles, tilesetImage, targetRenderLayer, tilesize)
	}

	return newLevel, nil
}

// ------ ENTITY ------
func (l *Level) Update(playerX, playerY, playerVelX, playerVelY float64) {
	for _, slambox := range l.slamboxEntities {
		slambox.Update()
	}
	for _, grassEntity := range l.grassEntities {
		grassEntity.Update(playerX, playerY, playerVelX, playerVelY)
	}
	for _, door := range l.doors {
		door.Update()
	}
	// for _, turret := range l.turrets {
	// hit, x, y := l.TilemapCollider.Raycast(turret.Hitbox.Left(), turret.Hitbox.Top(), turret.GetAimDir(), l.GetSlamboxRects())
	// if hit {
	// 	turret.RayEndX = x
	// 	turret.RayEndY = y
	// }
	// turret.Update((l.GetSlamboxRects()))
	// }
	for _, lantern := range l.lanterns {
		lantern.Update(playerX, playerY, playerVelX, playerVelY)
	}

	l.ambientParticles.Update()
	l.slamboxEnvironment.Update()
	// l.testSpeechBubble.Update()

	// if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	// 	cursorX, cursorY := ebiten.CursorPosition()
	// 	cursorX = cursorX / rendering.PIXEL_SCALE
	// 	cursorY = cursorY / rendering.PIXEL_SCALE
	// 	l.testSpeechBubble.SetAnchor(float64(cursorX), float64(cursorY))
	// }
}

func (l *Level) Draw(ctx rendering.Ctx, playerLight *shaders.Light) {
	l.frameLayers.Foreground.Clear()
	l.frameLayers.Playerspace.Clear()
	l.frameLayers.Midground.Clear()
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

	// Consider creating an entity interface
	for _, turretEntity := range l.turrets {
		turretEntity.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	for _, slamboxEntity := range l.slamboxEntities {
		slamboxEntity.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	for _, hazard := range l.hazards {
		hazard.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	for _, grassEntity := range l.grassEntities {
		grassEntity.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	for _, door := range l.doors {
		door.Draw(rendering.WithLayer(ctx, l.frameLayers.Midground))
	}

	for _, catcher := range l.catchers {
		catcher.Draw(rendering.WithLayer(ctx, l.frameLayers.Midground))
	}

	for _, lantern := range l.lanterns {
		lantern.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	}

	// for _, chainNode := range l.chainNodes {
	// 	chainNode.Draw(rendering.WithLayer(ctx, l.frameLayers.Playerspace))
	// }

	ebitenrenderutil.DrawAt(l.tileLayers.Playerspace, l.frameLayers.Playerspace, 0, 0)
	ebitenrenderutil.DrawAt(l.tileLayers.Midground, l.frameLayers.Midground, 0, 0)

	// TODO: Give a light to each turret guy
	shakeX, shakeY := camera.GetShake()
	camX, camY := camera.GetStablePos()
	shaderOp = shaders.MakeShaderOp(
		slices.Concat(
			arrays.MapSlice(l.turrets, func(turret *entities.Turret) *shaders.Light { return turret.Light }),
			arrays.MapSlice(l.lanterns, func(lantern *entities.Lantern) *shaders.Light { return lantern.Light }),
			arrays.MapSlice(l.slamboxEntities, func(slambox *SlamboxEntity) *shaders.Light { return slambox.Light }),
			[]*shaders.Light{playerLight},
		),
		rendering.GAME_WIDTH,
		rendering.GAME_HEIGHT,
		camX,
		camY,
		shakeX,
		shakeY,
		0.3,
		0.3,
		0.3,
		resources.Time/5,
		rendering.GAME_WIDTH,
		rendering.GAME_HEIGHT,
		l.frameLayers.Playerspace,
	)
	rendering.ScreenLayers.Playerspace.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.pixelLightShader, &shaderOp)

	shaderOp = shaders.ChangeSrc(shaderOp, camX, camY, rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.frameLayers.Midground)
	rendering.ScreenLayers.Midground.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.pixelLightShader, &shaderOp)

	shaderOp = shaders.ChangeSrc(shaderOp, camX, camY, rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.tileLayers.Background)
	rendering.ScreenLayers.Background.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.pixelLightShader, &shaderOp)

	// if l.testSpeechBubble != nil {
	// 	l.testSpeechBubble.Draw(rendering.WithLayer(ctx, l.frameLayers.Foreground))
	// }
	rendering.ScreenLayers.Foreground.DrawImage(l.frameLayers.Foreground, &ebiten.DrawImageOptions{})

	l.ambientParticles.Draw(rendering.WithLayer(ctx, rendering.ScreenLayers.Foreground))
}

func (l *Level) ProjectRect(rect *maths.Rect, dir maths.Direction) (maths.Rect, float64) {
	// Get platform rects
	platformHitboxes := make([]*maths.Rect, 0)
	if dir == maths.DirUp {
		platformHitboxes = l.GetPlatformHitboxes(false)
	} else if dir == maths.DirDown {
		platformHitboxes = l.GetPlatformHitboxes(true)
	}

	otherRects := slices.Concat(
		l.slamboxEnvironment.GetSlamboxRects(-1),
		l.slamboxEnvironment.GetSlamboxGroupRects(-1),
		l.slamboxEnvironment.GetSlamboxChainRects(-1),
		platformHitboxes,
	)
	projRect, dist := l.slamboxEnvironment.ProjectRect(*rect, dir, math.Inf(1), otherRects)
	return projRect, dist
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

func (l *Level) GetCatcherRects() []*maths.Rect {
	return arrays.MapSlice(l.catchers, func(c *entities.Catcher) *maths.Rect { return c.Hitbox })
}

// Get all the platforms matching the movement direction
func (l *Level) GetPlatformHitboxes(up bool) []*maths.Rect {
	hitboxes := make([]*maths.Rect, 0)
	for _, platform := range l.platforms {
		if platform.Up == up {
			hitboxes = append(hitboxes, platform.Hitbox)
		}
	}
	return hitboxes
}

func (l *Level) GetSlamboxPositions() []SlamboxPosition {
	return arrays.MapSlice(l.slamboxEntities, func(s *SlamboxEntity) SlamboxPosition {
		return SlamboxPosition{X: s.rect.Left(), Y: s.rect.Top()}
	})
}

// Queries the slambox environment for overlaps with the playercollider extended
// 1px in the slam direction. Returns the first hit, as well as an indication
// of whether this is a slambox, a
func (l *Level) GetSlamboxHit(playerCollider *maths.Rect, dir maths.Direction) slambox.QueryResult {
	extendedRect := playerCollider.Extended(dir, 1)
	return l.slamboxEnvironment.QuerySlamboxes(extendedRect)
}

func (l *Level) SlamSlambox(id int, dir maths.Direction) {
	slambox := l.GetSlamboxEntity(id)
	slambox.StartSlam(id, dir)
}

func (l *Level) SlamSlamboxGroup(id int, dir maths.Direction) {
	slambox := l.GetSlamboxEntity(id)
	slambox.StartSlam(id, dir)
}

func (l *Level) GetSlamboxEntity(id int) *SlamboxEntity {
	for _, slamboxEntity := range l.slamboxEntities {
		if slamboxEntity.backendID == id {
			return slamboxEntity
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
	entityLayer := errs.Must(l.levelLDTK.GetLayerByName(names.EntityLayer))
	for _, entity := range entityLayer.Entities {
		if entity.Name != names.SpawnPosEntity {
			continue
		}
		return entity.Px[0], entity.Px[1]
	}
	for _, entity := range entityLayer.Entities {
		if entity.Name != names.DoorEntity {
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

func (l *Level) GetTitle() string {
	return errs.Must(l.levelLDTK.GetFieldByName(names.LevelTitleField)).String
}

func (l *Level) GetGameEntryPos() (float64, float64) {
	entityLayer := errs.Must(l.levelLDTK.GetLayerByName(names.EntityLayer))
	for _, entity := range entityLayer.Entities {
		if entity.Name == names.GameEntryPos {
			return entity.Px[0], entity.Px[1]
		}
	}
	return 0, 0
}

// func (l *Level) GetChainNodes() []*entities.ChainNode {
// 	return l.chainNodes
// }

// ------ INTERNAL ------
func drawTiles(
	tiles []ebitenLDTK.Tile,
	tileset *ebiten.Image,
	targetLayer *ebiten.Image,
	tileSize float64,
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

// TODO: Respawn enemies
func (l *Level) reset() {
	entityLayer := errs.Must(l.levelLDTK.GetLayerByName(names.EntityLayer))
	for _, entity := range entityLayer.Entities {
		if entity.Name != names.SlamboxEntity {
			continue
		}

		for _, slambox := range l.slamboxEntities {
			slambox.SetPos(entity.Px[0], entity.Px[1])
		}
	}
}

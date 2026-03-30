package ldtkworld

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/actors/trigger"
	"mask_of_the_tomb/internal/engine/actors/vectorgraphic"
	ldtktilelayer "mask_of_the_tomb/internal/game/actors/LDTKTileLayer"
	"mask_of_the_tomb/internal/game/actors/autotilesprite"
	"mask_of_the_tomb/internal/game/actors/doorv2"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/actors/slamboxtilemap"
	"mask_of_the_tomb/internal/game/actors/tracker"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

var layerMap = map[string]string{
	"Foreground":       "Foreground",
	"PlayerspaceAlt":   "Playerspace",
	"Playerspace":      "Playerspace",
	"Props":            "Midground",
	"MidgroundSprites": "Midground",
	"BackgroundTiles":  "Background",
}

type LDTKLevel struct {
	*graphic.Graphic
	worldSrcPath string `debug:"auto"`
	levelName    string `debug:"auto"`
	LDTKData     *assetloader.AssetRef[assettypes.LDTKData]
}

func (l *LDTKLevel) OnTreeAdd(node *engine.Node, servers *engine.Commands) {
	l.Graphic.OnTreeAdd(node, servers)
	l.LDTKData = assetloader.StageAsset[assettypes.LDTKData](
		servers.AssetLoader(),
		l.worldSrcPath,
		assettypes.NewLDTKAsset(l.worldSrcPath),
	)

	// Pre-load this
	assetloader.StageAsset[assettypes.ImageAsset](
		servers.AssetLoader(),
		"sprites/environment/slambox_tilemap.png",
		assettypes.NewImageAsset("sprites/environment/slambox_tilemap.png"),
	)

	// Also pre-load this
	assetloader.StageAsset[assettypes.ImageAsset](
		servers.AssetLoader(),
		"sprites/environment/door_v2.png",
		assettypes.NewImageAsset("sprites/environment/door_v2.png"),
	)
}

func (l *LDTKLevel) Init(cmd *engine.Commands) {
	l.Graphic.Update(cmd)
	world := l.LDTKData.Value().World
	tilesetMap := l.LDTKData.Value().Tilesets
	defs := world.Defs
	level := utils.Must(world.GetLevelByName(l.levelName))

	playerspace, err := level.GetLayerByName("Playerspace")
	if err != nil {
		fmt.Println("Error when loading level:", err)
		return
	}

	var spikeIntGridID int
	for _, layerDef := range defs.LayerDefs {
		if layerDef.Name == "Playerspace" {
			spikeIntGridID = layerDef.GetIntGridValue("Spikes")
		}
	}

	intGridCSV := playerspace.ExtractLayerCSV([]int{spikeIntGridID})
	cmd.Scene().AddChild("SlamboxTilemap", slamboxtilemap.NewSlamboxTilemap(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		),
		intGridCSV,
		int(playerspace.GridSize),
	), l.GetNode(), cmd)

	for i := range level.Layers {
		layer := level.Layers[i]

		if layer.Type == ebitenLDTK.LayerTypeEntities {
			continue
		}

		tileset := utils.Must(defs.GetTilesetByUid(layer.TilesetUid))
		tileSize := tileset.TileGridSize
		tilesetImg := tilesetMap[tileset.Name]

		cmd.Scene().AddChild(layer.Name, ldtktilelayer.NewLDTKTilemapLayer(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			&layer, tilesetImg, "Playerspace",
			-i,
			int(tileSize),
			int(level.PxWid),
			int(level.PxHei),
		), l.GetNode(), cmd)
	}

	cmd.Scene().AddChild("BackgroundColor", vectorgraphic.NewVectorGraphic(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		),
		func(img *ebiten.Image) {
			vector64.FillRect(
				img,
				0, 0, level.PxWid, level.PxHei,
				utils.Must(colors.HexToRGB(level.BgColorHex)), false,
			)
		},
		"Background",
		-len(level.Layers)-1,
		int(level.PxWid),
		int(level.PxHei),
	), l.GetNode(), cmd)

	entityLayer := utils.Must(level.GetLayerByName("Entities"))
	for _, entity := range entityLayer.Entities {
		switch entity.Name {
		// case names.HazardEntity:
		// 	newLevel.hazards = append(newLevel.hazards, entities.NewHazard(&entity))
		// case names.DoorEntity:
		// 	newLevel.doors = append(newLevel.doors, entities.NewDoor(&entity))
		case "Slambox":
			// Spawn slambox bundle
			slamboxNode := cmd.Scene().SpawnActor("Slambox",
				slamboxactor.NewSlambox(
					tracker.NewTracker(
						graphic.NewGraphic(
							transform2D.NewTransform2D(
								nodeactor.NewNode(),
							),
						), 7.5, entity.Px[0], entity.Px[1],
					),
					slamboxactor.WithPos(entity.Px[0], entity.Px[1]),
					slamboxactor.WithSize(entity.Width, entity.Height),
				),
				cmd,
			)

			autotilesprite := cmd.Scene().AddChild("Sprite", autotilesprite.NewAutoTileSprite(
				graphic.NewGraphic(
					transform2D.NewTransform2D(
						nodeactor.NewNode(),
					),
				),
				autotilesprite.WithSize(entity.Width, entity.Height),
				autotilesprite.WithTilemap("sprites/environment/slambox_tilemap.png"),
			), slamboxNode, cmd)

			// This really is just horrible...
			// The quickest solution i could possibly find
			slamboxNode.GetValue().Init(cmd)
			autotilesprite.GetValue().Init(cmd)
		// newLevel.slamboxEntities = append(newLevel.slamboxEntities, NewSlamboxEntity(&entity, newLevel.slamboxEnvironment, levelLDTK))
		// case names.GrassEntity:
		// 	newLevel.grassEntities = append(newLevel.grassEntities, entities.NewGrass(&entity, 16, newLevel.grassTilemap, rendering.ScreenLayers.Playerspace))
		// case names.TurretEntity:
		// 	newLevel.turrets = append(newLevel.turrets, entities.NewTurret(&entity, entityLayer.GridSize))
		// case names.CatcherEntity:
		// 	newLevel.catchers = append(newLevel.catchers, entities.NewCatcher(&entity))
		// case names.PlatformEntity:
		// 	newLevel.platforms = append(newLevel.platforms, entities.NewPlatform(&entity, entityLayer.GridSize))
		// case names.LanternEntity:
		// 	newLevel.lanterns = append(newLevel.lanterns, entities.NewLantern(&entity, entityLayer.GridSize))
		case "DoorV2":
			// Spawn door bundle
			// Door actor - Just connects child logic, contains other side IDs
			// Children:
			// animatedsprite
			// triggerZone

			doorNode := cmd.Scene().SpawnActor("Door", doorv2.NewDoorV2(
				graphic.NewGraphic(
					transform2D.NewTransform2D(
						nodeactor.NewNode(),
						transform2D.WithPos(entity.Px[0], entity.Px[1]),
					),
				), &entity, &level,
			), cmd)

			doorSprite := cmd.Scene().AddChild("Sprite", animatedsprite.NewAnimator(
				graphic.NewGraphic(
					transform2D.NewTransform2D(
						nodeactor.NewNode(),
					),
				), map[int]*animatedsprite.Animation{
					0: animatedsprite.NewAnimation(
						animatedsprite.AnimationInfo{
							Name:              "Idle",
							SpriteSheetPath:   "assets/sprites/environment/door_v2-idle-Sheet.png",
							SpriteSheetFormat: animatedsprite.Strip,
							LoopMode:          animatedsprite.Loop,
							FrameDelay:        100,
							NextAnimationId:   -1,
						}, 48, 16,
					),
					1: animatedsprite.NewAnimation(
						animatedsprite.AnimationInfo{
							Name:              "Open",
							SpriteSheetPath:   "assets/sprites/environment/door_v2-open-Sheet.png",
							SpriteSheetFormat: animatedsprite.Strip,
							LoopMode:          animatedsprite.Once,
							FrameDelay:        100,
							NextAnimationId:   -1,
						}, 48, 16,
					),
					2: animatedsprite.NewAnimation(
						animatedsprite.AnimationInfo{
							Name:              "Close",
							SpriteSheetPath:   "assets/sprites/environment/door_v2-close-Sheet.png",
							SpriteSheetFormat: animatedsprite.Strip,
							LoopMode:          animatedsprite.Once,
							FrameDelay:        100,
							NextAnimationId:   -1,
						}, 48, 16,
					),
				}, "Playerspace", 5,
			), doorNode, cmd)

			doorTrigger := cmd.Scene().AddChild("Trigger", trigger.NewTrigger(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			), doorNode, cmd)

			doorNode.GetValue().Init(cmd)
			doorSprite.GetValue().Init(cmd)
			doorTrigger.GetValue().Init(cmd)

			// newLevel.doorsV2 = append(newLevel.doorsV2, entities.NewDoorV2(&entity, levelLDTK))
			// 	// case chainNodeEntityName:
			// // 	newLevel.chainNodes = append(newLevel.chainNodes, entities.NewChainNode(&entity))
			// case names.TestSpeechBubbleEntity:
			// 	fmt.Println(entity.Px[0], entity.Px[1])
			// 	newLevel.testSpeechBubble = speechbubble.NewSpeechBubble(
			// 		entity.Px[0], entity.Px[1], entity.Width, entity.Height, false,
			// 	)
			// }
		}
	}
}

func (l *LDTKLevel) DrawInspector(ctx *debugui.Context) {
	l.Transform2D.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, l)
}

func NewLDTKLevel(graphic *graphic.Graphic, levelName, worldSrcPath string) *LDTKLevel {
	return &LDTKLevel{
		Graphic:      graphic,
		levelName:    levelName,
		worldSrcPath: worldSrcPath,
	}
}

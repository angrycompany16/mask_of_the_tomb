package bundles

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/particles"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/actors/vectorgraphic"
	"mask_of_the_tomb/internal/engine/commands"
	ldtktilelayer "mask_of_the_tomb/internal/game/actors/LDTKTileLayer"
	"mask_of_the_tomb/internal/game/actors/autotilesprite"
	"mask_of_the_tomb/internal/game/actors/backgroundshader"
	"mask_of_the_tomb/internal/game/actors/doorv2"
	"mask_of_the_tomb/internal/game/actors/hazard"
	"mask_of_the_tomb/internal/game/actors/levelshader"
	"mask_of_the_tomb/internal/game/actors/platform"
	"mask_of_the_tomb/internal/game/actors/shaderactor"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/actors/slamboxtilemap"
	"mask_of_the_tomb/internal/game/actors/tracker"
	"mask_of_the_tomb/internal/game/actors/trigger"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
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

func MakeLDTKLevelBundle(levelIid string) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) {
		// gw, gh := cmd.Renderer.GetGameSize()

		// 1. Load and prepare data from LDTK
		LDTKData, ok := assetloader.GetAsset[assettypes.LDTKData](cmd.AssetLoader, "LDTK/world.ldtk")
		if !ok {
			fmt.Println("Unable to load LDTK world asset from assetloader when making level bundle. Returning.")
			return
		}

		world := LDTKData.Value().World
		tilesetMap := LDTKData.Value().Tilesets
		defs := world.Defs
		level := utils.Must(world.GetLevelByIid(levelIid))

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

		// 2. Spawn nodes for tilemap layers, including spikes
		envParentNode := scene.SpawnActor("Environment", transform2D.NewTransform2D(
			nodeactor.NewNode(),
		), cmd)

		intGridCSV := playerspace.ExtractLayerCSV([]int{spikeIntGridID})
		slamboxTilemapActor := slamboxtilemap.NewSlamboxTilemap(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			intGridCSV,
			int(playerspace.GridSize),
		)
		envParentNode.AddChild(slamboxTilemapActor, "SlamboxTilemap", engine.MakeOnTreeAdd(slamboxTilemapActor, cmd))

		for i := range level.Layers {
			layer := level.Layers[i]

			if layer.Type == ebitenLDTK.LayerTypeEntities {
				continue
			}

			tileset := utils.Must(defs.GetTilesetByUid(layer.TilesetUid))
			tileSize := tileset.TileGridSize
			tilesetImg := tilesetMap[tileset.Name]

			ldtkTileLayerActor := ldtktilelayer.NewLDTKTilemapLayer(
				graphic.NewGraphic(
					transform2D.NewTransform2D(
						nodeactor.NewNode(),
					),
				),
				&layer, tilesetImg, renderer.RenderTarget{
					Type: renderer.TEXTURE,
					Name: "LevelTextureRaw",
				},
				-i,
				int(tileSize),
				int(level.PxWid),
				int(level.PxHei),
				// false,
			)

			envParentNode.AddChild(ldtkTileLayerActor, layer.Name, engine.MakeOnTreeAdd(ldtkTileLayerActor, cmd))
		}

		scene.SpawnActor("BackgroundColor", vectorgraphic.NewVectorGraphic(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			vectorgraphic.WithDrawFunc(
				func(img *ebiten.Image) {
					vector64.FillRect(
						img,
						0, 0, level.PxWid, level.PxHei,
						utils.Must(colors.HexToRGB(level.BgColorHex)), false,
					)
				},
			),
			vectorgraphic.WithTarget(renderer.RenderTarget{
				Type: renderer.TEXTURE,
				Name: "BackgroundRaw",
			}),
			vectorgraphic.WithDrawOrder(-len(level.Layers)-1),
			vectorgraphic.WithImage(int(level.PxWid), int(level.PxHei)),
			vectorgraphic.WithPivot(0, 0),
		), cmd)

		// 3. Spawn entities (doors, slamboxes, platforms, etc...)
		entityLayer := utils.Must(level.GetLayerByName("Entities"))
		for _, entity := range entityLayer.Entities {
			switch entity.Name {
			case "Hazard":
				SpawnHazard(cmd, scene, &entity, envParentNode)
			case "Slambox":
				SpawnSlambox(cmd, scene, &entity, envParentNode)
			// case names.GrassEntity:
			// 	newLevel.grassEntities = append(newLevel.grassEntities, entities.NewGrass(&entity, 16, newLevel.grassTilemap, rendering.ScreenLayers.Playerspace))
			// case names.TurretEntity:
			// 	newLevel.turrets = append(newLevel.turrets, entities.NewTurret(&entity, entityLayer.GridSize))
			// case names.CatcherEntity:
			// 	newLevel.catchers = append(newLevel.catchers, entities.NewCatcher(&entity))
			case "Platform":
				SpawnPlatform(cmd, scene, &entity, envParentNode)
			// case names.LanternEntity:
			// 	newLevel.lanterns = append(newLevel.lanterns, entities.NewLantern(&entity, entityLayer.GridSize))
			case "DoorV2":
				SpawnDoor(cmd, scene, &entity, &level, envParentNode)
				// case chainNodeEntityName:
				// 	newLevel.chainNodes = append(newLevel.chainNodes, entities.NewChainNode(&entity))
				// case names.TestSpeechBubbleEntity:
				// 	fmt.Println(entity.Px[0], entity.Px[1])
				// 	newLevel.testSpeechBubble = speechbubble.NewSpeechBubble(
				// 		entity.Px[0], entity.Px[1], entity.Width, entity.Height, false,
				// 	)
				// }
			}
		}

		// 4. Spawn remaining actors (particlesystems, shaders, etc...)
		scene.SpawnActor("BackgroundParticles", particles.NewParticleSystemOld(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			), make([]*particles.Burst, 0), true, 0.5, 0, maths.RandomFloat64{0, 480}, maths.RandomFloat64{0, 270}, maths.RandomFloat64{0, 0}, maths.RandomFloat64{-5, 0}, maths.RandomFloat64{0, 4}, maths.RandomFloat64{0, 0.5}, maths.RandomFloat64{0, 0.01}, maths.RandomFloat64{0.7, 1.5}, maths.RandomFloat64{0.0, 0.0}, maths.RandomFloat64{3.0, 5.0}, maths.RandomFloat64{0.0, 2.0}, maths.RandomFloat64{0.0, 1.0}, [4]uint8{255, 255, 255, 255}, [4]uint8{255, 255, 255, 255}, 13, 13, "sprites/environment/star.png", renderer.RenderTarget{Type: renderer.SCREEN, Name: "Background"}, 1,
		), cmd)

		scene.SpawnActor("BackgroundShader",
			backgroundshader.NewBackgroundShader(
				shaderactor.NewShader(
					graphic.NewGraphic(
						transform2D.NewTransform2D(
							nodeactor.NewNode(),
						),
					), "shaders/fog.kage", cmd.Renderer.Textures["BackgroundRaw"], "Background", 0,
				),
			), cmd)

		scene.SpawnActor("LevelShader",
			levelshader.NewLevelShader(
				shaderactor.NewShader(
					graphic.NewGraphic(
						transform2D.NewTransform2D(
							nodeactor.NewNode(),
						),
					), "shaders/pixel_lights.kage", cmd.Renderer.Textures["LevelTextureRaw"], "Playerspace", 10,
				),
			), cmd)

		// scene.SpawnActor("Resetlistener", resetlistener.NewResetListener(nodeactor.NewNode()), cmd)
	}
}

func SpawnSlambox(cmd *commands.Commands, scene *engine.Scene, entity *ebitenLDTK.Entity, envParentNode *engine.Node) {
	// Spawn slambox bundle
	slamboxActor := slamboxactor.NewSlambox(
		tracker.NewTracker(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			), 7.5, entity.Px[0], entity.Px[1],
		),
		slamboxactor.WithPos(entity.Px[0], entity.Px[1]),
		slamboxactor.WithSize(entity.Width, entity.Height),
	)
	slamboxNode := envParentNode.AddChild(slamboxActor, "Slambox", engine.MakeOnTreeAdd(slamboxActor, cmd))

	autotileActor := autotilesprite.NewAutoTileSprite(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		), renderer.RenderTarget{
			Type: renderer.TEXTURE,
			Name: "LevelTextureRaw",
		},
		autotilesprite.WithSize(entity.Width, entity.Height),
		autotilesprite.WithTilemap("sprites/environment/slambox_tilemap.png"),
	)
	slamboxNode.AddChild(autotileActor, "Sprite", engine.MakeOnTreeAdd(autotileActor, cmd))
}

func SpawnDoor(cmd *commands.Commands, scene *engine.Scene, entity *ebitenLDTK.Entity, level *ebitenLDTK.Level, envParentNode *engine.Node) {
	directionField := utils.Must(entity.GetFieldByName("Direction"))
	direction := maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

	doorV2Actor := doorv2.NewDoorV2(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
				transform2D.WithPos(entity.Px[0], entity.Px[1]),
			),
		), entity, level,
	)
	doorNode := envParentNode.AddChild(doorV2Actor, "Door", engine.MakeOnTreeAdd(doorV2Actor, cmd))

	doorAnim := animatedsprite.NewAnimatedSprite(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
				transform2D.WithPos(entity.Width/2, entity.Height/2),
			),
		),
		map[string]*animatedsprite.Clip{
			"Idle": animatedsprite.NewClip(
				"sprites/environment/door_v2-idle-Sheet.png",
				48,
				16,
				animatedsprite.Loop,
				100,
				"",
			),
			"Open": animatedsprite.NewClip(
				"sprites/environment/door_v2-open-Sheet.png",
				48,
				16,
				animatedsprite.Once,
				100,
				"",
			),
			"Close": animatedsprite.NewClip(
				"sprites/environment/door_v2-close-Sheet.png", 48, 16,
				animatedsprite.Once,
				100,
				"",
			),
		}, renderer.RenderTarget{
			Type: renderer.TEXTURE,
			Name: "LevelTextureRaw",
		}, 5, 0.5, 0.5, "Idle",
	)

	doorAnimNode := doorNode.AddChild(doorAnim, "Sprite", engine.MakeOnTreeAdd(doorAnim, cmd))

	transform, ok := engine.As[*transform2D.Transform2D](doorAnimNode.GetValue())
	if ok {
		transform.SetAngle(maths.DirToRadians(direction))
	}

	doorV2Actor.AnimatedSprite = doorAnim

	triggerField := utils.Must(entity.GetFieldByName("InteractRegion"))
	triggerEntityIid := ebitenLDTK.As[ebitenLDTK.EntityRef](triggerField).EntityIid
	triggerEntity := utils.Must(level.GetEntityByIid(triggerEntityIid))

	relPosX := triggerEntity.Px[0] - entity.Px[0]
	relPosY := triggerEntity.Px[1] - entity.Px[1]
	triggerActor := trigger.NewTrigger(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
				transform2D.WithPos(relPosX, relPosY),
			),
		),
		trigger.WithRect(maths.NewRect(triggerEntity.Px[0], triggerEntity.Px[1], triggerEntity.Width, triggerEntity.Height)),
		trigger.WithName(fmt.Sprintf("Door-%s", triggerEntityIid)),
	)

	doorNode.AddChild(triggerActor, "Trigger", engine.MakeOnTreeAdd(triggerActor, cmd))

	doorV2Actor.Trigger = triggerActor
}

func SpawnPlatform(cmd *commands.Commands, scene *engine.Scene, entity *ebitenLDTK.Entity, envParentNode *engine.Node) {
	platformActor := platform.NewPlatform(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		), entity,
	)
	envParentNode.AddChild(platformActor, "Platform", engine.MakeOnTreeAdd(platformActor, cmd))
}

func SpawnHazard(cmd *commands.Commands, scene *engine.Scene, entity *ebitenLDTK.Entity, envParentNode *engine.Node) {
	hazardActor := hazard.NewHazard(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		), entity,
	)

	envParentNode.AddChild(hazardActor, "Hazard", engine.MakeOnTreeAdd(hazardActor, cmd))
}

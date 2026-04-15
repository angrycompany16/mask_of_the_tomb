package bundles

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/backend/maths"
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
	"mask_of_the_tomb/internal/game/actors/doorv2"
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
		LDTKData, ok := assetloader.GetAsset[assettypes.LDTKData](cmd.AssetLoader, "LDTK/world.ldtk")
		if !ok {
			panic("Grusom død")
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

		intGridCSV := playerspace.ExtractLayerCSV([]int{spikeIntGridID})
		scene.SpawnActor("SlamboxTilemap", slamboxtilemap.NewSlamboxTilemap(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			intGridCSV,
			int(playerspace.GridSize),
		), cmd)

		for i := range level.Layers {
			layer := level.Layers[i]

			if layer.Type == ebitenLDTK.LayerTypeEntities {
				continue
			}

			tileset := utils.Must(defs.GetTilesetByUid(layer.TilesetUid))
			tileSize := tileset.TileGridSize
			tilesetImg := tilesetMap[tileset.Name]

			// TODO: Find a way to render the tile layers to texture instead
			// of to screen. Then process them with shaders.
			scene.SpawnActor(layer.Name, ldtktilelayer.NewLDTKTilemapLayer(
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
				false,
			), cmd)

			// // Spawn the shader as a child
			// scene.AddChild("Shader",
			// 	shaderactor.NewShaderEffect(
			// 		graphic.NewGraphic(
			// 			transform2D.NewTransform2D(
			// 				nodeactor.NewNode(),
			// 			),
			// 		),
			// 	),
			// 	node, cmd)
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
			vectorgraphic.WithLayer("Background"),
			vectorgraphic.WithDrawOrder(-len(level.Layers)-1),
			vectorgraphic.WithImage(int(level.PxWid), int(level.PxHei)),
			vectorgraphic.WithPivot(0, 0),
		), cmd)

		entityLayer := utils.Must(level.GetLayerByName("Entities"))
		for _, entity := range entityLayer.Entities {
			switch entity.Name {
			// case names.HazardEntity:
			// 	newLevel.hazards = append(newLevel.hazards, entities.NewHazard(&entity))
			case "Slambox":
				SpawnSlambox(cmd, scene, &entity)
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
				SpawnDoor(cmd, scene, &entity, &level)
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

		// Spawn the appropriate particle system
		scene.SpawnActor("BackgroundParticles", particles.NewParticleSystem(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			), make([]*particles.Burst, 0), true, 0.5, 0, maths.RandomFloat64{0, 480}, maths.RandomFloat64{0, 270}, maths.RandomFloat64{0, 0}, maths.RandomFloat64{-5, 0}, maths.RandomFloat64{0, 4}, maths.RandomFloat64{0, 0.5}, maths.RandomFloat64{0, 0.01}, maths.RandomFloat64{0.7, 1.5}, maths.RandomFloat64{0.0, 0.0}, maths.RandomFloat64{3.0, 5.0}, maths.RandomFloat64{0.0, 2.0}, maths.RandomFloat64{0.0, 1.0}, [4]uint8{255, 255, 255, 255}, [4]uint8{255, 255, 255, 255}, 13, 13, "sprites/environment/star.png", "Background", 0,
		), cmd)
	}
}

func SpawnSlambox(cmd *commands.Commands, scene *engine.Scene, entity *ebitenLDTK.Entity) {
	// Spawn slambox bundle
	slamboxNode := scene.SpawnActor("Slambox",
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

	scene.AddChild("Sprite", autotilesprite.NewAutoTileSprite(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		),
		autotilesprite.WithSize(entity.Width, entity.Height),
		autotilesprite.WithTilemap("sprites/environment/slambox_tilemap.png"),
	), slamboxNode, cmd)
}

func SpawnDoor(cmd *commands.Commands, scene *engine.Scene, entity *ebitenLDTK.Entity, level *ebitenLDTK.Level) {
	directionField := utils.Must(entity.GetFieldByName("Direction"))
	direction := maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

	doorNode := scene.SpawnActor("Door", doorv2.NewDoorV2(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
				transform2D.WithPos(entity.Px[0], entity.Px[1]),
			),
		), entity, level,
	), cmd)

	doorV2Actor, ok := engine.As[*doorv2.DoorV2](doorNode.GetValue())

	doorSprite := scene.AddChild("Sprite", animatedsprite.NewAnimatedSprite(
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
		}, "Playerspace", 5, 0.5, 0.5, "Idle",
	), doorNode, cmd)

	transform, ok := engine.As[*transform2D.Transform2D](doorSprite.GetValue())
	if ok {
		transform.SetAngle(maths.DirToRadians(direction))
		doorV2Actor.SpriteTransform = transform
	}

	triggerField := utils.Must(entity.GetFieldByName("InteractRegion"))
	triggerEntityIid := ebitenLDTK.As[ebitenLDTK.EntityRef](triggerField).EntityIid
	triggerEntity := utils.Must(level.GetEntityByIid(triggerEntityIid))

	relPosX := triggerEntity.Px[0] - entity.Px[0]
	relPosY := triggerEntity.Px[1] - entity.Px[1]
	triggerNode := scene.AddChild("Trigger", trigger.NewTrigger(
		graphic.NewGraphic(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
				transform2D.WithPos(relPosX, relPosY),
			),
		),
		trigger.WithRect(maths.NewRect(triggerEntity.Px[0], triggerEntity.Px[1], triggerEntity.Width, triggerEntity.Height)),
		trigger.WithName(fmt.Sprintf("Door-%s", triggerEntityIid)),
	), doorNode, cmd)

	doorV2Actor.Trigger, ok = engine.As[*trigger.Trigger](triggerNode.GetValue())
}

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

func MakeLDTKLevelBundle(levelName string) engine.Bundle {
	return func(cmd *engine.Commands, scene *engine.Scene) {
		LDTKData, ok := assetloader.GetAsset[assettypes.LDTKData](cmd.AssetLoader(), "LDTK/world.ldtk")
		if !ok {
			panic("Grusom død")
		}

		world := LDTKData.Value().World
		tilesetMap := LDTKData.Value().Tilesets
		defs := world.Defs
		level := utils.Must(world.GetLevelByName(levelName))

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
		// fmt.Println(cmd.Scene())
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
			), cmd)
		}

		scene.SpawnActor("BackgroundColor", vectorgraphic.NewVectorGraphic(
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
		), cmd)

		entityLayer := utils.Must(level.GetLayerByName("Entities"))
		for _, entity := range entityLayer.Entities {
			switch entity.Name {
			// case names.HazardEntity:
			// 	newLevel.hazards = append(newLevel.hazards, entities.NewHazard(&entity))
			case "Slambox":
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

				// This really is just horrible...
				// The quickest solution i could possibly find
				// cmd.AssetLoader().LoadAll()
				// slamboxNode.GetValue().Init(cmd)
				// autotilesprite.GetValue().Init(cmd)
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
				directionField := utils.Must(entity.GetFieldByName("Direction"))
				direction := maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

				doorNode := scene.SpawnActor("Door", doorv2.NewDoorV2(
					graphic.NewGraphic(
						transform2D.NewTransform2D(
							nodeactor.NewNode(),
							transform2D.WithPos(entity.Px[0], entity.Px[1]),
						),
					), &entity, &level,
				), cmd)

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

				transform, ok := engine.GetActor[*transform2D.Transform2D](doorSprite.GetValue())
				if ok {
					transform.SetAngle(maths.DirToRadians(direction))
				}

				scene.AddChild("Trigger", trigger.NewTrigger(
					transform2D.NewTransform2D(
						nodeactor.NewNode(),
					),
				), doorNode, cmd)

				// cmd.AssetLoader().LoadAll()
				// doorNode.GetValue().Init(cmd)
				// doorSprite.GetValue().Init(cmd)
				// doorTrigger.GetValue().Init(cmd)
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

		// scene.SpawnActor("LDTKWorld",
		// 	ldtkworld.NewLDTKLevel(
		// 		graphic.NewGraphic(
		// 			transform2D.NewTransform2D(
		// 				nodeactor.NewNode(),
		// 			),
		// 		), "Level_3", "LDTK/world.ldtk",
		// 	), cmd,
		// )
	}
}

package bundles

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	sound_v2 "mask_of_the_tomb/internal/backend/sound"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/particles"
	"mask_of_the_tomb/internal/engine/actors/sound"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/actors/vectorgraphic"
	"mask_of_the_tomb/internal/engine/commands"
	ldtktilelayer "mask_of_the_tomb/internal/game/actors/LDTKTileLayer"
	"mask_of_the_tomb/internal/game/actors/backgroundshader"
	"mask_of_the_tomb/internal/game/actors/levelshader"
	"mask_of_the_tomb/internal/game/actors/shaderactor"
	"mask_of_the_tomb/internal/game/actors/slamboxtilemap"
	"mask_of_the_tomb/internal/game/actors/sounddebug"
	"mask_of_the_tomb/internal/game/globaldata"
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

var shaderMap = map[string]string{
	"Basement":                     "shaders/basement_background.kage",
	"Library":                      "shaders/basement_background.kage",
	"Grasslands":                   "shaders/grass_background.kage",
	"Strange_dark_blue_palm_trees": "shaders/basement_background.kage",
	"Royal_palace":                 "shaders/basement_background.kage",
	"Purple_rain":                  "shaders/basement_background.kage",
}

var songMap = map[string]string{
	"Basement":                     "music/basement.ogg",
	"Library":                      "music/library.ogg",
	"Grasslands":                   "music/grasslands.ogg",
	"Strange_dark_blue_palm_trees": "music/strange_dark_blue_palm_trees.ogg",
	"Royal_palace":                 "music/royal_palace.mp3",
	"Purple_rain":                  "music/purple_rain.mp3",
}

func MakeLDTKLevelBundle(levelIid string) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		// gw, gh := cmd.Renderer.GetGameSize()

		// 1. Load and prepare data from LDTK
		LDTKData, ok := assetloader.GetAsset[assettypes.LDTKData](cmd.AssetLoader, "LDTK/world.ldtk")
		if !ok {
			fmt.Println("Unable to load LDTK world asset from assetloader when making level bundle. Returning.")
			return scene.SpawnActor("LDTK-FAILURE", transform2D.NewTransform2D(), cmd)
		}

		world := LDTKData.Value().World
		tilesetMap := LDTKData.Value().Tilesets
		defs := world.Defs
		level := utils.Must(world.GetLevelByIid(levelIid))
		biomeField := utils.Must(level.GetFieldByName("Biome"))
		biome := ebitenLDTK.As[ebitenLDTK.Enum](biomeField).Value

		globaldata, _ := commands.Get[globaldata.GlobalData](cmd)
		if biome != globaldata.Temp.SceneSwitch.PreviousBiome {
			// Stop playing the current song (This is VERY UGLY!)
			sound_v2.StopSound(songMap[globaldata.Temp.SceneSwitch.PreviousBiome])
		}

		playerspace, err := level.GetLayerByName("Playerspace")
		if err != nil {
			fmt.Println("Error when loading level:", err)
			return scene.SpawnActor("LDTK-FAILURE", transform2D.NewTransform2D(), cmd)
		}

		var spikeIntGridID int
		for _, layerDef := range defs.LayerDefs {
			if layerDef.Name == "Playerspace" {
				spikeIntGridID = layerDef.GetIntGridValue("Spikes")
			}
		}

		// 2. Spawn nodes for tilemap layers, including spikes
		envParentNode := scene.SpawnActor("Environment", transform2D.NewTransform2D(), cmd)

		intGridCSV := playerspace.ExtractLayerCSV([]int{spikeIntGridID})
		slamboxTilemapActor := slamboxtilemap.NewSlamboxTilemap(
			graphic.NewGraphic(),
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
				graphic.NewGraphic(),
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
			var entityBundle engine.Bundle
			hasBundle := true

			switch entity.Name {
			case "Hazard":
				entityBundle = makeHazardBundle(&entity)
			case "Slambox":
				mainRect := maths.NewRect(entity.Px[0], entity.Px[1], entity.Width, entity.Height)
				subSlamboxesField, _ := entity.GetFieldByName("SubSlamboxes")
				subSlamboxes := ebitenLDTK.AsArray[ebitenLDTK.EntityRef](subSlamboxesField)
				subrects := make([]*maths.Rect, len(subSlamboxes))

				for i, entityRef := range subSlamboxes {
					entity, _ := level.GetEntityByIid(entityRef.EntityIid)
					subrects[i] = maths.NewRect(entity.Px[0], entity.Px[1], entity.Width, entity.Height)
				}

				entityBundle = makeSlamboxGroupBundle(mainRect, subrects)
			case "Grass":
				entityBundle = makeGrassBundle(&entity, &level)
			// case names.TurretEntity:
			// 	newLevel.turrets = append(newLevel.turrets, entities.NewTurret(&entity, entityLayer.GridSize))
			// case names.CatcherEntity:
			// 	newLevel.catchers = append(newLevel.catchers, entities.NewCatcher(&entity))
			case "Platform":
				entityBundle = MakePlatformBundle(&entity)
			// case names.LanternEntity:
			// 	newLevel.lanterns = append(newLevel.lanterns, entities.NewLantern(&entity, entityLayer.GridSize))
			case "Door":
				entityBundle = makeDoorBundle(&entity, &level)
				// case chainNodeEntityName:
				// 	newLevel.chainNodes = append(newLevel.chainNodes, entities.NewChainNode(&entity))
			case "DoorKey":
				entityBundle = makeKeyBundle(&entity)
			case "NPCSpawn":
				NPCField, _ := entity.GetFieldByName("NPC")
				NPCName := ebitenLDTK.As[ebitenLDTK.Enum](NPCField).Value
				switch NPCName {
				case "Lemma":
					entityBundle = makeLemmaBundle(&entity)
				case "Grok":
					entityBundle = makeGrokBundle(&entity)
				}
			default:
				hasBundle = false
			}

			if hasBundle {
				scene.SetParent(scene.SpawnBundleV2(cmd, entityBundle), envParentNode)
			}
		}

		// 4. Spawn remaining actors (particlesystems, shaders, etc...)
		scene.SpawnActor("SoundDebug", sounddebug.CreateSoundDebug(), cmd)

		scene.SpawnActor("MusicPlayer", sound.NewSoundPlayer(
			sound.WithSoundData(songMap[biome], true, songMap[biome]),
			sound.WithAutoPlay(true),
		), cmd)

		scene.SpawnActor("BackgroundParticles", particles.NewParticleSystem(
			particles.WithBursts(),
			particles.WithGlobalSpace(true),
			particles.WithEmission(0.5),
			particles.WithSpawnPos(0, 480, 0, 270),
			particles.WithSpawnVel(0, 0, -5, 0),
			particles.WithSpawnAngle(0, 4),
			particles.WithSpawnAngularVel(0, 0.5),
			particles.WithAirFriction(0, 0.01),
			particles.WithScale(0.7, 1.5, 0.0, 0.0),
			particles.WithLifetime(3, 5),
			particles.WithNoiseFactor(0, 2, 0, 1),
			particles.WithColors(
				[4]uint8{255, 255, 255, 255},
				[4]uint8{255, 255, 255, 255},
				[4]uint8{255, 255, 255, 255},
				[4]uint8{255, 255, 255, 255},
			),
			particles.WithImageSize(13, 13),
			particles.WithSprite("sprites/environment/star.png"),
			particles.WithRenderInfo(
				renderer.RenderTarget{Type: renderer.SCREEN, Name: "Background"},
				1,
			),
		), cmd)

		scene.SpawnActor("BackgroundShader",
			backgroundshader.NewBackgroundShader(
				shaderactor.NewShader(
					graphic.NewGraphic(), shaderMap[biome], cmd.Renderer.Textures["BackgroundRaw"], "Background", 0,
				),
			), cmd)

		scene.SpawnActor("LevelShader",
			levelshader.NewLevelShader(
				shaderactor.NewShader(
					graphic.NewGraphic(), "shaders/pixel_lights.kage", cmd.Renderer.Textures["LevelTextureRaw"], "Playerspace", 0,
				),
			), cmd)

		return envParentNode
	}
}

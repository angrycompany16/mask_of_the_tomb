package game

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/backend/triggerenv"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/scenes"
	"mask_of_the_tomb/internal/game/sceneswitch"

	"github.com/hajimehoshi/ebiten/v2"
)

func CreateGame(gw, gh, ps int) *engine.Game {
	cmd := commands.NewCommands(
		renderer.NewRenderer(gw, gh, ps, true, true),
		assetloader.NewAssetLoader(&assets.FS),
		input.NewInputHandler(),
	)

	commands.Set[triggerenv.TriggerEnv](cmd, triggerenv.NewTriggerEnv())
	commands.Set[slambox.SlamboxEnvironment](cmd, slambox.NewSlamboxEnvironment(8))
	commands.Set[sceneswitch.SceneSwitch](cmd, &sceneswitch.SceneSwitch{"", maths.DirUp})

	cmd.Renderer.Textures["ForegroundRaw"] = ebiten.NewImage(gw, gh)
	cmd.Renderer.Textures["LevelTextureRaw"] = ebiten.NewImage(gw, gh)
	cmd.Renderer.Textures["BackgroundRaw"] = ebiten.NewImage(gw, gh)

	game := engine.NewGame(cmd)

	ldtkDataRef := assetloader.StageAsset[assettypes.LDTKData](
		cmd.AssetLoader,
		"LDTK/world.ldtk",
		assettypes.NewLDTKAsset(
			"LDTK/world.ldtk",
		),
	)

	cmd.AssetLoader.LoadAll()

	sceneManager, _ := commands.Get[engine.SceneManager](cmd)
	LDTKWorld := ldtkDataRef.Value().World
	for _, level := range LDTKWorld.Levels {
		sceneManager.RegisterScene(level.Iid, scenes.MakeGamePlayeScene(level.Iid))
	}

	spawnScene, _ := LDTKWorld.GetLevelByName("Level_8")
	sceneManager.SpawnScene(spawnScene.Iid, cmd)
	return game
}

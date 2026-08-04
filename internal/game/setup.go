package game

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/backend/triggerenv"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/globaldata"
	"mask_of_the_tomb/internal/game/scenes"
	"mask_of_the_tomb/motb_assets/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

func CreateGame(gw, gh, ps int, bypassSave bool) *engine.Game {
	inputhandler := input.NewInputHandler()
	inputhandler.InputSchemes["PlayerControls"] = input.NewInputScheme()
	inputhandler.InputSchemes["UIControls"] = input.NewInputScheme()
	inputhandler.InputSchemes["EngineControls"] = input.NewInputScheme()
	cmd := commands.NewCommands(
		renderer.NewRenderer(gw, gh, ps, true, true),
		assetloader.NewAssetLoader(&assets.FS),
		inputhandler,
	)

	commands.Set[triggerenv.TriggerEnv](cmd, triggerenv.NewTriggerEnv())
	commands.Set[slambox.SlamboxEnvironment](cmd, slambox.NewSlamboxEnvironment(8))
	commands.Set[globaldata.GlobalData](cmd, globaldata.NewGlobalData(bypassSave))

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

	if !bypassSave {
		globaldata, _ := commands.Get[globaldata.GlobalData](cmd)
		globaldata.Persist.Load(0, true)
	}

	sceneManager, _ := commands.Get[engine.SceneManager](cmd)
	LDTKWorld := ldtkDataRef.Value().World
	for _, level := range LDTKWorld.Levels {
		sceneManager.RegisterScene(level.Iid, scenes.MakeGamePlayScene(0, 0, level.Iid))
	}

	sceneManager.RegisterScene("MainMenu", scenes.MakeMainMenuScene())
	sceneManager.RegisterScene("OptionsMenu", scenes.MakeOptionsScene())
	sceneManager.RegisterScene("SaveMenu", scenes.MakeSaveMenuScene())

	sceneManager.SpawnScene("MainMenu", cmd)
	return game
}

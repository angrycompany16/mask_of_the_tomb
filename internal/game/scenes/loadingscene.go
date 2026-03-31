package scenes

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/enginebundles"
	"mask_of_the_tomb/internal/game/bundles"
)

// Typical scene layout
// Load any important assets
// Create node tree and stuff

func LoadingScene(cmd *engine.Commands) *engine.Scene {
	scene := engine.NewScene("loadingScene", nodeactor.NewNode(), cmd)

	assetloader.StageAsset[assettypes.LDTKAsset](
		cmd.AssetLoader(),
		"LDTK/world.ldtk",
		assettypes.NewLDTKAsset(
			"LDTK/world.ldtk",
		),
	)

	cmd.AssetLoader().LoadAll()

	gameWidth, gameHeigth := cmd.Renderer().GetGameSize()
	pixelScale := cmd.Renderer().GetPixelScale()
	scene.SpawnBundle(cmd, enginebundles.MakeDefaultBundle(gameWidth, gameHeigth, pixelScale))
	scene.SpawnBundle(cmd, bundles.MakeLDTKLevelBundle("Level_3"))
	scene.SpawnBundle(cmd, bundles.MakePlayerBundle(0, 0, 16, 16))
	return scene
}

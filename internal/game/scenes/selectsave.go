package scenes

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/engine/enginebundles"
	"mask_of_the_tomb/internal/game/bundles"
)

func MakeSaveMenuScene() engine.SceneBuilder {
	return func(cmd *commands.Commands) *engine.Scene {
		scene := engine.NewScene("meaninglessName", nodeactor.NewNode(), cmd)

		gameWidth, gameHeigth := cmd.Renderer.GetGameSize()
		pixelScale := cmd.Renderer.GetPixelScale()
		scene.SpawnBundleV2(cmd, enginebundles.MakeDefaultBundle(gameWidth, gameHeigth, pixelScale))
		scene.SpawnBundleV2(cmd, bundles.MakeSaveMenuBundle())
		return scene
	}
}

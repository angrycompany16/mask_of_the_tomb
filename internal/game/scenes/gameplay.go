package scenes

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/engine/enginebundles"
	"mask_of_the_tomb/internal/game/actors/pausemenuspawner"
	"mask_of_the_tomb/internal/game/bundles"
)

func MakeGamePlayeScene(playerX, playerY float64, levelIid string) engine.SceneBuilder {
	return func(cmd *commands.Commands) *engine.Scene {
		scene := engine.NewScene("loadingScene", nodeactor.NewNode(), cmd)

		gameWidth, gameHeigth := cmd.Renderer.GetGameSize()
		pixelScale := cmd.Renderer.GetPixelScale()
		scene.SpawnBundle(cmd, enginebundles.MakeDefaultBundle(gameWidth, gameHeigth, pixelScale))
		scene.SpawnBundle(cmd, bundles.MakeLDTKLevelBundle(levelIid))
		scene.SpawnBundle(cmd, bundles.MakePlayerBundle(playerX, playerY, 16, 16))
		scene.SpawnActor("PauseMenuTrigger", pausemenuspawner.NewPauseMenuSpawner(
			nodeactor.NewNode(),
		), cmd)
		return scene
	}
}

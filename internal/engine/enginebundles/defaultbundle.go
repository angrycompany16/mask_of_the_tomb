package enginebundles

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/assetviewer"
	"mask_of_the_tomb/internal/engine/actors/camera"
	"mask_of_the_tomb/internal/engine/actors/inspector"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
)

func MakeDefaultBundle(gameWidth, gameHeight, pixelScale float64) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) {
		scene.SpawnActor("Camera", camera.NewCamera(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
			camera.WithSize(gameWidth, gameHeight),
		), cmd)

		scene.SpawnActor("Inspector", inspector.NewInspector(
			nodeactor.NewNode(),
			inspector.WithSize(int(gameWidth*pixelScale/3), int(gameHeight*pixelScale*0.8)),
		), cmd)

		scene.SpawnActor("AssetViewer", assetviewer.NewAssetViewer(
			nodeactor.NewNode(),
		), cmd)

	}
}

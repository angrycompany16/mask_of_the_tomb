package enginebundles

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/assetviewer"
	"mask_of_the_tomb/internal/engine/actors/camera"
	"mask_of_the_tomb/internal/engine/actors/inspector"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
)

// Bundle that spawns a camera, an inspector and an asset viewer.
func MakeDefaultBundle(gameWidth, gameHeight, pixelScale float64) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		defaultBundleRoot := scene.SpawnActor("defaultBundle", transform2D.NewTransform2D(), cmd)

		cameraActor := camera.NewCamera(camera.WithSize(gameWidth, gameHeight))
		defaultBundleRoot.AddChild(cameraActor, "Camera", engine.MakeOnTreeAdd(cameraActor, cmd))

		inspectorActor := inspector.NewInspector(inspector.WithSize(int(gameWidth*pixelScale/3), int(gameHeight*pixelScale*0.8)))
		defaultBundleRoot.AddChild(inspectorActor, "Inspector", engine.MakeOnTreeAdd(inspectorActor, cmd))

		assetViewerActor := assetviewer.NewAssetViewer()
		defaultBundleRoot.AddChild(assetViewerActor, "AssetViewer", engine.MakeOnTreeAdd(assetViewerActor, cmd))

		return defaultBundleRoot
	}
}

package scenes

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/inspector"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	ldtkworld "mask_of_the_tomb/internal/game/actors/LDTKworld"
)

func LoadingScene(servers *engine.Commands) *engine.Scene {
	scene := engine.NewScene("loadingScene", nodeactor.NewNode(), servers)

	scene.SpawnActor("Inspector", inspector.NewInspector(0, 0, 300, 400), servers)
	scene.SpawnActor("LDTKWorld",
		ldtkworld.NewLDTKLevel(
			*transform2D.NewTransform2D(
				*nodeactor.NewNode(),
			), "Level_3", "LDTK/world.ldtk",
		), servers,
	)

	return scene
}

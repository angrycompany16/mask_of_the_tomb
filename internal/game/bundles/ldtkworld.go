package bundles

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	ldtkworld "mask_of_the_tomb/internal/game/actors/LDTKworld"
)

func MakeLDTKWorldBundle() engine.Bundle {
	return func(cmd *engine.Commands, scene *engine.Scene) {
		scene.SpawnActor("LDTKWorld",
			ldtkworld.NewLDTKLevel(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				), "Level_3", "LDTK/world.ldtk",
			), cmd,
		)
	}
}

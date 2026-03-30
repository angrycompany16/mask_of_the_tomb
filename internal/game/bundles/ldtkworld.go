package bundles

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	ldtkworld "mask_of_the_tomb/internal/game/actors/LDTKworld"
)

// This is dumb. All the LDTKWorld actor does is spawn children and act
// as a parent, it doesn't really actually have any functionality?
// This is apparent if one takes a look at the actor object itself.
// It doesn't even have an update method...
func MakeLDTKWorldBundle() engine.Bundle {
	return func(cmd *engine.Commands, scene *engine.Scene) {
		scene.SpawnActor("LDTKWorld",
			ldtkworld.NewLDTKLevel(
				graphic.NewGraphic(
					transform2D.NewTransform2D(
						nodeactor.NewNode(),
					),
				), "Level_3", "LDTK/world.ldtk",
			), cmd,
		)
	}
}

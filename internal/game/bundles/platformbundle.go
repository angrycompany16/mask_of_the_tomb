package bundles

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/platform"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

func MakePlatformBundle(entity *ebitenLDTK.Entity) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		return scene.SpawnActor("Platform", platform.NewPlatform(graphic.NewGraphic(), entity), cmd)
	}
}

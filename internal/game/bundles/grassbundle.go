package bundles

import (
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/grass"
	"mask_of_the_tomb/internal/game/globaldata"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

func makeGrassBundle(entity *ebitenLDTK.Entity, level *ebitenLDTK.Level) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		playerspace := utils.Must(level.GetLayerByName("Playerspace"))
		globaldata, _ := commands.Get[globaldata.GlobalData](cmd)

		grassActor := grass.NewGrass(
			graphic.NewGraphic(
				graphic.WithTransform(
					transform2D.NewTransform2D(
						transform2D.WithPos(entity.Px[0], entity.Px[1]),
					),
				),
			),
			entity,
			playerspace.GridSize,
			"sprites/environment/grass.png",
			globaldata.Temp.GrassWindSeed,
			renderer.TextureTarget("LevelTextureRaw"),
			-10,
		)

		return scene.SpawnActor("Grass", grassActor, cmd)
	}
}

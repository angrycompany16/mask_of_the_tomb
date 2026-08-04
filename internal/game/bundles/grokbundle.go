package bundles

import (
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

func makeGrokBundle(entity *ebitenLDTK.Entity) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		return scene.SpawnActor("Grok", animatedsprite.NewAnimatedSprite(
			graphic.NewGraphic(
				graphic.WithTransform(
					transform2D.NewTransform2D(
						transform2D.WithPos(entity.Px[0], entity.Px[1]),
					),
				),
			),
			map[string]*animatedsprite.Clip{
				"Idle": animatedsprite.NewClip("sprites/NPC/grok.png", 38, 32, animatedsprite.Loop, 100, ""),
			},
			renderer.TextureTarget("LevelTextureRaw"), 1, 0.5, 0.5, "Idle",
		), cmd)
	}
}

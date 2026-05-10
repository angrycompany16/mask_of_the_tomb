package player

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/particles"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
)

// TODO: Need some autodestroy functionality...
// However, it's not causing any lag right now, fortunately.
func MakeJumpParticlesBundle(x, y float64, dir maths.Direction, halfSize float64) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) {
		dx, dy := maths.VectorFromDir(dir)
		x -= dx * halfSize
		y -= dy * halfSize
		rot := maths.DirToRadians(dir)
		particlesBroad := scene.SpawnActor("JumpParticlesBroad", particles.NewParticleSystem(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
					transform2D.WithPos(x, y),
					transform2D.WithAngle(rot),
				),
			),
			particles.WithBursts(
				&particles.Burst{Count: 5, Time: 0},
			),
			particles.WithSpawnPos(-8, 8, 0, 0),
			particles.WithSpawnVel(0, 0, -40, -10),
			particles.WithColors(
				[4]uint8{255, 255, 255, 255},
				[4]uint8{200, 200, 200, 255},
				[4]uint8{200, 200, 200, 255},
				[4]uint8{255, 0, 0, 255},
			),
			particles.WithGlobalSpace(false),
			particles.WithImageSize(64, 64),
			particles.WithEmission(0),
			particles.WithLifetime(0.1, 0.5),
			particles.WithScale(0.2, 0.01, 0.0, 0.0),
			particles.WithSprite("sprites/icons/square-16x16.png"),
		), cmd)

		particlesTight := scene.SpawnActor("JumpParticlesTight", particles.NewParticleSystem(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
					transform2D.WithPos(x, y),
					transform2D.WithAngle(rot),
				),
			),
			particles.WithBursts(
				&particles.Burst{Count: 6, Time: 0},
			),
			particles.WithSpawnPos(-3, 3, 0, 0),
			particles.WithSpawnVel(0, 0, -100, -40),
			particles.WithColors(
				[4]uint8{255, 255, 255, 255},
				[4]uint8{200, 200, 200, 255},
				[4]uint8{200, 200, 200, 255},
				[4]uint8{255, 0, 0, 255},
			),
			particles.WithGlobalSpace(false),
			particles.WithImageSize(64, 64),
			particles.WithEmission(0),
			particles.WithLifetime(0.1, 0.5),
			particles.WithScale(0.2, 0.01, 0.0, 0.0),
			particles.WithSprite("sprites/icons/square-16x16.png"),
			particles.WithRenderInfo(
				renderer.RenderTarget{
					Type: renderer.TEXTURE,
					Name: "LevelTextureRaw",
				}, 0),
		), cmd)

		cmd.AssetLoader.LoadAll()

		// This is not the worst, but it's not ideal either
		// ...
		particlesBroad.Traverse(
			func(n *node.Node[engine.Actor]) {
				n.GetValue().Init(cmd)
			})

		particlesTight.Traverse(
			func(n *node.Node[engine.Actor]) {
				n.GetValue().Init(cmd)
			})
	}
}

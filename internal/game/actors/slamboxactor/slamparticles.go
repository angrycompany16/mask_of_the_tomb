package slamboxactor

import (
	"fmt"
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

func MakeSlamboxParticlesBundle(x, y float64, dir maths.Direction, halfWidth, halfHeight float64) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) {
		dx, dy := maths.VecFromDir(dir)
		x += dx * halfWidth
		y += dy * halfHeight
		size := 0.0

		if dir == maths.DirDown || dir == maths.DirUp {
			size = halfWidth
		} else if dir == maths.DirLeft || dir == maths.DirRight {
			size = halfHeight
		}

		fmt.Println(dir)
		rot := maths.DirToRadians(dir)
		slamboxParticles := scene.SpawnActor("SlamboxParticles", particles.NewParticleSystem(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
					transform2D.WithPos(x, y),
					transform2D.WithAngle(rot),
				),
			),

			particles.WithBursts(
				&particles.Burst{Count: 25, Time: 0},
			),
			particles.WithAirFriction(3, 5),
			particles.WithSpawnPos(-size, size, 0, 0),
			particles.WithSpawnVel(-20, 20, 2, 120),
			particles.WithColors(
				[4]uint8{255, 230, 70, 255},
				[4]uint8{255, 200, 50, 255},
				[4]uint8{200, 100, 0, 255},
				[4]uint8{100, 50, 0, 255},
			),
			particles.WithGlobalSpace(false),
			particles.WithImageSize(256, 64),
			particles.WithEmission(0),
			particles.WithLifetime(0.2, 0.7),
			particles.WithScale(0.15, 0.3, 0.0, 0.0),
			particles.WithSprite("sprites/icons/circle-64x64.png"),
			particles.WithRenderInfo(renderer.RenderTarget{
				Type: renderer.TEXTURE,
				Name: "LevelTextureRaw",
			}, 30),
		), cmd)

		cmd.AssetLoader.LoadAll()

		// This is not the worst, but it's not ideal either
		// ...
		slamboxParticles.Traverse(
			func(n *node.Node[engine.Actor]) {
				n.GetValue().Init(cmd)
			})
	}
}

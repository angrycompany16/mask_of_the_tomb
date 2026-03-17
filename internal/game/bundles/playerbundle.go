package bundles

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/actors/tracker"
)

func MakePlayerBundle() engine.Bundle {
	return func(cmd *engine.Commands, scene *engine.Scene) {
		playerNode := scene.SpawnActor("Player", player.NewPlayer(
			slamboxactor.NewSlambox(
				tracker.NewTracker(
					transform2D.NewTransform2D(
						nodeactor.NewNode(),
					), 5.0, 0.0, 0.0,
				), maths.NewRect(0, 0, 16, 16),
			),
		), cmd)

		scene.AddChild("PlayerSprite", sprite.NewSprite(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			), "Playerspace", "sprites/player/player.png",
		), playerNode, cmd)
	}
}

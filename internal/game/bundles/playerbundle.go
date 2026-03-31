package bundles

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/actors/tracker"
)

func MakePlayerBundle(playerX, playerY, playerWidth, playerHeight float64) engine.Bundle {
	return func(cmd *engine.Commands, scene *engine.Scene) {
		gw, gh := cmd.Renderer().GetGameSize()
		tlX, tlY := playerX+gw/2, playerY+gh/2
		playerNode := scene.SpawnActor("Player", player.NewPlayer(
			slamboxactor.NewSlambox(
				tracker.NewTracker(
					graphic.NewGraphic(
						transform2D.NewTransform2D(
							nodeactor.NewNode(),
							transform2D.WithPos(playerX, playerY),
						),
					), 10.0, tlX, tlY,
				),
				slamboxactor.WithPos(tlX, tlY),
				slamboxactor.WithSize(playerWidth, playerHeight),
			),
			0.1,
		), cmd)

		pivotNode := scene.AddChild("PlayerPivot", transform2D.NewTransform2D(
			nodeactor.NewNode(),
			transform2D.WithPos(playerWidth/2, playerHeight/2),
		), playerNode, cmd)

		scene.AddChild("PlayerSprite", animatedsprite.NewAnimatedSprite(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			map[string]*animatedsprite.Clip{
				player.IDLE_ANIM: animatedsprite.NewClip(
					"sprites/player/player-idle-Sheet.png",
					16,
					16,
					animatedsprite.Loop,
					140,
					"",
				),
				player.DASH_INIT_ANIM: animatedsprite.NewClip(
					"sprites/player/player-init-jump-Sheet.png",
					16,
					16,
					animatedsprite.Once,
					80,
					player.DASH_LOOP_ANIM,
				),
				player.DASH_LOOP_ANIM: animatedsprite.NewClip(
					"sprites/player/player-loop-jump-Sheet.png",
					16,
					16,
					animatedsprite.Loop,
					80,
					"",
				),
				player.SLAM_ANIM: animatedsprite.NewClip(
					"sprites/player/player-slam-Sheet.png",
					16,
					16,
					animatedsprite.Once,
					80,
					player.IDLE_ANIM,
				),
			}, "Playerspace", 6, 0.5, 0.5, player.IDLE_ANIM,
		), pivotNode, cmd)
	}
}

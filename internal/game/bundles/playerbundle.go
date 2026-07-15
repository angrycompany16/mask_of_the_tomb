package bundles

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/sound"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/actors/tracker"
	"mask_of_the_tomb/internal/game/actors/trigger"
)

func MakePlayerBundle(playerX, playerY, playerWidth, playerHeight float64) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) {
		gw, gh := cmd.Renderer.GetGameSize()
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
				slamboxactor.WithHasParticles(false),
			),
			0.1,
		), cmd)

		pivotActor := transform2D.NewTransform2D(
			nodeactor.NewNode(),
			transform2D.WithPos(playerWidth/2, playerHeight/2),
		)

		pivotNode := playerNode.AddChild(pivotActor, "PlayerPivot", engine.MakeOnTreeAdd(pivotActor, cmd))
		
		spriteActor := animatedsprite.NewAnimatedSprite(
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
			}, renderer.RenderTarget{
				Type: renderer.TEXTURE,
				Name: "LevelTextureRaw",
			}, 6, 0.5, 0.5, player.IDLE_ANIM,
		)

		pivotNode.AddChild(spriteActor, "PlayerSprite", engine.MakeOnTreeAdd(spriteActor, cmd))

		triggerActor := trigger.NewTrigger(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			trigger.WithRect(maths.NewRect(playerX, playerY, playerWidth, playerHeight)),
			trigger.WithName("Player"),
		)

		playerNode.AddChild(triggerActor, "PlayerTrigger", engine.MakeOnTreeAdd(triggerActor, cmd))

		playerActor, _ := engine.As[*player.Player](playerNode.GetValue())

		onMoveEv := playerActor.OnMove

		playerSound := sound.NewSoundPlayer(
			nodeactor.NewNode(),
			sound.WithSoundData("sfx/dash.wav", false, "dash"),
			sound.WithDspChannel("master"),
			sound.WithEventBus(onMoveEv),
		)

		playerNode.AddChild(playerSound, "DashSound", engine.MakeOnTreeAdd(playerSound, cmd));
	}
}

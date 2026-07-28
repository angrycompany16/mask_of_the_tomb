package bundles

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/missile"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/key"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/game/globaldata"
)

func MakeCollectedKeysBundle(playerX, playerY float64) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		keysRootNode := scene.SpawnActor("CollectedKeys", nodeactor.NewNode(), cmd)
		globaldata, _ := commands.Get[globaldata.GlobalData](cmd)
		playerNode, _ := scene.GetNodeByName("Player")
		playerActor, _ := engine.As[*player.Player](playerNode.GetValue())

		for _, keyData := range globaldata.Persist.Profile.Inventory.Keys {
			if keyData.Used {
				continue
			}

			missileActor := key.NewCollectedKey(
				missile.NewMissile(
					missile.WithTransform(
						transform2D.NewTransform2D(transform2D.WithPos(playerX, playerY)),
					),
					missile.WithSpeed(4),
					missile.WithCircularOffset(40, 1, maths.RandomRange(0, 10)),
					missile.WithTargetTransform(playerActor.Transform2D),
				),
			)
			missileNode := keysRootNode.AddChild(missileActor, "Key", engine.MakeOnTreeAdd(missileActor, cmd))

			sprite := sprite.NewSprite(
				renderer.TextureTarget("LevelTextureRaw"),
				"sprites/environment/key.png",
			)

			missileActor.Sprite = sprite
			missileNode.AddChild(sprite, "sprite", engine.MakeOnTreeAdd(sprite, cmd))
		}
		return keysRootNode
	}
}

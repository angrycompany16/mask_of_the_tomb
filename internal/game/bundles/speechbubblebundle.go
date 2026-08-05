package bundles

import (
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/speechbubble"
)

func makeSpeechBubbleBundle(anchor *transform2D.Transform2D, width, height float64) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		tickSprite := sprite.NewSprite(
			renderer.ScreenTarget("WorldUI"), "sprites/icons/textbox-tick.png",
			sprite.WithScaling(cmd.Renderer.GetPixelScale()),
			sprite.WithDrawOrder(10),
		)

		speechBubble := scene.SpawnActor("Speech bubble", speechbubble.NewSpeechBubble(
			anchor, tickSprite, width, height, 16,
		), cmd)

		scene.SpawnActor("TickSprite", tickSprite, cmd)

		return speechBubble
	}
}

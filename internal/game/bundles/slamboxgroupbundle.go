package bundles

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/sound"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/autotilesprite"
	"mask_of_the_tomb/internal/game/actors/slamboxgroup"
	"mask_of_the_tomb/internal/game/actors/tracker"
)

func makeSlamboxGroupBundle(mainRect *maths.Rect, subrects []*maths.Rect) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		slamboxActor := slamboxgroup.NewSlamboxGroup(
			tracker.NewTracker(
				graphic.NewGraphic(), 7.5, mainRect.X, mainRect.Y,
			),
			slamboxgroup.WithRects(mainRect, subrects),
		)

		slamboxNode := scene.SpawnActor("Slambox", slamboxActor, cmd)

		subrects_ := make([]*maths.Rect, len(subrects))
		for i := range subrects_ {
			r := subrects[i]
			subrects_[i] = maths.NewRect(r.X, r.Y, r.Width, r.Height)
		}

		autotileActor := autotilesprite.NewAutoTileSprite(
			graphic.NewGraphic(), renderer.RenderTarget{
				Type: renderer.TEXTURE,
				Name: "LevelTextureRaw",
			},
			autotilesprite.WithRects(maths.NewRect(mainRect.X, mainRect.Y, mainRect.Width, mainRect.Height), subrects_),
			autotilesprite.WithTilemap("sprites/environment/slambox_tilemap.png"),
		)
		slamboxNode.AddChild(autotileActor, "Sprite", engine.MakeOnTreeAdd(autotileActor, cmd))

		slamboxSound := sound.NewSoundPlayer(
			sound.WithSoundData("sfx/stone-crash-trimmed.wav", false, "Slambox-land"),
			sound.WithStartTriggers(slamboxActor.OnMoveFinishEv),
		)

		slamboxNode.AddChild(slamboxSound, "slamboxSound", engine.MakeOnTreeAdd(slamboxSound, cmd))

		return slamboxNode
	}
}

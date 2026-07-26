package bundles

import (
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/sound"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/key"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

func makeKeyBundle(entity *ebitenLDTK.Entity) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		keyActor := key.NewKey(entity)

		keyActorNode := scene.SpawnActor("Key", keyActor, cmd)
		//		keyActorNode := envParentNode.AddChild(keyActor, "Key", engine.MakeOnTreeAdd(keyActor, cmd))

		sprite := sprite.NewSprite(
			renderer.TextureTarget("LevelTextureRaw"),
			"sprites/environment/key.png",
			sprite.WithPivot(0, 0),
		)

		keyActorNode.AddChild(sprite, "sprite", engine.MakeOnTreeAdd(sprite, cmd))

		collectSound := sound.NewSoundPlayer(
			sound.WithSoundData("sfx/key-collect.ogg", false, "key-collect"),
			sound.WithStartTriggers(keyActor.OnPickupEv),
			sound.WithVolume(0.6),
		)

		keyActorNode.AddChild(collectSound, "collectSound", engine.MakeOnTreeAdd(collectSound, cmd))

		return keyActorNode

	}
}

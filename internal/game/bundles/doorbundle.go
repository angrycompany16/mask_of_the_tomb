package bundles

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/sound"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/door"
	"mask_of_the_tomb/internal/game/actors/trigger"
	"mask_of_the_tomb/internal/game/sceneswitch"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

func makeDoorBundle(entity *ebitenLDTK.Entity, level *ebitenLDTK.Level) engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		directionField := utils.Must(entity.GetFieldByName("Direction"))
		direction := maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

		doorActor := door.NewDoor(
			graphic.NewGraphic(
				graphic.WithTransform(
					transform2D.NewTransform2D(
						transform2D.WithPos(entity.Px[0], entity.Px[1]),
					),
				),
			), entity, level,
		)
		doorNode := scene.SpawnActor("Door", doorActor, cmd)

		doorAnim := animatedsprite.NewAnimatedSprite(
			graphic.NewGraphic(
				graphic.WithTransform(
					transform2D.NewTransform2D(
						transform2D.WithPos(entity.Width/2, entity.Height/2),
					),
				),
			),
			map[string]*animatedsprite.Clip{
				"Idle":  animatedsprite.NewClip("sprites/environment/door-idle-Sheet.png", 48, 16, animatedsprite.Loop, 100, ""),
				"Open":  animatedsprite.NewClip("sprites/environment/door-open-Sheet.png", 48, 16, animatedsprite.Once, 100, ""),
				"Close": animatedsprite.NewClip("sprites/environment/door-close-Sheet.png", 48, 16, animatedsprite.Once, 100, ""),
			}, renderer.TextureTarget("LevelTextureRaw"), 5, 0.5, 0.5, "Idle",
		)

		doorAnimNode := doorNode.AddChild(doorAnim, "Sprite", engine.MakeOnTreeAdd(doorAnim, cmd))

		transform, ok := engine.As[*transform2D.Transform2D](doorAnimNode.GetValue())
		if ok {
			transform.SetAngle(maths.DirToRadians(direction))
		}

		doorActor.AnimatedSprite = doorAnim

		lockAnim := animatedsprite.NewAnimatedSprite(
			graphic.NewGraphic(
				graphic.WithTransform(
					transform2D.NewTransform2D(
						transform2D.WithPos(entity.Width/2, entity.Height/2),
					),
				),
			),
			map[string]*animatedsprite.Clip{
				"Idle":      animatedsprite.NewClip("sprites/environment/lock-idle-Sheet.png", 40, 62, animatedsprite.Loop, 100, ""),
				"Unlock":    animatedsprite.NewClip("sprites/environment/lock-unlock-Sheet.png", 40, 62, animatedsprite.Once, 80, ""),
				"TryUnlock": animatedsprite.NewClip("sprites/environment/lock-tryunlock-Sheet.png", 40, 62, animatedsprite.Once, 80, ""),
			}, renderer.TextureTarget("LevelTextureRaw"), 10, 0.5, 0.5, "Idle",
		)

		doorNode.AddChild(lockAnim, "LockAnim", engine.MakeOnTreeAdd(lockAnim, cmd))
		doorActor.LockAnim = lockAnim

		triggerField := utils.Must(entity.GetFieldByName("InteractRegion"))
		triggerEntityIid := ebitenLDTK.As[ebitenLDTK.EntityRef](triggerField).EntityIid
		triggerEntity := utils.Must(level.GetEntityByIid(triggerEntityIid))

		relPosX := triggerEntity.Px[0] - entity.Px[0]
		relPosY := triggerEntity.Px[1] - entity.Px[1]
		triggerActor := trigger.NewTrigger(
			graphic.NewGraphic(
				graphic.WithTransform(
					transform2D.NewTransform2D(
						transform2D.WithPos(relPosX, relPosY),
					),
				),
			),
			trigger.WithRect(maths.NewRect(triggerEntity.Px[0], triggerEntity.Px[1], triggerEntity.Width, triggerEntity.Height)),
			trigger.WithName(fmt.Sprintf("Door-%s", triggerEntityIid)),
		)

		doorNode.AddChild(triggerActor, "Trigger", engine.MakeOnTreeAdd(triggerActor, cmd))

		doorActor.Trigger = triggerActor

		doorOpenSound := sound.NewSoundPlayer(
			sound.WithSoundData("sfx/door-open.ogg", false, "door-open"),
			sound.WithStartTriggers(doorActor.OnOpen),
		)

		doorNode.AddChild(doorOpenSound, "door-open", engine.MakeOnTreeAdd(doorOpenSound, cmd))

		sceneswitch, _ := commands.Get[sceneswitch.SceneSwitch](cmd)

		doorCloseSound := sound.NewSoundPlayer(
			sound.WithSoundData("sfx/door-close.ogg", false, "door-close"),
			sound.WithAutoPlay(sceneswitch.SpawnEntityIid == doorActor.EntityIid),
		)

		doorNode.AddChild(doorCloseSound, "door-close", engine.MakeOnTreeAdd(doorCloseSound, cmd))

		doorLockedSound := sound.NewSoundPlayer(
			sound.WithSoundData("sfx/door-locked.ogg", false, "door-locked"),
			sound.WithStartTriggers(doorActor.OnTryOpenLocked),
		)

		doorNode.AddChild(doorLockedSound, "door-try-unlock", engine.MakeOnTreeAdd(doorLockedSound, cmd))

		doorUnlockSound := sound.NewSoundPlayer(
			sound.WithSoundData("sfx/door-unlock.ogg", false, "door-unlock"),
			sound.WithStartTriggers(doorActor.OnUnlockEv),
		)

		doorNode.AddChild(doorUnlockSound, "door-unlock", engine.MakeOnTreeAdd(doorUnlockSound, cmd))

		return doorNode
	}
}

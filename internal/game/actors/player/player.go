package player

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/inputbuffer"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/backend/shaders"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/door"
	"mask_of_the_tomb/internal/game/actors/hazard"
	"mask_of_the_tomb/internal/game/actors/key"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/actors/slamboxgroup"
	"mask_of_the_tomb/internal/game/globaldata"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	IDLE_ANIM      = "Idle"
	DASH_INIT_ANIM = "Dash_init"
	DASH_LOOP_ANIM = "Dash_loop"
	SLAM_ANIM      = "Slam"
)

//go:generate stringer -type=PlayerState
type PlayerState int

const (
	IDLE PlayerState = iota
	MOVING
	SLAMMING
	DYING
	LEAVING
	ENTERING
)

type Player struct {
	*slamboxactor.Slambox
	State                     PlayerState
	Direction                 maths.Direction
	spriteTransform           *transform2D.Transform2D
	animatedSprite            *animatedsprite.AnimatedSprite
	pivotTransform            *transform2D.Transform2D
	inputbuffer               *inputbuffer.InputBuffer
	OnMoveFinish              *events.EventBus
	OnClipFinish              *events.EventBus
	OnMove                    *events.Event
	OnSlam                    *events.Event
	jumpOffset, jumpOffsetvel float64
	slamboxIDBuffer           int
	slamDirBuffer             maths.Direction
	hasSlammedBox             bool
	doorOffset                float64
	trueDoorOffset            float64
	doorY                     float64
	Light                     *shaders.Light
	doorOnOpen                *events.EventBus
	doorMemory                *door.Door
}

func (p *Player) Init(cmd *commands.Commands) {
	p.Slambox.Init(cmd)

	scene, _ := commands.Get[engine.Scene](cmd)
	globaldata_, _ := commands.Get[globaldata.GlobalData](cmd)

	savedDoor := (globaldata_.Persist.Profile.LastDoor != "")
	enteringFromRoom := (globaldata_.Temp.SceneSwitch.SpawnEntityIid != "")

	var spawnDoorIid string
	if savedDoor {
		spawnDoorIid = globaldata_.Persist.Profile.LastDoor
	} else {
		spawnDoorIid = globaldata_.Temp.SceneSwitch.SpawnEntityIid
	}
	spawnDoorNode, _ := scene.GetNodeFunc(
		func(n *node.Node[engine.Actor]) bool {
			doorActor, ok := engine.As[*door.Door](n.GetValue())
			if !ok {
				return false
			}
			return doorActor.EntityIid == spawnDoorIid
		},
	)

	hasGameInitPos := false
	if enteringFromRoom {
		doorActor, _ := engine.As[*door.Door](spawnDoorNode.GetValue())
		p.SetPos(doorActor.GetSpawnPos())
		//		p.Transform2D.SetPos(doorActor.GetSpawnPos())
		p.Transform2D.Propagate()
		p.State = ENTERING
		p.jumpOffset = -2 * p.GetRect().Height
		p.jumpOffsetvel = 4.5
		p.trueDoorOffset = 0
		p.doorOffset = 0
		p.Direction = globaldata_.Temp.SceneSwitch.SpawnDirection
	} else if savedDoor {
		doorActor, _ := engine.As[*door.Door](spawnDoorNode.GetValue())
		p.SetPos(doorActor.GetSpawnPos())
		p.Transform2D.Propagate()
		p.Direction = doorActor.Direction
	} else if hasGameInitPos {

	} else {
		fmt.Println("All methods for finding spawn pos failed. Picking first available door")

		spawnDoorNode, _ := scene.GetNodeFunc(
			func(n *node.Node[engine.Actor]) bool {
				_, ok := engine.As[*door.Door](n.GetValue())
				return ok
			},
		)

		doorActor, _ := engine.As[*door.Door](spawnDoorNode.GetValue())
		p.SetPos(doorActor.GetSpawnPos())
		p.Transform2D.SetPos(doorActor.GetSpawnPos())
		p.Transform2D.Propagate()
		p.Direction = doorActor.Direction
	}

	x, y := p.GetPos()
	globaldata_.Temp.LevelStates[scene.GetName()] = globaldata.NewLevelState(x, y, p.Direction)
	globaldata_.Temp.SaveLevelState(scene)

	playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
	playerControls.RegisterAction("moveLeft", input.KeyJustPressedAction(ebiten.KeyA))
	playerControls.AddBinding("moveLeft", input.KeyJustPressedAction(ebiten.KeyLeft))
	playerControls.RegisterAction("moveRight", input.KeyJustPressedAction(ebiten.KeyD))
	playerControls.AddBinding("moveRight", input.KeyJustPressedAction(ebiten.KeyRight))
	playerControls.RegisterAction("moveUp", input.KeyJustPressedAction(ebiten.KeyW))
	playerControls.AddBinding("moveUp", input.KeyJustPressedAction(ebiten.KeyUp))
	playerControls.RegisterAction("moveDown", input.KeyJustPressedAction(ebiten.KeyS))
	playerControls.AddBinding("moveDown", input.KeyJustPressedAction(ebiten.KeyDown))
	playerControls.AddBinding("Reset", input.KeyJustPressedAction(ebiten.KeyR))
	playerControls.AddBinding("Use", input.KeyJustPressedAction(ebiten.KeyE))

	// Would be very nice to set up a reference like this in another
	// way
	// But how? I guess we would have to link them together somehow
	// in the bundle
	childNode, enteringFromRoom := scene.GetNodeByName("PlayerSprite")
	p.spriteTransform, enteringFromRoom = engine.As[*transform2D.Transform2D](childNode.GetValue())
	p.animatedSprite, enteringFromRoom = engine.As[*animatedsprite.AnimatedSprite](childNode.GetValue())
	p.OnClipFinish = events.NewBusFrom(p.animatedSprite.OnClipFinished)

	pivotNode, enteringFromRoom := scene.GetNodeByName("PlayerPivot")
	p.pivotTransform, enteringFromRoom = engine.As[*transform2D.Transform2D](pivotNode.GetValue())

	p.pivotTransform.Propagate()
	x, y = p.pivotTransform.GetPos(false)
	p.Light.X = x
	p.Light.Y = y

	if !enteringFromRoom {
		fmt.Println("død og jøde, markens grøde")
	}
}

func (p *Player) Update(cmd *commands.Commands) {
	p.Slambox.Update(cmd)

	scene, _ := commands.Get[engine.Scene](cmd)
	slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)
	globaldata_, _ := commands.Get[globaldata.GlobalData](cmd)

	x, y := p.pivotTransform.GetPos(false)

	p.Light.X = x
	p.Light.Y = y

	switch p.State {
	case SLAMMING:
		info, finished := p.OnClipFinish.Poll()
		if finished && info["clip"] == "Slam" {
			p.State = IDLE
			p.jumpOffset = 0
			p.jumpOffsetvel = 0
			p.animatedSprite.SetPos(0, 0)
		}

		if p.jumpOffsetvel > 0 {
			p.jumpOffsetvel -= 0.3
		} else {
			p.jumpOffsetvel -= 0.6
		}

		p.jumpOffset += p.jumpOffsetvel
		p.jumpOffset = math.Max(p.jumpOffset, 0)
		if p.jumpOffset == 0 && !p.hasSlammedBox {
			p.OnSlam.Raise()
			slamboxHits := scene.GetRoot().GetChildrenFunc(
				func(n *engine.Node) bool {
					slambox_, ok := n.GetValue().(*slamboxactor.Slambox)
					if !ok {
						return false
					}

					return slambox_.GetBackendID() == p.slamboxIDBuffer
				},
			)

			slamboxGroupHits := scene.GetRoot().GetChildrenFunc(
				func(n *engine.Node) bool {
					slamboxGroup, ok := n.GetValue().(*slamboxgroup.SlamboxGroup)
					if !ok {
						return false
					}
					return slamboxGroup.BackendIndex == p.slamboxIDBuffer
				},
			)

			if len(slamboxHits) != 0 {
				slamboxactor, ok := slamboxHits[0].GetValue().(*slamboxactor.Slambox)
				if !ok {
					fmt.Println("Slambox node could not be cast to slambox actor")
					return
				}
				slamboxactor.RequestSlam(p.slamDirBuffer)
				p.slamDirBuffer = maths.DirNone
				p.slamboxIDBuffer = -1
				p.hasSlammedBox = true
			} else if len(slamboxGroupHits) != 0 {
				slamboxgroup, ok := slamboxGroupHits[0].GetValue().(*slamboxgroup.SlamboxGroup)
				if !ok {
					fmt.Println("Slambox node could not be cast to slambox actor")
					return
				}
				slamboxgroup.RequestSlam(p.slamDirBuffer)
				p.slamDirBuffer = maths.DirNone
				p.slamboxIDBuffer = -1
				p.hasSlammedBox = true
			} else {
				fmt.Println("No slambox matching ID")
				return
			}
		}
		p.animatedSprite.SetPos(0, -p.jumpOffset)
	case IDLE:
		p.animatedSprite.SwitchClip(IDLE_ANIM)
		playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]

		if playerControls.PollAction("DoorInteract") {
			doorNode, ok := scene.GetNodeFunc(func(n *engine.Node) bool {
				door, ok := engine.As[*door.Door](n.GetValue())
				if !ok {
					return false
				}
				triggerRect := door.Trigger.GetRect()
				return triggerRect.Contains(p.GetCenterPos())
			})

			if ok {
				doorActor, _ := engine.As[*door.Door](doorNode.GetValue())
				if !doorActor.Locked {
					p.State = LEAVING
					p.Direction = doorActor.Direction

					p.setDoorOffset(doorActor.Hitbox)

					globaldata_.Temp.SaveLevelState(scene)
					cmd.InputHandler.InputSchemes["PlayerControls"].Active = false
					p.jumpOffsetvel = 3.5
				}
			}
		}
	case MOVING:
		_, finished := p.OnMoveFinish.Poll()
		if finished {
			p.Direction = maths.Opposite(p.Direction)
			p.State = IDLE
		}
	case DYING:
		p.animatedSprite.SwitchClip(IDLE_ANIM)
	case LEAVING:
		p.jumpOffsetvel -= 0.2
		p.jumpOffset += p.jumpOffsetvel
		p.trueDoorOffset = maths.Lerp(p.trueDoorOffset, p.doorOffset, 0.1)
		p.animatedSprite.SetPos(p.trueDoorOffset, -p.jumpOffset)
	case ENTERING:
		if p.jumpOffset > 0 {
			p.jumpOffsetvel -= 0.32
		} else if p.jumpOffset <= 0 && p.jumpOffsetvel <= 0 {
			p.jumpOffset = 0
			p.jumpOffsetvel = 0
			p.State = IDLE
		}
		p.animatedSprite.SetPos(0, -p.jumpOffset)
		p.jumpOffset += p.jumpOffsetvel
	}

	direction := p.readMoveInput(cmd)
	if direction != maths.DirNone {
		p.inputbuffer.Set(direction)
	}

	p.pivotTransform.SetAngle(maths.DirToRadians(p.Direction))

	p.inputbuffer.Update()

	if cmd.InputHandler.InputSchemes["PlayerControls"].PollAction("Reset") {
		p.ResetLevel(scene, globaldata_)
	}

	// TODO: This is not clean at all
	// Potentially just turn into a separate function, not a terrible "solution"...
	// Would be a lot better if we could use events for this
	_, ok := scene.GetNodeFunc(func(n *engine.Node) bool {
		hazard, ok := engine.As[*hazard.Hazard](n.GetValue())
		if !ok {
			return false
		}
		triggerRect := hazard.Trigger.GetRect()
		return triggerRect.Overlapping(p.GetRect())
	})
	if ok {
		p.ResetLevel(scene, globaldata_)
	}

	keyNode, ok := scene.GetNodeFunc(func(n *engine.Node) bool {
		key, ok := engine.As[*key.Key](n.GetValue())
		if !ok {
			return false
		}

		return key.Hitbox.Overlapping(p.GetRect())
	})

	if ok {
		key, _ := engine.As[*key.Key](keyNode.GetValue())
		if !key.PickedUp {
			gamestate_, _ := commands.Get[globaldata.GlobalData](cmd)
			inventory := gamestate_.Persist.Profile.Inventory
			inventory.Keys[key.EntityIid] = &globaldata.Key{key.DoorIid, false}
			key.OnPickupEv.Raise().WithData("Transform", p.Transform2D)
		}
	}

	moveDir := p.inputbuffer.Read()

	if moveDir == maths.DirNone || p.State != IDLE {
		return
	}

	// Check whether we should slam, do nothing or dash
	slamboxQuery := slamboxenv.QuerySlamboxes(p.GetRect().Extended(moveDir, 1.0), slambox.QueryFilter{p.Slambox.GetBackendID()})
	tilemapCollision := slamboxenv.CheckTileOverlap(p.GetRect().Extended(moveDir, 1.0))

	if slamboxQuery.HitKind == slambox.NONE && !tilemapCollision {
		p.OnMove.WithData("Direction", moveDir).Raise()
		p.Dash(moveDir)
		p.inputbuffer.Clear()
		x, y := p.GetCenterPos()
		scene.SpawnBundleV2(cmd, MakeJumpParticlesBundle(x, y, moveDir, p.GetRect().Width/2))
		return
	}

	if !tilemapCollision {
		//		p.OnMove.WithData("Direction", moveDir).Raise()
		p.hasSlammedBox = false
		p.slamboxIDBuffer = slamboxQuery.Index
		p.slamDirBuffer = moveDir
		p.inputbuffer.Clear()
		p.StartSlamming(moveDir)
	}
}

func (p *Player) GetCenterPos() (float64, float64) {
	x, y := p.Transform2D.GetPos(false)
	return x + p.GetRect().Width/2, y + p.GetRect().Height/2
}

func (p *Player) ResetLevel(scene *engine.Scene, gameState *globaldata.GlobalData) {
	levelstate := gameState.Temp.LevelStates[scene.GetName()]

	p.SetPos(levelstate.PlayerSpawnPos.X, levelstate.PlayerSpawnPos.Y)
	p.Direction = levelstate.PlayerSpawnDir

	slamboxes := scene.GetRoot().GetChildrenFunc(
		func(n *node.Node[engine.Actor]) bool {
			_, ok := engine.As[*slamboxgroup.SlamboxGroup](n.GetValue())
			return ok
		},
	)

	for _, slambox := range slamboxes {
		slamboxactor, _ := engine.As[*slamboxgroup.SlamboxGroup](slambox.GetValue())
		if _, ok := engine.As[*Player](slambox.GetValue()); ok {
			continue
		}
		slamboxPos := levelstate.SlamboxPositions[slambox.GetID()]
		slamboxactor.SetPos(slamboxPos.X, slamboxPos.Y)
	}
}

func (p *Player) setDoorOffset(doorRect *maths.Rect) {
	dx := doorRect.Cx() - p.GetRect().Cx()
	dy := doorRect.Cy() - p.GetRect().Cy()
	switch p.Direction {
	case maths.DirUp:
		p.doorOffset = dx
	case maths.DirDown:
		p.doorOffset = -dx
	case maths.DirRight:
		p.doorOffset = dy
	case maths.DirLeft:
		p.doorOffset = -dy
	}
}

func (p *Player) Dash(direction maths.Direction) {
	p.inputbuffer.Clear()
	p.Direction = direction
	p.State = MOVING
	p.Slambox.RequestSlam(direction)
	p.animatedSprite.SwitchClip(DASH_INIT_ANIM)
	// sound_v2.PlaySound("playerDash", "sfxMaster", 0.06)
	// p.playJumpParticles(direction)
}

func (p *Player) StartSlamming(direction maths.Direction) {
	// sound_v2.PlaySound("playerDash", "sfxMaster", 0.06)
	// p.canPlaySlamSound = true
	p.Direction = maths.Opposite(direction)
	p.animatedSprite.SwitchClip(SLAM_ANIM)
	p.State = SLAMMING
	p.jumpOffsetvel = 4
}

func (p *Player) readMoveInput(cmd *commands.Commands) maths.Direction {
	playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
	if playerControls.PollAction("moveLeft") {
		return maths.DirLeft
	} else if playerControls.PollAction("moveRight") {
		return maths.DirRight
	} else if playerControls.PollAction("moveUp") {
		return maths.DirUp
	} else if playerControls.PollAction("moveDown") {
		return maths.DirDown
	}
	return maths.DirNone
}

func NewPlayer(slambox *slamboxactor.Slambox, inputBufferDuration float64) *Player {
	player := &Player{
		Slambox:      slambox,
		inputbuffer:  inputbuffer.NewInputBuffer(inputBufferDuration),
		OnMoveFinish: events.NewBusFrom(slambox.OnMoveFinishEv),
		OnMove:       events.NewEvent(),
		OnSlam:       events.NewEvent(),

		// TODO: Change so that the parameters are more intuitive
		// noise factor should not be hard-coded
		Light: &shaders.Light{
			InnerRadius: 0,
			OuterRadius: 200,
			ZOffset:     0.2,
			Intensity:   0.6,
			R:           1.0,
			G:           1.0,
			B:           1.0,
		},
	}

	return player
}

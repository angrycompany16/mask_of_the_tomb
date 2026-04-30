package player

import (
	"fmt"
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
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
	"mask_of_the_tomb/internal/game/actors/doorv2"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/sceneswitch"
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
	direction                 maths.Direction
	spriteTransform           *transform2D.Transform2D
	animatedSprite            *animatedsprite.AnimatedSprite
	pivotTransform            *transform2D.Transform2D
	inputbuffer               *inputbuffer.InputBuffer
	OnMoveFinish              *eventsv2.EventBus
	OnClipFinish              *eventsv2.EventBus
	OnMove                    *eventsv2.Event
	jumpOffset, jumpOffsetvel float64
	slamboxIDBuffer           int
	slamDirBuffer             maths.Direction
	hasSlammedBox             bool
	doorOffset                float64
	trueDoorOffset            float64
	doorY                     float64
	Light                     *shaders.Light

	// Turn engine.bundle into some kind of generic thing (probably
	// just a function), then pass it in explicitly. Hopefully we can
	// then spawn the bundle after we've started playing the game, with
	// no import problems :D
	// jumpParticleSys engine.bundle

	// moveSpeed   float64 // 5.0
	// defaultPlayerHealth   = 5.0
	// invincibilityDuration = time.Second
	// inputBufferDuration   = 0.1
}

// Children
// jumpParticlesBroad *particles.ParticleSystem
// jumpParticlesTight *particles.ParticleSystem
// sprite                    *ebiten.Image

func (p *Player) Init(cmd *commands.Commands) {
	p.Slambox.Init(cmd)

	scene, ok := commands.Get[engine.Scene](cmd)
	if !ok {
		panic("Missing scene (Player)")
	}

	sceneswitch, ok := commands.Get[sceneswitch.SceneSwitch](cmd)
	if !ok {
		panic("Missing sceneswitch (Player)")
	}
	spawnDoorIid := sceneswitch.SpawnEntityIid
	spawnDoorNode, ok := scene.GetNodeFunc(
		func(n *node.Node[engine.Actor]) bool {
			doorActor, ok := engine.As[*doorv2.DoorV2](n.GetValue())
			if !ok {
				return false
			}
			return doorActor.EntityIid == spawnDoorIid
		},
	)
	if !ok {
		fmt.Println("Problem: Could not find door to spawn at...")
	} else {
		doorActor, ok := engine.As[*doorv2.DoorV2](spawnDoorNode.GetValue())
		if !ok {
			fmt.Println("Spawn door node could not convert to Door actor")
		}
		p.SetPos(doorActor.GetSpawnPos())
		p.Transform2D.SetPos(doorActor.GetSpawnPos())
		p.Transform2D.Propagate()
		p.State = ENTERING
		p.jumpOffset = -2 * p.GetRect().Height
		p.jumpOffsetvel = 4.5
		p.trueDoorOffset = 0
		p.doorOffset = 0
	}

	p.direction = sceneswitch.SpawnDirection

	playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
	playerControls.RegisterAction("moveLeft", input.KeyJustPressedAction(ebiten.KeyA))
	playerControls.AddBinding("moveLeft", input.KeyJustPressedAction(ebiten.KeyLeft))
	playerControls.RegisterAction("moveRight", input.KeyJustPressedAction(ebiten.KeyD))
	playerControls.AddBinding("moveRight", input.KeyJustPressedAction(ebiten.KeyRight))
	playerControls.RegisterAction("moveUp", input.KeyJustPressedAction(ebiten.KeyW))
	playerControls.AddBinding("moveUp", input.KeyJustPressedAction(ebiten.KeyUp))
	playerControls.RegisterAction("moveDown", input.KeyJustPressedAction(ebiten.KeyS))
	playerControls.AddBinding("moveDown", input.KeyJustPressedAction(ebiten.KeyDown))

	// Would be very nice to set up a reference like this in another
	// way
	// But how? I guess we would have to link them together somehow
	// in the bundle
	childNode, ok := scene.GetNodeByName("PlayerSprite")
	p.spriteTransform, ok = engine.As[*transform2D.Transform2D](childNode.GetValue())
	p.animatedSprite, ok = engine.As[*animatedsprite.AnimatedSprite](childNode.GetValue())
	p.OnClipFinish = eventsv2.NewBusFrom(p.animatedSprite.OnClipFinished)

	pivotNode, ok := scene.GetNodeByName("PlayerPivot")
	p.pivotTransform, ok = engine.As[*transform2D.Transform2D](pivotNode.GetValue())

	p.pivotTransform.Propagate()
	x, y := p.pivotTransform.GetPos(false)
	p.Light.X = x
	p.Light.Y = y

	if !ok {
		fmt.Println("død og jøde, markens grøde")
	}
}

func (p *Player) Update(cmd *commands.Commands) {
	p.Slambox.Update(cmd)

	scene, _ := commands.Get[engine.Scene](cmd)
	slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)

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
			slamboxHits := scene.GetRoot().GetChildrenFunc(
				func(n *engine.Node) bool {
					slambox_, ok := n.GetValue().(*slamboxactor.Slambox)
					if !ok {
						return false
					}
					return slambox_.GetBackendID() == p.slamboxIDBuffer
				},
			)
			if len(slamboxHits) == 0 {
				fmt.Println("No slambox matching ID")
				return
			}
			slamboxactor, ok := slamboxHits[0].GetValue().(*slamboxactor.Slambox)
			if !ok {
				fmt.Println("Slambox node could not be cast to slambox actor")
				return
			}
			slamboxactor.RequestSlam(p.slamDirBuffer)
			p.slamDirBuffer = maths.DirNone
			p.slamboxIDBuffer = -1
			p.hasSlammedBox = true
		}
		p.animatedSprite.SetPos(0, -p.jumpOffset)
	case IDLE:
		p.animatedSprite.SwitchClip(IDLE_ANIM)
		playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
		if playerControls.PollAction("DoorInteract") {
			doorNode, ok := scene.GetNodeFunc(func(n *engine.Node) bool {
				door, ok := engine.As[*doorv2.DoorV2](n.GetValue())
				if !ok {
					return false
				}
				triggerRect := door.Trigger.GetRect()
				return triggerRect.Contains(p.getCenterPos())
			})

			if ok {
				p.State = LEAVING
				door, _ := engine.As[*doorv2.DoorV2](doorNode.GetValue())
				p.direction = door.Direction

				p.setDoorOffset(door.Hitbox)

				cmd.InputHandler.InputSchemes["PlayerControls"].Active = false
				p.jumpOffsetvel = 3.5
			}
		}
	case MOVING:
		_, finished := p.OnMoveFinish.Poll()
		if finished {
			p.direction = maths.Opposite(p.direction)
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

	moveDir := p.inputbuffer.Read()

	p.pivotTransform.SetAngle(maths.DirToRadians(p.direction))

	p.inputbuffer.Update()

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
		// Gotta find the right x, y
		x, y := p.getCenterPos()
		scene.SpawnBundle(cmd, MakeJumpParticlesBundle(x, y, moveDir, p.GetRect().Width/2))
		return
	}

	if !tilemapCollision {
		p.hasSlammedBox = false
		p.slamboxIDBuffer = slamboxQuery.Index
		p.slamDirBuffer = moveDir
		p.inputbuffer.Clear()
		p.StartSlamming(moveDir)
	}
}

func (p *Player) getCenterPos() (float64, float64) {
	x, y := p.Transform2D.GetPos(false)
	return x + p.GetRect().Width/2, y + p.GetRect().Height/2
}

func (p *Player) setDoorOffset(doorRect *maths.Rect) {
	dx := doorRect.Cx() - p.GetRect().Cx()
	dy := doorRect.Cy() - p.GetRect().Cy()
	switch p.direction {
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
	p.direction = direction
	p.State = MOVING
	p.Slambox.RequestSlam(direction)
	p.animatedSprite.SwitchClip(DASH_INIT_ANIM)
	// sound_v2.PlaySound("playerDash", "sfxMaster", 0.06)
	// p.playJumpParticles(direction)
}

func (p *Player) StartSlamming(direction maths.Direction) {
	// sound_v2.PlaySound("playerDash", "sfxMaster", 0.06)
	// p.canPlaySlamSound = true
	p.direction = maths.Opposite(direction)
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
		OnMoveFinish: eventsv2.NewBusFrom(slambox.OnMoveFinishEv),
		OnMove:       eventsv2.NewEvent(),

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

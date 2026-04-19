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

type playerState int

const (
	Idle playerState = iota
	Moving
	Slamming
	Dying
)

// var (
// 	playerSpritePath       = filepath.Join(assets.PlayerFolder, "player.png")
// 	jumpParticlesBroadPath = filepath.Join("assets", "particlesystems", "jump-broad.yaml")
// 	jumpParticlesTightPath = filepath.Join("assets", "particlesystems", "jump-tight.yaml")
// )

type Player struct {
	*slamboxactor.Slambox
	State                     playerState
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
	}

	p.direction = sceneswitch.SpawnDirection

	cmd.InputHandler.RegisterAction("moveLeft", input.KeyJustPressedAction(ebiten.KeyA))
	cmd.InputHandler.AddBinding("moveLeft", input.KeyJustPressedAction(ebiten.KeyLeft))
	cmd.InputHandler.RegisterAction("moveRight", input.KeyJustPressedAction(ebiten.KeyD))
	cmd.InputHandler.AddBinding("moveRight", input.KeyJustPressedAction(ebiten.KeyRight))
	cmd.InputHandler.RegisterAction("moveUp", input.KeyJustPressedAction(ebiten.KeyW))
	cmd.InputHandler.AddBinding("moveUp", input.KeyJustPressedAction(ebiten.KeyUp))
	cmd.InputHandler.RegisterAction("moveDown", input.KeyJustPressedAction(ebiten.KeyS))
	cmd.InputHandler.AddBinding("moveDown", input.KeyJustPressedAction(ebiten.KeyDown))

	// Would be very nice to set up a reference like this in another
	// way
	// But how? I guess we would have to link them together somehow
	// in the bundle
	childNode, ok := scene.GetNodeByName("PlayerSprite")
	p.spriteTransform, ok = engine.As[*transform2D.Transform2D](childNode.GetValue())
	p.animatedSprite, ok = engine.As[*animatedsprite.AnimatedSprite](childNode.GetValue())
	p.OnClipFinish = eventsv2.NewEventBus(p.animatedSprite.OnClipFinished)

	pivotNode, ok := scene.GetNodeByName("PlayerPivot")
	p.pivotTransform, ok = engine.As[*transform2D.Transform2D](pivotNode.GetValue())

	// If relevant, spawn in a different position

	if !ok {
		fmt.Println("død og jøde, markens grøde")
	}
}

func (p *Player) Update(cmd *commands.Commands) {
	p.Slambox.Update(cmd)

	x, y := p.pivotTransform.GetPos(false)
	p.Light.X = x
	p.Light.Y = y

	scene, ok := commands.Get[engine.Scene](cmd)
	if !ok {
		panic("Missing scene (Player)")
	}

	slamboxenv, ok := commands.Get[slambox.SlamboxEnvironment](cmd)
	if !ok {
		panic("Missing slambox env (Player)")
	}

	switch p.State {
	case Slamming:
		info, finished := p.OnClipFinish.Poll()
		if finished && info["clip"] == "Slam" {
			p.State = Idle
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
			theOther, ok := scene.GetRoot().GetChildFunc(
				func(n *engine.Node) bool {
					slambox_, ok := n.GetValue().(*slamboxactor.Slambox)
					if !ok {
						return false
					}
					return slambox_.GetBackendID() == p.slamboxIDBuffer
				},
			)
			if !ok {
				fmt.Println("Major problems")
				return
			}
			slamboxactor, ok := theOther.GetValue().(*slamboxactor.Slambox)
			if !ok {
				fmt.Println("Major problems")
				return
			}
			slamboxactor.RequestSlam(p.slamDirBuffer)
			p.slamDirBuffer = maths.DirNone
			p.slamboxIDBuffer = -1
			p.hasSlammedBox = true
		}
		// p.jumpOffset = maths.Clamp(p.jumpOffset, 0, 1000000)
		p.animatedSprite.SetPos(0, -p.jumpOffset)
		// if p.jumpOffset == 0 && p.canPlaySlamSound {
		// 	sound_v2.PlaySound("playerSlam", "sfxMaster", 0.04)
		// 	p.canPlaySlamSound = false
		// 	camera.Shake(0.4, 7, 1)
		// }
	case Idle:
		p.animatedSprite.SwitchClip(IDLE_ANIM)
	case Moving:
		// p.movebox.Update()
		_, finished := p.OnMoveFinish.Poll()
		if finished {
			p.direction = maths.Opposite(p.direction)
			p.State = Idle
		}
	case Dying:
		p.animatedSprite.SwitchClip(IDLE_ANIM)
		// p.deathAnim.Update()
	}

	direction := p.readMoveInput(cmd)
	if direction != maths.DirNone {
		// Also include platforms

		// 1 - Get all platforms
		// 2 - Add to slambox env
		p.OnMove.WithData("Direction", direction).Raise()
		p.inputbuffer.Set(direction)
	}

	moveDir := p.inputbuffer.Read()

	p.pivotTransform.SetAngle(maths.DirToRadians(p.direction))

	p.inputbuffer.Update()

	if moveDir == maths.DirNone || p.State != Idle {
		return
	}

	// Check whether we should slam, do nothing or dash
	slamboxQuery := slamboxenv.QuerySlamboxes(p.GetRect().Extended(moveDir, 1.0), slambox.QueryFilter{p.Slambox.GetBackendID()})
	tilemapCollision := slamboxenv.CheckTileOverlap(p.GetRect().Extended(moveDir, 1.0))

	if slamboxQuery.HitKind == slambox.NONE && !tilemapCollision {
		p.Dash(moveDir)
		p.inputbuffer.Clear()
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

func (p *Player) Dash(direction maths.Direction) {
	p.inputbuffer.Clear()
	p.direction = direction
	p.State = Moving
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
	p.State = Slamming
	p.jumpOffsetvel = 4
}

func (p *Player) readMoveInput(cmd *commands.Commands) maths.Direction {
	if cmd.InputHandler.PollAction("moveLeft") {
		return maths.DirLeft
	} else if cmd.InputHandler.PollAction("moveRight") {
		return maths.DirRight
	} else if cmd.InputHandler.PollAction("moveUp") {
		return maths.DirUp
	} else if cmd.InputHandler.PollAction("moveDown") {
		return maths.DirDown
	}
	return maths.DirNone
}

func NewPlayer(slambox *slamboxactor.Slambox, inputBufferDuration float64) *Player {
	player := &Player{
		Slambox:      slambox,
		inputbuffer:  inputbuffer.NewInputBuffer(inputBufferDuration),
		OnMoveFinish: eventsv2.NewEventBus(slambox.OnMoveFinishEv),
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

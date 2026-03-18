package player

import (
	"fmt"
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/inputbuffer"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	IDLE_ANIM = iota
	DASH_INIT_ANIM
	DASH_LOOP_ANIM
	SLAM_ANIM
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
	State           playerState
	direction       maths.Direction
	spriteTransform *transform2D.Transform2D
	animatedSprite  *animatedsprite.AnimatedSprite
	inputbuffer     *inputbuffer.InputBuffer
	OnMoveFinish    *eventsv2.EventBus
	// moveSpeed   float64 // 5.0
	// defaultPlayerHealth   = 5.0
	// invincibilityDuration = time.Second
	// inputBufferDuration   = 0.1
}

// Children
// jumpParticlesBroad *particles.ParticleSystem
// jumpParticlesTight *particles.ParticleSystem
// sprite                    *ebiten.Image

func (p *Player) Update(cmd *engine.Commands) {
	p.Slambox.Update(cmd)

	switch p.State {
	case Slamming:
		// _, finished := p.clipFinishedListener.Poll()
		// if finished {
		// 	p.State = Idle
		// 	p.jumpOffset = 0
		// 	p.jumpOffsetvel = 0
		// }

		// if p.jumpOffsetvel > 0 {
		// 	p.jumpOffsetvel -= 0.3
		// } else {
		// 	p.jumpOffsetvel -= 0.6
		// }

		// p.jumpOffset += p.jumpOffsetvel
		// p.jumpOffset = maths.Clamp(p.jumpOffset, 0, 1000000)
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
		p.inputbuffer.Set(direction)
	}

	moveDir := p.inputbuffer.Read()

	// TODO: If we can slam or are already locked to a wall, we don't
	// want to do this
	if moveDir != maths.DirNone && p.State == Idle {
		p.Dash(moveDir)
	}

	p.spriteTransform.SetAngle(maths.DirToRadians(p.direction))

	p.inputbuffer.Update()
}

func (p *Player) Init(cmd *engine.Commands) {
	p.Slambox.Init(cmd)

	cmd.InputHandler().RegisterAction("moveLeft", input.KeyJustPressedAction(ebiten.KeyLeft))
	cmd.InputHandler().RegisterAction("moveRight", input.KeyJustPressedAction(ebiten.KeyRight))
	cmd.InputHandler().RegisterAction("moveUp", input.KeyJustPressedAction(ebiten.KeyUp))
	cmd.InputHandler().RegisterAction("moveDown", input.KeyJustPressedAction(ebiten.KeyDown))

	// Would be very nive to set up a reference like this in another
	// way
	// But how? I guess we would have to link them together somehow
	// in the bundle
	childNode, ok := cmd.Scene().GetNodeByName("PlayerSprite")
	p.spriteTransform, ok = engine.GetActor[*transform2D.Transform2D](childNode.GetValue())
	p.animatedSprite, ok = engine.GetActor[*animatedsprite.AnimatedSprite](childNode.GetValue())

	if !ok {
		fmt.Println("død og jøde, markens grøde")
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

func (p *Player) readMoveInput(cmd *engine.Commands) maths.Direction {
	if cmd.InputHandler().PollAction("moveLeft") {
		return maths.DirLeft
	} else if cmd.InputHandler().PollAction("moveRight") {
		return maths.DirRight
	} else if cmd.InputHandler().PollAction("moveUp") {
		return maths.DirUp
	} else if cmd.InputHandler().PollAction("moveDown") {
		return maths.DirDown
	}
	return maths.DirNone
}

func NewPlayer(slambox *slamboxactor.Slambox, inputBufferDuration float64) *Player {
	player := &Player{
		Slambox:      slambox,
		inputbuffer:  inputbuffer.NewInputBuffer(inputBufferDuration),
		OnMoveFinish: eventsv2.NewEventBus(slambox.OnMoveFinishEv),
	}

	return player
}

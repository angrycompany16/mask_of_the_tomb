package player

import (
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	IDLE_ANIM = iota
	DASH_INIT_ANIM
	DASH_LOOP_ANIM
	SLAM_ANIM
)

// var (
// 	playerSpritePath       = filepath.Join(assets.PlayerFolder, "player.png")
// 	jumpParticlesBroadPath = filepath.Join("assets", "particlesystems", "jump-broad.yaml")
// 	jumpParticlesTightPath = filepath.Join("assets", "particlesystems", "jump-tight.yaml")
// )

// TODO: Allow for sprites which aren't exactly 16x16

type Player struct {
	*slamboxactor.Slambox
	// inputbuffer *inputbuffer.InputBuffer
	// moveSpeed   float64 // 5.0
	// defaultPlayerHealth   = 5.0
	// invincibilityDuration = time.Second
	// inputBufferDuration   = 0.1
}

// Children
// jumpParticlesBroad *particles.ParticleSystem
// jumpParticlesTight *particles.ParticleSystem
// sprite                    *ebiten.Image

// ------ INIT ------
func (p *Player) Init(cmd *engine.Commands) {
	p.Slambox.Init(cmd)
	cmd.InputHandler().RegisterAction("moveLeft", input.KeyJustPressedAction(ebiten.KeyLeft))
	cmd.InputHandler().RegisterAction("moveRight", input.KeyJustPressedAction(ebiten.KeyRight))
	cmd.InputHandler().RegisterAction("moveUp", input.KeyJustPressedAction(ebiten.KeyUp))
	cmd.InputHandler().RegisterAction("moveDown", input.KeyJustPressedAction(ebiten.KeyDown))

	p.Slambox.Init(cmd)

	// spawn important children
}

func (p *Player) Update(cmd *engine.Commands) {
	p.Slambox.Update(cmd)
	if cmd.InputHandler().PollAction("moveLeft") {
		p.Slambox.RequestSlam(maths.DirLeft)
	} else if cmd.InputHandler().PollAction("moveRight") {
		p.Slambox.RequestSlam(maths.DirRight)
	} else if cmd.InputHandler().PollAction("moveUp") {
		p.Slambox.RequestSlam(maths.DirUp)
	} else if cmd.InputHandler().PollAction("moveDown") {
		p.Slambox.RequestSlam(maths.DirDown)
	}
}

// ------ GETTERS ------
func NewPlayer(slambox *slamboxactor.Slambox) *Player {
	player := &Player{
		Slambox: slambox,
	}

	return player
}

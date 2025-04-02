package player

import (
	"mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	"mask_of_the_tomb/internal/maths"
)

type playerState int

const (
	Idle playerState = iota
	Moving
	Slamming
	Dying
)

// FINALLY it feels good
func (p *Player) StartSlamming(direction maths.Direction) {
	// TODO: orient player properly
	p.direction = maths.Opposite(direction)
	p.animator.SwitchClip(slamAnim)
	p.State = Slamming
	p.jumpOffsetvel = 2.5
}

func (p *Player) Update() {
	switch p.State {
	case Slamming:
		_, finished := p.clipFinishedListener.Poll()
		if finished {
			p.State = Idle
			p.jumpOffset = 0
			p.jumpOffsetvel = 0
		}

		if p.jumpOffsetvel > 0 {
			p.jumpOffsetvel -= 0.1
		} else {
			p.jumpOffsetvel -= 0.25
		}

		p.jumpOffset += p.jumpOffsetvel
		p.jumpOffset = maths.Clamp(p.jumpOffset, 0, 1000000)
	case Idle:
		p.animator.SwitchClip(idleAnim)
	case Moving:
		p.movebox.Update()
		_, finished := p.moveFinishedListener.Poll()
		if finished {
			p.direction = maths.Opposite(p.direction)
			p.State = Idle
		}
	case Dying:
		p.animator.SwitchClip(idleAnim)
		p.deathAnim.Update()
		// Update death animation
	}

	direction := p.getMoveInput()
	if direction != maths.DirNone {
		p.InputBuffer.set(direction)
	}

	p.InputBuffer.update()
	p.Hitbox.SetPos(p.movebox.GetPos())
	p.animator.Update()
	p.jumpParticlesBroad.Update()
	p.jumpParticlesTight.Update()
}

// TODO: this will be changed back when we add some kind of death (sprite) animation
func (p *Player) Draw() {
	p.jumpParticlesBroad.Draw()
	p.jumpParticlesTight.Draw()
	if p.State == Dying {
		p.deathAnim.Draw()
	} else {
		posX, posY := p.movebox.GetPos()
		camX, camY := camera.GetPos()
		jumpOffsetX, jumpOffsetY := p.getJumpOffset()
		ebitenrenderutil.DrawAtRotated(
			p.animator.GetSprite(),
			rendering.RenderLayers.Playerspace,
			posX-camX-jumpOffsetX,
			posY-camY-jumpOffsetY,
			maths.ToRadians(p.direction),
			0.5,
			0.5,
		)
	}
}

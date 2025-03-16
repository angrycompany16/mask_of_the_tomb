package player

import (
	"mask_of_the_tomb/internal/maths"
	"math"
)

type playerState int

const (
	Idle playerState = iota
	Moving
	Slamming
)

// FINALLY it feels good
func (p *Player) StartSlamming(direction maths.Direction) {
	// TODO: orient player properly
	p.angle = maths.ToRadians(maths.Opposite(direction))
	p.animator.SwitchClip(slamAnim)
	p.State = Slamming
	p.jumpOffsetvel = 2.5
}

func (p *Player) PreUpdate() {}

func (p *Player) Update() {
	switch p.State {
	case Slamming:
		if p.finishedClipEventListener.Poll() {
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
		p.PosX += moveSpeed * p.moveDirX
		p.PosY += moveSpeed * p.moveDirY

		if p.moveDirX < 0 {
			p.PosX = maths.Clamp(p.PosX, p.targetPosX, p.PosX)
		} else if p.moveDirX > 0 {
			p.PosX = maths.Clamp(p.PosX, p.PosX, p.targetPosX)
		}
		if p.moveDirY < 0 {
			p.PosY = maths.Clamp(p.PosY, p.targetPosY, p.PosY)
		} else if p.moveDirY > 0 {
			p.PosY = maths.Clamp(p.PosY, p.PosY, p.targetPosY)
		}

		if p.PosX == p.targetPosX {
			p.moveDirX = 0
		}
		if p.PosY == p.targetPosY {
			p.moveDirY = 0
		}

		if p.PosX == p.targetPosX && p.PosY == p.targetPosY {
			p.angle = p.angle - math.Pi
			p.State = Idle
		}
	}

	direction := p.getMoveInput()
	if direction != maths.DirNone {
		p.InputBuffer.set(direction)
	}

	p.InputBuffer.update()

	p.damageOverlay.Update()
	p.Hitbox.SetPos(p.PosX, p.PosY)

	p.animator.Update()
}

func (p *Player) PostUpdate() {}

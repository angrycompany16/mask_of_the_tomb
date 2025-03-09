package player

import (
	"mask_of_the_tomb/internal/maths"
)

type playerState int

const (
	StateIdle playerState = iota
	StateMoving
)

func (p *Player) Update() {
	switch p.State {
	case StateIdle:

	case StateMoving:
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
			p.State = StateIdle
		}
	}

	direction := p.getMoveInput()
	if direction != maths.DirNone {
		p.InputBuffer.set(direction)
	}
	p.InputBuffer.update()

	p.damageOverlay.Update()
	p.Hitbox.SetPos(p.PosX, p.PosY)

	p.playerTestAnim.Update()
}

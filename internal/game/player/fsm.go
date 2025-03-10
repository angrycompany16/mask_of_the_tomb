package player

import (
	"mask_of_the_tomb/internal/maths"
	"math"
	"time"
)

type playerState int

const (
	Idle playerState = iota
	Moving
	Slamming
)

var (
	slamming       = false
	slamFinishChan = make(chan int, 1)
	SlamHitChan    = make(chan int, 1)
)

func (p *Player) Update() {
	switch p.State {
	case Slamming:
		if slamming {
			select {
			case <-slamFinishChan:
				slamming = false
				p.State = Idle
			default:
			}
		} else {
			slamming = true
			go p.Slam()
		}
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

func (p *Player) Slam() {
	p.animator.SwitchClip(slamAnim)
	time.Sleep(500 * time.Millisecond)
	SlamHitChan <- 1
	time.Sleep(500 * time.Millisecond)
	slamFinishChan <- 1
}

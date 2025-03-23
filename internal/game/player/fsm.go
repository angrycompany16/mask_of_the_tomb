package player

import (
	"mask_of_the_tomb/internal/maths"
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
	p.direction = maths.Opposite(direction)
	p.animator.SwitchClip(slamAnim)
	p.State = Slamming
	p.jumpOffsetvel = 2.5
}

func (p *Player) Update() {
	switch p.State {
	case Slamming:
		_, finished := p.finishedClipEventListener.Poll()
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
		_, finished := p.finishedMoveEventListener.Poll()
		if finished {
			p.direction = maths.Opposite(p.direction)
			p.State = Idle
		}
	}

	direction := p.getMoveInput()
	if direction != maths.DirNone {
		p.InputBuffer.set(direction)
	}

	p.InputBuffer.update()

	p.damageOverlay.Update()
	p.Hitbox.SetPos(p.movebox.GetPos())

	p.animator.Update()
}

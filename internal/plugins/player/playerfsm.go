package player

import (
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/camera"
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
	p.dashSound.Play()
	p.direction = maths.Opposite(direction)
	p.animator.SwitchClip(slamAnim)
	p.State = Slamming
	p.jumpOffsetvel = 2.5
	p.canPlaySlamSound = true
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
		if p.jumpOffset == 0 && p.canPlaySlamSound {
			p.slamSound.Play()
			p.canPlaySlamSound = false
			camera.Shake(0.4, 7, 1)
		}
	case Idle:
		p.animator.SwitchClip(idleAnim)
	case Moving:
		p.movebox.Update()
		_, finished := p.moveFinishedListener.Poll()
		if finished && p.InputBuffer.Read() == maths.DirNone {
			p.direction = maths.Opposite(p.direction)
			p.State = Idle
		}
	case Dying:
		p.animator.SwitchClip(idleAnim)
		p.deathAnim.Update()
	}

	direction := p.readMoveInput()
	if direction != maths.DirNone {
		p.InputBuffer.Set(direction)
	}

	playerMove := p.InputBuffer.Read()

	if playerMove != maths.DirNone && p.CanMove() && !p.Disabled {
		p.OnMove.Raise(events.EventInfo{Data: playerMove})
	}

	p.InputBuffer.Update()
	p.hitbox.SetPos(p.movebox.GetPos())
	p.animator.Update()
	p.jumpParticlesBroad.Update()
	p.jumpParticlesTight.Update()
}

// TODO: this will be changed back when we add some kind of death (sprite) animation
func (p *Player) Draw(ctx rendering.Ctx) {
	p.jumpParticlesBroad.Draw(rendering.WithLayer(ctx, rendering.ScreenLayers.Playerspace))
	p.jumpParticlesTight.Draw(rendering.WithLayer(ctx, rendering.ScreenLayers.Playerspace))

	if p.State == Dying {
		p.deathAnim.Draw(ctx)
	} else {
		posX, posY := p.movebox.GetPos()
		jumpOffsetX, jumpOffsetY := p.calculateJumpOffset()
		ebitenrenderutil.DrawAtRotated(
			p.animator.GetSprite(),
			ctx.Dst,
			posX-ctx.CamX-jumpOffsetX,
			posY-ctx.CamY-jumpOffsetY,
			maths.ToRadians(p.direction),
			0.5,
			0.5,
		)
	}
	// vector.StrokeRect(
	// 	rendering.RenderLayers.Playerspace,
	// 	float32(p.hitbox.Left()),
	// 	float32(p.hitbox.Top()),
	// 	float32(p.hitbox.Width()),
	// 	float32(p.hitbox.Height()),
	// 	1,
	// 	color.RGBA{255, 0, 0, 255},
	// 	false,
	// )
}

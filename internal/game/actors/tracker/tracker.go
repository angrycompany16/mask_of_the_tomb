package tracker

import (
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"math"
)

type Tracker struct {
	*transform2D.Transform2D
	OnMoveFinished         *events.Event
	isMoving               bool
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
	moveSpeed              float64
}

func (t *Tracker) Update(cmd *engine.Commands) {
	t.Transform2D.Update(cmd)
	t.posX += t.moveSpeed * t.moveDirX
	t.posY += t.moveSpeed * t.moveDirY

	if t.moveDirX < 0 {
		t.posX = maths.Clamp(t.posX, t.targetPosX, t.posX)
	} else if t.moveDirX > 0 {
		t.posX = maths.Clamp(t.posX, t.posX, t.targetPosX)
	}
	if t.moveDirY < 0 {
		t.posY = maths.Clamp(t.posY, t.targetPosY, t.posY)
	} else if t.moveDirY > 0 {
		t.posY = maths.Clamp(t.posY, t.posY, t.targetPosY)
	}

	if t.posX == t.targetPosX && t.posY == t.targetPosY && t.isMoving {
		t.OnMoveFinished.Raise(events.EventInfo{
			Data: maths.DirFromVector(t.moveDirX, t.moveDirY),
		})
		t.isMoving = false
	}

	if t.posX == t.targetPosX {
		t.moveDirX = 0
	}
	if t.posY == t.targetPosY {
		t.moveDirY = 0
	}

}

func (t *Tracker) SetTarget(x, y float64) {
	t.targetPosX = x
	t.targetPosY = y
	t.moveDirX = math.Copysign(1, t.targetPosX-t.posX)
	t.moveDirY = math.Copysign(1, t.targetPosY-t.posY)
	t.isMoving = true
}

func (t *Tracker) GetTarget() (float64, float64) {
	return t.targetPosX, t.targetPosY
}

func (t *Tracker) SetPos(x, y float64) {
	t.posX, t.posY = x, y
	t.targetPosX, t.targetPosY = x, y
}

func (t *Tracker) GetPos() (float64, float64) {
	return t.posX, t.posY
}

func (t *Tracker) GetMovedir() (float64, float64) {
	return t.moveDirX, t.moveDirY
}

func NewTracker(transform2D *transform2D.Transform2D, moveSpeed, x, y float64) *Tracker {
	_movebox := &Tracker{
		Transform2D:    transform2D,
		isMoving:       false,
		OnMoveFinished: events.NewEvent(),
		moveSpeed:      moveSpeed,
		posX:           x,
		posY:           y,
	}
	return _movebox
}

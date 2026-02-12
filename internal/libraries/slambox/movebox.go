package slambox

import (
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"math"
)

// An object that tracks a target position by moving at a fixed speed.
type Tracker struct {
	OnMoveFinished         *events.Event
	isMoving               bool
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
	moveSpeed              float64
}

func (t *Tracker) Update() {
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

func NewTracker(moveSpeed, x, y float64) *Tracker {
	_movebox := &Tracker{
		isMoving:       false,
		OnMoveFinished: events.NewEvent(),
		moveSpeed:      moveSpeed,
		posX:           x,
		posY:           y,
	}
	return _movebox
}

package tracker

import (
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"math"
)

type EventData struct {
	dir maths.Direction
}

type Tracker struct {
	*graphic.Graphic
	OnMoveFinishEv         *eventsv2.Event
	isMoving               bool
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
	moveSpeed              float64
}

func (t *Tracker) Update(cmd *engine.Commands) {
	t.Graphic.Update(cmd)
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
		t.OnMoveFinishEv.Raise().WithData("dir", maths.DirFromVector(t.moveDirX, t.moveDirY))
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

func NewTracker(graphic *graphic.Graphic, moveSpeed, x, y float64) *Tracker {
	_movebox := &Tracker{
		Graphic:  graphic,
		isMoving: false,
		// I feel like this is kinda dumb
		// Event listeners should store events, not the other way around
		OnMoveFinishEv: eventsv2.NewEvent(),
		moveSpeed:      moveSpeed,
		posX:           x,
		posY:           y,
	}
	return _movebox
}

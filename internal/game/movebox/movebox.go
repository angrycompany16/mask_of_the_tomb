package movebox

import (
	"mask_of_the_tomb/internal/game/events"
	"mask_of_the_tomb/internal/maths"
	"math"
)

type Movebox struct {
	OnMoveFinished         *events.Event
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
	moveSpeed              float64
}

func (m *Movebox) Update() {
	m.posX += m.moveSpeed * m.moveDirX
	m.posY += m.moveSpeed * m.moveDirY

	if m.moveDirX < 0 {
		m.posX = maths.Clamp(m.posX, m.targetPosX, m.posX)
	} else if m.moveDirX > 0 {
		m.posX = maths.Clamp(m.posX, m.posX, m.targetPosX)
	}
	if m.moveDirY < 0 {
		m.posY = maths.Clamp(m.posY, m.targetPosY, m.posY)
	} else if m.moveDirY > 0 {
		m.posY = maths.Clamp(m.posY, m.posY, m.targetPosY)
	}

	if m.posX == m.targetPosX {
		m.moveDirX = 0
	}
	if m.posY == m.targetPosY {
		m.moveDirY = 0
	}

	if m.posX == m.targetPosX && m.posY == m.targetPosY {
		m.OnMoveFinished.Raise(events.EventInfo{})
	}
}

func (m *Movebox) SetTarget(x, y float64) {
	m.targetPosX = x
	m.targetPosY = y
	m.moveDirX = math.Copysign(1, m.targetPosX-m.posX)
	m.moveDirY = math.Copysign(1, m.targetPosY-m.posY)
}

func (m *Movebox) SetPos(x, y float64) {
	m.posX, m.posY = x, y
	m.targetPosX, m.targetPosY = x, y
}

func (m *Movebox) GetPos() (float64, float64) {
	return m.posX, m.posY
}

func (m *Movebox) GetMovedir() (float64, float64) {
	return m.moveDirX, m.moveDirY
}

func NewMovebox(moveSpeed float64) *Movebox {
	_movebox := &Movebox{
		OnMoveFinished: events.NewEvent(),
		moveSpeed:      moveSpeed,
	}
	return _movebox
}

package moveBox

import (
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/engine/events"
	"mask_of_the_tomb/internal/entities/player/moveBox/pubmovebox"
	"mask_of_the_tomb/internal/libraries/maths"
	"math"
)

type movebox struct {
	*entities.Entity
	positionAdvertiser     pubmovebox.PositionAdvertiser
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
	moveProgress           float64
	moveSpeed              float64
}

func NewMovebox(posX, posY, moveSpeed float64, name string) *movebox {
	_movebox := &movebox{
		posX:       posX,
		posY:       posY,
		targetPosX: posX,
		targetPosY: posY,
		moveSpeed:  moveSpeed,
	}
	_movebox.Entity = entities.RegisterEntity(_movebox, name)
	return _movebox
}

func (m *movebox) Update() {
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
		pubmovebox.FinishedMoveEvent.Raise(events.EventInfo{
			EntityID: m.GetID(),
		})
	}

	m.positionAdvertiser.PosX = m.posX
	m.positionAdvertiser.PosY = m.posY
}

func (m *movebox) SetTarget(x, y float64) {
	m.targetPosX = x
	m.targetPosY = y
	m.moveDirX = math.Copysign(1, m.targetPosX-m.posX)
	m.moveDirY = math.Copysign(1, m.targetPosY-m.posY)
}

package moveBox

import (
	"mask_of_the_tomb/internal/libraries/maths"
	"math"
)

type MoveBox struct {
	PosX, PosY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
	moveProgress           float64
	moveSpeed              float64
}

func (m *MoveBox) Move() {
	m.PosX += m.moveSpeed * m.moveDirX
	m.PosY += m.moveSpeed * m.moveDirY

	if m.moveDirX < 0 {
		m.PosX = maths.Clamp(m.PosX, m.targetPosX, m.PosX)
	} else if m.moveDirX > 0 {
		m.PosX = maths.Clamp(m.PosX, m.PosX, m.targetPosX)
	}
	if m.moveDirY < 0 {
		m.PosY = maths.Clamp(m.PosY, m.targetPosY, m.PosY)
	} else if m.moveDirY > 0 {
		m.PosY = maths.Clamp(m.PosY, m.PosY, m.targetPosY)
	}

	if m.PosX == m.targetPosX {
		m.moveDirX = 0
	}
	if m.PosY == m.targetPosY {
		m.moveDirY = 0
	}

	// if m.PosX == m.targetPosX && m.PosY == m.targetPosY {
	// 	m.angle = m.angle - math.Pi
	// 	m.State = Idle
	// }
}

func (m *MoveBox) SetTarget(x, y float64) {
	m.targetPosX = x
	m.targetPosY = y
	// m.prevPosX = m.PosX
	// m.prevPosY = m.PosY
	m.moveDirX = math.Copysign(1, m.targetPosX-m.PosX)
	m.moveDirY = math.Copysign(1, m.targetPosY-m.PosY)
}

package missile

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"
	"math"
)

// Object that tracks targets
type Missile struct {
	*transform2D.Transform2D
	Active           bool
	TargetX, TargetY float64
	TargetTransform  *transform2D.Transform2D
	speed            float64
	radialOffset     float64
	angularOffset    float64
	rotationSpeed    float64
}

func (m *Missile) Update(cmd *commands.Commands) {
	m.Transform2D.Update(cmd)

	if !m.Active {
		return
	}

	m.angularOffset += m.rotationSpeed * 0.01
	circularOffsetX := m.radialOffset * math.Cos(m.angularOffset)
	circularOffsetY := m.radialOffset * math.Sin(m.angularOffset)

	if m.TargetTransform != nil {
		m.TargetX, m.TargetY = m.TargetTransform.GetPos(false)
	}

	x, y := m.GetPos(false)
	x = maths.Lerp(x, m.TargetX+circularOffsetX, 0.01*m.speed)
	y = maths.Lerp(y, m.TargetY+circularOffsetY, 0.01*m.speed)
	m.SetPos(x, y)
}

func defaultMissile() *Missile {
	return &Missile{
		Transform2D:   transform2D.NewTransform2D(),
		Active:        true,
		TargetX:       0,
		TargetY:       0,
		speed:         1,
		radialOffset:  0,
		angularOffset: 0,
		rotationSpeed: 0,
	}
}

func NewMissile(options ...utils.Option[Missile]) *Missile {
	missile := defaultMissile()

	for _, op := range options {
		op(missile)
	}

	return missile
}

func WithTransform(transform *transform2D.Transform2D) utils.Option[Missile] {
	return func(m *Missile) {
		m.Transform2D = transform
	}
}

func WithActive(active bool) utils.Option[Missile] {
	return func(m *Missile) {
		m.Active = active
	}
}

func WithSpeed(speed float64) utils.Option[Missile] {
	return func(m *Missile) {
		m.speed = speed
	}
}

func WithCircularOffset(radialOffset float64, rotationSpeed float64, angularOffset float64) utils.Option[Missile] {
	return func(m *Missile) {
		m.radialOffset = radialOffset
		m.rotationSpeed = rotationSpeed
		m.angularOffset = angularOffset
	}
}

func WithTargetTransform(transform *transform2D.Transform2D) utils.Option[Missile] {
	return func(m *Missile) {
		m.TargetTransform = transform
	}
}

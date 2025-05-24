package hazard

import (
	"mask_of_the_tomb/internal/core/maths"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

const (
	hazardDamageFieldName = "Damage"
)

type Hazard struct {
	Hitbox maths.Rect
	Damage float64
}

func NewHazard(
	entity *ebitenLDTK.Entity,
) Hazard {
	newHazard := Hazard{}
	newHazard.Hitbox = *maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)

	for _, fieldInstance := range entity.Fields {
		if fieldInstance.Name == hazardDamageFieldName {
			newHazard.Damage = fieldInstance.Float
		}
	}
	return newHazard
}

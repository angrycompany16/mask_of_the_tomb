package world

import (
	"mask_of_the_tomb/internal/maths"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

type hazard struct {
	hitbox maths.Rect
	damage float64
}

func newHazard(
	entity *ebitenLDTK.Entity,
) hazard {
	newHazard := hazard{}
	newHazard.hitbox = *maths.RectFromEntity(entity)

	for _, fieldInstance := range entity.Fields {
		if fieldInstance.Name == hazardDamageFieldName {
			newHazard.damage = fieldInstance.Float
		}
	}
	return newHazard
}

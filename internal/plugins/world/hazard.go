package world

import (
	"mask_of_the_tomb/internal/core/maths"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

const (
	HazardEntityName      = "Hazard"
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
	newHazard.Hitbox = *RectFromEntity(entity)

	for _, fieldInstance := range entity.Fields {
		if fieldInstance.Name == hazardDamageFieldName {
			newHazard.Damage = fieldInstance.Float
		}
	}
	return newHazard
}

package world

import "mask_of_the_tomb/ebitenLDTK"

type Hazard struct {
	posX, posY    float64
	width, height float64
	damage        float64
}

func newHazard(
	entity *ebitenLDTK.Entity,
) Hazard {
	newHazard := Hazard{}
	newHazard.posX = entity.Px[0]
	newHazard.posY = entity.Px[1]
	newHazard.width = entity.Width
	newHazard.height = entity.Height

	for _, fieldInstance := range entity.Fields {
		if fieldInstance.Name == hazardDamageFieldName {
			newHazard.damage = fieldInstance.Float
		}
	}
	return newHazard
}

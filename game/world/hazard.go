package world

import "mask_of_the_tomb/ebitenLDTK"

type Hazard struct {
	posX, posY    float64
	width, height float64
	damage        float64
}

func newHazard(
	entityInstance *ebitenLDTK.EntityInstance,
) Hazard {
	newHazard := Hazard{}
	newHazard.posX = entityInstance.Px[0]
	newHazard.posY = entityInstance.Px[1]
	newHazard.width = entityInstance.Width
	newHazard.height = entityInstance.Height

	for _, fieldInstance := range entityInstance.FieldInstances {
		if fieldInstance.Name == hazardDamageFieldName {
			newHazard.damage = fieldInstance.Float
		}
	}
	return newHazard
}

package world

import (
	"mask_of_the_tomb/ebitenLDTK"
	"mask_of_the_tomb/rect"
	. "mask_of_the_tomb/utils"
)

type Door struct {
	levelIid   string
	hitbox     rect.Rect
	posX, posY float64
}

func newDoor(
	entityInstance *ebitenLDTK.EntityInstance,
) Door {
	newDoor := Door{}
	newDoor.posX = entityInstance.Px[0]
	newDoor.posY = entityInstance.Px[0]
	newDoor.hitbox = *rect.FromEntity(entityInstance)

	fieldInstance, err := entityInstance.GetFieldInstanceByName(doorOtherSideFieldName)
	HandleLazy(err)
	newDoor.levelIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

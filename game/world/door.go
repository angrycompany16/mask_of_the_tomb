package world

import (
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/utils"
	"mask_of_the_tomb/utils/rect"
)

type Door struct {
	levelIid  string
	entityIid string
	hitbox    rect.Rect
}

func newDoor(
	entityInstance *ebitenLDTK.EntityInstance,
	level *ebitenLDTK.Level,
) Door {
	newDoor := Door{}
	newDoor.hitbox = *rect.FromEntity(entityInstance)

	fieldInstance, err := entityInstance.GetFieldInstanceByName(doorOtherSideFieldName)
	HandleLazy(err)
	newDoor.levelIid = fieldInstance.EntityRef.LevelIid
	newDoor.entityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

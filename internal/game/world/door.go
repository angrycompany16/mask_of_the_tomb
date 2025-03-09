package world

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/maths"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

type door struct {
	levelIid  string
	entityIid string
	hitbox    maths.Rect
}

func newDoor(
	entityInstance *ebitenLDTK.Entity,
	// level *ebitenLDTK.Level,
) door {
	newDoor := door{}
	newDoor.hitbox = *maths.RectFromEntity(entityInstance)

	fieldInstance := errs.Must(entityInstance.GetFieldByName(doorOtherSideFieldName))
	newDoor.levelIid = fieldInstance.EntityRef.LevelIid
	newDoor.entityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

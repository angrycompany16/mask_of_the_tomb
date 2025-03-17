package world

import (
	"mask_of_the_tomb/internal/libraries/assets/ldtknames"
	"mask_of_the_tomb/internal/libraries/errs"
	"mask_of_the_tomb/internal/libraries/maths"

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

	fieldInstance := errs.Must(entityInstance.GetFieldByName(ldtknames.DoorOtherSideFieldName))
	newDoor.levelIid = fieldInstance.EntityRef.LevelIid
	newDoor.entityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

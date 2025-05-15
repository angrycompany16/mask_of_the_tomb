package world

import (
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

const (
	doorOtherSideFieldName = "OtherSide"
)

type Door struct {
	LevelIid  string
	EntityIid string
	Hitbox    maths.Rect
}

func NewDoor(
	entityInstance *ebitenLDTK.Entity,
) Door {
	newDoor := Door{}
	newDoor.Hitbox = *RectFromEntity(entityInstance)

	fieldInstance := errs.Must(entityInstance.GetFieldByName(doorOtherSideFieldName))
	newDoor.LevelIid = fieldInstance.EntityRef.LevelIid
	newDoor.EntityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

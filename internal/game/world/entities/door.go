package entities

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/maths"

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
	newDoor.Hitbox = *maths.RectFromEntity(entityInstance)

	fieldInstance := errs.Must(entityInstance.GetFieldByName(doorOtherSideFieldName))
	newDoor.LevelIid = fieldInstance.EntityRef.LevelIid
	newDoor.EntityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

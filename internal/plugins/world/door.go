package world

import (
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	doorOtherSideFieldName = "OtherSide"
)

type Door struct {
	LevelIid  string
	EntityIid string
	Hitbox    maths.Rect
	sprite    *ebiten.Image
}

func (d *Door) Draw() {
	x, y := d.Hitbox.TopLeft()
	camX, camY := camera.GetPos()
	ebitenrenderutil.DrawAt(d.sprite, rendering.RenderLayers.Playerspace, x-camX, y-camY)
}

func NewDoor(
	entityInstance *ebitenLDTK.Entity,
) Door {
	newDoor := Door{}
<<<<<<< HEAD:internal/plugins/world/door.go
	newDoor.Hitbox = *RectFromEntity(entityInstance)
=======
	newDoor.Hitbox = *maths.RectFromEntity(entityInstance)
	newDoor.sprite = errs.MustNewImageFromFile(doorSpritePath)
>>>>>>> 22a8537ec3bcc53973c6a0f42fc8f2788df75d55:internal/game/world/entities/door.go

	fieldInstance := errs.Must(entityInstance.GetFieldByName(doorOtherSideFieldName))
	newDoor.LevelIid = fieldInstance.EntityRef.LevelIid
	newDoor.EntityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

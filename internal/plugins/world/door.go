package world

import (
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/libraries/rendering"

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

func (d *Door) Draw(camX, camY float64) {
	x, y := d.Hitbox.TopLeft()
	ebitenrenderutil.DrawAt(d.sprite, rendering.RenderLayers.Playerspace, x-camX, y-camY)
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

package entities

import (
	ebitenrenderutil "mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	"mask_of_the_tomb/internal/maths"

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
	newDoor.Hitbox = *maths.RectFromEntity(entityInstance)
	newDoor.sprite = errs.MustNewImageFromFile(doorSpritePath)

	fieldInstance := errs.Must(entityInstance.GetFieldByName(doorOtherSideFieldName))
	newDoor.LevelIid = fieldInstance.EntityRef.LevelIid
	newDoor.EntityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

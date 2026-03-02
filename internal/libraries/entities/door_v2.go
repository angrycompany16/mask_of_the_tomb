package entities

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	// TODO: This needs to be fixed

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	doorV2OtherSideFieldName = "OtherSide"
	doorDirectionFieldName   = "Direction"
)

type DoorV2 struct {
	LevelIid  string
	EntityIid string
	Hitbox    *maths.Rect
	sprite    *ebiten.Image
	direction maths.Direction
}

func (d *DoorV2) Update() {

}

func (d *DoorV2) Draw(ctx rendering.Ctx) {
	x, y := d.Hitbox.TopLeft()
	ebitenrenderutil.DrawAt(d.sprite, ctx.Dst, x, y)
}

func NewDoorV2(
	entity *ebitenLDTK.Entity,
) DoorV2 {
	newDoor := DoorV2{}
	newDoor.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)

	newDoor.sprite = errs.Must(assettypes.GetImageAsset("doorV2"))

	// directionField := errs.Must(entity.GetFieldByName(doorDirectionFieldName))
	fieldInstance := errs.Must(entity.GetFieldByName(doorV2OtherSideFieldName))
	entityRef := ebitenLDTK.As[ebitenLDTK.EntityRef](fieldInstance)
	newDoor.LevelIid = entityRef.LevelIid
	newDoor.EntityIid = entityRef.EntityIid

	return newDoor
}

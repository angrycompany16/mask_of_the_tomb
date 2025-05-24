package door

import (
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

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

func (d *Door) Draw(ctx rendering.Ctx) {
	x, y := d.Hitbox.TopLeft()
	ebitenrenderutil.DrawAt(d.sprite, ctx.Dst, x-ctx.CamX, y-ctx.CamY)
}

func NewDoor(
	entity *ebitenLDTK.Entity,
) Door {
	newDoor := Door{}
	newDoor.Hitbox = *maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)

	fieldInstance := errs.Must(entity.GetFieldByName(doorOtherSideFieldName))
	newDoor.LevelIid = fieldInstance.EntityRef.LevelIid
	newDoor.EntityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

package entities

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	doorV2OtherSideFieldName = "OtherSide"
	doorDirectionFieldName   = "Direction"
)

type DoorV2 struct {
	EntityIid          string
	OtherSideLevelIid  string
	OtherSideEntityIid string
	Hitbox             *maths.Rect
	InteractRegion     *maths.Rect
	sprite             *ebiten.Image
	direction          maths.Direction
	isReady            bool
}

func (d *DoorV2) Update(playerX, playerY float64) {
	d.isReady = d.InteractRegion.Contains(playerX, playerY)
}

func (d *DoorV2) Draw(ctx rendering.Ctx) {
	if d.isReady {
		cx, cy := d.Hitbox.Center()
		// x, y := d.Hitbox.TopLeft()
		rot := maths.DirToRadians(d.direction)
		ebitenrenderutil.DrawAtRotated(d.sprite, ctx.Dst, cx, cy, rot, 0.5, 0.5)
	} else {
		// x, y := d.Hitbox.TopLeft()
		cx, cy := d.Hitbox.Center()
		fmt.Println(cx, cy)

		rot := maths.DirToRadians(maths.Opposite(d.direction))
		ebitenrenderutil.DrawAtRotated(d.sprite, ctx.Dst, cx, cy, rot, 0.5, 0.5)
		// ebitenrenderutil.DrawAt(d.sprite, ctx.Dst, cx, cy, rot, 0.5, 0.5)
	}
}

func (d *DoorV2) IsReady() bool {
	return d.isReady
}

// Hard-coded for now. Not great but might have to do
func (d *DoorV2) GetSpawnPos() (float64, float64) {
	cx, cy := d.Hitbox.Center()
	switch d.direction {
	case maths.DirUp:
		return cx - 8, d.Hitbox.Top() - 16
	case maths.DirDown:
		return cx - 8, d.Hitbox.Bottom()
	case maths.DirLeft:
		return d.Hitbox.Left() - 16, cy - 8
	case maths.DirRight:
		return d.Hitbox.Right(), cy - 8
	}
	return 0, 0
}

func (d *DoorV2) GetDir() maths.Direction {
	return d.direction
}

func NewDoorV2(
	entity *ebitenLDTK.Entity,
	levelLDTK *ebitenLDTK.Level,
) *DoorV2 {
	newDoor := DoorV2{}
	newDoor.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)

	newDoor.EntityIid = entity.Iid
	newDoor.sprite = errs.Must(assettypes.GetImageAsset("doorV2"))

	directionField := errs.Must(entity.GetFieldByName(doorDirectionFieldName))
	newDoor.direction = maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

	doorOtherSideField := errs.Must(entity.GetFieldByName(doorV2OtherSideFieldName))
	doorOtherSide := ebitenLDTK.As[ebitenLDTK.EntityRef](doorOtherSideField)

	newDoor.OtherSideLevelIid = doorOtherSide.LevelIid
	newDoor.OtherSideEntityIid = doorOtherSide.EntityIid

	interactRegionField := errs.Must(entity.GetFieldByName(resources.LDTKNames.DoorInteractRegionField))
	interactRegion := errs.Must(levelLDTK.GetEntityByIid(ebitenLDTK.As[ebitenLDTK.EntityRef](interactRegionField).EntityIid))

	newDoor.InteractRegion = maths.NewRect(
		interactRegion.Px[0],
		interactRegion.Px[1],
		interactRegion.Width,
		interactRegion.Height,
	)

	return &newDoor
}

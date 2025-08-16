package entities

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	// TODO: This needs to be fixed
	"mask_of_the_tomb/internal/libraries/animation"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

var ()

const (
	doorOtherSideFieldName = "OtherSide"
)

type Door struct {
	LevelIid  string
	EntityIid string
	animator  *animation.Animator
	Hitbox    maths.Rect
	sprite    *ebiten.Image
}

func (d *Door) Update() {
	d.animator.Update()
}

func (d *Door) Draw(ctx rendering.Ctx) {
	x, y := d.Hitbox.TopLeft()
	ebitenrenderutil.DrawAt(d.animator.GetSprite(), ctx.Dst, x, y)
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
	animatorMap := make(map[int]*animation.Animation)
	animationInfo := errs.Must(assettypes.GetYamlAsset("teleporterAnimation")).(*animation.AnimationInfo)
	animatorMap[1] = animation.NewAnimation(*animationInfo)
	newDoor.animator = animation.MakeAnimator(animatorMap)
	newDoor.animator.SwitchClip(1)

	fieldInstance := errs.Must(entity.GetFieldByName(doorOtherSideFieldName))
	newDoor.LevelIid = fieldInstance.EntityRef.LevelIid
	newDoor.EntityIid = fieldInstance.EntityRef.EntityIid

	return newDoor
}

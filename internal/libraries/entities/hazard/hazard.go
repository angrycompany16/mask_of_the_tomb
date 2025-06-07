package hazard

import (
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	hazardDamageFieldName = "Damage"
)

type Hazard struct {
	Sprite     *ebiten.Image // NOTE: not always used, only for slamboxes
	Hitbox     *maths.Rect
	PosOffsetX float64
	PosOffsetY float64
	LinkID     string // For linking with slamboxes
}

func (h *Hazard) Draw(ctx rendering.Ctx) {
	ebitenrenderutil.DrawAt(h.Sprite, ctx.Dst, h.Hitbox.Left(), h.Hitbox.Top())
}

func NewHazard(
	entity *ebitenLDTK.Entity,
) *Hazard {
	newHazard := &Hazard{}
	newHazard.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)
	newHazard.LinkID = entity.Iid
	newHazard.Sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))

	return newHazard
}

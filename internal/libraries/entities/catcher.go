package entities

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

type Catcher struct {
	sprite *ebiten.Image
	Hitbox *maths.Rect
}

func (c *Catcher) Draw(ctx rendering.Ctx) {
	ebitenrenderutil.DrawAt(c.sprite, ctx.Dst, c.Hitbox.Left(), c.Hitbox.Top())
}

func NewCatcher(
	entity *ebitenLDTK.Entity,
) *Catcher {
	newCatcher := &Catcher{}
	newCatcher.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)
	newCatcher.sprite = errs.Must(assettypes.GetImageAsset("catcherSprite"))

	return newCatcher
}
